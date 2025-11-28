import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { quotaApi, type QuotaInfo, type FilesystemQuotaStatus } from '@/api/quota';
import toast from 'react-hot-toast';

type QuotaTab = 'users' | 'groups';

interface QuotaDialogData {
  type: 'user' | 'group';
  name?: string;
  existing?: QuotaInfo;
}

export function QuotaManager() {
  const [activeTab, setActiveTab] = useState<QuotaTab>('users');
  const [filesystem, setFilesystem] = useState<string>('/');
  const [userQuotas, setUserQuotas] = useState<QuotaInfo[]>([]);
  const [groupQuotas, setGroupQuotas] = useState<QuotaInfo[]>([]);
  const [quotaStatus, setQuotaStatus] = useState<FilesystemQuotaStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [dialogData, setDialogData] = useState<QuotaDialogData | null>(null);

  const tabs = [
    { id: 'users' as QuotaTab, name: 'User Quotas', icon: '👤' },
    { id: 'groups' as QuotaTab, name: 'Group Quotas', icon: '👥' },
  ];

  useEffect(() => {
    loadQuotaStatus();
  }, [filesystem]);

  useEffect(() => {
    if (quotaStatus?.quotasEnabled) {
      if (activeTab === 'users') {
        loadUserQuotas();
      } else {
        loadGroupQuotas();
      }
    }
  }, [activeTab, filesystem, quotaStatus]);

  const loadQuotaStatus = async () => {
    try {
      setLoading(true);
      const response = await quotaApi.getQuotaStatus(filesystem);
      if (response.success && response.data) {
        setQuotaStatus(response.data);
        if (!response.data.quotasEnabled) {
          toast.error(`Quotas are not enabled on ${filesystem}`);
        }
      }
    } catch (error: any) {
      toast.error(`Failed to check quota status: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const loadUserQuotas = async () => {
    try {
      setLoading(true);
      const response = await quotaApi.listUserQuotas(filesystem);
      if (response.success && response.data) {
        setUserQuotas(response.data);
      }
    } catch (error: any) {
      toast.error(`Failed to load user quotas: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const loadGroupQuotas = async () => {
    try {
      setLoading(true);
      const response = await quotaApi.listGroupQuotas(filesystem);
      if (response.success && response.data) {
        setGroupQuotas(response.data);
      }
    } catch (error: any) {
      toast.error(`Failed to load group quotas: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleSetQuota = (quota?: QuotaInfo) => {
    if (quota) {
      setDialogData({ type: activeTab === 'users' ? 'user' : 'group', existing: quota });
    } else {
      setDialogData({ type: activeTab === 'users' ? 'user' : 'group' });
    }
  };

  const handleRemoveQuota = async (quota: QuotaInfo) => {
    if (!confirm(`Remove quota for ${quota.name}?`)) {
      return;
    }

    try {
      const request = { name: quota.name, filesystem: quota.filesystem };
      const response = quota.type === 'user'
        ? await quotaApi.removeUserQuota(request)
        : await quotaApi.removeGroupQuota(request);

      if (response.success) {
        toast.success(`Quota removed for ${quota.name}`);
        if (activeTab === 'users') {
          loadUserQuotas();
        } else {
          loadGroupQuotas();
        }
      }
    } catch (error: any) {
      toast.error(`Failed to remove quota: ${error.message}`);
    }
  };

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const getUsagePercentage = (used: number, limit: number): number => {
    if (limit === 0) return 0;
    return Math.min((used / limit) * 100, 100);
  };

  const getUsageColor = (percentage: number): string => {
    if (percentage >= 90) return 'bg-red-500';
    if (percentage >= 75) return 'bg-yellow-500';
    return 'bg-green-500';
  };

  const renderQuotaTable = (quotas: QuotaInfo[]) => {
    if (quotas.length === 0) {
      return (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          <div className="text-5xl mb-4">📊</div>
          <p>No quotas set for this filesystem</p>
          <button
            onClick={() => handleSetQuota()}
            className="mt-4 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            + Set First Quota
          </button>
        </div>
      );
    }

    return (
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Name
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Disk Usage
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Block Quota
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                File Count
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                File Quota
              </th>
              <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-macos-dark-100 divide-y divide-gray-200 dark:divide-gray-700">
            {quotas.map((quota) => {
              const blockPercentage = getUsagePercentage(quota.blocksUsed, quota.blocksHard || quota.blocksSoft);
              const filePercentage = getUsagePercentage(quota.filesUsed, quota.filesHard || quota.filesSoft);

              return (
                <tr key={quota.name} className="hover:bg-gray-50 dark:hover:bg-gray-800">
                  <td className="px-4 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <span className="text-2xl mr-3">{activeTab === 'users' ? '👤' : '👥'}</span>
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-gray-100">
                          {quota.name}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400">
                          {quota.filesystem}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-4">
                    <div className="space-y-1">
                      <div className="text-sm text-gray-900 dark:text-gray-100">
                        {formatBytes(quota.blocksUsed * 1024)}
                      </div>
                      {(quota.blocksHard > 0 || quota.blocksSoft > 0) && (
                        <div className="w-32 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                          <div
                            className={`h-2 rounded-full transition-all ${getUsageColor(blockPercentage)}`}
                            style={{ width: `${blockPercentage}%` }}
                          />
                        </div>
                      )}
                      {quota.blocksGrace && quota.blocksGrace !== 'none' && (
                        <div className="text-xs text-red-600 dark:text-red-400">
                          Grace: {quota.blocksGrace}
                        </div>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-4">
                    <div className="text-sm text-gray-900 dark:text-gray-100">
                      {quota.blocksSoft > 0 && (
                        <div>Soft: {formatBytes(quota.blocksSoft * 1024)}</div>
                      )}
                      {quota.blocksHard > 0 && (
                        <div>Hard: {formatBytes(quota.blocksHard * 1024)}</div>
                      )}
                      {quota.blocksSoft === 0 && quota.blocksHard === 0 && (
                        <span className="text-gray-400">No limit</span>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-4">
                    <div className="space-y-1">
                      <div className="text-sm text-gray-900 dark:text-gray-100">
                        {quota.filesUsed.toLocaleString()}
                      </div>
                      {(quota.filesHard > 0 || quota.filesSoft > 0) && (
                        <div className="w-32 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                          <div
                            className={`h-2 rounded-full transition-all ${getUsageColor(filePercentage)}`}
                            style={{ width: `${filePercentage}%` }}
                          />
                        </div>
                      )}
                      {quota.filesGrace && quota.filesGrace !== 'none' && (
                        <div className="text-xs text-red-600 dark:text-red-400">
                          Grace: {quota.filesGrace}
                        </div>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-4">
                    <div className="text-sm text-gray-900 dark:text-gray-100">
                      {quota.filesSoft > 0 && (
                        <div>Soft: {quota.filesSoft.toLocaleString()}</div>
                      )}
                      {quota.filesHard > 0 && (
                        <div>Hard: {quota.filesHard.toLocaleString()}</div>
                      )}
                      {quota.filesSoft === 0 && quota.filesHard === 0 && (
                        <span className="text-gray-400">No limit</span>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap text-right text-sm">
                    <button
                      onClick={() => handleSetQuota(quota)}
                      className="text-macos-blue hover:text-blue-700 dark:hover:text-blue-400 mr-3"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleRemoveQuota(quota)}
                      className="text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-500"
                    >
                      Remove
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    );
  };

  return (
    <div className="h-full flex flex-col bg-gray-50 dark:bg-macos-dark-200">
      {/* Header */}
      <div className="bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Disk Quota Manager
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Manage disk quotas for users and groups
            </p>
          </div>
          <div className="flex items-center gap-4">
            {/* Filesystem Selector */}
            <div>
              <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">
                Filesystem
              </label>
              <select
                value={filesystem}
                onChange={(e) => setFilesystem(e.target.value)}
                className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
              >
                <option value="/">/</option>
                <option value="/home">/home</option>
                <option value="/var">/var</option>
              </select>
            </div>
            <button
              onClick={() => handleSetQuota()}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={!quotaStatus?.quotasEnabled}
            >
              + Set Quota
            </button>
          </div>
        </div>

        {/* Quota Status */}
        {quotaStatus && (
          <div className="mt-4 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <div className="flex items-center gap-4 text-sm">
              <div className="flex items-center gap-2">
                <span className={`w-2 h-2 rounded-full ${quotaStatus.quotasEnabled ? 'bg-green-500' : 'bg-red-500'}`} />
                <span className="text-gray-700 dark:text-gray-300">
                  Quotas: {quotaStatus.quotasEnabled ? 'Enabled' : 'Disabled'}
                </span>
              </div>
              {quotaStatus.quotasEnabled && (
                <>
                  <div className="flex items-center gap-2">
                    <span className={`w-2 h-2 rounded-full ${quotaStatus.userQuotaEnabled ? 'bg-green-500' : 'bg-gray-400'}`} />
                    <span className="text-gray-700 dark:text-gray-300">
                      User Quotas: {quotaStatus.userQuotaEnabled ? 'On' : 'Off'}
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className={`w-2 h-2 rounded-full ${quotaStatus.groupQuotaEnabled ? 'bg-green-500' : 'bg-gray-400'}`} />
                    <span className="text-gray-700 dark:text-gray-300">
                      Group Quotas: {quotaStatus.groupQuotaEnabled ? 'On' : 'Off'}
                    </span>
                  </div>
                </>
              )}
            </div>
          </div>
        )}
      </div>

      {/* Tabs */}
      <div className="bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700">
        <div className="flex space-x-1 px-6">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                activeTab === tab.id
                  ? 'border-macos-blue text-macos-blue'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
              }`}
            >
              <span className="mr-2">{tab.icon}</span>
              {tab.name}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
          </div>
        ) : !quotaStatus?.quotasEnabled ? (
          <div className="text-center py-12 text-gray-500 dark:text-gray-400">
            <div className="text-5xl mb-4">⚠️</div>
            <p className="text-lg font-medium mb-2">Quotas Not Enabled</p>
            <p>Disk quotas are not enabled on {filesystem}</p>
            <p className="text-sm mt-2">Enable quotas in filesystem mount options (usrquota, grpquota)</p>
          </div>
        ) : (
          <motion.div
            key={activeTab}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.2 }}
          >
            {activeTab === 'users' ? renderQuotaTable(userQuotas) : renderQuotaTable(groupQuotas)}
          </motion.div>
        )}
      </div>

      {/* Set Quota Dialog */}
      {dialogData && (
        <QuotaDialog
          type={dialogData.type}
          filesystem={filesystem}
          existing={dialogData.existing}
          onClose={() => setDialogData(null)}
          onSuccess={() => {
            setDialogData(null);
            if (activeTab === 'users') {
              loadUserQuotas();
            } else {
              loadGroupQuotas();
            }
          }}
        />
      )}
    </div>
  );
}

interface QuotaDialogProps {
  type: 'user' | 'group';
  filesystem: string;
  existing?: QuotaInfo;
  onClose: () => void;
  onSuccess: () => void;
}

function QuotaDialog({ type, filesystem, existing, onClose, onSuccess }: QuotaDialogProps) {
  const [name, setName] = useState(existing?.name || '');
  const [blocksSoft, setBlocksSoft] = useState(existing?.blocksSoft ? String(existing.blocksSoft) : '');
  const [blocksHard, setBlocksHard] = useState(existing?.blocksHard ? String(existing.blocksHard) : '');
  const [filesSoft, setFilesSoft] = useState(existing?.filesSoft ? String(existing.filesSoft) : '');
  const [filesHard, setFilesHard] = useState(existing?.filesHard ? String(existing.filesHard) : '');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      toast.error(`Please enter a ${type} name`);
      return;
    }

    try {
      setSubmitting(true);
      const request = {
        name: name.trim(),
        filesystem,
        blocksSoft: blocksSoft ? Number(blocksSoft) : undefined,
        blocksHard: blocksHard ? Number(blocksHard) : undefined,
        filesSoft: filesSoft ? Number(filesSoft) : undefined,
        filesHard: filesHard ? Number(filesHard) : undefined,
      };

      const response = type === 'user'
        ? await quotaApi.setUserQuota(request)
        : await quotaApi.setGroupQuota(request);

      if (response.success) {
        toast.success(`Quota ${existing ? 'updated' : 'set'} for ${name}`);
        onSuccess();
      }
    } catch (error: any) {
      toast.error(`Failed to set quota: ${error.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50" onClick={onClose}>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-xl shadow-2xl w-full max-w-md"
      >
        <form onSubmit={handleSubmit}>
          {/* Header */}
          <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                {existing ? 'Edit' : 'Set'} {type === 'user' ? 'User' : 'Group'} Quota
              </h2>
              <button
                type="button"
                onClick={onClose}
                className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
              >
                ✕
              </button>
            </div>
          </div>

          {/* Body */}
          <div className="p-6 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                {type === 'user' ? 'Username' : 'Group Name'}
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={!!existing}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 disabled:opacity-50"
                placeholder={type === 'user' ? 'Enter username' : 'Enter group name'}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Filesystem
              </label>
              <input
                type="text"
                value={filesystem}
                disabled
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100"
              />
            </div>

            <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
              <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                Block Limits (KB)
              </h3>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">
                    Soft Limit
                  </label>
                  <input
                    type="number"
                    value={blocksSoft}
                    onChange={(e) => setBlocksSoft(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                    placeholder="0"
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">
                    Hard Limit
                  </label>
                  <input
                    type="number"
                    value={blocksHard}
                    onChange={(e) => setBlocksHard(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                    placeholder="0"
                    min="0"
                  />
                </div>
              </div>
            </div>

            <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
              <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                File (Inode) Limits
              </h3>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">
                    Soft Limit
                  </label>
                  <input
                    type="number"
                    value={filesSoft}
                    onChange={(e) => setFilesSoft(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                    placeholder="0"
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">
                    Hard Limit
                  </label>
                  <input
                    type="number"
                    value={filesHard}
                    onChange={(e) => setFilesHard(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                    placeholder="0"
                    min="0"
                  />
                </div>
              </div>
            </div>

            <div className="text-xs text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800 p-3 rounded">
              <strong>Note:</strong> Set to 0 for no limit. Soft limits can be exceeded temporarily, hard limits cannot be exceeded.
            </div>
          </div>

          {/* Footer */}
          <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
            >
              {submitting ? 'Setting...' : existing ? 'Update Quota' : 'Set Quota'}
            </button>
          </div>
        </form>
      </motion.div>
    </div>
  );
}
