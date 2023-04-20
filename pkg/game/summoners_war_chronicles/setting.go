package summonerswar

import (
	"os"
	"time"

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

	Tasks []TaskSetting
}

type TaskStatus struct {
	Tasks []any `json:"tasks"`
}

func LoadSetting(fileName string) (Setting, error) {
	return laodYAMLFile[Setting](fileName)
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
	RepeatQuest        *TaskRepeatQuestSetting `json:"repeatQuest"`
	TaskChallengeArena *TaskChallengeArena     `json:"challengeArena"`
}

type TaskRepeatQuestSetting struct {
}

type TaskChallengeArena struct {
	Enable bool
	Times  int
}

type TaskBrawlArena struct {
}
