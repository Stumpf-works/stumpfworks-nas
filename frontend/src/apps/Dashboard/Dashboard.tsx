import { useEffect } from 'react';
import { useSystemStore } from '@/store';
import { systemApi } from '@/api/system';
import Card from '@/components/ui/Card';

export function Dashboard() {
  const metrics = useSystemStore((state) => state.metrics);
  const setMetrics = useSystemStore((state) => state.setMetrics);
  const isLoading = useSystemStore((state) => state.isLoading);
  const setLoading = useSystemStore((state) => state.setLoading);

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
    </div>
  );
}
