package summonerswar

import (
	"image"
	"time"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
)

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

func (b *SummonersWar) Back() {
	b.screen.Back()
}

func (b *SummonersWar) Click(x, y int) {
	b.screen.MouseMoveAndClick(x, y)
}

func (b *SummonersWar) ClickPt(pt image.Point) {
	b.screen.MouseMoveAndClickByPoint(pt)
}

func (b *SummonersWar) Drag(pt1, pt2 image.Point) {
	b.screen.MouseDrag(pt1.X, pt1.Y, pt2.X, pt2.Y)
}
func (b *SummonersWar) DragDuration(pt1, pt2 image.Point, waitMs int) {
	b.screen.MouseDragDuration(pt1.X, pt1.Y, pt2.X, pt2.Y, waitMs)
}

func (b *SummonersWar) GoToMainScreen(m gocv.Mat) (done bool) {
	foundCoin, _ := b.MatchInROI(m, roi.ROIMainScreen.CoinIcon, domain.MatchOption{
		Path: "icon_coin",
	})
	foundCrystal, _ := b.MatchInROI(m, roi.ROIMainScreen.CrystalIcon, domain.MatchOption{
		Path: "icon_crystal",
	})

	if foundCoin && foundCrystal {
		done = true
		return
	}

	// detect top right icon
	foundTopRightIcon, _ := b.MatchInROI(m, roi.ROITopRigthMenuBtn, domain.MatchOption{
		Path: "btn_top_right_menu",
	})
	if !foundTopRightIcon {
		b.Back()
		b.SleepMs(1000)
		return
	}

	return
}

func (b *SummonersWar) SleepMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (b *SummonersWar) GetImageManager() domain.ImageManager {
	return b.im
}
