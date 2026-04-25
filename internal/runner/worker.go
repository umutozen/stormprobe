package runner

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/umutozen/stormprobe/internal/config"
	"github.com/umutozen/stormprobe/internal/metrics"
)

func executeWorker(workerID int, jobs <-chan int, client *http.Client, targetURL string, endpoints []string,
	successful, failed, timeouts, resets, refused, http5xx, http4xx, other *int32,
	latencies *[]float64, latencyMu *lockable, statusDist map[int]int, statusMu *lockable) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))
	time.Sleep(time.Duration(workerID*8) * time.Millisecond)

	for range jobs {
		reqStart := time.Now()

		endpoint := endpoints[rng.Intn(len(endpoints))]
		var fullURL string
		if strings.HasPrefix(endpoint, "http") {
			fullURL = fmt.Sprintf("%s?cb=%d", endpoint, rng.Intn(999999))
		} else {
			fullURL = fmt.Sprintf("%s%s?cb=%d", strings.TrimRight(targetURL, "/"), endpoint, rng.Intn(999999))
		}

		req, err := http.NewRequest("GET", fullURL, nil)
		elapsedMs := float64(time.Since(reqStart).Milliseconds())

		if err != nil {
			latencyMu.Lock()
			*latencies = append(*latencies, elapsedMs)
			latencyMu.Unlock()
			atomic.AddInt32(failed, 1)
			atomic.AddInt32(other, 1)
			continue
		}

		req.Header.Set("User-Agent", config.UserAgents[rng.Intn(len(config.UserAgents))])
		req.Header.Set("Accept", "text/html,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")

		resp, doErr := client.Do(req)
		elapsedMs = float64(time.Since(reqStart).Milliseconds())

		latencyMu.Lock()
		*latencies = append(*latencies, elapsedMs)
		latencyMu.Unlock()

		if doErr != nil {
			atomic.AddInt32(failed, 1)
			switch metrics.ClassifyError(doErr) {
			case "timeout":
				atomic.AddInt32(timeouts, 1)
			case "reset":
				atomic.AddInt32(resets, 1)
			case "refused":
				atomic.AddInt32(refused, 1)
			default:
				atomic.AddInt32(other, 1)
			}
			continue
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		statusMu.Lock()
		statusDist[resp.StatusCode]++
		statusMu.Unlock()

		switch {
		case resp.StatusCode >= 500:
			atomic.AddInt32(failed, 1)
			atomic.AddInt32(http5xx, 1)
		case resp.StatusCode >= 400:
			atomic.AddInt32(failed, 1)
			atomic.AddInt32(http4xx, 1)
		default:
			atomic.AddInt32(successful, 1)
		}
	}
}

type lockable struct {
	ch chan struct{}
}

func newLock() *lockable {
	l := &lockable{ch: make(chan struct{}, 1)}
	l.ch <- struct{}{}
	return l
}

func (l *lockable) Lock()   { <-l.ch }
func (l *lockable) Unlock() { l.ch <- struct{}{} }
