package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/inokone/photostorage/common"
)

// Controller is a struct for web handles related to roles.
type Controller struct {
	roles Storer
}

// NewController creates a new `Controller`, based on the user persistence parameter.
func NewController(roles Storer) Controller {
	return Controller{
		roles: roles,
	}
}

// List lists all roles of the application.
// @Summary Role list endpoint
// @Schemes
// @Description Lists all roles of the application
// @Accept json
// @Produce json
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Router /roles/ [get]
func (c Controller) List(g *gin.Context) {
	var (
		roles []Role
		err   error
		res   []ProfileRole
	)
	roles, err = c.roles.List()
	if err != nil {
		log.Err(err).Msg("Failed to list roles")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Unknown error, please contact administrator!",
		})
		return
	}

	res = make([]ProfileRole, 0)

	for _, role := range roles {
		res = append(res, role.AsProfileRole())
	}
	g.JSON(http.StatusOK, res)
}

// Update saves changes to settings (e.g. quota) for a user role.
// @Summary Role update endpoint
// @Schemes
// @Description Updates the settings of a role
// @Accept json
// @Produce json
// @Param id path int true "ID of the role information to patch"
// @Param data body role.ProfileRole true "The new version of the role to use for update"
// @Success 200 {object} common.StatusMessage
// @Failure 400 {object} common.StatusMessage
// @Router /roles/:id [put]
func (c Controller) Update(g *gin.Context) {
	var (
		in  ProfileRole
		err error
	)
	if err = g.ShouldBindJSON(&in); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Malformed role data"})
		return
	}
	if err = c.roles.Update(in); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid role parameters provided!"})
		return
	}
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Role patched!",
	})
}
