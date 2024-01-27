package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/rs/zerolog/log"
)

// FacebookSource is the source we use for `user.User`s registered from Google OAuth
const FacebookSource = "Facebook"

// state should be regenerated per auth request
var (
	FacebookState = "facebook_random_csrf_string"
)

// FacebookController is a handler for endpoints of Facebook OAuth based autentication
type FacebookController struct {
	key         string
	secret      string
	users       user.Storer
	jwt         JWTHandler
	successURL  string
	redirectURL string
}

// NewFacebookController is a function creating an instance of `FacebookController`
func NewFacebookController(c common.AuthConfig, users user.Storer, jwt JWTHandler) *FacebookController {
	return &FacebookController{
		key:         c.FacebookKey,
		secret:      c.FacebookSecret,
		users:       users,
		jwt:         jwt,
		successURL:  c.FrontendRoot + "/home",
		redirectURL: c.BackendRoot + "/api/public/v1/auth/facebook/redirect",
	}
}

// Login endpoint
// @Summary Login is the authentication endpoint. Starts Facebook authentication process.
// @Schemes
// @Description Starts Facebook authentication process.
// @Accept json
// @Produce json
// @Router /auth/facebook [get]
func (c FacebookController) Login(g *gin.Context) {
	url := c.getConfig().AuthCodeURL(FacebookState)
	g.Redirect(http.StatusTemporaryRedirect, url)
}

// Redirect endpoint
// @Summary Redirect is the authentication callback endpoint. Authenticates/Registers users, sets up JWT token.
// @Schemes
// @Description Called by Facebook Auth when we have a result of the authentication process
// @Accept json
// @Produce text/html
// @Router /auth/facebook/redirect [get]
func (c FacebookController) Redirect(g *gin.Context) {
	email, err := c.authenticateCode(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: err.Error()})
		return
	}

	usr, err := c.users.ByEmail(email)
	if err != nil {
		log.Debug().Msg("Populating new user from Facebook.")
		userData := &user.User{
			Email:    email,
			PassHash: "",
			Source:   FacebookSource,
			RoleID:   2,
			Status:   user.Confirmed,
			Enabled:  true,
		}
		if err = c.users.Store(userData); err != nil {
			log.Err(err).Msg("Can not store Facebook user.")
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
	if usr.Source != FacebookSource {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "The provided email address is registered already with a different provider!"})
		return
	}
	c.jwt.Issue(g, usr.ID.String())
	g.Redirect(http.StatusTemporaryRedirect, c.successURL)
}

func (c FacebookController) authenticateCode(g *gin.Context) (string, error) {
	var u UserDetails
	state := g.Query("state")
	if state != FacebookState {
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
	response, err := client.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
	if err != nil {
		log.Err(err).Msg("error getting userinfo")
		return "", fmt.Errorf("error getting userinfo: %s", err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&u)
	if err != nil {
		log.Err(err).Msg("user details invalid")
		return "", fmt.Errorf("user details invalid: %s", err)
	}

	return u.Email, nil
}

func (c FacebookController) getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.key,
		ClientSecret: c.secret,
		RedirectURL:  c.redirectURL,
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}
}

// UserDetails is the user details structure for Facebook userdetails API
type UserDetails struct {
	ID      string
	Name    string
	Email   string
	Picture string
}
