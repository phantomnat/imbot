package summonerswar

import (
	"os"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

type BotMode string

var (
	BotModeStoryQuest       BotMode = "storyQuest"
	BotModeExplorationQuest BotMode = "explorationQuest"
	BotModeMonsterStory     BotMode = "monsterStory"
)

type Setting struct {
	Mode BotMode
}

func LoadSetting(fileName string) (Setting, error) {
	s := Setting{}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return Setting{}, errors.WithStack(err)
	}
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		return Setting{}, errors.WithStack(err)
	}
	return s, nil
}

type TaskSetting struct {

}

type TaskRepeatQuestSetting struct {

}