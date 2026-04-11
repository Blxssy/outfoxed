package domain

type FoxState struct {
	Track    int `json:"track"`
	EscapeAt int `json:"escapeAt"`
}

type FoxView struct {
	Track    int `json:"track"`
	EscapeAt int `json:"escapeAt"`
}
