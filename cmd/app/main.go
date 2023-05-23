package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Mickey327/rcsp-backend/internal/app/auth"
	"github.com/Mickey327/rcsp-backend/internal/app/cart"
	"github.com/Mickey327/rcsp-backend/internal/app/category"
	"github.com/Mickey327/rcsp-backend/internal/app/comment"
	"github.com/Mickey327/rcsp-backend/internal/app/company"
	appConfig "github.com/Mickey327/rcsp-backend/internal/app/config"
	"github.com/Mickey327/rcsp-backend/internal/app/order"
	"github.com/Mickey327/rcsp-backend/internal/app/orderItem"
	"github.com/Mickey327/rcsp-backend/internal/app/product"
	"github.com/Mickey327/rcsp-backend/internal/app/user"
	"github.com/Mickey327/rcsp-backend/internal/app/validator"
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
		TokenLookup: "header:Authorization:Bearer ,cookie:jwt",
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

	productHandler := product.NewHandler(product.NewService(product.NewRepository(db)))
	e.GET("/api/product/:id", productHandler.Read)
	e.GET("/api/product", productHandler.ReadAll) // ?categoryID&companyID
	e.DELETE("/api/product/:id", productHandler.Delete, jwtMiddleware)
	e.POST("/api/product", productHandler.Create, jwtMiddleware)
	e.PUT("/api/product", productHandler.Update, jwtMiddleware)

	valid := validator.NewValidator()
	userHandler := user.NewHandler(user.NewService(user.NewRepository(db)))
	e.Validator = valid
	e.POST("/api/register", userHandler.Register)
	e.POST("/api/login", userHandler.Login)
	e.GET("/api/logout", userHandler.Logout)
	e.GET("/api/user", userHandler.GetAuthenticatedUser, jwtMiddleware)

	cartHandler := cart.NewHandler(cart.NewService(order.NewRepository(db), orderItem.NewRepository(db)))
	e.POST("/api/cart", cartHandler.UpdateCart, jwtMiddleware)       // ?orderID&productID
	e.DELETE("/api/cart", cartHandler.RemoveFromCart, jwtMiddleware) // ?orderID&productID

	orderHandler := order.NewHandler(order.NewService(order.NewRepository(db)))
	e.GET("/api/order", orderHandler.ReadCurrentUserArrangingOrder, jwtMiddleware)
	e.POST("/api/order", orderHandler.Create, jwtMiddleware)
	e.GET("/api/order/:id", orderHandler.ReadByIdEager, jwtMiddleware)
	e.PUT("/api/order", orderHandler.Update, jwtMiddleware)

	commentHandler := comment.NewHandler(comment.NewService(comment.NewRepository(db)))
	e.POST("/api/comment", commentHandler.WriteComment, jwtMiddleware) //?productID
	e.GET("/api/comment", commentHandler.ReadComments)                 //?productID

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{appConf.ClientHost + ":" + appConf.ClientPort, "http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowCredentials: true,
	}))

	e.Static("/", "static")
	e.Logger.Fatal(e.Start(":" + appConf.ApiPort))
}
