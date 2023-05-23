package company

import (
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Create(c echo.Context, companyDTO *DTO) (uint64, error)
	Read(c echo.Context, id uint64) (*DTO, error)
	ReadAll(c echo.Context) ([]*DTO, error)
	Update(c echo.Context, companyDTO *DTO) (bool, error)
	Delete(c echo.Context, id uint64) (bool, error)
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
	_, err := auth.GetUserDataAndCheckRole(c, "admin")

	if err != nil {
		return err
	}

	companyDTO := DTO{}

	if err = c.Bind(&companyDTO); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка привязки данных из json")
	}

	if companyDTO.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	_, err = h.service.Create(c, &companyDTO)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, CompanyAlreadyExistsErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "компания была успешно создана",
	})
}

func (h *Handler) Read(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id компании")
	}
	if id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id компании должно быть положительным")
	}

	companyDTO, err := h.service.Read(c, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"company": companyDTO,
	})
}

func (h *Handler) ReadAll(c echo.Context) error {
	companyDTOs, err := h.service.ReadAll(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":      http.StatusOK,
		"companies": companyDTOs,
	})
}

func (h *Handler) Update(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "admin")

	if err != nil {
		return err
	}

	companyDTO := DTO{}

	if err = c.Bind(&companyDTO); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка парсинга id категории")
	}
	if companyDTO.ID <= 0 || companyDTO.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	isUpdated, err := h.service.Update(c, &companyDTO)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !isUpdated {
		return echo.NewHTTPError(http.StatusNotFound, CompanyNotFoundErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "компания была успешно обновлена",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "admin")

	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id компании")
	}
	if id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id компании должно быть положительным")
	}

	isDeleted, err := h.service.Delete(c, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !isDeleted {
		return echo.NewHTTPError(http.StatusNotFound, CompanyNotFoundErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "компания была успешно удалена",
	})
}
