import { useState } from 'react';
import { motion } from 'framer-motion';
import { useWindowStore } from '@/store';
import { registeredApps } from '@/apps';

interface DockIconProps {
  app: typeof registeredApps[0];
  isRunning: boolean;
  onClick: () => void;
}

function DockIcon({ app, isRunning, onClick }: DockIconProps) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <motion.div
      className="relative flex flex-col items-center group"
      onHoverStart={() => setIsHovered(true)}
      onHoverEnd={() => setIsHovered(false)}
      whileHover={{ scale: 1.4, y: -10 }}
      whileTap={{ scale: 0.9 }}
      transition={{ type: 'spring', stiffness: 300, damping: 20 }}
    >
      <button
        onClick={onClick}
        className="relative w-12 h-12 rounded-xl bg-white dark:bg-macos-dark-200 shadow-lg flex items-center justify-center text-2xl hover:shadow-xl transition-shadow"
      >
        {app.icon}
      </button>

      {/* Running Indicator */}
      {isRunning && (
        <div className="absolute -bottom-1 w-1 h-1 rounded-full bg-gray-800 dark:bg-white" />
      )}

      {/* Tooltip */}
      {isHovered && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="absolute -top-10 px-2 py-1 bg-gray-900/90 dark:bg-gray-100/90 text-white dark:text-gray-900 text-xs rounded whitespace-nowrap"
        >
          {app.name}
        </motion.div>
      )}
    </motion.div>
  );
}

export default function Dock() {
  const windows = useWindowStore((state) => state.windows);
  const openWindow = useWindowStore((state) => state.openWindow);

  const handleAppClick = (app: typeof registeredApps[0]) => {
    openWindow(app);
  };

  return (
    <motion.div
      initial={{ y: 100 }}
      animate={{ y: 0 }}
      transition={{ type: 'spring', stiffness: 100, damping: 20 }}
      className="fixed bottom-2 left-1/2 transform -translate-x-1/2 z-50"
    >
      <div className="px-4 py-2 glass-light dark:glass-dark rounded-2xl border border-gray-200/20 dark:border-gray-700/20 shadow-macos-xl">
        <div className="flex items-end space-x-3">
          {registeredApps.map((app) => (
            <DockIcon
              key={app.id}
              app={app}
              isRunning={windows.some((w) => w.appId === app.id)}
              onClick={() => handleAppClick(app)}
            />
          ))}
        </div>
      </div>
    </motion.div>
  );
}
