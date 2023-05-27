package summonerswar

import (
	"context"
	"image"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"sigs.k8s.io/yaml"

	"github.com/phantomnat/imbot/pkg/domain"
	area_exploration "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/area_exploration"
	auto_farm "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/auto_farm"
	challenge_arena "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/challenge_arena"
	fishing "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/fishing"
	main_story "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/main_story"
	monster_story "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/monster_story"
	rune_combination "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/rune_combination"
	"github.com/phantomnat/imbot/pkg/im"
	"github.com/phantomnat/imbot/pkg/screen/mumu"
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
	tasksByName      map[string]domain.Task
	taskStatus       TaskStatus
	currentTaskIndex int

	taskMainStory       domain.Task
	taskAreaExploration domain.Task
	taskMonsterStory    domain.Task
	taskFishing         domain.Task
}

var _ domain.Game = (*SummonersWar)(nil)

type Option struct {
	Ctx          context.Context
	ImageManager *im.ImageManager
}

func New(o Option) (*SummonersWar, error) {
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
		im:               o.ImageManager,
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

	// load task status
	status, err := LoadTaskStatus(StatusFile)
	if err != nil {
		b.log.Errorf("cannot load status from %s: %+v", StatusFile, err)
		b.taskStatus = TaskStatus{
			Names: make(map[string]any),
		}
	} else {
		b.taskStatus = status
		b.log.Infof("status reloaded")
	}

	b.InitTasks()

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
	for _, t := range b.tasks {
		if t == nil {
			continue
		}
		t.Reset()
	}
}

func (b *SummonersWar) InitTasks() {
	if b.setting.MainStory != nil {
		b.taskMainStory = main_story.NewMainStory(b, b.setting.MainStory)
	}
	if b.setting.AreaExploration != nil {
		b.taskAreaExploration = area_exploration.NewAreaExploration(0, b, *b.setting.AreaExploration)
	}
	if b.setting.MonsterStory != nil {
		b.taskMonsterStory = monster_story.NewMonsterStory(0, b, *b.setting.MonsterStory)
	}
	if b.setting.Fishing != nil {
		b.taskFishing = fishing.NewFishing(b, b.setting.Fishing)
	}

	b.tasks = make([]domain.Task, 0, len(b.setting.Tasks))
	b.tasksByName = make(map[string]domain.Task, len(b.setting.Tasks))

	index := 0
	for i := range b.setting.Tasks {
		var task domain.Task
		taskSetting := b.setting.Tasks[i]
		switch {
		case taskSetting.ChallengeArena != nil:
			task = challenge_arena.NewChallengeArenaTask(index, b, *taskSetting.ChallengeArena)
		case taskSetting.AutoFarm != nil:
			task = auto_farm.NewAutoFarm(index, b, *taskSetting.AutoFarm)
		case taskSetting.RuneCombination != nil:
			task = rune_combination.NewRuneCombination(b, taskSetting.RuneCombination)
		default:
			continue
		}
		b.tasks = append(b.tasks, task)
		if _, exist := b.tasksByName[task.GetName()]; exist {
			// TODO: handle duplicated task
		}
		b.tasksByName[task.GetName()] = task
		if v, found := b.taskStatus.Names[task.GetName()]; found {
			task.LoadStatus(v)
		}
		index++
	}

	b.log.Infof("%d tasks loaded", index)
}

// Run the main loop for bot
func (b *SummonersWar) Run(done <-chan struct{}) {
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

		// TODO: execute the command if any

		startTime := time.Now()

		b.handleState()

		processedTime := time.Since(startTime)
		if oneThirtiethFrameTime > processedTime {
			time.Sleep(oneThirtiethFrameTime - processedTime)
		}

		// b.log.Debugf("processed time: %v", processedTime)
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
		b.SetState(StateExecuteTask)

	case StateExecuteTask:
		// TODO: special task need to revise late
		doSpecialTask := true
		switch {
		case b.setting.MainStory != nil && b.setting.MainStory.IsEnabled():
			b.taskMainStory.Do(m)
		case b.setting.AreaExploration != nil && b.setting.AreaExploration.IsEnabled():
			b.taskAreaExploration.Do(m)
		case b.setting.MonsterStory != nil && b.setting.MonsterStory.IsEnabled():
			b.taskMonsterStory.Do(m)
		case b.setting.Fishing != nil && b.setting.Fishing.IsEnabled():
			b.taskFishing.Do(m)
		default:
			doSpecialTask = false
		}
		if doSpecialTask {
			return
		}

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
	case EndState:

	}
}

func waitMs(v int) {
	time.Sleep(time.Duration(v) * time.Millisecond)
}

func (b *SummonersWar) SetState(next BotState) {
	b.muCurrentState.Lock()
	defer b.muCurrentState.Unlock()
	if b.currentState == next {
		return
	}
	b.log.Debugf("changing bot state %s -> %s", b.currentState, next)
	b.currentState = next
}

func (b *SummonersWar) GetState() BotState {
	b.muCurrentState.RLock()
	defer b.muCurrentState.RUnlock()

	return b.currentState
}

func (b *SummonersWar) SaveStatus(key string, v any) {
	b.taskStatus.Names[key] = v

	data, err := yaml.Marshal(b.taskStatus)
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
