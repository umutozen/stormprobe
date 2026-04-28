package config

import (
	"fmt"
	"time"
)

const Version = "1.0.0"

type PhaseStep struct {
	Name         string
	Concurrency  int
	ReqPerWorker int
	Duration     time.Duration
}

type AlertConfig struct {
	MaxP99Ms     float64
	MaxErrorRate float64
	MinReqPerSec float64
}

type ErrorCounts struct {
	Timeout           int `json:"timeout"`
	ConnectionReset   int `json:"connection_reset"`
	ConnectionRefused int `json:"connection_refused"`
	HTTP5xx           int `json:"http_5xx"`
	HTTP4xx           int `json:"http_4xx"`
	BodyMismatch      int `json:"body_mismatch"`
	Other             int `json:"other"`
}

type PhaseResult struct {
	PhaseName     string      `json:"phase_name"`
	Concurrency   int         `json:"concurrency"`
	DurationSec   float64     `json:"duration_sec"`
	TotalRequests int         `json:"total_requests"`
	Successful    int         `json:"successful"`
	Failed        int         `json:"failed"`
	SuccessRate   float64     `json:"success_rate"`
	AvgMs         float64     `json:"avg_ms"`
	P50Ms         float64     `json:"p50_ms"`
	P95Ms         float64     `json:"p95_ms"`
	P99Ms         float64     `json:"p99_ms"`
	ReqPerSec     float64     `json:"req_per_sec"`
	Errors        ErrorCounts `json:"errors"`
	StatusDist    map[int]int `json:"status_distribution"`
}

type Config struct {
	Target               string
	RequestTimeout       time.Duration
	EndpointsFile        string
	NoDiscovery          bool
	Insecure             bool
	KatanaPath           string
	HttpxPath            string
	OutputDir            string
	Format               string
	RampPeak             int
	SustainedConc        int
	SpikePeak            int
	ReqPerWorker         int
	Headers              map[string]string
	PhaseDuration        time.Duration
	Alert                AlertConfig
	ConfigFile           string
	SuccessMatch         string
	FailMatch            string
	AutoProfile          bool
	Profile              string
	NoFingerprint        bool
	ProfileOverrideAudit string
}

var UserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/125.0.0.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) AppleWebKit/605.1.15 Mobile",
	"Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 Chrome/125.0.0.0 Mobile",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15",
}

func (c Config) RampSteps() []PhaseStep {
	return []PhaseStep{
		{Name: "Ramp L1 (5 vu)", Concurrency: 5, ReqPerWorker: c.ReqPerWorker, Duration: c.PhaseDuration},
		{Name: "Ramp L2 (15 vu)", Concurrency: 15, ReqPerWorker: c.ReqPerWorker, Duration: c.PhaseDuration},
		{Name: "Ramp L3 (30 vu)", Concurrency: 30, ReqPerWorker: c.ReqPerWorker, Duration: c.PhaseDuration},
		{Name: fmt.Sprintf("Ramp L4 (%d vu)", c.RampPeak), Concurrency: c.RampPeak, ReqPerWorker: c.ReqPerWorker, Duration: c.PhaseDuration},
	}
}

func (c Config) SustainedSteps() []PhaseStep {
	steps := make([]PhaseStep, 3)
	for i := range steps {
		steps[i] = PhaseStep{
			Name:         fmt.Sprintf("Sustained wave #%d", i+1),
			Concurrency:  c.SustainedConc,
			ReqPerWorker: c.ReqPerWorker,
			Duration:     c.PhaseDuration,
		}
	}
	return steps
}

func (c Config) SpikeSteps() []PhaseStep {
	start := int(float64(c.SpikePeak) * 0.6)
	if start < 10 {
		start = 10
	}
	return []PhaseStep{
		{Name: fmt.Sprintf("Spike start (%d vu)", start), Concurrency: start, ReqPerWorker: c.ReqPerWorker, Duration: c.PhaseDuration},
		{Name: fmt.Sprintf("Spike peak (%d vu)", c.SpikePeak), Concurrency: c.SpikePeak, ReqPerWorker: c.ReqPerWorker, Duration: c.PhaseDuration},
		{Name: "Spike cooldown (50 vu)", Concurrency: 50, ReqPerWorker: c.ReqPerWorker, Duration: c.PhaseDuration},
	}
}

func (c Config) RecoveryStep() PhaseStep {
	return PhaseStep{Name: "Recovery check", Concurrency: 10, ReqPerWorker: 10}
}
