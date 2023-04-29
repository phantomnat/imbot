package monster_story

import (
	"image"

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

func NewMainStory(manager domain.Manager, setting *TaskSetting) domain.Task {
	status := &TaskStatus{}
	t := &task{
		setting: setting,
		status:  status,
		BaseTask: tasks.NewBaseTask(
			manager, setting, status,
			map[domain.TaskState]string{},
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

			foundIcon, ptIcon := t.Im.MatchPoint(mLeftMenuDetail, domain.MatchOption{
				Path: t.Manager.GetImagePath("icon_main_story_quest"),
				// PrintVal: true,
			})
			if !foundIcon {
				t.Manager.ClickPt(roi.PtActiveQuest)
				t.WaitMs(1000)
				return true
			}

			// mQuestText := m.Region(roiQuestText)
			// defer mQuestText.Close()
			// out := gocv.NewMatWithSize(mQuestText.Rows(), mQuestText.Cols(), gocv.MatTypeCV8UC3)
			// gocv.CvtColor(mQuestText, &out, gocv.ColorRGBAToBGR)
			// gocv.IMWrite("cap/quest_text.png", mQuestText)

			x := roi.ROILeftMenuDetail.Min.X + ptIcon.X
			y := roi.ROILeftMenuDetail.Min.Y + ptIcon.Y
			roiQuestText := image.Rect(x-5, y, x+245, y+45)
			if t.SearchROI(m,
				tasks.WithROI(roiQuestText),
				tasks.WithPath("accept_quest"),
				tasks.WithClick(),
				tasks.WithDebugMatch(),
			) {

			}
		case t.Manager.IsOnMainMenu(m):
			t.Manager.ClickPt(roi.PtTopRightHomeBtn)
			t.WaitMs(1000)

		// case t.SearchROI(m,
		// 	tasks.WithPath(prefix, "txt_guard_journal"),
		// 	tasks.WithROI(roi.ROITopLeft),
		// 	tasks.WithNoWait(),
		// 	tasks.WithNextState(stateFindAndAcceptQuest),
		// ):
		// find and accept quest

		case t.Manager.HandleConversationDialog(m):
		case t.Manager.HandleQuestCompleted(m):
		case t.Manager.HandleVictory(m):
		}

	case domain.TaskStateEnd:
		t.Exit()
	}

	return false
}

func (t *task) IsReady() bool {
	return t.setting.Enable
}
