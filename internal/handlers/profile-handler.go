package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func (app *Application) Profile(c echo.Context) error {
	data := TemplateData{
		Title:           "Profile",
		PasswordEntropy: 0,
	}

	data.IsAuthenticated = app.SessionManager.Exists(c.Request().Context(), "authUserID")
	if !data.IsAuthenticated {
		return c.Redirect(http.StatusFound, "/login")
	}
	userId := app.getUserIdFromSession(c)
	data.UserData = app.Users.GetCurrentUser(userId)

	err := c.Render(http.StatusOK, "profile", data)
	if err != nil {
		println(err.Error())
	}
	return err
}

func (app *Application) ChangeUserPassword(c echo.Context) error {
	userId := app.getUserIdFromSession(c)

	oldPassword := c.FormValue("old-password")
	newPassword := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	fmt.Println(oldPassword)

	currentPassword := app.Users.GetCurrentPassword(userId)

	data := TemplateData{}

	switch {
	case bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(oldPassword)) != nil:
		data.ErrorMessage = "Current Password is not correct, please try again"
	case strings.Compare(newPassword, confirmPassword) != 0:
		data.ErrorMessage = "New Passwords do not match, please try again"
	case passwordvalidator.GetEntropy(newPassword) < minEntropyBits:
		data.ErrorMessage = "Password strength must be 100%"
	}

	if strings.Compare(data.ErrorMessage, "") != 0 {
		err := c.Render(http.StatusBadRequest, "error-message", data)
		return err
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = app.Users.UpdatePassword(userId, newPasswordHash)
	if err != nil {
		return err
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		err := app.SessionManager.RenewToken(c.Request().Context())
		if err != nil {
			return err
		}

		app.SessionManager.Remove(c.Request().Context(), "authUserID")
		app.SessionManager.Put(c.Request().Context(), "flash", "You have successfully logged out")

		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusOK)
	}

	return nil
}

func (app *Application) DeleteUser(c echo.Context) error {
	userId := app.getUserIdFromSession(c)

	fmt.Println(userId)

	err := app.Users.DeleteUser(userId)
	if err != nil {
		return err
	}

	app.SessionManager.Remove(c.Request().Context(), "authUserID")
	app.SessionManager.Put(c.Request().Context(), "flash", "You have successfully logged out")

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)

}

func (app *Application) getUserIdFromSession(c echo.Context) int {
	userID, ok := app.SessionManager.Get(c.Request().Context(), "authUserID").(int)
	if !ok {
		return 0
	}
	return userID
}
