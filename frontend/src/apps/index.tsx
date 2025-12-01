import { lazy } from 'react';
import type { App } from '@/types';

// Lazy load all apps for better initial load performance
const Dashboard = lazy(() => import('./Dashboard/Dashboard').then(m => ({ default: m.Dashboard })));
const UserManager = lazy(() => import('./UserManager/UserManager').then(m => ({ default: m.UserManager })));
const QuotaManager = lazy(() => import('./QuotaManager/QuotaManager').then(m => ({ default: m.QuotaManager })));
const Settings = lazy(() => import('./Settings/Settings').then(m => ({ default: m.Settings })));
const StorageManager = lazy(() => import('./StorageManager/StorageManager').then(m => ({ default: m.StorageManager })));
const FileManager = lazy(() => import('./FileManager/FileManager'));
const NetworkManager = lazy(() => import('./NetworkManager/NetworkManager').then(m => ({ default: m.NetworkManager })));
const DockerManager = lazy(() => import('./DockerManager/DockerManager').then(m => ({ default: m.DockerManager })));
const PluginManager = lazy(() => import('./PluginManager/PluginManager').then(m => ({ default: m.PluginManager })));
const SecurityCenter = lazy(() => import('./SecurityCenter').then(m => ({ default: m.SecurityCenter })));
const AppStore = lazy(() => import('./AppStore/AppStore').then(m => ({ default: m.AppStore })));
const Terminal = lazy(() => import('./Terminal/Terminal').then(m => ({ default: m.Terminal })));
const ADDCManager = lazy(() => import('./ADDCManager').then(m => ({ default: m.ADDCManager })));
const SystemManager = lazy(() => import('./SystemManager/SystemManager').then(m => ({ default: m.SystemManager })));
const HighAvailability = lazy(() => import('./HighAvailability/HighAvailability').then(m => ({ default: m.HighAvailability })));
const VMManager = lazy(() => import('./VMManager').then(m => ({ default: m.VMManager })));
const LXCManager = lazy(() => import('./LXCManager').then(m => ({ default: m.LXCManager })));

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
  {
    id: 'vm-manager',
    name: 'VM Manager',
    icon: '🖥️',
    component: VMManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
  },
  {
    id: 'lxc-manager',
    name: 'LXC Manager',
    icon: '📦',
    component: LXCManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
  },
];

export function getAppById(id: string): App | undefined {
  return registeredApps.find((app) => app.id === id);
}

// App categories for App Gallery
export const appCategories = {
  system: ['dashboard', 'system', 'settings', 'terminal'],
  management: ['users', 'quotas', 'network', 'storage', 'ad-dc', 'high-availability', 'vm-manager', 'lxc-manager'],
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
