package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/inokone/photostorage/common"
)

const (
	jwtTokenKey string = "Authorization"
)

var (
	// StatusInvalidCredentials is a `common.StatusMessage` is a response for non-existing users and invalid credentials.
	StatusInvalidCredentials common.StatusMessage = common.StatusMessage{Code: 404, Message: "User does not exist or password does not match!"}
	// StatusBadRequest is a `common.StatusMessage` is a response for invalid requests.
	StatusBadRequest common.StatusMessage = common.StatusMessage{Code: 400, Message: "Incorrect user data provided!"}
)

// Controller is a struct for web handles related to authentication and authorization.
type Controller struct {
	users Storer
	jwt   JWTHandler
}

// NewController creates a new `Controller`, based on the user persistence and the authentication configuration parameters.
func NewController(users Storer, jwt JWTHandler) Controller {
	return Controller{
		users: users,
		jwt:   jwt,
	}
}

// @BasePath /api/v1/auth

// Signup is a method of `Controller`. Signs the user up for the application with username/password credentials.
// @Summary User registration endpoint
// @Schemes
// @Description Signs the user up for the application
// @Accept json
// @Produce json
// @Success 201 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /signup [post]
func (c Controller) Signup(g *gin.Context) {
	var s Registration
	if err := g.Bind(&s); err != nil {
		g.JSON(http.StatusBadRequest, StatusBadRequest)
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
	if err = c.users.Store(*user); err != nil {
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

// Login is a method of `Controller`. Authenticates the user to the application, sets a JWT token on success in the cookies.
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
		g.JSON(http.StatusBadRequest, StatusBadRequest)
		return
	}

	user, err := c.users.ByEmail(s.Email)
	if err != nil {
		g.JSON(http.StatusBadRequest, StatusInvalidCredentials)
		return
	}

	verified := user.VerifyPassword(s.Password)
	if !verified {
		g.JSON(http.StatusBadRequest, StatusInvalidCredentials)
		return
	}

	c.jwt.Issue(g, user.ID.String())

	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Logged in!",
	})
}

// Profile is a method of `Controller`. Retrieves profile data of the user based on the JWT token in the request.
// @Summary Get user profile endpoint
// @Schemes
// @Description Gets the current logged in user
// @Accept json
// @Produce json
// @Success 200 {object} Profile
// @Failure 403 {object} common.StatusMessage
// @Router /profile [get]
func (c Controller) Profile(g *gin.Context) {
	userObj, _ := g.Get("user")
	user := userObj.(User)
	g.JSON(http.StatusOK, user.AsProfile())
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

// Logout is a method of `Controller`. Clears the JWT token from the cookies thus logging out the current user.
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
