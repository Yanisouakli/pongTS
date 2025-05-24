package models

type ConnectedEvent struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
}

type LeftEvent struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
}

type WsEvent[T any] struct {
	Type   string `json:"type"`
	Params T      `json:"params"`
}
