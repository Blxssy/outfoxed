package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"fox/internal/modules/game/domain"
)

// DecodeCommand превращает WSRequest в domain.Command.
// actorID берём с сервера из токена
func DecodeCommand(req WSRequest, actorID domain.PlayerID) (domain.Command, error) {
	if req.Type != "command" {
		return nil, fmt.Errorf("invalid message type: %s", req.Type)
	}

	switch req.Command {
	case string(domain.CmdChooseGoal):
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

	case string(domain.CmdRollAuto):
		// payload может быть пустым
		return domain.RollAutoCommand{
			Player: actorID,
		}, nil

	case string(domain.CmdMovePawn):
		var p struct {
			Steps int `json:"steps"`
		}
		if err := json.Unmarshal(req.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}

		return domain.MovePawnCommand{
			Player: actorID,
			Steps:  p.Steps,
		}, nil

	case string(domain.CmdTakeClue):
		return domain.TakeClueCommand{
			Player: actorID,
		}, nil

	case string(domain.CmdRevealSuspects):
		var p struct {
			SuspectIDs []string `json:"suspectIds"`
		}
		if err := json.Unmarshal(req.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}

		return domain.RevealSuspectsCommand{
			Player:     actorID,
			SuspectIDs: p.SuspectIDs,
		}, nil

	case string(domain.CmdEndTurn):
		return domain.EndTurnCommand{
			Player: actorID,
		}, nil

	case string(domain.CmdAccuse):
		var p struct {
			SuspectID string `json:"suspectId"`
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
