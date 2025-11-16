// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
import { useState, useEffect } from 'react';
import { useAuthStore } from '@/store';
import { systemApi } from '@/api/system';
import Input from '@/components/ui/Input';

// Section imports
import { GeneralSection } from './sections/GeneralSection';
import { AppearanceSection } from './sections/AppearanceSection';
import { StorageSection } from './sections/StorageSection';
import { SharesSection } from './sections/SharesSection';
import { NetworkSection } from './sections/NetworkSection';
import { UsersGroupsSection } from './sections/UsersGroupsSection';
import { BackupSection } from './sections/BackupSection';
import { TasksSection } from './sections/TasksSection';
import { MonitoringSection } from './sections/MonitoringSection';
import { BrandingSection } from './sections/BrandingSection';
import { ActiveDirectorySection } from './sections/ActiveDirectorySection';
import { UpdatesSection } from './sections/UpdatesSection';

type SettingsSection = {
  id: string;
  label: string;
  icon: string;
  component: React.ComponentType<{ user: any; systemInfo: any }>;
  searchTerms: string[];
};

const sections: SettingsSection[] = [
  {
    id: 'general',
    label: 'General',
    icon: '⚙️',
    component: GeneralSection,
    searchTerms: ['general', 'user', 'account', 'profile', '2fa', 'two-factor'],
  },
  {
    id: 'appearance',
    label: 'Appearance',
    icon: '🎨',
    component: AppearanceSection,
    searchTerms: ['appearance', 'theme', 'dark', 'light'],
  },
  {
    id: 'storage',
    label: 'Storage',
    icon: '💾',
    component: StorageSection,
    searchTerms: ['storage', 'zfs', 'pool', 'disk', 'raid', 'volume'],
  },
  {
    id: 'shares',
    label: 'Shares',
    icon: '📁',
    component: SharesSection,
    searchTerms: ['shares', 'samba', 'nfs', 'smb', 'iscsi', 'webdav', 'ftp'],
  },
  {
    id: 'network',
    label: 'Network',
    icon: '🌐',
    component: NetworkSection,
    searchTerms: ['network', 'interface', 'dns', 'firewall', 'ip', 'vlan', 'bonding'],
  },
  {
    id: 'users',
    label: 'Users & Groups',
    icon: '👥',
    component: UsersGroupsSection,
    searchTerms: ['users', 'groups', 'permissions', 'access'],
  },
  {
    id: 'backup',
    label: 'Backup',
    icon: '💼',
    component: BackupSection,
    searchTerms: ['backup', 'snapshot', 'rsync', 'cloud', 'restore'],
  },
  {
    id: 'tasks',
    label: 'Scheduled Tasks',
    icon: '⏱️',
    component: TasksSection,
    searchTerms: ['tasks', 'schedule', 'cron', 'job', 'automation'],
  },
  {
    id: 'monitoring',
    label: 'Monitoring',
    icon: '📊',
    component: MonitoringSection,
    searchTerms: ['monitoring', 'prometheus', 'grafana', 'datadog', 'metrics'],
  },
  {
    id: 'branding',
    label: 'Branding',
    icon: '🎭',
    component: BrandingSection,
    searchTerms: ['branding', 'logo', 'color', 'theme', 'customization'],
  },
  {
    id: 'ad',
    label: 'Active Directory',
    icon: '🔐',
    component: ActiveDirectorySection,
    searchTerms: ['active directory', 'ad', 'ldap', 'authentication'],
  },
  {
    id: 'updates',
    label: 'Updates',
    icon: '🔄',
    component: UpdatesSection,
    searchTerms: ['updates', 'version', 'upgrade', 'release'],
  },
];

export function Settings() {
  const user = useAuthStore((state) => state.user);
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const [activeSection, setActiveSection] = useState('general');
  const [searchQuery, setSearchQuery] = useState('');
  const [systemInfo, setSystemInfo] = useState<any>(null);

  useEffect(() => {
    const fetchSystemInfo = async () => {
      try {
        const response = await systemApi.getInfo();
        if (response.success && response.data) {
          setSystemInfo(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch system info:', error);
      }
    };

    fetchSystemInfo();
  }, []);

  const handleLogout = async () => {
    try {
      const { authApi } = await import('@/api/auth');
      await authApi.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      clearAuth();
      window.location.reload();
    }
  };

  // Filter sections based on search query
  const filteredSections = searchQuery
    ? sections.filter((section) =>
        section.searchTerms.some((term) =>
          term.toLowerCase().includes(searchQuery.toLowerCase())
        )
      )
    : sections;

  const ActiveComponent = sections.find((s) => s.id === activeSection)?.component || GeneralSection;

  return (
    <div className="flex h-full bg-gray-50 dark:bg-macos-dark-50">
      {/* Sidebar */}
      <div className="w-64 bg-white dark:bg-macos-dark-100 border-r border-gray-200 dark:border-macos-dark-300 flex flex-col">
        {/* Search */}
        <div className="p-4 border-b border-gray-200 dark:border-macos-dark-300">
          <Input
            type="text"
            placeholder="Search settings..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full"
          />
        </div>

        {/* Navigation */}
        <div className="flex-1 overflow-y-auto py-2">
          {filteredSections.map((section) => (
            <button
              key={section.id}
              onClick={() => {
                setActiveSection(section.id);
                setSearchQuery('');
              }}
              className={`w-full px-4 py-2.5 flex items-center gap-3 transition-colors ${
                activeSection === section.id
                  ? 'bg-macos-blue text-white'
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-macos-dark-200'
              }`}
            >
              <span className="text-xl">{section.icon}</span>
              <span className="text-sm font-medium">{section.label}</span>
            </button>
          ))}
        </div>

        {/* Footer - Logout */}
        <div className="p-4 border-t border-gray-200 dark:border-macos-dark-300">
          <button
            onClick={handleLogout}
            className="w-full px-4 py-2 text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
          >
            Logout
          </button>
        </div>
      </div>

      {/* Content Area */}
      <div className="flex-1 overflow-auto">
        <div className="p-6 max-w-5xl">
          <ActiveComponent user={user} systemInfo={systemInfo} />
        </div>
      </div>
    </div>
  );
}
