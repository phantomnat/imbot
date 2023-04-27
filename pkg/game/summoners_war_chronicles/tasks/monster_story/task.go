package monster_story

import (
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks"
)

type task struct {
	tasks.BaseTask
	setting *TaskSetting
	status  *TaskStatus
}

type TaskSetting struct {
	domain.TaskSettingBase
}
type TaskStatus struct {
	domain.TaskStatusBase
}

var _ domain.Task = (*task)(nil)

const (
	stateFindAndAcceptQuest domain.TaskState = iota + 1
)

func NewMonsterStory(index int, manager domain.Manager, setting TaskSetting) domain.Task {
	status := &TaskStatus{}

	t := &task{
		setting: &setting,
		status:  status,
		BaseTask: tasks.NewBaseTask(
			manager, &setting, status,
			map[domain.TaskState]string{
				stateFindAndAcceptQuest: "find_and_accept_quest",
			},
		),
	}
	return t
}

var prefix = "guard_journal"

func (t *task) Do(m gocv.Mat) bool {
	switch t.State {
	case domain.TaskStateBegin:
		switch {
		case t.Manager.IsOnMainScreen(m):
			mLeftMenuDetail := m.Region(roi.ROILeftMenuDetail)
			defer mLeftMenuDetail.Close()

			switch {
			case t.SearchROI(m,
				tasks.WithROI(roi.ROILeftMenuDetail),
				tasks.WithPath(prefix, "icon_ongoing_quest"),
				tasks.WithNoWait(),
			):
				t.Log.Infof("waiting for the quest to finish...")
				t.WaitMs(500)
				return true

			case t.SearchROI(m,
				tasks.WithROI(roi.ROILeftMenuDetail),
				tasks.WithPath(prefix, "icon_finish_quest"),
				tasks.WithClick(),
			):
				t.Log.Infof("quest finish")
				return true

			default:
				// click active quest
				t.Manager.ClickPt(roi.PtActiveQuest)
				t.WaitMs(1000)
			}

		case t.Manager.IsOnMainMenu(m):
			t.Manager.ClickPt(roi.PtTopRightHomeBtn)
			t.WaitMs(1000)

		case t.SearchROI(m,
			tasks.WithPath(prefix, "txt_guard_journal"),
			tasks.WithROI(roi.ROITopLeft),
			tasks.WithNoWait(),
			tasks.WithNextState(stateFindAndAcceptQuest),
		):
			// find and accept quest

		case t.Manager.HandleConversationDialog(m):
		case t.Manager.HandleQuestCompleted(m):
		}

	case stateFindAndAcceptQuest:
		switch {
		case t.Manager.IsOnMainScreen(m):
			t.SetState(domain.TaskStateBegin)

		case t.Manager.IsOnMainMenu(m):
			t.Manager.ClickPt(roi.PtTopRightHomeBtn)
			t.WaitMs(1000)

		case t.SearchROI(m,
			tasks.WithPath(prefix, "btn_next_stage"),
			tasks.WithROI(roi.ROIMonsterStory.Buttons),
			tasks.WithClick(),
		):

		case t.SearchROI(m,
			tasks.WithPath(prefix, "btn_claim"),
			tasks.WithROI(roi.ROIMonsterStory.Buttons),
			tasks.WithClick(),
		):

		case t.SearchROI(m,
			tasks.WithPath(prefix, "btn_search"),
			tasks.WithROI(roi.ROIMonsterStory.Buttons),
			tasks.WithClick(),
		):

		case t.SearchROI(m,
			tasks.WithPath(prefix, "txt_start_story"),
			tasks.WithROI(roi.ROIMonsterStory.ModalStartStory),
			tasks.WithNoWait(),
		):
			if t.SearchROI(m,
				tasks.WithPath(prefix, "btn_ok"),
				tasks.WithROI(roi.ROIMonsterStory.ModalStartStoryButtons),
				tasks.WithClick(),
			) {

			}

		case t.Manager.HandleQuestCompleted(m):

		default:
			t.WaitMs(500)
		}

	case domain.TaskStateEnd:
		t.Exit()
	}

	return false
}

func (t *task) IsReady() bool {
	return t.setting.Enable
}
