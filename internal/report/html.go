package report

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/umutozen/stormprobe/internal/config"
)

const (
	labelCharLimit     = 18
	maxEndpointDisplay = 20
	colorSuccess       = "#10b981"
	colorWarning       = "#f59e0b"
	colorDanger        = "#ef4444"
)

func pickSuccessColor(rate float64) string {
	if rate < 80 {
		return colorDanger
	}
	if rate < 95 {
		return colorWarning
	}
	return colorSuccess
}

func WriteHTML(target, outputDir string, results []config.PhaseResult, endpoints []string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output dir: %w", err)
	}

	filename := fmt.Sprintf("stormprobe_report_%s.html", time.Now().Format("20060102_150405"))
	path := filepath.Join(outputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create HTML report: %w", err)
	}
	defer f.Close()

	var totalRequests, totalFailed int
	var totalDuration float64
	for _, r := range results {
		totalRequests += r.TotalRequests
		totalFailed += r.Failed
		totalDuration += r.DurationSec
	}
	overallSuccess := 0.0
	if totalRequests > 0 {
		overallSuccess = float64(totalRequests-totalFailed) / float64(totalRequests) * 100
	}

	var maxP99 float64
	for _, r := range results {
		if r.P99Ms > maxP99 {
			maxP99 = r.P99Ms
		}
	}
	if maxP99 == 0 {
		maxP99 = 100
	}

	chartHTML := buildChart(results, maxP99)
	tableHTML := buildTable(results)
	errorHTML := buildErrorSection(results, totalFailed)
	endpointHTML := buildEndpointList(endpoints)

	successColor := pickSuccessColor(overallSuccess)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>StormProbe Report | %s</title>
<style>
:root {
  --bg-primary: #0a0e1a;
  --bg-card: #111827;
  --bg-card-alt: #1a2332;
  --border: #1e293b;
  --border-accent: #334155;
  --text-primary: #f1f5f9;
  --text-secondary: #94a3b8;
  --text-muted: #64748b;
  --accent-blue: #38bdf8;
  --accent-purple: #818cf8;
  --accent-pink: #f472b6;
  --accent-green: #10b981;
  --accent-amber: #f59e0b;
  --accent-red: #ef4444;
  --accent-orange: #f97316;
  --accent-cyan: #22d3ee;
  --gradient-header: linear-gradient(135deg, #38bdf8 0%%, #818cf8 50%%, #c084fc 100%%);
  --gradient-card: linear-gradient(145deg, #111827 0%%, #0f172a 100%%);
  --shadow-lg: 0 8px 32px rgba(0,0,0,0.4);
}

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: 'Inter', 'SF Pro Display', system-ui, -apple-system, sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  line-height: 1.6;
  min-height: 100vh;
}

.container {
  max-width: 1100px;
  margin: 0 auto;
  padding: 48px 24px;
}

.header {
  text-align: center;
  margin-bottom: 48px;
  position: relative;
}

.header::before {
  content: '';
  position: absolute;
  top: -48px;
  left: 50%%;
  transform: translateX(-50%%);
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(56,189,248,0.08) 0%%, transparent 70%%);
  pointer-events: none;
}

.header .logo {
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 3px;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 12px;
}

.header h1 {
  font-size: 36px;
  font-weight: 800;
  background: var(--gradient-header);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.header .meta {
  color: var(--text-muted);
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

.header .meta span {
  display: flex;
  align-items: center;
  gap: 6px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
  margin-bottom: 32px;
}

@media (max-width: 768px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
}

.stat-card {
  background: var(--gradient-card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 24px 20px;
  text-align: center;
  position: relative;
  overflow: hidden;
  transition: transform 0.2s, border-color 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  border-color: var(--border-accent);
}

.stat-card .value {
  font-size: 36px;
  font-weight: 800;
  margin-bottom: 4px;
  letter-spacing: -1px;
}

.stat-card .label {
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
}

.card {
  background: var(--gradient-card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 28px;
  margin-bottom: 24px;
  box-shadow: var(--shadow-lg);
}

.card h2 {
  font-size: 13px;
  font-weight: 700;
  margin-bottom: 20px;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 1.5px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.card h2::before {
  content: '';
  width: 3px;
  height: 16px;
  background: var(--gradient-header);
  border-radius: 2px;
}

.chart-container {
  overflow-x: auto;
  padding: 8px 0;
  -webkit-overflow-scrolling: touch;
}

.chart-container svg {
  display: block;
  margin: 0 auto;
}

.legend {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-secondary);
}

.legend-item .dot {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.legend-item .line {
  width: 20px;
  height: 2px;
  background: var(--accent-pink);
}

.legend-item .line-dashed {
  width: 20px;
  height: 0;
  border-top: 2px dashed var(--accent-cyan);
}

table {
  width: 100%%;
  border-collapse: separate;
  border-spacing: 0;
  font-size: 13px;
}

thead th {
  padding: 12px 14px;
  text-align: left;
  color: var(--text-muted);
  font-weight: 700;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 2px solid var(--border-accent);
  white-space: nowrap;
}

tbody td {
  padding: 12px 14px;
  border-bottom: 1px solid var(--border);
  font-variant-numeric: tabular-nums;
}

tbody tr {
  transition: background 0.15s;
}

tbody tr:hover td {
  background: rgba(56,189,248,0.04);
}

tbody tr:last-child td {
  border-bottom: none;
}

.phase-name {
  font-weight: 600;
  color: var(--text-primary);
}

.success-badge {
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
}

.error-bar-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.error-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.error-row .indicator {
  width: 14px;
  height: 14px;
  border-radius: 4px;
  flex-shrink: 0;
}

.error-row .name {
  min-width: 130px;
  font-size: 13px;
  font-weight: 500;
}

.error-row .bar-bg {
  flex: 1;
  background: var(--bg-primary);
  border-radius: 6px;
  height: 24px;
  overflow: hidden;
}

.error-row .bar-fill {
  height: 100%%;
  border-radius: 6px;
  opacity: 0.8;
  transition: width 0.3s;
}

.error-row .count {
  min-width: 90px;
  text-align: right;
  font-size: 13px;
  font-variant-numeric: tabular-nums;
  color: var(--text-secondary);
}

.no-errors {
  color: var(--accent-green);
  font-weight: 700;
  font-size: 15px;
}

.endpoints-list {
  list-style: none;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 6px;
}

.endpoints-list li {
  padding: 6px 0;
  font-size: 13px;
}

.endpoints-list li code {
  background: var(--bg-primary);
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 12px;
  color: var(--accent-cyan);
  border: 1px solid var(--border);
}

.footer {
  text-align: center;
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
  font-size: 12px;
  color: var(--text-muted);
  display: flex;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
}
</style>
</head>
<body>
<div class="container">

<div class="header">
  <div class="logo">StormProbe</div>
  <h1>Load Test Report</h1>
  <div class="meta">
    <span>%s</span>
    <span>%s</span>
    <span>%.1fs total</span>
  </div>
</div>

<div class="stats-grid">
  <div class="stat-card">
    <div class="value" style="color:%s">%.1f%%</div>
    <div class="label">Success Rate</div>
  </div>
  <div class="stat-card">
    <div class="value" style="color:var(--accent-blue)">%d</div>
    <div class="label">Total Requests</div>
  </div>
  <div class="stat-card">
    <div class="value" style="color:var(--accent-pink)">%d</div>
    <div class="label">Failed</div>
  </div>
  <div class="stat-card">
    <div class="value" style="color:var(--accent-purple)">%d</div>
    <div class="label">Phases</div>
  </div>
  <div class="stat-card">
    <div class="value" style="color:var(--accent-cyan)">%d</div>
    <div class="label">Endpoints</div>
  </div>
</div>

<div class="card">
  <h2>Latency Overview</h2>
  <div class="chart-container">%s</div>
  <div class="legend">
    <div class="legend-item"><div class="dot" style="background:var(--accent-green)"></div>Avg Latency (ms)</div>
    <div class="legend-item"><div class="line"></div>P95 Latency</div>
    <div class="legend-item"><div class="line-dashed"></div>P99 Latency</div>
  </div>
</div>

<div class="card">
  <h2>Phase Results</h2>
  <div style="overflow-x:auto">%s</div>
</div>

<div class="card">
  <h2>Error Distribution</h2>
  %s
</div>

<div class="card">
  <h2>Endpoints (%d)</h2>
  %s
</div>

<div class="footer">
  <span>StormProbe v%s</span>
  <span>%s</span>
  <span>github.com/umutozen/stormprobe</span>
</div>

</div>
</body>
</html>`,
		target,
		target, timestamp, totalDuration,
		successColor, overallSuccess,
		totalRequests, totalFailed, len(results), len(endpoints),
		chartHTML, tableHTML, errorHTML,
		len(endpoints), endpointHTML,
		config.Version, timestamp)

	_, err = f.WriteString(html)
	if err != nil {
		return err
	}
	fmt.Printf("[+] HTML report: %s\n", path)
	return nil
}

func buildChart(results []config.PhaseResult, maxP99 float64) string {
	barWidth := 48
	chartHeight := 260
	gap := 18
	leftPad := 65
	bottomPad := 72

	totalWidth := leftPad + len(results)*(barWidth+gap) + 40
	if totalWidth < 700 {
		totalWidth = 700
	}
	totalHeight := chartHeight + bottomPad

	svg := fmt.Sprintf(`<svg viewBox="0 0 %d %d" width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, totalWidth, totalHeight, totalWidth, totalHeight)

	svg += fmt.Sprintf(`<defs>
<linearGradient id="barGrad" x1="0" y1="0" x2="0" y2="1">
  <stop offset="0%%" stop-color="#10b981" stop-opacity="0.9"/>
  <stop offset="100%%" stop-color="#059669" stop-opacity="0.7"/>
</linearGradient>
<linearGradient id="barWarn" x1="0" y1="0" x2="0" y2="1">
  <stop offset="0%%" stop-color="#f59e0b" stop-opacity="0.9"/>
  <stop offset="100%%" stop-color="#d97706" stop-opacity="0.7"/>
</linearGradient>
<linearGradient id="barDanger" x1="0" y1="0" x2="0" y2="1">
  <stop offset="0%%" stop-color="#ef4444" stop-opacity="0.9"/>
  <stop offset="100%%" stop-color="#dc2626" stop-opacity="0.7"/>
</linearGradient>
<filter id="glow"><feGaussianBlur stdDeviation="2" result="blur"/><feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge></filter>
</defs>`)

	for i := 0; i <= 4; i++ {
		y := int(float64(i) / 4 * float64(chartHeight-20))
		val := maxP99 * (1 - float64(i)/4)
		svg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#1e293b" stroke-dasharray="4"/>`, leftPad-5, y+10, totalWidth-20, y+10)
		svg += fmt.Sprintf(`<text x="%d" y="%d" text-anchor="end" fill="#64748b" font-size="10" font-family="system-ui">%.0f</text>`, leftPad-10, y+14, val)
	}

	svg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#334155" stroke-width="1"/>`, leftPad-5, chartHeight-10, totalWidth-20, chartHeight-10)

	var p95Points, p99Points string

	for i, r := range results {
		x := leftPad + i*(barWidth+gap)
		drawArea := float64(chartHeight - 30)

		barHeight := int(r.AvgMs / maxP99 * drawArea)
		if barHeight < 4 {
			barHeight = 4
		}
		y := chartHeight - 10 - barHeight

		gradID := "barGrad"
		if r.SuccessRate < 95 {
			gradID = "barWarn"
		}
		if r.SuccessRate < 80 {
			gradID = "barDanger"
		}

		svg += fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" rx="6" fill="url(#%s)"><title>%s | Avg:%.0fms P50:%.0fms P95:%.0fms P99:%.0fms | %.1f%% success | %.1f req/s</title></rect>`,
			x, y, barWidth, barHeight, gradID, r.PhaseName, r.AvgMs, r.P50Ms, r.P95Ms, r.P99Ms, r.SuccessRate, r.ReqPerSec)

		svg += fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" fill="#e2e8f0" font-size="11" font-weight="700" font-family="system-ui">%.0f</text>`,
			x+barWidth/2, y-6, r.AvgMs)

		lbl := r.PhaseName
		if len(lbl) > labelCharLimit {
			lbl = lbl[:labelCharLimit]
		}
		svg += fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" fill="#64748b" font-size="9" font-family="system-ui" transform="rotate(-35 %d %d)">%s</text>`,
			x+barWidth/2, chartHeight+8, x+barWidth/2, chartHeight+8, lbl)

		mid := x + barWidth/2
		p95y := chartHeight - 10 - int(r.P95Ms/maxP99*drawArea)
		p99y := chartHeight - 10 - int(r.P99Ms/maxP99*drawArea)

		if i == 0 {
			p95Points = fmt.Sprintf("%d,%d", mid, p95y)
			p99Points = fmt.Sprintf("%d,%d", mid, p99y)
		} else {
			p95Points += fmt.Sprintf(" %d,%d", mid, p95y)
			p99Points += fmt.Sprintf(" %d,%d", mid, p99y)
		}
	}

	svg += fmt.Sprintf(`<polyline points="%s" fill="none" stroke="#f472b6" stroke-width="2" stroke-linejoin="round" filter="url(#glow)"/>`, p95Points)
	svg += fmt.Sprintf(`<polyline points="%s" fill="none" stroke="#22d3ee" stroke-width="1.5" stroke-dasharray="6,4" stroke-linejoin="round"/>`, p99Points)

	svg += `</svg>`
	return svg
}

func buildTable(results []config.PhaseResult) string {
	html := `<table>
<thead><tr>
<th>Phase</th><th>VU</th><th>Success</th><th>Avg ms</th><th>P50 ms</th><th>P95 ms</th><th>P99 ms</th><th>req/s</th><th>OK / Total</th><th>Duration</th>
</tr></thead><tbody>`

	for _, r := range results {
		color := pickSuccessColor(r.SuccessRate)
		bg := "rgba(245,158,11,0.08)"
		if r.SuccessRate >= 95 {
			bg = "rgba(16,185,129,0.08)"
		} else if r.SuccessRate < 80 {
			bg = "rgba(239,68,68,0.08)"
		}

		html += fmt.Sprintf(`<tr><td class="phase-name">%s</td><td>%d</td><td><span class="success-badge" style="color:%s;background:%s">%.1f%%</span></td><td>%.0f</td><td>%.0f</td><td>%.0f</td><td>%.0f</td><td>%.1f</td><td>%d / %d</td><td>%.1fs</td></tr>`,
			r.PhaseName, r.Concurrency, color, bg, r.SuccessRate,
			r.AvgMs, r.P50Ms, r.P95Ms, r.P99Ms, r.ReqPerSec,
			r.Successful, r.TotalRequests, r.DurationSec)
	}

	html += `</tbody></table>`
	return html
}

func buildErrorSection(results []config.PhaseResult, totalFailed int) string {
	var tTo, tRs, tRf, t5, t4, tOt int
	for _, r := range results {
		tTo += r.Errors.Timeout
		tRs += r.Errors.ConnectionReset
		tRf += r.Errors.ConnectionRefused
		t5 += r.Errors.HTTP5xx
		t4 += r.Errors.HTTP4xx
		tOt += r.Errors.Other
	}

	if totalFailed == 0 {
		return `<div class="no-errors">No errors recorded</div>`
	}

	html := `<div class="error-bar-container">`
	errorTypes := []struct {
		name  string
		count int
		color string
	}{
		{"Timeout", tTo, colorDanger},
		{"Connection Reset", tRs, "#f97316"},
		{"Connection Refused", tRf, colorWarning},
		{"HTTP 5xx", t5, "#ec4899"},
		{"HTTP 4xx", t4, "#8b5cf6"},
		{"Other", tOt, "#64748b"},
	}

	for _, e := range errorTypes {
		if e.count == 0 {
			continue
		}
		pct := float64(e.count) / float64(totalFailed) * 100
		html += fmt.Sprintf(`<div class="error-row">
<div class="indicator" style="background:%s"></div>
<span class="name">%s</span>
<div class="bar-bg"><div class="bar-fill" style="width:%.1f%%;background:%s"></div></div>
<span class="count">%d (%.1f%%)</span>
</div>`, e.color, e.name, pct, e.color, e.count, pct)
	}

	html += `</div>`
	return html
}

func buildEndpointList(endpoints []string) string {
	html := `<ul class="endpoints-list">`
	for i, ep := range endpoints {
		if i >= maxEndpointDisplay {
			html += fmt.Sprintf(`<li style="color:var(--text-muted)">... and %d more</li>`, len(endpoints)-maxEndpointDisplay)
			break
		}
		html += fmt.Sprintf(`<li><code>%s</code></li>`, ep)
	}
	html += `</ul>`
	return html
}
