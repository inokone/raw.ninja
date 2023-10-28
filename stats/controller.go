package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/photo"
)

// Controller is a collection of handlers for statistical and aggregated data.
type Controller struct {
	photos photo.Storer
	users  user.Storer
	config common.ImageStoreConfig
}

// NewController creates a new `Controller` instance based on the photo persistence provided in the parameters.
func NewController(photos photo.Storer, users user.Storer, config common.ImageStoreConfig) Controller {
	return Controller{
		photos: photos,
		users:  users,
		config: config,
	}
}

// UserStats is a method of `Controller` returning aggregated data on the photos of a user.
// @Summary User statistics endpoint
// @Schemes
// @Description Returns the user statistics on stored photos
// @Accept json
// @Produce json
// @Success 200 {object} UserStats
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /statistics/user [get]
func (c Controller) UserStats(g *gin.Context) {
	u, _ := g.Get("user")
	usr := u.(*user.User)
	stats := NewUserStats(*usr)
	ps, err := c.photos.UserStats(usr.ID.String())
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
	}
	stats.Photos = ps.Photos
	stats.Favorites = ps.Favorites
	stats.UsedSpace = ps.UsedSpace
	stats.AvailableSpace = -1
	if stats.Quota > 0 {
		stats.AvailableSpace = stats.Quota - ps.UsedSpace
	}
	g.JSON(http.StatusOK, stats)
}

// AppStats is a method of `Controller` returning aggregated data on the application for administrators.
// @Summary App statistics endpoint
// @Schemes
// @Description Returns the app statistics on stored photos
// @Accept json
// @Produce json
// @Success 200 {object} AppStats
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /statistics/app [get]
func (c Controller) AppStats(g *gin.Context) {
	stats := AppStats{}

	ps, err := c.photos.Stats()
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
	}
	stats.UsedSpace = ps.UsedSpace
	stats.Favorites = ps.Favorites
	stats.Photos = ps.Photos
	stats.Quota = c.config.Quota

	us, err := c.users.Stats()
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
	}
	stats.TotalUsers = us.TotalUsers
	stats.UserDistribution = us.Distribution

	g.JSON(http.StatusOK, stats)
}
