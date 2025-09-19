// Simplified TypeScript types to match Go backend

// Forward declare interfaces to avoid circular references
export interface ITask {
  id: string;
  text: string;
  createdAt: number;
}

export class Task implements ITask {
  id: string = "";
  text: string = "";
  createdAt: number = 0;
  constructor(data?: { id?: string; text?: string; createdAt?: number }) {
    if (data?.id) this.id = data.id;
    if (data?.text) this.text = data.text;
    if (data?.createdAt !== undefined) this.createdAt = data.createdAt;
  }
}

export class AddTaskRequest {
  text: string = "";
  constructor(data?: { text?: string }) {
    if (data?.text) this.text = data.text;
  }
}

export class AddTaskResponse {
  task?: Task;
  constructor(data?: { task?: Task }) {
    if (data?.task) this.task = new Task(data.task);
  }
}

export class GetTasksRequest {
  constructor() {}
}

export class GetTasksResponse {
  tasks: Task[] = [];
  constructor(data?: { tasks?: Task[] }) {
    if (data?.tasks) this.tasks = data.tasks.map(task => new Task(task));
  }
}

export class DeleteTaskRequest {
  id: string = "";
  constructor(data?: { id?: string }) {
    if (data?.id) this.id = data.id;
  }
}

export class DeleteTaskResponse {
  success: boolean = false;
  constructor(data?: { success?: boolean }) {
    if (data?.success !== undefined) this.success = data.success;
  }
}