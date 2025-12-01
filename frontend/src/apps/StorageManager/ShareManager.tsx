import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { storageApi, Share, CreateShareRequest, Volume } from '@/api/storage';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import FolderPicker from '@/components/FolderPicker';
import UserPicker from '@/components/UserPicker';
import GroupPicker from '@/components/GroupPicker';
import { Share2, Network, Upload, Folder, RefreshCw, Plus, Power, Edit2, Trash2, X, Users, UserCheck, CheckCircle, XCircle } from 'lucide-react';

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
      if (response.success && response.data) {
        setShares(response.data);
      } else {
        console.error('Failed to load shares:', response.error);
      }
    } catch (error) {
      console.error('Failed to load shares:', error);
      alert('Failed to load shares. Please check the console for details.');
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
    const iconClass = "w-8 h-8";
    switch (type) {
      case 'smb': return <Share2 className={`${iconClass} text-blue-500`} />;
      case 'nfs': return <Network className={`${iconClass} text-green-500`} />;
      case 'ftp': return <Upload className={`${iconClass} text-orange-500`} />;
      default: return <Folder className={`${iconClass} text-gray-500`} />;
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
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl shadow-lg">
            <Share2 className="w-6 h-6 text-white" />
          </div>
          <div>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Network Shares
            </h2>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {shares.length} share{shares.length !== 1 ? 's' : ''} configured
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button
            onClick={loadShares}
            variant="secondary"
            className="flex items-center gap-2"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </Button>
          <Button
            onClick={() => setShowCreate(true)}
            className="flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            Create Share
          </Button>
        </div>
      </div>

      {/* Shares Grid */}
      {shares.length > 0 && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {shares.map((share) => (
            <motion.div
              key={share.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              whileHover={{ y: -2 }}
              transition={{ duration: 0.2 }}
            >
              <Card className="h-full hover:shadow-xl transition-shadow duration-200">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center space-x-3 flex-1 min-w-0">
                    <div className="flex-shrink-0 p-2 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-xl">
                      {getShareIcon(share.type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-semibold text-gray-900 dark:text-gray-100 truncate">
                        {share.name}
                      </h3>
                      <p className="text-xs text-gray-600 dark:text-gray-400 uppercase">
                        {share.type}
                      </p>
                    </div>
                  </div>
                  <span className={`px-3 py-1 text-xs font-semibold rounded-full shadow-sm flex-shrink-0 ml-2 ${
                    share.enabled
                      ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                      : 'bg-gradient-to-r from-gray-400 to-gray-500 text-white'
                  }`}>
                    {share.enabled ? 'Enabled' : 'Disabled'}
                  </span>
                </div>

                {share.description && (
                  <p className="text-sm text-gray-600 dark:text-gray-400 mb-4 italic">
                    {share.description}
                  </p>
                )}

                {/* Details */}
                <div className="space-y-2 mb-4">
                  <div className="p-3 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-macos-dark-200/50 dark:to-macos-dark-300/50 rounded-lg">
                    <div className="flex items-center gap-2 mb-1">
                      <Folder className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                      <span className="text-xs text-gray-600 dark:text-gray-400">Path</span>
                    </div>
                    <p className="text-sm font-mono font-semibold text-gray-900 dark:text-gray-100 truncate">
                      {share.path}
                    </p>
                  </div>

                  <div className="grid grid-cols-2 gap-2 text-xs">
                    <div className="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                      <CheckCircle className="w-3 h-3" />
                      Read-Only: {share.readOnly ? 'Yes' : 'No'}
                    </div>
                    <div className="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                      <UserCheck className="w-3 h-3" />
                      Guest: {share.guestOk ? 'Allowed' : 'Denied'}
                    </div>
                  </div>

                  {share.validUsers && share.validUsers.length > 0 && (
                    <div>
                      <div className="flex items-center gap-1 text-xs text-gray-600 dark:text-gray-400 mb-1">
                        <UserCheck className="w-3 h-3" />
                        Valid Users:
                      </div>
                      <div className="flex flex-wrap gap-1">
                        {share.validUsers.map((user) => (
                          <span
                            key={user}
                            className="px-2 py-0.5 bg-gradient-to-r from-blue-100 to-blue-200 dark:from-blue-900/30 dark:to-blue-800/30 text-blue-800 dark:text-blue-400 rounded-md text-xs font-medium border border-blue-200 dark:border-blue-700"
                          >
                            {user}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                  {share.validGroups && share.validGroups.length > 0 && (
                    <div>
                      <div className="flex items-center gap-1 text-xs text-gray-600 dark:text-gray-400 mb-1">
                        <Users className="w-3 h-3" />
                        Valid Groups:
                      </div>
                      <div className="flex flex-wrap gap-1">
                        {share.validGroups.map((group) => (
                          <span
                            key={group}
                            className="px-2 py-0.5 bg-gradient-to-r from-purple-100 to-purple-200 dark:from-purple-900/30 dark:to-purple-800/30 text-purple-800 dark:text-purple-400 rounded-md text-xs font-medium border border-purple-200 dark:border-purple-700"
                          >
                            @{group}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                </div>

                {/* Actions */}
                <div className="grid grid-cols-3 gap-2">
                  <Button
                    onClick={() => handleToggle(share.id, share.enabled)}
                    variant="secondary"
                    size="sm"
                    className="flex items-center justify-center gap-1"
                  >
                    <Power className="w-3 h-3" />
                    {share.enabled ? 'Disable' : 'Enable'}
                  </Button>
                  <Button
                    onClick={() => {
                      setEditShare(share);
                      setShowCreate(true);
                    }}
                    variant="secondary"
                    size="sm"
                    className="flex items-center justify-center gap-1"
                  >
                    <Edit2 className="w-3 h-3" />
                    Edit
                  </Button>
                  <Button
                    onClick={() => handleDelete(share.id)}
                    variant="secondary"
                    size="sm"
                    className="flex items-center justify-center gap-1 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
                  >
                    <Trash2 className="w-3 h-3" />
                    Delete
                  </Button>
                </div>
              </Card>
            </motion.div>
          ))}
        </div>
      )}

      {/* Empty State */}
      {shares.length === 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex flex-col items-center justify-center py-16 px-4"
        >
          <div className="p-6 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-macos-dark-200 dark:to-macos-dark-300 rounded-2xl mb-6">
            <Share2 className="w-16 h-16 text-gray-400 dark:text-gray-600" />
          </div>
          <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
            No Shares Found
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-6 text-center max-w-md">
            Create your first network share to start sharing files
          </p>
          <Button
            onClick={() => setShowCreate(true)}
            className="flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            Create Share
          </Button>
        </motion.div>
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
    volumeId: share?.volumeId || '',
    path: share?.path || '',
    type: share?.type || 'smb',
    description: share?.description || '',
    readOnly: share?.readOnly || false,
    browseable: share?.browseable !== undefined ? share.browseable : true,
    guestOk: share?.guestOk || false,
    validUsers: share?.validUsers || [],
    validGroups: share?.validGroups || [],
  });
  const [volumes, setVolumes] = useState<Volume[]>([]);
  const [useManualPath, setUseManualPath] = useState(!share?.volumeId);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadVolumes();
  }, []);

  const loadVolumes = async () => {
    try {
      const response = await storageApi.listVolumes();
      if (response.success && response.data) {
        setVolumes(response.data);
      }
    } catch (error) {
      console.error('Failed to load volumes:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    // Validate that either volumeId or path is provided
    if (!formData.volumeId && !formData.path) {
      setError('Please select a volume or provide a manual path');
      setLoading(false);
      return;
    }

    try {
      const response = share
        ? await storageApi.updateShare(share.id, formData)
        : await storageApi.createShare(formData);

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
      className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-2xl shadow-2xl p-6 w-full max-w-2xl max-h-[80vh] overflow-auto"
      >
        {/* Header */}
        <div className="flex items-center justify-between mb-6 pb-4 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-gradient-to-br from-macos-blue to-macos-purple rounded-xl">
              {share ? <Edit2 className="w-6 h-6 text-white" /> : <Plus className="w-6 h-6 text-white" />}
            </div>
            <div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                {share ? 'Edit Share' : 'Create Share'}
              </h2>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {share ? 'Update share configuration' : 'Configure a new network share'}
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-macos-dark-200 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Share Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="my-share"
            required
          />

          {/* Path Source Toggle */}
          <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-3">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              Path Source
            </label>
            <div className="flex space-x-4 mb-3">
              <label className="flex items-center space-x-2 cursor-pointer">
                <input
                  type="radio"
                  checked={!useManualPath}
                  onChange={() => {
                    setUseManualPath(false);
                    setFormData({ ...formData, path: undefined });
                  }}
                  className="w-4 h-4"
                />
                <span className="text-sm text-gray-700 dark:text-gray-300">
                  Select Volume
                </span>
              </label>
              <label className="flex items-center space-x-2 cursor-pointer">
                <input
                  type="radio"
                  checked={useManualPath}
                  onChange={() => {
                    setUseManualPath(true);
                    setFormData({ ...formData, volumeId: undefined });
                  }}
                  className="w-4 h-4"
                />
                <span className="text-sm text-gray-700 dark:text-gray-300">
                  Manual Path
                </span>
              </label>
            </div>

            {!useManualPath ? (
              <div>
                <select
                  value={formData.volumeId || ''}
                  onChange={(e) => setFormData({ ...formData, volumeId: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100"
                  required={!useManualPath}
                >
                  <option value="">Select a volume...</option>
                  {volumes.filter(v => v.status === 'online').map((volume) => (
                    <option key={volume.id} value={volume.id}>
                      {volume.name} ({volume.mountPoint})
                    </option>
                  ))}
                </select>
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  Select a mounted volume for the share
                </p>
              </div>
            ) : (
              <FolderPicker
                label=""
                value={formData.path || ''}
                onChange={(path) => setFormData({ ...formData, path })}
                placeholder="/mnt/storage/share"
                required={useManualPath}
              />
            )}
          </div>

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

          <UserPicker
            label="Valid Users"
            value={formData.validUsers || []}
            onChange={(users) => setFormData({ ...formData, validUsers: users })}
            placeholder="Select users who can access this share"
            helperText="Individual users who can access this share"
          />

          <GroupPicker
            label="Valid Groups"
            value={formData.validGroups || []}
            onChange={(groups) => setFormData({ ...formData, validGroups: groups })}
            placeholder="Select groups who can access this share"
            helperText="User groups who can access this share"
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
            <div className="p-3 bg-gradient-to-br from-red-50 to-rose-50 dark:from-red-900/20 dark:to-rose-900/20 border border-red-200 dark:border-red-800 rounded-xl flex items-start gap-2">
              <XCircle className="w-4 h-4 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-red-600 dark:text-red-400">
                {error}
              </p>
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
