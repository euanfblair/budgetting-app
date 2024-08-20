package handlers

import (
	"euanfblair/budgeting-app/internal/models"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Application struct {
	Users          *models.UserModel
	SessionManager *scs.SessionManager
	//Transactions   *models.TransactionModel
}

type TemplateData struct {
	Title           string
	PasswordEntropy int
	ErrorMessage    string
	IsAuthenticated bool
	UserData        []string
}

// Home handler function
func (app *Application) Home(c echo.Context) error {

	data := TemplateData{
		Title: "Home",
	}

	data.IsAuthenticated = app.SessionManager.Exists(c.Request().Context(), "authUserID")
	return c.Render(http.StatusOK, "home", data)
}
