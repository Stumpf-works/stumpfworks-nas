import { useEffect, useState } from 'react';
import { useAuthStore, useSystemStore, useThemeStore } from '@/store';
import { systemApi } from '@/api/system';
import { motion } from 'framer-motion';
import UpdateNotification from '@/components/UpdateNotification';

export default function TopBar() {
  const user = useAuthStore((state) => state.user);
  const metrics = useSystemStore((state) => state.metrics);
  const setMetrics = useSystemStore((state) => state.setMetrics);
  const toggleTheme = useThemeStore((state) => state.toggleTheme);
  const isDark = useThemeStore((state) => state.isDark);

  const [time, setTime] = useState(new Date());

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

  return (
    <motion.div
      initial={{ y: -32 }}
      animate={{ y: 0 }}
      className="fixed top-0 left-0 right-0 h-8 flex items-center justify-between px-4 glass-light dark:glass-dark border-b border-gray-200/20 dark:border-gray-700/20 z-50"
    >
      {/* Left Section */}
      <div className="flex items-center space-x-4">
        <div className="flex items-center space-x-2">
          <svg
            className="w-5 h-5 text-macos-blue"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z" />
          </svg>
          <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
            Stumpf.Works NAS
          </span>
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

        {/* User */}
        {user && (
          <div className="flex items-center space-x-2">
            <div className="w-6 h-6 rounded-full bg-macos-blue text-white flex items-center justify-center text-xs font-semibold">
              {user.username.charAt(0).toUpperCase()}
            </div>
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
  );
}
