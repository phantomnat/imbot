package summonerswar

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
	"github.com/phantomnat/imbot/pkg/screen/mumu"
)

const (
	// WindowTitle = "Summoners War: Chronicles"
	WindowTitle = "Chronicles - MuMu Player"

	srcImageDir = "swc"
)

type BotRunningState string

type SummonersWar struct {
	log            *zap.SugaredLogger
	screen         domain.Screen
	currentState   BotState
	muCurrentState sync.RWMutex
	// nextState    BotState
	isRunning *atomic.Bool

	im *im.ImageManager

	muSendCaptureImage sync.RWMutex
	isSendCaptureImage bool
	cbSendCaptureImage func(image.Image)
}

var _ domain.Game = (*SummonersWar)(nil)

func New(imgManager *im.ImageManager) (*SummonersWar, error) {
	sc, err := mumu.NewFromTitle(WindowTitle)
	if err != nil {
		return nil, err
	}

	b := &SummonersWar{
		log:       zap.S().Named("summoners-war-chronicles"),
		screen:    sc,
		isRunning: new(atomic.Bool),
		im:        imgManager,
	}
	return b, nil
}

func (b *SummonersWar) Start() {
	b.isRunning.Store(true)
	b.log.Infof("bot is continue running")
}

func (b *SummonersWar) Pause() {
	b.isRunning.Store(false)
	b.log.Infof("bot is paused")
}

func (b *SummonersWar) Reset() {
	b.log.Debugf("reset!")
	b.Pause()
	b.SetState(StartState)
}

// Run the main loop for bot
func (b *SummonersWar) Run(done <-chan struct{}) {
	oneThirtiethFrameTime := 33 * time.Millisecond
	//
	// b.currentStage = NewStageExpertDarkCastle(b)
	b.SetState(DoQuest)
	//b.SetState(PlayingState)

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

		// TODO: execute the command if any

		startTime := time.Now()

		b.handleState()

		processedTime := time.Since(startTime)
		if oneThirtiethFrameTime > processedTime {
			time.Sleep(oneThirtiethFrameTime - processedTime)
		}
	}
}

func (b *SummonersWar) ToggleSendCaptureImage(isSend bool, cb ...func(image.Image)) {
	b.muSendCaptureImage.Lock()
	defer b.muSendCaptureImage.Unlock()

	b.isSendCaptureImage = true
	if len(cb) > 0 && cb[0] != nil {
		b.cbSendCaptureImage = cb[0]
	}
}

func (b *SummonersWar) sendCaptureImage(m gocv.Mat) {
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

func (b *SummonersWar) imMatchDefault(m gocv.Mat, path ...string) (bool, image.Point) {
	paths := append([]string{srcImageDir}, path...)
	return b.im.MatchDefault(m, paths...)
}

func (b *SummonersWar) imMatchDefaultInROI(m gocv.Mat, roi image.Rectangle, path ...string) (bool, image.Point) {
	mROI := m.Region(roi)
	defer mROI.Close()
	ok, pt := b.imMatchDefault(mROI, path...)
	if ok {
		return ok, image.Point{X: pt.X + roi.Min.X, Y: pt.Y + roi.Min.Y}
	}
	return ok, pt
}

func (b *SummonersWar) imMatchInROI(m gocv.Mat, roi image.Rectangle, o im.MatchOption) (bool, image.Point) {
	mROI := m.Region(roi)
	defer mROI.Close()
	ok, pt := b.im.Match(mROI, srcImageDir+"."+o.Path, o.Th, o)
	if ok {
		return ok, image.Point{X: pt.X + roi.Min.X, Y: pt.Y + roi.Min.Y}
	}
	return ok, pt
}

func ptFromROIandPt(roi image.Rectangle, pt image.Point) image.Point {
	return image.Point{X: pt.X + roi.Min.X, Y: pt.Y + roi.Min.Y}
}

func (b *SummonersWar) handleState() {
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
		b.SetState(ActivateQuest)

	case ActivateQuest:
		// activate quest
		mQuest := m.Region(roiQuest)
		defer mQuest.Close()
		avg := mQuest.Mean()
		if avg.Val1 > thQuestActive.Val1 && avg.Val2 > thQuestActive.Val2 && avg.Val3 > thQuestActive.Val3 {
			b.SetState(DoQuest)
		} else {
			b.screen.MouseMoveAndClickByRect(roiQuest)
			sleepMs(100)
		}
	case DoQuest:
		// detect quest complete dialog
		{
			// quest complete
			foundExp, _ := b.imMatchDefaultInROI(m, roiQuestCompleteExp, "quest_complete", "exp")
			mQuestCompleteBtns := m.Region(roiQuestCompleteBtns)
			defer mQuestCompleteBtns.Close()

			foundBtnOK, ptBtnOk := b.imMatchDefault(mQuestCompleteBtns, "quest_complete", "btn_ok")
			foundBtnNextStory, ptBtnNextStory := b.imMatchDefault(mQuestCompleteBtns, "quest_complete", "btn_next_story")
			var pt image.Point
			var btn string
			var found bool
			switch {
			case foundExp && foundBtnNextStory:
				pt = ptBtnNextStory
				btn = "next story"
				found = true
			case foundExp && foundBtnOK:
				pt = ptBtnOk
				btn = "ok"
				found = true
			}
			if found {
				b.log.Infof("quest completed, click %s (%v)", btn, pt)
				b.screen.MouseMoveAndClickByPoint(ptFromROIandPt(roiQuestCompleteBtns, pt))
				sleepMs(500)
			}
		}
		{
			// sleep mode
			isSleep, _ := b.imMatchDefaultInROI(m, roiSleepModeLogo, "sleep_logo")
			if isSleep {
				// swipe
				b.log.Infof("sleep mode found, swipe")
				b.screen.MouseDrag(650, 600, 650, 400)
			}
		}
		{
			// accept quest
			// TODO: only when in main screen
			canTapToAccept, ptAccept := b.imMatchInROI(m, roiActiveQuest, im.MatchOption{
				Path:     "accept_quest_2",
				Th:       0.03,
				// PrintVal: true,
			})
			canTapToAccept2, _ := b.imMatchInROI(m, roiActiveQuest, im.MatchOption{
				Path:     "accept_quest_1",
				Th:       0.03,
				// PrintVal: true,
			})
			if canTapToAccept && canTapToAccept2 {
				b.log.Infof("tap to accept quest found (%v)", ptAccept)
				b.screen.MouseMoveAndClickByPoint(ptAccept)
				sleepMs(500)
			}
		}
		{
			// Blue quest
			isInAreaExploration, _ := b.imMatchDefaultInROI(m, roiAreaExplorationTitle, "area_exploration", "title_area_exploration")
			if isInAreaExploration {
				foundBtnAccept, ptBtnAccept := b.imMatchDefaultInROI(m, roiAreaExplorationBtns, "area_exploration", "btn_accept")
				if foundBtnAccept {
					b.log.Infof("tap to accept quest found (%v)", ptBtnAccept)
				b.screen.MouseMoveAndClickByPoint(ptBtnAccept)
				sleepMs(500)
				}
			}
		}

		{
			// isQuestInactive, ptQuest := b.imMatchInROI(m, roiActiveQuestIcon, im.MatchOption{
			// 	Path: "inactive_quest",
			// 	Th:   0.1,
			// 	PrintVal: true,
			// })
			// if isQuestInactive {
			// 	b.log.Infof("quest inactive found (%v)", ptQuest)
			// 	b.screen.MouseMoveAndClickByPoint(ptQuest)
			// 	sleepMs(500)
			// }
		}

		// detect talking dialog

	case EndState:

	}
}

func sleepMs(v int) {
	time.Sleep(time.Duration(v) * time.Millisecond)
}

func (b *SummonersWar) SetState(next BotState) {
	b.muCurrentState.Lock()
	defer b.muCurrentState.Unlock()

	b.log.Debugf("changing bot state %s -> %s", b.currentState, next)
	b.currentState = next
}

func (b *SummonersWar) GetState() BotState {
	b.muCurrentState.RLock()
	defer b.muCurrentState.RUnlock()

	return b.currentState
}

func (b *SummonersWar) SentESC() {
	robotgo.KeyTap(robotgo.Escape)
}

func (b *SummonersWar) WindowSize() domain.Rect {
	r, _ := b.screen.GetRect()
	return r
}

func (b *SummonersWar) GetMat() (gocv.Mat, error) {
	if err := b.screen.CaptureToBuffer(); err != nil {
		b.log.Errorf("capture image to buffer: %+v", err)
		return gocv.NewMat(), err
	}

	return b.screen.GetMat()
}

func (b *SummonersWar) GetScreen() domain.Screen {
	return b.screen
}
