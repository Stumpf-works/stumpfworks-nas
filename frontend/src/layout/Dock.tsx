import { useState } from 'react';
import { motion } from 'framer-motion';
import { useWindowStore, useDockStore } from '@/store';
import { registeredApps, getAppById } from '@/apps';

interface DockIconProps {
  app: typeof registeredApps[0];
  isRunning: boolean;
  onClick: () => void;
  onRemove: () => void;
}

function DockIcon({ app, isRunning, onClick, onRemove }: DockIconProps) {
  const [isHovered, setIsHovered] = useState(false);
  const [showContextMenu, setShowContextMenu] = useState(false);
  const [contextMenuPos, setContextMenuPos] = useState({ x: 0, y: 0 });

  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault();
    setContextMenuPos({ x: e.clientX, y: e.clientY });
    setShowContextMenu(true);
  };

  return (
    <>
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
          onContextMenu={handleContextMenu}
          className="relative w-12 h-12 rounded-xl bg-white dark:bg-macos-dark-200 shadow-lg flex items-center justify-center text-2xl hover:shadow-xl transition-shadow"
        >
          {app.icon}
        </button>

        {/* Running Indicator */}
        {isRunning && (
          <div className="absolute -bottom-1 w-1 h-1 rounded-full bg-gray-800 dark:bg-white" />
        )}

        {/* Tooltip */}
        {isHovered && !showContextMenu && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="absolute -top-10 px-2 py-1 bg-gray-900/90 dark:bg-gray-100/90 text-white dark:text-gray-900 text-xs rounded whitespace-nowrap"
          >
            {app.name}
          </motion.div>
        )}
      </motion.div>

      {/* Context Menu */}
      {showContextMenu && (
        <>
          <div
            className="fixed inset-0 z-50"
            onClick={() => setShowContextMenu(false)}
          />
          <div
            className="fixed z-50 bg-white dark:bg-macos-dark-200 rounded-lg shadow-2xl border border-gray-200 dark:border-gray-700 py-1 min-w-[160px]"
            style={{
              left: `${contextMenuPos.x}px`,
              top: `${contextMenuPos.y}px`,
            }}
          >
            <button
              onClick={() => {
                onRemove();
                setShowContextMenu(false);
              }}
              className="w-full px-4 py-2 text-left text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-macos-dark-300 transition-colors"
            >
              Remove from Dock
            </button>
          </div>
        </>
      )}
    </>
  );
}

export default function Dock() {
  const windows = useWindowStore((state) => state.windows);
  const openWindow = useWindowStore((state) => state.openWindow);
  const dockApps = useDockStore((state) => state.dockApps);
  const removeFromDock = useDockStore((state) => state.removeFromDock);

  const handleAppClick = (app: typeof registeredApps[0]) => {
    openWindow(app);
  };

  const handleRemoveFromDock = (appId: string) => {
    removeFromDock(appId);
  };

  // Get apps that are in the dock (in order)
  const dockAppsList = dockApps
    .map((appId) => getAppById(appId))
    .filter((app): app is NonNullable<typeof app> => app !== undefined);

  return (
    <motion.div
      initial={{ y: 100 }}
      animate={{ y: 0 }}
      transition={{ type: 'spring', stiffness: 100, damping: 20 }}
      className="fixed bottom-2 left-1/2 transform -translate-x-1/2 z-50"
    >
      <div className="px-4 py-2 glass-light dark:glass-dark rounded-2xl border border-gray-200/20 dark:border-gray-700/20 shadow-macos-xl">
        <div className="flex items-end space-x-3">
          {dockAppsList.map((app) => (
            <DockIcon
              key={app.id}
              app={app}
              isRunning={windows.some((w) => w.appId === app.id)}
              onClick={() => handleAppClick(app)}
              onRemove={() => handleRemoveFromDock(app.id)}
            />
          ))}
        </div>
      </div>
    </motion.div>
  );
}
