import { useState } from 'react';
import { UserManager as UserList } from './UserManager';
import { UserGroupManager } from './UserGroupManager';
import Button from '@/components/ui/Button';

export function UserManager() {
  const [activeTab, setActiveTab] = useState<'users' | 'groups'>('users');

  return (
    <div className="h-full flex flex-col">
      {/* Tabs */}
      <div className="bg-white dark:bg-macos-dark-100 border-b border-gray-200 dark:border-gray-700 px-6">
        <div className="flex space-x-1">
          <button
            onClick={() => setActiveTab('users')}
            className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'users'
                ? 'border-macos-blue text-macos-blue dark:text-blue-400'
                : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
            }`}
          >
            Users
          </button>
          <button
            onClick={() => setActiveTab('groups')}
            className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'groups'
                ? 'border-macos-blue text-macos-blue dark:text-blue-400'
                : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
            }`}
          >
            Groups
          </button>
        </div>
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-hidden">
        {activeTab === 'users' ? <UserList /> : <UserGroupManager />}
      </div>
    </div>
  );
}
