package auth

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is the user representation for database storage
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PassHash  string    `gorm:"type:varchar(100)"`
	Phone     string    `gorm:"type:varchar(20)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// NewUser is a function to create a new `User` instance, hashing the password right off the bat
func NewUser(email string, password string, phone string) (*User, error) {
	u := new(User)
	u.Email = email
	u.Phone = phone
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.PassHash = string(hash)
	return u, nil
}

// VerifyPassword is a method of the `User` struct. It takes a password string as input
// and compares it with the hashed password stored in the `PassHash` field of the `User` struct.
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PassHash), []byte(password))
	print(err)
	return err == nil
}

// AsProfile is a method of the `User` struct. It converts a `User` object into a `Profile` object.
func (u *User) AsProfile() Profile {
	return Profile{
		ID:    u.ID.String(),
		Email: u.Email,
		Phone: u.Phone,
	}
}

// Credentials is the JSON user representation for logging in with username and password
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Profile is the JSON user representation for authenticated users
type Profile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Registration is the JSON user representation for registration/signup process
type Registration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}
