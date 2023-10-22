package statistics

import "github.com/inokone/photostorage/auth"

// UserStatistics is a JSON representation for aggregated data on a user's photos.
type UserStatistics struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Registration int64  `json:"registration_date"`
	Photos       int    `json:"photos"`
	Favorites    int    `json:"favorites"`
	UsedSpace    int64  `json:"used_space"`
}

// NewUserStatistics function creates a new `UserStatistics` entity for the user provided in the parameters.
func NewUserStatistics(u auth.User) UserStatistics {
	return UserStatistics{
		ID:           u.ID.String(),
		Email:        u.Email,
		Phone:        u.Phone,
		Registration: u.CreatedAt.Unix(),
	}
}
