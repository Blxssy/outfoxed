package http

type CreateGameRequest struct {
	Mode string `json:"mode"`
}

type JoinGameRequest struct {
	InviteCode string `json:"invite_code"`
}
