/**
 * StatusFilter — a row of pill buttons that let the user
 * filter the task list by status.
 */

const FILTERS = [
  { value: 'all', label: 'All' },
  { value: 'todo', label: 'To Do' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'completed', label: 'Completed' },
];

export default function StatusFilter({ current, onChange }) {
  return (
    <div className="filter-bar" role="group" aria-label="Filter tasks by status">
      {FILTERS.map((f) => (
        <button
          key={f.value}
          className={`filter-btn ${current === f.value ? 'active' : ''}`}
          onClick={() => onChange(f.value)}
          aria-pressed={current === f.value}
        >
          {f.label}
        </button>
      ))}
    </div>
  );
}
