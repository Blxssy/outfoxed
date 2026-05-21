package domain

import (
	"testing"
)

func TestApplyChooseGoal_StartsFirstRoll(t *testing.T) {
	rng := &FixedRNG{Values: []int{0, 1, 0}}

	state := NewActiveGameState("g1", []SetupPlayer{
		{UserID: "u1", Name: "A", Seat: 0},
		{UserID: "u2", Name: "B", Seat: 1},
	}, rng)

	newState, events, err := Apply(state, ChooseGoalCommand{
		Player: "u1",
		Goal:   GoalClue,
	}, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if newState.Phase != PhaseRolling {
		t.Fatalf("expected phase %s, got %s", PhaseRolling, newState.Phase)
	}

	if newState.TurnState.Roll == nil {
		t.Fatal("expected roll state")
	}

	if newState.TurnState.Roll.RollsUsed != 1 {
		t.Fatalf("expected first roll used = 1, got %d", newState.TurnState.Roll.RollsUsed)
	}

	if len(events) < 2 {
		t.Fatalf("expected goal+rolled events")
	}
}
