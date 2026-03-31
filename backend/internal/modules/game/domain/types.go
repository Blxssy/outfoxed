package domain

import "time"

type GameStatus string

const (
	StatusWaiting  GameStatus = "waiting"
	StatusActive   GameStatus = "active"
	StatusFinished GameStatus = "finished"
)

type Phase string

const (
	PhaseChooseGoal Phase = "choose_goal"
	PhaseRolling    Phase = "rolling"
	PhaseAction     Phase = "action"
	PhaseEndTurn    Phase = "end_turn"
)

type GoalType string

const (
	GoalClue    GoalType = "clue"
	GoalSuspect GoalType = "suspect"
)

type PlayerID string
type GameID string

type NowFunc func() time.Time
