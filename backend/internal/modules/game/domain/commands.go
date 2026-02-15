package domain

type CommandType string

const (
	CmdChooseGoal     CommandType = "choose_goal"
	CmdRollAuto       CommandType = "roll_auto" // MVP: сервер сам делает до 3 попыток
	CmdEndTurn        CommandType = "end_turn"
	CmdTakeClue       CommandType = "take_clue"
	CmdRevealSuspects CommandType = "reveal_suspects"
)

type Command interface {
	Type() CommandType
	Actor() PlayerID
}

// ChooseGoalCommand: игрок выбирает цель хода (подсказка или подозреваемые)
type ChooseGoalCommand struct {
	Player PlayerID `json:"player"`
	Goal   GoalType `json:"goal"`
}

func (c ChooseGoalCommand) Type() CommandType { return CmdChooseGoal }
func (c ChooseGoalCommand) Actor() PlayerID   { return c.Player }

// RollAutoCommand: MVP упрощение. Сервер сам делает до 3 “попыток” и решает успех/неуспех
type RollAutoCommand struct {
	Player PlayerID `json:"player"`
}

func (c RollAutoCommand) Type() CommandType { return CmdRollAuto }
func (c RollAutoCommand) Actor() PlayerID   { return c.Player }

// EndTurnCommand: завершить ход
type EndTurnCommand struct {
	Player PlayerID `json:"player"`
}

func (c EndTurnCommand) Type() CommandType { return CmdEndTurn }
func (c EndTurnCommand) Actor() PlayerID   { return c.Player }

// TakeClueCommand взять подсказку
type TakeClueCommand struct {
	Player PlayerID `json:"player"`
}

func (c TakeClueCommand) Type() CommandType { return CmdTakeClue }
func (c TakeClueCommand) Actor() PlayerID   { return c.Player }

type RevealSuspectsCommand struct {
	Player PlayerID `json:"player"`
}

func (c RevealSuspectsCommand) Type() CommandType { return CmdRevealSuspects }
func (c RevealSuspectsCommand) Actor() PlayerID   { return c.Player }
