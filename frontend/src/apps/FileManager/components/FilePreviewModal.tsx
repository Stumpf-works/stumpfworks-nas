import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FileInfo, isImageFile, isVideoFile, isTextFile } from '@/api/files';

interface FilePreviewModalProps {
  file: FileInfo;
  onClose: () => void;
}

const FilePreviewModal: React.FC<FilePreviewModalProps> = ({ file, onClose }) => {
  const [imageError, setImageError] = useState(false);
  const [textContent, setTextContent] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const fileUrl = `/api/v1/files/download?path=${encodeURIComponent(file.path)}`;

  // Load text content for text files
  useEffect(() => {
    if (isTextFile(file) && file.size < 1024 * 1024) { // Max 1MB for text preview
      setLoading(true);
      fetch(fileUrl)
        .then(res => res.text())
        .then(text => {
          setTextContent(text);
          setLoading(false);
        })
        .catch(err => {
          console.error('Failed to load text:', err);
          setLoading(false);
        });
    }
  }, [file, fileUrl]);

  const renderPreview = () => {
    // Image Preview
    if (isImageFile(file)) {
      return (
        <div className="relative w-full h-full flex items-center justify-center bg-gray-100 dark:bg-gray-900 rounded-lg">
          {!imageError ? (
            <img
              src={fileUrl}
              alt={file.name}
              className="max-w-full max-h-full object-contain"
              onError={() => setImageError(true)}
            />
          ) : (
            <div className="text-center p-12">
              <div className="text-6xl mb-4">🖼️</div>
              <p className="text-gray-600 dark:text-gray-400">Failed to load image</p>
            </div>
          )}
        </div>
      );
    }

    // Video Preview
    if (isVideoFile(file)) {
      return (
        <div className="relative w-full h-full flex items-center justify-center bg-black rounded-lg">
          <video
            src={fileUrl}
            controls
            className="max-w-full max-h-full"
          >
            Your browser does not support video playback.
          </video>
        </div>
      );
    }

    // Text Preview
    if (isTextFile(file)) {
      if (loading) {
        return (
          <div className="flex items-center justify-center h-full">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
          </div>
        );
      }

      return (
        <div className="w-full h-full bg-gray-50 dark:bg-gray-900 rounded-lg p-4 overflow-auto">
          <pre className="text-sm text-gray-800 dark:text-gray-200 font-mono whitespace-pre-wrap break-words">
            {textContent || 'Failed to load content'}
          </pre>
        </div>
      );
    }

    // PDF Preview
    if (file.mimeType === 'application/pdf' || file.name.toLowerCase().endsWith('.pdf')) {
      return (
        <div className="w-full h-full bg-gray-100 dark:bg-gray-900 rounded-lg">
          <iframe
            src={fileUrl}
            className="w-full h-full rounded-lg"
            title={file.name}
          />
        </div>
      );
    }

    // Default: No preview available
    return (
      <div className="flex flex-col items-center justify-center h-full text-center p-12">
        <div className="text-6xl mb-4">📄</div>
        <p className="text-gray-600 dark:text-gray-400 mb-2">No preview available</p>
        <p className="text-sm text-gray-500 dark:text-gray-500 mb-6">
          {file.mimeType || 'Unknown file type'}
        </p>
        <a
          href={fileUrl}
          download={file.name}
          className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          Download File
        </a>
      </div>
    );
  };

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4"
        onClick={onClose}
      >
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
          onClick={(e) => e.stopPropagation()}
          className="bg-white dark:bg-macos-dark-100 rounded-xl shadow-2xl w-full max-w-6xl h-[90vh] flex flex-col"
        >
          {/* Header */}
          <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <div className="flex items-center space-x-3 flex-1 min-w-0">
              <span className="text-3xl">{isImageFile(file) ? '🖼️' : isVideoFile(file) ? '🎬' : isTextFile(file) ? '📝' : '📄'}</span>
              <div className="flex-1 min-w-0">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-white truncate">
                  {file.name}
                </h2>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {file.mimeType || 'Unknown'} • {Math.round(file.size / 1024)} KB
                </p>
              </div>
            </div>

            <div className="flex items-center space-x-2">
              <a
                href={fileUrl}
                download={file.name}
                className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
                title="Download"
              >
                <svg className="w-6 h-6 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </a>
              <button
                onClick={onClose}
                className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
                title="Close"
              >
                <svg className="w-6 h-6 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          {/* Preview Content */}
          <div className="flex-1 overflow-hidden p-6">
            {renderPreview()}
          </div>

          {/* Footer */}
          <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-macos-dark-200">
            <div className="flex items-center justify-between text-sm text-gray-600 dark:text-gray-400">
              <div>
                Path: <span className="font-mono text-xs">{file.path}</span>
              </div>
              <div>
                Modified: {new Date(file.modTime).toLocaleString()}
              </div>
            </div>
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
};

export default FilePreviewModal;
