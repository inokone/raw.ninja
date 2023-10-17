package search

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

type Controller struct {
	rep photo.Repository
	p   bluemonday.Policy
}

func NewController(db *gorm.DB, ir image.Repository) Controller {
	rep := photo.Repository{
		DB: db,
		Ir: ir,
	}
	p := bluemonday.UGCPolicy()
	return Controller{
		rep: rep,
		p:   *p,
	}
}

// Search godoc
// @Summary Quick search user's photo descriptors endpoint
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
	user, error := currentUser(g)
	if error != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	unsafeText := g.DefaultQuery("query", "")
	searchText := c.p.Sanitize(unsafeText)
	result, error := c.rep.Search(user.ID.String(), searchText)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photos do not exist!"})
		return
	}

	images := make([]photo.Response, len(result))
	for i, photo := range result {
		images[i] = photo.AsResp("http://" + g.Request.Host + "/api/v1/photos/" + photo.ID.String())
	}
	if error != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Photos could not be exported!"})
	}

	g.JSON(http.StatusOK, images)
}

// Favorites godoc
// @Summary Search user's favorite photo descriptors endpoint
// @Schemes
// @Tags photos
// @Description Returns favorite photo descriptors
// @Accept json
// @Produce json
// @Success 200 {array} photo.Response
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /search/favorites [get]
func (c Controller) Favorites(g *gin.Context) {
	user, error := currentUser(g)
	if error != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	result, error := c.rep.Favorites(user.ID.String())
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photos do not exist!"})
		return
	}

	images := make([]photo.Response, len(result))
	for i, photo := range result {
		images[i] = photo.AsResp("http://" + g.Request.Host + "/api/v1/photos/" + photo.ID.String())
	}
	if error != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Photos could not be exported!"})
	}

	g.JSON(http.StatusOK, images)
}

func currentUser(g *gin.Context) (*auth.User, error) {
	user, ok := g.Get("user")
	if !ok {
		return nil, errors.New("user could not be extracted from session")
	}
	userObj := user.(auth.User)
	return &userObj, nil
}
