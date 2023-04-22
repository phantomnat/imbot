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
	c.log.Debugf("sending %s", cmd)
	io.WriteString(c.sessionWriter, cmd)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	// a := make([]string, 0, len(args)+4)
	// a = append(a, "-s", c.deviceURL(), "shell", in)
	// if len(args) > 0 {
	// 	a = append(a, args...)
	// }
	// cmd := exec.CommandContext(ctx, "adb", a...)
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	HideWindow: true,
	// }
	// _, _ = cmd.CombinedOutput()
	// c.log.Debugf("click took %d ms", time.Since(start).Milliseconds())
	// buf := &bytes.Buffer{}
	// session := exec.Command("adb", "shell")
	// session.SysProcAttr.HideWindow = true
	// session.Stdin = buf

	// session.Start()
	// go func() {
	// 	_, _ = io.Copy(os.Stdout, session.Process)
	// }()
	// fmt.Fprintf(buf, "")
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
