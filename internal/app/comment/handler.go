package comment

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
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
		return echo.NewHTTPError(http.StatusBadRequest, "комментарий не может быть пустой")
	}

	productIDString := c.QueryParam("productID")
	if productIDString != "" {
		productID, err = strconv.ParseUint(productIDString, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id товара")
		}
		if productID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "id товара должно быть положительным")
		}
	}

	_, err = h.service.WriteComment(c, userData.ID, productID, message)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка создания комментария")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "комментарий был успешно выложен",
	})

}

func (h *Handler) ReadComments(c echo.Context) error {
	var productID uint64
	var err error
	productIDString := c.QueryParam("productID")
	if productIDString != "" {
		productID, err = strconv.ParseUint(productIDString, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга id товара")
		}
		if productID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "id товара должно быть положительным")
		}
	}

	comments, err := h.service.ReadByProductID(c, productID)
	if err != nil {
		if errors.Is(err, CommentNotFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, "комментарии к данному товару не найдены")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "ошибка получения комментарий к товару")
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":     http.StatusOK,
		"comments": comments,
	})

}
