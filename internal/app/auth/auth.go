package auth

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Mickey327/rcsp-backend/internal/app/response"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
)

var (
	instance *JWTSecret
	once     sync.Once
)

type JWTSecret struct {
	Secret                string        `env:"JWT_SECRET"`
	RefreshSecret         string        `env:"JWT_REFRESH_SECRET"`
	ExpirationTimeInHours time.Duration `env:"EXPIRATION_TIME_IN_HOURS"`
}

type UserData struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Claims struct {
	UserData
	jwt.RegisteredClaims
}

func GetJWTSecret() *JWTSecret {
	once.Do(func() {
		log.Println("gather jwt-secret config")
		instance = &JWTSecret{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "Gametrade - the best gaming store"
			description, _ := cleanenv.GetDescription(instance, &helpText)
			log.Println(description)
			log.Fatal(err)
		}
	})
	return instance
}

//TODO: Add refresh token

func GenerateToken(user UserData, secret []byte) (string, error) {
	claims := &Claims{
		UserData: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(GetJWTSecret().ExpirationTimeInHours)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserDataFromToken(token *jwt.Token) UserData {
	claims := token.Claims.(*Claims)
	return claims.UserData
}

func GetUserToken(c echo.Context) (*jwt.Token, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetJWTSecret().Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetUserDataAndCheckRole(c echo.Context, roles ...string) (*UserData, error) {
	token, err := GetUserToken(c)
	log.Println(token)
	if err != nil {
		return nil, c.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Message: "can't get jwt token from cookie",
		})
	}
	check := false

	userData := GetUserDataFromToken(token)
	for _, role := range roles {
		if userData.Role == role {
			check = true
		}
	}
	log.Println(userData.ID, userData.Role, userData.Email)

	if check == false {
		return nil, c.JSON(http.StatusForbidden, response.Response{
			Code:    http.StatusForbidden,
			Message: "you don't have enough rights to do that action",
		})
	}

	if userData.ID <= 0 {
		return nil, c.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Message: "wrong user id value",
		})
	}

	return &userData, nil
}
