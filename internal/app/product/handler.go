package product

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/category"
	"github.com/Mickey327/rcsp-backend/internal/app/company"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Create(c echo.Context, productDTO *DTO, image *multipart.FileHeader) (uint64, error)
	Read(c echo.Context, id uint64) (*DTO, error)
	ReadEager(c echo.Context, id uint64) (*DTO, error)
	ReadAll(c echo.Context) ([]*DTO, error)
	ReadByCategoryID(c echo.Context, categoryID uint64) ([]*DTO, error)
	ReadByCompanyID(c echo.Context, companyID uint64) ([]*DTO, error)
	ReadByCompanyIDAndCategoryID(c echo.Context, companyID, categoryID uint64) ([]*DTO, error)
	Update(c echo.Context, productDTO *DTO) (bool, error)
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

	name := c.FormValue("name")
	description := c.FormValue("description")

	price, err := strconv.ParseUint(c.FormValue("price"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "error parsing price value from form",
		})
	}

	companyID, err := strconv.ParseUint(c.FormValue("companyID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "error parsing company value from form",
		})
	}

	categoryID, err := strconv.ParseUint(c.FormValue("categoryID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "error parsing category value from form",
		})
	}

	stock, err := strconv.ParseUint(c.FormValue("stock"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "error parsing stock value from form",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "error getting image from form",
		})
	}

	comp := &company.DTO{
		ID: companyID,
	}
	cat := &category.DTO{
		ID: categoryID,
	}

	productDTO := DTO{
		Name:        name,
		Description: description,
		Price:       price,
		Company:     comp,
		Category:    cat,
		Stock:       stock,
		Image:       file.Filename,
	}

	if productDTO.Name == "" || productDTO.Price <= 0 || productDTO.Company.ID <= 0 || productDTO.Category.ID <= 0 || productDTO.Stock < 0 {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	_, err = h.service.Create(c, &productDTO, file)
	if err != nil {
		if errors.Is(err, ProductAlreadyExistsErr) {
			return c.JSON(http.StatusConflict, response.Response{
				Code:    http.StatusConflict,
				Message: ProductAlreadyExistsErr.Error(),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, response.Response{
				Code:    http.StatusInternalServerError,
				Message: "error uploading product",
			})
		}
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "product was successfully created",
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

	var productDTO *DTO

	fetchType := c.QueryParam("fetch")

	if fetchType == "eager" {
		productDTO, err = h.service.ReadEager(c, id)
	} else {
		productDTO, err = h.service.Read(c, id)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"product": productDTO,
	})
}

func (h *Handler) ReadAll(c echo.Context) error {
	categoryIDString := c.QueryParam("categoryID")
	companyIDString := c.QueryParam("companyID")
	var categoryID, companyID uint64
	var productDTOs []*DTO
	var err error

	if companyIDString != "" {
		companyID, err = strconv.ParseUint(companyIDString, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "error parsing companyID query parameter",
			})
		}

		if companyID <= 0 {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "company id value must be positive",
			})
		}
	}

	if categoryIDString != "" {
		categoryID, err = strconv.ParseUint(categoryIDString, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "error parsing categoryID query parameter",
			})
		}
		if categoryID <= 0 {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "category id value must be positive",
			})
		}
	}

	if categoryID > 0 && companyID > 0 {
		productDTOs, err = h.service.ReadByCompanyIDAndCategoryID(c, companyID, categoryID)
	} else if categoryID > 0 {
		productDTOs, err = h.service.ReadByCategoryID(c, categoryID)
	} else if companyID > 0 {
		productDTOs, err = h.service.ReadByCompanyID(c, companyID)
	} else {
		productDTOs, err = h.service.ReadAll(c)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":     http.StatusOK,
		"products": productDTOs,
	})
}

func (h *Handler) Update(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "admin")

	if err != nil {
		return err
	}

	productDTO := DTO{}

	if err = c.Bind(&productDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}

	if productDTO.Company == nil || productDTO.Category == nil || productDTO.ID <= 0 || productDTO.Name == "" || productDTO.Price <= 0 || productDTO.Image == "" || productDTO.Company.ID <= 0 || productDTO.Category.ID <= 0 || productDTO.Stock < 0 {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong values format provided",
		})
	}

	isUpdated, err := h.service.Update(c, &productDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if !isUpdated {
		return c.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Message: ProductNotFoundErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "product was successfully updated",
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
			Message: ProductNotFoundErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "product was successfully deleted",
	})
}
