package domain

func BuildAutoCommand(st GameState, rng RNG) (Command, bool) {
	activePlayer, ok := st.ActivePlayer()
	if !ok {
		return nil, false
	}
	actor := activePlayer.UserID

	switch st.Phase {
	case PhaseChooseGoal:
		return ChooseGoalCommand{
			Player: actor,
			Goal:   pickAutoGoal(st),
		}, true

	case PhaseRolling:
		if st.TurnState.Roll == nil || !st.TurnState.Goal.Set {
			return FinishRollCommand{Player: actor}, true
		}

		roll := st.TurnState.Roll
		if roll.Success || roll.RollsUsed >= roll.MaxRolls {
			return FinishRollCommand{Player: actor}, true
		}

		keep := pickKeepIndicesForGoal(roll, st.TurnState.Goal.Type)
		return RerollDiceCommand{
			Player:      actor,
			KeepIndices: keep,
		}, true

	case PhaseMovePawn:
		if st.TurnState.Move == nil || len(st.TurnState.Move.ReachableCells) == 0 {
			return EndTurnCommand{Player: actor}, true
		}

		target, ok := pickBestAutoMoveTarget(st)
		if !ok {
			return EndTurnCommand{Player: actor}, true
		}

		return MovePawnCommand{
			Player:      actor,
			TargetIndex: target,
		}, true

	case PhaseResolveClue:
		if !canAutoTakeClue(st) {
			return EndTurnCommand{Player: actor}, true
		}
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

func pickAutoGoal(st GameState) GoalType {
	// Пока логика простая:
	// если ещё есть неоткрытые улики — бот в приоритете ищет улику,
	// иначе идёт в подозреваемых.
	if hasUnrevealedClues(st) {
		return GoalClue
	}
	return GoalSuspect
}

func hasUnrevealedClues(st GameState) bool {
	for _, clue := range st.Clues {
		if !clue.Revealed {
			return true
		}
	}
	return false
}

func pickKeepIndicesForGoal(roll *RollState, goal GoalType) []int {
	if roll == nil {
		return nil
	}

	want := faceForGoal(goal)
	keep := make([]int, 0, len(roll.Faces))

	for i, face := range roll.Faces {
		if face == want {
			keep = append(keep, i)
		}
	}

	return keep
}

func pickBestAutoMoveTarget(st GameState) (int, bool) {
	move := st.TurnState.Move
	if move == nil || len(move.ReachableCells) == 0 {
		return 0, false
	}

	// 1. Если можем сразу встать на клетку с неоткрытой уликой — идём туда.
	bestDirect := -1
	for _, idx := range move.ReachableCells {
		if isUnrevealedClueCell(st, idx) {
			if bestDirect == -1 || idx < bestDirect {
				bestDirect = idx
			}
		}
	}
	if bestDirect != -1 {
		return bestDirect, true
	}

	// 2. Иначе идём в клетку, которая минимизирует расстояние до ближайшей неоткрытой улики.
	bestIdx := move.ReachableCells[0]
	bestScore := distanceToNearestUnrevealedClue(st, bestIdx)

	for _, idx := range move.ReachableCells[1:] {
		score := distanceToNearestUnrevealedClue(st, idx)

		// Чем меньше score, тем лучше.
		// При равенстве берём меньший индекс, чтобы поведение было детерминированным.
		if score < bestScore || (score == bestScore && idx < bestIdx) {
			bestIdx = idx
			bestScore = score
		}
	}

	return bestIdx, true
}

func isUnrevealedClueCell(st GameState, index int) bool {
	cell, ok := st.Board.CellAt(index)
	if !ok {
		return false
	}
	if cell.Type != BoardCellClue || cell.ClueTokenID == "" {
		return false
	}

	clue, ok := findClueByID(st.Clues, cell.ClueTokenID)
	if !ok {
		return false
	}

	return !clue.Revealed
}

func distanceToNearestUnrevealedClue(st GameState, from int) int {
	best := -1

	for _, clue := range st.Clues {
		if clue.Revealed {
			continue
		}

		dist := ManhattanDistance(from, clue.BoardCell, st.Board.Width)
		if best == -1 || dist < best {
			best = dist
		}
	}

	// Если улик уже нет, просто возвращаем 0 — бот тогда выберет первую/минимальную клетку.
	if best == -1 {
		return 0
	}

	return best
}

func canAutoTakeClue(st GameState) bool {
	player, ok := st.ActivePlayer()
	if !ok {
		return false
	}

	cell, ok := st.Board.CellAt(player.PawnCell)
	if !ok {
		return false
	}
	if cell.Type != BoardCellClue || cell.ClueTokenID == "" {
		return false
	}

	clue, ok := findClueByID(st.Clues, cell.ClueTokenID)
	if !ok {
		return false
	}

	return !clue.Revealed
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
