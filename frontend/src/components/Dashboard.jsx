/**
 * Dashboard — summary statistics at the top of the page.
 */
export default function Dashboard({ tasks }) {
  const total = tasks.length;
  const todo = tasks.filter((t) => t.status === 'todo').length;
  const inProgress = tasks.filter((t) => t.status === 'in_progress').length;
  const completed = tasks.filter((t) => t.status === 'completed').length;

  return (
    <div className="dashboard-stats">
      <div className="stat-card">
        <div className="stat-value">{total}</div>
        <div className="stat-label">Total Tasks</div>
      </div>
      <div className="stat-card">
        <div className="stat-value" style={{ color: '#1e40af' }}>{todo}</div>
        <div className="stat-label">To Do</div>
      </div>
      <div className="stat-card">
        <div className="stat-value" style={{ color: '#92400e' }}>{inProgress}</div>
        <div className="stat-label">In Progress</div>
      </div>
      <div className="stat-card">
        <div className="stat-value" style={{ color: '#065f46' }}>{completed}</div>
        <div className="stat-label">Completed</div>
      </div>
    </div>
  );
}

