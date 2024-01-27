package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/role"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is the user representation for database storage.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PassHash  string    `gorm:"type:varchar(100)"`
	FirstName string    `gorm:"type:varchar(100)"`
	LastName  string    `gorm:"type:varchar(100)"`
	Role      role.Role `gorm:"foreignKey:RoleID"`
	Source    string    `gorm:"type:varchar(255)"`
	Enabled   bool      `gorm:"default:true"`
	RoleID    int
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// Status id the accound status of the user
type Status string

const (
	// Registered is the status for having a registration, but not yet confirmed account
	Registered Status = "registered"
	// Confirmed is the status for having full access to the application
	Confirmed Status = "confirmed"
	// Deactivated is the status for a deleted or unregistered user
	Deactivated Status = "deactivated"
)

// NewUser is a function to create a new `User` instance, hashing the password right off the bat
func NewUser(email string, password string) (*User, error) {
	u := new(User)
	u.Email = email
	u.Source = "credentials"
	u.Status = Registered
	err := u.SetPassword(password)
	return u, err
}

// SetPassword sets the password of the target user.
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassHash = string(hash)
	return nil
}

// VerifyPassword is a method of the `User` struct. It takes a password string as input
// and compares it with the hashed password stored in the `PassHash` field of the `User` struct.
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PassHash), []byte(password))
	print(err)
	return err == nil
}

// IsActive is a method of `User` returning whether the user is enabled, confirmed and can store data
func (u *User) IsActive() bool {
	return u.Enabled && u.Status == Confirmed
}

// AsProfile is a method of the `User` struct. It converts a `User` object into a `Profile` object.
func (u *User) AsProfile() Profile {
	return Profile{
		ID:        u.ID.String(),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role.AsProfileRole(),
		Status:    string(u.Status),
		Source:    u.Source,
	}
}

// AsAdminView is a method of the `User` struct. It converts a `User` object into a `AdminView` object.
func (u *User) AsAdminView() AdminView {
	var deleted *int
	if u.DeletedAt.Valid {
		d := int(u.DeletedAt.Time.Unix())
		deleted = &d
	}
	return AdminView{
		ID:        u.ID.String(),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Status:    string(u.Status),
		Enabled:   u.Enabled,
		Role:      u.Role.AsProfileRole(),
		Created:   int(u.CreatedAt.Unix()),
		Updated:   int(u.UpdatedAt.Unix()),
		Deleted:   deleted,
	}
}

// Credentials is the JSON user representation for logging in with username and password
type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha_token" binding:"required"`
}

// Profile is the JSON user representation for authenticated users
type Profile struct {
	ID        string           `json:"id"`
	Email     string           `json:"email"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Role      role.ProfileRole `json:"role"`
	Status    string           `json:"status"`
	Source    string           `json:"source"`
}

// Registration is the JSON user representation for registration/signup process
type Registration struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password" binding:"required"`
	Captcha   string `json:"captcha_token" binding:"required"`
}

// AdminView is the user representation for the admin view of the application.
type AdminView struct {
	ID        string           `json:"id"`
	Email     string           `json:"email"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Status    string           `json:"status"`
	Role      role.ProfileRole `json:"role"`
	Enabled   bool             `json:"enabled"`
	Created   int              `json:"created"`
	Updated   int              `json:"updated"`
	Deleted   *int             `json:"deleted"`
}

// Patch is the user representation for patching an admin view of the application.
type Patch struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	//	Role      role.ProfileRole `json:"role"` TODO: fix role patching
	Enabled bool `json:"enabled"`
}

// SetEnabled is the user representation for enabling/disabling user authentication.
type SetEnabled struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

// RoleUser is aggregated data on the role, with the user count.
type RoleUser struct {
	Role  string `json:"role"`
	Users int    `json:"users"`
}

// Stats is aggregated data on the storer.
type Stats struct {
	TotalUsers   int
	Distribution []RoleUser
}
