package todov1

type AddTaskRequest struct {
	Text string `json:"text"`
}

type AddTaskResponse struct {
	Task *Task `json:"task"`
}

type GetTasksRequest struct {
}

type GetTasksResponse struct {
	Tasks []*Task `json:"tasks"`
}

type DeleteTaskRequest struct {
	Id string `json:"id"`
}

type DeleteTaskResponse struct {
	Success bool `json:"success"`
}

type Task struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"createdAt"`
}