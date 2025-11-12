package files

import (
	"time"
)

// FileInfo represents file or directory information
type FileInfo struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	IsDir        bool      `json:"isDir"`
	ModTime      time.Time `json:"modTime"`
	Permissions  string    `json:"permissions"`
	Owner        string    `json:"owner"`
	Group        string    `json:"group"`
	MimeType     string    `json:"mimeType,omitempty"`
	Extension    string    `json:"extension,omitempty"`
	HasThumbnail bool      `json:"hasThumbnail"`
}

// BrowseRequest represents a directory browsing request
type BrowseRequest struct {
	Path      string `json:"path"`
	ShareID   string `json:"shareId,omitempty"`
	ShowHidden bool   `json:"showHidden"`
}

// BrowseResponse represents the directory browsing response
type BrowseResponse struct {
	Path       string     `json:"path"`
	Files      []FileInfo `json:"files"`
	TotalSize  int64      `json:"totalSize"`
	TotalFiles int        `json:"totalFiles"`
	TotalDirs  int        `json:"totalDirs"`
}

// CreateDirRequest represents a directory creation request
type CreateDirRequest struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Permissions string `json:"permissions,omitempty"` // e.g., "0755"
}

// RenameRequest represents a file/directory rename request
type RenameRequest struct {
	OldPath string `json:"oldPath"`
	NewName string `json:"newName"`
}

// CopyMoveRequest represents a file/directory copy or move request
type CopyMoveRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Overwrite   bool   `json:"overwrite"`
}

// DeleteRequest represents a deletion request
type DeleteRequest struct {
	Paths     []string `json:"paths"`
	Recursive bool     `json:"recursive"`
}

// PermissionsRequest represents a permissions change request
type PermissionsRequest struct {
	Path        string `json:"path"`
	Permissions string `json:"permissions"` // e.g., "0644"
	Owner       string `json:"owner,omitempty"`
	Group       string `json:"group,omitempty"`
	Recursive   bool   `json:"recursive"`
}

// ArchiveRequest represents an archive creation request
type ArchiveRequest struct {
	Paths      []string `json:"paths"`
	OutputPath string   `json:"outputPath"`
	Format     string   `json:"format"` // "zip", "tar", "tar.gz"
}

// ExtractRequest represents an archive extraction request
type ExtractRequest struct {
	ArchivePath string `json:"archivePath"`
	Destination string `json:"destination"`
}

// SearchRequest represents a file search request
type SearchRequest struct {
	BasePath   string `json:"basePath"`
	Query      string `json:"query"`
	FileType   string `json:"fileType,omitempty"`   // e.g., "image", "video", "document"
	MinSize    int64  `json:"minSize,omitempty"`
	MaxSize    int64  `json:"maxSize,omitempty"`
	ModifiedAfter  *time.Time `json:"modifiedAfter,omitempty"`
	ModifiedBefore *time.Time `json:"modifiedBefore,omitempty"`
}

// UploadSession represents an active upload session
type UploadSession struct {
	ID          string    `json:"id"`
	FileName    string    `json:"fileName"`
	TotalSize   int64     `json:"totalSize"`
	UploadedSize int64    `json:"uploadedSize"`
	ChunkSize   int64     `json:"chunkSize"`
	Chunks      []bool    `json:"chunks"`
	StartTime   time.Time `json:"startTime"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

// DiskUsageInfo represents disk usage information for a path
type DiskUsageInfo struct {
	Path       string  `json:"path"`
	TotalSize  int64   `json:"totalSize"`
	UsedSize   int64   `json:"usedSize"`
	FreeSize   int64   `json:"freeSize"`
	UsagePercent float64 `json:"usagePercent"`
}
