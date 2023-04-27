package challenge_arena

import (
	"image"
	"time"

	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks"
)

const (
	stateGoToMainMenu domain.TaskState = iota + 1
	stateGoToArena
	stateGoToChallengeArena
	stateScrollToTop
	stateRefreshList
	stateSearchForChallenge
	stateChallenge
	stateDoQuickBattle
	stateWaitForQuickBattle
)

var (
	resetCoolDown        = 30 * time.Minute
	prefixArena          = "arena"
	prefixChallengeArena = "arena.challenge"
	maxRepeat            = 3
)

type task struct {
	tasks.BaseTask
	setting TaskSetting
	status  TaskStatus
}

type TaskSetting struct {
	Enable bool
}

type TaskStatus struct {
	Name         string
	State        domain.TaskState
	LastExecuted time.Time
	NextExecuted time.Time
	NextReset    time.Time
	Repeat       int
	Points       []image.Point
	PointIdx     int
	Stats        []DailyStats
}

type DailyStats struct {
	Date        time.Time
	Refresh     int
	QuickBattle int
}

var _ domain.Task = (*task)(nil)

func NewChallengeArenaTask(index int, manager domain.Manager, setting TaskSetting) domain.Task {
	t := &task{
		setting: setting,
		BaseTask: tasks.NewBaseTask(index, manager, setting,
			map[domain.TaskState]string{
				stateGoToMainMenu:       "go_to_main_menu",
				stateGoToArena:          "go_to_arena",
				stateGoToChallengeArena: "go_to_challenge_arena",
				stateScrollToTop:        "scroll_to_top",
				stateRefreshList:        "refresh_list",
				stateSearchForChallenge: "search_for_challenge",
				stateChallenge:          "challenge",
				stateDoQuickBattle:      "do_quick_battle",
				stateWaitForQuickBattle: "wait_for_quick_battle",
			},
		),
	}
	return t
}

func (t *task) LoadStatus(in any) {
	err := t.ConvertTo(in, &t.status)
	if err != nil {
		t.Log.Warnf("reset status, cannot the current: %+v", err)
		t.status = TaskStatus{}
	} else {
		t.SetState(t.status.State)
	}
}

func (t *task) GetStatus() any {
	return t.status
}

func (t *task) Do(m gocv.Mat) bool {
	if t.Exiting {
		// TODO: handle exit request
	}

	switch t.State {
	case domain.TaskStateBegin:
		// load status
		if t.status.Name == "" {
			t.status.Name = t.GetName()
		}

		// check reset
		if time.Now().After(t.status.NextReset) {
			// reset
			today := time.Now().Truncate(time.Hour)
			t.status.Stats = append([]DailyStats{{Date: today}}, t.status.Stats...)
			if len(t.status.Stats) > 30 {
				t.status.Stats = t.status.Stats[:30]
			}
			t.status.NextReset = time.Now().Truncate(time.Hour).AddDate(0, 0, 1)
			t.SaveStatus()
		}

		switch {
		case t.SearchROI(m,
			tasks.WithROI(roi.MainMenu.OfficialForum),
			tasks.WithPath("menu", "official_forum"),
			tasks.WithNextState(stateGoToArena),
			tasks.WithNoWait(),
			// tasks.WithDebugMatch(),
		):
			return true
		case t.Manager.GoToMainScreen(m):
			t.SetState(stateGoToMainMenu)
			return true
		}

	case stateGoToMainMenu:
		// get status for this task
		// check if ready to do
		if t.SearchROI(m,
			tasks.WithROI(roi.MainMenu.OfficialForum),
			tasks.WithPath("menu", "official_forum"),
			tasks.WithNextState(stateGoToArena),
			tasks.WithNoWait(),
		) {
			return true
		}

		t.Manager.ClickPt(roi.PtTopRightMenu)
		t.WaitMs(1000)

	case stateGoToArena:
		if t.SearchROI(m,
			tasks.WithROI(roi.MainMenu.LeftSide),
			tasks.WithPath("menu", "btn_arena"),
			tasks.WithNextState(stateGoToChallengeArena),
			tasks.WithClick(),
			// tasks.WithWaitMs(1000),
		) {
			return true
		}

	case stateGoToChallengeArena:
		if t.SearchROI(m,
			tasks.WithROI(roi.Arena.Title),
			tasks.WithPath("arena", "txt_challenge_arena"),
			tasks.WithNextState(stateScrollToTop),
			tasks.WithClick(),
			// tasks.WithWaitMs(1000),
		) {
			return true
		}

	case stateScrollToTop:
		if !t.SearchROI(m,
			tasks.WithPath(prefixChallengeArena, "btn_refresh"),
			tasks.WithROI(roi.ChallengeArena.RefreshBtn),
			tasks.WithNoWait(),
		) {
			return false
		}

		// scroll to top
		t.Manager.DragDuration(roi.ChallengeArena.PtStopDrag, roi.ChallengeArena.PtStartDrag, 1000)
		// t.WaitMs(1000)
		t.Manager.DragDuration(roi.ChallengeArena.PtStopDrag, roi.ChallengeArena.PtStartDrag, 1000)
		t.WaitMs(1000)

		t.SetState(stateRefreshList)

	case stateRefreshList:
		// refresh
		mRefreshListDialog := m.Region(roi.ChallengeArena.RefreshDialog)
		defer mRefreshListDialog.Close()

		switch {
		case t.SearchROI(mRefreshListDialog,
			tasks.WithPath(prefixChallengeArena, "txt_refresh_list"),
			tasks.WithNoWait(),
		):
			// click refresh list or close dialog
			switch {
			case t.SearchROI(mRefreshListDialog,
				tasks.WithPath(prefixChallengeArena, "btn_refresh_list"),
				tasks.WithNextState(stateSearchForChallenge),
				tasks.WithClick(),
				tasks.WithClickOffset(roi.ChallengeArena.RefreshDialog),
			):
				// refresh
				t.status.NextExecuted = time.Now().Add(resetCoolDown)
				t.status.Repeat = 0
				t.status.Stats[0].Refresh++
				t.Log.Infof("next execute at %v", t.status.NextExecuted.Format(time.RFC3339))
				t.SaveStatus()

			case t.SearchROI(
				mRefreshListDialog,
				tasks.WithPath(prefixChallengeArena, "btn_refresh_list_wait"),
				tasks.WithNextState(stateSearchForChallenge),
				tasks.WithNoWait(),
			):
				// TODO: refresh with crystal
				// skip during development
				t.status.NextExecuted = time.Now().Add(resetCoolDown)
				t.status.Repeat = 0
				t.SaveStatus()
				t.Log.Infof("skip during development")
				t.Manager.ClickPt(roi.ChallengeArena.PtCloseRefreshListDialog)
			}

			t.WaitMs(2000)
			return true

		case t.SearchROI(m,
			tasks.WithPath(prefixChallengeArena, "btn_refresh"),
			tasks.WithROI(roi.ChallengeArena.RefreshBtn),
			tasks.WithClick(),
			tasks.WithWaitMs(2000),
		):
			// click refresh
		}

	case stateSearchForChallenge:
		if !t.SearchROI(m,
			tasks.WithPath(prefixChallengeArena, "btn_refresh"),
			tasks.WithROI(roi.ChallengeArena.RefreshBtn),
			tasks.WithNoWait(),
		) {
			return false
		}

		mChooseOpponent := m.Region(roi.ChallengeArena.ChooseOpponent)
		defer mChooseOpponent.Close()

		points := t.Im.MatchMultiPoints(mChooseOpponent, domain.MatchOption{
			Path:     t.Manager.GetImagePath(prefixChallengeArena, "btn_challenge"),
			PrintVal: true,
			Th:       0.03,
		})

		// t.Log.Infof("%#v", points)

		t.status.Points = points
		t.status.PointIdx = 0
		t.SaveStatus()

		t.SetState(stateChallenge)

	case stateChallenge:
		if t.status.PointIdx >= len(t.status.Points) {
			t.Manager.DragDuration(roi.ChallengeArena.PtStartDrag, roi.ChallengeArena.PtStopDrag, 1000)
			t.WaitMs(2000)
			t.SetState(stateSearchForChallenge)
			t.status.Repeat++
			if t.status.Repeat >= maxRepeat {
				// TODO: exit
				t.SetState(domain.TaskStateEnd)
			}
			t.SaveStatus()
			return true
		}

		if t.SearchROI(m,
			tasks.WithPath(prefixChallengeArena, "btn_refresh"),
			tasks.WithROI(roi.ChallengeArena.RefreshBtn),
			tasks.WithNoWait(),
		) {
			t.Manager.ClickPt(t.GetPtWithROI(roi.ChallengeArena.ChooseOpponent, t.status.Points[t.status.PointIdx]))
			t.WaitMs(2000)
			t.SetState(stateDoQuickBattle)
			return true
		}

	case stateDoQuickBattle:
		if t.SearchROI(m,
			tasks.WithPath(prefixChallengeArena, "btn_quick_battle"),
			tasks.WithROI(roi.ChallengeArena.CharSelectionBattleBtns),
			tasks.WithNextState(stateWaitForQuickBattle),
			tasks.WithClick(),
			tasks.WithWaitMs(2000),
		) {
			// do quick battle
			t.status.Stats[0].QuickBattle++
			t.SaveStatus()

		} else if t.SearchROI(m,
			tasks.WithPath(prefixChallengeArena, "btn_quick_battle_disable"),
			tasks.WithROI(roi.ChallengeArena.CharSelectionBattleBtns),
			tasks.WithNoWait(),
		) {
			// back and try again
			if t.SearchROI(m,
				tasks.WithPath(prefixChallengeArena, "btn_battle_start"),
				tasks.WithROI(roi.ChallengeArena.CharSelectionBattleBtns),
				tasks.WithNextState(stateChallenge),
				tasks.WithNoWait(),
			) {
				// back
				t.Manager.ClickPt(roi.PtTopLeftBackBtn)
				t.WaitMs(1000)
				t.status.PointIdx++
				t.SaveStatus()
			}
		}

	case stateWaitForQuickBattle:
		// looks for quick battle
		// remember the name
		// TODO: check if out of battle
		// check for the limit

		if t.SearchROI(m,
			tasks.WithPath(prefixChallengeArena, "icon_arena_coin"),
			tasks.WithROI(roi.ChallengeArena.VictoryReward),
			tasks.WithNextState(stateSearchForChallenge),
			tasks.WithNoWait(),
		) {
			t.Manager.ClickPt(roi.ChallengeArena.PtVictoryOKBtn)
			t.WaitMs(2000)
		}

	case domain.TaskStateEnd:
		t.Manager.ClickPt(roi.PtTopRightHomeBtn)
		t.WaitMs(3000)
		t.Exit()
		return true
	}

	return false
}

func (t *task) IsReady() bool {
	if !t.setting.Enable {
		return false
	}
	return t.status.NextExecuted.IsZero() || time.Now().After(t.status.NextExecuted)
}

func (t *task) SaveStatus() {
	t.status.State = t.State
	t.Manager.SaveStatusByIndex(t.Index, t.Name, t.status)
}
