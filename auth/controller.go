package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/common"
)

const (
	jwtTokenKey string = "Authorization"
)

type Controller struct {
	store  Store
	config common.AuthConfig
}

func NewController(db *gorm.DB, config *common.AuthConfig) Controller {
	return Controller{
		store:  Store{db: db},
		config: *config,
	}
}

// @BasePath /api/v1/auth

// Register godoc
// @Summary User registration endpoint
// @Schemes
// @Description Registers the user
// @Accept json
// @Produce json
// @Success 201 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /signup [post]
func (c Controller) Signup(g *gin.Context) {
	var s Registration
	err := g.Bind(&s)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Incorrect user registration data provided!",
		})
		return
	}
	user, err := NewUser(s.Email, s.Password, s.Phone)
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Could not create user.",
		})
		return
	}
	err = c.store.Store(*user)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "User with this email already exist.",
		})
		return
	}
	g.JSON(http.StatusCreated, common.StatusMessage{
		Code:    201,
		Message: "User has been created!",
	})
}

// Login godoc
// @Summary User login endpoint
// @Schemes
// @Description Logs in the user, sets up the JWT authorization
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /login [post]
func (c Controller) Login(g *gin.Context) {
	var s Credentials
	err := g.Bind(&s)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Incorrect user registration data provided!",
		})
		return
	}

	user, err := c.store.ByEmail(s.Email)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "User does not exist or password does not match!",
		})
		return
	}

	verified := user.VerifyPassword(s.Password)
	if !verified {
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "User does not exist or password does not match!",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(c.config.JWTSecret))

	g.SetSameSite(http.SameSiteLaxMode)
	g.SetCookie(jwtTokenKey, tokenString, 3600*24*30, "", "", false, true) // Max live time is 30 days
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Could not sign JWT token, please contact administrator!",
		})
		return
	}

	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Logged in!",
	})
}

func (c Controller) Validate(g *gin.Context) {
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "I am authorized!",
	})
}

// Reset godoc - not implemented yet
// @Summary Reset password endpoint
// @Schemes
// @Description Returns the status and version of the application
// @Accept json
// @Produce json
// @Failure 501 {object} common.StatusMessage
// @Router /reset [post]
func (c Controller) Reset(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, common.StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}

// Logout godoc
// @Summary Logout endpoint
// @Schemes
// @Description Logs out of the application, deletes the JWT token uased for authorization
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Router /logout [get]
func (c Controller) Logout(g *gin.Context) {
	g.SetCookie(jwtTokenKey, "", 0, "", "", true, true)
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Logged out successfully! See you!",
	})
}
