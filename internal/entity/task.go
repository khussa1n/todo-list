package entity

import "time"

type Tasks struct {
	ID       int64     `json:"id" db:"id"`
	title    string    `json:"title" db:"title" binding:"required"`
	activeAt time.Time `json:"active_at" db:"active_at" binding:"required"`
	status   string    `json:"status" db:"active_at"`
}
