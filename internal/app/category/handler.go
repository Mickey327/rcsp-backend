package category

import (
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
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
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}

	if categoryDTO.Name == "" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	_, err = h.service.Create(c, &categoryDTO)
	if err != nil {
		return c.JSON(http.StatusConflict, response.Response{
			Code:    http.StatusConflict,
			Message: CategoryAlreadyExistsErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "category was successfully created",
	})
}

func (h *Handler) Read(c echo.Context) error {
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

	categoryDTO, err := h.service.Read(c, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":     http.StatusOK,
		"category": categoryDTO,
	})
}

func (h *Handler) ReadAll(c echo.Context) error {
	categoryDTOs, err := h.service.ReadAll(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
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
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}
	if categoryDTO.ID <= 0 || categoryDTO.Name == "" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	isUpdated, err := h.service.Update(c, &categoryDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if !isUpdated {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: CategoryNotFoundErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "category was successfully updated",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "admin")

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

	isDeleted, err := h.service.Delete(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if !isDeleted {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: CategoryNotFoundErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "category was successfully deleted",
	})
}
