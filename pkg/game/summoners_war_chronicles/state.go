package summonerswar

type BotState string

var (
	UnknownState BotState = "unknown"
	StartState   BotState = "start"
	EndState     BotState = "end"

	StateExecuteTask     BotState = "execute_task"
	StateExitCurrentTask BotState = "exit_current_task"

	StateDoMainStoryQuest       BotState = "do_main_story_quest"
	StateDoAreaExplorationQuest BotState = "do_area_exploration_quest"
	StateDoMonsterStoryQuest    BotState = "do_monster_story_quest"
	StateDoRuneCombination      BotState = "do_rune_combination"
)

const (
	TaskUnknown = -1
)
