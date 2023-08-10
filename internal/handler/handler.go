package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/khussa1n/todo-list/internal/custom_error"
	"github.com/khussa1n/todo-list/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	srvs service.Service
}

func New(srvs service.Service) *Handler {
	return &Handler{
		srvs: srvs,
	}
}

func parseIdFromPath(c *gin.Context, param string) (primitive.ObjectID, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return primitive.ObjectID{}, custom_error.ErrEmptyID
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return primitive.ObjectID{}, custom_error.ErrInvalidIDParameter
	}

	return id, nil
}
