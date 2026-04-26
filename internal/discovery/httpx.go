package discovery

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func filterWithHttpx(httpxPath string, rawURLs []string, targetURL string, headers map[string]string) []string {
	if len(rawURLs) == 0 {
		return nil
	}
	fmt.Printf("  [Httpx] Probing %d URLs\n", len(rawURLs))

	tmp, err := os.CreateTemp("", "httpx-input-*.txt")
	if err != nil {
		fmt.Printf("  [Httpx] Temp file error: %v\n", err)
		return nil
	}
	defer os.Remove(tmp.Name())

	for _, u := range rawURLs {
		fmt.Fprintln(tmp, u)
	}
	tmp.Close()

	args := []string{"-l", tmp.Name(), "-silent", "-mc", "200,201,301,302,303", "-timeout", "8", "-threads", "25", "-rate-limit", "40"}
	for k, v := range headers {
		args = append(args, "-H", fmt.Sprintf("%s: %s", k, v))
	}
	cmd := exec.Command(httpxPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case <-done:
	case <-time.After(60 * time.Second):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		<-done
		fmt.Println("  [Httpx] Timeout, using partial results")
	}

	base, _ := url.Parse(targetURL)
	targetHost := base.Hostname()
	seen := make(map[string]struct{})
	var paths []string

	for _, line := range strings.Split(stdout.String(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		parsed, err := url.Parse(parts[0])
		if err != nil || parsed.Hostname() != targetHost {
			continue
		}
		path := parsed.Path
		if parsed.RawQuery != "" {
			path += "?" + parsed.RawQuery
		}
		if _, exists := seen[path]; exists {
			continue
		}
		seen[path] = struct{}{}
		paths = append(paths, path)
		if len(paths) >= 50 {
			break
		}
	}
	fmt.Printf("  [Httpx] %d live endpoints confirmed\n", len(paths))
	return paths
}
