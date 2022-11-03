package bloon_td6

import (
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/phantomnat/imbot/pkg/domain"
	screen "github.com/phantomnat/imbot/pkg/screen"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

const (
	processName = "BloonsTD6.exe"
	WindowTitle = "BloonsTD6"
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

type BloonTD6 struct {
	log          *zap.SugaredLogger
	screen       *screen.Screen
	currentState BotState
}

func New() (*BloonTD6, error) {
	sc, err := screen.NewFromTitle(WindowTitle)
	if err != nil {
		return nil, err
	}

	b := &BloonTD6{
		log:    zap.S().Named("bloon-td-6"),
		screen: sc,
	}
	return b, nil
}

func (b *BloonTD6) Run(done <-chan struct{}) {
	b.SetState(StartState)
	for {
		select {
		case <-done:
			b.log.Infof("exiting")
			return
		default:
		}

		switch b.currentState {
		case StartState:
			// detect main menu
			// click play
			time.Sleep(time.Second)
			b.SetState(StageSelectionState)
		case StageSelectionState:
			// select stage
			// choose easy
			// click play
			time.Sleep(time.Second)
			b.SetState(PlayingState)
		case PlayingState:
			// hardest
			// building
			// upgrading
			// ack the level up
			// wait for win or lose dialog
			// restart if lose
			// go to collect reward if win
			time.Sleep(time.Second * 10)
			b.SetState(CollectRewardState)
		case CollectRewardState:
			time.Sleep(time.Second)
			b.SetState(EndState)
		case EndState:
			// go to start stage
			time.Sleep(time.Second)
			b.SetState(StartState)
		}
	}
}

func (b *BloonTD6) SetState(next BotState) {
	b.log.Debugf("changing bot state %s -> %s", b.currentState, next)
	b.currentState = next
}

func (b *BloonTD6) SentESC() {
	robotgo.KeyTap(robotgo.Escape)
}

func (b *BloonTD6) WindowSize() *domain.Rect {
	r, _ := b.screen.GetRect()
	return r
}

func (b *BloonTD6) GetMat() (gocv.Mat, error) {
	if err := b.screen.CaptureToBuffer(); err != nil {
		b.log.Errorf("capture image to buffer: %+v", err)
		return gocv.NewMat(), err
	}

	return b.screen.GetMat()
}

func (b *BloonTD6) GetMatFromRobotgo() (gocv.Mat, error) {
	rect, err := b.screen.GetRect()
	if err != nil {
		return gocv.NewMat(), err
	}

	bit := robotgo.CaptureScreen(rect.X, rect.Y, rect.Width, rect.Height)
	defer robotgo.FreeBitmap(bit)

	img := robotgo.ToImage(bit)
	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return gocv.NewMat(), err
	}
	return mat, nil
}
