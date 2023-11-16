package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/inokone/photostorage/common"
)

// Controller is a struct for web handles related to application users.
type Controller struct {
	users Storer
}

// NewController creates a new `Controller`, based on the user persistence.
func NewController(users Storer) Controller {
	return Controller{
		users: users,
	}
}

// Profile is a method of `Controller`. Retrieves profile data of the user based on the JWT token in the request.
// @Summary Get user profile endpoint
// @Schemes
// @Description Gets the current logged in user
// @Accept json
// @Produce json
// @Success 200 {object} Profile
// @Failure 403 {object} common.StatusMessage
// @Router /users/profile [get]
func (c Controller) Profile(g *gin.Context) {
	u, _ := g.Get("user")
	usr := u.(*User)
	g.JSON(http.StatusOK, usr.AsProfile())
}

// Reset resets the password of the logged in user - not implemented yet
// @Summary Reset password endpoint
// @Schemes
// @Description Resets the password of the logged in user
// @Accept json
// @Produce json
// @Failure 501 {object} common.StatusMessage
// @Router /users/reset [post]
func (c Controller) Reset(g *gin.Context) {
	g.AbortWithStatusJSON(http.StatusNotImplemented, common.StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}

// List lists the users of the application.
// @Summary List users endpoint
// @Schemes
// @Description Lists the users of the application.
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Router /users [get]
func (c Controller) List(g *gin.Context) {
	var (
		users []User
		err   error
		res   []AdminView
	)
	users, err = c.users.List()
	if err != nil {
		log.Err(err).Msg("Failed to list users")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Unknown error, please contact administrator!",
		})
		return
	}

	res = make([]AdminView, 0)

	for _, usr := range users {
		res = append(res, usr.AsAdminView())
	}
	g.JSON(http.StatusOK, res)
}

// Patch updates settings (e.g. role) for a user.
// @Summary User update endpoint
// @Schemes
// @Description Updates the target user
// @Accept json
// @Produce json
// @Param id path int true "ID of the user information to patch"
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Router /users/:id [patch]
func (c Controller) Patch(g *gin.Context) {
	var (
		in  AdminView
		err error
	)
	if err = g.ShouldBindJSON(&in); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Malformed user data"})
		return
	}
	if err = c.users.Patch(in); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid user parameters provided!"})
		return
	}
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "User patched!",
	})
}

// Update details (firstname and lastname) for a user.
// @Summary User update endpoint
// @Schemes
// @Description Updates the target user
// @Accept json
// @Produce json
// @Param id path int true "ID of the user information to patch"
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Router /users/:id [put]
func (c Controller) Update(g *gin.Context) {
	var (
		in  Profile
		err error
	)

	u, _ := g.Get("user")
	usr := u.(*User)

	if err = g.ShouldBindJSON(&in); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Malformed user data"})
		return
	}

	usr.FirstName = in.FirstName
	usr.LastName = in.LastName

	if err = c.users.Update(usr); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid user parameters provided!"})
		return
	}
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "User patched!",
	})
}
