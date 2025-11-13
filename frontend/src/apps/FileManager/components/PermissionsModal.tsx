import React, { useState, useEffect } from 'react';
import { FileInfo } from '@/api/files';
import { getFilePermissions, changeFilePermissions, PermissionsInfo } from '@/api/files';
import { useAuthStore } from '@/stores/authStore';

interface PermissionsModalProps {
  file: FileInfo;
  onClose: () => void;
  onSuccess: () => void;
}

// Permission bits for easier manipulation
interface PermissionBits {
  owner: { read: boolean; write: boolean; execute: boolean };
  group: { read: boolean; write: boolean; execute: boolean };
  others: { read: boolean; write: boolean; execute: boolean };
}

// Parse Unix permissions string (e.g., "rwxr-xr-x") to bits
const parsePermissions = (perms: string): PermissionBits => {
  if (perms.length < 9) {
    perms = perms.padEnd(9, '-');
  }

  return {
    owner: {
      read: perms[0] === 'r',
      write: perms[1] === 'w',
      execute: perms[2] === 'x',
    },
    group: {
      read: perms[3] === 'r',
      write: perms[4] === 'w',
      execute: perms[5] === 'x',
    },
    others: {
      read: perms[6] === 'r',
      write: perms[7] === 'w',
      execute: perms[8] === 'x',
    },
  };
};

// Convert permission bits to Unix permissions string
const bitsToPermString = (bits: PermissionBits): string => {
  const ownerStr =
    (bits.owner.read ? 'r' : '-') +
    (bits.owner.write ? 'w' : '-') +
    (bits.owner.execute ? 'x' : '-');
  const groupStr =
    (bits.group.read ? 'r' : '-') +
    (bits.group.write ? 'w' : '-') +
    (bits.group.execute ? 'x' : '-');
  const othersStr =
    (bits.others.read ? 'r' : '-') +
    (bits.others.write ? 'w' : '-') +
    (bits.others.execute ? 'x' : '-');

  return ownerStr + groupStr + othersStr;
};

// Convert permission bits to octal notation (e.g., "755")
const bitsToOctal = (bits: PermissionBits): string => {
  const ownerValue = (bits.owner.read ? 4 : 0) + (bits.owner.write ? 2 : 0) + (bits.owner.execute ? 1 : 0);
  const groupValue = (bits.group.read ? 4 : 0) + (bits.group.write ? 2 : 0) + (bits.group.execute ? 1 : 0);
  const othersValue = (bits.others.read ? 4 : 0) + (bits.others.write ? 2 : 0) + (bits.others.execute ? 1 : 0);

  return `${ownerValue}${groupValue}${othersValue}`;
};

const PermissionsModal: React.FC<PermissionsModalProps> = ({ file, onClose, onSuccess }) => {
  const { user } = useAuthStore();
  const isAdmin = user?.role === 'admin';

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [permInfo, setPermInfo] = useState<PermissionsInfo | null>(null);
  const [permissions, setPermissions] = useState<PermissionBits>({
    owner: { read: true, write: true, execute: false },
    group: { read: true, write: false, execute: false },
    others: { read: true, write: false, execute: false },
  });
  const [recursive, setRecursive] = useState(false);
  const [owner, setOwner] = useState('');
  const [group, setGroup] = useState('');

  useEffect(() => {
    loadPermissions();
  }, [file.path]);

  const loadPermissions = async () => {
    try {
      setLoading(true);
      setError(null);
      const info = await getFilePermissions(file.path);
      setPermInfo(info);
      setPermissions(parsePermissions(info.permissions));
      setOwner(info.owner);
      setGroup(info.group);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load permissions');
    } finally {
      setLoading(false);
    }
  };

  const handleToggle = (category: 'owner' | 'group' | 'others', permission: 'read' | 'write' | 'execute') => {
    setPermissions(prev => ({
      ...prev,
      [category]: {
        ...prev[category],
        [permission]: !prev[category][permission],
      },
    }));
  };

  const handleSave = async () => {
    try {
      setSaving(true);
      setError(null);

      const permString = bitsToPermString(permissions);

      await changeFilePermissions(
        file.path,
        permString,
        isAdmin ? owner : undefined,
        isAdmin ? group : undefined,
        recursive
      );

      onSuccess();
      onClose();
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to change permissions');
    } finally {
      setSaving(false);
    }
  };

  const presetPermissions = (preset: string) => {
    switch (preset) {
      case '777':
        setPermissions({
          owner: { read: true, write: true, execute: true },
          group: { read: true, write: true, execute: true },
          others: { read: true, write: true, execute: true },
        });
        break;
      case '755':
        setPermissions({
          owner: { read: true, write: true, execute: true },
          group: { read: true, write: false, execute: true },
          others: { read: true, write: false, execute: true },
        });
        break;
      case '644':
        setPermissions({
          owner: { read: true, write: true, execute: false },
          group: { read: true, write: false, execute: false },
          others: { read: true, write: false, execute: false },
        });
        break;
      case '600':
        setPermissions({
          owner: { read: true, write: true, execute: false },
          group: { read: false, write: false, execute: false },
          others: { read: false, write: false, execute: false },
        });
        break;
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div
        className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">File Permissions</h2>

        <div className="mb-6">
          <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">File: {file.name}</p>
          <p className="text-xs text-gray-500 dark:text-gray-400">Path: {file.path}</p>
        </div>

        {loading && (
          <div className="text-center py-8">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
            <p className="mt-2 text-gray-600 dark:text-gray-400">Loading permissions...</p>
          </div>
        )}

        {error && (
          <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded">
            <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
          </div>
        )}

        {!loading && permInfo && (
          <>
            {/* Current Info */}
            <div className="mb-6 p-4 bg-gray-50 dark:bg-gray-900 rounded">
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div>
                  <span className="text-gray-500 dark:text-gray-400">Current:</span>{' '}
                  <span className="font-mono text-gray-900 dark:text-white">{permInfo.permissions}</span>
                </div>
                <div>
                  <span className="text-gray-500 dark:text-gray-400">Octal:</span>{' '}
                  <span className="font-mono text-gray-900 dark:text-white">{permInfo.mode}</span>
                </div>
                <div>
                  <span className="text-gray-500 dark:text-gray-400">Owner:</span>{' '}
                  <span className="text-gray-900 dark:text-white">{permInfo.owner}</span>
                </div>
                <div>
                  <span className="text-gray-500 dark:text-gray-400">Group:</span>{' '}
                  <span className="text-gray-900 dark:text-white">{permInfo.group}</span>
                </div>
              </div>
            </div>

            {/* Preset Buttons */}
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Quick Presets:</label>
              <div className="flex gap-2 flex-wrap">
                <button
                  onClick={() => presetPermissions('644')}
                  className="px-3 py-1 text-xs bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded"
                >
                  644 (rw-r--r--)
                </button>
                <button
                  onClick={() => presetPermissions('755')}
                  className="px-3 py-1 text-xs bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded"
                >
                  755 (rwxr-xr-x)
                </button>
                <button
                  onClick={() => presetPermissions('600')}
                  className="px-3 py-1 text-xs bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded"
                >
                  600 (rw-------)
                </button>
                <button
                  onClick={() => presetPermissions('777')}
                  className="px-3 py-1 text-xs bg-yellow-100 dark:bg-yellow-900/30 hover:bg-yellow-200 dark:hover:bg-yellow-900/50 rounded"
                >
                  777 (rwxrwxrwx) ⚠️
                </button>
              </div>
            </div>

            {/* Permission Grid */}
            <div className="mb-6 overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-200 dark:border-gray-700">
                    <th className="text-left py-2 px-3 text-gray-700 dark:text-gray-300">User Type</th>
                    <th className="text-center py-2 px-3 text-gray-700 dark:text-gray-300">Read</th>
                    <th className="text-center py-2 px-3 text-gray-700 dark:text-gray-300">Write</th>
                    <th className="text-center py-2 px-3 text-gray-700 dark:text-gray-300">Execute</th>
                  </tr>
                </thead>
                <tbody>
                  {/* Owner */}
                  <tr className="border-b border-gray-100 dark:border-gray-800">
                    <td className="py-3 px-3 font-medium text-gray-900 dark:text-white">Owner</td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.owner.read}
                        onChange={() => handleToggle('owner', 'read')}
                        className="w-4 h-4"
                      />
                    </td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.owner.write}
                        onChange={() => handleToggle('owner', 'write')}
                        className="w-4 h-4"
                      />
                    </td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.owner.execute}
                        onChange={() => handleToggle('owner', 'execute')}
                        className="w-4 h-4"
                      />
                    </td>
                  </tr>

                  {/* Group */}
                  <tr className="border-b border-gray-100 dark:border-gray-800">
                    <td className="py-3 px-3 font-medium text-gray-900 dark:text-white">Group</td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.group.read}
                        onChange={() => handleToggle('group', 'read')}
                        className="w-4 h-4"
                      />
                    </td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.group.write}
                        onChange={() => handleToggle('group', 'write')}
                        className="w-4 h-4"
                      />
                    </td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.group.execute}
                        onChange={() => handleToggle('group', 'execute')}
                        className="w-4 h-4"
                      />
                    </td>
                  </tr>

                  {/* Others */}
                  <tr>
                    <td className="py-3 px-3 font-medium text-gray-900 dark:text-white">Others</td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.others.read}
                        onChange={() => handleToggle('others', 'read')}
                        className="w-4 h-4"
                      />
                    </td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.others.write}
                        onChange={() => handleToggle('others', 'write')}
                        className="w-4 h-4"
                      />
                    </td>
                    <td className="text-center">
                      <input
                        type="checkbox"
                        checked={permissions.others.execute}
                        onChange={() => handleToggle('others', 'execute')}
                        className="w-4 h-4"
                      />
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            {/* Preview */}
            <div className="mb-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded">
              <p className="text-sm text-gray-700 dark:text-gray-300 mb-1">New Permissions:</p>
              <p className="font-mono text-lg text-gray-900 dark:text-white">
                {bitsToPermString(permissions)} ({bitsToOctal(permissions)})
              </p>
            </div>

            {/* Owner/Group (Admin only) */}
            {isAdmin && (
              <div className="mb-6 grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Owner
                  </label>
                  <input
                    type="text"
                    value={owner}
                    onChange={(e) => setOwner(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                    placeholder="username"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Group
                  </label>
                  <input
                    type="text"
                    value={group}
                    onChange={(e) => setGroup(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                    placeholder="groupname"
                  />
                </div>
              </div>
            )}

            {/* Recursive option (for directories) */}
            {file.isDir && (
              <div className="mb-6">
                <label className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={recursive}
                    onChange={(e) => setRecursive(e.target.checked)}
                    className="w-4 h-4"
                  />
                  <span className="text-sm text-gray-700 dark:text-gray-300">
                    Apply recursively to all files and subdirectories
                  </span>
                </label>
              </div>
            )}
          </>
        )}

        {/* Actions */}
        <div className="flex justify-end gap-2">
          <button
            onClick={onClose}
            disabled={saving}
            className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={saving || loading}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded disabled:opacity-50 flex items-center gap-2"
          >
            {saving ? (
              <>
                <div className="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
                Saving...
              </>
            ) : (
              'Save Changes'
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default PermissionsModal;
