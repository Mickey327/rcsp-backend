package user

import "errors"

var (
	UserNotFoundErr      = errors.New("пользователя с таким email не существует")
	UserAlreadyExistsErr = errors.New("пользователь с таким email уже существует")
	UserWrongPasswordErr = errors.New("неверный пароль")
	UserTokenErr         = errors.New("ошибка генерации токена для пользователя")
)
