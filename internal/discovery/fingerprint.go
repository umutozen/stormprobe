package discovery

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type StackType string

const (
	StackIIS     StackType = "IIS"
	StackTomcat  StackType = "Tomcat"
	StackPHP     StackType = "PHP"
	StackNginx   StackType = "Nginx"
	StackApache  StackType = "Apache"
	StackUnknown StackType = "Unknown"
)

type ConfidenceLevel string

const (
	ConfidenceHigh    ConfidenceLevel = "High"
	ConfidenceMedium  ConfidenceLevel = "Medium"
	ConfidenceLow     ConfidenceLevel = "Low"
	ConfidenceUnknown ConfidenceLevel = "Unknown"
)

type StackFingerprint struct {
	EdgeLayer         StackType
	EdgeConfidence    ConfidenceLevel
	AppLayer          StackType
	AppConfidence     ConfidenceLevel
	RawServerHeader   string
	RawPoweredBy      string
	SessionCookie     string
	BehindProxy       bool
	DiagnosticProfile string
}

func (f StackFingerprint) Summary() string {
	if f.BehindProxy {
		return fmt.Sprintf(
			"Edge: %s (%s) | App: %s (%s)",
			f.EdgeLayer, f.EdgeConfidence,
			f.AppLayer, f.AppConfidence,
		)
	}
	return fmt.Sprintf("%s (%s)", f.AppLayer, f.AppConfidence)
}

func (f StackFingerprint) EffectiveApp() (StackType, ConfidenceLevel) {
	if f.AppLayer != StackUnknown {
		return f.AppLayer, f.AppConfidence
	}
	return f.EdgeLayer, f.EdgeConfidence
}

func Fingerprint(targetURL string, insecure bool) StackFingerprint {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("HEAD", targetURL, nil)
	if err != nil {
		return StackFingerprint{AppLayer: StackUnknown, AppConfidence: ConfidenceUnknown}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; StormProbe/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		req2, _ := http.NewRequest("GET", targetURL, nil)
		req2.Header.Set("User-Agent", "Mozilla/5.0 (compatible; StormProbe/1.0)")
		resp, err = client.Do(req2)
		if err != nil {
			return StackFingerprint{AppLayer: StackUnknown, AppConfidence: ConfidenceUnknown}
		}
	}
	defer resp.Body.Close()

	server := resp.Header.Get("Server")
	poweredBy := resp.Header.Get("X-Powered-By")
	aspxVersion := resp.Header.Get("X-AspNet-Version")
	generator := resp.Header.Get("X-Generator")

	var sessionCookie string
	for _, c := range resp.Cookies() {
		sessionCookie += c.Name + ";"
	}

	fp := StackFingerprint{
		RawServerHeader: server,
		RawPoweredBy:    poweredBy,
		SessionCookie:   strings.TrimRight(sessionCookie, ";"),
	}

	fp.EdgeLayer, fp.EdgeConfidence = detectEdge(server)
	fp.AppLayer, fp.AppConfidence = detectApp(server, poweredBy, aspxVersion, generator, sessionCookie)
	fp.BehindProxy = fp.EdgeLayer != StackUnknown && fp.AppLayer != StackUnknown && fp.EdgeLayer != fp.AppLayer

	appStack, appConf := fp.EffectiveApp()
	fp.DiagnosticProfile = selectProfile(appStack, appConf)

	return fp
}

func detectEdge(server string) (StackType, ConfidenceLevel) {
	s := strings.ToLower(server)
	switch {
	case strings.Contains(s, "nginx"):
		return StackNginx, ConfidenceHigh
	case strings.Contains(s, "apache"):
		return StackApache, ConfidenceHigh
	case strings.Contains(s, "cloudflare"):
		return StackNginx, ConfidenceMedium
	case strings.Contains(s, "microsoft-iis"):
		return StackIIS, ConfidenceHigh
	}
	return StackUnknown, ConfidenceUnknown
}

func detectApp(server, poweredBy, aspxVersion, generator, cookies string) (StackType, ConfidenceLevel) {
	sl := strings.ToLower(server)
	pl := strings.ToLower(poweredBy)
	cl := strings.ToLower(cookies)
	gl := strings.ToLower(generator)

	switch {
	case strings.Contains(pl, "asp.net") || aspxVersion != "" || strings.Contains(cl, "asp.net_sessionid"):
		return StackIIS, ConfidenceHigh
	case strings.Contains(cl, "jsessionid"):
		return StackTomcat, ConfidenceHigh
	case strings.Contains(pl, "php") || strings.Contains(gl, "php") || strings.Contains(sl, "php"):
		return StackPHP, ConfidenceHigh
	case strings.Contains(pl, "express") || strings.Contains(pl, "node"):
		return StackUnknown, ConfidenceLow
	case strings.Contains(sl, "tomcat") || strings.Contains(sl, "jetty"):
		return StackTomcat, ConfidenceMedium
	case strings.Contains(pl, "asp.net") || strings.Contains(sl, "iis"):
		return StackIIS, ConfidenceMedium
	}
	return StackUnknown, ConfidenceLow
}

func selectProfile(stack StackType, conf ConfidenceLevel) string {
	if conf == ConfidenceLow || conf == ConfidenceUnknown {
		return "default"
	}
	switch stack {
	case StackIIS:
		return "iis"
	case StackTomcat:
		return "tomcat"
	case StackPHP:
		return "php"
	default:
		return "default"
	}
}

func PrintFingerprint(fp StackFingerprint) {
	fmt.Println()
	fmt.Println("+==================================================+")
	fmt.Println("|  PHASE 0: STACK FINGERPRINT                      |")
	fmt.Println("+==================================================+")

	if fp.BehindProxy {
		fmt.Printf("  Edge Layer       : %-20s [%s]\n", fp.EdgeLayer, fp.EdgeConfidence)
		fmt.Printf("  Application Layer: %-20s [%s]\n", fp.AppLayer, fp.AppConfidence)
	} else {
		app, conf := fp.EffectiveApp()
		fmt.Printf("  Detected Stack   : %-20s [%s]\n", app, conf)
	}

	if fp.RawServerHeader != "" {
		fmt.Printf("  Server Header    : %s\n", fp.RawServerHeader)
	}
	if fp.RawPoweredBy != "" {
		fmt.Printf("  X-Powered-By     : %s\n", fp.RawPoweredBy)
	}
	if fp.SessionCookie != "" {
		fmt.Printf("  Session Cookies  : %s\n", fp.SessionCookie)
	}

	fmt.Printf("  Diagnostic Profile: %s\n", fp.DiagnosticProfile)

	app, conf := fp.EffectiveApp()
	if conf == ConfidenceLow || conf == ConfidenceUnknown || app == StackUnknown {
		fmt.Println("  ⚠  Low confidence — using default profile to avoid wrong diagnosis")
	}
	fmt.Println()
}
