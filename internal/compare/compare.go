package compare

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type verdictData struct {
	Rating          string `json:"rating"`
	Confidence      string `json:"confidence"`
	SafeConcurrency int    `json:"safe_concurrency"`
	DegradationAtVU *int   `json:"degradation_at_vu"`
	FailurePointVU  *int   `json:"failure_point_vu"`
	RecoveryOK      bool   `json:"recovery_ok"`
	BottleneckCause string `json:"bottleneck_cause"`
	Recommendation  string `json:"recommendation"`
}

type Report struct {
	ReportID      string      `json:"report_id"`
	Target        string      `json:"target"`
	Timestamp     string      `json:"timestamp"`
	TotalRequests int         `json:"total_requests"`
	TotalFailed   int         `json:"total_failed"`
	Verdict       verdictData `json:"verdict"`
}

type CompareResult struct {
	Outcome  string
	Why      []string
	Score    int
	Before   Report
	After    Report
	Warnings []string
}

const (
	outcomeImproved           = "IMPROVED"
	outcomePartialImprovement = "PARTIAL IMPROVEMENT"
	outcomeStable             = "STABLE"
	outcomeStableWithRisk     = "STABLE WITH RISK"
	outcomeDegraded           = "DEGRADED"
	outcomeCriticalRegression = "CRITICAL REGRESSION"

	concImprovementThreshold = 1.05
	concDegradationThreshold = 0.90

	rankHealthy  = 3
	rankDegraded = 2
	rankCritical = 1
	rankUnstable = 0
)

var ratingRank = map[string]int{
	"Healthy":              rankHealthy,
	"Degraded Under Load":  rankDegraded,
	"Not Production Ready": rankCritical,
	"High Risk":            rankUnstable,
}

var confidenceRank = map[string]int{
	"High Confidence":                 3,
	"Confidence: Normal":              2,
	"Observed Minor Transient Errors": 1,
	"Recovery Impaired":               0,
	"Degraded Under Load":             0,
	"":                                0,
	"Low Confidence — Immediate Action Required": -1,
	"Unstable — Recovery Failed":                 -2,
}

func LoadReport(path string) (Report, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Report{}, fmt.Errorf("cannot read file %q: %w", path, err)
	}
	var r Report
	if err := json.Unmarshal(data, &r); err != nil {
		return Report{}, fmt.Errorf("invalid report format in %q: %w", path, err)
	}
	if r.Verdict.Rating == "" {
		return Report{}, fmt.Errorf("file %q does not contain a valid StormProbe verdict", path)
	}
	return r, nil
}

func Analyze(before, after Report) CompareResult {
	var warnings []string

	if before.Target != after.Target {
		warnings = append(warnings, fmt.Sprintf("targets differ (%s vs %s)", before.Target, after.Target))
	}

	if ts := checkTimestampOrder(before.Timestamp, after.Timestamp); ts != "" {
		warnings = append(warnings, ts)
	}

	outcome, score := computeOutcome(before.Verdict, after.Verdict)
	why := buildWhy(before, after)

	return CompareResult{
		Outcome:  outcome,
		Why:      why,
		Score:    score,
		Before:   before,
		After:    after,
		Warnings: warnings,
	}
}

func computeOutcome(b, a verdictData) (string, int) {
	score := 0

	rankBefore := ratingRank[b.Rating]
	rankAfter := ratingRank[a.Rating]

	if (!a.RecoveryOK && b.RecoveryOK) || (rankAfter <= rankUnstable && rankBefore > rankCritical) {
		return outcomeCriticalRegression, -100
	}

	concRatio := safeConcRatio(b.SafeConcurrency, a.SafeConcurrency)

	if concRatio < concDegradationThreshold {
		score += int((concRatio - 1) * 50)
		return outcomeDegraded, score
	}

	if concRatio >= concImprovementThreshold {
		score += int((concRatio - 1) * 50)
	}
	score += (rankAfter - rankBefore) * 20
	score += (confidenceRank[a.Confidence] - confidenceRank[b.Confidence]) * 10
	if a.RecoveryOK && !b.RecoveryOK {
		score += 20
	}
	if !a.RecoveryOK && b.RecoveryOK {
		score -= 20
	}

	switch {
	case score >= 10:
		if rankAfter < rankHealthy {
			return outcomePartialImprovement, score
		}
		return outcomeImproved, score
	case score <= -10 || (!a.RecoveryOK) || confidenceRank[a.Confidence] < confidenceRank[b.Confidence]:
		return outcomeStableWithRisk, score
	default:
		return outcomeStable, score
	}
}

func buildWhy(before, after Report) []string {
	b, a := before.Verdict, after.Verdict
	var reasons []string

	if b.SafeConcurrency > 0 && a.SafeConcurrency > 0 {
		delta := float64(a.SafeConcurrency-b.SafeConcurrency) / float64(b.SafeConcurrency) * 100
		switch {
		case delta > 1:
			reasons = append(reasons, fmt.Sprintf("+ Safe concurrency increased %.0f%% (%d → %d VU)", delta, b.SafeConcurrency, a.SafeConcurrency))
		case delta < -1:
			reasons = append(reasons, fmt.Sprintf("- Safe concurrency decreased %.0f%% (%d → %d VU)", -delta, b.SafeConcurrency, a.SafeConcurrency))
		default:
			reasons = append(reasons, fmt.Sprintf("~ Safe concurrency unchanged (%d VU)", a.SafeConcurrency))
		}
	}

	switch {
	case b.RecoveryOK && !a.RecoveryOK:
		reasons = append(reasons, "- Recovery degraded (now FAILED after spike)")
	case !b.RecoveryOK && a.RecoveryOK:
		reasons = append(reasons, "+ Recovery restored (was FAILED, now OK)")
	default:
		recov := "OK"
		if !a.RecoveryOK {
			recov = "FAILED"
		}
		reasons = append(reasons, fmt.Sprintf("~ Recovery unchanged (%s)", recov))
	}

	if b.Confidence != "" || a.Confidence != "" {
		bc, ac := b.Confidence, a.Confidence
		if bc == "" {
			bc = "N/A"
		}
		if ac == "" {
			ac = "N/A"
		}
		cs := confidenceRank[b.Confidence] - confidenceRank[a.Confidence]
		if cs < 0 {
			reasons = append(reasons, fmt.Sprintf("+ Confidence improved (%s → %s)", bc, ac))
		} else if cs > 0 {
			reasons = append(reasons, fmt.Sprintf("- Confidence degraded (%s → %s)", bc, ac))
		}
	}

	if b.BottleneckCause != a.BottleneckCause {
		reasons = append(reasons, fmt.Sprintf("~ Bottleneck changed: %q → %q", b.BottleneckCause, a.BottleneckCause))
	}

	if before.TotalFailed > 0 && after.TotalFailed > 0 {
		delta := before.TotalFailed - after.TotalFailed
		if delta > 0 {
			reasons = append(reasons, fmt.Sprintf("+ Lower transient error volume observed (%d → %d failed)", before.TotalFailed, after.TotalFailed))
		} else if delta < 0 {
			reasons = append(reasons, fmt.Sprintf("- Higher transient error volume observed (%d → %d failed)", before.TotalFailed, after.TotalFailed))
		}
	}

	return reasons
}

func safeConcRatio(before, after int) float64 {
	if before == 0 {
		return 1.0
	}
	return float64(after) / float64(before)
}

func checkTimestampOrder(beforeTS, afterTS string) string {
	if beforeTS == "" || afterTS == "" {
		return ""
	}
	tb, err1 := time.Parse(time.RFC3339, beforeTS)
	ta, err2 := time.Parse(time.RFC3339, afterTS)
	if err1 != nil || err2 != nil {
		return ""
	}
	if tb.After(ta) {
		return fmt.Sprintf("timestamp order reversed: 'before' (%s) is newer than 'after' (%s)",
			tb.Format("2006-01-02"), ta.Format("2006-01-02"))
	}
	return ""
}
