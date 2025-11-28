import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Activity, AlertCircle, RefreshCw } from 'lucide-react';
import { monitoringApi, type SystemMetrics, type HealthScore } from '@/api/monitoring';
import { getErrorMessage } from '@/api/client';
import Card from '@/components/ui/Card';
import MetricsOverview from './components/MetricsOverview';
import HealthScoreCard from './components/HealthScoreCard';
import MetricsCharts from './components/MetricsCharts';

export function Monitoring() {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
  const [healthScore, setHealthScore] = useState<HealthScore | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [autoRefresh, setAutoRefresh] = useState(true);

  useEffect(() => {
    loadData();

    // Auto-refresh every 5 seconds if enabled
    if (autoRefresh) {
      const interval = setInterval(loadData, 5000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  const loadData = async () => {
    try {
      const [metricsResponse, healthResponse] = await Promise.all([
        monitoringApi.getLatestMetrics(),
        monitoringApi.getLatestHealthScore(),
      ]);

      if (metricsResponse.success && metricsResponse.data) {
        setMetrics(metricsResponse.data);
      }

      if (healthResponse.success && healthResponse.data) {
        setHealthScore(healthResponse.data);
      }

      setError('');
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = () => {
    setLoading(true);
    loadData();
  };

  if (loading && !metrics) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="flex flex-col items-center gap-4">
          <RefreshCw className="w-8 h-8 animate-spin text-macos-blue" />
          <p className="text-gray-600 dark:text-gray-400">Loading metrics...</p>
        </div>
      </div>
    );
  }

  if (error && !metrics) {
    return (
      <div className="flex items-center justify-center h-full">
        <Card className="p-8 max-w-md">
          <div className="flex flex-col items-center gap-4 text-center">
            <AlertCircle className="w-12 h-12 text-red-500" />
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">Failed to Load Metrics</h3>
            <p className="text-gray-600 dark:text-gray-400">{error}</p>
            <button
              onClick={handleRefresh}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
            >
              Retry
            </button>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="h-full overflow-auto bg-gray-50 dark:bg-macos-dark-100">
      <div className="p-6 space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Activity className="w-8 h-8 text-macos-blue" />
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">System Monitoring</h1>
              <p className="text-sm text-gray-600 dark:text-gray-400">Real-time system metrics and health status</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
              <input
                type="checkbox"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
                className="rounded"
              />
              Auto-refresh (5s)
            </label>

            <button
              onClick={handleRefresh}
              disabled={loading}
              className="p-2 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-macos-dark-200 rounded-lg transition-colors disabled:opacity-50"
              title="Refresh metrics"
            >
              <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
            </button>
          </div>
        </div>

        {/* Error banner */}
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg"
          >
            <div className="flex items-center gap-2">
              <AlertCircle className="w-5 h-5 text-red-500" />
              <p className="text-sm text-red-700 dark:text-red-400">{error}</p>
            </div>
          </motion.div>
        )}

        {/* Health Score */}
        {healthScore && <HealthScoreCard healthScore={healthScore} />}

        {/* Metrics Overview */}
        {metrics && <MetricsOverview metrics={metrics} />}

        {/* Metrics Charts */}
        {metrics && <MetricsCharts />}
      </div>
    </div>
  );
}
