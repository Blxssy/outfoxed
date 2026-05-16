package domain

type GameStatus string

const (
	StatusWaiting  GameStatus = "waiting"
	StatusActive   GameStatus = "active"
	StatusFinished GameStatus = "finished"
)

type GamePhase string

const (
	PhaseChooseGoal     GamePhase = "choose_goal"
	PhaseRolling        GamePhase = "rolling"
	PhaseMovePawn       GamePhase = "move_pawn"
	PhaseResolveClue    GamePhase = "resolve_clue"
	PhaseRevealSuspects GamePhase = "reveal_suspects"
	PhaseAccusation     GamePhase = "accusation"
	PhaseEndTurn        GamePhase = "end_turn"
)

type GoalType string

const (
	GoalClue    GoalType = "clue"
	GoalSuspect GoalType = "suspect"
)

type PendingAction string

const (
	PendingNone           PendingAction = ""
	PendingMoveToClue     PendingAction = "move_to_clue"
	PendingResolveClue    PendingAction = "resolve_clue"
	PendingRevealSuspects PendingAction = "reveal_suspects"
)

type GameResult string

const (
	ResultNone GameResult = "none"
	ResultWin  GameResult = "win"
	ResultLose GameResult = "lose"
)

type ActionType string

const (
	ActionChooseGoal     ActionType = "choose_goal"
	ActionRerollDice     ActionType = "reroll_dice"
	ActionFinishRoll     ActionType = "finish_roll"
	ActionRollAuto       ActionType = "roll_auto"
	ActionMovePawn       ActionType = "move_pawn"
	ActionTakeClue       ActionType = "take_clue"
	ActionRevealSuspects ActionType = "reveal_suspects"
	ActionAccuse         ActionType = "accuse"
	ActionEndTurn        ActionType = "end_turn"
)

type ClueTrait string

const (
	ClueTraitGlasses  ClueTrait = "glasses"
	ClueTraitHat      ClueTrait = "hat"
	ClueTraitScarf    ClueTrait = "scarf"
	ClueTraitUmbrella ClueTrait = "umbrella"
	ClueTraitBag      ClueTrait = "bag"
	ClueTraitBoots    ClueTrait = "boots"
)

type TraitValue string

const (
	TraitUnknown TraitValue = "unknown"
	TraitYes     TraitValue = "yes"
	TraitNo      TraitValue = "no"
)

type SuspectCode string
