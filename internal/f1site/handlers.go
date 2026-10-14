package f1site

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func handleCustomHTML() (func(ctx *rod.Hijack), <-chan error) {
	var errChan = make(chan error, 1)
	return func(ctx *rod.Hijack) {
		if ctx.Request.Method() == http.MethodPost {
			if err := ctx.LoadResponse(http.DefaultClient, true); err != nil {
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
				return
			}
			var err error

			if code := ctx.Response.Payload().ResponseCode; code != http.StatusOK {
				err = fmt.Errorf("invalid password or email address (status: %d)", code)
			}

			errChan <- err
		} else {
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
		}
	}, errChan
}

func handleToken() (func(ctx *rod.Hijack), <-chan PageResult) {
	var loginResultChan = make(chan PageResult, 1)
	return func(ctx *rod.Hijack) {
		if ctx.Request.Method() == http.MethodPost {
			if err := ctx.LoadResponse(http.DefaultClient, true); err != nil {
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
				return
			}
			var loginResult = PageResult{}

			if code := ctx.Response.Payload().ResponseCode; code != http.StatusOK {
				loginResult.Err = fmt.Errorf("unexpected response code: %v", code)
			} else {
				loginResult.Response = ctx.Response
			}

			loginResultChan <- loginResult
		} else {
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
		}
	}, loginResultChan
}

func handleRejectAll(wg *sync.WaitGroup) func(ctx *rod.Hijack) {
	return func(ctx *rod.Hijack) {
		if ctx.Request.Method() == http.MethodGet {
			wg.Done()
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	}
}

func handleGDPR(wg *sync.WaitGroup) func(ctx *rod.Hijack) {
	return func(ctx *rod.Hijack) {
		if ctx.Request.Method() == http.MethodPost {
			wg.Done()
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	}
}

func handleConsentCookies(wg *sync.WaitGroup) func(ctx *rod.Hijack) {
	return func(ctx *rod.Hijack) {
		var consentUUID, consentDate bool
		cookieHeader := ctx.Request.Header("Cookie")
		cookieHeaders := strings.SplitSeq(cookieHeader, ";")
		for cookie := range cookieHeaders {
			if strings.Contains(cookie, "consentUUID=") {
				consentUUID = true
			}
			if strings.Contains(cookie, "consentDate=") {
				consentDate = true
			}
		}
		if consentUUID && consentDate {
			wg.Done()
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	}
}
