import React, { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import { FileInfo, browseFiles, BrowseResponse } from '@/api/files';
import Toolbar from './components/Toolbar';
import Breadcrumbs from './components/Breadcrumbs';
import FileBrowser from './components/FileBrowser';
import StatusBar from './components/StatusBar';
import FilePreviewModal from './components/FilePreviewModal';
import NewFolderModal from './components/NewFolderModal';
import UploadModal from './components/UploadModal';
import PermissionsModal from './components/PermissionsModal';
import { useAuthStore } from '@/store';

const FileManager: React.FC = () => {
  const user = useAuthStore((state) => state.user);
  const [currentPath, setCurrentPath] = useState<string>('/');
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [selectedFiles, setSelectedFiles] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [browseData, setBrowseData] = useState<BrowseResponse | null>(null);
  const [viewMode, setViewMode] = useState<'list' | 'grid'>('list');
  const [showHidden, setShowHidden] = useState<boolean>(false);

  // Modals
  const [previewFile, setPreviewFile] = useState<FileInfo | null>(null);
  const [showNewFolderModal, setShowNewFolderModal] = useState<boolean>(false);
  const [showUploadModal, setShowUploadModal] = useState<boolean>(false);
  const [permissionsFile, setPermissionsFile] = useState<FileInfo | null>(null);

  // Load files for current path
  const loadFiles = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await browseFiles(currentPath, showHidden);
      setBrowseData(data);
      setFiles(data.files);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load files');
      console.error('Failed to load files:', err);
    } finally {
      setLoading(false);
    }
  }, [currentPath, showHidden]);

  useEffect(() => {
    loadFiles();
  }, [loadFiles]);

  // Navigate to a path
  const navigateTo = (path: string) => {
    setCurrentPath(path);
    setSelectedFiles(new Set());
  };

  // Handle file/folder click
  const handleFileClick = (file: FileInfo, event: React.MouseEvent) => {
    if (file.isDir) {
      // Navigate into directory
      navigateTo(file.path);
    } else {
      // Handle multi-selection
      if (event.ctrlKey || event.metaKey) {
        const newSelected = new Set(selectedFiles);
        if (newSelected.has(file.path)) {
          newSelected.delete(file.path);
        } else {
          newSelected.add(file.path);
        }
        setSelectedFiles(newSelected);
      } else if (event.shiftKey && selectedFiles.size > 0) {
        // TODO: Implement shift-click range selection
        const newSelected = new Set(selectedFiles);
        newSelected.add(file.path);
        setSelectedFiles(newSelected);
      } else {
        // Single selection
        setSelectedFiles(new Set([file.path]));
      }
    }
  };

  // Handle file double-click
  const handleFileDoubleClick = (file: FileInfo) => {
    if (file.isDir) {
      navigateTo(file.path);
    } else {
      // Open preview for files
      setPreviewFile(file);
    }
  };

  // Handle file operations
  const handleRefresh = () => {
    loadFiles();
  };

  const handleNewFolder = () => {
    setShowNewFolderModal(true);
  };

  const handleUpload = () => {
    setShowUploadModal(true);
  };

  const handleDownload = () => {
    if (selectedFiles.size === 1) {
      const filePath = Array.from(selectedFiles)[0];
      // Download logic is in the API client
      const { downloadFile } = require('../../api/files');
      downloadFile(filePath);
    }
  };

  const handleDelete = async () => {
    if (selectedFiles.size === 0) return;

    if (!confirm(`Delete ${selectedFiles.size} item(s)?`)) return;

    try {
      const { deleteFiles } = require('../../api/files');
      await deleteFiles(Array.from(selectedFiles), true);
      setSelectedFiles(new Set());
      loadFiles();
    } catch (err: any) {
      alert('Failed to delete: ' + (err.response?.data?.error?.message || err.message));
    }
  };

  const handlePermissions = () => {
    if (selectedFiles.size === 1) {
      const filePath = Array.from(selectedFiles)[0];
      const file = files.find(f => f.path === filePath);
      if (file) {
        setPermissionsFile(file);
      }
    }
  };

  const handleToggleHidden = () => {
    setShowHidden(!showHidden);
  };

  return (
    <div className="flex flex-col h-full bg-gray-50 dark:bg-gray-900">
      {/* Toolbar */}
      <Toolbar
        onRefresh={handleRefresh}
        onNewFolder={handleNewFolder}
        onUpload={handleUpload}
        onDownload={handleDownload}
        onDelete={handleDelete}
        onPermissions={handlePermissions}
        onToggleHidden={handleToggleHidden}
        showHidden={showHidden}
        viewMode={viewMode}
        onViewModeChange={setViewMode}
        selectedCount={selectedFiles.size}
        canDelete={selectedFiles.size > 0}
        canDownload={selectedFiles.size === 1 && !files.find(f => f.path === Array.from(selectedFiles)[0])?.isDir}
        canPermissions={user?.role === 'admin' && selectedFiles.size === 1}
      />

      {/* Breadcrumbs */}
      <Breadcrumbs
        currentPath={currentPath}
        onNavigate={navigateTo}
      />

      {/* File Browser */}
      <div className="flex-1 overflow-auto">
        {loading ? (
          <div className="flex items-center justify-center h-full">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
          </div>
        ) : error ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <p className="text-red-500 dark:text-red-400 mb-4">{error}</p>
              <button
                onClick={handleRefresh}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                Retry
              </button>
            </div>
          </div>
        ) : (
          <FileBrowser
            files={files}
            selectedFiles={selectedFiles}
            viewMode={viewMode}
            onFileClick={handleFileClick}
            onFileDoubleClick={handleFileDoubleClick}
            onSelectionChange={setSelectedFiles}
            currentPath={currentPath}
            onRefresh={loadFiles}
          />
        )}
      </div>

      {/* Status Bar */}
      {browseData && (
        <StatusBar
          totalFiles={browseData.totalFiles}
          totalDirs={browseData.totalDirs}
          totalSize={browseData.totalSize}
          selectedCount={selectedFiles.size}
        />
      )}

      {/* Modals */}
      {previewFile && (
        <FilePreviewModal
          file={previewFile}
          onClose={() => setPreviewFile(null)}
        />
      )}

      {showNewFolderModal && (
        <NewFolderModal
          currentPath={currentPath}
          onClose={() => setShowNewFolderModal(false)}
          onSuccess={loadFiles}
        />
      )}

      {showUploadModal && (
        <UploadModal
          currentPath={currentPath}
          onClose={() => setShowUploadModal(false)}
          onSuccess={loadFiles}
        />
      )}

      {permissionsFile && (
        <PermissionsModal
          file={permissionsFile}
          onClose={() => setPermissionsFile(null)}
          onSuccess={loadFiles}
        />
      )}
    </div>
  );
};

export default FileManager;
