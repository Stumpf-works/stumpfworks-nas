import React, { useState, useEffect, useCallback } from 'react';
import { FileInfo, browseFiles, BrowseResponse, uploadFile, downloadFile, deleteFiles, copyFiles, moveFiles } from '@/api/files';
import Toolbar from './components/Toolbar';
import Breadcrumbs from './components/Breadcrumbs';
import FileBrowser from './components/FileBrowser';
import StatusBar from './components/StatusBar';
import FilePreviewModal from './components/FilePreviewModal';
import NewFolderModal from './components/NewFolderModal';
import UploadModal from './components/UploadModal';
import PermissionsModal from './components/PermissionsModal';
import ContextMenu from './components/ContextMenu';
import DropZone from './components/DropZone';
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
  const [searchQuery, setSearchQuery] = useState<string>('');

  // Modals
  const [previewFile, setPreviewFile] = useState<FileInfo | null>(null);
  const [showNewFolderModal, setShowNewFolderModal] = useState<boolean>(false);
  const [showUploadModal, setShowUploadModal] = useState<boolean>(false);
  const [permissionsFile, setPermissionsFile] = useState<FileInfo | null>(null);

  // Context Menu
  const [contextMenu, setContextMenu] = useState<{ x: number; y: number; file: FileInfo | null } | null>(null);

  // Clipboard for Copy/Cut/Paste
  const [clipboard, setClipboard] = useState<{ files: string[]; operation: 'copy' | 'cut' } | null>(null);

  // Filter files based on search query
  const filteredFiles = React.useMemo(() => {
    if (!searchQuery.trim()) {
      return files;
    }
    const query = searchQuery.toLowerCase();
    return files.filter(file =>
      file.name.toLowerCase().includes(query)
    );
  }, [files, searchQuery]);

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
      downloadFile(filePath);
    }
  };

  const handleDelete = async () => {
    if (selectedFiles.size === 0) return;

    if (!confirm(`Delete ${selectedFiles.size} item(s)?`)) return;

    try {
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

  // Drag & Drop Handler
  const handleFilesDropped = async (fileList: FileList) => {
    const filesArray = Array.from(fileList);
    console.log(`Dropped ${filesArray.length} files`);

    // Upload files one by one
    for (const file of filesArray) {
      try {
        await uploadFile(currentPath, file, (progress) => {
          console.log(`Uploading ${file.name}: ${progress}%`);
        });
      } catch (err: any) {
        alert(`Failed to upload ${file.name}: ${err.message}`);
        break;
      }
    }

    // Refresh after all uploads
    loadFiles();
  };

  // Context Menu Handlers
  const handleContextMenu = (event: React.MouseEvent, file: FileInfo | null) => {
    event.preventDefault();
    setContextMenu({ x: event.clientX, y: event.clientY, file });
  };

  const handleContextMenuOpen = () => {
    if (contextMenu?.file) {
      if (contextMenu.file.isDir) {
        navigateTo(contextMenu.file.path);
      } else {
        setPreviewFile(contextMenu.file);
      }
    }
  };

  const handleContextMenuDownload = () => {
    if (contextMenu?.file && !contextMenu.file.isDir) {
      downloadFile(contextMenu.file.path);
    }
  };

  const handleContextMenuDelete = async () => {
    if (contextMenu?.file) {
      if (confirm(`Delete "${contextMenu.file.name}"?`)) {
        try {
          await deleteFiles([contextMenu.file.path], true);
          loadFiles();
        } catch (err: any) {
          alert('Failed to delete: ' + err.message);
        }
      }
    }
  };

  const handleContextMenuPermissions = () => {
    if (contextMenu?.file) {
      setPermissionsFile(contextMenu.file);
    }
  };

  // Clipboard Operations
  const handleCopy = () => {
    if (selectedFiles.size > 0) {
      setClipboard({ files: Array.from(selectedFiles), operation: 'copy' });
      console.log(`Copied ${selectedFiles.size} items`);
    }
  };

  const handleCut = () => {
    if (selectedFiles.size > 0) {
      setClipboard({ files: Array.from(selectedFiles), operation: 'cut' });
      console.log(`Cut ${selectedFiles.size} items`);
    }
  };

  const handlePaste = async () => {
    if (!clipboard || clipboard.files.length === 0) return;

    try {
      for (const sourcePath of clipboard.files) {
        const fileName = sourcePath.split('/').pop() || '';
        const destination = `${currentPath}/${fileName}`;

        if (clipboard.operation === 'copy') {
          await copyFiles(sourcePath, destination, false);
        } else {
          await moveFiles(sourcePath, destination, false);
        }
      }

      if (clipboard.operation === 'cut') {
        setClipboard(null); // Clear clipboard after cut
      }

      loadFiles();
      console.log(`Pasted ${clipboard.files.length} items`);
    } catch (err: any) {
      alert('Paste failed: ' + err.message);
    }
  };

  const handleSelectAll = () => {
    const allPaths = new Set(filteredFiles.map(f => f.path));
    setSelectedFiles(allPaths);
  };

  const handleDeselectAll = () => {
    setSelectedFiles(new Set());
  };

  // Keyboard Shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Ignore if typing in input
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
        return;
      }

      const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
      const cmdOrCtrl = isMac ? e.metaKey : e.ctrlKey;

      // Ctrl+C / Cmd+C - Copy
      if (cmdOrCtrl && e.key === 'c') {
        e.preventDefault();
        handleCopy();
      }
      // Ctrl+X / Cmd+X - Cut
      else if (cmdOrCtrl && e.key === 'x') {
        e.preventDefault();
        handleCut();
      }
      // Ctrl+V / Cmd+V - Paste
      else if (cmdOrCtrl && e.key === 'v') {
        e.preventDefault();
        handlePaste();
      }
      // Ctrl+A / Cmd+A - Select All
      else if (cmdOrCtrl && e.key === 'a') {
        e.preventDefault();
        handleSelectAll();
      }
      // Delete / Backspace - Delete
      else if ((e.key === 'Delete' || e.key === 'Backspace') && selectedFiles.size > 0) {
        e.preventDefault();
        handleDelete();
      }
      // Escape - Deselect
      else if (e.key === 'Escape') {
        handleDeselectAll();
        setContextMenu(null);
      }
      // F5 - Refresh
      else if (e.key === 'F5') {
        e.preventDefault();
        handleRefresh();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [selectedFiles, clipboard, filteredFiles, currentPath]);

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
        onSearch={setSearchQuery}
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

      {/* File Browser with Drag & Drop */}
      <div className="flex-1 overflow-auto">
        <DropZone onFilesDropped={handleFilesDropped} disabled={loading}>
          <div
            className="h-full"
            onContextMenu={(e) => handleContextMenu(e, null)}
          >
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
                files={filteredFiles}
                selectedFiles={selectedFiles}
                viewMode={viewMode}
                onFileClick={handleFileClick}
                onFileDoubleClick={handleFileDoubleClick}
                onSelectionChange={setSelectedFiles}
                onContextMenu={handleContextMenu}
                currentPath={currentPath}
                onRefresh={loadFiles}
              />
            )}
          </div>
        </DropZone>
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

      {/* Context Menu */}
      {contextMenu && (
        <ContextMenu
          x={contextMenu.x}
          y={contextMenu.y}
          file={contextMenu.file}
          onClose={() => setContextMenu(null)}
          onOpen={handleContextMenuOpen}
          onDownload={handleContextMenuDownload}
          onCopy={handleCopy}
          onCut={handleCut}
          onDelete={handleContextMenuDelete}
          onPermissions={user?.role === 'admin' ? handleContextMenuPermissions : undefined}
          isAdmin={user?.role === 'admin'}
        />
      )}
    </div>
  );
};

export default FileManager;
