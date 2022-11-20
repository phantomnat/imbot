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

type BaseStageOption struct {
	Name  string
	Level StageLevel
}

type BaseStage struct {
	Game *BloonsTD6
	Log  *zap.SugaredLogger

	Name  string
	Level StageLevel

	StepIdx          int
	CurrentStepCount int
	Slots            []image.Point
	Monkeys          []*Monkey
	Strategies       []Strategy
}

var _ Stage = (*BaseStage)(nil)

func NewBaseStage(game *BloonsTD6, name string, level StageLevel, slots []image.Point) BaseStage {
	b := BaseStage{
		Game:  game,
		Log:   zap.S().Named("stage." + level.String() + "." + name),
		Name:  name,
		Level: level,
		Slots: slots,
	}
	b.Reset()
	return b
}

func (b *BaseStage) GetName() string {
	return b.Name
}

func (b *BaseStage) GetLevel() StageLevel {
	return b.Level
}

func (b *BaseStage) Run(m gocv.Mat) {
	if b.StepIdx >= len(b.Strategies) {
		// no strategy, wait idlely
		return
	}

	strategy := b.Strategies[b.StepIdx]
	b.CurrentStepCount++
	switch strategy.Action {
	case ActionTypePlay:
		b.Game.GetScreen().MouseMoveAndClick(ptBtnPlay.X, ptBtnPlay.Y)
		time.Sleep(250 * time.Millisecond)
		b.StepToNextStrategy()

	case ActionTypePlay2X:
		b.Game.GetScreen().MouseMoveAndClick(ptBtnPlay.X, ptBtnPlay.Y)
		time.Sleep(250 * time.Millisecond)
		b.Game.GetScreen().MouseMoveAndClick(ptBtnPlay.X, ptBtnPlay.Y)
		time.Sleep(250 * time.Millisecond)
		b.StepToNextStrategy()

	case ActionTypeBuy:
		if strategy.Slot >= len(b.Slots) {
			b.Log.With("slotIndex", strategy.Slot, "slotLimit", len(b.Slots)).
				Errorf("invalid slot in strategy index %d (%+v)", b.StepIdx, strategy)
			return
		}

		// click on the slot
		// check for info icon
		if b.Monkeys[strategy.Slot] == nil {
			m := strategy.Monkey
			b.Monkeys[strategy.Slot] = &m
		}

		monkey := b.Monkeys[strategy.Slot]
		if b.CurrentStepCount == 1 {
			// scroll
			if monkey.Page == 1 {
				// scroll up
				b.Game.GetScreen().MouseDrag(
					ptsShopDragUp[0].X,
					ptsShopDragUp[0].Y,
					ptsShopDragUp[1].X,
					ptsShopDragUp[1].Y,
				)
				time.Sleep(500 * time.Millisecond)
			} else {
				// scroll down
				b.Game.GetScreen().MouseDrag(
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
		found, _ := b.Game.imMatchDefaultInROI(m, roiMonkeyShopPanel, "monkey", "shop", monkey.Name)
		if found {
			coord := b.Slots[strategy.Slot]
			b.Game.GetScreen().MouseMove(coord.X, coord.Y)
			time.Sleep(100 * time.Millisecond)
			b.Game.GetScreen().KeyTap(monkey.Shortcut)
			time.Sleep(100 * time.Millisecond)
			b.Game.GetScreen().MouseMoveAndClick(coord.X, coord.Y)
			time.Sleep(100 * time.Millisecond)
			if strategy.Camo {
				// TODO: camo
				//time.Sleep(400 * time.Millisecond)
				//b.game.GetScreen().MouseMoveAndClick(coord.X, coord.Y)
				// detect panel
			}
			// check and go to next strategy
			b.StepToNextStrategy()
		}

	case ActionTypeUpgrade:
		coord := b.Slots[strategy.Slot]
		monkey := b.Monkeys[strategy.Slot]

		var monkeyPanel *ROIMonkeyPanel

		allLevelMatched := monkey.Path1Level >= strategy.Path1 && monkey.Path2Level >= strategy.Path2 && monkey.Path3Level >= strategy.Path3

		// check for info icon
		if okL, _ := b.Game.imMatchDefaultInROI(m, roiLeftMonkeyPanel.InfoIcon, "info-icon"); okL {
			if allLevelMatched {
				b.Game.GetScreen().MouseMoveAndClickByPoint(roiLeftMonkeyPanel.PtCloseBth)
				time.Sleep(500 * time.Millisecond)

				b.StepToNextStrategy()
				return
			}
			monkeyPanel = &roiLeftMonkeyPanel
		} else if okR, _ := b.Game.imMatchDefaultInROI(m, roiRightMonkeyPanel.InfoIcon, "info-icon"); okR {
			if allLevelMatched {
				b.Game.GetScreen().MouseMoveAndClickByPoint(roiRightMonkeyPanel.PtCloseBth)
				time.Sleep(500 * time.Millisecond)

				b.StepToNextStrategy()
				return
			}
			monkeyPanel = &roiRightMonkeyPanel
		} else if allLevelMatched {
			b.StepToNextStrategy()
			return
		} else {
			// click on the slot
			b.Game.GetScreen().MouseMoveAndClick(coord.X, coord.Y)
			time.Sleep(500 * time.Millisecond)
			return
		}

		// detect level
		// check for paths level
		// check for upgradable
		if monkey.Path1Level < strategy.Path1 {
			b.upgradeMonkey(m, monkey, 1, strategy.Path1, monkeyPanel.UpgradePath1)
		} else if monkey.Path2Level < strategy.Path2 {
			b.upgradeMonkey(m, monkey, 2, strategy.Path2, monkeyPanel.UpgradePath2)
		} else if monkey.Path3Level < strategy.Path3 {
			b.upgradeMonkey(m, monkey, 3, strategy.Path3, monkeyPanel.UpgradePath3)
		}
	}
}

func (b *BaseStage) upgradeMonkey(
	m gocv.Mat,
	monkey *Monkey,
	pathNo,
	wantedLv int,
	upgradePath ROIMonkeyUpgradePath,
) {
	detectedLv := b.checkUpgradePathLevel(m, upgradePath)
	if detectedLv < wantedLv {
		if b.isPathUpgradable(m, upgradePath.Buyable) {
			b.Game.GetScreen().MouseMoveAndClickByPoint(upgradePath.PtBtnUpgrade)
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
		b.Log.Warnf("get stucked at path (strategy lv: %d, detected lv: %d, monkey lv: %+v)", wantedLv, detectedLv, monkey)
	}
	//s.log.Debugf("path %d strategy lv: %d, detected lv: %d, monkey lv: %+v", pathNo, wantedLv, detectedLv, monkey)
}

func (b *BaseStage) isPathUpgradable(m gocv.Mat, roi image.Rectangle) bool {
	roiUpgradable := m.Region(roi)
	defer roiUpgradable.Close()
	avg := roiUpgradable.Mean()
	return avg.Val2 > thUpgradeBtn.Val2 && avg.Val3 < thUpgradeBtn.Val3
}

func (b *BaseStage) checkUpgradePathLevel(m gocv.Mat, path ROIMonkeyUpgradePath) int {
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

func (b *BaseStage) StepToNextStrategy() {
	b.StepIdx++
	b.CurrentStepCount = 0
	b.Log.Debugf("step to next action -> %d", b.StepIdx)
}

func (b *BaseStage) Reset() {
	b.StepIdx = 0
	b.CurrentStepCount = 0
	b.Monkeys = make([]*Monkey, len(b.Slots))
}

func (b *BaseStage) loadConfigFromFile(filePath string, inSlots []image.Point) (slots []image.Point, strategies []Strategy, err error) {
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
		return nil, nil, errors.New("empty strategy")
	}

	if len(strategiesConfig.Slots) > 0 {
		slots = make([]image.Point, 0, len(strategiesConfig.Slots))
		for i := range strategiesConfig.Slots {
			slots = append(slots, image.Pt(
				strategiesConfig.Slots[i].X,
				strategiesConfig.Slots[i].Y,
			))
		}
	} else if len(slots) == 0 {
		return nil, nil, errors.New("empty slots")
	} else {
		slots = inSlots
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
		if config.Slot >= len(slots) {
			err = errors.Errorf("invalid slot %d, the limit is %d", config.Slot, len(slots)-1)
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

	return slots, strategies, nil
}
