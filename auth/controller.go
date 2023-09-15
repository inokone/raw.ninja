package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/common"
)

type Controller struct {
	store Store
}

func NewController(db *gorm.DB) Controller {
	store := Store{db: db}

	return Controller{
		store: store,
	}
}

// @BasePath /api/v1/auth

// Register godoc
// @Summary User registration endpoint
// @Schemes
// @Description Registers the user
// @Accept json
// @Produce json
// @Success 200
// @Router /signup [post]
func (c Controller) Signup(g *gin.Context) {
	var s Registration
	err := g.Bind(&s)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Incorrect user registration dara provided!",
		})
		return
	}
	user, err := NewUser(s.Email, s.Password, s.Phone)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Could not encrypt password???",
		})
		return
	}
	err = c.store.Store(*user)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{
			Code:    400,
			Message: "Could not store user.",
		})
		return
	}
	g.JSON(http.StatusCreated, common.StatusMessage{
		Code:    201,
		Message: "User has been created!",
	})
}

// Login godoc
// @Summary User login endpoint
// @Schemes
// @Description Logs in the user, sets the necessary cookies
// @Accept json
// @Produce json
// @Success 200
// @Router /login [post]
func (c Controller) Login(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, common.StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}

// Reset godoc
// @Summary Reset password endpoint
// @Schemes
// @Description Returns the status and version of the application
// @Accept json
// @Produce json
// @Success 200
// @Router /reset [post]
func (c Controller) Reset(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, common.StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}

// Logout godoc
// @Summary Logout endpoint
// @Schemes
// @Description Logs out of the application
// @Accept json
// @Produce json
// @Success 200
// @Router /logout [get]
func (c Controller) Logout(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, common.StatusMessage{
		Code:    501,
		Message: "Functionality has not been implemented yet!",
	})
}
