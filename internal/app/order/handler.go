package order

import (
	"net/http"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
	"github.com/labstack/echo/v4"
)

type Service interface {
	ReadCurrentUserArrangingOrderLazy(c echo.Context, userID uint64) (*DTO, error)
	ReadCurrentUserArrangingOrderEager(c echo.Context, userID uint64) (*DTO, error)
	Create(c echo.Context, userID uint64) (uint64, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c echo.Context) error {
	userData, err := auth.GetUserDataAndCheckRole(c, "user")

	if err != nil {
		return err
	}

	_, err = h.service.Create(c, userData.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error creating order for user",
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "order was successfully created",
	})
}

func (h *Handler) ReadCurrentUserArrangingOrder(c echo.Context) error {
	userData, err := auth.GetUserDataAndCheckRole(c, "user")

	if err != nil {
		return err
	}

	fetch := c.QueryParam("fetch")
	var orderDTO *DTO

	if fetch == "eager" {
		orderDTO, err = h.service.ReadCurrentUserArrangingOrderEager(c, userData.ID)
	} else {
		orderDTO, err = h.service.ReadCurrentUserArrangingOrderLazy(c, userData.ID)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"order": orderDTO,
	})
}
