package dto

type TasksDTO struct {
	Title    string `json:"title" binding:"required"`
	ActiveAt string `json:"activeAt" binding:"required"`
}
