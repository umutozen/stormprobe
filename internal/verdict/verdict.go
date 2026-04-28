package verdict

import (
	"fmt"
	"strings"

	"github.com/umutozen/stormprobe/internal/config"
)

const (
	p99DegradedMs          = 2000.0
	p99CriticalMs          = 5000.0
	errorRateSafe          = 1.0
	errorRateFail          = 2.0
	recoveryRatio          = 1.5
	minorErrorThreshold    = 0.05
	criticalErrorThreshold = 1.0
)

type Rating string

const (
	RatingHealthy  Rating = "Healthy"
	RatingDegraded Rating = "Degraded Under Load"
	RatingCritical Rating = "Not Production Ready"
	RatingUnstable Rating = "High Risk"
)

type Verdict struct {
	Rating                   Rating
	Confidence               string
	SafeConcurrency          int
	DegradationAt            *int
	FailurePoint             *int
	ErrorZoneAt              int
	HasErrors                bool
	ThroughputPlateau        float64
	BottleneckCause          string
	BottleneckDetail         string
	PriorityChecks           []string
	RecoveryOK               bool
	Recommendation           string
	ProductionRecommendation string
	StackSummary             string
	PrimaryCause             string
	SecondaryContributors    []string
	WhyDiagnosis             []string
	ProfileOverrideAudit     string
}

func Analyze(results []config.PhaseResult, stackSummary, profileOverrideAudit string) Verdict {
	if len(results) == 0 {
		return Verdict{Rating: RatingHealthy, Recommendation: "No data to analyze."}
	}

	var (
		safeConcurrency int
		degradationAt   int
		failurePoint    int
		errorZoneAt     int
		baseline        config.PhaseResult
		recoveryResult  config.PhaseResult
		hasRecovery     bool
	)

	baseline = results[0]

	for _, r := range results {
		errRate := errorRate(r)
		if r.PhaseName == "Recovery check" {
			hasRecovery = true
			recoveryResult = r
			continue
		}

		if errRate < errorRateSafe && r.P99Ms < p99DegradedMs {
			if r.Concurrency > safeConcurrency {
				safeConcurrency = r.Concurrency
			}
		}

		if degradationAt == 0 && r.P99Ms >= p99DegradedMs {
			v := r.Concurrency
			degradationAt = v
		}

		if failurePoint == 0 && r.P99Ms >= p99CriticalMs {
			v := r.Concurrency
			failurePoint = v
		}

		if errorZoneAt == 0 && errRate >= errorRateFail {
			errorZoneAt = r.Concurrency
		}
	}

	recoveryOK := !hasRecovery || isRecoveryHealthy(recoveryResult, baseline)
	hasErrors := totalErrorCount(results) > 0

	cause, detail := diagnoseBottleneck(results)
	primary, secondaries := buildStackDiagnosis(cause, stackSummary)
	whyDiagnosis := buildWhyDiagnosis(results, degradationAt, failurePoint)
	priority := priorityChecks(cause)
	rating := computeRating(failurePoint, degradationAt, recoveryOK, hasErrors)
	recommendation := buildRecommendation(rating, cause, safeConcurrency)

	var degradationPtr, failurePtr *int
	if degradationAt > 0 {
		v := degradationAt
		degradationPtr = &v
	}
	if failurePoint > 0 {
		v := failurePoint
		failurePtr = &v
	}

	return Verdict{
		Rating:                   rating,
		Confidence:               confidenceLabel(rating, cause, recoveryOK),
		SafeConcurrency:          safeConcurrency,
		DegradationAt:            degradationPtr,
		FailurePoint:             failurePtr,
		ErrorZoneAt:              errorZoneAt,
		HasErrors:                hasErrors,
		ThroughputPlateau:        throughputPlateau(results, degradationAt),
		BottleneckCause:          cause,
		BottleneckDetail:         detail,
		PriorityChecks:           priority,
		RecoveryOK:               recoveryOK,
		Recommendation:           recommendation,
		ProductionRecommendation: productionRecommendation(rating),
		StackSummary:             stackSummary,
		PrimaryCause:             primary,
		SecondaryContributors:    secondaries,
		WhyDiagnosis:             whyDiagnosis,
		ProfileOverrideAudit:     profileOverrideAudit,
	}
}

const (
	recoveryMaxErrorRate = 1.0
	recoveryAvgRatio     = 2.0
	recoveryP95Ratio     = 2.0
)

func isRecoveryHealthy(recovery, baseline config.PhaseResult) bool {
	if errorRate(recovery) >= recoveryMaxErrorRate {
		return false
	}
	if baseline.AvgMs > 0 && recovery.AvgMs > baseline.AvgMs*recoveryAvgRatio {
		return false
	}
	if baseline.P95Ms > 0 && recovery.P95Ms > baseline.P95Ms*recoveryP95Ratio {
		return false
	}
	return true
}

func errorRate(r config.PhaseResult) float64 {
	if r.TotalRequests == 0 {
		return 0
	}
	return float64(r.Failed) / float64(r.TotalRequests) * 100
}

func totalErrorCount(results []config.PhaseResult) int {
	var total int
	for _, r := range results {
		total += r.Errors.Timeout + r.Errors.ConnectionReset +
			r.Errors.ConnectionRefused + r.Errors.HTTP5xx + r.Errors.HTTP4xx
	}
	return total
}

func throughputPlateau(results []config.PhaseResult, degradationAt int) float64 {
	if degradationAt == 0 {
		return 0
	}
	var sum float64
	var count int
	for _, r := range results {
		if r.PhaseName == "Recovery check" {
			continue
		}
		if r.Concurrency >= degradationAt {
			sum += r.ReqPerSec
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

func overallErrorRate(results []config.PhaseResult) float64 {
	var totalReq, totalFailed int
	for _, r := range results {
		if r.PhaseName == "Recovery check" {
			continue
		}
		totalReq += r.TotalRequests
		totalFailed += r.Failed
	}
	if totalReq == 0 {
		return 0
	}
	return float64(totalFailed) / float64(totalReq) * 100
}

func diagnoseBottleneck(results []config.PhaseResult) (cause, detail string) {
	rate := overallErrorRate(results)
	if rate >= criticalErrorThreshold {
		return diagnoseErrorBottleneck(results)
	}
	if rate >= minorErrorThreshold {
		return diagnoseTransient(results)
	}
	if cause, detail, ok := detectSessionLockPattern(results); ok {
		return cause, detail
	}
	return diagnoseThroughputBottleneck(results)
}

func diagnoseTransient(results []config.PhaseResult) (cause, detail string) {
	var timeouts, http5xx int
	for _, r := range results {
		timeouts += r.Errors.Timeout
		http5xx += r.Errors.HTTP5xx
	}
	return "Minor transient failures under peak load",
		fmt.Sprintf("Small number of timeouts (%d) and HTTP 5xx responses (%d) observed at peak concurrency, "+
			"but throughput continued scaling and recovery remained healthy.", timeouts, http5xx)
}

func diagnoseErrorBottleneck(results []config.PhaseResult) (cause, detail string) {
	var timeouts, resets, refused, http5xx int
	for _, r := range results {
		timeouts += r.Errors.Timeout
		resets += r.Errors.ConnectionReset
		refused += r.Errors.ConnectionRefused
		http5xx += r.Errors.HTTP5xx
	}
	total := timeouts + resets + refused + http5xx

	switch {
	case isdominant(timeouts, total):
		return "Application/DB processing bottleneck",
			"Requests queue up server-side under load. Likely cause: slow DB queries, blocking I/O, or insufficient worker threads."
	case isDominant2(resets, refused, total):
		return "Connection pool / infrastructure capacity exhausted",
			"Server is dropping connections before processing. Check reverse proxy limits, OS socket backlog, and connection pool size."
	case isdominant(http5xx, total):
		return "Server-side application errors under load",
			"HTTP 5xx responses indicate application crashes or unhandled exceptions at high concurrency. Review application logs."
	default:
		return "Mixed failure mode — likely cascading overload",
			fmt.Sprintf("Multiple error types: Timeouts: %d | Resets: %d | Refused: %d | 5xx: %d", timeouts, resets, refused, http5xx)
	}
}

func diagnoseThroughputBottleneck(results []config.PhaseResult) (cause, detail string) {
	var sane []config.PhaseResult
	for _, r := range results {
		if r.PhaseName != "Recovery check" {
			sane = append(sane, r)
		}
	}

	if len(sane) < 2 {
		return "None", "Insufficient phase data."
	}

	first := sane[0]
	last := sane[len(sane)-1]

	concurrencyGrowth := float64(last.Concurrency) / float64(first.Concurrency)
	throughputGrowth := last.ReqPerSec / first.ReqPerSec
	latencyGrowth := last.P99Ms / first.P99Ms

	if latencyGrowth > 3 && throughputGrowth < concurrencyGrowth*0.5 {
		return "Throughput saturation",
			fmt.Sprintf("Latency grew %.1fx while throughput grew only %.1fx for %.1fx concurrency increase. "+
				"Backend resources are fully contended — adding more users queues requests rather than increasing output.",
				latencyGrowth, throughputGrowth, concurrencyGrowth)
	}

	if latencyGrowth > 5 {
		return "Latency explosion under concurrency",
			fmt.Sprintf("P99 grew %.1fx from baseline with no errors. Suggests single-threaded bottleneck or severe lock contention.", latencyGrowth)
	}

	return "No sustained bottleneck detected", "System scaled acceptably under test load with no significant latency degradation or throughput plateau."
}

const (
	sessionLockP50P99Ratio    = 1.4
	sessionLockThroughputCap  = 0.6
	sessionLockMinConcurrency = 3
)

func detectSessionLockPattern(results []config.PhaseResult) (cause, detail string, detected bool) {
	var sane []config.PhaseResult
	for _, r := range results {
		if r.PhaseName != "Recovery check" && r.P50Ms > 0 && r.P99Ms > 0 {
			sane = append(sane, r)
		}
	}
	if len(sane) < sessionLockMinConcurrency {
		return "", "", false
	}

	// Session lock pattern: P99 yakın P50'ye → serialize edilmiş kuyruk
	// Tüm fazlarda bu pattern tutarlıysa session lock var
	serializedCount := 0
	for _, r := range sane {
		ratio := r.P99Ms / r.P50Ms
		if ratio < sessionLockP50P99Ratio {
			serializedCount++
		}
	}

	if serializedCount < len(sane)/2 {
		return "", "", false
	}

	// Throughput concurrency'e göre büyümüyor mu?
	first := sane[0]
	last := sane[len(sane)-1]
	if first.ReqPerSec == 0 || first.Concurrency == 0 {
		return "", "", false
	}

	cuGrowth := float64(last.Concurrency) / float64(first.Concurrency)
	tpGrowth := last.ReqPerSec / first.ReqPerSec

	if tpGrowth >= cuGrowth*sessionLockThroughputCap {
		return "", "", false
	}

	cause = "Session lock contention — serialized request processing"
	detail = fmt.Sprintf(
		"P99/P50 latency ratio is consistent (%.1f-%.1f×) across concurrency levels, "+
			"indicating requests are being serialized rather than processed in parallel. "+
			"Likely cause: ASP.NET exclusive session lock or single-threaded critical section. "+
			"Throughput grew only %.1f× for %.1f× concurrency increase.",
		sane[0].P99Ms/sane[0].P50Ms, sane[len(sane)-1].P99Ms/sane[len(sane)-1].P50Ms,
		tpGrowth, cuGrowth,
	)
	return cause, detail, true
}

func isdominant(count, total int) bool {
	return total > 0 && float64(count)/float64(total) >= 0.60
}

func isDominant2(a, b, total int) bool {
	return total > 0 && float64(a+b)/float64(total) >= 0.60
}

func buildStackDiagnosis(cause, stackSummary string) (primary string, secondaries []string) {
	if stackSummary == "" {
		return cause, nil
	}

	stack := strings.ToLower(stackSummary)

	switch {
	case strings.Contains(stack, "iis") || strings.Contains(stack, "asp.net"):
		switch cause {
		case "Application/DB processing bottleneck":
			return "IIS worker process queue saturation",
				[]string{"SQL connection pool contention", "Session lock contention", "Synchronous blocking I/O on hot paths"}
		case "Throughput saturation":
			return "IIS App Pool worker limit reached",
				[]string{"DB connection pool pressure", "Session state locking"}
		case "Connection pool / infrastructure capacity exhausted":
			return "IIS connection queue (AcceptEx backlog) overflow",
				[]string{"OS TCP socket backlog limit", "App Pool rapid recycle under pressure"}
		case "Session lock contention — serialized request processing":
			return "ASP.NET exclusive session lock detected",
				[]string{
					"Session.Abandon() or async session provider missing",
					"Authenticated endpoints serialize under IIS session mutex",
					"Recommendation: enable read-only session or use distributed cache (Redis)",
				}
		default:
			return cause, []string{"IIS thread pool exhaustion possible under sustained load"}
		}

	case strings.Contains(stack, "tomcat") || strings.Contains(stack, "java"):
		switch cause {
		case "Application/DB processing bottleneck":
			return "Tomcat thread pool starvation (maxThreads limit)",
				[]string{"JDBC connection pool exhaustion", "Servlet blocking I/O wait", "Long GC pauses under heap pressure"}
		case "Throughput saturation":
			return "Tomcat executor queue buildup",
				[]string{"JDBC pool wait time elevated", "Slow downstream service calls"}
		default:
			return cause, []string{"Java thread pool saturation possible under sustained load"}
		}

	case strings.Contains(stack, "php") || strings.Contains(stack, "nginx") || strings.Contains(stack, "apache"):
		switch cause {
		case "Application/DB processing bottleneck":
			return "PHP-FPM worker pool exhaustion",
				[]string{"FastCGI upstream queue pressure", "MySQL connection wait", "OpCache miss under high concurrency"}
		case "Connection pool / infrastructure capacity exhausted":
			return "Nginx worker_connections limit reached",
				[]string{"PHP-FPM max_children saturation", "KeepAlive exhaustion"}
		default:
			return cause, []string{"PHP-FPM pool pressure possible under burst traffic"}
		}
	}

	return cause, nil
}

func buildWhyDiagnosis(results []config.PhaseResult, degradationAt, failurePoint int) []string {
	var reasons []string
	var sane []config.PhaseResult
	for _, r := range results {
		if r.PhaseName != "Recovery check" {
			sane = append(sane, r)
		}
	}
	if len(sane) < 2 {
		return reasons
	}

	first := sane[0]
	last := sane[len(sane)-1]

	if degradationAt > 0 {
		reasons = append(reasons, fmt.Sprintf("Latency spike begins after %d VU", degradationAt))
	}

	tpGrowth := last.ReqPerSec / first.ReqPerSec
	cuGrowth := float64(last.Concurrency) / float64(first.Concurrency)
	if tpGrowth < cuGrowth*0.5 {
		reasons = append(reasons, fmt.Sprintf("Throughput plateaus at ~%.0f req/s despite %.0fx concurrency increase", last.ReqPerSec, cuGrowth))
	}

	if last.P99Ms > 5000 {
		reasons = append(reasons, fmt.Sprintf("P99 latency reaches %.0fms — queue buildup pattern, not hard limit", last.P99Ms))
	}

	var totalTimeouts, total5xx int
	for _, r := range sane {
		totalTimeouts += r.Errors.Timeout
		total5xx += r.Errors.HTTP5xx
	}

	if totalTimeouts > 0 && total5xx == 0 {
		reasons = append(reasons, fmt.Sprintf("%d timeouts with 0 × 5xx — server accepts connections but processing stalls", totalTimeouts))
	}
	if total5xx > 0 {
		reasons = append(reasons, fmt.Sprintf("%d × HTTP 5xx errors — application-level crashes or hard limits hit", total5xx))
	}
	if failurePoint > 0 {
		reasons = append(reasons, fmt.Sprintf("Hard failure threshold at %d VU concurrent", failurePoint))
	}

	return reasons
}

func computeRating(failurePoint, degradationAt int, recoveryOK, hasErrors bool) Rating {
	if failurePoint > 0 && !recoveryOK {
		return RatingUnstable
	}
	if failurePoint > 0 && hasErrors {
		return RatingCritical
	}
	if failurePoint > 0 || degradationAt > 0 {
		return RatingDegraded
	}
	return RatingHealthy
}

func priorityChecks(cause string) []string {
	switch cause {
	case "Minor transient failures under peak load":
		return []string{
			"Monitor transient upstream failures (upstream 502/504 in proxy logs)",
			"Verify per-request timeout thresholds match upstream SLA",
			"Review occasional 502 responses for proxy or CDN layer issues",
		}
	case "Application/DB processing bottleneck":
		return []string{
			"Review slow query logs — identify queries exceeding 1s",
			"Analyze DB connection pool size vs peak concurrency",
			"Check for blocking I/O or synchronous operations on hot paths",
			"Profile thread pool saturation under sustained load",
		}
	case "Connection pool / infrastructure capacity exhausted":
		return []string{
			"Increase reverse proxy (nginx/HAProxy) worker connections",
			"Tune OS-level TCP backlog (net.core.somaxconn)",
			"Review upstream service connection limits",
			"Check keep-alive settings on load balancer",
		}
	case "Server-side application errors under load":
		return []string{
			"Inspect application exception logs during high-concurrency window",
			"Check memory limits — OOM killer may be terminating worker processes",
			"Review rate limiting or circuit breaker configuration",
			"Analyze upstream dependency latency under load",
		}
	case "Throughput saturation", "Latency explosion under concurrency":
		return []string{
			"Analyze database query performance and slow query logs",
			"Check web server thread pool / worker process limits",
			"Review connection pool sizing vs concurrent request volume",
			"Profile application for lock contention or blocking synchronous calls",
		}
	default:
		return []string{
			"Correlate server logs with test timestamps",
			"Monitor CPU, memory, and network saturation during load",
		}
	}
}

func buildRecommendation(r Rating, cause string, safeConcurrency int) string {
	switch r {
	case RatingHealthy:
		return fmt.Sprintf("System is healthy up to %d virtual users. Monitor P99 under real traffic.", safeConcurrency)
	case RatingDegraded:
		return fmt.Sprintf("Cap production traffic at ~%d concurrent users. Priority: Backend resource analysis — %s.", safeConcurrency, cause)
	case RatingCritical:
		return fmt.Sprintf("Do not exceed %d concurrent users without fixes. Immediate action required: performance tuning before production scaling.", safeConcurrency)
	case RatingUnstable:
		return "Server did not fully recover after the spike. High risk of cascading failure under burst traffic."
	default:
		return "Review results manually."
	}
}

func productionRecommendation(r Rating) string {
	switch r {
	case RatingHealthy:
		return "Safe for production at tested concurrency levels"
	case RatingDegraded:
		return "Safe for low-to-medium traffic only"
	case RatingCritical:
		return "Not recommended for production without performance fixes"
	case RatingUnstable:
		return "Do not deploy — risk of cascading failure under burst traffic"
	default:
		return "Evaluate manually before deploying"
	}
}

func confidenceLabel(r Rating, cause string, recoveryOK bool) string {
	switch r {
	case RatingHealthy:
		if cause == "No sustained bottleneck detected" && recoveryOK {
			return "High Confidence"
		}
		if cause == "Minor transient failures under peak load" {
			return "Observed Minor Transient Errors"
		}
		return "Confidence: Normal"
	case RatingDegraded:
		if !recoveryOK {
			return "Recovery Impaired"
		}
		return "Degraded Under Load"
	case RatingCritical:
		return "Low Confidence — Immediate Action Required"
	case RatingUnstable:
		return "Unstable — Recovery Failed"
	default:
		return ""
	}
}
