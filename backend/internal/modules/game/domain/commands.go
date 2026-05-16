package domain

type CommandType string

const (
	CmdChooseGoal     CommandType = "choose_goal"
	CmdRerollDice     CommandType = "reroll_dice"
	CmdFinishRoll     CommandType = "finish_roll"
	CmdMovePawn       CommandType = "move_pawn"
	CmdTakeClue       CommandType = "take_clue"
	CmdRevealSuspects CommandType = "reveal_suspects"
	CmdEndTurn        CommandType = "end_turn"
	CmdAccuse         CommandType = "accuse"
	CmdRollAuto       CommandType = "roll_auto" // можно оставить для timeout/автохода
)

type Command interface {
	Type() CommandType
	Actor() PlayerID
}

type PlayerID string

func (id PlayerID) String() string {
	return string(id)
}

// ChooseGoalCommand: игрок выбирает цель хода (подсказка или подозреваемые)
type ChooseGoalCommand struct {
	Player PlayerID `json:"player"`
	Goal   GoalType `json:"goal"`
}

func (c ChooseGoalCommand) Type() CommandType { return CmdChooseGoal }
func (c ChooseGoalCommand) Actor() PlayerID   { return c.Player }

// RollAutoCommand: MVP-упрощение. Сервер сам делает до 3 попыток
// и решает успех/неуспех.
type RollAutoCommand struct {
	Player PlayerID `json:"player"`
}

func (c RollAutoCommand) Type() CommandType { return CmdRollAuto }
func (c RollAutoCommand) Actor() PlayerID   { return c.Player }

// MovePawnCommand: игрок тратит часть или все доступные шаги.
type MovePawnCommand struct {
	Player      PlayerID `json:"player"`
	TargetIndex int      `json:"targetIndex"`
}

func (c MovePawnCommand) Type() CommandType { return CmdMovePawn }
func (c MovePawnCommand) Actor() PlayerID   { return c.Player }

// EndTurnCommand: завершить ход
type EndTurnCommand struct {
	Player PlayerID `json:"player"`
}

func (c EndTurnCommand) Type() CommandType { return CmdEndTurn }
func (c EndTurnCommand) Actor() PlayerID   { return c.Player }

// TakeClueCommand: взять подсказку на текущей клетке
type TakeClueCommand struct {
	Player PlayerID `json:"player"`
}

func (c TakeClueCommand) Type() CommandType { return CmdTakeClue }
func (c TakeClueCommand) Actor() PlayerID   { return c.Player }

// RevealSuspectsCommand: открыть 2 выбранные карты подозреваемых
type RevealSuspectsCommand struct {
	Player     PlayerID `json:"player"`
	SuspectIDs []string `json:"suspectIds"`
}

func (c RevealSuspectsCommand) Type() CommandType { return CmdRevealSuspects }
func (c RevealSuspectsCommand) Actor() PlayerID   { return c.Player }

// AccuseCommand: обвинить конкретного подозреваемого
type AccuseCommand struct {
	Player    PlayerID `json:"player"`
	SuspectID string   `json:"suspectId"`
}

func (c AccuseCommand) Type() CommandType { return CmdAccuse }
func (c AccuseCommand) Actor() PlayerID   { return c.Player }

type RerollDiceCommand struct {
	Player      PlayerID `json:"player"`
	KeepIndices []int    `json:"keepIndices"`
}

func (c RerollDiceCommand) Type() CommandType { return CmdRerollDice }
func (c RerollDiceCommand) Actor() PlayerID   { return c.Player }

type FinishRollCommand struct {
	Player PlayerID `json:"player"`
}

func (c FinishRollCommand) Type() CommandType { return CmdFinishRoll }
func (c FinishRollCommand) Actor() PlayerID   { return c.Player }
