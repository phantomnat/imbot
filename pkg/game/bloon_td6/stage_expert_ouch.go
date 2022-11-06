package bloon_td6

import (
	"image"

	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

type StageExpertOuch struct {
	level            StageLevel
	slots            []image.Point
	monkeys          []*Monkey
	strategies       []Strategy
	stepIdx          int
	currentStepCount int
	log              *zap.SugaredLogger
}

var _ Stage = (*StageExpertOuch)(nil)

func NewStageExpertOuch() *StageExpertOuch {
	s := &StageExpertOuch{
		log:   zap.S().Named("stage.expert.ouch"),
		level: StageLevelExpert,
		strategies: []Strategy{
			{
				Action: ActionTypeBuy,
				Monkey: EngineerMonkey,
				Slot:   5,
			},
			{
				Action: ActionTypeBuy,
				Monkey: DartMonkey,
				Slot:   15,
			},
			{
				Action: ActionTypeBuy,
				Monkey: DartMonkey,
				Slot:   7,
			},
		},
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
	return s
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
	case ActionTypeBuy:
		if strategy.Slot >= len(s.slots) {
			s.log.With("slotIndex", strategy.Slot, "slotLimit", len(s.slots)).
				Errorf("invalid slot in strategy index %d (%+v)", s.stepIdx, strategy)
		}

		if s.monkeys[strategy.Slot] == nil {
			m := strategy.Monkey
			s.monkeys[strategy.Slot] = &m
		}

		// check and go to next strategy
		// click on the slot
		// check for info icon

		// search on the monkey shop panel
		// wait for green
		// buy and click on the slot coord
	case ActionTypeUpgrade:
		// click on the slot
		// check for info icon
		// check for paths level
		// check for upgradable
	}
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
