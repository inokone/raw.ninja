package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/mail"
)

const (
	jwtTokenKey string        = "Authorization"
	twoWeeks    time.Duration = time.Hour * 24 * 14
)

var (
	statusInvalidCredentials common.StatusMessage = common.StatusMessage{Code: 404, Message: "User does not exist or password does not match!"}
	statusBadRequest         common.StatusMessage = common.StatusMessage{Code: 400, Message: "Incorrect user data provided!"}
)

// Controller is a struct for web handles related to authentication and authorization.
type Controller struct {
	users  user.Storer
	auths  Storer
	jwt    JWTHandler
	sender mail.Service
	config common.AuthConfig
}

// NewController creates a new `Controller`, based on the user persistence and the authentication configuration parameters.
func NewController(users user.Storer, auths Storer, jwt JWTHandler, sender mail.Service, config common.AuthConfig) Controller {
	return Controller{
		users:  users,
		auths:  auths,
		jwt:    jwt,
		sender: sender,
		config: config,
	}
}

// Signup is a method of `Controller`. Signs the user up for the application with username/password credentials.
// @Summary User registration endpoint
// @Schemes
// @Description Signs the user up for the application
// @Accept json
// @Produce json
// @Success 201 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /auth/signup [post]
func (c Controller) Signup(g *gin.Context) {
	var s user.Registration
	if err := g.Bind(&s); err != nil {
		log.Err(err).Msg("invalid request")
		g.JSON(http.StatusBadRequest, statusBadRequest)
		return
	}
	usr, err := user.NewUser(s.Email, s.Password, s.Phone)
	if err != nil {
		log.Err(err).Msg("can not create new user")
		g.JSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Could not create user.",
		})
		return
	}
	if err = c.users.Store(*usr); err != nil {
		log.Err(err).Msg("can not store user")
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "User with this email already exist.",
		})
		return
	}

	err = c.sendConfirmation(usr)
	if err != nil {
		log.Err(err).Msg("can not send e-mail confirmation")
		g.JSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Could not create user.",
		})
		return
	}

	g.JSON(http.StatusCreated, common.StatusMessage{
		Code:    201,
		Message: "User has been created!",
	})
}

func (c Controller) sendConfirmation(usr *user.User) error {
	state := AuthenticationState{
		UserID:                usr.ID,
		EmailConfirmationHash: uuid.New().String(),
		EmailConfirmationTTL:  time.Now().Add(twoWeeks),
	}
	if err := c.auths.Store(state); err != nil {
		return err
	}
	url := c.config.DomainRoot + "/auth/confirmation?" + state.EmailConfirmationHash
	c.sender.EmailConfirmation(usr.Email, url)
	return nil
}

// Login is a method of `Controller`. Authenticates the user to the application, sets a JWT token on success in the cookies.
// @Summary User login endpoint
// @Schemes
// @Description Logs in the user, sets up the JWT authorization
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /auth/login [post]
func (c Controller) Login(g *gin.Context) {
	var s user.Credentials
	err := g.Bind(&s)
	if err != nil {
		g.JSON(http.StatusBadRequest, statusBadRequest)
		return
	}

	user, err := c.users.ByEmail(s.Email)
	if err != nil {
		g.JSON(http.StatusBadRequest, statusInvalidCredentials)
		return
	}

	verified := user.VerifyPassword(s.Password)
	if !verified {
		g.JSON(http.StatusBadRequest, statusInvalidCredentials)
		return
	}

	c.jwt.Issue(g, user.ID.String())

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
// @Router /auth/logout [get]
func (c Controller) Logout(g *gin.Context) {
	g.SetCookie(jwtTokenKey, "", 0, "", "", true, true)
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Logged out successfully! See you!",
	})
}
