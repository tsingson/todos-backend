package todos

import (
	"fmt"
)

// Job to add a Todo to the database
type SaveTodoJob struct {
	toSave   Todo
	saved    chan Todo
	exitChan chan error
}

func NewSaveTodoJob(todo Todo) *SaveTodoJob {
	return &SaveTodoJob{
		toSave:   todo,
		saved:    make(chan Todo, 1),
		exitChan: make(chan error, 1),
	}
}
func (j SaveTodoJob) ExitChan() chan error {
	return j.exitChan
}
func (j SaveTodoJob) Run(todos map[string]Todo) (map[string]Todo, error) {
	var todo Todo
	var id string

	currentId := len(todos) + 1

	if j.toSave.Id == "" || j.toSave.Code == "" { // new item
		id = fmt.Sprintf("%d", currentId)
		Code, err := newUUID()
		if err != nil {
			return nil, err
		}
		todo = Todo{Id: id, Code: Code, Description: j.toSave.Description}
	} else { // exited item
		todo = j.toSave
	}
	todos[todo.Id] = todo //  add to list or update the list

	j.saved <- todo
	return todos, nil
}
