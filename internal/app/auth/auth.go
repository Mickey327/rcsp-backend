package auth

import (
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
)

var (
	instance *JWTSecret
	once     sync.Once
)

type JWTSecret struct {
	Secret        string `env:"JWT_SECRET"`
	RefreshSecret string `env:"JWT_REFRESH_SECRET"`
}

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
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

func GenerateToken(email string, role string, secret []byte) (string, error) {
	claims := &Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Printf("New token generated: %v", token)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserEmailAndRole - helps to retrieve info about user from jwt token (will be helpful for other handlers)
func GetUserEmailAndRole(c echo.Context) (string, string) {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)
	return claims.Email, claims.Role
}
