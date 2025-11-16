// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  HardDrive,
  Plus,
  RefreshCw,
  Shield,
  Activity,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Server,
  Zap,
  X,
} from 'lucide-react';
import { syslibApi, type RAIDArray, type CreateRAIDArrayRequest } from '@/api/syslib';
import { storageApi, type Disk } from '@/api/storage';

export default function RAIDManager() {
  const [arrays, setArrays] = useState<RAIDArray[]>([]);
  const [selectedArray, setSelectedArray] = useState<RAIDArray | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [availableDisks, setAvailableDisks] = useState<Disk[]>([]);

  const [formData, setFormData] = useState<CreateRAIDArrayRequest>({
    name: '',
    level: '5',
    devices: [],
    spare: [],
  });

  // Fetch RAID arrays
  const fetchArrays = async () => {
    setIsLoading(true);
    try {
      const response = await syslibApi.raid.listArrays();
      if (response.success && response.data) {
        setArrays(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch RAID arrays:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch available disks for array creation
  const fetchAvailableDisks = async () => {
    try {
      const response = await storageApi.listDisks();
      if (response.success && response.data) {
        // Filter out system disks and already used disks
        const available = response.data.filter((disk) => !disk.isSystem);
        setAvailableDisks(available);
      }
    } catch (error) {
      console.error('Failed to fetch disks:', error);
    }
  };

  useEffect(() => {
    fetchArrays();
  }, []);

  useEffect(() => {
    if (showCreateDialog) {
      fetchAvailableDisks();
    }
  }, [showCreateDialog]);

  const handleCreateArray = async () => {
    if (!formData.name || formData.devices.length < 2) {
      alert('Please provide array name and select at least 2 devices');
      return;
    }

    try {
      const response = await syslibApi.raid.createArray(formData);
      if (response.success) {
        alert(`RAID array created: ${formData.name}`);
        setShowCreateDialog(false);
        setFormData({ name: '', level: '5', devices: [], spare: [] });
        fetchArrays();
      }
    } catch (error) {
      console.error('Failed to create RAID array:', error);
      alert('Failed to create RAID array');
    }
  };

  const getStateIcon = (state: string) => {
    switch (state.toLowerCase()) {
      case 'clean':
      case 'active':
        return <CheckCircle className="w-5 h-5 text-green-500" />;
      case 'degraded':
      case 'recovering':
        return <AlertTriangle className="w-5 h-5 text-yellow-500" />;
      case 'failed':
      case 'inactive':
        return <XCircle className="w-5 h-5 text-red-500" />;
      default:
        return <Activity className="w-5 h-5 text-gray-500" />;
    }
  };

  const getRaidLevelInfo = (level: string) => {
    const info: Record<string, { name: string; description: string; minDisks: number }> = {
      '0': { name: 'RAID 0 (Striping)', description: 'Performance, no redundancy', minDisks: 2 },
      '1': { name: 'RAID 1 (Mirroring)', description: 'Complete redundancy', minDisks: 2 },
      '5': { name: 'RAID 5 (Parity)', description: 'Good balance', minDisks: 3 },
      '6': { name: 'RAID 6 (Dual Parity)', description: 'High redundancy', minDisks: 4 },
      '10': { name: 'RAID 10 (1+0)', description: 'Performance + redundancy', minDisks: 4 },
    };
    return info[level] || { name: `RAID ${level}`, description: 'Unknown', minDisks: 2 };
  };

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const toggleDevice = (diskPath: string, type: 'device' | 'spare') => {
    if (type === 'device') {
      setFormData((prev) => ({
        ...prev,
        devices: prev.devices.includes(diskPath)
          ? prev.devices.filter((d) => d !== diskPath)
          : [...prev.devices, diskPath],
      }));
    } else {
      setFormData((prev) => ({
        ...prev,
        spare: prev.spare?.includes(diskPath)
          ? prev.spare.filter((d) => d !== diskPath)
          : [...(prev.spare || []), diskPath],
      }));
    }
  };

  return (
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-100">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center gap-3">
          <Shield className="w-6 h-6 text-macos-blue" />
          <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
            RAID Array Manager
          </h1>
        </div>
        <div className="flex gap-2">
          <button
            onClick={fetchArrays}
            className="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </button>
          <button
            onClick={() => setShowCreateDialog(true)}
            className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            <Plus className="w-4 h-4" />
            Create Array
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Arrays List */}
        <div className="w-1/3 border-r border-gray-200 dark:border-gray-700 overflow-y-auto">
          {isLoading ? (
            <div className="flex items-center justify-center p-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
            </div>
          ) : arrays.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center">
              <Shield className="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" />
              <p className="text-gray-500 dark:text-gray-400">No RAID arrays found</p>
              <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
                Create a new array to get started
              </p>
            </div>
          ) : (
            <div className="p-4 space-y-2">
              {arrays.map((array) => (
                <motion.div
                  key={array.name}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  onClick={() => setSelectedArray(array)}
                  className={`p-4 rounded-xl cursor-pointer transition-all ${
                    selectedArray?.name === array.name
                      ? 'bg-macos-blue/10 dark:bg-macos-blue/20 border-2 border-macos-blue'
                      : 'bg-gray-50 dark:bg-macos-dark-200 hover:bg-gray-100 dark:hover:bg-macos-dark-300 border-2 border-transparent'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      {getStateIcon(array.state)}
                      <span className="font-semibold text-gray-900 dark:text-gray-100">
                        {array.name}
                      </span>
                    </div>
                    <span className="text-xs px-2 py-1 bg-macos-blue/10 text-macos-blue rounded">
                      RAID {array.level}
                    </span>
                  </div>

                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">State:</span>
                      <span className="font-medium text-gray-900 dark:text-gray-100 capitalize">
                        {array.state}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Devices:</span>
                      <span className="font-medium text-gray-900 dark:text-gray-100">
                        {array.active}/{array.devices}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Size:</span>
                      <span className="font-medium text-gray-900 dark:text-gray-100">
                        {formatBytes(array.size)}
                      </span>
                    </div>
                  </div>

                  {/* Resync Progress */}
                  {array.resync > 0 && array.resync < 100 && (
                    <div className="mt-3">
                      <div className="flex justify-between text-xs mb-1">
                        <span className="text-gray-600 dark:text-gray-400">Resyncing...</span>
                        <span className="font-medium text-gray-900 dark:text-gray-100">
                          {array.resync}%
                        </span>
                      </div>
                      <div className="h-1.5 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden">
                        <motion.div
                          initial={{ width: 0 }}
                          animate={{ width: `${array.resync}%` }}
                          className="h-full bg-blue-500"
                        />
                      </div>
                    </div>
                  )}
                </motion.div>
              ))}
            </div>
          )}
        </div>

        {/* Array Details */}
        <div className="flex-1 overflow-y-auto">
          {selectedArray ? (
            <div className="p-6">
              {/* Array Info Card */}
              <div className="bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 rounded-2xl p-6 mb-6">
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                      {selectedArray.name}
                    </h2>
                    <div className="flex items-center gap-2 mt-1">
                      {getStateIcon(selectedArray.state)}
                      <span className="text-sm text-gray-600 dark:text-gray-400 capitalize">
                        {selectedArray.state}
                      </span>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-3xl font-bold text-macos-blue">
                      {getRaidLevelInfo(selectedArray.level).name}
                    </div>
                    <div className="text-sm text-gray-600 dark:text-gray-400">
                      {getRaidLevelInfo(selectedArray.level).description}
                    </div>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Total Size</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {formatBytes(selectedArray.size)}
                    </div>
                  </div>
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Used</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {formatBytes(selectedArray.usedSize)}
                    </div>
                  </div>
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Devices</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedArray.devices}
                    </div>
                  </div>
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Active</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedArray.active}
                    </div>
                  </div>
                </div>
              </div>

              {/* Device Status */}
              <div>
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Device Status
                </h3>
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-gray-50 dark:bg-macos-dark-200 rounded-lg p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <HardDrive className="w-4 h-4 text-green-500" />
                      <span className="font-medium text-gray-900 dark:text-gray-100">Working</span>
                    </div>
                    <div className="text-2xl font-bold text-green-500">
                      {selectedArray.working}
                    </div>
                  </div>

                  <div className="bg-gray-50 dark:bg-macos-dark-200 rounded-lg p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <XCircle className="w-4 h-4 text-red-500" />
                      <span className="font-medium text-gray-900 dark:text-gray-100">Failed</span>
                    </div>
                    <div className="text-2xl font-bold text-red-500">
                      {selectedArray.failed}
                    </div>
                  </div>

                  <div className="bg-gray-50 dark:bg-macos-dark-200 rounded-lg p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Server className="w-4 h-4 text-blue-500" />
                      <span className="font-medium text-gray-900 dark:text-gray-100">Spare</span>
                    </div>
                    <div className="text-2xl font-bold text-blue-500">
                      {selectedArray.spare}
                    </div>
                  </div>

                  <div className="bg-gray-50 dark:bg-macos-dark-200 rounded-lg p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Zap className="w-4 h-4 text-yellow-500" />
                      <span className="font-medium text-gray-900 dark:text-gray-100">Active</span>
                    </div>
                    <div className="text-2xl font-bold text-yellow-500">
                      {selectedArray.active}
                    </div>
                  </div>
                </div>
              </div>

              {/* Resync Progress (if active) */}
              {selectedArray.resync > 0 && selectedArray.resync < 100 && (
                <div className="mt-6 bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      Resync in Progress
                    </span>
                    <span className="text-sm font-bold text-blue-600 dark:text-blue-400">
                      {selectedArray.resync}%
                    </span>
                  </div>
                  <div className="h-2 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${selectedArray.resync}%` }}
                      className="h-full bg-blue-500"
                    />
                  </div>
                  <p className="text-xs text-gray-600 dark:text-gray-400 mt-2">
                    Array is being rebuilt. Performance may be degraded.
                  </p>
                </div>
              )}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center h-full text-center p-12">
              <Shield className="w-24 h-24 text-gray-300 dark:text-gray-600 mb-4" />
              <p className="text-lg text-gray-500 dark:text-gray-400">
                Select an array to view details
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Create Array Dialog */}
      {showCreateDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-2xl w-full m-4 max-h-[80vh] overflow-y-auto"
          >
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Create RAID Array
              </h3>
              <button
                onClick={() => setShowCreateDialog(false)}
                className="p-1 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              {/* Array Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Array Name
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="md0"
                  className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent"
                />
              </div>

              {/* RAID Level */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  RAID Level
                </label>
                <div className="grid grid-cols-2 gap-2">
                  {(['0', '1', '5', '6', '10'] as const).map((level) => {
                    const info = getRaidLevelInfo(level);
                    return (
                      <button
                        key={level}
                        onClick={() => setFormData({ ...formData, level })}
                        className={`p-3 rounded-lg border-2 transition-all text-left ${
                          formData.level === level
                            ? 'border-macos-blue bg-macos-blue/10 dark:bg-macos-blue/20'
                            : 'border-gray-300 dark:border-gray-600 hover:border-macos-blue/50'
                        }`}
                      >
                        <div className="font-semibold text-gray-900 dark:text-gray-100">
                          {info.name}
                        </div>
                        <div className="text-xs text-gray-600 dark:text-gray-400">
                          {info.description}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                          Min: {info.minDisks} disks
                        </div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Device Selection */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Select Devices ({formData.devices.length} selected)
                </label>
                <div className="space-y-2 max-h-48 overflow-y-auto">
                  {availableDisks.map((disk) => (
                    <div
                      key={disk.path}
                      className="flex items-center justify-between p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <input
                          type="checkbox"
                          checked={formData.devices.includes(disk.path)}
                          onChange={() => toggleDevice(disk.path, 'device')}
                          className="w-4 h-4 text-macos-blue rounded focus:ring-2 focus:ring-macos-blue"
                        />
                        <div>
                          <div className="font-medium text-gray-900 dark:text-gray-100">
                            {disk.name} - {disk.model}
                          </div>
                          <div className="text-xs text-gray-600 dark:text-gray-400">
                            {formatBytes(disk.size)} • {disk.type.toUpperCase()}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Spare Devices */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Hot Spare Devices (Optional)
                </label>
                <div className="space-y-2 max-h-32 overflow-y-auto">
                  {availableDisks
                    .filter((disk) => !formData.devices.includes(disk.path))
                    .map((disk) => (
                      <div
                        key={disk.path}
                        className="flex items-center gap-3 p-2 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                      >
                        <input
                          type="checkbox"
                          checked={formData.spare?.includes(disk.path)}
                          onChange={() => toggleDevice(disk.path, 'spare')}
                          className="w-4 h-4 text-macos-blue rounded focus:ring-2 focus:ring-macos-blue"
                        />
                        <div className="text-sm">
                          <div className="font-medium text-gray-900 dark:text-gray-100">
                            {disk.name}
                          </div>
                          <div className="text-xs text-gray-600 dark:text-gray-400">
                            {formatBytes(disk.size)}
                          </div>
                        </div>
                      </div>
                    ))}
                </div>
              </div>
            </div>

            <div className="mt-6 flex justify-end gap-2">
              <button
                onClick={() => setShowCreateDialog(false)}
                className="px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateArray}
                className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
              >
                Create Array
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}
