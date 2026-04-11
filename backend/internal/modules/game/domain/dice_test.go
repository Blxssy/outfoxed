package domain

import (
	"testing"
)

func TestRollForGoal_SuccessWithin3(t *testing.T) {
	// Хотим clue => footprint.
	// Intn(2): 0=footprint, 1=eye
	// Сценарий: три нуля подряд => успех за 1 попытку
	rng := &FixedRNG{Values: []int{0, 0, 0}}

	res := RollForGoal(GoalClue, rng)
	if !res.Success {
		t.Fatalf("expected success")
	}
	if res.Attempts != 1 {
		t.Fatalf("expected attempts=1, got %d", res.Attempts)
	}
}

func TestRollForGoal_FailAfter3(t *testing.T) {
	// Хотим suspects => eye (1)
	// Дадим значения так, чтобы за 3 попытки не собрать 3 глаза:
	// всего будет много следов
	rng := &FixedRNG{Values: []int{0, 0, 0, 0, 0, 0, 0, 0, 0}}

	res := RollForGoal(GoalSuspect, rng)
	if res.Success {
		t.Fatalf("expected failure")
	}
	if res.Attempts != 3 {
		t.Fatalf("expected attempts=3, got %d", res.Attempts)
	}
}
