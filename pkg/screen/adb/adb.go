package adb

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
)

type Screen struct {
	log       *zap.SugaredLogger
	buf       domain.ImageBuffer
	imgBuf    *bytes.Buffer
	o         Option
	lineBreak []byte

	session       *exec.Cmd
	sessionWriter io.Writer
	sessionReader io.Reader

	screenStream *exec.Cmd
	screenWriter io.Writer
	screenReader io.Reader
	eventCapture chan bool
}

type Option struct {
	Ctx           context.Context
	PackageName   string
	ActivityName  string
	Width         int
	Height        int
	MouseMargin   image.Point
	ADBPort       int
	ADBHost       string
	ADBDeviceName string
}

var _ domain.Screen = (*Screen)(nil)

func NewADBClient(o Option) (*Screen, error) {
	s := &Screen{
		o:            o,
		log:          zap.S().Named("adb"),
		eventCapture: make(chan bool),
		imgBuf:       &bytes.Buffer{},
		lineBreak:    []byte("---lb---"),
	}
	s.imgBuf.Grow(4 << 20)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "adb", "connect", s.deviceURL())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "adb connect")
	} else if !strings.Contains(string(out), "connected") {
		return nil, errors.Errorf("adb cannot connect to '%s'", s.deviceURL())
	}

	s.screenStream = exec.Command("adb", "-s", s.deviceURL(), "shell")
	s.screenStream.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	{
		// initialize adb shell writer
		r, w := io.Pipe()
		s.screenWriter = w
		s.screenStream.Stdin = r
		// s.screenStream.Stdout = os.Stdout
		s.screenStream.Stderr = os.Stdout
	}
	{
		// initialize adb shell reader
		s.screenReader, err = s.screenStream.StdoutPipe()
		if err != nil {
			s.log.Errorf("create stdout pipe: %+v", err)
		}
	}

	go func() {
		err := s.screenStream.Run()
		if err != nil {
			s.log.Errorf("create adb shell: %+v", err)
		}
	}()

	go s.capture()

	// go func() {
	// 	s.log.Debugf("start screen reader")
	// 	lineBreak := []byte("---lb---")
	// 	buf := bufio.NewReader(s.screenReader)
	// 	b64Buf := &bytes.Buffer{}
	// 	b64Buf.Grow(4 << 20)
	// 	rawGzip := make([]byte, 0, 2<<20)
	// 	for {
	// 		data, err := buf.ReadBytes('\n')
	// 		if err != nil {
	// 			s.log.Errorf("read data from screen reader: %+v", err)
	// 			continue
	// 		}
	// 		if bytes.Equal(bytes.Trim(data, "\r\n"), lineBreak) {
	// 			// line break found
	// 			func() {
	// 				if b64Buf.Len() == 0 {
	// 					return
	// 				}

	// 				s.log.Infof("b64 data '%d' bytes", b64Buf.Len())
	// 				defer b64Buf.Reset()

	// 				n, err := base64.StdEncoding.Decode(rawGzip[:cap(rawGzip)], b64Buf.Bytes())
	// 				if err != nil {
	// 					s.log.Errorf("decode base64: %+v", err)
	// 					return
	// 				}
	// 				rawGzip = rawGzip[:n]

	// 				s.log.Infof("gzip data '%d' bytes", len(rawGzip))

	// 				gzReader, err := gzip.NewReader(bytes.NewReader(rawGzip))
	// 				if err != nil {
	// 					s.log.Errorf("create gzip reader: %+v", err)
	// 					return
	// 				}

	// 				imgBuf, err := io.ReadAll(gzReader)
	// 				if err != nil {
	// 					s.log.Errorf("decompress gzip: %+v", err)
	// 					return
	// 				}
	// 				s.log.Infof("image data '%d' bytes", len(imgBuf))

	// 				src, err := gocv.NewMatFromBytes(
	// 					720,
	// 					1280,
	// 					gocv.MatTypeCV8UC4,
	// 					imgBuf[16:],
	// 				)
	// 				if err != nil {
	// 					s.log.Errorf("create gocv image: %+v", err)
	// 					return
	// 				}
	// 				defer src.Close()
	// 				out := gocv.NewMat()
	// 				defer out.Close()

	// 				gocv.CvtColor(src, &out, gocv.ColorRGBAToBGR)

	// 				gocv.IMWrite("t1.png", src)
	// 				gocv.IMWrite("t2.png", out)
	// 			}()

	// 		} else {
	// 			_, err = b64Buf.Write(data[:len(data)-2])
	// 			if err != nil {
	// 				s.log.Errorf("cannot write to gzip buffer: %+v", err)
	// 				b64Buf.Reset()
	// 				continue
	// 			}
	// 		}
	// 	}
	// }()

	return s, nil
}

func (s *Screen) deviceURL() string {
	return s.o.ADBHost + ":" + strconv.Itoa(s.o.ADBPort)
}

func (s *Screen) GetRect() (domain.Rect, error) {
	return domain.Rect{}, nil
}

func (s *Screen) GetMat() (gocv.Mat, error) {
	src, err := gocv.NewMatFromBytes(
		s.o.Height,
		s.o.Width,
		gocv.MatTypeCV8UC4,
		s.imgBuf.Bytes()[16:],
	)
	if err != nil {
		s.log.Errorf("create gocv image: %+v", err)
		return gocv.NewMat(), nil
	}
	defer src.Close()
	out := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV8UC3)
	gocv.CvtColor(src, &out, gocv.ColorRGBAToBGR)
	return out, nil
}

func (s *Screen) capture() {
	s.log.Debugf("start streaming screen")

	buf := bufio.NewReader(s.screenReader)
	gzipBuf := &bytes.Buffer{}
	gzipBuf.Grow(4 << 20)

	for {
		select {
		case <-s.o.Ctx.Done():
			s.log.Debugf("stop streaming screen")
			return
		default:
		}

		data, err := buf.ReadBytes('\n')
		if err != nil {
			s.log.Errorf("read data from screen reader: %+v", err)
			continue
		}

		if bytes.Equal(bytes.Trim(data, "\r\n"), s.lineBreak) {
			// line break found
			result := func() bool {
				if gzipBuf.Len() == 0 {
					return true
				}
				gzipBuf.Truncate(gzipBuf.Len() - 1)
				s.log.Infof("gzip data '%d' bytes", gzipBuf.Len())

				gzReader, err := gzip.NewReader(gzipBuf)
				if err != nil {
					s.log.Errorf("create gzip reader: %+v", err)
					return false
				}

				_, err = io.Copy(s.imgBuf, gzReader)
				if err != nil {
					s.log.Errorf("decompress gzip: %+v", err)
					return false
				}
				s.log.Infof("image data '%d' bytes", s.imgBuf.Len())
				return true
			}()
			s.eventCapture <- result

			gzipBuf.Reset()
			buf.Reset(s.screenReader)
		} else {
			data = data[:len(data)-1]
			data[len(data)-1] = '\n'
			_, err = gzipBuf.Write(data)
			if err != nil {
				s.log.Errorf("cannot write to gzip buffer: %+v", err)
				gzipBuf.Reset()
				continue
			}
		}
	}
}

func (s *Screen) CaptureToBuffer() error {
	fmt.Fprintf(s.screenWriter, "screencap | gzip -2 && echo\n")
	fmt.Fprintf(s.screenWriter, "echo ---lb---\n")
	// wait for finished
	result := <-s.eventCapture
	if !result {
		return errors.New("cannot capture")
	}
	return nil
}

func (s *Screen) CaptureMatAndSave(filePath string) {
	if err := s.CaptureToBuffer(); err != nil {
		s.log.Errorf("cannot capture screen to buffer: %+v", err)
		return
	}

	src, err := gocv.NewMatFromBytes(
		s.o.Height,
		s.o.Width,
		gocv.MatTypeCV8UC4,
		s.imgBuf.Bytes()[16:],
	)
	if err != nil {
		s.log.Errorf("create gocv image: %+v", err)
		return
	}
	defer src.Close()
	out := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV8UC3)
	defer out.Close()
	gocv.CvtColor(src, &out, gocv.ColorRGBAToBGR)
	gocv.IMWrite(filePath, out)
	s.imgBuf.Reset()
}

func (s *Screen) MouseMoveAndClickByRect(roi image.Rectangle, args ...any) {
}
func (s *Screen) MouseMoveAndClickByPoint(pt image.Point, args ...any) {
}
func (s *Screen) MouseMoveAndClick(x, y int, args ...any) {
}
func (s *Screen) MouseDrag(x1, y1, x2, y2 int) {
}
func (s *Screen) MouseDragDuration(x1, y1, x2, y2, waitMs int) {
}
func (s *Screen) KeyTap(key string, args ...any) {
}
func (s *Screen) MouseMove(x, y int) {
}
func (s *Screen) Back() {
}
