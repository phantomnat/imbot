package bloon_td6

import (
	"image"

	"gocv.io/x/gocv"
)

type StageLevel int

const (
	StageLevelBeginner StageLevel = iota + 1
	StageLevelIntermediate
	StageLevelAdvanced
	StageLevelExpert
)

func (s StageLevel) String() string {
	switch s {
	case StageLevelBeginner:
		return "beginner"
	case StageLevelIntermediate:
		return "intermediate"
	case StageLevelAdvanced:
		return "advanced"
	case StageLevelExpert:
		return "expert"
	}
	return ""
}

type Stage interface {
	GetName() string
	GetLevel() StageLevel
	Run(m gocv.Mat)
	Reset()
}

var (
	DartMonkey      = NewMonkey("dart-monkey", 1, "q")
	BoomerangMonkey = NewMonkey("boomerang-monkey", 1, "w")
	BombShooter     = NewMonkey("bomb-shooter", 1, "e")
	TackShooter     = NewMonkey("tack-shooter", 1, "r")
	IceMonkey       = NewMonkey("ice-monkey", 1, "t")
	GlueGunner      = NewMonkey("glue-gunner", 1, "y")

	SniperMonkey    = NewMonkey("sniper-monkey", 1, "z")
	MonkeySub       = NewMonkey("monkey-sub", 1, "x")
	MonkeyBuccaneer = NewMonkey("monkey-buccaneer", 1, "c")
	MonkeyAce       = NewMonkey("monkey-ace", 1, "v")
	HeliPilotMonkey = NewMonkey("heli-pilot", 1, "b")
	MortarMonkey    = NewMonkey("mortar-monkey", 2, "n")
	DartlingGunner  = NewMonkey("dartling-gunner", 2, "m")

	WizardMonkey    = NewMonkey("wizard-monkey", 2, "a")
	SuperMonkey     = NewMonkey("super-monkey", 2, "s")
	NinjaMonkey     = NewMonkey("ninja-monkey", 2, "d")
	AlchemistMonkey = NewMonkey("alchemist", 2, "f")
	DruidMonkey     = NewMonkey("druid", 2, "g")

	BananaFarm     = NewMonkey("banana-farm", 2, "h")
	SpikeFactory   = NewMonkey("spike-factory", 2, "j")
	MonkeyVillage  = NewMonkey("monkey-village", 2, "k")
	EngineerMonkey = NewMonkey("engineer-monkey", 2, "l")

	Monkeys = map[string]Monkey{
		DartMonkey.Name:      DartMonkey,
		BoomerangMonkey.Name: BoomerangMonkey,
		BombShooter.Name:     BombShooter,
		TackShooter.Name:     TackShooter,
		IceMonkey.Name:       IceMonkey,
		GlueGunner.Name:      GlueGunner,
		SniperMonkey.Name:    SniperMonkey,
		MonkeySub.Name:       MonkeySub,
		MonkeyBuccaneer.Name: MonkeyBuccaneer,
		MonkeyAce.Name:       MonkeyAce,
		HeliPilotMonkey.Name: HeliPilotMonkey,
		MortarMonkey.Name:    MortarMonkey,
		DartlingGunner.Name:  DartlingGunner,
		WizardMonkey.Name:    WizardMonkey,
		SuperMonkey.Name:     SuperMonkey,
		NinjaMonkey.Name:     NinjaMonkey,
		AlchemistMonkey.Name: AlchemistMonkey,
		DruidMonkey.Name:     DruidMonkey,
		BananaFarm.Name:      BananaFarm,
		SpikeFactory.Name:    SpikeFactory,
		MonkeyVillage.Name:   MonkeyVillage,
		EngineerMonkey.Name:  EngineerMonkey,
	}
)

type ActionType string

const (
	ActionTypeUnknown ActionType = "unknown"
	ActionTypeBuy     ActionType = "buy"
	ActionTypeUpgrade ActionType = "upgrade"
	ActionTypePlay    ActionType = "play"
	ActionTypePlay2X  ActionType = "play-2x"
)

var (
	Actions = map[ActionType]struct{}{
		ActionTypeBuy:     {},
		ActionTypeUpgrade: {},
		ActionTypePlay:    {},
		ActionTypePlay2X:  {},
	}
)

type BaseStage struct {
}

type Monkey struct {
	Name       string
	Page       int
	Shortcut   string
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

func NewMonkey(name string, page int, shortcut string, paths ...int) Monkey {
	m := Monkey{Name: name, Page: page, Shortcut: shortcut}
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

type Strategy struct {
	Action ActionType
	Monkey Monkey
	Slot   int
	Path1  int
	Path2  int
	Path3  int
	Camo   bool
}

type StrategiesConfig struct {
	Strategies []StrategyConfig `yaml:"strategies"`
}
type StrategyConfig struct {
	Action ActionType `yaml:"action"`
	Monkey string     `yaml:"monkey"`
	Slot   int        `yaml:"slot"`
	Path1  int        `yaml:"path1"`
	Path2  int        `yaml:"path2"`
	Path3  int        `yaml:"path3"`
	Camo   bool       `yaml:"camo"`
}

func NewStrategy(action ActionType, monkey Monkey, pos image.Point) {

}
