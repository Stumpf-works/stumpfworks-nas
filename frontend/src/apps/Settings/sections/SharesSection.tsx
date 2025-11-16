// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState, useEffect } from 'react';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { syslibApi, type SambaShare, type NFSExport } from '@/api/syslib';
import { getErrorMessage } from '@/api/client';

export function SharesSection({ user, systemInfo }: { user: any; systemInfo: any }) {
  const [activeTab, setActiveTab] = useState<'samba' | 'nfs'>('samba');
  const [sambaShares, setSambaShares] = useState<SambaShare[]>([]);
  const [nfsExports, setNfsExports] = useState<NFSExport[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showCreateSamba, setShowCreateSamba] = useState(false);
  const [newShare, setNewShare] = useState({ name: '', path: '', comment: '' });

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

  const handleCreateSambaShare = async () => {
    try {
      const response = await syslibApi.samba.createShare({
        name: newShare.name,
        path: newShare.path,
        comment: newShare.comment,
      });
      if (response.success) {
        setShowCreateSamba(false);
        setNewShare({ name: '', path: '', comment: '' });
        await loadSambaShares();
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleDeleteSambaShare = async (name: string) => {
    if (!confirm(`Delete share "${name}"?`)) return;
    try {
      const response = await syslibApi.samba.deleteShare(name);
      if (response.success) {
        await loadSambaShares();
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Shares Management</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure Samba, NFS, iSCSI, WebDAV, and FTP shares
        </p>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 dark:border-macos-dark-300">
        <button
          onClick={() => setActiveTab('samba')}
          className={`px-4 py-2 font-medium border-b-2 transition-colors ${
            activeTab === 'samba'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          Samba (SMB)
        </button>
        <button
          onClick={() => setActiveTab('nfs')}
          className={`px-4 py-2 font-medium border-b-2 transition-colors ${
            activeTab === 'nfs'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          NFS
        </button>
      </div>

      {/* Samba Shares */}
      {activeTab === 'samba' && (
        <Card>
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                Samba Shares
              </h2>
              <Button variant="primary" onClick={() => setShowCreateSamba(true)}>
                Create Share
              </Button>
            </div>

            {showCreateSamba && (
              <div className="mb-4 p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg">
                <h3 className="font-medium text-gray-900 dark:text-gray-100 mb-3">New Samba Share</h3>
                <div className="space-y-3">
                  <Input
                    placeholder="Share name"
                    value={newShare.name}
                    onChange={(e) => setNewShare({ ...newShare, name: e.target.value })}
                  />
                  <Input
                    placeholder="Path (e.g., /mnt/pool/share)"
                    value={newShare.path}
                    onChange={(e) => setNewShare({ ...newShare, path: e.target.value })}
                  />
                  <Input
                    placeholder="Comment (optional)"
                    value={newShare.comment}
                    onChange={(e) => setNewShare({ ...newShare, comment: e.target.value })}
                  />
                  <div className="flex gap-2">
                    <Button variant="primary" onClick={handleCreateSambaShare}>Create</Button>
                    <Button variant="secondary" onClick={() => setShowCreateSamba(false)}>Cancel</Button>
                  </div>
                </div>
              </div>
            )}

            {loading ? (
              <p className="text-gray-600 dark:text-gray-400">Loading shares...</p>
            ) : sambaShares.length === 0 ? (
              <p className="text-gray-600 dark:text-gray-400">No Samba shares configured</p>
            ) : (
              <div className="space-y-3">
                {sambaShares.map((share) => (
                  <div
                    key={share.name}
                    className="p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg flex items-center justify-between"
                  >
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100">{share.name}</h3>
                      <p className="text-sm text-gray-600 dark:text-gray-400">{share.path}</p>
                      {share.comment && (
                        <p className="text-sm text-gray-500 dark:text-gray-500">{share.comment}</p>
                      )}
                    </div>
                    <Button size="sm" variant="danger" onClick={() => handleDeleteSambaShare(share.name)}>
                      Delete
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </Card>
      )}

      {/* NFS Exports */}
      {activeTab === 'nfs' && (
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              NFS Exports
            </h2>
            {loading ? (
              <p className="text-gray-600 dark:text-gray-400">Loading exports...</p>
            ) : nfsExports.length === 0 ? (
              <p className="text-gray-600 dark:text-gray-400">No NFS exports configured</p>
            ) : (
              <div className="space-y-3">
                {nfsExports.map((exp, idx) => (
                  <div
                    key={idx}
                    className="p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg"
                  >
                    <h3 className="font-semibold text-gray-900 dark:text-gray-100">{exp.path}</h3>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      Clients: {exp.clients.join(', ')}
                    </p>
                  </div>
                ))}
              </div>
            )}
          </div>
        </Card>
      )}
    </div>
  );
}
