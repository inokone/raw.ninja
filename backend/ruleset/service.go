package ruleset

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/ruleset/rule"
	"github.com/rs/zerolog/log"
)

// Service is a type for encapsulating business logic of rule sets
type Service struct {
	sets  Storer
	rules rule.Storer
}

// NewService creates a `Service` instance based on the rule persistence provided in the parameter.
func NewService(sets Storer, rules rule.Storer) Service {
	return Service{
		sets:  sets,
		rules: rules,
	}
}

// InvalidRuleID is an error for invalid, malformed of non-existing IDs of rules
type InvalidRuleID struct {
	ID string
}

// Error is the string representation of an `InvalidRuleID`
func (e InvalidRuleID) Error() string { return fmt.Sprintf("invalid rule ID [%v]", e.ID) }

// InvalidRuleSetID is an error for invalid, malformed of non-existing IDs of rule sets
type InvalidRuleSetID struct {
	ID string
}

// Error is the string representation of an `InvalidRuleSetID`
func (e InvalidRuleSetID) Error() string { return fmt.Sprintf("invalid rule set ID [%v]", e.ID) }

// Update is a method of `Service`. Handles updating rule sets and rules in the rule set.
func (s Service) Update(usr *user.User, set *Resp) (*RuleSet, error) {
	var (
		r        []rule.Rule
		original *RuleSet
		updated  *RuleSet
		id       uuid.UUID
		err      error
	)
	log.Debug().Str("ID", set.ID).Msg("Updating rule set")
	id, err = uuid.Parse(set.ID)
	if err != nil {
		return nil, InvalidRuleSetID{set.ID}
	}

	original, err = s.sets.ByID(id)
	if err != nil {
		return nil, InvalidRuleSetID{set.ID}
	}

	r, err = s.updateRules(original.Rules, set.Rules, usr)
	if err != nil {
		return nil, err
	}

	updated = &RuleSet{
		ID:          id,
		UserID:      usr.ID,
		Name:        set.Name,
		Description: set.Description,
		Rules:       r,
	}
	log.Debug().Str("Name", updated.Name).Str("Description", updated.Description).Int("Rule count", len(r)).Msg("Rule set updated")

	if err = s.sets.Update(updated); err != nil {
		log.Err(err).Msg("Failed to save rule set!")
		return nil, err
	}
	return updated, nil
}

func (s Service) updateRules(original []rule.Rule, updated []rule.Resp, usr *user.User) ([]rule.Rule, error) {
	var (
		r       []rule.Rule
		err     error
		id      uuid.UUID
		ruleMap map[uuid.UUID]rule.Rule = make(map[uuid.UUID]rule.Rule)
	)

	r = make([]rule.Rule, len(updated))

	for i, rl := range updated {
		if len(rl.ID) > 0 {
			id, err = uuid.Parse(rl.ID)
			if err != nil {
				return nil, InvalidRuleID{rl.ID}
			}
		} else {
			id = uuid.Nil
		}
		r[i] = rule.Rule{
			ID:          id,
			Name:        rl.Name,
			Description: rl.Description,
			Timing:      rl.Timing,
			ActionID:    rl.Action.ID,
			TargetID:    rl.Target.ID,
			UserID:      usr.ID,
		}
		// if the ID is 0, we have to create the rule, otherwise update
		if len(rl.ID) == 0 {
			s.rules.Store(&r[i])
			log.Debug().Str("ID", r[i].ID.String()).Str("Name", r[i].Name).Str("Description", r[i].Description).Int("Action", r[i].ActionID).Int("Timing", r[i].Timing).Msg("Created rule")
		} else {
			s.rules.Update(&r[i])
			log.Debug().Str("ID", r[i].ID.String()).Str("Name", r[i].Name).Str("Description", r[i].Description).Int("Action", r[i].ActionID).Int("Timing", r[i].Timing).Msg("Updated rule")
		}
		ruleMap[id] = r[i]
	}
	// Delete rules, that we do not need anymore
	for _, rl := range original {
		if _, ok := ruleMap[rl.ID]; !ok {
			log.Debug().Str("ID", rl.ID.String()).Str("Name", rl.Name).Str("Description", rl.Description).Int("Action", rl.ActionID).Int("Timing", rl.Timing).Msg("Deleted rule")
			err = s.rules.Delete(rl.ID)
			if err != nil {
				return nil, err
			}
		}
	}
	return r, nil
}
