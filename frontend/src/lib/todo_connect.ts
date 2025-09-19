// Simplified ConnectRPC client to match Go backend
import { AddTaskRequest, AddTaskResponse, GetTasksRequest, GetTasksResponse, DeleteTaskRequest, DeleteTaskResponse } from "./todo_pb";

interface TaskData {
  id: string;
  text: string;
  createdAt: number;
}

interface ApiResponse {
  Task?: TaskData;
  Tasks?: TaskData[];
  Success?: boolean;
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
 * parsed as JSON and mapped from the backend's PascalCase fields to the
 * client's camelCase shape.
 *
 * @param baseUrl - Base URL of the backend (e.g. "https://api.example.com"); used as the prefix for service endpoints.
 * @returns A TodoService whose methods perform HTTP requests to the backend and return the corresponding response message objects.
 *
 * Note: HTTP/network or JSON parsing errors from fetch will propagate as rejected promises.
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
      const data = await response.json() as ApiResponse;
      // Convert PascalCase to camelCase for TypeScript compatibility
      const convertedData = {
        task: data.Task 
      };
      return new AddTaskResponse(convertedData);
    },

    async getTasks(): Promise<GetTasksResponse> {
      const response = await fetch(`${baseUrl}/todo.v1.TodoService/GetTasks`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      const data = await response.json() as ApiResponse;
      // Convert PascalCase to camelCase for TypeScript compatibility
      const convertedData = {
        tasks: data.Tasks
      };
      return new GetTasksResponse(convertedData);
    },

    async deleteTask(request: DeleteTaskRequest): Promise<DeleteTaskResponse> {
      const response = await fetch(`${baseUrl}/todo.v1.TodoService/DeleteTask`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });
      const data = await response.json() as ApiResponse;
      // Convert PascalCase to camelCase for TypeScript compatibility
      const convertedData = {
        success: data.Success || false
      };
      return new DeleteTaskResponse(convertedData);
    },
  };
} /** Note: HTTP/network or JSON parsing errors from fetch will propagate as rejected