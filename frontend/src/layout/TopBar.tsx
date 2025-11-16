// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import { useEffect, useState, useRef } from 'react';
import { useAuthStore, useSystemStore, useThemeStore } from '@/store';
import { systemApi } from '@/api/system';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronDown, LogOut, Info } from 'lucide-react';
import UpdateNotification from '@/components/UpdateNotification';
import AboutDialog from '@/components/AboutDialog';

export default function TopBar() {
  const user = useAuthStore((state) => state.user);
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const metrics = useSystemStore((state) => state.metrics);
  const setMetrics = useSystemStore((state) => state.setMetrics);
  const toggleTheme = useThemeStore((state) => state.toggleTheme);
  const isDark = useThemeStore((state) => state.isDark);

  const [time, setTime] = useState(new Date());
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
  const [isAboutDialogOpen, setIsAboutDialogOpen] = useState(false);

  const menuRef = useRef<HTMLDivElement>(null);
  const userMenuRef = useRef<HTMLDivElement>(null);

  // Update time every second
  useEffect(() => {
    const timer = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(timer);
  }, []);

  // Fetch metrics every 5 seconds
  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const response = await systemApi.getMetrics();
        if (response.success && response.data) {
          setMetrics(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch metrics:', error);
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 5000);
    return () => clearInterval(interval);
  }, [setMetrics]);

  // Handle click outside to close menus
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false);
      }
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setIsUserMenuOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    });
  };

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
    });
  };

  const handleLogout = () => {
    clearAuth();
    setIsUserMenuOpen(false);
  };

  const handleAboutClick = () => {
    setIsAboutDialogOpen(true);
    setIsMenuOpen(false);
  };

  return (
    <>
      <motion.div
        initial={{ y: -32 }}
        animate={{ y: 0 }}
        className="fixed top-0 left-0 right-0 h-8 flex items-center justify-between px-4 glass-light dark:glass-dark border-b border-gray-200/20 dark:border-gray-700/20 z-50"
      >
        {/* Left Section */}
        <div className="flex items-center space-x-4">
          <div className="relative" ref={menuRef}>
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="flex items-center space-x-2 hover:bg-gray-200/30 dark:hover:bg-gray-700/30 px-2 py-1 rounded transition-colors"
            >
              <svg
                className="w-5 h-5 text-macos-blue"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
              </svg>
              <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                Stumpf.Works NAS
              </span>
              <ChevronDown className="w-3 h-3 text-gray-600 dark:text-gray-400" />
            </button>

            {/* System Menu Dropdown */}
            <AnimatePresence>
              {isMenuOpen && (
                <motion.div
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  transition={{ duration: 0.15 }}
                  className="absolute top-full left-0 mt-1 w-56 bg-white/95 dark:bg-macos-dark-100/95 backdrop-blur-xl rounded-lg shadow-macos-lg border border-gray-200/20 dark:border-gray-700/20 py-1 z-50"
                >
                  <button
                    onClick={handleAboutClick}
                    className="w-full px-4 py-2 text-left text-sm text-gray-900 dark:text-gray-100 hover:bg-macos-blue/10 dark:hover:bg-macos-blue/20 flex items-center gap-2 transition-colors"
                  >
                    <Info className="w-4 h-4" />
                    About This NAS
                  </button>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </div>

        {/* Right Section */}
        <div className="flex items-center space-x-4 text-xs">
          {/* System Metrics */}
          {metrics && (
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-1">
                <svg
                  className="w-3 h-3 text-gray-600 dark:text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"
                  />
                </svg>
                <span className="text-gray-700 dark:text-gray-300">
                  {metrics.cpu.usagePercent.toFixed(1)}%
                </span>
              </div>

              <div className="flex items-center space-x-1">
                <svg
                  className="w-3 h-3 text-gray-600 dark:text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                  />
                </svg>
                <span className="text-gray-700 dark:text-gray-300">
                  {metrics.memory.usedPercent.toFixed(1)}%
                </span>
              </div>
            </div>
          )}

          {/* Theme Toggle */}
          <button
            onClick={toggleTheme}
            className="p-1 rounded hover:bg-gray-200/50 dark:hover:bg-gray-700/50 transition-colors"
            aria-label="Toggle theme"
          >
            {isDark ? (
              <svg
                className="w-4 h-4 text-yellow-500"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z"
                  clipRule="evenodd"
                />
              </svg>
            ) : (
              <svg
                className="w-4 h-4 text-gray-700"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
              </svg>
            )}
          </button>

          {/* Update Notification */}
          <UpdateNotification />

          {/* User Menu */}
          {user && (
            <div className="relative" ref={userMenuRef}>
              <button
                onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                className="flex items-center space-x-2 hover:bg-gray-200/30 dark:hover:bg-gray-700/30 px-2 py-1 rounded transition-colors"
              >
                <div className="w-6 h-6 rounded-full bg-macos-blue text-white flex items-center justify-center text-xs font-semibold">
                  {user.username.charAt(0).toUpperCase()}
                </div>
              </button>

              {/* User Dropdown */}
              <AnimatePresence>
                {isUserMenuOpen && (
                  <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                    transition={{ duration: 0.15 }}
                    className="absolute top-full right-0 mt-1 w-64 bg-white/95 dark:bg-macos-dark-100/95 backdrop-blur-xl rounded-lg shadow-macos-lg border border-gray-200/20 dark:border-gray-700/20 overflow-hidden z-50"
                  >
                    {/* User Info Header */}
                    <div className="px-4 py-3 bg-gradient-to-br from-macos-blue/10 to-macos-purple/10 dark:from-macos-blue/20 dark:to-macos-purple/20 border-b border-gray-200/20 dark:border-gray-700/20">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-macos-blue text-white flex items-center justify-center text-lg font-semibold">
                          {user.username.charAt(0).toUpperCase()}
                        </div>
                        <div>
                          <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                            {user.username}
                          </div>
                          <div className="text-xs text-gray-600 dark:text-gray-400">
                            {user.email || 'Administrator'}
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Menu Items */}
                    <div className="py-1">
                      <button
                        onClick={handleLogout}
                        className="w-full px-4 py-2 text-left text-sm text-gray-900 dark:text-gray-100 hover:bg-red-50 dark:hover:bg-red-900/20 flex items-center gap-2 transition-colors"
                      >
                        <LogOut className="w-4 h-4 text-red-600 dark:text-red-400" />
                        <span>Log Out</span>
                      </button>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          )}

          {/* Time & Date */}
          <div className="flex flex-col items-end leading-tight">
            <span className="text-gray-900 dark:text-gray-100 font-medium">
              {formatTime(time)}
            </span>
            <span className="text-gray-600 dark:text-gray-400 text-xs">
              {formatDate(time)}
            </span>
          </div>
        </div>
      </motion.div>

      {/* About Dialog */}
      <AboutDialog
        isOpen={isAboutDialogOpen}
        onClose={() => setIsAboutDialogOpen(false)}
      />
    </>
  );
}
