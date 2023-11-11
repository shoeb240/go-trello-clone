package handler

import "github.com/shoeb240/go-trello-clone/repository"

type Handler struct {
	Board *BoardHandler
	List  *ListHandler
}

func NewHandlers(repository *repository.Repository) *Handler {
	return &Handler{
		Board: NewBoardHandler(repository),
		List:  NewListHandler(repository),
	}
}
