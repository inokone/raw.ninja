package role

const (
	// RoleAdmin is the administrator role type. Has value 0.
	RoleAdmin = iota
	// RoleFreeTier is the role for free tier users. Has value of 1.
	RoleFreeTier
)

// Role is a struct representing the user role representation for database storage.
type Role struct {
	RoleType    int `gorm:"primary_key"`
	Quota       int64
	DisplayName string `gorm:"type:varchar(100)"`
}

// ProfileRole is astruct, the JSON representation of the `Role` entity for profile and admin views.
type ProfileRole struct {
	ID    int    `json:"id" binding:"required"`
	Quota int64  `json:"quota" binding:"required"`
	Name  string `json:"name" binding:"required,len=3"`
}

// AsProfileRole is a method of the `Role` struct. It converts a `Role` object into a `ProfileRole` object.
func (u *Role) AsProfileRole() ProfileRole {
	return ProfileRole{
		ID:    u.RoleType,
		Quota: u.Quota,
		Name:  u.DisplayName,
	}
}
