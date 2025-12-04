package f1site

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func handleByPassword() (func(ctx *rod.Hijack), <-chan PageResult) {
	var loginResultChan = make(chan PageResult, 1)
	return func(ctx *rod.Hijack) {
		if ctx.Request.Method() == http.MethodPost {
			if err := ctx.LoadResponse(http.DefaultClient, true); err != nil {
				return
			}
			var loginResult = PageResult{}

			switch ctx.Response.Payload().ResponseCode {
			case http.StatusOK:
				loginResult.Response = ctx.Response
			case http.StatusUnauthorized:
				loginResult.Err = errors.New(`{"401":"unauthorized"}`)
			case http.StatusForbidden:
				loginResult.Err = errors.New(`{"403": "forbidden"}`)
			default:
				loginResult.Err = errors.New("no login response found")
			}
			loginResultChan <- loginResult
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
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
