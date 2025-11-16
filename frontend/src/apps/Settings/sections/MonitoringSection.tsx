// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState } from 'react';
import Card from '@/components/ui/Card';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';

export function MonitoringSection() {
  const [prometheusEnabled, setPrometheusEnabled] = useState(true);
  const [grafanaUrl, setGrafanaUrl] = useState('http://localhost:3000');
  const [datadogApiKey, setDatadogApiKey] = useState('');

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Monitoring Configuration</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure Prometheus, Grafana, and Datadog integration
        </p>
      </div>

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
              onClick={() => setPrometheusEnabled(!prometheusEnabled)}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                prometheusEnabled ? 'bg-macos-blue' : 'bg-gray-300'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  prometheusEnabled ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
          </div>

          {prometheusEnabled && (
            <div className="space-y-2">
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Metrics endpoint: <span className="font-mono text-gray-900 dark:text-gray-100">/metrics</span>
              </p>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Available metrics: CPU, Memory, Disk, Network, ZFS, SMART, Services
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
                value={grafanaUrl}
                onChange={(e) => setGrafanaUrl(e.target.value)}
                placeholder="http://localhost:3000"
              />
            </div>
            <Button variant="secondary" onClick={() => window.open(grafanaUrl, '_blank')}>
              Open Grafana
            </Button>
          </div>
        </div>
      </Card>

      {/* Datadog */}
      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Datadog Integration
          </h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                API Key
              </label>
              <Input
                type="password"
                value={datadogApiKey}
                onChange={(e) => setDatadogApiKey(e.target.value)}
                placeholder="Enter Datadog API key"
              />
            </div>
            <Button variant="primary">
              Save Configuration
            </Button>
          </div>
        </div>
      </Card>

      {/* Info */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          StumpfWorks NAS exports comprehensive metrics including ZFS pool health, SMART disk status,
          service status, and share connections for monitoring with Prometheus and Grafana.
        </p>
      </div>
    </div>
  );
}
