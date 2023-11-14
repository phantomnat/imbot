package area_exploration

import (
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks"
)

type task struct {
	tasks.BaseTask
	setting *TaskSetting
	// status  TaskStatus
}

type TaskSetting struct {
	domain.TaskSettingBase
}

var _ domain.Task = (*task)(nil)

const (
	stateFindAndAcceptQuest domain.TaskState = iota + 1
)

func NewAreaExploration(index int, manager domain.Manager, setting TaskSetting) domain.Task {
	status := &domain.TaskStatusBase{}
	t := &task{
		setting: &setting,
		BaseTask: tasks.NewBaseTask(manager, &setting, status,
			map[domain.TaskState]string{
				stateFindAndAcceptQuest: "find_and_accept_quest",
			},
		),
	}
	return t
}

var prefix = "area_exploration"

func (t *task) Do(m gocv.Mat) bool {
	if t.Exiting {
		return true
	}

	switch t.State {
	case domain.TaskStateBegin:
		switch {
		case t.Manager.IsOnMainScreen(m):
			//mLeftMenuDetail := m.Region(roi.ROILeftMenuDetail)
			//defer mLeftMenuDetail.Close()
			//
			//foundIcon, ptIcon := t.Im.MatchPoint(mLeftMenuDetail, domain.MatchOption{
			//	Path:     t.Manager.GetImagePath("icon_area_exploration_quest"),
			//	PrintVal: true,
			//})
			//if !foundIcon {
			//	t.Manager.ClickPt(roi.PtActiveQuest)
			//	t.WaitMs(1000)
			//	return true
			//}
			//
			//x := roi.ROILeftMenuDetail.Min.X + ptIcon.X
			//y := roi.ROILeftMenuDetail.Min.Y + ptIcon.Y
			//roiQuestText := image.Rect(x-5, y, x+245, y+45)
			//if t.SearchROI(m,
			//	tasks.WithROI(roiQuestText),
			//	tasks.WithPath("exploration_quest"),
			//	tasks.WithClick(),
			//) {
			//
			//}

		case t.Manager.IsOnMainMenu(m):
			t.Manager.ClickPt(roi.PtTopRightHomeBtn)
			t.WaitMs(1000)

		case t.SearchROI(m,
			tasks.WithPath(prefix, "title_area_exploration"),
			tasks.WithROI(roi.AreaExploration.Title),
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
			tasks.WithNextState(domain.TaskStateBegin)
		case t.Manager.IsOnMainMenu(m):
			t.Manager.ClickPt(roi.PtTopRightHomeBtn)
			t.WaitMs(1000)

		case t.SearchROI(m,
			tasks.WithPath(prefix, "btn_accept"),
			tasks.WithROI(roi.AreaExploration.Buttons),
			tasks.WithClick(),
			// tasks.WithWaitMs(500),
			tasks.WithNextState(domain.TaskStateBegin),
			tasks.WithDebugMatch(),
		):
			t.Log.Infof("accept exploration quest")
		case t.SearchROI(m,
			tasks.WithPath(prefix, "icon_exclamation"),
			tasks.WithROI(roi.AreaExploration.QuestList),
			tasks.WithClick(),
			// tasks.WithWaitMs(500),
			tasks.WithDebugMatch(),
		):
			t.Log.Infof("choose new exploration quest")
		default:
			// drag up
			t.Manager.DragDuration(
				roi.AreaExploration.PtStartDragQuestList,
				roi.AreaExploration.PtStopDragQuestList,
				1500,
			)
			t.WaitMs(1500)
		}

	case domain.TaskStateEnd:
		t.Exit()
	}

	return false
}

func (t *task) IsReady() bool {
	return t.setting.Enable
}
