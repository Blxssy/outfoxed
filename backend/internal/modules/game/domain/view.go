package domain

import "time"

type PlayerView struct {
	UserID    PlayerID `json:"userId"`
	Seat      int      `json:"seat"`
	Name      string   `json:"name"`
	PawnCell  int      `json:"pawnCell"`
	Connected bool     `json:"connected"`
}

type RollView struct {
	RollsUsed int      `json:"rollsUsed"`
	MaxRolls  int      `json:"maxRolls"`
	Faces     []string `json:"faces"`
	Kept      []bool   `json:"kept"`
	Success   bool     `json:"success"`
}

type MoveView struct {
	StepsTotal     int   `json:"stepsTotal"`
	StepsRemaining int   `json:"stepsRemaining"`
	ReachableCells []int `json:"reachableCells,omitempty"`
}

type GameView struct {
	ID             string     `json:"id"`
	Status         GameStatus `json:"status"`
	Phase          GamePhase  `json:"phase"`
	Result         GameResult `json:"result,omitempty"`
	Version        int        `json:"version"`
	Turn           int        `json:"turn"`
	ActiveSeat     int        `json:"activeSeat"`
	TurnDeadlineAt *time.Time `json:"turnDeadlineAt,omitempty"`

	Me       *PlayerView       `json:"me,omitempty"`
	Players  []PlayerView      `json:"players"`
	Board    BoardView         `json:"board"`
	Fox      FoxView           `json:"fox"`
	Suspects []SuspectCardView `json:"suspects"`
	Clues    []ClueTokenView   `json:"clues"`

	Journal []JournalEntry `json:"journal"`

	Move *MoveView `json:"move,omitempty"`
	Roll *RollView `json:"roll,omitempty"`

	AvailableActions []ActionType `json:"availableActions"`
}

func BuildGameView(st GameState, userID PlayerID) GameView {
	meState, ok := findPlayerByID(st.Players, userID)

	var me *PlayerView
	if ok {
		me = &PlayerView{
			UserID:    meState.UserID,
			Seat:      meState.Seat,
			Name:      meState.Name,
			PawnCell:  meState.PawnCell,
			Connected: meState.Connected,
		}
	}

	view := GameView{
		ID:             st.ID,
		Status:         st.Status,
		Phase:          st.Phase,
		Result:         st.Result,
		Version:        st.Version,
		Turn:           st.Turn,
		ActiveSeat:     st.ActiveSeat,
		TurnDeadlineAt: st.TurnDeadlineAt,

		Me:       me,
		Players:  make([]PlayerView, 0, len(st.Players)),
		Board:    buildBoardView(st.Board),
		Fox:      FoxView{Track: st.Fox.Track, EscapeAt: st.Fox.EscapeAt},
		Suspects: make([]SuspectCardView, 0, len(st.Suspects)),
		Clues:    make([]ClueTokenView, 0, len(st.Clues)),

		Journal: copyJournal(st.Journal),

		Move:             buildMoveView(st.TurnState.Move),
		Roll:             buildRollView(st.TurnState.Roll),
		AvailableActions: AvailableActionsFor(st, userID),
	}

	for _, p := range st.Players {
		view.Players = append(view.Players, PlayerView{
			UserID:    p.UserID,
			Seat:      p.Seat,
			Name:      p.Name,
			PawnCell:  p.PawnCell,
			Connected: p.Connected,
		})
	}

	for _, s := range st.Suspects {
		item := SuspectCardView{
			ID:       s.ID,
			Name:     string(s.Code),
			Revealed: s.Revealed,
			Excluded: s.Excluded,
		}

		if s.Revealed {
			traits := s.Traits
			item.Traits = &traits
		}

		view.Suspects = append(view.Suspects, item)
	}

	for _, c := range st.Clues {
		item := ClueTokenView{
			ID:        c.ID,
			Revealed:  c.Revealed,
			BoardCell: c.BoardCell,
		}

		if c.Revealed {
			trait := c.Trait
			item.Trait = &trait
			item.Result = c.Result
		}

		view.Clues = append(view.Clues, item)
	}

	return view
}

func buildBoardView(board BoardState) BoardView {
	out := BoardView{
		Cells: make([]BoardCellView, 0, len(board.Cells)),
	}

	for _, c := range board.Cells {
		out.Cells = append(out.Cells, BoardCellView{
			Index:       c.Index,
			Type:        c.Type,
			HasClue:     c.ClueTokenID != "",
			ClueTokenID: c.ClueTokenID,
		})
	}

	return out
}

func buildMoveView(move *MoveState) *MoveView {
	if move == nil {
		return nil
	}

	cells := make([]int, len(move.ReachableCells))
	copy(cells, move.ReachableCells)

	return &MoveView{
		StepsTotal:     move.StepsTotal,
		StepsRemaining: move.StepsRemaining,
		ReachableCells: cells,
	}
}

func buildRollView(roll *RollState) *RollView {
	if roll == nil {
		return nil
	}

	faces := make([]string, len(roll.Faces))
	copy(faces, roll.Faces)

	kept := make([]bool, len(roll.Kept))
	copy(kept, roll.Kept)

	return &RollView{
		RollsUsed: roll.RollsUsed,
		MaxRolls:  roll.MaxRolls,
		Faces:     faces,
		Kept:      kept,
		Success:   roll.Success,
	}
}

func findPlayerByID(players []PlayerState, userID PlayerID) (PlayerState, bool) {
	for _, p := range players {
		if p.UserID == userID {
			return p, true
		}
	}
	return PlayerState{}, false
}

func AvailableActionsFor(st GameState, userID PlayerID) []ActionType {
	if st.Status != StatusActive {
		return nil
	}

	activePlayer, ok := st.ActivePlayer()
	if !ok || activePlayer.UserID != userID {
		return nil
	}

	switch st.Phase {
	case PhaseChooseGoal:
		return []ActionType{
			ActionChooseGoal,
			ActionAccuse,
		}

	case PhaseRolling:
		return []ActionType{
			ActionRerollDice,
			ActionFinishRoll,
			ActionAccuse,
		}

	case PhaseMovePawn:
		return []ActionType{
			ActionMovePawn,
			ActionAccuse,
		}

	case PhaseResolveClue:
		return []ActionType{
			ActionTakeClue,
			ActionAccuse,
		}

	case PhaseRevealSuspects:
		return []ActionType{
			ActionRevealSuspects,
			ActionAccuse,
		}

	case PhaseEndTurn:
		return []ActionType{
			ActionEndTurn,
			ActionAccuse,
		}

	default:
		return nil
	}
}

func copyJournal(items []JournalEntry) []JournalEntry {
	if len(items) == 0 {
		return []JournalEntry{}
	}

	out := make([]JournalEntry, len(items))
	copy(out, items)
	return out
}
