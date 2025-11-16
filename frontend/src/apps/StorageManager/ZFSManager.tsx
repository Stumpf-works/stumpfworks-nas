// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Database,
  Plus,
  RefreshCw,
  Trash2,
  HardDrive,
  Activity,
  Camera,
  AlertTriangle,
  CheckCircle,
  XCircle,
} from 'lucide-react';
import { syslibApi, type ZFSPool, type ZFSDataset, type CreateZFSPoolRequest } from '@/api/syslib';

export default function ZFSManager() {
  const [pools, setPools] = useState<ZFSPool[]>([]);
  const [selectedPool, setSelectedPool] = useState<ZFSPool | null>(null);
  const [datasets, setDatasets] = useState<ZFSDataset[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
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
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-100">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center gap-3">
          <Database className="w-6 h-6 text-macos-blue" />
          <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
            ZFS Pool Manager
          </h1>
        </div>
        <div className="flex gap-2">
          <button
            onClick={fetchPools}
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
            Create Pool
          </button>
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
              {pools.map((pool) => (
                <motion.div
                  key={pool.name}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  onClick={() => setSelectedPool(pool)}
                  className={`p-4 rounded-xl cursor-pointer transition-all ${
                    selectedPool?.name === pool.name
                      ? 'bg-macos-blue/10 dark:bg-macos-blue/20 border-2 border-macos-blue'
                      : 'bg-gray-50 dark:bg-macos-dark-200 hover:bg-gray-100 dark:hover:bg-macos-dark-300 border-2 border-transparent'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      {getHealthIcon(pool.health)}
                      <span className="font-semibold text-gray-900 dark:text-gray-100">
                        {pool.name}
                      </span>
                    </div>
                    <span className="text-xs text-gray-500 dark:text-gray-400 uppercase">
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
                  <div className="mt-3 h-2 bg-gray-200 dark:bg-macos-dark-300 rounded-full overflow-hidden">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${pool.capacity}%` }}
                      className={`h-full rounded-full ${
                        pool.capacity > 90
                          ? 'bg-red-500'
                          : pool.capacity > 75
                          ? 'bg-yellow-500'
                          : 'bg-macos-blue'
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
              <div className="bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 rounded-2xl p-6 mb-6">
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                      {selectedPool.name}
                    </h2>
                    <div className="flex items-center gap-2 mt-1">
                      {getHealthIcon(selectedPool.health)}
                      <span className="text-sm text-gray-600 dark:text-gray-400">
                        {selectedPool.health}
                      </span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleScrub(selectedPool.name)}
                      className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-macos-dark-300 transition-colors"
                    >
                      <Activity className="w-4 h-4" />
                      Scrub
                    </button>
                    <button
                      onClick={() => handleDestroyPool(selectedPool.name)}
                      className="flex items-center gap-2 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
                      Destroy
                    </button>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Total Size</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {formatBytes(selectedPool.size)}
                    </div>
                  </div>
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Used</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {formatBytes(selectedPool.allocated)}
                    </div>
                  </div>
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Free</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {formatBytes(selectedPool.free)}
                    </div>
                  </div>
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">Fragmentation</div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedPool.fragmentation}%
                    </div>
                  </div>
                </div>
              </div>

              {/* Datasets */}
              <div>
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Datasets
                </h3>
                {datasets.length === 0 ? (
                  <div className="text-center py-12 text-gray-500 dark:text-gray-400">
                    No datasets found
                  </div>
                ) : (
                  <div className="space-y-2">
                    {datasets.map((dataset) => (
                      <div
                        key={dataset.name}
                        className="bg-gray-50 dark:bg-macos-dark-200 rounded-lg p-4"
                      >
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-2">
                            <HardDrive className="w-4 h-4 text-macos-blue" />
                            <span className="font-medium text-gray-900 dark:text-gray-100">
                              {dataset.name}
                            </span>
                          </div>
                          <span className="text-xs px-2 py-1 bg-macos-blue/10 text-macos-blue rounded">
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
                            <span className="font-mono text-xs text-gray-900 dark:text-gray-100">
                              {dataset.mountpoint}
                            </span>
                          </div>
                        </div>
                      </div>
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
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-md w-full m-4"
          >
            <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
              Create ZFS Pool
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              Pool creation UI coming soon. Use CLI for now:
            </p>
            <code className="block bg-gray-100 dark:bg-macos-dark-200 p-3 rounded text-sm font-mono">
              zpool create tank raidz sda sdb sdc
            </code>
            <div className="mt-6 flex justify-end">
              <button
                onClick={() => setShowCreateDialog(false)}
                className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
              >
                Close
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}
