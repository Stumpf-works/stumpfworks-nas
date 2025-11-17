// Revision: 2025-11-17 | Author: Claude | Version: 1.3.0
import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { usersApi } from '@/api/users';
import { groupsApi, UserGroup } from '@/api/groups';
import { User } from '@/api/auth';
import { getErrorMessage } from '@/api/client';

export function UsersGroupsSection() {
  const [activeTab, setActiveTab] = useState<'users' | 'groups'>('users');
  const [users, setUsers] = useState<User[]>([]);
  const [groups, setGroups] = useState<UserGroup[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // User Dialog State
  const [userDialog, setUserDialog] = useState<{
    open: boolean;
    mode: 'create' | 'edit';
    user?: User;
  }>({ open: false, mode: 'create' });

  const [userForm, setUserForm] = useState({
    username: '',
    email: '',
    password: '',
    role: 'user' as 'admin' | 'user',
  });

  // Group Dialog State
  const [groupDialog, setGroupDialog] = useState<{
    open: boolean;
    mode: 'create' | 'edit';
    group?: UserGroup;
  }>({ open: false, mode: 'create' });

  const [groupForm, setGroupForm] = useState({
    name: '',
    description: '',
  });

  // Group Members Dialog State
  const [membersDialog, setMembersDialog] = useState<{
    open: boolean;
    group?: UserGroup;
    members: any[];
  }>({ open: false, members: [] });

  const [selectedUserId, setSelectedUserId] = useState<string>('');

  useEffect(() => {
    if (activeTab === 'users') {
      loadUsers();
    } else {
      loadGroups();
    }
  }, [activeTab]);

  const loadUsers = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await usersApi.list();
      if (response.success && response.data) {
        setUsers(response.data);
      } else {
        setError(response.error?.message || 'Failed to load users');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const loadGroups = async () => {
    setLoading(true);
    setError(null);
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
      setLoading(false);
    }
  };

  // User CRUD Operations
  const handleCreateUser = () => {
    setUserForm({ username: '', email: '', password: '', role: 'user' });
    setUserDialog({ open: true, mode: 'create' });
  };

  const handleEditUser = (user: User) => {
    setUserForm({
      username: user.username,
      email: user.email || '',
      password: '',
      role: user.role as 'admin' | 'user',
    });
    setUserDialog({ open: true, mode: 'edit', user });
  };

  const handleSaveUser = async () => {
    setError(null);
    setSuccess(null);

    try {
      if (userDialog.mode === 'create') {
        const response = await usersApi.create(userForm);
        if (response.success) {
          setSuccess('User created successfully');
          setUserDialog({ open: false, mode: 'create' });
          loadUsers();
        } else {
          setError(response.error?.message || 'Failed to create user');
        }
      } else if (userDialog.user) {
        const updateData: any = {
          email: userForm.email,
          role: userForm.role,
        };
        if (userForm.password) {
          updateData.password = userForm.password;
        }

        const response = await usersApi.update(userDialog.user.id, updateData);
        if (response.success) {
          setSuccess('User updated successfully');
          setUserDialog({ open: false, mode: 'edit' });
          loadUsers();
        } else {
          setError(response.error?.message || 'Failed to update user');
        }
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleDeleteUser = async (user: User) => {
    if (!confirm(`Are you sure you want to delete user "${user.username}"?`)) {
      return;
    }

    try {
      const response = await usersApi.delete(user.id);
      if (response.success) {
        setSuccess('User deleted successfully');
        loadUsers();
      } else {
        setError(response.error?.message || 'Failed to delete user');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  // Group CRUD Operations
  const handleCreateGroup = () => {
    setGroupForm({ name: '', description: '' });
    setGroupDialog({ open: true, mode: 'create' });
  };

  const handleEditGroup = (group: UserGroup) => {
    setGroupForm({
      name: group.name,
      description: group.description || '',
    });
    setGroupDialog({ open: true, mode: 'edit', group });
  };

  const handleSaveGroup = async () => {
    setError(null);
    setSuccess(null);

    try {
      if (groupDialog.mode === 'create') {
        const response = await groupsApi.create(groupForm);
        if (response.success) {
          setSuccess('Group created successfully');
          setGroupDialog({ open: false, mode: 'create' });
          loadGroups();
        } else {
          setError(response.error?.message || 'Failed to create group');
        }
      } else if (groupDialog.group) {
        const response = await groupsApi.update(groupDialog.group.id, groupForm);
        if (response.success) {
          setSuccess('Group updated successfully');
          setGroupDialog({ open: false, mode: 'edit' });
          loadGroups();
        } else {
          setError(response.error?.message || 'Failed to update group');
        }
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleDeleteGroup = async (group: UserGroup) => {
    if (!confirm(`Are you sure you want to delete group "${group.name}"?`)) {
      return;
    }

    try {
      const response = await groupsApi.delete(group.id);
      if (response.success) {
        setSuccess('Group deleted successfully');
        loadGroups();
      } else {
        setError(response.error?.message || 'Failed to delete group');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  // Group Members Management
  const handleManageMembers = async (group: UserGroup) => {
    try {
      const response = await groupsApi.getMembers(group.id);
      if (response.success && response.data) {
        setMembersDialog({ open: true, group, members: response.data });
      } else {
        setError(response.error?.message || 'Failed to load group members');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleAddMember = async () => {
    if (!membersDialog.group || !selectedUserId) return;

    try {
      const response = await groupsApi.addMember(membersDialog.group.id, parseInt(selectedUserId));
      if (response.success) {
        setSuccess('Member added successfully');
        handleManageMembers(membersDialog.group);
        setSelectedUserId('');
      } else {
        setError(response.error?.message || 'Failed to add member');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleRemoveMember = async (userId: number) => {
    if (!membersDialog.group) return;

    if (!confirm('Are you sure you want to remove this member?')) {
      return;
    }

    try {
      const response = await groupsApi.removeMember(membersDialog.group.id, userId);
      if (response.success) {
        setSuccess('Member removed successfully');
        handleManageMembers(membersDialog.group);
      } else {
        setError(response.error?.message || 'Failed to remove member');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  // Get available users for adding to group (not already in group)
  const availableUsers = users.filter(
    (u) => !membersDialog.members.some((m) => m.id === u.id)
  );

  return (
    <div className="space-y-4 md:space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-xl md:text-2xl font-bold text-gray-900 dark:text-gray-100">
          Users & Groups
        </h1>
        <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mt-1">
          Manage local users, groups, and permissions
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
          onClick={() => setActiveTab('users')}
          className={`px-3 md:px-4 py-2 font-medium border-b-2 transition-colors text-sm md:text-base ${
            activeTab === 'users'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          Users ({users.length})
        </button>
        <button
          onClick={() => setActiveTab('groups')}
          className={`px-3 md:px-4 py-2 font-medium border-b-2 transition-colors text-sm md:text-base ${
            activeTab === 'groups'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          Groups ({groups.length})
        </button>
      </div>

      {/* Users Tab */}
      {activeTab === 'users' && (
        <div className="space-y-4">
          <div className="flex justify-end">
            <Button onClick={handleCreateUser}>
              <span className="hidden sm:inline">+ Create User</span>
              <span className="sm:hidden">+ User</span>
            </Button>
          </div>

          <Card>
            <div className="p-4 md:p-6">
              {loading ? (
                <p className="text-sm md:text-base text-gray-600 dark:text-gray-400">
                  Loading users...
                </p>
              ) : users.length === 0 ? (
                <div className="text-center py-8 md:py-12">
                  <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mb-4">
                    No users found
                  </p>
                  <Button onClick={handleCreateUser}>Create your first user</Button>
                </div>
              ) : (
                <div className="space-y-3">
                  {users.map((user) => (
                    <motion.div
                      key={user.id}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="p-3 md:p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg hover:border-macos-blue transition-colors"
                    >
                      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 md:gap-3">
                            <h3 className="text-base md:text-lg font-semibold text-gray-900 dark:text-gray-100">
                              {user.username}
                            </h3>
                            <span
                              className={`px-2 py-1 text-xs rounded ${
                                user.role === 'admin'
                                  ? 'bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-200'
                                  : 'bg-gray-100 dark:bg-macos-dark-200 text-gray-800 dark:text-gray-300'
                              }`}
                            >
                              {user.role}
                            </span>
                          </div>
                          <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400 mt-1">
                            {user.email || 'No email'}
                          </p>
                        </div>
                        <div className="flex gap-2 w-full sm:w-auto">
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() => handleEditUser(user)}
                            className="flex-1 sm:flex-initial"
                          >
                            Edit
                          </Button>
                          <Button
                            variant="danger"
                            size="sm"
                            onClick={() => handleDeleteUser(user)}
                            disabled={user.username === 'admin'}
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

      {/* Groups Tab */}
      {activeTab === 'groups' && (
        <div className="space-y-4">
          <div className="flex justify-end">
            <Button onClick={handleCreateGroup}>
              <span className="hidden sm:inline">+ Create Group</span>
              <span className="sm:hidden">+ Group</span>
            </Button>
          </div>

          <Card>
            <div className="p-4 md:p-6">
              {loading ? (
                <p className="text-sm md:text-base text-gray-600 dark:text-gray-400">
                  Loading groups...
                </p>
              ) : groups.length === 0 ? (
                <div className="text-center py-8 md:py-12">
                  <p className="text-sm md:text-base text-gray-600 dark:text-gray-400 mb-4">
                    No groups found
                  </p>
                  <Button onClick={handleCreateGroup}>Create your first group</Button>
                </div>
              ) : (
                <div className="space-y-3">
                  {groups.map((group) => (
                    <motion.div
                      key={group.id}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="p-3 md:p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg hover:border-macos-blue transition-colors"
                    >
                      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
                        <div className="flex-1">
                          <h3 className="text-base md:text-lg font-semibold text-gray-900 dark:text-gray-100">
                            {group.name}
                          </h3>
                          <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400 mt-1">
                            {group.description || 'No description'}
                          </p>
                          <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                            {group.memberCount || 0} members
                          </p>
                        </div>
                        <div className="flex gap-2 w-full sm:w-auto flex-wrap sm:flex-nowrap">
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() => handleManageMembers(group)}
                            className="flex-1 sm:flex-initial"
                          >
                            Members
                          </Button>
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() => handleEditGroup(group)}
                            className="flex-1 sm:flex-initial"
                          >
                            Edit
                          </Button>
                          <Button
                            variant="danger"
                            size="sm"
                            onClick={() => handleDeleteGroup(group)}
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

      {/* User Create/Edit Dialog */}
      <AnimatePresence>
        {userDialog.open && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
            onClick={() => setUserDialog({ open: false, mode: 'create' })}
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
                  {userDialog.mode === 'create' ? 'Create User' : 'Edit User'}
                </h2>

                <div className="space-y-4">
                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Username
                    </label>
                    <Input
                      value={userForm.username}
                      onChange={(e) => setUserForm({ ...userForm, username: e.target.value })}
                      disabled={userDialog.mode === 'edit'}
                      placeholder="Enter username"
                    />
                  </div>

                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Email
                    </label>
                    <Input
                      type="email"
                      value={userForm.email}
                      onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
                      placeholder="Enter email (optional)"
                    />
                  </div>

                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Password {userDialog.mode === 'edit' && '(leave empty to keep current)'}
                    </label>
                    <Input
                      type="password"
                      value={userForm.password}
                      onChange={(e) => setUserForm({ ...userForm, password: e.target.value })}
                      placeholder={
                        userDialog.mode === 'create' ? 'Enter password' : 'Leave empty to keep current'
                      }
                    />
                  </div>

                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Role
                    </label>
                    <select
                      value={userForm.role}
                      onChange={(e) =>
                        setUserForm({ ...userForm, role: e.target.value as 'admin' | 'user' })
                      }
                      className="w-full px-3 py-2 text-sm bg-white dark:bg-macos-dark-200 border border-gray-300 dark:border-macos-dark-300 rounded-md text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-macos-blue"
                    >
                      <option value="user">User</option>
                      <option value="admin">Admin</option>
                    </select>
                  </div>
                </div>

                <div className="flex gap-2 mt-6">
                  <Button
                    variant="secondary"
                    onClick={() => setUserDialog({ open: false, mode: 'create' })}
                    className="flex-1"
                  >
                    Cancel
                  </Button>
                  <Button onClick={handleSaveUser} className="flex-1">
                    {userDialog.mode === 'create' ? 'Create' : 'Save'}
                  </Button>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Group Create/Edit Dialog */}
      <AnimatePresence>
        {groupDialog.open && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
            onClick={() => setGroupDialog({ open: false, mode: 'create' })}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-md w-full"
            >
              <div className="p-4 md:p-6">
                <h2 className="text-lg md:text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                  {groupDialog.mode === 'create' ? 'Create Group' : 'Edit Group'}
                </h2>

                <div className="space-y-4">
                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Group Name
                    </label>
                    <Input
                      value={groupForm.name}
                      onChange={(e) => setGroupForm({ ...groupForm, name: e.target.value })}
                      placeholder="Enter group name"
                    />
                  </div>

                  <div>
                    <label className="block text-xs md:text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Description
                    </label>
                    <textarea
                      value={groupForm.description}
                      onChange={(e) => setGroupForm({ ...groupForm, description: e.target.value })}
                      placeholder="Enter description (optional)"
                      rows={3}
                      className="w-full px-3 py-2 text-sm bg-white dark:bg-macos-dark-200 border border-gray-300 dark:border-macos-dark-300 rounded-md text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-macos-blue"
                    />
                  </div>
                </div>

                <div className="flex gap-2 mt-6">
                  <Button
                    variant="secondary"
                    onClick={() => setGroupDialog({ open: false, mode: 'create' })}
                    className="flex-1"
                  >
                    Cancel
                  </Button>
                  <Button onClick={handleSaveGroup} className="flex-1">
                    {groupDialog.mode === 'create' ? 'Create' : 'Save'}
                  </Button>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Group Members Dialog */}
      <AnimatePresence>
        {membersDialog.open && membersDialog.group && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
            onClick={() => setMembersDialog({ open: false, members: [] })}
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-auto"
            >
              <div className="p-4 md:p-6">
                <h2 className="text-lg md:text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                  Manage Members: {membersDialog.group.name}
                </h2>

                {/* Add Member Section */}
                {availableUsers.length > 0 && (
                  <div className="mb-6 p-3 md:p-4 bg-gray-50 dark:bg-macos-dark-200 rounded-lg">
                    <h3 className="text-sm md:text-base font-semibold text-gray-900 dark:text-gray-100 mb-3">
                      Add Member
                    </h3>
                    <div className="flex gap-2">
                      <select
                        value={selectedUserId}
                        onChange={(e) => setSelectedUserId(e.target.value)}
                        className="flex-1 px-3 py-2 text-sm bg-white dark:bg-macos-dark-100 border border-gray-300 dark:border-macos-dark-300 rounded-md text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-macos-blue"
                      >
                        <option value="">Select a user...</option>
                        {availableUsers.map((user) => (
                          <option key={user.id} value={user.id}>
                            {user.username} ({user.email || 'no email'})
                          </option>
                        ))}
                      </select>
                      <Button onClick={handleAddMember} disabled={!selectedUserId} size="sm">
                        Add
                      </Button>
                    </div>
                  </div>
                )}

                {/* Current Members */}
                <div>
                  <h3 className="text-sm md:text-base font-semibold text-gray-900 dark:text-gray-100 mb-3">
                    Current Members ({membersDialog.members.length})
                  </h3>
                  {membersDialog.members.length === 0 ? (
                    <p className="text-sm text-gray-600 dark:text-gray-400 text-center py-4">
                      No members in this group
                    </p>
                  ) : (
                    <div className="space-y-2">
                      {membersDialog.members.map((member: any) => (
                        <div
                          key={member.id}
                          className="flex items-center justify-between p-3 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                        >
                          <div>
                            <p className="text-sm md:text-base font-medium text-gray-900 dark:text-gray-100">
                              {member.username}
                            </p>
                            <p className="text-xs md:text-sm text-gray-600 dark:text-gray-400">
                              {member.email || 'No email'}
                            </p>
                          </div>
                          <Button
                            variant="danger"
                            size="sm"
                            onClick={() => handleRemoveMember(member.id)}
                          >
                            Remove
                          </Button>
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                <div className="mt-6">
                  <Button
                    variant="secondary"
                    onClick={() => setMembersDialog({ open: false, members: [] })}
                    className="w-full"
                  >
                    Close
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
