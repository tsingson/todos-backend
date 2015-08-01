package main

import (
	"fmt"
	"guid/guid"
)

func main() {

	objID := guid.NewObjectId()
	fmt.Println(objID)
	fmt.Println(objID.Hex())

}
