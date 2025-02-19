package models

type Score struct {
	ID     string      `json:"id"`
	Score  interface{} `json:"score"`
	GameID string      `json:"game_id"`
	UserID string      `json:"user_id"`
}
