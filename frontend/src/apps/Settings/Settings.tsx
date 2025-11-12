import { useState, useEffect } from 'react';
import { useAuthStore, useThemeStore } from '@/store';
import { systemApi } from '@/api/system';
import { authApi } from '@/api/auth';
import { adApi, type ADConfig } from '@/api/ad';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

export function Settings() {
  const user = useAuthStore((state) => state.user);
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const isDark = useThemeStore((state) => state.isDark);
  const toggleTheme = useThemeStore((state) => state.toggleTheme);

  const [systemInfo, setSystemInfo] = useState<any>(null);
  const [adConfig, setAdConfig] = useState<ADConfig | null>(null);
  const [adConfigEditing, setAdConfigEditing] = useState(false);
  const [adTestResult, setAdTestResult] = useState<{ success: boolean; message: string } | null>(null);
  const [adLoading, setAdLoading] = useState(false);

  useEffect(() => {
    const fetchSystemInfo = async () => {
      try {
        const response = await systemApi.getInfo();
        if (response.success && response.data) {
          setSystemInfo(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch system info:', error);
      }
    };

    const fetchAdConfig = async () => {
      try {
        const response = await adApi.getConfig();
        if (response.success && response.data) {
          setAdConfig(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch AD config:', error);
      }
    };

    fetchSystemInfo();
    fetchAdConfig();
  }, []);

  const handleLogout = async () => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      clearAuth();
      window.location.reload();
    }
  };

  const handleAdConfigSave = async () => {
    if (!adConfig) return;

    setAdLoading(true);
    try {
      const response = await adApi.updateConfig(adConfig);
      if (response.success && response.data) {
        setAdConfig(response.data);
        setAdConfigEditing(false);
        setAdTestResult({ success: true, message: 'Configuration saved successfully' });
      }
    } catch (error) {
      console.error('Failed to save AD config:', error);
      setAdTestResult({ success: false, message: 'Failed to save configuration' });
    } finally {
      setAdLoading(false);
    }
  };

  const handleAdTestConnection = async () => {
    setAdLoading(true);
    setAdTestResult(null);
    try {
      const response = await adApi.testConnection();
      if (response.success && response.data) {
        setAdTestResult(response.data);
      }
    } catch (error) {
      console.error('Failed to test AD connection:', error);
      setAdTestResult({ success: false, message: 'Connection test failed' });
    } finally {
      setAdLoading(false);
    }
  };

  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${days}d ${hours}h ${minutes}m`;
  };

  return (
    <div className="p-6 h-full overflow-auto bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Settings
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          System configuration and preferences
        </p>
      </div>

      <div className="space-y-6 max-w-4xl">
        {/* User Information */}
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              User Information
            </h2>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">Username:</span>
                <span className="font-medium text-gray-900 dark:text-gray-100">
                  {user?.username}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">Email:</span>
                <span className="font-medium text-gray-900 dark:text-gray-100">
                  {user?.email}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">Role:</span>
                <span className="font-medium text-gray-900 dark:text-gray-100 capitalize">
                  {user?.role}
                </span>
              </div>
            </div>
          </div>
        </Card>

        {/* Appearance */}
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Appearance
            </h2>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-gray-900 dark:text-gray-100">
                  Dark Mode
                </p>
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  Toggle dark/light theme
                </p>
              </div>
              <button
                onClick={toggleTheme}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  isDark ? 'bg-macos-blue' : 'bg-gray-300'
                }`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    isDark ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>
          </div>
        </Card>

        {/* Active Directory Configuration */}
        {adConfig && (
          <Card>
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  Active Directory
                </h2>
                {!adConfigEditing && (
                  <Button
                    variant="secondary"
                    onClick={() => setAdConfigEditing(true)}
                    size="sm"
                  >
                    Edit
                  </Button>
                )}
              </div>

              <div className="space-y-4">
                {/* Enabled Toggle */}
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium text-gray-900 dark:text-gray-100">
                      AD Integration
                    </p>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      Enable Active Directory authentication
                    </p>
                  </div>
                  <button
                    onClick={() =>
                      adConfigEditing &&
                      setAdConfig({ ...adConfig, enabled: !adConfig.enabled })
                    }
                    disabled={!adConfigEditing}
                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                      adConfig.enabled ? 'bg-macos-blue' : 'bg-gray-300'
                    } ${!adConfigEditing ? 'opacity-50 cursor-not-allowed' : ''}`}
                  >
                    <span
                      className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                        adConfig.enabled ? 'translate-x-6' : 'translate-x-1'
                      }`}
                    />
                  </button>
                </div>

                {/* AD Server Configuration */}
                {adConfigEditing ? (
                  <>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Server
                      </label>
                      <Input
                        value={adConfig.server}
                        onChange={(e) =>
                          setAdConfig({ ...adConfig, server: e.target.value })
                        }
                        placeholder="ldap.example.com"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Port
                      </label>
                      <Input
                        type="number"
                        value={adConfig.port}
                        onChange={(e) =>
                          setAdConfig({ ...adConfig, port: parseInt(e.target.value) })
                        }
                        placeholder="389"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Base DN
                      </label>
                      <Input
                        value={adConfig.baseDN}
                        onChange={(e) =>
                          setAdConfig({ ...adConfig, baseDN: e.target.value })
                        }
                        placeholder="dc=example,dc=com"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Bind User
                      </label>
                      <Input
                        value={adConfig.bindUser}
                        onChange={(e) =>
                          setAdConfig({ ...adConfig, bindUser: e.target.value })
                        }
                        placeholder="cn=admin,dc=example,dc=com"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Bind Password
                      </label>
                      <Input
                        type="password"
                        value={adConfig.bindPassword || ''}
                        onChange={(e) =>
                          setAdConfig({ ...adConfig, bindPassword: e.target.value })
                        }
                        placeholder="••••••••"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        User Filter
                      </label>
                      <Input
                        value={adConfig.userFilter}
                        onChange={(e) =>
                          setAdConfig({ ...adConfig, userFilter: e.target.value })
                        }
                        placeholder="(objectClass=user)"
                      />
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-medium text-gray-900 dark:text-gray-100">
                          Use TLS
                        </p>
                        <p className="text-sm text-gray-600 dark:text-gray-400">
                          Enable secure connection
                        </p>
                      </div>
                      <button
                        onClick={() =>
                          setAdConfig({ ...adConfig, useTLS: !adConfig.useTLS })
                        }
                        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                          adConfig.useTLS ? 'bg-macos-blue' : 'bg-gray-300'
                        }`}
                      >
                        <span
                          className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                            adConfig.useTLS ? 'translate-x-6' : 'translate-x-1'
                          }`}
                        />
                      </button>
                    </div>
                  </>
                ) : (
                  <>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <span className="text-sm text-gray-600 dark:text-gray-400">Server:</span>
                        <p className="font-medium text-gray-900 dark:text-gray-100">
                          {adConfig.server}:{adConfig.port}
                        </p>
                      </div>
                      <div>
                        <span className="text-sm text-gray-600 dark:text-gray-400">Base DN:</span>
                        <p className="font-medium text-gray-900 dark:text-gray-100">
                          {adConfig.baseDN}
                        </p>
                      </div>
                      <div>
                        <span className="text-sm text-gray-600 dark:text-gray-400">TLS:</span>
                        <p className="font-medium text-gray-900 dark:text-gray-100">
                          {adConfig.useTLS ? 'Enabled' : 'Disabled'}
                        </p>
                      </div>
                    </div>
                  </>
                )}

                {/* Test Result */}
                {adTestResult && (
                  <div
                    className={`p-3 rounded-lg ${
                      adTestResult.success
                        ? 'bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800'
                        : 'bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800'
                    }`}
                  >
                    <p
                      className={`text-sm ${
                        adTestResult.success
                          ? 'text-green-800 dark:text-green-200'
                          : 'text-red-800 dark:text-red-200'
                      }`}
                    >
                      {adTestResult.message}
                    </p>
                  </div>
                )}

                {/* Action Buttons */}
                <div className="flex gap-2">
                  {adConfigEditing ? (
                    <>
                      <Button
                        variant="primary"
                        onClick={handleAdConfigSave}
                        disabled={adLoading}
                        className="flex-1"
                      >
                        {adLoading ? 'Saving...' : 'Save'}
                      </Button>
                      <Button
                        variant="secondary"
                        onClick={() => setAdConfigEditing(false)}
                        disabled={adLoading}
                        className="flex-1"
                      >
                        Cancel
                      </Button>
                    </>
                  ) : (
                    <Button
                      variant="secondary"
                      onClick={handleAdTestConnection}
                      disabled={adLoading}
                      className="flex-1"
                    >
                      {adLoading ? 'Testing...' : 'Test Connection'}
                    </Button>
                  )}
                </div>
              </div>
            </div>
          </Card>
        )}

        {/* System Information */}
        {systemInfo && (
          <Card>
            <div className="p-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                System Information
              </h2>
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Hostname:</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {systemInfo.hostname}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Platform:</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {systemInfo.platform}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400">OS:</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {systemInfo.os}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Architecture:</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {systemInfo.architecture}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400">CPU Cores:</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {systemInfo.cpuCores}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Uptime:</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {formatUptime(systemInfo.uptime)}
                  </span>
                </div>
              </div>
            </div>
          </Card>
        )}

        {/* About */}
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              About
            </h2>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">Application:</span>
                <span className="font-medium text-gray-900 dark:text-gray-100">
                  Stumpf.Works NAS
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">Version:</span>
                <span className="font-medium text-gray-900 dark:text-gray-100">
                  0.2.1-alpha
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">License:</span>
                <span className="font-medium text-gray-900 dark:text-gray-100">
                  MIT
                </span>
              </div>
            </div>
          </div>
        </Card>

        {/* Actions */}
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Actions
            </h2>
            <div className="space-y-3">
              <Button variant="danger" onClick={handleLogout} className="w-full">
                Logout
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}
