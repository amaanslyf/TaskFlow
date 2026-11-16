import TaskCard from './TaskCard.jsx';

/**
 * TaskList — renders the array of tasks or a helpful message
 * when there are none (or when an error occurred).
 */
export default function TaskList({ tasks, loading, error, onEdit, onDelete, onStatusChange }) {
  if (loading) {
    return (
      <div className="state-message">
        <div className="spinner" />
        <p style={{ marginTop: '0.75rem' }}>Loading tasks…</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="state-message error">
        <p>⚠️ {error}</p>
        <p style={{ fontSize: '0.85rem', marginTop: '0.5rem' }}>
          Make sure the Go backend is running and the database is reachable.
        </p>
      </div>
    );
  }

  if (tasks.length === 0) {
    return (
      <div className="state-message">
        <p>📋 No tasks yet</p>
        <p style={{ fontSize: '0.85rem', marginTop: '0.5rem' }}>
          Create your first task above!
        </p>
      </div>
    );
  }

  return (
    <div className="task-list">
      {tasks.map((task) => (
        <TaskCard
          key={task.id}
          task={task}
          onEdit={onEdit}
          onDelete={onDelete}
          onStatusChange={onStatusChange}
        />
      ))}
    </div>
  );
}

