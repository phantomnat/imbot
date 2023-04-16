package summonerswar

import (
	"image"

	"gocv.io/x/gocv"
)

var (
	roiQuest      = Rect(17, 226, 5, 15)
	thQuestActive = gocv.NewScalar(200, 200, 200, 0)

	// Quest complete dialog
	roiQuestCompleteExp        = Rect(488, 318, 75, 50)
	roiQuestCompleteBtns       = Rect(365, 524, 550, 70)
	roiQuestCompleteTapToClose = Rect(482, 477, 320, 20)
	roiModalComplete           = Rect(380, 598, 530, 64)
	// roiQuestCompleteBtns = Rect(506, 524, 270, 70)

	roiActiveQuestIcon = Rect(55, 5, 35, 320)

	// Area exploration
	roiAreaExplorationTitle    = Rect(48, 5, 215, 50)
	roiAreaExplorationBtns     = Rect(770, 600, 450, 70)
	roiAreaExplorationNewQuest = Rect(196, 214, 30, 480)

	// Main menu
	// roiTopRigthMenuBtn = Rect(1222, 5, 50, 50)
	roiTopRigthMenuBtn = Rect(1227, 10, 40, 40)

	roiLeftMenu   = Rect(0, 10, 55, 315)
	ptActiveQuest = image.Pt(25, 235)

	roiLeftMenuDetail = Rect(56, 10, 255, 360)
	roiActiveQuest    = Rect(55, 190, 120, 65)

	// Sleep screen
	roiSleepModeLogo    = Rect(10, 70, 300, 100)
	ptSleepModeWakeFrom = image.Pt(650, 600)
	ptSleepModeWakeTo   = image.Pt(650, 400)

	// Dialog
	roiBtnBack        = Rect(17, 12, 34, 40)
	roiTxtAutoAndIcon = Rect(1009, 533, 71, 27)
	roiTxtAuto        = Rect(1035, 536, 45, 21)
	ptContinue        = image.Pt(320, 600)

	// Guard Journal
	roiTopLeft = Rect(0, 0, 300, 60)

	// Monster story
	roiMonsterStory = MonsterStoryROI{
		Buttons: Rect(970, 555, 230, 80),

		ModalStartStory:        Rect(580, 168, 120, 35),
		ModalStartStoryButtons: Rect(389, 484, 503, 70),
	}
)

type MonsterStoryROI struct {
	Buttons                image.Rectangle
	ModalStartStory        image.Rectangle
	ModalStartStoryButtons image.Rectangle
}
