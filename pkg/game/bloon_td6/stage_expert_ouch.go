package bloon_td6

import "image"

type Strategy struct {
	Action ActionType
	Monkey Monkey
	Slot   int
}

func NewStrategy(action ActionType, monkey Monkey, pos image.Point) {

}

type StageExpertOuch struct {
	level      StageLevel
	strategies []Strategy
	slots      []image.Point
	monkeys    []*Monkey
}

var _ Stage = (*StageExpertOuch)(nil)

func NewStageExpertOuch() *StageExpertOuch {
	s := &StageExpertOuch{
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
				Slot:   3,
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
	return s
}
func (s *StageExpertOuch) GetLevel() StageLevel {
	return s.level
}
