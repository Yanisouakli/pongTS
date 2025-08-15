package models

import "time"

type ConnectedEvent struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
}

type Player struct {
	PlayerID  string `json:"player_id"`
	Score     int64  `json:"score"`
	XPos      int64  `json:"x_pos"`
	YPos      int64  `json:"y_pos"`
	Height    int64  `json:"height"`
	Width     int64  `json:"width"`
	Direction string `json:"direction"`
	PreviousY int64  `json:"previousY"`
	VelocityY int64  `json:"velocityY"`
}

type BallState struct {
	XPos   int64 `json:"x_pos"`
	YPos   int64 `json:"y_pos"`
	Height int64 `json:"height"`
	Width  int64 `json:"width"`
}

type GameState struct {
	Ball BallState `json:"ball"`
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

type InitEvent struct {
	GameID     string    `json:"game_id"`
	PlayerInit Player    `json:"player_init"`
	BallInit   BallState `json:"ball_init"`
}

type SuccesInitEvent struct {
	Message string `json:"message"`
}
type ErrorEvent struct {
	Error string `json:"error"`
}

type WsEvent[T any] struct {
	Type   string `json:"type"`
	Params T      `json:"params"`
}
type InputEvent struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
	Key      string `json:"key"`
}

type GoalReturn struct {
	Goal   bool   `json:"goal"`
	Player string `json:"player"`
}

type UpdatesBody struct {
	Update string `json:"updated"`
}
