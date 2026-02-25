package service

import (
	"errors"

	"fox/internal/modules/game/domain"
)

func ToWSError(err error) WSError {
	switch {
	case errors.Is(err, domain.ErrNotYourTurn):
		return WSError{Code: "not_your_turn", Message: "It is not your turn."}

	case errors.Is(err, domain.ErrInvalidPhase):
		return WSError{Code: "invalid_phase", Message: "This action is not allowed in the current game phase."}

	case errors.Is(err, domain.ErrGameNotActive):
		return WSError{Code: "game_not_active", Message: "The game is not active."}

	case errors.Is(err, domain.ErrGameFinished):
		return WSError{Code: "game_finished", Message: "The game has already finished."}

	case errors.Is(err, domain.ErrGoalAlreadySet):
		return WSError{Code: "goal_already_set", Message: "Goal has already been chosen for this turn."}

	case errors.Is(err, domain.ErrGoalNotSet):
		return WSError{Code: "goal_not_set", Message: "Choose a goal before rolling."}

	case errors.Is(err, domain.ErrNoPendingAction):
		return WSError{Code: "no_pending_action", Message: "No action is pending right now."}

	case errors.Is(err, domain.ErrPendingNotClue):
		return WSError{Code: "pending_not_clue", Message: "The pending action is not a clue action."}

	case errors.Is(err, domain.ErrPendingNotSuspect):
		return WSError{Code: "pending_not_suspect", Message: "The pending action is not a suspect action."}

	case errors.Is(err, domain.ErrAllCluesCollected):
		return WSError{Code: "all_clues_collected", Message: "All clues have already been collected."}

	case errors.Is(err, domain.ErrSuspectNotRevealed):
		return WSError{Code: "suspect_not_revealed", Message: "You can only accuse a revealed suspect."}

	case errors.Is(err, domain.ErrSuspectExcluded):
		return WSError{Code: "suspect_excluded", Message: "This suspect has been excluded."}

	case errors.Is(err, ErrForbidden):
		return WSError{Code: "forbidden", Message: "You do not have access to this game."}
	}

	// Любая другая ошибка (500).
	return WSError{Code: "internal_error", Message: "Internal server error."}
}
