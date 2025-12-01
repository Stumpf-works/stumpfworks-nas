import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { storageApi, Disk } from '@/api/storage';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { HardDrive, Zap, Database, Usb, RefreshCw, Activity, Edit2, X, Thermometer, Layers, AlertCircle, CheckCircle, XCircle } from 'lucide-react';

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
    const iconClass = "w-8 h-8";
    switch (type) {
      case 'nvme': return <Zap className={`${iconClass} text-yellow-500`} />;
      case 'ssd': return <Database className={`${iconClass} text-purple-500`} />;
      case 'usb': return <Usb className={`${iconClass} text-blue-500`} />;
      default: return <HardDrive className={`${iconClass} text-gray-500`} />;
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'healthy':
        return 'bg-gradient-to-r from-green-500 to-emerald-500 text-white';
      case 'warning':
        return 'bg-gradient-to-r from-yellow-500 to-orange-500 text-white';
      case 'critical':
        return 'bg-gradient-to-r from-orange-500 to-red-500 text-white';
      case 'failed':
        return 'bg-gradient-to-r from-red-500 to-rose-600 text-white';
      default:
        return 'bg-gradient-to-r from-gray-400 to-gray-500 text-white';
    }
  };

  const getTemperatureColor = (temp: number) => {
    if (temp >= 60) return 'text-red-500';
    if (temp >= 50) return 'text-orange-500';
    if (temp >= 40) return 'text-yellow-500';
    return 'text-green-500';
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl shadow-lg">
            <HardDrive className="w-6 h-6 text-white" />
          </div>
          <div>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Physical Disks
            </h2>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {disks.length} disk{disks.length !== 1 ? 's' : ''} detected
            </p>
          </div>
        </div>
        <Button
          onClick={loadDisks}
          variant="secondary"
          className="flex items-center gap-2"
        >
          <RefreshCw className="w-4 h-4" />
          Refresh
        </Button>
      </div>

      {/* Disk Grid */}
      {disks.length > 0 ? (
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
              getStatusBadge={getStatusBadge}
              getTemperatureColor={getTemperatureColor}
            />
          ))}
        </div>
      ) : (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex flex-col items-center justify-center py-16 px-4"
        >
          <div className="p-6 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-2xl mb-6">
            <HardDrive className="w-16 h-16 text-gray-400 dark:text-gray-600" />
          </div>
          <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
            No Disks Found
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-6 text-center max-w-md">
            No physical disks were detected in this environment. Check your hardware connections.
          </p>
          <Button
            onClick={loadDisks}
            variant="secondary"
            className="flex items-center gap-2"
          >
            <RefreshCw className="w-4 h-4" />
            Retry
          </Button>
        </motion.div>
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
  getDiskIcon: (type: string) => JSX.Element;
  getStatusBadge: (status: string) => string;
  getTemperatureColor: (temp: number) => string;
}

function DiskCard({ disk, onSelect, onRename, formatBytes, getDiskIcon, getStatusBadge, getTemperatureColor }: DiskCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ y: -2 }}
      transition={{ duration: 0.2 }}
    >
      <Card className="h-full hover:shadow-xl transition-shadow duration-200">
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center space-x-3 flex-1 min-w-0">
            <div className="flex-shrink-0 p-2 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl">
              {getDiskIcon(disk.type)}
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <h3 className="font-semibold text-gray-900 dark:text-gray-100 truncate">
                  {disk.label || disk.model || disk.name}
                </h3>
                <button
                  onClick={onRename}
                  className="text-gray-400 hover:text-macos-blue dark:hover:text-macos-purple transition-colors flex-shrink-0"
                  title="Rename disk"
                >
                  <Edit2 className="w-4 h-4" />
                </button>
                {disk.isSystem && (
                  <span className="px-2 py-0.5 text-xs bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-full flex-shrink-0 shadow-sm">
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
          <span className={`px-3 py-1 text-xs font-semibold rounded-full shadow-sm flex-shrink-0 ml-2 ${getStatusBadge(disk.status)}`}>
            {disk.status.toUpperCase()}
          </span>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-2 gap-3">
          <div className="p-3 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <Database className="w-4 h-4 text-gray-500 dark:text-gray-400" />
              <span className="text-xs text-gray-600 dark:text-gray-400">Size</span>
            </div>
            <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">
              {formatBytes(disk.size)}
            </p>
          </div>
          <div className="p-3 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <HardDrive className="w-4 h-4 text-gray-500 dark:text-gray-400" />
              <span className="text-xs text-gray-600 dark:text-gray-400">Type</span>
            </div>
            <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 uppercase">
              {disk.type}
            </p>
          </div>
          <div className="p-3 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <Thermometer className="w-4 h-4 text-gray-500 dark:text-gray-400" />
              <span className="text-xs text-gray-600 dark:text-gray-400">Temp</span>
            </div>
            <p className={`text-sm font-semibold ${getTemperatureColor(disk.temperature)}`}>
              {disk.temperature}°C
            </p>
          </div>
          <div className="p-3 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <Layers className="w-4 h-4 text-gray-500 dark:text-gray-400" />
              <span className="text-xs text-gray-600 dark:text-gray-400">Partitions</span>
            </div>
            <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">
              {disk.partitions?.length || 0}
            </p>
          </div>
        </div>

        {/* Partitions */}
        {disk.partitions && disk.partitions.length > 0 && (
          <div className="mt-4 space-y-2">
            <div className="flex items-center gap-2 mb-2">
              <Layers className="w-4 h-4 text-gray-500 dark:text-gray-400" />
              <span className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">
                Partitions
              </span>
            </div>
            {disk.partitions.map((part) => (
              <div
                key={part.name}
                className="flex justify-between items-center p-3 bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-lg border border-gray-200 dark:border-gray-700"
              >
                <div>
                  <span className="font-mono font-medium text-gray-900 dark:text-gray-100 text-sm">
                    {part.name}
                  </span>
                  {part.filesystem && (
                    <span className="ml-2 text-xs text-gray-600 dark:text-gray-400">
                      ({part.filesystem})
                    </span>
                  )}
                </div>
                <div className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  {formatBytes(part.size)}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Actions */}
        {disk.smartEnabled && (
          <div className="mt-4">
            <Button
              onClick={onSelect}
              variant="secondary"
              size="sm"
              className="w-full flex items-center justify-center gap-2"
            >
              <Activity className="w-4 h-4" />
              View SMART Data
            </Button>
          </div>
        )}
      </Card>
    </motion.div>
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
      className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-2xl shadow-2xl p-6 w-full max-w-2xl max-h-[80vh] overflow-auto"
      >
        {/* Header */}
        <div className="flex items-center justify-between mb-6 pb-4 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl">
              <Activity className="w-6 h-6 text-white" />
            </div>
            <div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                SMART Data
              </h2>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {disk.name}
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="space-y-4">
          {/* Health Status */}
          <div className={`p-4 rounded-xl ${smart.healthy ? 'bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20' : 'bg-gradient-to-br from-red-50 to-rose-50 dark:from-red-900/20 dark:to-rose-900/20'}`}>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                {smart.healthy ? (
                  <CheckCircle className="w-5 h-5 text-green-600 dark:text-green-400" />
                ) : (
                  <XCircle className="w-5 h-5 text-red-600 dark:text-red-400" />
                )}
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                  Health Status
                </span>
              </div>
              <span className={`text-lg font-bold ${smart.healthy ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
                {smart.healthy ? 'PASSED' : 'FAILED'}
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
    <div className={`p-4 rounded-xl transition-all ${warning ? 'bg-gradient-to-br from-yellow-50 to-orange-50 dark:from-yellow-900/20 dark:to-orange-900/20 border border-yellow-200 dark:border-yellow-800' : 'bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200 dark:to-macos-dark-300'}`}>
      <div className="flex items-center gap-2 mb-2">
        {warning && <AlertCircle className="w-4 h-4 text-yellow-600 dark:text-yellow-400" />}
        <div className="text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">
          {label}
        </div>
      </div>
      <div className={`text-lg font-bold ${warning ? 'text-yellow-600 dark:text-yellow-400' : 'text-gray-900 dark:text-gray-100'}`}>
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
      className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-2xl shadow-2xl p-6 w-full max-w-md"
      >
        {/* Header */}
        <div className="flex items-center justify-between mb-6 pb-4 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl">
              <Edit2 className="w-5 h-5 text-white" />
            </div>
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
              Rename Disk
            </h2>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Disk Info */}
          <div className="p-4 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl">
            <div className="space-y-2 text-sm">
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400 font-medium">Disk:</span>
                <span className="font-mono text-gray-900 dark:text-gray-100">{disk.name}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400 font-medium">Model:</span>
                <span className="text-gray-900 dark:text-gray-100">{disk.model}</span>
              </div>
              {disk.serial && (
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400 font-medium">Serial:</span>
                  <span className="font-mono text-xs text-gray-600 dark:text-gray-400">{disk.serial}</span>
                </div>
              )}
            </div>
            {!disk.serial && (
              <div className="mt-3 p-2 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg flex items-start gap-2">
                <AlertCircle className="w-4 h-4 text-yellow-600 dark:text-yellow-400 flex-shrink-0 mt-0.5" />
                <p className="text-xs text-yellow-700 dark:text-yellow-300">
                  This disk has no serial number and cannot be labeled
                </p>
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
            <div className="p-3 bg-gradient-to-br from-red-50 to-rose-50 dark:from-red-900/20 dark:to-rose-900/20 border border-red-200 dark:border-red-800 rounded-xl flex items-start gap-2">
              <XCircle className="w-4 h-4 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-red-600 dark:text-red-400">
                {error}
              </p>
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
