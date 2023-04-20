package summonerswar

import (
	"image"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
)

func (b *SummonersWar) SentESC() {
}

func (b *SummonersWar) WindowSize() domain.Rect {
	r, _ := b.screen.GetRect()
	return r
}

func (b *SummonersWar) GetMat() (gocv.Mat, error) {
	if err := b.screen.CaptureToBuffer(); err != nil {
		if errors.Is(err, domain.ErrNeedToSkipFrame) {
			b.log.Info(err.Error())
		} else {
			b.log.Errorf("capture image to buffer: %+v", err)
		}
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

func (b *SummonersWar) Click(x, y int) {
	b.screen.MouseMoveAndClick(x, y)
}

func (b *SummonersWar) ClickPt(pt image.Point) {
	b.screen.MouseMoveAndClickByPoint(pt)
}
