package domain

import (
	"reflect"
	"testing"
)

type testRNG struct {
	values []int
	pos    int
}

func (r *testRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	if len(r.values) == 0 {
		return 0
	}

	v := r.values[r.pos%len(r.values)]
	r.pos++
	if v < 0 {
		v = -v
	}
	return v % n
}

func buildTestBoard() BoardState {
	cells := make([]BoardCell, 0, BoardSize)

	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			idx := XYToIndex(x, y, BoardWidth)
			cells = append(cells, BoardCell{
				Index: idx,
				X:     x,
				Y:     y,
				Type:  BoardCellPath,
			})
		}
	}

	return BoardState{
		Width:  BoardWidth,
		Height: BoardHeight,
		Cells:  cells,
	}
}

func setClueCell(board *BoardState, index int, clueID string) {
	board.Cells[index].Type = BoardCellClue
	board.Cells[index].ClueTokenID = clueID
}

func buildBaseState(phase GamePhase) GameState {
	board := buildTestBoard()

	return GameState{
		ID:         "g1",
		Status:     StatusActive,
		Phase:      phase,
		Result:     ResultNone,
		Version:    1,
		Turn:       1,
		ActiveSeat: 0,
		Players: []PlayerState{
			{
				UserID:    PlayerID("u1"),
				Seat:      0,
				Name:      "Bot",
				PawnCell:  XYToIndex(7, 7, BoardWidth),
				Connected: true,
			},
		},
		Board:     board,
		Suspects:  nil,
		Clues:     nil,
		TurnState: NewTurnState(),
	}
}

func TestBuildAutoCommand_ChooseGoal_PicksClue_WhenUnrevealedCluesExist(t *testing.T) {
	st := buildBaseState(PhaseChooseGoal)
	st.Clues = []ClueToken{
		{ID: "clue_1", Revealed: false},
		{ID: "clue_2", Revealed: true},
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	c, ok := cmd.(ChooseGoalCommand)
	if !ok {
		t.Fatalf("expected ChooseGoalCommand, got %T", cmd)
	}

	if c.Goal != GoalClue {
		t.Fatalf("expected goal %q, got %q", GoalClue, c.Goal)
	}
}

func TestBuildAutoCommand_ChooseGoal_PicksSuspect_WhenAllCluesRevealed(t *testing.T) {
	st := buildBaseState(PhaseChooseGoal)
	st.Clues = []ClueToken{
		{ID: "clue_1", Revealed: true},
		{ID: "clue_2", Revealed: true},
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	c, ok := cmd.(ChooseGoalCommand)
	if !ok {
		t.Fatalf("expected ChooseGoalCommand, got %T", cmd)
	}

	if c.Goal != GoalSuspect {
		t.Fatalf("expected goal %q, got %q", GoalSuspect, c.Goal)
	}
}

func TestBuildAutoCommand_Rolling_PicksRerollKeepingMatchingDice(t *testing.T) {
	st := buildBaseState(PhaseRolling)
	st.TurnState.Goal = TurnGoal{
		Set:  true,
		Type: GoalClue,
	}
	st.TurnState.Roll = &RollState{
		RollsUsed: 1,
		MaxRolls:  3,
		Faces:     []string{FaceFootprint, FaceEye, FaceFootprint},
		Kept:      []bool{false, false, false},
		Success:   false,
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	c, ok := cmd.(RerollDiceCommand)
	if !ok {
		t.Fatalf("expected RerollDiceCommand, got %T", cmd)
	}

	want := []int{0, 2}
	if !reflect.DeepEqual(c.KeepIndices, want) {
		t.Fatalf("expected keep indices %v, got %v", want, c.KeepIndices)
	}
}

func TestBuildAutoCommand_Rolling_FinishesWhenSuccess(t *testing.T) {
	st := buildBaseState(PhaseRolling)
	st.TurnState.Goal = TurnGoal{
		Set:  true,
		Type: GoalClue,
	}
	st.TurnState.Roll = &RollState{
		RollsUsed: 1,
		MaxRolls:  3,
		Faces:     []string{FaceFootprint, FaceFootprint, FaceFootprint},
		Kept:      []bool{true, true, true},
		Success:   true,
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	if _, ok := cmd.(FinishRollCommand); !ok {
		t.Fatalf("expected FinishRollCommand, got %T", cmd)
	}
}

func TestBuildAutoCommand_Rolling_FinishesWhenNoRollsLeft(t *testing.T) {
	st := buildBaseState(PhaseRolling)
	st.TurnState.Goal = TurnGoal{
		Set:  true,
		Type: GoalClue,
	}
	st.TurnState.Roll = &RollState{
		RollsUsed: 3,
		MaxRolls:  3,
		Faces:     []string{FaceFootprint, FaceEye, FaceFootprint},
		Kept:      []bool{true, false, true},
		Success:   false,
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	if _, ok := cmd.(FinishRollCommand); !ok {
		t.Fatalf("expected FinishRollCommand, got %T", cmd)
	}
}

func TestBuildAutoCommand_MovePawn_GoesDirectlyToClue_WhenReachable(t *testing.T) {
	st := buildBaseState(PhaseMovePawn)

	clueIndex := XYToIndex(7, 9, BoardWidth)
	setClueCell(&st.Board, clueIndex, "clue_1")

	st.Clues = []ClueToken{
		{
			ID:        "clue_1",
			Revealed:  false,
			BoardCell: clueIndex,
		},
	}

	st.TurnState.Move = &MoveState{
		StepsTotal:     2,
		StepsRemaining: 2,
		ReachableCells: []int{
			XYToIndex(7, 8, BoardWidth),
			clueIndex,
			XYToIndex(8, 8, BoardWidth),
		},
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	c, ok := cmd.(MovePawnCommand)
	if !ok {
		t.Fatalf("expected MovePawnCommand, got %T", cmd)
	}

	if c.TargetIndex != clueIndex {
		t.Fatalf("expected target %d, got %d", clueIndex, c.TargetIndex)
	}
}

func TestBuildAutoCommand_MovePawn_GoesCloserToNearestClue(t *testing.T) {
	st := buildBaseState(PhaseMovePawn)

	// Игрок стоит в центре старта (7,7)
	// Улика находится ниже.
	clueIndex := XYToIndex(7, 12, BoardWidth)
	setClueCell(&st.Board, clueIndex, "clue_1")

	st.Clues = []ClueToken{
		{
			ID:        "clue_1",
			Revealed:  false,
			BoardCell: clueIndex,
		},
	}

	candidateA := XYToIndex(7, 8, BoardWidth) // ближе к улике
	candidateB := XYToIndex(6, 7, BoardWidth) // дальше от улики
	candidateC := XYToIndex(8, 7, BoardWidth) // дальше от улики

	st.TurnState.Move = &MoveState{
		StepsTotal:     1,
		StepsRemaining: 1,
		ReachableCells: []int{candidateB, candidateA, candidateC},
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	c, ok := cmd.(MovePawnCommand)
	if !ok {
		t.Fatalf("expected MovePawnCommand, got %T", cmd)
	}

	if c.TargetIndex != candidateA {
		t.Fatalf("expected target %d, got %d", candidateA, c.TargetIndex)
	}
}

func TestBuildAutoCommand_ResolveClue_TakesClue_WhenAvailable(t *testing.T) {
	st := buildBaseState(PhaseResolveClue)

	playerCell := st.Players[0].PawnCell
	setClueCell(&st.Board, playerCell, "clue_1")

	st.Clues = []ClueToken{
		{
			ID:        "clue_1",
			Revealed:  false,
			BoardCell: playerCell,
		},
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	if _, ok := cmd.(TakeClueCommand); !ok {
		t.Fatalf("expected TakeClueCommand, got %T", cmd)
	}
}

func TestBuildAutoCommand_ResolveClue_EndsTurn_WhenClueAlreadyRevealed(t *testing.T) {
	st := buildBaseState(PhaseResolveClue)

	playerCell := st.Players[0].PawnCell
	setClueCell(&st.Board, playerCell, "clue_1")

	st.Clues = []ClueToken{
		{
			ID:        "clue_1",
			Revealed:  true,
			BoardCell: playerCell,
		},
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	if _, ok := cmd.(EndTurnCommand); !ok {
		t.Fatalf("expected EndTurnCommand, got %T", cmd)
	}
}

func TestBuildAutoCommand_RevealSuspects_PicksFirstTwoUnrevealed(t *testing.T) {
	st := buildBaseState(PhaseRevealSuspects)
	st.Suspects = []SuspectCard{
		{ID: "suspect_1", Revealed: true},
		{ID: "suspect_2", Revealed: false},
		{ID: "suspect_3", Revealed: false},
		{ID: "suspect_4", Revealed: false},
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	c, ok := cmd.(RevealSuspectsCommand)
	if !ok {
		t.Fatalf("expected RevealSuspectsCommand, got %T", cmd)
	}

	want := []string{"suspect_2", "suspect_3"}
	if !reflect.DeepEqual(c.SuspectIDs, want) {
		t.Fatalf("expected suspect ids %v, got %v", want, c.SuspectIDs)
	}
}

func TestBuildAutoCommand_RevealSuspects_EndsTurn_WhenLessThanTwoRemain(t *testing.T) {
	st := buildBaseState(PhaseRevealSuspects)
	st.Suspects = []SuspectCard{
		{ID: "suspect_1", Revealed: true},
		{ID: "suspect_2", Revealed: true},
		{ID: "suspect_3", Revealed: false},
	}

	cmd, ok := BuildAutoCommand(st, &testRNG{})
	if !ok {
		t.Fatal("expected command")
	}

	if _, ok := cmd.(EndTurnCommand); !ok {
		t.Fatalf("expected EndTurnCommand, got %T", cmd)
	}
}
