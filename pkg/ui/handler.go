package ui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

const (
	AppName = "imbot"
)

type uiHandler struct {
	log *zap.SugaredLogger

	app         fyne.App
	mainWindow  fyne.Window
	debugWindow fyne.Window

	// vboxMain  *fyne.Container
	// lblStatus *widget.Label

	btnToggleRun    *widget.Button
	btnReset        *widget.Button
	btnShowDebugWin *widget.Button
}

func New() *uiHandler {
	handler := &uiHandler{
		log: zap.S().Named("ui"),
		app: app.New(),
	}

	handler.init()

	return handler
}

func (h *uiHandler) Run() {
	h.mainWindow.ShowAndRun()
}

func (h *uiHandler) init() {
	h.mainWindow = h.app.NewWindow(AppName)
	h.mainWindow.Resize(fyne.Size{Width: 300})
	// h.mainWindow.SetFixedSize(true)

	h.debugWindow = h.app.NewWindow("debug")

	h.btnToggleRun = widget.NewButton("Start", func() {})
	h.btnReset = widget.NewButton("reset", func() {})
	h.btnShowDebugWin = widget.NewButton("win", h.onShowDebugWinClicked)
	// h.btnShowDebugWin.Resize(fyne.NewSize(100, 40))

	// h.lblStatus = widget.NewLabel("test")
	// h.lblStatus.TextStyle = fyne.TextStyle{Monospace: true}
	// h.lblStatus.Wrapping = fyne.TextTruncate

	h.mainWindow.SetContent(
		container.New(
			layout.NewVBoxLayout(),
			container.NewHBox(h.btnToggleRun, h.btnReset, h.btnShowDebugWin),
		),
	)
}

func (h *uiHandler) OnImageUpdated(v image.Image) {
	h.debugWindow.Resize(fyne.NewSize(float32(v.Bounds().Dx()), float32(v.Bounds().Dy())))
	h.debugWindow.SetContent(canvas.NewImageFromImage(v))
}

func (h *uiHandler) onShowDebugWinClicked() {
	h.debugWindow.Show()
}
