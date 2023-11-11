package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shoeb240/go-trello-clone/handler"
	"github.com/shoeb240/go-trello-clone/repository"
)

type appConfig struct {
	router     *gin.Engine
	repository *repository.Repository
	handler    *handler.Handler
}

func NewAppConfig() *appConfig {
	DB, err := repository.NewDBConnection()
	if err != nil {
		panic("Could not connect to database")
	}

	repository := repository.NewRepositories(DB)

	return &appConfig{
		router:     gin.Default(),
		repository: repository,
		handler:    handler.NewHandlers(repository),
	}
}

func (appConfig *appConfig) StartServer() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	appConfig.router.Run(":" + portString)
}
