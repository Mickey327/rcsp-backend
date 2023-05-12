package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/response"
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
		return c.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Message: "error binding json data",
		})
	}

	if userDTO.Email == "" || userDTO.Password == "" {
		return c.JSON(
			http.StatusBadRequest, response.Response{
				Code:    http.StatusBadRequest,
				Message: "wrong email or password format provided",
			})
	}

	if err := h.service.Register(c, userDTO); err != nil {
		return c.JSON(
			http.StatusInternalServerError, response.Response{
				Code:    http.StatusInternalServerError,
				Message: "error happened during user registration",
			})
	}

	return c.JSON(
		http.StatusOK, response.Response{
			Code:    http.StatusOK,
			Message: "user was successfully registered",
		},
	)
}

func (h *Handler) Login(c echo.Context) error {
	userDTO := &DTO{}

	if err := c.Bind(userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong request body",
		})
	}

	token, err := h.service.Login(c, userDTO)
	if err != nil {
		if errors.Is(err, UserNotFoundErr) {
			return c.JSON(http.StatusNotFound, response.Response{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
		}

		if errors.Is(err, UserWrongPasswordErr) {
			return c.JSON(http.StatusUnauthorized, response.Response{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
		}

		if errors.Is(err, UserTokenErr) {
			return c.JSON(http.StatusUnauthorized, response.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(auth.GetJWTSecret().ExpirationTimeInHours),
		HttpOnly: true,
		Secure:   true,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"message": "user was successfully logged in",
		"token":   token,
	})
}

func (h *Handler) GetAuthenticatedUser(c echo.Context) error {
	token, err := auth.GetUserToken(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "can't get jwt token from cookie",
		})
	}

	userData := auth.GetUserDataFromToken(token)

	if token.Raw == "" || userData.ID <= 0 || userData.Email == "" || userData.Role == "" {
		return c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong user jwt token",
		})
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

	return c.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Message: "user successfully logout",
	})
}
