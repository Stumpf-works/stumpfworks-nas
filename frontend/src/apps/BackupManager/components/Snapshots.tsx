import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { backupApi, Snapshot } from '../../../api/backup';

const Snapshots: React.FC = () => {
  const [snapshots, setSnapshots] = useState<Snapshot[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showRestoreModal, setShowRestoreModal] = useState(false);
  const [selectedSnapshot, setSelectedSnapshot] = useState<Snapshot | null>(null);
  const [filesystem, setFilesystem] = useState('');
  const [snapshotName, setSnapshotName] = useState('');
  const [restoreDestination, setRestoreDestination] = useState('');

  const fetchSnapshots = async () => {
    try {
      setLoading(true);
      const response = await backupApi.listSnapshots();
      setSnapshots(response.data || []);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch snapshots');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSnapshots();
    const interval = setInterval(fetchSnapshots, 15000); // Refresh every 15s
    return () => clearInterval(interval);
  }, []);

  const handleCreate = async () => {
    try {
      await backupApi.createSnapshot({ filesystem, name: snapshotName });
      setShowCreateModal(false);
      setFilesystem('');
      setSnapshotName('');
      fetchSnapshots();
    } catch (err: any) {
      setError(err.message || 'Failed to create snapshot');
    }
  };

  const handleDelete = async () => {
    if (!selectedSnapshot) return;

    try {
      await backupApi.deleteSnapshot(selectedSnapshot.id);
      setShowDeleteModal(false);
      setSelectedSnapshot(null);
      fetchSnapshots();
    } catch (err: any) {
      setError(err.message || 'Failed to delete snapshot');
    }
  };

  const handleRestore = async () => {
    if (!selectedSnapshot) return;

    try {
      await backupApi.restoreSnapshot(selectedSnapshot.id, {
        destination: restoreDestination,
      });
      setShowRestoreModal(false);
      setSelectedSnapshot(null);
      setRestoreDestination('');
      fetchSnapshots();
    } catch (err: any) {
      setError(err.message || 'Failed to restore snapshot');
    }
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  if (loading && snapshots.length === 0) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-gray-400">Loading snapshots...</div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-4">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">Snapshots</h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            {snapshots.length} snapshot(s)
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
        >
          Create Snapshot
        </button>
      </div>

      {/* Error Display */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="p-4 bg-red-500/10 border border-red-500/30 rounded-lg"
          >
            <div className="flex justify-between items-start">
              <p className="text-red-400">{error}</p>
              <button
                onClick={() => setError(null)}
                className="text-red-400 hover:text-red-300"
              >
                ✕
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Snapshots List */}
      {snapshots.length === 0 ? (
        <div className="text-center py-12 text-gray-400">
          No snapshots available. Create your first snapshot to get started.
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {snapshots.map((snapshot) => (
            <motion.div
              key={snapshot.id}
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-4"
            >
              {/* Snapshot Header */}
              <div className="flex justify-between items-start mb-3">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                    {snapshot.name}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                    {snapshot.filesystem}
                  </p>
                </div>
                <span className="px-2 py-1 text-xs rounded-md border bg-blue-500/20 text-blue-400 border-blue-500/30">
                  {snapshot.type}
                </span>
              </div>

              {/* Snapshot Details */}
              <div className="space-y-2 text-sm mb-4">
                <div className="flex justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Created:</span>
                  <span className="text-gray-900 dark:text-gray-100">
                    {new Date(snapshot.createdAt).toLocaleString()}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Size:</span>
                  <span className="text-gray-900 dark:text-gray-100">
                    {formatBytes(snapshot.size)}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Used:</span>
                  <span className="text-gray-900 dark:text-gray-100">
                    {formatBytes(snapshot.used)}
                  </span>
                </div>
                {snapshot.description && (
                  <div className="pt-2 border-t border-gray-200 dark:border-gray-700">
                    <p className="text-gray-600 dark:text-gray-400 text-xs">
                      {snapshot.description}
                    </p>
                  </div>
                )}
              </div>

              {/* Snapshot Actions */}
              <div className="flex gap-2">
                <button
                  onClick={() => {
                    setSelectedSnapshot(snapshot);
                    setShowRestoreModal(true);
                  }}
                  className="flex-1 px-3 py-2 text-sm bg-green-500/20 text-green-600 dark:text-green-400 hover:bg-green-500/30 rounded-lg transition-colors"
                >
                  Restore
                </button>
                <button
                  onClick={() => {
                    setSelectedSnapshot(snapshot);
                    setShowDeleteModal(true);
                  }}
                  className="px-3 py-2 text-sm bg-red-500/20 text-red-600 dark:text-red-400 hover:bg-red-500/30 rounded-lg transition-colors"
                >
                  Delete
                </button>
              </div>
            </motion.div>
          ))}
        </div>
      )}

      {/* Create Modal */}
      <AnimatePresence>
        {showCreateModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowCreateModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-6 max-w-md w-full"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                Create Snapshot
              </h3>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                    Filesystem
                  </label>
                  <input
                    type="text"
                    value={filesystem}
                    onChange={(e) => setFilesystem(e.target.value)}
                    placeholder="tank/data"
                    className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                  />
                </div>

                <div>
                  <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                    Snapshot Name
                  </label>
                  <input
                    type="text"
                    value={snapshotName}
                    onChange={(e) => setSnapshotName(e.target.value)}
                    placeholder="backup-2024-01-15"
                    className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                  />
                </div>
              </div>

              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => {
                    setShowCreateModal(false);
                    setFilesystem('');
                    setSnapshotName('');
                  }}
                  className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreate}
                  disabled={!filesystem.trim() || !snapshotName.trim()}
                  className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Create
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Restore Modal */}
      <AnimatePresence>
        {showRestoreModal && selectedSnapshot && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowRestoreModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-6 max-w-md w-full"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                Restore Snapshot
              </h3>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Restore snapshot "{selectedSnapshot.name}" from {selectedSnapshot.filesystem}
              </p>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm text-gray-600 dark:text-gray-400 mb-2">
                    Destination Path
                  </label>
                  <input
                    type="text"
                    value={restoreDestination}
                    onChange={(e) => setRestoreDestination(e.target.value)}
                    placeholder="/path/to/restore"
                    className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
                  />
                </div>
              </div>

              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => {
                    setShowRestoreModal(false);
                    setSelectedSnapshot(null);
                    setRestoreDestination('');
                  }}
                  className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleRestore}
                  className="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors"
                >
                  Restore
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Delete Confirmation Modal */}
      <AnimatePresence>
        {showDeleteModal && selectedSnapshot && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={() => setShowDeleteModal(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg p-6 max-w-md w-full"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                Delete Snapshot
              </h3>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Are you sure you want to delete snapshot "{selectedSnapshot.name}"? This cannot be
                undone.
              </p>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => {
                    setShowDeleteModal(false);
                    setSelectedSnapshot(null);
                  }}
                  className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleDelete}
                  className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors"
                >
                  Delete
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default Snapshots;
