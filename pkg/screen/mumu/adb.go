package mumu

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ADBClient struct {
	devicePort    string
	deviceHost    string
	log           *zap.SugaredLogger
	session       *exec.Cmd
	sessionWriter io.Writer
}

func NewADBClient(parentCtx context.Context, adbPort int) (*ADBClient, error) {
	c := &ADBClient{
		deviceHost: "127.0.0.1",
		devicePort: strconv.Itoa(adbPort),
		log:        zap.S().Named("adb-client"),
	}

	ctx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "adb", "connect", c.deviceURL())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "adb connect")
	} else if !strings.Contains(string(out), "connected") {
		return nil, errors.Errorf("adb cannot connect to '%s'", c.deviceURL())
	}

	c.session = exec.Command("adb", "-s", c.deviceURL(), "shell")
	c.session.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	reader, writer := io.Pipe()
	c.sessionWriter = writer
	c.session.Stdin = reader
	c.session.Stdout = os.Stdout
	c.session.Stderr = os.Stderr
	go func() {
		err = c.session.Run()
		if err != nil {
			c.log.Errorf("create shell session: %+v", err)
		}
	}()
	return c, nil
}

func (c *ADBClient) deviceURL() string {
	return c.deviceHost + ":" + c.devicePort
}

func (c *ADBClient) RunShell(in string, args ...string) {
	// start := time.Now()
	cmd := in + " " + strings.Join(args, " ") + "\n"
	// c.log.Debugf("sending %s", cmd)
	_, _ = io.WriteString(c.sessionWriter, cmd)
}

func (c *ADBClient) Tap(x, y int) {
	c.RunShell("input", "tap", strconv.Itoa(x), strconv.Itoa(y))
}

func (c *ADBClient) Swipe(x1, y1, x2, y2, durationMs int) {
	c.RunShell("input", "draganddrop", strconv.Itoa(x1), strconv.Itoa(y1), strconv.Itoa(x2), strconv.Itoa(y2), strconv.Itoa(durationMs))
	time.Sleep(time.Duration(durationMs) * time.Millisecond)
}

func (c *ADBClient) Back() {
	c.RunShell("input", "keyevent", "4")
}
