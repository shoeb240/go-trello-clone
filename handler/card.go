package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"github.com/shoeb240/go-trello-clone/service"
)

type CardHandler struct {
	service  *service.CardService
	validate *validator.Validate
}

func NewCardHandler(repository *repository.Repository) *CardHandler {
	return &CardHandler{
		service:  service.NewCardService(repository),
		validate: validator.New(),
	}
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	cardModel := model.Card{}
	c.BindJSON(&cardModel)

	err := h.validate.Struct(cardModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cardModelNew, err := h.service.Create(c.Request.Context(), cardModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //"Internal server error"
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created", "data": cardModelNew})
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	cardModel := model.Card{}
	c.BindJSON(&cardModel)

	ctx := context.WithValue(c.Request.Context(), boardIDKey, c.Param("boardID"))
	ctx = context.WithValue(ctx, listIDKey, c.Param("listID"))

	msg, err := h.service.Update(ctx, cardModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
