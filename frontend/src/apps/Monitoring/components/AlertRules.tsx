import { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Bell, BellOff, AlertTriangle } from 'lucide-react';
import { alertRulesApi, type AlertRule } from '@/api/alertrules';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';
import { RuleModal } from './RuleModal';

export function AlertRules() {
  const [rules, setRules] = useState<AlertRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [modalOpen, setModalOpen] = useState(false);
  const [selectedRule, setSelectedRule] = useState<AlertRule | null>(null);

  useEffect(() => {
    loadRules();
  }, []);

  const loadRules = async () => {
    try {
      const response = await alertRulesApi.listRules();
      if (response.success && response.data) {
        setRules(response.data);
      }
      setError('');
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setSelectedRule(null);
    setModalOpen(true);
  };

  const handleEdit = (rule: AlertRule) => {
    setSelectedRule(rule);
    setModalOpen(true);
  };

  const handleDelete = async (rule: AlertRule) => {
    if (!confirm(`Delete alert rule "${rule.name}"?`)) return;

    try {
      await alertRulesApi.deleteRule(rule.id);
      loadRules();
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleToggleEnabled = async (rule: AlertRule) => {
    try {
      await alertRulesApi.updateRule(rule.id, { enabled: !rule.enabled });
      loadRules();
    } catch (err) {
      alert(getErrorMessage(err));
    }
  };

  const handleSaved = () => {
    setModalOpen(false);
    loadRules();
  };

  const getMetricLabel = (type: string) => {
    const labels: Record<string, string> = {
      cpu: 'CPU Usage',
      memory: 'Memory Usage',
      disk: 'Disk Usage',
      network: 'Network Traffic',
      health: 'Health Score',
      temperature: 'Temperature',
      iops: 'Disk IOPS',
    };
    return labels[type] || type;
  };

  const getConditionSymbol = (condition: string) => {
    const symbols: Record<string, string> = {
      gt: '>',
      lt: '<',
      eq: '=',
      gte: '≥',
      lte: '≤',
    };
    return symbols[condition] || condition;
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20';
      case 'warning':
        return 'text-yellow-600 dark:text-yellow-400 bg-yellow-50 dark:bg-yellow-900/20';
      default:
        return 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-gray-500">Loading alert rules...</div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Alert Rules</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Configure custom alert rules based on system metrics
          </p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="w-4 h-4 mr-2" />
          Create Rule
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

      {/* Rules List */}
      {rules.length === 0 ? (
        <Card className="p-12 text-center">
          <Bell className="w-12 h-12 mx-auto text-gray-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            No Alert Rules
          </h3>
          <p className="text-gray-500 dark:text-gray-400 mb-4">
            Get started by creating your first alert rule
          </p>
          <Button onClick={handleCreate}>
            <Plus className="w-4 h-4 mr-2" />
            Create Alert Rule
          </Button>
        </Card>
      ) : (
        <div className="grid gap-4">
          {rules.map((rule) => (
            <Card key={rule.id} className="p-6">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      {rule.name}
                    </h3>
                    <span
                      className={`px-2 py-1 rounded text-xs font-medium ${getSeverityColor(
                        rule.severity
                      )}`}
                    >
                      {rule.severity.toUpperCase()}
                    </span>
                    {rule.isActive && (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400">
                        ACTIVE ALERT
                      </span>
                    )}
                    {!rule.enabled && (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400">
                        DISABLED
                      </span>
                    )}
                  </div>

                  {rule.description && (
                    <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
                      {rule.description}
                    </p>
                  )}

                  <div className="flex items-center gap-4 text-sm">
                    <div className="flex items-center gap-2">
                      <span className="text-gray-500 dark:text-gray-400">Condition:</span>
                      <code className="px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded font-mono">
                        {getMetricLabel(rule.metricType)} {getConditionSymbol(rule.condition)}{' '}
                        {rule.threshold}
                        {rule.metricType === 'cpu' || rule.metricType === 'memory' || rule.metricType === 'disk' ? '%' : ''}
                      </code>
                    </div>
                    {rule.duration > 0 && (
                      <div className="flex items-center gap-2">
                        <span className="text-gray-500 dark:text-gray-400">Duration:</span>
                        <span className="font-medium">{rule.duration}s</span>
                      </div>
                    )}
                    <div className="flex items-center gap-2">
                      <span className="text-gray-500 dark:text-gray-400">Cooldown:</span>
                      <span className="font-medium">{rule.cooldownMins}min</span>
                    </div>
                  </div>

                  {rule.lastTriggered && (
                    <div className="mt-3 text-xs text-gray-500 dark:text-gray-400">
                      Last triggered: {new Date(rule.lastTriggered).toLocaleString()} (
                      {rule.triggerCount} times)
                    </div>
                  )}
                </div>

                <div className="flex items-center gap-2 ml-4">
                  <button
                    onClick={() => handleToggleEnabled(rule)}
                    className="p-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
                    title={rule.enabled ? 'Disable rule' : 'Enable rule'}
                  >
                    {rule.enabled ? (
                      <Bell className="w-5 h-5" />
                    ) : (
                      <BellOff className="w-5 h-5" />
                    )}
                  </button>
                  <button
                    onClick={() => handleEdit(rule)}
                    className="p-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
                    title="Edit rule"
                  >
                    <Edit className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(rule)}
                    className="p-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors"
                    title="Delete rule"
                  >
                    <Trash2 className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Rule Modal */}
      {modalOpen && (
        <RuleModal rule={selectedRule} onClose={() => setModalOpen(false)} onSaved={handleSaved} />
      )}
    </div>
  );
}
