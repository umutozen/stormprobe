package report

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/umutozen/stormprobe/internal/config"
	"github.com/umutozen/stormprobe/internal/verdict"
)

var safeNameRe = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func sanitizeLabel(label string) string {
	label = strings.ToLower(strings.TrimSpace(label))
	label = strings.ReplaceAll(label, " ", "_")
	return safeNameRe.ReplaceAllString(label, "")
}

func WriteJSON(target, outputDir string, results []config.PhaseResult, endpoints []string, v verdict.Verdict) error {
	_, err := WriteJSONReturn(target, outputDir, results, endpoints, v)
	return err
}

func WriteJSONReturn(target, outputDir string, results []config.PhaseResult, endpoints []string, v verdict.Verdict) (string, error) {
	return WriteJSONWithLabel(target, outputDir, "", results, endpoints, v)
}

func WriteJSONWithLabel(target, outputDir, label string, results []config.PhaseResult, endpoints []string, v verdict.Verdict) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create output dir: %w", err)
	}

	var filename string
	if label != "" {
		filename = fmt.Sprintf("stormprobe_report_%s.json", sanitizeLabel(label))
	} else {
		filename = fmt.Sprintf("stormprobe_report_%s.json", time.Now().Format("20060102_150405"))
	}
	path := filepath.Join(outputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("cannot create JSON report: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	now := time.Now()
	host := target
	if u, err2 := url.Parse(target); err2 == nil && u.Hostname() != "" {
		host = u.Hostname()
	}
	reportID := fmt.Sprintf("%s-%s", now.Format("20060102-150405"), host)

	totalRequests, totalFailed := 0, 0
	for _, r := range results {
		totalRequests += r.TotalRequests
		totalFailed += r.Failed
	}

	err = enc.Encode(map[string]any{
		"report_id":            reportID,
		"target":               target,
		"timestamp":            now.Format(time.RFC3339),
		"tool_version":         config.Version,
		"total_requests":       totalRequests,
		"total_failed":         totalFailed,
		"discovered_endpoints": endpoints,
		"phases":               results,
		"verdict": map[string]any{
			"rating":            v.Rating,
			"confidence":        v.Confidence,
			"safe_concurrency":  v.SafeConcurrency,
			"degradation_at_vu": v.DegradationAt,
			"failure_point_vu":  v.FailurePoint,
			"bottleneck_cause":  v.BottleneckCause,
			"bottleneck_detail": v.BottleneckDetail,
			"priority_checks":   v.PriorityChecks,
			"recovery_ok":       v.RecoveryOK,
			"recommendation":    v.Recommendation,
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Printf("[+] JSON report: %s\n", path)
	return path, nil
}
