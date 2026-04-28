package domain

func BuildAutoCommand(st GameState, rng RNG) (Command, bool) {
	activePlayer, ok := st.ActivePlayer()
	if !ok {
		return nil, false
	}
	actor := activePlayer.UserID

	switch st.Phase {
	case PhaseChooseGoal:
		goals := []GoalType{GoalClue, GoalSuspect}
		return ChooseGoalCommand{
			Player: actor,
			Goal:   goals[rng.Intn(len(goals))],
		}, true

	case PhaseRolling:
		return RollAutoCommand{Player: actor}, true

	case PhaseMovePawn:
		if st.TurnState.Move == nil || st.TurnState.Move.StepsRemaining <= 0 {
			return nil, false
		}
		return MovePawnCommand{
			Player: actor,
			Steps:  st.TurnState.Move.StepsRemaining,
		}, true

	case PhaseResolveClue:
		return TakeClueCommand{Player: actor}, true

	case PhaseRevealSuspects:
		ids := pickTwoAutoSuspects(st)
		if len(ids) < 2 {
			return EndTurnCommand{Player: actor}, true
		}
		return RevealSuspectsCommand{
			Player:     actor,
			SuspectIDs: ids,
		}, true

	case PhaseEndTurn:
		return EndTurnCommand{Player: actor}, true

	default:
		return nil, false
	}
}

func pickTwoAutoSuspects(st GameState) []string {
	ids := make([]string, 0, 2)
	for _, s := range st.Suspects {
		if s.Revealed {
			continue
		}
		ids = append(ids, s.ID)
		if len(ids) == 2 {
			return ids
		}
	}
	return ids
}
