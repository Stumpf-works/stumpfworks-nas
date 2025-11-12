import React from 'react';
import { FileInfo } from '@/api/files';

interface FilePreviewModalProps {
  file: FileInfo;
  onClose: () => void;
}

const FilePreviewModal: React.FC<FilePreviewModalProps> = ({ file, onClose }) => {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-2xl w-full mx-4" onClick={(e) => e.stopPropagation()}>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Preview</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-gray-400 mb-4">Preview for {file.name}</p>
          <p className="text-sm text-gray-400 dark:text-gray-500">Preview functionality coming soon...</p>
        </div>
      </div>
    </div>
  );
};

export default FilePreviewModal;
