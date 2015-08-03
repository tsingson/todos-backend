package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Twitter Twitter `json:"twitter"`
}

type Twitter struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func main() {

	raw := `{
        "name":"Hodor",
        "email":"hodor@hodor.io",
        "twitter": {"id": 123,"username": "hodor"}
    }`

	user := &User{}
	json.Unmarshal([]byte(raw), &user)

	fmt.Println(user.Email)
	// hodor@hodor.io
	fmt.Println(user.Twitter.Id)
	// 123

	raw1 := `{"pi":3.14,"langs":["go","c","haskell","erlang","python"]}`

	var yolo map[string]interface{}

	if err := json.Unmarshal([]byte(raw1), &yolo); err != nil {
		panic(err)
	}

	pi := yolo["pi"].(float64)
	fmt.Println(pi)
	// 3.14

	langs := yolo["langs"].([]interface{})
	fmt.Println(len(langs))
	// 5
	lang1 := langs[0].(string)
	fmt.Println(lang1)
	// go

	type Country struct {
		Name              string `json:"name"`
		NationalDance     string `json:"national_dance,omitempty"`
		NationalDrink     string `json:"-"`
		nuclearLaunchCode int
	}

	russia := &Country{"Russia", "troika", "vodka", 4321}
	r, _ := json.Marshal(russia)
	fmt.Println(string(r))
	// {"name":"Russia","national_dance":"troika"}

	usa := &Country{"USA", "", "bourbon", 1234}
	u, _ := json.Marshal(usa)
	fmt.Println(string(u))
	//{"name":"USA"}

}
