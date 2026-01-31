package service

import "fmt"

var (
	ErrPasswordTooLong error = fmt.Errorf("password is too long")
)
