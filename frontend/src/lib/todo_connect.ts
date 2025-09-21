// ConnectRPC Web Client for Todo Service
import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { create } from '@bufbuild/protobuf';
import {
  AddTaskRequest,
  GetTasksRequest,
  DeleteTaskRequest,
  DeleteTaskResponse,
  Task,
  TodoService as TodoServiceDef,
  AddTaskRequestSchema,
  GetTasksRequestSchema,
  DeleteTaskRequestSchema,
} from './todo_pb';

// Re-export generated types for convenience
export type {
  AddTaskRequest,
  GetTasksRequest,
  DeleteTaskRequest,
  DeleteTaskResponse,
  Task,
};

// Define application-level types derived from generated types
// This provides cleaner interfaces for React components while maintaining type safety
// Consider updating this type if the Protocol Buffer definition changes
export type AppTask = {
  id: string;
  text: string;
  createdAt: number; // Convert bigint to number for easier use in React
};

// Define the TodoClient interface using application-level types
export interface TodoClient {
  addTask(request: AddTaskRequest): Promise<{
    task?: AppTask;
  }>;
  getTasks(request: GetTasksRequest): Promise<{
    tasks: AppTask[];
  }>;
  deleteTask(request: DeleteTaskRequest): Promise<{
    success: boolean;
  }>;
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
export function createTodoService(baseUrl: string): TodoClient {
  // Create the ConnectRPC transport
  const transport = createConnectTransport({
    baseUrl,
    useBinaryFormat: false,
  });

  // Create the true ConnectRPC client using service definitions
  const client = createClient(TodoServiceDef, transport);

  // Helper function to convert Task to AppTask
  const toAppTask = (task: Task): AppTask => {
    const n = Number(task.createdAt);
    if (!Number.isSafeInteger(n)) {
      throw new Error('Task.createdAt exceeds Number.MAX_SAFE_INTEGER; confirm units (expected seconds or ms).');
    }
    return { id: task.id, text: task.text, createdAt: n };
  };

  // Return a typed interface that matches our expected API
  return {
    async addTask(request: AddTaskRequest) {
      const response = await client.addTask(request);
      return {
        task: response.task ? toAppTask(response.task) : undefined,
      };
    },

    async getTasks(request: GetTasksRequest) {
      const response = await client.getTasks(request);
      return {
        tasks: response.tasks.map(toAppTask),
      };
    },

    async deleteTask(request: DeleteTaskRequest) {
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
  addTask: (text: string): AddTaskRequest => {
    const t = text.trim();

    if (t.length === 0) {
      throw new Error('Task text cannot be empty');
    }
    if (t.length > 500) {
      throw new Error('Task text cannot exceed 500 characters');
    }
    return create(AddTaskRequestSchema, { text: t });
  },
  getTasks: (): GetTasksRequest => create(GetTasksRequestSchema, {}),
  deleteTask: (id: string): DeleteTaskRequest => {
    if (!id || id.trim() === '') {
      throw new Error('Task ID cannot be empty');
    }
    return create(DeleteTaskRequestSchema, { id: id.trim() });
  },
};