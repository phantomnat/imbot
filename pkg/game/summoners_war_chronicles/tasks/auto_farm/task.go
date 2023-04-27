package auto_farm

import (
	"time"

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

	MiniBoss string
	Interval domain.Duration
}

type TaskStatus struct {
	domain.TaskStatusBase

	// State         domain.TaskState
	// NextExecution time.Time
	LastMoving time.Time
}

var _ domain.Task = (*task)(nil)

const (
	stateGoToMainScreen domain.TaskState = iota + 1

	stateOpenMap
	stateFindCreature
	stateMoving

	stateEnsureAutoBattle
)
const (
	prefix = "auto_farm"
)

func NewAutoFarm(index int, manager domain.Manager, setting TaskSetting) domain.Task {
	status := &TaskStatus{}
	t := &task{
		setting: &setting,
		status:  status,
		BaseTask: tasks.NewBaseTask(
			manager, &setting, status,
			map[domain.TaskState]string{
				stateGoToMainScreen:   "go_to_main_screen",
				stateOpenMap:          "open_map",
				stateFindCreature:     "find_creature",
				stateMoving:           "moving",
				stateEnsureAutoBattle: "ensure_auto_battle",
			},
		),
	}
	return t
}

func (t *task) Do(m gocv.Mat) (triggered bool) {

	triggered = true

	switch t.State {
	case domain.TaskStateBegin:
		t.SetState(stateGoToMainScreen)

	case stateGoToMainScreen:
		switch {
		case t.setting.MiniBoss != "" && t.SearchROI(m,
			tasks.WithROI(roi.AutoFarm.MonsterName),
			tasks.WithPath(prefix, "title_"+t.setting.MiniBoss),
			tasks.WithNoWait(),
			tasks.WithNextState(stateEnsureAutoBattle),
		):
			t.Log.Infof("boss '%s' detected, skipping", t.setting.MiniBoss)

		case time.Since(t.StateChangedAt).Seconds() > 5:
			if t.setting.MiniBoss != "" {
				t.Manager.ClickPt(roi.MainScreen.PtMinimap)
				t.WaitMs(1000)
				t.SetState(stateOpenMap)
				return
			}
			if t.Manager.IsOnMainScreen(m) {
				t.SetState(stateEnsureAutoBattle)
			}
		case t.Manager.IsOnMainMenu(m):
		default:
		}

	case stateOpenMap:
		switch {
		case t.SearchROI(m,
			tasks.WithROI(roi.ROITopLeft),
			tasks.WithPath(prefix, "title_map"),
			tasks.WithNoWait(),
		):
			// TODO: might need to search creature list label instead
			t.Manager.ClickPt(roi.AutoFarm.PtOpenCreatureList)
			t.WaitMs(1000)
			t.SetState(stateFindCreature)

		case time.Since(t.StateChangedAt).Seconds() > 5:
			t.SetState(stateGoToMainScreen)
		}

	case stateFindCreature:
		switch {
		case t.SearchROI(m,
			tasks.WithROI(roi.ROITopLeft),
			tasks.WithPath(prefix, "title_map"),
			tasks.WithNoWait(),
		):
			mCreatureList := m.Region(roi.AutoFarm.CreatureList)
			defer mCreatureList.Close()

			found, pt := t.Im.MatchPoint(mCreatureList, domain.MatchOption{
				Path: t.Manager.GetImagePath(prefix, "text_creature_list_"+t.setting.MiniBoss),
			})
			if found {
				newPt := t.GetPtWithROI(roi.AutoFarm.CreatureList, pt)
				t.Manager.Click(newPt.X+330, newPt.Y+15)
				t.WaitMs(1000)
				t.status.LastMoving = time.Now()
				t.SetState(stateMoving)
			}
		}

	case stateMoving:
		switch {
		case t.SearchROI(m,
			tasks.WithROI(roi.AutoFarm.Moving),
			tasks.WithPath(prefix, "text_moving"),
			tasks.WithNoWait(),
		):
			t.WaitMs(500)
			t.status.LastMoving = time.Now()
			t.Log.Debug("moving...")
		case time.Since(t.status.LastMoving).Seconds() > 1:
			t.SetState(stateEnsureAutoBattle)
			t.Manager.ClickPt(roi.MainScreen.PtBasicAttack)
			t.WaitMs(500)
		}

	case stateEnsureAutoBattle:
		if t.SearchROI(m,
			tasks.WithPath("icon_auto_battle"),
			tasks.WithROI(roi.MainScreen.AutoBattleIcon),
			tasks.WithClick(),
			tasks.WithWaitMs(800),
			// tasks.WithDebugMatch(),
			// tasks.WithNextState(domain.TaskStateEnd),
		) {
			return
		} else if time.Since(t.StateChangedAt).Seconds() > 6 {
			t.SetState(domain.TaskStateEnd)
		}

	case domain.TaskStateEnd:
		t.status.NextExecution = time.Now().Add(time.Duration(t.setting.Interval))
		t.SaveStatus()
		t.Exit()

	default:
		triggered = false
	}

	return
}

func (t *task) IsReady() bool {
	if !t.setting.Enable {
		return false
	}
	return t.status.NextExecution.IsZero() || time.Now().After(t.status.NextExecution)
}

func (t *task) SaveStatus() {
	t.Manager.SaveStatus(t.Name, t.status)
}
