import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { browseFiles, FileInfo } from '@/api/files';
import Button from '@/components/ui/Button';

interface FolderPickerProps {
  value: string;
  onChange: (path: string) => void;
  label?: string;
  placeholder?: string;
  required?: boolean;
}

export default function FolderPicker({ value, onChange, label, placeholder, required }: FolderPickerProps) {
  const [showBrowser, setShowBrowser] = useState(false);

  return (
    <div>
      {label && (
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      <div className="flex space-x-2">
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          required={required}
          className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-macos-dark-200 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-macos-blue"
        />
        <Button
          type="button"
          onClick={() => setShowBrowser(true)}
          variant="secondary"
        >
          📁 Browse
        </Button>
      </div>

      <AnimatePresence>
        {showBrowser && (
          <FolderBrowserModal
            currentPath={value || '/'}
            onSelect={(path) => {
              onChange(path);
              setShowBrowser(false);
            }}
            onClose={() => setShowBrowser(false)}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

interface FolderBrowserModalProps {
  currentPath: string;
  onSelect: (path: string) => void;
  onClose: () => void;
}

function FolderBrowserModal({ currentPath, onSelect, onClose }: FolderBrowserModalProps) {
  const [path, setPath] = useState(currentPath);
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadFolder(path);
  }, [path]);

  const loadFolder = async (folderPath: string) => {
    setLoading(true);
    setError('');
    try {
      const response = await browseFiles(folderPath, false);
      // Filter to only show directories
      const directories = response.files.filter((f) => f.isDir);
      setFiles(directories);
    } catch (err: any) {
      setError(err.message || 'Failed to load folder');
    } finally {
      setLoading(false);
    }
  };

  const navigateUp = () => {
    if (path === '/') return;
    const parts = path.split('/').filter(Boolean);
    parts.pop();
    const newPath = '/' + parts.join('/');
    setPath(newPath);
  };

  const navigateToPath = (targetPath: string) => {
    setPath(targetPath);
  };

  const getBreadcrumbs = () => {
    if (path === '/') return [{ name: 'Root', path: '/' }];

    const parts = path.split('/').filter(Boolean);
    const breadcrumbs = [{ name: 'Root', path: '/' }];

    let currentPath = '';
    parts.forEach((part) => {
      currentPath += '/' + part;
      breadcrumbs.push({ name: part, path: currentPath });
    });

    return breadcrumbs;
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="bg-white dark:bg-macos-dark-100 rounded-lg shadow-2xl w-full max-w-3xl max-h-[80vh] flex flex-col"
      >
        {/* Header */}
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            Select Folder
          </h2>

          {/* Breadcrumbs */}
          <div className="flex items-center space-x-2 text-sm overflow-x-auto">
            {getBreadcrumbs().map((crumb, index) => (
              <div key={crumb.path} className="flex items-center space-x-2">
                {index > 0 && (
                  <span className="text-gray-400 dark:text-gray-600">/</span>
                )}
                <button
                  onClick={() => navigateToPath(crumb.path)}
                  className="px-2 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 whitespace-nowrap"
                >
                  {crumb.name}
                </button>
              </div>
            ))}
          </div>

          {/* Current Path Display */}
          <div className="mt-3 p-2 bg-gray-100 dark:bg-macos-dark-200 rounded font-mono text-sm text-gray-900 dark:text-gray-100">
            {path}
          </div>
        </div>

        {/* File Browser */}
        <div className="flex-1 overflow-auto p-6">
          {/* Up Button */}
          {path !== '/' && (
            <button
              onClick={navigateUp}
              className="w-full p-3 mb-2 flex items-center space-x-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 text-left"
            >
              <span className="text-2xl">⬆️</span>
              <span className="text-gray-900 dark:text-gray-100 font-medium">
                Parent Directory
              </span>
            </button>
          )}

          {/* Loading State */}
          {loading && (
            <div className="flex justify-center items-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-macos-blue" />
            </div>
          )}

          {/* Error State */}
          {error && (
            <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400">
              {error}
            </div>
          )}

          {/* Folders List */}
          {!loading && !error && (
            <>
              {files.length === 0 ? (
                <div className="text-center py-12 text-gray-600 dark:text-gray-400">
                  <div className="text-6xl mb-4">📁</div>
                  <p className="text-lg font-medium">No subdirectories</p>
                  <p className="text-sm mt-2">This folder is empty or contains only files</p>
                </div>
              ) : (
                <div className="space-y-1">
                  {files.map((file) => (
                    <button
                      key={file.name}
                      onClick={() => setPath(file.path)}
                      className="w-full p-3 flex items-center space-x-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 text-left transition-colors"
                    >
                      <span className="text-2xl">📁</span>
                      <div className="flex-1 min-w-0">
                        <div className="text-gray-900 dark:text-gray-100 font-medium truncate">
                          {file.name}
                        </div>
                        <div className="text-xs text-gray-600 dark:text-gray-400 font-mono truncate">
                          {file.path}
                        </div>
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </>
          )}
        </div>

        {/* Footer */}
        <div className="p-6 border-t border-gray-200 dark:border-gray-700 flex justify-between items-center">
          <Button type="button" onClick={onClose} variant="secondary">
            Cancel
          </Button>
          <Button type="button" onClick={() => onSelect(path)}>
            Select "{path}"
          </Button>
        </div>
      </motion.div>
    </motion.div>
  );
}
