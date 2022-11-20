package bloon_td6

import (
	"image"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
	"gopkg.in/yaml.v3"
)

type StageExpertOuch struct {
	game *BloonsTD6
	log  *zap.SugaredLogger

	name  string
	level StageLevel

	slots            []image.Point
	monkeys          []*Monkey
	strategies       []Strategy
	stepIdx          int
	currentStepCount int
	strategyFile     string
}

var _ Stage = (*StageExpertOuch)(nil)

func NewStageExpertOuch(game *BloonsTD6) *StageExpertOuch {
	stageName := "ouch"
	stageLevel := StageLevelExpert
	s := &StageExpertOuch{
		log:   zap.S().Named("stage." + stageLevel.String() + "." + stageName),
		game:  game,
		name:  stageName,
		level: stageLevel,
		strategies: []Strategy{
			{Action: ActionTypeBuy, Monkey: EngineerMonkey, Slot: 5},
			{Action: ActionTypeBuy, Monkey: DartMonkey, Slot: 15},
			{Action: ActionTypePlay2X},
			{Action: ActionTypeBuy, Monkey: DartMonkey, Slot: 7},
			{Action: ActionTypeUpgrade, Slot: 5, Path1: 1},
			{Action: ActionTypeBuy, Monkey: EngineerMonkey, Slot: 10},
			{Action: ActionTypeUpgrade, Slot: 10, Path1: 1},
			{Action: ActionTypeBuy, Slot: 0, Monkey: DruidMonkey},
			{Action: ActionTypeUpgrade, Slot: 0, Path1: 1},
			{Action: ActionTypeBuy, Slot: 6, Monkey: DruidMonkey},
			{Action: ActionTypeUpgrade, Slot: 6, Path1: 1},
			{Action: ActionTypeBuy, Slot: 8, Monkey: DruidMonkey},
			{Action: ActionTypeUpgrade, Slot: 8, Path1: 1},
			{Action: ActionTypeBuy, Slot: 4, Monkey: NinjaMonkey, Camo: true},
			{Action: ActionTypeBuy, Slot: 9, Monkey: NinjaMonkey, Camo: true},
			{Action: ActionTypeBuy, Slot: 11, Monkey: NinjaMonkey, Camo: true},
			{Action: ActionTypeUpgrade, Slot: 4, Path1: 1, Path3: 1},
			{Action: ActionTypeUpgrade, Slot: 9, Path1: 1, Path3: 1},
			{Action: ActionTypeUpgrade, Slot: 11, Path1: 1, Path3: 1},
			{Action: ActionTypeUpgrade, Slot: 4, Path1: 2},
			{Action: ActionTypeUpgrade, Slot: 9, Path1: 2},
			{Action: ActionTypeUpgrade, Slot: 11, Path1: 2}, // ninja 2-1-0
			{Action: ActionTypeUpgrade, Slot: 5, Path1: 2},
			{Action: ActionTypeUpgrade, Slot: 10, Path1: 2},
			{Action: ActionTypeUpgrade, Slot: 5, Path3: 2},
			{Action: ActionTypeUpgrade, Slot: 10, Path3: 2}, // engineer 2-0-2
			{Action: ActionTypeUpgrade, Slot: 5, Path1: 3},
			{Action: ActionTypeUpgrade, Slot: 10, Path1: 3}, // engineer 3-0-2
			{Action: ActionTypeBuy, Slot: 1, Monkey: BombShooter},
			{Action: ActionTypeUpgrade, Slot: 0, Path2: 1},
			{Action: ActionTypeUpgrade, Slot: 6, Path2: 1},
			{Action: ActionTypeUpgrade, Slot: 8, Path2: 1},           // alchemist 1-1-0
			{Action: ActionTypeUpgrade, Slot: 1, Path1: 2, Path2: 3}, // bomb shooter 2-3-0
			{Action: ActionTypeUpgrade, Slot: 4, Path1: 3, Path3: 1}, // ninja 3-0-1
			{Action: ActionTypeUpgrade, Slot: 8, Path1: 3, Path2: 2}, // alchemist 3-1-0
			{Action: ActionTypeUpgrade, Slot: 0, Path1: 2, Path2: 2}, // alchemist 2-1-0
			{Action: ActionTypeUpgrade, Slot: 9, Path1: 3, Path3: 1}, // ninja 3-0-1
		},
		strategyFile: "./configs/bloons_td_6/expert_ouch_2.yaml",
		slots: []image.Point{
			image.Pt(365, 210),
			image.Pt(448, 210),
			image.Pt(679, 210),
			image.Pt(758, 210),

			image.Pt(365, 280),
			image.Pt(454, 280),
			image.Pt(676, 280),
			image.Pt(760, 280),

			image.Pt(365, 444),
			image.Pt(450, 444),
			image.Pt(676, 444),
			image.Pt(760, 444),

			image.Pt(365, 510),
			image.Pt(450, 510),
			image.Pt(676, 510),
			image.Pt(760, 510),
		},
	}
	s.Reset()
	if s.strategyFile != "" {
		// load config from file
		if strategies, err := s.loadConfigFromFile(s.strategyFile); err != nil {
			s.log.Errorf("cannot load config file from file: %+v", err)
		} else {
			s.log.Infof("startegy successfully loaded from %s", s.strategyFile)
			s.strategies = strategies
		}
	}
	return s
}

func (s *StageExpertOuch) loadConfigFromFile(filePath string) (strategies []Strategy, err error) {
	var data []byte
	data, err = os.ReadFile(filePath)
	if err != nil {
		err = errors.Wrapf(err, "read config file %s", filePath)
		return
	}
	var strategiesConfig StrategiesConfig
	err = yaml.Unmarshal(data, &strategiesConfig)
	if err != nil {
		err = errors.Wrapf(err, "parse yaml config file %s", filePath)
		return
	}
	if len(strategiesConfig.Strategies) == 0 {
		return nil, nil
	}

	strategies = make([]Strategy, 0, len(strategiesConfig.Strategies))
	for i := range strategiesConfig.Strategies {
		config := strategiesConfig.Strategies[i]

		if _, exist := Actions[config.Action]; !exist {
			err = errors.Errorf("invalid action %s", config.Action)
			return
		}

		if _, exist := Monkeys[config.Monkey]; config.Action == ActionTypeBuy && !exist {
			err = errors.Errorf("monkey %s not exist", config.Monkey)
			return
		}
		if config.Slot >= len(s.slots) {
			err = errors.Errorf("invalid slot %d, the limit is %d", config.Slot, len(s.slots)-1)
			return
		}

		strategies = append(strategies, Strategy{
			Action: config.Action,
			Monkey: Monkeys[config.Monkey],
			Slot:   config.Slot,
			Path1:  config.Path1,
			Path2:  config.Path2,
			Path3:  config.Path3,
			Camo:   config.Camo,
		})
	}

	return strategies, nil
}

func (s *StageExpertOuch) GetName() string {
	return s.name
}

func (s *StageExpertOuch) GetLevel() StageLevel {
	return s.level
}

func (s *StageExpertOuch) Run(m gocv.Mat) {
	if s.stepIdx >= len(s.strategies) {
		// no strategy, wait idlely
		return
	}

	strategy := s.strategies[s.stepIdx]
	s.currentStepCount++
	switch strategy.Action {
	case ActionTypePlay:
		s.game.GetScreen().MouseMoveAndClick(ptBtnPlay.X, ptBtnPlay.Y)
		time.Sleep(250 * time.Millisecond)
		s.StepToNextStrategy()

	case ActionTypePlay2X:
		s.game.GetScreen().MouseMoveAndClick(ptBtnPlay.X, ptBtnPlay.Y)
		time.Sleep(250 * time.Millisecond)
		s.game.GetScreen().MouseMoveAndClick(ptBtnPlay.X, ptBtnPlay.Y)
		time.Sleep(250 * time.Millisecond)
		s.StepToNextStrategy()

	case ActionTypeBuy:
		if strategy.Slot >= len(s.slots) {
			s.log.With("slotIndex", strategy.Slot, "slotLimit", len(s.slots)).
				Errorf("invalid slot in strategy index %d (%+v)", s.stepIdx, strategy)
			return
		}

		// click on the slot
		// check for info icon
		if s.monkeys[strategy.Slot] == nil {
			m := strategy.Monkey
			s.monkeys[strategy.Slot] = &m
		}

		monkey := s.monkeys[strategy.Slot]
		if s.currentStepCount == 1 {
			// scroll
			if monkey.Page == 1 {
				// scroll up
				s.game.GetScreen().MouseDrag(
					ptsShopDragUp[0].X,
					ptsShopDragUp[0].Y,
					ptsShopDragUp[1].X,
					ptsShopDragUp[1].Y,
				)
				time.Sleep(500 * time.Millisecond)
			} else {
				// scroll down
				s.game.GetScreen().MouseDrag(
					ptsShopDragDown[0].X,
					ptsShopDragDown[0].Y,
					ptsShopDragDown[1].X,
					ptsShopDragDown[1].Y,
				)
				time.Sleep(500 * time.Millisecond)
			}
			return
		}

		// search on the monkey shop panel
		// wait for green
		// buy and click on the slot coord
		found, _ := s.game.imMatchDefaultInROI(m, roiMonkeyShopPanel, "monkey", "shop", monkey.Name)
		if found {
			coord := s.slots[strategy.Slot]
			s.game.GetScreen().MouseMove(coord.X, coord.Y)
			time.Sleep(100 * time.Millisecond)
			s.game.GetScreen().KeyTap(monkey.Shortcut)
			time.Sleep(100 * time.Millisecond)
			s.game.GetScreen().MouseMoveAndClick(coord.X, coord.Y)
			time.Sleep(100 * time.Millisecond)
			if strategy.Camo {
				// TODO: camo
				//time.Sleep(400 * time.Millisecond)
				//s.game.GetScreen().MouseMoveAndClick(coord.X, coord.Y)
				// detect panel
			}
			// check and go to next strategy
			s.StepToNextStrategy()
		}

	case ActionTypeUpgrade:
		coord := s.slots[strategy.Slot]
		monkey := s.monkeys[strategy.Slot]

		var monkeyPanel *ROIMonkeyPanel

		allLevelMatched := monkey.Path1Level >= strategy.Path1 && monkey.Path2Level >= strategy.Path2 && monkey.Path3Level >= strategy.Path3

		// check for info icon
		if okL, _ := s.game.imMatchDefaultInROI(m, roiLeftMonkeyPanel.InfoIcon, "info-icon"); okL {
			if allLevelMatched {
				s.game.GetScreen().MouseMoveAndClickByPoint(roiLeftMonkeyPanel.PtCloseBth)
				time.Sleep(500 * time.Millisecond)

				s.StepToNextStrategy()
				return
			}
			monkeyPanel = &roiLeftMonkeyPanel
		} else if okR, _ := s.game.imMatchDefaultInROI(m, roiRightMonkeyPanel.InfoIcon, "info-icon"); okR {
			if allLevelMatched {
				s.game.GetScreen().MouseMoveAndClickByPoint(roiRightMonkeyPanel.PtCloseBth)
				time.Sleep(500 * time.Millisecond)

				s.StepToNextStrategy()
				return
			}
			monkeyPanel = &roiRightMonkeyPanel
		} else if allLevelMatched {
			s.StepToNextStrategy()
			return
		} else {
			// click on the slot
			s.game.GetScreen().MouseMoveAndClick(coord.X, coord.Y)
			time.Sleep(500 * time.Millisecond)
			return
		}

		// detect level
		// check for paths level
		// check for upgradable
		if monkey.Path1Level < strategy.Path1 {
			s.upgradeMonkey(m, monkey, 1, strategy.Path1, monkeyPanel.UpgradePath1)
		} else if monkey.Path2Level < strategy.Path2 {
			s.upgradeMonkey(m, monkey, 2, strategy.Path2, monkeyPanel.UpgradePath2)
		} else if monkey.Path3Level < strategy.Path3 {
			s.upgradeMonkey(m, monkey, 3, strategy.Path3, monkeyPanel.UpgradePath3)
		}
	}
}

func (s *StageExpertOuch) upgradeMonkey(
	m gocv.Mat,
	monkey *Monkey,
	pathNo,
	wantedLv int,
	upgradePath ROIMonkeyUpgradePath,
) {
	detectedLv := s.checkUpgradePathLevel(m, upgradePath)
	if detectedLv < wantedLv {
		if s.isPathUpgradable(m, upgradePath.Buyable) {
			s.game.GetScreen().MouseMoveAndClickByPoint(upgradePath.PtBtnUpgrade)
			time.Sleep(500 * time.Millisecond)
		}
	} else if detectedLv == wantedLv {
		switch pathNo {
		case 1:
			monkey.Path1Level = wantedLv
		case 2:
			monkey.Path2Level = wantedLv
		default:
			monkey.Path3Level = wantedLv
		}
	} else {
		s.log.Warnf("get stucked at path (strategy lv: %d, detected lv: %d, monkey lv: %+v)", wantedLv, detectedLv, monkey)
	}
	//s.log.Debugf("path %d strategy lv: %d, detected lv: %d, monkey lv: %+v", pathNo, wantedLv, detectedLv, monkey)
}

func (s *StageExpertOuch) isPathUpgradable(m gocv.Mat, roi image.Rectangle) bool {
	roiUpgradable := m.Region(roi)
	defer roiUpgradable.Close()
	avg := roiUpgradable.Mean()
	return avg.Val2 > thUpgradeBtn.Val2 && avg.Val3 < thUpgradeBtn.Val3
}

func (s *StageExpertOuch) checkUpgradePathLevel(m gocv.Mat, path ROIMonkeyUpgradePath) int {
	lv := 0
	for i := 0; i < len(path.Levels); i++ {
		if func() bool {
			mLV := m.Region(path.Levels[i])
			defer mLV.Close()
			avg := mLV.Mean()
			return avg.Val2 > 150
		}() {
			lv++
		} else {
			break
		}
	}
	return lv
}

func (s *StageExpertOuch) StepToNextStrategy() {
	s.stepIdx++
	s.currentStepCount = 0
	s.log.Debugf("step to next action -> %d", s.stepIdx)
}

func (s *StageExpertOuch) Reset() {
	s.stepIdx = 0
	s.currentStepCount = 0
	s.monkeys = make([]*Monkey, len(s.slots))
}
