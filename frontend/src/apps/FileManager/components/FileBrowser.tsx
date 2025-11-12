import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { FileInfo, formatFileSize, getFileIcon } from '@/api/files';

interface FileBrowserProps {
  files: FileInfo[];
  selectedFiles: Set<string>;
  viewMode: 'list' | 'grid';
  onFileClick: (file: FileInfo, event: React.MouseEvent) => void;
  onFileDoubleClick: (file: FileInfo) => void;
  onSelectionChange: (selected: Set<string>) => void;
  onContextMenu?: (event: React.MouseEvent, file: FileInfo) => void;
  currentPath: string;
  onRefresh: () => void;
}

const FileBrowser: React.FC<FileBrowserProps> = ({
  files,
  selectedFiles,
  viewMode,
  onFileClick,
  onFileDoubleClick,
  onContextMenu,
  currentPath,
  onRefresh,
}) => {
  const [sortBy, setSortBy] = useState<'name' | 'size' | 'date'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  const sortedFiles = [...files].sort((a, b) => {
    // Directories first
    if (a.isDir && !b.isDir) return -1;
    if (!a.isDir && b.isDir) return 1;

    let compare = 0;
    switch (sortBy) {
      case 'name':
        compare = a.name.localeCompare(b.name);
        break;
      case 'size':
        compare = a.size - b.size;
        break;
      case 'date':
        compare = new Date(a.modTime).getTime() - new Date(b.modTime).getTime();
        break;
    }

    return sortOrder === 'asc' ? compare : -compare;
  });

  const handleSort = (column: 'name' | 'size' | 'date') => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortOrder('asc');
    }
  };

  if (files.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-gray-500 dark:text-gray-400">
        <svg className="w-24 h-24 mb-4 text-gray-300 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        <p className="text-lg font-medium mb-2">This folder is empty</p>
        <p className="text-sm">Upload files or create a new folder to get started</p>
      </div>
    );
  }

  if (viewMode === 'grid') {
    return (
      <div className="p-4 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-8 gap-4">
        {sortedFiles.map((file) => (
          <motion.div
            key={file.path}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            onClick={(e) => onFileClick(file, e)}
            onDoubleClick={() => onFileDoubleClick(file)}
            onContextMenu={(e) => onContextMenu?.(e, file)}
            className={`flex flex-col items-center p-4 rounded-lg border-2 cursor-pointer transition-colors ${
              selectedFiles.has(file.path)
                ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                : 'border-transparent hover:border-gray-300 dark:hover:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800'
            }`}
          >
            <div className="text-5xl mb-2">{getFileIcon(file)}</div>
            <div className="text-sm text-center break-all line-clamp-2 text-gray-700 dark:text-gray-300">
              {file.name}
            </div>
            {!file.isDir && (
              <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {formatFileSize(file.size)}
              </div>
            )}
          </motion.div>
        ))}
      </div>
    );
  }

  // List view
  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center px-4 py-2 bg-gray-100 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 text-sm font-medium text-gray-700 dark:text-gray-300">
        <div className="w-8"></div>
        <button
          onClick={() => handleSort('name')}
          className="flex-1 flex items-center gap-1 hover:text-blue-600 dark:hover:text-blue-400"
        >
          Name
          {sortBy === 'name' && (
            <svg className={`w-4 h-4 transform ${sortOrder === 'desc' ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
            </svg>
          )}
        </button>
        <button
          onClick={() => handleSort('size')}
          className="w-24 flex items-center gap-1 justify-end hover:text-blue-600 dark:hover:text-blue-400"
        >
          Size
          {sortBy === 'size' && (
            <svg className={`w-4 h-4 transform ${sortOrder === 'desc' ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
            </svg>
          )}
        </button>
        <button
          onClick={() => handleSort('date')}
          className="w-48 flex items-center gap-1 justify-end hover:text-blue-600 dark:hover:text-blue-400"
        >
          Modified
          {sortBy === 'date' && (
            <svg className={`w-4 h-4 transform ${sortOrder === 'desc' ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
            </svg>
          )}
        </button>
      </div>

      {/* File List */}
      <div className="flex-1 overflow-y-auto">
        {sortedFiles.map((file) => (
          <motion.div
            key={file.path}
            whileHover={{ backgroundColor: 'rgba(0, 0, 0, 0.02)' }}
            onClick={(e) => onFileClick(file, e)}
            onDoubleClick={() => onFileDoubleClick(file)}
            onContextMenu={(e) => onContextMenu?.(e, file)}
            className={`flex items-center px-4 py-2 border-b border-gray-100 dark:border-gray-800 cursor-pointer ${
              selectedFiles.has(file.path)
                ? 'bg-blue-50 dark:bg-blue-900/20'
                : 'hover:bg-gray-50 dark:hover:bg-gray-800/50'
            }`}
          >
            <div className="w-8 text-2xl">{getFileIcon(file)}</div>
            <div className="flex-1 truncate text-gray-700 dark:text-gray-300">{file.name}</div>
            <div className="w-24 text-right text-sm text-gray-500 dark:text-gray-400">
              {file.isDir ? '—' : formatFileSize(file.size)}
            </div>
            <div className="w-48 text-right text-sm text-gray-500 dark:text-gray-400">
              {new Date(file.modTime).toLocaleString()}
            </div>
          </motion.div>
        ))}
      </div>
    </div>
  );
};

export default FileBrowser;
