package main

import (
	"html/template"
	"io"
	"net/http"
	"test/go_server/controllers"
	"test/go_server/db_client"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", "")
}

func Test(c echo.Context) error {
	return c.Render(200, "test.html", "")
}

func main() {
	db_client.InitialiseDBConnection()

	e := echo.New()
	e.Static("/", "public")
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.GET("/", Hello)
	e.GET("/test", Test)
	e.GET("/post", controllers.GetPosts)
	e.GET("/post/:id", controllers.GetPost)
	e.POST("/post", controllers.CreatePost)
	e.Logger.Fatal(e.Start(":8000"))
}
