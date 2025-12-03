import { useState, useEffect } from 'react';
import { CheckCircle, AlertTriangle, Clock, User } from 'lucide-react';
import { alertRulesApi, type AlertRuleExecution } from '@/api/alertrules';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

export function AlertExecutions() {
  const [executions, setExecutions] = useState<AlertRuleExecution[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [ackModal, setAckModal] = useState<{ execution: AlertRuleExecution; note: string } | null>(
    null
  );

  useEffect(() => {
    loadExecutions();
  }, []);

  const loadExecutions = async () => {
    try {
      const response = await alertRulesApi.getRecentExecutions(100);
      if (response.success && response.data) {
        setExecutions(response.data);
      }
      setError('');
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleAcknowledge = (execution: AlertRuleExecution) => {
    setAckModal({ execution, note: '' });
  };

  const submitAcknowledgment = async () => {
    if (!ackModal) return;

    try {
      // Get current user from session/context - for now use a placeholder
      const username = 'Admin'; // TODO: Get from auth context

      await alertRulesApi.acknowledgeExecution(
        ackModal.execution.id,
        username,
        ackModal.note || undefined
      );

      setAckModal(null);
      loadExecutions();
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString();
  };

  const getSeverityColor = (execution: AlertRuleExecution) => {
    if (!execution.rule) return 'gray';
    switch (execution.rule.severity) {
      case 'critical':
        return 'red';
      case 'warning':
        return 'yellow';
      default:
        return 'blue';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-gray-500">Loading alert history...</div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Alert History</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            View and acknowledge triggered alerts
          </p>
        </div>
        <Button onClick={loadExecutions} variant="secondary">
          Refresh
        </Button>
      </div>

      {/* Error Banner */}
      {error && (
        <Card className="p-4 bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800">
          <div className="flex items-center gap-2 text-red-600 dark:text-red-400">
            <AlertTriangle className="w-5 h-5" />
            <span>{error}</span>
          </div>
        </Card>
      )}

      {/* Executions List */}
      {executions.length === 0 ? (
        <Card className="p-12 text-center">
          <Clock className="w-12 h-12 mx-auto text-gray-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            No Alert History
          </h3>
          <p className="text-gray-500 dark:text-gray-400">
            No alerts have been triggered yet
          </p>
        </Card>
      ) : (
        <div className="space-y-3">
          {executions.map((execution) => {
            const color = getSeverityColor(execution);
            const colorClasses = {
              red: 'border-l-red-500 bg-red-50/50 dark:bg-red-900/10',
              yellow: 'border-l-yellow-500 bg-yellow-50/50 dark:bg-yellow-900/10',
              blue: 'border-l-blue-500 bg-blue-50/50 dark:bg-blue-900/10',
              gray: 'border-l-gray-500 bg-gray-50/50 dark:bg-gray-900/10',
            };

            return (
              <Card
                key={execution.id}
                className={`p-4 border-l-4 ${colorClasses[color as keyof typeof colorClasses]}`}
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-2">
                      <h4 className="font-semibold text-gray-900 dark:text-gray-100">
                        {execution.rule?.name || 'Unknown Rule'}
                      </h4>
                      {execution.rule && (
                        <span
                          className={`px-2 py-0.5 rounded text-xs font-medium ${
                            color === 'red'
                              ? 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400'
                              : color === 'yellow'
                              ? 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600 dark:text-yellow-400'
                              : 'bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400'
                          }`}
                        >
                          {execution.rule.severity.toUpperCase()}
                        </span>
                      )}
                      {execution.acknowledged && (
                        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400">
                          <CheckCircle className="w-3 h-3" />
                          ACKNOWLEDGED
                        </span>
                      )}
                    </div>

                    <p className="text-sm text-gray-700 dark:text-gray-300 mb-2">
                      {execution.message}
                    </p>

                    <div className="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
                      <div className="flex items-center gap-1">
                        <Clock className="w-3 h-3" />
                        {formatDate(execution.createdAt)}
                      </div>
                      <div>
                        Value: <span className="font-mono">{execution.metricValue.toFixed(2)}</span>{' '}
                        / Threshold:{' '}
                        <span className="font-mono">{execution.threshold.toFixed(2)}</span>
                      </div>
                      {execution.acknowledgedBy && (
                        <div className="flex items-center gap-1">
                          <User className="w-3 h-3" />
                          {execution.acknowledgedBy}
                        </div>
                      )}
                    </div>

                    {execution.acknowledgeNote && (
                      <div className="mt-2 p-2 bg-white dark:bg-gray-800 rounded text-sm">
                        <span className="font-medium">Note:</span> {execution.acknowledgeNote}
                      </div>
                    )}
                  </div>

                  {!execution.acknowledged && (
                    <Button
                      onClick={() => handleAcknowledge(execution)}
                      variant="secondary"
                      className="shrink-0"
                    >
                      Acknowledge
                    </Button>
                  )}
                </div>
              </Card>
            );
          })}
        </div>
      )}

      {/* Acknowledgment Modal */}
      {ackModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50">
          <Card className="w-full max-w-md p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Acknowledge Alert
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
              {ackModal.execution.message}
            </p>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Note (optional)
              </label>
              <textarea
                value={ackModal.note}
                onChange={(e) => setAckModal({ ...ackModal, note: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                rows={3}
                placeholder="Add a note about this acknowledgment..."
              />
            </div>
            <div className="flex gap-3">
              <Button onClick={submitAcknowledgment}>Acknowledge</Button>
              <Button variant="secondary" onClick={() => setAckModal(null)}>
                Cancel
              </Button>
            </div>
          </Card>
        </div>
      )}
    </div>
  );
}
