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
	PtTopRightMenu = image.Pt(1240, 26)

	MainMenu = struct {
		OfficialForum image.Rectangle
		LeftSide      image.Rectangle
		RightSide     image.Rectangle

		PtRune image.Point
	}{
		OfficialForum: Rect(1110, 657, 156, 52),
		LeftSide:      Rect(11, 249, 378, 398),
		RightSide:     Rect(798, 211, 471, 428),

		PtRune: Pt(940, 390),
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

	Rune = struct {
		PtRuneAlchemy image.Point
	}{
		PtRuneAlchemy: Pt(930, 660),
	}
	RuneAlchemy = struct {
		PtRuneCombination image.Point
		RuneCombination   RuneCombination
	}{
		PtRuneCombination: Pt(100, 280),
		RuneCombination: RuneCombination{
			PtSimpleSetting:        Pt(1040, 105),
			RuneCombinationButtons: Rect(304, 616, 450, 75),
			PtDeselectAll:          Pt(400, 650),

			SimpleSetting:        Rect(543, 69, 320, 630),
			SimpleSettingButtons: Rect(543, 622, 320, 75),
			PtSimpleSettingApply: Pt(750, 660),

			RuneList: Rect(870, 137, 370, 460),

			CheckRuneCombinationModal: Rect(312, 152, 655, 415),
			RuneCombinedRune:          Rect(474, 254, 99, 102),
			RuneCombinedButtons:       Rect(387, 594, 507, 79),

			PtReset: Pt(600, 660),
			PtApply: Pt(760, 660),

			PtRuneSet: map[RuneSet]image.Point{
				EnergyRuneSet: Pt(600, 110),
				GuardRuneSet:  Pt(668, 110),
				BladeRuneSet:  Pt(736, 110),
				RageRuneSet:   Pt(804, 110),

				FatalRuneSet:  Pt(600, 170),
				SwiftRuneSet:  Pt(668, 170),
				FocusRuneSet:  Pt(736, 170),
				EndureRuneSet: Pt(804, 170),

				ForesightRuneSet: Pt(600, 225),
				AssembleRuneSet:  Pt(668, 225),
				DespairRuneSet:   Pt(736, 225),
				VampireRuneSet:   Pt(804, 225),
			},

			PtRuneSlots: [7]image.Point{
				Pt(0, 0), // 0 - no used
				Pt(610, 340),
				Pt(700, 340),
				Pt(790, 340),
				Pt(610, 390),
				Pt(700, 390),
				Pt(790, 390),
			},

			PtRuneStars: [7]image.Point{
				Pt(0, 0), // 0 - no used
				Pt(610, 450),
				Pt(700, 450),
				Pt(790, 450),
				Pt(610, 510),
				Pt(700, 510),
				Pt(790, 510),
			},
		},
	}
)

var Pt = image.Pt

type MonsterStoryROI struct {
	Buttons                image.Rectangle
	ModalStartStory        image.Rectangle
	ModalStartStoryButtons image.Rectangle
}

type RuneCombination struct {
	PtSimpleSetting        image.Point
	RuneCombinationButtons image.Rectangle
	PtDeselectAll          image.Point

	SimpleSetting        image.Rectangle
	SimpleSettingButtons image.Rectangle
	PtSimpleSettingApply image.Point

	RuneList image.Rectangle

	CheckRuneCombinationModal image.Rectangle
	RuneCombinedRune          image.Rectangle
	RuneCombinedButtons       image.Rectangle

	PtReset image.Point
	PtApply image.Point

	PtRuneSet map[RuneSet]image.Point

	// slot 1-6
	PtRuneSlots [7]image.Point

	// star 1-6
	PtRuneStars [7]image.Point
}

type RuneSet string

const (
	EnergyRuneSet RuneSet = "energy"
	GuardRuneSet  RuneSet = "guard"
	BladeRuneSet  RuneSet = "blade"
	RageRuneSet   RuneSet = "rage"

	FatalRuneSet  RuneSet = "fatal"
	SwiftRuneSet  RuneSet = "swift"
	FocusRuneSet  RuneSet = "focus"
	EndureRuneSet RuneSet = "endure"

	ForesightRuneSet RuneSet = "foresight"
	AssembleRuneSet  RuneSet = "assemble"
	DespairRuneSet   RuneSet = "despair"
	VampireRuneSet   RuneSet = "vampire"
)
