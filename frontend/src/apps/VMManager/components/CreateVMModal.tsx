import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Server, AlertCircle } from 'lucide-react';
import { vmsApi, type VMCreateRequest } from '@/api/vms';
import { networkApi } from '@/api/network';
import { getErrorMessage } from '@/api/client';

interface CreateVMModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export function CreateVMModal({ isOpen, onClose, onSuccess }: CreateVMModalProps) {
  const [formData, setFormData] = useState<VMCreateRequest>({
    name: '',
    memory: 2048,
    vcpus: 2,
    disk_size: 20,
    os_type: 'linux',
    os_variant: 'ubuntu22.04',
    iso_path: '',
    network: 'default',
    autostart: false,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [bridges, setBridges] = useState<string[]>(['default', 'br0']); // Default options

  // Fetch available bridges when component mounts
  useEffect(() => {
    const fetchBridges = async () => {
      try {
        const response = await networkApi.listBridges();
        if (response.success && response.data && response.data.length > 0) {
          // Include 'default' and available bridges
          const availableBridges = ['default', ...response.data!];
          setBridges(availableBridges);
        }
      } catch (err) {
        // If fetching bridges fails, keep the default options
        console.error('Failed to fetch bridges:', err);
      }
    };

    if (isOpen) {
      fetchBridges();
    }
  }, [isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name.trim()) {
      setError('VM name is required');
      return;
    }

    try {
      setLoading(true);
      setError('');
      const response = await vmsApi.createVM(formData);

      if (response.success) {
        onSuccess();
        onClose();
        // Reset form
        setFormData({
          name: '',
          memory: 2048,
          vcpus: 2,
          disk_size: 20,
          os_type: 'linux',
          os_variant: 'ubuntu22.04',
          iso_path: '',
          network: 'default',
          autostart: false,
        });
      } else {
        setError(response.error?.message || 'Failed to create VM');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
        onClick={onClose}
      >
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
          className="bg-white dark:bg-macos-dark-100 rounded-xl shadow-2xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-auto"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
            <div className="flex items-center gap-3">
              <Server className="w-6 h-6 text-macos-blue" />
              <h2 className="text-xl font-bold text-gray-900 dark:text-white">Create Virtual Machine</h2>
            </div>
            <button
              onClick={onClose}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
            >
              <X className="w-5 h-5 text-gray-600 dark:text-gray-400" />
            </button>
          </div>

          {/* Error Display */}
          {error && (
            <div className="mx-6 mt-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg flex items-start gap-3">
              <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <h3 className="font-semibold text-red-900 dark:text-red-200">Error</h3>
                <p className="text-sm text-red-700 dark:text-red-300 mt-1">{error}</p>
              </div>
            </div>
          )}

          {/* Form */}
          <form onSubmit={handleSubmit} className="p-6 space-y-6">
            {/* General Section */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">General</h3>
              <div className="space-y-4">
                {/* VM Name */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    VM Name *
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="my-vm"
                    className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                    required
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  {/* OS Type */}
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      OS Type
                    </label>
                    <select
                      value={formData.os_type}
                      onChange={(e) => setFormData({ ...formData, os_type: e.target.value })}
                      className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                    >
                      <option value="linux">Linux</option>
                      <option value="windows">Windows</option>
                      <option value="unix">Unix</option>
                      <option value="other">Other</option>
                    </select>
                  </div>

                  {/* OS Variant */}
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      OS Variant
                    </label>
                    <input
                      type="text"
                      value={formData.os_variant}
                      onChange={(e) => setFormData({ ...formData, os_variant: e.target.value })}
                      placeholder="ubuntu22.04, win10"
                      className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                    />
                  </div>
                </div>

                {/* ISO Path */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    ISO Path (optional)
                  </label>
                  <input
                    type="text"
                    value={formData.iso_path}
                    onChange={(e) => setFormData({ ...formData, iso_path: e.target.value })}
                    placeholder="/path/to/iso/file.iso"
                    className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                  />
                </div>
              </div>
            </div>

            {/* Resources Section */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Resources</h3>
              <div className="space-y-4">
                {/* Memory */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Memory (MB)
                  </label>
                  <div className="flex items-center gap-4">
                    <input
                      type="range"
                      min="512"
                      max="32768"
                      step="512"
                      value={formData.memory}
                      onChange={(e) => setFormData({ ...formData, memory: parseInt(e.target.value) })}
                      className="flex-1"
                    />
                    <input
                      type="number"
                      value={formData.memory}
                      onChange={(e) => setFormData({ ...formData, memory: parseInt(e.target.value) || 2048 })}
                      className="w-24 px-3 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white"
                      min="512"
                      max="32768"
                      step="512"
                    />
                    <span className="text-sm text-gray-600 dark:text-gray-400">MB</span>
                  </div>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    {(formData.memory / 1024).toFixed(1)} GB - Recommended: 2048 MB minimum for modern OS
                  </p>
                </div>

                {/* vCPUs */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Virtual CPUs
                  </label>
                  <div className="flex items-center gap-4">
                    <input
                      type="range"
                      min="1"
                      max="16"
                      step="1"
                      value={formData.vcpus}
                      onChange={(e) => setFormData({ ...formData, vcpus: parseInt(e.target.value) })}
                      className="flex-1"
                    />
                    <input
                      type="number"
                      value={formData.vcpus}
                      onChange={(e) => setFormData({ ...formData, vcpus: parseInt(e.target.value) || 1 })}
                      className="w-24 px-3 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white"
                      min="1"
                      max="16"
                    />
                    <span className="text-sm text-gray-600 dark:text-gray-400">Cores</span>
                  </div>
                </div>

                {/* Disk Size */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Disk Size (GB)
                  </label>
                  <div className="flex items-center gap-4">
                    <input
                      type="range"
                      min="10"
                      max="500"
                      step="10"
                      value={formData.disk_size}
                      onChange={(e) => setFormData({ ...formData, disk_size: parseInt(e.target.value) })}
                      className="flex-1"
                    />
                    <input
                      type="number"
                      value={formData.disk_size}
                      onChange={(e) => setFormData({ ...formData, disk_size: parseInt(e.target.value) || 20 })}
                      className="w-24 px-3 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white"
                      min="10"
                      max="500"
                      step="5"
                    />
                    <span className="text-sm text-gray-600 dark:text-gray-400">GB</span>
                  </div>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    Recommended: 20 GB minimum for Linux, 40 GB+ for Windows
                  </p>
                </div>
              </div>
            </div>

            {/* Network Section */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Network</h3>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Network Bridge
                </label>
                <select
                  value={formData.network}
                  onChange={(e) => setFormData({ ...formData, network: e.target.value })}
                  className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                >
                  {bridges.map((bridge) => (
                    <option key={bridge} value={bridge}>
                      {bridge}
                    </option>
                  ))}
                </select>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Available bridges: {bridges.join(', ')}
                </p>
              </div>
            </div>

            {/* Options Section */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Options</h3>
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="autostart"
                  checked={formData.autostart}
                  onChange={(e) => setFormData({ ...formData, autostart: e.target.checked })}
                  className="w-4 h-4 text-macos-blue bg-white dark:bg-macos-dark-50 border-gray-300 dark:border-gray-600 rounded focus:ring-macos-blue"
                />
                <label htmlFor="autostart" className="ml-2 text-sm text-gray-700 dark:text-gray-300">
                  Start VM automatically on boot
                </label>
              </div>
            </div>

            {/* Buttons */}
            <div className="flex gap-3 pt-4">
              <button
                type="button"
                onClick={onClose}
                disabled={loading}
                className="flex-1 px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={loading}
                className="flex-1 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
              >
                {loading ? (
                  <>
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                    Creating...
                  </>
                ) : (
                  <>
                    <Server className="w-4 h-4" />
                    Create VM
                  </>
                )}
              </button>
            </div>
          </form>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
}
