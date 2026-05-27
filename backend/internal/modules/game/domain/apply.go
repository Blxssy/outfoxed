package domain

import "time"

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
		return applyChooseGoal(s, c, rng)

	case RerollDiceCommand:
		return applyRerollDice(s, c, rng)

	case FinishRollCommand:
		return applyFinishRoll(s, c)

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

func applyChooseGoal(s GameState, c ChooseGoalCommand, rng RNG) (GameState, []Event, error) {
	if s.Phase != PhaseChooseGoal {
		return s, nil, ErrInvalidPhase
	}
	if s.TurnState.Goal.Set {
		return s, nil, ErrGoalAlreadySet
	}

	faces := rollInitialFaces(rng)
	success := isRollSuccessful(c.Goal, faces)

	s.TurnState.Goal = TurnGoal{
		Set:  true,
		Type: c.Goal,
	}
	s.TurnState.Pending = PendingNone
	s.TurnState.Move = nil
	s.TurnState.Roll = &RollState{
		RollsUsed: 1,
		MaxRolls:  MaxRolls,
		Faces:     faces,
		Kept:      make([]bool, len(faces)),
		Success:   success,
	}

	s.Phase = PhaseRolling
	s.Version++

	return s, []Event{
		{
			Type: EvGoalChosen,
			Data: map[string]any{
				"goal": c.Goal,
			},
		},
		{
			Type: EvRolled,
			Data: map[string]any{
				"goal":      c.Goal,
				"faces":     faces,
				"rollsUsed": 1,
				"maxRolls":  MaxRolls,
				"success":   success,
			},
		},
	}, nil
}

func applyRerollDice(s GameState, c RerollDiceCommand, rng RNG) (GameState, []Event, error) {
	if s.Phase != PhaseRolling {
		return s, nil, ErrInvalidPhase
	}
	if !s.TurnState.Goal.Set {
		return s, nil, ErrGoalNotSet
	}
	if s.TurnState.Roll == nil {
		return s, nil, ErrInvalidPhase
	}

	roll := s.TurnState.Roll
	if roll.RollsUsed >= roll.MaxRolls {
		return s, nil, ErrNoRollsLeft
	}

	keep := make([]bool, DiceCount)
	for _, idx := range c.KeepIndices {
		if idx < 0 || idx >= DiceCount {
			return s, nil, ErrInvalidMove
		}
		keep[idx] = true
	}

	roll.Faces = rerollFaces(roll.Faces, keep, rng)
	roll.Kept = keep
	roll.RollsUsed++
	roll.Success = isRollSuccessful(s.TurnState.Goal.Type, roll.Faces)

	s.Version++

	return s, []Event{
		{
			Type: EvRolled,
			Data: map[string]any{
				"goal":      s.TurnState.Goal.Type,
				"faces":     roll.Faces,
				"rollsUsed": roll.RollsUsed,
				"maxRolls":  roll.MaxRolls,
				"kept":      keep,
				"success":   roll.Success,
			},
		},
	}, nil
}

func applyFinishRoll(s GameState, c FinishRollCommand) (GameState, []Event, error) {
	if s.Phase != PhaseRolling {
		return s, nil, ErrInvalidPhase
	}
	if !s.TurnState.Goal.Set {
		return s, nil, ErrGoalNotSet
	}
	if s.TurnState.Roll == nil {
		return s, nil, ErrInvalidPhase
	}

	res := s.TurnState.Roll

	events := []Event{
		{
			Type: "roll_finished",
			Data: map[string]any{
				"success": res.Success,
				"goal":    s.TurnState.Goal.Type,
				"faces":   res.Faces,
			},
		},
	}

	if !res.Success {
		s.TurnState.Pending = PendingNone
		s.TurnState.Move = nil

		s.Fox.Track += FoxStepPerFailure
		events = append(events, Event{
			Type: EvFoxMoved,
			Data: map[string]any{
				"by":    FoxStepPerFailure,
				"track": s.Fox.Track,
			},
		})

		if s.Fox.EscapeAt > 0 && s.Fox.Track >= s.Fox.EscapeAt {
			s.Status = StatusFinished
			s.Result = ResultLose
			s.TurnDeadlineAt = nil
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

		idx := activePlayerIndex(s.Players, s.ActiveSeat)
		if idx < 0 {
			return s, nil, ErrInvalidMove
		}

		from := s.Players[idx].PawnCell
		occupied := occupiedCells(s.Players, s.ActiveSeat)

		move := &MoveState{
			StepsTotal:     steps,
			StepsRemaining: steps,
		}
		move.ReachableCells = s.Board.ReachableWithin(from, steps, occupied)

		s.TurnState.Pending = PendingMoveToClue
		s.TurnState.Move = move
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

func applyRollAuto(s GameState, c RollAutoCommand, rng RNG) (GameState, []Event, error) {
	if s.Phase != PhaseRolling {
		return s, nil, ErrInvalidPhase
	}
	if !s.TurnState.Goal.Set {
		return s, nil, ErrGoalNotSet
	}

	res := RollForGoal(s.TurnState.Goal.Type, rng)

	s.TurnState.Roll = &RollState{
		RollsUsed: res.Attempts,
		MaxRolls:  MaxRolls,
		Faces:     res.Faces,
		Kept:      []bool{true, true, true},
		Success:   res.Success,
	}

	return applyFinishRoll(s, FinishRollCommand{Player: c.Player})
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

	playerIdx := activePlayerIndex(s.Players, s.ActiveSeat)
	if playerIdx < 0 {
		return s, nil, ErrInvalidMove
	}

	from := s.Players[playerIdx].PawnCell
	to := c.TargetIndex

	if to < 0 || to >= len(s.Board.Cells) {
		return s, nil, ErrInvalidMove
	}
	if from == to {
		return s, nil, ErrInvalidMove
	}

	move := s.TurnState.Move
	if !containsInt(move.ReachableCells, to) {
		return s, nil, ErrInvalidMove
	}

	dist := ManhattanDistance(from, to, s.Board.Width)
	if dist <= 0 || dist > move.StepsRemaining {
		return s, nil, ErrInvalidMove
	}

	s.Players[playerIdx].PawnCell = to
	move.StepsRemaining -= dist

	events := []Event{
		{
			Type: EvPawnMoved,
			Data: map[string]any{
				"seat":      s.ActiveSeat,
				"fromCell":  from,
				"toCell":    to,
				"cost":      dist,
				"stepsLeft": move.StepsRemaining,
			},
		},
	}

	cell, ok := s.Board.CellAt(to)
	if !ok {
		return s, nil, ErrInvalidMove
	}

	if cell.Type == BoardCellClue && cell.ClueTokenID != "" {
		s.TurnState.Pending = PendingResolveClue
		s.TurnState.Move = nil
		s.Phase = PhaseResolveClue
		s.Version++
		return s, events, nil
	}

	if move.StepsRemaining == 0 {
		s.TurnState.Pending = PendingNone
		s.TurnState.Move = nil
		s.Phase = PhaseEndTurn
		s.Version++
		return s, events, nil
	}

	occupied := occupiedCells(s.Players, s.ActiveSeat)
	move.ReachableCells = s.Board.ReachableWithin(to, move.StepsRemaining, occupied)

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
		s.TurnState.Pending = PendingNone
		s.TurnState.Move = nil
		s.Phase = PhaseEndTurn
		s.Version++

		return s, []Event{
			{
				Type: "clue_already_taken",
				Data: map[string]any{
					"clueId": clue.ID,
				},
			},
		}, nil
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

	s.TurnDeadlineAt = nil

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

	activeIdx := activePlayerIndex(s.Players, s.ActiveSeat)
	if activeIdx >= 0 {
		s.TurnDeadlineAt = computeTurnDeadlineForPlayer(s.Players[activeIdx])
	} else {
		s.TurnDeadlineAt = nil
	}

	s.Version++

	var deadline time.Time
	if s.TurnDeadlineAt != nil {
		deadline = *s.TurnDeadlineAt
	}

	ev := Event{
		Type: EvTurnEnded,
		Data: map[string]any{
			"activeSeat":     s.ActiveSeat,
			"turn":           s.Turn,
			"turnDeadlineAt": deadline,
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
