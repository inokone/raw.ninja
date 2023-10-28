package search

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/photo"
	"github.com/microcosm-cc/bluemonday"
)

// Controller is a struct containing all handlers about searching for a photo.
type Controller struct {
	photos photo.Storer
	p      bluemonday.Policy
}

// NewController is a function to create a new `Controller` instance based on the photo persistence.
func NewController(photos photo.Storer) Controller {
	p := bluemonday.UGCPolicy()
	return Controller{
		photos: photos,
		p:      *p,
	}
}

// Search is a handler for searching the authenticated user's photo descriptors by file name by prefix
// @Summary Quick search user's photo descriptors endpoint, case sensitive prefix search
// @Schemes
// @Tags photos
// @Description Returns all photo descriptors matching the provided search text
// @Accept json
// @Produce json
// @Success 200 {array} photo.Response
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /search/quick [get]
func (c Controller) Search(g *gin.Context) {
	user, err := currentUser(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	unsafeText := g.DefaultQuery("query", "")
	searchText := c.p.Sanitize(unsafeText)
	result, err := c.photos.Search(user.ID.String(), searchText)
	if err != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photos do not exist!"})
		return
	}

	images := make([]photo.Response, len(result))
	for i, photo := range result {
		images[i] = photo.AsResp("http://" + g.Request.Host + "/api/v1/photos/" + photo.ID.String())
	}

	g.JSON(http.StatusOK, images)
}

// Favorites is a handler for listing the authenticaed user's favorite photo descriptors endpoint
// @Summary Search user's favorite photo descriptors endpoint
// @Schemes
// @Tags photos
// @Description Returns favorite photo descriptors for the authenticated user
// @Accept json
// @Produce json
// @Success 200 {array} photo.Response
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /search/favorites [get]
func (c Controller) Favorites(g *gin.Context) {
	user, err := currentUser(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	result, err := c.photos.Favorites(user.ID.String())
	if err != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photos do not exist!"})
		return
	}

	images := make([]photo.Response, len(result))
	for i, photo := range result {
		images[i] = photo.AsResp("http://" + g.Request.Host + "/api/v1/photos/" + photo.ID.String())
	}

	g.JSON(http.StatusOK, images)
}

func currentUser(g *gin.Context) (*user.User, error) {
	u, ok := g.Get("user")
	if !ok {
		return nil, errors.New("user could not be extracted from session")
	}
	usr := u.(*user.User)
	return usr, nil
}
