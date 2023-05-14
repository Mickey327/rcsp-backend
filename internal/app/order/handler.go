package order

import (
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
	"github.com/labstack/echo/v4"
)

type Service interface {
	ReadCurrentUserArrangingOrderLazy(c echo.Context, userID uint64) (*DTO, error)
	ReadCurrentUserArrangingOrderEager(c echo.Context, userID uint64) (*DTO, error)
	Create(c echo.Context, userID uint64) (*DTO, error)
	ReadByIdEager(c echo.Context, id uint64) (*DTO, error)
	Update(c echo.Context, dto *DTO) (bool, error)
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

	o, err := h.service.Create(c, userData.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error creating order for user",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"order": o,
	})
}

func (h *Handler) Update(c echo.Context) error {
	userData, err := auth.GetUserDataAndCheckRole(c, "user", "admin")

	if err != nil {
		return err
	}

	databaseOrderDTO, err := h.service.ReadCurrentUserArrangingOrderLazy(c, userData.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	if databaseOrderDTO.Total == 0 {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "user can't update zero total order",
		})
	}

	if databaseOrderDTO.UserID != userData.ID && userData.Role == "user" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "user can't update other user's order",
		})
	}

	orderDTO := DTO{}

	if err = c.Bind(&orderDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}
	if orderDTO.Status == "" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	if userData.Role == "admin" && orderDTO.ID == 0 {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	databaseOrderDTO.Status = orderDTO.Status
	databaseOrderDTO.IsArranged = true

	isUpdated, err := h.service.Update(c, databaseOrderDTO)
	if !isUpdated || err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error during update order",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"order": databaseOrderDTO,
	})
}

func (h *Handler) ReadByIdEager(c echo.Context) error {
	userData, err := auth.GetUserDataAndCheckRole(c, "user", "admin")

	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "error parsing id path parameter",
		})
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "id value must be positive",
		})
	}

	order, err := h.service.ReadByIdEager(c, id)

	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: "order not found",
		})
	}

	if order.UserID != userData.ID && userData.Role == "user" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "user can't get other user's order",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"order": order,
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
