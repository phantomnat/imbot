package summonerswar

import (
	"image"
	"sync"
	"sync/atomic"
	"time"

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

	configFile = "configs/swc/config.yaml"
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

	setting Setting
}

var _ domain.Game = (*SummonersWar)(nil)

func New(imgManager *im.ImageManager) (*SummonersWar, error) {
	sc, err := mumu.NewFromTitle(WindowTitle, mumu.Option{
		AutoResize: true,
		Width:      1280,
		Height:     720,
	})
	if err != nil {
		return nil, err
	}

	setting, err := LoadSetting(configFile)
	if err != nil {
		return nil, err
	}
	b := &SummonersWar{
		log:       zap.S().Named("summoners-war-chronicles"),
		screen:    sc,
		isRunning: new(atomic.Bool),
		im:        imgManager,
		setting:   setting,
	}
	return b, nil
}

func (b *SummonersWar) Start() {
	b.log.Infof("reloading config...")
	setting, err := LoadSetting(configFile)
	if err != nil {
		b.log.Errorf("cannot load config from %s: %+v", configFile, err)
	} else {
		b.setting = setting
		b.log.Infof("config reloaded")
	}
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

func (b *SummonersWar) handleState() {
	// capture
	m, err := b.GetMat()
	if err != nil {
		return
	}
	defer m.Close()

	// b.sendCaptureImage(m)

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

		switch {
		case b.handleMainScreen(m):
		case b.handleQuestComplete(m):
		case b.handleSleepScreen(m):
		case b.handleAreaExploration(m):
		case b.handleMonsterStory(m):
		case b.handleDialog(m):
		}

	case EndState:

	}
}

func (b *SummonersWar) handleMonsterStory(m gocv.Mat) bool {
	if b.setting.Mode != BotModeMonsterStory {
		return false
	}

	prefix := "guard_journal"
	log := b.log.Named("monster-story")

	// accept quest
	foundGuardJournal, _ := b.imMatchDefaultInROI(m, roiTopLeft, prefix, "txt_guard_journal")
	if foundGuardJournal {
		mBtns := m.Region(roiMonsterStory.Buttons)
		defer mBtns.Close()

		var isFound = true
		var pt image.Point
		var txt string

		if foundNextStageBtn, ptNextStageBtn := b.imMatchDefault(mBtns, prefix, "btn_next_stage"); foundNextStageBtn {
			pt = ptNextStageBtn
			txt = "next stage"

		} else if foundClaimBtn, ptClaimBtn := b.imMatchDefault(mBtns, prefix, "btn_claim"); foundClaimBtn {
			pt = ptClaimBtn
			txt = "claim"
		} else if foundSearchBtn, ptSearcBtn := b.imMatchDefault(mBtns, prefix, "btn_search"); foundSearchBtn {
			pt = ptSearcBtn
			txt = "search"
		} else {
			isFound = false
		}

		if isFound {
			log.Infof("%s (%v)", txt, ptFromROIandPt(roiMonsterStory.Buttons, pt))
			b.screen.MouseMoveAndClickByPoint(ptFromROIandPt(roiMonsterStory.Buttons, pt))
			sleepMs(500)
		}
		return true
	}

	if foundModalStartStory, _ := b.imMatchDefaultInROI(m, roiMonsterStory.ModalStartStory, prefix, "txt_start_story"); foundModalStartStory {
		foundOkBtn, ptOKBtn := b.imMatchDefaultInROI(m, roiMonsterStory.ModalStartStoryButtons, prefix, "btn_ok")
		if foundOkBtn {
			log.Infof("start story (%v)", ptOKBtn)
			b.screen.MouseMoveAndClickByPoint(ptOKBtn)
			sleepMs(500)
		}
		return true
	}

	return false
}

func (b *SummonersWar) handleMonsterStoryQuestState(m gocv.Mat) bool {
	prefix := "guard_journal"
	log := b.log.Named("monster-story")

	isQuestActive, _ := b.imMatchDefaultInROI(m, roiLeftMenuDetail, prefix, "icon_ongoing_quest")
	if isQuestActive {
		log.Infof("waiting for the quest to finish...")
		sleepMs(500)
		return true
	}

	// quest done
	isQuestFinish, ptQuestFinishBtn := b.imMatchDefaultInROI(m, roiLeftMenuDetail, prefix, "icon_finish_quest")
	if isQuestFinish {
		log.Infof("quest finish (%v)", ptQuestFinishBtn)
		b.screen.MouseMoveAndClickByPoint(ptQuestFinishBtn)
		sleepMs(1000)
		return true
	}

	return true
}

func (b *SummonersWar) handleMainScreen(m gocv.Mat) bool {
	foundMainScreen, _ := b.imMatchInROI(m, roiTopRigthMenuBtn, im.MatchOption{
		Path: "btn_top_right_menu",
		Th:   0.01,
		// PrintVal: true,
	})
	if !foundMainScreen {
		return false
	}

	isQuestActive, _ := b.imMatchInROI(m, roiLeftMenu, im.MatchOption{
		Path: "btn_quest_active",
		Th:   0.01,
		// PrintVal: true,
		// Normalize: true,
	})
	if !isQuestActive {
		b.log.Infof("tap to active quest")
		b.screen.MouseMoveAndClickByPoint(ptActiveQuest)
		sleepMs(1000)
		return true
	}

	switch b.setting.Mode {
	case BotModeStoryQuest:
		canTapToAccept, ptAccept := b.imMatchInROI(m, roiLeftMenuDetail, im.MatchOption{
			Path: "accept_quest",
			Th:   0.01,
			// PrintVal: true,
			// Normalize: true,
		})
		if canTapToAccept {
			b.log.Infof("tap to accept quest found (%v)", ptAccept)
			b.screen.MouseMoveAndClickByPoint(ptAccept)
			sleepMs(500)
			return true
		}
	case BotModeExplorationQuest:
		canAcceptExploration, ptAcceptExploration := b.imMatchInROI(m, roiLeftMenuDetail, im.MatchOption{
			Path: "exploration_quest",
			Th:   0.01,
			// PrintVal: true,
			// Normalize: true,
		})
		if canAcceptExploration {
			b.log.Infof("tap to accept exploration quest found (%v)", ptAcceptExploration)
			b.screen.MouseMoveAndClickByPoint(ptAcceptExploration)
			sleepMs(500)
			return true
		}
	case BotModeMonsterStory:
		return b.handleMonsterStoryQuestState(m)
	default:
		return false
	}

	// accept quest
	// TODO: only when in main screen
	return true
}

func (b *SummonersWar) handleQuestComplete(m gocv.Mat) bool {
	// quest complete
	prefix := "quest_complete"

	foundTapToClose, _ := b.imMatchInROI(m, roiQuestCompleteTapToClose, im.MatchOption{
		Path: prefix + ".txt_tab_to_close",
		Th:   0.01,
		// PrintVal: true,
	})
	if foundTapToClose {
		b.log.Infof("quest completed, click anywhere to close")
		b.screen.MouseMoveAndClickByPoint(ptContinue)
		sleepMs(1000)
		return true
	}

	foundModal, ptModalBtnOK := b.imMatchInROI(m, roiModalComplete, im.MatchOption{
		Path: prefix + ".btn_ok_modal",
	})
	if foundModal {
		b.log.Infof("modal found, click %v to close", ptModalBtnOK)
		b.screen.MouseMoveAndClickByPoint(ptModalBtnOK)
		sleepMs(500)
		return true
	}

	foundExp, _ := b.imMatchDefaultInROI(m, roiQuestCompleteExp, prefix, "exp")
	mQuestCompleteBtns := m.Region(roiQuestCompleteBtns)
	defer mQuestCompleteBtns.Close()

	foundBtnOK, ptBtnOk := b.imMatchDefault(mQuestCompleteBtns, prefix, "btn_ok")
	foundBtnNextStory, ptBtnNextStory := b.imMatchDefault(mQuestCompleteBtns, prefix, "btn_next_story")

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
	if !found {
		return false
	}

	b.log.Infof("quest completed, click %s (%v)", btn, pt)
	b.screen.MouseMoveAndClickByPoint(ptFromROIandPt(roiQuestCompleteBtns, pt))
	sleepMs(500)
	return true
}

func (b *SummonersWar) handleSleepScreen(m gocv.Mat) bool {
	found, _ := b.imMatchDefaultInROI(m, roiSleepModeLogo, "sleep_logo")
	if !found {
		return false
	}

	b.screen.MouseDrag(ptSleepModeWakeFrom.X, ptSleepModeWakeFrom.Y, ptSleepModeWakeTo.X, ptSleepModeWakeTo.Y)
	sleepMs(3000)
	return true
}

func (b *SummonersWar) handleAreaExploration(m gocv.Mat) (found bool) {
	prefix := "area_exploration"
	found, _ = b.imMatchDefaultInROI(m, roiAreaExplorationTitle, prefix, "title_area_exploration")
	if !found {
		return
	}

	canAccept, ptAccept := b.imMatchDefaultInROI(m, roiAreaExplorationBtns, prefix, "btn_accept")
	if canAccept {
		b.log.Infof("accept exploration quest (%v)", ptAccept)
		b.screen.MouseMoveAndClickByPoint(ptAccept)
		sleepMs(500)
		return
	}

	isNewQuest, ptNewQuest := b.imMatchInROI(m, roiAreaExplorationNewQuest, im.MatchOption{
		Path: prefix + ".icon_exclamation",
	})
	if isNewQuest {
		b.log.Infof("choose new exploration quest (%v)", ptNewQuest)
		b.screen.MouseMoveAndClickByPoint(ptNewQuest)
		sleepMs(1000)
	}

	return
}

func (b *SummonersWar) handleDialog(m gocv.Mat) bool {
	prefix := "dialog"
	// foundBtnBack, _ := b.imMatchDefaultInROI(m, roiSleepModeLogo, prefix, "btn_back")
	// foundTxtAuto, _ := b.imMatchDefaultInROI(m, roiSleepModeLogo, prefix, "txt_auto")
	foundBtnBack, _ := b.imMatchInROI(m, roiBtnBack, im.MatchOption{
		Path: prefix + ".btn_back",
		Th:   0.01,
		// PrintVal: true,
	})
	foundTxtAuto, _ := b.imMatchInROI(m, roiTxtAutoAndIcon, im.MatchOption{
		Path:     prefix + ".txt_auto_and_icon",
		Th:       0.01,
		PrintVal: true,
	})

	if !(foundTxtAuto && foundBtnBack) {
		return false
	}
	b.log.Infof("dialog detected")
	b.screen.MouseMoveAndClickByPoint(ptContinue)
	sleepMs(600)
	return true
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
