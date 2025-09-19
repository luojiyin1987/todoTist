package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
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

	var task *todov1.Task
	// Try to generate a unique ID (retry on collision)
	for i := 0; i < 10; i++ {
		task = &todov1.Task{
			Id:        generateID(),
			Text:      strings.TrimSpace(req.Msg.Text),
			CreatedAt: time.Now().Unix(),
		}
		
		if _, exists := s.tasks[task.Id]; !exists {
			s.tasks[task.Id] = task
			return connect.NewResponse(&todov1.AddTaskResponse{
				Task: task,
			}), nil
		}
	}
	
	// If we get here, we couldn't generate a unique ID after 10 attempts
	return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate unique task ID"))
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
	})

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
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func main() {
	todoServer := NewTodoServer()
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Println("Server starting on :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	fmt.Println("\nShutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}