package discovery

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func findTool(name, customPath string) string {
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			return customPath
		}
	}
	if path, err := exec.LookPath(name); err == nil {
		return path
	}
	if runtime.GOOS == "windows" {
		if path, err := exec.LookPath(name + ".exe"); err == nil {
			return path
		}
	}
	cwd, _ := os.Getwd()
	for _, c := range []string{
		filepath.Join(cwd, name),
		filepath.Join(cwd, name+".exe"),
	} {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

func Discover(targetURL, katanaCustom, httpxCustom string, noDiscovery bool, headers map[string]string) []string {
	fmt.Println("\n+==================================================+")
	fmt.Println("|  PHASE 0: DISCOVERY (Katana + Httpx)             |")
	fmt.Println("+==================================================+")

	if noDiscovery {
		fmt.Println("  [!] Discovery disabled, using root path only")
		return []string{"/"}
	}

	katanaPath := findTool("katana", katanaCustom)
	httpxPath := findTool("httpx", httpxCustom)

	if katanaPath == "" {
		fmt.Println("  [!] Katana not found, using root path only")
		return []string{"/"}
	}

	allURLs, err := crawlWithKatana(katanaPath, targetURL, headers)
	if err != nil || len(allURLs) == 0 {
		fmt.Println("  [!] Katana yielded no results, using root path only")
		return []string{"/"}
	}

	var paths []string
	if httpxPath != "" {
		paths = filterWithHttpx(httpxPath, allURLs, targetURL, headers)
	}

	if len(paths) == 0 {
		fmt.Println("  [!] Httpx unavailable or no results, extracting from katana output")
		base, _ := url.Parse(targetURL)
		targetHost := base.Hostname()
		seen := make(map[string]struct{})
		for _, u := range allURLs {
			parsed, parseErr := url.Parse(u)
			if parseErr != nil || parsed.Hostname() != targetHost {
				continue
			}
			if _, exists := seen[parsed.Path]; exists {
				continue
			}
			seen[parsed.Path] = struct{}{}
			paths = append(paths, parsed.Path)
			if len(paths) >= 50 {
				break
			}
		}
	}

	if len(paths) == 0 {
		fmt.Println("  [!] Discovery yielded nothing, using root path only")
		return []string{"/"}
	}

	before := len(paths)
	scored := ScoreAndFilter(paths)
	fmt.Printf("  [Scorer] %d → %d endpoints (static assets filtered, top %d by value)\n", before, len(scored), maxScoredEndpoints)

	var finalPaths []string
	for _, s := range scored {
		finalPaths = append(finalPaths, s.Path)
	}

	result := append([]string{"/"}, finalPaths...)
	fmt.Printf("  [+] %d endpoints ready for testing\n", len(result))
	fmt.Printf("  %-42s  %5s  %s\n", "Endpoint", "Score", "Reason")
	fmt.Printf("  %-42s  %5s  %s\n", "--------", "-----", "------")
	fmt.Printf("  %-42s  %5d  %s\n", "/", 5, "root path")
	for i, s := range scored {
		if i >= 9 {
			fmt.Printf("  ... and %d more\n", len(scored)-9)
			break
		}
		label := s.Path
		if len(label) > 41 {
			label = label[:38] + "..."
		}
		fmt.Printf("  %-42s  %5d  %s\n", label, s.Score, s.Reason)
	}
	return result
}
