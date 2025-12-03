import { useState, useEffect } from 'react';
import { upsApi, UPSConfig, UPSStatus, UPSEvent } from '@/api/ups';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';
import { Battery, Zap, Power, AlertTriangle, CheckCircle, Clock } from 'lucide-react';

export function UPSSection(_props: { user: any; systemInfo: any }) {
  const [config, setConfig] = useState<UPSConfig | null>(null);
  const [status, setStatus] = useState<UPSStatus | null>(null);
  const [events, setEvents] = useState<UPSEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadConfig();
    loadStatus();
    loadEvents();
  }, []);

  const loadConfig = async () => {
    try {
      const response = await upsApi.getConfig();
      if (response.success && response.data) {
        setConfig(response.data);
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const loadStatus = async () => {
    try {
      const response = await upsApi.getStatus();
      if (response.success && response.data) {
        setStatus(response.data);
      }
    } catch (err) {
      // Silently fail if UPS is not enabled
    }
  };

  const loadEvents = async () => {
    try {
      const response = await upsApi.getEvents(10);
      if (response.success && response.data) {
        setEvents(response.data);
      }
    } catch (err) {
      // Silently fail
    }
  };

  const handleSave = async () => {
    if (!config) return;

    setSaving(true);
    try {
      const response = await upsApi.updateConfig(config);
      if (response.success) {
        alert('UPS configuration saved successfully');
        loadStatus();
        loadEvents();
      } else {
        alert(response.error?.message || 'Failed to save configuration');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setSaving(false);
    }
  };

  const handleTest = async () => {
    if (!config) return;

    setTesting(true);
    try {
      const response = await upsApi.testConnection(config);
      if (response.success && response.data) {
        alert(`Connection successful!\n\nUPS Model: ${response.data.model}\nManufacturer: ${response.data.manufacturer}\nBattery: ${response.data.battery_charge}%`);
      } else {
        alert(response.error?.message || 'Connection test failed');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setTesting(false);
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'text-red-600 dark:text-red-400';
      case 'warning':
        return 'text-yellow-600 dark:text-yellow-400';
      default:
        return 'text-blue-600 dark:text-blue-400';
    }
  };

  const getEventIcon = (eventType: string) => {
    switch (eventType) {
      case 'POWER_LOSS':
        return <AlertTriangle className="w-5 h-5 text-yellow-500" />;
      case 'BATTERY_LOW':
        return <Battery className="w-5 h-5 text-red-500" />;
      case 'POWER_RESTORED':
        return <CheckCircle className="w-5 h-5 text-green-500" />;
      case 'SHUTDOWN_INITIATED':
        return <Power className="w-5 h-5 text-red-500" />;
      default:
        return <Zap className="w-5 h-5 text-gray-500" />;
    }
  };

  if (loading) {
    return <div className="text-center py-8">Loading UPS configuration...</div>;
  }

  if (!config) {
    return <div className="text-center py-8">Failed to load UPS configuration</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">UPS Management</h2>
        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
          Configure UPS monitoring and automatic shutdown on power loss
        </p>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Current Status */}
      {config.enabled && status && (
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center gap-2">
              <Battery className="w-5 h-5" />
              Current UPS Status
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <div>
                <div className="text-sm text-gray-600 dark:text-gray-400">Status</div>
                <div className="flex items-center gap-2 mt-1">
                  {status.online ? (
                    <CheckCircle className="w-4 h-4 text-green-500" />
                  ) : (
                    <AlertTriangle className="w-4 h-4 text-yellow-500" />
                  )}
                  <span className="font-semibold">{status.online ? 'Online' : 'On Battery'}</span>
                </div>
              </div>
              <div>
                <div className="text-sm text-gray-600 dark:text-gray-400">Battery Charge</div>
                <div className="text-xl font-semibold text-gray-900 dark:text-gray-100 mt-1">
                  {status.battery_charge}%
                </div>
              </div>
              <div>
                <div className="text-sm text-gray-600 dark:text-gray-400">Runtime</div>
                <div className="text-xl font-semibold text-gray-900 dark:text-gray-100 mt-1">
                  {Math.floor(status.runtime / 60)} min
                </div>
              </div>
              <div>
                <div className="text-sm text-gray-600 dark:text-gray-400">Load</div>
                <div className="text-xl font-semibold text-gray-900 dark:text-gray-100 mt-1">
                  {status.load_percent}%
                </div>
              </div>
            </div>
            <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
              <div className="text-xs text-gray-500 dark:text-gray-400">
                {status.manufacturer} {status.model} • Last Update: {new Date(status.last_update).toLocaleString()}
              </div>
            </div>
          </div>
        </Card>
      )}

      {/* Connection Settings */}
      <Card>
        <div className="p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Connection Settings
          </h3>
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="enabled"
                checked={config.enabled}
                onChange={(e) => setConfig({ ...config, enabled: e.target.checked })}
                className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
              />
              <label htmlFor="enabled" className="text-sm font-medium text-gray-700 dark:text-gray-300">
                Enable UPS Monitoring
              </label>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="UPS Name"
                value={config.ups_name}
                onChange={(e) => setConfig({ ...config, ups_name: e.target.value })}
                placeholder="ups"
                disabled={!config.enabled}
              />
              <Input
                label="Host"
                value={config.ups_host}
                onChange={(e) => setConfig({ ...config, ups_host: e.target.value })}
                placeholder="localhost"
                disabled={!config.enabled}
              />
              <Input
                label="Port"
                type="number"
                value={config.ups_port.toString()}
                onChange={(e) => setConfig({ ...config, ups_port: parseInt(e.target.value) || 3493 })}
                disabled={!config.enabled}
              />
              <Input
                label="Poll Interval (seconds)"
                type="number"
                value={config.poll_interval.toString()}
                onChange={(e) => setConfig({ ...config, poll_interval: parseInt(e.target.value) || 30 })}
                disabled={!config.enabled}
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Username (optional)"
                value={config.ups_username || ''}
                onChange={(e) => setConfig({ ...config, ups_username: e.target.value })}
                disabled={!config.enabled}
              />
              <Input
                label="Password (optional)"
                type="password"
                value={config.ups_password || ''}
                onChange={(e) => setConfig({ ...config, ups_password: e.target.value })}
                disabled={!config.enabled}
              />
            </div>
          </div>
        </div>
      </Card>

      {/* Shutdown Settings */}
      <Card>
        <div className="p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Automatic Shutdown
          </h3>
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="low_battery_shutdown"
                checked={config.low_battery_shutdown}
                onChange={(e) => setConfig({ ...config, low_battery_shutdown: e.target.checked })}
                disabled={!config.enabled}
                className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
              />
              <label htmlFor="low_battery_shutdown" className="text-sm font-medium text-gray-700 dark:text-gray-300">
                Enable automatic shutdown on low battery
              </label>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Low Battery Threshold (%)"
                type="number"
                value={config.low_battery_threshold.toString()}
                onChange={(e) => setConfig({ ...config, low_battery_threshold: parseInt(e.target.value) || 20 })}
                disabled={!config.enabled || !config.low_battery_shutdown}
                min="5"
                max="50"
              />
              <Input
                label="Shutdown Delay (seconds)"
                type="number"
                value={config.shutdown_delay.toString()}
                onChange={(e) => setConfig({ ...config, shutdown_delay: parseInt(e.target.value) || 120 })}
                disabled={!config.enabled || !config.low_battery_shutdown}
                min="0"
                max="600"
              />
            </div>
          </div>
        </div>
      </Card>

      {/* Notifications */}
      <Card>
        <div className="p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Notifications
          </h3>
          <div className="space-y-3">
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="notify_on_power_loss"
                checked={config.notify_on_power_loss}
                onChange={(e) => setConfig({ ...config, notify_on_power_loss: e.target.checked })}
                disabled={!config.enabled}
                className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
              />
              <label htmlFor="notify_on_power_loss" className="text-sm text-gray-700 dark:text-gray-300">
                Notify on power loss
              </label>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="notify_on_battery_low"
                checked={config.notify_on_battery_low}
                onChange={(e) => setConfig({ ...config, notify_on_battery_low: e.target.checked })}
                disabled={!config.enabled}
                className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
              />
              <label htmlFor="notify_on_battery_low" className="text-sm text-gray-700 dark:text-gray-300">
                Notify on low battery
              </label>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="notify_on_power_restored"
                checked={config.notify_on_power_restored}
                onChange={(e) => setConfig({ ...config, notify_on_power_restored: e.target.checked })}
                disabled={!config.enabled}
                className="rounded border-gray-300 text-macos-blue focus:ring-macos-blue"
              />
              <label htmlFor="notify_on_power_restored" className="text-sm text-gray-700 dark:text-gray-300">
                Notify on power restored
              </label>
            </div>
          </div>
        </div>
      </Card>

      {/* Recent Events */}
      {config.enabled && events.length > 0 && (
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center gap-2">
              <Clock className="w-5 h-5" />
              Recent Events
            </h3>
            <div className="space-y-2">
              {events.map((event) => (
                <div
                  key={event.id}
                  className="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
                >
                  {getEventIcon(event.event_type)}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-sm text-gray-900 dark:text-gray-100">
                        {event.event_type.replace(/_/g, ' ')}
                      </span>
                      <span className={`text-xs font-medium ${getSeverityColor(event.severity)}`}>
                        {event.severity}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      {event.description}
                    </p>
                    {event.battery_level && (
                      <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                        Battery: {event.battery_level}%
                      </p>
                    )}
                    <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                      {new Date(event.created_at).toLocaleString()}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </Card>
      )}

      {/* Actions */}
      <div className="flex gap-3">
        <Button onClick={handleTest} variant="secondary" disabled={!config.enabled || testing}>
          {testing ? 'Testing...' : 'Test Connection'}
        </Button>
        <Button onClick={handleSave} disabled={saving}>
          {saving ? 'Saving...' : 'Save Configuration'}
        </Button>
      </div>
    </div>
  );
}
