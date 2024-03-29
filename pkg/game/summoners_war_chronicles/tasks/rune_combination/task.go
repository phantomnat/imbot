package rune_combination

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks"
)

type task struct {
	tasks.BaseTask

	setting *TaskSetting
	status  *TaskStatus

	ExitableStates map[domain.TaskState]struct{}
}

type TaskSetting struct {
	domain.TaskSettingBase

	Cooldown domain.Duration
	Steps    []RuneCombineStep
}

type RuneCombineStep struct {
	RuneSet []roi.RuneSet
	Stars   int
}

type TaskStatus struct {
	domain.TaskStatusBase

	StepIdx     int
	RuneLimit   map[roi.RuneSet]int
	RuneChoices []roi.RuneSet
	RuneCount   int
	RunesCount  map[roi.RuneSet]int

	CurrentRuneSet roi.RuneSet
	LastRuneSet    roi.RuneSet
	CurrentStars   int
	LastStars      int

	Stats []DailyStats
}

type DailyStats struct {
	Date          time.Time
	CombineCount  int
	FourStarRunes int
	FiveStarRunes int
	SixStarRunes  int
}

var _ domain.Task = (*task)(nil)

const (
	stateEnsureRuneAlchemy domain.TaskState = iota + 1
	stateGoToRune

	stateInitStep
	stateConfigStep
	stateApplyStep
	statePickRune
	stateCombineRune
	stateCheckResult
	stateNextStep

	stateGoToMainScreen
)

const (
	runeLimit = 3
)

func NewRuneCombination(manager domain.Manager, setting *TaskSetting) domain.Task {
	status := &TaskStatus{}
	t := &task{
		setting: setting,
		status:  status,
		BaseTask: tasks.NewBaseTask(
			manager, setting, status,
			map[domain.TaskState]string{
				stateGoToRune:          "go_to_rune",
				stateEnsureRuneAlchemy: "ensure_rune_alchemy",
				stateInitStep:          "init_step",
				stateConfigStep:        "config_step",
				stateApplyStep:         "apply_step",
				statePickRune:          "pick_rune",
				stateCombineRune:       "combine_rune",
				stateCheckResult:       "check_result",
				stateNextStep:          "next_step",
				stateGoToMainScreen:    "go_to_main_screen",
			},
		),
		ExitableStates: map[domain.TaskState]struct{}{
			stateInitStep: {},
			stateNextStep: {},
		},
	}
	return t
}

// func (t *task) LoadStatus(in any) {
// 	err := t.ConvertTo(in, &t.status)
// 	if err != nil {
// 		t.Log.Warnf("reset status, cannot the current: %+v", err)
// 		t.status = &TaskStatus{}
// 	} else {
// 		t.SetState(t.status.State)
// 	}
// }

var prefix = "rune"

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (t *task) handleExit() {
	if !t.Exiting {
		return
	}
	if _, exist := t.ExitableStates[t.State]; exist {
		t.SetState(stateGoToMainScreen)
	}
}

func (t *task) Reset() {
	t.SetState(domain.TaskStateBegin)
	t.reset()
}

func (t *task) reset() {
	t.status.StepIdx = 0
	t.status.CurrentRuneSet = ""
	t.status.CurrentStars = 0
	t.status.LastRuneSet = ""
	t.status.LastStars = 0
	t.status.RunesCount = make(map[roi.RuneSet]int)
	t.SaveStatus()
}

func (t *task) Do(m gocv.Mat) (triggered bool) {
	t.handleExit()

	triggered = true

	switch t.State {
	case domain.TaskStateBegin:

		t.status.StepIdx = 0
		// check reset
		t.status.Reset(func(today time.Time) {
			if len(t.status.Stats) == 0 || !today.Equal(t.status.Stats[0].Date) {
				t.status.Stats = append([]DailyStats{{Date: today}}, t.status.Stats...)
			}
			if len(t.status.Stats) > 30 {
				t.status.Stats = t.status.Stats[:30]
			}
		})
		t.SaveStatus()

		// TODO: go to rune alchemy from main screen
		t.SetState(stateGoToRune)

	case stateGoToRune:

		switch {
		case t.Manager.IsOnMainMenu(m):
			t.Manager.ClickPt(roi.MainMenu.PtRune)
			t.WaitMs(1000)

		case t.SearchROI(m,
			tasks.WithROI(roi.ROITopLeft),
			tasks.WithPath(prefix, "title_rune"),
			tasks.WithNoWait(),
			tasks.WithNextState(stateEnsureRuneAlchemy),
		):
			t.Manager.ClickPt(roi.Rune.PtRuneAlchemy)
			t.WaitMs(1000)

		case t.Manager.IsOnMainScreen(m):
			t.Manager.ClickPt(roi.PtTopRightMenu)
			t.WaitMs(1000)

		}

	case stateEnsureRuneAlchemy:
		invalidConfig := t.status.StepIdx >= len(t.setting.Steps)

		switch {
		case invalidConfig:
			// exit
			t.SetState(stateGoToMainScreen)
			return

		case t.SearchROI(m,
			tasks.WithROI(roi.ROITopLeft),
			tasks.WithPath(prefix, "title_rune_alchemy"),
			tasks.WithNoWait(),
			tasks.WithNextState(stateInitStep),
			tasks.WithDebugMatch(),
		):
			t.Manager.ClickPt(roi.RuneAlchemy.PtRuneCombination)
			t.WaitMs(1000)

			// deselect rune
			t.Manager.ClickPt(roi.RuneAlchemy.RuneCombination.PtDeselectAll)
			t.WaitMs(500)

			t.reset()
		}

	case stateNextStep:
		t.status.StepIdx++
		if t.status.StepIdx >= len(t.setting.Steps) {
			t.status.StepIdx = 0
			t.status.NextExecution = time.Now().Add(time.Duration(t.setting.Cooldown))
			t.SetState(stateGoToMainScreen)
		} else {
			// deselect rune
			t.Manager.ClickPt(roi.RuneAlchemy.RuneCombination.PtDeselectAll)
			t.WaitMs(500)

			t.SetState(stateInitStep)
		}
		t.SaveStatus()
		return

	case stateInitStep:
		invalidConfig := t.status.StepIdx >= len(t.setting.Steps)
		if invalidConfig {
			t.SetState(stateGoToMainScreen)
			return
		}
		stepSetting := t.setting.Steps[t.status.StepIdx]
		if stepSetting.Stars >= len(roi.RuneAlchemy.RuneCombination.PtRuneStars) {
			t.SetState(stateGoToMainScreen)
			return
		}

		// create rune choices
		// create rune limit
		rand.Seed(time.Now().Unix())

		t.status.CurrentStars = stepSetting.Stars
		t.status.RuneCount = 0
		t.status.RuneChoices = make([]roi.RuneSet, len(stepSetting.RuneSet))
		t.status.RuneLimit = make(map[roi.RuneSet]int, len(stepSetting.RuneSet))
		for i, runeSet := range stepSetting.RuneSet {
			t.status.RuneChoices[i] = runeSet
			t.status.RuneLimit[runeSet] = 1
		}

		switch len(t.status.RuneChoices) {
		case 1:
			t.status.RuneLimit[t.status.RuneChoices[0]] = 3
		case 2:
			rc := t.status.RuneChoices[rand.Intn(len(t.status.RuneChoices))]
			rs1 := t.status.RuneChoices[0]
			rs2 := t.status.RuneChoices[1]
			if t.status.RunesCount[rs1] > 0 && t.status.RunesCount[rs2] > 0 {
				rc = rs1
				if t.status.RunesCount[rs2] > t.status.RunesCount[rs1] {
					rc = rs2
				}
			}
			t.status.RuneLimit[rc] = 2
		}

		t.SetState(stateConfigStep)

	case stateConfigStep:
		switch {
		case len(t.status.RuneChoices) > 2:
			idx := rand.Intn(len(t.status.RuneChoices))
			t.status.CurrentRuneSet = t.status.RuneChoices[idx]
			t.status.RuneChoices = remove(t.status.RuneChoices, idx)

		case len(t.status.RuneChoices) == 2:
			idx := 0
			if t.status.LastRuneSet == t.status.RuneChoices[1] {
				idx = 1
			}
			t.status.CurrentRuneSet = t.status.RuneChoices[idx]
			t.status.RuneChoices = remove(t.status.RuneChoices, idx)

		case len(t.status.RuneChoices) == 0:
			// shouldn't reach this
			// next step
			t.SetState(stateNextStep)
			return

		default:
			t.status.CurrentRuneSet = t.status.RuneChoices[0]
			t.status.RuneChoices = nil
		}

		isConfigChanged := t.status.CurrentStars != t.status.LastStars ||
			t.status.CurrentRuneSet != t.status.LastRuneSet

		if !isConfigChanged {
			t.SetState(statePickRune)
			return
		}

		t.SetState(stateApplyStep)

	case stateApplyStep:
		switch {
		case t.SearchROI(m,
			tasks.WithROI(roi.RuneAlchemy.RuneCombination.SimpleSettingButtons),
			tasks.WithPath(prefix, "btn_simple_setting_apply"),
			tasks.WithNoWait(),
			// tasks.WithClick(),
			// tasks.WithNextState(statePickRune),
			// tasks.WithWaitMs(800),
		):
			t.Manager.ClickPt(roi.RuneAlchemy.RuneCombination.PtReset)
			t.WaitMs(500)

			t.Manager.ClickPt(roi.RuneAlchemy.RuneCombination.PtRuneStars[t.status.CurrentStars])
			t.WaitMs(500)

			// rune set
			t.Manager.ClickPt(roi.RuneAlchemy.RuneCombination.PtRuneSet[t.status.CurrentRuneSet])
			t.WaitMs(500)

			t.Manager.ClickPt(roi.RuneAlchemy.RuneCombination.PtSimpleSettingApply)
			t.WaitMs(1000)

			t.status.LastRuneSet = t.status.CurrentRuneSet
			t.status.LastStars = t.status.CurrentStars
			t.SetState(statePickRune)

		default:
			t.Manager.ClickPt(roi.RuneAlchemy.RuneCombination.PtSetting)
			t.WaitMs(800)
		}

	case statePickRune:
		mRuneList := m.Region(roi.RuneAlchemy.RuneCombination.RuneList)
		defer mRuneList.Close()
		stepSetting := t.setting.Steps[t.status.StepIdx]
		runeSet := t.status.CurrentRuneSet
		runes := t.Im.MatchMultiPoints(mRuneList, domain.MatchOption{
			Path: t.Manager.GetImagePath(prefix, "icon_rune_list_"+strconv.Itoa(stepSetting.Stars)+"_stars"),
		})

		t.status.RunesCount[runeSet] = len(runes)

		if len(runes) == 0 && time.Since(t.StateChangedAt).Seconds() < 2 {
			// retry for two seconds
			return
		} else if len(runes) < t.status.RuneLimit[runeSet] {
			// try new rune choice
			t.SetState(stateConfigStep)
			return
		}

		if t.status.RuneChoices == nil && len(runes) >= 3 {
			t.Manager.ClickPt(roi.RuneAlchemy.RuneCombination.PtSelectAll)
			t.WaitMs(500)
			t.status.RuneCount += t.status.RuneLimit[runeSet]
		} else {
			// TODO: ensure that we click the rune
			for i := 0; i < t.status.RuneLimit[runeSet]; i++ {
				pt := t.GetPtWithROI(roi.RuneAlchemy.RuneCombination.RuneList, runes[i])
				t.status.RuneCount++
				t.Manager.Click(pt.X+50, pt.Y+50)
				t.WaitMs(500)
			}
		}

		t.status.RuneLimit[runeSet] = 0

		t.SetState(stateCombineRune)

	case stateCombineRune:
		if t.status.RuneCount < runeLimit {
			t.SetState(stateConfigStep)
			return
		}

		switch {
		case t.SearchROI(m,
			tasks.WithROI(roi.RuneAlchemy.RuneCombination.RuneCombinationButtons),
			tasks.WithPath(prefix, "btn_combine_rune"),
			tasks.WithNextState(stateCheckResult),
			tasks.WithClick(),
		):
			t.Log.Infof("combining rune...")
			t.status.Stats[0].CombineCount++
			t.SaveStatus()

		case time.Since(t.StateChangedAt).Seconds() > 10:
			t.SetState(stateNextStep)
			return
		}

	case stateCheckResult:

		switch {
		case t.SearchROI(m,
			tasks.WithROI(roi.RuneAlchemy.RuneCombination.RuneCombinedButtons),
			tasks.WithPath(prefix, "btn_rune_combined_ok"),
			tasks.WithClick(),
			tasks.WithNextState(stateInitStep),
			tasks.WithWaitMs(500),
		):
			// TODO: check stars
			mRune := m.Region(roi.RuneAlchemy.RuneCombination.RuneCombinedRune)
			defer mRune.Close()

			is4Stars := t.SearchROI(mRune,
				tasks.WithPath(prefix, "icon_half_4_stars"),
			)
			is5Stars := t.SearchROI(mRune,
				tasks.WithPath(prefix, "icon_half_5_stars"),
			)

			if is4Stars {
				t.status.Stats[0].FourStarRunes++
			} else {
				if is5Stars {
					t.status.Stats[0].FiveStarRunes++
				}

				// save
				out := gocv.NewMatWithSize(mRune.Rows(), mRune.Cols(), gocv.MatTypeCV8UC4)
				gocv.CvtColor(mRune, &out, gocv.ColorBGRToRGBA)
				defer out.Close()
				today := strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "-")
				filePath := filepath.Join("cap", today+".png")
				gocv.IMWrite(filePath, out)
			}
			t.SaveStatus()

		case t.SearchROI(m,
			tasks.WithROI(roi.RuneAlchemy.RuneCombination.CheckRuneCombinationModal),
			tasks.WithPath(prefix, "btn_check_rune_combination_ok"),
			tasks.WithClick(),
		):

		}

	case stateGoToMainScreen:
		t.Manager.ClickPt(roi.PtTopRightHomeBtn)
		t.WaitMs(1000)
		t.SetState(domain.TaskStateEnd)
		t.SaveStatus()

	case domain.TaskStateEnd:
		t.Exiting = false
		t.SetState(domain.TaskStateBegin)
		t.SaveStatus()
		t.Exit()

	default:
		triggered = false
	}
	return
}
