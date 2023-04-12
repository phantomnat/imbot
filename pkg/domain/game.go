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
	KeyTap(key string, args ...any)
	MouseMove(x, y int)
}
