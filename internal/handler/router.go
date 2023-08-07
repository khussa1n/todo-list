package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/todo-list")

	task := api.Group("/tasks")
	task.POST("/", h.createTask)
	task.GET("/", h.getAllTasks)
	task.PUT("/", h.updateTask)
	task.DELETE("/:id", h.deleteTask)

	return router
}
