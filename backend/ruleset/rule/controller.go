package rule

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/rs/zerolog/log"
)

// Controller is a struct for all REST handlers related to lifecycle rules in the application.
type Controller struct {
	rules Storer
}

// NewController creates a new `Controller` instance based on the rule persistence provided in the parameter.
func NewController(rules Storer) Controller {
	return Controller{
		rules: rules,
	}
}

// Create is a method of `Controller`. Handles creating rule.
// @Summary Rule creation endpoint
// @Schemes
// @Tags rule
// @Description Creates a lifecycle rule for the user
// @Accept json
// @Produce json
// @Param data body rule.CreateRule true "Data provided for creating the lifecycle rule"
// @Success 201 {object} rule.Resp
// @Failure 400 {object} common.StatusMessage
// @Failure 415 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rules/ [post]
func (c Controller) Create(g *gin.Context) {
	var (
		cr  CreateRule
		r   Rule
		err error
		usr *user.User
		res Resp
	)

	usr, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}

	if err := g.ShouldBindJSON(&cr); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationMessage(err))
		return
	}

	r = Rule{
		UserID:      usr.ID,
		Name:        cr.Name,
		Description: cr.Description,
		Timing:      cr.Timing,
		TargetID:    cr.TargetID,
		ActionID:    cr.ActionID,
	}

	if err = c.rules.Store(&r); err != nil {
		log.Err(err).Msg("Failed to create rule!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	res, err = r.AsResp()
	if err != nil {
		log.Err(err).Msg("Failed to create rule!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	g.JSON(http.StatusCreated, res)
}

// Get is the REST handler for retrieving a rule by ID.
// @Summary Endpoint for retrieving a rule by ID.
// @Schemes
// @Description Returns a rule by the ID
// @Accept json
// @Produce json
// @Param id path int true "ID of rule to retrieve"
// @Success 200 {object} rule.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rules/:id [get]
func (c Controller) Get(g *gin.Context) {
	var (
		err error
		r   *Rule
		id  uuid.UUID
		res Resp
	)
	id, err = uuid.Parse(g.Param("id"))
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid identifier!"})
		return
	}
	r, err = c.rules.ByID(id)
	if err != nil {
		log.Err(err).Msg("Failed to retrieve rule!")
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Failed to retrieve rule!"})
		return
	}
	if err = authorize(g, r.UserID); err != nil {
		log.Err(err).Msg("Failed to retrieve rule!")
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}
	res, err = r.AsResp()
	if err != nil {
		log.Err(err).Msg("Failed to retrieve rule!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	g.JSON(http.StatusOK, res)
}

// List is a REST handler for retrieving rules of a user
// @Summary endpoint for retrieving rules of a user
// @Schemes
// @Description Returns the list of rules for the user
// @Accept json
// @Produce json
// @Success 200 {array} rule.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rules/ [get]
func (c Controller) List(g *gin.Context) {
	var (
		err      error
		ruleSets []Rule
		res      []Resp
		user     *user.User
	)
	user, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}
	ruleSets, err = c.rules.ByUser(user.ID)
	if err != nil {
		log.Err(err).Msg("Failed to list rules!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to list rules!"})
		return
	}

	res = make([]Resp, len(ruleSets))
	for i, ruleSet := range ruleSets {
		res[i], err = ruleSet.AsResp()
		if err != nil {
			log.Err(err).Msg("Failed to list rules!")
			g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
			return
		}
	}
	g.JSON(http.StatusOK, res)
}

// Constants is the REST handler for retrieving constants, that can be used for rule creation.
// @Summary Endpoint for retrieving constants for rule creation.
// @Schemes
// @Description Returns the constants used for rule creation
// @Accept json
// @Produce json
// @Success 200 {object} rule.Constants
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rules/constants [get]
func (c Controller) Constants(g *gin.Context) {
	g.JSON(http.StatusOK, Constants{
		Actions: actions,
		Targets: targets,
	})
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
