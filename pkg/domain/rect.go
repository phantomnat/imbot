package domain

import (
	"fmt"

	"github.com/lxn/win"
)

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r Rect) String() string {
	var zero Rect
	if r == zero {
		return "{}"
	}
	return fmt.Sprintf("{x: %d, y: %d, w: %d, h: %d}", r.X, r.Y, r.Width, r.Height)
}

func (r *Rect) FromRect(rect win.RECT) *Rect {
	if r == nil {
		return nil
	}
	r.X = int(rect.Left)
	r.Y = int(rect.Top)
	r.Width = int(rect.Right - rect.Left)
	r.Height = int(rect.Bottom - rect.Top)
	return r
}
