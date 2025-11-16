// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { systemApi } from '@/api/system';
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
    totalFormatted: string;
  };
  storage: {
    total: number;
    totalFormatted: string;
  };
  uptime: string;
}

export default function AboutDialog({ isOpen, onClose }: AboutDialogProps) {
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (isOpen) {
      fetchSystemInfo();
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

        setSystemInfo({
          hostname: info.hostname || 'StumpfWorks NAS',
          version: '1.1.0',
          cpu: {
            model: `${info.cpuCores}-Core CPU`,
            cores: info.cpuCores || 0,
          },
          memory: {
            total: metrics.memory?.total || 0,
            totalFormatted: formatBytes(metrics.memory?.total || 0),
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

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 GB';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
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
            <div className="w-[500px] bg-white dark:bg-macos-dark-100 rounded-2xl shadow-macos-xl border border-gray-200/50 dark:border-gray-700/50 overflow-hidden">
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
                    <p className="text-sm opacity-90 mt-0.5">Version 1.1.0</p>
                  </div>
                </div>
              </div>

              {/* Content */}
              <div className="p-6">
                {isLoading ? (
                  <div className="flex items-center justify-center py-8">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
                  </div>
                ) : systemInfo ? (
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

                    {/* Buttons */}
                    <div className="flex gap-2 pt-4">
                      <button className="flex-1 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors text-sm font-medium">
                        Software Update
                      </button>
                      <button className="flex-1 px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors text-sm font-medium">
                        Support
                      </button>
                    </div>

                    {/* Footer */}
                    <div className="pt-4 text-center text-xs text-gray-500 dark:text-gray-400">
                      <p>© 2025 Stumpf.Works. All rights reserved.</p>
                      <p className="mt-1">
                        Built with ❤️ for the homelab community
                      </p>
                    </div>
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                    Failed to load system information
                  </div>
                )}
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
