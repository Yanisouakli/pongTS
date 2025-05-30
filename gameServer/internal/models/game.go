package models

type ConnectedEvent struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
}

type Player struct {
	PlayerID string `json:"player_id"`
	score    int64  `json:"x_pos"`
	XPos     int64  `json:"x_pos"`
	YPos     int64  `json:"y_pos"`
}

type Game struct {
	GameID  string   `json:"game_id"`
	Players []Player `json:"players"`
}

type LeftEvent struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
}

type WsEvent[T any] struct {
	Type   string `json:"type"`
	Params T      `json:"params"`
}
