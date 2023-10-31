package auth

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
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
	statusBadRequest         common.StatusMessage = common.StatusMessage{Code: 400, Message: "Invalid user data provided!"}
)

// Controller is a struct for web handles related to authentication and authorization.
type Controller struct {
	users  user.Storer
	auths  Storer
	jwt    JWTHandler
	sender mail.Service
	config common.AuthConfig
	p      bluemonday.Policy
}

// NewController creates a new `Controller`, based on the user persistence and the authentication configuration parameters.
func NewController(users user.Storer, auths Storer, jwt JWTHandler, sender mail.Service, config common.AuthConfig) Controller {
	p := bluemonday.StrictPolicy()
	return Controller{
		users:  users,
		auths:  auths,
		jwt:    jwt,
		sender: sender,
		config: config,
		p:      *p,
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
	state := AuthenticationState{
		UserID:                usr.ID,
		EmailConfirmationHash: uuid.New().String(),
		EmailConfirmationTTL:  time.Now().Add(twoWeeks),
	}
	if err := c.auths.Store(&state); err != nil {
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
// @Router /auth/resend [post]
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
	s, err := c.auths.ByUser(usr.ID)
	if err != nil {
		return err
	}
	s.EmailConfirmationHash = uuid.New().String()
	s.EmailConfirmationTTL = time.Now().Add(twoWeeks)
	if err := c.auths.Update(&s); err != nil {
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
// @Router /auth/confirm [get]
func (c Controller) Confirm(g *gin.Context) {
	var (
		token string
		state AuthenticationState
		err   error
		usr   *user.User
	)
	token = g.Query("token")
	state, err = c.auths.ByConfirmToken(token)
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
	if err = c.auths.Update(&state); err != nil {
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

// Login is a method of `Controller`. Authenticates the user to the application, sets a JWT token on success in the cookies.
// @Summary User login endpoint
// @Schemes
// @Description Logs in the user, sets up the JWT authorization
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 403 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /auth/login [post]
func (c Controller) Login(g *gin.Context) {
	var (
		s        user.Credentials
		err      error
		usr      *user.User
		verified bool
		secs     int64
	)

	err = g.ShouldBindJSON(&s)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, statusBadRequest)
		return
	}
	s.Email = c.p.Sanitize(s.Email)

	usr, err = c.users.ByEmail(s.Email)
	if err != nil {
		log.Err(err).Msg("Failed to collect user.")
		g.AbortWithStatusJSON(http.StatusBadRequest, statusInvalidCredentials)
		return
	}

	secs, err = c.checkTimeout(usr)
	if err != nil {
		log.Err(err).Str("UserID", usr.ID.String()).Msg("Failed to collect login timeout.")
		g.AbortWithStatusJSON(http.StatusBadRequest, statusInvalidCredentials)
		return
	}
	if secs > 0 {
		g.AbortWithStatusJSON(http.StatusForbidden, common.StatusMessage{
			Code:    403,
			Message: fmt.Sprintf("You have been locked out for failed credentials. You have to wait %v more seconds.", secs),
		})
		return
	}

	verified = usr.VerifyPassword(s.Password)
	if !verified {
		err = c.increaseTimeout(usr)
		if err != nil {
			log.Warn().Str("user", usr.ID.String()).Msg("Failed to increase timeout for user")
		}
		g.AbortWithStatusJSON(http.StatusBadRequest, statusInvalidCredentials)
		return
	}

	c.clearTimeout(usr)
	c.jwt.Issue(g, usr.ID.String())

	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Logged in!",
	})
}

func (c Controller) checkTimeout(usr *user.User) (int64, error) {
	var (
		s   AuthenticationState
		err error
	)
	s, err = c.auths.ByUser(usr.ID)
	if err != nil {
		return 0, err
	}
	if s.FailedLoginLock.After(time.Now()) {
		return s.FailedLoginLock.Unix() - time.Now().Unix(), nil
	}
	return 0, nil
}

func (c Controller) increaseTimeout(usr *user.User) error {
	var (
		s       AuthenticationState
		err     error
		timeout int
	)
	s, err = c.auths.ByUser(usr.ID)
	if err != nil {
		return err
	}
	s.FailedLoginCounter++
	s.LastFailedLogin = time.Now()
	if s.FailedLoginCounter > 2 {
		timeout = int(math.Pow(10, float64(s.FailedLoginCounter-2))) // exponential backoff - 10 sec, 10 sec, 1000 sec, ...
		s.FailedLoginLock = time.Now().Add(time.Second * time.Duration(timeout))
	}
	return c.auths.Update(&s)
}

func (c Controller) clearTimeout(usr *user.User) error {
	var (
		s   AuthenticationState
		err error
	)
	s, err = c.auths.ByUser(usr.ID)
	if err != nil {
		return err
	}
	s.FailedLoginCounter = 0
	s.FailedLoginLock = time.Now()
	return c.auths.Update(&s)
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
