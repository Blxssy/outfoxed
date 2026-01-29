package domain

import "errors"

var (
	ErrGameNotActive  = errors.New("game is not active")
	ErrNotYourTurn    = errors.New("not your turn")
	ErrInvalidPhase   = errors.New("invalid phase for command")
	ErrGoalAlreadySet = errors.New("goal already set")
	ErrGoalNotSet     = errors.New("goal not set")
)
