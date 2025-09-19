// Simplified ConnectRPC client to match Go backend
import { AddTaskRequest, AddTaskResponse, GetTasksRequest, GetTasksResponse, DeleteTaskRequest, DeleteTaskResponse } from "./todo_pb";

interface TaskData {
  id: string;
  text: string;
  createdAt: number;
}

interface ApiResponse {
  task?: TaskData;
  tasks?: TaskData[];
  success?: boolean;
}

export interface TodoService {
  addTask(request: AddTaskRequest): Promise<AddTaskResponse>;
  getTasks(request: GetTasksRequest): Promise<GetTasksResponse>;
  deleteTask(request: DeleteTaskRequest): Promise<DeleteTaskResponse>;
}

export const TodoServiceName = "todo.v1.TodoService";

/**
 * Creates a TodoService implementation that talks to a backend over HTTP.
 *
 * The returned service implements addTask, getTasks, and deleteTask by calling
 * the corresponding HTTP endpoints under the provided baseUrl. Responses are
 * parsed as JSON with unified camelCase field names between backend and frontend.
 *
 * @param baseUrl - Base URL of the backend (e.g. "https://api.example.com"); used as the prefix for service endpoints.
 * @returns A TodoService whose methods perform HTTP requests to the backend and return the corresponding response message objects.
 *
 * Note: HTTP/network or JSON parsing errors from fetch will propagate as rejected promises.
 */
export function createTodoService(baseUrl: string): TodoService {
  return {
    async addTask(request: AddTaskRequest): Promise<AddTaskResponse> {
      const response = await fetch(`${baseUrl}/todo.v1.TodoService/AddTask`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });
      
      if (!response.ok) {
        // Handle ConnectRPC error responses (plain text)
        const errorText = await response.text();
        if (errorText.includes('invalid task text')) {
          throw new Error('Invalid task text');
        } else {
          throw new Error(`HTTP error! status: ${response.status}: ${errorText}`);
        }
      }
      
      const data = await response.json() as ApiResponse;
      return new AddTaskResponse(data);
    },

    async getTasks(_: GetTasksRequest): Promise<GetTasksResponse> {
      const response = await fetch(`${baseUrl}/todo.v1.TodoService/GetTasks`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      
      if (!response.ok) {
        // Handle ConnectRPC error responses (plain text)
        const raw = await response.text();
        let msg = raw;
        try{ const j = JSON.parse(raw); msg = j.message ?? j.error ?? raw; } catch {}
        throw new Error(`HTTP error! status: ${response.status}: ${msg}`);
      }
      
      const data = await response.json() as ApiResponse;
      return new GetTasksResponse(data);
    },

    async deleteTask(request: DeleteTaskRequest): Promise<DeleteTaskResponse> {
      const response = await fetch(`${baseUrl}/todo.v1.TodoService/DeleteTask`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });
      
      if (!response.ok) {
        // Handle ConnectRPC error responses (plain text)
        const raw = await response.text();
        let msg = raw;

        try { const j = JSON.parse(raw); msg = j.message ?? j.error ?? raw; } catch {}
        if (response.status === 404) throw new Error('Task not found');
        if (response.status === 400) throw new Error('Invalid task ID');
        throw new Error(`HTTP ${response.status}: ${msg}`);
      }
      
      const data = await response.json() as ApiResponse;
      return new DeleteTaskResponse(data);
    },
  };
} /** Note: HTTP/network or JSON parsing errors from fetch will propagate as rejected */