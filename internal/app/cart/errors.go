package cart

import "errors"

var (
	WrongCartErr           = errors.New("user can't change other user's cart")
	NotPositiveQuantityErr = errors.New("user can't make quantity below 1")
)
