package handler

import "github.com/shoeb240/go-trello-clone/repository"

type Handler struct {
	Board *BoardHandler
	List  *ListHandler
	Card  *CardHandler
}

func NewHandlers(repository *repository.Repository) *Handler {
	return &Handler{
		Board: NewBoardHandler(repository),
		List:  NewListHandler(repository),
		Card:  NewCardHandler(repository),
	}
}
