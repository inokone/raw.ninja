package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"

	"github.com/inokone/photostorage/auth/account"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
)

const (
	jwtTokenKey string = "Authorization"
)

var (
	statusInvalidCredentials common.StatusMessage = common.StatusMessage{Code: 404, Message: "User does not exist or password does not match!"}
	statusBadRequest         common.StatusMessage = common.StatusMessage{Code: 400, Message: "Invalid user data provided!"}
)

// Controller is a struct for web handles related to authentication and authorization.
type Controller struct {
	users   user.Storer
	auths   account.Storer
	jwt     JWTHandler
	p       bluemonday.Policy
	captcha *common.RecaptchaValidator
	service *Service
}

// NewController creates a new `Controller`, based on the user persistence.
func NewController(users user.Storer, auths account.Storer, jwt JWTHandler, c *common.AuthConfig) Controller {
	p := bluemonday.StrictPolicy()
	return Controller{
		users:   users,
		auths:   auths,
		jwt:     jwt,
		p:       *p,
		captcha: common.NewRecaptchaValidator(c.RecaptchaSecret),
		service: NewService(users, auths, jwt),
	}
}

// Login is a method of `Controller`. Authenticates the user to the application, sets a JWT token on success in the cookies.
// @Summary User login endpoint
// @Schemes
// @Description Logs in the user, sets up the JWT authorization
// @Accept json
// @Produce json
// @Param data body user.Credentials true "Credentials provided for the login"
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 403 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /account/login [post]
func (c Controller) Login(g *gin.Context) {
	var (
		s   user.Credentials
		err error
		usr *user.User
	)

	err = g.ShouldBindJSON(&s)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, statusBadRequest)
		return
	}

	isValid, err := c.captcha.Verify(s.Captcha)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationMessage(err))
		return
	}
	if !isValid {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Captcha verification failed!"})
		return
	}

	s.Email = c.p.Sanitize(s.Email)

	usr, err = c.users.ByEmail(s.Email)
	if err != nil {
		log.Err(err).Msg("Failed to collect user.")
		g.AbortWithStatusJSON(http.StatusBadRequest, statusInvalidCredentials)
		return
	}

	if !usr.Enabled {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Your account has been deactivated. Please contact our administrators!"})
		return
	}

	if usr.Source != "credentials" {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Your account can not be used with credentials!"})
		return
	}

	if err = c.service.ValidateCredentials(usr, s.Password); err != nil {
		switch err.(type) {
		default:
			g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error"})
		case *InvalidCredentials:
			g.AbortWithStatusJSON(http.StatusBadRequest, statusInvalidCredentials)
		case *LockedUser:
			g.AbortWithStatusJSON(http.StatusForbidden, common.StatusMessage{
				Code:    403,
				Message: fmt.Sprintf("You have been locked out for failed credentials. You have to wait %v more seconds.", err.(LockedUser).seconds),
			})
		}
		return
	}

	c.jwt.Issue(g, usr.ID.String())

	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Logged in!",
	})
}

// Logout is a method of `Controller`. Clears the JWT token from the cookies thus logging out the current user.
// @Summary Logout endpoint
// @Schemes
// @Description Logs out of the application, deletes the JWT token uased for authorization
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Router /account/logout [get]
func (c Controller) Logout(g *gin.Context) {
	g.SetCookie(jwtTokenKey, "", 0, "", "", true, true)
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Logged out successfully! See you!",
	})
}
