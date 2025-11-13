import { useEffect, useState } from 'react';
import { usersApi } from '@/api/users';
import { User } from '@/api/auth';

interface UserPickerProps {
  label?: string;
  value: string[]; // Array of usernames
  onChange: (usernames: string[]) => void;
  placeholder?: string;
  error?: string;
  helperText?: string;
}

export default function UserPicker({
  label,
  value,
  onChange,
  placeholder = 'Select users...',
  error,
  helperText,
}: UserPickerProps) {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      const response = await usersApi.list();
      if (response.success) {
        setUsers(response.data);
      } else {
        console.error('Failed to load users:', response.error);
      }
    } catch (err) {
      console.error('Failed to load users:', err);
    } finally {
      setLoading(false);
    }
  };

  const filteredUsers = users.filter(
    (user) =>
      user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
      user.fullName?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const toggleUser = (username: string) => {
    if (value.includes(username)) {
      onChange(value.filter((u) => u !== username));
    } else {
      onChange([...value, username]);
    }
  };

  const removeUser = (username: string) => {
    onChange(value.filter((u) => u !== username));
  };

  const selectAll = () => {
    onChange(users.map((u) => u.username));
  };

  const clearAll = () => {
    onChange([]);
  };

  return (
    <div className="relative">
      {label && (
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          {label}
        </label>
      )}

      {/* Selected Users (Chips) */}
      {value.length > 0 && (
        <div className="mb-2 flex flex-wrap gap-2">
          {value.map((username) => (
            <span
              key={username}
              className="inline-flex items-center gap-1 px-2.5 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 rounded-full text-sm font-medium"
            >
              👤 {username}
              <button
                type="button"
                onClick={() => removeUser(username)}
                className="ml-1 hover:text-blue-600 dark:hover:text-blue-300"
              >
                ✕
              </button>
            </span>
          ))}
        </div>
      )}

      {/* Dropdown Trigger */}
      <div className="relative">
        <button
          type="button"
          onClick={() => setShowDropdown(!showDropdown)}
          className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 text-left flex items-center justify-between hover:border-gray-400 dark:hover:border-gray-500 transition-colors"
        >
          <span className="text-sm">
            {value.length === 0
              ? placeholder
              : `${value.length} user${value.length !== 1 ? 's' : ''} selected`}
          </span>
          <span className="text-gray-500">
            {showDropdown ? '▲' : '▼'}
          </span>
        </button>

        {/* Dropdown Menu */}
        {showDropdown && (
          <div className="absolute z-10 mt-1 w-full bg-white dark:bg-macos-dark-200 border border-gray-300 dark:border-gray-600 rounded-lg shadow-lg max-h-64 overflow-hidden flex flex-col">
            {/* Search Box */}
            <div className="p-2 border-b border-gray-200 dark:border-gray-700">
              <input
                type="text"
                placeholder="Search users..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-macos-dark-300 text-gray-900 dark:text-gray-100"
              />
            </div>

            {/* Bulk Actions */}
            <div className="p-2 border-b border-gray-200 dark:border-gray-700 flex justify-between text-xs">
              <button
                type="button"
                onClick={selectAll}
                className="text-blue-600 dark:text-blue-400 hover:underline"
              >
                Select All
              </button>
              <button
                type="button"
                onClick={clearAll}
                className="text-gray-600 dark:text-gray-400 hover:underline"
              >
                Clear All
              </button>
            </div>

            {/* User List */}
            <div className="overflow-y-auto">
              {loading ? (
                <div className="p-4 text-center text-gray-500">
                  <div className="animate-spin inline-block w-5 h-5 border-2 border-current border-t-transparent rounded-full" />
                  <p className="mt-2 text-sm">Loading users...</p>
                </div>
              ) : filteredUsers.length === 0 ? (
                <div className="p-4 text-center text-gray-500 text-sm">
                  {searchTerm ? 'No users found' : 'No users available'}
                </div>
              ) : (
                filteredUsers.map((user) => (
                  <label
                    key={user.id}
                    className="flex items-center gap-3 px-3 py-2 hover:bg-gray-100 dark:hover:bg-macos-dark-300 cursor-pointer transition-colors"
                  >
                    <input
                      type="checkbox"
                      checked={value.includes(user.username)}
                      onChange={() => toggleUser(user.username)}
                      className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <div className="flex-1 min-w-0">
                      <div className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                        {user.username}
                      </div>
                      {user.fullName && (
                        <div className="text-xs text-gray-500 dark:text-gray-400 truncate">
                          {user.fullName}
                        </div>
                      )}
                    </div>
                    {user.role === 'admin' && (
                      <span className="text-xs px-2 py-0.5 bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-400 rounded">
                        Admin
                      </span>
                    )}
                  </label>
                ))
              )}
            </div>
          </div>
        )}
      </div>

      {/* Helper Text / Error */}
      {(helperText || error) && (
        <p className={`mt-1 text-sm ${error ? 'text-red-600 dark:text-red-400' : 'text-gray-500 dark:text-gray-400'}`}>
          {error || helperText}
        </p>
      )}
    </div>
  );
}
