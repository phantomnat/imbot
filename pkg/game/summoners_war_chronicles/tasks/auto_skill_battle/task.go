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
		mMonster1Skill1 := m.Region(roi.MainScreen.Monster1Skills[0])
		mMonster1Skill2 := m.Region(roi.MainScreen.Monster1Skills[1])
		mMonster2Skill1 := m.Region(roi.MainScreen.Monster2Skills[0])
		mMonster2Skill2 := m.Region(roi.MainScreen.Monster2Skills[1])
		defer mMonster1Skill1.Close()
		defer mMonster1Skill2.Close()
		defer mMonster2Skill1.Close()
		defer mMonster2Skill2.Close()

		if ok, tpl := t.Im.Get("icon_monster_skill_mask"); !ok {
			colorMonster1Skill1 := mMonster1Skill1.MeanWithMask(*tpl)

			t.Log.Infof("monster 1 skill 1 (bgr: %.1f,%.1f,%.1f)",
				colorMonster1Skill1.Val1,
				colorMonster1Skill1.Val2,
				colorMonster1Skill1.Val3,
			)
		}

		// 	opt.Mask = tpl
		// 	opt.HasMask = true
		// }
		//

	case domain.TaskStateEnd:
		t.Exit()
	}

	return false
}
