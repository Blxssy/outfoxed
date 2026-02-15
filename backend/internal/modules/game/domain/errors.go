package domain

import "errors"

var (
	ErrGameNotActive  = errors.New("game is not active")
	ErrNotYourTurn    = errors.New("not your turn")
	ErrInvalidPhase   = errors.New("invalid phase for command")
	ErrGoalAlreadySet = errors.New("goal already set")
	ErrGoalNotSet     = errors.New("goal not set")

	ErrNoPendingAction   = errors.New("no pending action to resolve")
	ErrPendingNotClue    = errors.New("pending action is not clue")
	ErrAllCluesCollected = errors.New("all clues already collected")

	ErrPendingNotSuspect  = errors.New("pending action is not suspect")
	ErrNoSuspectsToReveal = errors.New("no more suspects to reveal")
)
