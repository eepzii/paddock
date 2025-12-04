package browser

import "fmt"

var baseFlags = []string{
	fmt.Sprintf("--remote-debugging-port=%d", LOOPBACK_PORT),
	"--no-first-run",
	"--no-default-browser-check",
	"--disable-background-networking",
	"--start-maximized",
}

var brandFlags = map[Brand][]string{
	CHROME: {},
	EDGE: {
		"--disable-sync",
		"--no-welcome-page",
	},
	UNKNOWN: {
		"--no-sandbox",
	},
}
