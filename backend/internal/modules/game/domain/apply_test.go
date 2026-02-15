package domain

import (
	"errors"
	"testing"
)

func TestChooseGoal_WrongPhase(t *testing.T) {
	s := baseState()
	s.Phase = PhaseRolling

	_, _, err := Apply(s, ChooseGoalCommand{Player: "p1", Goal: GoalClue}, &FixedRNG{Values: []int{1}})
	if !errors.Is(err, ErrInvalidPhase) {
		t.Fatalf("expected ErrInvalidPhase, got %v", err)
	}
}

func TestRollAuto_NotYourTurn(t *testing.T) {
	s := baseState()
	s.ActiveSeat = 1 // ходит p2

	_, _, err := Apply(s, RollAutoCommand{Player: "p1"}, &FixedRNG{Values: []int{1}})
	if !errors.Is(err, ErrNotYourTurn) {
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
	s3, evs, err := Apply(s2, RollAutoCommand{Player: "p1"}, rngAllEyes())
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
func TestTakeClue_HappyPath(t *testing.T) {
	s := baseState()
	s.CluesTotal = 12

	// choose goal = clue
	s2, _, err := Apply(s, ChooseGoalCommand{Player: "p1", Goal: GoalClue}, &FixedRNG{Values: []int{0}})
	if err != nil {
		t.Fatalf("choose goal err: %v", err)
	}

	// roll success for clue: footprints => FixedRNG 0 gives footprint
	s3, _, err := Apply(s2, RollAutoCommand{Player: "p1"}, &FixedRNG{Values: []int{0, 0, 0}})
	if err != nil {
		t.Fatalf("roll err: %v", err)
	}
	if s3.Phase != PhaseAction {
		t.Fatalf("expected PhaseAction, got %s", s3.Phase)
	}
	if s3.Pending != PendingClue {
		t.Fatalf("expected pending clue, got %s", s3.Pending)
	}

	// take clue
	s4, evs, err := Apply(s3, TakeClueCommand{Player: "p1"}, &FixedRNG{Values: []int{0}})
	if err != nil {
		t.Fatalf("take clue err: %v", err)
	}
	if s4.CluesFound != 1 {
		t.Fatalf("expected cluesFound=1, got %d", s4.CluesFound)
	}
	if s4.Pending != PendingNone {
		t.Fatalf("expected pending none, got %s", s4.Pending)
	}
	if s4.Phase != PhaseEndTurn {
		t.Fatalf("expected PhaseEndTurn, got %s", s4.Phase)
	}
	if len(evs) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evs))
	}
}

func TestTakeClue_WrongPhase(t *testing.T) {
	s := baseState()
	s.Phase = PhaseRolling

	_, _, err := Apply(s, TakeClueCommand{Player: "p1"}, &FixedRNG{Values: []int{0}})
	if err != ErrInvalidPhase {
		t.Fatalf("expected ErrInvalidPhase, got %v", err)
	}
}

func TestTakeClue_PendingNotClue(t *testing.T) {
	s := baseState()
	s.Phase = PhaseAction
	s.Pending = PendingSuspect

	_, _, err := Apply(s, TakeClueCommand{Player: "p1"}, &FixedRNG{Values: []int{0}})
	if err != ErrPendingNotClue {
		t.Fatalf("expected ErrPendingNotClue, got %v", err)
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
		Pending:  PendingNone,
		Suspects: NewSuspects(16),
	}
	return s
}

func TestRevealSuspects_HappyPath(t *testing.T) {
	s := baseState()

	// choose goal = suspect
	s2, _, err := Apply(s, ChooseGoalCommand{Player: "p1", Goal: GoalSuspect}, &FixedRNG{Values: []int{1}})
	if err != nil {
		t.Fatalf("choose goal err: %v", err)
	}

	// success for suspect: eyes => 1
	s3, _, err := Apply(s2, RollAutoCommand{Player: "p1"}, &FixedRNG{Values: []int{1, 1, 1}})
	if err != nil {
		t.Fatalf("roll err: %v", err)
	}
	if s3.Pending != PendingSuspect {
		t.Fatalf("expected pending suspect, got %s", s3.Pending)
	}

	s4, evs, err := Apply(s3, RevealSuspectsCommand{Player: "p1"}, &FixedRNG{Values: []int{0}})
	if err != nil {
		t.Fatalf("reveal err: %v", err)
	}

	count := 0
	for _, sp := range s4.Suspects {
		if sp.Revealed {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("expected 2 revealed suspects, got %d", count)
	}

	if s4.Phase != PhaseEndTurn {
		t.Fatalf("expected PhaseEndTurn, got %s", s4.Phase)
	}

	if len(evs) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evs))
	}
}

func TestRevealSuspects_WrongPending(t *testing.T) {
	s := baseState()
	s.Phase = PhaseAction
	s.Pending = PendingClue

	_, _, err := Apply(s, RevealSuspectsCommand{Player: "p1"}, &FixedRNG{Values: []int{0}})
	if err != ErrPendingNotSuspect {
		t.Fatalf("expected ErrPendingNotSuspect, got %v", err)
	}
}

func rngAllEyes() *FixedRNG       { return &FixedRNG{Values: []int{1}} }
func rngAllFootprints() *FixedRNG { return &FixedRNG{Values: []int{0}} }
