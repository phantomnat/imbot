package summonerswar

type BotState string

var (
	UnknownState BotState = "unknown"
	StartState   BotState = "start"
	EndState     BotState = "end"

	DoQuest BotState = "do-quest"

	StateExecuteTask     BotState = "execute_task"
	StateExitCurrentTask BotState = "exit_current_task"

	StateDoAreaExplorationQuest BotState = "do_area_exploration_quest"
)

const (
	TaskUnknown = -1
)
