package photo

import (
	"errors"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"

	"github.com/rs/zerolog/log"
)

var (
	statusNotFound       = common.StatusMessage{Code: 404, Message: "Photo does not exist!"}
	statusMalformedPhoto = common.StatusMessage{Code: 400, Message: "Malformed photo data!"}
	// ErrMalformedRequest is an error for invalid or inconsistent photo data
	ErrMalformedRequest = errors.New("photo data inconsistent")
)

// Controller is a struct for all REST handlers related to photos in the application.
type Controller struct {
	photos Storer
	images image.Storer
	cfg    common.ImageStoreConfig
	s      uploadService
	l      LoadService
}

// NewController creates a new `Controller` instance based on the photo persistence provided in the parameter.
func NewController(photos Storer, images image.Storer, cfg common.ImageStoreConfig) Controller {
	return Controller{
		photos: photos,
		images: images,
		cfg:    cfg,
		s:      *newUploadService(photos, images, cfg),
		l:      *NewLoadService(photos, images, cfg),
	}
}

// Upload is a method of `Controller`. Handles RAW and photo upload requests. Capable of handling multiple files
// uploaded within a single request.
// @Summary Photo upload endpoint
// @Schemes
// @Tags photos
// @Description Upload RAW files to store
// @Accept multipart/form-data
// @Produce json
// @Param files[] formData file true "Photos to store"
// @Success 201 {object} UploadSuccess
// @Failure 400 {object} common.StatusMessage
// @Failure 415 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /photos/upload [post]
func (c Controller) Upload(g *gin.Context) {
	var (
		usr   *user.User
		err   error
		form  *multipart.Form
		files []*multipart.FileHeader
		ids   []string
		ch    chan uploadResult
		wg    *sync.WaitGroup
	)

	form, err = g.MultipartForm()
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Could not extract uploaded file from request!"})
		return
	}

	files = form.File["files[]"]

	if len(files) == 0 {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "You have to upload at least 1 file!"})
		return
	}

	usr, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	ch = make(chan uploadResult, len(files))
	wg = new(sync.WaitGroup)
	for _, file := range files {
		wg.Add(1)
		go c.s.upload(usr, file, ch, wg)
	}
	wg.Wait()
	close(ch)
	for result := range ch {
		if result.err != nil {
			g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Uploaded file is corrupt!"})
		}
		ids = append(ids, result.id.String())
	}

	g.JSON(http.StatusCreated, UploadSuccess{
		PhotoIDs: ids,
		UserID:   usr.ID.String(),
	})
}

// List is a method of `Controller`. Handles listing all photos and RAW files of the authenticated user.
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
	var (
		user     *user.User
		err      error
		result   []Photo
		protocol string
		baseURL  string
	)

	user, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	result, err = c.photos.All(user.ID.String())
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	protocol = "http"
	if g.Request.TLS != nil {
		protocol = "https"
	}
	baseURL = protocol + "://" + g.Request.Host + g.Request.URL.Path
	imgs, err := c.l.AsResponse(result, baseURL)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to collect images!"})
		return
	}
	g.JSON(http.StatusOK, imgs)
}

// Get is a method of `Controller`. Handles requests for retrieving metadata for a single photo or RAW file
// of the authenticated user. The target photo is specified by the photo ID in the URL parameter.
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
	var (
		id       string = g.Param("id")
		result   *Photo
		err      error
		protocol string
		resp     Response
		baseURL  string
	)

	result, err = c.photos.Load(id)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	if err = authorize(g, result.UserID); err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	protocol = "http"
	if g.Request.TLS != nil {
		protocol = "https"
	}
	baseURL = protocol + "://" + g.Request.Host + g.Request.URL.Path + id

	resp = result.AsResp()
	if err = c.l.decorateWithRequest(&resp, baseURL); err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to collect image!"})
		return
	}

	g.JSON(http.StatusOK, resp)
}

// Update is a method of `Controller`. Handles requests for updating a single photo or RAW file of the authenticated user.
// The target photo specified by the photo ID in the URL parameter. Update is limited only for a subset of the metadata.
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
	id := g.Param("id")
	persisted, err := c.photos.Load(id)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	if err = authorize(g, persisted.UserID); err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	var newVersion Response
	if err := g.ShouldBindJSON(&newVersion); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, statusMalformedPhoto)
		return
	}

	if err = applyChange(persisted, newVersion); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, statusMalformedPhoto)
		return
	}

	c.photos.Update(persisted)
	g.JSON(http.StatusOK, common.StatusMessage{Code: 200, Message: "Photo updated!"})
}

// Delete is a method of `Controller`. Handles requests for deleting a single photo or RAW file of the authenticated user.
// The target photo specified by the photo ID in the URL parameter.
// @Summary Delete photo endpoint
// @Schemes
// @Tags photos
// @Description Deletes the photo with the provided ID
// @Accept json
// @Produce json
// @Param id path int true "ID of the photo to delete"
// @Success 200 {object} photo.Response
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /photos/:id [delete]
func (c Controller) Delete(g *gin.Context) {
	id := g.Param("id")
	result, err := c.photos.Load(id)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	if err = authorize(g, result.UserID); err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	if err = c.photos.Delete(id); err != nil {
		log.Err(err).Msg("Failed to delete photo")
		g.AbortWithStatusJSON(http.StatusInternalServerError, statusNotFound)
		return
	}
	g.JSON(http.StatusOK, common.StatusMessage{Code: 200, Message: "Photo deleted!"})
}

func applyChange(persisted *Photo, newVersion Response) error {
	if persisted.ID.String() != newVersion.ID {
		return ErrMalformedRequest
	}
	persisted.Desc.Tags = newVersion.Desc.Tags
	persisted.Desc.Favorite = newVersion.Desc.Favorite
	return nil
}

// Raw is a method of `Controller`. Handles requests for downloding binary for a single photo or RAW file of the
// authenticated user. The target photo specified by the photo ID in the URL parameter.
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
// @Router /photos/:id/raw [get]
func (c Controller) Raw(g *gin.Context) {
	id := g.Param("id")
	img, err := c.photos.Load(id)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	if err = authorize(g, img.UserID); err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	fileName := img.Desc.FileName
	raw, err := c.images.LoadImage(id)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	g.Header("Content-Description", "File Transfer")
	g.Header("Content-Disposition", "attachment; filename="+fileName)
	g.Data(http.StatusOK, "application/octet-stream", raw)
}

// Thumbnail is a method of `Controller`. Handles requests for downloding thumbnail binary for a single photo or RAW file of the
// authenticated user. The target photo specified by the photo ID in the URL parameter.
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
	id := g.Param("id")
	img, err := c.photos.Load(id)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	if err = authorize(g, img.UserID); err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	fileName := img.Desc.FileName
	thumbnail, err := c.images.LoadThumbnail(id)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, statusNotFound)
		return
	}

	g.Header("Content-Description", "File Transfer")
	g.Header("Content-Disposition", "attachment; filename="+fileName)
	g.Data(http.StatusOK, "application/octet-stream", thumbnail)
}

func authorize(g *gin.Context, userID string) error {
	user, err := currentUser(g)
	if err != nil {
		return err
	}
	if userID != user.ID.String() {
		return errors.New("user is not authorized")
	}
	return nil
}

func currentUser(g *gin.Context) (*user.User, error) {
	u, ok := g.Get("user")
	if !ok {
		return nil, errors.New("user could not be extracted from session")
	}
	usr := u.(*user.User)
	return usr, nil
}
