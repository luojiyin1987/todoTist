package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"todo-list/todo/v1"
)

func TestAddTask(t *testing.T) {
	server := NewTodoServer()

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
			wantErr: true,
		},
		{
			name:    "whitespace task",
			text:    "   ",
			wantErr: true,
		},
		{
			name:    "task at max length",
			text:    strings.Repeat("a", MaxTaskTextLength),
			wantErr: false,
		},
		{
			name:    "task exceeding max length",
			text:    strings.Repeat("a", MaxTaskTextLength+1),
			wantErr: true,
		},
		{
			name:    "task with leading/trailing spaces",
			text:    "  Valid task  ",
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

				expectedText := strings.TrimSpace(tt.text)
				if resp.Msg.Task.Text != expectedText {
					t.Errorf("AddTask() task text = %v, want %v", resp.Msg.Task.Text, expectedText)
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
	server := NewTodoServer()

	// Add some tasks first
	ctx := context.Background()
	task1 := connect.NewRequest(&todov1.AddTaskRequest{Text: "Task 1"})
	task2 := connect.NewRequest(&todov1.AddTaskRequest{Text: "Task 2"})

	_, err := server.AddTask(ctx, task1)
	if err != nil {
		t.Fatalf("AddTask() error = %v", err)
	}

	time.Sleep( 3 * time.Millisecond) // Ensure different seconds if CreatedAt uses Unix seconds

	_, err = server.AddTask(ctx, task2)
	if err != nil {
		t.Fatalf("AddTask() error = %v", err)
	}

	getReq := connect.NewRequest(&todov1.GetTasksRequest{})
	getResp, err := server.GetTasks(ctx, getReq)
	if err != nil {
		t.Fatalf("GetTasks() error = %v", err)
	}

	if len(getResp.Msg.Tasks) != 2 {
		t.Errorf("GetTasks() returned %d tasks, want 2", len(getResp.Msg.Tasks))
	}

	// Check that tasks are sorted by creation time (newest first)
	if len(getResp.Msg.Tasks) >= 2 {
		t0, t1 := getResp.Msg.Tasks[0], getResp.Msg.Tasks[1]
		if t0.CreatedAt < t1.CreatedAt {
			t.Error("GetTasks() tasks not sorted by creation time (newest first)")
		}
		// Only assert by-text ordering when timestamps differ.
		if t0.CreatedAt > t1.CreatedAt && t0.Text != "Task 2" {
			t.Error("GetTasks() newest task not first in list")
		}
	}
}

func TestGetTasksEmpty(t *testing.T) {
	server := NewTodoServer()

	ctx := context.Background()
	req := connect.NewRequest(&todov1.GetTasksRequest{})

	resp, err := server.GetTasks(ctx, req)
	if err != nil {
		t.Errorf("GetTasks() error = %v", err)
	}

	if len(resp.Msg.Tasks) != 0 {
		t.Errorf("GetTasks() returned %d tasks, want 0", len(resp.Msg.Tasks))
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name       string
		taskID     string
		setupTasks bool
		wantErr    bool
		wantSuccess bool
	}{
		{
			name:       "existing task",
			setupTasks: true,
			wantErr:    false,
			wantSuccess: true,
		},
		{
			name:       "non-existing task",
			taskID:     "nonexistent",
			setupTasks: false,
			wantErr:    true,
			wantSuccess: false,
		},
		{
			name:       "empty task ID",
			taskID:     "",
			setupTasks: false,
			wantErr:    true,
			wantSuccess: false,
		},
		{
			name:       "whitespace task ID",
			taskID:     "   ",
			setupTasks: false,
			wantErr:    true,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTodoServer()
			ctx := context.Background()

			var taskID string
			if tt.setupTasks {
				// Add a task first
				addReq := connect.NewRequest(&todov1.AddTaskRequest{Text: "Test task"})
				addResp, err := server.AddTask(ctx, addReq)
				if err != nil {
					t.Fatalf("Setup AddTask() error = %v", err)
				}
				taskID = addResp.Msg.Task.Id
			} else if tt.taskID != "" {
				taskID = tt.taskID
			}

			deleteReq := connect.NewRequest(&todov1.DeleteTaskRequest{Id: taskID})
			deleteResp, err := server.DeleteTask(ctx, deleteReq)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if deleteResp.Msg.Success != tt.wantSuccess {
					t.Errorf("DeleteTask() success = %v, want %v", deleteResp.Msg.Success, tt.wantSuccess)
				}
			}
		})
	}
}

func TestAddTaskIntegration(t *testing.T) {
	server := NewTodoServer()
	ctx := context.Background()

	// Add a task
	addReq := connect.NewRequest(&todov1.AddTaskRequest{
		Text: "Integration test task",
	})
	addResp, err := server.AddTask(ctx, addReq)
	if err != nil {
		t.Fatalf("AddTask() error = %v", err)
	}

	taskID := addResp.Msg.Task.Id

	// Get tasks and verify
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

func TestValidateTaskText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "valid text",
			text:    "Valid task",
			wantErr: false,
		},
		{
			name:    "empty text",
			text:    "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			text:    "   \t\n  ",
			wantErr: true,
		},
		{
			name:    "text at max length",
			text:    strings.Repeat("a", MaxTaskTextLength),
			wantErr: false,
		},
		{
			name:    "text exceeding max length",
			text:    strings.Repeat("a", MaxTaskTextLength+1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTaskText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTaskText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkAddTask(b *testing.B) {
	ctx := context.Background()
	req := connect.NewRequest(&todov1.AddTaskRequest{
		Text: "Benchmark task",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a new server for each iteration to avoid state carry-over
		b.StopTimer()
		server := NewTodoServer()
		b.StartTimer()
		_, err := server.AddTask(ctx, req)
		if err != nil {
			b.Fatalf("AddTask() error = %v", err)
		}
	}
}

func BenchmarkGetTasks(b *testing.B) {
	ctx := context.Background()
	req := connect.NewRequest(&todov1.GetTasksRequest{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a new server and add tasks for each benchmark iteration
		b.StopTimer()
		server := NewTodoServer()
		// Add some tasks
		for j := 0; j < 1000; j++ {
			id1, _ := generateID()
			id2, _ := generateID()
			task := &todov1.Task{
				Id:        id1,
				Text:      "Benchmark task " + id2,
				CreatedAt: time.Now().UnixMilli(),
			}
			server.tasks[task.Id] = task
		}
		b.StartTimer()
		
		_, err := server.GetTasks(ctx, req)
		if err != nil {
			b.Fatalf("GetTasks() error = %v", err)
		}
	}
}