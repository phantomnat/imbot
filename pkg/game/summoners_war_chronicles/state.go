package summonerswar

type BotState string

var (
	UnknownState BotState = "unknown"
	StartState   BotState = "start"
	EndState     BotState = "end"

	ActivateQuest BotState = "activate-quest"
	DoQuest BotState = "do-quest"
)
