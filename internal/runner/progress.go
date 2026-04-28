package runner

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/umutozen/stormprobe/internal/config"
)

const (
	barWidth  = 20
	barFilled = '█'
	barEmpty  = '░'
)

func showProgress(step config.PhaseStep, successful, failed *int32, start time.Time, done chan struct{}) {
	ticker := time.NewTicker(progressInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			s := atomic.LoadInt32(successful)
			f := atomic.LoadInt32(failed)
			total := s + f
			elapsed := time.Since(start)

			rps := 0.0
			if elapsed.Seconds() > 0.5 {
				rps = float64(total) / elapsed.Seconds()
			}

			if step.Duration > 0 {
				pct := elapsed.Seconds() / step.Duration.Seconds()
				if pct > 1.0 {
					pct = 1.0
				}
				filled := int(pct * barWidth)
				bar := make([]rune, barWidth)
				for i := range bar {
					if i < filled {
						bar[i] = barFilled
					} else {
						bar[i] = barEmpty
					}
				}
				fmt.Printf("\r  ⏳ %s %3.0f%% | %d req | %.1f req/s | %s elapsed",
					string(bar), pct*100, total, rps, fmtDuration(elapsed))
			} else {
				expected := int32(step.Concurrency * step.ReqPerWorker)
				pct := 0.0
				if expected > 0 {
					pct = float64(total) / float64(expected)
					if pct > 1.0 {
						pct = 1.0
					}
				}
				filled := int(pct * barWidth)
				bar := make([]rune, barWidth)
				for i := range bar {
					if i < filled {
						bar[i] = barFilled
					} else {
						bar[i] = barEmpty
					}
				}
				fmt.Printf("\r  ⏳ %s %3.0f%% | %d/%d req | %.1f req/s",
					string(bar), pct*100, total, expected, rps)
			}
		}
	}
}

func fmtDuration(d time.Duration) string {
	s := int(d.Seconds())
	if s < 60 {
		return fmt.Sprintf("%ds", s)
	}
	return fmt.Sprintf("%dm%02ds", s/60, s%60)
}
