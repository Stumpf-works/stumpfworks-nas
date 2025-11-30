import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Server, AlertCircle } from 'lucide-react';
import { vmsApi, type VMCreateRequest } from '@/api/vms';
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
          <form onSubmit={handleSubmit} className="p-6 space-y-4">
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

            {/* Memory */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Memory (MB)
              </label>
              <input
                type="number"
                value={formData.memory}
                onChange={(e) => setFormData({ ...formData, memory: parseInt(e.target.value) || 0 })}
                min="512"
                step="512"
                className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
              />
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {(formData.memory / 1024).toFixed(1)} GB
              </p>
            </div>

            {/* vCPUs */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Virtual CPUs
              </label>
              <input
                type="number"
                value={formData.vcpus}
                onChange={(e) => setFormData({ ...formData, vcpus: parseInt(e.target.value) || 1 })}
                min="1"
                max="32"
                className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
              />
            </div>

            {/* Disk Size */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Disk Size (GB)
              </label>
              <input
                type="number"
                value={formData.disk_size}
                onChange={(e) => setFormData({ ...formData, disk_size: parseInt(e.target.value) || 0 })}
                min="5"
                step="5"
                className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
              />
            </div>

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
                placeholder="ubuntu22.04, debian11, win10, etc."
                className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
              />
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

            {/* Network */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Network
              </label>
              <input
                type="text"
                value={formData.network}
                onChange={(e) => setFormData({ ...formData, network: e.target.value })}
                placeholder="default"
                className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
              />
            </div>

            {/* Autostart */}
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
