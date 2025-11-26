import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { TrendingUp, AlertCircle } from 'lucide-react';
import { monitoringApi, type SystemMetrics } from '@/api/monitoring';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';

export default function MetricsCharts() {
  const [history, setHistory] = useState<SystemMetrics[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadHistory();
  }, []);

  const loadHistory = async () => {
    try {
      // Get last hour of data
      const end = new Date();
      const start = new Date(end.getTime() - 60 * 60 * 1000); // 1 hour ago

      const response = await monitoringApi.getMetricsHistory({
        start: start.toISOString(),
        end: end.toISOString(),
        limit: 100,
      });

      if (response.success && response.data?.metrics) {
        setHistory(response.data.metrics);
      }

      setError('');
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center h-40">
          <p className="text-gray-600 dark:text-gray-400">Loading historical data...</p>
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="p-6">
        <div className="flex items-center gap-2 text-red-500">
          <AlertCircle className="w-5 h-5" />
          <p>{error}</p>
        </div>
      </Card>
    );
  }

  if (history.length === 0) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center h-40">
          <p className="text-gray-600 dark:text-gray-400">No historical data available yet</p>
        </div>
      </Card>
    );
  }

  // Calculate statistics
  const avgCpu = history.reduce((sum, m) => sum + m.cpu_usage_percent, 0) / history.length;
  const avgMem = history.reduce((sum, m) => sum + m.memory_usage_percent, 0) / history.length;
  const avgDisk = history.reduce((sum, m) => sum + m.disk_usage_percent, 0) / history.length;

  const maxCpu = Math.max(...history.map((m) => m.cpu_usage_percent));
  const maxMem = Math.max(...history.map((m) => m.memory_usage_percent));
  const maxDisk = Math.max(...history.map((m) => m.disk_usage_percent));

  const stats = [
    { label: 'CPU', avg: avgCpu, max: maxCpu, color: 'text-blue-500' },
    { label: 'Memory', avg: avgMem, max: maxMem, color: 'text-purple-500' },
    { label: 'Disk', avg: avgDisk, max: maxDisk, color: 'text-orange-500' },
  ];

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.3 }}
    >
      <Card className="p-6">
        <div className="flex items-center gap-2 mb-6">
          <TrendingUp className="w-6 h-6 text-macos-blue" />
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
            Metrics Trends (Last Hour)
          </h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {stats.map((stat) => (
            <div key={stat.label} className="space-y-2">
              <div className="flex items-center justify-between">
                <span className={`text-sm font-medium ${stat.color}`}>{stat.label}</span>
                <span className="text-xs text-gray-500 dark:text-gray-500">{history.length} samples</span>
              </div>

              <div className="space-y-1">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Average</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">{stat.avg.toFixed(1)}%</span>
                </div>
                <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                  <div
                    className={`h-full bg-gradient-to-r ${stat.color.replace('text-', 'from-')} to-transparent`}
                    style={{ width: `${stat.avg}%` }}
                  />
                </div>
              </div>

              <div className="space-y-1">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Peak</span>
                  <span className="font-medium text-gray-900 dark:text-gray-100">{stat.max.toFixed(1)}%</span>
                </div>
                <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                  <div
                    className={`h-full ${stat.color.replace('text-', 'bg-')}`}
                    style={{ width: `${stat.max}%` }}
                  />
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
          <p className="text-xs text-gray-500 dark:text-gray-500 text-center">
            Showing {history.length} data points from the last hour • Auto-refreshes with parent component
          </p>
        </div>
      </Card>
    </motion.div>
  );
}
