package photo

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/image"
	"gorm.io/gorm"
)

type Controller struct {
	rep Repository
}

func NewController(db *gorm.DB, ir image.Repository) Controller {
	rep := Repository{
		DB: db,
		Ir: ir,
	}

	return Controller{
		rep: rep,
	}
}

type UploadSuccess struct {
	PhotoID string `json:"photoId"`
	UserID  string `json:"userId"`
}

// @BasePath /api/v1/photo

// Upload godoc
// @Summary Photo upload endpoint
// @Schemes
// @Tags photos
// @Description Upload a RAW file with descriptor
// @Accept json
// @Produce json
// @Param photo formData file true "Photo to store"
// @Success 201 {object} UploadSuccess
// @Failure 400 {object} common.StatusMessage
// @Failure 415 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /upload [post]
func (c Controller) Upload(g *gin.Context) {
	user, error := currentUser(g)
	if error != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	file, err := g.FormFile("file")
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Could not extract uploaded file from request!"})
		return
	}

	mp, err := file.Open()
	if err != nil {
		g.Error(err)
		return
	}
	defer mp.Close()
	raw, err := io.ReadAll(mp)
	if err != nil {
		g.JSON(http.StatusUnsupportedMediaType, common.StatusMessage{Code: 415, Message: "Uploaded file is corrupt!"})
		return
	}
	target, err := createPhoto(
		*user,
		filepath.Base(file.Filename),
		filepath.Ext(file.Filename)[1:],
		raw,
	)
	if err != nil {
		g.JSON(http.StatusUnsupportedMediaType, common.StatusMessage{Code: 415, Message: fmt.Sprintf("Uploaded file format is not supported! Cause: %v", err)})
		return
	}

	id, err := c.rep.Create(*target)
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Uploaded file could not be stored!"})
		return
	}

	g.JSON(http.StatusCreated, UploadSuccess{
		PhotoID: id.String(),
		UserID:  user.ID.String(),
	})
}

func createPhoto(user auth.User, filename, extension string, raw []byte) (*Photo, error) {
	i := image.NewImporter()
	thumbnail, err := i.Thumbnail(raw)
	if err != nil {
		return nil, err
	}
	metadata, err := i.Describe(raw)
	if err != nil {
		return nil, err
	}
	return &Photo{
		Desc: descriptor.Descriptor{
			FileName:  filename,
			Format:    descriptor.ParseFormat(extension),
			Uploaded:  time.Now(),
			Thumbnail: thumbnail,
			Metadata:  *metadata,
		},
		User: user,
		Raw:  raw,
	}, nil
}

// List godoc
// @Summary List user's photo descriptors endpoint
// @Schemes
// @Tags photos
// @Description Returns all photo descriptors for the current user
// @Accept json
// @Produce json
// @Success 200 {array} photo.Response
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /photos [get]
func (c Controller) List(g *gin.Context) {
	user, error := currentUser(g)
	if error != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	result, error := c.rep.All(user.ID.String())
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photos do not exist!"})
		return
	}

	images := make([]Response, len(result))
	for i, photo := range result {
		images[i] = photo.AsResp("http://" + g.Request.Host + g.Request.URL.Path + photo.ID.String())
	}
	if error != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Photos could not be exported!"})
	}

	g.JSON(http.StatusOK, images)
}

// Get godoc
// @Summary Get photo endpoint
// @Schemes
// @Tags photos
// @Description Returns the photo descriptor with the provided ID
// @Accept json
// @Produce json
// @Param id path int true "ID of the photo information to collect"
// @Success 200 {object} photo.Response
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /photos/:id [get]
func (c Controller) Get(g *gin.Context) {
	user, error := currentUser(g)
	if error != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	id := g.Param("id")

	result, error := c.rep.Get(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Photo does not exist!",
		})
		return
	}

	// Accessing other user's photos is forbidden
	if result.UserID != user.ID.String() {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photo does not exist!"})
		return
	}

	g.JSON(http.StatusOK, result.AsResp("http://"+g.Request.Host+g.Request.URL.Path))
}

// Update godoc
// @Summary Update photo endpoint for tags and favorite setting
// @Schemes
// @Tags photos
// @Description Updates tags and favorite setting for RAW file
// @Accept json
// @Produce json
// @Param id path int true "ID of the photo information to collect"
// @Success 200 {object} common.StatusMessage
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /photos/:id [put]
func (c Controller) Update(g *gin.Context) {
	user, error := currentUser(g)
	if error != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	id := g.Param("id")
	persisted, error := c.rep.Get(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photo does not exist!"})
		return
	}

	// Changing other user's photos is forbidden
	if persisted.UserID != user.ID.String() {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photo does not exist!"})
		return
	}

	var newVersion Response
	if err := g.Bind(&newVersion); err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{Code: 200, Message: "Malformed photo data!"})
		return
	}

	err := applyChange(&persisted, newVersion)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{Code: 200, Message: err.Error()})
		return
	}

	c.rep.Update(persisted)
	g.JSON(http.StatusOK, common.StatusMessage{Code: 200, Message: "Photo updated!"})
}

func applyChange(persisted *Photo, newVersion Response) error {
	if persisted.ID.String() != newVersion.ID {
		return errors.New("photo data inconsistent, ID does not match the path")
	}
	persisted.Desc.Tags = newVersion.Desc.Tags
	persisted.Desc.Favorite = newVersion.Desc.Favorite
	return nil
}

// Download godoc
// @Summary Download RAW file endpoint
// @Schemes
// @Tags photos
// @Description Returns the RAW file for the provided ID
// @Accept json
// @Produce json
// @Param id path int true "ID of the RAW photo to download"
// @Success 200 {array} byte
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /photos/:id/download [get]
func (c Controller) Download(g *gin.Context) {
	user, error := currentUser(g)
	if error != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	id := g.Param("id")
	img, error := c.rep.Get(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Raw file does not exist!",
		})
		return
	}

	if img.UserID != user.ID.String() {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Raw file does not exist!"})
		return
	}
	fileName := img.Desc.FileName

	raw, error := c.rep.Raw(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Raw file does not exist!",
		})
		return
	}

	g.Header("Content-Description", "File Transfer")
	g.Header("Content-Disposition", "attachment; filename="+fileName)
	g.Data(http.StatusOK, "application/octet-stream", raw)
}

// Thumbnail godoc
// @Summary Thumbnail image endpoint
// @Schemes
// @Tags photos
// @Description Returns the thumbnail for the provided ID
// @Accept json
// @Produce json
// @Param id path int true "ID of the thumbnail to download"
// @Success 200 {array} byte
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /photos/:id/thumbnail [get]
func (c Controller) Thumbnail(g *gin.Context) {
	user, error := currentUser(g)
	if error != nil {
		g.JSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	id := g.Param("id")
	img, error := c.rep.Get(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Image file does not exist!",
		})
		return
	}
	if img.UserID != user.ID.String() {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Image file does not exist!"})
		return
	}
	fileName := img.Desc.FileName

	thumbnail, error := c.rep.Thumbnail(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Image file does not exist!",
		})
		return
	}

	g.Header("Content-Description", "File Transfer")
	g.Header("Content-Disposition", "attachment; filename="+fileName)
	g.Data(http.StatusOK, "application/octet-stream", thumbnail)
}

func currentUser(g *gin.Context) (*auth.User, error) {
	user, ok := g.Get("user")
	if !ok {
		return nil, errors.New("user could not be extracted from session")
	}
	userObj := user.(auth.User)
	return &userObj, nil
}
