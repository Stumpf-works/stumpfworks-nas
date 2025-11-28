import { useState, useEffect } from 'react';
import { addcApi, ADOU } from '../../../api/addc';
import { motion, AnimatePresence } from 'framer-motion';
import { FolderTree, Plus, Trash2, RefreshCw, AlertCircle } from 'lucide-react';

export default function OUManagement() {
  const [ous, setOUs] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  const [createForm, setCreateForm] = useState<ADOU>({
    name: '',
    description: '',
    parent_dn: '',
  });

  const loadOUs = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await addcApi.listOUs();
      if (response.success && response.data) {
        setOUs(response.data);
      } else {
        setError(response.error?.message || 'Failed to load OUs');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load OUs');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadOUs();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!createForm.name) {
      setError('OU name is required');
      return;
    }

    try {
      setActionLoading('create');
      setError('');
      const response = await addcApi.createOU(createForm);

      if (response.success) {
        alert(`OU ${createForm.name} created successfully!`);
        setShowCreateForm(false);
        setCreateForm({
          name: '',
          description: '',
          parent_dn: '',
        });
        loadOUs();
      } else {
        setError(response.error?.message || 'Failed to create OU');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create OU');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async (dn: string) => {
    if (!confirm(`Are you sure you want to delete OU "${dn}"?`)) {
      return;
    }

    try {
      setActionLoading(`delete-${dn}`);
      setError('');
      const response = await addcApi.deleteOU(dn);

      if (response.success) {
        alert(`OU ${dn} deleted successfully`);
        loadOUs();
      } else {
        setError(response.error?.message || 'Failed to delete OU');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to delete OU');
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <FolderTree className="w-6 h-6 text-macos-blue" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">Organizational Units</h2>
        </div>
        <div className="flex gap-3">
          <button
            onClick={loadOUs}
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
            Create OU
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

      {/* OUs List */}
      {loading && ous.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <RefreshCw className="w-6 h-6 animate-spin text-gray-400" />
        </div>
      ) : ous.length === 0 ? (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          No organizational units found
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-3">
          {ous.map((ou) => (
            <div
              key={ou}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3 flex-1 min-w-0">
                  <FolderTree className="w-5 h-5 text-macos-blue flex-shrink-0" />
                  <span className="font-mono text-sm text-gray-900 dark:text-gray-100 truncate">{ou}</span>
                </div>
                <button
                  onClick={() => handleDelete(ou)}
                  disabled={actionLoading === `delete-${ou}`}
                  className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors disabled:opacity-50 flex-shrink-0"
                  title="Delete OU"
                >
                  {actionLoading === `delete-${ou}` ? (
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
                  Create Organizational Unit
                </h2>

                <form onSubmit={handleCreate} className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      OU Name <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={createForm.name}
                      onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                      placeholder="IT Department"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Description
                    </label>
                    <input
                      type="text"
                      value={createForm.description}
                      onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                      placeholder="IT Department organizational unit"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Parent DN
                    </label>
                    <input
                      type="text"
                      value={createForm.parent_dn}
                      onChange={(e) => setCreateForm({ ...createForm, parent_dn: e.target.value })}
                      placeholder="DC=example,DC=com"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100 font-mono text-sm"
                    />
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      Leave empty to create at domain root
                    </p>
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
                          Creating...
                        </>
                      ) : (
                        <>
                          <Plus className="w-4 h-4" />
                          Create OU
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
