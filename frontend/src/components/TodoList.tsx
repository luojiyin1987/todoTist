'use client';

import { useState, useEffect } from 'react';
import { createTodoService } from '@/lib/todo_connect';
import { AddTaskRequest, GetTasksRequest, DeleteTaskRequest, Task } from '@/lib/todo_pb';

interface TodoItem {
  id: string;
  text: string;
  createdAt: number;
}

export default function TodoList() {
  const client = createTodoService('http://localhost:8080');
  const [tasks, setTasks] = useState<TodoItem[]>([]);
  const [newTask, setNewTask] = useState('');
  const [loading, setLoading] = useState(false);

  const fetchTasks = async () => {
    try {
      const response = await client.getTasks(new GetTasksRequest());
      const todoItems = response.tasks.map(task => ({
        id: task.id,
        text: task.text,
        createdAt: task.createdAt,
      }));
      setTasks(todoItems);
    } catch (error) {
      console.error('Error fetching tasks:', error);
    }
  };

  const addTask = async () => {
    if (!newTask.trim()) return;

    setLoading(true);
    try {
      const request = new AddTaskRequest({ text: newTask.trim() });
      await client.addTask(request);
      setNewTask('');
      await fetchTasks();
    } catch (error) {
      console.error('Error adding task:', error);
    } finally {
      setLoading(false);
    }
  };

  const deleteTask = async (id: string) => {
    try {
      const request = new DeleteTaskRequest({ id });
      await client.deleteTask(request);
      await fetchTasks();
    } catch (error) {
      console.error('Error deleting task:', error);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, []);

  const formatDate = (timestamp: number) => {
    return new Date(timestamp * 1000).toLocaleString();
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6 text-gray-900">Todo List</h1>
      
      <div className="mb-6 flex gap-2">
        <input
          type="text"
          value={newTask}
          onChange={(e) => setNewTask(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && addTask()}
          placeholder="Add a new task..."
          className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          disabled={loading}
        />
        <button
          onClick={addTask}
          disabled={loading || !newTask.trim()}
          className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {loading ? 'Adding...' : 'Add Task'}
        </button>
      </div>

      <div className="space-y-3">
        {tasks.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            No tasks yet. Add one above!
          </div>
        ) : (
          tasks.map((task) => (
            <div
              key={task.id}
              className="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="flex-1">
                <div className="text-gray-900 font-medium">{task.text}</div>
                <div className="text-sm text-gray-500">
                  {formatDate(task.createdAt)}
                </div>
              </div>
              <button
                onClick={() => deleteTask(task.id)}
                className="ml-4 px-3 py-1 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors"
              >
                Delete
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}