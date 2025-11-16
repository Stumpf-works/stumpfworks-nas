// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState, useEffect } from 'react';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { adApi, type ADConfig } from '@/api/ad';

export function ActiveDirectorySection({ user, systemInfo }: { user: any; systemInfo: any }) {
  const [adConfig, setAdConfig] = useState<ADConfig | null>(null);
  const [adConfigEditing, setAdConfigEditing] = useState(false);
  const [adTestResult, setAdTestResult] = useState<{ success: boolean; message: string } | null>(null);
  const [adLoading, setAdLoading] = useState(false);
  const [showAdSetup, setShowAdSetup] = useState(false);

  useEffect(() => {
    fetchAdConfig();
  }, []);

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

  const handleSetupAd = () => {
    const defaultConfig: ADConfig = {
      enabled: false,
      server: '',
      port: 389,
      baseDN: '',
      bindUser: '',
      bindPassword: '',
      userFilter: '(&(objectClass=user)(sAMAccountName=%s))',
      groupFilter: '(&(objectClass=group)(member=%s))',
      useTLS: false,
      skipVerify: false,
    };
    setAdConfig(defaultConfig);
    setAdConfigEditing(true);
    setShowAdSetup(false);
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

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Active Directory</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure Active Directory integration for user authentication
        </p>
      </div>

      <Card>
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Active Directory Configuration
            </h2>
            {!adConfig && !showAdSetup && (
              <Button variant="primary" onClick={handleSetupAd} size="sm">
                Setup AD
              </Button>
            )}
            {adConfig && !adConfigEditing && (
              <Button variant="secondary" onClick={() => setAdConfigEditing(true)} size="sm">
                Edit
              </Button>
            )}
          </div>

          {!adConfig && !showAdSetup ? (
            <div className="text-center py-8">
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Active Directory integration is not configured yet.
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-500 mb-4">
                Configure AD to enable user authentication and synchronization with your directory services.
              </p>
              <Button variant="primary" onClick={handleSetupAd}>
                Configure Active Directory
              </Button>
            </div>
          ) : (
            <div className="space-y-4">
              {/* Enabled Toggle */}
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-900 dark:text-gray-100">AD Integration</p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Enable Active Directory authentication
                  </p>
                </div>
                <button
                  onClick={() =>
                    adConfigEditing &&
                    adConfig &&
                    setAdConfig({ ...adConfig, enabled: !adConfig.enabled })
                  }
                  disabled={!adConfigEditing}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                    adConfig?.enabled ? 'bg-macos-blue' : 'bg-gray-300'
                  } ${!adConfigEditing ? 'opacity-50 cursor-not-allowed' : ''}`}
                >
                  <span
                    className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                      adConfig?.enabled ? 'translate-x-6' : 'translate-x-1'
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
                      value={adConfig?.server || ''}
                      onChange={(e) =>
                        adConfig && setAdConfig({ ...adConfig, server: e.target.value })
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
                      value={adConfig?.port || 389}
                      onChange={(e) =>
                        adConfig && setAdConfig({ ...adConfig, port: parseInt(e.target.value) })
                      }
                      placeholder="389"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Base DN
                    </label>
                    <Input
                      value={adConfig?.baseDN || ''}
                      onChange={(e) =>
                        adConfig && setAdConfig({ ...adConfig, baseDN: e.target.value })
                      }
                      placeholder="dc=example,dc=com"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Bind User
                    </label>
                    <Input
                      value={adConfig?.bindUser || ''}
                      onChange={(e) =>
                        adConfig && setAdConfig({ ...adConfig, bindUser: e.target.value })
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
                      value={adConfig?.bindPassword || ''}
                      onChange={(e) =>
                        adConfig && setAdConfig({ ...adConfig, bindPassword: e.target.value })
                      }
                      placeholder="••••••••"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      User Filter
                    </label>
                    <Input
                      value={adConfig?.userFilter || ''}
                      onChange={(e) =>
                        adConfig && setAdConfig({ ...adConfig, userFilter: e.target.value })
                      }
                      placeholder="(&(objectClass=user)(sAMAccountName=%s))"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Group Filter
                    </label>
                    <Input
                      value={adConfig?.groupFilter || ''}
                      onChange={(e) =>
                        adConfig && setAdConfig({ ...adConfig, groupFilter: e.target.value })
                      }
                      placeholder="(&(objectClass=group)(member=%s))"
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium text-gray-900 dark:text-gray-100">Use TLS</p>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        Enable secure connection
                      </p>
                    </div>
                    <button
                      onClick={() =>
                        adConfig && setAdConfig({ ...adConfig, useTLS: !adConfig.useTLS })
                      }
                      className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                        adConfig?.useTLS ? 'bg-macos-blue' : 'bg-gray-300'
                      }`}
                    >
                      <span
                        className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                          adConfig?.useTLS ? 'translate-x-6' : 'translate-x-1'
                        }`}
                      />
                    </button>
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium text-gray-900 dark:text-gray-100">
                        Skip TLS Verification
                      </p>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        Skip certificate verification (use for self-signed certs)
                      </p>
                    </div>
                    <button
                      onClick={() =>
                        adConfig && setAdConfig({ ...adConfig, skipVerify: !adConfig.skipVerify })
                      }
                      className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                        adConfig?.skipVerify ? 'bg-macos-blue' : 'bg-gray-300'
                      }`}
                    >
                      <span
                        className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                          adConfig?.skipVerify ? 'translate-x-6' : 'translate-x-1'
                        }`}
                      />
                    </button>
                  </div>
                </>
              ) : (
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <span className="text-sm text-gray-600 dark:text-gray-400">Server:</span>
                    <p className="font-medium text-gray-900 dark:text-gray-100">
                      {adConfig?.server}:{adConfig?.port}
                    </p>
                  </div>
                  <div>
                    <span className="text-sm text-gray-600 dark:text-gray-400">Base DN:</span>
                    <p className="font-medium text-gray-900 dark:text-gray-100">
                      {adConfig?.baseDN}
                    </p>
                  </div>
                  <div>
                    <span className="text-sm text-gray-600 dark:text-gray-400">TLS:</span>
                    <p className="font-medium text-gray-900 dark:text-gray-100">
                      {adConfig?.useTLS ? 'Enabled' : 'Disabled'}
                    </p>
                  </div>
                </div>
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
          )}
        </div>
      </Card>
    </div>
  );
}
