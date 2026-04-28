package runner

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/umutozen/stormprobe/internal/config"
	"github.com/umutozen/stormprobe/internal/metrics"
)

const bodyReadLimit = 4096

func executeWorker(workerID int, jobs <-chan int, client *http.Client, targetURL string, endpoints []string,
	successful, failed, timeouts, resets, refused, http5xx, http4xx, bodyMismatch, other *int32,
	latencies *[]float64, latencyMu *lockable, statusDist map[int]int, statusMu *lockable,
	customHeaders map[string]string, successRe, failRe *regexp.Regexp) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))
	time.Sleep(time.Duration(workerID*8) * time.Millisecond)

	for range jobs {
		gonderIstek(rng, client, targetURL, endpoints, successful, failed,
			timeouts, resets, refused, http5xx, http4xx, bodyMismatch, other,
			latencies, latencyMu, statusDist, statusMu, customHeaders, successRe, failRe)
	}
}

func executeWorkerDuration(ctx context.Context, workerID int, client *http.Client, targetURL string, endpoints []string,
	successful, failed, timeouts, resets, refused, http5xx, http4xx, bodyMismatch, other *int32,
	latencies *[]float64, latencyMu *lockable, statusDist map[int]int, statusMu *lockable,
	customHeaders map[string]string, successRe, failRe *regexp.Regexp) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))
	time.Sleep(time.Duration(workerID*8) * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			gonderIstek(rng, client, targetURL, endpoints, successful, failed,
				timeouts, resets, refused, http5xx, http4xx, bodyMismatch, other,
				latencies, latencyMu, statusDist, statusMu, customHeaders, successRe, failRe)
		}
	}
}

func gonderIstek(rng *rand.Rand, client *http.Client, targetURL string, endpoints []string,
	successful, failed, timeouts, resets, refused, http5xx, http4xx, bodyMismatch, other *int32,
	latencies *[]float64, latencyMu *lockable, statusDist map[int]int, statusMu *lockable,
	customHeaders map[string]string, successRe, failRe *regexp.Regexp) {

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
		return
	}

	req.Header.Set("User-Agent", config.UserAgents[rng.Intn(len(config.UserAgents))])
	req.Header.Set("Accept", "text/html,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	for k, v := range customHeaders {
		req.Header.Set(k, v)
	}

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
		return
	}

	needsBodyCheck := (successRe != nil || failRe != nil) && resp.StatusCode < 400
	var bodyBytes []byte

	if needsBodyCheck {
		bodyBytes, _ = io.ReadAll(io.LimitReader(resp.Body, bodyReadLimit))
		resp.Body.Close()
	} else {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

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
		if needsBodyCheck && bodyMatchFailed(bodyBytes, successRe, failRe) {
			atomic.AddInt32(failed, 1)
			atomic.AddInt32(bodyMismatch, 1)
		} else {
			atomic.AddInt32(successful, 1)
		}
	}
}

func bodyMatchFailed(body []byte, successRe, failRe *regexp.Regexp) bool {
	if failRe != nil && failRe.Match(body) {
		return true
	}
	if successRe != nil && !successRe.Match(body) {
		return true
	}
	return false
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
