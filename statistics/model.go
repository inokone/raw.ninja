package statistics

import "github.com/inokone/photostorage/auth"

type UserStatistics struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Registration int64  `json:"registration_date"`
	Photos       int    `json:"photos"`
	Favorites    int    `json:"favorites"`
	UsedSpace    int64  `json:"used_space"`
}

func NewUserStatistics(u auth.User) UserStatistics {
	return UserStatistics{
		ID:           u.ID.String(),
		Email:        u.Email,
		Phone:        u.Phone,
		Registration: u.CreatedAt.Unix(),
	}
}
