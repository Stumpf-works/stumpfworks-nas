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
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-100">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center gap-3">
          <Share2 className="w-6 h-6 text-macos-blue" />
          <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
            NFS Export Manager
          </h1>
        </div>
        <div className="flex gap-2">
          <button
            onClick={handleRestartNFS}
            className="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
          >
            <Settings className="w-4 h-4" />
            Restart NFS
          </button>
          <button
            onClick={fetchExports}
            className="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </button>
          <button
            onClick={() => setShowCreateDialog(true)}
            className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            <Plus className="w-4 h-4" />
            Create Export
          </button>
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
                className="bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl p-6 border border-gray-200 dark:border-gray-700"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div className="p-2 bg-macos-blue/10 rounded-lg">
                      <Folder className="w-6 h-6 text-macos-blue" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                        {exp.path}
                      </h3>
                      <p className="text-xs text-gray-500 dark:text-gray-400">NFS Export</p>
                    </div>
                  </div>
                  <button
                    onClick={() => {
                      setSelectedExport(exp);
                      setShowDeleteDialog(true);
                    }}
                    className="p-2 hover:bg-red-100 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                  >
                    <Trash2 className="w-4 h-4 text-red-500" />
                  </button>
                </div>

                {/* Clients */}
                <div className="mb-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Globe className="w-4 h-4 text-gray-600 dark:text-gray-400" />
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                      Allowed Clients ({exp.clients.length})
                    </span>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {exp.clients.map((client, idx) => (
                      <span
                        key={idx}
                        className="px-2 py-1 bg-blue-100 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 text-xs rounded"
                      >
                        {client}
                      </span>
                    ))}
                  </div>
                </div>

                {/* Options */}
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    <Settings className="w-4 h-4 text-gray-600 dark:text-gray-400" />
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                      Export Options
                    </span>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {exp.options.map((option, idx) => (
                      <span
                        key={idx}
                        className="px-2 py-1 bg-gray-200 dark:bg-macos-dark-100 text-gray-700 dark:text-gray-300 text-xs font-mono rounded"
                      >
                        {option}
                      </span>
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
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-2xl w-full m-4 max-h-[80vh] overflow-y-auto"
          >
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Create NFS Export
              </h3>
              <button
                onClick={() => setShowCreateDialog(false)}
                className="p-1 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-6">
              {/* Export Path */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Export Path
                </label>
                <input
                  type="text"
                  value={formData.path}
                  onChange={(e) => setFormData({ ...formData, path: e.target.value })}
                  placeholder="/srv/nfs/share"
                  className="w-full px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent"
                />
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Absolute path to the directory to export
                </p>
              </div>

              {/* Clients */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Allowed Clients
                </label>
                <div className="flex gap-2 mb-2">
                  <input
                    type="text"
                    value={clientInput}
                    onChange={(e) => setClientInput(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && addClient()}
                    placeholder="192.168.1.0/24 or hostname"
                    className="flex-1 px-3 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent"
                  />
                  <button
                    onClick={addClient}
                    className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
                  >
                    Add
                  </button>
                </div>
                <div className="flex flex-wrap gap-2">
                  {formData.clients.map((client, idx) => (
                    <span
                      key={idx}
                      className="inline-flex items-center gap-1 px-2 py-1 bg-blue-100 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 text-sm rounded"
                    >
                      {client}
                      <button
                        onClick={() => removeClient(client)}
                        className="hover:text-red-500"
                      >
                        <X className="w-3 h-3" />
                      </button>
                    </span>
                  ))}
                </div>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Use * for all hosts, IP addresses, CIDR notation, or hostnames
                </p>
              </div>

              {/* Options */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Export Options
                </label>
                <div className="grid grid-cols-2 gap-2">
                  {commonOptions.map((opt) => {
                    const isSelected = formData.options?.includes(opt.value);
                    return (
                      <button
                        key={opt.value}
                        onClick={() => toggleOption(opt.value)}
                        className={`p-3 rounded-lg border-2 transition-all text-left ${
                          isSelected
                            ? 'border-macos-blue bg-macos-blue/10 dark:bg-macos-blue/20'
                            : 'border-gray-300 dark:border-gray-600 hover:border-macos-blue/50'
                        }`}
                      >
                        <div className="font-medium text-gray-900 dark:text-gray-100 text-sm">
                          {opt.label}
                        </div>
                        <div className="text-xs text-gray-600 dark:text-gray-400">
                          {opt.description}
                        </div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Info */}
              <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
                <div className="flex gap-2">
                  <AlertCircle className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <div className="text-sm text-gray-700 dark:text-gray-300">
                    <p className="font-medium mb-1">Important Notes:</p>
                    <ul className="list-disc list-inside space-y-1 text-xs">
                      <li>Ensure the directory exists and has proper permissions</li>
                      <li>Use <code className="bg-white dark:bg-macos-dark-100 px-1 rounded">no_root_squash</code> carefully - it allows root access</li>
                      <li>The export will be added to /etc/exports</li>
                      <li>Changes take effect after NFS restart</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-6 flex justify-end gap-2">
              <button
                onClick={() => setShowCreateDialog(false)}
                className="px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateExport}
                className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
              >
                Create Export
              </button>
            </div>
          </motion.div>
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      {showDeleteDialog && selectedExport && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-md w-full m-4"
          >
            <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
              Delete NFS Export
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              Are you sure you want to delete the export for{' '}
              <span className="font-semibold text-gray-900 dark:text-gray-100">
                {selectedExport.path}
              </span>
              ? This action cannot be undone.
            </p>
            <div className="flex justify-end gap-2">
              <button
                onClick={() => {
                  setShowDeleteDialog(false);
                  setSelectedExport(null);
                }}
                className="px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleDeleteExport}
                className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
              >
                Delete
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}
