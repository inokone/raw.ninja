package photo

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/image"
)

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
func Upload(g *gin.Context) {
	file, err := g.FormFile("file")

	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Could not extract uploaded file from request!"})
	}
	var raw string
	err = g.SaveUploadedFile(file, raw)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Uploaded file is damaged!"})
	}
	target, err := createPhoto(
		*currentUser(g),
		filepath.Ext(file.Filename),
		filepath.Base(file.Filename),
		raw,
	)
	if err != nil {
		g.JSON(http.StatusUnsupportedMediaType, common.StatusMessage{Code: 415, Message: "Uploaded file format is not supported!"})
	}
	store := store()
	store.Store(*target)
	if err != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Uploaded file could not be stored!"})
	}
	g.JSON(http.StatusCreated, common.StatusMessage{
		Code:    201,
		Message: fmt.Sprintf("File upload successful for %s.", file.Filename),
	})
}

func createPhoto(user auth.User, filename, extension, raw string) (*Photo, error) {
	i, err := image.Factory(extension)
	if err != nil {
		return nil, err
	}
	imported, err := i.Process(bytes.NewReader([]byte(raw)))
	if err != nil {
		return nil, err
	}
	thumbnail, err := image.Thumbnail(imported)
	if err != nil {
		return nil, err
	}
	return &Photo{
		Desc: descriptor.Descriptor{
			FileName:  filename,
			Format:    descriptor.Formats[extension],
			Uploaded:  time.Now(),
			Thumbnail: thumbnail,
		},
		User: user,
		Raw:  []byte(raw),
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
func List(g *gin.Context) {
	store := store()
	user := currentUser(g)

	result, error := store.List(user.ID.String())

	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photos do not exist!"})
	}

	images := make([]Response, len(result))
	for i, photo := range result {
		p, error := photo.AsResp()
		if error != nil {
			break
		}
		images[i] = *p
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
func Get(g *gin.Context) {
	id := g.Param("id")
	store := store()

	result, error := store.Get(id)

	if error != nil {
		g.JSON(http.StatusNotFound, common.StatusMessage{
			Code:    404,
			Message: "Photos does not exist!",
		})
	}

	exported, error := result.AsResp()

	if error != nil {
		g.JSON(http.StatusInternalServerError, common.StatusMessage{
			Code:    500,
			Message: "Photo could not be exported!",
		})
	}

	g.JSON(http.StatusOK, exported)
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
func Download(g *gin.Context) {
	id := g.Param("id")
	store := store()

	raw, error := store.Raw(id)
	if error != nil {
		g.JSON(http.StatusNotFound, "Raw file does not exist!")
	}

	img, error := store.Get(id)
	if error != nil {
		g.JSON(http.StatusNotFound, "Raw file does not exist!")
	}

	fileName := img.Desc.FileName
	g.Header("Content-Disposition", "attachment; filename="+fileName)
	g.Header("Content-Type", "application/text/plain")
	g.Header("Accept-Length", fmt.Sprintf("%d", len(raw)))
	g.Writer.Write(raw)
	g.JSON(http.StatusOK, gin.H{
		"msg": "File download successful.",
	})
}

func store() *Store {
	return nil
}

func currentUser(g *gin.Context) *auth.User {
	return nil
}
