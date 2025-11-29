import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { haApi, VIPStatus, VIPConfig } from '@/api/ha';
import toast from 'react-hot-toast';

export default function VIPPanel() {
  const [vips, setVips] = useState<VIPStatus[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateDialog, setShowCreateDialog] = useState(false);

  const loadVIPs = async () => {
    try {
      const response = await haApi.listVIPs();
      if (response.success && response.data) {
        setVips(response.data);
      } else {
        toast.error(response.error?.message || 'Failed to load VIPs');
      }
    } catch (error) {
      toast.error('Failed to load VIPs');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadVIPs();
    const interval = setInterval(loadVIPs, 5000); // Refresh every 5 seconds
    return () => clearInterval(interval);
  }, []);

  const handleVIPAction = async (action: string, actionFn: () => Promise<any>) => {
    try {
      const response = await actionFn();
      if (response.success) {
        toast.success(`${action} successful`);
        await loadVIPs();
      } else {
        toast.error(response.error?.message || `${action} failed`);
      }
    } catch (error) {
      toast.error(`${action} failed`);
      console.error(error);
    }
  };

  const handleDeleteVIP = async (vipId: string) => {
    if (!confirm(`Are you sure you want to delete VIP "${vipId}"?`)) {
      return;
    }

    try {
      const response = await haApi.deleteVIP(vipId);
      if (response.success) {
        toast.success('VIP deleted successfully');
        await loadVIPs();
      } else {
        toast.error(response.error?.message || 'Failed to delete VIP');
      }
    } catch (error) {
      toast.error('Failed to delete VIP');
      console.error(error);
    }
  };

  const getStateColor = (state: string) => {
    switch (state.toUpperCase()) {
      case 'MASTER': return 'text-green-600 bg-green-100';
      case 'BACKUP': return 'text-blue-600 bg-blue-100';
      case 'FAULT': return 'text-red-600 bg-red-100';
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
            Virtual IPs (Keepalived)
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            Manage virtual IP addresses for high availability
          </p>
        </div>
        <button
          onClick={() => setShowCreateDialog(true)}
          className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2"
        >
          <span>+</span>
          Create VIP
        </button>
      </div>

      {/* VIPs List */}
      {vips.length === 0 ? (
        <div className="text-center py-12 bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700">
          <div className="text-6xl mb-4">🌐</div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            No Virtual IPs
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-4">
            Create your first Virtual IP for automatic failover
          </p>
          <button
            onClick={() => setShowCreateDialog(true)}
            className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            Create VIP
          </button>
        </div>
      ) : (
        <div className="grid gap-4">
          {vips.map((vip) => (
            <motion.div
              key={vip.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-white dark:bg-macos-dark-100 rounded-lg border border-gray-200 dark:border-gray-700 p-6"
            >
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    {vip.virtual_ip}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                    Interface: {vip.interface} • Priority: {vip.priority}
                  </p>
                </div>
                <button
                  onClick={() => handleDeleteVIP(vip.id)}
                  className="text-red-600 hover:text-red-700 px-3 py-1 rounded hover:bg-red-50 transition-colors"
                >
                  Delete
                </button>
              </div>

              {/* Status Badges */}
              <div className="flex flex-wrap gap-2 mb-4">
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStateColor(vip.state)}`}>
                  {vip.state}
                </span>
                {vip.is_active && (
                  <span className="px-3 py-1 rounded-full text-sm font-medium text-green-600 bg-green-100">
                    VIP Active
                  </span>
                )}
                {vip.is_master && (
                  <span className="px-3 py-1 rounded-full text-sm font-medium text-purple-600 bg-purple-100">
                    This Node is Master
                  </span>
                )}
              </div>

              {/* Actions */}
              <div className="flex flex-wrap gap-2">
                {!vip.is_master && (
                  <button
                    onClick={() => handleVIPAction('Promote to MASTER', () => haApi.promoteVIPToMaster(vip.id))}
                    className="px-3 py-1.5 bg-green-100 text-green-700 rounded hover:bg-green-200 transition-colors text-sm"
                  >
                    Promote to MASTER
                  </button>
                )}
                {vip.is_master && (
                  <button
                    onClick={() => handleVIPAction('Demote to BACKUP', () => haApi.demoteVIPToBackup(vip.id))}
                    className="px-3 py-1.5 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors text-sm"
                  >
                    Demote to BACKUP
                  </button>
                )}
              </div>
            </motion.div>
          ))}
        </div>
      )}

      {/* Create VIP Dialog */}
      <AnimatePresence>
        {showCreateDialog && (
          <CreateVIPDialog
            onClose={() => setShowCreateDialog(false)}
            onCreated={() => {
              setShowCreateDialog(false);
              loadVIPs();
            }}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

function CreateVIPDialog({ onClose, onCreated }: { onClose: () => void; onCreated: () => void }) {
  const [formData, setFormData] = useState<VIPConfig>({
    virtual_ip: '',
    interface: 'eth0',
    router_id: 51,
    priority: 100,
    state: 'BACKUP',
    auth_pass: 'StumpfWorks',
  });
  const [creating, setCreating] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);

    try {
      const response = await haApi.createVIP(formData);
      if (response.success) {
        toast.success('VIP created successfully');
        onCreated();
      } else {
        toast.error(response.error?.message || 'Failed to create VIP');
      }
    } catch (error) {
      toast.error('Failed to create VIP');
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
            Create Virtual IP
          </h3>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Virtual IP Address *
            </label>
            <input
              type="text"
              required
              value={formData.virtual_ip}
              onChange={(e) => setFormData({ ...formData, virtual_ip: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
              placeholder="192.168.1.100"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Network Interface *
              </label>
              <input
                type="text"
                required
                value={formData.interface}
                onChange={(e) => setFormData({ ...formData, interface: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                placeholder="eth0"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Priority
              </label>
              <input
                type="number"
                min="1"
                max="255"
                value={formData.priority}
                onChange={(e) => setFormData({ ...formData, priority: parseInt(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
              />
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Higher priority becomes MASTER (1-255)
              </p>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Router ID
              </label>
              <input
                type="number"
                min="1"
                max="255"
                value={formData.router_id}
                onChange={(e) => setFormData({ ...formData, router_id: parseInt(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Initial State
              </label>
              <select
                value={formData.state}
                onChange={(e) => setFormData({ ...formData, state: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
              >
                <option value="MASTER">MASTER</option>
                <option value="BACKUP">BACKUP</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Authentication Password
            </label>
            <input
              type="password"
              value={formData.auth_pass}
              onChange={(e) => setFormData({ ...formData, auth_pass: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
            />
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Used for VRRP authentication between nodes
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
              {creating ? 'Creating...' : 'Create VIP'}
            </button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
}
