package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/category"
	"github.com/Mickey327/rcsp-backend/internal/app/company"
	appConfig "github.com/Mickey327/rcsp-backend/internal/app/config"
	"github.com/Mickey327/rcsp-backend/internal/app/user"
	dbConfig "github.com/Mickey327/rcsp-backend/internal/db/config"
	"github.com/Mickey327/rcsp-backend/internal/db/repository/postgres"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	appConf := appConfig.GetConfig()
	dbConf := dbConfig.GetConfig()
	db, err := postgres.New(ctx, dbConf.GenerateConnectPath())
	if err != nil {
		log.Fatal(err)
	}

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(auth.GetJWTSecret().Secret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.Claims)
		},
	})

	categoryHandler := category.NewHandler(category.NewService(category.NewRepository(db)))
	e.GET("/api/category/:id", categoryHandler.Read)
	e.GET("/api/category", categoryHandler.ReadAll)
	e.DELETE("/api/category/:id", categoryHandler.Delete, jwtMiddleware)
	e.POST("/api/category", categoryHandler.Create, jwtMiddleware)
	e.PUT("/api/category", categoryHandler.Update, jwtMiddleware)

	companyHandler := company.NewHandler(company.NewService(company.NewRepository(db)))
	e.GET("/api/company/:id", companyHandler.Read)
	e.GET("/api/company", companyHandler.ReadAll)
	e.DELETE("/api/company/:id", companyHandler.Delete, jwtMiddleware)
	e.POST("/api/company", companyHandler.Create, jwtMiddleware)
	e.PUT("/api/company", companyHandler.Update, jwtMiddleware)

	e.GET("/hello", func(c echo.Context) error {
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*auth.Claims)
		email := claims.Email
		role := claims.Role
		return c.String(http.StatusOK, "Welcome "+email+", Your Role "+role+"!")
	}, jwtMiddleware)

	userHandler := user.NewHandler(user.NewService(user.NewRepository(db)))
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{appConf.ClientHost + ":" + appConf.ClientPort},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowCredentials: true,
	}))
	e.Logger.Fatal(e.Start(":" + appConf.ApiPort))
}
