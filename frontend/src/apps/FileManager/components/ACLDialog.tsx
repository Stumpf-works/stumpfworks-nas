import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { filesystemACLApi, type ACLEntry } from '@/api/filesystem-acl';
import toast from 'react-hot-toast';

interface ACLDialogProps {
  path: string;
  isDirectory: boolean;
  onClose: () => void;
}

export function ACLDialog({ path, isDirectory, onClose }: ACLDialogProps) {
  const [entries, setEntries] = useState<ACLEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [newEntry, setNewEntry] = useState<ACLEntry>({
    type: 'user',
    name: '',
    permissions: 'r-x',
  });

  useEffect(() => {
    loadACLs();
  }, [path]);

  const loadACLs = async () => {
    try {
      setLoading(true);
      const response = await filesystemACLApi.getACL(path);
      if (response.success && response.data) {
        setEntries(response.data.entries);
      }
    } catch (error: any) {
      toast.error(`Failed to load ACLs: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleAddEntry = async () => {
    if (!newEntry.name && newEntry.type !== 'other' && newEntry.type !== 'mask') {
      toast.error('Please enter a name for this ACL entry');
      return;
    }

    try {
      const response = await filesystemACLApi.setACL({
        path,
        entries: [newEntry],
      });

      if (response.success) {
        toast.success('ACL entry added successfully');
        setNewEntry({ type: 'user', name: '', permissions: 'r-x' });
        loadACLs();
      }
    } catch (error: any) {
      toast.error(`Failed to add ACL entry: ${error.message}`);
    }
  };

  const handleRemoveEntry = async (entry: ACLEntry) => {
    try {
      const response = await filesystemACLApi.removeACL({
        path,
        type: entry.type,
        name: entry.name,
      });

      if (response.success) {
        toast.success('ACL entry removed successfully');
        loadACLs();
      }
    } catch (error: any) {
      toast.error(`Failed to remove ACL entry: ${error.message}`);
    }
  };

  const handleRemoveAll = async () => {
    if (!confirm('Remove all ACL entries? This will revert to standard Unix permissions.')) {
      return;
    }

    try {
      const response = await filesystemACLApi.removeAllACLs(path);

      if (response.success) {
        toast.success('All ACL entries removed');
        loadACLs();
      }
    } catch (error: any) {
      toast.error(`Failed to remove ACLs: ${error.message}`);
    }
  };

  const handleSetDefault = async () => {
    if (!isDirectory) {
      toast.error('Default ACLs can only be set on directories');
      return;
    }

    if (!confirm('Set current ACLs as default for new files in this directory?')) {
      return;
    }

    try {
      const response = await filesystemACLApi.setDefaultACL({
        dir_path: path,
        entries: entries,
      });

      if (response.success) {
        toast.success('Default ACLs set successfully');
      }
    } catch (error: any) {
      toast.error(`Failed to set default ACLs: ${error.message}`);
    }
  };

  const handleApplyRecursive = async () => {
    if (!isDirectory) {
      toast.error('Recursive ACL application can only be done on directories');
      return;
    }

    if (!confirm('Apply current ACLs recursively to all files and subdirectories? This may take a while.')) {
      return;
    }

    try {
      const response = await filesystemACLApi.applyRecursive({
        dir_path: path,
        entries: entries,
      });

      if (response.success) {
        toast.success('ACLs applied recursively');
      }
    } catch (error: any) {
      toast.error(`Failed to apply ACLs recursively: ${error.message}`);
    }
  };

  const getPermissionDisplay = (perms: string) => {
    return perms.split('').map((p, i) => (
      <span key={i} className={p !== '-' ? 'text-green-600 dark:text-green-400' : 'text-gray-400'}>
        {p}
      </span>
    ));
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'user': return 'User';
      case 'group': return 'Group';
      case 'mask': return 'Mask';
      case 'other': return 'Other';
      default: return type;
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50" onClick={onClose}>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-xl shadow-2xl w-full max-w-3xl max-h-[90vh] overflow-hidden flex flex-col"
      >
        {/* Header */}
        <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Manage ACLs
              </h2>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                {path}
              </p>
            </div>
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
            >
              ✕
            </button>
          </div>
        </div>

        {/* Body */}
        <div className="flex-1 overflow-y-auto p-6">
          {loading ? (
            <div className="flex items-center justify-center h-40">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-macos-blue" />
            </div>
          ) : (
            <>
              {/* Current ACLs */}
              <div className="mb-6">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-3">
                  Current ACL Entries
                </h3>
                {entries.length === 0 ? (
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    No ACL entries set. Using standard Unix permissions.
                  </p>
                ) : (
                  <div className="space-y-2">
                    {entries.map((entry, index) => (
                      <div
                        key={index}
                        className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
                      >
                        <div className="flex items-center gap-4">
                          <span className="text-sm font-medium text-gray-700 dark:text-gray-300 w-16">
                            {getTypeLabel(entry.type)}
                          </span>
                          <span className="text-sm text-gray-600 dark:text-gray-400 w-32">
                            {entry.name || <em>(owner)</em>}
                          </span>
                          <code className="text-sm font-mono">
                            {getPermissionDisplay(entry.permissions)}
                          </code>
                        </div>
                        <button
                          onClick={() => handleRemoveEntry(entry)}
                          className="px-3 py-1 text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20 rounded transition-colors"
                        >
                          Remove
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Add New Entry */}
              <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-3">
                  Add ACL Entry
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-4 gap-3 mb-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Type
                    </label>
                    <select
                      value={newEntry.type}
                      onChange={(e) => setNewEntry({ ...newEntry, type: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                    >
                      <option value="user">User</option>
                      <option value="group">Group</option>
                      <option value="mask">Mask</option>
                      <option value="other">Other</option>
                    </select>
                  </div>
                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Name
                    </label>
                    <input
                      type="text"
                      value={newEntry.name}
                      onChange={(e) => setNewEntry({ ...newEntry, name: e.target.value })}
                      placeholder={newEntry.type === 'other' || newEntry.type === 'mask' ? 'Leave empty' : 'username or groupname'}
                      disabled={newEntry.type === 'other' || newEntry.type === 'mask'}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 disabled:opacity-50"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Permissions
                    </label>
                    <select
                      value={newEntry.permissions}
                      onChange={(e) => setNewEntry({ ...newEntry, permissions: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                    >
                      <option value="rwx">rwx (7)</option>
                      <option value="rw-">rw- (6)</option>
                      <option value="r-x">r-x (5)</option>
                      <option value="r--">r-- (4)</option>
                      <option value="-wx">-wx (3)</option>
                      <option value="-w-">-w- (2)</option>
                      <option value="--x">--x (1)</option>
                      <option value="---">--- (0)</option>
                    </select>
                  </div>
                </div>
                <button
                  onClick={handleAddEntry}
                  className="px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                  + Add Entry
                </button>
              </div>
            </>
          )}
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
          <div className="flex items-center justify-between">
            <div className="flex gap-2">
              {isDirectory && (
                <>
                  <button
                    onClick={handleSetDefault}
                    className="px-4 py-2 text-sm bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-lg hover:bg-blue-200 dark:hover:bg-blue-900/50 transition-colors"
                  >
                    Set as Default
                  </button>
                  <button
                    onClick={handleApplyRecursive}
                    className="px-4 py-2 text-sm bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 rounded-lg hover:bg-purple-200 dark:hover:bg-purple-900/50 transition-colors"
                  >
                    Apply Recursively
                  </button>
                </>
              )}
              <button
                onClick={handleRemoveAll}
                className="px-4 py-2 text-sm bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 rounded-lg hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors"
              >
                Remove All
              </button>
            </div>
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
