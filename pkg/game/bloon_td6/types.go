package bloon_td6

type StageLevel int

const (
	StageLevelBeginner StageLevel = iota + 1
	StageLevelIntermediate
	StageLevelAdvanced
	StageLevelExpert
)

type Stage interface {
	GetLevel() StageLevel
}

var (
	DartMonkey     = NewMonkey("dart_monkey")
	EngineerMonkey = NewMonkey("engineer_monkey")
)

type ActionType string

const (
	ActionTypeUnknown ActionType = "unknown"
	ActionTypeBuy     ActionType = "buy"
	ActionTypeUpgrade ActionType = "upgrade"
)

type BaseStage struct {
}

type Monkey struct {
	Name       string
	Path1Level int
	Path2Level int
	Path3Level int
}

func (m *Monkey) WithLevels(p1, p2, p3 int) *Monkey {
	m.Path1Level = p1
	m.Path2Level = p2
	m.Path3Level = p3
	return m
}

func NewMonkey(name string, paths ...int) Monkey {
	m := Monkey{Name: name}
	if len(paths) > 0 {
		m.Path1Level = paths[0]
	}
	if len(paths) > 1 {
		m.Path1Level = paths[1]
	}
	if len(paths) > 2 {
		m.Path1Level = paths[2]
	}
	return m
}
