package main

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"todo-list/todo/v1"
)

func TestAddTask(t *testing.T) {
	server := &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}

	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "valid task",
			text:    "Test task",
			wantErr: false,
		},
		{
			name:    "empty task",
			text:    "",
			wantErr: false,
		},
		{
			name:    "whitespace task",
			text:    "   ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := connect.NewRequest(&todov1.AddTaskRequest{
				Text: tt.text,
			})

			resp, err := server.AddTask(ctx, req)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.Msg.Task == nil {
					t.Error("AddTask() returned nil task")
					return
				}

				if resp.Msg.Task.Text != tt.text {
					t.Errorf("AddTask() task text = %v, want %v", resp.Msg.Task.Text, tt.text)
				}

				if resp.Msg.Task.Id == "" {
					t.Error("AddTask() task ID is empty")
				}

				if resp.Msg.Task.CreatedAt == 0 {
					t.Error("AddTask() task CreatedAt is zero")
				}
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	server := &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}

	// Add some test tasks
	task1 := &todov1.Task{
		Id:        "test1",
		Text:      "Task 1",
		CreatedAt: time.Now().Unix(),
	}
	task2 := &todov1.Task{
		Id:        "test2",
		Text:      "Task 2",
		CreatedAt: time.Now().Unix(),
	}
	server.tasks["test1"] = task1
	server.tasks["test2"] = task2

	ctx := context.Background()
	req := connect.NewRequest(&todov1.GetTasksRequest{})

	resp, err := server.GetTasks(ctx, req)
	if err != nil {
		t.Fatalf("GetTasks() error = %v", err)
	}

	if len(resp.Msg.Tasks) != 2 {
		t.Errorf("GetTasks() returned %d tasks, want 2", len(resp.Msg.Tasks))
	}

	// Verify tasks are returned correctly
	foundTask1 := false
	foundTask2 := false
	for _, task := range resp.Msg.Tasks {
		if task.Id == "test1" && task.Text == "Task 1" {
			foundTask1 = true
		}
		if task.Id == "test2" && task.Text == "Task 2" {
			foundTask2 = true
		}
	}

	if !foundTask1 {
		t.Error("GetTasks() did not return task1")
	}
	if !foundTask2 {
		t.Error("GetTasks() did not return task2")
	}
}

func TestGetTasksEmpty(t *testing.T) {
	server := &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}

	ctx := context.Background()
	req := connect.NewRequest(&todov1.GetTasksRequest{})

	resp, err := server.GetTasks(ctx, req)
	if err != nil {
		t.Fatalf("GetTasks() error = %v", err)
	}

	if len(resp.Msg.Tasks) != 0 {
		t.Errorf("GetTasks() returned %d tasks, want 0", len(resp.Msg.Tasks))
	}
}

func TestDeleteTask(t *testing.T) {
	server := &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}

	// Add a test task
	task := &todov1.Task{
		Id:        "test-delete",
		Text:      "Task to delete",
		CreatedAt: time.Now().Unix(),
	}
	server.tasks["test-delete"] = task

	tests := []struct {
		name        string
		taskID      string
		wantSuccess bool
		wantErr     bool
	}{
		{
			name:        "existing task",
			taskID:      "test-delete",
			wantSuccess: true,
			wantErr:     false,
		},
		{
			name:        "non-existing task",
			taskID:      "non-existent",
			wantSuccess: false,
			wantErr:     false,
		},
		{
			name:        "empty task ID",
			taskID:      "",
			wantSuccess: false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := connect.NewRequest(&todov1.DeleteTaskRequest{
				Id: tt.taskID,
			})

			resp, err := server.DeleteTask(ctx, req)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resp.Msg.Success != tt.wantSuccess {
				t.Errorf("DeleteTask() success = %v, want %v", resp.Msg.Success, tt.wantSuccess)
			}

			// Verify task is actually deleted from the map
			if tt.wantSuccess {
				if _, exists := server.tasks[tt.taskID]; exists {
					t.Error("DeleteTask() did not remove task from map")
				}
			}
		})
	}
}

func TestAddTaskIntegration(t *testing.T) {
	server := &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}

	ctx := context.Background()

	// Add a task
	addReq := connect.NewRequest(&todov1.AddTaskRequest{
		Text: "Integration test task",
	})
	addResp, err := server.AddTask(ctx, addReq)
	if err != nil {
		t.Fatalf("AddTask() error = %v", err)
	}

	if addResp.Msg.Task == nil {
		t.Fatal("AddTask() returned nil task")
	}

	taskID := addResp.Msg.Task.Id

	// Get tasks to verify it was added
	getReq := connect.NewRequest(&todov1.GetTasksRequest{})
	getResp, err := server.GetTasks(ctx, getReq)
	if err != nil {
		t.Fatalf("GetTasks() error = %v", err)
	}

	if len(getResp.Msg.Tasks) != 1 {
		t.Errorf("GetTasks() returned %d tasks, want 1", len(getResp.Msg.Tasks))
	}

	if getResp.Msg.Tasks[0].Id != taskID {
		t.Errorf("GetTasks() returned task with ID %v, want %v", getResp.Msg.Tasks[0].Id, taskID)
	}

	// Delete the task
	deleteReq := connect.NewRequest(&todov1.DeleteTaskRequest{
		Id: taskID,
	})
	deleteResp, err := server.DeleteTask(ctx, deleteReq)
	if err != nil {
		t.Fatalf("DeleteTask() error = %v", err)
	}

	if !deleteResp.Msg.Success {
		t.Error("DeleteTask() returned success = false, want true")
	}

	// Verify task is gone
	getResp2, err := server.GetTasks(ctx, getReq)
	if err != nil {
		t.Fatalf("GetTasks() error = %v", err)
	}

	if len(getResp2.Msg.Tasks) != 0 {
		t.Errorf("GetTasks() returned %d tasks after deletion, want 0", len(getResp2.Msg.Tasks))
	}
}

func BenchmarkAddTask(b *testing.B) {
	server := &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}

	ctx := context.Background()
	req := connect.NewRequest(&todov1.AddTaskRequest{
		Text: "Benchmark task",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := server.AddTask(ctx, req)
		if err != nil {
			b.Fatalf("AddTask() error = %v", err)
		}
	}
}

func BenchmarkGetTasks(b *testing.B) {
	server := &TodoServer{
		tasks: make(map[string]*todov1.Task),
	}

	// Add some tasks
	for i := 0; i < 1000; i++ {
		task := &todov1.Task{
			Id:        generateID(),
			Text:      "Benchmark task " + generateID(),
			CreatedAt: time.Now().Unix(),
		}
		server.tasks[task.Id] = task
	}

	ctx := context.Background()
	req := connect.NewRequest(&todov1.GetTasksRequest{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := server.GetTasks(ctx, req)
		if err != nil {
			b.Fatalf("GetTasks() error = %v", err)
		}
	}
}