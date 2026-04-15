package domain

import "fmt"

type BoardCellType string

const (
	BoardCellStart BoardCellType = "start"
	BoardCellPath  BoardCellType = "path"
	BoardCellClue  BoardCellType = "clue"
)

type BoardState struct {
	Cells []BoardCell `json:"cells"`
}

type BoardCell struct {
	Index       int           `json:"index"`
	Type        BoardCellType `json:"type"`
	ClueTokenID string        `json:"clueTokenId,omitempty"`
}

type BoardView struct {
	Cells []BoardCellView `json:"cells"`
}

type BoardCellView struct {
	Index       int           `json:"index"`
	Type        BoardCellType `json:"type"`
	HasClue     bool          `json:"hasClue"`
	ClueTokenID string        `json:"clueTokenId,omitempty"`
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
