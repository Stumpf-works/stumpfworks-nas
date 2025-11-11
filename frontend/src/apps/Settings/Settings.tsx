import { useState, useEffect } from 'react';
import { useAuthStore, useThemeStore } from '@/store';
import { systemApi } from '@/api/system';
import { authApi } from '@/api/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

export function Settings() {
  const user = useAuthStore((state) => state.user);
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const isDark = useThemeStore((state) => state.isDark);
  const toggleTheme = useThemeStore((state) => state.toggleTheme);

  const [systemInfo, setSystemInfo] = useState<any>(null);
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

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
    fetchSystemInfo();
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
                  0.1.0-alpha
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
