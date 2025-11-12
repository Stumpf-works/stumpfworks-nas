import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { dockerApi, DockerVolume } from '@/api/docker';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';

export default function VolumeManager() {
  const [volumes, setVolumes] = useState<DockerVolume[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [deleteModal, setDeleteModal] = useState<DockerVolume | null>(null);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  useEffect(() => {
    loadVolumes();
    const interval = setInterval(loadVolumes, 5000); // Refresh every 5s
    return () => clearInterval(interval);
  }, []);

  const loadVolumes = async () => {
    try {
      const response = await dockerApi.listVolumes();
      if (response.success && response.data) {
        setVolumes(response.data);
        setError('');
      } else {
        setError(response.error?.message || 'Failed to load volumes');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (volume: DockerVolume) => {
    setActionLoading(volume.name);
    try {
      const response = await dockerApi.removeVolume(volume.name);
      if (response.success) {
        setDeleteModal(null);
        loadVolumes();
      } else {
        alert(response.error?.message || 'Failed to remove volume');
      }
    } catch (err) {
      alert(getErrorMessage(err));
    } finally {
      setActionLoading(null);
    }
  };

  const formatDate = (dateStr: string) => {
    if (!dateStr) return 'N/A';
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  const getDriverColor = (driver: string) => {
    switch (driver.toLowerCase()) {
      case 'local':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400';
      case 'nfs':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'cifs':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400';
    }
  };

  const getScopeColor = (scope: string) => {
    switch (scope.toLowerCase()) {
      case 'local':
        return 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900/30 dark:text-cyan-400';
      case 'global':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400';
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
    <div className="p-6 space-y-6">
      {/* Error Display */}
      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-600 dark:text-gray-400">
          {volumes.length} volume{volumes.length !== 1 ? 's' : ''}
        </div>
      </div>

      {/* Volumes Grid */}
      {volumes.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">💾</div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            No volumes found
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            Docker volumes will appear here when created
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
          {volumes.map((volume) => (
            <Card key={volume.name} hoverable>
              <div className="p-6">
                {/* Header */}
                <div className="mb-4">
                  <h3 className="font-bold text-lg text-gray-900 dark:text-gray-100 mb-2">
                    {volume.name}
                  </h3>
                  <div className="flex gap-2">
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${getDriverColor(
                        volume.driver
                      )}`}
                    >
                      {volume.driver}
                    </span>
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${getScopeColor(
                        volume.scope
                      )}`}
                    >
                      {volume.scope}
                    </span>
                  </div>
                </div>

                {/* Details */}
                <div className="space-y-2 mb-4">
                  <div className="text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Mountpoint:</span>
                    <div className="mt-1 p-2 bg-gray-50 dark:bg-gray-800 rounded font-mono text-xs text-gray-900 dark:text-gray-100 break-all">
                      {volume.mountpoint}
                    </div>
                  </div>
                  {volume.createdAt && (
                    <div className="text-sm">
                      <span className="text-gray-600 dark:text-gray-400">Created:</span>
                      <span className="ml-2 font-medium text-gray-900 dark:text-gray-100">
                        {formatDate(volume.createdAt)}
                      </span>
                    </div>
                  )}
                </div>

                {/* Labels */}
                {volume.labels && Object.keys(volume.labels).length > 0 && (
                  <div className="mb-4">
                    <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                      Labels:
                    </div>
                    <div className="space-y-1">
                      {Object.entries(volume.labels).map(([key, value]) => (
                        <div
                          key={key}
                          className="text-xs p-2 bg-gray-50 dark:bg-gray-800 rounded"
                        >
                          <span className="font-medium text-gray-900 dark:text-gray-100">
                            {key}:
                          </span>{' '}
                          <span className="text-gray-600 dark:text-gray-400">{value}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Options */}
                {volume.options && Object.keys(volume.options).length > 0 && (
                  <div className="mb-4">
                    <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                      Options:
                    </div>
                    <div className="flex flex-wrap gap-1">
                      {Object.entries(volume.options).map(([key, value]) => (
                        <span
                          key={key}
                          className="px-2 py-0.5 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs rounded"
                        >
                          {key}={value}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Actions */}
                <Button
                  size="sm"
                  variant="danger"
                  onClick={() => setDeleteModal(volume)}
                  disabled={actionLoading === volume.name}
                  className="w-full"
                >
                  🗑️ Delete
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Delete Confirmation Modal */}
      <AnimatePresence>
        {deleteModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => setDeleteModal(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Delete Volume
              </h2>
              <p className="text-gray-600 dark:text-gray-400 mb-6">
                Are you sure you want to delete volume <strong>{deleteModal.name}</strong>? This
                action cannot be undone and all data will be lost.
              </p>
              <div className="flex gap-3">
                <Button
                  variant="secondary"
                  onClick={() => setDeleteModal(null)}
                  className="flex-1"
                  disabled={actionLoading === deleteModal.name}
                >
                  Cancel
                </Button>
                <Button
                  variant="danger"
                  onClick={() => handleDelete(deleteModal)}
                  className="flex-1"
                  disabled={actionLoading === deleteModal.name}
                >
                  {actionLoading === deleteModal.name ? 'Deleting...' : 'Delete'}
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
