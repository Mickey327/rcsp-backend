package cart

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/order"
	"github.com/Mickey327/rcsp-backend/internal/app/orderItem"
	"github.com/Mickey327/rcsp-backend/internal/app/product"
	"github.com/labstack/echo/v4"
)

type Service interface {
	UpdateCart(c echo.Context, dto *orderItem.DTO) (*order.DTO, error)
	RemoveFromCart(c echo.Context, dto *orderItem.DTO) (*order.DTO, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) UpdateCart(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "user")

	if err != nil {
		return err
	}

	var productID uint64
	var orderID uint64

	productIDString := c.QueryParam("productID")
	if productIDString != "" {
		productID, err = strconv.ParseUint(productIDString, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id товара")
		}
		if productID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "id товара должно быть положительным")
		}
	}

	orderIDString := c.QueryParam("orderID")
	if orderIDString != "" {
		orderID, err = strconv.ParseUint(orderIDString, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id заказа")
		}

		if orderID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "id заказа должно быть положительным")
		}
	}

	dto := &orderItem.DTO{}

	if err = c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка привязки данных из json")
	}

	dto.OrderID = orderID
	dto.Product = &product.DTO{
		ID: productID,
	}

	o, err := h.service.UpdateCart(c, dto)

	if err != nil {
		if errors.Is(err, WrongCartErr) || errors.Is(err, NotPositiveQuantityErr) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		} else if errors.Is(err, order.OrderNotFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "ошибка произошла во время обновления корзины")
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"order": o,
	})
}

func (h *Handler) RemoveFromCart(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "user")

	if err != nil {
		return err
	}

	var productID uint64
	var orderID uint64

	productIDString := c.QueryParam("productID")
	if productIDString != "" {
		productID, err = strconv.ParseUint(productIDString, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id товара")
		}
		if productID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "id товара должно быть положительным")
		}
	}

	orderIDString := c.QueryParam("orderID")
	if orderIDString != "" {
		orderID, err = strconv.ParseUint(orderIDString, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id заказа")
		}

		if orderID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "id заказа должно быть положительным")
		}
	}

	dto := &orderItem.DTO{}

	dto.OrderID = orderID
	dto.Product = &product.DTO{
		ID: productID,
	}

	o, err := h.service.RemoveFromCart(c, dto)
	if err != nil {
		if errors.Is(err, orderItem.OrderItemNotFound) || errors.Is(err, order.OrderNotFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, WrongCartErr) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"order": o,
	})
}
