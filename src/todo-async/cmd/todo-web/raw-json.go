package main

import (
	"github.com/gin-gonic/gin"
)

//*****************************************************************************
// test only
//*****************************************************************************

/*

	nextPageUrl := "<https://api.github.com/user/tsingson/?page2>;rel=\"next\""

	token := "mytoken1234123412341234"
	jsonData := []byte(`{
  "login": "tsingson",
  "id": 3875274,
  "avatar_url": "https://avatars.githubusercontent.com/u/3875274?v=3"
}
`)

*/

type JwtPayload struct {
	Header  map[string]string
	Payload []byte
}

func (j *JwtPayload) RawJSON(c *gin.Context) {
	/*
		c.Writer.Header().Set("link", nextPageUrl)
		c.Writer.Header().Set("token", token)
		c.Writer.Header().Set("X-GitHub-Media-Type", "github.v3")
	*/

	if len(j.Header) > 0 {
		for key, value := range j.Header {
			c.Writer.Header().Set(key, value)
		}
	}

	c.Data(200, "application/json; charset=utf-8", j.Payload)
}
