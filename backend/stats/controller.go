package stats

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/collection"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/photo"
	"github.com/rs/zerolog/log"
)

// Controller is a collection of handlers for statistical and aggregated data.
type Controller struct {
	photos photo.Storer
	users  user.Storer
	albums collection.Storer
	config *common.ImageStoreConfig
}

// NewController creates a new `Controller` instance based on the photo persistence provided in the parameters.
func NewController(photos photo.Storer, users user.Storer, albums collection.Storer, config *common.ImageStoreConfig) Controller {
	return Controller{
		photos: photos,
		users:  users,
		albums: albums,
		config: config,
	}
}

// Users is a method of `Controller` returning aggregated data on all users.
// @Summary Users statistics endpoint
// @Schemes
// @Description Returns the users statistics on stored photos
// @Accept json
// @Produce json
// @Success 200 {array} UserPreview
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /statistics/users [get]
func (c Controller) Users(g *gin.Context) {
	var (
		stats *UserStats
		users []user.User
		err   error
		res   []UserPreview
	)
	users, err = c.users.List()
	if err != nil {
		log.Err(err).Msg("Failed to collect users stats")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
	}
	res = make([]UserPreview, len(users))
	for idx, usr := range users {
		stats, err = c.userStats(&usr)
		if err != nil {
			log.Err(err).Msg("Failed to collect users stats")
			g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
		}
		res[idx] = UserPreview{
			Stats: *stats,
			User:  usr.AsAdminView(),
		}
	}
	g.JSON(http.StatusOK, res)
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
	var (
		stats *UserStats
		usr   *user.User
		err   error
	)
	u, _ := g.Get("user")
	usr = u.(*user.User)
	stats, err = c.userStats(usr)
	if err != nil {
		log.Err(err).Msg("Failed to collect user stats")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
	}
	g.JSON(http.StatusOK, stats)
}

func (c Controller) userStats(usr *user.User) (*UserStats, error) {
	var (
		stats UserStats
		ps    photo.UserStats
		as    []collection.Stat
		err   error
	)

	stats = NewUserStats(*usr)
	ps, err = c.photos.UserStats(usr.ID.String())
	if err != nil {
		return nil, err
	}
	stats.Photos = ps.Photos
	stats.Favorites = ps.Favorites
	stats.UsedSpace = ps.UsedSpace
	stats.AvailableSpace = -1
	if stats.Quota > 0 {
		stats.AvailableSpace = stats.Quota - ps.UsedSpace
	}
	as, err = c.albums.Stats(usr)
	if err != nil {
		return nil, err
	}
	stats.Uploads = make(map[time.Time]int)
	for _, stat := range as {
		if stat.Type == collection.Album {
			stats.Albums++
		} else {
			stats.Uploads[stat.Created] = stat.Photos
		}
	}
	return &stats, nil
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
		log.Err(err).Msg("Failed to collect photo stats")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
		return
	}
	stats.UsedSpace = ps.UsedSpace
	stats.Favorites = ps.Favorites
	stats.Photos = ps.Photos
	stats.Quota = c.config.Quota

	albums, uploads, err := c.albums.StorerStats()
	if err != nil {
		log.Err(err).Msg("Failed to collect album and upload stats")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
		return
	}
	stats.Uploads = uploads
	stats.Albums = albums

	us, err := c.users.Stats()
	if err != nil {
		log.Err(err).Msg("Failed to collect user stats")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
		return
	}
	stats.TotalUsers = us.TotalUsers
	stats.UserDistribution = us.Distribution

	g.JSON(http.StatusOK, stats)
}
