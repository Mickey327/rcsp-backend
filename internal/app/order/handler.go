package order

import (
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
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
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка создания заказа для пользователя")
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
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if databaseOrderDTO.Total == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "пользователь не может обновить пустой заказ")
	}

	if databaseOrderDTO.UserID != userData.ID && userData.Role == "user" {
		return echo.NewHTTPError(http.StatusForbidden, "пользователь не может обновить чужой заказ")
	}

	orderDTO := DTO{}

	if err = c.Bind(&orderDTO); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка привязки данных из json")
	}
	if orderDTO.Status == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	if userData.Role == "admin" && orderDTO.ID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	databaseOrderDTO.Status = orderDTO.Status
	databaseOrderDTO.IsArranged = true

	isUpdated, err := h.service.Update(c, databaseOrderDTO)
	if !isUpdated || err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка обновления заказа")
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
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id заказа")
	}
	if id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id заказа должно быть положительным")
	}

	order, err := h.service.ReadByIdEager(c, id)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, OrderNotFoundErr.Error())
	}

	if order.UserID != userData.ID && userData.Role == "user" {
		return echo.NewHTTPError(http.StatusForbidden, "пользователь не может получить чужой заказ")
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
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"order": orderDTO,
	})
}
