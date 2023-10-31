package account

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
	twoWeeks time.Duration = time.Hour * 24 * 14
)

var statusBadRequest common.StatusMessage = common.StatusMessage{Code: 400, Message: "Invalid user data provided!"}

// Controller is a struct for web handles related to authentication and authorization.
type Controller struct {
	users    user.Storer
	accounts Storer
	sender   mail.Service
	config   common.AuthConfig
}

// NewController creates a new `Controller`, based on the user persistence and the authentication configuration parameters.
func NewController(users user.Storer, accounts Storer, sender mail.Service, config common.AuthConfig) Controller {
	return Controller{
		users:    users,
		accounts: accounts,
		sender:   sender,
		config:   config,
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
// @Router /account/signup [post]
func (c Controller) Signup(g *gin.Context) {
	var s user.Registration
	if err := g.ShouldBindJSON(&s); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationMessage(err))
		return
	}

	usr, err := user.NewUser(s.Email, s.Password, s.FirstName, s.LastName)
	if err != nil {
		log.Err(err).Msg("Could not create new user")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Could not create user.",
		})
		return
	}

	if err = c.users.Store(usr); err != nil {
		log.Err(err).Msg("Could not store user")
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "User with this email already exist.",
		})
		return
	}

	err = c.sendMail(usr)
	if err != nil {
		log.Err(err).Msg("Could not send e-mail confirmation")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{
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

func (c Controller) sendMail(usr *user.User) error {
	state := Account{
		UserID:                usr.ID,
		EmailConfirmationHash: uuid.New().String(),
		EmailConfirmationTTL:  time.Now().Add(twoWeeks),
	}
	if err := c.accounts.Store(&state); err != nil {
		return err
	}
	url := c.config.FrontendRoot + "/confirm?token=" + state.EmailConfirmationHash
	c.sender.EmailConfirmation(usr.Email, url)
	return nil
}

// Resend is a method of `Controller`. Resends email confirmation for an email address.
// @Summary Resends email confirmation endpoint
// @Schemes
// @Description Resends email confirmation for an email address.
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /account/resend [put]
func (c Controller) Resend(g *gin.Context) {
	var (
		s   ConfirmationResend
		err error
		usr *user.User
	)

	if err = g.ShouldBindJSON(&s); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, "Invalid confirmation resend data.")
		return
	}

	usr, err = c.users.ByEmail(s.Email)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, "Invalid confirmation resend e-mail.")
		return
	}

	err = c.resendMail(usr)
	if err != nil {
		log.Err(err).Msg("Could not send e-mail confirmation")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Could not send confirmation email.",
		})
		return
	}

	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Confirmation sent!",
	})
}

func (c Controller) resendMail(usr *user.User) error {
	s, err := c.accounts.ByUser(usr.ID)
	if err != nil {
		return err
	}
	s.EmailConfirmationHash = uuid.New().String()
	s.EmailConfirmationTTL = time.Now().Add(twoWeeks)
	if err := c.accounts.Update(&s); err != nil {
		return err
	}
	url := c.config.FrontendRoot + "/confirm?token=" + s.EmailConfirmationHash
	c.sender.EmailConfirmation(usr.Email, url)
	return nil
}

// Confirm is a method of `Controller`. Confirms the email of the user for the hash provided as URL parameter.
// @Summary Email confirmation endpoint
// @Schemes
// @Description Confirms the email address of the user
// @Accept json
// @Produce json
// @Param   token    query     string  true  "Token for the email confirmation"  Format(uuid)
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /account/confirm [get]
func (c Controller) Confirm(g *gin.Context) {
	var (
		token string
		state Account
		err   error
		usr   *user.User
	)
	token = g.Query("token")
	state, err = c.accounts.ByConfirmToken(token)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid token!"})
		return
	}
	if state.EmailConfirmationTTL.Before(time.Now()) {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Expired token, resend it please!"})
		return
	}
	state.EmailConfirmationTTL = time.Now()
	state.EmailConfirmationHash = ""
	state.EmailConfirmed = true
	if err = c.accounts.Update(&state); err != nil {
		log.Err(err).Msg("Failed to update auth state.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	usr, err = c.users.ByID(state.UserID)
	if err != nil {
		log.Err(err).Msg("Failed to collect user.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	usr.Status = user.Confirmed
	if err = c.users.Update(usr); err != nil {
		log.Err(err).Msg("Failed to update user.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "E-mail is confirmed!",
	})
}

// RequestReset creates a reset request for a password of based on the email of a user - not implemented yet
// @Summary Reset password endpoint
// @Schemes
// @Description Resets the password of the logged in user
// @Accept json
// @Produce json
// @Failure 501 {object} common.StatusMessage
// @Router /account/reset [put]
func (c Controller) RequestReset(g *gin.Context) {
	g.AbortWithStatusJSON(http.StatusNotImplemented, common.StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}

// Reset resets the password of the logged in user - not implemented yet
// @Summary Reset password endpoint
// @Schemes
// @Description Resets the password of the logged in user
// @Accept json
// @Produce json
// @Failure 501 {object} common.StatusMessage
// @Router /account/reset [put]
func (c Controller) Reset(g *gin.Context) {
	g.AbortWithStatusJSON(http.StatusNotImplemented, common.StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}
