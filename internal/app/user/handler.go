package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Register(c echo.Context, userDTO *DTO) error
	Login(c echo.Context, userDTO *DTO) (string, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c echo.Context) error {
	userDTO := &DTO{}

	if err := c.Bind(userDTO); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка привязки данных из json")
	}

	if err := c.Validate(userDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "данные представлены в неверном формате")
	}

	if err := h.service.Register(c, userDTO); err != nil {
		return echo.NewHTTPError(http.StatusConflict, UserAlreadyExistsErr.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "пользователь успешно зарегистрирован",
	})
}

func (h *Handler) Login(c echo.Context) error {
	userDTO := &DTO{}

	if err := c.Bind(userDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка привязки данных из json")
	}

	token, err := h.service.Login(c, userDTO)
	if err != nil {
		if errors.Is(err, UserNotFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		if errors.Is(err, UserWrongPasswordErr) || errors.Is(err, UserTokenErr) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(auth.GetJWTSecret().ExpirationTimeInHours),
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "пользователь успешно залогинился",
		"token":   token,
	})
}

func (h *Handler) GetAuthenticatedUser(c echo.Context) error {
	token, err := auth.GetUserToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "невозможно получить jwt token из cookie")
	}

	userData := auth.GetUserDataFromToken(token)

	if token.Raw == "" || userData.ID <= 0 || userData.Email == "" || userData.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный jwt token пользователя")
	}

	userDTO := &DTO{ID: userData.ID, Email: userData.Email, Role: userData.Role}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"user":  userDTO,
		"token": token.Raw,
	})
}

func (h *Handler) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "пользователь успешно вышел",
	})
}
