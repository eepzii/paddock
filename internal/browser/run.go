package browser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func (b *Browser) Run(tasks func(page *rod.Page) error) error {
	pageDone := make(chan struct{})
	timer := time.NewTimer(PAGE_TIMEOUT_DURATION)

	args := make([]string, len(baseFlags))
	copy(args, baseFlags)

	if b.headless {
		args = append(args, "--headless=new")
	}
	if extras, ok := brandFlags[b.brand]; ok {
		args = append(args, extras...)
	}
	args = append(args, "--user-data-dir="+b.profilePath)

	cmd := exec.Command(b.commandPath, args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	defer cmd.Wait()
	defer cmd.Process.Kill()

	scanner := bufio.NewScanner(stderr)
	var port = ""
	for scanner.Scan() {
		if index := strings.Index(scanner.Text(), fmt.Sprintf("%s:%d/", LOOPBACK_ADDRESS, LOOPBACK_PORT)); index != -1 {
			if err := scanner.Err(); err != nil && err != io.EOF {
				return err
			}
			port = scanner.Text()[index:]
			break
		}
	}

	browser := rod.New().NoDefaultDevice().ControlURL(port).MustConnect()
	defer browser.Close()

	page := browser.MustPage("")

	var pageError error
	go func() {
		defer close(pageDone)
		pageError = tasks(page)
	}()

	select {
	case <-pageDone:
		timer.Stop()
	case <-timer.C:
		return errors.New("page timed out for unknown reason")
	}
	return pageError
}
