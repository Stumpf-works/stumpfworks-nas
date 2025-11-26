import { Dashboard } from './Dashboard/Dashboard';
import { UserManager } from './UserManager/UserManager';
import { Settings } from './Settings/Settings';
import { StorageManager } from './StorageManager/StorageManager';
import FileManager from './FileManager/FileManager';
import { NetworkManager } from './NetworkManager/NetworkManager';
import { DockerManager } from './DockerManager/DockerManager';
import { PluginManager } from './PluginManager/PluginManager';
import { BackupManager } from './BackupManager/BackupManager';
import { AuditLogs } from './AuditLogs/AuditLogs';
import { Security } from './Security/Security';
import { Alerts } from './Alerts/Alerts';
import { Tasks } from './Tasks/Tasks';
import { AppStore } from './AppStore/AppStore';
import { Terminal } from './Terminal/Terminal';
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
    id: 'audit-logs',
    name: 'Audit Logs',
    icon: '🔒',
    component: AuditLogs,
    defaultSize: { width: 1400, height: 800 },
    minSize: { width: 1000, height: 600 },
  },
  {
    id: 'security',
    name: 'Security',
    icon: '🛡️',
    component: Security,
    defaultSize: { width: 1400, height: 800 },
    minSize: { width: 1000, height: 600 },
  },
  {
    id: 'alerts',
    name: 'Alerts',
    icon: '🔔',
    component: Alerts,
    defaultSize: { width: 1000, height: 800 },
    minSize: { width: 800, height: 600 },
  },
  {
    id: 'tasks',
    name: 'Scheduled Tasks',
    icon: '📅',
    component: Tasks,
    defaultSize: { width: 1400, height: 800 },
    minSize: { width: 1000, height: 600 },
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
    id: 'backups',
    name: 'Backups',
    icon: '⏱️',
    component: BackupManager,
    defaultSize: { width: 1200, height: 800 },
    minSize: { width: 900, height: 600 },
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
