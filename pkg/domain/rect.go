package domain

import "fmt"

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
