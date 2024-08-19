package handlers

import (
	"euanfblair/budgeting-app/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Application struct {
	Users *models.UserModel
	//SessionManager *scs.SessionManager
	//Transactions   *models.TransactionModel
}

type TemplateData struct {
	Title           string
	PasswordEntropy int
}

// Home handler function
func (app *Application) Home(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "Home",
	}
	return c.Render(http.StatusOK, "home", data)
}
