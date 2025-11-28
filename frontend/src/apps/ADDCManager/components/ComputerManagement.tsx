import { useState, useEffect } from 'react';
import { addcApi, ADComputer } from '../../../api/addc';
import { motion, AnimatePresence } from 'framer-motion';
import { Monitor, Plus, Trash2, RefreshCw, AlertCircle } from 'lucide-react';

export default function ComputerManagement() {
  const [computers, setComputers] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  const [createForm, setCreateForm] = useState<ADComputer>({
    name: '',
    description: '',
    ou: '',
    ip: '',
    enabled: true,
  });

  const loadComputers = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await addcApi.listComputers();
      if (response.success && response.data) {
        setComputers(response.data);
      } else {
        setError(response.error?.message || 'Failed to load computers');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load computers');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadComputers();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!createForm.name) {
      setError('Computer name is required');
      return;
    }

    try {
      setActionLoading('create');
      setError('');
      const response = await addcApi.createComputer(createForm);

      if (response.success) {
        alert(`Computer ${createForm.name} created successfully!`);
        setShowCreateForm(false);
        setCreateForm({
          name: '',
          description: '',
          ou: '',
          ip: '',
          enabled: true,
        });
        loadComputers();
      } else {
        setError(response.error?.message || 'Failed to create computer');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create computer');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async (name: string) => {
    if (!confirm(`Are you sure you want to delete computer "${name}"?`)) {
      return;
    }

    try {
      setActionLoading(`delete-${name}`);
      setError('');
      const response = await addcApi.deleteComputer(name);

      if (response.success) {
        alert(`Computer ${name} deleted successfully`);
        loadComputers();
      } else {
        setError(response.error?.message || 'Failed to delete computer');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to delete computer');
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Monitor className="w-6 h-6 text-macos-blue" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">Computer Management</h2>
        </div>
        <div className="flex gap-3">
          <button
            onClick={loadComputers}
            disabled={loading}
            className="p-2 text-gray-600 dark:text-gray-400 hover:text-macos-blue dark:hover:text-macos-blue transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
          >
            <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
          </button>
          <button
            onClick={() => setShowCreateForm(true)}
            className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            <Plus className="w-4 h-4" />
            Add Computer
          </button>
        </div>
      </div>

      {/* Error Message */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex items-start gap-3"
          >
            <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <p className="font-medium text-red-900 dark:text-red-100">Error</p>
              <p className="text-red-700 dark:text-red-300 text-sm mt-1">{error}</p>
            </div>
            <button onClick={() => setError('')} className="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300">
              ×
            </button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Computers List */}
      {loading && computers.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <RefreshCw className="w-6 h-6 animate-spin text-gray-400" />
        </div>
      ) : computers.length === 0 ? (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          No computers found
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {computers.map((computerName) => (
            <div
              key={computerName}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-2">
                  <Monitor className="w-5 h-5 text-macos-blue" />
                  <span className="font-medium text-gray-900 dark:text-gray-100">{computerName}</span>
                </div>
              </div>

              <div className="flex flex-wrap gap-2">
                <button
                  onClick={() => handleDelete(computerName)}
                  disabled={actionLoading === `delete-${computerName}`}
                  className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors disabled:opacity-50"
                  title="Delete Computer"
                >
                  {actionLoading === `delete-${computerName}` ? (
                    <RefreshCw className="w-4 h-4 animate-spin" />
                  ) : (
                    <Trash2 className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create Form Modal */}
      <AnimatePresence>
        {showCreateForm && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !actionLoading && setShowCreateForm(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-auto"
            >
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Add Computer
                </h2>

                <form onSubmit={handleCreate} className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Computer Name <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={createForm.name}
                        onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                        placeholder="COMPUTER01"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                        required
                      />
                    </div>

                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Description
                      </label>
                      <input
                        type="text"
                        value={createForm.description}
                        onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                        placeholder="Workstation for IT department"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        IP Address
                      </label>
                      <input
                        type="text"
                        value={createForm.ip}
                        onChange={(e) => setCreateForm({ ...createForm, ip: e.target.value })}
                        placeholder="192.168.1.100"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>

                    <div className="flex items-center">
                      <label className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                        <input
                          type="checkbox"
                          checked={createForm.enabled}
                          onChange={(e) => setCreateForm({ ...createForm, enabled: e.target.checked })}
                          className="w-4 h-4 text-macos-blue border-gray-300 rounded focus:ring-macos-blue"
                        />
                        Enabled
                      </label>
                    </div>

                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Organizational Unit (OU)
                      </label>
                      <input
                        type="text"
                        value={createForm.ou}
                        onChange={(e) => setCreateForm({ ...createForm, ou: e.target.value })}
                        placeholder="OU=Computers,DC=example,DC=com"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>
                  </div>

                  <div className="flex justify-end gap-3 mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <button
                      type="button"
                      onClick={() => setShowCreateForm(false)}
                      disabled={!!actionLoading}
                      className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors disabled:opacity-50"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={!!actionLoading}
                      className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
                    >
                      {actionLoading === 'create' ? (
                        <>
                          <RefreshCw className="w-4 h-4 animate-spin" />
                          Adding...
                        </>
                      ) : (
                        <>
                          <Plus className="w-4 h-4" />
                          Add Computer
                        </>
                      )}
                    </button>
                  </div>
                </form>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
