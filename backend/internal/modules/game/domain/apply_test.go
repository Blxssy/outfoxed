package domain

import "testing"

func TestApplyMovePawn(t *testing.T) {
	rng := &FixedRNG{Values: []int{0}}

	state := NewActiveGameState("g1", []SetupPlayer{
		{UserID: "u1", Name: "A", Seat: 0},
		{UserID: "u2", Name: "B", Seat: 1},
	}, rng)

	idx := activePlayerIndex(state.Players, state.ActiveSeat)
	if idx < 0 {
		t.Fatal("active player not found")
	}

	from := state.Players[idx].PawnCell

	state.Phase = PhaseMovePawn
	state.TurnState.Pending = PendingMoveToClue
	state.TurnState.Move = &MoveState{
		StepsTotal:     2,
		StepsRemaining: 2,
	}

	state.TurnState.Move.ReachableCells = state.Board.ReachableWithin(
		from,
		2,
		occupiedCells(state.Players, state.ActiveSeat),
	)

	if len(state.TurnState.Move.ReachableCells) == 0 {
		t.Fatal("expected reachable cells, got none")
	}

	target := state.TurnState.Move.ReachableCells[0]

	newState, events, err := applyMovePawn(state, MovePawnCommand{
		Player:      state.Players[idx].UserID,
		TargetIndex: target,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) == 0 {
		t.Fatalf("expected move event")
	}

	if newState.Players[idx].PawnCell != target {
		t.Fatalf("expected pawn at %d, got %d", target, newState.Players[idx].PawnCell)
	}

	if newState.TurnState.Move == nil && newState.Phase == PhaseMovePawn {
		t.Fatalf("expected move state while still in move phase")
	}

	dist := ManhattanDistance(from, target, state.Board.Width)

	// Если не попали на улику и не закончили ход, должно остаться steps - dist
	if newState.Phase == PhaseMovePawn && newState.TurnState.Move != nil {
		want := 2 - dist
		if newState.TurnState.Move.StepsRemaining != want {
			t.Fatalf("expected %d step(s) remaining, got %d", want, newState.TurnState.Move.StepsRemaining)
		}
	}
}

func TestApplyMovePawn_ToOccupiedCell_ReturnsError(t *testing.T) {
	rng := &FixedRNG{Values: []int{0}}

	state := NewActiveGameState("g1", []SetupPlayer{
		{UserID: "u1", Name: "A", Seat: 0},
		{UserID: "u2", Name: "B", Seat: 1},
	}, rng)

	idx := activePlayerIndex(state.Players, state.ActiveSeat)
	if idx < 0 {
		t.Fatal("active player not found")
	}

	state.Phase = PhaseMovePawn
	state.TurnState.Pending = PendingMoveToClue
	state.TurnState.Move = &MoveState{
		StepsTotal:     2,
		StepsRemaining: 2,
		ReachableCells: []int{},
	}

	occupiedTarget := state.Players[1].PawnCell

	_, _, err := applyMovePawn(state, MovePawnCommand{
		Player:      state.Players[idx].UserID,
		TargetIndex: occupiedTarget,
	})
	if err == nil {
		t.Fatal("expected error for occupied target")
	}
}
