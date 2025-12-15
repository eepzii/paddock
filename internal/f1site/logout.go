package f1site

import (
	"errors"
	"time"

	"github.com/go-rod/rod"
)

func Logout(page *rod.Page) error {

	var timeoutTimer = time.NewTimer(LOGOUT_TIMEOUT_DURATION)
	var done = make(chan struct{})

	go func() {
		defer close(done)
		page.MustNavigate(HOMEPAGE_URL).MustWaitRequestIdle()
		page.MustNavigate(LOGOUT_URL).MustWaitRequestIdle()
		page.MustWait(WAIT_COOKIE_GONE)
	}()

	select {
	case <-done:
		timeoutTimer.Stop()
	case <-timeoutTimer.C:
		return errors.New("timed out")
	}

	return nil
}
