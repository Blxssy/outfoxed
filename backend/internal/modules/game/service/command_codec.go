package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"fox/internal/modules/game/domain"
)

// WSRequest — сообщение от клиента по WebSocket.
type WSRequest struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`    // "command"
	Command string          `json:"command"` // choose_goal, roll_auto, take_clue, reveal_suspects, end_turn, accuse
	Payload json.RawMessage `json:"payload"`
}

// WSError — стандартная ошибка для клиента
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
	State  domain.GameState `json:"state"`
	Events []domain.Event   `json:"events"`
}

// DecodeCommand превращает WSRequest в `domain.Command`.
func DecodeCommand(req WSRequest, actorID domain.PlayerID) (domain.Command, error) {
	if req.Type != "command" {
		return nil, fmt.Errorf("invalid message type: %s", req.Type)
	}

	switch req.Command {
	case "choose_goal":
		var p struct {
			Goal string `json:"goal"`
		}
		if err := json.Unmarshal(req.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		goal, err := parseGoalType(p.Goal)
		if err != nil {
			return nil, err
		}
		return domain.ChooseGoalCommand{
			Player: actorID,
			Goal:   goal,
		}, nil

	case "roll_auto":
		// payload может быть пустым
		return domain.RollAutoCommand{Player: actorID}, nil

	case "take_clue":
		return domain.TakeClueCommand{Player: actorID}, nil

	case "reveal_suspects":
		return domain.RevealSuspectsCommand{Player: actorID}, nil

	case "end_turn":
		return domain.EndTurnCommand{Player: actorID}, nil

	case "accuse":
		var p struct {
			SuspectID int `json:"suspectId"`
		}
		if err := json.Unmarshal(req.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		return domain.AccuseCommand{
			Player:    actorID,
			SuspectID: p.SuspectID,
		}, nil

	default:
		return nil, fmt.Errorf("unknown command: %s", req.Command)
	}
}

func parseGoalType(s string) (domain.GoalType, error) {
	switch s {
	case string(domain.GoalClue):
		return domain.GoalClue, nil
	case string(domain.GoalSuspect):
		return domain.GoalSuspect, nil
	default:
		return "", errors.New("invalid goal type")
	}
}
