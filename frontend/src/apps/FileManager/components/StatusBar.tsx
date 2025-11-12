import React from 'react';
import { formatFileSize } from '@/api/files';

interface StatusBarProps {
  totalFiles: number;
  totalDirs: number;
  totalSize: number;
  selectedCount: number;
}

const StatusBar: React.FC<StatusBarProps> = ({ totalFiles, totalDirs, totalSize, selectedCount }) => {
  return (
    <div className="flex items-center justify-between px-4 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 text-xs text-gray-600 dark:text-gray-400">
      <div className="flex items-center gap-4">
        <span>{totalDirs} folder{totalDirs !== 1 ? 's' : ''}</span>
        <span>{totalFiles} file{totalFiles !== 1 ? 's' : ''}</span>
        {selectedCount > 0 && (
          <span className="text-blue-600 dark:text-blue-400 font-medium">
            {selectedCount} selected
          </span>
        )}
      </div>
      <div>
        Total: {formatFileSize(totalSize)}
      </div>
    </div>
  );
};

export default StatusBar;
