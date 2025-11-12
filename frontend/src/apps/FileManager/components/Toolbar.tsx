import React from 'react';
import { motion } from 'framer-motion';

interface ToolbarProps {
  onRefresh: () => void;
  onNewFolder: () => void;
  onUpload: () => void;
  onDownload: () => void;
  onDelete: () => void;
  onPermissions: () => void;
  onToggleHidden: () => void;
  showHidden: boolean;
  viewMode: 'list' | 'grid';
  onViewModeChange: (mode: 'list' | 'grid') => void;
  selectedCount: number;
  canDelete: boolean;
  canDownload: boolean;
  canPermissions: boolean;
}

const Toolbar: React.FC<ToolbarProps> = ({
  onRefresh,
  onNewFolder,
  onUpload,
  onDownload,
  onDelete,
  onPermissions,
  onToggleHidden,
  showHidden,
  viewMode,
  onViewModeChange,
  selectedCount,
  canDelete,
  canDownload,
  canPermissions,
}) => {
  return (
    <div className="flex items-center gap-2 px-4 py-3 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      {/* Primary Actions */}
      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={onRefresh}
        className="p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300"
        title="Refresh"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </motion.button>

      <div className="w-px h-6 bg-gray-300 dark:bg-gray-600"></div>

      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={onNewFolder}
        className="flex items-center gap-2 px-3 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 text-sm font-medium"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
        </svg>
        New Folder
      </motion.button>

      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={onUpload}
        className="flex items-center gap-2 px-3 py-2 bg-green-500 text-white rounded hover:bg-green-600 text-sm font-medium"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
        Upload
      </motion.button>

      <div className="w-px h-6 bg-gray-300 dark:bg-gray-600"></div>

      {/* File Operations */}
      <motion.button
        whileHover={{ scale: canDownload ? 1.05 : 1 }}
        whileTap={{ scale: canDownload ? 0.95 : 1 }}
        onClick={canDownload ? onDownload : undefined}
        disabled={!canDownload}
        className={`p-2 rounded ${
          canDownload
            ? 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
            : 'text-gray-400 dark:text-gray-600 cursor-not-allowed'
        }`}
        title="Download"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
      </motion.button>

      <motion.button
        whileHover={{ scale: canDelete ? 1.05 : 1 }}
        whileTap={{ scale: canDelete ? 0.95 : 1 }}
        onClick={canDelete ? onDelete : undefined}
        disabled={!canDelete}
        className={`p-2 rounded ${
          canDelete
            ? 'hover:bg-red-100 dark:hover:bg-red-900/20 text-red-600 dark:text-red-400'
            : 'text-gray-400 dark:text-gray-600 cursor-not-allowed'
        }`}
        title="Delete"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </motion.button>

      {canPermissions && (
        <motion.button
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          onClick={onPermissions}
          className="p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300"
          title="Permissions"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </motion.button>
      )}

      <div className="flex-1"></div>

      {/* View Options */}
      <div className="flex items-center gap-1 bg-gray-100 dark:bg-gray-700 rounded p-1">
        <button
          onClick={() => onViewModeChange('list')}
          className={`p-1.5 rounded ${
            viewMode === 'list'
              ? 'bg-white dark:bg-gray-600 text-blue-600 dark:text-blue-400'
              : 'text-gray-600 dark:text-gray-400'
          }`}
          title="List View"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <button
          onClick={() => onViewModeChange('grid')}
          className={`p-1.5 rounded ${
            viewMode === 'grid'
              ? 'bg-white dark:bg-gray-600 text-blue-600 dark:text-blue-400'
              : 'text-gray-600 dark:text-gray-400'
          }`}
          title="Grid View"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
        </button>
      </div>

      <button
        onClick={onToggleHidden}
        className={`p-2 rounded ${
          showHidden
            ? 'bg-blue-100 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400'
            : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400'
        }`}
        title="Toggle Hidden Files"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
      </button>

      {selectedCount > 0 && (
        <div className="ml-2 px-3 py-1 bg-blue-100 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 text-sm rounded-full">
          {selectedCount} selected
        </div>
      )}
    </div>
  );
};

export default Toolbar;
