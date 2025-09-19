package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2/h2c"
	http2 "golang.org/x/net/http2"

	"todo-list/todo/v1"
)

type TodoServer struct {
	mu     sync.Mutex
	tasks  map[string]*todov1.Task
}

func NewTodoServer() *TodoServer {
	return &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}
}

func (s *TodoServer) AddTask(
	ctx context.Context,
	req *connect.Request[todov1.AddTaskRequest],
) (*connect.Response[todov1.AddTaskResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &todov1.Task{
		Id:        generateID(),
		Text:      req.Msg.Text,
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
	s.mu.Lock()
	defer s.mu.Unlock()

	var tasks []*todov1.Task
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return connect.NewResponse(&todov1.GetTasksResponse{
		Tasks: tasks,
	}), nil
}

func (s *TodoServer) DeleteTask(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTaskRequest],
) (*connect.Response[todov1.DeleteTaskResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[req.Msg.Id]; exists {
		delete(s.tasks, req.Msg.Id)
		return connect.NewResponse(&todov1.DeleteTaskResponse{
			Success: true,
		}), nil
	}

	return connect.NewResponse(&todov1.DeleteTaskResponse{
		Success: false,
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