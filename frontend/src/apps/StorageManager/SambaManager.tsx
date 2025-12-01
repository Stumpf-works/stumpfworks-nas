// Revision: 2025-12-01 | Author: StumpfWorks AI | Version: 1.2.0
import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Share2,
  Plus,
  RefreshCw,
  Trash2,
  Folder,
  Users,
  Shield,
  Lock,
  Eye,
  EyeOff,
  Power,
  CheckCircle,
  XCircle,
  AlertCircle,
} from 'lucide-react';
import { syslibApi, type SambaShare, type CreateSambaShareRequest } from '@/api/syslib';
import { storageApi, type Share } from '@/api/storage';

export default function SambaManager() {
  const [shares, setShares] = useState<SambaShare[]>([]);
  const [shareIdMap, setShareIdMap] = useState<Map<string, string>>(new Map()); // Maps share name to database ID
  const [selectedShare, setSelectedShare] = useState<SambaShare | null>(null);
  const [serviceStatus, setServiceStatus] = useState<{ active: boolean; enabled: boolean } | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  // Form state for creating share
  const [formData, setFormData] = useState<CreateSambaShareRequest>({
    name: '',
    path: '',
    comment: '',
    valid_users: [],
    valid_groups: [],
    read_only: false,
    browseable: true,
    guest_ok: false,
    recycle_bin: false,
  });

  // Fetch shares from database (filtered for SMB type)
  const fetchShares = async () => {
    setIsLoading(true);
    try {
      const response = await storageApi.listShares();
      if (response.success && response.data) {
        // Filter for SMB shares
        const smbShares = response.data.filter((share: Share) => share.type === 'smb');

        // Build ID map (share name -> database ID)
        const idMap = new Map<string, string>();
        smbShares.forEach((share: Share) => {
          idMap.set(share.name, share.id);
        });
        setShareIdMap(idMap);

        // Map to SambaShare format for display
        const sambaShares: SambaShare[] = smbShares.map((share: Share) => ({
          name: share.name,
          path: share.path,
          comment: share.description || '',
          validUsers: share.validUsers || [],
          validGroups: share.validGroups || [],
          readOnly: share.readOnly,
          browseable: share.browseable,
          guestOK: share.guestOk,
          recycleBin: false, // Not stored in database yet, default to false
        }));
        setShares(sambaShares);
      }
    } catch (error) {
      console.error('Failed to fetch Samba shares:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch service status
  const fetchServiceStatus = async () => {
    try {
      const response = await syslibApi.samba.getStatus();
      if (response.success && response.data) {
        setServiceStatus(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch Samba status:', error);
    }
  };

  useEffect(() => {
    fetchShares();
    fetchServiceStatus();
  }, []);

  const handleCreateShare = async () => {
    if (!formData.name || !formData.path) {
      alert('Please fill in required fields (name and path)');
      return;
    }

    try {
      // Create share using storage API (saves to database)
      const response = await storageApi.createShare({
        name: formData.name,
        path: formData.path,
        type: 'smb',
        description: formData.comment || '',
        readOnly: formData.read_only || false,
        browseable: formData.browseable !== false,
        guestOk: formData.guest_ok || false,
        validUsers: formData.valid_users,
        validGroups: formData.valid_groups,
      });
      if (response.success) {
        alert(`Share "${formData.name}" created successfully`);
        setShowCreateDialog(false);
        setFormData({
          name: '',
          path: '',
          comment: '',
          valid_users: [],
          valid_groups: [],
          read_only: false,
          browseable: true,
          guest_ok: false,
          recycle_bin: false,
        });
        fetchShares();
      }
    } catch (error) {
      console.error('Failed to create share:', error);
      alert('Failed to create share');
    }
  };

  const handleDeleteShare = async () => {
    if (!selectedShare) return;

    // Get the database ID for this share
    const shareId = shareIdMap.get(selectedShare.name);
    if (!shareId) {
      alert('Failed to delete share: ID not found');
      return;
    }

    try {
      const response = await storageApi.deleteShare(shareId);
      if (response.success) {
        alert('Share deleted successfully');
        setShowDeleteConfirm(false);
        setSelectedShare(null);
        fetchShares();
      }
    } catch (error) {
      console.error('Failed to delete share:', error);
      alert('Failed to delete share');
    }
  };

  const handleRestartService = async () => {
    try {
      const response = await syslibApi.samba.restart();
      if (response.success) {
        alert('Samba service restarted successfully');
        fetchServiceStatus();
      }
    } catch (error) {
      console.error('Failed to restart Samba:', error);
      alert('Failed to restart Samba service');
    }
  };

  return (
    <div className="flex flex-col h-full bg-white dark:bg-macos-dark-100">
      {/* Header */}
      <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center gap-3">
          <Share2 className="w-6 h-6 text-macos-blue" />
          <div>
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
              Samba Share Manager
            </h1>
            {serviceStatus && (
              <div className="flex items-center gap-2 mt-1">
                {serviceStatus.active ? (
                  <>
                    <CheckCircle className="w-4 h-4 text-green-500" />
                    <span className="text-sm text-green-600 dark:text-green-400">
                      Service Active
                    </span>
                  </>
                ) : (
                  <>
                    <XCircle className="w-4 h-4 text-red-500" />
                    <span className="text-sm text-red-600 dark:text-red-400">
                      Service Inactive
                    </span>
                  </>
                )}
              </div>
            )}
          </div>
        </div>
        <div className="flex gap-2">
          <button
            onClick={handleRestartService}
            className="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
          >
            <Power className="w-4 h-4" />
            Restart Service
          </button>
          <button
            onClick={fetchShares}
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
            Create Share
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Shares List */}
        <div className="w-1/3 border-r border-gray-200 dark:border-gray-700 overflow-y-auto">
          {isLoading ? (
            <div className="flex items-center justify-center p-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
            </div>
          ) : shares.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center">
              <Share2 className="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" />
              <p className="text-gray-500 dark:text-gray-400">No Samba shares configured</p>
              <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
                Create a new share to get started
              </p>
            </div>
          ) : (
            <div className="p-4 space-y-2">
              {shares.map((share) => (
                <motion.div
                  key={share.name}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  onClick={() => setSelectedShare(share)}
                  className={`p-4 rounded-xl cursor-pointer transition-all ${
                    selectedShare?.name === share.name
                      ? 'bg-macos-blue/10 dark:bg-macos-blue/20 border-2 border-macos-blue'
                      : 'bg-gray-50 dark:bg-macos-dark-200 hover:bg-gray-100 dark:hover:bg-macos-dark-300 border-2 border-transparent'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <Folder className="w-5 h-5 text-macos-blue" />
                      <span className="font-semibold text-gray-900 dark:text-gray-100">
                        {share.name}
                      </span>
                    </div>
                    <div className="flex gap-1">
                      {share.readOnly && (
                        <Lock className="w-4 h-4 text-yellow-500" />
                      )}
                      {share.guestOK && (
                        <Users className="w-4 h-4 text-green-500" />
                      )}
                      {!share.browseable && (
                        <EyeOff className="w-4 h-4 text-gray-400" />
                      )}
                    </div>
                  </div>

                  <div className="space-y-1 text-xs">
                    <div className="text-gray-600 dark:text-gray-400 font-mono truncate">
                      {share.path}
                    </div>
                    {share.comment && (
                      <div className="text-gray-500 dark:text-gray-500 italic">
                        {share.comment}
                      </div>
                    )}
                  </div>
                </motion.div>
              ))}
            </div>
          )}
        </div>

        {/* Share Details */}
        <div className="flex-1 overflow-y-auto">
          {selectedShare ? (
            <div className="p-6">
              {/* Share Info Card */}
              <div className="bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 rounded-2xl p-6 mb-6">
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                      {selectedShare.name}
                    </h2>
                    <p className="text-sm text-gray-600 dark:text-gray-400 font-mono mt-1">
                      {selectedShare.path}
                    </p>
                  </div>
                  <button
                    onClick={() => setShowDeleteConfirm(true)}
                    className="flex items-center gap-2 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                    Delete
                  </button>
                </div>

                {selectedShare.comment && (
                  <p className="text-gray-700 dark:text-gray-300 italic mb-4">
                    {selectedShare.comment}
                  </p>
                )}

                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Lock className="w-4 h-4 text-macos-blue" />
                      <span className="text-xs text-gray-600 dark:text-gray-400">Access Mode</span>
                    </div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedShare.readOnly ? 'Read-Only' : 'Read-Write'}
                    </div>
                  </div>

                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-2">
                      {selectedShare.browseable ? (
                        <Eye className="w-4 h-4 text-macos-blue" />
                      ) : (
                        <EyeOff className="w-4 h-4 text-gray-400" />
                      )}
                      <span className="text-xs text-gray-600 dark:text-gray-400">Visibility</span>
                    </div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedShare.browseable ? 'Browseable' : 'Hidden'}
                    </div>
                  </div>

                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Users className="w-4 h-4 text-macos-blue" />
                      <span className="text-xs text-gray-600 dark:text-gray-400">Guest Access</span>
                    </div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedShare.guestOK ? 'Enabled' : 'Disabled'}
                    </div>
                  </div>

                  <div className="bg-white/50 dark:bg-macos-dark-200/50 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Shield className="w-4 h-4 text-macos-blue" />
                      <span className="text-xs text-gray-600 dark:text-gray-400">Recycle Bin</span>
                    </div>
                    <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
                      {selectedShare.recycleBin ? 'Enabled' : 'Disabled'}
                    </div>
                  </div>
                </div>
              </div>

              {/* Access Control */}
              <div className="mb-6">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                  Access Control
                </h3>
                <div className="space-y-4">
                  {/* Valid Users */}
                  <div className="bg-gray-50 dark:bg-macos-dark-200 rounded-lg p-4">
                    <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
                      Valid Users
                    </h4>
                    {selectedShare.validUsers.length > 0 ? (
                      <div className="flex flex-wrap gap-2">
                        {selectedShare.validUsers.map((user) => (
                          <span
                            key={user}
                            className="px-3 py-1 bg-macos-blue/10 text-macos-blue rounded-full text-sm"
                          >
                            {user}
                          </span>
                        ))}
                      </div>
                    ) : (
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        No user restrictions
                      </p>
                    )}
                  </div>

                  {/* Valid Groups */}
                  <div className="bg-gray-50 dark:bg-macos-dark-200 rounded-lg p-4">
                    <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
                      Valid Groups
                    </h4>
                    {selectedShare.validGroups.length > 0 ? (
                      <div className="flex flex-wrap gap-2">
                        {selectedShare.validGroups.map((group) => (
                          <span
                            key={group}
                            className="px-3 py-1 bg-macos-purple/10 text-macos-purple rounded-full text-sm"
                          >
                            {group}
                          </span>
                        ))}
                      </div>
                    ) : (
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        No group restrictions
                      </p>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center h-full text-center p-12">
              <Share2 className="w-24 h-24 text-gray-300 dark:text-gray-600 mb-4" />
              <p className="text-lg text-gray-500 dark:text-gray-400">
                Select a share to view details
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Create Share Dialog */}
      <AnimatePresence>
        {showCreateDialog && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-2xl w-full max-h-[90vh] overflow-y-auto"
            >
              <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-6">
                Create Samba Share
              </h3>

              <div className="space-y-4">
                {/* Name */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Share Name *
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="w-full px-4 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-200 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-gray-100"
                    placeholder="myshare"
                  />
                </div>

                {/* Path */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Path *
                  </label>
                  <input
                    type="text"
                    value={formData.path}
                    onChange={(e) => setFormData({ ...formData, path: e.target.value })}
                    className="w-full px-4 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-200 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-gray-100"
                    placeholder="/mnt/storage/share"
                  />
                </div>

                {/* Comment */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Comment
                  </label>
                  <input
                    type="text"
                    value={formData.comment}
                    onChange={(e) => setFormData({ ...formData, comment: e.target.value })}
                    className="w-full px-4 py-2 bg-gray-50 dark:bg-macos-dark-200 border border-gray-200 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-macos-blue focus:border-transparent text-gray-900 dark:text-gray-100"
                    placeholder="Description of this share"
                  />
                </div>

                {/* Checkboxes */}
                <div className="grid grid-cols-2 gap-4">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={formData.read_only}
                      onChange={(e) => setFormData({ ...formData, read_only: e.target.checked })}
                      className="w-4 h-4 text-macos-blue rounded focus:ring-macos-blue"
                    />
                    <span className="text-sm text-gray-700 dark:text-gray-300">Read-Only</span>
                  </label>

                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={formData.browseable}
                      onChange={(e) => setFormData({ ...formData, browseable: e.target.checked })}
                      className="w-4 h-4 text-macos-blue rounded focus:ring-macos-blue"
                    />
                    <span className="text-sm text-gray-700 dark:text-gray-300">Browseable</span>
                  </label>

                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={formData.guest_ok}
                      onChange={(e) => setFormData({ ...formData, guest_ok: e.target.checked })}
                      className="w-4 h-4 text-macos-blue rounded focus:ring-macos-blue"
                    />
                    <span className="text-sm text-gray-700 dark:text-gray-300">Guest Access</span>
                  </label>

                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={formData.recycle_bin}
                      onChange={(e) => setFormData({ ...formData, recycle_bin: e.target.checked })}
                      className="w-4 h-4 text-macos-blue rounded focus:ring-macos-blue"
                    />
                    <span className="text-sm text-gray-700 dark:text-gray-300">Recycle Bin</span>
                  </label>
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
                  onClick={handleCreateShare}
                  className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                  Create Share
                </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* Delete Confirmation Dialog */}
      <AnimatePresence>
        {showDeleteConfirm && selectedShare && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="bg-white dark:bg-macos-dark-100 rounded-2xl p-6 max-w-md w-full"
            >
              <div className="flex items-center gap-3 mb-4">
                <AlertCircle className="w-6 h-6 text-red-500" />
                <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                  Delete Share
                </h3>
              </div>
              <p className="text-gray-600 dark:text-gray-400 mb-6">
                Are you sure you want to delete the share <strong>"{selectedShare.name}"</strong>?
                This action cannot be undone.
              </p>
              <div className="flex justify-end gap-2">
                <button
                  onClick={() => setShowDeleteConfirm(false)}
                  className="px-4 py-2 bg-gray-100 dark:bg-macos-dark-200 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-macos-dark-300 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleDeleteShare}
                  className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
                >
                  Delete
                </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </div>
  );
}
