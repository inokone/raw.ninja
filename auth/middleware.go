package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JWTHandler struct {
	secret string
	users  Store
}

func NewJWTHandler(db *gorm.DB, secret string) JWTHandler {
	return JWTHandler{
		secret: secret,
		users: Store{
			db: db,
		},
	}
}

func (h JWTHandler) ValidateJWTToken(g *gin.Context) {
	log.Print("Validating JWT token...")
	tokenString, err := g.Cookie(jwtTokenKey)
	if err != nil {
		g.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, g.AbortWithError(http.StatusBadRequest, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(h.secret), nil
	})
	if err != nil {
		g.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		g.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	expired := float64(time.Now().Unix()) > claims["exp"].(float64)
	if expired {
		g.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID := claims["sub"]
	uuid, err := uuid.Parse(userID.(string))
	if err != nil {
		g.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := h.users.ByID(uuid)
	if err != nil || user.Email == "" {
		g.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	g.Set("user", user)
	g.Next()
}
