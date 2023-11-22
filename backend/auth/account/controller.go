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
	aDay     time.Duration = time.Hour * 24
	twoWeeks time.Duration = aDay * 14
)

var statusBadRequest common.StatusMessage = common.StatusMessage{Code: 400, Message: "Invalid user data provided!"}

// Controller is a struct for web handles related to authentication and authorization.
type Controller struct {
	users    user.Storer
	accounts Storer
	sender   mail.Service
	config   common.AuthConfig
	captcha  common.RecaptchaValidator
}

// NewController creates a new `Controller`, based on the user persistence and the authentication configuration parameters.
func NewController(users user.Storer, accounts Storer, sender mail.Service, config common.AuthConfig) Controller {
	return Controller{
		users:    users,
		accounts: accounts,
		sender:   sender,
		config:   config,
		captcha:  common.NewRecaptchaValidator(config.RecaptchaSecret),
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

	isValid, err := c.captcha.Verify(s.Captcha)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Capthca verification failed!"})
		return
	}
	if !isValid {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Captcha verification failed!"})
		return
	}

	usr, err := user.NewUser(s.Email, s.Password)
	if err != nil {
		log.Err(err).Msg("Could not create new user")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Could not create user.",
		})
		return
	}
	usr.RoleID = 1
	if err = c.users.Store(usr); err != nil {
		log.Err(err).Msg("Could not store user")
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "User with this email already exist.",
		})
		return
	}

	err = c.confirmMail(usr)
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

func (c Controller) confirmMail(usr *user.User) error {
	state := Account{
		UserID:            usr.ID,
		ConfirmationToken: uuid.New().String(),
		ConfirmationTTL:   time.Now().Add(twoWeeks),
	}
	if err := c.accounts.Store(&state); err != nil {
		return err
	}
	url := c.config.FrontendRoot + "/confirm?token=" + state.ConfirmationToken
	c.sender.EmailConfirmation(usr.Email, url)
	return nil
}

// ResendConfirmation is a method of `Controller`. Resends email confirmation for an email address.
// @Summary Resends email confirmation endpoint
// @Schemes
// @Description Resends email confirmation for an email address.
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /account/resend [put]
func (c Controller) ResendConfirmation(g *gin.Context) {
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
	s.ConfirmationToken = uuid.New().String()
	s.ConfirmationTTL = time.Now().Add(twoWeeks)
	if err := c.accounts.Update(&s); err != nil {
		return err
	}
	url := c.config.FrontendRoot + "/confirm?token=" + s.ConfirmationToken
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
		token   string
		account Account
		err     error
		usr     *user.User
	)
	token = g.Query("token")
	account, err = c.accounts.ByConfirmToken(token)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid token!"})
		return
	}
	if account.ConfirmationTTL.Before(time.Now()) {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Expired token, resend it please!"})
		return
	}
	account.ConfirmationTTL = time.Now()
	account.ConfirmationToken = ""
	account.Confirmed = true
	if err = c.accounts.Update(&account); err != nil {
		log.Err(err).Msg("Failed to update account.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	usr, err = c.users.ByID(account.UserID)
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

// Recover initiates a password reset - sends an email to a user
// @Summary Recover account endpoint
// @Schemes
// @Description Send a password reset email to a user
// @Accept json
// @Produce json
// @Success 202 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /account/recover [put]
func (c Controller) Recover(g *gin.Context) {
	var (
		s   Recovery
		err error
		usr *user.User
	)

	if err = g.ShouldBindJSON(&s); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Invalid recovery data.",
		})
		return
	}

	usr, err = c.users.ByEmail(s.Email)
	if err == nil {
		err = c.recoverMail(usr)
		if err != nil {
			log.Err(err).Msg("Could not send recovery e-mail")
		}
	}

	g.JSON(http.StatusAccepted, common.StatusMessage{
		Code:    202,
		Message: "Recover request accepted!",
	})
}

func (c Controller) recoverMail(usr *user.User) error {
	s, err := c.accounts.ByUser(usr.ID)
	if err != nil {
		return err
	}
	s.RecoveryToken = uuid.New().String()
	s.RecoveryTTL = time.Now().Add(aDay)
	if err := c.accounts.Update(&s); err != nil {
		return err
	}
	url := c.config.FrontendRoot + "/password/reset?token=" + s.RecoveryToken
	c.sender.PasswordReset(usr.Email, url)
	return nil
}

// ResetPassword resets the password of the logged in user - not implemented yet
// @Summary Reset password endpoint
// @Schemes
// @Description Resets the password of the logged in user
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /account/password/reset [put]
func (c Controller) ResetPassword(g *gin.Context) {
	var (
		reset PasswordReset
		state Account
		err   error
		usr   *user.User
	)

	if err = g.ShouldBindJSON(&reset); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Invalid recovery data.",
		})
		return
	}

	state, err = c.accounts.ByRecoveryToken(reset.Token)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid token!"})
		return
	}
	if state.RecoveryTTL.Before(time.Now()) {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Expired token, please restart the recovery!"})
		return
	}
	state.RecoveryToken = ""
	state.LastRecovery = time.Now()
	if err = c.accounts.Update(&state); err != nil {
		log.Err(err).Msg("Failed to update account.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	usr, err = c.users.ByID(state.UserID)
	if err != nil {
		log.Err(err).Msg("Failed to collect user.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	if err = usr.SetPassword(reset.Password); err != nil {
		log.Err(err).Msg("Failed to set password for the user.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	if err = c.users.Update(usr); err != nil {
		log.Err(err).Msg("Failed to update user.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Password updated!",
	})
}

// ChangePassword resets the password of the logged in user - not implemented yet
// @Summary Reset password endpoint
// @Schemes
// @Description Resets the password of the logged in user
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /account/password/change [put]
func (c Controller) ChangePassword(g *gin.Context) {
	var (
		chg PasswordChange
		err error
		usr *user.User
	)

	u, _ := g.Get("user")
	usr = u.(*user.User)

	if err = g.ShouldBindJSON(&chg); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Invalid change data.",
		})
		return
	}

	if !usr.VerifyPassword(chg.Old) {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Incorrect old password.",
		})
		return
	}

	if err = usr.SetPassword(chg.New); err != nil {
		log.Err(err).Msg("Failed to set password for the user.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	if err = c.users.Update(usr); err != nil {
		log.Err(err).Msg("Failed to update user.")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusBadRequest)
		return
	}

	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Password updated!",
	})
}
