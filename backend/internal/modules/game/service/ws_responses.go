package service

import "fox/internal/modules/game/domain"

func NewUpdateResponse(reqID string, state domain.GameState, events []domain.Event) WSResponse {
	return WSResponse{
		ID:   reqID,
		Type: "update",
		Payload: UpdatePayload{
			State:  state,
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
