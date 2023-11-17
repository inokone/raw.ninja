package auth

import (
	"fmt"
	"math"
	"net/http"
	"time"

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
	captcha common.RecaptchaValidator
}

// NewController creates a new `Controller`, based on the user persistence.
func NewController(users user.Storer, auths account.Storer, jwt JWTHandler, c common.AuthConfig) Controller {
	p := bluemonday.StrictPolicy()
	return Controller{
		users:   users,
		auths:   auths,
		jwt:     jwt,
		p:       *p,
		captcha: common.NewRecaptchaValidator(c.RecaptchaSecret),
	}
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
// @Router /account/login [post]
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
		s   account.Account
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
		s       account.Account
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
		s   account.Account
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
// @Router /account/logout [get]
func (c Controller) Logout(g *gin.Context) {
	g.SetCookie(jwtTokenKey, "", 0, "", "", true, true)
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Logged out successfully! See you!",
	})
}
