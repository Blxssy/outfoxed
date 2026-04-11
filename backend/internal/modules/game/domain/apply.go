package domain

func Apply(s GameState, cmd Command, rng RNG) (GameState, []Event, error) {
	if s.Status == StatusFinished {
		return s, nil, ErrGameFinished
	}
	if s.Status != StatusActive {
		return s, nil, ErrGameNotActive
	}

	activePlayer, ok := s.ActivePlayer()
	if !ok || string(activePlayer.UserID) != string(cmd.Actor()) {
		return s, nil, ErrNotYourTurn
	}

	switch c := cmd.(type) {
	case ChooseGoalCommand:
		return applyChooseGoal(s, c)

	case RollAutoCommand:
		return applyRollAuto(s, c, rng)

	case MovePawnCommand:
		return applyMovePawn(s, c)

	case TakeClueCommand:
		return applyTakeClue(s, c)

	case RevealSuspectsCommand:
		return applyRevealSuspects(s, c)

	case AccuseCommand:
		return applyAccuse(s, c)

	case EndTurnCommand:
		return applyEndTurn(s, c)

	default:
		return s, nil, ErrInvalidPhase
	}
}

func applyChooseGoal(s GameState, c ChooseGoalCommand) (GameState, []Event, error) {
	if s.Phase != PhaseChooseGoal {
		return s, nil, ErrInvalidPhase
	}
	if s.TurnState.Goal.Set {
		return s, nil, ErrGoalAlreadySet
	}

	s.TurnState.Goal = TurnGoal{
		Set:  true,
		Type: c.Goal,
	}
	s.TurnState.Pending = PendingNone
	s.TurnState.Roll = nil
	s.TurnState.Move = nil

	s.Phase = PhaseRolling
	s.Version++

	return s, []Event{
		{
			Type: EvGoalChosen,
			Data: map[string]any{
				"goal": c.Goal,
			},
		},
	}, nil
}

func applyRollAuto(s GameState, c RollAutoCommand, rng RNG) (GameState, []Event, error) {
	if s.Phase != PhaseRolling {
		return s, nil, ErrInvalidPhase
	}
	if !s.TurnState.Goal.Set {
		return s, nil, ErrGoalNotSet
	}

	res := RollForGoal(s.TurnState.Goal.Type, rng)

	s.TurnState.Roll = &RollState{
		Attempts: res.Attempts,
		Faces:    res.Faces,
		Success:  res.Success,
	}

	events := []Event{
		{
			Type: EvRolled,
			Data: map[string]any{
				"success":  res.Success,
				"goal":     res.Goal,
				"attempts": res.Attempts,
				"faces":    res.Faces,
			},
		},
	}

	if !res.Success {
		s.TurnState.Pending = PendingNone
		s.TurnState.Move = nil

		s.Fox.Track += 3

		events = append(events, Event{
			Type: EvFoxMoved,
			Data: map[string]any{
				"by":    3,
				"track": s.Fox.Track,
			},
		})

		if s.Fox.EscapeAt > 0 && s.Fox.Track >= s.Fox.EscapeAt {
			s.Status = StatusFinished
			s.Result = ResultLose
			s.Phase = PhaseEndTurn
			s.Version++

			events = append(events, Event{
				Type: EvGameFinished,
				Data: map[string]any{
					"result": s.Result,
				},
			})
			return s, events, nil
		}

		s.Phase = PhaseEndTurn
		s.Version++
		return s, events, nil
	}

	switch s.TurnState.Goal.Type {
	case GoalClue:
		steps := countMoveSteps(res.Faces)
		if steps <= 0 {
			return s, nil, ErrInvalidMove
		}

		s.TurnState.Pending = PendingMoveToClue
		s.TurnState.Move = &MoveState{
			StepsTotal:     steps,
			StepsRemaining: steps,
		}
		s.Phase = PhaseMovePawn

	case GoalSuspect:
		s.TurnState.Pending = PendingRevealSuspects
		s.TurnState.Move = nil
		s.Phase = PhaseRevealSuspects

	default:
		return s, nil, ErrInvalidPhase
	}

	s.Version++
	return s, events, nil
}

func applyMovePawn(s GameState, c MovePawnCommand) (GameState, []Event, error) {
	if s.Phase != PhaseMovePawn {
		return s, nil, ErrInvalidPhase
	}
	if s.TurnState.Pending != PendingMoveToClue {
		return s, nil, ErrNoPendingAction
	}
	if s.TurnState.Move == nil {
		return s, nil, ErrInvalidMove
	}

	player, ok := s.ActivePlayer()
	if !ok {
		return s, nil, ErrInvalidMove
	}

	move := s.TurnState.Move
	if c.Steps <= 0 || c.Steps > move.StepsRemaining {
		return s, nil, ErrInvalidMove
	}

	newCell := s.Board.ClampIndex(player.PawnCell + c.Steps)

	for i := range s.Players {
		if s.Players[i].Seat == s.ActiveSeat {
			s.Players[i].PawnCell = newCell
			break
		}
	}

	move.StepsRemaining -= c.Steps

	cell, ok := s.Board.CellAt(newCell)
	if !ok {
		return s, nil, ErrInvalidMove
	}

	events := []Event{
		{
			Type: EvPawnMoved,
			Data: map[string]any{
				"seat":      s.ActiveSeat,
				"toCell":    newCell,
				"stepsLeft": move.StepsRemaining,
			},
		},
	}

	// Дошли до клетки с подсказкой — можно брать улику сразу.
	if cell.Type == BoardCellClue && cell.ClueTokenID != "" {
		s.TurnState.Pending = PendingResolveClue
		s.Phase = PhaseResolveClue
		s.Version++
		return s, events, nil
	}

	// Шаги закончились, но до улики не дошли — ход заканчивается.
	if move.StepsRemaining == 0 {
		s.TurnState.Pending = PendingNone
		s.TurnState.Move = nil
		s.Phase = PhaseEndTurn
		s.Version++
		return s, events, nil
	}

	s.Version++
	return s, events, nil
}

func applyTakeClue(s GameState, c TakeClueCommand) (GameState, []Event, error) {
	if s.Phase != PhaseResolveClue {
		return s, nil, ErrInvalidPhase
	}
	if s.TurnState.Pending != PendingResolveClue {
		return s, nil, ErrNoPendingAction
	}

	player, ok := s.ActivePlayer()
	if !ok {
		return s, nil, ErrInvalidPhase
	}

	cell, ok := s.Board.CellAt(player.PawnCell)
	if !ok {
		return s, nil, ErrInvalidMove
	}
	if cell.Type != BoardCellClue || cell.ClueTokenID == "" {
		return s, nil, ErrNoPendingAction
	}

	clue, ok := findClueByID(s.Clues, cell.ClueTokenID)
	if !ok {
		return s, nil, ErrNoPendingAction
	}
	if clue.Revealed {
		return s, nil, ErrAllCluesCollected
	}

	result, ok := s.Secret.ClueTruth[clue.ID]
	if !ok {
		return s, nil, ErrInvalidPhase
	}

	clue.Revealed = true
	clue.Result = ptrTraitValue(result)

	s.TurnState.Pending = PendingNone
	s.TurnState.Move = nil
	s.Phase = PhaseEndTurn
	s.Version++

	ev := Event{
		Type: EvClueTaken,
		Data: map[string]any{
			"clueId":    clue.ID,
			"trait":     clue.Trait,
			"result":    result,
			"boardCell": clue.BoardCell,
		},
	}

	return s, []Event{ev}, nil
}

func applyRevealSuspects(s GameState, c RevealSuspectsCommand) (GameState, []Event, error) {
	if s.Phase != PhaseRevealSuspects {
		return s, nil, ErrInvalidPhase
	}
	if s.TurnState.Pending != PendingRevealSuspects {
		return s, nil, ErrNoPendingAction
	}

	if len(c.SuspectIDs) != 2 {
		return s, nil, ErrInvalidRevealSelection
	}
	if c.SuspectIDs[0] == c.SuspectIDs[1] {
		return s, nil, ErrInvalidRevealSelection
	}

	revealed := make([]string, 0, 2)

	for _, suspectID := range c.SuspectIDs {
		suspect, ok := findSuspectByID(s.Suspects, suspectID)
		if !ok {
			return s, nil, ErrSuspectNotFound
		}
		if suspect.Revealed {
			return s, nil, ErrSuspectAlreadyRevealed
		}

		suspect.Revealed = true
		revealed = append(revealed, suspect.ID)
	}

	applyAutoExclusion(&s)

	s.TurnState.Pending = PendingNone
	s.TurnState.Move = nil
	s.Phase = PhaseEndTurn
	s.Version++

	ev := Event{
		Type: EvSuspectsRevealed,
		Data: map[string]any{
			"ids": revealed,
		},
	}

	return s, []Event{ev}, nil
}

func applyAccuse(s GameState, c AccuseCommand) (GameState, []Event, error) {
	suspect, ok := findSuspectByID(s.Suspects, c.SuspectID)
	if !ok {
		return s, nil, ErrSuspectNotFound
	}
	if !suspect.Revealed {
		return s, nil, ErrSuspectNotRevealed
	}
	if suspect.Excluded {
		return s, nil, ErrSuspectExcluded
	}

	correct := c.SuspectID == s.Secret.CulpritSuspectID

	s.Status = StatusFinished
	if correct {
		s.Result = ResultWin
	} else {
		s.Result = ResultLose
	}

	s.TurnState.ResetForNextTurn()
	s.Phase = PhaseEndTurn
	s.Version++

	evs := []Event{
		{
			Type: EvAccused,
			Data: map[string]any{
				"suspectId": c.SuspectID,
				"correct":   correct,
			},
		},
		{
			Type: EvGameFinished,
			Data: map[string]any{
				"result": s.Result,
			},
		},
	}

	return s, evs, nil
}

func applyEndTurn(s GameState, c EndTurnCommand) (GameState, []Event, error) {
	if s.Phase != PhaseEndTurn {
		return s, nil, ErrInvalidPhase
	}

	s.TurnState.ResetForNextTurn()
	s.ActiveSeat = nextSeat(s)
	s.Turn++
	s.Phase = PhaseChooseGoal
	s.Version++

	ev := Event{
		Type: EvTurnEnded,
		Data: map[string]any{
			"activeSeat": s.ActiveSeat,
			"turn":       s.Turn,
		},
	}

	return s, []Event{ev}, nil
}

func nextSeat(s GameState) int {
	n := len(s.Players)
	if n == 0 {
		return 0
	}
	return (s.ActiveSeat + 1) % n
}

func countMoveSteps(faces []string) int {
	steps := 0
	for _, face := range faces {
		if face == "footprint" || face == "move" || face == "step" {
			steps++
		}
	}
	return steps
}

func findClueByID(clues []ClueToken, id string) (*ClueToken, bool) {
	for i := range clues {
		if clues[i].ID == id {
			return &clues[i], true
		}
	}
	return nil, false
}

func findSuspectByID(suspects []SuspectCard, id string) (*SuspectCard, bool) {
	for i := range suspects {
		if suspects[i].ID == id {
			return &suspects[i], true
		}
	}
	return nil, false
}

func ptrTraitValue(v TraitValue) *TraitValue {
	return &v
}

// applyAutoExclusion - пока заглушка
// Позже тут будет настоящая дедукция по уликам.
func applyAutoExclusion(s *GameState) {
	_ = s
}
