package fishing

import (
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks"
)

type task struct {
	tasks.BaseTask
	setting *TaskSetting
}

type TaskSetting struct {
	domain.TaskSettingBase
}

var _ domain.Task = (*task)(nil)

func NewFishing(manager domain.Manager, setting *TaskSetting) domain.Task {
	status := &domain.TaskStatusBase{}
	t := &task{
		setting: setting,
		BaseTask: tasks.NewBaseTask(manager, setting, status,
			map[domain.TaskState]string{},
		),
	}
	return t
}

func (t *task) Do(m gocv.Mat) bool {
	if t.Exiting {
		return true
	}

	switch t.State {
	case domain.TaskStateBegin:
		switch {
		case t.SearchROI(m,
			tasks.WithROI(roi.Fishing.Icons),
			tasks.WithPath("fishing", "icon_fish"),
			tasks.WithWaitMs(500),
			tasks.WithClick(),
			// tasks.WithDebugMatch(),
		):
		default:
			t.WaitMs(50)
		}

		// t.Manager.Click(470, 560)
		// t.WaitMs(800)
		// t.Manager.Click(585, 560)
		// t.WaitMs(800)
		// t.Manager.Click(700, 560)
		// t.WaitMs(800)

	case domain.TaskStateEnd:
		t.Exit()
	}

	return false
}
