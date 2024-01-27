package auth

import (
	"fmt"
	"math"
	"time"

	"github.com/inokone/photostorage/auth/account"
	"github.com/inokone/photostorage/auth/user"
	"github.com/rs/zerolog/log"
)

// Service is a worker for authentication and authorization.
type Service struct {
	users    user.Storer
	accounts account.Storer
	jwt      JWTHandler
}

// NewService creates a new `Service`, based on the user and account persistence.
func NewService(users user.Storer, auths account.Storer, jwt JWTHandler) *Service {
	return &Service{
		users:    users,
		accounts: auths,
		jwt:      jwt,
	}
}

// InvalidCredentials is an error for a bad email address
type InvalidCredentials string

// Error is the string representation of an `CredentialsError`
func (e InvalidCredentials) Error() string { return string(e) }

// LockedUser is an error for a user locked out of the system
type LockedUser struct {
	seconds int64
}

// Error is the string representation of an `LockedUserError`
func (e LockedUser) Error() string { return fmt.Sprintf("%v", e.seconds) }

// ValidateCredentials validates the user credentials sets and clears retry timeout for failed creds
func (s Service) ValidateCredentials(usr *user.User, password string) error {
	var (
		secs     int64
		err      error
		verified bool
	)
	secs, err = s.checkTimeout(usr)
	if err != nil {
		log.Err(err).Str("UserID", usr.ID.String()).Msg("Failed to collect login timeout.")
		return InvalidCredentials("")
	}
	if secs > 0 {
		return LockedUser{
			seconds: secs,
		}
	}

	verified = usr.VerifyPassword(password)
	if !verified {
		err = s.increaseTimeout(usr)
		if err != nil {
			log.Warn().Str("user", usr.ID.String()).Msg("Failed to increase timeout for user")
		}
		return InvalidCredentials("")
	}
	s.clearTimeout(usr)
	return nil
}

func (s Service) checkTimeout(usr *user.User) (int64, error) {
	var (
		a   account.Account
		err error
	)
	a, err = s.accounts.ByUser(usr.ID)
	if err != nil {
		return 0, err
	}
	if a.FailedLoginLock.After(time.Now()) {
		return a.FailedLoginLock.Unix() - time.Now().Unix(), nil
	}
	return 0, nil
}

func (s Service) increaseTimeout(usr *user.User) error {
	var (
		a       account.Account
		err     error
		timeout int
	)
	a, err = s.accounts.ByUser(usr.ID)
	if err != nil {
		return err
	}
	a.FailedLoginCounter++
	a.LastFailedLogin = time.Now()
	if a.FailedLoginCounter > 2 {
		timeout = int(math.Pow(10, float64(a.FailedLoginCounter-2))) // exponential backoff - 10 sec, 10 sec, 1000 sec, ...
		a.FailedLoginLock = time.Now().Add(time.Second * time.Duration(timeout))
	}
	return s.accounts.Update(&a)
}

func (s Service) clearTimeout(usr *user.User) error {
	var (
		a   account.Account
		err error
	)
	a, err = s.accounts.ByUser(usr.ID)
	if err != nil {
		return err
	}
	a.FailedLoginCounter = 0
	a.FailedLoginLock = time.Now()
	return s.accounts.Update(&a)
}
