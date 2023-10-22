package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/inokone/photostorage/common"
	"github.com/rs/zerolog/log"
)

// JWTHandler is a struct for issuing and validating JWT tokens.
type JWTHandler struct {
	conf  common.AuthConfig
	users Storer
}

// NewJWTHandler creates a new `JWTHandler`.
func NewJWTHandler(users Storer, conf common.AuthConfig) JWTHandler {
	return JWTHandler{
		conf:  conf,
		users: users,
	}
}

// Validate is a method of `JWTHandler`. Validates the authentication token in the Gin context provided as a parameter.
func (h *JWTHandler) Validate(g *gin.Context) {
	log.Debug().Msg("Validating JWT token...")
	tokenString, err := g.Cookie(jwtTokenKey)
	if err != nil {
		g.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, g.AbortWithError(http.StatusBadRequest, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(h.conf.JWTSecret), nil
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

// Issue is a method of `JWTHandler`. Issues a new authentication token for a user ID into the Gin context provided as parameters.
// The JWT token is set as a http-only cookie. The JWTSecure option of the AuthConfig controle "secure" option for the cookie.
// For production deployment this must be set to true.
func (h *JWTHandler) Issue(g *gin.Context, userID string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * time.Duration(h.conf.JWTExp)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.conf.JWTSecret))
	if err != nil {
		log.Warn().Err(err).Str("User", userID).Msg("JWT token could not be signed!")
		g.JSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Failed to sign JWT token, please contact administrator!",
		})
		return
	}

	g.SetSameSite(http.SameSiteLaxMode)
	g.SetCookie(jwtTokenKey, tokenString, 3600*24*30, "", "", h.conf.JWTSecure, true) // Max live time is 30 days
}
