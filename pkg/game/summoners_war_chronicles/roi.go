package summonerswar

import (
	"gocv.io/x/gocv"
)

var (
	roiQuest      = Rect(17, 226, 5, 15)
	thQuestActive = gocv.NewScalar(200, 200, 200, 0)

	roiQuestCompleteExp = Rect(488, 318, 75, 50)

	roiQuestCompleteBtns = Rect(365, 524, 550, 70)
	// roiQuestCompleteBtns = Rect(506, 524, 270, 70)

	roiSleepModeLogo = Rect(10, 70, 300, 100)

	roiActiveQuestIcon = Rect(55, 5, 35, 320)

	roiActiveQuest = Rect(55, 190, 120, 65)

	// Area exploration
	roiAreaExplorationTitle = Rect(48, 5, 215, 50)
	roiAreaExplorationBtns = Rect(770, 600, 450, 70)
)
