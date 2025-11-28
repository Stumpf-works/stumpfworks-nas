import { useState, useEffect } from 'react';
import { addcApi, CreateUserRequest } from '../../../api/addc';
import { motion, AnimatePresence } from 'framer-motion';
import { Users, Plus, Trash2, RefreshCw, AlertCircle, UserCheck, UserX } from 'lucide-react';

export default function UserManagement() {
  const [users, setUsers] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  const [createForm, setCreateForm] = useState<CreateUserRequest>({
    user: {
      username: '',
      given_name: '',
      surname: '',
      display_name: '',
      email: '',
      description: '',
      department: '',
      company: '',
      title: '',
      telephone: '',
      ou: '',
      enabled: true,
      password_expired: false,
    },
    password: '',
  });

  const loadUsers = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await addcApi.listUsers();
      if (response.success && response.data) {
        setUsers(response.data);
      } else {
        setError(response.error?.message || 'Failed to load users');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load users');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadUsers();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!createForm.user.username || !createForm.password) {
      setError('Username and password are required');
      return;
    }

    try {
      setActionLoading('create');
      setError('');
      const response = await addcApi.createUser(createForm);

      if (response.success) {
        alert(`User ${createForm.user.username} created successfully!`);
        setShowCreateForm(false);
        setCreateForm({
          user: {
            username: '',
            given_name: '',
            surname: '',
            display_name: '',
            email: '',
            description: '',
            department: '',
            company: '',
            title: '',
            telephone: '',
            ou: '',
            enabled: true,
            password_expired: false,
          },
          password: '',
        });
        loadUsers();
      } else {
        setError(response.error?.message || 'Failed to create user');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create user');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async (username: string) => {
    if (!confirm(`Are you sure you want to delete user "${username}"?`)) {
      return;
    }

    try {
      setActionLoading(`delete-${username}`);
      setError('');
      const response = await addcApi.deleteUser(username);

      if (response.success) {
        alert(`User ${username} deleted successfully`);
        loadUsers();
      } else {
        setError(response.error?.message || 'Failed to delete user');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to delete user');
    } finally {
      setActionLoading(null);
    }
  };

  const handleEnable = async (username: string) => {
    try {
      setActionLoading(`enable-${username}`);
      setError('');
      const response = await addcApi.enableUser(username);

      if (response.success) {
        alert(`User ${username} enabled successfully`);
        loadUsers();
      } else {
        setError(response.error?.message || 'Failed to enable user');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to enable user');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDisable = async (username: string) => {
    try {
      setActionLoading(`disable-${username}`);
      setError('');
      const response = await addcApi.disableUser(username);

      if (response.success) {
        alert(`User ${username} disabled successfully`);
        loadUsers();
      } else {
        setError(response.error?.message || 'Failed to disable user');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to disable user');
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
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">User Management</h2>
        </div>
        <div className="flex gap-3">
          <button
            onClick={loadUsers}
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
            Create User
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

      {/* Users List */}
      {loading && users.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <RefreshCw className="w-6 h-6 animate-spin text-gray-400" />
        </div>
      ) : users.length === 0 ? (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          No users found
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {users.map((username) => (
            <div
              key={username}
              className="bg-white dark:bg-macos-dark-100 border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-2">
                  <Users className="w-5 h-5 text-macos-blue" />
                  <span className="font-medium text-gray-900 dark:text-gray-100">{username}</span>
                </div>
              </div>

              <div className="flex flex-wrap gap-2">
                <button
                  onClick={() => handleEnable(username)}
                  disabled={actionLoading === `enable-${username}`}
                  className="p-1.5 text-green-600 dark:text-green-400 hover:bg-green-50 dark:hover:bg-green-900/20 rounded transition-colors disabled:opacity-50"
                  title="Enable User"
                >
                  <UserCheck className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDisable(username)}
                  disabled={actionLoading === `disable-${username}`}
                  className="p-1.5 text-orange-600 dark:text-orange-400 hover:bg-orange-50 dark:hover:bg-orange-900/20 rounded transition-colors disabled:opacity-50"
                  title="Disable User"
                >
                  <UserX className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(username)}
                  disabled={actionLoading === `delete-${username}`}
                  className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors disabled:opacity-50"
                  title="Delete User"
                >
                  <Trash2 className="w-4 h-4" />
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
                  Create New User
                </h2>

                <form onSubmit={handleCreate} className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Username <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={createForm.user.username}
                        onChange={(e) => setCreateForm({
                          ...createForm,
                          user: { ...createForm.user, username: e.target.value },
                        })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                        required
                      />
                    </div>

                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Password <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="password"
                        value={createForm.password}
                        onChange={(e) => setCreateForm({ ...createForm, password: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                        required
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        First Name
                      </label>
                      <input
                        type="text"
                        value={createForm.user.given_name}
                        onChange={(e) => setCreateForm({
                          ...createForm,
                          user: { ...createForm.user, given_name: e.target.value },
                        })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Last Name
                      </label>
                      <input
                        type="text"
                        value={createForm.user.surname}
                        onChange={(e) => setCreateForm({
                          ...createForm,
                          user: { ...createForm.user, surname: e.target.value },
                        })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Email
                      </label>
                      <input
                        type="email"
                        value={createForm.user.email}
                        onChange={(e) => setCreateForm({
                          ...createForm,
                          user: { ...createForm.user, email: e.target.value },
                        })}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-50 text-gray-900 dark:text-gray-100"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Department
                      </label>
                      <input
                        type="text"
                        value={createForm.user.department}
                        onChange={(e) => setCreateForm({
                          ...createForm,
                          user: { ...createForm.user, department: e.target.value },
                        })}
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
                          Create User
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
    </div>
  );
}
