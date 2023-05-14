package cart

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/order"
	"github.com/Mickey327/rcsp-backend/internal/app/orderItem"
	"github.com/Mickey327/rcsp-backend/internal/app/product"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
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
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "error parsing productID query parameter",
			})
		}
		if productID <= 0 {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "product id value must be positive",
			})
		}
	}

	orderIDString := c.QueryParam("orderID")
	if orderIDString != "" {
		orderID, err = strconv.ParseUint(orderIDString, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "error parsing orderID query parameter",
			})
		}

		if orderID <= 0 {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "order id value must be positive",
			})
		}
	}

	dto := &orderItem.DTO{}

	if err = c.Bind(&dto); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}

	dto.OrderID = orderID
	dto.Product = &product.DTO{
		ID: productID,
	}

	o, err := h.service.UpdateCart(c, dto)

	if err != nil {
		if errors.Is(err, WrongCartErr) || errors.Is(err, NotPositiveQuantityErr) {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		} else if errors.Is(err, order.OrderNotFoundErr) {
			return c.JSON(http.StatusNotFound, response.Response{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, response.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
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
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "error parsing productID query parameter",
			})
		}
		if productID <= 0 {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "product id value must be positive",
			})
		}
	}

	orderIDString := c.QueryParam("orderID")
	if orderIDString != "" {
		orderID, err = strconv.ParseUint(orderIDString, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "error parsing orderID query parameter",
			})
		}

		if orderID <= 0 {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "order id value must be positive",
			})
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
			return c.JSON(http.StatusNotFound, response.Response{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
		}
		if errors.Is(err, WrongCartErr) {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, response.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"order": o,
	})
}
