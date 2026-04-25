package metrics

import "strings"

func ClassifyError(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "Timeout") || strings.Contains(msg, "deadline"):
		return "timeout"
	case strings.Contains(msg, "reset") || strings.Contains(msg, "forcibly closed"):
		return "reset"
	case strings.Contains(msg, "refused") || strings.Contains(msg, "actively refused"):
		return "refused"
	default:
		return "other"
	}
}
