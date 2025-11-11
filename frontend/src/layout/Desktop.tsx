import { useEffect } from 'react';
import { useThemeStore } from '@/store';
import TopBar from './TopBar';
import Dock from './Dock';
import WindowManager from './WindowManager';

export default function Desktop() {
  const setTheme = useThemeStore((state) => state.setTheme);
  const isDark = useThemeStore((state) => state.isDark);

  useEffect(() => {
    // Initialize theme
    setTheme(isDark);
  }, []);

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
    </div>
  );
}
