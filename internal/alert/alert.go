package alert

import (
	"fmt"

	"github.com/umutozen/stormprobe/internal/config"
)

type Ihlal struct {
	Message string
}

func Check(sonuclar []config.PhaseResult, cfg config.AlertConfig) []Ihlal {
	if cfg.MaxP99Ms == 0 && cfg.MaxErrorRate == 0 && cfg.MinReqPerSec == 0 {
		return nil
	}

	var ihlaller []Ihlal

	for _, r := range sonuclar {
		if cfg.MaxP99Ms > 0 && r.P99Ms > cfg.MaxP99Ms {
			ihlaller = append(ihlaller, Ihlal{
				Message: fmt.Sprintf("[%s] P99 %.0fms > threshold %.0fms", r.PhaseName, r.P99Ms, cfg.MaxP99Ms),
			})
		}

		if cfg.MaxErrorRate > 0 {
			errorRate := 0.0
			if r.TotalRequests > 0 {
				errorRate = float64(r.Failed) / float64(r.TotalRequests) * 100
			}
			if errorRate > cfg.MaxErrorRate {
				ihlaller = append(ihlaller, Ihlal{
					Message: fmt.Sprintf("[%s] Error rate %.1f%% > threshold %.1f%%", r.PhaseName, errorRate, cfg.MaxErrorRate),
				})
			}
		}

		if cfg.MinReqPerSec > 0 && r.ReqPerSec < cfg.MinReqPerSec {
			ihlaller = append(ihlaller, Ihlal{
				Message: fmt.Sprintf("[%s] req/s %.1f < min threshold %.1f", r.PhaseName, r.ReqPerSec, cfg.MinReqPerSec),
			})
		}
	}

	return ihlaller
}
