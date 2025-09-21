// Simplified ConnectRPC handler implementation
package todov1

import (
	"context"
	"fmt"
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

func writeConnectError(w http.ResponseWriter, err *connect.Error) {
	// Ensure protocol version is always present
	w.Header().Set("Connect-Protocol-Version", "1")
	w.Header().Set("Content-Type", "application/json")
	
	// Map ConnectRPC codes to HTTP status codes
	var statusCode int
	switch err.Code() {
	case connect.CodeInvalidArgument:
		statusCode = http.StatusBadRequest
	case connect.CodeNotFound:
		statusCode = http.StatusNotFound
	case connect.CodeInternal:
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusInternalServerError
	}
	
	w.WriteHeader(statusCode)
	
	// Use simple error response for now - can be enhanced later
	fmt.Fprintf(w, `{"code":"%s","message":"%s"}`, err.Code(), err.Message())
}

func (h *todoServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connect-Protocol-Version", "1")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Connect-Protocol-Version, Connect-Timeout-Ms, Connect-Content-Encoding, Connect-Accept-Encoding, Accept")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Connect-Protocol-Version, Connect-Content-Encoding, Connect-Accept-Encoding, Connect-Error-Code")
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")

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
		writeConnectError(w, connect.NewError(connect.CodeInvalidArgument, err))
		return
	}

	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.AddTask(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// propagate headers from svc
	for k, vv := range resp.Header() {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	data, err = h.pjm.Marshal(resp.Msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		// can't recover after write starts; optionally log
		return
	}
}

func (h *todoServiceHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	req := GetTasksRequest{}
	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.GetTasks(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// propagate headers from svc
	for k, vv := range resp.Header() {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	data, err := h.pjm.Marshal(resp.Msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		// can't recover after write starts; optionally log
		return
	}
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
		writeConnectError(w, connect.NewError(connect.CodeInvalidArgument, err))
		return
	}

	connectReq := connect.NewRequest(&req)
	resp, err := h.svc.DeleteTask(r.Context(), connectReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// propagate headers from svc
	for k, vv := range resp.Header() {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	data, err = h.pjm.Marshal(resp.Msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		// can't recover after write starts; optionally log
		return
	}
}