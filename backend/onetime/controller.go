package onetime

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/rs/zerolog/log"
)

const (
	expiry = time.Minute * 10
)

// Controller is a struct for all REST handlers related to one time accesses in the application.
type Controller struct {
	accesses Storer
	images   image.Storer
}

// NewController creates a new `Controller` instance based on the one time access persistence provided in the parameter.
func NewController(acceses Storer, images image.Storer) Controller {
	return Controller{
		accesses: acceses,
		images:   images,
	}
}

// Create is a method of `Controller`. Handles creating one time access.
// @Summary One time access creation endpoint
// @Schemes
// @Tags rule
// @Description Creates a one time access for the user
// @Accept json
// @Produce json
// @Param data body onetime.CreateAccess true "Data provided for creating the one time access"
// @Success 201 {object} onetime.Resp
// @Failure 400 {object} common.StatusMessage
// @Failure 415 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /onetime/ [post]
func (c Controller) Create(g *gin.Context) {
	var (
		ca  CreateAccess
		id  uuid.UUID
		a   Access
		err error
	)

	if err := g.ShouldBindJSON(&ca); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationMessage(err))
		return
	}

	id, err = uuid.Parse(ca.OriginalID)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid identifier!"})
		return
	}

	a = Access{
		OriginalID: id,
		TTL:        time.Now().Add(expiry),
	}

	if err = c.accesses.Store(&a); err != nil {
		log.Err(err).Msg("Failed to create one time access!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	g.JSON(http.StatusCreated, Resp{ID: a.ID.String()})
}

// Raw is a method of `Controller`. Handles requests for downloding binary for a single photo or RAW file via one time access
// authenticated user. The target photo specified by the photo ID in the URL parameter.
// @Summary Download RAW file endpoint via one time access
// @Schemes
// @Tags photos
// @Description Returns the RAW file for the provided one time access ID
// @Accept json
// @Produce json
// @Param id path int true "one time access ID of the RAW photo to download"
// @Success 200 {array} byte
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /onetime/raw/:id [get]
func (c Controller) Raw(g *gin.Context) {
	var (
		id     uuid.UUID
		access *Access
		err    error
	)

	id, err = uuid.Parse(g.Param("id"))
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid identifier!"})
		return
	}

	access, err = c.accesses.ByID(id)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Resource not found or expired!"})
		return
	}

	raw, err := c.images.LoadImage(access.OriginalID.String())
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Resource not found or expired!"})
		return
	}

	g.Header("Content-Description", "File Transfer")
	g.Header("Content-Disposition", "attachment; filename=edited")
	g.Data(http.StatusOK, "application/octet-stream", raw)
}
