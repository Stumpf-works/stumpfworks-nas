import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { storageApi, StorageStats, DiskHealth } from '@/api/storage';
import Card from '@/components/ui/Card';

export default function StorageOverview() {
  const [stats, setStats] = useState<StorageStats | null>(null);
  const [health, setHealth] = useState<DiskHealth[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 10000); // Refresh every 10s
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [statsRes, healthRes] = await Promise.all([
        storageApi.getStats(),
        storageApi.getHealth(),
      ]);

      if (statsRes.success) setStats(statsRes.data);
      if (healthRes.success) setHealth(healthRes.data);
    } catch (error) {
      console.error('Failed to load storage data:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="text-center text-gray-600 dark:text-gray-400">
        Failed to load storage statistics
      </div>
    );
  }

  const usagePercent = stats.totalCapacity > 0
    ? (stats.usedCapacity / stats.totalCapacity) * 100
    : 0;

  return (
    <div className="space-y-6">
      {/* Capacity Overview */}
      <Card>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
          Storage Capacity
        </h3>
        <div className="space-y-4">
          <div className="flex justify-between text-sm">
            <span className="text-gray-600 dark:text-gray-400">Total</span>
            <span className="font-medium text-gray-900 dark:text-gray-100">
              {formatBytes(stats.totalCapacity)}
            </span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-600 dark:text-gray-400">Used</span>
            <span className="font-medium text-gray-900 dark:text-gray-100">
              {formatBytes(stats.usedCapacity)}
            </span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-600 dark:text-gray-400">Available</span>
            <span className="font-medium text-gray-900 dark:text-gray-100">
              {formatBytes(stats.availableCapacity)}
            </span>
          </div>

          {/* Usage Bar */}
          <div className="relative pt-1">
            <div className="overflow-hidden h-4 text-xs flex rounded-full bg-gray-200 dark:bg-gray-700">
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: `${usagePercent}%` }}
                transition={{ duration: 1, ease: 'easeOut' }}
                className={`flex flex-col text-center whitespace-nowrap text-white justify-center ${
                  usagePercent > 90
                    ? 'bg-red-500'
                    : usagePercent > 75
                    ? 'bg-yellow-500'
                    : 'bg-macos-blue'
                }`}
              />
            </div>
            <div className="text-center text-sm mt-2 font-medium text-gray-700 dark:text-gray-300">
              {usagePercent.toFixed(1)}% Used
            </div>
          </div>
        </div>
      </Card>

      {/* Statistics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Total Disks"
          value={stats.totalDisks}
          icon="💿"
          color="blue"
        />
        <StatCard
          title="Volumes"
          value={stats.totalVolumes}
          icon="📦"
          color="purple"
        />
        <StatCard
          title="Shares"
          value={stats.totalShares}
          icon="📁"
          color="green"
        />
        <StatCard
          title="Healthy Disks"
          value={stats.healthyDisks}
          icon="✅"
          color="green"
        />
      </div>

      {/* Disk Health Status */}
      {stats.warningDisks > 0 || stats.criticalDisks > 0 ? (
        <Card>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            ⚠️ Disk Health Alerts
          </h3>
          <div className="space-y-3">
            {stats.warningDisks > 0 && (
              <div className="p-3 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
                <span className="font-medium text-yellow-700 dark:text-yellow-400">
                  {stats.warningDisks} disk(s) with warnings
                </span>
              </div>
            )}
            {stats.criticalDisks > 0 && (
              <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
                <span className="font-medium text-red-700 dark:text-red-400">
                  {stats.criticalDisks} disk(s) in critical condition
                </span>
              </div>
            )}
          </div>
        </Card>
      ) : null}

      {/* Recent Health Issues */}
      {health.filter((h) => h.issues.length > 0).length > 0 && (
        <Card>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Health Issues
          </h3>
          <div className="space-y-3">
            {health
              .filter((h) => h.issues.length > 0)
              .map((disk) => (
                <div
                  key={disk.diskName}
                  className="p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                >
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {disk.diskName}
                    </span>
                    <StatusBadge status={disk.status} />
                  </div>
                  <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                    {disk.issues.map((issue, index) => (
                      <li key={index}>• {issue}</li>
                    ))}
                  </ul>
                </div>
              ))}
          </div>
        </Card>
      )}
    </div>
  );
}

interface StatCardProps {
  title: string;
  value: number;
  icon: string;
  color: 'blue' | 'purple' | 'green' | 'yellow' | 'red';
}

function StatCard({ title, value, icon, color }: StatCardProps) {
  const colorClasses = {
    blue: 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400',
    purple: 'bg-purple-50 dark:bg-purple-900/20 text-purple-600 dark:text-purple-400',
    green: 'bg-green-50 dark:bg-green-900/20 text-green-600 dark:text-green-400',
    yellow: 'bg-yellow-50 dark:bg-yellow-900/20 text-yellow-600 dark:text-yellow-400',
    red: 'bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400',
  };

  return (
    <Card>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm text-gray-600 dark:text-gray-400">{title}</p>
          <p className="text-2xl font-bold text-gray-900 dark:text-gray-100 mt-1">
            {value}
          </p>
        </div>
        <div className={`text-3xl p-3 rounded-lg ${colorClasses[color]}`}>
          {icon}
        </div>
      </div>
    </Card>
  );
}

function StatusBadge({ status }: { status: string }) {
  const colorClasses = {
    healthy: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
    warning: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
    critical: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
    failed: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
    unknown: 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400',
  };

  return (
    <span
      className={`px-2 py-1 rounded-full text-xs font-medium ${
        colorClasses[status as keyof typeof colorClasses] || colorClasses.unknown
      }`}
    >
      {status}
    </span>
  );
}
