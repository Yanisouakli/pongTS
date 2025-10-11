package models

import "time"

type ConnectedEvent struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
}

type Player struct {
	PlayerID  string `json:"player_id"`
	XPos      int64  `json:"x_pos"`
	YPos      int64  `json:"y_pos"`
	Direction string `json:"direction"`
	PreviousY int64  `json:"previous_y"`
	VelocityY int64  `json:"velocity_y"`
	Width     int64  `json:"width"`
	Height    int64  `json:"height"`
}

type BallState struct {
	XPos      int64 `json:"x_pos"`
	YPos      int64 `json:"y_pos"`
	Height    int64 `json:"height"`
	Width     int64 `json:"width"`
	VelocityY int64 `json:"velocity_y"`
	VelocityX int64 `json:"velocity_x"`
}

type GameState struct {
	Ball   BallState     `json:"ball"`
	Score  int64         `json:"score"`
	Timer  time.Duration `json:"timer"`
	Canvas Canvas        `json:"canvas"`
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

type Canvas struct {
	CanvasWidth  int64 `json:"canvas_width"`
	CanvasHeight int64 `json:"canvas_height"`
}

type InitEvent struct {
	GameID     string `json:"game_id"`
	PlayerInit Player `json:"player_init"`
	CanvasInit Canvas `json:"canvas"`
}

type StartEvent struct {
	GameID     string    `json:"game_id"`
	PlayerInit Player    `json:"player_init"`
	CanvasInit Canvas    `json:"canvas"`
	Ball       BallState `json:"ball"`
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
	Type    string        `json:"type"`
	Players []Player      `json:"players"`
	Ball    BallState     `json:"ball"`
	Score   int64         `json:"score"`
	Timer   time.Duration `json:"timer"`
	Canvas  Canvas        `json:"canvas"`
}
