package errors

import "fmt"

var (
	ErrWordAlreadyExists = fmt.Errorf("word already exists")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrInternal          = fmt.Errorf("internal error")
	ErrUserNotFound      = fmt.Errorf("user not found")
)
