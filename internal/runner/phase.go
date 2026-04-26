package runner

import (
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/umutozen/stormprobe/internal/config"
	"github.com/umutozen/stormprobe/internal/metrics"
)

func RunPhase(step config.PhaseStep, targetURL string, client *http.Client, endpoints []string, headers map[string]string) config.PhaseResult {
	total := step.Concurrency * step.ReqPerWorker
	var successful, failed int32
	var timeouts, resets, refused, http5xx, http4xx, other int32

	statusMu := newLock()
	statusDist := make(map[int]int)
	latencyMu := newLock()
	latencies := make([]float64, 0, total)

	jobs := make(chan int, total)
	var wg sync.WaitGroup
	start := time.Now()

	for w := 0; w < step.Concurrency; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			executeWorker(id, jobs, client, targetURL, endpoints,
				&successful, &failed, &timeouts, &resets, &refused, &http5xx, &http4xx, &other,
				&latencies, latencyMu, statusDist, statusMu, headers)
		}(w)
	}

	for j := 0; j < total; j++ {
		jobs <- j
	}
	close(jobs)
	wg.Wait()

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
			Other:             int(other),
		},
		StatusDist: statusDist,
	}
}
