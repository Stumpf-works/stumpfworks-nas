import { useState, useEffect } from 'react';
import { X } from 'lucide-react';
import { alertRulesApi, type AlertRule } from '@/api/alertrules';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';

interface RuleModalProps {
  rule: AlertRule | null;
  onClose: () => void;
  onSaved: () => void;
}

export function RuleModal({ rule, onClose, onSaved }: RuleModalProps) {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    metricType: 'cpu' as const,
    condition: 'gt' as const,
    threshold: 80,
    duration: 0,
    cooldownMins: 5,
    severity: 'warning' as const,
    enabled: true,
    notifyEmail: true,
    notifyWebhook: false,
    notifyChannels: '',
    triggerCount: 0,
  });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (rule) {
      setFormData({
        name: rule.name,
        description: rule.description || '',
        metricType: rule.metricType as any,
        condition: rule.condition as any,
        threshold: rule.threshold,
        duration: rule.duration,
        cooldownMins: rule.cooldownMins,
        severity: rule.severity as any,
        enabled: rule.enabled,
        notifyEmail: rule.notifyEmail,
        notifyWebhook: rule.notifyWebhook,
        notifyChannels: rule.notifyChannels || '',
        triggerCount: rule.triggerCount || 0,
      });
    }
  }, [rule]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSaving(true);

    try {
      if (rule) {
        await alertRulesApi.updateRule(rule.id, formData);
      } else {
        await alertRulesApi.createRule(formData);
      }
      onSaved();
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setSaving(false);
    }
  };

  const metricTypes = [
    { value: 'cpu', label: 'CPU Usage (%)' },
    { value: 'memory', label: 'Memory Usage (%)' },
    { value: 'disk', label: 'Disk Usage (%)' },
    { value: 'temperature', label: 'CPU Temperature (�C)' },
    { value: 'health', label: 'Health Score' },
    { value: 'iops', label: 'Disk IOPS' },
    { value: 'network', label: 'Network Traffic (bytes/s)' },
  ];

  const conditions = [
    { value: 'gt', label: 'Greater than (>)' },
    { value: 'gte', label: 'Greater than or equal (e)' },
    { value: 'lt', label: 'Less than (<)' },
    { value: 'lte', label: 'Less than or equal (d)' },
    { value: 'eq', label: 'Equal to (=)' },
  ];

  const severities = [
    { value: 'info', label: 'Info', color: 'text-blue-600' },
    { value: 'warning', label: 'Warning', color: 'text-yellow-600' },
    { value: 'critical', label: 'Critical', color: 'text-red-600' },
  ];

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50">
      <Card className="w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 p-6 flex items-center justify-between">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
            {rule ? 'Edit Alert Rule' : 'Create Alert Rule'}
          </h3>
          <button
            onClick={onClose}
            className="p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {/* Error Banner */}
          {error && (
            <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
              <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
            </div>
          )}

          {/* Basic Information */}
          <div className="space-y-4">
            <h4 className="text-sm font-semibold text-gray-900 dark:text-gray-100 uppercase tracking-wide">
              Basic Information
            </h4>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Rule Name <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                placeholder="e.g., High CPU Usage"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Description
              </label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                rows={2}
                placeholder="Describe what this alert monitors"
              />
            </div>
          </div>

          {/* Alert Condition */}
          <div className="space-y-4">
            <h4 className="text-sm font-semibold text-gray-900 dark:text-gray-100 uppercase tracking-wide">
              Alert Condition
            </h4>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Metric Type <span className="text-red-500">*</span>
                </label>
                <select
                  value={formData.metricType}
                  onChange={(e) => setFormData({ ...formData, metricType: e.target.value as any })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                  required
                >
                  {metricTypes.map((type) => (
                    <option key={type.value} value={type.value}>
                      {type.label}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Condition <span className="text-red-500">*</span>
                </label>
                <select
                  value={formData.condition}
                  onChange={(e) => setFormData({ ...formData, condition: e.target.value as any })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                  required
                >
                  {conditions.map((cond) => (
                    <option key={cond.value} value={cond.value}>
                      {cond.label}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Threshold <span className="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  step="0.01"
                  value={formData.threshold}
                  onChange={(e) => setFormData({ ...formData, threshold: parseFloat(e.target.value) })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                  required
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Duration (seconds)
                </label>
                <input
                  type="number"
                  min="0"
                  value={formData.duration}
                  onChange={(e) => setFormData({ ...formData, duration: parseInt(e.target.value) || 0 })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                  placeholder="0 for instant alert"
                />
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  Condition must be met for this duration before triggering (0 = instant)
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Cooldown (minutes) <span className="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  min="1"
                  value={formData.cooldownMins}
                  onChange={(e) => setFormData({ ...formData, cooldownMins: parseInt(e.target.value) || 5 })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                  required
                />
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  Minimum time between repeat alerts
                </p>
              </div>
            </div>
          </div>

          {/* Severity & Status */}
          <div className="space-y-4">
            <h4 className="text-sm font-semibold text-gray-900 dark:text-gray-100 uppercase tracking-wide">
              Severity & Status
            </h4>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Severity Level <span className="text-red-500">*</span>
              </label>
              <div className="grid grid-cols-3 gap-3">
                {severities.map((severity) => (
                  <button
                    key={severity.value}
                    type="button"
                    onClick={() => setFormData({ ...formData, severity: severity.value as any })}
                    className={`
                      px-4 py-3 rounded-lg border-2 transition-all font-medium text-sm
                      ${
                        formData.severity === severity.value
                          ? 'border-macos-blue bg-blue-50 dark:bg-blue-900/20 text-macos-blue'
                          : 'border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:border-gray-400'
                      }
                    `}
                  >
                    <span className={severity.color}>{severity.label}</span>
                  </button>
                ))}
              </div>
            </div>

            <div className="flex items-center">
              <input
                type="checkbox"
                id="enabled"
                checked={formData.enabled}
                onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
                className="w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
              />
              <label htmlFor="enabled" className="ml-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                Enable this rule
              </label>
            </div>
          </div>

          {/* Notifications */}
          <div className="space-y-4">
            <h4 className="text-sm font-semibold text-gray-900 dark:text-gray-100 uppercase tracking-wide">
              Notifications
            </h4>

            <div className="space-y-3">
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="notifyEmail"
                  checked={formData.notifyEmail}
                  onChange={(e) => setFormData({ ...formData, notifyEmail: e.target.checked })}
                  className="w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
                />
                <label htmlFor="notifyEmail" className="ml-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Send email notifications
                </label>
              </div>

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="notifyWebhook"
                  checked={formData.notifyWebhook}
                  onChange={(e) => setFormData({ ...formData, notifyWebhook: e.target.checked })}
                  className="w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
                />
                <label htmlFor="notifyWebhook" className="ml-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Send webhook notifications
                </label>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving...' : rule ? 'Update Rule' : 'Create Rule'}
            </Button>
            <Button type="button" variant="secondary" onClick={onClose} disabled={saving}>
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
