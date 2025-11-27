import { useState, useEffect } from 'react';
import { addcApi, ADGroup } from '../../../api/addc';
import { motion, AnimatePresence } from 'framer-motion';
import { Users, Plus, Trash2, RefreshCw, AlertCircle, UserPlus, UserMinus } from 'lucide-react';

export default function GroupManagement() {
  const [groups, setGroups] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showMembersModal, setShowMembersModal] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<string>('');
  const [groupMembers, setGroupMembers] = useState<string[]>([]);
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [newMemberUsername, setNewMemberUsername] = useState('');

  const [createForm, setCreateForm] = useState<ADGroup>({
    name: '',
    description: '',
    ou: '',
    group_scope: 'Global',
    group_type: 'Security',
  });

  const loadGroups = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await addcApi.listGroups();
      if (response.success && response.data) {
        setGroups(response.data);
      } else {
        setError(response.error?.message || 'Failed to load groups');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load groups');
    } finally {
      setLoading(false);
    }
  };

  const loadGroupMembers = async (groupName: string) => {
    try {
      setActionLoading(`members-${groupName}`);
      const response = await addcApi.listGroupMembers(groupName);
      if (response.success && response.data) {
        setGroupMembers(response.data);
        setSelectedGroup(groupName);
        setShowMembersModal(true);
      } else {
        setError(response.error?.message || 'Failed to load group members');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load group members');
    } finally {
      setActionLoading(null);
    }
  };

  useEffect(() => {
    loadGroups();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!createForm.name) {
      setError('Group name is required');
      return;
    }

    try {
      setActionLoading('create');
      setError('');
      const response = await addcApi.createGroup(createForm);

      if (response.success) {
        alert(`Group ${createForm.name} created successfully!`);
        setShowCreateForm(false);
        setCreateForm({
          name: '',
          description: '',
          ou: '',
          group_scope: 'Global',
          group_type: 'Security',
        });
        loadGroups();
      } else {
        setError(response.error?.message || 'Failed to create group');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create group');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async (name: string) => {
    if (!confirm(`Are you sure you want to delete group "${name}"?`)) {
      return;
    }

    try {
      setActionLoading(`delete-${name}`);
      setError('');
      const response = await addcApi.deleteGroup(name);

      if (response.success) {
        alert(`Group ${name} deleted successfully`);
        loadGroups();
      } else {
        setError(response.error?.message || 'Failed to delete group');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to delete group');
    } finally {
      setActionLoading(null);
    }
  };

  const handleAddMember = async () => {
    if (!newMemberUsername.trim()) {
      setError('Username is required');
      return;
    }

    try {
      setActionLoading(`add-member-${selectedGroup}`);
      setError('');
      const response = await addcApi.addGroupMember(selectedGroup, newMemberUsername);

      if (response.success) {
        alert(`User ${newMemberUsername} added to group ${selectedGroup}`);
        setNewMemberUsername('');
        loadGroupMembers(selectedGroup);
      } else {
        setError(response.error?.message || 'Failed to add member');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to add member');
    } finally {
      setActionLoading(null);
    }
  };

  const handleRemoveMember = async (username: string) => {
    if (!confirm(`Remove ${username} from group ${selectedGroup}?`)) {
      return;
    }

    try {
      setActionLoading(`remove-member-${username}`);
      setError('');
      const response = await addcApi.removeGroupMember(selectedGroup, username);

      if (response.success) {
        alert(`User ${username} removed from group ${selectedGroup}`);
        loadGroupMembers(selectedGroup);
      } else {
        setError(response.error?.message || 'Failed to remove member');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to remove member');
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Users className="w-6 h-6 text-macos-blue" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">Group Management</h2>
        </div>
        <div className="flex gap-3">
          <button
            onClick={loadGroups}
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
            Create Group
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

      {/* Groups List */}
      {loading && groups.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <RefreshCw className="w-6 h-6 animate-spin text-gray-400" />
        </div>
      ) : groups.length === 0 ? (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          No groups found
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {groups.map((groupName) => (
            <div
              key={groupName}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-2">
                  <Users className="w-5 h-5 text-macos-blue" />
                  <span className="font-medium text-gray-900 dark:text-gray-100">{groupName}</span>
                </div>
              </div>

              <div className="flex flex-wrap gap-2">
                <button
                  onClick={() => loadGroupMembers(groupName)}
                  disabled={actionLoading === `members-${groupName}`}
                  className="p-1.5 text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded transition-colors disabled:opacity-50"
                  title="Manage Members"
                >
                  {actionLoading === `members-${groupName}` ? (
                    <RefreshCw className="w-4 h-4 animate-spin" />
                  ) : (
                    <UserPlus className="w-4 h-4" />
                  )}
                </button>
                <button
                  onClick={() => handleDelete(groupName)}
                  disabled={actionLoading === `delete-${groupName}`}
                  className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors disabled:opacity-50"
                  title="Delete Group"
                >
                  {actionLoading === `delete-${groupName}` ? (
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
                  Create New Group
                </h2>

                <form onSubmit={handleCreate} className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Group Name <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={createForm.name}
                        onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                        required
                      />
                    </div>

                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Description
                      </label>
                      <input
                        type="text"
                        value={createForm.description}
                        onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Group Scope
                      </label>
                      <select
                        value={createForm.group_scope}
                        onChange={(e) => setCreateForm({ ...createForm, group_scope: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      >
                        <option value="Global">Global</option>
                        <option value="Domain">Domain Local</option>
                        <option value="Universal">Universal</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Group Type
                      </label>
                      <select
                        value={createForm.group_type}
                        onChange={(e) => setCreateForm({ ...createForm, group_type: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      >
                        <option value="Security">Security</option>
                        <option value="Distribution">Distribution</option>
                      </select>
                    </div>

                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Organizational Unit (OU)
                      </label>
                      <input
                        type="text"
                        value={createForm.ou}
                        onChange={(e) => setCreateForm({ ...createForm, ou: e.target.value })}
                        placeholder="OU=Groups,DC=example,DC=com"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>
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
                          Create Group
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

      {/* Members Modal */}
      <AnimatePresence>
        {showMembersModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => !actionLoading && setShowMembersModal(false)}
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
                  Manage Members - {selectedGroup}
                </h2>

                {/* Add Member */}
                <div className="mb-6 p-4 bg-gray-50 dark:bg-macos-dark-50 rounded-lg">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Add Member
                  </label>
                  <div className="flex gap-2">
                    <input
                      type="text"
                      value={newMemberUsername}
                      onChange={(e) => setNewMemberUsername(e.target.value)}
                      placeholder="Username"
                      className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-100 text-gray-900 dark:text-gray-100"
                    />
                    <button
                      onClick={handleAddMember}
                      disabled={actionLoading === `add-member-${selectedGroup}`}
                      className="flex items-center gap-2 px-4 py-2 bg-macos-blue text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
                    >
                      {actionLoading === `add-member-${selectedGroup}` ? (
                        <RefreshCw className="w-4 h-4 animate-spin" />
                      ) : (
                        <UserPlus className="w-4 h-4" />
                      )}
                      Add
                    </button>
                  </div>
                </div>

                {/* Members List */}
                <div className="space-y-2">
                  <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Current Members ({groupMembers.length})
                  </h3>
                  {groupMembers.length === 0 ? (
                    <p className="text-gray-500 dark:text-gray-400 text-sm text-center py-4">
                      No members in this group
                    </p>
                  ) : (
                    <div className="max-h-64 overflow-y-auto space-y-2">
                      {groupMembers.map((member) => (
                        <div
                          key={member}
                          className="flex items-center justify-between p-3 bg-gray-50 dark:bg-macos-dark-50 rounded-lg"
                        >
                          <span className="text-gray-900 dark:text-gray-100">{member}</span>
                          <button
                            onClick={() => handleRemoveMember(member)}
                            disabled={actionLoading === `remove-member-${member}`}
                            className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors disabled:opacity-50"
                            title="Remove from group"
                          >
                            {actionLoading === `remove-member-${member}` ? (
                              <RefreshCw className="w-4 h-4 animate-spin" />
                            ) : (
                              <UserMinus className="w-4 h-4" />
                            )}
                          </button>
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                <div className="flex justify-end gap-3 mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
                  <button
                    onClick={() => {
                      setShowMembersModal(false);
                      setNewMemberUsername('');
                    }}
                    disabled={!!actionLoading}
                    className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors disabled:opacity-50"
                  >
                    Close
                  </button>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
