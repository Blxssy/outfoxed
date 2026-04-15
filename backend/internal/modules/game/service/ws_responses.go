package service

import (
	"encoding/json"

	"fox/internal/modules/game/domain"
)

// WSRequest — сообщение от клиента по WebSocket.
type WSRequest struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`    // "command"
	Command string          `json:"command"` // choose_goal, roll_auto, move_pawn, take_clue, reveal_suspects, end_turn, accuse
	Payload json.RawMessage `json:"payload"`
}

// WSError — стандартная ошибка для клиента.
type WSError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type WSResponse struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"` // "update" | "error"
	Payload any    `json:"payload"`
}

type UpdatePayload struct {
	State  domain.GameView `json:"state"`
	Events []domain.Event  `json:"events"`
}

func NewUpdateResponse(reqID string, view domain.GameView, events []domain.Event) WSResponse {
	return WSResponse{
		ID:   reqID,
		Type: "update",
		Payload: UpdatePayload{
			State:  view,
			Events: events,
		},
	}
}

func NewErrorResponse(reqID string, code, message string) WSResponse {
	return WSResponse{
		ID:   reqID,
		Type: "error",
		Payload: WSError{
			Code:    code,
			Message: message,
		},
	}
}

func ErrorResponse(reqID string, err error) WSResponse {
	wsErr := ToWSError(err)
	return NewErrorResponse(reqID, wsErr.Code, wsErr.Message)
}
