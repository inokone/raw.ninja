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
		db: db,
		ir: ir,
	}

	return Controller{
		rep: rep,
	}
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
// @Success 201 {object} common.StatusMessage
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

	if err = c.rep.Create(*target); err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Uploaded file could not be stored!"})
		return
	}

	g.JSON(http.StatusCreated, common.StatusMessage{
		Code:    201,
		Message: fmt.Sprintf("File upload successful for %s.", file.Filename),
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
		images[i] = photo.AsResp()
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
	id := g.Param("id")

	result, error := c.rep.Get(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Photos does not exist!",
		})
		return
	}

	g.JSON(http.StatusOK, result.AsResp())
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
	id := g.Param("id")

	raw, error := c.rep.Raw(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Raw file does not exist!",
		})
		return
	}

	img, error := c.rep.Get(id)
	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Raw file does not exist!",
		})
		return
	}

	fileName := img.Desc.FileName
	g.Header("Content-Disposition", "attachment; filename="+fileName)
	g.Header("Content-Type", "application/text/plain")
	g.Header("Accept-Length", fmt.Sprintf("%d", len(raw)))
	g.Writer.Write(raw)
	g.JSON(http.StatusOK, common.StatusMessage{
		Code:    200,
		Message: "Download successful!",
	})
}

func currentUser(g *gin.Context) (*auth.User, error) {
	user, ok := g.Get("user")
	if !ok {
		return nil, errors.New("user could not be extracted from session")
	}
	userObj := user.(auth.User)
	return &userObj, nil
}
