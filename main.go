package main

import (
	"github.com/shoeb240/go-trello-clone/app"
)

func main() {
	appConfig := app.NewAppConfig()
	appConfig.InitRouter()
	appConfig.StartServer()
}
