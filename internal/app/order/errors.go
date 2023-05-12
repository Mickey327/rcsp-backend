package order

import "errors"

var (
	OrderNotFoundErr = errors.New("order not found")
	OrderEmptyErr    = errors.New("order is empty")
)
