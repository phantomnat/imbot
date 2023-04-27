package tasks_test

import (
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gocv.io/x/gocv"
	"sigs.k8s.io/yaml"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks"
)

type testTask struct {
	tasks.BaseTask

	setting *taskSetting
	status  *taskStatus
}

type taskSetting struct {
	domain.TaskSettingBase
}

type taskStatus struct {
	domain.TaskStatusBase
	Test int
	Next time.Time
}

type tasksStatus struct {
	Names map[string]any
}

func TestTaskBase(t *testing.T) {
	a := assert.New(t)

	setting := taskSetting{}
	status := &taskStatus{}
	manager := &mockManager{
		status: make(map[string]any),
	}

	config := []byte(`enable: true`)
	txtStatus := []byte(`Next: "2023-04-27T00:00:00+07:00"
NextExecution: "2023-04-27T09:43:14.2483705+07:00"
NextReset: "2023-04-28T00:00:00+07:00"
State: 9999
Test: 26`)
	var inStatus any
	err := yaml.Unmarshal(config, &setting)
	a.Nilf(err, "%+v", err)
	err = yaml.Unmarshal(txtStatus, &inStatus)
	a.Nilf(err, "%+v", err)

	task := &testTask{
		setting: &setting,
		status:  status,
		BaseTask: tasks.NewBaseTask(
			manager, &setting, status,
			nil,
		),
	}
	task.LoadStatus(inStatus)

	a.Equal(26, task.status.Test)

	a.True(task.IsReady(), "task should ready")

	task.status.NextExecution = time.Now().Add(time.Second)
	a.False(task.IsReady(), "task should not be ready")

	time.Sleep(time.Second)
	a.True(task.IsReady(), "task should ready")

	task.SetState(domain.TaskStateEnd)
	task.SaveStatus()

	fmt.Printf("%s\n", manager.status[task.GetName()])
	fmt.Println("------")

	if task.status.Reset(func(today time.Time) {
		task.status.Next = today
		task.status.Test = today.Day()
	}) {
		task.SaveStatus()
	}

	fmt.Printf("%s\n", manager.status[task.GetName()])
	fmt.Println("------")

}

type mockManager struct {
	mock.Mock

	status map[string]any
}

var _ domain.Manager = (*mockManager)(nil)

func (m *mockManager) ExitTask() {
}
func (m *mockManager) GetImageManager() domain.ImageManager {
	return nil
}
func (m *mockManager) GetImagePath(path ...string) string {
	return ""
}
func (m *mockManager) MatchInROI(_ gocv.Mat, roi image.Rectangle, o domain.MatchOption) (bool, image.Point) {
	return false, image.Pt(0, 0)
}
func (m *mockManager) Back() {
}
func (m *mockManager) Click(x, y int) {
}
func (m *mockManager) ClickPt(pt image.Point) {
}
func (m *mockManager) Drag(pt1, pt2 image.Point) {
}
func (m *mockManager) DragDuration(pt1, pt2 image.Point, waitMs int) {
}
func (m *mockManager) GoToMainScreen(_ gocv.Mat) (done bool) {
	return false
}
func (m *mockManager) IsOnMainScreen(_ gocv.Mat) (done bool) {
	return false
}
func (m *mockManager) IsOnMainMenu(_ gocv.Mat) (done bool) {
	return false
}
func (m *mockManager) HandleConversationDialog(_ gocv.Mat) (done bool) {
	return false
}
func (m *mockManager) HandleQuestCompleted(_ gocv.Mat) (done bool) {
	return false
}
func (m *mockManager) SaveStatus(key string, v any) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return
	}
	m.status[key] = data
}
