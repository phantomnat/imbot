package bloon_td6

type StageExpertDarkCastle struct {
	BaseStage
}

func NewStageExpertDarkCastle(game *BloonsTD6) *StageExpertDarkCastle {
	stageName := "dark-castle"
	stageLevel := StageLevelExpert
	configFile := "./configs/bloons_td_6/expert_dark_castle_1.yaml"

	s := &StageExpertDarkCastle{
		BaseStage: NewBaseStage(game, stageName, stageLevel, nil),
	}
	slots, strategies, err := s.loadConfigFromFile(configFile, nil)
	if err != nil {
		s.Log.Errorf("cannot load config file from file: %+v", err)
	} else {
		s.Log.Infof("startegy successfully loaded from %s", configFile)
		s.Strategies = strategies
		s.Slots = slots
	}

	s.Reset()
	return s
}
