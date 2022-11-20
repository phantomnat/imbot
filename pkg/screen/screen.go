package screen

import (
	"image"
	"reflect"
	"time"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

type imageBuffer struct {
	data   []byte
	width  int
	height int
}
type Screen struct {
	hwnd       win.HWND
	windowRect win.RECT
	clientRect win.RECT
	screenRect domain.Rect
	log        *zap.SugaredLogger
	buf        imageBuffer
}

func NewFromTitle(title string) (*Screen, error) {
	hwnd := robotgo.FindWindow(title)
	if hwnd == 0 {
		return nil, errors.Errorf("cannot find window '%s'", title)
	}
	s := &Screen{
		hwnd: hwnd,
		log:  zap.S().Named("screen"),
	}
	return s, nil
}

func (s *Screen) GetRect() (*domain.Rect, error) {
	if !s.getRect() {
		return nil, errors.Errorf("cannot get window rect (hwnd: %x)", s.hwnd)
	}

	windowRect := (&domain.Rect{}).FromRect(s.windowRect)
	clientRect := (&domain.Rect{}).FromRect(s.clientRect)

	// s.log.Debugf("window: %v", windowRect)
	// s.log.Debugf("client: %v", clientRect)

	marginLeft := (windowRect.Width - clientRect.Width) / 2
	marginTop := windowRect.Height - clientRect.Height - marginLeft

	s.screenRect = domain.Rect{
		X:      windowRect.X + marginLeft,
		Y:      windowRect.Y + marginTop,
		Width:  clientRect.Width,
		Height: clientRect.Height,
	}
	return &s.screenRect, nil
}

func (s *Screen) GetCurrentCursorPos() image.Point {
	if !s.getRect() {
		return image.Point{}
	}

	pt := win.POINT{}
	if !win.GetCursorPos(&pt) {
		return image.Point{}
	}
	s.log.Infof("mouse pos image.Pt(%d, %d)", pt.X, pt.Y)
	if !win.ScreenToClient(s.hwnd, &pt) {
		return image.Point{}
	}
	s.log.Infof("mouse client pos image.Pt(%d, %d)", pt.X, pt.Y)
	return image.Pt(int(pt.X), int(pt.Y))
}

func (s *Screen) getRect() bool {
	if s.hwnd == 0 {
		return false
	}
	if !win.GetWindowRect(s.hwnd, &s.windowRect) {
		return false
	}
	if !win.GetClientRect(s.hwnd, &s.clientRect) {
		return false
	}
	return true
}

func (s *Screen) GetMat() (gocv.Mat, error) {
	src, err := gocv.NewMatFromBytes(
		s.buf.height,
		s.buf.width,
		gocv.MatTypeCV8UC4,
		s.buf.data,
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

func (s *Screen) CaptureMatAndSave(filePath string) {
	if err := s.CaptureToBuffer(); err != nil {
		s.log.Errorf("cannot capture screen to buffer: %+v", err)
		return
	}
	src, err := gocv.NewMatFromBytes(
		s.buf.height,
		s.buf.width,
		gocv.MatTypeCV8UC4,
		s.buf.data,
	)
	if err != nil {
		s.log.Errorf("cannot create new image from buffer: %+v", err)
		return
	}
	defer src.Close()
	gocv.IMWrite(filePath, src)
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

	if s.buf.width != rect.Width || s.buf.height != rect.Height {
		s.buf.data = nil
		s.buf = imageBuffer{
			data:   make([]byte, rect.Width*rect.Height*4),
			width:  rect.Width,
			height: rect.Height,
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
	copy(s.buf.data, slice)
	return nil
}

func (s *Screen) KeyTap(key string, args ...any) {
	//robotgo.KeyTap(key, args...)
	s.toggleKeyCode(uint(key[0]), false)
	time.Sleep(5 * time.Millisecond)
	s.toggleKeyCode(uint(key[0]), true)
}

//func (s *Screen) KeyPress(key string, args ...any) {
//	s.log.Debugf("pressing key %s", key)
//	robotgo.KeyPress(key, args...)
//}
//
//func (s *Screen) UnicodeType(key string) {
//	s.unicodeType(key)
//	time.Sleep(50 * time.Millisecond)
//}

func (s *Screen) toggleKeyCode(code uint, up bool) {
	var dwFlags uint
	if up {
		dwFlags = win.KEYEVENTF_KEYUP
	}
	s.keyEvent(code, dwFlags)
}

func (s *Screen) keyEvent(key, dwFlags uint) {

	//aLayout := GetKeyboardLayout(win.GetWindowThreadProcessId(win.GetForegroundWindow(), nil))
	//layoutName := GetKeyboardLayoutNameW()
	vk := VkKeyScanW(uint16(key))
	// force us layout
	usLayout := HKL(0x4090409)
	scan := MapVirtualKeyEx(uint(vk)&0xFF, MAPVK_VK_TO_VSC, usLayout)

	//s.log.Infof("layout: %x, vk: %c (%x), scan: %c (%x)", usLayout, vk, vk, scan, scan)

	input := win.KEYBD_INPUT{
		Type: win.INPUT_KEYBOARD,
		Ki: win.KEYBDINPUT{
			WScan:   uint16(scan),
			WVk:     uint16(vk),
			DwFlags: uint32(dwFlags),
		},
	}
	win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
}

func LoByte(x uint32) byte { return uint8(x & 0xff) }

func (s *Screen) unicodeType(key string) {
	var input = make([]win.KEYBD_INPUT, 0, len(key)*2)

	for _, r := range key {
		input = append(input, win.KEYBD_INPUT{
			Type: win.INPUT_KEYBOARD,
			Ki: win.KEYBDINPUT{
				WScan:   uint16(r),
				DwFlags: win.KEYEVENTF_UNICODE,
			},
		})

		input = append(input, win.KEYBD_INPUT{
			Type: win.INPUT_KEYBOARD,
			Ki: win.KEYBDINPUT{
				WScan:   uint16(r),
				DwFlags: win.KEYEVENTF_UNICODE | win.KEYEVENTF_KEYUP,
			},
		})
	}

	win.SendInput(uint32(len(input)), unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
}

func (s *Screen) MouseMoveAndClickByPoint(pt image.Point, args ...any) {
	s.MouseMoveAndClick(pt.X, pt.Y, args...)
}

func (s *Screen) MouseMoveAndClick(x, y int, args ...any) {
	nx := x + s.screenRect.X
	ny := y + s.screenRect.Y
	s.log.Debugf("move and click at %d, %d", nx, ny)
	s.move(nx, ny)
	robotgo.MilliSleep(50)
	robotgo.Click(args...)
}

func (s *Screen) MouseMove(x, y int) {
	nx := x + s.screenRect.X
	ny := y + s.screenRect.Y
	s.move(nx, ny)
}

func (s *Screen) MouseDrag(x1, y1, x2, y2 int) {
	x1 += s.screenRect.X
	y1 += s.screenRect.Y
	x2 += s.screenRect.X
	y2 += s.screenRect.Y

	s.move(x1, y1)
	time.Sleep(100 * time.Millisecond)
	robotgo.Toggle("left")
	time.Sleep(100 * time.Millisecond)
	s.moveMouseIncremental(x1, y1, x2, y2)
	time.Sleep(50 * time.Millisecond)
	robotgo.Toggle("left", "up")
}

func (s *Screen) moveMouseIncremental(x1, y1, x2, y2 int) {
	speed := 50
	minMouseSpeed := 32
	calcDelta := func(v1, v2 int) int {
		delta := 0
		if v1 < v2 {
			delta = (v2 - v1) / speed
			if delta == 0 || delta < minMouseSpeed {
				delta = minMouseSpeed
			}
			if v1+delta > v2 {
				delta = v2 - v1
			}
		} else if v1 > v2 {
			delta = (v1 - v2) / speed
			if delta == 0 || delta < minMouseSpeed {
				delta = minMouseSpeed
			}
			if v1-delta < v2 {
				delta = v1 - v2
			}
			return -delta
		}
		return delta
	}

	for x1 != x2 || y1 != y2 {
		x1 += calcDelta(x1, x2)
		y1 += calcDelta(y1, y2)

		s.move(x1, y1)
		s.doMouseDelay()
	}
}

func (s *Screen) doMouseDelay() {
	time.Sleep(5 * time.Millisecond)
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
