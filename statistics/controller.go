package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/photo"
)

// @BasePath /api/v1/statistics

// Controller is a collection of handlers for statistical and aggregated data.
type Controller struct {
	photos photo.Storer
}

// NewController creates a new `Controller` instance based on the photo persistence provided in the parameters.
func NewController(photos photo.Storer) Controller {
	return Controller{
		photos: photos,
	}
}

// UserStatistics is a method of `Controller` returning aggregated data on the photos of a user.
// @Summary User statistics endpoint
// @Schemes
// @Description Returns the user statistics on stored photos
// @Accept json
// @Produce json
// @Success 200 {object} UserStatistics
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /user [get]
func (c Controller) UserStatistics(g *gin.Context) {
	userObj, _ := g.Get("user")
	user := userObj.(auth.User)
	statistics := NewUserStatistics(user)
	ps, err := c.photos.Statistics(user.ID.String())
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
	}
	statistics.Photos = ps.Photos
	statistics.Favorites = ps.Favorites
	statistics.UsedSpace = ps.UsedSpace
	g.JSON(http.StatusOK, statistics)
}
