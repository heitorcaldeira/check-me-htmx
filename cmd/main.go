package main

import (
	"fmt"
	"html/template"
	"io"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	todo "github.com/heitorcaldeira/check-me-htmx/pkg"
)

type Templates struct {
  templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
  return &Templates{
    templates: template.Must(template.ParseGlob("views/*.html")),
  }
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

  e.Renderer = newTemplate()

  todo.StartConnection()

  e.GET("/", func(c echo.Context) error {
    list := todo.GetTodos()

    return c.Render(200, "index", todo.TodoList{List: list})
  })

  e.POST("/add", func(c echo.Context) error {
    title := c.FormValue("title")

    item := todo.NewTodo(title)
    item.Save()

    list := todo.GetTodos()
    return c.Render(200, "list", todo.TodoList{List: list})
  })

  e.POST("/delete", func(c echo.Context) error {
    err := todo.DeleteAll()

    if err != nil {
      fmt.Println("error deleting all")
    }

    list := todo.GetTodos()
    return c.Render(200, "list", todo.TodoList{List: list})
  })

  e.POST("/toggle/:id/:done", func(c echo.Context) error {
    id := c.Param("id")
    idParsed, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
      fmt.Println("failed parsing int")
    }

    done := c.Param("done")
    doneParsed, err := strconv.ParseBool(done)
    if err != nil {
      fmt.Println("failed parsing bool")
    }

    doneParsed = !doneParsed
    err = todo.UpdateById(idParsed, doneParsed)

    if err != nil {
      fmt.Println("error updating by id")
    }

    list := todo.GetTodos()
    return c.Render(200, "list", todo.TodoList{List: list})
  })


  e.Logger.Fatal(e.Start(":42069"))
}
