package rule

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"gorm.io/gorm"
)

var (
	// Delete is an action to delete a photo from storage
	Delete = Action{1, "Delete", false}
	// MoveTo is an action to move photo to the storage specified by the "Action target"
	MoveTo  = Action{2, "Move to", true}
	actions = []Action{Delete, MoveTo}

	// StandardStorage is a target of an action, planned to count 1x from quota
	StandardStorage = Target{1, "Standard storage"}
	// FrozenStorage is a target of an action, planned to count 0.5x from quota
	FrozenStorage = Target{2, "Frozen storage"}
	// Bin is storage for files marked as deleted (not sure this is needed at all)
	Bin     = Target{3, "Bin"}
	targets = []Target{StandardStorage, FrozenStorage, Bin}
)

// Rule is a struct for lifecycle events on stored photos. A rule is a single action that needs to be executed at a specified time
type Rule struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	User        user.User `gorm:"foreignKey:UserID"`
	UserID      uuid.UUID `gorm:"index"`
	Name        string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"type:varchar(1000)"`
	Timing      int       // number of days from creation
	ActionID    int
	TargetID    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

// Resp is a JSON type for rule responses
type Resp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Timing      int    `json:"timing"`
	Action      Action `json:"action"`
	Target      Target `json:"target"`
}

// AsResp is a method of `Rule` to convert to JSON representation.
func (r Rule) AsResp() (Resp, error) {
	var (
		a   Action
		t   Target
		err error
		res Resp
	)
	a, err = ActionFor(r.ActionID)
	if err != nil {
		return Resp{}, err
	}

	res = Resp{
		ID:          r.ID.String(),
		Name:        r.Name,
		Description: r.Description,
		Timing:      r.Timing,
		Action:      a,
	}

	if r.TargetID != 0 {
		t, err = TargetFor(r.TargetID)
		if err != nil {
			return Resp{}, err
		}
		res.Target = t
	}
	return res, nil
}

// Action is a struct for events that can be defined by rules, e.g. delete, move to
type Action struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Targeted bool   `json:"targeted"`
}

// Target is a struct defining target for an `Action`, mostly storage types. Can be extended later. e.g. Standard storage, Frozen storage, Bin
type Target struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Constants is a struct containing the resources that can be used for creating a new rule
type Constants struct {
	Actions []Action `json:"actions"`
	Targets []Target `json:"targets"`
}

// CreateRule is a type for REST representation of data for creating a rule
type CreateRule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Timing      int    `json:"timing"`
	ActionID    int    `json:"action_id"`
	TargetID    int    `json:"target_id"`
}

// TargetFor returns a Target object for an ID
func TargetFor(ID int) (Target, error) {
	for _, target := range targets {
		if target.ID == ID {
			return target, nil
		}
	}
	return Target{}, fmt.Errorf("no target defined for ID %v", ID)
}

// ActionFor returns an Action object for an ID
func ActionFor(ID int) (Action, error) {
	for _, action := range actions {
		if action.ID == ID {
			return action, nil
		}
	}
	return Action{}, fmt.Errorf("no action defined for ID %v", ID)
}
