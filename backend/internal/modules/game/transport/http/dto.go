package http

type CreateGameRequest struct {
	Title      string `json:"title"`
	Visibility string `json:"visibility"` // public | private
}

type JoinGameRequest struct {
	InviteCode string `json:"invite_code"`
}

type GamesListResponse struct {
	Games []LobbyItem `json:"games"`
}

type LobbyItem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	HostUserID   string `json:"hostUserId"`
	PlayersCount int    `json:"playersCount"`
	MaxPlayers   int    `json:"maxPlayers"`
	Status       string `json:"status"`
}

type JoinByCodeRequest struct {
	Code string `json:"code"`
}

type PublicGameItem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	HostUsername string `json:"host_username"`
	PlayersCount int    `json:"players_count"`
	MaxPlayers   int    `json:"max_players"`
	Status       string `json:"status"`
}

type GamesListResult struct {
	Games []PublicGameItem `json:"games"`
}
