// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { systemApi } from '@/api/system';
import { syslibApi } from '@/api/syslib';
import { X } from 'lucide-react';

interface AboutDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

interface SystemInfo {
  hostname: string;
  version: string;
  cpu: {
    model: string;
    cores: number;
  };
  memory: {
    total: number;
    used: number;
    free: number;
    totalFormatted: string;
    usedFormatted: string;
    freeFormatted: string;
  };
  storage: {
    total: number;
    totalFormatted: string;
  };
  uptime: string;
}

type TabId = 'overview' | 'storage' | 'memory' | 'support';

interface Tab {
  id: TabId;
  label: string;
}

const tabs: Tab[] = [
  { id: 'overview', label: 'Overview' },
  { id: 'storage', label: 'Storage' },
  { id: 'memory', label: 'Memory' },
  { id: 'support', label: 'Support' },
];

export default function AboutDialog({ isOpen, onClose }: AboutDialogProps) {
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<TabId>('overview');
  const [pools, setPools] = useState<any[]>([]);
  const [rawMetrics, setRawMetrics] = useState<any>(null);

  useEffect(() => {
    if (isOpen) {
      fetchSystemInfo();
      fetchZFSPools();
    }
  }, [isOpen]);

  const fetchSystemInfo = async () => {
    setIsLoading(true);
    try {
      const [infoResponse, metricsResponse] = await Promise.all([
        systemApi.getInfo(),
        systemApi.getMetrics(),
      ]);

      if (infoResponse.success && infoResponse.data && metricsResponse.success && metricsResponse.data) {
        const info = infoResponse.data;
        const metrics = metricsResponse.data;
        setRawMetrics(metrics);

        // Calculate uptime
        const uptimeSeconds = info.uptime || 0;
        const days = Math.floor(uptimeSeconds / 86400);
        const hours = Math.floor((uptimeSeconds % 86400) / 3600);
        const minutes = Math.floor((uptimeSeconds % 3600) / 60);

        let uptimeString = '';
        if (days > 0) uptimeString += `${days}d `;
        if (hours > 0) uptimeString += `${hours}h `;
        uptimeString += `${minutes}m`;

        // Get total storage from first disk
        const totalStorage = metrics.disk && metrics.disk.length > 0 ? metrics.disk[0].total : 0;

        // Memory info
        const memTotal = metrics.memory?.total || 0;
        const memUsed = metrics.memory?.used || 0;
        const memFree = memTotal - memUsed;

        setSystemInfo({
          hostname: info.hostname || 'StumpfWorks NAS',
          version: '1.1.1',
          cpu: {
            model: `${info.cpuCores}-Core CPU`,
            cores: info.cpuCores || 0,
          },
          memory: {
            total: memTotal,
            used: memUsed,
            free: memFree,
            totalFormatted: formatBytes(memTotal),
            usedFormatted: formatBytes(memUsed),
            freeFormatted: formatBytes(memFree),
          },
          storage: {
            total: totalStorage,
            totalFormatted: formatBytes(totalStorage),
          },
          uptime: uptimeString,
        });
      }
    } catch (error) {
      console.error('Failed to fetch system info:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchZFSPools = async () => {
    try {
      const response = await syslibApi.zfs.listPools();
      if (response.success && response.data) {
        setPools(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch ZFS pools:', error);
    }
  };

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 GB';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const getHealthColor = (health: string) => {
    switch (health) {
      case 'ONLINE':
        return 'text-green-600 dark:text-green-400';
      case 'DEGRADED':
        return 'text-yellow-600 dark:text-yellow-400';
      case 'OFFLINE':
      case 'FAULTED':
        return 'text-red-600 dark:text-red-400';
      default:
        return 'text-gray-600 dark:text-gray-400';
    }
  };

  const renderTabContent = () => {
    if (isLoading) {
      return (
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
        </div>
      );
    }

    if (!systemInfo) {
      return (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          Failed to load system information
        </div>
      );
    }

    switch (activeTab) {
      case 'overview':
        return (
          <div className="space-y-4">
            {/* System Name */}
            <div>
              <div className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
                {systemInfo.hostname}
              </div>
            </div>

            {/* Specs */}
            <div className="space-y-2.5 pt-2">
              <div className="flex justify-between items-center py-2 border-b border-gray-200 dark:border-gray-700">
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                  Chip
                </span>
                <span className="text-sm text-gray-900 dark:text-gray-100">
                  {systemInfo.cpu.model}
                </span>
              </div>

              <div className="flex justify-between items-center py-2 border-b border-gray-200 dark:border-gray-700">
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                  CPU Cores
                </span>
                <span className="text-sm text-gray-900 dark:text-gray-100">
                  {systemInfo.cpu.cores} Cores
                </span>
              </div>

              <div className="flex justify-between items-center py-2 border-b border-gray-200 dark:border-gray-700">
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                  Memory
                </span>
                <span className="text-sm text-gray-900 dark:text-gray-100">
                  {systemInfo.memory.totalFormatted}
                </span>
              </div>

              <div className="flex justify-between items-center py-2 border-b border-gray-200 dark:border-gray-700">
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                  Storage
                </span>
                <span className="text-sm text-gray-900 dark:text-gray-100">
                  {systemInfo.storage.totalFormatted}
                </span>
              </div>

              <div className="flex justify-between items-center py-2">
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                  Uptime
                </span>
                <span className="text-sm text-gray-900 dark:text-gray-100">
                  {systemInfo.uptime}
                </span>
              </div>
            </div>
          </div>
        );

      case 'storage':
        return (
          <div className="space-y-4">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              ZFS Storage Pools
            </h3>

            {pools.length === 0 ? (
              <p className="text-sm text-gray-600 dark:text-gray-400">
                No ZFS pools configured
              </p>
            ) : (
              <div className="space-y-3">
                {pools.map((pool) => (
                  <div
                    key={pool.name}
                    className="p-4 bg-gray-50 dark:bg-macos-dark-200 rounded-lg border border-gray-200 dark:border-gray-700"
                  >
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <h4 className="font-semibold text-gray-900 dark:text-gray-100">
                          {pool.name}
                        </h4>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                          {pool.mountpoint || '-'}
                        </p>
                      </div>
                      <span className={`text-sm font-medium ${getHealthColor(pool.health)}`}>
                        {pool.health}
                      </span>
                    </div>
                    <div className="space-y-1.5 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">Capacity:</span>
                        <span className="text-gray-900 dark:text-gray-100 font-mono">
                          {formatBytes(pool.used)} / {formatBytes(pool.size)}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">Free:</span>
                        <span className="text-gray-900 dark:text-gray-100 font-mono">
                          {formatBytes(pool.free)}
                        </span>
                      </div>
                      {/* Usage bar */}
                      <div className="mt-2">
                        <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-macos-blue rounded-full transition-all"
                            style={{ width: `${((pool.used / pool.size) * 100).toFixed(1)}%` }}
                          />
                        </div>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1 text-right">
                          {((pool.used / pool.size) * 100).toFixed(1)}% used
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        );

      case 'memory':
        return (
          <div className="space-y-4">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Memory Details
            </h3>

            <div className="space-y-3">
              {/* RAM */}
              <div className="p-4 bg-gray-50 dark:bg-macos-dark-200 rounded-lg border border-gray-200 dark:border-gray-700">
                <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-3">
                  RAM
                </h4>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Total:</span>
                    <span className="text-gray-900 dark:text-gray-100 font-mono">
                      {systemInfo.memory.totalFormatted}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Used:</span>
                    <span className="text-gray-900 dark:text-gray-100 font-mono">
                      {systemInfo.memory.usedFormatted}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600 dark:text-gray-400">Free:</span>
                    <span className="text-gray-900 dark:text-gray-100 font-mono">
                      {systemInfo.memory.freeFormatted}
                    </span>
                  </div>
                  {/* Usage bar */}
                  <div className="mt-2">
                    <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-macos-blue rounded-full transition-all"
                        style={{
                          width: `${((systemInfo.memory.used / systemInfo.memory.total) * 100).toFixed(1)}%`,
                        }}
                      />
                    </div>
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1 text-right">
                      {((systemInfo.memory.used / systemInfo.memory.total) * 100).toFixed(1)}% used
                    </p>
                  </div>
                </div>
              </div>

              {/* Swap (if available) */}
              {rawMetrics?.memory?.swap && (
                <div className="p-4 bg-gray-50 dark:bg-macos-dark-200 rounded-lg border border-gray-200 dark:border-gray-700">
                  <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-3">
                    Swap
                  </h4>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Total:</span>
                      <span className="text-gray-900 dark:text-gray-100 font-mono">
                        {formatBytes(rawMetrics.memory.swap.total || 0)}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Used:</span>
                      <span className="text-gray-900 dark:text-gray-100 font-mono">
                        {formatBytes(rawMetrics.memory.swap.used || 0)}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Free:</span>
                      <span className="text-gray-900 dark:text-gray-100 font-mono">
                        {formatBytes(rawMetrics.memory.swap.free || 0)}
                      </span>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        );

      case 'support':
        return (
          <div className="space-y-4">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Support & Updates
            </h3>

            {/* Buttons */}
            <div className="flex gap-2 pt-2">
              <button className="flex-1 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors text-sm font-medium">
                Check for Updates
              </button>
              <button
                onClick={() => window.open('https://github.com/Stumpf-works/stumpfworks-nas', '_blank')}
                className="flex-1 px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors text-sm font-medium"
              >
                GitHub
              </button>
            </div>

            {/* Links */}
            <div className="pt-4 space-y-2">
              <a
                href="https://github.com/Stumpf-works/stumpfworks-nas/wiki"
                target="_blank"
                rel="noopener noreferrer"
                className="block text-sm text-macos-blue hover:underline"
              >
                📚 Documentation
              </a>
              <a
                href="https://github.com/Stumpf-works/stumpfworks-nas/issues"
                target="_blank"
                rel="noopener noreferrer"
                className="block text-sm text-macos-blue hover:underline"
              >
                🐛 Report an Issue
              </a>
              <a
                href="https://github.com/Stumpf-works/stumpfworks-nas/discussions"
                target="_blank"
                rel="noopener noreferrer"
                className="block text-sm text-macos-blue hover:underline"
              >
                💬 Community Discussions
              </a>
            </div>

            {/* Footer */}
            <div className="pt-6 text-center text-xs text-gray-500 dark:text-gray-400 border-t border-gray-200 dark:border-gray-700">
              <p>© 2025 Stumpf.Works. All rights reserved.</p>
              <p className="mt-1">
                Built with ❤️ for the homelab community
              </p>
              <p className="mt-2 font-mono text-gray-400 dark:text-gray-500">
                Version {systemInfo.version}
              </p>
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/30 dark:bg-black/50 backdrop-blur-sm z-[100]"
          />

          {/* Dialog */}
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            transition={{ type: 'spring', stiffness: 300, damping: 30 }}
            className="fixed left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 z-[101]"
          >
            <div className="w-[550px] bg-white dark:bg-macos-dark-100 rounded-2xl shadow-macos-xl border border-gray-200/50 dark:border-gray-700/50 overflow-hidden">
              {/* Header */}
              <div className="relative h-40 bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500 flex items-center justify-center">
                <button
                  onClick={onClose}
                  className="absolute top-4 right-4 p-1.5 rounded-full hover:bg-white/20 transition-colors"
                >
                  <X className="w-5 h-5 text-white" />
                </button>

                {/* Logo */}
                <div className="flex flex-col items-center">
                  <div className="w-20 h-20 rounded-2xl bg-white dark:bg-macos-dark-200 shadow-lg flex items-center justify-center">
                    <svg
                      className="w-12 h-12 text-macos-blue"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"
                      />
                    </svg>
                  </div>
                  <div className="mt-3 text-white text-center">
                    <h2 className="text-xl font-bold">Stumpf.Works NAS</h2>
                    <p className="text-sm opacity-90 mt-0.5">Version 1.1.1</p>
                  </div>
                </div>
              </div>

              {/* Tabs */}
              <div className="border-b border-gray-200 dark:border-gray-700">
                <div className="flex">
                  {tabs.map((tab) => (
                    <button
                      key={tab.id}
                      onClick={() => setActiveTab(tab.id)}
                      className={`flex-1 px-4 py-3 text-sm font-medium transition-colors ${
                        activeTab === tab.id
                          ? 'text-macos-blue border-b-2 border-macos-blue'
                          : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
                      }`}
                    >
                      {tab.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Content */}
              <div className="p-6 max-h-[400px] overflow-y-auto">
                {renderTabContent()}
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
