import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { storageApi, Share, CreateShareRequest } from '@/api/storage';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

export default function ShareManager() {
  const [shares, setShares] = useState<Share[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [editShare, setEditShare] = useState<Share | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadShares();
  }, []);

  const loadShares = async () => {
    try {
      const response = await storageApi.listShares();
      if (response.success) {
        setShares(response.data);
      }
    } catch (error) {
      console.error('Failed to load shares:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggle = async (id: string, enabled: boolean) => {
    try {
      if (enabled) {
        await storageApi.disableShare(id);
      } else {
        await storageApi.enableShare(id);
      }
      loadShares();
    } catch (error) {
      console.error('Failed to toggle share:', error);
      alert('Failed to toggle share');
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this share?')) {
      return;
    }

    try {
      await storageApi.deleteShare(id);
      loadShares();
    } catch (error) {
      console.error('Failed to delete share:', error);
      alert('Failed to delete share');
    }
  };

  const getShareIcon = (type: string) => {
    switch (type) {
      case 'smb': return '🖥️';
      case 'nfs': return '🌐';
      case 'ftp': return '📤';
      default: return '📁';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Network Shares ({shares.length})
        </h2>
        <div className="flex space-x-2">
          <Button onClick={loadShares} variant="secondary">
            🔄 Refresh
          </Button>
          <Button onClick={() => setShowCreate(true)}>
            ➕ Create Share
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {shares.map((share) => (
          <Card key={share.id}>
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center space-x-3">
                <div className="text-3xl">{getShareIcon(share.type)}</div>
                <div>
                  <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                    {share.name}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400 uppercase">
                    {share.type}
                  </p>
                </div>
              </div>
              <div className="flex items-center space-x-2">
                {share.enabled ? (
                  <span className="px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-400 rounded-full text-xs font-medium">
                    Enabled
                  </span>
                ) : (
                  <span className="px-2 py-1 bg-gray-100 dark:bg-gray-900/30 text-gray-800 dark:text-gray-400 rounded-full text-xs font-medium">
                    Disabled
                  </span>
                )}
              </div>
            </div>

            {share.description && (
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                {share.description}
              </p>
            )}

            <div className="space-y-2 mb-4 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-400">Path:</span>
                <span className="font-mono text-gray-900 dark:text-gray-100">
                  {share.path}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-400">Read-Only:</span>
                <span className="text-gray-900 dark:text-gray-100">
                  {share.readOnly ? 'Yes' : 'No'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-400">Guest Access:</span>
                <span className="text-gray-900 dark:text-gray-100">
                  {share.guestOk ? 'Allowed' : 'Denied'}
                </span>
              </div>
              {share.validUsers && share.validUsers.length > 0 && (
                <div>
                  <span className="text-gray-600 dark:text-gray-400">Valid Users:</span>
                  <div className="mt-1 flex flex-wrap gap-1">
                    {share.validUsers.map((user) => (
                      <span
                        key={user}
                        className="px-2 py-0.5 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 rounded text-xs"
                      >
                        {user}
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>

            <div className="flex space-x-2">
              <Button
                onClick={() => handleToggle(share.id, share.enabled)}
                variant="secondary"
                size="sm"
                className="flex-1"
              >
                {share.enabled ? '⏸️ Disable' : '▶️ Enable'}
              </Button>
              <Button
                onClick={() => {
                  setEditShare(share);
                  setShowCreate(true);
                }}
                variant="secondary"
                size="sm"
                className="flex-1"
              >
                ✏️ Edit
              </Button>
              <Button
                onClick={() => handleDelete(share.id)}
                variant="secondary"
                size="sm"
                className="text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
              >
                🗑️
              </Button>
            </div>
          </Card>
        ))}
      </div>

      {shares.length === 0 && (
        <div className="text-center py-12 text-gray-600 dark:text-gray-400">
          <div className="text-6xl mb-4">📁</div>
          <p className="text-lg font-medium mb-2">No shares found</p>
          <p className="text-sm mb-4">Create your first network share to get started</p>
          <Button onClick={() => setShowCreate(true)}>
            ➕ Create Share
          </Button>
        </div>
      )}

      {/* Create/Edit Share Modal */}
      <AnimatePresence>
        {showCreate && (
          <ShareModal
            share={editShare}
            onClose={() => {
              setShowCreate(false);
              setEditShare(null);
            }}
            onSuccess={() => {
              setShowCreate(false);
              setEditShare(null);
              loadShares();
            }}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

interface ShareModalProps {
  share: Share | null;
  onClose: () => void;
  onSuccess: () => void;
}

function ShareModal({ share, onClose, onSuccess }: ShareModalProps) {
  const [formData, setFormData] = useState<CreateShareRequest>({
    name: share?.name || '',
    path: share?.path || '',
    type: share?.type || 'smb',
    description: share?.description || '',
    readOnly: share?.readOnly || false,
    browseable: share?.browseable !== undefined ? share.browseable : true,
    guestOk: share?.guestOk || false,
    validUsers: share?.validUsers || [],
  });
  const [validUsersInput, setValidUsersInput] = useState(
    share?.validUsers?.join(', ') || ''
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    // Parse valid users
    const validUsers = validUsersInput
      .split(',')
      .map((u) => u.trim())
      .filter((u) => u.length > 0);

    const data = { ...formData, validUsers };

    try {
      const response = share
        ? await storageApi.updateShare(share.id, data)
        : await storageApi.createShare(data);

      if (response.success) {
        onSuccess();
      } else {
        setError(response.error?.message || `Failed to ${share ? 'update' : 'create'} share`);
      }
    } catch (err: any) {
      setError(err.message || `Failed to ${share ? 'update' : 'create'} share`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-2xl max-h-[80vh] overflow-auto"
      >
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">
          {share ? 'Edit Share' : 'Create Share'}
        </h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Share Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="my-share"
            required
          />

          <Input
            label="Path"
            value={formData.path}
            onChange={(e) => setFormData({ ...formData, path: e.target.value })}
            placeholder="/mnt/storage/share"
            required
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Share Type
            </label>
            <select
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value as any })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100"
            >
              <option value="smb">SMB/CIFS (Windows)</option>
              <option value="nfs">NFS (Linux/Unix)</option>
            </select>
          </div>

          <Input
            label="Description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="Optional description"
          />

          <Input
            label="Valid Users (comma-separated)"
            value={validUsersInput}
            onChange={(e) => setValidUsersInput(e.target.value)}
            placeholder="user1, user2, user3"
          />

          <div className="space-y-3">
            <label className="flex items-center space-x-3 cursor-pointer">
              <input
                type="checkbox"
                checked={formData.readOnly}
                onChange={(e) => setFormData({ ...formData, readOnly: e.target.checked })}
                className="w-4 h-4"
              />
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                Read-Only
              </span>
            </label>

            <label className="flex items-center space-x-3 cursor-pointer">
              <input
                type="checkbox"
                checked={formData.browseable}
                onChange={(e) => setFormData({ ...formData, browseable: e.target.checked })}
                className="w-4 h-4"
              />
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                Browseable
              </span>
            </label>

            <label className="flex items-center space-x-3 cursor-pointer">
              <input
                type="checkbox"
                checked={formData.guestOk}
                onChange={(e) => setFormData({ ...formData, guestOk: e.target.checked })}
                className="w-4 h-4"
              />
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                Allow Guest Access
              </span>
            </label>
          </div>

          {error && (
            <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400 text-sm">
              {error}
            </div>
          )}

          <div className="flex space-x-3 pt-4">
            <Button type="button" onClick={onClose} variant="secondary" className="flex-1">
              Cancel
            </Button>
            <Button type="submit" isLoading={loading} className="flex-1">
              {share ? 'Update' : 'Create'} Share
            </Button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
}
