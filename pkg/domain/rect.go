package domain

import (
	"fmt"
	"image"

	"github.com/lxn/win"
)

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func NewRect(x, y, w, h int) Rect {
	return Rect{X: x, Y: y, Width: w, Height: h}
}

func (r *Rect) String() string {
	var zero Rect
	if r == nil || *r == zero {
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

func (r *Rect) ToImage() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
}
