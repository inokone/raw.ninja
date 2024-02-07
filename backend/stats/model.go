package stats

import (
	"time"

	"github.com/inokone/photostorage/auth/user"
)

// UserStats is a JSON representation for aggregated data on a user's photos.
type UserStats struct {
	ID             string            `json:"id"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	Email          string            `json:"email"`
	Role           string            `json:"role"`
	Registration   int64             `json:"registration_date"`
	Photos         int               `json:"photos"`
	Favorites      int               `json:"favorites"`
	Albums         int               `json:"albums"`
	Uploads        map[time.Time]int `json:"uploads"`
	UsedSpace      int64             `json:"used_space"`
	AvailableSpace int64             `json:"available_space"`
	Quota          int64             `json:"quota"`
}

// UserPreview is a view on the User with statistical data
type UserPreview struct {
	User  user.AdminView `json:"user"`
	Stats UserStats      `json:"stats"`
}

// NewUserStats function creates a new `UserStats` entity for the user provided in the parameters.
func NewUserStats(u user.User) UserStats {
	return UserStats{
		ID:           u.ID.String(),
		Email:        u.Email,
		Role:         u.Role.DisplayName,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Registration: u.CreatedAt.Unix(),
		Quota:        u.Role.Quota,
	}
}

// AppStats is a JSON representation for aggregated data on a user's photos.
type AppStats struct {
	TotalUsers       int             `json:"total_users"`
	UserDistribution []user.RoleUser `json:"user_distribution"`
	Photos           int             `json:"photos"`
	Favorites        int             `json:"favorites"`
	UsedSpace        int64           `json:"used_space"`
	Quota            int64           `json:"quota"`
	Uploads          int             `json:"uploads"`
	Albums           int             `json:"albums"`
}
