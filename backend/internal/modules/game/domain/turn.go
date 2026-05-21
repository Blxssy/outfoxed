package domain

type TurnState struct {
	Goal    TurnGoal      `json:"goal"`
	Pending PendingAction `json:"pending"`

	Roll *RollState `json:"roll,omitempty"`
	Move *MoveState `json:"move,omitempty"`
}

type TurnGoal struct {
	Set  bool     `json:"set"`
	Type GoalType `json:"type"`
}

type RollState struct {
	RollsUsed int      `json:"rollsUsed"`
	MaxRolls  int      `json:"maxRolls"`
	Faces     []string `json:"faces"`
	Kept      []bool   `json:"kept"`
	Success   bool     `json:"success"`
}

type MoveState struct {
	StepsTotal     int   `json:"stepsTotal"`
	StepsRemaining int   `json:"stepsRemaining"`
	ReachableCells []int `json:"reachableCells,omitempty"`
}

func NewTurnState() TurnState {
	return TurnState{
		Goal:    TurnGoal{},
		Pending: PendingNone,
		Roll:    nil,
		Move:    nil,
	}
}

func (t *TurnState) ResetForNextTurn() {
	t.Goal = TurnGoal{}
	t.Pending = PendingNone
	t.Roll = nil
	t.Move = nil
}
