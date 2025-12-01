// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  Share2,
  Plus,
  RefreshCw,
  Trash2,
  Folder,
  Globe,
  Settings,
  X,
  AlertCircle,
} from 'lucide-react';
import { syslibApi, type NFSExport, type CreateNFSExportRequest } from '@/api/syslib';

export default function NFSManager() {
  const [exports, setExports] = useState<NFSExport[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [selectedExport, setSelectedExport] = useState<NFSExport | null>(null);

  const [formData, setFormData] = useState<CreateNFSExportRequest>({
    path: '',
    clients: [],
    options: ['rw', 'sync', 'no_subtree_check'],
  });

  const [clientInput, setClientInput] = useState('');

  // Common NFS options
  const commonOptions = [
    { value: 'rw', label: 'Read/Write', description: 'Allow read and write access' },
    { value: 'ro', label: 'Read-Only', description: 'Allow only read access' },
    { value: 'sync', label: 'Sync', description: 'Write changes before responding' },
    { value: 'async', label: 'Async', description: 'Write changes asynchronously' },
    { value: 'no_subtree_check', label: 'No Subtree Check', description: 'Disable subtree checking (faster)' },
    { value: 'no_root_squash', label: 'No Root Squash', description: 'Allow root access from clients' },
    { value: 'root_squash', label: 'Root Squash', description: 'Map root user to anonymous' },
    { value: 'all_squash', label: 'All Squash', description: 'Map all users to anonymous' },
  ];

  // Fetch NFS exports
  const fetchExports = async () => {
    setIsLoading(true);
    try {
      const response = await syslibApi.nfs.listExports();
      if (response.success && response.data) {
        setExports(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch NFS exports:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchExports();
  }, []);

  const handleCreateExport = async () => {
    if (!formData.path || formData.clients.length === 0) {
      alert('Please provide export path and at least one client');
      return;
    }

    try {
      const response = await syslibApi.nfs.createExport(formData);
      if (response.success) {
        alert('NFS export created successfully');
        setShowCreateDialog(false);
        setFormData({ path: '', clients: [], options: ['rw', 'sync', 'no_subtree_check'] });
        setClientInput('');
        fetchExports();
      }
    } catch (error) {
      console.error('Failed to create NFS export:', error);
      alert('Failed to create NFS export');
    }
  };

  const handleDeleteExport = async () => {
    if (!selectedExport) return;

    try {
      const response = await syslibApi.nfs.deleteExport(selectedExport.path);
      if (response.success) {
        alert('NFS export deleted successfully');
        setShowDeleteDialog(false);
        setSelectedExport(null);
        fetchExports();
      }
    } catch (error) {
      console.error('Failed to delete NFS export:', error);
      alert('Failed to delete NFS export');
    }
  };

  const handleRestartNFS = async () => {
    try {
      const response = await syslibApi.nfs.restart();
      if (response.success) {
        alert('NFS service restarted successfully');
      }
    } catch (error) {
      console.error('Failed to restart NFS:', error);
      alert('Failed to restart NFS service');
    }
  };

  const addClient = () => {
    if (clientInput && !formData.clients.includes(clientInput)) {
      setFormData({
        ...formData,
        clients: [...formData.clients, clientInput],
      });
      setClientInput('');
    }
  };

  const removeClient = (client: string) => {
    setFormData({
      ...formData,
      clients: formData.clients.filter((c) => c !== client),
    });
  };

  const toggleOption = (option: string) => {
    const currentOptions = formData.options || [];

    // Handle mutually exclusive options
    let newOptions: string[];
    if (option === 'rw' || option === 'ro') {
      newOptions = currentOptions.filter((o) => o !== 'rw' && o !== 'ro');
      newOptions.push(option);
    } else if (option === 'sync' || option === 'async') {
      newOptions = currentOptions.filter((o) => o !== 'sync' && o !== 'async');
      newOptions.push(option);
    } else if (option === 'root_squash' || option === 'no_root_squash' || option === 'all_squash') {
      newOptions = currentOptions.filter((o) => !['root_squash', 'no_root_squash', 'all_squash'].includes(o));
      newOptions.push(option);
    } else {
      if (currentOptions.includes(option)) {
        newOptions = currentOptions.filter((o) => o !== option);
      } else {
        newOptions = [...currentOptions, option];
      }
    }

    setFormData({ ...formData, options: newOptions });
  };

  return (
    <div className="flex flex-col h-full bg-gradient-to-br from-gray-50 to-white dark:from-macos-dark-100 dark:to-macos-dark-200">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200/50 dark:border-gray-700/50 bg-white/50 dark:bg-macos-dark-100/50 backdrop-blur-sm">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl shadow-lg">
            <Share2 className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              NFS Export Manager
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Network file system sharing
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <motion.button
            whileHover={{ scale: 1.05, y: -2 }}
            whileTap={{ scale: 0.95 }}
            onClick={handleRestartNFS}
            className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-xl hover:shadow-lg transition-all border border-gray-200 dark:border-gray-700"
          >
            <Settings className="w-4 h-4" />
            Restart NFS
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.05, y: -2 }}
            whileTap={{ scale: 0.95 }}
            onClick={fetchExports}
            className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-xl hover:shadow-lg transition-all border border-gray-200 dark:border-gray-700"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.05, y: -2 }}
            whileTap={{ scale: 0.95 }}
            onClick={() => setShowCreateDialog(true)}
            className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-macos-blue to-macos-purple text-white rounded-xl hover:shadow-lg transition-all"
          >
            <Plus className="w-4 h-4" />
            Create Export
          </motion.button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
          </div>
        ) : exports.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-64 text-center">
            <Share2 className="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" />
            <p className="text-gray-500 dark:text-gray-400">No NFS exports configured</p>
            <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
              Create an export to share directories via NFS
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {exports.map((exp, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.05 }}
                whileHover={{ y: -4, scale: 1.02 }}
                className="bg-gradient-to-br from-white to-gray-50 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-2xl p-6 border border-gray-200 dark:border-gray-700 shadow-md hover:shadow-xl transition-all"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div className="p-3 bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 rounded-xl">
                      <Folder className="w-6 h-6 text-macos-blue" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                        {exp.path}
                      </h3>
                      <p className="text-xs px-2 py-1 bg-macos-blue/10 text-macos-blue rounded-full inline-block mt-1">
                        NFS Export
                      </p>
                    </div>
                  </div>
                  <motion.button
                    whileHover={{ scale: 1.1, rotate: 10 }}
                    whileTap={{ scale: 0.9 }}
                    onClick={() => {
                      setSelectedExport(exp);
                      setShowDeleteDialog(true);
                    }}
                    className="p-2 hover:bg-red-100 dark:hover:bg-red-900/20 rounded-xl transition-colors"
                  >
                    <Trash2 className="w-4 h-4 text-red-500" />
                  </motion.button>
                </div>

                {/* Clients */}
                <div className="mb-4">
                  <div className="flex items-center gap-2 mb-2">
                    <div className="p-1.5 bg-green-500/10 rounded-lg">
                      <Globe className="w-4 h-4 text-green-600 dark:text-green-400" />
                    </div>
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                      Allowed Clients ({exp.clients.length})
                    </span>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {exp.clients.map((client, idx) => (
                      <motion.span
                        key={idx}
                        initial={{ opacity: 0, scale: 0.8 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ delay: 0.1 + idx * 0.05 }}
                        className="px-3 py-1 bg-gradient-to-r from-blue-100 to-blue-50 dark:from-blue-900/20 dark:to-blue-800/20 text-blue-700 dark:text-blue-300 text-xs rounded-full font-medium border border-blue-200 dark:border-blue-700"
                      >
                        {client}
                      </motion.span>
                    ))}
                  </div>
                </div>

                {/* Options */}
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    <div className="p-1.5 bg-purple-500/10 rounded-lg">
                      <Settings className="w-4 h-4 text-purple-600 dark:text-purple-400" />
                    </div>
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                      Export Options
                    </span>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {exp.options.map((option, idx) => (
                      <motion.span
                        key={idx}
                        initial={{ opacity: 0, scale: 0.8 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ delay: 0.2 + idx * 0.05 }}
                        className="px-2 py-1 bg-gray-100 dark:bg-macos-dark-100 text-gray-700 dark:text-gray-300 text-xs font-mono rounded border border-gray-300 dark:border-gray-600"
                      >
                        {option}
                      </motion.span>
                    ))}
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </div>

      {/* Create Export Dialog */}
      {showCreateDialog && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-2xl w-full m-4 max-h-[80vh] overflow-y-auto shadow-2xl border border-gray-200 dark:border-gray-700"
          >
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-gradient-to-br from-macos-blue to-macos-purple rounded-lg">
                  <Share2 className="w-5 h-5 text-white" />
                </div>
                <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                  Create NFS Export
                </h3>
              </div>
              <motion.button
                whileHover={{ scale: 1.1, rotate: 90 }}
                whileTap={{ scale: 0.9 }}
                onClick={() => setShowCreateDialog(false)}
                className="p-1 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded-lg transition-colors"
              >
                <X className="w-5 h-5" />
              </motion.button>
            </div>

            <div className="space-y-6">
              {/* Export Path */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Export Path
                </label>
                <input
                  type="text"
                  value={formData.path}
                  onChange={(e) => setFormData({ ...formData, path: e.target.value })}
                  placeholder="/srv/nfs/share"
                  className="w-full px-4 py-2.5 bg-gray-50 dark:bg-macos-dark-200 border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-macos-blue focus:border-macos-blue transition-all font-mono text-sm"
                />
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Absolute path to the directory to export
                </p>
              </div>

              {/* Clients */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Allowed Clients
                </label>
                <div className="flex gap-2 mb-3">
                  <input
                    type="text"
                    value={clientInput}
                    onChange={(e) => setClientInput(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && addClient()}
                    placeholder="192.168.1.0/24 or hostname"
                    className="flex-1 px-4 py-2.5 bg-gray-50 dark:bg-macos-dark-200 border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-macos-blue focus:border-macos-blue transition-all font-mono text-sm"
                  />
                  <motion.button
                    whileHover={{ scale: 1.05, y: -2 }}
                    whileTap={{ scale: 0.95 }}
                    onClick={addClient}
                    className="px-5 py-2.5 bg-gradient-to-r from-macos-blue to-macos-purple text-white rounded-xl hover:shadow-lg transition-all font-medium"
                  >
                    Add
                  </motion.button>
                </div>
                <div className="flex flex-wrap gap-2">
                  {formData.clients.map((client, idx) => (
                    <motion.span
                      key={idx}
                      initial={{ opacity: 0, scale: 0.8 }}
                      animate={{ opacity: 1, scale: 1 }}
                      className="inline-flex items-center gap-2 px-3 py-1.5 bg-gradient-to-r from-blue-100 to-blue-50 dark:from-blue-900/20 dark:to-blue-800/20 text-blue-700 dark:text-blue-300 text-sm rounded-full border border-blue-200 dark:border-blue-700"
                    >
                      {client}
                      <motion.button
                        whileHover={{ scale: 1.2, rotate: 90 }}
                        whileTap={{ scale: 0.8 }}
                        onClick={() => removeClient(client)}
                        className="hover:text-red-500 transition-colors"
                      >
                        <X className="w-3 h-3" />
                      </motion.button>
                    </motion.span>
                  ))}
                </div>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
                  Use * for all hosts, IP addresses, CIDR notation, or hostnames
                </p>
              </div>

              {/* Options */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                  Export Options
                </label>
                <div className="grid grid-cols-2 gap-3">
                  {commonOptions.map((opt, index) => {
                    const isSelected = formData.options?.includes(opt.value);
                    return (
                      <motion.button
                        key={opt.value}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.05 }}
                        whileHover={{ scale: 1.02, y: -2 }}
                        whileTap={{ scale: 0.98 }}
                        onClick={() => toggleOption(opt.value)}
                        className={`p-3 rounded-xl border-2 transition-all text-left ${
                          isSelected
                            ? 'border-macos-blue bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 shadow-md'
                            : 'border-gray-300 dark:border-gray-600 hover:border-macos-blue/50 bg-gray-50 dark:bg-macos-dark-200'
                        }`}
                      >
                        <div className="font-medium text-gray-900 dark:text-gray-100 text-sm">
                          {opt.label}
                        </div>
                        <div className="text-xs text-gray-600 dark:text-gray-400">
                          {opt.description}
                        </div>
                      </motion.button>
                    );
                  })}
                </div>
              </div>

              {/* Info */}
              <div className="bg-gradient-to-br from-blue-50 to-blue-100/50 dark:from-blue-900/20 dark:to-blue-800/20 rounded-xl p-4 border border-blue-200 dark:border-blue-700">
                <div className="flex gap-3">
                  <div className="p-2 bg-blue-500/10 rounded-lg">
                    <AlertCircle className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  </div>
                  <div className="text-sm text-gray-700 dark:text-gray-300">
                    <p className="font-semibold mb-2">Important Notes:</p>
                    <ul className="list-disc list-inside space-y-1 text-xs">
                      <li>Ensure the directory exists and has proper permissions</li>
                      <li>Use <code className="bg-white dark:bg-macos-dark-100 px-2 py-0.5 rounded font-mono">no_root_squash</code> carefully - it allows root access</li>
                      <li>The export will be added to /etc/exports</li>
                      <li>Changes take effect after NFS restart</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-6 flex justify-end gap-2">
              <motion.button
                whileHover={{ scale: 1.05, y: -2 }}
                whileTap={{ scale: 0.95 }}
                onClick={() => setShowCreateDialog(false)}
                className="px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-xl hover:shadow-lg transition-all border border-gray-200 dark:border-gray-700"
              >
                Cancel
              </motion.button>
              <motion.button
                whileHover={{ scale: 1.05, y: -2 }}
                whileTap={{ scale: 0.95 }}
                onClick={handleCreateExport}
                className="px-4 py-2 bg-gradient-to-r from-macos-blue to-macos-purple text-white rounded-xl hover:shadow-lg transition-all"
              >
                Create Export
              </motion.button>
            </div>
          </motion.div>
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      {showDeleteDialog && selectedExport && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-md w-full m-4 shadow-2xl border border-gray-200 dark:border-gray-700"
          >
            <div className="flex items-center gap-3 mb-4">
              <div className="p-2 bg-gradient-to-br from-red-500 to-rose-600 rounded-lg">
                <AlertCircle className="w-5 h-5 text-white" />
              </div>
              <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Delete NFS Export
              </h3>
            </div>
            <p className="text-gray-600 dark:text-gray-400 mb-6">
              Are you sure you want to delete the export for{' '}
              <span className="font-semibold px-2 py-1 bg-gray-100 dark:bg-macos-dark-200 rounded font-mono text-gray-900 dark:text-gray-100">
                {selectedExport.path}
              </span>
              ? This action cannot be undone.
            </p>
            <div className="flex justify-end gap-2">
              <motion.button
                whileHover={{ scale: 1.05, y: -2 }}
                whileTap={{ scale: 0.95 }}
                onClick={() => {
                  setShowDeleteDialog(false);
                  setSelectedExport(null);
                }}
                className="px-4 py-2 bg-white dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-xl hover:shadow-lg transition-all border border-gray-200 dark:border-gray-700"
              >
                Cancel
              </motion.button>
              <motion.button
                whileHover={{ scale: 1.05, y: -2 }}
                whileTap={{ scale: 0.95 }}
                onClick={handleDeleteExport}
                className="px-4 py-2 bg-gradient-to-r from-red-500 to-rose-600 text-white rounded-xl hover:shadow-lg transition-all"
              >
                Delete
              </motion.button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}
