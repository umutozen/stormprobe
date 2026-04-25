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

func Discover(targetURL, katanaCustom, httpxCustom string, noDiscovery bool) []string {
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

	allURLs, err := crawlWithKatana(katanaPath, targetURL)
	if err != nil || len(allURLs) == 0 {
		fmt.Println("  [!] Katana yielded no results, using root path only")
		return []string{"/"}
	}

	var paths []string
	if httpxPath != "" {
		paths = filterWithHttpx(httpxPath, allURLs, targetURL)
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

	result := append([]string{"/"}, paths...)
	fmt.Printf("  [+] %d endpoints ready for testing\n", len(result))
	for i, p := range result {
		if i >= 10 {
			fmt.Printf("  ... and %d more\n", len(result)-10)
			break
		}
		fmt.Printf("  > %s\n", p)
	}
	return result
}
