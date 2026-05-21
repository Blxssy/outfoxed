package domain

import "fmt"

const (
	BoardWidth  = 16
	BoardHeight = 16
	BoardSize   = BoardWidth * BoardHeight
	ClueCount   = 6
)

type BoardCellType string

const (
	BoardCellStart BoardCellType = "start"
	BoardCellPath  BoardCellType = "path"
	BoardCellClue  BoardCellType = "clue"
)

type BoardState struct {
	Width  int         `json:"width"`
	Height int         `json:"height"`
	Cells  []BoardCell `json:"cells"`
}

type BoardCell struct {
	Index       int           `json:"index"`
	X           int           `json:"x"`
	Y           int           `json:"y"`
	Type        BoardCellType `json:"type"`
	ClueTokenID string        `json:"clueTokenId,omitempty"`
}

type BoardView struct {
	Width  int             `json:"width"`
	Height int             `json:"height"`
	Cells  []BoardCellView `json:"cells"`
}

type BoardCellView struct {
	Index       int           `json:"index"`
	X           int           `json:"x"`
	Y           int           `json:"y"`
	Type        BoardCellType `json:"type"`
	HasClue     bool          `json:"hasClue"`
	ClueTokenID string        `json:"clueTokenId,omitempty"`
}

func NewBoard16x16(clueIDs []string, rng RNG) (BoardState, error) {
	if len(clueIDs) != ClueCount {
		return BoardState{}, fmt.Errorf("expected %d clue ids, got %d", ClueCount, len(clueIDs))
	}

	board := BoardState{
		Width:  BoardWidth,
		Height: BoardHeight,
		Cells:  make([]BoardCell, 0, BoardSize),
	}

	startSet := make(map[int]struct{}, len(CenterStartZoneIndexes()))
	for _, idx := range CenterStartZoneIndexes() {
		startSet[idx] = struct{}{}
	}

	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			index := XYToIndex(x, y, BoardWidth)

			cellType := BoardCellPath
			if _, ok := startSet[index]; ok {
				cellType = BoardCellStart
			}

			board.Cells = append(board.Cells, BoardCell{
				Index: index,
				X:     x,
				Y:     y,
				Type:  cellType,
			})
		}
	}

	if err := placeCluesRandomly(&board, clueIDs, rng); err != nil {
		return BoardState{}, err
	}

	return board, nil
}

func placeCluesRandomly(board *BoardState, clueIDs []string, rng RNG) error {
	forbidden := make(map[int]struct{}, len(CenterStartZoneIndexes()))
	for _, idx := range CenterStartZoneIndexes() {
		forbidden[idx] = struct{}{}
	}

	available := make([]int, 0, len(board.Cells))
	for _, cell := range board.Cells {
		if _, bad := forbidden[cell.Index]; bad {
			continue
		}
		available = append(available, cell.Index)
	}

	if len(available) < len(clueIDs) {
		return fmt.Errorf("not enough cells to place clues")
	}

	shuffleInts(available, rng)

	for i, clueID := range clueIDs {
		idx := available[i]
		board.Cells[idx].Type = BoardCellClue
		board.Cells[idx].ClueTokenID = clueID
	}

	return nil
}

func shuffleInts(items []int, rng RNG) {
	for i := len(items) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

func CenterStartZoneIndexes() []int {
	// Центральный квадрат 4x4:
	// x = 6..9, y = 6..9
	result := make([]int, 0, 16)
	for y := 6; y <= 9; y++ {
		for x := 6; x <= 9; x++ {
			result = append(result, XYToIndex(x, y, BoardWidth))
		}
	}
	return result
}

func DefaultSpawnCells() []int {
	// Игроки стартуют в центральном квадрате 2x2 внутри стартовой зоны 4x4:
	// (7,7), (8,7), (7,8), (8,8)
	return []int{
		XYToIndex(7, 7, BoardWidth), // 119
		XYToIndex(8, 7, BoardWidth), // 120
		XYToIndex(7, 8, BoardWidth), // 135
		XYToIndex(8, 8, BoardWidth), // 136
	}
}

func XYToIndex(x, y, width int) int {
	return y*width + x
}

func IndexToXY(index, width int) (x, y int) {
	return index % width, index / width
}

func (b BoardState) IsInside(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

func (b BoardState) CellAt(index int) (BoardCell, bool) {
	if index < 0 || index >= len(b.Cells) {
		return BoardCell{}, false
	}
	return b.Cells[index], true
}

func (b BoardState) MustCellAt(index int) (BoardCell, error) {
	cell, ok := b.CellAt(index)
	if !ok {
		return BoardCell{}, fmt.Errorf("board cell %d out of range", index)
	}
	return cell, nil
}

func (b BoardState) LastIndex() int {
	if len(b.Cells) == 0 {
		return 0
	}
	return len(b.Cells) - 1
}

func (b BoardState) ClampIndex(index int) int {
	if index < 0 {
		return 0
	}
	last := b.LastIndex()
	if index > last {
		return last
	}
	return index
}

func (b BoardState) HasClueAt(index int) bool {
	cell, ok := b.CellAt(index)
	if !ok {
		return false
	}
	return cell.Type == BoardCellClue && cell.ClueTokenID != ""
}

func ManhattanDistance(fromIndex, toIndex, width int) int {
	fx, fy := IndexToXY(fromIndex, width)
	tx, ty := IndexToXY(toIndex, width)

	dx := fx - tx
	if dx < 0 {
		dx = -dx
	}

	dy := fy - ty
	if dy < 0 {
		dy = -dy
	}

	return dx + dy
}

func (b BoardState) ReachableWithin(fromIndex, steps int, occupied map[int]struct{}) []int {
	if steps <= 0 {
		return nil
	}

	result := make([]int, 0)

	for _, cell := range b.Cells {
		if cell.Index == fromIndex {
			continue
		}

		if _, blocked := occupied[cell.Index]; blocked {
			continue
		}

		dist := ManhattanDistance(fromIndex, cell.Index, b.Width)
		if dist > 0 && dist <= steps {
			result = append(result, cell.Index)
		}
	}

	return result
}
func occupiedCells(players []PlayerState, exceptSeat int) map[int]struct{} {
	res := make(map[int]struct{})
	for _, p := range players {
		if p.Seat == exceptSeat {
			continue
		}
		res[p.PawnCell] = struct{}{}
	}
	return res
}

func computeReachableCells(board BoardState, from int, steps int, occupied map[int]struct{}) []int {
	return board.ReachableWithin(from, steps, occupied)
}

func containsInt(items []int, target int) bool {
	for _, v := range items {
		if v == target {
			return true
		}
	}
	return false
}

func activePlayerIndex(players []PlayerState, activeSeat int) int {
	for i := range players {
		if players[i].Seat == activeSeat {
			return i
		}
	}
	return -1
}
