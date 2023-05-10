package company

import (
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
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
	_, role := auth.GetUserEmailAndRole(c)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, response.Response{
			Code:    http.StatusForbidden,
			Message: "only admins can create companies",
		})
	}

	companyDTO := &DTO{}

	if err := c.Bind(companyDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}

	if companyDTO.Name == "" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	_, err := h.service.Create(c, companyDTO)
	if err != nil {
		return c.JSON(http.StatusConflict, response.Response{
			Code:    http.StatusConflict,
			Message: CompanyAlreadyExistsErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "company was successfully created",
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

	companyDTO, err := h.service.Read(c, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"company": companyDTO,
	})
}

func (h *Handler) ReadAll(c echo.Context) error {
	companyDTOs, err := h.service.ReadAll(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":      http.StatusOK,
		"companies": companyDTOs,
	})
}

func (h *Handler) Update(c echo.Context) error {
	_, role := auth.GetUserEmailAndRole(c)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, response.Response{
			Code:    http.StatusForbidden,
			Message: "only admins can update companies",
		})
	}

	companyDTO := &DTO{}

	if err := c.Bind(companyDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}
	if companyDTO.ID <= 0 || companyDTO.Name == "" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	isUpdated, err := h.service.Update(c, companyDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if !isUpdated {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: CompanyNotFoundErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "company was successfully updated",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	_, role := auth.GetUserEmailAndRole(c)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, response.Response{
			Code:    http.StatusForbidden,
			Message: "only admins can delete companies",
		})
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
			Message: CompanyNotFoundErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "company was successfully deleted",
	})
}
