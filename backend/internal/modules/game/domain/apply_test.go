package domain

import "testing"

func TestChooseGoal_WrongPhase(t *testing.T) {
	s := baseState()
	s.Phase = PhaseRolling

	_, _, err := Apply(s, ChooseGoalCommand{Player: "p1", Goal: GoalClue}, &FixedRNG{Values: []int{1}})
	if err != ErrInvalidPhase {
		t.Fatalf("expected ErrInvalidPhase, got %v", err)
	}
}

func TestRollAuto_NotYourTurn(t *testing.T) {
	s := baseState()
	s.ActiveSeat = 1 // ходит p2

	_, _, err := Apply(s, RollAutoCommand{Player: "p1"}, &FixedRNG{Values: []int{1}})
	if err != ErrNotYourTurn {
		t.Fatalf("expected ErrNotYourTurn, got %v", err)
	}
}

func TestRollAuto_FailureMovesFox(t *testing.T) {
	s := baseState()
	// choose goal
	s2, _, err := Apply(s, ChooseGoalCommand{Player: "p1", Goal: GoalClue}, &FixedRNG{Values: []int{0}})
	if err != nil {
		t.Fatalf("choose goal err: %v", err)
	}
	// force failure: Intn(2)==0
	s3, evs, err := Apply(s2, RollAutoCommand{Player: "p1"}, &FixedRNG{Values: []int{0}})
	if err != nil {
		t.Fatalf("roll err: %v", err)
	}
	if s3.FoxTrack != 3 {
		t.Fatalf("expected foxTrack=3, got %d", s3.FoxTrack)
	}
	if s3.Phase != PhaseEndTurn {
		t.Fatalf("expected PhaseEndTurn, got %s", s3.Phase)
	}
	if len(evs) < 2 {
		t.Fatalf("expected 2 events, got %d", len(evs))
	}
}

func baseState() GameState {
	s := GameState{
		ID:         "g1",
		Status:     StatusActive,
		Phase:      PhaseChooseGoal,
		Turn:       1,
		Version:    1,
		ActiveSeat: 0,
		FoxTrack:   0,
		Players: []PlayerState{
			{ID: "p1", Seat: 0, Position: 0},
			{ID: "p2", Seat: 1, Position: 0},
		},
	}
	return s
}
