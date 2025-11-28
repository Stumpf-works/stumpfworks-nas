import { useEffect, useState } from 'react';
import { useThemeStore, useWindowStore } from '@/store';
import { getAppById } from '@/apps';
import TopBar from './TopBar';
import Dock from './Dock';
import WindowManager from './WindowManager';
import WidgetSidebar from '@/components/WidgetSidebar';
import { AppGallery } from '@/components/AppGallery';

export default function Desktop() {
  const setTheme = useThemeStore((state) => state.setTheme);
  const isDark = useThemeStore((state) => state.isDark);
  const openWindow = useWindowStore((state) => state.openWindow);
  const [isWidgetSidebarOpen, setIsWidgetSidebarOpen] = useState(false);
  const [isAppGalleryOpen, setIsAppGalleryOpen] = useState(false);

  useEffect(() => {
    // Initialize theme
    setTheme(isDark);
  }, []);

  // Expose global function to open App Gallery
  useEffect(() => {
    (window as any).openAppGallery = () => setIsAppGalleryOpen(true);
    return () => {
      delete (window as any).openAppGallery;
    };
  }, []);

  const handleLaunchApp = (appId: string) => {
    const app = getAppById(appId);
    if (app) {
      openWindow(app);
    }
  };

  return (
    <div className="relative w-full h-full overflow-hidden bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50 dark:from-macos-dark-50 dark:via-macos-dark-100 dark:to-macos-dark-200">
      {/* Background Wallpaper */}
      <div
        className="absolute inset-0 opacity-30 dark:opacity-20"
        style={{
          backgroundImage: 'url(https://images.unsplash.com/photo-1557683316-973673baf926?w=1920)',
          backgroundSize: 'cover',
          backgroundPosition: 'center',
        }}
      />

      {/* Top Bar */}
      <TopBar />

      {/* Window Manager */}
      <div className="absolute inset-0 top-8 bottom-20">
        <WindowManager />
      </div>

      {/* Dock */}
      <Dock />

      {/* Widget Sidebar */}
      <WidgetSidebar
        isOpen={isWidgetSidebarOpen}
        onToggle={() => setIsWidgetSidebarOpen(!isWidgetSidebarOpen)}
      />

      {/* App Gallery */}
      <AppGallery
        isOpen={isAppGalleryOpen}
        onClose={() => setIsAppGalleryOpen(false)}
        onLaunchApp={handleLaunchApp}
      />
    </div>
  );
}
