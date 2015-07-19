package todos

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// Generate a uuid to use as a unique identifier for each Todo
// http://play.golang.org/p/4FkNSiUDMg
func newUUID() (string, error) {
	u1 := uuid.NewV4()
	return fmt.Sprintf("%s", u1), nil
	//return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
