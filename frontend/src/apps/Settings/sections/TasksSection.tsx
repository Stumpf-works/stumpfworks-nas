// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState, useEffect } from 'react';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { tasksApi } from '@/api/tasks';
import { getErrorMessage } from '@/api/client';

export function TasksSection() {
  const [tasks, setTasks] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadTasks();
  }, []);

  const loadTasks = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await tasksApi.listTasks();
      setTasks(response.tasks || []);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Scheduled Tasks</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Manage cron jobs and automated tasks
        </p>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Scheduled Jobs
            </h2>
            <Button variant="primary" onClick={loadTasks}>
              Refresh
            </Button>
          </div>

          {loading ? (
            <p className="text-gray-600 dark:text-gray-400">Loading tasks...</p>
          ) : tasks.length === 0 ? (
            <p className="text-gray-600 dark:text-gray-400">No scheduled tasks configured</p>
          ) : (
            <div className="space-y-3">
              {tasks.map((task: any) => (
                <div
                  key={task.id}
                  className="p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg"
                >
                  <div className="flex items-start justify-between">
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100">{task.name}</h3>
                      <p className="text-sm text-gray-600 dark:text-gray-400">{task.schedule}</p>
                      <p className="text-sm text-gray-500 dark:text-gray-500 font-mono mt-1">
                        {task.command}
                      </p>
                    </div>
                    <span
                      className={`px-2 py-1 text-xs rounded ${
                        task.enabled
                          ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-200'
                          : 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
                      }`}
                    >
                      {task.enabled ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </Card>

      {/* Note */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          For creating and editing scheduled tasks, use the dedicated Tasks app.
        </p>
      </div>
    </div>
  );
}
