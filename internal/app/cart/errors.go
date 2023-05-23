package cart

import "errors"

var (
	WrongCartErr           = errors.New("пользователь не может изменять чужую корзину")
	NotPositiveQuantityErr = errors.New("пользователь не может сделать количество позиции менее 1")
)
