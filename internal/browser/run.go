package browser

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/go-rod/rod"
)

func (b *Browser) Run(tasks func(page *rod.Page) error) error {
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
	defer cmd.Process.Kill()
	defer cmd.Wait()

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

	return tasks(page)
}
