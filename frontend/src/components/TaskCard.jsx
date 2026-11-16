import { useState } from 'react';

/**
 * StatusBadge — a small coloured pill showing the task status.
 *
 * The CSS class name matches the status value
 * (e.g. "status-in_progress") so each status gets its own
 * colour automatically.
 */
const STATUS_LABELS = {
  todo: 'To Do',
  in_progress: 'In Progress',
  completed: 'Completed',
};

export default function TaskCard({ task, onEdit, onDelete, onStatusChange }) {
  const [deleting, setDeleting] = useState(false);

  const handleDelete = async () => {
    if (!window.confirm(`Delete this task?`)) return;
    setDeleting(true);
    try {
      await onDelete(task.id);
    } finally {
      setDeleting(false);
    }
  };

  // Cycle through statuses on click.
  const nextStatus = {
    todo: 'in_progress',
    in_progress: 'completed',
    completed: 'todo',
  };

  const handleStatusCycle = () => {
    onStatusChange(task.id, nextStatus[task.status]);
  };

  const createdDate = new Date(task.created_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  const statusText = STATUS_LABELS[task.status] || task.status;

  return (
    <div className="task-card" aria-label={`Task: ${task.title}`}>
      <div className="task-card-body">
        <div className="task-card-title">{task.title}</div>
        {task.description && (
          <div className="task-card-desc">{task.description}</div>
        )}
        <div className="task-card-meta">
          <span
            role="button"
            tabIndex={0}
            className={`status-badge status-${task.status}`}
            onClick={handleStatusCycle}
            onKeyDown={(e) => e.key === 'Enter' && handleStatusCycle()}
            title="Click to change status"
            aria-label={`Status: ${statusText}. Click to change.`}
            style={{ cursor: 'pointer' }}
          >
            {statusText}
          </span>
          <span style={{ marginLeft: '0.75rem' }}>Created {createdDate}</span>
        </div>
      </div>
      <div className="task-card-actions">
        <button
          className="btn btn-secondary btn-sm"
          onClick={() => onEdit(task)}
          aria-label={`Edit task ${task.title}`}
        >
          Edit
        </button>
        <button
          className="btn btn-danger btn-sm"
          onClick={handleDelete}
          disabled={deleting}
          aria-label={`Delete task ${task.title}`}
        >
          {deleting ? '…' : 'Delete'}
        </button>
      </div>
    </div>
  );
}
