import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { storageApi, Disk, SMARTData } from '@/api/storage';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

export default function DiskManager() {
  const [disks, setDisks] = useState<Disk[]>([]);
  const [selectedDisk, setSelectedDisk] = useState<Disk | null>(null);
  const [showSMART, setShowSMART] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDisks();
  }, []);

  const loadDisks = async () => {
    try {
      const response = await storageApi.listDisks();
      if (response.success) {
        setDisks(response.data);
      }
    } catch (error) {
      console.error('Failed to load disks:', error);
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

  const getDiskIcon = (type: string) => {
    switch (type) {
      case 'nvme': return '⚡';
      case 'ssd': return '💎';
      case 'usb': return '🔌';
      default: return '💿';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy': return 'text-green-600 dark:text-green-400';
      case 'warning': return 'text-yellow-600 dark:text-yellow-400';
      case 'critical': return 'text-red-600 dark:text-red-400';
      case 'failed': return 'text-red-600 dark:text-red-400';
      default: return 'text-gray-600 dark:text-gray-400';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Disks ({disks.length})
        </h2>
        <Button onClick={loadDisks} variant="secondary">
          🔄 Refresh
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {disks.map((disk) => (
          <DiskCard
            key={disk.name}
            disk={disk}
            onSelect={() => {
              setSelectedDisk(disk);
              setShowSMART(true);
            }}
            formatBytes={formatBytes}
            getDiskIcon={getDiskIcon}
            getStatusColor={getStatusColor}
          />
        ))}
      </div>

      {/* SMART Data Modal */}
      <AnimatePresence>
        {showSMART && selectedDisk && (
          <SMARTModal
            disk={selectedDisk}
            onClose={() => {
              setShowSMART(false);
              setSelectedDisk(null);
            }}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

interface DiskCardProps {
  disk: Disk;
  onSelect: () => void;
  formatBytes: (bytes: number) => string;
  getDiskIcon: (type: string) => string;
  getStatusColor: (status: string) => string;
}

function DiskCard({ disk, onSelect, formatBytes, getDiskIcon, getStatusColor }: DiskCardProps) {
  return (
    <Card>
      <div className="flex items-start justify-between">
        <div className="flex items-center space-x-3">
          <div className="text-3xl">{getDiskIcon(disk.type)}</div>
          <div>
            <h3 className="font-semibold text-gray-900 dark:text-gray-100">
              {disk.name}
              {disk.isSystem && (
                <span className="ml-2 px-2 py-0.5 text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 rounded">
                  System
                </span>
              )}
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">{disk.model}</p>
          </div>
        </div>
        <div className={`text-sm font-medium ${getStatusColor(disk.status)}`}>
          {disk.status.toUpperCase()}
        </div>
      </div>

      <div className="mt-4 grid grid-cols-2 gap-3 text-sm">
        <div>
          <span className="text-gray-600 dark:text-gray-400">Size:</span>
          <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
            {formatBytes(disk.size)}
          </span>
        </div>
        <div>
          <span className="text-gray-600 dark:text-gray-400">Type:</span>
          <span className="ml-2 font-medium text-gray-900 dark:text-gray-100 uppercase">
            {disk.type}
          </span>
        </div>
        <div>
          <span className="text-gray-600 dark:text-gray-400">Temp:</span>
          <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
            {disk.temperature}°C
          </span>
        </div>
        <div>
          <span className="text-gray-600 dark:text-gray-400">Partitions:</span>
          <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
            {disk.partitions.length}
          </span>
        </div>
      </div>

      {disk.partitions.length > 0 && (
        <div className="mt-4 space-y-2">
          <div className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">
            Partitions
          </div>
          {disk.partitions.map((part) => (
            <div
              key={part.name}
              className="flex justify-between items-center p-2 bg-gray-50 dark:bg-macos-dark-200 rounded text-xs"
            >
              <div>
                <span className="font-mono font-medium text-gray-900 dark:text-gray-100">
                  {part.name}
                </span>
                {part.filesystem && (
                  <span className="ml-2 text-gray-600 dark:text-gray-400">
                    ({part.filesystem})
                  </span>
                )}
              </div>
              <div className="text-gray-600 dark:text-gray-400">
                {formatBytes(part.size)}
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="mt-4 flex space-x-2">
        {disk.smartEnabled && (
          <Button onClick={onSelect} variant="secondary" size="sm" className="flex-1">
            📊 View SMART
          </Button>
        )}
      </div>
    </Card>
  );
}

interface SMARTModalProps {
  disk: Disk;
  onClose: () => void;
}

function SMARTModal({ disk, onClose }: SMARTModalProps) {
  if (!disk.smart) {
    return null;
  }

  const smart = disk.smart;

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-2xl max-h-[80vh] overflow-auto"
      >
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
            SMART Data - {disk.name}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="space-y-4">
          {/* Health Status */}
          <div className="p-4 rounded-lg bg-gray-50 dark:bg-macos-dark-200">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                Health Status
              </span>
              <span className={`text-lg font-bold ${smart.healthy ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
                {smart.healthy ? '✅ PASSED' : '❌ FAILED'}
              </span>
            </div>
          </div>

          {/* SMART Attributes */}
          <div className="grid grid-cols-2 gap-4">
            <SMARTAttribute
              label="Temperature"
              value={`${smart.temperature}°C`}
              warning={smart.temperature > 60}
            />
            <SMARTAttribute
              label="Power On Hours"
              value={smart.powerOnHours.toLocaleString()}
            />
            <SMARTAttribute
              label="Power Cycles"
              value={smart.powerCycleCount.toLocaleString()}
            />
            <SMARTAttribute
              label="Reallocated Sectors"
              value={smart.reallocatedSectors.toString()}
              warning={smart.reallocatedSectors > 0}
            />
            <SMARTAttribute
              label="Pending Sectors"
              value={smart.pendingSectors.toString()}
              warning={smart.pendingSectors > 0}
            />
            <SMARTAttribute
              label="Uncorrectable Errors"
              value={smart.uncorrectableErrors.toString()}
              warning={smart.uncorrectableErrors > 0}
            />
            {smart.percentLifeUsed > 0 && (
              <SMARTAttribute
                label="Life Used (SSD)"
                value={`${smart.percentLifeUsed}%`}
                warning={smart.percentLifeUsed > 80}
              />
            )}
          </div>

          <div className="text-xs text-gray-500 dark:text-gray-400 text-center pt-4">
            Last Updated: {new Date(smart.lastUpdated).toLocaleString()}
          </div>
        </div>
      </motion.div>
    </motion.div>
  );
}

interface SMARTAttributeProps {
  label: string;
  value: string;
  warning?: boolean;
}

function SMARTAttribute({ label, value, warning }: SMARTAttributeProps) {
  return (
    <div className="p-3 bg-gray-50 dark:bg-macos-dark-200 rounded">
      <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">{label}</div>
      <div className={`text-lg font-semibold ${warning ? 'text-yellow-600 dark:text-yellow-400' : 'text-gray-900 dark:text-gray-100'}`}>
        {value}
      </div>
    </div>
  );
}
