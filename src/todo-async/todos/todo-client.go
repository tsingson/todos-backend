package todos

type TodoClient struct {
	Jobs chan Job
}

// getTodoHash
func (c *TodoClient) getTodoHash() (map[string]Todo, error) {
	job := NewReadTodosJob()
	c.Jobs <- job

	if err := <-job.ExitChan(); err != nil {
		return make(map[string]Todo, 0), err
	}
	return <-job.todos, nil
}

// get list ( array ) of todos
func (c *TodoClient) GetTodos() ([]Todo, error) {
	arr := make([]Todo, 0)

	todos, err := c.getTodoHash()
	if err != nil {
		return arr, err
	}

	for _, value := range todos {
		arr = append(arr, value)
	}
	return arr, nil
}

// get single todo
func (c *TodoClient) GetTodo(id string) (Todo, error) {
	todos, err := c.getTodoHash()
	if err != nil {
		return Todo{}, err
	}
	return todos[id], nil
}

// save todo ( create new one or update old one )
func (c *TodoClient) SaveTodo(todo Todo) (Todo, error) {
	job := NewSaveTodoJob(todo)
	c.Jobs <- job

	if err := <-job.ExitChan(); err != nil {
		return Todo{}, err
	}
	return <-job.saved, nil
}

// delete single todo with id
func (c *TodoClient) DeleteTodo(id string) error {
	job := NewDeleteTodoJob(id)
	c.Jobs <- job

	if err := <-job.ExitChan(); err != nil {
		return err
	}
	return nil
}
