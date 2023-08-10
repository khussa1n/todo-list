package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/khussa1n/todo-list/internal/custom_error"
	"github.com/khussa1n/todo-list/internal/entity/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

// createTask 	Create new task
// @Summary      Create task
// @Description  Create new task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param request body dto.TasksDTO true "req body"
// @Success      201  {object}  entity.Tasks
// @Failure      400  {object}  dto.Error
// @Failure      500  {object}  dto.Error
// @Router       /tasks [post]
func (h *Handler) createTask(ctx *gin.Context) {
	var req dto.TasksDTO
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("bind json err: %s \n", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.srvs.CreateTask(ctx, &req)
	if err != nil {
		log.Printf("can not create task: %s \n", err.Error())
		switch err {
		case custom_error.ErrMessageTooLong, custom_error.ErrInvalidActiveAtFormat, custom_error.ErrDuplicateTask:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusCreated, task)
}

// updateTask 	Update task
// @Summary      Update task
// @Description  Update task
// @Tags         task
// @Accept       json
// @Param req body dto.TasksDTO true "req body"
// @Success      204
// @Failure      400  {object}  dto.Error
// @Failure      500  {object}  dto.Error
// @Router       /tasks/{id} [put]
func (h *Handler) updateTask(ctx *gin.Context) {
	id, err := parseIdFromPath(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var req dto.TasksDTO
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("bind json err: %s \n", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.srvs.UpdateTask(ctx, &req, id)
	if err != nil {
		log.Printf("can not update task: %s \n", err.Error())
		switch err {
		case custom_error.ErrTaskNotFound, custom_error.ErrInvalidActiveAtFormat, custom_error.ErrDuplicateTask:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		case custom_error.ErrMessageTooLong:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusNoContent, "")
}

// updateTaskStatus 	Update task status
// @Summary      Update task status to done
// @Description  Update task status to done
// @Tags         task
// @Param 		 id   path      string  true  "Task ID"
// @Success      204
// @Failure      400  {object}  dto.Error
// @Failure      500  {object}  dto.Error
// @Router       /tasks/{id}/done [put]
func (h *Handler) updateTaskStatus(ctx *gin.Context) {
	id, err := parseIdFromPath(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.srvs.UpdateTaskStatus(ctx, id, "done")
	if err != nil {
		log.Printf("can not update status task: %s \n", err.Error())
		switch err {
		case mongo.ErrNoDocuments:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusNoContent, "")
}

// deleteTask 	Delete task
// @Summary      Delete task
// @Description  Delete task
// @Tags         task
// @Param 		 id   path      string  true  "Task ID"
// @Success      204
// @Failure      400  {object}  dto.Error
// @Failure      500  {object}  dto.Error
// @Router       /tasks/{id} [delete]
func (h *Handler) deleteTask(ctx *gin.Context) {
	id, err := parseIdFromPath(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.srvs.DeleteTask(ctx, id)
	if err != nil {
		log.Printf("can not delete task: %s \n", err.Error())
		switch err {
		case custom_error.ErrTaskNotFound:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusNoContent, "")
}

// getAllTasks 	Get all tasks
// @Summary      Get all tasks by status
// @Description  Get all tasks by status
// @Tags         task
// @Produce      json
// @Param		 status    query     string false "name search by status"
// @Success      200  {array}  entity.Tasks
// @Failure      400  {object}  dto.Error
// @Failure      500  {object}  dto.Error
// @Router       /tasks [get]
func (h *Handler) getAllTasks(ctx *gin.Context) {
	status := ctx.Query("status")

	tasks, err := h.srvs.GetAllTasks(ctx, status)
	if err != nil {
		log.Printf("can not get task: %s \n", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
