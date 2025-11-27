import { useState, useEffect } from 'react';
import { addcApi, ADGPO } from '../../../api/addc';
import { motion, AnimatePresence } from 'framer-motion';
import { FileText, Plus, Trash2, RefreshCw, AlertCircle, Link, Unlink } from 'lucide-react';

export default function GPOManagement() {
  const [gpos, setGPOs] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showLinkForm, setShowLinkForm] = useState(false);
  const [selectedGPO, setSelectedGPO] = useState('');
  const [ouDN, setOUDN] = useState('');
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  const [createForm, setCreateForm] = useState<ADGPO>({
    name: '',
    display_name: '',
    description: '',
  });

  const loadGPOs = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await addcApi.listGPOs();
      if (response.success && response.data) {
        setGPOs(response.data);
      } else {
        setError(response.error?.message || 'Failed to load GPOs');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load GPOs');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadGPOs();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!createForm.name) {
      setError('GPO name is required');
      return;
    }

    try {
      setActionLoading('create');
      setError('');
      const response = await addcApi.createGPO(createForm);

      if (response.success) {
        alert(`GPO ${createForm.name} created successfully!`);
        setShowCreateForm(false);
        setCreateForm({
          name: '',
          display_name: '',
          description: '',
        });
        loadGPOs();
      } else {
        setError(response.error?.message || 'Failed to create GPO');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create GPO');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async (name: string) => {
    if (!confirm(`Are you sure you want to delete GPO "${name}"?`)) {
      return;
    }

    try {
      setActionLoading(`delete-${name}`);
      setError('');
      const response = await addcApi.deleteGPO(name);

      if (response.success) {
        alert(`GPO ${name} deleted successfully`);
        loadGPOs();
      } else {
        setError(response.error?.message || 'Failed to delete GPO');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to delete GPO');
    } finally {
      setActionLoading(null);
    }
  };

  const handleLink = async () => {
    if (!ouDN.trim()) {
      setError('OU DN is required');
      return;
    }

    try {
      setActionLoading(`link-${selectedGPO}`);
      setError('');
      const response = await addcApi.linkGPO(selectedGPO, ouDN);

      if (response.success) {
        alert(`GPO ${selectedGPO} linked to ${ouDN}`);
        setShowLinkForm(false);
        setOUDN('');
      } else {
        setError(response.error?.message || 'Failed to link GPO');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to link GPO');
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <FileText className="w-6 h-6 text-macos-blue" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">Group Policy Management</h2>
        </div>
        <div className="flex gap-3">
          <button
            onClick={loadGPOs}
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
            Create GPO
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

      {/* GPOs List */}
      {loading && gpos.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <RefreshCw className="w-6 h-6 animate-spin text-gray-400" />
        </div>
      ) : gpos.length === 0 ? (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          No GPOs found
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {gpos.map((gpoName) => (
            <div
              key={gpoName}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-2">
                  <FileText className="w-5 h-5 text-macos-blue" />
                  <span className="font-medium text-gray-900 dark:text-gray-100">{gpoName}</span>
                </div>
              </div>

              <div className="flex flex-wrap gap-2">
                <button
                  onClick={() => {
                    setSelectedGPO(gpoName);
                    setShowLinkForm(true);
                  }}
                  className="p-1.5 text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded transition-colors"
                  title="Link to OU"
                >
                  <Link className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(gpoName)}
                  disabled={actionLoading === `delete-${gpoName}`}
                  className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors disabled:opacity-50"
                  title="Delete GPO"
                >
                  {actionLoading === `delete-${gpoName}` ? (
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
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-2xl w-full"
            >
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Create Group Policy Object
                </h2>

                <form onSubmit={handleCreate} className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      GPO Name <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={createForm.name}
                      onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                      placeholder="Security Policy"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Display Name
                    </label>
                    <input
                      type="text"
                      value={createForm.display_name}
                      onChange={(e) => setCreateForm({ ...createForm, display_name: e.target.value })}
                      placeholder="Default Security Policy"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Description
                    </label>
                    <textarea
                      value={createForm.description}
                      onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                      placeholder="Policy description"
                      rows={3}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                    />
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
                          Create GPO
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

      {/* Link Form Modal */}
      <AnimatePresence>
        {showLinkForm && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !actionLoading && setShowLinkForm(false)}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-lg w-full"
            >
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Link GPO to OU
                </h2>

                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      GPO
                    </label>
                    <input
                      type="text"
                      value={selectedGPO}
                      disabled
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      OU Distinguished Name <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={ouDN}
                      onChange={(e) => setOUDN(e.target.value)}
                      placeholder="OU=IT,DC=example,DC=com"
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100 font-mono text-sm"
                    />
                  </div>

                  <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <button
                      onClick={() => setShowLinkForm(false)}
                      disabled={!!actionLoading}
                      className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors disabled:opacity-50"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={handleLink}
                      disabled={!!actionLoading}
                      className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
                    >
                      {actionLoading === `link-${selectedGPO}` ? (
                        <>
                          <RefreshCw className="w-4 h-4 animate-spin" />
                          Linking...
                        </>
                      ) : (
                        <>
                          <Link className="w-4 h-4" />
                          Link GPO
                        </>
                      )}
                    </button>
                  </div>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
