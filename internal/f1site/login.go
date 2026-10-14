package f1site

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

func Login(email, password string) (func(page *rod.Page) error, <-chan PageResult) {
	var loginResult = make(chan PageResult, 1)
	return func(page *rod.Page) error {
		defer close(loginResult)

		var wg = sync.WaitGroup{}
		var wgDoneChan = make(chan struct{})
		wg.Add(3)

		var timeoutTimer = time.NewTimer(LOGIN_EVENT_TIMEOUT_DURATION)
		defer timeoutTimer.Stop()

		err := rod.Try(func() {
			page.MustNavigate(HOMEPAGE_URL).MustWaitRequestIdle()
		})
		if err != nil {
			return errors.New("couldn't navigate to homepage")
		}

		router := page.HijackRequests()
		router.MustAdd(REJECT_ALL_URL, handleRejectAll(&wg))
		router.MustAdd(GDPR_URL, handleGDPR(&wg))
		router.MustAdd(HOMEPAGE_URL, handleConsentCookies(&wg))
		customHTMLFunc, customHTMLErr := handleCustomHTML()
		router.MustAdd(CUSTOM_HTML_URL, customHTMLFunc)
		tokenFunc, tokenChan := handleToken()
		router.MustAdd(TOKEN_URL, tokenFunc)
		go router.Run()

		go func() {
			defer close(wgDoneChan)

			rod.Try(func() {
				f1CookieFrame := page.Timeout(LOGIN_EVENT_TIMEOUT_DURATION).MustElement(COOKIE_BANNER_SELECTORS.I_FRAME).MustFrame()
				f1CookieFrame.Timeout(LOGIN_EVENT_TIMEOUT_DURATION).MustElementX(COOKIE_BANNER_SELECTORS.ESSENTIAL_ONLY_BTN).MustWaitInteractable().MustClick()
			})

			wg.Wait()
		}()

		select {
		case <-wgDoneChan:
			timeoutTimer.Reset(LOGIN_EVENT_TIMEOUT_DURATION)
			// continue
		case <-timeoutTimer.C:
			return errors.New("cookie consent timeout: banner was unresponsive or network requests stalled")
		}
		page.MustNavigate(LOGIN_URL).MustWaitRequestIdle()

		page.MustElement(LOGIN_FORM_SELECTORS.EMAIL_INPUT).MustWaitInteractable().MustInput(email)
		page.MustElement(LOGIN_FORM_SELECTORS.PASSWORD_INPUT).MustWaitInteractable().MustInput(password)
		page.MustElement(LOGIN_FORM_SELECTORS.SUBMIT_BTN).MustWaitInteractable().MustClick()

		var before bool
		for {
			select {
			case err := <-customHTMLErr:
				if err != nil {
					return err
				}

				before = true

			case result := <-tokenChan:
				if result.Err != nil {
					return result.Err
				}

				if !before {
					return fmt.Errorf("failed to hijack or process: %q", CUSTOM_HTML_URL)
				}

				loginResult <- result
				timeoutTimer.Stop()

				page.MustWait(WAIT_COOKIE_EXISTS)
				return nil
			case <-timeoutTimer.C:
				return errors.New("timed out while filling form")
			}
		}
	}, loginResult
}
