package main

import (
	// web
	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
	//	"fmt"
	"runtime"
	"todo-async/todos"
	// debug
	//"github.com/DeanThompson/ginpprof"
	_ "github.com/icattlecoder/godaemon"
	// run current program as terminate-stay-resident  or daemonizing current program
	//  -d=true

	//"runtime/pprof"
)

const Db = "/Users/qinshen/git/go-project/go-todos/bin/todo-data/todos.json"

var (
	logger log.Logger
)

func init() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu * 2) // 尝试使用所有可用的CPU
}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	/*
		// init postgres DB connection
		dbhost := "localhost"
		dbuser := "postgres"
		dbpassword := "postgres"
		dbname := "tsingcloud"

		taskPgdb := new(postgres.PostgresDB)
		taskPgdb.InitDb(dbhost, dbuser, dbpassword, dbname)
		defer taskPgdb.Pool.Close()

		taskPgdb.ListTasks()
	*/

	// create channel to communicate over
	jobs := make(chan todos.Job)

	log.Crit("start process job")
	// start watching jobs channel for work
	go todos.ProcessJobs(jobs, Db)

	// create dependencies
	client := &todos.TodoClient{Jobs: jobs}
	handlers := &TodoHandlers{Client: client}

	//var staticHtmlPath string
	staticHtmlPath := "/Users/qinshen/git/web-project/redux-learning/redux-async-learning/dist"
	// configure routes
	router := gin.Default()
	router.Static("/static", staticHtmlPath)
	//*****************************************************************************
	//  jwt token handle
	//*****************************************************************************

	//router.Use( CommHeade)

	router.GET("/user/token", JwtGetToken)
	router.POST("/user/balance", JwtCheckToken)

	//*****************************************************************************
	//  action for todos
	//*****************************************************************************

	v1 := router.Group("/v1")
	{
		v1.POST("/todo", handlers.AddTodo)
		v1.GET("/todo", handlers.GetTodos)
		v1.GET("/todo/:id", handlers.GetTodo)
		v1.PUT("/todo/:id", handlers.SaveTodo)
		v1.DELETE("/todo/:id", handlers.DeleteTodo)
	}

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/github")
	})

	//*****************************************************************************
	// test only  end
	//*****************************************************************************

	// start web server
	router.Run(":8080")

	/*
		routerAdmin := gin.Default()
		// start web server
		// debug
		// automatically add routers for net/http/pprof
		// e.g. /debug/pprof, /debug/pprof/heap, etc.
		routerAdmin.GET("/test", func(c *gin.Context) {
			c.Writer.Header().Set("link", nextPageUrl)
			c.Writer.Header().Set("token", token)
			c.Writer.Header().Set("X-GitHub-Media-Type", "github.v3")

			c.Data(200, "application/json; charset=utf-8", jsonData)
		})

		ginpprof.Wrapper(routerAdmin)
		//routerAdmin.Run(":8091")

	*/
}
