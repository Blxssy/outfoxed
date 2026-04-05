package domain

func BuildGameView(st GameState, userID string) GameView {
	me, _ := st.FindPlayer(userID)

	view := GameView{
		ID:         st.ID,
		Status:     st.Status,
		Phase:      st.Phase,
		Result:     st.Result,
		Version:    st.Version,
		Turn:       st.Turn,
		ActiveSeat: st.ActiveSeat,
		Me: PlayerView{
			UserID:    me.UserID,
			Seat:      me.Seat,
			Name:      me.Name,
			PawnCell:  me.PawnCell,
			Connected: me.Connected,
		},
		Players:          make([]PlayerView, 0, len(st.Players)),
		Board:            buildBoardView(st.Board),
		Fox:              FoxView{Track: st.Fox.Track, EscapeAt: st.Fox.EscapeAt},
		Suspects:         make([]SuspectCardView, 0, len(st.Suspects)),
		Clues:            make([]ClueTokenView, 0, len(st.Clues)),
		TurnState:        TurnView(st.TurnState),
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

func AvailableActionsFor(st GameState, userID string) []ActionType {
	player, ok := st.ActivePlayer()
	if !ok || player.UserID != userID || st.Status != StatusActive {
		return nil
	}

	switch st.Phase {
	case PhaseChooseGoal:
		return []ActionType{ActionChooseGoal, ActionAccuse}
	case PhaseRolling:
		return []ActionType{ActionRollAuto, ActionAccuse}
	case PhaseMovePawn:
		return []ActionType{ActionMovePawn, ActionAccuse}
	case PhaseResolveClue:
		return []ActionType{ActionTakeClue, ActionAccuse}
	case PhaseRevealSuspects:
		return []ActionType{ActionRevealSuspects, ActionAccuse}
	case PhaseEndTurn:
		return []ActionType{ActionEndTurn, ActionAccuse}
	default:
		return []ActionType{ActionAccuse}
	}
}
