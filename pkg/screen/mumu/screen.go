package mumu

import (
	"image"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
)

type Screen struct {
	hwnd       win.HWND
	childHwnd  win.HWND
	windowRect domain.Rect
	clientRect domain.Rect
	screenRect domain.Rect
	log        *zap.SugaredLogger
	buf        domain.ImageBuffer
	o          Option
}

type Option struct {
	PackageName  string
	ActivityName string
	AutoResize   bool
	Width        int
	Height       int
}

var _ domain.Screen = (*Screen)(nil)

func NewFromTitle(title string, o Option) (*Screen, error) {
	hwnd := robotgo.FindWindow(title)
	if hwnd == 0 {
		return nil, errors.Errorf("cannot find window '%s'", title)
	}
	s := &Screen{
		hwnd: hwnd,
		log:  zap.S().Named("mumu"),
		o:    o,
	}

	printme := func(hwnd uintptr, lParam uintptr) uintptr {
		// spew.Dump(hwnd)
		s.childHwnd = win.HWND(hwnd)
		return 0
	}

	win.EnumChildWindows(hwnd, syscall.NewCallback(printme), 0)

	if s.childHwnd == 0 {
		return nil, errors.Errorf("cannot find child window '%s'", title)
	}

	rect, err := s.GetRect()
	if err != nil {
		return nil, errors.Errorf("connect get rect: %+v", err)
	}
	s.log.Infof("screen: %v", rect)
	return s, nil
}

func (s *Screen) CaptureMatAndSave(filePath string) {
	if err := s.CaptureToBuffer(); err != nil {
		s.log.Errorf("cannot capture screen to buffer: %+v", err)
		return
	}
	src, err := gocv.NewMatFromBytes(
		s.buf.Height,
		s.buf.Width,
		gocv.MatTypeCV8UC4,
		s.buf.Data,
	)
	if err != nil {
		s.log.Errorf("cannot create new image from buffer: %+v", err)
		return
	}
	defer src.Close()
	gocv.IMWrite(filePath, src)
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
	hwnd := s.childHwnd
	if hwnd == 0 {
		return
	}
	nx := x + s.screenRect.X
	ny := y + s.screenRect.Y
	s.log.Debugf("move and click at %d, %d", nx, ny)
	// s.move(nx, ny)
	// robotgo.MilliSleep(80)
	// robotgo.Click(args...)
	win.PostMessage(hwnd,
		win.WM_LBUTTONDOWN,
		0,
		uintptr(win.MAKELONG(uint16(x), uint16(y))))
	time.Sleep(60 * time.Millisecond)
	win.PostMessage(hwnd,
		win.WM_LBUTTONUP,
		0,
		uintptr(win.MAKELONG(uint16(x), uint16(y))))
}

func makeLongFromP(p image.Point) uintptr {
	return uintptr(win.MAKELONG(uint16(p.X), uint16(p.Y)))
}

func (s *Screen) MouseDrag(x1, y1, x2, y2 int) {
	hwnd := s.childHwnd
	if hwnd == 0 {
		return
	}
	win.PostMessage(hwnd, win.WM_LBUTTONDOWN, win.MK_LBUTTON, uintptr(win.MAKELONG(uint16(x1), uint16(y1))))
	time.Sleep(10 * time.Millisecond)
	for !(x1 == x2 && y1 == y2) {
		if x1 > x2 {
			x1--
		} else if x1 < x2 {
			x1++
		}
		if y1 > y2 {
			y1--
		} else if y1 < y2 {
			y1++
		}
		win.PostMessage(hwnd, win.WM_MOUSEMOVE, win.MK_LBUTTON, uintptr(win.MAKELONG(uint16(x1), uint16(y1))))
		time.Sleep(3 * time.Millisecond)
	}
	win.PostMessage(hwnd, win.WM_LBUTTONUP, win.MK_LBUTTON, uintptr(win.MAKELONG(uint16(x1), uint16(y1))))
	time.Sleep(50 * time.Millisecond)
}

func (s *Screen) KeyTap(key string, args ...any) {

}

func (s *Screen) MouseMove(x, y int) {

}

func (s *Screen) GetRect() (domain.Rect, error) {
	if !s.getRect() {
		return domain.Rect{}, errors.Errorf("cannot get window rect (hwnd: %x)", s.hwnd)
	}

	if s.ensureScreenSize() {
		return domain.Rect{}, errors.Wrap(domain.ErrNeedToSkipFrame, "ensure screen size")
	}

	// s.log.Debugf("window: %v", windowRect)
	// s.log.Debugf("client: %v", clientRect)

	marginLeft := (s.windowRect.Width - s.clientRect.Width) / 2
	marginTop := s.windowRect.Height - s.clientRect.Height - marginLeft

	s.screenRect = domain.Rect{
		X:      s.windowRect.X + marginLeft,
		Y:      s.windowRect.Y + marginTop,
		Width:  s.clientRect.Width,
		Height: s.clientRect.Height,
	}
	return s.screenRect, nil
}

func (s *Screen) getRect() bool {
	if s.hwnd == 0 {
		return false
	}
	winRect := win.RECT{}
	if !win.GetWindowRect(s.childHwnd, &winRect) {
		return false
	}
	clientRect := win.RECT{}
	if !win.GetClientRect(s.childHwnd, &clientRect) {
		return false
	}
	s.windowRect.FromRect(winRect)
	s.clientRect.FromRect(clientRect)
	return true
}

func (s *Screen) CaptureToBuffer() error {
	rect, err := s.GetRect()
	if err != nil {
		return err
	}

	dHWND := win.GetDesktopWindow()
	if dHWND == 0 {
		return nil
	}

	srcDC := win.GetDC(0)
	if srcDC == 0 {
		return nil
	}
	defer win.ReleaseDC(0, srcDC)
	dstDC := win.CreateCompatibleDC(srcDC)
	if dstDC == 0 {
		return nil
	}
	defer win.DeleteDC(dstDC)

	if s.buf.Width != rect.Width || s.buf.Height != rect.Height {
		s.buf.Data = nil
		s.buf = domain.ImageBuffer{
			Data:   make([]byte, rect.Width*rect.Height*4),
			Width:  rect.Width,
			Height: rect.Height,
		}
	}

	var biHeader = win.BITMAPINFOHEADER{
		BiSize:        uint32(reflect.TypeOf(win.BITMAPINFOHEADER{}).Size()),
		BiWidth:       int32(rect.Width),
		BiHeight:      -int32(rect.Height),
		BiPlanes:      1,
		BiBitCount:    32,
		BiCompression: win.BI_RGB,
	}
	var bitmapData = unsafe.Pointer(uintptr(0))
	bitmap := win.CreateDIBSection(dstDC, &biHeader, 0, &bitmapData, 0, 0)
	if bitmap == 0 {
		return nil
	}
	defer win.DeleteObject(win.HGDIOBJ(bitmap))

	win.SelectObject(dstDC, win.HGDIOBJ(bitmap))
	// | win.CAPTUREBLT
	if !win.BitBlt(dstDC, 0, 0, int32(rect.Width), int32(rect.Height), srcDC, int32(rect.X), int32(rect.Y), win.SRCCOPY|win.CAPTUREBLT) {
		//if !win.BitBlt(dstDC, 0, 0, width, height, srcDC, 0,0, win.SRCCOPY ) {
		return nil
	}

	// Convert the bitmap to an image.Image. We first start by directly
	// creating a slice. This is unsafe but we know the underlying structure
	// directly.
	var slice []byte
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sliceHdr.Data = uintptr(bitmapData)
	sliceHdr.Len = int(rect.Width * rect.Height * 4)
	sliceHdr.Cap = sliceHdr.Len
	copy(s.buf.Data, slice)
	return nil
}

func (s *Screen) GetMat() (gocv.Mat, error) {
	src, err := gocv.NewMatFromBytes(
		s.buf.Height,
		s.buf.Width,
		gocv.MatTypeCV8UC4,
		s.buf.Data,
	)
	if err != nil {
		return gocv.NewMat(), nil
	}
	// gocv.IMWrite("get-mat-rgba.png", src)
	// fmt.Printf("src ch: %d, type: %v\n", src.Channels(), src.Type())
	defer src.Close()

	// out := gocv.NewMat()
	// src.ConvertTo(&out, gocv.MatTypeCV8UC3)

	out := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV8UC3)
	gocv.CvtColor(src, &out, gocv.ColorRGBAToBGR)
	// gocv.IMWrite("get-mat-bgr.png", out)
	// fmt.Printf("src ch: %d, type: %v\n", out.Channels(), out.Type())

	return out, nil
}

// https://github.com/AutoHotkey/AutoHotkey/blob/master/Source/keyboard_mouse.cpp#L2285
func (s *Screen) move(x, y int) {
	screenWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	screenHeight := win.GetSystemMetrics(win.SM_CYSCREEN)

	aX := ((65536 * int32(x)) / screenWidth) + 1
	aY := (((65536 * int32(y)) / screenHeight) + 1)

	input := win.MOUSE_INPUT{
		Type: win.INPUT_MOUSE,
		Mi: win.MOUSEINPUT{
			Dx:      aX,
			Dy:      aY,
			DwFlags: win.MOUSEEVENTF_MOVE | win.MOUSEEVENTF_ABSOLUTE,
		},
	}
	win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
}
