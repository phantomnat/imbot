package summonerswar

import (
	"os"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	area_exploration "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/area_exploration"
	auto_farm "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/auto_farm"
	challenge_arena "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/challenge_arena"
	fishing "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/fishing"
	main_story "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/main_story"
	monster_story "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/monster_story"
	rune_combination "github.com/phantomnat/imbot/pkg/game/summoners_war_chronicles/tasks/rune_combination"
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
	EmuRedfinger     EmuType = "redfinger"

	SupportedEmulators = map[EmuType]struct{}{
		EmuTypeBlueStack: {},
		EmuTypeMumu:      {},
		EmuRedfinger:     {},
	}
)

type Setting struct {
	Emu  EmuType
	Mode BotMode

	MainStory       *main_story.TaskSetting
	AreaExploration *area_exploration.TaskSetting
	MonsterStory    *monster_story.TaskSetting
	Fishing         *fishing.TaskSetting

	Tasks []TaskSetting
}

type TaskSetting struct {
	ChallengeArena  *challenge_arena.TaskSetting
	AutoFarm        *auto_farm.TaskSetting
	RuneCombination *rune_combination.TaskSetting
}

func LoadSetting(fileName string) (Setting, error) {
	s, err := laodYAMLFile[Setting](fileName)
	if err != nil {
		return s, err
	}

	// validate apply default
	if _, found := SupportedEmulators[s.Emu]; !found {
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

type TaskStatus struct {
	Tasks []any `json:"tasks"`
	Names map[string]any
}

func LoadTaskStatus(fileName string) (TaskStatus, error) {
	return laodYAMLFile[TaskStatus](fileName)
}
