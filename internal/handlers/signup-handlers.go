package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (app *Application) Signup(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "Signup",
	}
	return c.Render(http.StatusOK, "signup", data)
}

func (app *Application) CreateUser(c echo.Context) error {

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = app.Users.Insert(username, email, passwordHash)
	if err != nil {
		return err
	}

	fmt.Println(username, email, password)
	return nil
}

func (app *Application) PasswordStrengthPost(c echo.Context) error {
	passwordString := c.FormValue("password")

	entropy := passwordvalidator.GetEntropy(passwordString)
	const minEntropyBits = 70
	fmt.Println(entropy)

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

//func (app *Application) UserSignupPost(w http.ResponseWriter, r *http.Request) {
//
//	err := r.ParseForm()
//	if err != nil {
//		helpers.ServeError(w, r, err)
//		return
//	}
//
//	var errMessage string
//
//	username := r.Form.Get("username")
//	email := r.Form.Get("email")
//	password := r.Form.Get("password")
//
//	fmt.Println(username, email, password)
//
//	switch {
//	case validators.LengthValidate(password, 8, 9999):
//		errMessage = "Password must be at least 8 characters in length"
//	case validators.LengthValidate(username, 6, 12):
//		errMessage = "Username must be 6-12 characters in length"
//	case !validators.Matches(email, emailRegex):
//		errMessage = "Please enter a valid email address"
//	case app.Users.Validate(username, email):
//		errMessage = "Username or Email already in system"
//	case !validators.PasswordRequirements(password):
//		errMessage = "Password must meet requirements below"
//	}
//
//	if errMessage != "" {
//		w.Header().Set("Content-Type", "text/html")
//		w.WriteHeader(http.StatusBadRequest)
//
//		// Render the error message template with the error message
//		data := struct {
//			ErrorMessage string
//		}{
//			ErrorMessage: errMessage,
//		}
//
//		app.defaultRenderHandler(w, r, "signup.html", data, "error_template")
//		return
//	}
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
//
//	err = app.Users.Insert(username, email, hashedPassword)
//	if err != nil {
//		helpers.ServeError(w, r, err)
//		return
//	}
//
//	//http.Redirect(w, r, "/user/login", http.StatusFound)
//
//}
