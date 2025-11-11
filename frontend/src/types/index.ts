export interface Window {
  id: string;
  appId: string;
  title: string;
  icon?: string;
  position: { x: number; y: number };
  size: { width: number; height: number };
  state: 'normal' | 'minimized' | 'maximized' | 'fullscreen';
  zIndex: number;
  isFocused: boolean;
  isResizable: boolean;
  minSize?: { width: number; height: number };
  maxSize?: { width: number; height: number };
}

export interface App {
  id: string;
  name: string;
  icon: string;
  component: React.ComponentType;
  defaultSize: { width: number; height: number };
  minSize?: { width: number; height: number };
  isResizable?: boolean;
}

export interface DockApp {
  id: string;
  name: string;
  icon: string;
  isRunning: boolean;
  badge?: number;
}
