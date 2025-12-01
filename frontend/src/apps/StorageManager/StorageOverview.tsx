// Revision: 2025-12-01 | Author: StumpfWorks AI | Version: 2.0.0
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  HardDrive,
  Database,
  Share2,
  Activity,
  AlertTriangle,
  CheckCircle,
  XCircle,
  TrendingUp,
  Thermometer,
  Shield,
} from 'lucide-react';
import { storageApi, StorageStats, DiskHealth } from '@/api/storage';

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

      if (statsRes.success && statsRes.data) setStats(statsRes.data);
      if (healthRes.success && healthRes.data) setHealth(healthRes.data);
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
      {/* Capacity Overview - Hero Card */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-macos-blue/10 via-macos-purple/10 to-macos-pink/10 dark:from-macos-blue/20 dark:via-macos-purple/20 dark:to-macos-pink/20 p-8 backdrop-blur-sm border border-white/20 dark:border-gray-700/30"
      >
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Left: Circular Progress */}
          <div className="flex flex-col items-center justify-center">
            <div className="relative w-48 h-48">
              {/* Background Circle */}
              <svg className="w-48 h-48 transform -rotate-90">
                <circle
                  cx="96"
                  cy="96"
                  r="88"
                  stroke="currentColor"
                  strokeWidth="12"
                  fill="none"
                  className="text-gray-200 dark:text-gray-700"
                />
                {/* Progress Circle */}
                <motion.circle
                  cx="96"
                  cy="96"
                  r="88"
                  stroke="url(#gradient)"
                  strokeWidth="12"
                  fill="none"
                  strokeLinecap="round"
                  initial={{ strokeDasharray: '0 600' }}
                  animate={{
                    strokeDasharray: `${(usagePercent / 100) * 552} 552`,
                  }}
                  transition={{ duration: 1.5, ease: 'easeOut' }}
                />
                <defs>
                  <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop
                      offset="0%"
                      stopColor={usagePercent > 90 ? '#ef4444' : usagePercent > 75 ? '#f59e0b' : '#3b82f6'}
                    />
                    <stop
                      offset="100%"
                      stopColor={usagePercent > 90 ? '#dc2626' : usagePercent > 75 ? '#d97706' : '#8b5cf6'}
                    />
                  </linearGradient>
                </defs>
              </svg>
              {/* Center Text */}
              <div className="absolute inset-0 flex flex-col items-center justify-center">
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  transition={{ delay: 0.5, type: 'spring' }}
                  className="text-center"
                >
                  <div className="text-4xl font-bold text-gray-900 dark:text-gray-100">
                    {usagePercent.toFixed(0)}%
                  </div>
                  <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">Used</div>
                </motion.div>
              </div>
            </div>
          </div>

          {/* Right: Details */}
          <div className="flex flex-col justify-center space-y-6">
            <div>
              <h2 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">
                Storage Capacity
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                System-wide storage usage overview
              </p>
            </div>

            <div className="space-y-4">
              <div className="flex items-center justify-between p-4 bg-white/50 dark:bg-macos-dark-200/50 rounded-xl backdrop-blur-sm">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-macos-blue/10 rounded-lg">
                    <Database className="w-5 h-5 text-macos-blue" />
                  </div>
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                    Total Capacity
                  </span>
                </div>
                <span className="text-lg font-bold text-gray-900 dark:text-gray-100">
                  {formatBytes(stats.totalCapacity)}
                </span>
              </div>

              <div className="flex items-center justify-between p-4 bg-white/50 dark:bg-macos-dark-200/50 rounded-xl backdrop-blur-sm">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-macos-purple/10 rounded-lg">
                    <Activity className="w-5 h-5 text-macos-purple" />
                  </div>
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                    Used Space
                  </span>
                </div>
                <span className="text-lg font-bold text-gray-900 dark:text-gray-100">
                  {formatBytes(stats.usedCapacity)}
                </span>
              </div>

              <div className="flex items-center justify-between p-4 bg-white/50 dark:bg-macos-dark-200/50 rounded-xl backdrop-blur-sm">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-green-500/10 rounded-lg">
                    <TrendingUp className="w-5 h-5 text-green-500" />
                  </div>
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                    Available Space
                  </span>
                </div>
                <span className="text-lg font-bold text-gray-900 dark:text-gray-100">
                  {formatBytes(stats.availableCapacity)}
                </span>
              </div>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Statistics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Total Disks"
          value={stats.totalDisks}
          icon={<HardDrive className="w-6 h-6" />}
          color="blue"
          delay={0.1}
        />
        <StatCard
          title="Volumes"
          value={stats.totalVolumes}
          icon={<Database className="w-6 h-6" />}
          color="purple"
          delay={0.2}
        />
        <StatCard
          title="Shares"
          value={stats.totalShares}
          icon={<Share2 className="w-6 h-6" />}
          color="indigo"
          delay={0.3}
        />
        <StatCard
          title="Healthy Disks"
          value={stats.healthyDisks}
          icon={<CheckCircle className="w-6 h-6" />}
          color="green"
          delay={0.4}
        />
      </div>

      {/* Disk Health Status */}
      {stats.warningDisks > 0 || stats.criticalDisks > 0 ? (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="bg-gradient-to-br from-red-50/50 to-yellow-50/50 dark:from-red-900/10 dark:to-yellow-900/10 rounded-2xl p-6 border border-red-200/30 dark:border-red-800/30 backdrop-blur-sm"
        >
          <div className="flex items-center gap-3 mb-4">
            <div className="p-3 bg-red-500/10 rounded-xl">
              <AlertTriangle className="w-6 h-6 text-red-600 dark:text-red-400" />
            </div>
            <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
              Disk Health Alerts
            </h3>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {stats.warningDisks > 0 && (
              <motion.div
                initial={{ scale: 0.9, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ delay: 0.6 }}
                className="relative overflow-hidden p-5 bg-gradient-to-br from-yellow-100/80 to-yellow-50/80 dark:from-yellow-900/30 dark:to-yellow-800/20 border-2 border-yellow-300/50 dark:border-yellow-700/30 rounded-xl"
              >
                <div className="flex items-center gap-3">
                  <Shield className="w-5 h-5 text-yellow-600 dark:text-yellow-400" />
                  <div>
                    <div className="text-2xl font-bold text-yellow-700 dark:text-yellow-300">
                      {stats.warningDisks}
                    </div>
                    <div className="text-sm font-medium text-yellow-600 dark:text-yellow-400">
                      Warning(s)
                    </div>
                  </div>
                </div>
              </motion.div>
            )}
            {stats.criticalDisks > 0 && (
              <motion.div
                initial={{ scale: 0.9, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ delay: 0.7 }}
                className="relative overflow-hidden p-5 bg-gradient-to-br from-red-100/80 to-red-50/80 dark:from-red-900/30 dark:to-red-800/20 border-2 border-red-300/50 dark:border-red-700/30 rounded-xl"
              >
                <div className="flex items-center gap-3">
                  <XCircle className="w-5 h-5 text-red-600 dark:text-red-400" />
                  <div>
                    <div className="text-2xl font-bold text-red-700 dark:text-red-300">
                      {stats.criticalDisks}
                    </div>
                    <div className="text-sm font-medium text-red-600 dark:text-red-400">
                      Critical
                    </div>
                  </div>
                </div>
              </motion.div>
            )}
          </div>
        </motion.div>
      ) : null}

      {/* Recent Health Issues */}
      {health.filter((h) => h.issues.length > 0).length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6 }}
          className="bg-white/50 dark:bg-macos-dark-100/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-200/50 dark:border-gray-700/50"
        >
          <div className="flex items-center gap-3 mb-6">
            <div className="p-3 bg-orange-500/10 rounded-xl">
              <Activity className="w-6 h-6 text-orange-600 dark:text-orange-400" />
            </div>
            <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
              Disk Health Details
            </h3>
          </div>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {health
              .filter((h) => h.issues.length > 0)
              .map((disk, index) => (
                <motion.div
                  key={disk.diskName}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.7 + index * 0.1 }}
                  className="relative overflow-hidden p-5 bg-gradient-to-br from-gray-50 to-gray-100/50 dark:from-macos-dark-200 dark:to-macos-dark-300/50 rounded-xl border border-gray-200/50 dark:border-gray-700/50 hover:shadow-lg transition-shadow"
                >
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <HardDrive className="w-5 h-5 text-gray-600 dark:text-gray-400" />
                      <span className="font-bold text-gray-900 dark:text-gray-100">
                        {disk.diskName}
                      </span>
                    </div>
                    <StatusBadge status={disk.status} />
                  </div>

                  {/* Temperature if available */}
                  {disk.temperature > 0 && (
                    <div className="flex items-center gap-2 mb-3 p-2 bg-white/50 dark:bg-macos-dark-100/50 rounded-lg">
                      <Thermometer className="w-4 h-4 text-orange-500" />
                      <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                        {disk.temperature}°C
                      </span>
                    </div>
                  )}

                  <ul className="space-y-2">
                    {disk.issues.map((issue, issueIndex) => (
                      <li
                        key={issueIndex}
                        className="flex items-start gap-2 text-sm text-gray-600 dark:text-gray-400"
                      >
                        <span className="text-red-500 mt-0.5">•</span>
                        <span>{issue}</span>
                      </li>
                    ))}
                  </ul>
                </motion.div>
              ))}
          </div>
        </motion.div>
      )}
    </div>
  );
}

interface StatCardProps {
  title: string;
  value: number;
  icon: React.ReactNode;
  color: 'blue' | 'purple' | 'green' | 'yellow' | 'red' | 'indigo';
  delay?: number;
}

function StatCard({ title, value, icon, color, delay = 0 }: StatCardProps) {
  const colorClasses = {
    blue: {
      bg: 'bg-gradient-to-br from-blue-500/10 to-blue-600/10 dark:from-blue-500/20 dark:to-blue-600/20',
      icon: 'bg-blue-500/10 text-blue-600 dark:text-blue-400',
      border: 'border-blue-200/50 dark:border-blue-800/30',
    },
    purple: {
      bg: 'bg-gradient-to-br from-purple-500/10 to-purple-600/10 dark:from-purple-500/20 dark:to-purple-600/20',
      icon: 'bg-purple-500/10 text-purple-600 dark:text-purple-400',
      border: 'border-purple-200/50 dark:border-purple-800/30',
    },
    indigo: {
      bg: 'bg-gradient-to-br from-indigo-500/10 to-indigo-600/10 dark:from-indigo-500/20 dark:to-indigo-600/20',
      icon: 'bg-indigo-500/10 text-indigo-600 dark:text-indigo-400',
      border: 'border-indigo-200/50 dark:border-indigo-800/30',
    },
    green: {
      bg: 'bg-gradient-to-br from-green-500/10 to-green-600/10 dark:from-green-500/20 dark:to-green-600/20',
      icon: 'bg-green-500/10 text-green-600 dark:text-green-400',
      border: 'border-green-200/50 dark:border-green-800/30',
    },
    yellow: {
      bg: 'bg-gradient-to-br from-yellow-500/10 to-yellow-600/10 dark:from-yellow-500/20 dark:to-yellow-600/20',
      icon: 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400',
      border: 'border-yellow-200/50 dark:border-yellow-800/30',
    },
    red: {
      bg: 'bg-gradient-to-br from-red-500/10 to-red-600/10 dark:from-red-500/20 dark:to-red-600/20',
      icon: 'bg-red-500/10 text-red-600 dark:text-red-400',
      border: 'border-red-200/50 dark:border-red-800/30',
    },
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay, duration: 0.3 }}
      className={`relative overflow-hidden rounded-xl ${colorClasses[color].bg} ${colorClasses[color].border} border backdrop-blur-sm p-6 hover:scale-[1.02] transition-transform cursor-pointer group`}
    >
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <p className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">
            {title}
          </p>
          <motion.p
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: delay + 0.2, type: 'spring' }}
            className="text-3xl font-bold text-gray-900 dark:text-gray-100"
          >
            {value}
          </motion.p>
        </div>
        <div className={`p-4 rounded-xl ${colorClasses[color].icon} group-hover:scale-110 transition-transform`}>
          {icon}
        </div>
      </div>

      {/* Decorative gradient overlay */}
      <div className="absolute -right-4 -bottom-4 w-24 h-24 bg-gradient-to-br from-white/5 to-transparent dark:from-white/10 rounded-full blur-2xl group-hover:scale-150 transition-transform" />
    </motion.div>
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
