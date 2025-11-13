import { useState, useEffect } from 'react';
import { tasksApi, ScheduledTask, TaskExecution } from '@/api/tasks';
import './Tasks.css';

export function Tasks() {
  const [tasks, setTasks] = useState<ScheduledTask[]>([]);
  const [selectedTask, setSelectedTask] = useState<ScheduledTask | null>(null);
  const [executions, setExecutions] = useState<TaskExecution[]>([]);
  const [loading, setLoading] = useState(false);
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showEditDialog, setShowEditDialog] = useState(false);
  const [showExecutionsDialog, setShowExecutionsDialog] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [cronValidation, setCronValidation] = useState<{
    valid: boolean;
    error?: string;
    nextRuns?: string[];
  } | null>(null);

  const [formData, setFormData] = useState<Omit<ScheduledTask, 'id'>>({
    name: '',
    description: '',
    taskType: 'cleanup',
    cronExpression: '0 2 * * *',
    enabled: true,
    config: '',
    timeoutSeconds: 300,
    retryOnFailure: false,
  });

  useEffect(() => {
    loadTasks();
  }, []);

  const loadTasks = async () => {
    try {
      setLoading(true);
      const response = await tasksApi.listTasks(0, 100);
      setTasks(response.tasks || []);
    } catch (error: any) {
      showMessage('error', `Failed to load tasks: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const loadExecutions = async (taskId: number) => {
    try {
      const response = await tasksApi.getTaskExecutions(taskId, 0, 50);
      setExecutions(response.executions || []);
      setShowExecutionsDialog(true);
    } catch (error: any) {
      showMessage('error', `Failed to load executions: ${error.message}`);
    }
  };

  const validateCronExpression = async (expression: string) => {
    if (!expression) {
      setCronValidation(null);
      return;
    }

    try {
      const result = await tasksApi.validateCron(expression);
      setCronValidation(result);
    } catch (error: any) {
      setCronValidation({ valid: false, error: error.message });
    }
  };

  const handleCreateTask = async () => {
    try {
      setLoading(true);
      await tasksApi.createTask(formData);
      showMessage('success', 'Task created successfully');
      setShowCreateDialog(false);
      resetForm();
      loadTasks();
    } catch (error: any) {
      showMessage('error', `Failed to create task: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateTask = async () => {
    if (!selectedTask?.id) return;

    try {
      setLoading(true);
      await tasksApi.updateTask(selectedTask.id, { ...formData, id: selectedTask.id });
      showMessage('success', 'Task updated successfully');
      setShowEditDialog(false);
      setSelectedTask(null);
      resetForm();
      loadTasks();
    } catch (error: any) {
      showMessage('error', `Failed to update task: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteTask = async (id: number) => {
    if (!confirm('Are you sure you want to delete this task?')) return;

    try {
      setLoading(true);
      await tasksApi.deleteTask(id);
      showMessage('success', 'Task deleted successfully');
      loadTasks();
    } catch (error: any) {
      showMessage('error', `Failed to delete task: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleRunNow = async (id: number, name: string) => {
    if (!confirm(`Run task "${name}" now?`)) return;

    try {
      await tasksApi.runTaskNow(id);
      showMessage('success', 'Task execution started');
      setTimeout(loadTasks, 2000); // Reload after 2 seconds
    } catch (error: any) {
      showMessage('error', `Failed to run task: ${error.message}`);
    }
  };

  const handleToggleEnabled = async (task: ScheduledTask) => {
    if (!task.id) return;

    try {
      await tasksApi.updateTask(task.id, { ...task, enabled: !task.enabled });
      showMessage('success', `Task ${!task.enabled ? 'enabled' : 'disabled'}`);
      loadTasks();
    } catch (error: any) {
      showMessage('error', `Failed to toggle task: ${error.message}`);
    }
  };

  const openEditDialog = (task: ScheduledTask) => {
    setSelectedTask(task);
    setFormData({
      name: task.name,
      description: task.description || '',
      taskType: task.taskType,
      cronExpression: task.cronExpression,
      enabled: task.enabled,
      config: task.config || '',
      timeoutSeconds: task.timeoutSeconds,
      retryOnFailure: task.retryOnFailure,
    });
    validateCronExpression(task.cronExpression);
    setShowEditDialog(true);
  };

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
      taskType: 'cleanup',
      cronExpression: '0 2 * * *',
      enabled: true,
      config: '',
      timeoutSeconds: 300,
      retryOnFailure: false,
    });
    setCronValidation(null);
  };

  const showMessage = (type: 'success' | 'error', text: string) => {
    setMessage({ type, text });
    setTimeout(() => setMessage(null), 5000);
  };

  const getStatusBadgeClass = (status?: string) => {
    switch (status) {
      case 'success':
        return 'status-badge status-success';
      case 'failed':
        return 'status-badge status-failed';
      case 'running':
        return 'status-badge status-running';
      default:
        return 'status-badge';
    }
  };

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return 'Never';
    return new Date(dateStr).toLocaleString();
  };

  const formatDuration = (ms: number) => {
    if (ms < 1000) return `${ms}ms`;
    const seconds = Math.floor(ms / 1000);
    if (seconds < 60) return `${seconds}s`;
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}m ${remainingSeconds}s`;
  };

  return (
    <div className="tasks-container">
      <div className="tasks-header">
        <h2>📅 Scheduled Tasks</h2>
        <button
          className="btn btn-primary"
          onClick={() => {
            resetForm();
            setShowCreateDialog(true);
          }}
        >
          + Create Task
        </button>
      </div>

      {message && (
        <div className={`message message-${message.type}`}>
          {message.text}
        </div>
      )}

      {loading && <div className="loading">Loading...</div>}

      <div className="tasks-list">
        {tasks.length === 0 ? (
          <div className="empty-state">
            <p>No scheduled tasks configured</p>
            <button className="btn btn-primary" onClick={() => setShowCreateDialog(true)}>
              Create your first task
            </button>
          </div>
        ) : (
          <table className="tasks-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Schedule</th>
                <th>Last Run</th>
                <th>Next Run</th>
                <th>Status</th>
                <th>Runs</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {tasks.map((task) => (
                <tr key={task.id} className={!task.enabled ? 'task-disabled' : ''}>
                  <td>
                    <div className="task-name">{task.name}</div>
                    {task.description && <div className="task-description">{task.description}</div>}
                  </td>
                  <td>
                    <span className="task-type">{task.taskType}</span>
                  </td>
                  <td>
                    <code className="cron-expression">{task.cronExpression}</code>
                  </td>
                  <td>{formatDate(task.lastRun)}</td>
                  <td>{formatDate(task.nextRun)}</td>
                  <td>
                    <span className={getStatusBadgeClass(task.lastStatus)}>
                      {task.lastStatus || 'pending'}
                    </span>
                  </td>
                  <td>{task.runCount || 0}</td>
                  <td className="task-actions">
                    <button
                      className="btn btn-sm btn-icon"
                      onClick={() => handleToggleEnabled(task)}
                      title={task.enabled ? 'Disable' : 'Enable'}
                    >
                      {task.enabled ? '⏸️' : '▶️'}
                    </button>
                    <button
                      className="btn btn-sm btn-icon"
                      onClick={() => task.id && handleRunNow(task.id, task.name)}
                      title="Run Now"
                    >
                      ▶️
                    </button>
                    <button
                      className="btn btn-sm btn-icon"
                      onClick={() => task.id && loadExecutions(task.id)}
                      title="View History"
                    >
                      📊
                    </button>
                    <button
                      className="btn btn-sm btn-icon"
                      onClick={() => openEditDialog(task)}
                      title="Edit"
                    >
                      ✏️
                    </button>
                    <button
                      className="btn btn-sm btn-icon btn-danger"
                      onClick={() => task.id && handleDeleteTask(task.id)}
                      title="Delete"
                    >
                      🗑️
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Create Task Dialog */}
      {showCreateDialog && (
        <div className="dialog-overlay" onClick={() => setShowCreateDialog(false)}>
          <div className="dialog" onClick={(e) => e.stopPropagation()}>
            <div className="dialog-header">
              <h3>Create Scheduled Task</h3>
              <button className="dialog-close" onClick={() => setShowCreateDialog(false)}>
                ×
              </button>
            </div>
            <div className="dialog-body">
              <TaskForm
                formData={formData}
                setFormData={setFormData}
                cronValidation={cronValidation}
                onCronChange={validateCronExpression}
              />
            </div>
            <div className="dialog-footer">
              <button className="btn" onClick={() => setShowCreateDialog(false)}>
                Cancel
              </button>
              <button
                className="btn btn-primary"
                onClick={handleCreateTask}
                disabled={loading || !cronValidation?.valid}
              >
                Create Task
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Edit Task Dialog */}
      {showEditDialog && (
        <div className="dialog-overlay" onClick={() => setShowEditDialog(false)}>
          <div className="dialog" onClick={(e) => e.stopPropagation()}>
            <div className="dialog-header">
              <h3>Edit Task</h3>
              <button className="dialog-close" onClick={() => setShowEditDialog(false)}>
                ×
              </button>
            </div>
            <div className="dialog-body">
              <TaskForm
                formData={formData}
                setFormData={setFormData}
                cronValidation={cronValidation}
                onCronChange={validateCronExpression}
              />
            </div>
            <div className="dialog-footer">
              <button className="btn" onClick={() => setShowEditDialog(false)}>
                Cancel
              </button>
              <button
                className="btn btn-primary"
                onClick={handleUpdateTask}
                disabled={loading || !cronValidation?.valid}
              >
                Update Task
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Executions Dialog */}
      {showExecutionsDialog && (
        <div className="dialog-overlay" onClick={() => setShowExecutionsDialog(false)}>
          <div className="dialog dialog-large" onClick={(e) => e.stopPropagation()}>
            <div className="dialog-header">
              <h3>Execution History</h3>
              <button className="dialog-close" onClick={() => setShowExecutionsDialog(false)}>
                ×
              </button>
            </div>
            <div className="dialog-body">
              {executions.length === 0 ? (
                <p>No execution history available</p>
              ) : (
                <table className="executions-table">
                  <thead>
                    <tr>
                      <th>Started</th>
                      <th>Duration</th>
                      <th>Status</th>
                      <th>Triggered By</th>
                      <th>Output/Error</th>
                    </tr>
                  </thead>
                  <tbody>
                    {executions.map((exec) => (
                      <tr key={exec.id}>
                        <td>{formatDate(exec.startedAt)}</td>
                        <td>{formatDuration(exec.duration)}</td>
                        <td>
                          <span className={getStatusBadgeClass(exec.status)}>
                            {exec.status}
                          </span>
                        </td>
                        <td>{exec.triggeredBy}</td>
                        <td className="execution-output">
                          {exec.error ? (
                            <span className="error-text">{exec.error}</span>
                          ) : (
                            <span className="success-text">{exec.output || 'No output'}</span>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
            <div className="dialog-footer">
              <button className="btn" onClick={() => setShowExecutionsDialog(false)}>
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

interface TaskFormProps {
  formData: Omit<ScheduledTask, 'id'>;
  setFormData: (data: Omit<ScheduledTask, 'id'>) => void;
  cronValidation: { valid: boolean; error?: string; nextRuns?: string[] } | null;
  onCronChange: (expression: string) => void;
}

function TaskForm({ formData, setFormData, cronValidation, onCronChange }: TaskFormProps) {
  const handleChange = (field: keyof typeof formData, value: any) => {
    const updated = { ...formData, [field]: value };
    setFormData(updated);

    if (field === 'cronExpression') {
      onCronChange(value);
    }
  };

  return (
    <div className="task-form">
      <div className="form-group">
        <label>Task Name *</label>
        <input
          type="text"
          value={formData.name}
          onChange={(e) => handleChange('name', e.target.value)}
          placeholder="e.g., Daily Cleanup"
          required
        />
      </div>

      <div className="form-group">
        <label>Description</label>
        <textarea
          value={formData.description}
          onChange={(e) => handleChange('description', e.target.value)}
          placeholder="Optional description"
          rows={2}
        />
      </div>

      <div className="form-group">
        <label>Task Type *</label>
        <select
          value={formData.taskType}
          onChange={(e) => handleChange('taskType', e.target.value)}
        >
          <option value="cleanup">Cleanup</option>
          <option value="maintenance">Database Maintenance</option>
          <option value="log_rotation">Log Rotation</option>
        </select>
      </div>

      <div className="form-group">
        <label>Cron Expression *</label>
        <input
          type="text"
          value={formData.cronExpression}
          onChange={(e) => handleChange('cronExpression', e.target.value)}
          placeholder="0 2 * * *"
          required
        />
        <small className="form-help">
          Format: minute hour day month weekday (e.g., "0 2 * * *" = daily at 2 AM)
        </small>
        {cronValidation && (
          <div className={`cron-validation ${cronValidation.valid ? 'valid' : 'invalid'}`}>
            {cronValidation.valid ? (
              <>
                <div className="validation-status">✓ Valid cron expression</div>
                {cronValidation.nextRuns && (
                  <div className="next-runs">
                    <strong>Next 5 runs:</strong>
                    <ul>
                      {cronValidation.nextRuns.map((run, i) => (
                        <li key={i}>{run}</li>
                      ))}
                    </ul>
                  </div>
                )}
              </>
            ) : (
              <div className="validation-status error">✗ {cronValidation.error}</div>
            )}
          </div>
        )}
      </div>

      <div className="form-group">
        <label>Configuration (JSON)</label>
        <textarea
          value={formData.config}
          onChange={(e) => handleChange('config', e.target.value)}
          placeholder='{"retentionDays": 30}'
          rows={3}
        />
        <small className="form-help">
          Optional JSON configuration for the task (e.g., retentionDays for cleanup)
        </small>
      </div>

      <div className="form-row">
        <div className="form-group">
          <label>Timeout (seconds)</label>
          <input
            type="number"
            value={formData.timeoutSeconds}
            onChange={(e) => handleChange('timeoutSeconds', parseInt(e.target.value) || 300)}
            min="1"
            max="3600"
          />
        </div>

        <div className="form-group">
          <label className="checkbox-label">
            <input
              type="checkbox"
              checked={formData.enabled}
              onChange={(e) => handleChange('enabled', e.target.checked)}
            />
            Enabled
          </label>
        </div>

        <div className="form-group">
          <label className="checkbox-label">
            <input
              type="checkbox"
              checked={formData.retryOnFailure}
              onChange={(e) => handleChange('retryOnFailure', e.target.checked)}
            />
            Retry on Failure
          </label>
        </div>
      </div>
    </div>
  );
}
