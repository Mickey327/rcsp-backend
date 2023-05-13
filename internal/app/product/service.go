package product

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Create(ctx context.Context, product *Product) (uint64, error)
	Read(ctx context.Context, id uint64) (*Product, error)
	ReadEager(ctx context.Context, id uint64) (*Product, error)
	ReadAll(ctx context.Context) ([]*Product, error)
	ReadByCategoryID(ctx context.Context, categoryID uint64) ([]*Product, error)
	ReadByCompanyID(ctx context.Context, companyID uint64) ([]*Product, error)
	ReadByCompanyIDAndCategoryID(ctx context.Context, companyID, categoryID uint64) ([]*Product, error)
	Update(ctx context.Context, product *Product) (bool, error)
	Delete(ctx context.Context, id uint64) (bool, error)
}

type ProductService struct {
	repository Repository
}

func NewService(repository Repository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) Create(c echo.Context, productDTO *DTO, file *multipart.FileHeader) (uint64, error) {
	src, err := file.Open()
	if err != nil {
		log.Println("file open", err)
		return 0, err
	}
	defer src.Close()

	path, _ := os.Getwd()
	path += "/static/" + file.Filename
	// Destination
	dst, err := os.Create(path)
	if err != nil {
		log.Println(os.Getwd())
		log.Println("Os create static:", err)
		return 0, err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println("io copy dst", err)
		return 0, err
	}

	id, err := s.repository.Create(c.Request().Context(), productDTO.ToProduct())

	if err != nil {
		return 0, ProductAlreadyExistsErr
	}

	return id, nil
}

func (s *ProductService) Read(c echo.Context, id uint64) (*DTO, error) {
	product, err := s.repository.Read(c.Request().Context(), id)

	if err != nil {
		return nil, err
	}

	return product.ToDTO(), nil
}

func (s *ProductService) ReadEager(c echo.Context, id uint64) (*DTO, error) {
	product, err := s.repository.ReadEager(c.Request().Context(), id)

	if err != nil {
		return nil, err
	}

	return product.ToDTO(), nil
}

func (s *ProductService) ReadAll(c echo.Context) ([]*DTO, error) {
	products, err := s.repository.ReadAll(c.Request().Context())

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, ProductNotFoundErr
	}

	return ToDTOs(products), nil
}

func (s *ProductService) ReadByCategoryID(c echo.Context, categoryID uint64) ([]*DTO, error) {
	products, err := s.repository.ReadByCategoryID(c.Request().Context(), categoryID)

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, ProductNotFoundErr
	}

	return ToDTOs(products), nil
}

func (s *ProductService) ReadByCompanyID(c echo.Context, companyID uint64) ([]*DTO, error) {
	products, err := s.repository.ReadByCompanyID(c.Request().Context(), companyID)

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, ProductNotFoundErr
	}

	return ToDTOs(products), nil
}

func (s *ProductService) ReadByCompanyIDAndCategoryID(c echo.Context, companyID, categoryID uint64) ([]*DTO, error) {
	products, err := s.repository.ReadByCompanyIDAndCategoryID(c.Request().Context(), companyID, categoryID)

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, ProductNotFoundErr
	}

	return ToDTOs(products), nil
}

func (s *ProductService) Update(c echo.Context, productDTO *DTO) (bool, error) {
	isUpdated, err := s.repository.Update(c.Request().Context(), productDTO.ToProduct())

	if err != nil {
		return false, err
	}

	return isUpdated, nil
}

func (s *ProductService) Delete(c echo.Context, id uint64) (bool, error) {
	isDeleted, err := s.repository.Delete(c.Request().Context(), id)

	if err != nil {
		return false, err
	}

	return isDeleted, nil
}
