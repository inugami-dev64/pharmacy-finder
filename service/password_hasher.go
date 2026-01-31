package service

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	// Creates a new password hash from
	// provided password and returns it
	//
	// Possible error values are:
	//   - ErrPasswordTooLong
	CreatePasswordHash(password string) (string, error)

	// Verifies if given password matches the
	// provided password hash and if it does,
	// returns true
	Verify(password string, hash string) bool
}

func ProvidePasswordHasher() PasswordHasher {
	return BCryptPasswordHasher{}
}

type BCryptPasswordHasher struct{}

func (hasher BCryptPasswordHasher) CreatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == bcrypt.ErrPasswordTooLong {
		return "", ErrPasswordTooLong
	} else if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (hasher BCryptPasswordHasher) Verify(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
