import { useEffect, useState } from 'react';
import { useSystemStore } from '@/store';
import { systemApi } from '@/api/system';
import { securityApi } from '@/api/security';
import { auditApi, type AuditLog } from '@/api/audit';
import Card from '@/components/ui/Card';
import { MonitoringWidgets } from '@/components/MonitoringWidgets/MonitoringWidgets';

export function Dashboard() {
  const metrics = useSystemStore((state) => state.metrics);
  const setMetrics = useSystemStore((state) => state.setMetrics);
  const isLoading = useSystemStore((state) => state.isLoading);
  const setLoading = useSystemStore((state) => state.setLoading);

  // Security stats state
  const [securityStats, setSecurityStats] = useState<any>(null);
  const [recentCriticalLogs, setRecentCriticalLogs] = useState<AuditLog[]>([]);

  useEffect(() => {
    const fetchMetrics = async () => {
      setLoading(true);
      try {
        const response = await systemApi.getMetrics();
        if (response.success && response.data) {
          setMetrics(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch metrics:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 3000);
    return () => clearInterval(interval);
  }, [setMetrics, setLoading]);

  // Fetch security stats
  useEffect(() => {
    const fetchSecurityStats = async () => {
      try {
        const [statsResponse, logsResponse] = await Promise.all([
          securityApi.getStats(),
          auditApi.listLogs({ severity: 'critical', limit: 5 }),
        ]);

        if (statsResponse.success && statsResponse.data) {
          setSecurityStats(statsResponse.data);
        }

        if (logsResponse.success && logsResponse.data) {
          setRecentCriticalLogs(logsResponse.data.logs);
        }
      } catch (error) {
        console.error('Failed to fetch security stats:', error);
      }
    };

    fetchSecurityStats();
    const interval = setInterval(fetchSecurityStats, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
  };

  if (isLoading && !metrics) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  if (!metrics) {
    return (
      <div className="flex items-center justify-center h-full text-gray-500">
        No metrics available
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6 overflow-auto h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          System Dashboard
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Real-time system metrics and monitoring
        </p>
      </div>

      {/* Metrics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* CPU */}
        <Card>
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                CPU Usage
              </h3>
              <svg
                className="w-8 h-8 text-macos-blue"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"
                />
              </svg>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-3xl font-bold text-gray-900 dark:text-gray-100">
                  {metrics.cpu.usagePercent.toFixed(1)}%
                </span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-macos-blue rounded-full h-2 transition-all duration-500"
                  style={{ width: `${metrics.cpu.usagePercent}%` }}
                />
              </div>
              {metrics.cpu.perCore && (
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  {metrics.cpu.perCore.length} cores
                </div>
              )}
            </div>
          </div>
        </Card>

        {/* Memory */}
        <Card>
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                Memory
              </h3>
              <svg
                className="w-8 h-8 text-macos-green"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                />
              </svg>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-3xl font-bold text-gray-900 dark:text-gray-100">
                  {metrics.memory.usedPercent.toFixed(1)}%
                </span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-macos-green rounded-full h-2 transition-all duration-500"
                  style={{ width: `${metrics.memory.usedPercent}%` }}
                />
              </div>
              <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400">
                <span>{formatBytes(metrics.memory.used)}</span>
                <span>{formatBytes(metrics.memory.total)}</span>
              </div>
            </div>
          </div>
        </Card>

        {/* Network */}
        <Card>
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                Network
              </h3>
              <svg
                className="w-8 h-8 text-macos-purple"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0"
                />
              </svg>
            </div>
            <div className="space-y-3">
              <div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Sent</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {formatBytes(metrics.network.bytesSent)}
                  </span>
                </div>
              </div>
              <div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Received</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {formatBytes(metrics.network.bytesRecv)}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </Card>
      </div>

      {/* Security Overview */}
      {securityStats && (
        <div>
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            Security Overview
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Failed Logins (Last 24h) */}
            <Card>
              <div className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                    Failed Logins
                  </h3>
                  <svg
                    className="w-6 h-6 text-macos-red"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                    />
                  </svg>
                </div>
                <div className="space-y-1">
                  <div className="text-3xl font-bold text-gray-900 dark:text-gray-100">
                    {securityStats.last24hAttempts || 0}
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    Last 24 hours
                  </div>
                </div>
              </div>
            </Card>

            {/* Blocked IPs */}
            <Card>
              <div className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                    Blocked IPs
                  </h3>
                  <svg
                    className="w-6 h-6 text-macos-orange"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                    />
                  </svg>
                </div>
                <div className="space-y-1">
                  <div className="text-3xl font-bold text-gray-900 dark:text-gray-100">
                    {securityStats.blockedIPsCount || 0}
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    Currently blocked
                  </div>
                </div>
              </div>
            </Card>

            {/* Total Attempts */}
            <Card>
              <div className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                    Total Attempts
                  </h3>
                  <svg
                    className="w-6 h-6 text-macos-blue"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                </div>
                <div className="space-y-1">
                  <div className="text-3xl font-bold text-gray-900 dark:text-gray-100">
                    {securityStats.totalAttempts || 0}
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    All time
                  </div>
                </div>
              </div>
            </Card>

            {/* Security Status */}
            <Card>
              <div className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                    Security Status
                  </h3>
                  <svg
                    className={`w-6 h-6 ${
                      (securityStats.blockedIPsCount || 0) > 0 ||
                      (securityStats.last24hAttempts || 0) > 5
                        ? 'text-macos-orange'
                        : 'text-macos-green'
                    }`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                    />
                  </svg>
                </div>
                <div className="space-y-1">
                  <div
                    className={`text-2xl font-bold ${
                      (securityStats.blockedIPsCount || 0) > 0 ||
                      (securityStats.last24hAttempts || 0) > 5
                        ? 'text-macos-orange'
                        : 'text-macos-green'
                    }`}
                  >
                    {(securityStats.blockedIPsCount || 0) > 0 ||
                    (securityStats.last24hAttempts || 0) > 5
                      ? 'Active Threats'
                      : 'Secure'}
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    System protection
                  </div>
                </div>
              </div>
            </Card>
          </div>
        </div>
      )}

      {/* Recent Critical Events */}
      {recentCriticalLogs.length > 0 && (
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Recent Critical Events
            </h3>
            <div className="space-y-3">
              {recentCriticalLogs.map((log) => (
                <div
                  key={log.id}
                  className="flex items-start justify-between p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg"
                >
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="font-medium text-red-900 dark:text-red-100">
                        {log.action}
                      </span>
                      <span
                        className={`px-2 py-0.5 text-xs rounded-full ${
                          log.status === 'success'
                            ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200'
                            : log.status === 'failure'
                            ? 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-200'
                            : 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200'
                        }`}
                      >
                        {log.status}
                      </span>
                    </div>
                    <div className="text-sm text-red-800 dark:text-red-200">
                      <span className="font-medium">{log.username}</span>
                      {log.ipAddress && (
                        <span className="text-red-700 dark:text-red-300">
                          {' '}
                          from {log.ipAddress}
                        </span>
                      )}
                    </div>
                    {log.message && (
                      <div className="text-xs text-red-700 dark:text-red-300 mt-1">
                        {log.message}
                      </div>
                    )}
                  </div>
                  <div className="text-xs text-red-600 dark:text-red-400 ml-4">
                    {new Date(log.createdAt).toLocaleTimeString()}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </Card>
      )}

      {/* Disk Usage */}
      {metrics.disk.length > 0 && (
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Disk Usage
            </h3>
            <div className="space-y-4">
              {metrics.disk.map((disk, index) => (
                <div key={index} className="space-y-2">
                  <div className="flex items-center justify-between text-sm">
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {disk.mountpoint}
                    </span>
                    <span className="text-gray-600 dark:text-gray-400">
                      {disk.usedPercent.toFixed(1)}% used
                    </span>
                  </div>
                  <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                    <div
                      className={`rounded-full h-2 transition-all duration-500 ${
                        disk.usedPercent > 90
                          ? 'bg-macos-red'
                          : disk.usedPercent > 70
                          ? 'bg-macos-orange'
                          : 'bg-macos-blue'
                      }`}
                      style={{ width: `${disk.usedPercent}%` }}
                    />
                  </div>
                  <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400">
                    <span>
                      {formatBytes(disk.used)} / {formatBytes(disk.total)}
                    </span>
                    <span>{formatBytes(disk.free)} free</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </Card>
      )}

      {/* Advanced Monitoring Widgets */}
      <MonitoringWidgets />
    </div>
  );
}
