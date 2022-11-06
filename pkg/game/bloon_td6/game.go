package bloon_td6

import (
	"image"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-vgo/robotgo"
	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/im"
	screen "github.com/phantomnat/imbot/pkg/screen"
)

const (
	processName = "BloonsTD6.exe"
	WindowTitle = "BloonsTD6"

	srcImageDir = "bloons_td_6"
)

type BotState string

var (
	UnknownState        BotState = "unknown"
	StartState          BotState = "start"
	StageSelectionState BotState = "stage selection"
	PlayingState        BotState = "playing"
	CollectRewardState  BotState = "collect reward"
	EndState            BotState = "end"
)

type BotRunningState string

type BloonsTD6 struct {
	log            *zap.SugaredLogger
	screen         *screen.Screen
	currentState   BotState
	muCurrentState sync.RWMutex
	// nextState    BotState
	isRunning *atomic.Bool

	im *im.ImageManager

	muSendCaptureImage sync.RWMutex
	isSendCaptureImage bool
	cbSendCaptureImage func(image.Image)
}

var _ domain.Game = (*BloonsTD6)(nil)

func New(imgManager *im.ImageManager) (*BloonsTD6, error) {
	sc, err := screen.NewFromTitle(WindowTitle)
	if err != nil {
		return nil, err
	}

	b := &BloonsTD6{
		log:       zap.S().Named("bloons-td-6"),
		screen:    sc,
		isRunning: new(atomic.Bool),
		im:        imgManager,
	}
	return b, nil
}

func (b *BloonsTD6) Start() {
	b.isRunning.Store(true)
	b.log.Infof("bot is continue running")
}

func (b *BloonsTD6) Pause() {
	b.isRunning.Store(false)
	b.log.Infof("bot is paused")
}

func (b *BloonsTD6) Reset() {
	b.log.Debugf("reset!")
	b.Pause()
	b.SetState(StartState)
}

// Run the main loop for bot
func (b *BloonsTD6) Run(done <-chan struct{}) {
	oneThirtiethFrameTime := 33 * time.Millisecond

	b.SetState(StartState)
	for {
		select {
		case <-done:
			b.log.Infof("exiting")
			return
		default:
		}

		// pause
		if !b.isRunning.Load() {
			time.Sleep(oneThirtiethFrameTime)
			continue
		}

		startTime := time.Now()

		b.handleState()

		processedTime := time.Since(startTime)
		if oneThirtiethFrameTime > processedTime {
			time.Sleep(oneThirtiethFrameTime - processedTime)
		}
	}
}

func (b *BloonsTD6) ToggleSendCaptureImage(isSend bool, cb ...func(image.Image)) {
	b.muSendCaptureImage.Lock()
	defer b.muSendCaptureImage.Unlock()

	b.isSendCaptureImage = true
	if len(cb) > 0 && cb[0] != nil {
		b.cbSendCaptureImage = cb[0]
	}
}

func (b *BloonsTD6) sendCaptureImage(m gocv.Mat) {
	b.muSendCaptureImage.RLock()
	defer b.muSendCaptureImage.RUnlock()

	if !b.isSendCaptureImage {
		return
	}

	img, err := m.ToImage()
	if err != nil {
		// ignore error
		return
	}
	b.cbSendCaptureImage(img)
}

func Rect(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

var (
	ptBtnAdvanced = image.Pt(730, 650)
	ptBtnExpert   = image.Pt(890, 650)
	ptBtnEasy     = image.Pt(430, 270)
	ptBtnStandard = image.Pt(430, 390)
	roiStage      = Rect(236, 84, 808, 400)

	// playing
	roiVictory = Rect(442, 63, 397, 108)

	// collect reward
	ptBtnNext = image.Pt(640, 600)
	ptBtnHome = image.Pt(465, 565)
)

func (b *BloonsTD6) imMatchDefault(m gocv.Mat, path ...string) (bool, image.Point) {
	paths := append([]string{srcImageDir}, path...)
	return b.im.MatchDefault(m, paths...)
}

func (b *BloonsTD6) imMatchDefaultInROI(m gocv.Mat, roi image.Rectangle, path ...string) (bool, image.Point) {
	mROI := m.Region(roi)
	defer mROI.Close()
	ok, pt := b.imMatchDefault(mROI, path...)
	if ok {
		return ok, image.Point{X: pt.X + roi.Min.X, Y: pt.Y + roi.Min.Y}
	}
	return ok, pt
}

func (b *BloonsTD6) handleState() {
	// capture
	m, err := b.GetMat()
	if err != nil {
		b.log.Errorf("cannot capture screen: %+v", err)
		return
	}
	defer m.Close()

	b.sendCaptureImage(m)

	switch b.GetState() {
	case StartState:
		// detect main menu
		ok, pt := b.imMatchDefault(m, "btn-play")
		if ok {
			b.screen.MouseMoveAndClick(pt.X, pt.Y)
			b.SetState(StageSelectionState)
		}
	case StageSelectionState:
		ok, pt := b.imMatchDefaultInROI(m, roiStage, "lv-expert", "ouch")
		// ok, pt := b.imMatchDefaultInROI(m, roiStage, "lv-advanced", "midnight-mansion")
		if !ok {
			b.screen.MouseMoveAndClick(ptBtnExpert.X, ptBtnExpert.Y)
			// b.screen.MouseMoveAndClick(ptBtnAdvanced.X, ptBtnAdvanced.Y)
			time.Sleep(time.Millisecond * 500)
		} else {
			b.screen.MouseMoveAndClick(pt.X, pt.Y)
			time.Sleep(time.Millisecond * 500)
			b.screen.MouseMoveAndClick(ptBtnEasy.X, ptBtnEasy.Y)
			time.Sleep(time.Millisecond * 500)
			b.screen.MouseMoveAndClick(ptBtnStandard.X, ptBtnStandard.Y)
			b.SetState(PlayingState)
		}
		// select stage

		// choose easy
		// click play
		// b.SetState(PlayingState)
	case PlayingState:
		if ok, _ := b.imMatchDefaultInROI(m, roiVictory, "victory"); ok {
			b.SetState(CollectRewardState)
			return
		}
		// hardest
		// building
		// upgrading
		// ack the level up
		// wait for win or lose dialog
		// restart if lose
		// go to collect reward if win
	case CollectRewardState:
		if ok, _ := b.imMatchDefault(m, "victory"); ok {
			b.screen.MouseMoveAndClick(ptBtnNext.X, ptBtnNext.Y)
			time.Sleep(time.Millisecond * 1000)
			b.screen.MouseMoveAndClick(ptBtnHome.X, ptBtnHome.Y)
			time.Sleep(time.Millisecond * 1000)

			b.SetState(EndState)
		}
	case EndState:
		if ok, _ := b.imMatchDefault(m, "btn-play"); ok {
			b.SetState(StartState)
		}
	}
}

func (b *BloonsTD6) SetState(next BotState) {
	b.muCurrentState.Lock()
	defer b.muCurrentState.Unlock()

	b.log.Debugf("changing bot state %s -> %s", b.currentState, next)
	b.currentState = next
}

func (b *BloonsTD6) GetState() BotState {
	b.muCurrentState.RLock()
	defer b.muCurrentState.RUnlock()

	return b.currentState
}

func (b *BloonsTD6) SentESC() {
	robotgo.KeyTap(robotgo.Escape)
}

func (b *BloonsTD6) WindowSize() *domain.Rect {
	r, _ := b.screen.GetRect()
	return r
}

func (b *BloonsTD6) GetMat() (gocv.Mat, error) {
	if err := b.screen.CaptureToBuffer(); err != nil {
		b.log.Errorf("capture image to buffer: %+v", err)
		return gocv.NewMat(), err
	}

	return b.screen.GetMat()
}

func (b *BloonsTD6) GetScreen() *screen.Screen {
	return b.screen
}
