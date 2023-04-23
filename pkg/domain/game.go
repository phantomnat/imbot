package domain

import (
	"image"

	"gocv.io/x/gocv"
)

type Game interface {
	Start()
	Pause()
	Reset()
	ToggleSendCaptureImage(isSend bool, cb ...func(image.Image))
	GetScreen() Screen
}

type Screen interface {
	GetRect() (Rect, error)
	GetMat() (gocv.Mat, error)
	CaptureToBuffer() error
	CaptureMatAndSave(filePath string)
	MouseMoveAndClickByRect(roi image.Rectangle, args ...any)
	MouseMoveAndClickByPoint(pt image.Point, args ...any)
	MouseMoveAndClick(x, y int, args ...any)
	MouseDrag(x1, y1, x2, y2 int)
	MouseDragDuration(x1, y1, x2, y2, waitMs int)
	KeyTap(key string, args ...any)
	MouseMove(x, y int)
	Back()
}

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

	// UpdateSetting(v any)
}

type Manager interface {
	// ExitTask resets the current task index to unknown
	ExitTask()

	// im
	GetImageManager() ImageManager
	GetImagePath(path ...string) string
	MatchInROI(m gocv.Mat, roi image.Rectangle, o MatchOption) (bool, image.Point)

	// emu
	Back()
	Click(x, y int)
	ClickPt(pt image.Point)
	Drag(pt1, pt2 image.Point)
	DragDuration(pt1, pt2 image.Point, waitMs int)

	// general
	GoToMainScreen(m gocv.Mat) (done bool)
	IsOnMainScreen(m gocv.Mat) (done bool)
	IsOnMainMenu(m gocv.Mat) (done bool)
	HandleConversationDialog(m gocv.Mat) (done bool)
	HandleQuestCompleted(m gocv.Mat) (done bool)

	StatusManager
}

type StatusManager interface {
	LoadStatus(index int, key string) any
	SaveStatus(index int, key string, v any)
}
