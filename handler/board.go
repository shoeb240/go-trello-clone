package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"github.com/shoeb240/go-trello-clone/service"
)

type BoardHandler struct {
	service *service.BoardService
}

func NewBoardHandler(repository *repository.Repository) *BoardHandler {
	return &BoardHandler{
		service: service.NewBoardService(repository),
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

	boardID, err := h.service.Create(c.Request.Context(), boardModel)
	boardModel.ID = boardID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created", "data": boardModel})
}
