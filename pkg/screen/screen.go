package screen

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Screen struct {
	hwnd       win.HWND
	windowRect win.RECT
	clientRect win.RECT
	log        *zap.SugaredLogger
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

func GetStringRECT(r win.RECT) string {
	w := r.Right - r.Left
	h := r.Bottom - r.Top
	return fmt.Sprintf("x: %d, y: %d, w: %d, h: %d", r.Left, r.Top, w, h)
}

func (s *Screen) GetRect() *domain.Rect {
	if !s.getRect() {
		return nil
	}

	s.log.Debugf("window: %v", GetStringRECT(s.windowRect))
	s.log.Debugf("client: %v", GetStringRECT(s.clientRect))

	clientWidth := int(s.clientRect.Right - s.clientRect.Left)
	windowWidth := int(s.windowRect.Right - s.windowRect.Left)
	clientHeight := int(s.clientRect.Bottom - s.clientRect.Top)
	windowHeight := int(s.windowRect.Bottom - s.windowRect.Top)
	marginLeft := (windowWidth - clientWidth) / 2
	marginTop := windowHeight - clientHeight - marginLeft

	return &domain.Rect{
		X:      int(s.windowRect.Left) + marginLeft,
		Y:      int(s.windowRect.Top) + marginTop,
		Width:  clientWidth,
		Height: clientHeight,
	}
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
