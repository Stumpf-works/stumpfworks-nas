import { useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FileInfo } from '@/api/files';

interface ContextMenuProps {
  x: number;
  y: number;
  file: FileInfo | null;
  onClose: () => void;
  onOpen?: () => void;
  onDownload?: () => void;
  onRename?: () => void;
  onCopy?: () => void;
  onCut?: () => void;
  onDelete?: () => void;
  onPermissions?: () => void;
  onCompress?: () => void;
  isAdmin?: boolean;
}

export default function ContextMenu({
  x,
  y,
  file,
  onClose,
  onOpen,
  onDownload,
  onRename,
  onCopy,
  onCut,
  onDelete,
  onPermissions,
  onCompress,
  isAdmin = false,
}: ContextMenuProps) {
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        onClose();
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    document.addEventListener('keydown', handleEscape);

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.removeEventListener('keydown', handleEscape);
    };
  }, [onClose]);

  // Adjust position to keep menu on screen
  useEffect(() => {
    if (menuRef.current) {
      const rect = menuRef.current.getBoundingClientRect();
      const viewportWidth = window.innerWidth;
      const viewportHeight = window.innerHeight;

      let adjustedX = x;
      let adjustedY = y;

      if (x + rect.width > viewportWidth) {
        adjustedX = viewportWidth - rect.width - 10;
      }

      if (y + rect.height > viewportHeight) {
        adjustedY = viewportHeight - rect.height - 10;
      }

      menuRef.current.style.left = `${adjustedX}px`;
      menuRef.current.style.top = `${adjustedY}px`;
    }
  }, [x, y]);

  const handleAction = (action?: () => void) => {
    if (action) {
      action();
    }
    onClose();
  };

  const menuItems = [
    ...(file?.isDir && onOpen ? [
      { label: 'Open', icon: '📂', action: onOpen, divider: true }
    ] : []),
    ...(!file?.isDir && onOpen ? [
      { label: 'Open', icon: '👁️', action: onOpen, divider: true }
    ] : []),
    ...(onDownload && !file?.isDir ? [
      { label: 'Download', icon: '⬇️', action: onDownload }
    ] : []),
    ...(onRename ? [
      { label: 'Rename', icon: '✏️', action: onRename, shortcut: 'F2' }
    ] : []),
    ...(onCopy ? [
      { label: 'Copy', icon: '📋', action: onCopy, shortcut: 'Ctrl+C' }
    ] : []),
    ...(onCut ? [
      { label: 'Cut', icon: '✂️', action: onCut, shortcut: 'Ctrl+X', divider: true }
    ] : []),
    ...(onCompress ? [
      { label: 'Compress', icon: '🗜️', action: onCompress, divider: true }
    ] : []),
    ...(isAdmin && onPermissions ? [
      { label: 'Permissions', icon: '🔒', action: onPermissions }
    ] : []),
    ...(onDelete ? [
      { label: 'Delete', icon: '🗑️', action: onDelete, danger: true, shortcut: 'Del' }
    ] : []),
  ];

  return (
    <AnimatePresence>
      <motion.div
        ref={menuRef}
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        transition={{ duration: 0.1 }}
        className="fixed z-50 min-w-[200px] bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl border border-gray-200 dark:border-gray-700 py-1"
        style={{ left: x, top: y }}
      >
        {file && (
          <div className="px-3 py-2 border-b border-gray-200 dark:border-gray-700">
            <div className="flex items-center space-x-2">
              <span className="text-2xl">{file.isDir ? '📁' : '📄'}</span>
              <div className="flex-1 min-w-0">
                <div className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                  {file.name}
                </div>
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  {file.isDir ? 'Folder' : file.mimeType || 'File'}
                </div>
              </div>
            </div>
          </div>
        )}

        <div className="py-1">
          {menuItems.map((item, index) => (
            <div key={index}>
              <button
                onClick={() => handleAction(item.action)}
                className={`w-full px-3 py-2 text-left text-sm flex items-center justify-between space-x-3 transition-colors ${
                  item.danger
                    ? 'text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                }`}
              >
                <div className="flex items-center space-x-2">
                  <span>{item.icon}</span>
                  <span>{item.label}</span>
                </div>
                {item.shortcut && (
                  <span className="text-xs text-gray-400 dark:text-gray-500">
                    {item.shortcut}
                  </span>
                )}
              </button>
              {item.divider && (
                <div className="my-1 h-px bg-gray-200 dark:bg-gray-700" />
              )}
            </div>
          ))}
        </div>
      </motion.div>
    </AnimatePresence>
  );
}
