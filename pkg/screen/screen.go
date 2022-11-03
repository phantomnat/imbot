package screen

import (
	"fmt"
	"reflect"
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

	s.log.Debugf("window: %v", windowRect)
	s.log.Debugf("client: %v", clientRect)

	marginLeft := (windowRect.Width - clientRect.Width) / 2
	marginTop := windowRect.Height - clientRect.Height - marginLeft

	return &domain.Rect{
		X:      windowRect.X + marginLeft,
		Y:      windowRect.Y + marginTop,
		Width:  clientRect.Width,
		Height: clientRect.Height,
	}, nil
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
	gocv.IMWrite("get-mat-rgba.png", src)
	fmt.Printf("src ch: %d, type: %v\n", src.Channels(), src.Type())
	defer src.Close()

	// out := gocv.NewMat()
	// src.ConvertTo(&out, gocv.MatTypeCV8UC3)

	
	out := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV8UC3)
	gocv.CvtColor(src, &out, gocv.ColorRGBAToBGR)
	gocv.IMWrite("get-mat-bgr.png", out)
	fmt.Printf("src ch: %d, type: %v\n", out.Channels(), out.Type())
	
	return out, nil
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
