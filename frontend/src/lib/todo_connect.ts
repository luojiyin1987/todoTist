// Simplified ConnectRPC client to match Go backend
import { AddTaskRequest, AddTaskResponse, GetTasksRequest, GetTasksResponse, DeleteTaskRequest, DeleteTaskResponse } from "./todo_pb";

interface TaskData {
  Id: string;
  Text: string;
  CreatedAt: number;
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
        task: data.Task ? {
          id: data.Task.Id,
          text: data.Task.Text,
          createdAt: data.Task.CreatedAt
        } : undefined
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
        tasks: data.Tasks?.map((task: TaskData) => ({
          id: task.Id,
          text: task.Text,
          createdAt: task.CreatedAt
        })) || []
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
}