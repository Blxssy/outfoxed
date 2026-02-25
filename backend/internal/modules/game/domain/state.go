package domain

type PendingAction string

const (
	PendingNone    PendingAction = "none"
	PendingClue    PendingAction = "clue"
	PendingSuspect PendingAction = "suspect"
)

type PlayerState struct {
	ID       PlayerID `json:"id"`
	Seat     int      `json:"seat"` // 0..3
	Position int      `json:"position"`
}

type GameState struct {
	ID      GameID     `json:"id"`
	Status  GameStatus `json:"status"`
	Phase   Phase      `json:"phase"`
	Turn    int        `json:"turn"`    // номер хода
	Version int        `json:"version"` // версия состояния

	Players    []PlayerState `json:"players"`
	ActiveSeat int           `json:"active_seat"` // чей ход

	FoxTrack int `json:"fox_track"` // позиция лиса на дорожке
	Goal     struct {
		Type GoalType `json:"type"`
		Set  bool     `json:"set"`
	} `json:"goal"`

	CluesFound int `json:"clues_found"`
	CluesTotal int `json:"clues_total"`

	// Что нужно сделать после успешного броска
	Pending PendingAction `json:"pending"`

	Suspects []SuspectState
	// Clues    []ClueState

	FoxEscapeAt int        `json:"foxEscapeAt"` // порог, когда лис убежал
	CulpritID   int        `json:"culpritId"`   // кто виновник
	Result      GameResult `json:"result"`
}

type GameResult string

const (
	ResultNone GameResult = "none"
	ResultWin  GameResult = "win"
	ResultLose GameResult = "lose"
)

func (gs GameState) ActivePlayer() (PlayerState, bool) {
	for _, p := range gs.Players {
		if p.Seat == gs.ActiveSeat {
			return p, true
		}
	}
	return PlayerState{}, false
}
