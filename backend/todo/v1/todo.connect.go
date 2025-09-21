// Simplified ConnectRPC handler implementation
package todov1

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"connectrpc.com/connect"
)

type TodoServiceHandler interface {
	AddTask(context.Context, *connect.Request[AddTaskRequest]) (*connect.Response[AddTaskResponse], error)
	GetTasks(context.Context, *connect.Request[GetTasksRequest]) (*connect.Response[GetTasksResponse], error)
	DeleteTask(context.Context, *connect.Request[DeleteTaskRequest]) (*connect.Response[DeleteTaskResponse], error)
}

const TodoServiceName = "todo.v1.TodoService"

func NewTodoServiceHandler(svc TodoServiceHandler) (string, http.Handler) {
	return "/" + TodoServiceName, &todoServiceHandler{svc: svc}
}

type todoServiceHandler struct {
	svc TodoServiceHandler
}

func (h *todoServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms, Accept, Accept-Encoding")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract the method from the URL path
	// ConnectRPC client sends requests to /todo.v1.TodoService/MethodName
	path := r.URL.Path
	if !strings.HasPrefix(path, "/todo.v1.TodoService/") {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	methodName := strings.TrimPrefix(path, "/todo.v1.TodoService/")
	
	switch methodName {
	case "AddTask":
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleAddTask(w, r)
	case "GetTasks":
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleGetTasks(w, r)
	case "DeleteTask":
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleDeleteTask(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *todoServiceHandler) handleAddTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req AddTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.AddTask(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp.Msg)
}

func (h *todoServiceHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	req := GetTasksRequest{}
	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.GetTasks(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp.Msg)
}

func (h *todoServiceHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req DeleteTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.DeleteTask(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp.Msg)
}