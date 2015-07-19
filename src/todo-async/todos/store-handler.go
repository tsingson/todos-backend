package todos

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
)

type Job interface {
	ExitChan() chan error
	Run(todos map[string]Todo) (map[string]Todo, error)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("error: ", err)
	}
	return
}
func ProcessJobs(jobs <-chan Job, db string) {
	todos := make(map[string]Todo, 0) // initial a null map
	var err error
	// file storage
	// todo:   monitor the file change and reload in
	todos, err = loadTodoListfromFile(db)
	checkError(err)

	for {
		j := <-jobs

		// store the data to  todos map ( hash list )
		todosMod, err := j.Run(todos) //  transfer todo data to channel and check change or not in return

		if err == nil && todosMod != nil { // identify the change
			saveTodoListtoFile(db, todosMod) // save to file
		}

		j.ExitChan() <- err
	}
}
