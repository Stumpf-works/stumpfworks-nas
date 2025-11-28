import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import { User } from '@/api/auth';
import { SystemMetrics } from '@/api/system';
import type { Window as WindowType, App } from '@/types';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, accessToken: string, refreshToken: string) => void;
  clearAuth: () => void;
  setUser: (user: User) => void;
}

interface SystemState {
  metrics: SystemMetrics | null;
  isLoading: boolean;
  setMetrics: (metrics: SystemMetrics) => void;
  setLoading: (loading: boolean) => void;
}

interface WindowState {
  windows: WindowType[];
  focusedWindowId: string | null;
  openWindow: (app: App) => void;
  closeWindow: (windowId: string) => void;
  focusWindow: (windowId: string) => void;
  minimizeWindow: (windowId: string) => void;
  maximizeWindow: (windowId: string) => void;
  updateWindowPosition: (windowId: string, position: { x: number; y: number }) => void;
  updateWindowSize: (windowId: string, size: { width: number; height: number }) => void;
}

// @ts-ignore - AppState defined for future use
interface AppState {
  apps: App[];
  runningApps: string[];
  registerApp: (app: App) => void;
  launchApp: (appId: string) => void;
}

interface ThemeState {
  isDark: boolean;
  toggleTheme: () => void;
  setTheme: (isDark: boolean) => void;
}

export interface DockFolder {
  id: string;
  name: string;
  icon: string;
  apps: string[]; // Array of app IDs
}

export type DockItem = string | DockFolder;

export function isDockFolder(item: DockItem): item is DockFolder {
  return typeof item === 'object' && 'apps' in item;
}

interface DockState {
  dockItems: DockItem[]; // Array of app IDs or folders
  addToDock: (appId: string) => void;
  removeFromDock: (appId: string) => void;
  reorderDock: (from: number, to: number) => void;
  resetToDefault: () => void;
  isInDock: (appId: string) => boolean;
  // Folder management
  createFolder: (name: string, icon: string, appIds: string[]) => string;
  deleteFolder: (folderId: string) => void;
  addAppToFolder: (folderId: string, appId: string) => void;
  removeAppFromFolder: (folderId: string, appId: string) => void;
  renameFolder: (folderId: string, name: string) => void;
  getFolderById: (folderId: string) => DockFolder | undefined;
}

// Default dock apps (essential apps only)
const DEFAULT_DOCK_APPS = [
  'dashboard',
  'files',
  'storage',
  'network',
  'docker',
  'terminal',
  'settings',
];

// Auth Store
export const useAuthStore = create<AuthState>()(
  devtools(
    persist(
      (set) => ({
        user: null,
        accessToken: null,
        refreshToken: null,
        isAuthenticated: false,
        setAuth: (user, accessToken, refreshToken) => {
          localStorage.setItem('accessToken', accessToken);
          localStorage.setItem('refreshToken', refreshToken);
          set({ user, accessToken, refreshToken, isAuthenticated: true });
        },
        clearAuth: () => {
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
          set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false });
        },
        setUser: (user) => set({ user }),
      }),
      { name: 'auth-storage' }
    ),
    { name: 'AuthStore' }
  )
);

// System Store
export const useSystemStore = create<SystemState>()(
  devtools(
    (set) => ({
      metrics: null,
      isLoading: false,
      setMetrics: (metrics) => set({ metrics }),
      setLoading: (loading) => set({ isLoading: loading }),
    }),
    { name: 'SystemStore' }
  )
);

// Window Store
export const useWindowStore = create<WindowState>()(
  devtools(
    (set, get) => ({
      windows: [],
      focusedWindowId: null,

      openWindow: (app) => {
        const windows = get().windows;
        const existingWindow = windows.find((w) => w.appId === app.id);

        if (existingWindow) {
          // Bring existing window to front
          get().focusWindow(existingWindow.id);
        } else {
          // Create new window
          const newWindow: WindowType = {
            id: `window-${app.id}-${Date.now()}`,
            appId: app.id,
            title: app.name,
            icon: app.icon,
            position: { x: 100 + windows.length * 30, y: 80 + windows.length * 30 },
            size: app.defaultSize,
            state: 'normal',
            zIndex: windows.length + 1,
            isFocused: true,
            isResizable: app.isResizable !== false,
            minSize: app.minSize,
          };

          set({
            windows: [...windows.map((w) => ({ ...w, isFocused: false })), newWindow],
            focusedWindowId: newWindow.id,
          });
        }
      },

      closeWindow: (windowId) => {
        const windows = get().windows.filter((w) => w.id !== windowId);
        const focusedWindowId = get().focusedWindowId === windowId ? null : get().focusedWindowId;
        set({ windows, focusedWindowId });
      },

      focusWindow: (windowId) => {
        const windows = get().windows;
        const maxZ = Math.max(...windows.map((w) => w.zIndex), 0);

        set({
          windows: windows.map((w) =>
            w.id === windowId
              ? { ...w, isFocused: true, zIndex: maxZ + 1, state: w.state === 'minimized' ? 'normal' : w.state }
              : { ...w, isFocused: false }
          ),
          focusedWindowId: windowId,
        });
      },

      minimizeWindow: (windowId) => {
        set({
          windows: get().windows.map((w) =>
            w.id === windowId ? { ...w, state: 'minimized', isFocused: false } : w
          ),
          focusedWindowId: null,
        });
      },

      maximizeWindow: (windowId) => {
        const window = get().windows.find((w) => w.id === windowId);
        if (!window) return;

        const newState = window.state === 'maximized' ? 'normal' : 'maximized';
        set({
          windows: get().windows.map((w) => (w.id === windowId ? { ...w, state: newState } : w)),
        });
      },

      updateWindowPosition: (windowId, position) => {
        set({
          windows: get().windows.map((w) => (w.id === windowId ? { ...w, position } : w)),
        });
      },

      updateWindowSize: (windowId, size) => {
        set({
          windows: get().windows.map((w) => (w.id === windowId ? { ...w, size } : w)),
        });
      },
    }),
    { name: 'WindowStore' }
  )
);

// Theme Store
export const useThemeStore = create<ThemeState>()(
  devtools(
    persist(
      (set) => ({
        isDark: false,
        toggleTheme: () =>
          set((state) => {
            const newDark = !state.isDark;
            if (newDark) {
              document.documentElement.classList.add('dark');
            } else {
              document.documentElement.classList.remove('dark');
            }
            return { isDark: newDark };
          }),
        setTheme: (isDark) => {
          if (isDark) {
            document.documentElement.classList.add('dark');
          } else {
            document.documentElement.classList.remove('dark');
          }
          set({ isDark });
        },
      }),
      { name: 'theme-storage' }
    ),
    { name: 'ThemeStore' }
  )
);

// Dock Store
export const useDockStore = create<DockState>()(
  devtools(
    persist(
      (set, get) => ({
        dockItems: DEFAULT_DOCK_APPS,

        addToDock: (appId) => {
          const { dockItems } = get();
          // Check if app is already in dock (including in folders)
          const isInDock = get().isInDock(appId);
          if (!isInDock) {
            set({ dockItems: [...dockItems, appId] });
          }
        },

        removeFromDock: (appId) => {
          const { dockItems } = get();
          // Remove from top level or from folders
          const newDockItems = dockItems
            .map((item) => {
              if (isDockFolder(item)) {
                return {
                  ...item,
                  apps: item.apps.filter((id) => id !== appId),
                };
              }
              return item;
            })
            .filter((item) => {
              // Remove the app if it's at top level
              if (typeof item === 'string') {
                return item !== appId;
              }
              // Remove folders that become empty
              return item.apps.length > 0;
            });
          set({ dockItems: newDockItems });
        },

        reorderDock: (from, to) => {
          const { dockItems } = get();
          const newDockItems = [...dockItems];
          const [removed] = newDockItems.splice(from, 1);
          newDockItems.splice(to, 0, removed);
          set({ dockItems: newDockItems });
        },

        resetToDefault: () => {
          set({ dockItems: DEFAULT_DOCK_APPS });
        },

        isInDock: (appId) => {
          const { dockItems } = get();
          return dockItems.some((item) => {
            if (typeof item === 'string') {
              return item === appId;
            }
            return item.apps.includes(appId);
          });
        },

        // Folder management
        createFolder: (name, icon, appIds) => {
          const { dockItems } = get();
          const folderId = `folder-${Date.now()}`;
          const folder: DockFolder = {
            id: folderId,
            name,
            icon,
            apps: appIds,
          };
          // Remove apps from dock that are now in folder
          const newDockItems = dockItems.filter(
            (item) => typeof item === 'object' || !appIds.includes(item)
          );
          set({ dockItems: [...newDockItems, folder] });
          return folderId;
        },

        deleteFolder: (folderId) => {
          const { dockItems } = get();
          const folder = get().getFolderById(folderId);
          if (!folder) return;

          // Remove folder and add its apps back to dock
          const newDockItems = dockItems.filter(
            (item) => !(isDockFolder(item) && item.id === folderId)
          );
          set({ dockItems: [...newDockItems, ...folder.apps] });
        },

        addAppToFolder: (folderId, appId) => {
          const { dockItems } = get();
          const newDockItems = dockItems.map((item) => {
            if (isDockFolder(item) && item.id === folderId) {
              if (!item.apps.includes(appId)) {
                return { ...item, apps: [...item.apps, appId] };
              }
            }
            return item;
          });
          set({ dockItems: newDockItems });
        },

        removeAppFromFolder: (folderId, appId) => {
          const { dockItems } = get();
          const newDockItems = dockItems
            .map((item) => {
              if (isDockFolder(item) && item.id === folderId) {
                return {
                  ...item,
                  apps: item.apps.filter((id) => id !== appId),
                };
              }
              return item;
            })
            .filter((item) => {
              // Remove folders that become empty
              if (isDockFolder(item)) {
                return item.apps.length > 0;
              }
              return true;
            });
          set({ dockItems: newDockItems });
        },

        renameFolder: (folderId, name) => {
          const { dockItems } = get();
          const newDockItems = dockItems.map((item) => {
            if (isDockFolder(item) && item.id === folderId) {
              return { ...item, name };
            }
            return item;
          });
          set({ dockItems: newDockItems });
        },

        getFolderById: (folderId) => {
          const { dockItems } = get();
          const folder = dockItems.find(
            (item) => isDockFolder(item) && item.id === folderId
          );
          return folder && isDockFolder(folder) ? folder : undefined;
        },
      }),
      { name: 'dock-storage' }
    ),
    { name: 'DockStore' }
  )
);
