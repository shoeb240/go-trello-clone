package app

func (appConfig *appConfig) InitRouter() {
	handler := appConfig.handler
	appConfig.router.GET("/:boardID", handler.Board.GetBoard)
	appConfig.router.POST("/board", handler.Board.CreateBoard)
	appConfig.router.PATCH("/board/:boardID", handler.Board.UpdateBoard)

	appConfig.router.POST("/list/:boardID", handler.List.CreateList)
	appConfig.router.PATCH("/list/:boardID/:listID", handler.List.UpdateList)

	appConfig.router.POST("/card", handler.Card.CreateCard)
	appConfig.router.PATCH("/card/:cardID", handler.Card.UpdateCard)
	appConfig.router.PATCH("/card/:cardID/move", handler.Card.MoveCard)
}
