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
	stateGoToBrawlArena
	stateStartPlay
	stateWaitForConfirm
	stateInTheLobby
	stateWaitBeforeExit
)

var (
	prefix = "arena"
)

type task struct {
	tasks.BaseTask

	setting *TaskSetting
	status  *TaskStatus
}

type TaskSetting struct {
	domain.TaskSettingBase

	Times int
}

type TaskStatus struct {
	domain.TaskStatusBase

	NextReset time.Time
	Stats     []DailyStats
}

type DailyStats struct {
	Date  time.Time
	Count int
}

var _ domain.Task = (*task)(nil)

func NewBrawlArena(index int, manager domain.Manager, setting TaskSetting) domain.Task {
	status := &TaskStatus{}

	t := &task{
		setting: &setting,
		status:  status,
		BaseTask: tasks.NewBaseTask(
			manager, &setting, status,
			map[domain.TaskState]string{
				stateGoToMainMenu:   "go_to_main_menu",
				stateGoToArena:      "go_to_arena",
				stateGoToBrawlArena: "go_to_brawl_arena",
				//stateStartPlay: "start"
				//stateWaitForConfirm
				//stateInTheLobby
				//stateWaitBeforeExit
			},
		),
	}
	return t
}

func (t *task) Do(m gocv.Mat) bool {
	if t.Exiting {
		// TODO: handle exit request
	}

	// checking event time
	// checking daily limit
	// go to arena
	// go to brawl arena
	// start challenge
	// hey this is vsocde , can you skip the scanning please
	// is it ok right nown

	switch t.State {
	case domain.TaskStateBegin:
		// load status
		if t.status.Reset(func(today time.Time) {
			if len(t.status.Stats) == 0 || !today.Equal(t.status.Stats[0].Date) {
				t.status.Stats = append([]DailyStats{{Date: today}}, t.status.Stats...)
			}
			if len(t.status.Stats) > 30 {
				t.status.Stats = t.status.Stats[:30]
			}
		}) {
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
			tasks.WithROI(roi.MainMenu.RightSide),
			tasks.WithPath("menu", "btn_arena"),
			tasks.WithNextState(stateGoToBrawlArena),
			tasks.WithClick(),
			// tasks.WithWaitMs(1000),
		) {
			return true
		}

	case stateGoToBrawlArena:
		if t.SearchROI(m,
			tasks.WithROI(roi.Arena.Title),
			tasks.WithPath("arena", "txt_brawl_arena"),
			tasks.WithNextState(domain.TaskStateEnd),
			tasks.WithClick(),
		) {
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

func (t *task) Reset() {
	t.SetState(domain.TaskStateBegin)
}

func (t *task) IsReady() bool {
	if !t.setting.IsEnabled() {
		return false
	}
	hour := time.Now().Hour()
	isOnPlayingTime := (hour >= 12 && hour < 13) || (hour >= 16 && hour < 17) || (hour >= 21 && hour < 22)
	isLimitReached := len(t.status.Stats) == 0 || (len(t.status.Stats) > 0 && t.status.Stats[0].Count >= t.setting.Times)
	return !isLimitReached && isOnPlayingTime && time.Now().After(t.status.GetNextExecution())
}
