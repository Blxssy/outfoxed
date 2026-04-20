package service

import "errors"

var (
	ErrForbidden            = errors.New("forbidden")
	ErrGameNotFound         = errors.New("game not found")
	ErrGameFull             = errors.New("game is full")
	ErrGameAlreadyStarted   = errors.New("game already started")
	ErrNotEnoughPlayers     = errors.New("not enough players")
	ErrOnlyCreatorCanStart  = errors.New("only the game creator can start the game")
	ErrAlreadyInAnotherGame = errors.New("user is already in another game")
)
