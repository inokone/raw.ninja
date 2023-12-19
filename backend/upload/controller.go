package upload

import (
	"errors"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/collection"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/photo"
	"github.com/rs/zerolog/log"
)

// Controller is a struct for all REST handlers related to uploads in the application.
type Controller struct {
	uploads  collection.Storer
	service  collection.Service
	uploader *photo.UploadService
	loader   *photo.LoadService
}

// NewController creates a new `Controller` instance based on the collection persistence provided in the parameter.
func NewController(uploads collection.Storer, uploader *photo.UploadService, loader *photo.LoadService) Controller {
	return Controller{
		uploads:  uploads,
		service:  *collection.NewService(uploads),
		uploader: uploader,
		loader:   loader,
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
// @Success 201 {object} string
// @Failure 400 {object} common.StatusMessage
// @Failure 415 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /uploads/ [post]
func (c Controller) Upload(g *gin.Context) {
	var (
		usr   *user.User
		err   error
		form  *multipart.Form
		files []*multipart.FileHeader
		ids   []uuid.UUID
		ch    chan photo.UploadResult
		wg    *sync.WaitGroup
		u     *collection.Collection
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

	ch = make(chan photo.UploadResult, len(files))
	wg = new(sync.WaitGroup)
	for _, file := range files {
		wg.Add(1)
		go c.uploader.Upload(usr, file, ch, wg)
	}
	wg.Wait()
	close(ch)
	for result := range ch {
		if result.Err != nil {
			log.Err(result.Err).Msg("Failed to upload file!")
			g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Uploaded file is corrupt!"})
			return
		}
		ids = append(ids, result.ID)
	}

	u, err = c.service.CreateUpload(*usr, ids)
	if err != nil {
		log.Err(err).Msg("Failed to create upload collection!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Error with the upload. Please try again!"})
		return
	}

	g.JSON(http.StatusCreated, u.ID)
}

// Get is the REST handler for retrieving an upload by ID.
// @Summary Endpoint fore retrieving an upload by ID.
// @Schemes
// @Description Returns an upload by the ID
// @Accept json
// @Produce json
// @Param id path int true "ID of Collection to retrieve"
// @Success 200 {object} collection.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /uploads/:id [get]
func (c Controller) Get(g *gin.Context) {
	var (
		err      error
		cl       *collection.Collection
		id       uuid.UUID
		protocol string
		baseURL  string
		res      collection.Resp
	)
	id, err = uuid.Parse(g.Param("id"))
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid identifier!"})
		return
	}
	cl, err = c.uploads.ByID(id)
	if err != nil {
		log.Err(err).Msg("Failed to retrieve upload!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 404, Message: "Failed to retrieve upload!"})
		return
	}
	if err = authorize(g, cl.UserID); err != nil {
		log.Err(err).Msg("Failed to retrieve upload!")
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}
	res = cl.AsResp()

	protocol = "http"
	if g.Request.TLS != nil {
		protocol = "https"
	}
	baseURL = protocol + "://" + g.Request.Host + "/api/v1/photos/"

	res.Photos, err = c.loader.AsResponse(cl.Photos, baseURL)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to collect images!"})
		return
	}

	g.JSON(http.StatusOK, res)
}

// List is a REST handler for retrieving uploads of a user
// @Summary endpoint for retrieving uploads of a user
// @Schemes
// @Description Returns a list of uploads for the user
// @Accept json
// @Produce json
// @Success 200 {array} collection.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /uploads/ [get]
func (c Controller) List(g *gin.Context) {
	var (
		err     error
		uploads []collection.Collection
		res     []collection.Resp
		user    *user.User
	)
	user, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}
	uploads, err = c.uploads.ByUserAndType(user, collection.Upload)
	if err != nil {
		log.Err(err).Msg("Failed to list uploads!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 404, Message: "Failed to retrieve uploads!"})
		return
	}
	for i, upload := range uploads {
		res[i] = upload.AsResp()
	}
	g.JSON(http.StatusOK, res)
}

func authorize(g *gin.Context, userID uuid.UUID) error {
	user, err := currentUser(g)
	if err != nil {
		return err
	}
	if userID != user.ID {
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
