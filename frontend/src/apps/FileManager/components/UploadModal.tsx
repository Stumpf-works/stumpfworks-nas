import React, { useState } from 'react';
import { uploadFile } from '@/api/files';

interface UploadModalProps {
  currentPath: string;
  onClose: () => void;
  onSuccess: () => void;
}

const UploadModal: React.FC<UploadModalProps> = ({ currentPath, onClose, onSuccess }) => {
  const [files, setFiles] = useState<File[]>([]);
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState<number>(0);
  const [error, setError] = useState<string | null>(null);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFiles(Array.from(e.target.files));
    }
  };

  const handleUpload = async () => {
    if (files.length === 0) return;

    try {
      setUploading(true);
      setError(null);

      // Upload files one by one for simplicity
      for (let i = 0; i < files.length; i++) {
        await uploadFile(currentPath, files[i], (prog) => {
          const totalProgress = (i / files.length) * 100 + (prog / files.length);
          setProgress(totalProgress);
        });
      }

      onSuccess();
      onClose();
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to upload files');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">Upload Files</h2>

        <div className="mb-4">
          <input
            type="file"
            multiple
            onChange={handleFileSelect}
            className="w-full text-sm text-gray-500 dark:text-gray-400
              file:mr-4 file:py-2 file:px-4
              file:rounded file:border-0
              file:text-sm file:font-semibold
              file:bg-blue-50 file:text-blue-700
              dark:file:bg-blue-900/20 dark:file:text-blue-400
              hover:file:bg-blue-100 dark:hover:file:bg-blue-900/30"
          />
        </div>

        {files.length > 0 && (
          <div className="mb-4 p-3 bg-gray-50 dark:bg-gray-700 rounded">
            <p className="text-sm text-gray-700 dark:text-gray-300 mb-2">
              {files.length} file(s) selected
            </p>
            <div className="max-h-32 overflow-y-auto text-xs text-gray-500 dark:text-gray-400 space-y-1">
              {files.map((file, i) => (
                <div key={i}>{file.name}</div>
              ))}
            </div>
          </div>
        )}

        {uploading && (
          <div className="mb-4">
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
              <div
                className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">{Math.round(progress)}% uploaded</p>
          </div>
        )}

        {error && (
          <div className="mb-4 p-3 bg-red-100 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded">
            {error}
          </div>
        )}

        <div className="flex justify-end gap-2">
          <button
            onClick={onClose}
            disabled={uploading}
            className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={handleUpload}
            disabled={uploading || files.length === 0}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {uploading ? 'Uploading...' : 'Upload'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default UploadModal;
