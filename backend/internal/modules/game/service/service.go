package service

import (
	"context"
	"encoding/json"
	"fmt"

	"fox/internal/modules/game/domain"
	"fox/internal/modules/game/repo"
)

type RNGFactory func() domain.RNG

type Service struct {
	repo repo.GameRepo
	rng  RNGFactory
}

func New(r repo.GameRepo, rng RNGFactory) *Service {
	return &Service{repo: r, rng: rng}
}

// ApplyCommand — главная операция: применить команду игрока к игре атомарно.
func (s *Service) ApplyCommand(ctx context.Context, gameID string, userID string, cmd domain.Command) (domain.GameState, []domain.Event, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return domain.GameState{}, nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Проверяем доступ
	ok, err := s.repo.IsPlayerInGame(ctx, tx, gameID, userID)
	if err != nil {
		return domain.GameState{}, nil, fmt.Errorf("check player in game: %w", err)
	}
	if !ok {
		return domain.GameState{}, nil, ErrForbidden
	}

	// Блокируем строку игры, чтобы два хода не применились параллельно
	row, err := s.repo.GetGameForUpdate(ctx, tx, gameID)
	if err != nil {
		return domain.GameState{}, nil, fmt.Errorf("get game for update: %w", err)
	}

	// Десериализуем state
	var st domain.GameState
	if err := json.Unmarshal(row.StateJSON, &st); err != nil {
		return domain.GameState{}, nil, fmt.Errorf("unmarshal state: %w", err)
	}

	// Применяем домен
	rng := s.rng()
	newState, events, err := domain.Apply(st, cmd, rng)
	if err != nil {
		return domain.GameState{}, nil, err
	}

	// Сериализуем новый state
	stateJSON, err := json.Marshal(newState)
	if err != nil {
		return domain.GameState{}, nil, fmt.Errorf("marshal state: %w", err)
	}

	// Сохраняем новый state + версию
	if err = s.repo.UpdateState(ctx, tx, gameID, stateJSON, newState.Version); err != nil {
		return domain.GameState{}, nil, fmt.Errorf("update state: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return domain.GameState{}, nil, fmt.Errorf("commit: %w", err)
	}

	return newState, events, nil
}

func (s *Service) GetState(ctx context.Context, gameID string, userID string) (domain.GameState, error) {
	ok, err := s.repo.IsPlayerInGameReadonly(ctx, gameID, userID)
	if err != nil {
		return domain.GameState{}, fmt.Errorf("check player in game: %w", err)
	}
	if !ok {
		return domain.GameState{}, ErrForbidden
	}

	row, err := s.repo.GetGame(ctx, gameID)
	if err != nil {
		return domain.GameState{}, fmt.Errorf("get game: %w", err)
	}

	var st domain.GameState
	if err := json.Unmarshal(row.StateJSON, &st); err != nil {
		return domain.GameState{}, fmt.Errorf("unmarshal state: %w", err)
	}

	return st, nil
}

func (s *Service) GetView(ctx context.Context, gameID string, userID string) (domain.GameView, error) {
	st, err := s.GetState(ctx, gameID, userID)
	if err != nil {
		return domain.GameView{}, err
	}

	return domain.BuildGameView(st, domain.PlayerID(userID)), nil
}
