package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"
	"github.com/inokone/photostorage/user"
)

// @BasePath /api/v1/photo

// Upload godoc
// @Summary Photo upload endpoint
// @Schemes
// @Description Upload a RAW file with other parameters
// @Accept json
// @Produce json
// @Success 200 {}
// @Router /upload [post]
func Upload(g *gin.Context) {
	file, err := g.FormFile("file")

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"msg": err})
	}
	var raw string
	err = g.SaveUploadedFile(file, raw)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Uploaded file damaged!")
	}
	target, err := createPhoto(
		*currentUser(g),
		filepath.Ext(file.Filename),
		filepath.Base(file.Filename),
		raw,
	)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Uploaded file damaged!")
	}
	store := store()
	store.Store(*target)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Uploaded file could not be stored!")
	}
	g.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("File upload successful for %s.", file.Filename),
	})
}

func createPhoto(user user.User, filename, extension, raw string) (*photo.Photo, error) {
	i, err := image.Factory(extension)
	if err != nil {
		return nil, err
	}
	imported, err := i.Process(bytes.NewReader([]byte(raw)))
	thumbnail, err := image.Thumbnail(imported)
	return &photo.Photo{
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
// @Description Returns all photo descriptors for the current user
// @Accept json
// @Produce json
// @Success 200 {[]Get} List of photo descriptors
// @Router /photos [get]
func List(g *gin.Context) {
	store := store()
	user := currentUser(g)

	result, error := store.List(user.ID.String())

	if error != nil {
		g.JSON(http.StatusNotFound, "Photos do not exist!")
	}

	images := make([]photo.Get, len(result))
	for i, photo := range result {
		p, error := photo.AsGet()
		if error != nil {
			break
		}
		images[i] = *p
	}
	if error != nil {
		g.JSON(http.StatusInternalServerError, "Photos could not be exported!")
	}

	g.JSON(http.StatusOK, images)
}

// Get godoc
// @Summary Get photo endpoint
// @Schemes
// @Description Returns the photo descriptor with the provided ID
// @Accept json
// @Produce json
// @Success 200 {Get} The photo with id
// @Router /photos/:id [get]
func Get(g *gin.Context) {
	id := g.Param("id")
	store := store()

	result, error := store.Get(id)

	if error != nil {
		g.JSON(http.StatusNotFound, "Photos does not exist!")
	}

	exported, error := result.AsGet()

	if error != nil {
		g.JSON(http.StatusInternalServerError, "Photo could not be exported!")
	}

	g.JSON(http.StatusOK, exported)
}

// Download godoc
// @Summary Download RAW file endpoint
// @Schemes
// @Description Returns the RAW file for the provided ID
// @Accept json
// @Produce json
// @Success 200 {[]byte} the RAW file
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

func store() *photo.Store {
	return nil
}

func currentUser(g *gin.Context) *user.User {
	return nil
}
