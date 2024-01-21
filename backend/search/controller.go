package search

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/collection"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/photo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

// Controller is a struct containing all handlers about searching for a photo.
type Controller struct {
	photos photo.Storer
	p      bluemonday.Policy
	l      photo.LoadService
	c      *collection.Service
}

// NewController is a function to create a new `Controller` instance based on the photo persistence.
func NewController(photos photo.Storer, loader photo.LoadService, collections *collection.Service) Controller {
	p := bluemonday.StrictPolicy()
	return Controller{
		photos: photos,
		p:      *p,
		l:      loader,
		c:      collections,
	}
}

// Search is a handler for searching the authenticated user's photo descriptors by file name by prefix
// @Summary Quick search user's photo descriptors endpoint, case sensitive prefix search
// @Schemes
// @Tags photos
// @Description Returns all photo descriptors matching the provided search text
// @Accept json
// @Produce json
// @Success 200 {array} search.QuickSearchResp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /search/quick [get]
func (c Controller) Search(g *gin.Context) {
	var (
		usr        *user.User
		phs        []photo.Photo
		als        []collection.ListItem
		ups        []collection.ListItem
		unsafeText string
		searchText string
		err        error
	)

	usr, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}
	unsafeText = g.DefaultQuery("query", "")
	searchText = c.p.Sanitize(unsafeText)

	phs, err = c.photos.Search(usr.ID.String(), searchText)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photos do not exist!"})
		return
	}

	als, err = c.c.SearchAlbums(usr.ID, searchText)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Albums do not exist!"})
		return
	}

	ups, err = c.c.SearchUploads(usr.ID, searchText)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Uploads do not exist!"})
		return
	}

	c.searchJSON(g, unsafeText, phs, als, ups)
}

func (c Controller) searchJSON(g *gin.Context, query string, photos []photo.Photo, albums []collection.ListItem, uploads []collection.ListItem) {
	var (
		baseURL string
		imgs    []photo.Response
		err     error
		als     []collection.ListResp
		ups     []collection.ListResp
	)

	imgs, baseURL, err = c.photosJSON(g, photos)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to collect images!"})
		return
	}

	als, err = c.collectionJSON(albums, baseURL)
	if err != nil {
		log.Err(err).Msg("Failed to list albums!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to list albums!"})
		return
	}

	ups, err = c.collectionJSON(uploads, baseURL)
	if err != nil {
		log.Err(err).Msg("Failed to list uploads!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to list uploads!"})
		return
	}

	g.JSON(http.StatusOK, QuickSearchResp{
		Query:   query,
		Photos:  imgs,
		Albums:  als,
		Uploads: ups,
	})
}

func (c Controller) photosJSON(g *gin.Context, photos []photo.Photo) ([]photo.Response, string, error) {
	var (
		protocol = "http"
		err      error
		baseURL  string
		imgs     []photo.Response
	)
	if g.Request.TLS != nil {
		protocol = "https"
	}
	baseURL = protocol + "://" + g.Request.Host + "/api/v1/photos/"
	imgs, err = c.l.AsResponse(photos, baseURL)
	return imgs, baseURL, err
}

func (c Controller) collectionJSON(cls []collection.ListItem, baseURL string) ([]collection.ListResp, error) {
	var (
		res = make([]collection.ListResp, len(cls))
		err error
	)
	for i, cl := range cls {
		res[i] = cl.AsListResp()
		if cl.Thumbnail == uuid.Nil {
			continue
		}
		res[i].Thumbnail, err = c.l.ThumbnailURL(cl.Thumbnail, baseURL)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// Favorites is a handler for listing the authenticaed user's favorite photo descriptors endpoint
// @Summary Search user's favorite photo descriptors endpoint
// @Schemes
// @Tags photos
// @Description Returns favorite photo descriptors for the authenticated user
// @Accept json
// @Produce json
// @Success 200 {array} photo.Response
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /search/favorites [get]
func (c Controller) Favorites(g *gin.Context) {
	var (
		usr *user.User
		phs []photo.Photo
		res []photo.Response
		err error
	)

	usr, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Error with the session. Please log in again!"})
		return
	}

	phs, err = c.photos.Favorites(usr.ID.String())
	if err != nil {
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Photos do not exist!"})
		return
	}

	res, _, err = c.photosJSON(g, phs)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to collect images!"})
		return
	}
	g.JSON(http.StatusOK, res)
}

func currentUser(g *gin.Context) (*user.User, error) {
	u, ok := g.Get("user")
	if !ok {
		return nil, errors.New("user could not be extracted from session")
	}
	usr := u.(*user.User)
	return usr, nil
}
