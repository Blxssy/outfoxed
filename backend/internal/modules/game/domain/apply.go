package domain

func Apply(s GameState, cmd Command, rng RNG) (GameState, []Event, error) {
	if s.Status == StatusFinished {
		return s, nil, ErrGameFinished
	}
	if s.Status != StatusActive {
		return s, nil, ErrGameNotActive
	}

	activePlayer, ok := s.ActivePlayer()
	if !ok || activePlayer.ID != cmd.Actor() {
		return GameState{}, nil, ErrNotYourTurn
	}

	switch c := cmd.(type) {
	case ChooseGoalCommand:
		return applyChoseGoal(s, c)
	case RollAutoCommand:
		return applyRollAuto(s, c, rng)
	case EndTurnCommand:
		return applyEndTurn(s, c)
	case TakeClueCommand:
		return applyTakeClue(s, c)
	case RevealSuspectsCommand:
		return applyRevealSuspects(s, c)
	case AccuseCommand:
		return applyAccuse(s, c)

	default:
		return s, nil, ErrInvalidPhase
	}
}

func applyChoseGoal(s GameState, c ChooseGoalCommand) (GameState, []Event, error) {
	if s.Phase != PhaseChooseGoal {
		return s, nil, ErrInvalidPhase
	}

	if s.Goal.Set {
		return s, nil, ErrGoalAlreadySet
	}

	s.Goal.Type = c.Goal
	s.Goal.Set = true
	s.Phase = PhaseRolling
	s.Version++

	ev := Event{
		Type: EvGoalChosen,
		Data: map[string]any{"goal": c.Goal},
	}

	return s, []Event{ev}, nil
}

func applyRollAuto(s GameState, c RollAutoCommand, rng RNG) (GameState, []Event, error) {
	if s.Phase != PhaseRolling {
		return s, nil, ErrInvalidPhase
	}
	if !s.Goal.Set {
		return s, nil, ErrGoalNotSet
	}

	res := RollForGoal(s.Goal.Type, rng)

	events := make([]Event, 0, 2)
	events = append(events, Event{
		Type: EvRolled,
		Data: map[string]any{
			"success":  res.Success,
			"goal":     res.Goal,
			"attempts": res.Attempts,
			"faces":    res.Faces,
		},
	})

	if !res.Success {
		s.FoxTrack += 3
		if s.FoxEscapeAt > 0 && s.FoxTrack >= s.FoxEscapeAt {
			s.Status = StatusFinished
			s.Result = ResultLose
			events = append(events, Event{
				Type: EvGameFinished,
				Data: map[string]any{"result": s.Result},
			})
			s.Version++
			return s, events, nil
		}
		s.Pending = PendingNone
		events = append(events, Event{
			Type: EvFoxMoved,
			Data: map[string]any{"by": 3, "foxTrack": s.FoxTrack},
		})
		s.Phase = PhaseEndTurn
	} else {
		if s.Goal.Type == GoalClue {
			s.Pending = PendingClue
		} else {
			s.Pending = PendingSuspect
		}
		s.Phase = PhaseAction
	}

	s.Version++
	return s, events, nil
}

func applyEndTurn(s GameState, c EndTurnCommand) (GameState, []Event, error) {
	if s.Phase != PhaseEndTurn && s.Phase != PhaseAction {
		// разрешим завершать ход из Action тоже (пока MVP)
		return s, nil, ErrInvalidPhase
	}

	// очистим цель, подготовим следующего игрока
	s.Goal.Set = false
	s.Goal.Type = ""

	s.ActiveSeat = nextSeat(s)
	s.Turn++
	s.Phase = PhaseChooseGoal
	s.Version++

	ev := Event{Type: EvTurnEnded, Data: map[string]any{"activeSeat": s.ActiveSeat, "turn": s.Turn}}
	return s, []Event{ev}, nil
}

func applyTakeClue(s GameState, c TakeClueCommand) (GameState, []Event, error) {
	if s.Phase != PhaseAction {
		return s, nil, ErrInvalidPhase
	}

	if s.Pending == PendingNone {
		return s, nil, ErrNoPendingAction
	}
	if s.Pending != PendingClue {
		return s, nil, ErrPendingNotClue
	}
	if s.CluesTotal > 0 && s.CluesFound >= s.CluesTotal {
		return s, nil, ErrAllCluesCollected
	}

	s.CluesFound++
	s.Pending = PendingNone
	s.Phase = PhaseEndTurn
	s.Version++

	ev := Event{
		Type: EvClueTaken,
		Data: map[string]any{
			"cluesFound": s.CluesFound,
			"cluesTotal": s.CluesTotal,
		},
	}
	return s, []Event{ev}, nil
}

func applyRevealSuspects(s GameState, c RevealSuspectsCommand) (GameState, []Event, error) {
	if s.Phase != PhaseAction {
		return s, nil, ErrInvalidPhase
	}
	if s.Pending == PendingNone {
		return s, nil, ErrNoPendingAction
	}
	if s.Pending != PendingSuspect {
		return s, nil, ErrPendingNotSuspect
	}

	revealed := make([]int, 0, 2)

	for i := range s.Suspects {
		if !s.Suspects[i].Revealed {
			s.Suspects[i].Revealed = true
			revealed = append(revealed, s.Suspects[i].ID)
			if len(revealed) == 2 {
				break
			}
		}
	}

	if len(revealed) == 0 {
		return s, nil, ErrNoSuspectsToReveal
	}

	s.Pending = PendingNone
	s.Phase = PhaseEndTurn
	s.Version++

	ev := Event{
		Type: "suspects_revealed",
		Data: map[string]any{
			"ids": revealed,
		},
	}

	return s, []Event{ev}, nil
}

func applyAccuse(s GameState, c AccuseCommand) (GameState, []Event, error) {
	// В текущем MVP обвинение допускается в любой активной фазе хода.

	// Базовая защита: обвинять можно только раскрытого и не исключённого
	if c.SuspectID < 0 || c.SuspectID >= len(s.Suspects) {
		return s, nil, ErrInvalidPhase
	}
	sp := s.Suspects[c.SuspectID]
	if !sp.Revealed {
		return s, nil, ErrSuspectNotRevealed
	}
	if sp.Excluded {
		return s, nil, ErrSuspectExcluded
	}

	correct := (c.SuspectID == s.CulpritID)

	s.Version++

	evs := []Event{
		{
			Type: EvAccused,
			Data: map[string]any{
				"suspectId": c.SuspectID,
				"correct":   correct,
			},
		},
	}

	// Завершение игры
	s.Status = StatusFinished
	if correct {
		s.Result = ResultWin
	} else {
		s.Result = ResultLose
	}
	s.Phase = ""
	s.Pending = PendingNone

	evs = append(evs, Event{
		Type: EvGameFinished,
		Data: map[string]any{
			"result": s.Result,
		},
	})

	return s, evs, nil
}

func nextSeat(s GameState) int {
	n := len(s.Players)
	if n == 0 {
		return 0
	}
	return (s.ActiveSeat + 1) % n
}
