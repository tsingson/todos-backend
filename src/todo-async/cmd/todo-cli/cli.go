package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	log "github.com/inconshreveable/log15"
	"os"
	"todo-async/todos"

//"encoding/json"
)

const Db = "/Users/qinshen/git/go-project/go-todos/bin/todo-data/todos.json"

var (
	logger log.Logger
)

func main() {

	// create channel to communicate over
	jobs := make(chan todos.Job)

	//log.Crit("start process job")
	// start watching jobs channel for work
	go todos.ProcessJobs(jobs, Db)

	// create dependencies
	client := &todos.TodoClient{Jobs: jobs}
	handlers := &TodoHandlers{Client: client}

	//************************************************************

	app := cli.NewApp()
	app.Name = "todo-cli"
	app.Usage = "todo cli"
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "lang, l",
			Value:  "english",
			Usage:  "language for the greeting",
			EnvVar: "LEGACY_COMPAT_LANG,APP_LANG,LANG",
		},
	}

	app.Action = func(c *cli.Context) {
		//name := "someone"
		if len(c.Args()) > 0 {
			name := c.Args()[0]
			fmt.Println("name is ", name)
		} else {
			cli.ShowAppHelp(c)
		}
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) {
				handlers.AddTodo(c)
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a todo to the list",
			Action: func(c *cli.Context) {
				handlers.AddTodo(c)
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update todo by id ",
			Action: func(c *cli.Context) {
				handlers.SaveTodo(c)
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete a todo by id",
			Action: func(c *cli.Context) {
				handlers.DeleteTodo(c)
			},
		},
		{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "query a todo by id",
			Action: func(c *cli.Context) {
				handlers.GetTodo(c)
			},
		},

		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all todos",
			Action: func(c *cli.Context) {
				handlers.GetTodos(c)
			},
		},
	}

	app.Run(os.Args)
}
