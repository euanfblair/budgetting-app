package routers

import (
	"euanfblair/budgeting-app/internal/handlers"
	"github.com/labstack/echo/v4"
)

func InitGetRoutes(e *echo.Echo, app *handlers.Application) {
	e.Static("/static", "ui/static")
	e.GET("/", app.Home)
	e.GET("/signup", app.Signup)
}

func InitPostRoutes(e *echo.Echo, app *handlers.Application) {
	e.POST("/signup/create", app.CreateUser)
	e.POST("/signup/password_check", app.PasswordStrengthPost)
}
