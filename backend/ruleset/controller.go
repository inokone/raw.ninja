package ruleset

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/ruleset/rule"
	"github.com/rs/zerolog/log"
)

// Controller is a struct for all REST handlers related to lifecycle rules in the application.
type Controller struct {
	sets    Storer
	rules   rule.Storer
	service Service
}

// NewController creates a new `Controller` instance based on the rule persistence provided in the parameter.
func NewController(sets Storer, rules rule.Storer) Controller {
	return Controller{
		sets:    sets,
		rules:   rules,
		service: NewService(sets, rules),
	}
}

// Create is a method of `Controller`. Handles creating rule sets.
// @Summary Rule set creation endpoint
// @Schemes
// @Tags rule
// @Description Creates a lifecycle rule set for the user
// @Accept json
// @Produce json
// @Param data body ruleset.CreateRuleSet true "The data to use for creating the ruleset"
// @Success 201 {object} ruleset.Resp
// @Failure 400 {object} common.StatusMessage
// @Failure 415 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rulesets/ [post]
func (c Controller) Create(g *gin.Context) {
	var (
		cr  CreateRuleSet
		r   RuleSet
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

	r = RuleSet{
		UserID:      usr.ID,
		Name:        cr.Name,
		Description: cr.Description,
	}

	if err = c.sets.Store(&r); err != nil {
		log.Err(err).Msg("Failed to create rule set!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	res, err = r.AsResp()
	if err != nil {
		log.Err(err).Msg("Failed to create rule set!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	g.JSON(http.StatusCreated, res)
}

// Update is a method of `Controller`. Handles updating rule sets.
// @Summary Rule set updating endpoint
// @Schemes
// @Tags rule
// @Description Updates a lifecycle rule set for the user
// @Accept json
// @Produce json
// @Param id path int true "ID of the role to update"
// @Param data body ruleset.Resp true "The ruleset to be updated"
// @Success 201 {object} ruleset.Resp
// @Failure 400 {object} common.StatusMessage
// @Failure 415 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rulesets/:id [put]
func (c Controller) Update(g *gin.Context) {
	var (
		ur  Resp
		rs  *RuleSet
		err error
		usr *user.User
		res Resp
		id  uuid.UUID
	)

	usr, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}

	if err = g.ShouldBindJSON(&ur); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationMessage(err))
		return
	}

	id, err = uuid.Parse(ur.ID)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: fmt.Sprintf("Invalid rule set ID: %v", ur.ID)})
		return
	}

	rs, err = c.sets.ByID(id)
	if err != nil {
		log.Err(err).Msg("Failed to retrieve rule set!")
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Failed to retrieve rule set!"})
		return
	}
	if err = authorize(g, rs.UserID); err != nil {
		log.Err(err).Msg("Failed to retrieve rule set!")
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}

	rs, err = c.service.Update(usr, &ur)
	if err != nil {
		switch err.(type) {
		default:
			g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error"})
		case *InvalidRuleSetID:
			g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: fmt.Sprintf("Invalid rule set ID: %v", ur.ID)})
		case *InvalidRuleID:
			g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: fmt.Sprintf("Invalid rule set ID: %v", ur.ID)})
		}
		return
	}

	res, err = rs.AsResp()
	if err != nil {
		log.Err(err).Msg("Failed to save rule set!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	g.JSON(http.StatusCreated, res)
}

// Get is the REST handler for retrieving a rule set by ID.
// @Summary Endpoint fore retrieving a rule set by ID.
// @Schemes
// @Description Returns a rule set by the ID
// @Accept json
// @Produce json
// @Param id path int true "ID of rule set to retrieve"
// @Success 200 {object} ruleset.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rulesets/:id [get]
func (c Controller) Get(g *gin.Context) {
	var (
		err error
		rs  *RuleSet
		id  uuid.UUID
		res Resp
	)
	id, err = uuid.Parse(g.Param("id"))
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid identifier!"})
		return
	}
	rs, err = c.sets.ByID(id)
	if err != nil {
		log.Err(err).Msg("Failed to retrieve rule set!")
		g.AbortWithStatusJSON(http.StatusNotFound, common.StatusMessage{Code: 404, Message: "Failed to retrieve rule set!"})
		return
	}
	if err = authorize(g, rs.UserID); err != nil {
		log.Err(err).Msg("Failed to retrieve rule set!")
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}

	res, err = rs.AsResp()
	if err != nil {
		log.Err(err).Msg("Failed to retrieve rule set!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
		return
	}

	g.JSON(http.StatusOK, res)
}

// List is a REST handler for retrieving rule sets of a user
// @Summary endpoint for retrieving rule sets of a user
// @Schemes
// @Description Returns a list of rule sets for the user
// @Accept json
// @Produce json
// @Success 200 {array} ruleset.Resp
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rulesets/ [get]
func (c Controller) List(g *gin.Context) {
	var (
		err      error
		ruleSets []RuleSet
		res      []Resp
		user     *user.User
	)
	user, err = currentUser(g)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}
	ruleSets, err = c.sets.ByUser(user.ID)
	if err != nil {
		log.Err(err).Msg("Failed to list rule sets!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to list rule sets!"})
		return
	}

	res = make([]Resp, len(ruleSets))
	for i, ruleSet := range ruleSets {
		res[i], err = ruleSet.AsResp()
		if err != nil {
			log.Err(err).Msg("Failed to list rule sets!")
			g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Unknown error!"})
			return
		}
	}
	g.JSON(http.StatusOK, res)
}

// Delete is the REST handler for deleting a rule set by ID.
// @Summary Endpoint fore deleting a rule set by ID.
// @Schemes
// @Description Deletes a rule set by the ID
// @Accept json
// @Produce json
// @Param id path int true "ID of rule set to delete"
// @Success 200 {object} common.StatusMessage
// @Failure 404 {object} common.StatusMessage
// @Failure 500 {object} common.StatusMessage
// @Router /rulesets/:id [delete]
func (c Controller) Delete(g *gin.Context) {
	var (
		err error
		rs  *RuleSet
		id  uuid.UUID
	)
	id, err = uuid.Parse(g.Param("id"))
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, common.StatusMessage{Code: 400, Message: "Invalid identifier!"})
		return
	}
	rs, err = c.sets.ByID(id)
	if err != nil {
		log.Err(err).Msg("Failed to delete rule set!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to delete rule set!"})
		return
	}
	if err = authorize(g, rs.UserID); err != nil {
		log.Err(err).Msg("Failed to delete rule set!")
		g.AbortWithStatusJSON(http.StatusUnauthorized, common.StatusMessage{Code: 401, Message: "Unauthorized!"})
		return
	}

	if err = c.sets.Delete(id); err != nil {
		log.Err(err).Msg("Failed to delete rule set!")
		g.AbortWithStatusJSON(http.StatusInternalServerError, common.StatusMessage{Code: 500, Message: "Failed to delete rule set!"})
		return
	}

	g.JSON(http.StatusOK, common.StatusMessage{Code: 200, Message: "Success!"})
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
