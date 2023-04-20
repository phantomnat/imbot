package roi

import (
	"image"
)

var (
	// Quest complete dialog
	ROIQuestCompleteExp        = Rect(488, 318, 75, 50)
	ROIQuestCompleteBtns       = Rect(365, 524, 550, 70)
	ROIQuestCompleteTapToClose = Rect(482, 477, 320, 20)
	ROIModalComplete           = Rect(380, 598, 530, 64)
	// ROIQuestCompleteBtns = Rect(506, 524, 270, 70)

	ROIActiveQuestIcon = Rect(55, 5, 35, 320)

	// Area exploration
	ROIAreaExplorationTitle    = Rect(48, 5, 215, 50)
	ROIAreaExplorationBtns     = Rect(770, 600, 450, 70)
	ROIAreaExplorationNewQuest = Rect(196, 214, 30, 480)

	// Main menu
	// ROITopRigthMenuBtn = Rect(1222, 5, 50, 50)
	ROITopRigthMenuBtn = Rect(1227, 10, 40, 40)

	ROILeftMenu   = Rect(0, 10, 55, 315)
	PtActiveQuest = image.Pt(25, 235)

	ROILeftMenuDetail = Rect(56, 10, 255, 360)
	ROIActiveQuest    = Rect(55, 190, 120, 65)

	// Sleep screen
	ROISleepModeLogo    = Rect(10, 70, 300, 100)
	PtSleepModeWakeFrom = image.Pt(650, 600)
	PtSleepModeWakeTo   = image.Pt(650, 400)

	// Dialog
	ROIBtnBack        = Rect(17, 12, 34, 40)
	ROITxtAutoAndIcon = Rect(1009, 533, 71, 27)
	ROITxtAuto        = Rect(1035, 536, 45, 21)
	PtContinue        = image.Pt(320, 600)

	// Guard Journal
	ROITopLeft = Rect(0, 0, 300, 60)

	// Monster story
	ROIMonsterStory = MonsterStoryROI{
		Buttons: Rect(970, 555, 230, 80),

		ModalStartStory:        Rect(580, 168, 120, 35),
		ModalStartStoryButtons: Rect(389, 484, 503, 70),
	}

	// Menu
	PtMenu           = image.Pt(1240, 26)
	ROIOfficialForum = Rect(1110, 657, 156, 52)
)

type MonsterStoryROI struct {
	Buttons                image.Rectangle
	ModalStartStory        image.Rectangle
	ModalStartStoryButtons image.Rectangle
}
