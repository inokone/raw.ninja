package collection

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/rs/zerolog/log"
)

// Controller is a struct for all REST handlers related to collections in the application.
type Controller struct {
	collections Storer
	service     Service
}

// NewController creates a new `Controller` instance based on the collection persistence provided in the parameter.
func NewController(collections Storer) Controller {
	return Controller{
		collections: collections,
		service:     *NewService(collections),
	}
}

// CreateAlbum is the REST handler for creating an album collection.
// @Summary Endpoint for creating an album collection.
// @Schemes
// @Description Creates an album collection
// @Accept json
// @Produce json
// @Success 200 {object} collection.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /collections/ [post]
func (c Controller) CreateAlbum(g *gin.Context) {
	var (
		ca     CreateAlbum
		err    error
		usr    *user.User
		result *Collection
	)

	usr, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}

	if err := g.ShouldBindJSON(&ca); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationMessage(err))
		return
	}

	result, err = c.service.CreateAlbum(*usr, ca.Name, ca.PhotoIDs)
	if err != nil {
		log.Err(err).Msg("Failed to create album!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	g.JSON(http.StatusOK, result.AsResp())
}

// Get is the REST handler for retrieving a collection by ID.
// @Summary Endpoint fore retrieving a collection by ID.
// @Schemes
// @Description Returns a collection by the ID
// @Accept json
// @Produce json
// @Param id path int true "ID of Collection to retrieve"
// @Success 200 {object} collection.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /collections/:id [get]
func (c Controller) Get(g *gin.Context) {
	var (
		err    error
		result *Collection
		id     uuid.UUID
	)
	id, err = uuid.Parse(g.Param("id"))
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid identifier!"})
		return
	}
	result, err = c.collections.ByID(id)
	if err != nil {
		log.Err(err).Msg("Failed to retrieve collection!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 404, Message: "Failed to retrieve collection!"})
		return
	}
	if err = authorize(g, result.Owner.ID); err != nil {
		log.Err(err).Msg("Failed to retrieve collection!")
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}
	g.JSON(http.StatusOK, result.AsResp())
}

// List is a REST handler for retrieving collections of a user
// @Summary endpoint for retrieving collections of a user
// @Schemes
// @Description Returns a list of collection for the user
// @Accept json
// @Produce json
// @Success 200 {array} collection.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /collections/ [get]
func (c Controller) List(g *gin.Context) {
	var (
		err    error
		albums []Collection
		res    []Resp
		user   *user.User
	)
	user, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}
	albums, err = c.collections.ByUserAndType(user, Album)
	if err != nil {
		log.Err(err).Msg("Failed to list albums!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 404, Message: "Failed to retrieve collection!"})
		return
	}
	for i, album := range albums {
		res[i] = album.AsResp()
	}
	g.JSON(http.StatusOK, res)
}

// Delete is the REST handler for deleting a collection by ID.
// @Summary Endpoint fore deleting a collection by ID.
// @Schemes
// @Description Deletes a collection by the ID
// @Accept json
// @Produce json
// @Param id path int true "ID of Collection to delete"
// @Success 200 {object} collection.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /collections/:id [delete]
func (c Controller) Delete(g *gin.Context) {
	var (
		err    error
		result *Collection
		id     uuid.UUID
	)
	id, err = uuid.Parse(g.Param("id"))
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid identifier!"})
		return
	}
	result, err = c.collections.ByID(id)
	if err != nil {
		log.Err(err).Msg("Failed to delete album!")
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Failed to retrieve collection!"})
		return
	}
	if err = authorize(g, result.Owner.ID); err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}
	if err = c.collections.Delete(id); err != nil {
		log.Err(err).Msg("Failed to delete album!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}
	g.JSON(http.StatusOK, result.AsResp())
}

func authorize(g *gin.Context, userID uuid.UUID) error {
	user, err := currentUser(g)
	if err != nil {
		return err
	}
	if userID.String() != user.ID.String() {
		return errors.New("user is not authorized")
	}
	return nil
}

func currentUser(g *gin.Context) (*user.User, error) {
	u, ok := g.Get("user")
	if !ok {
		return nil, errors.New("user could not be extracted from session")
	}
	usr := u.(*user.User)
	return usr, nil
}
