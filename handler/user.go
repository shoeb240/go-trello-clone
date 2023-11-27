package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
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

func (h *UserHandler) Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userModel, err := h.service.Login(c.Request.Context(), body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = bcrypt.CompareHashAndPassword(userModel.Password, []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, "invalid email or password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userModel.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	godotenv.Load(".env")
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	c.JSON(http.StatusCreated, gin.H{"token": tokenString})
}
