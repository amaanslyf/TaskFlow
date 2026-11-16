import { useState, useEffect } from 'react';

/**
 * TaskForm — handles both "create new task" and "edit
 * existing task".
 *
 * When `editingTask` is provided the form pre-fills with
 * that task's data and the submit button says "Update".
 */
export default function TaskForm({ onSubmit, editingTask, onCancel }) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState('todo');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  // When the parent passes a task to edit, fill the form.
  useEffect(() => {
    if (editingTask) {
      setTitle(editingTask.title);
      setDescription(editingTask.description || '');
      setStatus(editingTask.status);
    } else {
      setTitle('');
      setDescription('');
      setStatus('todo');
    }
    setError('');
  }, [editingTask]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    const trimmedTitle = title.trim();
    if (!trimmedTitle) {
      setError('Title is required');
      return;
    }

    setSubmitting(true);
    setError('');

    try {
      await onSubmit({
        title: trimmedTitle,
        description: description.trim(),
        status,
      });
      // Clear form on success (only for create; edit is
      // handled by the parent clearing editingTask).
      if (!editingTask) {
        setTitle('');
        setDescription('');
        setStatus('todo');
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form className="task-form" onSubmit={handleSubmit}>
      <h2>{editingTask ? 'Edit Task' : 'Create New Task'}</h2>

      {error && <p style={{ color: '#dc2626', marginBottom: '0.75rem' }}>{error}</p>}

      <div className="form-group">
        <label htmlFor="title">Title *</label>
        <input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="What needs to be done?"
          maxLength={200}
        />
      </div>

      <div className="form-group">
        <label htmlFor="description">Description</label>
        <textarea
          id="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Add more details (optional)"
          maxLength={2000}
        />
      </div>

      <div className="form-group">
        <label htmlFor="status">Status</label>
        <select
          id="status"
          value={status}
          onChange={(e) => setStatus(e.target.value)}
        >
          <option value="todo">To Do</option>
          <option value="in_progress">In Progress</option>
          <option value="completed">Completed</option>
        </select>
      </div>

      <div className="form-actions">
        <button className="btn btn-primary" type="submit" disabled={submitting}>
          {submitting ? 'Saving…' : editingTask ? 'Update Task' : 'Create Task'}
        </button>
        {editingTask && (
          <button className="btn btn-secondary" type="button" onClick={onCancel}>
            Cancel
          </button>
        )}
      </div>
    </form>
  );
}

