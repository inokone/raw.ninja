package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v2"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/rs/zerolog/log"
)

// GoogleSource is the source we use for `user.User`s registered from Google OAuth
const GoogleSource = "Google"

// state should be regenerated per auth request
var (
	GoogleState = "google_random_csrf_string"
)

// GoogleController is a handler for endpoints of Google OAuth based autentication
type GoogleController struct {
	key         string
	secret      string
	users       user.Storer
	jwt         JWTHandler
	successURL  string
	redirectURL string
}

// NewGoogleController is a function creating an instance of `GoogleController`
func NewGoogleController(c common.AuthConfig, users user.Storer, jwt JWTHandler) *GoogleController {
	return &GoogleController{
		key:         c.GoogleKey,
		secret:      c.GoogleSecret,
		users:       users,
		jwt:         jwt,
		successURL:  c.FrontendRoot + "/home",
		redirectURL: c.BackendRoot + "/api/public/v1/auth/google/redirect",
	}
}

// Login endpoint
// @Summary Login is the authentication endpoint. Starts Google authentication process.
// @Schemes
// @Description Starts Google authentication process.
// @Accept json
// @Produce json
// @Router /auth/google [get]
func (c GoogleController) Login(g *gin.Context) {
	url := c.getConfig().AuthCodeURL(GoogleState)
	g.Redirect(http.StatusTemporaryRedirect, url)
}

// Redirect endpoint
// @Summary Redirect is the authentication callback endpoint. Authenticates/Registers users, sets up JWT token.
// @Schemes
// @Description Called by Google Auth when we have a result of the authentication process
// @Accept json
// @Produce text/html
// @Router /auth/google/redirect [get]
func (c GoogleController) Redirect(g *gin.Context) {
	email, err := c.authenticateCode(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: err.Error()})
		return
	}

	usr, err := c.users.ByEmail(email)
	if err != nil {
		log.Debug().Msg("Populating new user from Google.")
		userData := &user.User{
			// Name:  google_user.Name,
			// Photo: google_user.Picture,
			Email:    email,
			PassHash: "",
			Source:   GoogleSource,
			RoleID:   2,
			Status:   user.Confirmed,
			Enabled:  true,
		}
		if err = c.users.Store(userData); err != nil {
			log.Err(err).Msg("Can not store Google user.")
			g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Something went wrong. Please contact our administrators!"})
			return
		}
		usr, err = c.users.ByEmail(email)
		if err != nil {
			log.Err(err).Msg("Can not load Google user.")
			g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Something went wrong. Please contact our administrators!"})
			return
		}
	}
	if !usr.Enabled {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Your account has been deactivated. Please contact our administrators!"})
		return
	}
	if usr.Source != GoogleSource {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "The provided email address is registered already with a different provider!"})
		return
	}
	c.jwt.Issue(g, usr.ID.String())
	g.Redirect(http.StatusTemporaryRedirect, c.successURL)
}

func (c GoogleController) authenticateCode(g *gin.Context) (string, error) {
	state := g.Query("state")
	if state != GoogleState {
		log.Warn().Msg("invalid oauth state")
		return "", errors.New("invalid oauth state")
	}

	code := g.Query("code")
	token, err := c.getConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Err(err).Msg("token exchange error")
		return "", fmt.Errorf("token exchange error: %s", err)
	}

	client := c.getConfig().Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Err(err).Msg("error getting userinfo")
		return "", fmt.Errorf("error getting userinfo: %s", err)
	}

	//nolint:staticcheck
	userinfoService, err := goauth2.New(client)
	if err != nil {
		log.Err(err).Msg("error creating userinfo")
		return "", fmt.Errorf("error creating userinfo service: %s", err)
	}

	userinfo, err := goauth2.NewUserinfoV2MeService(userinfoService).Get().Context(g).Do()
	if err != nil {
		log.Err(err).Msg("error getting userinfo")
		return "", fmt.Errorf("error getting userinfo: %s", err)
	}

	defer response.Body.Close()

	return userinfo.Email, nil
}

func (c GoogleController) getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.key,
		ClientSecret: c.secret,
		RedirectURL:  c.redirectURL,
		Scopes:       []string{goauth2.UserinfoEmailScope}, // goauth2.UserinfoProfileScope},
		Endpoint:     google.Endpoint,
	}
}
