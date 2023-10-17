package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"
	"gorm.io/gorm"
)

type Controller struct {
	rep photo.Repository
}

func NewController(db *gorm.DB, ir image.Repository) Controller {
	rep := photo.Repository{
		DB: db,
		Ir: ir,
	}
	return Controller{
		rep: rep,
	}
}

// Statistics godoc
// @Summary User statistics endpoint
// @Schemes
// @Tags photos
// @Description Returns the user statistics on stored photos
// @Accept json
// @Produce json
// @Success 200 {object} UserStatistics
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /statistics/user [get]
func (c Controller) UserStatistics(g *gin.Context) {
	userObj, _ := g.Get("user")
	user := userObj.(auth.User)
	statistics := NewUserStatistics(user)
	ps, err := c.rep.Statistics(user.ID.String())
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error, please contact an administrator!"})
	}
	statistics.Photos = ps.Photos
	statistics.Favorites = ps.Favorites
	statistics.UsedSpace = ps.UsedSpace
	g.JSON(http.StatusOK, statistics)
}
