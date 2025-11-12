import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { backupApi, BackupHistory as BackupHistoryType } from '../../../api/backup';

const BackupHistory: React.FC = () => {
  const [history, setHistory] = useState<BackupHistoryType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [limit, setLimit] = useState(50);

  const fetchHistory = async () => {
    try {
      setLoading(true);
      const response = await backupApi.getHistory(undefined, limit);
      setHistory(response.data || []);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch backup history');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHistory();
    const interval = setInterval(fetchHistory, 15000); // Refresh every 15s
    return () => clearInterval(interval);
  }, [limit]);

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'running':
        return 'bg-blue-500/20 text-blue-400 border-blue-500/30';
      case 'success':
        return 'bg-green-500/20 text-green-400 border-green-500/30';
      case 'failed':
        return 'bg-red-500/20 text-red-400 border-red-500/30';
      default:
        return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
    }
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDuration = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    if (hours > 0) {
      return `${hours}h ${minutes}m ${secs}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${secs}s`;
    }
    return `${secs}s`;
  };

  if (loading && history.length === 0) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-gray-400">Loading backup history...</div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-4">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
            Backup History
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            {history.length} backup execution(s)
          </p>
        </div>
        <select
          value={limit}
          onChange={(e) => setLimit(parseInt(e.target.value))}
          className="px-3 py-2 bg-gray-50 dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-gray-100"
        >
          <option value={25}>Last 25</option>
          <option value={50}>Last 50</option>
          <option value={100}>Last 100</option>
        </select>
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

      {/* History List */}
      {history.length === 0 ? (
        <div className="text-center py-12 text-gray-400">
          No backup history available. Run a backup job to see history here.
        </div>
      ) : (
        <div className="space-y-3">
          {history.map((entry) => (
            <motion.div
              key={entry.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-4"
            >
              <div className="flex justify-between items-start mb-3">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                    {entry.jobName}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                    {new Date(entry.startTime).toLocaleString()}
                  </p>
                </div>
                <span
                  className={`px-2 py-1 text-xs rounded-md border ${getStatusBadge(
                    entry.status
                  )}`}
                >
                  {entry.status}
                </span>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                <div>
                  <p className="text-gray-600 dark:text-gray-400 mb-1">Duration</p>
                  <p className="text-gray-900 dark:text-gray-100 font-medium">
                    {formatDuration(entry.duration)}
                  </p>
                </div>
                <div>
                  <p className="text-gray-600 dark:text-gray-400 mb-1">Size</p>
                  <p className="text-gray-900 dark:text-gray-100 font-medium">
                    {formatBytes(entry.bytesBackup)}
                  </p>
                </div>
                <div>
                  <p className="text-gray-600 dark:text-gray-400 mb-1">Files</p>
                  <p className="text-gray-900 dark:text-gray-100 font-medium">
                    {entry.filesBackup.toLocaleString()}
                  </p>
                </div>
                <div>
                  <p className="text-gray-600 dark:text-gray-400 mb-1">Path</p>
                  <p className="text-gray-900 dark:text-gray-100 font-medium truncate">
                    {entry.backupPath}
                  </p>
                </div>
              </div>

              {entry.error && (
                <div className="mt-3 p-3 bg-red-500/10 border border-red-500/30 rounded">
                  <p className="text-sm text-red-400">{entry.error}</p>
                </div>
              )}
            </motion.div>
          ))}
        </div>
      )}
    </div>
  );
};

export default BackupHistory;
