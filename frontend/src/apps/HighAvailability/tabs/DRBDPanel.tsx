import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { haApi, DRBDStatus, CreateDRBDResourceRequest } from '@/api/ha';
import toast from 'react-hot-toast';

export default function DRBDPanel() {
  const [resources, setResources] = useState<string[]>([]);
  const [resourceStatuses, setResourceStatuses] = useState<Map<string, DRBDStatus>>(new Map());
  const [loading, setLoading] = useState(true);
  const [showCreateDialog, setShowCreateDialog] = useState(false);

  const loadResources = async () => {
    try {
      const response = await haApi.listDRBDResources();
      if (response.success && response.data) {
        setResources(response.data);

        // Load status for each resource
        const statuses = new Map<string, DRBDStatus>();
        for (const resourceName of response.data) {
          const statusResponse = await haApi.getDRBDResourceStatus(resourceName);
          if (statusResponse.success && statusResponse.data) {
            statuses.set(resourceName, statusResponse.data);
          }
        }
        setResourceStatuses(statuses);
      } else {
        toast.error(response.error?.message || 'Failed to load DRBD resources');
      }
    } catch (error) {
      toast.error('Failed to load DRBD resources');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadResources();
    const interval = setInterval(loadResources, 5000); // Refresh every 5 seconds
    return () => clearInterval(interval);
  }, []);

  const handleResourceAction = async (action: string, actionFn: () => Promise<any>) => {
    try {
      const response = await actionFn();
      if (response.success) {
        toast.success(`${action} successful`);
        await loadResources();
      } else {
        toast.error(response.error?.message || `${action} failed`);
      }
    } catch (error) {
      toast.error(`${action} failed`);
      console.error(error);
    }
  };

  const handleDeleteResource = async (resourceName: string) => {
    if (!confirm(`Are you sure you want to delete DRBD resource "${resourceName}"?`)) {
      return;
    }

    try {
      const response = await haApi.deleteDRBDResource(resourceName);
      if (response.success) {
        toast.success('Resource deleted successfully');
        await loadResources();
      } else {
        toast.error(response.error?.message || 'Failed to delete resource');
      }
    } catch (error) {
      toast.error('Failed to delete resource');
      console.error(error);
    }
  };

  const getConnectionStateColor = (state: string) => {
    switch (state.toLowerCase()) {
      case 'connected': return 'text-green-600 bg-green-100';
      case 'disconnected': return 'text-yellow-600 bg-yellow-100';
      case 'standalone': return 'text-red-600 bg-red-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getRoleColor = (role: string) => {
    switch (role.toLowerCase()) {
      case 'primary': return 'text-blue-600 bg-blue-100';
      case 'secondary': return 'text-purple-600 bg-purple-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getDiskStateColor = (state: string) => {
    switch (state.toLowerCase()) {
      case 'uptodate': return 'text-green-600 bg-green-100';
      case 'inconsistent': return 'text-yellow-600 bg-yellow-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue"></div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
            DRBD Resources
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            Manage distributed replicated block devices for high availability
          </p>
        </div>
        <button
          onClick={() => setShowCreateDialog(true)}
          className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2"
        >
          <span>+</span>
          Create Resource
        </button>
      </div>

      {/* Resources List */}
      {resources.length === 0 ? (
        <div className="text-center py-12 bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700">
          <div className="text-6xl mb-4">💿</div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            No DRBD Resources
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-4">
            Create your first DRBD resource to enable block-level replication
          </p>
          <button
            onClick={() => setShowCreateDialog(true)}
            className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            Create Resource
          </button>
        </div>
      ) : (
        <div className="grid gap-4">
          {resources.map((resourceName) => {
            const status = resourceStatuses.get(resourceName);
            return (
              <motion.div
                key={resourceName}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 p-6"
              >
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      {resourceName}
                    </h3>
                    {status && (
                      <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                        {status.device}
                      </p>
                    )}
                  </div>
                  <button
                    onClick={() => handleDeleteResource(resourceName)}
                    className="text-red-600 hover:text-red-700 px-3 py-1 rounded hover:bg-red-50 transition-colors"
                  >
                    Delete
                  </button>
                </div>

                {status ? (
                  <div className="space-y-4">
                    {/* Status Badges */}
                    <div className="flex flex-wrap gap-2">
                      <span className={`px-3 py-1 rounded-full text-sm font-medium ${getConnectionStateColor(status.connection_state)}`}>
                        {status.connection_state}
                      </span>
                      <span className={`px-3 py-1 rounded-full text-sm font-medium ${getRoleColor(status.role)}`}>
                        {status.role}
                      </span>
                      <span className={`px-3 py-1 rounded-full text-sm font-medium ${getDiskStateColor(status.disk_state)}`}>
                        {status.disk_state}
                      </span>
                    </div>

                    {/* Peer Info */}
                    {status.peer_role && status.peer_role !== 'Unknown' && (
                      <div className="grid grid-cols-2 gap-4 p-4 bg-gray-50 dark:bg-macos-dark-50 rounded-lg">
                        <div>
                          <p className="text-sm text-gray-600 dark:text-gray-400">Peer Role</p>
                          <p className="text-sm font-medium text-gray-900 dark:text-gray-100">{status.peer_role}</p>
                        </div>
                        <div>
                          <p className="text-sm text-gray-600 dark:text-gray-400">Peer Disk State</p>
                          <p className="text-sm font-medium text-gray-900 dark:text-gray-100">{status.peer_disk_state}</p>
                        </div>
                      </div>
                    )}

                    {/* Sync Progress */}
                    {status.resyncing && (
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <p className="text-sm text-gray-600 dark:text-gray-400">Synchronization Progress</p>
                          <p className="text-sm font-medium text-gray-900 dark:text-gray-100">{status.sync_progress}%</p>
                        </div>
                        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                          <div
                            className="bg-macos-blue h-2 rounded-full transition-all duration-300"
                            style={{ width: `${status.sync_progress}%` }}
                          />
                        </div>
                      </div>
                    )}

                    {/* Actions */}
                    <div className="flex flex-wrap gap-2 pt-2">
                      {status.role === 'Secondary' && (
                        <button
                          onClick={() => handleResourceAction('Promote to Primary', () => haApi.promoteDRBDResource(resourceName))}
                          className="px-3 py-1.5 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors text-sm"
                        >
                          Promote to Primary
                        </button>
                      )}
                      {status.role === 'Primary' && (
                        <button
                          onClick={() => handleResourceAction('Demote to Secondary', () => haApi.demoteDRBDResource(resourceName))}
                          className="px-3 py-1.5 bg-purple-100 text-purple-700 rounded hover:bg-purple-200 transition-colors text-sm"
                        >
                          Demote to Secondary
                        </button>
                      )}
                      {status.connection_state === 'Disconnected' && (
                        <button
                          onClick={() => handleResourceAction('Connect', () => haApi.connectDRBDResource(resourceName))}
                          className="px-3 py-1.5 bg-green-100 text-green-700 rounded hover:bg-green-200 transition-colors text-sm"
                        >
                          Connect
                        </button>
                      )}
                      {status.connection_state === 'Connected' && (
                        <button
                          onClick={() => handleResourceAction('Disconnect', () => haApi.disconnectDRBDResource(resourceName))}
                          className="px-3 py-1.5 bg-yellow-100 text-yellow-700 rounded hover:bg-yellow-200 transition-colors text-sm"
                        >
                          Disconnect
                        </button>
                      )}
                      <button
                        onClick={() => handleResourceAction('Start Sync', () => haApi.startDRBDSync(resourceName))}
                        className="px-3 py-1.5 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors text-sm"
                      >
                        Start Sync
                      </button>
                      <button
                        onClick={() => handleResourceAction('Verify Data', () => haApi.verifyDRBDData(resourceName))}
                        className="px-3 py-1.5 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors text-sm"
                      >
                        Verify Data
                      </button>
                      {(status.connection_state === 'StandAlone' || status.connection_state === 'Disconnected') && (
                        <button
                          onClick={() => {
                            if (confirm('Force primary is a dangerous operation that should only be used for split-brain recovery. Continue?')) {
                              handleResourceAction('Force Primary', () => haApi.forcePrimaryDRBDResource(resourceName));
                            }
                          }}
                          className="px-3 py-1.5 bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors text-sm"
                        >
                          Force Primary
                        </button>
                      )}
                    </div>
                  </div>
                ) : (
                  <div className="text-sm text-gray-500 dark:text-gray-400">
                    Loading status...
                  </div>
                )}
              </motion.div>
            );
          })}
        </div>
      )}

      {/* Create Resource Dialog */}
      <AnimatePresence>
        {showCreateDialog && (
          <CreateResourceDialog
            onClose={() => setShowCreateDialog(false)}
            onCreated={() => {
              setShowCreateDialog(false);
              loadResources();
            }}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

function CreateResourceDialog({ onClose, onCreated }: { onClose: () => void; onCreated: () => void }) {
  const [formData, setFormData] = useState<CreateDRBDResourceRequest>({
    name: '',
    device: '/dev/drbd0',
    disk: '',
    meta_disk: 'internal',
    local_address: '',
    peer_address: '',
    protocol: 'C',
  });
  const [creating, setCreating] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);

    try {
      const response = await haApi.createDRBDResource(formData);
      if (response.success) {
        toast.success('DRBD resource created successfully');
        onCreated();
      } else {
        toast.error(response.error?.message || 'Failed to create resource');
      }
    } catch (error) {
      toast.error('Failed to create resource');
      console.error(error);
    } finally {
      setCreating(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
      >
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
            Create DRBD Resource
          </h3>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Resource Name
            </label>
            <input
              type="text"
              required
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
              placeholder="r0"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Device
              </label>
              <input
                type="text"
                required
                value={formData.device}
                onChange={(e) => setFormData({ ...formData, device: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                placeholder="/dev/drbd0"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Disk
              </label>
              <input
                type="text"
                required
                value={formData.disk}
                onChange={(e) => setFormData({ ...formData, disk: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                placeholder="/dev/sda1"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Meta Disk
            </label>
            <input
              type="text"
              value={formData.meta_disk}
              onChange={(e) => setFormData({ ...formData, meta_disk: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
              placeholder="internal"
            />
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Use "internal" for internal metadata, or specify a device path
            </p>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Local Address
              </label>
              <input
                type="text"
                required
                value={formData.local_address}
                onChange={(e) => setFormData({ ...formData, local_address: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                placeholder="192.168.1.10:7788"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Peer Address
              </label>
              <input
                type="text"
                required
                value={formData.peer_address}
                onChange={(e) => setFormData({ ...formData, peer_address: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                placeholder="192.168.1.11:7788"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Protocol
            </label>
            <select
              value={formData.protocol}
              onChange={(e) => setFormData({ ...formData, protocol: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
            >
              <option value="A">Protocol A (Asynchronous)</option>
              <option value="B">Protocol B (Semi-synchronous)</option>
              <option value="C">Protocol C (Synchronous)</option>
            </select>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Protocol C is recommended for high availability (synchronous replication)
            </p>
          </div>

          <div className="flex justify-end gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
              disabled={creating}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={creating}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
            >
              {creating ? 'Creating...' : 'Create Resource'}
            </button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
}
