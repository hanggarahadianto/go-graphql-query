package models

import "time"

type Laptop struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
}