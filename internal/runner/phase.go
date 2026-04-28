package runner

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/umutozen/stormprobe/internal/config"
	"github.com/umutozen/stormprobe/internal/metrics"
)

const progressInterval = 500 * time.Millisecond

func RunPhase(step config.PhaseStep, targetURL string, client *http.Client, endpoints []string, headers map[string]string, successRe, failRe *regexp.Regexp) config.PhaseResult {
	var successful, failed int32
	var timeouts, resets, refused, http5xx, http4xx, bodyMismatch, other int32

	statusMu := newLock()
	statusDist := make(map[int]int)
	latencyMu := newLock()
	latencies := make([]float64, 0, step.Concurrency*step.ReqPerWorker)

	var wg sync.WaitGroup
	start := time.Now()

	done := make(chan struct{})
	go showProgress(step, &successful, &failed, start, done)

	if step.Duration > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), step.Duration)
		defer cancel()

		for w := 0; w < step.Concurrency; w++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				executeWorkerDuration(ctx, id, client, targetURL, endpoints,
					&successful, &failed, &timeouts, &resets, &refused, &http5xx, &http4xx, &bodyMismatch, &other,
					&latencies, latencyMu, statusDist, statusMu, headers, successRe, failRe)
			}(w)
		}
		wg.Wait()
	} else {
		total := step.Concurrency * step.ReqPerWorker
		jobs := make(chan int, total)

		for w := 0; w < step.Concurrency; w++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				executeWorker(id, jobs, client, targetURL, endpoints,
					&successful, &failed, &timeouts, &resets, &refused, &http5xx, &http4xx, &bodyMismatch, &other,
					&latencies, latencyMu, statusDist, statusMu, headers, successRe, failRe)
			}(w)
		}

		for j := 0; j < total; j++ {
			jobs <- j
		}
		close(jobs)
		wg.Wait()
	}

	close(done)
	fmt.Print("\r\033[K")

	elapsed := time.Since(start)
	totalCount := int(successful) + int(failed)

	var avgMs, p50, p95, p99 float64
	if len(latencies) > 0 {
		sum := 0.0
		for _, l := range latencies {
			sum += l
		}
		avgMs = sum / float64(len(latencies))
		sort.Float64s(latencies)
		p50 = metrics.PercentileSorted(latencies, 0.50)
		p95 = metrics.PercentileSorted(latencies, 0.95)
		p99 = metrics.PercentileSorted(latencies, 0.99)
	}

	successRate := 0.0
	if totalCount > 0 {
		successRate = float64(successful) / float64(totalCount) * 100
	}
	rps := 0.0
	if elapsed.Seconds() > 0 {
		rps = float64(totalCount) / elapsed.Seconds()
	}

	return config.PhaseResult{
		PhaseName:     step.Name,
		Concurrency:   step.Concurrency,
		DurationSec:   elapsed.Seconds(),
		TotalRequests: totalCount,
		Successful:    int(successful),
		Failed:        int(failed),
		SuccessRate:   successRate,
		AvgMs:         avgMs,
		P50Ms:         p50,
		P95Ms:         p95,
		P99Ms:         p99,
		ReqPerSec:     rps,
		Errors: config.ErrorCounts{
			Timeout:           int(timeouts),
			ConnectionReset:   int(resets),
			ConnectionRefused: int(refused),
			HTTP5xx:           int(http5xx),
			HTTP4xx:           int(http4xx),
			BodyMismatch:      int(bodyMismatch),
			Other:             int(other),
		},
		StatusDist: statusDist,
	}
}
