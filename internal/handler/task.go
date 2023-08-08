package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/khussa1n/todo-list/internal/entity/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func (h *Handler) createTask(ctx *gin.Context) {
	var req dto.TasksDTO
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("bind json err: %s \n", err.Error())
		ctx.JSON(http.StatusBadRequest, &dto.Error{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	task, err := h.srvs.CreateTask(ctx, &req)
	if err != nil {
		log.Printf("can not create task: %s \n", err.Error())
		if err.Error() == "more than 200 char" || err.Error() == "activeAt invalid format" || err.Error() == "a task "+
			"with the same title and activeAt already exists" || err.Error() == "the provided hex string is "+
			"not a valid ObjectID" {
			ctx.JSON(http.StatusBadRequest, &dto.Error{
				Code:    -2,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, &dto.Error{
			Code:    -3,
			Message: "invalid to create task",
		})
		return
	}

	ctx.JSON(http.StatusCreated, task)
}

func (h *Handler) updateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.TasksDTO
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("bind json err: %s \n", err.Error())
		ctx.JSON(http.StatusBadRequest, &dto.Error{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	err = h.srvs.UpdateTask(ctx, &req, id)
	if err != nil {
		log.Printf("can not update task: %s \n", err.Error())
		if err.Error() == "task not found" || err.Error() == "a task with the same title and activeAt already exists"+
			"" || err.Error() == "the provided hex string is not a valid ObjectID" {
			ctx.JSON(http.StatusBadRequest, &dto.Error{
				Code:    -2,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, &dto.Error{
			Code:    -3,
			Message: "invalid update task",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *Handler) updateTaskStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.srvs.UpdateTaskStatus(ctx, id, "done")
	if err != nil {
		log.Printf("can not update status task: %s \n", err.Error())
		if err == mongo.ErrNoDocuments || err.Error() == "the provided hex string is not a valid ObjectID" {
			ctx.JSON(http.StatusBadRequest, &dto.Error{
				Code:    -1,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, &dto.Error{
			Code:    -2,
			Message: "invalid update status task",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *Handler) deleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.srvs.DeleteTask(ctx, id)
	if err != nil {
		log.Printf("can not delete task: %s \n", err.Error())
		if err.Error() == "task not found" || err.Error() == "the provided hex string is not a valid ObjectID" {
			ctx.JSON(http.StatusBadRequest, &dto.Error{
				Code:    -1,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, &dto.Error{
			Code:    -2,
			Message: "invalid delete task",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *Handler) getAllTasks(ctx *gin.Context) {
	status := ctx.Query("status")

	tasks, err := h.srvs.GetAllTasks(ctx, status)
	if err != nil {
		log.Printf("can not get task: %s \n", err.Error())
		ctx.JSON(http.StatusInternalServerError, &dto.Error{
			Code:    -1,
			Message: "invalid to get task",
		})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
