package main

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/*
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request
		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

*/

type User struct {
	Name string `json: "name"`
	Mail string `json: "mail"`
	Pass string `json: "pass"`
}
type Login struct {
	Name string `json: "name" binding: "required"`
	Pass string `json: "pass" binding: "required"`
}

var (
	validUser    = User{Name: "rainy", Mail: "me@rainy.im", Pass: "123"}
	mySigningKey = "USEHERE"
)

func CommHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Run this on all requests
		// Should be moved to a proper middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Next()
	}
}

/*
	router.OPTIONS("/*cors", func(c *gin.Context) {
		// Empty 200 response
	})
*/

// generate jwt token and return to client
func JwtGetToken(c *gin.Context) {
	var login Login
	val := c.Bind(&login)
	if !val {
		c.JSON(200, gin.H{"code": 401, "msg": "Both name & password are required"})
		return
	}
	//	if login.Name == validUser.Name && login.Pass == validUser.Pass {
	token := jwt.New(jwt.SigningMethodHS256)
	// Headers
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"

	// Claims
	token.Claims["name"] = validUser.Name
	token.Claims["mail"] = validUser.Mail
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte(mySigningKey))

	fmt.Println("jwt token raw: ", tokenString)

	if err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": "Server error!"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "OK", "jwt": tokenString})
	//	} else {
	//		c.JSON(200, gin.H{"code": 400, "msg": "Error username or password!"})
	//	}
}

//  get client's jwt token and verify
func JwtCheckToken(c *gin.Context) {
	token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
		b := ([]byte(mySigningKey))
		return b, nil
	})
	fmt.Println(err)
	if err != nil {
		c.JSON(200, gin.H{"code": 403, "msg": err.Error()})
	} else {
		if token.Valid {
			token.Claims["balance"] = 49
			tokenString, err := token.SignedString([]byte(mySigningKey))
			if err != nil {
				c.JSON(200, gin.H{"code": 500, "msg": "Server error!"})
				return
			}

			c.JSON(200, gin.H{"code": 200, "msg": "OK", "jwt": tokenString})
		} else {
			c.JSON(200, gin.H{"code": 401, "msg": "Sorry, you are not validate"})
		}
	}
}

/*
// Try to find the token in an http.Request.
// This method will call ParseMultipartForm if there's no token in the header.
// Currently, it looks in the Authorization header as well as
// looking for an 'access_token' request parameter in req.Form.
func ParseFromRequest(req *http.Request, keyFunc Keyfunc) (token *Token, err error) {

	// Look for an Authorization header
	if ah := req.Header.Get("Authorization"); ah != "" {
		// Should be a bearer token
		if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
			return Parse(ah[7:], keyFunc)
		}
	}

	// Look for "access_token" parameter
	req.ParseMultipartForm(10e6)
	if tokStr := req.Form.Get("access_token"); tokStr != "" {
		return Parse(tokStr, keyFunc)
	}
	return nil, ErrNoTokenInRequest
}
*/
