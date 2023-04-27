package challenge_arena

import (
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
	resetCoolDown = 30 * time.Minute
	prefix        = "arena"
	maxRepeat     = 3
)

type task struct {
	tasks.BaseTask
	setting TaskSetting
	status  TaskStatus
}

type TaskSetting struct {
	Enable bool
	Times  int
}

type TaskStatus struct {
	Name         string
	NextExecuted time.Time
	NextReset    time.Time
	Stats        []DailyStats
}

type DailyStats struct {
	Date  time.Time
	Count int
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

	case domain.TaskStateEnd:
		t.Manager.ClickPt(roi.PtTopRightHomeBtn)
		t.WaitMs(3000)
		t.Exit()
		return true
	}

	return false
}

func (t *task) IsReady() bool {
	if t.status.NextExecuted.IsZero() {
		return true
	}
	return time.Now().After(t.status.NextExecuted)
}

func (t *task) SaveStatus() {
	t.Manager.SaveStatusByIndex(t.Index, t.Name, t.status)
}
