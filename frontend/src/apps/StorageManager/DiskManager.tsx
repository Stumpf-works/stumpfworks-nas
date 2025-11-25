import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { storageApi, Disk } from '@/api/storage';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

export default function DiskManager() {
  const [disks, setDisks] = useState<Disk[]>([]);
  const [selectedDisk, setSelectedDisk] = useState<Disk | null>(null);
  const [showSMART, setShowSMART] = useState(false);
  const [showRename, setShowRename] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDisks();
  }, []);

  const loadDisks = async () => {
    try {
      const response = await storageApi.listDisks();
      if (response.success && response.data) {
        setDisks(response.data);
      } else {
        console.error('Failed to load disks:', response.error);
      }
    } catch (error) {
      console.error('Failed to load disks:', error);
      alert('Failed to load disks. Please check the console for details.');
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
            onRename={() => {
              setSelectedDisk(disk);
              setShowRename(true);
            }}
            formatBytes={formatBytes}
            getDiskIcon={getDiskIcon}
            getStatusColor={getStatusColor}
          />
        ))}
      </div>

      {disks.length === 0 && (
        <div className="text-center py-12 text-gray-600 dark:text-gray-400">
          <div className="text-6xl mb-4">💿</div>
          <p className="text-lg font-medium mb-2">No disks found</p>
          <p className="text-sm mb-4">No physical disks were detected in this environment</p>
          <Button onClick={loadDisks} variant="secondary">
            🔄 Retry
          </Button>
        </div>
      )}

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

      {/* Rename Disk Modal */}
      <AnimatePresence>
        {showRename && selectedDisk && (
          <RenameDiskModal
            disk={selectedDisk}
            onClose={() => {
              setShowRename(false);
              setSelectedDisk(null);
            }}
            onSuccess={() => {
              loadDisks();
              setShowRename(false);
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
  onRename: () => void;
  formatBytes: (bytes: number) => string;
  getDiskIcon: (type: string) => string;
  getStatusColor: (status: string) => string;
}

function DiskCard({ disk, onSelect, onRename, formatBytes, getDiskIcon, getStatusColor }: DiskCardProps) {
  return (
    <Card>
      <div className="flex items-start justify-between">
        <div className="flex items-center space-x-3 flex-1">
          <div className="text-3xl">{getDiskIcon(disk.type)}</div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <h3 className="font-semibold text-gray-900 dark:text-gray-100 truncate">
                {disk.label || disk.model || disk.name}
              </h3>
              <button
                onClick={onRename}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors flex-shrink-0"
                title="Rename disk"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                </svg>
              </button>
              {disk.isSystem && (
                <span className="px-2 py-0.5 text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 rounded flex-shrink-0">
                  System
                </span>
              )}
            </div>
            {disk.label ? (
              <p className="text-xs text-gray-500 dark:text-gray-500 truncate">
                {disk.model} • {disk.name}
              </p>
            ) : (
              <p className="text-xs text-gray-500 dark:text-gray-500 truncate">
                {disk.name}
              </p>
            )}
          </div>
        </div>
        <div className={`text-sm font-medium ${getStatusColor(disk.status)} flex-shrink-0 ml-2`}>
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
            {disk.partitions?.length || 0}
          </span>
        </div>
      </div>

      {disk.partitions && disk.partitions.length > 0 && (
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

interface RenameDiskModalProps {
  disk: Disk;
  onClose: () => void;
  onSuccess: () => void;
}

function RenameDiskModal({ disk, onClose, onSuccess }: RenameDiskModalProps) {
  const [label, setLabel] = useState(disk.label || '');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await storageApi.setDiskLabel(disk.name, label.trim());
      if (response.success) {
        onSuccess();
      } else {
        setError(response.error?.message || 'Failed to update disk label');
      }
    } catch (err) {
      setError('Failed to update disk label');
      console.error('Failed to update disk label:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleClear = async () => {
    setLoading(true);
    setError('');

    try {
      const response = await storageApi.setDiskLabel(disk.name, '');
      if (response.success) {
        onSuccess();
      } else {
        setError(response.error?.message || 'Failed to clear disk label');
      }
    } catch (err) {
      setError('Failed to clear disk label');
      console.error('Failed to clear disk label:', err);
    } finally {
      setLoading(false);
    }
  };

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
        className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
      >
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
            Rename Disk
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

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="p-3 bg-gray-50 dark:bg-macos-dark-200 rounded text-sm">
            <div className="flex items-center justify-between mb-1">
              <span className="text-gray-600 dark:text-gray-400">Disk:</span>
              <span className="font-mono text-gray-900 dark:text-gray-100">{disk.name}</span>
            </div>
            <div className="flex items-center justify-between mb-1">
              <span className="text-gray-600 dark:text-gray-400">Model:</span>
              <span className="text-gray-900 dark:text-gray-100">{disk.model}</span>
            </div>
            {disk.serial && (
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">Serial:</span>
                <span className="font-mono text-xs text-gray-600 dark:text-gray-400">{disk.serial}</span>
              </div>
            )}
            {!disk.serial && (
              <div className="mt-2 text-xs text-yellow-600 dark:text-yellow-400">
                ⚠️ This disk has no serial number and cannot be labeled
              </div>
            )}
          </div>

          <div>
            <Input
              label="Custom Label"
              value={label}
              onChange={(e) => setLabel(e.target.value)}
              placeholder="e.g., Main Storage, Backup Drive"
              disabled={loading}
            />
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Give this disk a friendly name. Leave empty to show model name.
            </p>
          </div>

          {error && (
            <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded text-sm">
              {error}
            </div>
          )}

          <div className="flex gap-2 pt-2">
            {disk.label && disk.serial && (
              <Button
                type="button"
                onClick={handleClear}
                variant="secondary"
                isLoading={loading}
                className="flex-1"
              >
                Clear Label
              </Button>
            )}
            <Button
              type="submit"
              isLoading={loading}
              disabled={!disk.serial}
              className="flex-1"
            >
              Save Label
            </Button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
}
