package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/shoeb240/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Auth(repositories *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["uid"].(string)
			if claims["exp"].(float64) < float64(time.Now().Unix()) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			userObjID, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			user, err := repositories.User.FindByID(c, userObjID)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("user", user)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}
