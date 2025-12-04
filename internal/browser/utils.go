package browser

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

func findBrowser() (string, error) {
	var paths []string

	switch OperatingSystem(runtime.GOOS) {
	case LINUX:
		paths = []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",

			"/usr/bin/microsoft-edge",
			"/usr/bin/microsoft-edge-stable",
		}
	case WINDOWS:
		paths = []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,

			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
		}
	case DARWIN:
		paths = []string{
			`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`,

			`/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge`,
		}
	default:
		return "", errors.New("cannot find a suitable os")
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		return path, nil
	}

	return "", errors.New("no compatible browser was found")
}

func findBrowserBrand(browserPath string) Brand {
	browserPath = strings.ToLower(browserPath)

	if strings.Contains(browserPath, string(CHROME)) {
		return CHROME
	} else if strings.Contains(browserPath, string(EDGE)) {
		return EDGE
	}
	return UNKNOWN
}
