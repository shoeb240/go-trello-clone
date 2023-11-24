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

type BoardHandler struct {
	service  *service.BoardService
	validate *validator.Validate
}

func NewBoardHandler(repository *repository.Repository) *BoardHandler {
	return &BoardHandler{
		service:  service.NewBoardService(repository),
		validate: validator.New(),
	}
}

func (h *BoardHandler) GetBoard(c *gin.Context) {
	board, err := h.service.FullBoard(c.Request.Context(), c.Param("boardID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	boardModel := model.Board{}
	c.BindJSON(&boardModel)

	err := h.validate.Struct(boardModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	boardID, err := h.service.Create(c.Request.Context(), boardModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	boardModel.ID = boardID

	c.JSON(http.StatusCreated, gin.H{"message": "Created", "data": boardModel})
}

func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	boardModel := model.Board{}
	if err := c.BindJSON(&boardModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.WithValue(c.Request.Context(), boardIDKey, c.Param("boardID"))

	msg, err := h.service.Update(ctx, boardModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	boardObjID, err := primitive.ObjectIDFromHex(c.Param("boardID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	msg, err := h.service.Delete(c.Request.Context(), boardObjID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
