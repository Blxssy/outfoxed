package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"fox/internal/modules/game/domain"
	"fox/internal/modules/game/repo"
)

type GameSummary struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type PlayerSummary struct {
	UserID string `json:"user_id"`
	Seat   int    `json:"seat"`
}

type CreateGameResult struct {
	Game   GameSummary   `json:"game"`
	Player PlayerSummary `json:"player"`
}

type JoinGameResult struct {
	Game   GameSummary   `json:"game"`
	Player PlayerSummary `json:"player"`
}

type LobbyPlayer struct {
	UserID      string `json:"user_id"`
	Seat        int    `json:"seat"`
	DisplayName string `json:"display_name"`
	IsMe        bool   `json:"is_me"`
}

type LobbyGame struct {
	ID         string        `json:"id"`
	Status     string        `json:"status"`
	Players    []LobbyPlayer `json:"players"`
	CanStart   bool          `json:"can_start"`
	MinPlayers int           `json:"min_players"`
	MaxPlayers int           `json:"max_players"`
}

type LobbySnapshot struct {
	Game LobbyGame `json:"game"`
}

type StartResult struct {
	Game     GameSummary `json:"game"`
	Redirect struct {
		Route string `json:"route"`
	} `json:"redirect"`
}

type StateResult struct {
	State domain.GameView `json:"state"`
}

func (s *Service) CreateGame(ctx context.Context, userID string) (CreateGameResult, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return CreateGameResult{}, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := s.repo.FindUnfinishedGameForUser(ctx, tx, userID); err == nil {
		return CreateGameResult{}, ErrAlreadyInAnotherGame
	} else if err != sql.ErrNoRows {
		return CreateGameResult{}, fmt.Errorf("check user unfinished game: %w", err)
	}

	row, err := s.repo.CreateGame(ctx, tx, userID, []byte(`{}`))
	if err != nil {
		return CreateGameResult{}, fmt.Errorf("create game: %w", err)
	}

	if err := s.repo.AddPlayer(ctx, tx, row.ID, userID, 0); err != nil {
		return CreateGameResult{}, fmt.Errorf("add creator to game: %w", err)
	}

	players, err := s.repo.GetPlayersForUpdate(ctx, tx, row.ID)
	if err != nil {
		return CreateGameResult{}, fmt.Errorf("get game players: %w", err)
	}

	stateJSON, waitingState, err := buildWaitingStateJSON(row.ID, players)
	if err != nil {
		return CreateGameResult{}, err
	}

	if err := s.repo.UpdateState(ctx, tx, row.ID, string(waitingState.Status), stateJSON, waitingState.Version); err != nil {
		return CreateGameResult{}, fmt.Errorf("update waiting state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return CreateGameResult{}, fmt.Errorf("commit: %w", err)
	}

	return CreateGameResult{
		Game: GameSummary{
			ID:     row.ID,
			Status: string(waitingState.Status),
		},
		Player: PlayerSummary{
			UserID: userID,
			Seat:   0,
		},
	}, nil
}

func (s *Service) JoinGame(ctx context.Context, gameID string, userID string) (JoinGameResult, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return JoinGameResult{}, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	row, err := s.repo.GetGameForUpdate(ctx, tx, gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return JoinGameResult{}, ErrGameNotFound
		}
		return JoinGameResult{}, fmt.Errorf("get game for update: %w", err)
	}

	players, err := s.repo.GetPlayersForUpdate(ctx, tx, gameID)
	if err != nil {
		return JoinGameResult{}, fmt.Errorf("get game players: %w", err)
	}

	for _, player := range players {
		if player.UserID == userID {
			if err := tx.Commit(); err != nil {
				return JoinGameResult{}, fmt.Errorf("commit: %w", err)
			}
			return JoinGameResult{
				Game: GameSummary{
					ID:     row.ID,
					Status: row.Status,
				},
				Player: PlayerSummary{
					UserID: userID,
					Seat:   player.Seat,
				},
			}, nil
		}
	}

	existingGame, err := s.repo.FindUnfinishedGameForUser(ctx, tx, userID)
	if err == nil && existingGame.ID != gameID {
		return JoinGameResult{}, ErrAlreadyInAnotherGame
	} else if err != nil && err != sql.ErrNoRows {
		return JoinGameResult{}, fmt.Errorf("check user unfinished game: %w", err)
	}

	if row.Status != string(domain.StatusWaiting) {
		return JoinGameResult{}, ErrGameAlreadyStarted
	}
	if len(players) >= domain.MaxPlayers {
		return JoinGameResult{}, ErrGameFull
	}

	seat := nextFreeSeat(players)
	if err := s.repo.AddPlayer(ctx, tx, gameID, userID, seat); err != nil {
		return JoinGameResult{}, fmt.Errorf("add player to game: %w", err)
	}

	updatedPlayers, err := s.repo.GetPlayersForUpdate(ctx, tx, gameID)
	if err != nil {
		return JoinGameResult{}, fmt.Errorf("get updated players: %w", err)
	}

	stateJSON, waitingState, err := buildWaitingStateJSON(gameID, updatedPlayers)
	if err != nil {
		return JoinGameResult{}, err
	}

	if err := s.repo.UpdateState(ctx, tx, gameID, string(waitingState.Status), stateJSON, waitingState.Version); err != nil {
		return JoinGameResult{}, fmt.Errorf("update waiting state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return JoinGameResult{}, fmt.Errorf("commit: %w", err)
	}

	return JoinGameResult{
		Game: GameSummary{
			ID:     gameID,
			Status: string(waitingState.Status),
		},
		Player: PlayerSummary{
			UserID: userID,
			Seat:   seat,
		},
	}, nil
}

func (s *Service) GetLobby(ctx context.Context, gameID string, userID string) (LobbySnapshot, error) {
	row, err := s.repo.GetGame(ctx, gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return LobbySnapshot{}, ErrGameNotFound
		}
		return LobbySnapshot{}, fmt.Errorf("get game: %w", err)
	}

	players, err := s.repo.GetPlayers(ctx, gameID)
	if err != nil {
		return LobbySnapshot{}, fmt.Errorf("get players: %w", err)
	}

	if !hasPlayer(players, userID) {
		return LobbySnapshot{}, ErrForbidden
	}

	lobbyPlayers := make([]LobbyPlayer, 0, len(players))
	for _, player := range players {
		lobbyPlayers = append(lobbyPlayers, LobbyPlayer{
			UserID:      player.UserID,
			Seat:        player.Seat,
			DisplayName: player.Username,
			IsMe:        player.UserID == userID,
		})
	}

	canStart := row.Status == string(domain.StatusWaiting) &&
		len(players) >= domain.MinPlayers &&
		row.CreatedBy.Valid &&
		row.CreatedBy.String == userID

	return LobbySnapshot{
		Game: LobbyGame{
			ID:         row.ID,
			Status:     row.Status,
			Players:    lobbyPlayers,
			CanStart:   canStart,
			MinPlayers: domain.MinPlayers,
			MaxPlayers: domain.MaxPlayers,
		},
	}, nil
}

func (s *Service) StartGame(ctx context.Context, gameID string, userID string) (StartResult, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return StartResult{}, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	row, err := s.repo.GetGameForUpdate(ctx, tx, gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return StartResult{}, ErrGameNotFound
		}
		return StartResult{}, fmt.Errorf("get game for update: %w", err)
	}
	if row.Status != string(domain.StatusWaiting) {
		return StartResult{}, ErrGameAlreadyStarted
	}
	if !row.CreatedBy.Valid || row.CreatedBy.String != userID {
		return StartResult{}, ErrOnlyCreatorCanStart
	}

	players, err := s.repo.GetPlayersForUpdate(ctx, tx, gameID)
	if err != nil {
		return StartResult{}, fmt.Errorf("get players for update: %w", err)
	}
	if !hasPlayer(players, userID) {
		return StartResult{}, ErrForbidden
	}
	if len(players) < domain.MinPlayers {
		return StartResult{}, ErrNotEnoughPlayers
	}

	activeState := domain.NewActiveGameState(gameID, toSetupPlayers(players))
	stateJSON, err := json.Marshal(activeState)
	if err != nil {
		return StartResult{}, fmt.Errorf("marshal active state: %w", err)
	}

	if err := s.repo.UpdateState(ctx, tx, gameID, string(activeState.Status), stateJSON, activeState.Version); err != nil {
		return StartResult{}, fmt.Errorf("update active state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return StartResult{}, fmt.Errorf("commit: %w", err)
	}

	result := StartResult{
		Game: GameSummary{
			ID:     gameID,
			Status: string(activeState.Status),
		},
	}
	result.Redirect.Route = "/game/" + gameID
	return result, nil
}

func (s *Service) GetViewState(ctx context.Context, gameID string, userID string) (StateResult, error) {
	view, err := s.GetView(ctx, gameID, userID)
	if err != nil {
		return StateResult{}, err
	}

	return StateResult{State: view}, nil
}

func buildWaitingStateJSON(gameID string, players []repo.GamePlayerRow) ([]byte, domain.GameState, error) {
	waitingState := domain.NewWaitingGameState(gameID, toSetupPlayers(players))
	stateJSON, err := json.Marshal(waitingState)
	if err != nil {
		return nil, domain.GameState{}, fmt.Errorf("marshal waiting state: %w", err)
	}
	return stateJSON, waitingState, nil
}

func toSetupPlayers(players []repo.GamePlayerRow) []domain.SetupPlayer {
	setupPlayers := make([]domain.SetupPlayer, 0, len(players))
	for _, player := range players {
		setupPlayers = append(setupPlayers, domain.SetupPlayer{
			UserID: domain.PlayerID(player.UserID),
			Name:   player.Username,
			Seat:   player.Seat,
		})
	}
	return setupPlayers
}

func nextFreeSeat(players []repo.GamePlayerRow) int {
	occupied := make(map[int]struct{}, len(players))
	for _, player := range players {
		occupied[player.Seat] = struct{}{}
	}
	for seat := 0; seat < domain.MaxPlayers; seat++ {
		if _, ok := occupied[seat]; !ok {
			return seat
		}
	}
	return len(players)
}

func hasPlayer(players []repo.GamePlayerRow, userID string) bool {
	for _, player := range players {
		if player.UserID == userID {
			return true
		}
	}
	return false
}
