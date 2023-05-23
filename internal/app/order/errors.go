package order

import "errors"

var (
	OrderNotFoundErr = errors.New("заказ не найден")
	OrderEmptyErr    = errors.New("заказ пустой")
)
