import { useState, useEffect } from 'react';
import { alertsApi, type AlertConfig } from '@/api/alerts';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

export function Alerts() {
  const [config, setConfig] = useState<AlertConfig | null>(null);
  const [loading, setLoading] = useState(false);
  const [testingEmail, setTestingEmail] = useState(false);
  const [testingWebhook, setTestingWebhook] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    try {
      const response = await alertsApi.getConfig();
      if (response.success && response.data) {
        setConfig(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch alert config:', error);
      // Set default config
      setConfig({
        enabled: false,
        smtpHost: '',
        smtpPort: 587,
        smtpUsername: '',
        smtpFromEmail: '',
        smtpFromName: 'Stumpf.Works NAS',
        smtpUseTLS: true,
        alertRecipient: '',
        webhookEnabled: false,
        webhookType: 'discord',
        webhookURL: '',
        webhookUsername: 'Stumpf.Works NAS',
        webhookAvatarURL: '',
        onFailedLogin: true,
        onIPBlock: true,
        onCriticalEvent: true,
        failedLoginThreshold: 3,
        rateLimitMinutes: 15,
      });
    }
  };

  const handleSave = async () => {
    if (!config) return;

    setLoading(true);
    setMessage(null);

    try {
      const response = await alertsApi.updateConfig(config);
      if (response.success) {
        setMessage({ type: 'success', text: 'Alert configuration saved successfully' });
        if (response.data) {
          setConfig(response.data);
        }
      } else {
        setMessage({ type: 'error', text: response.error?.message || 'Failed to save configuration' });
      }
    } catch (error: any) {
      console.error('Failed to save config:', error);
      setMessage({ type: 'error', text: error.message || 'Failed to save configuration' });
    } finally {
      setLoading(false);
    }
  };

  const handleTestEmail = async () => {
    if (!config) return;

    setTestingEmail(true);
    setMessage(null);

    try {
      const response = await alertsApi.testEmail(config);
      if (response.success) {
        setMessage({ type: 'success', text: 'Test email sent successfully! Check your inbox.' });
      } else {
        setMessage({ type: 'error', text: response.error?.message || 'Failed to send test email' });
      }
    } catch (error: any) {
      console.error('Failed to send test email:', error);
      setMessage({ type: 'error', text: error.message || 'Failed to send test email' });
    } finally {
      setTestingEmail(false);
    }
  };

  const handleTestWebhook = async () => {
    if (!config) return;

    setTestingWebhook(true);
    setMessage(null);

    try {
      const response = await alertsApi.testWebhook(config);
      if (response.success) {
        setMessage({ type: 'success', text: 'Test webhook sent successfully! Check your channel.' });
      } else {
        setMessage({ type: 'error', text: response.error?.message || 'Failed to send test webhook' });
      }
    } catch (error: any) {
      console.error('Failed to send test webhook:', error);
      setMessage({ type: 'error', text: error.message || 'Failed to send test webhook' });
    } finally {
      setTestingWebhook(false);
    }
  };

  if (!config) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="p-6 h-full overflow-auto bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Alert Configuration
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure email and webhook notifications for security events
        </p>
      </div>

      <div className="space-y-6 max-w-4xl">
        {/* Message */}
        {message && (
          <div
            className={`p-4 rounded-lg ${
              message.type === 'success'
                ? 'bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 text-green-800 dark:text-green-200'
                : 'bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-800 dark:text-red-200'
            }`}
          >
            {message.text}
          </div>
        )}

        {/* Email Configuration */}
        <Card>
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  Email Notifications
                </h2>
                <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                  Send alerts via email using SMTP
                </p>
              </div>
              <button
                onClick={() => setConfig({ ...config, enabled: !config.enabled })}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  config.enabled ? 'bg-macos-blue' : 'bg-gray-300 dark:bg-gray-600'
                }`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    config.enabled ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {config.enabled && (
              <div className="space-y-4 mt-4">
                <div className="grid grid-cols-2 gap-4">
                  <Input
                    label="SMTP Host"
                    value={config.smtpHost}
                    onChange={(e) => setConfig({ ...config, smtpHost: e.target.value })}
                    placeholder="smtp.gmail.com"
                  />
                  <Input
                    label="SMTP Port"
                    type="number"
                    value={config.smtpPort}
                    onChange={(e) => setConfig({ ...config, smtpPort: parseInt(e.target.value) || 587 })}
                    placeholder="587"
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <Input
                    label="SMTP Username"
                    value={config.smtpUsername}
                    onChange={(e) => setConfig({ ...config, smtpUsername: e.target.value })}
                    placeholder="your-email@example.com"
                  />
                  <Input
                    label="SMTP Password"
                    type="password"
                    value={config.smtpPassword || ''}
                    onChange={(e) => setConfig({ ...config, smtpPassword: e.target.value })}
                    placeholder="••••••••"
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <Input
                    label="From Email"
                    value={config.smtpFromEmail}
                    onChange={(e) => setConfig({ ...config, smtpFromEmail: e.target.value })}
                    placeholder="alerts@stumpfworks.local"
                  />
                  <Input
                    label="From Name"
                    value={config.smtpFromName}
                    onChange={(e) => setConfig({ ...config, smtpFromName: e.target.value })}
                    placeholder="Stumpf.Works NAS"
                  />
                </div>

                <Input
                  label="Alert Recipient"
                  value={config.alertRecipient}
                  onChange={(e) => setConfig({ ...config, alertRecipient: e.target.value })}
                  placeholder="admin@example.com"
                />

                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id="useTLS"
                    checked={config.smtpUseTLS}
                    onChange={(e) => setConfig({ ...config, smtpUseTLS: e.target.checked })}
                    className="rounded"
                  />
                  <label htmlFor="useTLS" className="text-sm text-gray-700 dark:text-gray-300">
                    Use TLS/SSL
                  </label>
                </div>

                <Button
                  variant="secondary"
                  onClick={handleTestEmail}
                  disabled={testingEmail || !config.alertRecipient}
                  className="w-full sm:w-auto"
                >
                  {testingEmail ? 'Sending...' : 'Send Test Email'}
                </Button>
              </div>
            )}
          </div>
        </Card>

        {/* Webhook Configuration */}
        <Card>
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  Webhook Notifications
                </h2>
                <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                  Send alerts to Discord, Slack, or custom webhooks
                </p>
              </div>
              <button
                onClick={() => setConfig({ ...config, webhookEnabled: !config.webhookEnabled })}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  config.webhookEnabled ? 'bg-macos-blue' : 'bg-gray-300 dark:bg-gray-600'
                }`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    config.webhookEnabled ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {config.webhookEnabled && (
              <div className="space-y-4 mt-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Webhook Type
                  </label>
                  <div className="grid grid-cols-3 gap-2">
                    {['discord', 'slack', 'custom'].map((type) => (
                      <button
                        key={type}
                        onClick={() => setConfig({ ...config, webhookType: type as any })}
                        className={`p-3 rounded-lg border-2 transition-colors ${
                          config.webhookType === type
                            ? 'border-macos-blue bg-blue-50 dark:bg-blue-900/20 text-blue-900 dark:text-blue-100'
                            : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                        }`}
                      >
                        <span className="font-medium capitalize">{type}</span>
                      </button>
                    ))}
                  </div>
                </div>

                <Input
                  label="Webhook URL"
                  value={config.webhookURL}
                  onChange={(e) => setConfig({ ...config, webhookURL: e.target.value })}
                  placeholder={
                    config.webhookType === 'discord'
                      ? 'https://discord.com/api/webhooks/...'
                      : config.webhookType === 'slack'
                      ? 'https://hooks.slack.com/services/...'
                      : 'https://your-webhook-endpoint.com'
                  }
                />

                <div className="grid grid-cols-2 gap-4">
                  <Input
                    label="Display Name (Optional)"
                    value={config.webhookUsername}
                    onChange={(e) => setConfig({ ...config, webhookUsername: e.target.value })}
                    placeholder="Stumpf.Works NAS"
                  />
                  <Input
                    label="Avatar URL (Optional)"
                    value={config.webhookAvatarURL}
                    onChange={(e) => setConfig({ ...config, webhookAvatarURL: e.target.value })}
                    placeholder="https://..."
                  />
                </div>

                <Button
                  variant="secondary"
                  onClick={handleTestWebhook}
                  disabled={testingWebhook || !config.webhookURL}
                  className="w-full sm:w-auto"
                >
                  {testingWebhook ? 'Sending...' : 'Send Test Webhook'}
                </Button>
              </div>
            )}
          </div>
        </Card>

        {/* Alert Triggers */}
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Alert Triggers
            </h2>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-900 dark:text-gray-100">
                    Failed Login Attempts
                  </p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Alert when multiple failed login attempts are detected
                  </p>
                </div>
                <button
                  onClick={() => setConfig({ ...config, onFailedLogin: !config.onFailedLogin })}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                    config.onFailedLogin ? 'bg-macos-blue' : 'bg-gray-300 dark:bg-gray-600'
                  }`}
                >
                  <span
                    className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                      config.onFailedLogin ? 'translate-x-6' : 'translate-x-1'
                    }`}
                  />
                </button>
              </div>

              {config.onFailedLogin && (
                <div>
                  <Input
                    label="Failed Login Threshold"
                    type="number"
                    value={config.failedLoginThreshold}
                    onChange={(e) =>
                      setConfig({ ...config, failedLoginThreshold: parseInt(e.target.value) || 3 })
                    }
                    placeholder="3"
                  />
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    Alert after this many failed attempts
                  </p>
                </div>
              )}

              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-900 dark:text-gray-100">IP Blocking</p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Alert when an IP address is automatically blocked
                  </p>
                </div>
                <button
                  onClick={() => setConfig({ ...config, onIPBlock: !config.onIPBlock })}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                    config.onIPBlock ? 'bg-macos-blue' : 'bg-gray-300 dark:bg-gray-600'
                  }`}
                >
                  <span
                    className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                      config.onIPBlock ? 'translate-x-6' : 'translate-x-1'
                    }`}
                  />
                </button>
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-900 dark:text-gray-100">Critical Events</p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Alert on critical security events
                  </p>
                </div>
                <button
                  onClick={() => setConfig({ ...config, onCriticalEvent: !config.onCriticalEvent })}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                    config.onCriticalEvent ? 'bg-macos-blue' : 'bg-gray-300 dark:bg-gray-600'
                  }`}
                >
                  <span
                    className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                      config.onCriticalEvent ? 'translate-x-6' : 'translate-x-1'
                    }`}
                  />
                </button>
              </div>
            </div>
          </div>
        </Card>

        {/* Rate Limiting */}
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Rate Limiting
            </h2>
            <div>
              <Input
                label="Rate Limit (minutes)"
                type="number"
                value={config.rateLimitMinutes}
                onChange={(e) => setConfig({ ...config, rateLimitMinutes: parseInt(e.target.value) || 15 })}
                placeholder="15"
              />
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Minimum time between alerts of the same type
              </p>
            </div>
          </div>
        </Card>

        {/* Save Button */}
        <div className="flex justify-end gap-3">
          <Button variant="secondary" onClick={fetchConfig} disabled={loading}>
            Reset
          </Button>
          <Button variant="primary" onClick={handleSave} disabled={loading}>
            {loading ? 'Saving...' : 'Save Configuration'}
          </Button>
        </div>
      </div>
    </div>
  );
}
