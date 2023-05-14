package comment

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
	"github.com/labstack/echo/v4"
)

type Service interface {
	WriteComment(c echo.Context, userID, productID uint64, message string) (uint64, error)
	ReadByProductID(c echo.Context, productID uint64) ([]*DTO, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) WriteComment(c echo.Context) error {
	userData, err := auth.GetUserDataAndCheckRole(c, "user")

	if err != nil {
		return err
	}

	var productID uint64

	message := c.FormValue("message")

	if message == "" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "comment message must be not empty",
		})
	}

	productIDString := c.QueryParam("productID")
	if productIDString != "" {
		productID, err = strconv.ParseUint(productIDString, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "error parsing productID query parameter",
			})
		}
		if productID <= 0 {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "product id value must be positive",
			})
		}
	}

	_, err = h.service.WriteComment(c, userData.ID, productID, message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error creating user message",
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "comment was successfully created",
	})

}

func (h *Handler) ReadComments(c echo.Context) error {
	var productID uint64
	var err error
	productIDString := c.QueryParam("productID")
	if productIDString != "" {
		productID, err = strconv.ParseUint(productIDString, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "error parsing productID query parameter",
			})
		}
		if productID <= 0 {
			return c.JSON(http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "product id value must be positive",
			})
		}
	}

	comments, err := h.service.ReadByProductID(c, productID)
	if err != nil {
		if errors.Is(err, CommentNotFoundErr) {
			return c.JSON(http.StatusNotFound, response.Response{
				Code:    http.StatusNotFound,
				Message: "comments for this product not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, response.Response{
				Code:    http.StatusInternalServerError,
				Message: "error getting product comments",
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":     http.StatusOK,
		"comments": comments,
	})

}
