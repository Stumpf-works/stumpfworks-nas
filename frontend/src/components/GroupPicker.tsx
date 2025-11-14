import { useEffect, useState } from 'react';
import { groupsApi } from '@/api/groups';
import { UserGroup } from '@/api/groups';

interface GroupPickerProps {
  label?: string;
  value: string[]; // Array of group names
  onChange: (groupNames: string[]) => void;
  placeholder?: string;
  error?: string;
  helperText?: string;
}

export default function GroupPicker({
  label,
  value,
  onChange,
  placeholder = 'Select groups...',
  error,
  helperText,
}: GroupPickerProps) {
  const [groups, setGroups] = useState<UserGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    loadGroups();
  }, []);

  const loadGroups = async () => {
    try {
      const response = await groupsApi.list();
      if (response.success && response.data) {
        setGroups(response.data);
      } else {
        console.error('Failed to load groups:', response.error);
      }
    } catch (err) {
      console.error('Failed to load groups:', err);
    } finally {
      setLoading(false);
    }
  };

  const filteredGroups = groups.filter(
    (group) =>
      group.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      group.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const toggleGroup = (groupName: string) => {
    if (value.includes(groupName)) {
      onChange(value.filter((g) => g !== groupName));
    } else {
      onChange([...value, groupName]);
    }
  };

  const removeGroup = (groupName: string) => {
    onChange(value.filter((g) => g !== groupName));
  };

  const selectAll = () => {
    onChange(groups.map((g) => g.name));
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

      {/* Selected Groups (Chips) */}
      {value.length > 0 && (
        <div className="mb-2 flex flex-wrap gap-2">
          {value.map((groupName) => (
            <span
              key={groupName}
              className="inline-flex items-center gap-1 px-2.5 py-1 bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-400 rounded-full text-sm font-medium"
            >
              👥 @{groupName}
              <button
                type="button"
                onClick={() => removeGroup(groupName)}
                className="ml-1 hover:text-purple-600 dark:hover:text-purple-300"
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
              : `${value.length} group${value.length !== 1 ? 's' : ''} selected`}
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
                placeholder="Search groups..."
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
                className="text-purple-600 dark:text-purple-400 hover:underline"
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

            {/* Group List */}
            <div className="overflow-y-auto">
              {loading ? (
                <div className="p-4 text-center text-gray-500">
                  <div className="animate-spin inline-block w-5 h-5 border-2 border-current border-t-transparent rounded-full" />
                  <p className="mt-2 text-sm">Loading groups...</p>
                </div>
              ) : filteredGroups.length === 0 ? (
                <div className="p-4 text-center text-gray-500 text-sm">
                  {searchTerm ? 'No groups found' : 'No groups available'}
                </div>
              ) : (
                filteredGroups.map((group) => (
                  <label
                    key={group.id}
                    className="flex items-center gap-3 px-3 py-2 hover:bg-gray-100 dark:hover:bg-macos-dark-300 cursor-pointer transition-colors"
                  >
                    <input
                      type="checkbox"
                      checked={value.includes(group.name)}
                      onChange={() => toggleGroup(group.name)}
                      className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                    />
                    <div className="flex-1 min-w-0">
                      <div className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                        @{group.name}
                      </div>
                      {group.description && (
                        <div className="text-xs text-gray-500 dark:text-gray-400 truncate">
                          {group.description}
                        </div>
                      )}
                    </div>
                    <span className="text-xs text-gray-500 dark:text-gray-400">
                      {group.memberCount} {group.memberCount === 1 ? 'member' : 'members'}
                    </span>
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
