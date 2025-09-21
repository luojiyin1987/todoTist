// ConnectRPC Web Client for Todo Service
import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { create } from '@bufbuild/protobuf';
import {
  AddTaskRequest,
  AddTaskResponse,
  GetTasksRequest,
  GetTasksResponse,
  DeleteTaskRequest,
  DeleteTaskResponse,
  Task,
  TodoService,
  AddTaskRequestSchema,
  GetTasksRequestSchema,
  DeleteTaskRequestSchema,
} from './todo_pb';

// Re-export generated types for convenience
export type {
  AddTaskRequest,
  AddTaskResponse,
  GetTasksRequest,
  GetTasksResponse,
  DeleteTaskRequest,
  DeleteTaskResponse,
  Task,
};

// Define application-level types derived from generated types
// This provides cleaner interfaces for React components while maintaining type safety
export type AppTask = {
  id: string;
  text: string;
  createdAt: number; // Convert bigint to number for easier use in React
};

export type AppAddTaskResponse = {
  task?: AppTask;
};

export type AppGetTasksResponse = {
  tasks: AppTask[];
};

export type AppDeleteTaskResponse = {
  success: boolean;
};

// Define the TodoService interface using generated types
export interface TodoService {
  addTask(request: AddTaskRequest): Promise<AppAddTaskResponse>;
  getTasks(request: GetTasksRequest): Promise<AppGetTasksResponse>;
  deleteTask(request: DeleteTaskRequest): Promise<AppDeleteTaskResponse>;
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

  // Helper function to convert Task to AppTask
  const toAppTask = (task: Task): AppTask => ({
    id: task.id,
    text: task.text,
    createdAt: Number(task.createdAt),
  });

  // Return a typed interface that matches our expected API
  return {
    async addTask(request: AddTaskRequest): Promise<AppAddTaskResponse> {
      const response = await client.addTask(request);
      return {
        task: response.task ? toAppTask(response.task) : undefined,
      };
    },

    async getTasks(request: GetTasksRequest): Promise<AppGetTasksResponse> {
      const response = await client.getTasks(request);
      return {
        tasks: response.tasks.map(toAppTask),
      };
    },

    async deleteTask(request: DeleteTaskRequest): Promise<AppDeleteTaskResponse> {
      const response = await client.deleteTask(request);
      return {
        success: response.success,
      };
    },
  };
}

export type TodoServiceClient = ReturnType<typeof createTodoService>;

// Export helper functions for creating requests using generated types
export const createRequests = {
  addTask: (text: string): AddTaskRequest => create(AddTaskRequestSchema, { text }),
  getTasks: (): GetTasksRequest => create(GetTasksRequestSchema, {}),
  deleteTask: (id: string): DeleteTaskRequest => create(DeleteTaskRequestSchema, { id }),
};