import { motion } from 'framer-motion';
import { Cpu, MemoryStick, HardDrive, Network, Clock } from 'lucide-react';
import type { SystemMetrics } from '@/api/monitoring';
import Card from '@/components/ui/Card';

interface MetricsOverviewProps {
  metrics: SystemMetrics;
}

export default function MetricsOverview({ metrics }: MetricsOverviewProps) {
  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
  };

  const formatUptime = (seconds: number): string => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    if (days > 0) return `${days}d ${hours}h`;
    if (hours > 0) return `${hours}h ${minutes}m`;
    return `${minutes}m`;
  };

  const getUsageColor = (usage: number): string => {
    if (usage >= 90) return 'bg-red-500';
    if (usage >= 75) return 'bg-yellow-500';
    return 'bg-green-500';
  };

  const metricCards = [
    {
      icon: Cpu,
      title: 'CPU Usage',
      value: `${metrics.cpu_usage_percent.toFixed(1)}%`,
      usage: metrics.cpu_usage_percent,
      details: `Load: ${metrics.load_average_1m.toFixed(2)} / ${metrics.load_average_5m.toFixed(2)} / ${metrics.load_average_15m.toFixed(2)}`,
      color: 'text-blue-500',
    },
    {
      icon: MemoryStick,
      title: 'Memory',
      value: `${metrics.memory_usage_percent.toFixed(1)}%`,
      usage: metrics.memory_usage_percent,
      details: `${metrics.memory_used_mb.toFixed(0)} MB / ${metrics.memory_total_mb.toFixed(0)} MB`,
      color: 'text-purple-500',
    },
    {
      icon: HardDrive,
      title: 'Disk',
      value: `${metrics.disk_usage_percent.toFixed(1)}%`,
      usage: metrics.disk_usage_percent,
      details: `${metrics.disk_used_gb.toFixed(1)} GB / ${metrics.disk_total_gb.toFixed(1)} GB`,
      color: 'text-orange-500',
    },
    {
      icon: Network,
      title: 'Network',
      value: formatBytes(metrics.network_rx_bytes + metrics.network_tx_bytes),
      usage: 0,
      details: `↓ ${formatBytes(metrics.network_rx_bytes)} | ↑ ${formatBytes(metrics.network_tx_bytes)}`,
      color: 'text-green-500',
      hideProgressBar: true,
    },
    {
      icon: Clock,
      title: 'Uptime',
      value: formatUptime(metrics.uptime_seconds),
      usage: 0,
      details: new Date(metrics.timestamp).toLocaleString(),
      color: 'text-gray-500',
      hideProgressBar: true,
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
      {metricCards.map((metric, index) => {
        const Icon = metric.icon;
        return (
          <motion.div
            key={metric.title}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.05 }}
          >
            <Card className="p-4 h-full">
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-2">
                  <Icon className={`w-5 h-5 ${metric.color}`} />
                  <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">{metric.title}</h3>
                </div>
              </div>

              <div className="mb-2">
                <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">{metric.value}</p>
              </div>

              {!metric.hideProgressBar && (
                <div className="mb-2">
                  <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${metric.usage}%` }}
                      transition={{ duration: 0.5 }}
                      className={`h-full ${getUsageColor(metric.usage)}`}
                    />
                  </div>
                </div>
              )}

              <p className="text-xs text-gray-500 dark:text-gray-500">{metric.details}</p>
            </Card>
          </motion.div>
        );
      })}
    </div>
  );
}
