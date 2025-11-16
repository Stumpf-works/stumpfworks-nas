// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { systemApi } from '@/api/system';
import { ChevronRight, Cpu, HardDrive, Network, MemoryStick } from 'lucide-react';
import type { SystemMetrics } from '@/api/system';

interface WidgetSidebarProps {
  isOpen: boolean;
  onToggle: () => void;
}

function CPUWidget({ metrics }: { metrics: SystemMetrics }) {
  const usage = metrics.cpu?.usagePercent || 0;
  const perCore = metrics.cpu?.perCore || [];

  return (
    <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4 backdrop-blur-sm">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <Cpu className="w-4 h-4 text-macos-blue" />
          <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
            CPU
          </span>
        </div>
        <span className="text-lg font-bold text-macos-blue">
          {usage.toFixed(1)}%
        </span>
      </div>

      {/* Progress Bar */}
      <div className="h-2 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden">
        <motion.div
          initial={{ width: 0 }}
          animate={{ width: `${Math.min(usage, 100)}%` }}
          transition={{ type: 'spring', stiffness: 100, damping: 20 }}
          className={`h-full rounded-full ${
            usage > 80
              ? 'bg-red-500'
              : usage > 60
              ? 'bg-yellow-500'
              : 'bg-macos-blue'
          }`}
        />
      </div>

      {perCore.length > 0 && (
        <div className="mt-2 text-xs text-gray-600 dark:text-gray-400">
          {perCore.length} Cores
        </div>
      )}
    </div>
  );
}

function MemoryWidget({ metrics }: { metrics: SystemMetrics }) {
  const usage = metrics.memory?.usedPercent || 0;
  const total = metrics.memory?.total || 0;
  const used = metrics.memory?.used || 0;

  const formatBytes = (bytes: number) => {
    const gb = bytes / (1024 ** 3);
    return gb.toFixed(1) + ' GB';
  };

  return (
    <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4 backdrop-blur-sm">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <MemoryStick className="w-4 h-4 text-macos-purple" />
          <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
            Memory
          </span>
        </div>
        <span className="text-lg font-bold text-macos-purple">
          {usage.toFixed(1)}%
        </span>
      </div>

      {/* Progress Bar */}
      <div className="h-2 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden">
        <motion.div
          initial={{ width: 0 }}
          animate={{ width: `${Math.min(usage, 100)}%` }}
          transition={{ type: 'spring', stiffness: 100, damping: 20 }}
          className={`h-full rounded-full ${
            usage > 80
              ? 'bg-red-500'
              : usage > 60
              ? 'bg-yellow-500'
              : 'bg-macos-purple'
          }`}
        />
      </div>

      <div className="mt-2 text-xs text-gray-600 dark:text-gray-400">
        {formatBytes(used)} / {formatBytes(total)}
      </div>
    </div>
  );
}

function StorageWidget({ metrics }: { metrics: SystemMetrics }) {
  const disk = metrics.disk && metrics.disk.length > 0 ? metrics.disk[0] : null;
  const usage = disk?.usedPercent || 0;
  const total = disk?.total || 0;
  const free = disk?.free || 0;

  const formatBytes = (bytes: number) => {
    const gb = bytes / (1024 ** 3);
    if (gb >= 1000) {
      return (gb / 1024).toFixed(1) + ' TB';
    }
    return gb.toFixed(1) + ' GB';
  };

  return (
    <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4 backdrop-blur-sm">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <HardDrive className="w-4 h-4 text-macos-green" />
          <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
            Storage
          </span>
        </div>
        <span className="text-lg font-bold text-macos-green">
          {usage.toFixed(1)}%
        </span>
      </div>

      {/* Progress Bar */}
      <div className="h-2 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden">
        <motion.div
          initial={{ width: 0 }}
          animate={{ width: `${Math.min(usage, 100)}%` }}
          transition={{ type: 'spring', stiffness: 100, damping: 20 }}
          className={`h-full rounded-full ${
            usage > 90
              ? 'bg-red-500'
              : usage > 70
              ? 'bg-yellow-500'
              : 'bg-macos-green'
          }`}
        />
      </div>

      <div className="mt-2 text-xs text-gray-600 dark:text-gray-400">
        {formatBytes(free)} available of {formatBytes(total)}
      </div>
    </div>
  );
}

function NetworkWidget({ metrics }: { metrics: SystemMetrics }) {
  const bytesSent = metrics.network?.bytesSent || 0;
  const bytesRecv = metrics.network?.bytesRecv || 0;

  const formatBytes = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 ** 2) return (bytes / 1024).toFixed(1) + ' KB';
    if (bytes < 1024 ** 3) return (bytes / (1024 ** 2)).toFixed(1) + ' MB';
    return (bytes / (1024 ** 3)).toFixed(1) + ' GB';
  };

  return (
    <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4 backdrop-blur-sm">
      <div className="flex items-center gap-2 mb-3">
        <Network className="w-4 h-4 text-macos-orange" />
        <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
          Network
        </span>
      </div>

      <div className="space-y-2">
        <div className="flex justify-between items-center">
          <span className="text-xs text-gray-600 dark:text-gray-400">Sent</span>
          <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
            {formatBytes(bytesSent)}
          </span>
        </div>
        <div className="flex justify-between items-center">
          <span className="text-xs text-gray-600 dark:text-gray-400">Received</span>
          <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
            {formatBytes(bytesRecv)}
          </span>
        </div>
      </div>
    </div>
  );
}

export default function WidgetSidebar({ isOpen, onToggle }: WidgetSidebarProps) {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const response = await systemApi.getMetrics();
        if (response.success && response.data) {
          setMetrics(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch metrics:', error);
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <>
      {/* Toggle Button */}
      <motion.button
        initial={{ x: 300 }}
        animate={{ x: isOpen ? 300 : 0 }}
        onClick={onToggle}
        className="fixed right-0 top-1/2 transform -translate-y-1/2 z-40 bg-white/80 dark:bg-macos-dark-200/80 backdrop-blur-sm rounded-l-lg shadow-lg px-2 py-4 hover:bg-white dark:hover:bg-macos-dark-200 transition-colors"
      >
        <motion.div
          animate={{ rotate: isOpen ? 0 : 180 }}
          transition={{ duration: 0.2 }}
        >
          <ChevronRight className="w-5 h-5 text-gray-600 dark:text-gray-400" />
        </motion.div>
      </motion.button>

      {/* Sidebar */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ x: 300, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            exit={{ x: 300, opacity: 0 }}
            transition={{ type: 'spring', stiffness: 300, damping: 30 }}
            className="fixed right-0 top-8 bottom-20 w-72 z-30"
          >
            <div className="h-full p-4 overflow-y-auto no-scrollbar">
              <div className="space-y-3">
                <div className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-4 px-1">
                  System Widgets
                </div>

                {metrics ? (
                  <>
                    <CPUWidget metrics={metrics} />
                    <MemoryWidget metrics={metrics} />
                    <StorageWidget metrics={metrics} />
                    <NetworkWidget metrics={metrics} />
                  </>
                ) : (
                  <div className="flex items-center justify-center py-12">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
                  </div>
                )}
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </>
  );
}
