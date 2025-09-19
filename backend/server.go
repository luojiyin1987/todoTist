package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2/h2c"
	http2 "golang.org/x/net/http2"

	"todo-list/todo/v1"
)

const (
	MaxTaskTextLength = 500
	MinTaskTextLength = 1
)

var (
	ErrTaskTextEmpty    = errors.New("task text cannot be empty")
	ErrTaskTextTooLong  = errors.New("task text exceeds maximum length")
	ErrTaskNotFound     = errors.New("task not found")
	ErrInvalidTaskID    = errors.New("invalid task ID")
)

type TodoServer struct {
	mu     sync.RWMutex
	tasks  map[string]*todov1.Task
}

// NewTodoServer returns a pointer to a TodoServer with its tasks map initialized.
// The returned server is ready for use; its zero-value sync.RWMutex is valid for guarding access to the tasks map.
func NewTodoServer() *TodoServer {
	return &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}
}

// validateTaskText trims leading and trailing whitespace and checks that the
// remaining text length is within allowed bounds.
//
// It returns ErrTaskTextEmpty if the trimmed text is shorter than MinTaskTextLength,
// ErrTaskTextTooLong if it exceeds MaxTaskTextLength, or nil if the text is valid.
func validateTaskText(text string) error {
	text = strings.TrimSpace(text)
	if len(text) < MinTaskTextLength {
		return ErrTaskTextEmpty
	}
	if len(text) > MaxTaskTextLength {
		return ErrTaskTextTooLong
	}
	return nil
}

func (s *TodoServer) AddTask(
	ctx context.Context,
	req *connect.Request[todov1.AddTaskRequest],
) (*connect.Response[todov1.AddTaskResponse], error) {
	if err := validateTaskText(req.Msg.Text); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := &todov1.Task{
		Id:        generateID(),
		Text:      strings.TrimSpace(req.Msg.Text),
		CreatedAt: time.Now().Unix(),
	}

	s.tasks[task.Id] = task

	return connect.NewResponse(&todov1.AddTaskResponse{
		Task: task,
	}), nil
}

func (s *TodoServer) GetTasks(
	ctx context.Context,
	req *connect.Request[todov1.GetTasksRequest],
) (*connect.Response[todov1.GetTasksResponse], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*todov1.Task
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	// Sort tasks by creation time (newest first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt > tasks[j].CreatedAt
	}

	return connect.NewResponse(&todov1.GetTasksResponse{
		Tasks: tasks,
	}), nil
}

func (s *TodoServer) DeleteTask(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTaskRequest],
) (*connect.Response[todov1.DeleteTaskResponse], error) {
	if strings.TrimSpace(req.Msg.Id) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrInvalidTaskID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[req.Msg.Id]; !exists {
		return nil, connect.NewError(connect.CodeNotFound, ErrTaskNotFound)
	}

	delete(s.tasks, req.Msg.Id)
	return connect.NewResponse(&todov1.DeleteTaskResponse{
		Success: true,
	}), nil
}

func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	todoServer := &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}

	mux := http.NewServeMux()
	_, handler := todov1.NewTodoServiceHandler(todoServer)
	mux.Handle("/", handler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Connect-Protocol-Version"},
		AllowCredentials: true,
	})

	finalHandler := corsHandler.Handler(h2c.NewHandler(mux, &http2.Server{}))

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", finalHandler))
}