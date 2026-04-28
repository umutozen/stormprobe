package discovery

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func headlessBrowserAvailable() bool {
	for _, name := range []string{"chromium-browser", "chromium", "google-chrome", "chrome"} {
		if _, err := exec.LookPath(name); err == nil {
			return true
		}
	}
	if runtime.GOOS == "windows" {
		for _, name := range []string{"chrome.exe", "chromium.exe", "msedge.exe"} {
			if _, err := exec.LookPath(name); err == nil {
				return true
			}
		}
		commonPaths := []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files\Chromium\Application\chrome.exe`,
		}
		for _, p := range commonPaths {
			if _, err := os.Stat(p); err == nil {
				return true
			}
		}
	}
	return false
}

func crawlWithKatana(katanaPath, targetURL string, headers map[string]string) ([]string, error) {
	fmt.Printf("  [Katana] Crawling %s\n", targetURL)

	args := []string{"-u", targetURL, "-d", "3", "-fs", "fqdn", "-silent", "-timeout", "10", "-rate-limit", "30"}
	for k, v := range headers {
		args = append(args, "-H", fmt.Sprintf("%s: %s", k, v))
	}
	if headlessBrowserAvailable() {
		args = append(args, "-jc")
		fmt.Println("  [Katana] Headless browser detected, JS crawl enabled")
	} else {
		fmt.Println("  [Katana] No headless browser, static crawl only")
	}

	cmd := exec.Command(katanaPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		if err != nil && stdout.Len() == 0 {
			return nil, fmt.Errorf("katana error: %v | stderr: %s", err, stderr.String())
		}
	case <-time.After(120 * time.Second):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		<-done
		fmt.Println("  [Katana] Timeout, using partial results")
	}

	var urls []string
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			urls = append(urls, line)
		}
	}
	fmt.Printf("  [Katana] %d URLs found\n", len(urls))
	return urls, nil
}
