package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/todo-list")

	task := api.Group("/tasks")
	task.POST("/", h.createTask)
	task.PUT("/:id", h.updateTask)
	task.DELETE("/:id", h.deleteTask)
	task.PUT("/:id/done", h.updateTaskStatus)
	task.GET("", h.getAllTasks)

	return router
}
