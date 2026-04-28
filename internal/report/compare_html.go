package report

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/umutozen/stormprobe/internal/compare"
	"github.com/umutozen/stormprobe/internal/config"
)

func WriteCompareHTML(result compare.CompareResult, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create output dir: %w", err)
	}

	filename := fmt.Sprintf("stormprobe_compare_%s.html", time.Now().Format("20060102_150405"))
	path := filepath.Join(outputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("cannot create compare HTML report: %w", err)
	}
	defer f.Close()

	b := result.Before.Verdict
	a := result.After.Verdict

	outcomeColor := outcomeToColor(result.Outcome)
	recovBefore := recoveryLabel(b.RecoveryOK)
	recovAfter := recoveryLabel(a.RecoveryOK)

	var whyRows string
	for _, w := range result.Why {
		icon := "~"
		rowColor := "#94a3b8"
		escaped := html.EscapeString(w)
		if strings.HasPrefix(w, "+") {
			icon = "▲"
			rowColor = "#22c55e"
		} else if strings.HasPrefix(w, "-") {
			icon = "▼"
			rowColor = "#ef4444"
		}
		whyRows += fmt.Sprintf(`<div style="padding:6px 12px;color:%s;font-size:14px;">%s %s</div>`, rowColor, icon, escaped)
	}

	var warningHTML string
	if len(result.Warnings) > 0 {
		warningHTML = `<div style="background:#fef3c7;border-left:4px solid #f59e0b;padding:12px 16px;margin-bottom:20px;border-radius:6px;">`
		for _, w := range result.Warnings {
			warningHTML += fmt.Sprintf(`<div style="color:#92400e;font-size:13px;">⚠ %s</div>`, html.EscapeString(w))
		}
		warningHTML += `</div>`
	}

	confBefore := b.Confidence
	if confBefore == "" {
		confBefore = "N/A"
	}
	confAfter := a.Confidence
	if confAfter == "" {
		confAfter = "N/A"
	}

	content := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>StormProbe Compare Report</title>
<style>
*{margin:0;padding:0;box-sizing:border-box;}
body{font-family:'Segoe UI',system-ui,-apple-system,sans-serif;background:#0f172a;color:#e2e8f0;padding:32px;min-height:100vh;}
.container{max-width:900px;margin:0 auto;}
.header{text-align:center;margin-bottom:32px;}
.header h1{font-size:28px;font-weight:700;color:#f8fafc;margin-bottom:4px;}
.header .sub{color:#64748b;font-size:13px;}
.outcome-card{text-align:center;padding:32px;border-radius:12px;margin-bottom:28px;background:linear-gradient(135deg,%s22,%s44);border:2px solid %s;}
.outcome-card h2{font-size:36px;font-weight:800;color:%s;margin-bottom:4px;}
.outcome-card .score{font-size:14px;color:#94a3b8;}
.compare-grid{display:grid;grid-template-columns:1fr 60px 1fr;gap:0;margin-bottom:28px;}
.side{background:#1e293b;border-radius:10px;padding:20px;}
.side h3{font-size:12px;text-transform:uppercase;letter-spacing:1px;color:#64748b;margin-bottom:14px;}
.metric{display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid #334155;}
.metric:last-child{border-bottom:none;}
.metric .label{color:#94a3b8;font-size:13px;}
.metric .value{color:#f8fafc;font-size:14px;font-weight:600;}
.arrow-col{display:flex;align-items:center;justify-content:center;font-size:24px;color:#64748b;}
.section{background:#1e293b;border-radius:10px;padding:20px;margin-bottom:20px;}
.section h3{font-size:14px;color:#64748b;text-transform:uppercase;letter-spacing:1px;margin-bottom:12px;}
.footer{text-align:center;color:#475569;font-size:12px;margin-top:32px;}
</style>
</head>
<body>
<div class="container">
<div class="header">
<h1>StormProbe Compare Report</h1>
<div class="sub">v%s · %s</div>
</div>
%s
<div class="outcome-card">
<h2>%s</h2>
<div class="score">Score: %d</div>
</div>

<div class="compare-grid">
<div class="side">
<h3>Before</h3>
<div class="metric"><span class="label">Target</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Timestamp</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Rating</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Confidence</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Safe VU</span><span class="value">%d</span></div>
<div class="metric"><span class="label">Recovery</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Bottleneck</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Requests</span><span class="value">%d</span></div>
<div class="metric"><span class="label">Failed</span><span class="value">%d</span></div>
</div>
<div class="arrow-col">→</div>
<div class="side">
<h3>After</h3>
<div class="metric"><span class="label">Target</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Timestamp</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Rating</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Confidence</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Safe VU</span><span class="value">%d</span></div>
<div class="metric"><span class="label">Recovery</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Bottleneck</span><span class="value">%s</span></div>
<div class="metric"><span class="label">Requests</span><span class="value">%d</span></div>
<div class="metric"><span class="label">Failed</span><span class="value">%d</span></div>
</div>
</div>

<div class="section">
<h3>Analysis</h3>
%s
</div>

<div class="footer">Generated by StormProbe v%s · github.com/umutozen/stormprobe</div>
</div>
</body>
</html>`,
		outcomeColor, outcomeColor, outcomeColor, outcomeColor,
		config.Version, time.Now().Format("2006-01-02 15:04"),
		warningHTML,
		result.Outcome, result.Score,
		html.EscapeString(result.Before.Target), result.Before.Timestamp,
		b.Rating, confBefore, b.SafeConcurrency, recovBefore,
		html.EscapeString(b.BottleneckCause),
		result.Before.TotalRequests, result.Before.TotalFailed,
		html.EscapeString(result.After.Target), result.After.Timestamp,
		a.Rating, confAfter, a.SafeConcurrency, recovAfter,
		html.EscapeString(a.BottleneckCause),
		result.After.TotalRequests, result.After.TotalFailed,
		whyRows,
		config.Version,
	)

	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}

	fmt.Printf("[+] Compare HTML report: %s\n", path)
	return path, nil
}

func outcomeToColor(outcome string) string {
	switch outcome {
	case "IMPROVED":
		return "#22c55e"
	case "PARTIAL IMPROVEMENT":
		return "#eab308"
	case "STABLE":
		return "#3b82f6"
	case "STABLE WITH RISK":
		return "#f59e0b"
	case "DEGRADED":
		return "#f97316"
	case "CRITICAL REGRESSION":
		return "#ef4444"
	default:
		return "#64748b"
	}
}

func recoveryLabel(ok bool) string {
	if ok {
		return "OK"
	}
	return "FAILED"
}
