package role

const (
	// RoleAdmin is the administrator role type. Has value 1.
	RoleAdmin = iota + 1
	// RoleFreeTier is the role for free tier users. Has value of 2.
	RoleFreeTier
)

// Role is a struct representing the user role representation for database storage.
type Role struct {
	RoleType    int `gorm:"primary_key"`
	Quota       int64
	DisplayName string `gorm:"type:varchar(100)"`
}

// ProfileRole is a struct, the JSON representation of the `Role` entity for profile and admin views.
type ProfileRole struct {
	RoleType    int    `json:"id"`
	Quota       int64  `json:"quota"`
	DisplayName string `json:"name"`
}

// AsProfileRole is a method of the `Role` struct. It converts a `Role` object into a `ProfileRole` object.
func (u *Role) AsProfileRole() ProfileRole {
	return ProfileRole{
		RoleType:    u.RoleType,
		Quota:       u.Quota,
		DisplayName: u.DisplayName,
	}
}
