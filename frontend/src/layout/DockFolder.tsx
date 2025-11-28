import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FolderOpen, Edit, Trash2 } from 'lucide-react';
import { useWindowStore } from '@/store';
import { getAppById } from '@/apps';
import type { DockFolder as DockFolderType } from '@/store';
import type { App } from '@/types';

interface DockFolderProps {
  folder: DockFolderType;
  onRemoveApp: (appId: string) => void;
  onDeleteFolder: () => void;
  onRenameFolder: () => void;
}

export default function DockFolder({
  folder,
  onRemoveApp,
  onDeleteFolder,
  onRenameFolder,
}: DockFolderProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [isHovered, setIsHovered] = useState(false);
  const [showContextMenu, setShowContextMenu] = useState(false);
  const [contextMenuPos, setContextMenuPos] = useState({ x: 0, y: 0 });

  const windows = useWindowStore((state) => state.windows);
  const openWindow = useWindowStore((state) => state.openWindow);

  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault();
    setContextMenuPos({ x: e.clientX, y: e.clientY });
    setShowContextMenu(true);
  };

  const handleAppClick = (app: App) => {
    openWindow(app);
    setIsExpanded(false);
  };

  const handleAppContextMenu = (e: React.MouseEvent, appId: string) => {
    e.preventDefault();
    e.stopPropagation();
    if (confirm(`Remove ${getAppById(appId)?.name} from folder?`)) {
      onRemoveApp(appId);
    }
  };

  // Get apps in folder
  const folderApps = folder.apps
    .map((appId) => getAppById(appId))
    .filter((app): app is NonNullable<typeof app> => app !== undefined);

  return (
    <>
      <motion.div
        className="relative flex flex-col items-center group"
        onHoverStart={() => setIsHovered(true)}
        onHoverEnd={() => {
          setIsHovered(false);
          // Keep expanded state if clicked
        }}
        whileHover={{ scale: 1.4, y: -10 }}
        whileTap={{ scale: 0.9 }}
        transition={{ type: 'spring', stiffness: 300, damping: 20 }}
      >
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          onContextMenu={handleContextMenu}
          className="relative w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500/90 to-purple-600/90 shadow-lg flex items-center justify-center text-2xl hover:shadow-xl transition-shadow"
        >
          {folder.icon}
          <div className="absolute -bottom-1 -right-1 w-4 h-4 rounded-full bg-white dark:bg-macos-dark-200 flex items-center justify-center">
            <FolderOpen className="w-3 h-3 text-blue-500" />
          </div>
        </button>

        {/* Folder indicator - show number of apps */}
        <div className="absolute -bottom-1 w-1 h-1 rounded-full bg-gray-800 dark:bg-white" />

        {/* Tooltip */}
        {isHovered && !showContextMenu && !isExpanded && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="absolute -top-10 px-2 py-1 bg-gray-900/90 dark:bg-gray-100/90 text-white dark:text-gray-900 text-xs rounded whitespace-nowrap z-50"
          >
            {folder.name} ({folder.apps.length})
          </motion.div>
        )}
      </motion.div>

      {/* Expanded Folder Grid */}
      <AnimatePresence>
        {isExpanded && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="fixed inset-0 z-40"
              onClick={() => setIsExpanded(false)}
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.8, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.8, y: 20 }}
              transition={{ type: 'spring', stiffness: 300, damping: 25 }}
              className="fixed bottom-20 left-1/2 transform -translate-x-1/2 z-50 p-6 glass-light dark:glass-dark rounded-2xl border border-gray-200/20 dark:border-gray-700/20 shadow-macos-xl"
              style={{ minWidth: '300px', maxWidth: '500px' }}
            >
              <div className="mb-3 flex items-center justify-between">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                  <span>{folder.icon}</span>
                  <span>{folder.name}</span>
                </h3>
                <button
                  onClick={() => setIsExpanded(false)}
                  className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
                >
                  ✕
                </button>
              </div>

              <div className="grid grid-cols-4 gap-4">
                {folderApps.map((app) => (
                  <motion.div
                    key={app.id}
                    whileHover={{ scale: 1.1 }}
                    whileTap={{ scale: 0.95 }}
                    className="flex flex-col items-center gap-1 cursor-pointer group"
                  >
                    <button
                      onClick={() => handleAppClick(app)}
                      onContextMenu={(e) => handleAppContextMenu(e, app.id)}
                      className="relative w-12 h-12 rounded-xl bg-white dark:bg-macos-dark-200 shadow-md flex items-center justify-center text-2xl hover:shadow-lg transition-shadow"
                    >
                      {app.icon}
                      {windows.some((w) => w.appId === app.id) && (
                        <div className="absolute -bottom-1 w-1 h-1 rounded-full bg-gray-800 dark:bg-white" />
                      )}
                    </button>
                    <span className="text-xs text-gray-700 dark:text-gray-300 text-center line-clamp-2 max-w-[60px]">
                      {app.name}
                    </span>
                  </motion.div>
                ))}
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>

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
                onRenameFolder();
                setShowContextMenu(false);
              }}
              className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-macos-dark-300 transition-colors flex items-center gap-2"
            >
              <Edit className="w-4 h-4" />
              Rename Folder
            </button>
            <button
              onClick={() => {
                onDeleteFolder();
                setShowContextMenu(false);
              }}
              className="w-full px-4 py-2 text-left text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-macos-dark-300 transition-colors flex items-center gap-2"
            >
              <Trash2 className="w-4 h-4" />
              Delete Folder
            </button>
          </div>
        </>
      )}
    </>
  );
}
