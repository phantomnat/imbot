package summonerswar

import (
	"image"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"gocv.io/x/gocv"
	"sigs.k8s.io/yaml"

	"github.com/phantomnat/imbot/pkg/domain"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/roi"
	"github.com/phantomnat/imbot/pkg/im"
	"github.com/phantomnat/imbot/pkg/screen/mumu"

	area_exploration "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/area_exploration"
	auto_farm "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/auto_farm"
	challenge_arena "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/challenge_arena"
	monster_story "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/monster_story"
)

const (
	// WindowTitle = "Summoners War: Chronicles"
	// WindowTitle = "Chronicles - MuMu Player"
	// WindowTitle = "BlueStacks App Player"

	srcImageDir = "swc"

	configFile = "configs/swc/config.yaml"
	StatusFile = "configs/swc/status.yaml"
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

	setting          Setting
	tasks            []domain.Task
	taskStatuses     []any
	currentTaskIndex int

	taskAreaExploration domain.Task
	taskMonsterStory    domain.Task
}

var _ domain.Game = (*SummonersWar)(nil)

func New(imgManager *im.ImageManager) (*SummonersWar, error) {
	var err error
	setting, err := LoadSetting(configFile)
	if err != nil {
		return nil, err
	}

	var sc domain.Screen
	emuOpt := mumu.Option{
		AutoResize: true,
		Width:      1280,
		Height:     720,
	}

	switch setting.Emu {
	case EmuTypeBlueStack:
		emuOpt.Title = "BlueStacks App Player"
		emuOpt.ADBPort = 5555
		sc, err = mumu.NewBlueStack(emuOpt)
	case EmuTypeMumu:
		emuOpt.Title = "Chronicles - MuMu Player"
		emuOpt.ADBPort = 7555
		sc, err = mumu.NewMumu(emuOpt)
	}
	if err != nil {
		return nil, err
	}

	b := &SummonersWar{
		log:              zap.S().Named("summoners-war-chronicles"),
		screen:           sc,
		isRunning:        new(atomic.Bool),
		im:               imgManager,
		setting:          setting,
		currentTaskIndex: TaskUnknown,
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

	b.LoadTasks()

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

func (b *SummonersWar) LoadTasks() {
	if b.setting.AreaExploration != nil {
		b.taskAreaExploration = area_exploration.NewAreaExploration(0, b, *b.setting.AreaExploration)
	}
	if b.setting.MonsterStory != nil {
		b.taskMonsterStory = monster_story.NewMonsterStory(0, b, *b.setting.MonsterStory)
	}

	b.tasks = make([]domain.Task, 0, len(b.setting.Tasks))
	index := 0
	for i := range b.setting.Tasks {
		task := b.setting.Tasks[i]
		switch {
		case task.ChallengeArena != nil:
			b.tasks = append(b.tasks, challenge_arena.NewChallengeArenaTask(index, b, *task.ChallengeArena))
		case task.AutoFarm != nil:
			b.tasks = append(b.tasks, auto_farm.NewAutoFarm(index, b, *task.AutoFarm))
		default:
			continue
		}
		index++
	}

	// load task status
	status, err := LoadTaskStatus(StatusFile)
	if err != nil {
		b.log.Errorf("cannot load status from %s: %+v", StatusFile, err)
		b.taskStatuses = make([]any, len(b.tasks))
	} else {
		b.taskStatuses = status.Tasks

		for i := 0; i < len(b.tasks) && i < len(b.taskStatuses); i++ {
			b.tasks[i].LoadStatus(b.taskStatuses[i])
		}

		// add more statuses
		if len(b.taskStatuses) < len(b.tasks) {
			diff := len(b.tasks) - len(b.taskStatuses)
			for i := 0; i < diff; i++ {
				b.taskStatuses = append(b.taskStatuses, nil)
			}
		}

		b.log.Infof("status reloaded")
	}

	b.log.Infof("%d tasks loaded", index)
}

// Run the main loop for bot
func (b *SummonersWar) Run(done <-chan struct{}) {
	oneThirtiethFrameTime := 33 * time.Millisecond
	//
	// b.currentStage = NewStageExpertDarkCastle(b)
	b.SetState(StartState)
	// b.SetState(DoQuest)
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
		switch {
		case b.setting.AreaExploration != nil && b.setting.AreaExploration.Enable:
			b.SetState(StateDoAreaExplorationQuest)
		case b.setting.MonsterStory != nil && b.setting.MonsterStory.Enable:
			b.SetState(StateDoMonsterStoryQuest)
		default:
			b.SetState(StateExecuteTask)
		}

	case StateDoAreaExplorationQuest:
		b.taskAreaExploration.Do(m)

	case StateDoMonsterStoryQuest:
		b.taskMonsterStory.Do(m)

	case StateExecuteTask:
		// check all tasks
		for i := range b.tasks {
			task := b.tasks[i]
			if !task.IsReady() {
				continue
			}
			if b.currentTaskIndex == TaskUnknown {
				b.SetCurrentTask(i)
				break
			}

			if i < b.currentTaskIndex {
				// need to exit the current task
				b.tasks[b.currentTaskIndex].RequestExit()
				b.log.Infof("task '%s' is ready, exit the current task '%s'", task.GetName(), b.tasks[b.currentTaskIndex].GetName())
				b.SetState(StateExitCurrentTask)
			}
		}

		fallthrough

	case StateExitCurrentTask:

		if b.currentTaskIndex == TaskUnknown {
			b.SetState(StateExecuteTask)
			return
		}

		// execute the current task
		b.tasks[b.currentTaskIndex].Do(m)

	case DoQuest:

		// execute task
		switch {
		case b.handleMainScreen(m):
		case b.handleQuestComplete(m):
		case b.handleSleepScreen(m):
		case b.handleAreaExploration(m):
		case b.handleDialog(m):
		}

	case EndState:

	}
}

func (b *SummonersWar) handleMonsterStoryQuestState(m gocv.Mat) bool {
	prefix := "guard_journal"
	log := b.log.Named("monster-story")

	isQuestActive, _ := b.imMatchDefaultInROI(m, roi.ROILeftMenuDetail, prefix, "icon_ongoing_quest")
	if isQuestActive {
		log.Infof("waiting for the quest to finish...")
		waitMs(500)
		return true
	}

	// quest done
	isQuestFinish, ptQuestFinishBtn := b.imMatchDefaultInROI(m, roi.ROILeftMenuDetail, prefix, "icon_finish_quest")
	if isQuestFinish {
		log.Infof("quest finish (%v)", ptQuestFinishBtn)
		b.screen.MouseMoveAndClickByPoint(ptQuestFinishBtn)
		waitMs(1000)
		return true
	}

	return true
}

func (b *SummonersWar) handleMainScreen(m gocv.Mat) bool {
	foundMainScreen, _ := b.MatchInROI(m, roi.ROITopRigthMenuBtn, domain.MatchOption{
		Path: "btn_top_right_menu",
		Th:   0.01,
		// PrintVal: true,
	})
	if !foundMainScreen {
		return false
	}

	isQuestActive, _ := b.MatchInROI(m, roi.ROILeftMenu, domain.MatchOption{
		Path: "btn_quest_active",
		Th:   0.01,
		// PrintVal: true,
		// Normalize: true,
	})
	if !isQuestActive {
		b.log.Infof("tap to active quest")
		b.screen.MouseMoveAndClickByPoint(roi.PtActiveQuest)
		waitMs(1000)
		return true
	}

	switch b.setting.Mode {
	case BotModeStoryQuest:
		canTapToAccept, ptAccept := b.MatchInROI(m, roi.ROILeftMenuDetail, domain.MatchOption{
			Path: "accept_quest",
			// Th:   0.01,
			// PrintVal: true,
			// Normalize: true,
		})
		if canTapToAccept {
			b.log.Infof("tap to accept quest found (%v)", ptAccept)
			b.screen.MouseMoveAndClickByPoint(ptAccept)
			waitMs(500)
			return true
		}
	case BotModeExplorationQuest:
		canAcceptExploration, ptAcceptExploration := b.MatchInROI(m, roi.ROILeftMenuDetail, domain.MatchOption{
			Path: "exploration_quest",
			Th:   0.01,
			// PrintVal: true,
			// Normalize: true,
		})
		if canAcceptExploration {
			b.log.Infof("tap to accept exploration quest found (%v)", ptAcceptExploration)
			b.screen.MouseMoveAndClickByPoint(ptAcceptExploration)
			waitMs(500)
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

	foundTapToClose, _ := b.MatchInROI(m, roi.ROIQuestCompleteTapToClose, domain.MatchOption{
		Path: prefix + ".txt_tab_to_close",
		Th:   0.01,
		// PrintVal: true,
	})
	if foundTapToClose {
		b.log.Infof("quest completed, click anywhere to close")
		b.screen.MouseMoveAndClickByPoint(roi.PtContinue)
		waitMs(1000)
		return true
	}

	foundModal, ptModalBtnOK := b.MatchInROI(m, roi.ROIModalComplete, domain.MatchOption{
		Path: prefix + ".btn_ok_modal",
		// PrintVal: true,
	})
	if foundModal {
		b.log.Infof("modal found, click %v to close", ptModalBtnOK)
		b.screen.MouseMoveAndClickByPoint(ptModalBtnOK)
		waitMs(500)
		return true
	}

	foundExp, _ := b.imMatchDefaultInROI(m, roi.ROIQuestCompleteExp, prefix, "exp")
	mQuestCompleteBtns := m.Region(roi.ROIQuestCompleteBtns)
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
	b.screen.MouseMoveAndClickByPoint(ptFromROIandPt(roi.ROIQuestCompleteBtns, pt))
	waitMs(500)
	return true
}

func (b *SummonersWar) handleSleepScreen(m gocv.Mat) bool {
	found, _ := b.imMatchDefaultInROI(m, roi.ROISleepModeLogo, "sleep_logo")
	if !found {
		return false
	}

	b.screen.MouseDrag(roi.PtSleepModeWakeFrom.X, roi.PtSleepModeWakeFrom.Y, roi.PtSleepModeWakeTo.X, roi.PtSleepModeWakeTo.Y)
	waitMs(3000)
	return true
}

func (b *SummonersWar) handleAreaExploration(m gocv.Mat) (found bool) {
	prefix := "area_exploration"
	found, _ = b.imMatchDefaultInROI(m, roi.ROIAreaExplorationTitle, prefix, "title_area_exploration")
	if !found {
		return
	}

	canAccept, ptAccept := b.imMatchDefaultInROI(m, roi.ROIAreaExplorationBtns, prefix, "btn_accept")
	if canAccept {
		b.log.Infof("accept exploration quest (%v)", ptAccept)
		b.screen.MouseMoveAndClickByPoint(ptAccept)
		waitMs(500)
		return
	}

	isNewQuest, ptNewQuest := b.MatchInROI(m, roi.ROIAreaExplorationNewQuest, domain.MatchOption{
		Path: prefix + ".icon_exclamation",
	})
	if isNewQuest {
		b.log.Infof("choose new exploration quest (%v)", ptNewQuest)
		b.screen.MouseMoveAndClickByPoint(ptNewQuest)
		waitMs(1000)
	}

	return
}

func (b *SummonersWar) handleDialog(m gocv.Mat) bool {
	prefix := "dialog"
	// foundBtnBack, _ := b.imMatchDefaultInROI(m, roi.ROISleepModeLogo, prefix, "btn_back")
	// foundTxtAuto, _ := b.imMatchDefaultInROI(m, roi.ROISleepModeLogo, prefix, "txt_auto")
	foundBtnBack, _ := b.MatchInROI(m, roi.ROIBtnBack, domain.MatchOption{
		Path: prefix + ".btn_back",
		// Th:   0.01,
		// PrintVal: true,
	})
	foundTxtAuto, _ := b.MatchInROI(m, roi.ROITxtAutoAndIcon, domain.MatchOption{
		Path: prefix + ".txt_auto_and_icon",
		// Th:       0.01,
		// PrintVal: true,
	})

	if !(foundTxtAuto && foundBtnBack) {
		return false
	}
	b.log.Infof("dialog detected")
	b.screen.MouseMoveAndClickByPoint(roi.PtContinue)
	waitMs(600)
	return true
}

func waitMs(v int) {
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

func (b *SummonersWar) LoadStatus(index int, key string) any {
	return b.taskStatuses[index]
	// return nil
}

func (b *SummonersWar) SaveStatus(index int, key string, v any) {
	b.taskStatuses[index] = v
	s := TaskStatus{Tasks: b.taskStatuses}
	data, err := yaml.Marshal(s)
	if err != nil {
		b.log.Errorf("cannot marshal yaml: %+v", err)
	}
	err = os.WriteFile(StatusFile, data, 0644)
	if err != nil {
		b.log.Errorf("cannot write status file '%s': %+v", StatusFile, err)
	}
	b.log.Info("save status")
}

func (b *SummonersWar) ExitTask() {
	b.currentTaskIndex = TaskUnknown
}

func (b *SummonersWar) SetCurrentTask(index int) {
	if index >= len(b.tasks) {
		b.log.Panicf("invalid task index %d (len: %d)", index, len(b.tasks))
	}
	b.log.Infof("executing task '%s'", b.tasks[index].GetName())
	b.currentTaskIndex = index
}
