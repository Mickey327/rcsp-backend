package user

import "errors"

var (
	UserNotFoundErr      = errors.New("no user with such email")
	UserAlreadyExistsErr = errors.New("user with such email already exists")
	UserWrongPasswordErr = errors.New("wrong password")
	UserTokenErr         = errors.New("error generating token for user")
)
