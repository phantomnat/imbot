package roi

import (
	"image"
)

var (
	PtTopLeftBackBtn  = Pt(35, 35)
	PtTopRightHomeBtn = Pt(1235, 30)

	// Quest complete dialog
	ROIQuestCompleteExp        = Rect(488, 318, 75, 50)
	ROIQuestCompleteBtns       = Rect(365, 524, 550, 70)
	ROIQuestCompleteTapToClose = Rect(482, 477, 320, 20)
	ROIModalComplete           = Rect(380, 598, 530, 64)
	// ROIQuestCompleteBtns = Rect(506, 524, 270, 70)

	QuestCompleted = struct {
		Buttons    image.Rectangle
		TapToClose image.Rectangle

		PtOK image.Point
	}{
		Buttons:    Rect(365, 524, 550, 70),
		TapToClose: Rect(482, 477, 320, 20),

		PtOK: Pt(670, 560),
	}

	ROIActiveQuestIcon = Rect(55, 5, 35, 320)

	// Area exploration
	AreaExploration = struct {
		Title   image.Rectangle
		Buttons image.Rectangle

		QuestList            image.Rectangle
		PtStartDragQuestList image.Point
		PtStopDragQuestList  image.Point
	}{
		Title:   Rect(48, 5, 215, 50),
		Buttons: Rect(770, 600, 450, 70),

		QuestList:            Rect(196, 214, 30, 480),
		PtStartDragQuestList: Pt(110, 640),
		PtStopDragQuestList:  Pt(110, 240),
	}

	ROIAreaExplorationTitle    = Rect(48, 5, 215, 50)
	ROIAreaExplorationBtns     = Rect(770, 600, 450, 70)
	ROIAreaExplorationNewQuest = Rect(196, 214, 30, 480)

	// Main screen
	ROIMainScreen = struct {
		CoinIcon    image.Rectangle
		CrystalIcon image.Rectangle

		AutoBattleIcon image.Rectangle
		PtAutoBattle   image.Point
	}{
		CoinIcon:    Rect(517, 15, 37, 38),
		CrystalIcon: Rect(646, 13, 36, 41),

		AutoBattleIcon: Rect(925, 654, 47, 48),
		PtAutoBattle:   Pt(950, 680),
	}

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
	PtMenu = image.Pt(1240, 26)

	MainMenu = struct {
		OfficialForum image.Rectangle
		LeftSide      image.Rectangle
		RightSide     image.Rectangle
	}{
		OfficialForum: Rect(1110, 657, 156, 52),
		LeftSide:      Rect(11, 249, 378, 398),
		RightSide:     Rect(798, 211, 471, 428),
	}

	// Arena
	Arena = struct {
		Title image.Rectangle
	}{
		Title: Rect(210, 84, 860, 98),
	}

	ChallengeArena = struct {
		ChooseOpponent image.Rectangle
		RefreshBtn     image.Rectangle

		PtStartDrag image.Point
		PtStopDrag  image.Point

		// refresh list dialog
		RefreshDialog            image.Rectangle
		PtCloseRefreshListDialog image.Point

		// CharSelectionBattle
		CharSelectionBattleBtns image.Rectangle

		// Victory Dialog
		VictoryReward  image.Rectangle
		PtVictoryOKBtn image.Point
	}{
		ChooseOpponent: Rect(510, 182, 736, 535),

		PtStartDrag: Pt(543, 702),
		PtStopDrag:  Pt(543, 200),

		RefreshBtn: Rect(1088, 129, 164, 41),

		RefreshDialog:            Rect(344, 173, 592, 373),
		PtCloseRefreshListDialog: Pt(900, 208),

		CharSelectionBattleBtns: Rect(879, 584, 401, 136),

		VictoryReward:  Rect(522, 368, 235, 133),
		PtVictoryOKBtn: Pt(640, 590),
	}
)

var Pt = image.Pt

type MonsterStoryROI struct {
	Buttons                image.Rectangle
	ModalStartStory        image.Rectangle
	ModalStartStoryButtons image.Rectangle
}
