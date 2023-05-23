package category

import (
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Create(c echo.Context, categoryDTO *DTO) (uint64, error)
	Read(c echo.Context, id uint64) (*DTO, error)
	ReadAll(c echo.Context) ([]*DTO, error)
	Update(c echo.Context, categoryDTO *DTO) (bool, error)
	Delete(c echo.Context, id uint64) (bool, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "admin")
	if err != nil {
		return err
	}

	categoryDTO := DTO{}

	if err = c.Bind(&categoryDTO); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка привязки данных из json")
	}

	if categoryDTO.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	_, err = h.service.Create(c, &categoryDTO)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, CategoryAlreadyExistsErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "категория была успешно создана",
	})
}

func (h *Handler) Read(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id категории")
	}
	if id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id категории должно быть положительным")
	}

	categoryDTO, err := h.service.Read(c, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":     http.StatusOK,
		"category": categoryDTO,
	})
}

func (h *Handler) ReadAll(c echo.Context) error {
	categoryDTOs, err := h.service.ReadAll(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":       http.StatusOK,
		"categories": categoryDTOs,
	})
}

func (h *Handler) Update(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "admin")

	if err != nil {
		return err
	}

	categoryDTO := DTO{}

	if err = c.Bind(&categoryDTO); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка привязки данных из json")
	}
	if categoryDTO.ID <= 0 || categoryDTO.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	isUpdated, err := h.service.Update(c, &categoryDTO)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !isUpdated {
		return echo.NewHTTPError(http.StatusNotFound, CategoryNotFoundErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "категория была успешно обновлена",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "admin")

	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id категории")
	}
	if id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id категории должно быть положительным")
	}

	isDeleted, err := h.service.Delete(c, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !isDeleted {
		return echo.NewHTTPError(http.StatusNotFound, CategoryNotFoundErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "категория была успешно удалена",
	})
}
