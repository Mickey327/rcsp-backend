package user

import (
	"context"
	"log"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Register(ctx context.Context, user *User) (uint64, error)
	GetByEmail(ctx context.Context, user *User) (*User, error)
}

type UserService struct {
	repository Repository
}

func NewService(repository Repository) *UserService {
	return &UserService{repository: repository}
}

// Register creates new user (with 'user' role)
func (u *UserService) Register(c echo.Context, userDTO *DTO) error {
	userDTO.Role = "user" //by default user has 'user' role.

	password, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 10)
	if err != nil {
		log.Printf("error during password encrypt: %v", err)
		return err
	}
	userDTO.Password = string(password)

	id, err := u.repository.Register(c.Request().Context(), userDTO.ToUser())
	if err != nil {
		log.Printf("user with such email already exists: %v", err)
		return UserAlreadyExistsErr
	}
	userDTO.ID = id

	return nil
}

// Login - returns jwt token if success, otherwise error
func (u *UserService) Login(c echo.Context, userDTO *DTO) (string, error) {
	user, err := u.repository.GetByEmail(c.Request().Context(), userDTO.ToUser())

	if err != nil {
		log.Printf("no user with such email: %v", err)
		return "", UserNotFoundErr
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password)); err != nil {
		log.Printf("wrong password: %v", err)
		return "", UserWrongPasswordErr
	}

	userData := auth.UserData{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	token, err := auth.GenerateToken(userData, []byte(auth.GetJWTSecret().Secret))

	if err != nil {
		return "", UserTokenErr
	}

	return token, nil
}
