import { useState, useEffect, useCallback } from 'react';
import Dashboard from './components/Dashboard.jsx';
import TaskForm from './components/TaskForm.jsx';
import TaskList from './components/TaskList.jsx';
import StatusFilter from './components/StatusFilter.jsx';
import * as api from './services/api.js';

/**
 * App — the root component.
 *
 * It manages the global task state and passes data + callbacks
 * down to the child components.
 *
 * Data flow:
 *   App  ──(tasks)──▸  Dashboard / TaskList
 *   App  ◂──(events)──  TaskForm / TaskCard
 */
export default function App() {
  // ── State ───────────────────────────────────────────────
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [editingTask, setEditingTask] = useState(null);
  const [filter, setFilter] = useState('all');

  // ── Fetch all tasks ─────────────────────────────────────
  const fetchTasks = useCallback(async () => {
    try {
      setLoading(true);
      setError('');
      const data = await api.getTasks();
      setTasks(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  // Load tasks when the component first mounts.
  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  // ── Create or Update ────────────────────────────────────
  const handleSubmit = async (taskData) => {
    if (editingTask) {
      // Update existing task.
      await api.updateTask(editingTask.id, taskData);
      setEditingTask(null);
    } else {
      // Create new task.
      await api.createTask(taskData);
    }
    // Re-fetch the list so the UI is always in sync.
    await fetchTasks();
  };

  // ── Delete ──────────────────────────────────────────────
  const handleDelete = async (id) => {
    await api.deleteTask(id);
    await fetchTasks();
  };

  // ── Status change ───────────────────────────────────────
  const handleStatusChange = async (id, newStatus) => {
    await api.updateTask(id, { status: newStatus });
    await fetchTasks();
  };

  // ── Edit mode ───────────────────────────────────────────
  const handleEdit = (task) => {
    setEditingTask(task);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleCancelEdit = () => {
    setEditingTask(null);
  };

  // ── Filtering ───────────────────────────────────────────
  const filteredTasks =
    filter === 'all'
      ? tasks
      : tasks.filter((t) => t.status === filter);

  // ── Render ──────────────────────────────────────────────
  return (
    <div className="app">
      <header className="app-header">
        <h1>📋 TaskFlow</h1>
        <p>A simple task manager — built to learn DevOps</p>
      </header>

      <Dashboard tasks={tasks} />

      <TaskForm
        onSubmit={handleSubmit}
        editingTask={editingTask}
        onCancel={handleCancelEdit}
      />

      <StatusFilter current={filter} onChange={setFilter} />

      <TaskList
        tasks={filteredTasks}
        loading={loading}
        error={error}
        onEdit={handleEdit}
        onDelete={handleDelete}
        onStatusChange={handleStatusChange}
      />
    </div>
  );
}

