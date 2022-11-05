package ui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/phantomnat/imbot/pkg/domain"
	"go.uber.org/zap"
)

const (
	AppName = "imbot"
)

type uiHandler struct {
	log *zap.SugaredLogger

	game domain.Game

	app            fyne.App
	mainWindow     fyne.Window
	debugWindow    fyne.Window
	isDebugCapture bool
	// vboxMain  *fyne.Container
	// lblStatus *widget.Label
	isBotRunning    bool
	btnToggleRun    *widget.Button
	btnReset        *widget.Button
	btnShowDebugWin *widget.Button
}

func New(g domain.Game) *uiHandler {
	handler := &uiHandler{
		log:  zap.S().Named("ui"),
		app:  app.New(),
		game: g,
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

	h.debugWindow = h.app.NewWindow("debug")

	h.btnToggleRun = widget.NewButton("start", h.OnBtnToggleRunClicked)
	h.btnReset = widget.NewButton("reset", h.onBtnResetClicked)
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
	if h.isDebugCapture {
		h.game.ToggleSendCaptureImage(false)
		h.debugWindow.Hide()
	} else {
		h.game.ToggleSendCaptureImage(true, h.OnImageUpdated)
		h.debugWindow.Show()
	}
}

func (h *uiHandler) OnBtnToggleRunClicked() {
	if h.isBotRunning {
		h.game.Pause()
		h.btnToggleRun.SetText("pause")
	} else {
		h.game.Start()
		h.btnToggleRun.SetText("start")
	}
	h.isBotRunning = !h.isBotRunning
}

func (h *uiHandler) onBtnResetClicked() {
	h.isBotRunning = false
	h.btnToggleRun.SetText("pause")
	h.game.Reset()
}
