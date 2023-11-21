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

type CardHandler struct {
	service  *service.CardService
	validate *validator.Validate
}

const (
	cardIDKey repository.CustomKeyType = "cardID"
)

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

	ctx := context.WithValue(c.Request.Context(), cardIDKey, c.Param("cardID"))

	msg, err := h.service.Update(ctx, cardModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *CardHandler) MoveCard(c *gin.Context) {
	cardMoveReq := repository.CardMoveReq{}
	if err := c.BindJSON(&cardMoveReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.validate.Struct(cardMoveReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cardID, err := primitive.ObjectIDFromHex(c.Param("cardID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	cardMoveReq.CardID = cardID

	ctx := context.WithValue(c.Request.Context(), cardIDKey, c.Param("cardID"))

	msg, err := h.service.MoveCard(ctx, cardMoveReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	cardObjID, err := primitive.ObjectIDFromHex(c.Param("cardID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	msg, err := h.service.DeleteCard(c.Request.Context(), cardObjID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
