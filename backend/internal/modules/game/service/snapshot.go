package service

import (
	"fox/internal/modules/game/domain"
	"fox/internal/modules/game/repo"
)

func BuildLobbySnapshot(row repo.GameRow, players []repo.GamePlayerRow, userID string) LobbySnapshot {
	lobbyPlayers := make([]LobbyPlayer, 0, len(players))
	hostUsername := ""

	for _, player := range players {
		lobbyPlayers = append(lobbyPlayers, LobbyPlayer{
			UserID:      player.UserID,
			Seat:        player.Seat,
			DisplayName: player.Username,
			IsMe:        player.UserID == userID,
		})

		if row.CreatedBy.Valid && row.CreatedBy.String == player.UserID {
			hostUsername = player.Username
		}
	}

	canStart := row.Status == string(domain.StatusWaiting) &&
		len(players) >= domain.MinPlayers &&
		row.CreatedBy.Valid &&
		row.CreatedBy.String == userID

	game := LobbyGame{
		ID:           row.ID,
		Title:        row.Title,
		Status:       row.Status,
		Visibility:   row.Visibility,
		HostUsername: hostUsername,
		Players:      lobbyPlayers,
		CanStart:     canStart,
		MinPlayers:   domain.MinPlayers,
		MaxPlayers:   domain.MaxPlayers,
	}

	if row.JoinCode.Valid && row.Visibility == "private" {
		code := row.JoinCode.String
		game.JoinCode = &code
	}

	return LobbySnapshot{
		Game: game,
	}
}
