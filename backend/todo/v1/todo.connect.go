// Simplified ConnectRPC handler implementation
package todov1

import (
	"context"
	"io"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
)

type TodoServiceHandler interface {
	AddTask(context.Context, *connect.Request[AddTaskRequest]) (*connect.Response[AddTaskResponse], error)
	GetTasks(context.Context, *connect.Request[GetTasksRequest]) (*connect.Response[GetTasksResponse], error)
	DeleteTask(context.Context, *connect.Request[DeleteTaskRequest]) (*connect.Response[DeleteTaskResponse], error)
}

const TodoServiceName = "todo.v1.TodoService"

func NewTodoServiceHandler(svc TodoServiceHandler) (string, http.Handler) {
	return "/" + TodoServiceName + "/", &todoServiceHandler{
		svc: svc,
		pjm: protojson.MarshalOptions{},
		pju: protojson.UnmarshalOptions{},
	}
}

type todoServiceHandler struct {
	svc  TodoServiceHandler
	pjm  protojson.MarshalOptions
	pju  protojson.UnmarshalOptions
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
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB limit
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req AddTaskRequest
	if err := h.pju.Unmarshal(data, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.AddTask(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err = h.pjm.Marshal(resp.Msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *todoServiceHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	req := GetTasksRequest{}
	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.GetTasks(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := h.pjm.Marshal(resp.Msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *todoServiceHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB limit
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req DeleteTaskRequest
	if err := h.pju.Unmarshal(data, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.DeleteTask(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err = h.pjm.Marshal(resp.Msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}