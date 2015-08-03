package main

import (
	"fmt"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

// Binding from JSON
type LoginJSON struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Binding from form values
type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

var (
	mysupersecretpassword = "unicornsAreAwesome"
)

type User struct {
	Id        int64  `db:"id" json:"id"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
}

func GetUsers(c *gin.Context) {
	type Users []User
	var users = Users{
		User{Id: 1, Firstname: "Oliver", Lastname: "Queen"},
		User{Id: 2, Firstname: "Malcom", Lastname: "Merlyn"},
	}
	c.JSON(200, users)
	// curl -i http://localhost:8080/api/v1/users
}

func main() {
	router := gin.Default()

	v1 := router.Group("api/v1")
	{
		v1.GET("/users", GetUsers)
	}

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", func(c *gin.Context) {

		fmt.Println( "/loginJSon been call ")
		var json LoginJSON


		c.Bind(&json) // This will infer what binder to use depending on the content-type header.

		req := c.Request
		requestToken := req.Header.Get("Authorization")

		/*
				if ah := req.Header.Get("Authorization"); ah != "" {
				// Should be a bearer token
				if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
					return Parse(ah[7:], keyFunc)
				}
			}
		*/

		fmt.Println("user is:", json.User)
		fmt.Println("password", json.Password)
		fmt.Println("requestToken", requestToken)
		/* debug  javascript fetch */
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in", "user": json.User, "password": json.Password, "token": requestToken})
		fmt.Println("gin.H: ", gin.H{"status": "you are logged in"})

		/* debug end */

		/*
			if json.User == "user" && json.Password == "password" {

				token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
				// Set some claims
				token.Claims["ID"] = json.User
				token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
				// Sign and get the complete encoded token as a string
				tokenString, _ := token.SignedString([]byte(mysupersecretpassword))

				c.JSON(http.StatusOK, gin.H{"status": "you are logged in", "id_token": tokenString})
				fmt.Println("gin.H: ", gin.H{"status": "you are logged in", "id_token": tokenString})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		*/
	})

	router.GET("/auth", func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		fmt.Println("Authorization: ", token)
		c.JSON(200, gin.H{"Authorization": token})
	})

	/*
		// Example for binding a HTML form (user=manu&password=123)
		router.POST("/loginHTML", func(c *gin.Context) {
			var form LoginForm

			c.BindWith(&form, binding.Form) // You can also specify which binder to use. We support binding.Form, binding.JSON and binding.XML.
			if form.User == "manu" && form.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		})
	*/

	staticHtmlPath := "/Users/qinshen/git/personal_project/todos/todos-font/layout-demo"
	// "/Users/qinshen/git/web-project/dashboard-ui-es6/build"

	router.Static("/static", staticHtmlPath)

	public := router.Group("/api")

	public.GET("/", func(c *gin.Context) {
		// Create the token
		token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
		// Set some claims
		token.Claims["ID"] = "Christopher"
		token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte(mysupersecretpassword))
		if err != nil {
			c.JSON(500, gin.H{"message": "Could not generate token"})
		}
		c.JSON(201, gin.H{"id_token": tokenString})
	})

	private := router.Group("/api/private")
	private.Use(jwt.Auth(mysupersecretpassword))

	/*
		Set this header in your request to get here.
		Authorization: Bearer `token`
	*/

	private.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from private"})
	})



















	router.Run("localhost:8080")
}
