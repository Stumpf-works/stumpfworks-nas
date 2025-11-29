import { Dashboard } from './Dashboard/Dashboard';
import { UserManager } from './UserManager/UserManager';
import { QuotaManager } from './QuotaManager/QuotaManager';
import { Settings } from './Settings/Settings';
import { StorageManager } from './StorageManager/StorageManager';
import FileManager from './FileManager/FileManager';
import { NetworkManager } from './NetworkManager/NetworkManager';
import { DockerManager } from './DockerManager/DockerManager';
import { PluginManager } from './PluginManager/PluginManager';
import { SecurityCenter } from './SecurityCenter';
import { AppStore } from './AppStore/AppStore';
import { Terminal } from './Terminal/Terminal';
import { ADDCManager } from './ADDCManager';
import { SystemManager } from './SystemManager/SystemManager';
import { HighAvailability } from './HighAvailability/HighAvailability';
import type { App } from '@/types';

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
    component: FileManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
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
    id: 'quotas',
    name: 'Quotas',
    icon: '📊',
    component: QuotaManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
  },
  {
    id: 'security-center',
    name: 'Security Center',
    icon: '🛡️',
    component: SecurityCenter,
    defaultSize: { width: 1400, height: 900 },
    minSize: { width: 1000, height: 700 },
  },
  {
    id: 'system',
    name: 'System',
    icon: '🖥️',
    component: SystemManager,
    defaultSize: { width: 1400, height: 900 },
    minSize: { width: 1000, height: 700 },
  },
  {
    id: 'network',
    name: 'Network',
    icon: '🌐',
    component: NetworkManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
  },
  {
    id: 'ad-dc',
    name: 'AD Domain Controller',
    icon: '🏢',
    component: ADDCManager,
    defaultSize: { width: 1400, height: 900 },
    minSize: { width: 1000, height: 700 },
  },
  {
    id: 'high-availability',
    name: 'High Availability',
    icon: '⚡',
    component: HighAvailability,
    defaultSize: { width: 1400, height: 900 },
    minSize: { width: 1000, height: 700 },
  },
  {
    id: 'docker',
    name: 'Docker',
    icon: '🐳',
    component: DockerManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
  },
  {
    id: 'plugins',
    name: 'Plugins',
    icon: '🔌',
    component: PluginManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
  },
  {
    id: 'app-store',
    name: 'App Store',
    icon: '🛒',
    component: AppStore,
    defaultSize: { width: 1400, height: 900 },
    minSize: { width: 1000, height: 700 },
  },
  {
    id: 'terminal',
    name: 'Terminal',
    icon: '💻',
    component: Terminal,
    defaultSize: { width: 1000, height: 700 },
    minSize: { width: 800, height: 500 },
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

// App categories for App Gallery
export const appCategories = {
  system: ['dashboard', 'system', 'settings', 'terminal'],
  management: ['users', 'quotas', 'network', 'storage', 'ad-dc', 'high-availability'],
  security: ['security-center'],
  tools: ['files'],
  development: ['docker', 'plugins', 'app-store'],
} as const;

export type AppCategory = keyof typeof appCategories;

export const categoryNames: Record<AppCategory, string> = {
  system: 'System',
  management: 'Management',
  security: 'Security',
  tools: 'Tools',
  development: 'Development',
};

export const categoryIcons: Record<AppCategory, string> = {
  system: '⚙️',
  management: '📊',
  security: '🛡️',
  tools: '🔧',
  development: '💻',
};
