package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type TargetEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FileConfig struct {
	Target         string        `json:"target"`
	Targets        []TargetEntry `json:"targets"`
	Insecure       *bool         `json:"insecure"`
	Format         *string       `json:"format"`
	Output         *string       `json:"output"`
	Duration       *string       `json:"duration"`
	RequestTimeout *string       `json:"request_timeout"`
	NoDiscovery    *bool         `json:"no_discovery"`
	EndpointsFile  *string       `json:"endpoints"`
	Headers        []string      `json:"headers"`

	AlertP99       *float64 `json:"alert_p99"`
	AlertErrorRate *float64 `json:"alert_error_rate"`
	AlertMinRPS    *float64 `json:"alert_rps"`

	ConcurrencyRamp      *int `json:"concurrency_ramp"`
	ConcurrencySustained *int `json:"concurrency_sustained"`
	ConcurrencySpike     *int `json:"concurrency_spike"`
	ReqPerWorker         *int `json:"req_per_worker"`

	KatanaPath *string `json:"katana_path"`
	HttpxPath  *string `json:"httpx_path"`

	SuccessMatch *string `json:"success_match"`
	FailMatch    *string `json:"fail_match"`
}

func LoadFile(path string) (*FileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var fc FileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("invalid config file syntax: %w", err)
	}

	return &fc, nil
}

func MergeFileIntoConfig(fc *FileConfig, cfg *Config, flagsSet map[string]bool) {
	if fc == nil {
		return
	}

	if fc.Target != "" && !flagsSet["target"] && cfg.Target == "" {
		cfg.Target = fc.Target
	}

	if fc.Insecure != nil && !flagsSet["insecure"] {
		cfg.Insecure = *fc.Insecure
	}
	if fc.NoDiscovery != nil && !flagsSet["no-discovery"] {
		cfg.NoDiscovery = *fc.NoDiscovery
	}
	if fc.Format != nil && !flagsSet["format"] {
		cfg.Format = *fc.Format
	}
	if fc.Output != nil && !flagsSet["output"] {
		cfg.OutputDir = *fc.Output
	}
	if fc.EndpointsFile != nil && !flagsSet["endpoints"] {
		cfg.EndpointsFile = *fc.EndpointsFile
	}
	if fc.RequestTimeout != nil && !flagsSet["timeout"] {
		if d, err := time.ParseDuration(*fc.RequestTimeout); err == nil {
			cfg.RequestTimeout = d
		}
	}
	if fc.Duration != nil && !flagsSet["duration"] {
		if *fc.Duration != "" {
			if d, err := time.ParseDuration(*fc.Duration); err == nil {
				cfg.PhaseDuration = d
			}
		}
	}

	if len(fc.Headers) > 0 && !flagsSet["header"] && !flagsSet["H"] {
		for _, h := range fc.Headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				if key != "" {
					cfg.Headers[key] = val
				}
			}
		}
	}

	if fc.AlertP99 != nil && !flagsSet["alert-p99"] {
		cfg.Alert.MaxP99Ms = *fc.AlertP99
	}
	if fc.AlertErrorRate != nil && !flagsSet["alert-error-rate"] {
		cfg.Alert.MaxErrorRate = *fc.AlertErrorRate
	}
	if fc.AlertMinRPS != nil && !flagsSet["alert-rps"] {
		cfg.Alert.MinReqPerSec = *fc.AlertMinRPS
	}

	if fc.ConcurrencyRamp != nil && !flagsSet["concurrency-ramp"] {
		cfg.RampPeak = *fc.ConcurrencyRamp
	}
	if fc.ConcurrencySustained != nil && !flagsSet["concurrency-sustained"] {
		cfg.SustainedConc = *fc.ConcurrencySustained
	}
	if fc.ConcurrencySpike != nil && !flagsSet["concurrency-spike"] {
		cfg.SpikePeak = *fc.ConcurrencySpike
	}
	if fc.ReqPerWorker != nil && !flagsSet["req-per-worker"] {
		cfg.ReqPerWorker = *fc.ReqPerWorker
	}

	if fc.KatanaPath != nil && !flagsSet["katana-path"] {
		cfg.KatanaPath = *fc.KatanaPath
	}
	if fc.HttpxPath != nil && !flagsSet["httpx-path"] {
		cfg.HttpxPath = *fc.HttpxPath
	}
	if fc.SuccessMatch != nil && !flagsSet["success-match"] {
		cfg.SuccessMatch = *fc.SuccessMatch
	}
	if fc.FailMatch != nil && !flagsSet["fail-match"] {
		cfg.FailMatch = *fc.FailMatch
	}
}
