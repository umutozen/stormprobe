package discovery

import (
	"bufio"
	"bytes"
	"fmt"
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
	}
	return false
}

func crawlWithKatana(katanaPath, targetURL string) ([]string, error) {
	fmt.Printf("  [Katana] Crawling %s\n", targetURL)

	args := []string{"-u", targetURL, "-d", "3", "-fs", "fqdn", "-silent", "-timeout", "10", "-rate-limit", "30"}
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
			cmd.Process.Kill()
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
