package routers

import (
	"euanfblair/budgeting-app/internal/handlers"
	"github.com/labstack/echo/v4"
)

func InitGetRoutes(e *echo.Echo, app *handlers.Application) {
	e.Static("/static", "ui/static")
	e.GET("/", app.Home)
	e.GET("/signup", app.Signup)
	e.GET("/login", app.Login)
	e.GET("/user-profile", app.Profile)
	e.GET("/transactions", app.GetTransactions)
	e.GET("/transactions/tab", app.FilteredTransactions)
	e.GET("/transactions/category", app.FilterCategory)
	e.GET("/transactions/next", app.NextPage)
	e.GET("/transactions/prev", app.PrevPage)
}

func InitPostRoutes(e *echo.Echo, app *handlers.Application) {
	e.POST("/signup/create", app.CreateUser)
	e.POST("/signup/password_check", app.PasswordStrengthPost)
	e.POST("/login/validate", app.ValidateUser)
	e.POST("/logout", app.Logout)
	e.POST("/profile/password_check", app.PasswordStrengthPost)
}

func InitPutRoutes(e *echo.Echo, app *handlers.Application) {
	e.PUT("/profile/change_password", app.ChangeUserPassword)
}

func InitDeleteRoutes(e *echo.Echo, app *handlers.Application) {
	e.DELETE("/profile/delete", app.DeleteUser)
	e.DELETE("/transactions/delete", app.DeleteTransaction)
}
