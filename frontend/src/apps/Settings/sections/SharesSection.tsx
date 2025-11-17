// Revision: 2025-11-17 | Author: Claude | Version: 1.3.0
import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { syslibApi, type SambaShare, type NFSExport } from '@/api/syslib';
import { getErrorMessage } from '@/api/client';

export function SharesSection() {
  const [activeTab, setActiveTab] = useState<'samba' | 'nfs'>('samba');
  const [sambaShares, setSambaShares] = useState<SambaShare[]>([]);
  const [nfsExports, setNfsExports] = useState<NFSExport[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Samba Dialog State
  const [sambaDialog, setSambaDialog] = useState<{
    open: boolean;
    mode: 'create' | 'edit';
    share?: SambaShare;
  }>({ open: false, mode: 'create' });

  const [sambaForm, setSambaForm] = useState({
    name: '',
    path: '',
    comment: '',
    readOnly: false,
    browseable: true,
    guestOK: false,
  });

  // NFS Dialog State
  const [nfsDialog, setNfsDialog] = useState<{
    open: boolean;
    mode: 'create' | 'edit';
    export?: NFSExport;
    index?: number;
  }>({ open: false, mode: 'create' });

  const [nfsForm, setNfsForm] = useState({
    path: '',
    clients: '',
    options: 'rw,sync,no_subtree_check',
  });

  useEffect(() => {
    if (activeTab === 'samba') {
      loadSambaShares();
    } else {
      loadNFSExports();
    }
  }, [activeTab]);

  const loadSambaShares = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await syslibApi.samba.listShares();
      if (response.success && response.data) {
        setSambaShares(response.data);
      } else {
        setError(response.error?.message || 'Failed to load Samba shares');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const loadNFSExports = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await syslibApi.nfs.listExports();
      if (response.success && response.data) {
        setNfsExports(response.data);
      } else {
        setError(response.error?.message || 'Failed to load NFS exports');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  // Samba CRUD Operations
  const handleCreateSambaShare = () => {
    setSambaForm({
      name: '',
      path: '',
      comment: '',
      readOnly: false,
      browseable: true,
      guestOK: false,
    });
    setSambaDialog({ open: true, mode: 'create' });
  };

  const handleEditSambaShare = (share: SambaShare) => {
    setSambaForm({
      name: share.name,
      path: share.path,
      comment: share.comment || '',
      readOnly: share.readOnly || false,
      browseable: share.browseable !== false,
      guestOK: share.guestOK || false,
    });
    setSambaDialog({ open: true, mode: 'edit', share });
  };

  const handleSaveSambaShare = async () => {
    setError(null);
    setSuccess(null);

    try {
      const response = await syslibApi.samba.createShare(sambaForm);
      if (response.success) {
        setSuccess('Samba share created successfully');
        setSambaDialog({ open: false, mode: 'create' });
        loadSambaShares();
      } else {
        setError(response.error?.message || 'Failed to create share');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleDeleteSambaShare = async (name: string) => {
    if (!confirm(`Are you sure you want to delete share "${name}"?`)) {
      return;
    }

    try {
      const response = await syslibApi.samba.deleteShare(name);
      if (response.success) {
        setSuccess('Samba share deleted successfully');
        loadSambaShares();
      } else {
        setError(response.error?.message || 'Failed to delete share');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleRestartSamba = async () => {
    try {
      const response = await syslibApi.samba.restart();
      if (response.success) {
        setSuccess('Samba service restarted successfully');
      } else {
        setError(response.error?.message || 'Failed to restart Samba');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  // NFS CRUD Operations
  const handleCreateNFSExport = () => {
    setNfsForm({
      path: '',
      clients: '*',
      options: 'rw,sync,no_subtree_check',
    });
    setNfsDialog({ open: true, mode: 'create' });
  };

  const handleEditNFSExport = (nfsExport: NFSExport, index: number) => {
    setNfsForm({
      path: nfsExport.path,
      clients: nfsExport.clients.join(' '),
      options: (nfsExport.options || []).join(',') || 'rw,sync,no_subtree_check',
    });
    setNfsDialog({ open: true, mode: 'edit', export: nfsExport, index });
  };

  const handleSaveNFSExport = async () => {
    setError(null);
    setSuccess(null);

    const payload = {
      path: nfsForm.path,
      clients: nfsForm.clients.split(/[\s,]+/).filter((c) => c),
      options: nfsForm.options.split(',').map((o) => o.trim()).filter((o) => o),
    };

    try {
      const response = await syslibApi.nfs.createExport(payload);
      if (response.success) {
        setSuccess(`NFS export ${nfsDialog.mode === 'create' ? 'created' : 'updated'} successfully`);
        setNfsDialog({ open: false, mode: 'create' });
        loadNFSExports();
      } else {
        setError(response.error?.message || 'Failed to save export');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleDeleteNFSExport = async (nfsExport: NFSExport) => {
    if (!confirm(`Are you sure you want to delete export "${nfsExport.path}"?`)) {
      return;
    }

    try {
      const response = await syslibApi.nfs.deleteExport(nfsExport.path);
      if (response.success) {
        setSuccess('NFS export deleted successfully');
        loadNFSExports();
      } else {
        setError(response.error?.message || 'Failed to delete export');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleRestartNFS = async () => {
    try {
      const response = await syslibApi.nfs.restart();
      if (response.success) {
        setSuccess('NFS service restarted successfully');
      } else {
        setError(response.error?.message || 'Failed to restart NFS');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  return (
    <div className="space-y-4 md:space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-xl md:text-2xl font-bold text-gray-900 dark:text-gray-100">
          Shares Management
        </h1>
        <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mt-1">
          Configure SMB/CIFS (Samba) and NFS network shares
        </p>
      </div>

      {/* Success/Error Messages */}
      {error && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="p-3 md:p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg"
        >
          <p className="text-sm md:text-base text-red-800 dark:text-red-200">{error}</p>
        </motion.div>
      )}

      {success && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="p-3 md:p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg"
        >
          <p className="text-sm md:text-base text-green-800 dark:text-green-200">{success}</p>
        </motion.div>
      )}

      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 dark:border-macos-dark-300">
        <button
          onClick={() => setActiveTab('samba')}
          className={`px-3 md:px-4 py-2 font-medium border-b-2 transition-colors text-sm md:text-base ${
            activeTab === 'samba'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          Samba (SMB/CIFS) <span className="hidden sm:inline">({sambaShares.length})</span>
        </button>
        <button
          onClick={() => setActiveTab('nfs')}
          className={`px-3 md:px-4 py-2 font-medium border-b-2 transition-colors text-sm md:text-base ${
            activeTab === 'nfs'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          NFS <span className="hidden sm:inline">({nfsExports.length})</span>
        </button>
      </div>

      {/* Samba Shares Tab */}
      {activeTab === 'samba' && (
        <div className="space-y-4">
          <div className="flex flex-col sm:flex-row justify-end gap-2">
            <Button variant="secondary" onClick={handleRestartSamba} size="sm">
              <span className="hidden sm:inline">Restart Samba Service</span>
              <span className="sm:hidden">Restart</span>
            </Button>
            <Button onClick={handleCreateSambaShare}>
              <span className="hidden sm:inline">+ Create Share</span>
              <span className="sm:hidden">+ Share</span>
            </Button>
          </div>

          <Card>
            <div className="p-4 md:p-6">
              {loading ? (
                <p className="text-sm md:text-base text-gray-600 dark:text-gray-400">
                  Loading Samba shares...
                </p>
              ) : sambaShares.length === 0 ? (
                <div className="text-center py-8 md:py-12">
                  <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mb-4">
                    No Samba shares configured
                  </p>
                  <Button onClick={handleCreateSambaShare}>Create your first share</Button>
                </div>
              ) : (
                <div className="space-y-3">
                  {sambaShares.map((share) => (
                    <motion.div
                      key={share.name}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="p-3 md:p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg hover:border-macos-blue transition-colors"
                    >
                      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 md:gap-3 flex-wrap">
                            <h3 className="text-base md:text-lg font-semibold text-gray-900 dark:text-gray-100">
                              {share.name}
                            </h3>
                            {share.readOnly && (
                              <span className="px-2 py-1 text-xs bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-200 rounded">
                                Read-Only
                              </span>
                            )}
                            {share.guestOK && (
                              <span className="px-2 py-1 text-xs bg-yellow-100 dark:bg-yellow-900/20 text-yellow-800 dark:text-yellow-200 rounded">
                                Guest Access
                              </span>
                            )}
                          </div>
                          <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400 mt-1">
                            Path: {share.path}
                          </p>
                          {share.comment && (
                            <p className="text-xs md:text-sm text-gray-500 dark:text-gray-500 mt-1">
                              {share.comment}
                            </p>
                          )}
                        </div>
                        <div className="flex gap-2 w-full sm:w-auto">
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() => handleEditSambaShare(share)}
                            className="flex-1 sm:flex-initial"
                          >
                            Edit
                          </Button>
                          <Button
                            variant="danger"
                            size="sm"
                            onClick={() => handleDeleteSambaShare(share.name)}
                            className="flex-1 sm:flex-initial"
                          >
                            Delete
                          </Button>
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              )}
            </div>
          </Card>
        </div>
      )}

      {/* NFS Exports Tab */}
      {activeTab === 'nfs' && (
        <div className="space-y-4">
          <div className="flex flex-col sm:flex-row justify-end gap-2">
            <Button variant="secondary" onClick={handleRestartNFS} size="sm">
              <span className="hidden sm:inline">Restart NFS Service</span>
              <span className="sm:hidden">Restart</span>
            </Button>
            <Button onClick={handleCreateNFSExport}>
              <span className="hidden sm:inline">+ Create Export</span>
              <span className="sm:hidden">+ Export</span>
            </Button>
          </div>

          <Card>
            <div className="p-4 md:p-6">
              {loading ? (
                <p className="text-sm md:text-base text-gray-600 dark:text-gray-400">
                  Loading NFS exports...
                </p>
              ) : nfsExports.length === 0 ? (
                <div className="text-center py-8 md:py-12">
                  <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mb-4">
                    No NFS exports configured
                  </p>
                  <Button onClick={handleCreateNFSExport}>Create your first export</Button>
                </div>
              ) : (
                <div className="space-y-3">
                  {nfsExports.map((nfsExport, idx) => (
                    <motion.div
                      key={idx}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="p-3 md:p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg hover:border-macos-blue transition-colors"
                    >
                      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
                        <div className="flex-1">
                          <h3 className="text-base md:text-lg font-semibold text-gray-900 dark:text-gray-100">
                            {nfsExport.path}
                          </h3>
                          <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400 mt-1">
                            Clients: {nfsExport.clients.join(', ')}
                          </p>
                          {nfsExport.options && (
                            <p className="text-xs md:text-sm text-gray-500 dark:text-gray-500 mt-1">
                              Options: {nfsExport.options}
                            </p>
                          )}
                        </div>
                        <div className="flex gap-2 w-full sm:w-auto">
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() => handleEditNFSExport(nfsExport, idx)}
                            className="flex-1 sm:flex-initial"
                          >
                            Edit
                          </Button>
                          <Button
                            variant="danger"
                            size="sm"
                            onClick={() => handleDeleteNFSExport(nfsExport)}
                            className="flex-1 sm:flex-initial"
                          >
                            Delete
                          </Button>
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              )}
            </div>
          </Card>
        </div>
      )}

      {/* Samba Create/Edit Dialog */}
      <AnimatePresence>
        {sambaDialog.open && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
            onClick={() => setSambaDialog({ open: false, mode: 'create' })}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-auto"
            >
              <div className="p-4 md:p-6">
                <h2 className="text-lg md:text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                  {sambaDialog.mode === 'create' ? 'Create Samba Share' : 'Edit Samba Share'}
                </h2>

                <div className="space-y-4">
                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Share Name
                    </label>
                    <Input
                      value={sambaForm.name}
                      onChange={(e) => setSambaForm({ ...sambaForm, name: e.target.value })}
                      disabled={sambaDialog.mode === 'edit'}
                      placeholder="e.g., Documents, Media"
                    />
                  </div>

                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Path
                    </label>
                    <Input
                      value={sambaForm.path}
                      onChange={(e) => setSambaForm({ ...sambaForm, path: e.target.value })}
                      placeholder="/mnt/storage/share"
                    />
                  </div>

                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Comment (Optional)
                    </label>
                    <Input
                      value={sambaForm.comment}
                      onChange={(e) => setSambaForm({ ...sambaForm, comment: e.target.value })}
                      placeholder="Description of this share"
                    />
                  </div>

                  <div className="space-y-3 p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={sambaForm.readOnly}
                        onChange={(e) => setSambaForm({ ...sambaForm, readOnly: e.target.checked })}
                        className="w-4 h-4 text-macos-blue bg-white dark:bg-macos-dark-100 border-gray-300 dark:border-macos-dark-300 rounded focus:ring-macos-blue"
                      />
                      <span className="text-xs md:text-sm text-gray-900 dark:text-gray-100">
                        Read-Only
                      </span>
                    </label>

                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={sambaForm.browseable}
                        onChange={(e) =>
                          setSambaForm({ ...sambaForm, browseable: e.target.checked })
                        }
                        className="w-4 h-4 text-macos-blue bg-white dark:bg-macos-dark-100 border-gray-300 dark:border-macos-dark-300 rounded focus:ring-macos-blue"
                      />
                      <span className="text-xs md:text-sm text-gray-900 dark:text-gray-100">
                        Browseable
                      </span>
                    </label>

                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={sambaForm.guestOK}
                        onChange={(e) => setSambaForm({ ...sambaForm, guestOK: e.target.checked })}
                        className="w-4 h-4 text-macos-blue bg-white dark:bg-macos-dark-100 border-gray-300 dark:border-macos-dark-300 rounded focus:ring-macos-blue"
                      />
                      <span className="text-xs md:text-sm text-gray-900 dark:text-gray-100">
                        Guest Access (No Password)
                      </span>
                    </label>
                  </div>
                </div>

                <div className="flex gap-2 mt-6">
                  <Button
                    variant="secondary"
                    onClick={() => setSambaDialog({ open: false, mode: 'create' })}
                    className="flex-1"
                  >
                    Cancel
                  </Button>
                  <Button onClick={handleSaveSambaShare} className="flex-1">
                    {sambaDialog.mode === 'create' ? 'Create' : 'Save'}
                  </Button>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* NFS Create/Edit Dialog */}
      <AnimatePresence>
        {nfsDialog.open && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
            onClick={() => setNfsDialog({ open: false, mode: 'create' })}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-auto"
            >
              <div className="p-4 md:p-6">
                <h2 className="text-lg md:text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                  {nfsDialog.mode === 'create' ? 'Create NFS Export' : 'Edit NFS Export'}
                </h2>

                <div className="space-y-4">
                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Export Path
                    </label>
                    <Input
                      value={nfsForm.path}
                      onChange={(e) => setNfsForm({ ...nfsForm, path: e.target.value })}
                      placeholder="/mnt/storage/nfs-export"
                    />
                  </div>

                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Allowed Clients
                    </label>
                    <Input
                      value={nfsForm.clients}
                      onChange={(e) => setNfsForm({ ...nfsForm, clients: e.target.value })}
                      placeholder="* or 192.168.1.0/24 or host1 host2"
                    />
                    <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                      Use * for all, IP ranges, or space-separated hostnames
                    </p>
                  </div>

                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Mount Options
                    </label>
                    <Input
                      value={nfsForm.options}
                      onChange={(e) => setNfsForm({ ...nfsForm, options: e.target.value })}
                      placeholder="rw,sync,no_subtree_check"
                    />
                    <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                      Common: rw (read-write), ro (read-only), sync, async
                    </p>
                  </div>

                  <div className="p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
                    <p className="text-xs md:text-sm text-blue-800 dark:text-blue-200">
                      <strong>Note:</strong> Changes require NFS service restart to take effect.
                    </p>
                  </div>
                </div>

                <div className="flex gap-2 mt-6">
                  <Button
                    variant="secondary"
                    onClick={() => setNfsDialog({ open: false, mode: 'create' })}
                    className="flex-1"
                  >
                    Cancel
                  </Button>
                  <Button onClick={handleSaveNFSExport} className="flex-1">
                    {nfsDialog.mode === 'create' ? 'Create' : 'Save'}
                  </Button>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
