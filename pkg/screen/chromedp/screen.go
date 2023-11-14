package chromedp

import (
	"bytes"
	"context"
	"image"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
)

type Screen struct {
	allocCtx    context.Context
	allocCancel context.CancelFunc
	taskCtx     context.Context
	taskCancel  context.CancelFunc

	log *zap.SugaredLogger

	buf []byte

	screenRect domain.Rect
}

var _ domain.Screen = (*Screen)(nil)

type Option struct {
	Username string
	Password string
}

func NewChromeDP(ctx context.Context, o Option) (*Screen, error) {
	currentDir, _ := os.Getwd()
	dataDir := path.Join(currentDir, "chromedp", "data")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dataDir),
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1534, 986),
	)
	s := &Screen{
		log: zap.S().Named("chromedp"),
		// TODO: auto find the rect
		screenRect: domain.NewRect(119, 64, 1280, 720),
	}

	s.allocCtx, s.allocCancel = chromedp.NewExecAllocator(ctx, opts...)
	s.taskCtx, s.taskCancel = chromedp.NewContext(s.allocCtx, chromedp.WithLogf(s.log.Debugf))

	err := chromedp.Run(s.taskCtx)
	if err != nil {
		s.log.Errorf("start chromedp: %+v", err)
		return nil, err
	}

	// TODO: go to redfinger and login
	// https://www.cloudemulator.net/app/sign-in/email
	ctx, cancel := context.WithTimeout(s.taskCtx, 5*time.Second)
	defer cancel()
	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate("https://www.cloudemulator.net/app/phone"),
		//chromedp.Navigate("https://www.cloudemulator.net/app/sign-in/email"),

		chromedp.WaitVisible(`div#phone-grid-style`, chromedp.ByQuery),

		chromedp.WaitVisible(`#phone-grid-style > ion-grid > ion-row > ion-col.md.hydrated.phone-col.row-last > div > div.preview`, chromedp.ByQuery),
	})
	if err != nil {
		// need to login
		s.log.Errorf("maybe need to login")
		//return nil, err
	} else {
		var iframes []*cdp.Node
		ctx2, cancel2 := context.WithTimeout(s.taskCtx, 10*time.Second)
		defer cancel2()
		_ = chromedp.Run(ctx2, chromedp.Tasks{
			//chromedp.Nodes(`#phone-grid-style > ion-grid > ion-row > ion-col.md.hydrated.phone-col.row-last > div > div.preview`, &screens, chromedp.ByQuery),
			chromedp.Click(`#phone-grid-style > ion-grid > ion-row > ion-col.md.hydrated.phone-col.row-last > div > div.preview`, chromedp.ByQuery),
			//chromedp.WaitVisible(`#phoneVideo`, chromedp.ByQuery),
			//chromedp.Screenshot(sel, res, chromedp.NodeVisible),
			chromedp.Sleep(2 * time.Second),
			chromedp.Nodes(`#webplay-container`, &iframes, chromedp.ByQuery),
		})
		if len(iframes) > 0 {
			err = chromedp.Run(ctx2, chromedp.Tasks{
				chromedp.WaitVisible(`#phoneVideo`, chromedp.ByID, chromedp.FromNode(iframes[0])),
				chromedp.Sleep(time.Second),
			})
			if err != nil {
				s.log.Errorf("cannot get phone screen")
			} else {
				s.CaptureMatAndSave(filepath.Join("cap", "chromedp-1.png"))
				//s.CaptureMatAndSave("cap/chromedp-2")
				//s.capturePhoneScreen(filepath.Join("cap", "chromedp-2.png"), iframes[0])

				// close the screen
				//_ = chromedp.Run(ctx2, chromedp.Tasks{
				//	chromedp.Sleep(500 * time.Millisecond),
				//	chromedp.Click(`#right-btns > ul.func-action > li:nth-child(9)`, chromedp.NodeVisible, chromedp.ByQuery, chromedp.FromNode(iframes[0])),
				//})
			}
			//document.querySelector("#right-btns > ul.func-action > li:nth-child(9)")
		}
		//s.log.Debugf("%#v", screens)
	}

	return s, nil
}

func (s *Screen) GetRect() (domain.Rect, error) {
	return s.screenRect, nil
}

func (s *Screen) GetMat() (gocv.Mat, error) {
	img, _, err := image.Decode(bytes.NewReader(s.buf))
	if err != nil {
		s.log.Errorf("cannot decode png image: %+v", err)
		return gocv.NewMat(), nil
	}
	src, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		s.log.Errorf("cannot convert png to gocv mat: %+v", err)
		return gocv.NewMat(), nil
	}
	defer src.Close()
	src2 := src.Region(s.screenRect.ToImage())
	defer src2.Close()

	out := gocv.NewMatWithSize(src2.Rows(), src2.Cols(), gocv.MatTypeCV8UC3)
	gocv.CvtColor(src2, &out, gocv.ColorRGBAToBGR)
	return out, nil
}

func (s *Screen) CaptureToBuffer() error {
	s.log.Debugf("capture to buffer")
	if err := chromedp.Run(s.taskCtx, chromedp.Tasks{
		//chromedp.FullScreenshot(&s.buf, 100),
		chromedp.CaptureScreenshot(&s.buf),
		//chromedp.Screenshot(`#phoneVideo`, &s.buf),
		chromedp.Sleep(50 * time.Millisecond),
	}); err != nil {
		return err
	}
	return nil
}

func (s *Screen) CaptureMatAndSave(filePath string) {
	if filePath == "" {
		filePath = filepath.Join("cap", "chromedp")
	}
	if err := s.CaptureToBuffer(); err != nil {
		s.log.Errorf("cannot capture screen to buffer: %+v", err)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(s.buf))
	if err != nil {
		s.log.Errorf("cannot decode image: %+v", err)
		return
	}
	m, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		s.log.Errorf("cannot convert png to gocv mat: %+v", err)
		return
	}
	defer m.Close()

	//crop
	m2 := m.Region(s.screenRect.ToImage())
	defer m2.Close()

	gocv.IMWrite(filePath, m2)
}

func (s *Screen) capturePhoneScreen(filePath string, iframe *cdp.Node) {
	var buf []byte
	//filePath := filepath.Join("cap", "chromedp.png")
	if err := chromedp.Run(s.taskCtx, chromedp.Tasks{
		chromedp.Screenshot(`#phoneVideo`, &buf),
	}); err != nil {
		s.log.Errorf("cannot capture screen: %+v", err)
		return
	}
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		s.log.Errorf("cannot decode image: %+v", err)
		return
	}
	m, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		s.log.Errorf("cannot convert png to gocv mat: %+v", err)
		return
	}
	gocv.IMWrite(filePath, m)
}

func (s *Screen) captureScreen(fp string) {
	var buf []byte
	if fp == "" {

		fp = filepath.Join("cap", "chromedp.png")
	}
	if err := chromedp.Run(s.taskCtx, chromedp.Tasks{
		chromedp.CaptureScreenshot(&buf),
	}); err != nil {
		s.log.Errorf("cannot capture screen: %+v", err)
		return
	}
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		s.log.Errorf("cannot decode image: %+v", err)
		return
	}
	m, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		s.log.Errorf("cannot convert png to gocv mat: %+v", err)
		return
	}
	gocv.IMWrite(fp, m)
}

func (s *Screen) MouseMoveAndClickByRect(roi image.Rectangle, args ...any) {
	x := roi.Min.X + (roi.Dx() / 2)
	y := roi.Min.Y + (roi.Dy() / 2)
	s.MouseMoveAndClick(x, y)
}

func (s *Screen) MouseMoveAndClickByPoint(pt image.Point, args ...any) {
	s.MouseMoveAndClick(pt.X, pt.Y)
}

func (s *Screen) MouseMoveAndClick(x, y int, args ...any) {
	chromedp.Run(s.taskCtx, chromedp.Tasks{
		chromedp.MouseClickXY(float64(x)+float64(s.screenRect.X), float64(y)+float64(s.screenRect.Y)),
	})
}

func (s *Screen) MouseDrag(x1, y1, x2, y2 int) {
	sx := float64(s.screenRect.X)
	sy := float64(s.screenRect.Y)
	chromedp.Run(s.taskCtx, chromedp.Tasks{
		MouseDrag(float64(x1)+sx, float64(y1)+sy, float64(x2)+sx, float64(y2)+sy),
	})
}

func (s *Screen) MouseDragDuration(x1, y1, x2, y2, waitMs int) {
	s.MouseDrag(x1, y1, x2, y2)
	time.Sleep(time.Duration(waitMs) * time.Millisecond)
}

func MouseDrag(x1, y1, x2, y2 float64) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		p := &input.DispatchMouseEventParams{
			Type:       input.MousePressed,
			X:          x1,
			Y:          y1,
			Button:     input.Left,
			ClickCount: 1,
		}

		if err := p.Do(ctx); err != nil {
			return err
		}

		// Mouse Move
		p.Type = input.MouseMoved

		for x1 != x2 || y1 != y2 {
			if x1 < x2 {
				x1++
			} else if x1 > x2 {
				x1--
			}
			if y1 < y2 {
				y1++
			} else if y1 > y2 {
				y1--
			}
			p.X = x1
			p.Y = y1
			if err := p.Do(ctx); err != nil {
				return err
			}
		}

		p.Type = input.MouseReleased
		return p.Do(ctx)
	}
}

func (s *Screen) KeyTap(key string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (s *Screen) MouseMove(x, y int) {
	//TODO implement me
	panic("implement me")
}

func (s *Screen) Back() {
	//TODO implement me
	panic("implement me")
}
