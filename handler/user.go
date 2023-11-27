package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"github.com/shoeb240/go-trello-clone/service"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service  *service.UserService
	validate *validator.Validate
}

func NewUserHandler(repository *repository.Repository) *UserHandler {
	return &UserHandler{
		service:  service.NewUserService(repository),
		validate: validator.New(),
	}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userModel := model.User{
		Email:    body.Email,
		Password: hash,
	}
	err = h.validate.Struct(userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.service.Signup(c.Request.Context(), userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userModel.ID = userID

	c.JSON(http.StatusCreated, gin.H{"message": "Created", "data": userModel})
}
