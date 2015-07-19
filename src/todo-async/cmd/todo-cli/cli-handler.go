//*****************************************************************
// function define
//*****************************************************************

package main

import (
	//log "github.com/inconshreveable/log15"
	"fmt"
	"github.com/codegangsta/cli"
	"todo-async/todos"
)

type TodoHandlers struct {
	Client *todos.TodoClient
}

// Add a new todo
func (h *TodoHandlers) AddTodo(c *cli.Context) {
	var todo todos.Todo
	todo.Id = ""
	todo.Description = c.Args().First()

	created, err := h.Client.SaveTodo(todo)
	if err != nil {
		////log.Crit(err)
		fmt.Println("resultCode: ", 500, "resultInfo:", "problem decoding body")
		return
	}
	fmt.Println("resultCode: ", 201, "resultInfo: ", created)
}

// Get all todos as an array
func (h *TodoHandlers) GetTodos(c *cli.Context) {
	todos, err := h.Client.GetTodos()
	if err != nil {
		//log.Crit(err)
		fmt.Println(500, "problem decoding body")
		return
	}

	for index, todo := range todos {
		fmt.Println(index)
		fmt.Println(todo)
	}
	//fmt.Println(200, todos)
}

// Get a specific todo by id
func (h *TodoHandlers) GetTodo(c *cli.Context) {
	id := c.Args().First()
	todo, err := h.Client.GetTodo(id)
	if err != nil {
		//log.Crit(err)
		fmt.Println(500, "problem decoding body")
		return
	}

	fmt.Println(200, todo)
}

// Add a new todo
func (h *TodoHandlers) SaveTodo(c *cli.Context) {
	id := c.Args().First()
	Description := c.Args().Get(1)

	var todo todos.Todo
	todo.Id = id
	todo.Description = Description

	saved, err := h.Client.SaveTodo(todo)
	if err != nil {
		//log.Crit(err)
		fmt.Println(500, "problem decoding body")
		return
	}

	fmt.Println(200, saved)
}

// Delete a todo by id
func (h *TodoHandlers) DeleteTodo(c *cli.Context) {
	id := c.Args().First() //

	if err := h.Client.DeleteTodo(id); err != nil {
		//log.Crit(err)
		fmt.Println(500, "problem decoding body")
		return
	}

	fmt.Println(204, "application/json", make([]byte, 0))
}
