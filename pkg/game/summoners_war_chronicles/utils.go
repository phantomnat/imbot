package summonerswar

import (
	"image"

	"github.com/go-vgo/robotgo"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
)

func (b *SummonersWar) SentESC() {
	robotgo.KeyTap(robotgo.Escape)
}

func (b *SummonersWar) WindowSize() domain.Rect {
	r, _ := b.screen.GetRect()
	return r
}

func (b *SummonersWar) GetMat() (gocv.Mat, error) {
	if err := b.screen.CaptureToBuffer(); err != nil {
		b.log.Errorf("capture image to buffer: %+v", err)
		return gocv.NewMat(), err
	}

	return b.screen.GetMat()
}

func (b *SummonersWar) GetScreen() domain.Screen {
	return b.screen
}

func ptFromROIandPt(roi image.Rectangle, pt image.Point) image.Point {
	return image.Point{X: pt.X + roi.Min.X, Y: pt.Y + roi.Min.Y}
}
