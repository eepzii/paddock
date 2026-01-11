package f1site

import (
	"errors"
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

		var hijackResult = make(chan PageResult)
		defer close(hijackResult)

		var timeoutTimer = time.NewTimer(LOGIN_EVENT_TIMEOUT_DURATION)
		defer timeoutTimer.Stop()

		page.MustNavigate(HOMEPAGE_URL).MustWaitRequestIdle()

		router := page.HijackRequests()
		router.MustAdd(REJECT_ALL_URL, handleRejectAll(&wg))
		router.MustAdd(GDPR_URL, handleGDPR(&wg))
		router.MustAdd(HOMEPAGE_URL, handleConsentCookies(&wg))
		byPasswordFunc, byPasswordChan := handleByPassword()
		router.MustAdd(BY_PASSWORD_URL, byPasswordFunc)
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

		select {
		case result := <-byPasswordChan:
			if result.Err != nil {
				return result.Err
			}
			loginResult <- result
			timeoutTimer.Stop()
		case <-timeoutTimer.C:
			return errors.New("timed out while filling form")
		}

		page.MustWait(WAIT_COOKIE_EXISTS)

		return nil
	}, loginResult
}
