package discovery

import (
	"sort"
	"strings"
)

const maxScoredEndpoints = 20

var staticExtensions = []string{
	".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico",
	".woff", ".woff2", ".ttf", ".eot", ".otf", ".map",
	".pdf", ".zip", ".tar", ".gz", ".xml", ".txt",
}

var staticPrefixes = []string{
	"/static/", "/assets/", "/vendor/", "/fonts/",
	"/images/", "/img/", "/media/", "/Content/", "/content/",
}

var negativePatterns = []string{
	"/logout", "/signout", "/sign-out", "/favicon",
	"utm_", "ga_", "__utm", "analytics",
}

type ScoredResult struct {
	Path   string
	Score  int
	Reason string
}

var scoringRules = []struct {
	pattern string
	points  int
	label   string
}{
	{"/api/", 10, "API endpoint"},
	{"/payment", 10, "financial endpoint"},
	{"/pay/", 10, "financial endpoint"},
	{"/checkout", 10, "financial endpoint"},
	{"/borc", 9, "financial endpoint"},
	{"/tahakkuk", 9, "financial endpoint"},
	{"/odeme", 9, "financial endpoint"},
	{"/admin", 8, "admin panel"},
	{"/login", 8, "authenticated path"},
	{"/giris", 8, "authenticated path"},
	{"/signin", 8, "authenticated path"},
	{"/auth", 8, "authenticated path"},
	{"/account", 7, "user account"},
	{"/dashboard", 7, "user dashboard"},
	{"/search", 7, "search endpoint"},
	{"/sorgula", 7, "query endpoint"},
	{"/query", 7, "query endpoint"},
	{"/ruhsat", 7, "document endpoint"},
	{"/user", 6, "user endpoint"},
	{"/profile", 6, "user endpoint"},
	{"/form", 6, "form submission"},
	{"/submit", 6, "form submission"},
	{"/register", 6, "registration flow"},
	{"/upload", 5, "file upload"},
}

func isStaticAsset(path string) bool {
	lower := strings.ToLower(path)
	for _, ext := range staticExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	for _, prefix := range staticPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

func isNegative(path string) bool {
	lower := strings.ToLower(path)
	for _, neg := range negativePatterns {
		if strings.Contains(lower, neg) {
			return true
		}
	}
	return false
}

func scoreAndReason(path string) (int, string) {
	lower := strings.ToLower(path)
	score := 2
	seen := make(map[string]struct{})
	var labels []string

	for _, rule := range scoringRules {
		if strings.Contains(lower, rule.pattern) {
			score += rule.points
			if _, exists := seen[rule.label]; !exists {
				seen[rule.label] = struct{}{}
				labels = append(labels, rule.label)
			}
		}
	}

	if len(labels) == 0 {
		return score, "general endpoint"
	}
	return score, strings.Join(labels, " + ")
}

func ScoreAndFilter(paths []string) []ScoredResult {
	var candidates []ScoredResult
	for _, p := range paths {
		if strings.TrimSpace(p) == "" || isStaticAsset(p) || isNegative(p) {
			continue
		}
		score, reason := scoreAndReason(p)
		candidates = append(candidates, ScoredResult{Path: p, Score: score, Reason: reason})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	limit := maxScoredEndpoints
	if len(candidates) < limit {
		limit = len(candidates)
	}
	return candidates[:limit]
}
