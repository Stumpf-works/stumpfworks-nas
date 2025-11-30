import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Box, AlertCircle, Network } from 'lucide-react';
import { lxcApi, type ContainerCreateRequest } from '@/api/lxc';
import { getErrorMessage } from '@/api/client';

interface CreateContainerModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export function CreateContainerModal({ isOpen, onClose, onSuccess }: CreateContainerModalProps) {
  const [formData, setFormData] = useState<ContainerCreateRequest>({
    name: '',
    template: 'ubuntu',
    release: 'jammy',
    architecture: 'amd64',
    memory_limit: 512,
    cpu_limit: 1,
    autostart: false,
    network_mode: 'internal',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name.trim()) {
      setError('Container name is required');
      return;
    }

    try {
      setLoading(true);
      setError('');
      const response = await lxcApi.createContainer(formData);

      if (response.success) {
        onSuccess();
        onClose();
        // Reset form
        setFormData({
          name: '',
          template: 'ubuntu',
          release: 'jammy',
          architecture: 'amd64',
          memory_limit: 512,
          cpu_limit: 1,
          autostart: false,
          network_mode: 'internal',
        });
      } else {
        setError(response.error?.message || 'Failed to create container');
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
              <Box className="w-6 h-6 text-macos-blue" />
              <h2 className="text-xl font-bold text-gray-900 dark:text-white">Create LXC Container</h2>
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
                {/* Container Name */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Container Name *
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="my-container"
                    className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                    required
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  {/* Template */}
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Template
                    </label>
                    <select
                      value={formData.template}
                      onChange={(e) => setFormData({ ...formData, template: e.target.value })}
                      className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                    >
                      <option value="ubuntu">Ubuntu</option>
                      <option value="debian">Debian</option>
                      <option value="alpine">Alpine</option>
                      <option value="centos">CentOS</option>
                      <option value="fedora">Fedora</option>
                    </select>
                  </div>

                  {/* Release */}
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Release
                    </label>
                    <input
                      type="text"
                      value={formData.release}
                      onChange={(e) => setFormData({ ...formData, release: e.target.value })}
                      placeholder="jammy, bookworm"
                      className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                    />
                  </div>
                </div>

                {/* Architecture */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Architecture
                  </label>
                  <select
                    value={formData.architecture}
                    onChange={(e) => setFormData({ ...formData, architecture: e.target.value })}
                    className="w-full px-4 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-white"
                  >
                    <option value="amd64">amd64 (x86_64)</option>
                    <option value="arm64">arm64 (aarch64)</option>
                    <option value="armhf">armhf</option>
                  </select>
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
                      min="128"
                      max="16384"
                      step="128"
                      value={formData.memory_limit}
                      onChange={(e) => setFormData({ ...formData, memory_limit: parseInt(e.target.value) })}
                      className="flex-1"
                    />
                    <input
                      type="number"
                      value={formData.memory_limit}
                      onChange={(e) => setFormData({ ...formData, memory_limit: parseInt(e.target.value) || 512 })}
                      className="w-24 px-3 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white"
                      min="128"
                      max="16384"
                    />
                    <span className="text-sm text-gray-600 dark:text-gray-400">MB</span>
                  </div>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    Recommended: 512 MB minimum, 2048 MB for typical workloads
                  </p>
                </div>

                {/* CPU */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    CPU Cores
                  </label>
                  <div className="flex items-center gap-4">
                    <input
                      type="range"
                      min="1"
                      max="16"
                      step="1"
                      value={formData.cpu_limit}
                      onChange={(e) => setFormData({ ...formData, cpu_limit: parseInt(e.target.value) })}
                      className="flex-1"
                    />
                    <input
                      type="number"
                      value={formData.cpu_limit}
                      onChange={(e) => setFormData({ ...formData, cpu_limit: parseInt(e.target.value) || 1 })}
                      className="w-24 px-3 py-2 bg-white dark:bg-macos-dark-50 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white"
                      min="1"
                      max="16"
                    />
                    <span className="text-sm text-gray-600 dark:text-gray-400">Cores</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Network Section */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                <Network className="w-5 h-5 text-macos-blue" />
                Network
              </h3>
              <div className="space-y-3">
                <label className="flex items-start p-3 border border-gray-300 dark:border-gray-600 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-macos-dark-50 transition-colors">
                  <input
                    type="radio"
                    name="network_mode"
                    value="internal"
                    checked={formData.network_mode === 'internal'}
                    onChange={(e) => setFormData({ ...formData, network_mode: e.target.value })}
                    className="w-4 h-4 text-macos-blue bg-white dark:bg-macos-dark-50 border-gray-300 dark:border-gray-600 mt-0.5"
                  />
                  <div className="ml-3 flex-1">
                    <div className="font-medium text-gray-900 dark:text-white">Internal Network (lxcbr0)</div>
                    <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      Container gets an IP from the internal network (10.0.3.x)
                    </div>
                  </div>
                </label>
                <label className="flex items-start p-3 border border-gray-300 dark:border-gray-600 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-macos-dark-50 transition-colors">
                  <input
                    type="radio"
                    name="network_mode"
                    value="bridged"
                    checked={formData.network_mode === 'bridged'}
                    onChange={(e) => setFormData({ ...formData, network_mode: e.target.value })}
                    className="w-4 h-4 text-macos-blue bg-white dark:bg-macos-dark-50 border-gray-300 dark:border-gray-600 mt-0.5"
                  />
                  <div className="ml-3 flex-1">
                    <div className="font-medium text-gray-900 dark:text-white">Bridged Network (br0)</div>
                    <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      Container gets an IP via DHCP from your router (192.168.178.x)
                    </div>
                  </div>
                </label>
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
                  Start container automatically on boot
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
                    <Box className="w-4 h-4" />
                    Create Container
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
