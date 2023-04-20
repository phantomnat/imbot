package summonerswar

type BotState string

var (
	UnknownState BotState = "unknown"
	StartState   BotState = "start"
	EndState     BotState = "end"

	DoQuest BotState = "do-quest"

	StateExecuteTask     BotState = "execute_task"
	StateExitCurrentTask BotState = "exit_current_task"
)

const (
	TaskUnknown = -1
)
