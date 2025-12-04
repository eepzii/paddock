package f1site

import "github.com/go-rod/rod"

type PageResult struct {
	Response *rod.HijackResponse
	Err      error
}
