import { Dashboard } from './Dashboard/Dashboard';
import { UserManager } from './UserManager/UserManager';
import { Settings } from './Settings/Settings';
import { StorageManager } from './StorageManager/StorageManager';
import type { App } from '@/types';

// Placeholder components for apps not yet implemented
const PlaceholderApp = ({ name }: { name: string }) => (
  <div className="flex items-center justify-center h-full bg-gray-50 dark:bg-macos-dark-50">
    <div className="text-center">
      <div className="text-6xl mb-4">🚧</div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
        {name}
      </h2>
      <p className="text-gray-600 dark:text-gray-400">Coming Soon</p>
    </div>
  </div>
);

export const registeredApps: App[] = [
  {
    id: 'dashboard',
    name: 'Dashboard',
    icon: '📊',
    component: Dashboard,
    defaultSize: { width: 900, height: 600 },
    minSize: { width: 600, height: 400 },
  },
  {
    id: 'storage',
    name: 'Storage',
    icon: '💾',
    component: StorageManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
  },
  {
    id: 'files',
    name: 'Files',
    icon: '📁',
    component: () => <PlaceholderApp name="File Station" />,
    defaultSize: { width: 900, height: 700 },
    minSize: { width: 700, height: 500 },
  },
  {
    id: 'users',
    name: 'Users',
    icon: '👥',
    component: UserManager,
    defaultSize: { width: 1000, height: 700 },
    minSize: { width: 800, height: 600 },
  },
  {
    id: 'network',
    name: 'Network',
    icon: '🌐',
    component: () => <PlaceholderApp name="Network Manager" />,
    defaultSize: { width: 800, height: 600 },
    minSize: { width: 600, height: 400 },
  },
  {
    id: 'settings',
    name: 'Settings',
    icon: '⚙️',
    component: Settings,
    defaultSize: { width: 800, height: 700 },
    minSize: { width: 600, height: 500 },
  },
];

export function getAppById(id: string): App | undefined {
  return registeredApps.find((app) => app.id === id);
}
