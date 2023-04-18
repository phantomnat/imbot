package mumu

import (
	"time"

	"github.com/lxn/win"
)

func (s *Screen) restoreFromMaximize() (skip bool) {
	placement := win.WINDOWPLACEMENT{}
	if win.GetWindowPlacement(s.hwnd, &placement) {
		if placement.ShowCmd == win.SW_SHOWMAXIMIZED {
			// restore window
			win.ShowWindow(s.hwnd, win.SW_RESTORE)
			time.Sleep(10 * time.Millisecond)
			skip = true
		}
	}
	return
}

func (s *Screen) restoreFromMinimize() (skip bool) {
	if win.IsIconic(s.hwnd) {
		win.ShowWindow(s.hwnd, win.SW_RESTORE)
		time.Sleep(10 * time.Millisecond)
		skip = true
	}
	return
}

func (s *Screen) resizeWindow() {
	if s.o.Height == 0 || s.o.Width == 0 {
		return
	}

	if !s.getRect() {
		return
	}

	borderW := s.windowRect.Width - s.clientRect.Width
	borderH := s.windowRect.Height - s.clientRect.Height

	win.SetWindowPos(
		s.hwnd, 
		0,
		0, 0,
		int32(s.o.Width+borderW), int32(s.o.Height+borderH),
		win.SWP_NOMOVE|win.SWP_NOOWNERZORDER|win.SWP_NOZORDER)
	time.Sleep(10 * time.Millisecond)
}

func (s *Screen) ensureScreenSize() (skip bool) {
	if !s.o.AutoResize {
		return
	}
	if s.restoreFromMaximize() {
		return true
	}
	if s.restoreFromMinimize() {
		return true
	}
	if s.clientRect.Width == s.o.Width && s.clientRect.Height == s.o.Height {
		return
	}

	s.resizeWindow()
	return true
}
