package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/umutozen/stormprobe/internal/config"
)

func WriteJSON(target, outputDir string, results []config.PhaseResult, endpoints []string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output dir: %w", err)
	}

	filename := fmt.Sprintf("stormprobe_report_%s.json", time.Now().Format("20060102_150405"))
	path := filepath.Join(outputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create JSON report: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(map[string]any{
		"target":               target,
		"timestamp":            time.Now().Format(time.RFC3339),
		"tool_version":         config.Version,
		"discovered_endpoints": endpoints,
		"phases":               results,
	})
	if err != nil {
		return err
	}

	fmt.Printf("[+] JSON report: %s\n", path)
	return nil
}
