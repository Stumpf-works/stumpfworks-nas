import { useEffect, useState } from 'react';
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from 'recharts';
import { metricsApi, type SystemMetric, type HealthScore, type MetricsTrend } from '@/api/metrics';
import Card from '@/components/ui/Card';

interface MonitoringWidgetsProps {
  timeRange?: '24h' | '7d' | '30d';
}

export function MonitoringWidgets({ timeRange = '24h' }: MonitoringWidgetsProps) {
  const [metrics, setMetrics] = useState<SystemMetric[]>([]);
  const [healthScore, setHealthScore] = useState<HealthScore | null>(null);
  const [trends, setTrends] = useState<MetricsTrend[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedRange, setSelectedRange] = useState(timeRange);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 60000); // Update every 60 seconds
    return () => clearInterval(interval);
  }, [selectedRange]);

  const fetchData = async () => {
    try {
      setIsLoading(true);

      // Calculate time range
      const end = new Date();
      const start = new Date();
      let limit = 144; // Default for 24h (every 10 min)

      switch (selectedRange) {
        case '24h':
          start.setHours(end.getHours() - 24);
          limit = 144;
          break;
        case '7d':
          start.setDate(end.getDate() - 7);
          limit = 168; // One per hour
          break;
        case '30d':
          start.setDate(end.getDate() - 30);
          limit = 360; // One per 2 hours
          break;
      }

      const [metricsResponse, healthScoreResponse, trendsResponse] = await Promise.all([
        metricsApi.getHistory(start.toISOString(), end.toISOString(), limit),
        metricsApi.getLatestHealthScore(),
        metricsApi.getTrends('1h'),
      ]);

      setMetrics(metricsResponse.metrics.reverse()); // Reverse to show oldest first
      setHealthScore(healthScoreResponse);
      setTrends(trendsResponse.trends);
    } catch (error) {
      console.error('Failed to fetch monitoring data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    if (selectedRange === '24h') {
      return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
    } else if (selectedRange === '7d') {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    } else {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    }
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B/s';
    const k = 1024;
    const sizes = ['B/s', 'KB/s', 'MB/s', 'GB/s'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
  };

  const getTrendIcon = (direction: string) => {
    switch (direction) {
      case 'up':
        return '↑';
      case 'down':
        return '↓';
      default:
        return '→';
    }
  };

  const getTrendColor = (direction: string, metricName: string) => {
    // For usage metrics, up is bad, down is good
    const isUsageMetric = metricName.toLowerCase().includes('usage');

    if (direction === 'up') {
      return isUsageMetric ? 'text-red-500' : 'text-green-500';
    } else if (direction === 'down') {
      return isUsageMetric ? 'text-green-500' : 'text-red-500';
    }
    return 'text-gray-500';
  };

  const getHealthScoreColor = (score: number) => {
    if (score >= 80) return '#10b981'; // green
    if (score >= 60) return '#f59e0b'; // yellow
    if (score >= 40) return '#f97316'; // orange
    return '#ef4444'; // red
  };

  const getHealthScoreLabel = (score: number) => {
    if (score >= 80) return 'Excellent';
    if (score >= 60) return 'Good';
    if (score >= 40) return 'Fair';
    return 'Poor';
  };

  if (isLoading && metrics.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  // Prepare chart data
  const cpuData = metrics.map((m) => ({
    time: formatTimestamp(m.timestamp),
    usage: parseFloat(m.cpuUsage.toFixed(1)),
    load1: parseFloat(m.cpuLoadAvg1.toFixed(2)),
  }));

  const memoryData = metrics.map((m) => ({
    time: formatTimestamp(m.timestamp),
    usage: parseFloat(m.memoryUsage.toFixed(1)),
    swap: parseFloat(m.swapUsage.toFixed(1)),
  }));

  const diskData = metrics.map((m) => ({
    time: formatTimestamp(m.timestamp),
    usage: parseFloat(m.diskUsage.toFixed(1)),
    read: m.diskReadBytesPerSec,
    write: m.diskWriteBytesPerSec,
  }));

  const networkData = metrics.map((m) => ({
    time: formatTimestamp(m.timestamp),
    rx: m.networkRxBytesPerSec,
    tx: m.networkTxBytesPerSec,
  }));

  return (
    <div className="space-y-6">
      {/* Time Range Selector */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
          Advanced Monitoring
        </h2>
        <div className="flex gap-2">
          {(['24h', '7d', '30d'] as const).map((range) => (
            <button
              key={range}
              onClick={() => setSelectedRange(range)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                selectedRange === range
                  ? 'bg-macos-blue text-white'
                  : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
              }`}
            >
              {range}
            </button>
          ))}
        </div>
      </div>

      {/* Health Score & Trends */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Health Score */}
        {healthScore && (
          <Card>
            <div className="p-6">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                System Health Score
              </h3>
              <div className="flex flex-col items-center">
                <div className="relative w-40 h-40">
                  <svg className="transform -rotate-90 w-40 h-40">
                    <circle
                      cx="80"
                      cy="80"
                      r="70"
                      stroke="currentColor"
                      strokeWidth="10"
                      fill="transparent"
                      className="text-gray-200 dark:text-gray-700"
                    />
                    <circle
                      cx="80"
                      cy="80"
                      r="70"
                      stroke={getHealthScoreColor(healthScore.score)}
                      strokeWidth="10"
                      fill="transparent"
                      strokeDasharray={`${(healthScore.score / 100) * 440} 440`}
                      strokeLinecap="round"
                      className="transition-all duration-1000"
                    />
                  </svg>
                  <div className="absolute inset-0 flex flex-col items-center justify-center">
                    <span className="text-4xl font-bold text-gray-900 dark:text-gray-100">
                      {healthScore.score}
                    </span>
                    <span className="text-sm text-gray-600 dark:text-gray-400">
                      {getHealthScoreLabel(healthScore.score)}
                    </span>
                  </div>
                </div>

                <div className="mt-4 grid grid-cols-2 gap-3 w-full">
                  <div className="text-center p-2 bg-gray-50 dark:bg-gray-800 rounded">
                    <div className="text-xs text-gray-600 dark:text-gray-400">CPU</div>
                    <div className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      {healthScore.cpuScore}
                    </div>
                  </div>
                  <div className="text-center p-2 bg-gray-50 dark:bg-gray-800 rounded">
                    <div className="text-xs text-gray-600 dark:text-gray-400">Memory</div>
                    <div className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      {healthScore.memoryScore}
                    </div>
                  </div>
                  <div className="text-center p-2 bg-gray-50 dark:bg-gray-800 rounded">
                    <div className="text-xs text-gray-600 dark:text-gray-400">Disk</div>
                    <div className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      {healthScore.diskScore}
                    </div>
                  </div>
                  <div className="text-center p-2 bg-gray-50 dark:bg-gray-800 rounded">
                    <div className="text-xs text-gray-600 dark:text-gray-400">Network</div>
                    <div className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      {healthScore.networkScore}
                    </div>
                  </div>
                </div>

                {healthScore.issues && (
                  <div className="mt-4 w-full p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded">
                    <div className="text-xs font-semibold text-red-900 dark:text-red-100 mb-1">
                      Issues Detected:
                    </div>
                    <div className="text-xs text-red-800 dark:text-red-200">
                      {healthScore.issues}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </Card>
        )}

        {/* Trends */}
        <Card className="lg:col-span-2">
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Metric Trends (Last Hour)
            </h3>
            <div className="space-y-3">
              {trends.map((trend, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
                >
                  <div className="flex items-center gap-3">
                    <div
                      className={`text-2xl font-bold ${getTrendColor(
                        trend.direction,
                        trend.metricName
                      )}`}
                    >
                      {getTrendIcon(trend.direction)}
                    </div>
                    <div>
                      <div className="font-medium text-gray-900 dark:text-gray-100">
                        {trend.metricName}
                      </div>
                      <div className="text-sm text-gray-600 dark:text-gray-400">
                        {trend.currentValue.toFixed(2)}% (
                        {trend.change > 0 ? '+' : ''}
                        {trend.changePercent.toFixed(1)}%)
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm text-gray-600 dark:text-gray-400">Previous</div>
                    <div className="font-medium text-gray-900 dark:text-gray-100">
                      {trend.previousValue.toFixed(2)}%
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </Card>
      </div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* CPU Usage Chart */}
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              CPU Usage
            </h3>
            <ResponsiveContainer width="100%" height={250}>
              <LineChart data={cpuData}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-gray-300 dark:stroke-gray-700" />
                <XAxis
                  dataKey="time"
                  tick={{ fontSize: 12 }}
                  className="text-gray-600 dark:text-gray-400"
                />
                <YAxis
                  tick={{ fontSize: 12 }}
                  className="text-gray-600 dark:text-gray-400"
                  domain={[0, 100]}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                    border: 'none',
                    borderRadius: '8px',
                    color: '#fff',
                  }}
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="usage"
                  stroke="#3b82f6"
                  strokeWidth={2}
                  dot={false}
                  name="Usage %"
                />
                <Line
                  type="monotone"
                  dataKey="load1"
                  stroke="#10b981"
                  strokeWidth={2}
                  dot={false}
                  name="Load Avg (1m)"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </Card>

        {/* Memory Usage Chart */}
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Memory Usage
            </h3>
            <ResponsiveContainer width="100%" height={250}>
              <LineChart data={memoryData}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-gray-300 dark:stroke-gray-700" />
                <XAxis
                  dataKey="time"
                  tick={{ fontSize: 12 }}
                  className="text-gray-600 dark:text-gray-400"
                />
                <YAxis
                  tick={{ fontSize: 12 }}
                  className="text-gray-600 dark:text-gray-400"
                  domain={[0, 100]}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                    border: 'none',
                    borderRadius: '8px',
                    color: '#fff',
                  }}
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="usage"
                  stroke="#10b981"
                  strokeWidth={2}
                  dot={false}
                  name="Memory %"
                />
                <Line
                  type="monotone"
                  dataKey="swap"
                  stroke="#f59e0b"
                  strokeWidth={2}
                  dot={false}
                  name="Swap %"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </Card>

        {/* Disk Usage Chart */}
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Disk Usage
            </h3>
            <ResponsiveContainer width="100%" height={250}>
              <LineChart data={diskData}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-gray-300 dark:stroke-gray-700" />
                <XAxis
                  dataKey="time"
                  tick={{ fontSize: 12 }}
                  className="text-gray-600 dark:text-gray-400"
                />
                <YAxis
                  tick={{ fontSize: 12 }}
                  className="text-gray-600 dark:text-gray-400"
                  domain={[0, 100]}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                    border: 'none',
                    borderRadius: '8px',
                    color: '#fff',
                  }}
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="usage"
                  stroke="#8b5cf6"
                  strokeWidth={2}
                  dot={false}
                  name="Disk %"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </Card>

        {/* Network Bandwidth Chart */}
        <Card>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              Network Bandwidth
            </h3>
            <ResponsiveContainer width="100%" height={250}>
              <AreaChart data={networkData}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-gray-300 dark:stroke-gray-700" />
                <XAxis
                  dataKey="time"
                  tick={{ fontSize: 12 }}
                  className="text-gray-600 dark:text-gray-400"
                />
                <YAxis
                  tick={{ fontSize: 12 }}
                  className="text-gray-600 dark:text-gray-400"
                  tickFormatter={(value) => formatBytes(value)}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                    border: 'none',
                    borderRadius: '8px',
                    color: '#fff',
                  }}
                  formatter={(value: number) => formatBytes(value)}
                />
                <Legend />
                <Area
                  type="monotone"
                  dataKey="rx"
                  stackId="1"
                  stroke="#3b82f6"
                  fill="#3b82f6"
                  fillOpacity={0.6}
                  name="Received"
                />
                <Area
                  type="monotone"
                  dataKey="tx"
                  stackId="2"
                  stroke="#10b981"
                  fill="#10b981"
                  fillOpacity={0.6}
                  name="Transmitted"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </Card>
      </div>

      {/* Disk I/O Chart */}
      <Card>
        <div className="p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Disk I/O Operations
          </h3>
          <ResponsiveContainer width="100%" height={250}>
            <AreaChart data={diskData}>
              <CartesianGrid strokeDasharray="3 3" className="stroke-gray-300 dark:stroke-gray-700" />
              <XAxis
                dataKey="time"
                tick={{ fontSize: 12 }}
                className="text-gray-600 dark:text-gray-400"
              />
              <YAxis
                tick={{ fontSize: 12 }}
                className="text-gray-600 dark:text-gray-400"
                tickFormatter={(value) => formatBytes(value)}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: 'rgba(0, 0, 0, 0.8)',
                  border: 'none',
                  borderRadius: '8px',
                  color: '#fff',
                }}
                formatter={(value: number) => formatBytes(value)}
              />
              <Legend />
              <Area
                type="monotone"
                dataKey="read"
                stackId="1"
                stroke="#f59e0b"
                fill="#f59e0b"
                fillOpacity={0.6}
                name="Read"
              />
              <Area
                type="monotone"
                dataKey="write"
                stackId="2"
                stroke="#ef4444"
                fill="#ef4444"
                fillOpacity={0.6}
                name="Write"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </Card>
    </div>
  );
}
