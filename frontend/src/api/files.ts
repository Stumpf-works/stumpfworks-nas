import client from './client';

// ===== Types =====

export interface FileInfo {
  name: string;
  path: string;
  size: number;
  isDir: boolean;
  modTime: string;
  permissions: string;
  owner: string;
  group: string;
  mimeType?: string;
  extension?: string;
  hasThumbnail: boolean;
}

export interface BrowseResponse {
  path: string;
  files: FileInfo[];
  totalSize: number;
  totalFiles: number;
  totalDirs: number;
}

export interface PermissionsInfo {
  path: string;
  permissions: string;
  mode: string;
  owner: string;
  group: string;
  uid: number;
  gid: number;
}

export interface DiskUsageInfo {
  path: string;
  totalSize: number;
  usedSize: number;
  freeSize: number;
  usagePercent: number;
}

export interface UploadSession {
  id: string;
  fileName: string;
  totalSize: number;
  uploadedSize: number;
  chunkSize: number;
  chunks: boolean[];
  startTime: string;
  lastUpdate: string;
}

// ===== Browse & Info =====

export const browseFiles = async (path: string, showHidden: boolean = false): Promise<BrowseResponse> => {
  const response = await client.get('/files/browse', {
    params: { path, showHidden }
  });
  return response.data.data;
};

export const getFileInfo = async (path: string): Promise<FileInfo> => {
  const response = await client.get('/files/info', {
    params: { path }
  });
  return response.data.data;
};

export const getDiskUsage = async (path: string): Promise<DiskUsageInfo> => {
  const response = await client.get('/files/usage', {
    params: { path }
  });
  return response.data.data;
};

// ===== File Operations =====

export const createDirectory = async (path: string, name: string, permissions?: string): Promise<void> => {
  await client.post('/files/mkdir', {
    path,
    name,
    permissions
  });
};

export const deleteFiles = async (paths: string[], recursive: boolean = false): Promise<void> => {
  await client.delete('/files/delete', {
    data: { paths, recursive }
  });
};

export const renameFile = async (oldPath: string, newName: string): Promise<void> => {
  await client.post('/files/rename', {
    oldPath,
    newName
  });
};

export const copyFiles = async (source: string, destination: string, overwrite: boolean = false): Promise<void> => {
  await client.post('/files/copy', {
    source,
    destination,
    overwrite
  });
};

export const moveFiles = async (source: string, destination: string, overwrite: boolean = false): Promise<void> => {
  await client.post('/files/move', {
    source,
    destination,
    overwrite
  });
};

// ===== Upload =====

export const uploadFile = async (
  path: string,
  file: File,
  onProgress?: (progress: number) => void
): Promise<void> => {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('path', path);

  await client.post('/files/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const progress = (progressEvent.loaded / progressEvent.total) * 100;
        onProgress(progress);
      }
    }
  });
};

// ===== Chunked Upload =====

export const startChunkedUpload = async (fileName: string, totalSize: number): Promise<UploadSession> => {
  const response = await client.post('/files/upload/start', {
    fileName,
    totalSize
  });
  return response.data.data;
};

export const uploadChunk = async (
  sessionId: string,
  chunkIndex: number,
  chunk: Blob,
  onProgress?: (progress: number) => void
): Promise<void> => {
  await client.post(`/files/upload/${sessionId}/chunk/${chunkIndex}`, chunk, {
    headers: {
      'Content-Type': 'application/octet-stream'
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const progress = (progressEvent.loaded / progressEvent.total) * 100;
        onProgress(progress);
      }
    }
  });
};

export const finalizeUpload = async (sessionId: string, destinationPath: string): Promise<void> => {
  await client.post('/files/upload/finalize', {
    sessionId,
    destinationPath
  });
};

export const cancelUpload = async (sessionId: string): Promise<void> => {
  await client.delete(`/files/upload/${sessionId}`);
};

export const getUploadSession = async (sessionId: string): Promise<UploadSession> => {
  const response = await client.get(`/files/upload/${sessionId}`);
  return response.data.data;
};

// ===== Download =====

export const downloadFile = (path: string): void => {
  // Create a download link
  const token = localStorage.getItem('token');
  const url = `${client.defaults.baseURL}/files/download?path=${encodeURIComponent(path)}`;

  // Create temporary link and trigger download
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', '');

  // Add authorization header via fetch
  fetch(url, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  })
    .then(response => response.blob())
    .then(blob => {
      const blobUrl = window.URL.createObjectURL(blob);
      link.href = blobUrl;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(blobUrl);
    });
};

// ===== Permissions =====

export const getFilePermissions = async (path: string): Promise<PermissionsInfo> => {
  const response = await client.get('/files/permissions', {
    params: { path }
  });
  return response.data.data;
};

export const changeFilePermissions = async (
  path: string,
  permissions: string,
  owner?: string,
  group?: string,
  recursive: boolean = false
): Promise<void> => {
  await client.post('/files/permissions', {
    path,
    permissions,
    owner,
    group,
    recursive
  });
};

// ===== Archives =====

export const createArchive = async (
  paths: string[],
  outputPath: string,
  format: 'zip' | 'tar' | 'tar.gz' = 'zip'
): Promise<void> => {
  await client.post('/files/archive/create', {
    paths,
    outputPath,
    format
  });
};

export const extractArchive = async (archivePath: string, destination: string): Promise<void> => {
  await client.post('/files/archive/extract', {
    archivePath,
    destination
  });
};

// ===== Helpers =====

export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
};

export const getFileIcon = (file: FileInfo): string => {
  if (file.isDir) return '📁';

  const ext = file.extension?.toLowerCase();
  const mime = file.mimeType?.toLowerCase();

  // Images
  if (mime?.startsWith('image/')) return '🖼️';

  // Videos
  if (mime?.startsWith('video/')) return '🎬';

  // Audio
  if (mime?.startsWith('audio/')) return '🎵';

  // Documents
  if (ext === '.pdf') return '📄';
  if (['.doc', '.docx'].includes(ext || '')) return '📝';
  if (['.xls', '.xlsx'].includes(ext || '')) return '📊';
  if (['.ppt', '.pptx'].includes(ext || '')) return '📽️';

  // Code
  if (['.js', '.ts', '.tsx', '.jsx', '.go', '.py', '.java', '.c', '.cpp', '.rs'].includes(ext || '')) return '💻';

  // Archives
  if (['.zip', '.tar', '.gz', '.rar', '.7z'].includes(ext || '')) return '📦';

  // Text
  if (mime?.startsWith('text/')) return '📃';

  // Default
  return '📄';
};

export const isImageFile = (file: FileInfo): boolean => {
  return file.mimeType?.startsWith('image/') || false;
};

export const isVideoFile = (file: FileInfo): boolean => {
  return file.mimeType?.startsWith('video/') || false;
};

export const isTextFile = (file: FileInfo): boolean => {
  return file.mimeType?.startsWith('text/') || false;
};

export const canPreview = (file: FileInfo): boolean => {
  return isImageFile(file) || isVideoFile(file) || isTextFile(file) || file.extension === '.pdf';
};
