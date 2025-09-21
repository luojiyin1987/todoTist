'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { createTodoService, createRequests } from '@/lib/todo_connect';

interface TodoItem {
  id: string;
  text: string;
  createdAt: number;
}

/**
 * TodoList React component — client UI for managing todo tasks.
 *
 * Renders an add-task form and a list of tasks backed by a Todo service at http://localhost:8080.
 * Handles fetching, adding, and deleting tasks via the service client, and surfaces user-facing
 * error messages and a loading state. Input is validated (non-empty, max 500 characters) before
 * creating a task. Tasks are displayed newest-first and each item shows its creation timestamp.
 *
 * @returns The TodoList component as JSX.
 */
export default function TodoList() {
  const client = useMemo(() => createTodoService('http://localhost:8080'), []);
  const [tasks, setTasks] = useState<TodoItem[]>([]);
  const [newTask, setNewTask] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');

  const fetchTasks = useCallback(async () => {
    try {
      setError('');
      const response = await client.getTasks(createRequests.getTasks());
      const todoItems = (response.tasks || []).map((task: { id: string; text: string; createdAt: number }) => ({
        id: task.id,
        text: task.text,
        createdAt: task.createdAt,
      }));
      setTasks(todoItems);
    } catch (err) {
      console.error('Error fetching tasks:', err);
      setError('Failed to fetch tasks. Please try again.');
    }
  }, [client]);

  const addTask = async () => {
    const taskText = newTask.trim();
    if (!taskText) {
      setError('Task text cannot be empty.');
      return;
    }

    if (taskText.length > 500) {
      setError('Task text cannot exceed 500 characters.');
      return;
    }

    setLoading(true);
    setError('');
    try {
      const request = createRequests.addTask(taskText);
      await client.addTask(request);
      setNewTask('');
      await fetchTasks();
    } catch (err) {
      console.error('Error adding task:', err);
      setError('Failed to add task. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const deleteTask = async (id: string) => {
    // Validate ID before sending request
    if (!id || id.trim() === '') {
      setError('Invalid task ID: Cannot delete task with empty ID.');
      return;
    }
    
    try {
      setError('');
      console.log('Deleting task with ID:', id);
      const request = createRequests.deleteTask(id);
      await client.deleteTask(request);
      await fetchTasks();
    } catch (err) {
      console.error('Error deleting task:', err);
      if (err instanceof Error) {
        if (err.message === 'Task not found') {
          setError('Task not found. It may have been already deleted.');
        } else if (err.message === 'Invalid task ID') {
          setError('Invalid task ID. Please refresh the page and try again.');
        } else {
          setError(`Failed to delete task: ${err.message}`);
        }
      } else {
        setError('Failed to delete task. Please try again.');
      }
    }
  };

  const handleDelete = (id: string, taskText: string) => {
    if (window.confirm(`Are you sure you want to delete this task?\n"${taskText}"`)) {
      deleteTask(id);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  const formatDate = (timestamp: number) => {
    return new Date(timestamp * 1000).toLocaleString();
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !loading) {
      addTask();
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6 text-gray-900">TodoTist - Your Personal Task Manager</h1>
      
      {/* Error Message */}
      {error && (
        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded-md">
          {error}
        </div>
      )}
      
      {/* Add Task Form */}
      <div className="mb-6">
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={newTask}
            onChange={(e) => setNewTask(e.target.value)}
            onKeyDown={handleKeyPress}
            placeholder="Add a new task..."
            className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
            disabled={loading}
            maxLength={500}
          />
          <button
            onClick={addTask}
            disabled={loading || !newTask.trim()}
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {loading ? 'Adding...' : 'Add Task'}
          </button>
        </div>
        <div className="text-sm text-gray-500">
          {newTask.length}/500 characters
        </div>
      </div>

      {/* Tasks List */}
      <div className="space-y-3">
        {tasks.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <div className="text-lg mb-2">No tasks yet!</div>
            <div className="text-sm">Add your first task above to get started.</div>
          </div>
        ) : (
          <>
            <div className="text-sm text-gray-600 mb-3">
              {tasks.length} task{tasks.length !== 1 ? 's' : ''} • Sorted by newest first
            </div>
            {tasks.map((task) => (
              <div
                key={task.id}
                className="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="flex-1 min-w-0">
                  <div className="text-gray-900 font-medium break-words">{task.text}</div>
                  <div className="text-sm text-gray-500 mt-1">
                    Created: {formatDate(task.createdAt)}
                  </div>
                </div>
                <button
                  onClick={() => handleDelete(task.id, task.text)}
                  className="ml-4 px-3 py-1 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors flex-shrink-0"
                  title="Delete task"
                >
                  Delete
                </button>
              </div>
            ))}
          </>
        )}
      </div>
    </div>
  );
}