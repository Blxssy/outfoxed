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
	if len(clueIDs) == 0 {
		return BoardState{}, fmt.Errorf("expected at least 1 clue id, got 0")
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

type quadrant int

const (
	qTopLeft quadrant = iota
	qTopRight
	qBottomLeft
	qBottomRight
)

type clueRing int

const (
	ringNear clueRing = iota // 2..3
	ringMid                  // 4..5
	ringFar                  // 6..7
)

type clueSlot struct {
	ring clueRing
	quad quadrant
}

func placeCluesRandomly(board *BoardState, clueIDs []string, rng RNG) error {
	if len(clueIDs) != 6 {
		return fmt.Errorf("expected 6 clue ids, got %d", len(clueIDs))
	}

	forbidden := make(map[int]struct{}, len(CenterStartZoneIndexes()))
	for _, idx := range CenterStartZoneIndexes() {
		forbidden[idx] = struct{}{}
	}

	buckets := map[clueRing]map[quadrant][]int{
		ringNear: {
			qTopLeft:     {},
			qTopRight:    {},
			qBottomLeft:  {},
			qBottomRight: {},
		},
		ringMid: {
			qTopLeft:     {},
			qTopRight:    {},
			qBottomLeft:  {},
			qBottomRight: {},
		},
		ringFar: {
			qTopLeft:     {},
			qTopRight:    {},
			qBottomLeft:  {},
			qBottomRight: {},
		},
	}

	for _, cell := range board.Cells {
		if _, bad := forbidden[cell.Index]; bad {
			continue
		}

		dist := distanceToStartZone(cell.Index)
		q := cellQuadrant(cell.X, cell.Y)

		switch {
		case dist >= 2 && dist <= 3:
			buckets[ringNear][q] = append(buckets[ringNear][q], cell.Index)
		case dist >= 4 && dist <= 5:
			buckets[ringMid][q] = append(buckets[ringMid][q], cell.Index)
		case dist >= 6 && dist <= 7:
			buckets[ringFar][q] = append(buckets[ringFar][q], cell.Index)
		}
	}

	slots, err := buildClueSlots(buckets, rng)
	if err != nil {
		return err
	}

	picked := make([]int, 0, len(clueIDs))

	for _, slot := range slots {
		candidates := buckets[slot.ring][slot.quad]
		if len(candidates) == 0 {
			return fmt.Errorf("no candidates for ring=%d quadrant=%d", slot.ring, slot.quad)
		}

		idx := pickBestCell(candidates, picked, board.Width, rng)
		picked = append(picked, idx)

		buckets[slot.ring][slot.quad] = removeCell(buckets[slot.ring][slot.quad], idx)
	}

	for i, clueID := range clueIDs {
		idx := picked[i]
		board.Cells[idx].Type = BoardCellClue
		board.Cells[idx].ClueTokenID = clueID
	}

	return nil
}

func buildClueSlots(
	buckets map[clueRing]map[quadrant][]int,
	rng RNG,
) ([]clueSlot, error) {
	rings := []clueRing{
		ringNear, ringNear,
		ringMid, ringMid,
		ringFar, ringFar,
	}

	qOrder := []quadrant{qTopLeft, qTopRight, qBottomLeft, qBottomRight}
	shuffleQuadrants(qOrder, rng)

	usedQuadrants := map[quadrant]int{
		qTopLeft:     0,
		qTopRight:    0,
		qBottomLeft:  0,
		qBottomRight: 0,
	}

	slots := make([]clueSlot, 0, len(rings))

	var dfs func(pos int) bool
	dfs = func(pos int) bool {
		if pos == len(rings) {
			for _, q := range qOrder {
				if usedQuadrants[q] == 0 {
					return false
				}
			}
			return true
		}

		slotsLeft := len(rings) - pos
		uncovered := 0
		for _, q := range qOrder {
			if usedQuadrants[q] == 0 {
				uncovered++
			}
		}
		if uncovered > slotsLeft {
			return false
		}

		r := rings[pos]

		// Сначала пробуем квадранты, в которых ещё нет улик
		for pass := 0; pass < 2; pass++ {
			for _, q := range qOrder {
				if pass == 0 && usedQuadrants[q] > 0 {
					continue
				}
				if pass == 1 && usedQuadrants[q] == 0 {
					continue
				}
				if len(buckets[r][q]) == 0 {
					continue
				}

				usedQuadrants[q]++
				slots = append(slots, clueSlot{ring: r, quad: q})

				if dfs(pos + 1) {
					return true
				}

				slots = slots[:len(slots)-1]
				usedQuadrants[q]--
			}
		}

		return false
	}

	if !dfs(0) {
		return nil, fmt.Errorf("failed to build clue slots with ring/quadrant constraints")
	}

	return slots, nil
}

func distanceToStartZone(index int) int {
	x, y := IndexToXY(index, BoardWidth)

	dx := 0
	switch {
	case x < 6:
		dx = 6 - x
	case x > 9:
		dx = x - 9
	}

	dy := 0
	switch {
	case y < 6:
		dy = 6 - y
	case y > 9:
		dy = y - 9
	}

	return dx + dy
}

func cellQuadrant(x, y int) quadrant {
	left := x < BoardWidth/2
	top := y < BoardHeight/2

	switch {
	case left && top:
		return qTopLeft
	case !left && top:
		return qTopRight
	case left && !top:
		return qBottomLeft
	default:
		return qBottomRight
	}
}

func shuffleQuadrants(items []quadrant, rng RNG) {
	for i := len(items) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}

func pickBestCell(candidates []int, picked []int, width int, rng RNG) int {
	if len(candidates) == 0 {
		return 0
	}

	if len(picked) == 0 {
		return candidates[rng.Intn(len(candidates))]
	}

	bestIdx := candidates[0]
	bestScore := -1

	for _, candidate := range candidates {
		score := minDistanceToPicked(candidate, picked, width)
		if score > bestScore {
			bestScore = score
			bestIdx = candidate
		}
	}

	return bestIdx
}

func minDistanceToPicked(candidate int, picked []int, width int) int {
	minDist := ManhattanDistance(candidate, picked[0], width)

	for _, p := range picked[1:] {
		d := ManhattanDistance(candidate, p, width)
		if d < minDist {
			minDist = d
		}
	}

	return minDist
}

func removeCell(items []int, target int) []int {
	for i, v := range items {
		if v == target {
			return append(items[:i], items[i+1:]...)
		}
	}
	return items
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
