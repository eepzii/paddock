package f1site

import (
	"errors"
	"time"

	"github.com/go-rod/rod"
)

func Logout(page *rod.Page) error {

	var timeoutTimer = time.NewTimer(LOGOUT_TIMEOUT_DURATION)
	var done = make(chan struct{})
	var errChan = make(chan error, 1)
	go func() {
		defer close(done)
		err := rod.Try(func() {
			page.MustNavigate(HOMEPAGE_URL).MustWaitRequestIdle()
			page.MustNavigate(LOGOUT_URL).MustWaitRequestIdle()
		})
		if err != nil {
			errChan <- err
			return
		}
		page.MustWait(WAIT_COOKIE_GONE)
	}()

	select {
	case <-done:
		timeoutTimer.Stop()
	case <-errChan:
		return errors.New("couldn't navigate to homepage or logout url")
	case <-timeoutTimer.C:
		return errors.New("timed out")
	}

	return nil
}
