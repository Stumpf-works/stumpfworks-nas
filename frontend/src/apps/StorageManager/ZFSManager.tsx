// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  Database,
  Plus,
  RefreshCw,
  Trash2,
  HardDrive,
  Activity,
  AlertTriangle,
  CheckCircle,
  XCircle,
} from 'lucide-react';
import { syslibApi, type ZFSPool, type ZFSDataset } from '@/api/syslib';

export default function ZFSManager() {
  const [pools, setPools] = useState<ZFSPool[]>([]);
  const [selectedPool, setSelectedPool] = useState<ZFSPool | null>(null);
  const [datasets, setDatasets] = useState<ZFSDataset[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateDialog, setShowCreateDialog] = useState(false);

  // Fetch pools
  const fetchPools = async () => {
    setIsLoading(true);
    try {
      const response = await syslibApi.zfs.listPools();
      if (response.success && response.data) {
        setPools(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch ZFS pools:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch datasets for selected pool
  const fetchDatasets = async (poolName: string) => {
    try {
      const response = await syslibApi.zfs.listDatasets(poolName);
      if (response.success && response.data) {
        setDatasets(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch datasets:', error);
    }
  };

  useEffect(() => {
    fetchPools();
  }, []);

  useEffect(() => {
    if (selectedPool) {
      fetchDatasets(selectedPool.name);
    } else {
      setDatasets([]);
    }
  }, [selectedPool]);

  const handleScrub = async (poolName: string) => {
    try {
      const response = await syslibApi.zfs.scrubPool(poolName);
      if (response.success) {
        alert(`Scrub started for pool: ${poolName}`);
      }
    } catch (error) {
      console.error('Failed to start scrub:', error);
      alert('Failed to start scrub');
    }
  };

  const handleDestroyPool = async (poolName: string) => {
    if (!confirm(`Are you sure you want to destroy pool "${poolName}"? This action cannot be undone!`)) {
      return;
    }

    try {
      const response = await syslibApi.zfs.destroyPool(poolName, false);
      if (response.success) {
        alert('Pool destroyed successfully');
        setSelectedPool(null);
        fetchPools();
      }
    } catch (error) {
      console.error('Failed to destroy pool:', error);
      alert('Failed to destroy pool');
    }
  };

  const getHealthIcon = (health: string) => {
    switch (health.toLowerCase()) {
      case 'online':
        return <CheckCircle className="w-5 h-5 text-green-500" />;
      case 'degraded':
        return <AlertTriangle className="w-5 h-5 text-yellow-500" />;
      case 'faulted':
      case 'offline':
        return <XCircle className="w-5 h-5 text-red-500" />;
      default:
        return <Activity className="w-5 h-5 text-gray-500" />;
    }
  };

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="flex flex-col h-full bg-gradient-to-br from-gray-50 to-white dark:from-macos-dark-100 dark:to-macos-dark-200">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200/50 dark:border-gray-700/50 bg-white/50 dark:bg-macos-dark-100/50 backdrop-blur-sm">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl shadow-lg">
            <Database className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              ZFS Pool Manager
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Advanced filesystem management
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <motion.button
            whileHover={{ scale: 1.05, y: -2 }}
            whileTap={{ scale: 0.95 }}
            onClick={fetchPools}
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
            Create Pool
          </motion.button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Pools List */}
        <div className="w-1/3 border-r border-gray-200 dark:border-gray-700 overflow-y-auto">
          {isLoading ? (
            <div className="flex items-center justify-center p-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
            </div>
          ) : pools.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center">
              <Database className="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" />
              <p className="text-gray-500 dark:text-gray-400">No ZFS pools found</p>
              <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
                Create a new pool to get started
              </p>
            </div>
          ) : (
            <div className="p-4 space-y-2">
              {pools.map((pool, index) => (
                <motion.div
                  key={pool.name}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.05 }}
                  whileHover={{ x: 4, scale: 1.02 }}
                  onClick={() => setSelectedPool(pool)}
                  className={`p-4 rounded-xl cursor-pointer transition-all ${
                    selectedPool?.name === pool.name
                      ? 'bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 border-2 border-macos-blue shadow-lg'
                      : 'bg-white dark:bg-macos-dark-200 hover:bg-gray-50 dark:hover:bg-macos-dark-300 border-2 border-gray-200 dark:border-gray-700 hover:shadow-md'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      {getHealthIcon(pool.health)}
                      <span className="font-semibold text-gray-900 dark:text-gray-100">
                        {pool.name}
                      </span>
                    </div>
                    <span className={`text-xs px-2 py-1 rounded-full font-medium ${
                      pool.health.toLowerCase() === 'online'
                        ? 'bg-green-500/10 text-green-600 dark:text-green-400'
                        : pool.health.toLowerCase() === 'degraded'
                        ? 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400'
                        : 'bg-red-500/10 text-red-600 dark:text-red-400'
                    }`}>
                      {pool.health}
                    </span>
                  </div>

                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Capacity:</span>
                      <span className="font-medium text-gray-900 dark:text-gray-100">
                        {pool.capacity}%
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Size:</span>
                      <span className="font-medium text-gray-900 dark:text-gray-100">
                        {formatBytes(pool.size)}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600 dark:text-gray-400">Free:</span>
                      <span className="font-medium text-gray-900 dark:text-gray-100">
                        {formatBytes(pool.free)}
                      </span>
                    </div>
                  </div>

                  {/* Progress Bar */}
                  <div className="mt-3 h-2 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden shadow-inner">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${pool.capacity}%` }}
                      transition={{ duration: 0.8, ease: 'easeOut' }}
                      className={`h-full rounded-full ${
                        pool.capacity > 90
                          ? 'bg-gradient-to-r from-red-500 to-rose-600'
                          : pool.capacity > 75
                          ? 'bg-gradient-to-r from-yellow-500 to-orange-500'
                          : 'bg-gradient-to-r from-macos-blue to-macos-purple'
                      }`}
                    />
                  </div>
                </motion.div>
              ))}
            </div>
          )}
        </div>

        {/* Pool Details */}
        <div className="flex-1 overflow-y-auto">
          {selectedPool ? (
            <div className="p-6">
              {/* Pool Info Card */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="bg-gradient-to-br from-macos-blue/10 via-macos-purple/10 to-pink-500/5 dark:from-macos-blue/20 dark:via-macos-purple/20 dark:to-pink-500/10 rounded-2xl p-6 mb-6 border border-gray-200/50 dark:border-gray-700/50 shadow-xl"
              >
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                      {selectedPool.name}
                    </h2>
                    <div className="flex items-center gap-2 mt-1">
                      {getHealthIcon(selectedPool.health)}
                      <span className={`text-sm font-medium px-2 py-1 rounded-full ${
                        selectedPool.health.toLowerCase() === 'online'
                          ? 'bg-green-500/10 text-green-600 dark:text-green-400'
                          : selectedPool.health.toLowerCase() === 'degraded'
                          ? 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400'
                          : 'bg-red-500/10 text-red-600 dark:text-red-400'
                      }`}>
                        {selectedPool.health}
                      </span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <motion.button
                      whileHover={{ scale: 1.05, y: -2 }}
                      whileTap={{ scale: 0.95 }}
                      onClick={() => handleScrub(selectedPool.name)}
                      className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-xl hover:shadow-lg transition-all border border-gray-200 dark:border-gray-700"
                    >
                      <Activity className="w-4 h-4" />
                      Scrub
                    </motion.button>
                    <motion.button
                      whileHover={{ scale: 1.05, y: -2 }}
                      whileTap={{ scale: 0.95 }}
                      onClick={() => handleDestroyPool(selectedPool.name)}
                      className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-red-500 to-rose-600 text-white rounded-xl hover:shadow-lg transition-all"
                    >
                      <Trash2 className="w-4 h-4" />
                      Destroy
                    </motion.button>
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
                      {formatBytes(selectedPool.size)}
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
                      {formatBytes(selectedPool.allocated)}
                    </div>
                  </motion.div>
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.3 }}
                    className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                  >
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Free</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {formatBytes(selectedPool.free)}
                    </div>
                  </motion.div>
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.4 }}
                    className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
                  >
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Fragmentation</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedPool.fragmentation}%
                    </div>
                  </motion.div>
                </div>
              </motion.div>

              {/* Datasets */}
              <div>
                <div className="flex items-center gap-2 mb-4">
                  <HardDrive className="w-5 h-5 text-macos-blue" />
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    Datasets
                  </h3>
                </div>
                {datasets.length === 0 ? (
                  <div className="text-center py-12 text-gray-500 dark:text-gray-400">
                    No datasets found
                  </div>
                ) : (
                  <div className="space-y-2">
                    {datasets.map((dataset, index) => (
                      <motion.div
                        key={dataset.name}
                        initial={{ opacity: 0, x: -20 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: index * 0.05 }}
                        whileHover={{ x: 4, scale: 1.01 }}
                        className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-4 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-all"
                      >
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-2">
                            <div className="p-2 bg-macos-blue/10 rounded-lg">
                              <HardDrive className="w-4 h-4 text-macos-blue" />
                            </div>
                            <span className="font-medium text-gray-900 dark:text-gray-100">
                              {dataset.name}
                            </span>
                          </div>
                          <span className="text-xs px-3 py-1 bg-gradient-to-r from-macos-blue/10 to-macos-purple/10 text-macos-blue rounded-full font-medium">
                            {dataset.type}
                          </span>
                        </div>
                        <div className="grid grid-cols-3 gap-4 text-sm">
                          <div>
                            <span className="text-gray-600 dark:text-gray-400">Used: </span>
                            <span className="font-medium text-gray-900 dark:text-gray-100">
                              {formatBytes(dataset.used)}
                            </span>
                          </div>
                          <div>
                            <span className="text-gray-600 dark:text-gray-400">Available: </span>
                            <span className="font-medium text-gray-900 dark:text-gray-100">
                              {formatBytes(dataset.available)}
                            </span>
                          </div>
                          <div>
                            <span className="text-gray-600 dark:text-gray-400">Mountpoint: </span>
                            <span className="font-mono text-xs text-gray-900 dark:text-gray-100 bg-gray-100 dark:bg-macos-dark-100 px-2 py-1 rounded">
                              {dataset.mountpoint}
                            </span>
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center h-full text-center p-12">
              <Database className="w-24 h-24 text-gray-300 dark:text-gray-600 mb-4" />
              <p className="text-lg text-gray-500 dark:text-gray-400">
                Select a pool to view details
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Create Pool Dialog - Placeholder */}
      {showCreateDialog && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-md w-full m-4 shadow-2xl border border-gray-200 dark:border-gray-700"
          >
            <div className="flex items-center gap-3 mb-4">
              <div className="p-2 bg-gradient-to-br from-macos-blue to-macos-purple rounded-lg">
                <Database className="w-5 h-5 text-white" />
              </div>
              <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Create ZFS Pool
              </h3>
            </div>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              Pool creation UI coming soon. Use CLI for now:
            </p>
            <code className="block bg-gradient-to-br from-gray-100 to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 p-3 rounded-lg text-sm font-mono border border-gray-200 dark:border-gray-700">
              zpool create tank raidz sda sdb sdc
            </code>
            <div className="mt-6 flex justify-end">
              <motion.button
                whileHover={{ scale: 1.05, y: -2 }}
                whileTap={{ scale: 0.95 }}
                onClick={() => setShowCreateDialog(false)}
                className="px-4 py-2 bg-gradient-to-r from-macos-blue to-macos-purple text-white rounded-xl hover:shadow-lg transition-all"
              >
                Close
              </motion.button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}
