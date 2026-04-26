package verdict

import (
	"fmt"

	"github.com/umutozen/stormprobe/internal/config"
)

// Eşik sabitleri — iş mantığı burada toplanır.
const (
	p99DegradedMs  = 2000.0 // P99 bu değeri aşarsa degraded sayılır
	p99CriticalMs  = 5000.0 // P99 bu değeri aşarsa critical sayılır
	errorRateSafe  = 1.0    // % hata oranı güvenli sınır
	errorRateFail  = 2.0    // % hata oranı kritik sınır
	recoveryRatio  = 1.5    // recovery P99 / baseline P99 oranı — bu aşılırsa toparlanma tam değil
)

// Rating sunucunun genel değerlendirmesidir.
type Rating string

const (
	RatingHealthy  Rating = "Healthy"
	RatingDegraded Rating = "Degraded"
	RatingCritical Rating = "Critical"
	RatingUnstable Rating = "Unstable" // spike sonrası toparlanamadı
)

// Verdict testin bütünsel yorumudur.
type Verdict struct {
	Rating           Rating
	SafeConcurrency  int
	DegradationAt    int    // 0 = görülmedi
	FailurePoint     int    // 0 = görülmedi
	BottleneckCause  string
	BottleneckDetail string
	RecoveryOK       bool
	Recommendation   string
}

// Analyse tüm faz sonuçlarını değerlendirerek bir Verdict üretir.
func Analyse(results []config.PhaseResult) Verdict {
	if len(results) == 0 {
		return Verdict{Rating: RatingHealthy, Recommendation: "No data to analyse."}
	}

	var (
		safeConcurrency int
		degradationAt   int
		failurePoint    int
		baselineP99     float64
		recoveryP99     float64
		hasRecovery     bool
	)

	// Baseline: ilk fazın P99'u
	baselineP99 = results[0].P99Ms

	for _, r := range results {
		errRate := errorRate(r)
		isRecovery := r.PhaseName == "Recovery check"

		if isRecovery {
			hasRecovery = true
			recoveryP99 = r.P99Ms
			continue
		}

		// Güvenli bölge: hata < %1 ve P99 < 2s
		if errRate < errorRateSafe && r.P99Ms < p99DegradedMs {
			if r.Concurrency > safeConcurrency {
				safeConcurrency = r.Concurrency
			}
		}

		// Bozulma başlangıcı
		if degradationAt == 0 && r.P99Ms >= p99DegradedMs {
			degradationAt = r.Concurrency
		}

		// Kritik eşik
		if failurePoint == 0 && (r.P99Ms >= p99CriticalMs || errRate >= errorRateFail) {
			failurePoint = r.Concurrency
		}
	}

	recoveryOK := !hasRecovery || (baselineP99 > 0 && recoveryP99 <= baselineP99*recoveryRatio)

	cause, detail := diagnoseBottleneck(results)
	rating := computeRating(failurePoint, degradationAt, recoveryOK)
	recommendation := buildRecommendation(rating, cause, safeConcurrency)

	return Verdict{
		Rating:           rating,
		SafeConcurrency:  safeConcurrency,
		DegradationAt:    degradationAt,
		FailurePoint:     failurePoint,
		BottleneckCause:  cause,
		BottleneckDetail: detail,
		RecoveryOK:       recoveryOK,
		Recommendation:   recommendation,
	}
}

func errorRate(r config.PhaseResult) float64 {
	if r.TotalRequests == 0 {
		return 0
	}
	return float64(r.Failed) / float64(r.TotalRequests) * 100
}

// diagnoseBottleneck hata tipine göre olası nedeni tahmin eder.
func diagnoseBottleneck(results []config.PhaseResult) (cause, detail string) {
	var timeouts, resets, refused, http5xx, http4xx int
	for _, r := range results {
		timeouts += r.Errors.Timeout
		resets += r.Errors.ConnectionReset
		refused += r.Errors.ConnectionRefused
		http5xx += r.Errors.HTTP5xx
		http4xx += r.Errors.HTTP4xx
	}

	totalErrors := timeouts + resets + refused + http5xx + http4xx
	if totalErrors == 0 {
		return "None", "No errors detected across all phases."
	}

	switch {
	case isdominant(timeouts, totalErrors):
		return "Application/DB processing bottleneck",
			"Requests queue up server-side under load. Likely cause: slow DB queries, blocking I/O, or insufficient worker threads."
	case isDominant2(resets, refused, totalErrors):
		return "Connection pool / infrastructure capacity exhausted",
			"Server is dropping connections before processing. Check reverse proxy limits, OS socket backlog, and connection pool size."
	case isdominant(http5xx, totalErrors):
		return "Server-side application errors under load",
			"HTTP 5xx responses indicate application crashes or unhandled exceptions at high concurrency. Review application logs."
	case isDominant2(resets, timeouts, totalErrors):
		return "Mixed failure mode — likely cascading overload",
			"Multiple error types suggest resource exhaustion across layers (network, application, and DB)."
	default:
		return "Mixed / undetermined",
			fmt.Sprintf("Timeouts: %d | Resets: %d | Refused: %d | 5xx: %d | 4xx: %d", timeouts, resets, refused, http5xx, http4xx)
	}
}

// isdominant bir hata tipinin toplam hataların %60'ından fazlasını oluşturup oluşturmadığını kontrol eder.
func isdominant(count, total int) bool {
	return total > 0 && float64(count)/float64(total) >= 0.60
}

func isDominant2(a, b, total int) bool {
	return total > 0 && float64(a+b)/float64(total) >= 0.60
}

func computeRating(failurePoint, degradationAt int, recoveryOK bool) Rating {
	if failurePoint > 0 && !recoveryOK {
		return RatingUnstable
	}
	if failurePoint > 0 {
		return RatingCritical
	}
	if degradationAt > 0 {
		return RatingDegraded
	}
	return RatingHealthy
}

func buildRecommendation(r Rating, cause string, safeConcurrency int) string {
	switch r {
	case RatingHealthy:
		return fmt.Sprintf("System is healthy up to %d virtual users. Monitor P99 under real traffic.", safeConcurrency)
	case RatingDegraded:
		return fmt.Sprintf("Cap production traffic at ~%d concurrent users. Investigate: %s", safeConcurrency, cause)
	case RatingCritical:
		return fmt.Sprintf("Do not exceed %d concurrent users without fixes. Priority: %s", safeConcurrency, cause)
	case RatingUnstable:
		return "Server did not fully recover after the spike. High risk of cascading failure under burst traffic."
	default:
		return "Review results manually."
	}
}
