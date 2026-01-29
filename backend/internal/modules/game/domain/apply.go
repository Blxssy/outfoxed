package domain

func Apply(s GameState, cmd Command, rng RNG) (GameState, []Event, error) {
	if s.Status != StatusActive {
		return GameState{}, nil, ErrGameNotActive
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

	success := rng.Intn(2) == 1

	events := make([]Event, 0, 2)
	events = append(events, Event{
		Type: EvRolled,
		Data: map[string]any{
			"success": success,
			"goal":    s.Goal.Type,
		},
	})

	if !success {
		s.FoxTrack += 3
		events = append(events, Event{
			Type: EvFoxMoved,
			Data: map[string]any{"by": 3, "foxTrack": s.FoxTrack},
		})
		// при неуспехе ход по сути заканчивается
		s.Phase = PhaseEndTurn
	} else {
		// при успехе будет действие (движение/открытие) — пока переходим в PhaseAction
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

func nextSeat(s GameState) int {
	n := len(s.Players)
	if n == 0 {
		return 0
	}
	return (s.ActiveSeat + 1) % n
}
