package models

import "time"

type ConnectedEvent struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
}

type Player struct {
	PlayerID string `json:"player_id"`
	Score    int64  `json:"score"`
	XPos     int64  `json:"x_pos"`
	YPos     int64  `json:"y_pos"`
}


type BallState struct{
	XPos     int64  `json:"x_pos"`
	YPos     int64  `json:"y_pos"`
}

type GameState struct{
  Ball    BallState  `json:"ball"`
}

type Game struct {
	GameID    string    `json:"game_id"`
	Players   []Player  `json:"players"`
	CreatedAt time.Time `json:"created_at"`
	State     GameState `json:"game_state"`
}

type LeftEvent struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
}

type WsEvent[T any] struct {
	Type   string `json:"type"`
	Params T      `json:"params"`
}
