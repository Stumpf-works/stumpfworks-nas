import { useEffect, useState } from 'react';
import { usersApi, CreateUserRequest, UpdateUserRequest } from '@/api/users';
import { User } from '@/api/auth';
import { adApi, type ADUser } from '@/api/ad';
import { getErrorMessage } from '@/api/client';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Card from '@/components/ui/Card';
import { motion, AnimatePresence } from 'framer-motion';

export function UserManager() {
  const [users, setUsers] = useState<User[]>([]);
  const [adUsers, setAdUsers] = useState<ADUser[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [adLoading, setAdLoading] = useState(false);
  const [error, setError] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showAdSyncModal, setShowAdSyncModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  // Form state
  const [formData, setFormData] = useState<CreateUserRequest>({
    username: '',
    email: '',
    password: '',
    fullName: '',
    role: 'user',
  });

  const loadUsers = async () => {
    setIsLoading(true);
    setError('');
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
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadUsers();
  }, []);

  const handleCreate = async () => {
    setError('');
    try {
      const response = await usersApi.create(formData);
      if (response.success) {
        setShowCreateModal(false);
        resetForm();
        loadUsers();
      } else {
        setError(response.error?.message || 'Failed to create user');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleUpdate = async () => {
    if (!editingUser) return;
    setError('');
    try {
      const updateData: UpdateUserRequest = {
        email: formData.email,
        fullName: formData.fullName,
        role: formData.role,
      };
      if (formData.password) {
        updateData.password = formData.password;
      }
      const response = await usersApi.update(editingUser.id, updateData);
      if (response.success) {
        setEditingUser(null);
        resetForm();
        loadUsers();
      } else {
        setError(response.error?.message || 'Failed to update user');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this user?')) return;
    setError('');
    try {
      const response = await usersApi.delete(id);
      if (response.success) {
        loadUsers();
      } else {
        setError(response.error?.message || 'Failed to delete user');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  const loadAdUsers = async () => {
    setAdLoading(true);
    setError('');
    try {
      const response = await adApi.listUsers();
      if (response.success && response.data) {
        setAdUsers(response.data);
        setShowAdSyncModal(true);
      } else {
        setError(response.error?.message || 'Failed to load AD users');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setAdLoading(false);
    }
  };

  const handleSyncAdUser = async (username: string) => {
    setAdLoading(true);
    setError('');
    try {
      const response = await adApi.syncUser(username);
      if (response.success) {
        loadUsers();
        setShowAdSyncModal(false);
      } else {
        setError(response.error?.message || 'Failed to sync user');
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setAdLoading(false);
    }
  };

  const resetForm = () => {
    setFormData({
      username: '',
      email: '',
      password: '',
      fullName: '',
      role: 'user',
    });
  };

  const openEditModal = (user: User) => {
    setEditingUser(user);
    setFormData({
      username: user.username,
      email: user.email,
      password: '',
      fullName: user.fullName,
      role: user.role as 'admin' | 'user' | 'guest',
    });
  };

  const getRoleBadge = (role: string) => {
    const colors = {
      admin: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
      user: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
      guest: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400',
    };
    return colors[role as keyof typeof colors] || colors.user;
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
            User Manager
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage system users and permissions
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={loadAdUsers} disabled={adLoading}>
            {adLoading ? 'Loading...' : 'Sync from AD'}
          </Button>
          <Button onClick={() => setShowCreateModal(true)}>
            + Create User
          </Button>
        </div>
      </div>

      {/* Error Display */}
      {error && (
        <div className="mb-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Users Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {users.map((user) => (
          <Card key={user.id} hoverable>
            <div className="p-6">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center space-x-3">
                  <div className="w-12 h-12 rounded-full bg-macos-blue text-white flex items-center justify-center text-xl font-bold">
                    {user.username.charAt(0).toUpperCase()}
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                      {user.username}
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      {user.email}
                    </p>
                  </div>
                </div>
              </div>

              <div className="space-y-2 mb-4">
                {user.fullName && (
                  <p className="text-sm text-gray-700 dark:text-gray-300">
                    {user.fullName}
                  </p>
                )}
                <div className="flex items-center space-x-2">
                  <span
                    className={`px-2 py-1 rounded text-xs font-medium ${getRoleBadge(
                      user.role
                    )}`}
                  >
                    {user.role}
                  </span>
                  {user.isActive ? (
                    <span className="px-2 py-1 rounded text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400">
                      Active
                    </span>
                  ) : (
                    <span className="px-2 py-1 rounded text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400">
                      Inactive
                    </span>
                  )}
                </div>
              </div>

              <div className="flex space-x-2">
                <Button
                  size="sm"
                  variant="secondary"
                  onClick={() => openEditModal(user)}
                  className="flex-1"
                >
                  Edit
                </Button>
                <Button
                  size="sm"
                  variant="danger"
                  onClick={() => handleDelete(user.id)}
                  className="flex-1"
                >
                  Delete
                </Button>
              </div>
            </div>
          </Card>
        ))}
      </div>

      {/* Create/Edit Modal */}
      <AnimatePresence>
        {(showCreateModal || editingUser) && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
            onClick={() => {
              setShowCreateModal(false);
              setEditingUser(null);
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
                {editingUser ? 'Edit User' : 'Create User'}
              </h2>

              <div className="space-y-4">
                <Input
                  label="Username"
                  value={formData.username}
                  onChange={(e) =>
                    setFormData({ ...formData, username: e.target.value })
                  }
                  disabled={!!editingUser}
                  required
                />

                <Input
                  label="Email"
                  type="email"
                  value={formData.email}
                  onChange={(e) =>
                    setFormData({ ...formData, email: e.target.value })
                  }
                  required
                />

                <Input
                  label="Full Name"
                  value={formData.fullName}
                  onChange={(e) =>
                    setFormData({ ...formData, fullName: e.target.value })
                  }
                />

                <Input
                  label={editingUser ? 'New Password (leave blank to keep current)' : 'Password'}
                  type="password"
                  value={formData.password}
                  onChange={(e) =>
                    setFormData({ ...formData, password: e.target.value })
                  }
                  required={!editingUser}
                />

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Role
                  </label>
                  <select
                    value={formData.role}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        role: e.target.value as 'admin' | 'user' | 'guest',
                      })
                    }
                    className="w-full px-3 py-2 border rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-macos-blue"
                  >
                    <option value="user">User</option>
                    <option value="admin">Admin</option>
                    <option value="guest">Guest</option>
                  </select>
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
                    setEditingUser(null);
                    resetForm();
                  }}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button
                  onClick={editingUser ? handleUpdate : handleCreate}
                  className="flex-1"
                >
                  {editingUser ? 'Update' : 'Create'}
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* AD Sync Modal */}
      <AnimatePresence>
        {showAdSyncModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
            onClick={() => {
              setShowAdSyncModal(false);
              setAdUsers([]);
            }}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl p-6 w-full max-w-2xl mx-4 max-h-[80vh] overflow-hidden flex flex-col"
            >
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                Sync Users from Active Directory
              </h2>

              <div className="flex-1 overflow-auto mb-4">
                {adUsers.length > 0 ? (
                  <div className="space-y-2">
                    {adUsers.map((adUser) => (
                      <div
                        key={adUser.username}
                        className="flex items-center justify-between p-4 bg-gray-50 dark:bg-macos-dark-200 rounded-lg"
                      >
                        <div className="flex-1">
                          <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                            {adUser.displayName}
                          </h3>
                          <p className="text-sm text-gray-600 dark:text-gray-400">
                            {adUser.username} • {adUser.email}
                          </p>
                          {adUser.groups.length > 0 && (
                            <div className="flex flex-wrap gap-1 mt-2">
                              {adUser.groups.slice(0, 3).map((group) => (
                                <span
                                  key={group}
                                  className="px-2 py-0.5 text-xs rounded bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400"
                                >
                                  {group}
                                </span>
                              ))}
                              {adUser.groups.length > 3 && (
                                <span className="px-2 py-0.5 text-xs text-gray-600 dark:text-gray-400">
                                  +{adUser.groups.length - 3} more
                                </span>
                              )}
                            </div>
                          )}
                        </div>
                        <Button
                          size="sm"
                          onClick={() => handleSyncAdUser(adUser.username)}
                          disabled={adLoading}
                        >
                          Sync
                        </Button>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                    No AD users found
                  </div>
                )}
              </div>

              <div className="flex justify-end">
                <Button
                  variant="secondary"
                  onClick={() => {
                    setShowAdSyncModal(false);
                    setAdUsers([]);
                  }}
                >
                  Close
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
