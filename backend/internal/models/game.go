package models

import "time"

type Game struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
