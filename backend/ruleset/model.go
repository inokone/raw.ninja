package ruleset

import (
	"time"

	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/ruleset/rule"
	"gorm.io/gorm"
)

const (
	// MaxRuleCount is the limit of rules for a rule set - it is not rational to have more tha 3 lifercycle rules at the moment
	MaxRuleCount = 3
)

// RuleSet is a struct for defining a set of rules for photos or collection of photos
type RuleSet struct {
	ID          uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	User        user.User   `gorm:"foreignKey:UserID"`
	UserID      uuid.UUID   `gorm:"index"`
	Name        string      `gorm:"type:varchar(255)"`
	Description string      `gorm:"type:varchar(1000)"`
	Rules       []rule.Rule `gorm:"many2many:ruleset_rules;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

// AsResp is a method of `RuleSet` to convert to JSON representation.
func (r RuleSet) AsResp() (Resp, error) {
	var (
		rules []rule.Resp = make([]rule.Resp, len(r.Rules))
		err   error
	)

	for i, r := range r.Rules {
		rules[i], err = r.AsResp()
		if err != nil {
			return Resp{}, err
		}
	}

	return Resp{
		ID:          r.ID.String(),
		Name:        r.Name,
		Description: r.Description,
		Rules:       rules,
		CreatedAt:   r.CreatedAt,
	}, nil
}

// Resp is a JSON type for rule set responses
type Resp struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Rules       []rule.Resp `json:"rules"`
	CreatedAt   time.Time   `json:"created_at" time_format:"unix"`
}

// CreateRuleSet is a type for REST representation of data for creating a rule set
type CreateRuleSet struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
