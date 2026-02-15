package domain

type EventType string

const (
	EvGoalChosen EventType = "goal_chosen"
	EvRolled     EventType = "rolled"
	EvFoxMoved   EventType = "fox_moved"
	EvTurnEnded  EventType = "turn_ended"
	EvClueTaken  EventType = "clue_taken"
)

type Event struct {
	Type EventType      `json:"type"`
	Data map[string]any `json:"data,omitempty"`
}
