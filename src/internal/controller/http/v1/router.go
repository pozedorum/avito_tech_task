package v1

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pozedorum/user-balance-service/internal/service"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	logsFilePath = "/logs/requests.log"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := handler.Group("/auth")

	newAuthRoutes(auth, services.Auth)

	authMiddleware := &authMiddleware{services.Auth}
	v1 := handler.Group("/api/v1", authMiddleware.UserIdentity)

	newAccountRoutes(v1.Group("/accounts"), services.Account)
	newReservationRoutes(v1.Group("/reservations"), services.Reservation)
	newProductRoutes(v1.Group("/products"), services.Product)
	newOperationRoutes(v1.Group("/products"), services.Operation)

}

func setLogsFile() *os.File {
	file, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
