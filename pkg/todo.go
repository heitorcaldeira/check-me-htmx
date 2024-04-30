package todo

import "fmt"

type TodoList struct {
	List []Todo
}

type Todo struct {
	Id    int
	Title string
	Done  bool
}

func NewTodo(title string) Todo {
	return Todo{Title: title, Done: false}
}

func (t *Todo) Save() {
  _, err := Insert(t)
  if err != nil {
    fmt.Print("insert failed")
  }
}
