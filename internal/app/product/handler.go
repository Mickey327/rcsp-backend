package product

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/category"
	"github.com/Mickey327/rcsp-backend/internal/app/company"
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
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга цены из формы")
	}

	companyID, err := strconv.ParseUint(c.FormValue("companyID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга компании из формы")
	}

	categoryID, err := strconv.ParseUint(c.FormValue("categoryID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга категории из формы")
	}

	stock, err := strconv.ParseUint(c.FormValue("stock"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга количества из формы")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга изображения из формы")
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
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	_, err = h.service.Create(c, &productDTO, file)
	if err != nil {
		if errors.Is(err, ProductAlreadyExistsErr) {
			return echo.NewHTTPError(http.StatusConflict, ProductAlreadyExistsErr.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "ошибка создания товара")
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "товар был успешно создан",
	})
}

func (h *Handler) Read(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id товара")
	}
	if id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id товара должно быть положительным")
	}

	var productDTO *DTO

	fetchType := c.QueryParam("fetch")

	if fetchType == "eager" {
		productDTO, err = h.service.ReadEager(c, id)
	} else {
		productDTO, err = h.service.Read(c, id)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
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
			return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id компании")
		}

		if companyID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "id компании должно быть положительным")
		}
	}

	if categoryIDString != "" {
		categoryID, err = strconv.ParseUint(categoryIDString, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id категории")
		}
		if categoryID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "id категории должно быть положительным")
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
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
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
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка привязки данных из json")
	}

	if productDTO.Company == nil || productDTO.Category == nil || productDTO.ID <= 0 || productDTO.Name == "" || productDTO.Price <= 0 || productDTO.Image == "" || productDTO.Company.ID <= 0 || productDTO.Category.ID <= 0 || productDTO.Stock < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	isUpdated, err := h.service.Update(c, &productDTO)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !isUpdated {
		return echo.NewHTTPError(http.StatusNotFound, ProductNotFoundErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "товар был успешно обновлен",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	_, err := auth.GetUserDataAndCheckRole(c, "admin")

	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id товара")
	}
	if id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "id товара должно быть положительным")
	}

	isDeleted, err := h.service.Delete(c, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !isDeleted {
		return echo.NewHTTPError(http.StatusNotFound, ProductNotFoundErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "товар был успешно удален",
	})
}
