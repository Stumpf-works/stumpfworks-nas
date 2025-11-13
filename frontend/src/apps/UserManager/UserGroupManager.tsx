import { useEffect, useState } from 'react';
import { groupsApi, UserGroup, CreateGroupRequest, UpdateGroupRequest } from '@/api/groups';
import { usersApi } from '@/api/users';
import { User } from '@/api/auth';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Card from '@/components/ui/Card';
import { motion, AnimatePresence } from 'framer-motion';
import { Users, UserPlus, UserMinus, Edit2, Trash2 } from 'lucide-react';

export function UserGroupManager() {
  const [groups, setGroups] = useState<UserGroup[]>([]);
  const [allUsers, setAllUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingGroup, setEditingGroup] = useState<UserGroup | null>(null);
  const [managingGroupId, setManagingGroupId] = useState<number | null>(null);

  // Form state
  const [formData, setFormData] = useState<CreateGroupRequest>({
    name: '',
    description: '',
  });

  const loadGroups = async () => {
    setIsLoading(true);
    setError('');
    try {
      const response = await groupsApi.list();
      if (response.success && response.data) {
        setGroups(response.data);
      } else {
        setError(response.error?.message || 'Failed to load groups');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  const loadUsers = async () => {
    try {
      const response = await usersApi.list();
      if (response.success && response.data) {
        setAllUsers(response.data);
      }
    } catch (err) {
      console.error('Failed to load users:', err);
    }
  };

  useEffect(() => {
    loadGroups();
    loadUsers();
  }, []);

  const handleCreate = async () => {
    setError('');
    try {
      const response = await groupsApi.create(formData);
      if (response.success) {
        setShowCreateModal(false);
        resetForm();
        loadGroups();
      } else {
        setError(response.error?.message || 'Failed to create group');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleUpdate = async () => {
    if (!editingGroup) return;
    setError('');
    try {
      const updateData: UpdateGroupRequest = {
        name: formData.name,
        description: formData.description,
      };
      const response = await groupsApi.update(editingGroup.id, updateData);
      if (response.success) {
        setEditingGroup(null);
        resetForm();
        loadGroups();
      } else {
        setError(response.error?.message || 'Failed to update group');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleDelete = async (id: number, name: string) => {
    if (!confirm(`Are you sure you want to delete the group "${name}"?`)) return;
    setError('');
    try {
      const response = await groupsApi.delete(id);
      if (response.success) {
        loadGroups();
      } else {
        setError(response.error?.message || 'Failed to delete group');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleAddMember = async (groupId: number, userId: number) => {
    setError('');
    try {
      const response = await groupsApi.addMember(groupId, userId);
      if (response.success) {
        loadGroups();
      } else {
        setError(response.error?.message || 'Failed to add member');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleRemoveMember = async (groupId: number, userId: number, username: string) => {
    if (!confirm(`Remove ${username} from this group?`)) return;
    setError('');
    try {
      const response = await groupsApi.removeMember(groupId, userId);
      if (response.success) {
        loadGroups();
      } else {
        setError(response.error?.message || 'Failed to remove member');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
    });
  };

  const openEditModal = (group: UserGroup) => {
    setEditingGroup(group);
    setFormData({
      name: group.name,
      description: group.description || '',
    });
  };

  const getAvailableUsers = (group: UserGroup) => {
    const memberIds = new Set(group.members?.map((m) => m.id) || []);
    return allUsers.filter((user) => !memberIds.has(user.id));
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
      </div>
    );
  }

  return (
    <div className="p-6 h-full overflow-auto bg-gray-50 dark:bg-macos-dark-50">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
            User Groups
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage user groups for share access control
          </p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          + Create Group
        </Button>
      </div>

      {/* Error Display */}
      {error && (
        <div className="mb-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Groups Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {groups.map((group) => (
          <Card key={group.id} hoverable>
            <div className="p-6">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center space-x-3">
                  <div className="w-12 h-12 rounded-full bg-purple-500 text-white flex items-center justify-center">
                    <Users className="w-6 h-6" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 dark:text-gray-100 text-lg">
                      {group.name}
                    </h3>
                    {group.description && (
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        {group.description}
                      </p>
                    )}
                  </div>
                </div>
                {group.isSystem && (
                  <span className="px-2 py-1 rounded text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400">
                    System
                  </span>
                )}
              </div>

              {/* Member Count */}
              <div className="mb-4 p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg">
                <div className="text-sm text-gray-600 dark:text-gray-400">
                  {group.memberCount} {group.memberCount === 1 ? 'member' : 'members'}
                </div>
                {group.members && group.members.length > 0 && (
                  <div className="mt-2 flex flex-wrap gap-1">
                    {group.members.slice(0, 5).map((member) => (
                      <span
                        key={member.id}
                        className="px-2 py-1 rounded text-xs bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400"
                      >
                        {member.username}
                      </span>
                    ))}
                    {group.members.length > 5 && (
                      <span className="px-2 py-1 text-xs text-gray-600 dark:text-gray-400">
                        +{group.members.length - 5} more
                      </span>
                    )}
                  </div>
                )}
              </div>

              {/* Actions */}
              <div className="flex space-x-2">
                <Button
                  size="sm"
                  variant="secondary"
                  onClick={() => setManagingGroupId(group.id)}
                  className="flex-1"
                >
                  <UserPlus className="w-4 h-4 mr-1" />
                  Members
                </Button>
                {!group.isSystem && (
                  <>
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => openEditModal(group)}
                    >
                      <Edit2 className="w-4 h-4" />
                    </Button>
                    <Button
                      size="sm"
                      variant="danger"
                      onClick={() => handleDelete(group.id, group.name)}
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </>
                )}
              </div>
            </div>
          </Card>
        ))}
      </div>

      {groups.length === 0 && (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          <Users className="w-16 h-16 mx-auto mb-4 opacity-50" />
          <p>No groups created yet</p>
        </div>
      )}

      {/* Create/Edit Modal */}
      <AnimatePresence>
        {(showCreateModal || editingGroup) && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
            onClick={() => {
              setShowCreateModal(false);
              setEditingGroup(null);
              resetForm();
            }}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-md mx-4"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                {editingGroup ? 'Edit Group' : 'Create Group'}
              </h2>

              <div className="space-y-4">
                <Input
                  label="Group Name"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  placeholder="e.g., buchhaltung"
                  required
                />

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Description
                  </label>
                  <textarea
                    value={formData.description}
                    onChange={(e) =>
                      setFormData({ ...formData, description: e.target.value })
                    }
                    placeholder="Optional description"
                    rows={3}
                    className="w-full px-3 py-2 border rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-macos-blue"
                  />
                </div>

                {error && (
                  <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400 text-sm">
                    {error}
                  </div>
                )}
              </div>

              <div className="flex space-x-3 mt-6">
                <Button
                  variant="secondary"
                  onClick={() => {
                    setShowCreateModal(false);
                    setEditingGroup(null);
                    resetForm();
                  }}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button
                  onClick={editingGroup ? handleUpdate : handleCreate}
                  className="flex-1"
                >
                  {editingGroup ? 'Update' : 'Create'}
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Manage Members Modal */}
      <AnimatePresence>
        {managingGroupId !== null && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
            onClick={() => setManagingGroupId(null)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-2xl mx-4 max-h-[80vh] overflow-hidden flex flex-col"
            >
              {(() => {
                const group = groups.find((g) => g.id === managingGroupId);
                if (!group) return null;

                const availableUsers = getAvailableUsers(group);

                return (
                  <>
                    <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                      Manage Members: {group.name}
                    </h2>

                    <div className="flex-1 overflow-auto">
                      {/* Current Members */}
                      <div className="mb-6">
                        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
                          Current Members ({group.members?.length || 0})
                        </h3>
                        {group.members && group.members.length > 0 ? (
                          <div className="space-y-2">
                            {group.members.map((member) => (
                              <div
                                key={member.id}
                                className="flex items-center justify-between p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                              >
                                <div>
                                  <div className="font-medium text-gray-900 dark:text-gray-100">
                                    {member.username}
                                  </div>
                                  <div className="text-sm text-gray-600 dark:text-gray-400">
                                    {member.email}
                                  </div>
                                </div>
                                <Button
                                  size="sm"
                                  variant="danger"
                                  onClick={() =>
                                    handleRemoveMember(group.id, member.id, member.username)
                                  }
                                >
                                  <UserMinus className="w-4 h-4" />
                                </Button>
                              </div>
                            ))}
                          </div>
                        ) : (
                          <div className="text-center py-4 text-gray-500 dark:text-gray-400 text-sm">
                            No members yet
                          </div>
                        )}
                      </div>

                      {/* Available Users */}
                      <div>
                        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
                          Add Members ({availableUsers.length} available)
                        </h3>
                        {availableUsers.length > 0 ? (
                          <div className="space-y-2">
                            {availableUsers.map((user) => (
                              <div
                                key={user.id}
                                className="flex items-center justify-between p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                              >
                                <div>
                                  <div className="font-medium text-gray-900 dark:text-gray-100">
                                    {user.username}
                                  </div>
                                  <div className="text-sm text-gray-600 dark:text-gray-400">
                                    {user.email}
                                  </div>
                                </div>
                                <Button
                                  size="sm"
                                  onClick={() => handleAddMember(group.id, user.id)}
                                >
                                  <UserPlus className="w-4 h-4" />
                                </Button>
                              </div>
                            ))}
                          </div>
                        ) : (
                          <div className="text-center py-4 text-gray-500 dark:text-gray-400 text-sm">
                            All users are already members
                          </div>
                        )}
                      </div>
                    </div>

                    <div className="flex justify-end mt-6">
                      <Button
                        variant="secondary"
                        onClick={() => setManagingGroupId(null)}
                      >
                        Close
                      </Button>
                    </div>
                  </>
                );
              })()}
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
