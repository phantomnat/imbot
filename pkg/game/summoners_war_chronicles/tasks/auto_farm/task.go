package auto_farm

import (
	"fmt"

	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks"
)

type task struct {
	tasks.BaseTask
	setting TaskSetting
	// status  TaskStatus
}

type TaskSetting struct {
	Enable bool
}

// type TaskStatus struct {
// 	Name string
// }

var _ domain.Task = (*task)(nil)

const (
	stateGoToMainMenu domain.TaskState = iota + 1
	stateEnsureAutoBattle
)

func NewAutoFarm(index int, manager domain.Manager, setting TaskSetting) domain.Task {
	t := &task{
		setting: setting,
		BaseTask: tasks.BaseTask{
			Im:      manager.GetImageManager(),
			Index:   index,
			Name:    fmt.Sprintf("%T", setting),
			Manager: manager,
			Log:     zap.S().Named("task").Named("challenge-arena"),
			StateTexts: map[domain.TaskState]string{
				domain.TaskStateBegin: "begin",
				stateGoToMainMenu:     "go_to_main_menu",
				stateEnsureAutoBattle: "ensure_auto_battle",
				domain.TaskStateEnd:   "end",
			},
		},
	}
	return t
}

func (t *task) Do(m gocv.Mat) bool {
	if t.Exiting {
		return true
	}

	switch t.State {
	case domain.TaskStateBegin:
		t.SetState(stateGoToMainMenu)
	case stateGoToMainMenu:
		if t.Manager.GoToMainScreen(m) {
			t.SetState(stateEnsureAutoBattle)
		}
	case stateEnsureAutoBattle:
		if t.SearchROI(m,
			tasks.WithPath("icon_auto_battle"),
			tasks.WithROI(roi.ROIMainScreen.AutoBattleIcon),
			tasks.WithClick(),
			tasks.WithDebugMatch(),
			tasks.WithNextState(domain.TaskStateEnd),
		) {
			return true
		}
		t.SetState(domain.TaskStateEnd)
	case domain.TaskStateEnd:
		t.Exit()
	}

	return false
}

func (t *task) IsReady() bool {
	return t.setting.Enable
}
