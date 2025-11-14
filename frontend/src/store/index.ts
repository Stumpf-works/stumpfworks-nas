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
