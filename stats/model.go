package stats

import "github.com/inokone/photostorage/auth/user"

// UserStats is a JSON representation for aggregated data on a user's photos.
type UserStats struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Registration int64  `json:"registration_date"`
	Photos       int    `json:"photos"`
	Favorites    int    `json:"favorites"`
	UsedSpace    int64  `json:"used_space"`
}

// NewUserStats function creates a new `UserStats` entity for the user provided in the parameters.
func NewUserStats(u user.User) UserStats {
	return UserStats{
		ID:           u.ID.String(),
		Email:        u.Email,
		Phone:        u.Phone,
		Registration: u.CreatedAt.Unix(),
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
}
