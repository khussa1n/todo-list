package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/khussa1n/todo-list/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/todo-list")

	task := api.Group("/tasks")
	task.POST("/", h.createTask)
	task.PUT("/:id", h.updateTask)
	task.DELETE("/:id", h.deleteTask)
	task.PUT("/:id/done", h.updateTaskStatus)
	task.GET("", h.getAllTasks)

	return router
}
