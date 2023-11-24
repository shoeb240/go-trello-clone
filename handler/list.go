package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"github.com/shoeb240/go-trello-clone/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListHandler struct {
	service  *service.ListService
	validate *validator.Validate
}

type CustomKeyType string

const (
	boardIDKey repository.CustomKeyType = "boardID"
	listIDKey  repository.CustomKeyType = "listID"
)

func NewListHandler(repository *repository.Repository) *ListHandler {
	return &ListHandler{
		service:  service.NewListService(repository),
		validate: validator.New(),
	}
}

func (h *ListHandler) CreateList(c *gin.Context) {
	listModel := model.List{}
	c.BindJSON(&listModel)
	listModel.ID = primitive.NewObjectID()

	err := h.validate.Struct(listModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.WithValue(c.Request.Context(), boardIDKey, c.Param("boardID"))

	err = h.service.Create(ctx, listModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created", "data": listModel})
}

func (h *ListHandler) UpdateList(c *gin.Context) {
	listModel := model.List{}
	c.BindJSON(&listModel)

	ctx := context.WithValue(c.Request.Context(), boardIDKey, c.Param("boardID"))
	ctx = context.WithValue(ctx, listIDKey, c.Param("listID"))

	msg, err := h.service.Update(ctx, listModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *ListHandler) DeleteList(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), boardIDKey, c.Param("boardID"))
	ctx = context.WithValue(ctx, listIDKey, c.Param("listID"))

	msg, err := h.service.Delete(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
