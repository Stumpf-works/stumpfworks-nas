import React from 'react';
import { FileInfo } from '../../../api/files';

interface PermissionsModalProps {
  file: FileInfo;
  onClose: () => void;
  onSuccess: () => void;
}

const PermissionsModal: React.FC<PermissionsModalProps> = ({ file, onClose, onSuccess }) => {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">Permissions</h2>

        <div className="mb-4">
          <p className="text-sm text-gray-700 dark:text-gray-300 mb-2">File: {file.name}</p>
          <p className="text-xs text-gray-500 dark:text-gray-400">Current permissions: {file.permissions}</p>
        </div>

        <div className="text-center py-8">
          <p className="text-gray-500 dark:text-gray-400">Permission editing coming soon...</p>
        </div>

        <div className="flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default PermissionsModal;
