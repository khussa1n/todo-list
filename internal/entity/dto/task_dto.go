package dto

import "time"

type TasksDTO struct {
	title    string    `json:"title" db:"title" binding:"required"`
	activeAt time.Time `json:"activeAt" db:"active_at" binding:"required"`
}
