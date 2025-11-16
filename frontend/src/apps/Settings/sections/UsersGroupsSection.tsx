// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState, useEffect } from 'react';
import Card from '@/components/ui/Card';
import { usersApi } from '@/api/users';
import { groupsApi } from '@/api/groups';
import { getErrorMessage } from '@/api/client';

export function UsersGroupsSection() {
  const [activeTab, setActiveTab] = useState<'users' | 'groups'>('users');
  const [users, setUsers] = useState<any[]>([]);
  const [groups, setGroups] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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
      }
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Users & Groups</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Manage local users, groups, and permissions
        </p>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 dark:border-macos-dark-300">
        <button
          onClick={() => setActiveTab('users')}
          className={`px-4 py-2 font-medium border-b-2 transition-colors ${
            activeTab === 'users'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          Users
        </button>
        <button
          onClick={() => setActiveTab('groups')}
          className={`px-4 py-2 font-medium border-b-2 transition-colors ${
            activeTab === 'groups'
              ? 'border-macos-blue text-macos-blue'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100'
          }`}
        >
          Groups
        </button>
      </div>

      {/* Users Tab */}
      {activeTab === 'users' && (
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Local Users</h2>
            {loading ? (
              <p className="text-gray-600 dark:text-gray-400">Loading users...</p>
            ) : users.length === 0 ? (
              <p className="text-gray-600 dark:text-gray-400">No users found</p>
            ) : (
              <div className="space-y-3">
                {users.map((u: any) => (
                  <div
                    key={u.id || u.username}
                    className="p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg"
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-gray-100">{u.username}</h3>
                        <p className="text-sm text-gray-600 dark:text-gray-400">{u.email || 'No email'}</p>
                      </div>
                      <span className="px-2 py-1 text-xs bg-gray-100 dark:bg-macos-dark-200 text-gray-800 dark:text-gray-300 rounded">
                        {u.role || 'user'}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </Card>
      )}

      {/* Groups Tab */}
      {activeTab === 'groups' && (
        <Card>
          <div className="p-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Groups</h2>
            {loading ? (
              <p className="text-gray-600 dark:text-gray-400">Loading groups...</p>
            ) : groups.length === 0 ? (
              <p className="text-gray-600 dark:text-gray-400">No groups found</p>
            ) : (
              <div className="space-y-3">
                {groups.map((g: any) => (
                  <div
                    key={g.id || g.name}
                    className="p-4 border border-gray-200 dark:border-macos-dark-300 rounded-lg"
                  >
                    <h3 className="font-semibold text-gray-900 dark:text-gray-100">{g.name}</h3>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      {g.memberCount || 0} members
                    </p>
                  </div>
                ))}
              </div>
            )}
          </div>
        </Card>
      )}

      {/* Note */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          For advanced user management including creating/editing users and managing permissions,
          use the dedicated User Manager app.
        </p>
      </div>
    </div>
  );
}
