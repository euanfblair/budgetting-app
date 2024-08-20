package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *Application) Profile(c echo.Context) error {
	data := TemplateData{
		Title: "Profile",
	}

	data.IsAuthenticated = app.SessionManager.Exists(c.Request().Context(), "authUserID")
	if !data.IsAuthenticated {
		return c.Redirect(http.StatusFound, "/login")
	}
	userId := app.getUserIdFromSession(c)
	data.UserData = app.Users.GetCurrentUser(userId)

	fmt.Println(userId)
	for _, data := range data.UserData {
		fmt.Println(data)
	}

	err := c.Render(http.StatusOK, "profile", data)
	if err != nil {
		println(err.Error())
	}
	return err
}

func (app *Application) getUserIdFromSession(c echo.Context) int {
	userID, ok := app.SessionManager.Get(c.Request().Context(), "authUserID").(int)
	if !ok {
		return 0
	}
	return userID
}
