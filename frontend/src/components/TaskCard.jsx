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
    if (!window.confirm('Delete this task?')) return;
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

  return (
    <div className="task-card">
      <div className="task-card-body">
        <div className="task-card-title">{task.title}</div>
        {task.description && (
          <div className="task-card-desc">{task.description}</div>
        )}
        <div className="task-card-meta">
          <span
            className={`status-badge status-${task.status}`}
            onClick={handleStatusCycle}
            title="Click to change status"
            style={{ cursor: 'pointer' }}
          >
            {STATUS_LABELS[task.status] || task.status}
          </span>
          <span style={{ marginLeft: '0.75rem' }}>Created {createdDate}</span>
        </div>
      </div>
      <div className="task-card-actions">
        <button className="btn btn-secondary btn-sm" onClick={() => onEdit(task)}>
          Edit
        </button>
        <button
          className="btn btn-danger btn-sm"
          onClick={handleDelete}
          disabled={deleting}
        >
          {deleting ? '…' : 'Delete'}
        </button>
      </div>
    </div>
  );
}

