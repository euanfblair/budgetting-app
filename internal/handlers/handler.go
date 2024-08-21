package handlers

import (
	"euanfblair/budgeting-app/internal/models"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"net/http"
)

type Application struct {
	Users          *models.UserModel
	SessionManager *scs.SessionManager
	Transactions   *models.TransactionModel
}

type TemplateData struct {
	Title           string
	PasswordEntropy int
	ErrorMessage    string
	IsAuthenticated bool
	UserData        []string
	ActiveTab       string
	TableData       []tableData
	TotalAmount     models.Money
}

// Home handler function
func (app *Application) Home(c echo.Context) error {

	data := TemplateData{
		Title: "Home",
	}

	data.IsAuthenticated = app.SessionManager.Exists(c.Request().Context(), "authUserID")
	return c.Render(http.StatusOK, "home", data)
}

// PasswordStrengthPost Render Password Strength Bar
func (app *Application) PasswordStrengthPost(c echo.Context) error {
	passwordString := c.FormValue("password")

	entropy := passwordvalidator.GetEntropy(passwordString)

	strengthPercent := (entropy / minEntropyBits) * 100
	if strengthPercent > 100 {
		strengthPercent = 100
	}

	data := TemplateData{
		PasswordEntropy: int(strengthPercent),
	}

	err := c.Render(http.StatusOK, "password-strength", data)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return nil
}
