package domain

import (
	"time"

	"gocv.io/x/gocv"
)

type TaskState int

const (
	TaskStateBegin TaskState = 0
	TaskStateEnd   TaskState = 9999
)

type Task interface {
	GetName() string

	IsReady() bool

	// NeedMainScreen indicates that task required to go to main screen before start
	IsNeedMainScreen() bool

	// CanInterrupt returns true if task can be cancel
	CanInterrupt() bool

	// Do the task
	Do(m gocv.Mat) bool

	// Exit the current task and go to main screen
	RequestExit()

	// GetState returns state in string
	GetState() string

	LoadStatus(any)

	Reset()
	// UpdateSetting(v any)
}

type TaskSetting interface {
	IsEnabled() bool
}

type TaskSettingBase struct {
	Enable bool
}

var _ TaskSetting = (*TaskSettingBase)(nil)

func (t *TaskSettingBase) IsEnabled() bool {
	if t == nil {
		return false
	}
	return t.Enable
}

// TaskStatus
type TaskStatus interface {
	GetState() TaskState
	SetState(state TaskState)
	GetNextExecution() time.Time
	SetNextExecution(next time.Time)

	IsReady() bool
}

type TaskStatusBase struct {
	NextExecution time.Time
	NextReset     time.Time
	State         TaskState
}

var _ TaskStatus = (*TaskStatusBase)(nil)

func (t *TaskStatusBase) GetNextExecution() time.Time {
	return t.NextExecution
}
func (t *TaskStatusBase) SetNextExecution(next time.Time) {
	t.NextExecution = next
}
func (t *TaskStatusBase) GetState() TaskState {
	return t.State
}
func (t *TaskStatusBase) SetState(state TaskState) {
	t.State = state
}
func (t *TaskStatusBase) IsReady() bool {
	return t.NextExecution.IsZero() || time.Now().After(t.NextExecution)
}
func (t *TaskStatusBase) Reset(cb func(today time.Time)) (triggered bool) {
	now := time.Now()
	triggered = now.After(t.NextReset)

	if triggered {
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		cb(today)
		t.NextReset = today.AddDate(0, 0, 1)
	}
	return
}
