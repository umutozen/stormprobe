package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/umutozen/stormprobe/internal/compare"
	"github.com/umutozen/stormprobe/internal/config"
)

func WriteCompareJSON(result compare.CompareResult, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create output dir: %w", err)
	}

	filename := fmt.Sprintf("stormprobe_compare_%s.json", time.Now().Format("20060102_150405"))
	path := filepath.Join(outputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("cannot create compare JSON report: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	err = enc.Encode(map[string]any{
		"report_type":  "compare",
		"timestamp":    time.Now().Format(time.RFC3339),
		"tool_version": config.Version,
		"outcome":      result.Outcome,
		"score":        result.Score,
		"why":          result.Why,
		"warnings":     result.Warnings,
		"before": map[string]any{
			"report_id":        result.Before.ReportID,
			"target":           result.Before.Target,
			"timestamp":        result.Before.Timestamp,
			"total_requests":   result.Before.TotalRequests,
			"total_failed":     result.Before.TotalFailed,
			"rating":           result.Before.Verdict.Rating,
			"confidence":       result.Before.Verdict.Confidence,
			"safe_concurrency": result.Before.Verdict.SafeConcurrency,
			"recovery_ok":      result.Before.Verdict.RecoveryOK,
			"bottleneck_cause": result.Before.Verdict.BottleneckCause,
		},
		"after": map[string]any{
			"report_id":        result.After.ReportID,
			"target":           result.After.Target,
			"timestamp":        result.After.Timestamp,
			"total_requests":   result.After.TotalRequests,
			"total_failed":     result.After.TotalFailed,
			"rating":           result.After.Verdict.Rating,
			"confidence":       result.After.Verdict.Confidence,
			"safe_concurrency": result.After.Verdict.SafeConcurrency,
			"recovery_ok":      result.After.Verdict.RecoveryOK,
			"bottleneck_cause": result.After.Verdict.BottleneckCause,
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Printf("[+] Compare JSON report: %s\n", path)
	return path, nil
}
