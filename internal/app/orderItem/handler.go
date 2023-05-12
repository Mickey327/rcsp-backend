package orderItem

import (
	"net/http"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
	"github.com/labstack/echo/v4"
)

type Service interface {
	ChangeOrderItemQuantity(c echo.Context, orderItemDTO *DTO) (bool, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ChangeOrderItemQuantity(c echo.Context) error {
	token, err := auth.GetUserToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Message: "can't get jwt token from cookie",
		})
	}

	userData := auth.GetUserDataFromToken(token)
	if userData.Role != "user" {
		return c.JSON(http.StatusForbidden, response.Response{
			Code:    http.StatusForbidden,
			Message: "only users can change cart products quantities",
		})
	}

	orderItemDTO := &DTO{}

	if err = c.Bind(orderItemDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}
	if orderItemDTO.OrderID <= 0 || orderItemDTO.ProductID <= 0 || orderItemDTO.QuantityDelta == 0 {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	isUpdated, err := h.service.ChangeOrderItemQuantity(c, orderItemDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if !isUpdated {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: "error updating cart product's quantity",
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "cart product's quantity was successfully updated",
	})
}
