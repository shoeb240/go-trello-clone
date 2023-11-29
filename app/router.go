package app

import (
	"github.com/shoeb240/go-trello-clone/middleware"
)

func (appConfig *appConfig) InitRouter() {
	handler := appConfig.handler

	appConfig.router.POST("/signup", handler.User.Signup)
	appConfig.router.POST("/login", handler.User.Login)

	authGroup := appConfig.router.Group("")
	authGroup.Use(middleware.Auth(appConfig.Repository))

	{
		authGroup.GET("/board/:boardID", handler.Board.GetBoard)
		authGroup.POST("/board", handler.Board.CreateBoard)
		authGroup.PATCH("/board/:boardID", handler.Board.UpdateBoard)
		authGroup.DELETE("/board/:boardID", handler.Board.DeleteBoard)

		authGroup.POST("/list/:boardID", handler.List.CreateList)
		authGroup.PATCH("/list/:boardID/:listID", handler.List.UpdateList)
		authGroup.DELETE("/list/:boardID/:listID", handler.List.DeleteList)

		authGroup.POST("/card", handler.Card.CreateCard)
		authGroup.PATCH("/card/:cardID", handler.Card.UpdateCard)
		authGroup.PATCH("/card/:cardID/move", handler.Card.MoveCard)
		authGroup.DELETE("/card/:cardID", handler.Card.DeleteCard)
	}
}
