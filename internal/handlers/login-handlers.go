package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func (app *Application) Login(c echo.Context) error {
	data := TemplateData{
		Title: "Login",
	}
	data.IsAuthenticated = app.SessionManager.Exists(c.Request().Context(), "authUserID")
	if data.IsAuthenticated {
		return c.Redirect(http.StatusFound, "/")
	}
	return c.Render(http.StatusOK, "login", data)
}

func (app *Application) ValidateUser(c echo.Context) error {

	email := c.FormValue("email")
	password := c.FormValue("password")

	fmt.Println(email, password)

	data := TemplateData{
		ErrorMessage: "",
	}

	userId, inputHash := app.Users.Login(email)

	switch {
	case userId == 0:
		data.ErrorMessage = "Invalid Email or Password"
	case !app.passwordsMatch(password, inputHash):
		data.ErrorMessage = "Invalid Email or Password"
	}

	if strings.Compare(data.ErrorMessage, "") != 0 {
		err := c.Render(http.StatusOK, "error-message", data)
		return err
	}

	err := app.SessionManager.RenewToken(c.Request().Context())
	if err != nil {
		return err
	}

	app.SessionManager.Put(c.Request().Context(), "authUserID", userId)

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}

	return nil
}

func (app *Application) Logout(c echo.Context) error {
	err := app.SessionManager.RenewToken(c.Request().Context())
	if err != nil {
		return err
	}

	app.SessionManager.Remove(c.Request().Context(), "authUserID")
	app.SessionManager.Put(c.Request().Context(), "flash", "You have successfully logged out")

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}

	return nil
}

func (app *Application) passwordsMatch(password, storedHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
