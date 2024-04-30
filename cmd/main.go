package main

import (
	"html/template"
	"io"
	"math/rand"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

type TodoList struct {
  List []Todo
}

type Todo struct {
  Id int
  Title string
  Done bool
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

  e.Renderer = newTemplate()

  todoList := TodoList{ List: make([]Todo, 0) }
  todoList.List = append(todoList.List, Todo{ Id: 10, Title: "do this shit", Done: false })

  e.GET("/", func(c echo.Context) error {
    return c.Render(200, "index", todoList)
  })

  e.POST("/add", func(c echo.Context) error {
    title := c.FormValue("title")
    id := rand.Int()
    todoList.List = append(todoList.List, Todo{ Id: id, Title: title, Done: false })
    return c.Render(200, "list", todoList)
  })

  e.Logger.Fatal(e.Start(":42069"))
}
