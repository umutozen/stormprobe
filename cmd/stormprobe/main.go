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
	"strings"
	"time"

	"github.com/umutozen/stormprobe/internal/alert"
	"github.com/umutozen/stormprobe/internal/config"
	"github.com/umutozen/stormprobe/internal/discovery"
	"github.com/umutozen/stormprobe/internal/report"
	"github.com/umutozen/stormprobe/internal/runner"
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
	var headers headerFlag
	flag.Var(&headers, "header", "Custom HTTP header (repeatable): -header 'Authorization: Bearer TOKEN'")
	flag.Var(&headers, "H", "Alias for --header")
	flag.Parse()

	if !validFormats[*format] {
		fmt.Fprintf(os.Stderr, "Invalid format: %q (valid: json, html, both)\n", *format)
		os.Exit(1)
	}

	fmt.Println("===========================================================")
	fmt.Printf("  StormProbe v%s -- HTTP Load Tester\n", config.Version)
	fmt.Println("  github.com/umutozen/stormprobe")
	fmt.Println("===========================================================")

	targetURL := ""
	if flag.NArg() > 0 {
		targetURL = flag.Arg(0)
	}
	if targetURL == "" {
		fmt.Print("\n  Enter target URL: ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		targetURL = strings.TrimSpace(line)
	}
	if targetURL == "" {
		fmt.Fprintln(os.Stderr, "Usage: stormprobe [flags] <target-url>")
		os.Exit(1)
	}
	if !strings.HasPrefix(targetURL, "http") {
		targetURL = "https://" + targetURL
	}
	parsed, err := url.Parse(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid URL: %s\n", err)
		os.Exit(1)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		fmt.Fprintf(os.Stderr, "Unsupported scheme: %s (only http and https)\n", parsed.Scheme)
		os.Exit(1)
	}
	fmt.Printf("\n  Target: %s\n", targetURL)

	cfg := config.Config{
		Target:         targetURL,
		RequestTimeout: *timeout,
		EndpointsFile:  *endpointsFile,
		NoDiscovery:    *noDiscovery,
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
		Alert: config.AlertConfig{
			MaxP99Ms:     *alertP99,
			MaxErrorRate: *alertErrorRate,
			MinReqPerSec: *alertMinRPS,
		},
	}

	tlsConfig := &tls.Config{}
	if *insecure {
		tlsConfig.InsecureSkipVerify = true
		fmt.Println("  [!] TLS verification disabled (--insecure)")
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
		os.Exit(1)
	}

	fmt.Println("\n[*] Warming up...")
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

	phaseHeader := func(name string) {
		fmt.Printf("\n+==================================================+\n")
		fmt.Printf("|  %-47s |\n", name)
		fmt.Printf("+==================================================+\n")
	}
	tableHeader := func() {
		fmt.Println("+------------------------+------+---------+--------+--------+--------+--------+---------+")
		fmt.Println("| Scenario               | Conc | Success%| Avg ms | P50 ms | P95 ms | P99 ms |  req/s  |")
		fmt.Println("+------------------------+------+---------+--------+--------+--------+--------+---------+")
	}
	tableRow := func(r config.PhaseResult) {
		fmt.Printf("| %-22s | %4d | %5.1f%%  | %6.0f | %6.0f | %6.0f | %6.0f | %7.1f |\n",
			r.PhaseName, r.Concurrency, r.SuccessRate, r.AvgMs, r.P50Ms, r.P95Ms, r.P99Ms, r.ReqPerSec)
	}
	tableFooter := func() {
		fmt.Println("+------------------------+------+---------+--------+--------+--------+--------+---------+")
	}

	phaseHeader("PHASE 1: RAMP-UP")
	tableHeader()
	for i, step := range cfg.RampSteps() {
		if i > 0 {
			time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
		}
		r := runner.RunPhase(step, targetURL, client, endpoints, cfg.Headers)
		results = append(results, r)
		tableRow(r)
		time.Sleep(2 * time.Second)
	}
	tableFooter()

	phaseHeader("PHASE 2: SUSTAINED")
	tableHeader()
	for _, step := range cfg.SustainedSteps() {
		r := runner.RunPhase(step, targetURL, client, endpoints, cfg.Headers)
		results = append(results, r)
		tableRow(r)
		time.Sleep(2 * time.Second)
	}
	tableFooter()

	phaseHeader("PHASE 3: SPIKE")
	tableHeader()
	for _, step := range cfg.SpikeSteps() {
		r := runner.RunPhase(step, targetURL, client, endpoints, cfg.Headers)
		results = append(results, r)
		tableRow(r)
		time.Sleep(3 * time.Second)
	}
	tableFooter()

	phaseHeader("PHASE 4: RECOVERY")
	fmt.Println("[*] Waiting 10 seconds post-spike...")
	time.Sleep(10 * time.Second)
	tableHeader()
	r := runner.RunPhase(cfg.RecoveryStep(), targetURL, client, endpoints, cfg.Headers)
	results = append(results, r)
	tableRow(r)
	tableFooter()

	printErrorSummary(results)

	if cfg.Format == "json" || cfg.Format == "both" {
		if writeErr := report.WriteJSON(targetURL, cfg.OutputDir, results, endpoints); writeErr != nil {
			fmt.Fprintf(os.Stderr, "[!] %v\n", writeErr)
		}
	}
	if cfg.Format == "html" || cfg.Format == "both" {
		if writeErr := report.WriteHTML(targetURL, cfg.OutputDir, results, endpoints); writeErr != nil {
			fmt.Fprintf(os.Stderr, "[!] %v\n", writeErr)
		}
	}

	ihlaller := alert.Check(results, cfg.Alert)
	if len(ihlaller) > 0 {
		fmt.Println("\n===========================================================")
		fmt.Println("                   ALERT — THRESHOLD BREACHED")
		fmt.Println("===========================================================")
		for _, ih := range ihlaller {
			fmt.Fprintf(os.Stderr, "  [!] %s\n", ih.Message)
		}
		fmt.Println("===========================================================")
		os.Exit(1)
	}

	fmt.Println("===========================================================")
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
	fmt.Println("\n===========================================================")
	fmt.Println("                    ERROR ANALYSIS")
	fmt.Println("===========================================================")

	var tTo, tRs, tRf, t5, t4, tOt, tFail, tAll int
	for _, r := range results {
		tTo += r.Errors.Timeout
		tRs += r.Errors.ConnectionReset
		tRf += r.Errors.ConnectionRefused
		t5 += r.Errors.HTTP5xx
		t4 += r.Errors.HTTP4xx
		tOt += r.Errors.Other
		tFail += r.Failed
		tAll += r.TotalRequests
	}

	if tFail > 0 && tAll > 0 {
		fmt.Printf("  Total Requests    : %d | Failed: %d (%.1f%% of total)\n", tAll, tFail, float64(tFail)/float64(tAll)*100)
		fmt.Println("  -- Failure Breakdown --")
		pf := func(label string, count int) {
			fmt.Printf("  %-20s: %d (%.1f%% of failures, %.1f%% of total)\n",
				label, count, float64(count)/float64(tFail)*100, float64(count)/float64(tAll)*100)
		}
		pf("Timeout", tTo)
		pf("Connection Reset", tRs)
		pf("Connection Refused", tRf)
		pf("HTTP 5xx", t5)
		pf("HTTP 4xx", t4)
		pf("Other", tOt)
	} else {
		fmt.Println("  No errors recorded.")
	}

	fmt.Println("\n  HTTP Status Distribution:")
	dist := make(map[int]int)
	for _, r := range results {
		for code, count := range r.StatusDist {
			dist[code] += count
		}
	}
	for code, count := range dist {
		fmt.Printf("    %d : %d requests\n", code, count)
	}
}
