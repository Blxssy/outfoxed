package domain

import (
	"fmt"
	"strconv"
	"time"
)

const (
	MinPlayers = 2
	MaxPlayers = 4
)

type SetupPlayer struct {
	UserID PlayerID
	Name   string
	Seat   int
}

func NewWaitingGameState(gameID string, players []SetupPlayer, rng RNG) GameState {
	statePlayers := buildInitialPlayers(players)

	clues := newDefaultClues(rng)
	clueIDs := clueIDsFromTokens(clues)

	board, err := NewBoard16x16(clueIDs, rng)
	if err != nil {
		panic(fmt.Errorf("build board: %w", err))
	}

	clues = bindCluesToBoard(clues, board)

	return GameState{
		ID:         gameID,
		Status:     StatusWaiting,
		Phase:      PhaseChooseGoal,
		Result:     ResultNone,
		Version:    1,
		Turn:       0,
		ActiveSeat: firstSeat(players),

		Players: statePlayers,

		Board: board,
		Fox: FoxState{
			Track:    0,
			EscapeAt: 15,
		},

		Suspects:  newDefaultSuspects(rng),
		Clues:     clues,
		TurnState: NewTurnState(),
	}
}

func NewActiveGameState(gameID string, players []SetupPlayer, rng RNG) GameState {
	state := NewWaitingGameState(gameID, players, rng)
	state.Status = StatusActive
	state.Turn = 1
	state.ActiveSeat = firstSeat(players)

	culprit := pickCulprit(state.Suspects, rng)
	state.Secret = SecretState{
		CulpritSuspectID: culprit.ID,
		ClueTruth:        buildClueTruth(culprit, state.Clues),
	}

	deadline := time.Now().UTC().Add(time.Minute)
	state.TurnDeadlineAt = &deadline

	return state
}

func buildInitialPlayers(players []SetupPlayer) []PlayerState {
	spawns := DefaultSpawnCells()
	statePlayers := make([]PlayerState, 0, len(players))

	for _, player := range players {
		spawn := spawnCellForSeat(player.Seat, spawns)

		statePlayers = append(statePlayers, PlayerState{
			UserID:    player.UserID,
			Seat:      player.Seat,
			Name:      player.Name,
			PawnCell:  spawn,
			Connected: false,
		})
	}

	return statePlayers
}

func spawnCellForSeat(seat int, spawns []int) int {
	if len(spawns) == 0 {
		return 0
	}
	if seat < 0 {
		seat = 0
	}
	if seat >= len(spawns) {
		return spawns[seat%len(spawns)]
	}
	return spawns[seat]
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

func pickCulprit(suspects []SuspectCard, rng RNG) SuspectCard {
	if len(suspects) == 0 {
		return SuspectCard{}
	}
	idx := rng.Intn(len(suspects))
	return suspects[idx]
}

func clueIDsFromTokens(clues []ClueToken) []string {
	ids := make([]string, 0, len(clues))
	for _, clue := range clues {
		ids = append(ids, clue.ID)
	}
	return ids
}

func bindCluesToBoard(clues []ClueToken, board BoardState) []ClueToken {
	indexByID := make(map[string]int, len(board.Cells))

	for _, cell := range board.Cells {
		if cell.ClueTokenID != "" {
			indexByID[cell.ClueTokenID] = cell.Index
		}
	}

	for i := range clues {
		if idx, ok := indexByID[clues[i].ID]; ok {
			clues[i].BoardCell = idx
		}
	}

	return clues
}

func newDefaultSuspects(rng RNG) []SuspectCard {
	namePool := []SuspectCode{
		"Анна", "Иван", "Мария", "Павел", "Елена", "Дмитрий", "Ольга", "Кирилл",
		"Наталья", "Алексей", "Вера", "Максим", "София", "Роман", "Дарья", "Тимур",
		"Лев", "Полина", "Никита", "Алина", "Глеб", "Виктория", "Степан", "Зоя",
		"Михаил", "Лидия", "Татьяна", "Арсений", "Валерия", "Игорь", "Нина", "Фёдор",
	}

	shuffleSuspectCodes(namePool, rng)

	allCombos := allSuspectTraitCombos()
	shuffleSuspectTraits(allCombos, rng)

	selectedCombos := make([]SuspectTraits, 16)
	copy(selectedCombos, allCombos[:16])

	suspects := make([]SuspectCard, 0, 16)
	for i := 0; i < 16; i++ {
		suspects = append(suspects, SuspectCard{
			ID:       suspectIDFromIndex(i),
			Code:     namePool[i],
			Revealed: false,
			Excluded: false,
			Traits:   selectedCombos[i],
		})
	}

	return suspects
}

func allSuspectTraitCombos() []SuspectTraits {
	res := make([]SuspectTraits, 0, 64)

	for mask := 0; mask < 64; mask++ {
		res = append(res, SuspectTraits{
			Glasses:  bitToTrait((mask >> 0) & 1),
			Hat:      bitToTrait((mask >> 1) & 1),
			Scarf:    bitToTrait((mask >> 2) & 1),
			Umbrella: bitToTrait((mask >> 3) & 1),
			Bag:      bitToTrait((mask >> 4) & 1),
			Boots:    bitToTrait((mask >> 5) & 1),
		})
	}

	return res
}

func bitToTrait(bit int) TraitValue {
	if bit == 1 {
		return TraitYes
	}
	return TraitNo
}

func shuffleSuspectCodes(items []SuspectCode, rng RNG) {
	for i := len(items) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

func shuffleSuspectTraits(items []SuspectTraits, rng RNG) {
	for i := len(items) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

func newDefaultClues(rng RNG) []ClueToken {
	traits := []ClueTrait{
		ClueTraitGlasses, ClueTraitGlasses,
		ClueTraitHat, ClueTraitHat,
		ClueTraitScarf, ClueTraitScarf,
		ClueTraitUmbrella, ClueTraitUmbrella,
		ClueTraitBag, ClueTraitBag,
		ClueTraitBoots, ClueTraitBoots,
	}

	shuffleClueTraits(traits, rng)

	clues := make([]ClueToken, 0, len(traits))
	for i, trait := range traits {
		clues = append(clues, ClueToken{
			ID:       "clue_" + strconv.Itoa(i+1),
			Trait:    trait,
			Revealed: false,
		})
	}

	return clues
}

func shuffleClueTraits(items []ClueTrait, rng RNG) {
	for i := len(items) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

func buildClueTruth(culprit SuspectCard, clues []ClueToken) map[string]TraitValue {
	res := make(map[string]TraitValue, len(clues))
	for _, clue := range clues {
		res[clue.ID] = traitValueForClue(culprit, clue.Trait)
	}
	return res
}

func traitValueForClue(culprit SuspectCard, trait ClueTrait) TraitValue {
	switch trait {
	case ClueTraitGlasses:
		return culprit.Traits.Glasses
	case ClueTraitHat:
		return culprit.Traits.Hat
	case ClueTraitScarf:
		return culprit.Traits.Scarf
	case ClueTraitUmbrella:
		return culprit.Traits.Umbrella
	case ClueTraitBag:
		return culprit.Traits.Bag
	case ClueTraitBoots:
		return culprit.Traits.Boots
	default:
		return TraitUnknown
	}
}
