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
    <div className="flex flex-col h-full bg-gradient-to-br from-gray-50 to-white dark:from-macos-dark-100 dark:to-macos-dark-200">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200/50 dark:border-gray-700/50 bg-white/50 dark:bg-macos-dark-100/50 backdrop-blur-sm">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl shadow-lg">
            <Shield className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              RAID Array Manager
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Redundant array management
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <motion.button
            whileHover={{ scale: 1.05, y: -2 }}
            whileTap={{ scale: 0.95 }}
            onClick={fetchArrays}
            className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-xl hover:shadow-lg transition-all border border-gray-200 dark:border-gray-700"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.05, y: -2 }}
            whileTap={{ scale: 0.95 }}
            onClick={() => setShowCreateDialog(true)}
            className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-macos-blue to-macos-purple text-white rounded-xl hover:shadow-lg transition-all"
          >
            <Plus className="w-4 h-4" />
            Create Array
          </motion.button>
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
              {arrays.map((array, index) => (
                <motion.div
                  key={array.name}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.05 }}
                  whileHover={{ x: 4, scale: 1.02 }}
                  onClick={() => setSelectedArray(array)}
                  className={`p-4 rounded-xl cursor-pointer transition-all ${
                    selectedArray?.name === array.name
                      ? 'bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 border-2 border-macos-blue shadow-lg'
                      : 'bg-white dark:bg-macos-dark-200 hover:bg-gray-50 dark:hover:bg-macos-dark-300 border-2 border-gray-200 dark:border-gray-700 hover:shadow-md'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      {getStateIcon(array.state)}
                      <span className="font-semibold text-gray-900 dark:text-gray-100">
                        {array.name}
                      </span>
                    </div>
                    <span className="text-xs px-3 py-1 bg-gradient-to-r from-macos-blue/10 to-macos-purple/10 text-macos-blue rounded-full font-medium">
                      RAID {array.level}
                    </span>
                  </div>

                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">State:</span>
                      <span className={`font-medium capitalize px-2 py-0.5 rounded ${
                        array.state.toLowerCase() === 'clean' || array.state.toLowerCase() === 'active'
                          ? 'bg-green-500/10 text-green-600 dark:text-green-400'
                          : array.state.toLowerCase() === 'degraded' || array.state.toLowerCase() === 'recovering'
                          ? 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400'
                          : 'bg-red-500/10 text-red-600 dark:text-red-400'
                      }`}>
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
                      <div className="h-2 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden shadow-inner">
                        <motion.div
                          initial={{ width: 0 }}
                          animate={{ width: `${array.resync}%` }}
                          transition={{ duration: 0.8, ease: 'easeOut' }}
                          className="h-full bg-gradient-to-r from-blue-500 to-blue-600"
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
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="bg-gradient-to-br from-macos-blue/10 via-macos-purple/10 to-pink-500/5 dark:from-macos-blue/20 dark:via-macos-purple/20 dark:to-pink-500/10 rounded-2xl p-6 mb-6 border border-gray-200/50 dark:border-gray-700/50 shadow-xl"
              >
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                      {selectedArray.name}
                    </h2>
                    <div className="flex items-center gap-2 mt-1">
                      {getStateIcon(selectedArray.state)}
                      <span className={`text-sm font-medium capitalize px-2 py-1 rounded-full ${
                        selectedArray.state.toLowerCase() === 'clean' || selectedArray.state.toLowerCase() === 'active'
                          ? 'bg-green-500/10 text-green-600 dark:text-green-400'
                          : selectedArray.state.toLowerCase() === 'degraded' || selectedArray.state.toLowerCase() === 'recovering'
                          ? 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400'
                          : 'bg-red-500/10 text-red-600 dark:text-red-400'
                      }`}>
                        {selectedArray.state}
                      </span>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-3xl font-bold bg-gradient-to-r from-macos-blue to-macos-purple bg-clip-text text-transparent">
                      {getRaidLevelInfo(selectedArray.level).name}
                    </div>
                    <div className="text-sm text-gray-600 dark:text-gray-400">
                      {getRaidLevelInfo(selectedArray.level).description}
                    </div>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.1 }}
                    className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                  >
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Total Size</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {formatBytes(selectedArray.size)}
                    </div>
                  </motion.div>
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.2 }}
                    className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                  >
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Used</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {formatBytes(selectedArray.usedSize)}
                    </div>
                  </motion.div>
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.3 }}
                    className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                  >
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Devices</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedArray.devices}
                    </div>
                  </motion.div>
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.4 }}
                    className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                  >
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Active</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedArray.active}
                    </div>
                  </motion.div>
                </div>
              </motion.div>

              {/* Device Status */}
              <div>
                <div className="flex items-center gap-2 mb-4">
                  <HardDrive className="w-5 h-5 text-macos-blue" />
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    Device Status
                  </h3>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.1 }}
                    whileHover={{ scale: 1.05, y: -2 }}
                    className="bg-gradient-to-br from-green-50 to-green-100/50 dark:from-green-900/20 dark:to-green-800/20 rounded-xl p-4 border border-green-200 dark:border-green-700 shadow-md hover:shadow-lg transition-all"
                  >
                    <div className="flex items-center gap-2 mb-2">
                      <div className="p-2 bg-green-500/10 rounded-lg">
                        <HardDrive className="w-4 h-4 text-green-600 dark:text-green-400" />
                      </div>
                      <span className="font-medium text-gray-900 dark:text-gray-100">Working</span>
                    </div>
                    <div className="text-2xl font-bold text-green-600 dark:text-green-400">
                      {selectedArray.working}
                    </div>
                  </motion.div>

                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.2 }}
                    whileHover={{ scale: 1.05, y: -2 }}
                    className="bg-gradient-to-br from-red-50 to-red-100/50 dark:from-red-900/20 dark:to-red-800/20 rounded-xl p-4 border border-red-200 dark:border-red-700 shadow-md hover:shadow-lg transition-all"
                  >
                    <div className="flex items-center gap-2 mb-2">
                      <div className="p-2 bg-red-500/10 rounded-lg">
                        <XCircle className="w-4 h-4 text-red-600 dark:text-red-400" />
                      </div>
                      <span className="font-medium text-gray-900 dark:text-gray-100">Failed</span>
                    </div>
                    <div className="text-2xl font-bold text-red-600 dark:text-red-400">
                      {selectedArray.failed}
                    </div>
                  </motion.div>

                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.3 }}
                    whileHover={{ scale: 1.05, y: -2 }}
                    className="bg-gradient-to-br from-blue-50 to-blue-100/50 dark:from-blue-900/20 dark:to-blue-800/20 rounded-xl p-4 border border-blue-200 dark:border-blue-700 shadow-md hover:shadow-lg transition-all"
                  >
                    <div className="flex items-center gap-2 mb-2">
                      <div className="p-2 bg-blue-500/10 rounded-lg">
                        <Server className="w-4 h-4 text-blue-600 dark:text-blue-400" />
                      </div>
                      <span className="font-medium text-gray-900 dark:text-gray-100">Spare</span>
                    </div>
                    <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">
                      {selectedArray.spare}
                    </div>
                  </motion.div>

                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.4 }}
                    whileHover={{ scale: 1.05, y: -2 }}
                    className="bg-gradient-to-br from-yellow-50 to-yellow-100/50 dark:from-yellow-900/20 dark:to-yellow-800/20 rounded-xl p-4 border border-yellow-200 dark:border-yellow-700 shadow-md hover:shadow-lg transition-all"
                  >
                    <div className="flex items-center gap-2 mb-2">
                      <div className="p-2 bg-yellow-500/10 rounded-lg">
                        <Zap className="w-4 h-4 text-yellow-600 dark:text-yellow-400" />
                      </div>
                      <span className="font-medium text-gray-900 dark:text-gray-100">Active</span>
                    </div>
                    <div className="text-2xl font-bold text-yellow-600 dark:text-yellow-400">
                      {selectedArray.active}
                    </div>
                  </motion.div>
                </div>
              </div>

              {/* Resync Progress (if active) */}
              {selectedArray.resync > 0 && selectedArray.resync < 100 && (
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="mt-6 bg-gradient-to-br from-blue-50 to-blue-100/50 dark:from-blue-900/20 dark:to-blue-800/20 rounded-xl p-4 border border-blue-200 dark:border-blue-700"
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <Activity className="w-5 h-5 text-blue-600 dark:text-blue-400" />
                      <span className="font-medium text-gray-900 dark:text-gray-100">
                        Resync in Progress
                      </span>
                    </div>
                    <span className="text-sm font-bold text-blue-600 dark:text-blue-400">
                      {selectedArray.resync}%
                    </span>
                  </div>
                  <div className="h-2 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden shadow-inner">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${selectedArray.resync}%` }}
                      transition={{ duration: 0.8, ease: 'easeOut' }}
                      className="h-full bg-gradient-to-r from-blue-500 to-blue-600"
                    />
                  </div>
                  <p className="text-xs text-gray-600 dark:text-gray-400 mt-2">
                    Array is being rebuilt. Performance may be degraded.
                  </p>
                </motion.div>
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
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-2xl w-full m-4 max-h-[80vh] overflow-y-auto shadow-2xl border border-gray-200 dark:border-gray-700"
          >
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-gradient-to-br from-macos-blue to-macos-purple rounded-lg">
                  <Shield className="w-5 h-5 text-white" />
                </div>
                <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                  Create RAID Array
                </h3>
              </div>
              <motion.button
                whileHover={{ scale: 1.1, rotate: 90 }}
                whileTap={{ scale: 0.9 }}
                onClick={() => setShowCreateDialog(false)}
                className="p-1 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded-lg transition-colors"
              >
                <X className="w-5 h-5" />
              </motion.button>
            </div>

            <div className="space-y-4">
              {/* Array Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Array Name
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="md0"
                  className="w-full px-4 py-2.5 bg-gray-50 dark:bg-macos-dark-200 border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-macos-blue focus:border-macos-blue transition-all"
                />
              </div>

              {/* RAID Level */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  RAID Level
                </label>
                <div className="grid grid-cols-2 gap-3">
                  {(['0', '1', '5', '6', '10'] as const).map((level, index) => {
                    const info = getRaidLevelInfo(level);
                    return (
                      <motion.button
                        key={level}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.05 }}
                        whileHover={{ scale: 1.02, y: -2 }}
                        whileTap={{ scale: 0.98 }}
                        onClick={() => setFormData({ ...formData, level })}
                        className={`p-3 rounded-xl border-2 transition-all text-left ${
                          formData.level === level
                            ? 'border-macos-blue bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 shadow-md'
                            : 'border-gray-300 dark:border-gray-600 hover:border-macos-blue/50 bg-gray-50 dark:bg-macos-dark-200'
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
                      </motion.button>
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
              <motion.button
                whileHover={{ scale: 1.05, y: -2 }}
                whileTap={{ scale: 0.95 }}
                onClick={() => setShowCreateDialog(false)}
                className="px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-xl hover:shadow-lg transition-all border border-gray-200 dark:border-gray-700"
              >
                Cancel
              </motion.button>
              <motion.button
                whileHover={{ scale: 1.05, y: -2 }}
                whileTap={{ scale: 0.95 }}
                onClick={handleCreateArray}
                className="px-4 py-2 bg-gradient-to-r from-macos-blue to-macos-purple text-white rounded-xl hover:shadow-lg transition-all"
              >
                Create Array
              </motion.button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}
