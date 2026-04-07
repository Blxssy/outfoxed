package domain

type PlayerState struct {
	UserID    PlayerID `json:"userId"`
	Seat      int      `json:"seat"`
	Name      string   `json:"name"`
	PawnCell  int      `json:"pawnCell"`
	Connected bool     `json:"connected"`
}

type GameState struct {
	ID         string     `json:"id"`
	Status     GameStatus `json:"status"`
	Phase      GamePhase  `json:"phase"`
	Result     GameResult `json:"result"`
	Version    int        `json:"version"`
	Turn       int        `json:"turn"`
	ActiveSeat int        `json:"activeSeat"`

	Players []PlayerState `json:"players"`

	Board BoardState `json:"board"`
	Fox   FoxState   `json:"fox"`

	Suspects []SuspectCard `json:"suspects"`
	Clues    []ClueToken   `json:"clues"`

	TurnState TurnState   `json:"turnState"`
	Secret    SecretState `json:"-"`
}

type SecretState struct {
	CulpritSuspectID string                `json:"-"`
	ClueTruth        map[string]TraitValue `json:"-"`
}

func (gs GameState) ActivePlayer() (PlayerState, bool) {
	for _, p := range gs.Players {
		if p.Seat == gs.ActiveSeat {
			return p, true
		}
	}
	return PlayerState{}, false
}

func (gs GameState) FindPlayer(userID string) (PlayerState, bool) {
	for _, p := range gs.Players {
		if string(p.UserID) == userID {
			return p, true
		}
	}
	return PlayerState{}, false
}

func (gs GameState) FindSuspectByID(id string) (*SuspectCard, bool) {
	for i := range gs.Suspects {
		if gs.Suspects[i].ID == id {
			return &gs.Suspects[i], true
		}
	}
	return nil, false
}

func (gs GameState) FindClueByID(id string) (*ClueToken, bool) {
	for i := range gs.Clues {
		if gs.Clues[i].ID == id {
			return &gs.Clues[i], true
		}
	}
	return nil, false
}
