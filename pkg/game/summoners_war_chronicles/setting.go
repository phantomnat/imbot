package summonerswar

import (
	"os"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	area_exploration "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/area_exploration"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/auto_farm"
	"github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/challenge_arena"
	monster_story "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/monster_story"
)

type BotMode string

var (
	BotModeStoryQuest       BotMode = "storyQuest"
	BotModeExplorationQuest BotMode = "explorationQuest"
	BotModeMonsterStory     BotMode = "monsterStory"
)

type EmuType string

var (
	EmuTypeBlueStack EmuType = "bluestack"
	EmuTypeMumu      EmuType = "mumu"
)

type Setting struct {
	Emu  EmuType
	Mode BotMode

	AreaExploration *area_exploration.TaskSetting
	MonsterStory    *monster_story.TaskSetting

	Tasks []TaskSetting
}

type TaskStatus struct {
	Tasks []any `json:"tasks"`
}

func LoadSetting(fileName string) (Setting, error) {
	s, err := laodYAMLFile[Setting](fileName)
	if err != nil {
		return s, err
	}

	// validate apply default
	switch s.Emu {
	case EmuTypeBlueStack, EmuTypeMumu:
	default:
		s.Emu = EmuTypeMumu
	}

	return s, nil
}

func laodYAMLFile[T any](fileName string) (T, error) {
	var v T
	data, err := os.ReadFile(fileName)
	if err != nil {
		return v, errors.WithStack(err)
	}
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		return v, errors.WithStack(err)
	}
	return v, nil
}

func LoadTaskStatus(fileName string) (TaskStatus, error) {
	return laodYAMLFile[TaskStatus](fileName)
}

type TaskSetting struct {
	ChallengeArena *challenge_arena.TaskSetting
	AutoFarm       *auto_farm.TaskSetting
}
