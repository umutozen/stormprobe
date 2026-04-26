package alert

import (
	"fmt"

	"github.com/umutozen/stormprobe/internal/config"
)

type Ihlal struct {
	Mesaj string
}

func Check(sonuclar []config.PhaseResult, cfg config.AlertConfig) []Ihlal {
	if cfg.MaxP99Ms == 0 && cfg.MaxErrorRate == 0 && cfg.MinReqPerSec == 0 {
		return nil
	}

	var ihlaller []Ihlal

	for _, r := range sonuclar {
		if cfg.MaxP99Ms > 0 && r.P99Ms > cfg.MaxP99Ms {
			ihlaller = append(ihlaller, Ihlal{
				Mesaj: fmt.Sprintf("[%s] P99 %.0fms > eşik %.0fms", r.PhaseName, r.P99Ms, cfg.MaxP99Ms),
			})
		}

		if cfg.MaxErrorRate > 0 {
			errorRate := 0.0
			if r.TotalRequests > 0 {
				errorRate = float64(r.Failed) / float64(r.TotalRequests) * 100
			}
			if errorRate > cfg.MaxErrorRate {
				ihlaller = append(ihlaller, Ihlal{
					Mesaj: fmt.Sprintf("[%s] Hata oranı %.1f%% > eşik %.1f%%", r.PhaseName, errorRate, cfg.MaxErrorRate),
				})
			}
		}

		if cfg.MinReqPerSec > 0 && r.ReqPerSec < cfg.MinReqPerSec {
			ihlaller = append(ihlaller, Ihlal{
				Mesaj: fmt.Sprintf("[%s] req/s %.1f < min eşik %.1f", r.PhaseName, r.ReqPerSec, cfg.MinReqPerSec),
			})
		}
	}

	return ihlaller
}
