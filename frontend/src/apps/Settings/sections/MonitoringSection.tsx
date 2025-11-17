// Revision: 2025-11-17 | Author: Claude | Version: 1.2.0
import { useState, useEffect } from 'react';
import Card from '@/components/ui/Card';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';
import { monitoringApi, MonitoringConfig } from '@/api/monitoring';

export function MonitoringSection() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [config, setConfig] = useState<MonitoringConfig>({
    prometheus_enabled: true,
    grafana_url: 'http://localhost:3000',
    datadog_enabled: false,
    datadog_api_key: '',
  });
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');

  // Load configuration on mount
  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const response = await monitoringApi.getConfig();
      if (response.success && response.data) {
        setConfig(response.data);
      }
    } catch (err) {
      console.error('Failed to load monitoring config:', err);
      setError('Failed to load configuration');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    setSuccess('');
    setError('');

    try {
      const response = await monitoringApi.updateConfig(config);
      if (response.success) {
        setSuccess('Configuration saved successfully');
        // Reload to get updated data (e.g., api_key_set status)
        await loadConfig();
      } else {
        setError('Failed to save configuration');
      }
    } catch (err) {
      console.error('Failed to save monitoring config:', err);
      setError('Failed to save configuration');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-12">
        <div className="text-gray-500 dark:text-gray-400">Loading configuration...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Monitoring Configuration</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure Prometheus, Grafana, and Datadog integration
        </p>
      </div>

      {/* Success/Error Messages */}
      {success && (
        <div className="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
          <p className="text-sm text-green-800 dark:text-green-200">{success}</p>
        </div>
      )}
      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {/* Prometheus */}
      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                Prometheus Metrics
              </h2>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Export metrics at /metrics endpoint
              </p>
            </div>
            <button
              onClick={() => setConfig({ ...config, prometheus_enabled: !config.prometheus_enabled })}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                config.prometheus_enabled ? 'bg-macos-blue' : 'bg-gray-300'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  config.prometheus_enabled ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
          </div>

          {config.prometheus_enabled && (
            <div className="space-y-2">
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Metrics endpoint: <span className="font-mono text-gray-900 dark:text-gray-100">/metrics</span>
              </p>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Available metrics: CPU, Memory, Disk, Network, ZFS, SMART, Services, Share Connections
              </p>
            </div>
          )}
        </div>
      </Card>

      {/* Grafana */}
      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Grafana Integration
          </h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Grafana URL
              </label>
              <Input
                type="url"
                value={config.grafana_url}
                onChange={(e) => setConfig({ ...config, grafana_url: e.target.value })}
                placeholder="http://localhost:3000"
              />
            </div>
            <Button variant="secondary" onClick={() => window.open(config.grafana_url, '_blank')}>
              Open Grafana
            </Button>
          </div>
        </div>
      </Card>

      {/* Datadog */}
      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                Datadog Integration
              </h2>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Send metrics to Datadog for monitoring
              </p>
            </div>
            <button
              onClick={() => setConfig({ ...config, datadog_enabled: !config.datadog_enabled })}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                config.datadog_enabled ? 'bg-macos-blue' : 'bg-gray-300'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  config.datadog_enabled ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
          </div>

          {config.datadog_enabled && (
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  API Key
                </label>
                <Input
                  type="password"
                  value={config.datadog_api_key || ''}
                  onChange={(e) => setConfig({ ...config, datadog_api_key: e.target.value })}
                  placeholder={config.datadog_api_key_set ? '••••••••••••••••' : 'Enter Datadog API key'}
                />
                {config.datadog_api_key_set && (
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    API key is already set. Enter a new key to change it.
                  </p>
                )}
              </div>
            </div>
          )}
        </div>
      </Card>

      {/* Save Button */}
      <div className="flex justify-end gap-3">
        <Button variant="secondary" onClick={loadConfig} disabled={saving}>
          Reset
        </Button>
        <Button variant="primary" onClick={handleSave} disabled={saving}>
          {saving ? 'Saving...' : 'Save Configuration'}
        </Button>
      </div>

      {/* Info */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          StumpfWorks NAS exports comprehensive metrics including CPU temperature, ZFS pool health,
          SMART disk status, service status, network health scores, and active Samba/NFS connections
          for monitoring with Prometheus and Grafana.
        </p>
      </div>
    </div>
  );
}
