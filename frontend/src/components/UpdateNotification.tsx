import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { systemApi, UpdateInfo } from '@/api/system';
import { useAuthStore } from '@/store';
import Button from '@/components/ui/Button';

export default function UpdateNotification() {
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [error, setError] = useState('');
  const user = useAuthStore((state) => state.user);

  const checkForUpdates = async () => {
    try {
      const response = await systemApi.checkUpdates();
      if (response.success && response.data) {
        setUpdateInfo(response.data);
      }
    } catch (err) {
      console.error('Failed to check for updates:', err);
    }
  };

  useEffect(() => {
    // Check for updates on mount
    checkForUpdates();

    // Check every hour
    const interval = setInterval(checkForUpdates, 60 * 60 * 1000);

    return () => clearInterval(interval);
  }, []);

  const handleApplyUpdate = async () => {
    setIsUpdating(true);
    setError('');
    try {
      const response = await systemApi.applyUpdates();
      if (response.success) {
        alert(
          'Update applied successfully! Please restart the server for changes to take effect.'
        );
        setShowModal(false);
        setUpdateInfo(null);
      } else {
        setError(response.error?.message || 'Failed to apply update');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to apply update');
    } finally {
      setIsUpdating(false);
    }
  };

  if (!updateInfo?.available) {
    return null;
  }

  const isAdmin = user?.role === 'admin';

  return (
    <>
      {/* Update Badge */}
      <motion.button
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        className="relative p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-200 transition-colors"
        onClick={() => setShowModal(true)}
        title="Update available"
      >
        <svg
          className="w-5 h-5 text-gray-700 dark:text-gray-300"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
        {/* Notification dot */}
        <div className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full animate-pulse" />
      </motion.button>

      {/* Update Modal */}
      <AnimatePresence>
        {showModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
            onClick={() => setShowModal(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-2xl mx-4 max-h-[80vh] overflow-auto"
            >
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                    Update Available
                  </h2>
                  <p className="text-gray-600 dark:text-gray-400 mt-1">
                    {updateInfo.behindBy} commit{updateInfo.behindBy > 1 ? 's' : ''}{' '}
                    behind
                  </p>
                </div>
                <button
                  onClick={() => setShowModal(false)}
                  className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                >
                  <svg
                    className="w-6 h-6"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              </div>

              {/* Version Info */}
              <div className="grid grid-cols-2 gap-4 mb-6 p-4 bg-gray-50 dark:bg-macos-dark-200 rounded-lg">
                <div>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Current Version
                  </p>
                  <p className="font-mono text-sm font-medium text-gray-900 dark:text-gray-100">
                    {updateInfo.currentCommit}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Latest Version
                  </p>
                  <p className="font-mono text-sm font-medium text-gray-900 dark:text-gray-100">
                    {updateInfo.latestCommit}
                  </p>
                </div>
              </div>

              {/* Changelog */}
              <div className="mb-6">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-3">
                  What's New
                </h3>
                <div className="space-y-2 max-h-60 overflow-y-auto">
                  {updateInfo.changeLog.length > 0 ? (
                    updateInfo.changeLog.map((commit, index) => (
                      <div
                        key={index}
                        className="p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                      >
                        <p className="text-sm font-mono text-gray-900 dark:text-gray-100">
                          {commit}
                        </p>
                      </div>
                    ))
                  ) : (
                    <p className="text-gray-600 dark:text-gray-400 text-sm">
                      No changelog available
                    </p>
                  )}
                </div>
              </div>

              {/* Error Display */}
              {error && (
                <div className="mb-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
                  {error}
                </div>
              )}

              {/* Actions */}
              <div className="flex space-x-3">
                <Button
                  variant="secondary"
                  onClick={() => setShowModal(false)}
                  className="flex-1"
                  disabled={isUpdating}
                >
                  Later
                </Button>
                {isAdmin ? (
                  <Button
                    onClick={handleApplyUpdate}
                    className="flex-1"
                    disabled={isUpdating}
                  >
                    {isUpdating ? 'Updating...' : 'Update Now'}
                  </Button>
                ) : (
                  <div className="flex-1 p-3 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg text-yellow-700 dark:text-yellow-400 text-sm text-center">
                    Admin permissions required to update
                  </div>
                )}
              </div>

              {isAdmin && (
                <p className="mt-4 text-xs text-gray-500 dark:text-gray-400 text-center">
                  Note: Server restart required after update
                </p>
              )}
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </>
  );
}
