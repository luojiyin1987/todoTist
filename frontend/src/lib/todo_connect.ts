// ConnectRPC Web Client for Todo Service
import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { TodoService } from './todo_pb';

// Define request/response types
interface AddTaskRequest {
  text: string;
}

// eslint-disable-next-line @typescript-eslint/no-empty-object-type
interface GetTasksRequest {
  // Empty request interface
}

interface DeleteTaskRequest {
  id: string;
}

interface AddTaskResponse {
  task?: Task;
}

interface GetTasksResponse {
  tasks: Task[];
}

interface DeleteTaskResponse {
  success: boolean;
}

interface Task {
  id: string;
  text: string;
  createdAt: number;
}

// Export types
export type { AddTaskRequest, GetTasksRequest, DeleteTaskRequest, AddTaskResponse, GetTasksResponse, DeleteTaskResponse, Task };

// Define the TodoService interface
export interface TodoService {
  addTask(request: AddTaskRequest): Promise<AddTaskResponse>;
  getTasks(request: GetTasksRequest): Promise<GetTasksResponse>;
  deleteTask(request: DeleteTaskRequest): Promise<DeleteTaskResponse>;
}

/**
 * Creates a true ConnectRPC client using the official createClient pattern.
 *
 * This implementation uses the proper ConnectRPC client with service definitions,
 * eliminating all fetch calls and providing true ConnectRPC protocol support.
 *
 * @param baseUrl - Base URL of the backend (e.g. "http://localhost:8080")
 * @returns A TodoService client with true ConnectRPC protocol support
 */
export function createTodoService(baseUrl: string): TodoService {
  // Create the ConnectRPC transport
  const transport = createConnectTransport({
    baseUrl,
    useBinaryFormat: false,
  });

  // Create the true ConnectRPC client using service definitions
  const client = createClient(TodoService, transport);

  // Return a typed interface that matches our expected API
  return {
    async addTask(request: AddTaskRequest): Promise<AddTaskResponse> {
      const response = await client.addTask({ text: request.text });
      return {
        task: response.task ? {
          id: response.task.id,
          text: response.task.text,
          createdAt: Number(response.task.createdAt)
        } : undefined
      };
    },

    async getTasks(request: GetTasksRequest): Promise<GetTasksResponse> {
      const response = await client.getTasks(request);
      return {
        tasks: response.tasks.map(task => ({
          id: task.id,
          text: task.text,
          createdAt: Number(task.createdAt)
        }))
      };
    },

    async deleteTask(request: DeleteTaskRequest): Promise<DeleteTaskResponse> {
      const response = await client.deleteTask(request);
      return {
        success: response.success
      };
    },
  };
}

export type TodoServiceClient = ReturnType<typeof createTodoService>;

// Export helper functions for creating requests
export const createRequests = {
  addTask: (text: string) => ({ text }),
  getTasks: () => ({}),
  deleteTask: (id: string) => ({ id }),
};