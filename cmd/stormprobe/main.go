package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/umutozen/stormprobe/internal/alert"
	"github.com/umutozen/stormprobe/internal/compare"
	"github.com/umutozen/stormprobe/internal/config"
	"github.com/umutozen/stormprobe/internal/discovery"
	"github.com/umutozen/stormprobe/internal/report"
	"github.com/umutozen/stormprobe/internal/runner"
	"github.com/umutozen/stormprobe/internal/verdict"
)

var validFormats = map[string]bool{
	"json": true,
	"html": true,
	"both": true,
}

type headerFlag []string

func (h *headerFlag) String() string { return strings.Join(*h, ", ") }
func (h *headerFlag) Set(v string) error {
	if !strings.Contains(v, ":") {
		return fmt.Errorf("header must be in 'Key: Value' format, got: %q", v)
	}
	*h = append(*h, v)
	return nil
}

func parseHeaders(raw headerFlag) map[string]string {
	out := make(map[string]string, len(raw))
	for _, h := range raw {
		parts := strings.SplitN(h, ":", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key != "" {
			out[key] = val
		}
	}
	return out
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "compare" {
		runCompare(os.Args[2:])
		return
	}

	configFile := flag.String("config", "", "YAML config file path")
	concRamp := flag.Int("concurrency-ramp", 50, "Peak concurrency for ramp-up phase")
	concSustained := flag.Int("concurrency-sustained", 75, "Concurrency for sustained phase")
	concSpike := flag.Int("concurrency-spike", 250, "Peak concurrency for spike phase")
	reqPerWorker := flag.Int("req-per-worker", 15, "Requests per worker per phase step")
	timeout := flag.Duration("timeout", 10*time.Second, "Per-request timeout")
	phaseDuration := flag.Duration("duration", 0, "Per-phase duration (e.g. 30s, 1m). Overrides req-per-worker when set")
	endpointsFile := flag.String("endpoints", "", "Endpoints file (one path per line, skips discovery)")
	noDiscovery := flag.Bool("no-discovery", false, "Skip katana+httpx, use / only")
	outputDir := flag.String("output", "./outputs", "Output directory for reports")
	format := flag.String("format", "both", "Report format: json, html, both")
	katanaPath := flag.String("katana-path", "", "Custom katana binary path")
	httpxPath := flag.String("httpx-path", "", "Custom httpx binary path")
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification")
	alertP99 := flag.Float64("alert-p99", 0, "Fail (exit 1) if P99 latency exceeds Xms in any phase")
	alertErrorRate := flag.Float64("alert-error-rate", 0, "Fail (exit 1) if error rate exceeds X%% in any phase")
	alertMinRPS := flag.Float64("alert-rps", 0, "Fail (exit 1) if req/s falls below X in any phase")
	successMatch := flag.String("success-match", "", "Regex: body must match for success (e.g. 'borç|ödeme')")
	failMatch := flag.String("fail-match", "", "Regex: body match means failure (e.g. 'error|fault|bakım')")
	autoProfile := flag.Bool("auto-profile", true, "Auto-detect server stack and adjust diagnostic profile")
	profile := flag.String("profile", "", "Manual diagnostic profile override: iis, tomcat, php")
	noFingerprint := flag.Bool("no-fingerprint", false, "Skip stack fingerprinting")
	showProfile := flag.Bool("show-profile", false, "Show detected stack and diagnostic profile, then exit")
	var headers headerFlag
	flag.Var(&headers, "header", "Custom HTTP header (repeatable): -header 'Authorization: Bearer TOKEN'")
	flag.Var(&headers, "H", "Alias for --header")
	flag.Parse()

	flagsSet := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagsSet[f.Name] = true })

	if !validFormats[*format] {
		fmt.Fprintf(os.Stderr, "Invalid format: %q (valid: json, html, both)\n", *format)
		os.Exit(1)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  StormProbe v%s  |  Crawl. Discover. Break.\n", config.Version)
	fmt.Println("  github.com/umutozen/stormprobe")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	targetURL := ""
	if flag.NArg() > 0 {
		targetURL = flag.Arg(0)
	}
	if targetURL == "" && *configFile == "" {
		fmt.Print("\n  Enter target URL: ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		targetURL = strings.TrimSpace(line)
	}

	cfg := config.Config{
		Target:         targetURL,
		RequestTimeout: *timeout,
		EndpointsFile:  *endpointsFile,
		NoDiscovery:    *noDiscovery,
		Insecure:       *insecure,
		KatanaPath:     *katanaPath,
		HttpxPath:      *httpxPath,
		OutputDir:      *outputDir,
		Format:         *format,
		RampPeak:       *concRamp,
		SustainedConc:  *concSustained,
		SpikePeak:      *concSpike,
		ReqPerWorker:   *reqPerWorker,
		Headers:        parseHeaders(headers),
		PhaseDuration:  *phaseDuration,
		ConfigFile:     *configFile,
		SuccessMatch:   *successMatch,
		FailMatch:      *failMatch,
		AutoProfile:    *autoProfile,
		Profile:        *profile,
		NoFingerprint:  *noFingerprint,
		Alert: config.AlertConfig{
			MaxP99Ms:     *alertP99,
			MaxErrorRate: *alertErrorRate,
			MinReqPerSec: *alertMinRPS,
		},
	}

	if cfg.ConfigFile != "" {
		fc, loadErr := config.LoadFile(cfg.ConfigFile)
		if loadErr != nil {
			fmt.Fprintf(os.Stderr, "[!] %v\n", loadErr)
			os.Exit(1)
		}
		config.MergeFileIntoConfig(fc, &cfg, flagsSet)
		fmt.Printf("  Config  : loaded from %s\n", cfg.ConfigFile)

		if len(fc.Targets) >= 2 {
			runMultiTarget(fc.Targets, cfg)
			return
		}
	}

	if cfg.Target == "" {
		fmt.Fprintln(os.Stderr, "Usage: stormprobe [flags] <target-url>")
		os.Exit(1)
	}
	if !strings.HasPrefix(cfg.Target, "http") {
		cfg.Target = "https://" + cfg.Target
	}
	parsed, err := url.Parse(cfg.Target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid URL: %s\n", err)
		os.Exit(1)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		fmt.Fprintf(os.Stderr, "Unsupported scheme: %s (only http and https)\n", parsed.Scheme)
		os.Exit(1)
	}
	targetURL = cfg.Target
	fmt.Printf("\n  Target  : %s\n", targetURL)

	if *showProfile {
		fp := discovery.Fingerprint(targetURL, cfg.Insecure)
		discovery.PrintFingerprint(fp)
		return
	}

	jsonPath := runSingleTarget(cfg, "")

	if jsonPath != "" {
		_ = jsonPath
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  Done. Reports saved to", cfg.OutputDir)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func runMultiTarget(targets []config.TargetEntry, baseCfg config.Config) {
	sep := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

	if len(targets) < 2 {
		fmt.Fprintln(os.Stderr, "[!] Multi-target mode requires at least 2 targets")
		os.Exit(1)
	}
	for _, t := range targets {
		if t.Name == "" {
			fmt.Fprintln(os.Stderr, "[!] Each target must have a \"name\" field in multi-target config")
			os.Exit(1)
		}
		if t.URL == "" {
			fmt.Fprintf(os.Stderr, "[!] Target %q is missing \"url\" field\n", t.Name)
			os.Exit(1)
		}
	}
	if len(targets) > 2 {
		fmt.Println("  [!] Multi-target mode currently supports 2 targets. Using first two.")
		targets = targets[:2]
	}

	fmt.Printf("\n  Multi-Target Mode: %d targets\n", len(targets))
	for i, t := range targets {
		fmt.Printf("    [%d] %s → %s\n", i+1, t.Name, t.URL)
	}

	type targetResult struct {
		name     string
		jsonPath string
		rating   string
		safeConc int
	}

	var targetResults []targetResult
	for i, t := range targets {
		fmt.Println("\n" + sep)
		fmt.Printf("  TARGET %d/%d: %s\n", i+1, len(targets), t.Name)
		fmt.Println(sep)

		cfg := baseCfg
		targetURL := t.URL
		if !strings.HasPrefix(targetURL, "http") {
			targetURL = "https://" + targetURL
		}
		cfg.Target = targetURL

		jsonPath := runSingleTarget(cfg, t.Name)

		tr := targetResult{name: t.Name, jsonPath: jsonPath}
		if jsonPath != "" {
			if rpt, err := compare.LoadReport(jsonPath); err == nil {
				tr.rating = rpt.Verdict.Rating
				tr.safeConc = rpt.Verdict.SafeConcurrency
			}
		}
		targetResults = append(targetResults, tr)

		if i < len(targets)-1 {
			fmt.Println("\n  Cooldown between targets: 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Println("\n" + sep)
	fmt.Println("  TARGET SUMMARY")
	fmt.Println(sep)
	fmt.Println("  ┌──────────────────────┬──────────────────────────────┬──────────┐")
	fmt.Println("  │ Target               │ Rating                       │ Safe VU  │")
	fmt.Println("  ├──────────────────────┼──────────────────────────────┼──────────┤")
	for _, tr := range targetResults {
		rating := tr.rating
		if rating == "" {
			rating = "N/A"
		}
		fmt.Printf("  │ %-20s │ %-28s │ %6d   │\n", tr.name, rating, tr.safeConc)
	}
	fmt.Println("  └──────────────────────┴──────────────────────────────┴──────────┘")

	var jsonPaths []string
	for _, tr := range targetResults {
		if tr.jsonPath != "" {
			jsonPaths = append(jsonPaths, tr.jsonPath)
		}
	}

	if len(jsonPaths) == 2 {
		fmt.Println("\n" + sep)
		fmt.Printf("  AUTO-COMPARE: %s vs %s\n", targets[0].Name, targets[1].Name)
		fmt.Println(sep)

		before, err1 := compare.LoadReport(jsonPaths[0])
		after, err2 := compare.LoadReport(jsonPaths[1])

		if err1 != nil || err2 != nil {
			if err1 != nil {
				fmt.Fprintf(os.Stderr, "  [!] Cannot load report 1: %v\n", err1)
			}
			if err2 != nil {
				fmt.Fprintf(os.Stderr, "  [!] Cannot load report 2: %v\n", err2)
			}
			return
		}

		result := compare.Analyze(before, after)
		printCompareResult(result, targets[0].Name, targets[1].Name)
	}

	fmt.Println("\n" + sep)
	fmt.Println("  Done. Reports saved to", baseCfg.OutputDir)
	fmt.Println(sep)
}

func runSingleTarget(cfg config.Config, label string) string {
	targetURL := cfg.Target

	tlsConfig := &tls.Config{}
	if cfg.Insecure {
		tlsConfig.InsecureSkipVerify = true
		fmt.Println("  Warning : TLS verification disabled (--insecure)")
	}

	client := &http.Client{
		Timeout: cfg.RequestTimeout,
		Transport: &http.Transport{
			TLSClientConfig:     tlsConfig,
			MaxIdleConns:        600,
			MaxIdleConnsPerHost: 600,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	var endpoints []string
	if cfg.EndpointsFile != "" {
		endpoints = loadEndpointsFile(cfg.EndpointsFile)
	}
	if len(endpoints) == 0 {
		endpoints = discovery.Discover(targetURL, cfg.KatanaPath, cfg.HttpxPath, cfg.NoDiscovery, cfg.Headers)
	}
	if len(endpoints) == 0 {
		fmt.Fprintln(os.Stderr, "No endpoints available. Use --endpoints to provide a list.")
		return ""
	}

	fmt.Println("\n  Warmup  : sending 3 requests to prime connection pool...")
	for i := 0; i < 3; i++ {
		req, reqErr := http.NewRequest("GET", targetURL, nil)
		if reqErr != nil {
			continue
		}
		for k, v := range cfg.Headers {
			req.Header.Set(k, v)
		}
		if resp, doErr := client.Do(req); doErr == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}

	var results []config.PhaseResult

	var successRe, failRe *regexp.Regexp
	if cfg.SuccessMatch != "" {
		re, reErr := regexp.Compile("(?i)" + cfg.SuccessMatch)
		if reErr != nil {
			fmt.Fprintf(os.Stderr, "[!] Invalid success-match regex: %v\n", reErr)
			return ""
		}
		successRe = re
		fmt.Printf("  Body check : success-match = %s\n", cfg.SuccessMatch)
	}
	if cfg.FailMatch != "" {
		re, reErr := regexp.Compile("(?i)" + cfg.FailMatch)
		if reErr != nil {
			fmt.Fprintf(os.Stderr, "[!] Invalid fail-match regex: %v\n", reErr)
			return ""
		}
		failRe = re
		fmt.Printf("  Body check : fail-match = %s\n", cfg.FailMatch)
	}

	phaseHeader := func(name string) {
		fmt.Printf("\n┌─────────────────────────────────────────────────────────┐\n")
		fmt.Printf("│  %-55s│\n", name)
		fmt.Printf("└─────────────────────────────────────────────────────────┘\n")
	}
	tableHeader := func() {
		fmt.Println("  ┌────────────────────────┬──────┬─────────┬────────┬────────┬────────┬────────┬─────────┐")
		fmt.Println("  │ Scenario               │ Conc │ Success │ Avg ms │ P50 ms │ P95 ms │ P99 ms │  req/s  │")
		fmt.Println("  ├────────────────────────┼──────┼─────────┼────────┼────────┼────────┼────────┼─────────┤")
	}
	tableRow := func(r config.PhaseResult) {
		fmt.Printf("  │ %-22s │ %4d │ %5.1f%%  │ %6.0f │ %6.0f │ %6.0f │ %6.0f │ %7.1f │\n",
			r.PhaseName, r.Concurrency, r.SuccessRate, r.AvgMs, r.P50Ms, r.P95Ms, r.P99Ms, r.ReqPerSec)
	}
	tableFooter := func() {
		fmt.Println("  └────────────────────────┴──────┴─────────┴────────┴────────┴────────┴────────┴─────────┘")
	}

	phaseHeader("PHASE 1: RAMP-UP")
	tableHeader()
	for i, step := range cfg.RampSteps() {
		if i > 0 {
			time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
		}
		r := runner.RunPhase(step, targetURL, client, endpoints, cfg.Headers, successRe, failRe)
		results = append(results, r)
		tableRow(r)
		time.Sleep(2 * time.Second)
	}
	tableFooter()

	phaseHeader("PHASE 2: SUSTAINED")
	tableHeader()
	for _, step := range cfg.SustainedSteps() {
		r := runner.RunPhase(step, targetURL, client, endpoints, cfg.Headers, successRe, failRe)
		results = append(results, r)
		tableRow(r)
		time.Sleep(2 * time.Second)
	}
	tableFooter()

	phaseHeader("PHASE 3: SPIKE")
	tableHeader()
	for _, step := range cfg.SpikeSteps() {
		r := runner.RunPhase(step, targetURL, client, endpoints, cfg.Headers, successRe, failRe)
		results = append(results, r)
		tableRow(r)
		time.Sleep(3 * time.Second)
	}
	tableFooter()

	phaseHeader("PHASE 4: RECOVERY")
	fmt.Println("  Cooldown : waiting 10 seconds post-spike...")
	time.Sleep(10 * time.Second)
	tableHeader()
	r := runner.RunPhase(cfg.RecoveryStep(), targetURL, client, endpoints, cfg.Headers, successRe, failRe)
	results = append(results, r)
	tableRow(r)
	tableFooter()

	printErrorSummary(results)

	var stackSummary, overrideAudit string
	if !cfg.NoFingerprint && cfg.AutoProfile {
		fp := discovery.Fingerprint(targetURL, cfg.Insecure)
		discovery.PrintFingerprint(fp)
		stackSummary = fp.Summary()
		if cfg.Profile != "" && cfg.Profile != fp.DiagnosticProfile {
			overrideAudit = fmt.Sprintf("Manual profile override used (%s). Auto-detected fingerprint: %s.", cfg.Profile, fp.Summary())
			fmt.Printf("  ⚠  %s\n", overrideAudit)
		}
	}

	v := verdict.Analyze(results, stackSummary, overrideAudit)
	printVerdict(v)

	var jsonPath string
	if cfg.Format == "json" || cfg.Format == "both" {
		path, writeErr := report.WriteJSONWithLabel(targetURL, cfg.OutputDir, label, results, endpoints, v)
		if writeErr != nil {
			fmt.Fprintf(os.Stderr, "[!] %v\n", writeErr)
		} else {
			jsonPath = path
		}
	}
	if cfg.Format == "html" || cfg.Format == "both" {
		if writeErr := report.WriteHTML(targetURL, cfg.OutputDir, results, endpoints, v); writeErr != nil {
			fmt.Fprintf(os.Stderr, "[!] %v\n", writeErr)
		}
	}

	ihlaller := alert.Check(results, cfg.Alert)
	if len(ihlaller) > 0 {
		fmt.Println("\n" + strings.Repeat("=", 57))
		fmt.Println("  ALERT  : threshold breached -- CI/CD gate failed")
		fmt.Println(strings.Repeat("=", 57))
		for _, ih := range ihlaller {
			fmt.Printf("  [FAIL] %s\n", ih.Message)
		}
		fmt.Println(strings.Repeat("=", 57))
	}

	return jsonPath
}

func loadEndpointsFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Cannot read endpoints file: %v\n", err)
		return nil
	}
	defer f.Close()
	var eps []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			eps = append(eps, line)
		}
	}
	fmt.Printf("  [Endpoints] %d paths loaded from %s\n", len(eps), path)
	return eps
}

func printErrorSummary(results []config.PhaseResult) {
	var tTo, tRs, tRf, t5, t4, tBm, tOt, tFail, tAll int
	for _, r := range results {
		tTo += r.Errors.Timeout
		tRs += r.Errors.ConnectionReset
		tRf += r.Errors.ConnectionRefused
		t5 += r.Errors.HTTP5xx
		t4 += r.Errors.HTTP4xx
		tBm += r.Errors.BodyMismatch
		tOt += r.Errors.Other
		tFail += r.Failed
		tAll += r.TotalRequests
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  ERROR ANALYSIS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if tFail > 0 && tAll > 0 {
		fmt.Printf("  Total     : %d requests | %d failed (%.1f%%)\n", tAll, tFail, float64(tFail)/float64(tAll)*100)
		fmt.Println()
		pf := func(label string, count int) {
			if count == 0 {
				return
			}
			fmt.Printf("  %-22s %d (%.1f%% of failures)\n",
				label+":", count, float64(count)/float64(tFail)*100)
		}
		pf("Timeout", tTo)
		pf("Connection Reset", tRs)
		pf("Connection Refused", tRf)
		pf("HTTP 5xx", t5)
		pf("HTTP 4xx", t4)
		pf("Body Mismatch", tBm)
		pf("Other", tOt)
	} else {
		fmt.Printf("  Total     : %d requests | No errors recorded.\n", tAll)
	}

	dist := make(map[int]int)
	for _, r := range results {
		for code, count := range r.StatusDist {
			dist[code] += count
		}
	}
	if len(dist) > 0 {
		fmt.Println()
		fmt.Println("  HTTP Status Distribution:")
		for code, count := range dist {
			fmt.Printf("    %d  %d requests\n", code, count)
		}
	}
}

func printVerdict(v verdict.Verdict) {
	sep := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	fmt.Println("\n" + sep)
	fmt.Println("  FINAL VERDICT")
	fmt.Println(sep)

	fmt.Printf("  Rating                          : %s (%s)\n", v.Rating, v.Confidence)

	if v.SafeConcurrency > 0 {
		fmt.Printf("  Safe concurrency                : up to %d virtual users\n", v.SafeConcurrency)
	} else {
		fmt.Println("  Safe concurrency                : could not determine")
	}

	if v.DegradationAt != nil {
		fmt.Printf("  Perf. degradation starts at     : %d vu\n", *v.DegradationAt)
	}
	if v.FailurePoint != nil {
		fmt.Printf("  Severe degradation starts at    : %d vu\n", *v.FailurePoint)
	}
	if v.ErrorZoneAt > 0 {
		fmt.Printf("  Failure zone begins at          : %d vu\n", v.ErrorZoneAt)
	}
	if v.ThroughputPlateau > 0 {
		fmt.Printf("  Throughput plateau              : ~%.0f req/s\n", v.ThroughputPlateau)
	}

	recoveryStr := "OK"
	if !v.RecoveryOK {
		recoveryStr = "FAILED — server did not recover after spike"
	}
	fmt.Printf("  Recovery                        : %s\n", recoveryStr)

	if v.BottleneckCause != "None" {
		fmt.Println()
		if v.StackSummary != "" {
			fmt.Printf("  Detected Stack : %s\n", v.StackSummary)
		}
		if v.PrimaryCause != "" {
			fmt.Printf("  Primary Cause  : %s\n", v.PrimaryCause)
		} else {
			fmt.Printf("  Bottleneck     : %s\n", v.BottleneckCause)
		}
		if len(v.SecondaryContributors) > 0 {
			fmt.Println("  Secondary Contributors:")
			for _, s := range v.SecondaryContributors {
				fmt.Printf("    - %s\n", s)
			}
		}
		if len(v.WhyDiagnosis) > 0 {
			fmt.Println("  Why this diagnosis?")
			for _, w := range v.WhyDiagnosis {
				fmt.Printf("    - %s\n", w)
			}
		}
		if v.ProfileOverrideAudit != "" {
			fmt.Println()
			fmt.Printf("  ⚠  %s\n", v.ProfileOverrideAudit)
		}
	}

	fmt.Println()
	fmt.Printf("  Recommendation : %s\n", v.Recommendation)

	if len(v.PriorityChecks) > 0 {
		fmt.Println()
		fmt.Println("  Priority checks:")
		for _, check := range v.PriorityChecks {
			fmt.Printf("    - %s\n", check)
		}
	}
	fmt.Println()
	fmt.Printf("  Production Recommendation : %s\n", v.ProductionRecommendation)
	fmt.Println(sep)
}

func runCompare(args []string) {
	sep := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	fmt.Println("\n" + sep)
	fmt.Printf("  StormProbe v%s  |  Compare Mode\n", config.Version)
	fmt.Println(sep)

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "  Usage: stormprobe compare <before.json> <after.json>")
		os.Exit(1)
	}

	before, err := compare.LoadReport(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [!] %v\n", err)
		os.Exit(1)
	}
	after, err := compare.LoadReport(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [!] %v\n", err)
		os.Exit(1)
	}

	result := compare.Analyze(before, after)
	printCompareResult(result, args[0], args[1])

	outputDir := "./outputs"
	if _, writeErr := report.WriteCompareJSON(result, outputDir); writeErr != nil {
		fmt.Fprintf(os.Stderr, "  [!] Compare JSON report error: %v\n", writeErr)
	}
	if _, writeErr := report.WriteCompareHTML(result, outputDir); writeErr != nil {
		fmt.Fprintf(os.Stderr, "  [!] Compare HTML report error: %v\n", writeErr)
	}
}

func printCompareResult(result compare.CompareResult, beforeFile, afterFile string) {
	sep := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	b := result.Before.Verdict
	a := result.After.Verdict

	if len(result.Warnings) > 0 {
		for _, w := range result.Warnings {
			fmt.Printf("  [!] Warning: %s\n", w)
		}
		fmt.Println()
	}

	befLabel := beforeFile
	if result.Before.ReportID != "" {
		befLabel = result.Before.ReportID
	}
	aftLabel := afterFile
	if result.After.ReportID != "" {
		aftLabel = result.After.ReportID
	}

	fmt.Printf("  Before : %s\n", befLabel)
	if result.Before.Target != "" {
		fmt.Printf("           %s | %s\n", result.Before.Timestamp, result.Before.Target)
	}
	fmt.Printf("  After  : %s\n", aftLabel)
	if result.After.Target != "" {
		fmt.Printf("           %s | %s\n", result.After.Timestamp, result.After.Target)
	}
	fmt.Println()

	fmt.Printf("  Rating           : %s → %s\n", b.Rating, a.Rating)
	if b.Confidence != "" || a.Confidence != "" {
		bc, ac := b.Confidence, a.Confidence
		if bc == "" {
			bc = "N/A"
		}
		if ac == "" {
			ac = "N/A"
		}
		fmt.Printf("  Confidence       : %s → %s\n", bc, ac)
	}
	fmt.Printf("  Safe concurrency : %d VU → %d VU\n", b.SafeConcurrency, a.SafeConcurrency)
	recovBefore, recovAfter := "OK", "OK"
	if !b.RecoveryOK {
		recovBefore = "FAILED"
	}
	if !a.RecoveryOK {
		recovAfter = "FAILED"
	}
	fmt.Printf("  Recovery         : %s → %s\n", recovBefore, recovAfter)

	fmt.Println()
	fmt.Println("  Bottleneck:")
	fmt.Printf("    Before: %s\n", b.BottleneckCause)
	if b.BottleneckCause == a.BottleneckCause {
		fmt.Printf("    After:  %s (unchanged)\n", a.BottleneckCause)
	} else {
		fmt.Printf("    After:  %s\n", a.BottleneckCause)
	}

	fmt.Println()
	verdictLine := result.Outcome
	if result.Outcome == "PARTIAL IMPROVEMENT" {
		verdictLine += " \u2014 Still " + a.Rating
	}
	content := "  Verdict: " + verdictLine
	boxWidth := len(content) + 4
	if boxWidth < 57 {
		boxWidth = 57
	}
	border := strings.Repeat("\u2500", boxWidth)
	fmt.Printf("  \u250c%s\u2510\n", border)
	fmt.Printf("  \u2502  Verdict: %-*s\u2502\n", boxWidth-12, verdictLine)
	fmt.Printf("  \u2514%s\u2518\n", border)

	if len(result.Why) > 0 {
		fmt.Println()
		fmt.Println("  Why:")
		for _, w := range result.Why {
			fmt.Printf("    %s\n", w)
		}
	}
	fmt.Println(sep)
}
