package todov1

type AddTaskRequest struct {
	Text string
}

type AddTaskResponse struct {
	Task *Task
}

type GetTasksRequest struct {
}

type GetTasksResponse struct {
	Tasks []*Task
}

type DeleteTaskRequest struct {
	Id string
}

type DeleteTaskResponse struct {
	Success bool
}

type Task struct {
	Id        string
	Text      string
	CreatedAt int64
}