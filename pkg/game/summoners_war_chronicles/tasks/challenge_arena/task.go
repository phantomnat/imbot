package challenge_arena

import (
	"time"

	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks"
)

const (
	stateGoToMainMenu domain.TaskState = iota + 1
	stateGoToArena
	stateGoToChallengeArena
	stateChallenge
	stateDoQuickBattle
)

var (
	resetCoolDown = 30 * time.Minute
)

type task struct {
	tasks.BaseTask
	status TaskStatus
}

type TaskStatus struct {
	Name         string
	LastExecuted time.Time `json:"lastExecuted"`
	NextExecuted time.Time `json:"nextExecuted"`
}

var _ domain.Task = (*task)(nil)

func NewChallengeArenaTask(index int, name string, manager domain.Manager) domain.Task {
	return &task{
		BaseTask: tasks.BaseTask{
			Index:   index,
			Name:    name,
			Manager: manager,
			Log:     zap.S().Named("task").Named("challenge-arena"),
			StateTexts: map[domain.TaskState]string{
				domain.TaskStateBegin:   "begin",
				stateGoToMainMenu:       "",
				stateGoToArena:          "",
				stateGoToChallengeArena: "",
				stateChallenge:          "",
				stateDoQuickBattle:      "",
				domain.TaskStateEnd:     "end",
				// TODO: readable task text
			},
		},
	}
}

func (t *task) Do(m gocv.Mat) bool {
	if t.Exiting {
		// TODO: handle exit request
	}

	switch t.State {
	case domain.TaskStateBegin:
		// load status
		t.LoadStatus(&t.status)
		if t.status.Name == "" {
			t.status.Name = t.GetName()
		}
		t.SetState(stateGoToMainMenu)
	case stateGoToMainMenu:
		// get status for this task
		// check if ready to do
		inMainMenu, _ := t.Manager.MatchInROI(m, roi.ROIOfficialForum, domain.MatchOption{
			Path: "menu.official_forum",
		})
		if inMainMenu {
			t.SetState(stateGoToArena)
		} else {
			// t.Manager.ClickPt(roi.PtMenu)
			time.Sleep(1000 * time.Millisecond)

			// TODO: remove
			t.SetState(stateGoToArena)
		}
	case stateGoToArena:

		time.Sleep(1000 * time.Millisecond)
		t.SetState(stateGoToChallengeArena)

	case stateGoToChallengeArena:
		time.Sleep(1000 * time.Millisecond)
		t.SetState(stateChallenge)
	case stateChallenge:
		// find challenge btn
		// pass the name to next state
		// TODO: need multiple points detection here
		// skip if name already done
		// if all name skipped, scroll down
		t.status.NextExecuted = time.Now().Add(10 * time.Second)
		t.Log.Infof("next execute at %v", t.status.NextExecuted.Format(time.RFC3339))

		t.SaveStatus(t.status)

		time.Sleep(1000 * time.Millisecond)
		t.SetState(domain.TaskStateEnd)
	case stateDoQuickBattle:
		// looks for quick battle
		// remember the name
		// TODO: check if out of battle
		// check for the limit
	case domain.TaskStateEnd:
		t.Exit()
	}
	return false
}

func (t *task) IsReady() bool {
	if t.status.NextExecuted.IsZero() {
		return true
	}
	return time.Now().After(t.status.NextExecuted)
}
