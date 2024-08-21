package main

import (
	"database/sql"
	"euanfblair/budgeting-app/cmd/web/routers"
	"euanfblair/budgeting-app/internal/handlers"
	"euanfblair/budgeting-app/internal/models"
	"fmt"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	session "github.com/spazzymoto/echo-scs-session"
	"html/template"
	"io"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const componentRoot = "ui/html/components/"

var sharedComponents = []string{
	componentRoot + "footer.html",
	componentRoot + "nav.html",
	"ui/html/base.html",
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

const postgresCred = "postgres://postgres:Pass@localhost/budgetting_app?sslmode=disable"

const templates = "ui/html/**/*.html"

func main() {

	db, err := openDB(postgresCred)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := handlers.Application{
		Users:          &models.UserModel{DB: db},
		SessionManager: sessionManager,
		Transactions:   &models.TransactionModel{DB: db},
	}

	e := echo.New()

	e.Use(session.LoadAndSave(sessionManager))

	t := &Template{
		template.Must(template.ParseGlob(templates)),
	}

	e.Renderer = t

	routers.InitGetRoutes(e, &app)
	routers.InitPostRoutes(e, &app)
	routers.InitPutRoutes(e, &app)
	routers.InitDeleteRoutes(e, &app)
	e.Logger.Fatal(e.StartTLS(":4000", "cert.pem", "key.pem"))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
