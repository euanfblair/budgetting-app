package handlers

import (
	"github.com/labstack/echo/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"strings"
)

const minEntropyBits = 70

func (app *Application) Signup(c echo.Context) error {
	data := TemplateData{
		Title: "Signup",
	}
	data.IsAuthenticated = app.SessionManager.Exists(c.Request().Context(), "authUserID")
	if data.IsAuthenticated {
		return c.Redirect(http.StatusFound, "/user-profile")
	}
	return c.Render(http.StatusOK, "signup", data)
}

func (app *Application) CreateUser(c echo.Context) error {

	firstname := c.FormValue("first-name")
	surname := c.FormValue("surname")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirm := c.FormValue("confirm-password")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	passwordStrongEnough := passwordvalidator.GetEntropy(password) < minEntropyBits

	data := TemplateData{
		ErrorMessage: "",
	}

	switch {
	case passwordStrongEnough:
		data.ErrorMessage = "Password not strong enough"
	case !validEmail(email):
		data.ErrorMessage = "Please enter a valid email address format"
	case strings.Compare(password, confirm) != 0:
		data.ErrorMessage = "Passwords do not match"
	case app.Users.ExistingEmail(email):
		data.ErrorMessage = "Email already registered with account"
	}

	if strings.Compare(data.ErrorMessage, "") != 0 {
		err := c.Render(http.StatusBadRequest, "error-message", data)
		return err
	}

	err = app.Users.Insert(firstname, surname, email, passwordHash)
	if err != nil {
		return err
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusOK)
	}

	return nil
}

func validEmail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true

}
