package domain

import "time"

const (
	MinPlayers = 2
	MaxPlayers = 4
)

type SetupPlayer struct {
	UserID PlayerID
	Name   string
	Seat   int
}

func NewWaitingGameState(gameID string, players []SetupPlayer) GameState {
	statePlayers := make([]PlayerState, 0, len(players))
	for _, player := range players {
		statePlayers = append(statePlayers, PlayerState{
			UserID:    player.UserID,
			Seat:      player.Seat,
			Name:      player.Name,
			PawnCell:  0,
			Connected: false,
		})
	}

	return GameState{
		ID:         gameID,
		Status:     StatusWaiting,
		Phase:      PhaseChooseGoal,
		Result:     ResultNone,
		Version:    1,
		Turn:       0,
		ActiveSeat: 0,
		Players:    statePlayers,
		Board:      newBoardState(),
		Fox: FoxState{
			Track:    0,
			EscapeAt: 15,
		},
		Suspects:  newDefaultSuspects(),
		Clues:     newDefaultClues(),
		TurnState: NewTurnState(),
	}
}

func NewActiveGameState(gameID string, players []SetupPlayer) GameState {
	state := NewWaitingGameState(gameID, players)
	state.Status = StatusActive
	state.Turn = 1
	state.ActiveSeat = firstSeat(players)

	culprit := state.Suspects[0]
	state.Secret = SecretState{
		CulpritSuspectID: culprit.ID,
		ClueTruth:        buildClueTruth(culprit),
	}

	deadline := time.Now().UTC().Add(time.Minute)
	state.TurnDeadlineAt = &deadline

	return state
}

func firstSeat(players []SetupPlayer) int {
	if len(players) == 0 {
		return 0
	}

	minSeat := players[0].Seat
	for _, player := range players[1:] {
		if player.Seat < minSeat {
			minSeat = player.Seat
		}
	}
	return minSeat
}

func newBoardState() BoardState {
	return BoardState{
		Cells: []BoardCell{
			{Index: 0, Type: BoardCellStart},
			{Index: 1, Type: BoardCellPath},
			{Index: 2, Type: BoardCellPath},
			{Index: 3, Type: BoardCellClue, ClueTokenID: "clue_1"},
			{Index: 4, Type: BoardCellPath},
			{Index: 5, Type: BoardCellPath},
			{Index: 6, Type: BoardCellClue, ClueTokenID: "clue_2"},
			{Index: 7, Type: BoardCellPath},
			{Index: 8, Type: BoardCellPath},
			{Index: 9, Type: BoardCellClue, ClueTokenID: "clue_3"},
		},
	}
}

func newDefaultSuspects() []SuspectCard {
	return []SuspectCard{
		{
			ID:       "suspect_1",
			Code:     SuspectCode("scarlet"),
			Revealed: false,
			Excluded: false,
			Traits: SuspectTraits{
				Glasses:  TraitYes,
				Hat:      TraitNo,
				Scarf:    TraitYes,
				Umbrella: TraitNo,
				Color:    TraitRed,
			},
		},
		{
			ID:       "suspect_2",
			Code:     SuspectCode("azure"),
			Revealed: false,
			Excluded: false,
			Traits: SuspectTraits{
				Glasses:  TraitNo,
				Hat:      TraitYes,
				Scarf:    TraitYes,
				Umbrella: TraitNo,
				Color:    TraitBlue,
			},
		},
		{
			ID:       "suspect_3",
			Code:     SuspectCode("moss"),
			Revealed: false,
			Excluded: false,
			Traits: SuspectTraits{
				Glasses:  TraitYes,
				Hat:      TraitYes,
				Scarf:    TraitNo,
				Umbrella: TraitYes,
				Color:    TraitGreen,
			},
		},
		{
			ID:       "suspect_4",
			Code:     SuspectCode("amber"),
			Revealed: false,
			Excluded: false,
			Traits: SuspectTraits{
				Glasses:  TraitNo,
				Hat:      TraitNo,
				Scarf:    TraitNo,
				Umbrella: TraitYes,
				Color:    TraitYellow,
			},
		},
	}
}

func newDefaultClues() []ClueToken {
	return []ClueToken{
		{ID: "clue_1", Trait: ClueTraitGlasses, Revealed: false, BoardCell: 3},
		{ID: "clue_2", Trait: ClueTraitHat, Revealed: false, BoardCell: 6},
		{ID: "clue_3", Trait: ClueTraitColor, Revealed: false, BoardCell: 9},
	}
}

func buildClueTruth(culprit SuspectCard) map[string]TraitValue {
	return map[string]TraitValue{
		"clue_1": culprit.Traits.Glasses,
		"clue_2": culprit.Traits.Hat,
		"clue_3": culprit.Traits.Color,
	}
}
