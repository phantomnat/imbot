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

func (b *SummonersWar) WaitMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (b *SummonersWar) GetImageManager() domain.ImageManager {
	return b.im
}

func (b *SummonersWar) GoToMainScreen(m gocv.Mat) (done bool) {
	foundCoin, _ := b.MatchInROI(m, roi.MainScreen.CoinIcon, domain.MatchOption{
		Path: "icon_coin",
	})
	foundCrystal, _ := b.MatchInROI(m, roi.MainScreen.CrystalIcon, domain.MatchOption{
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
		b.WaitMs(1000)
		return
	}

	return
}

func (b *SummonersWar) IsOnMainScreen(m gocv.Mat) bool {
	foundCoin, _ := b.MatchInROI(m, roi.MainScreen.CoinIcon, domain.MatchOption{
		Path: "icon_coin",
		//PrintVal: true,
	})
	foundCrystal, _ := b.MatchInROI(m, roi.MainScreen.CrystalIcon, domain.MatchOption{
		Path: "icon_crystal",
		//PrintVal: true,
	})

	if foundCoin && foundCrystal {
		return true
	}

	return false
}

func (b *SummonersWar) IsOnMainMenu(m gocv.Mat) bool {
	found, _ := b.MatchInROI(m, roi.MainMenu.OfficialForum, domain.MatchOption{
		Path: "menu.official_forum",
	})

	return found
}

func (b *SummonersWar) HandleConversationDialog(m gocv.Mat) bool {
	prefix := "dialog"
	foundBtnBack, _ := b.MatchInROI(m, roi.ROIBtnBack, domain.MatchOption{
		Path: prefix + ".btn_back",
		// Th:   0.01,
		//PrintVal: true,
	})
	foundTxtAuto, _ := b.MatchInROI(m, roi.ROITxtAutoAndIcon, domain.MatchOption{
		Path: prefix + ".txt_auto_and_icon",
		// Th:       0.01,
		//PrintVal: true,
	})

	if !(foundTxtAuto && foundBtnBack) {
		return false
	}

	b.log.Infof("dialog detected")
	b.screen.MouseMoveAndClickByPoint(roi.PtSkipDialog)
	b.WaitMs(600)
	return true
}

func (b *SummonersWar) HandleQuestCompleted(m gocv.Mat) bool {
	prefix := "quest_complete"

	// complete modal
	{
		foundTapToClose, _ := b.MatchInROI(m, roi.QuestCompleted.TapToClose, domain.MatchOption{
			Path: prefix + ".txt_tab_to_close",
			// Th:   0.01,
			// PrintVal: true,
		})
		if foundTapToClose {
			b.log.Infof("quest completed, click anywhere to close")
			b.ClickPt(roi.QuestCompleted.PtOK)
			waitMs(1500)
			return true
		}
	}

	// complete dialog

	mQuestCompleteButtons := m.Region(roi.QuestCompleted.Buttons)
	defer mQuestCompleteButtons.Close()

	foundBtnOK, ptBtnOk := b.imMatchDefault(mQuestCompleteButtons, prefix, "btn_ok")
	foundBtnNextStory, ptBtnNextStory := b.imMatchDefault(mQuestCompleteButtons, prefix, "btn_next_story")

	var pt image.Point
	var btn string
	var found bool

	switch {
	case foundBtnNextStory:
		pt = ptBtnNextStory
		btn = "next story"
		found = true
	case foundBtnOK:
		pt = ptBtnOk
		btn = "ok"
		found = true
	}
	if !found {
		return false
	}

	b.log.Infof("quest completed, click %s (%v)", btn, pt)
	b.ClickPt(ptFromROIandPt(roi.QuestCompleted.Buttons, pt))
	b.WaitMs(500)

	return true
}

func (b *SummonersWar) HandleVictory(m gocv.Mat) (done bool) {
	prefix := "quest_complete"
	done = true

	// victory
	{
		foundTapToClose, pt := b.MatchInROI(m, roi.ROIVictoryButtons, domain.MatchOption{
			Path: prefix + ".btn_victory_exit",
			// Th:   0.01,
			// PrintVal: true,
		})
		if foundTapToClose {
			b.log.Infof("quest completed, click to exit")
			b.ClickPt(pt)
			waitMs(1500)
			return
		}
	}

	done = false
	return
}
