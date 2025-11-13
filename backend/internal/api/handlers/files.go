package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	mw "github.com/Stumpf-works/stumpfworks-nas/internal/api/middleware"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/files"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

var (
	fileService   *files.Service
	uploadManager *files.UploadManager
)

// InitFileService initializes the file service with allowed paths from shares
func InitFileService() error {
	// Get all shares to determine allowed paths
	var shares []*models.Share
	if err := database.DB.Find(&shares).Error; err != nil {
		return err
	}

	// Extract paths
	allowedPaths := make([]string, len(shares))
	for i, share := range shares {
		allowedPaths[i] = share.Path
	}

	// Create permission checker
	permChecker := files.NewPermissionChecker(shares)

	// Initialize file service
	fileService = files.NewService(allowedPaths, permChecker)

	// Initialize upload manager
	uploadManager = files.NewUploadManager("/tmp/stumpfworks-uploads")

	logger.Info("File service initialized", zap.Int("allowedPaths", len(allowedPaths)))
	return nil
}

// getSecurityContext extracts security context from request
func getSecurityContext(r *http.Request) (*files.SecurityContext, error) {
	// Get user from context (set by auth middleware)
	user := mw.GetUserFromContext(r.Context())
	if user == nil {
		return nil, errors.Unauthorized("User not authenticated", nil)
	}

	// Get all shares
	var shares []*models.Share
	if err := database.DB.Find(&shares).Error; err != nil {
		return nil, errors.InternalServerError("Failed to load shares", err)
	}

	// Determine allowed paths based on user role and share permissions
	allowedPaths := files.GetAllowedPathsForUser(user, shares)

	return &files.SecurityContext{
		User:         user,
		IsAdmin:      user.Role == "admin",
		AllowedPaths: allowedPaths,
	}, nil
}

// ===== File Browsing Handlers =====

// BrowseFiles lists files in a directory
func BrowseFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path parameter", nil))
		return
	}

	showHidden := r.URL.Query().Get("showHidden") == "true"

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	req := &files.BrowseRequest{
		Path:       path,
		ShowHidden: showHidden,
	}

	result, err := fileService.Browse(ctx, req)
	if err != nil {
		logger.Error("Failed to browse files", zap.String("path", path), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, result)
}

// GetFileInfo returns information about a specific file
func GetFileInfo(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path parameter", nil))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	info, err := fileService.GetFileInfo(ctx, path)
	if err != nil {
		logger.Error("Failed to get file info", zap.String("path", path), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, info)
}

// ===== File Upload Handlers =====

// UploadFile handles file uploads (simple single-file upload)
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 100MB in memory)
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		utils.RespondError(w, errors.BadRequest("Failed to parse multipart form", err))
		return
	}

	// Get destination path
	destPath := r.FormValue("path")
	if destPath == "" {
		utils.RespondError(w, errors.BadRequest("Missing destination path", nil))
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Failed to get file from form", err))
		return
	}
	defer file.Close()

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Upload file
	if err := fileService.UploadSingleFile(ctx, destPath, file, header); err != nil {
		logger.Error("Failed to upload file", zap.String("filename", header.Filename), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message":  "File uploaded successfully",
		"filename": header.Filename,
	})
}

// StartChunkedUpload starts a chunked upload session
func StartChunkedUpload(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FileName  string `json:"fileName"`
		TotalSize int64  `json:"totalSize"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	session, err := uploadManager.StartUploadSession(req.FileName, req.TotalSize)
	if err != nil {
		logger.Error("Failed to start upload session", zap.String("fileName", req.FileName), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, session)
}

// UploadChunk uploads a chunk of a file
func UploadChunk(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	chunkIndexStr := chi.URLParam(r, "chunkIndex")

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid chunk index", err))
		return
	}

	// Read chunk data from request body
	if err := uploadManager.UploadChunk(sessionID, chunkIndex, r.Body); err != nil {
		logger.Error("Failed to upload chunk", zap.String("sessionID", sessionID), zap.Int("chunkIndex", chunkIndex), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Chunk uploaded successfully",
	})
}

// FinalizeUpload finalizes a chunked upload
func FinalizeUpload(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID       string `json:"sessionId"`
		DestinationPath string `json:"destinationPath"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Check write permissions on destination directory BEFORE finalizing
	// This prevents users from uploading chunks to paths they don't have access to
	destPath := req.DestinationPath
	if destPath == "" {
		utils.RespondError(w, errors.BadRequest("Destination path is required", nil))
		return
	}

	// Check permissions on the destination directory (parent of the file)
	if err := fileService.CheckWritePermission(ctx, destPath); err != nil {
		// Permission denied - clean up the upload session
		uploadManager.CancelUpload(req.SessionID)
		logger.Warn("Upload blocked due to insufficient permissions",
			zap.String("sessionID", req.SessionID),
			zap.String("user", ctx.User.Username),
			zap.String("destinationPath", destPath),
			zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	if err := uploadManager.FinalizeUpload(req.SessionID, req.DestinationPath); err != nil {
		logger.Error("Failed to finalize upload", zap.String("sessionID", req.SessionID), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	logger.Info("Upload finalized", zap.String("sessionID", req.SessionID), zap.String("user", ctx.User.Username))

	utils.RespondSuccess(w, map[string]string{
		"message": "Upload completed successfully",
	})
}

// CancelUpload cancels an upload session
func CancelUpload(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	if err := uploadManager.CancelUpload(sessionID); err != nil {
		logger.Error("Failed to cancel upload", zap.String("sessionID", sessionID), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Upload cancelled",
	})
}

// GetUploadSession returns information about an upload session
func GetUploadSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	session, err := uploadManager.GetUploadSession(sessionID)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, session)
}

// ===== File Download Handler =====

// DownloadFile handles file downloads
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path parameter", nil))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Get file info first
	info, err := fileService.GetFileInfo(ctx, path)
	if err != nil {
		logger.Error("Failed to get file info for download", zap.String("path", path), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	// Cannot download directories
	if info.IsDir {
		utils.RespondError(w, errors.BadRequest("Cannot download directory (create archive first)", nil))
		return
	}

	// Open file
	file, err := http.Dir("/").Open(path)
	if err != nil {
		logger.Error("Failed to open file for download", zap.String("path", path), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to open file", err))
		return
	}
	defer file.Close()

	// Set headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", info.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size))

	// Stream file
	if _, err := io.Copy(w, file); err != nil {
		logger.Error("Failed to stream file", zap.String("path", path), zap.Error(err))
	}

	logger.Info("File downloaded", zap.String("path", path), zap.String("user", ctx.User.Username))
}

// ===== Directory Operations =====

// CreateDirectory creates a new directory
func CreateDirectory(w http.ResponseWriter, r *http.Request) {
	var req files.CreateDirRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	if err := fileService.CreateDirectory(ctx, &req); err != nil {
		logger.Error("Failed to create directory", zap.String("path", req.Path), zap.String("name", req.Name), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Directory created successfully",
	})
}

// ===== File Operations =====

// DeleteFiles deletes files or directories
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
	var req files.DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	if err := fileService.Delete(ctx, &req); err != nil {
		logger.Error("Failed to delete files", zap.Strings("paths", req.Paths), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Files deleted successfully",
	})
}

// RenameFile renames a file or directory
func RenameFile(w http.ResponseWriter, r *http.Request) {
	var req files.RenameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	if err := fileService.Rename(ctx, &req); err != nil {
		logger.Error("Failed to rename file", zap.String("oldPath", req.OldPath), zap.String("newName", req.NewName), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "File renamed successfully",
	})
}

// CopyFiles copies files or directories
func CopyFiles(w http.ResponseWriter, r *http.Request) {
	var req files.CopyMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	if err := fileService.Copy(ctx, &req); err != nil {
		logger.Error("Failed to copy files", zap.String("source", req.Source), zap.String("destination", req.Destination), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Files copied successfully",
	})
}

// MoveFiles moves files or directories
func MoveFiles(w http.ResponseWriter, r *http.Request) {
	var req files.CopyMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	if err := fileService.Move(ctx, &req); err != nil {
		logger.Error("Failed to move files", zap.String("source", req.Source), zap.String("destination", req.Destination), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Files moved successfully",
	})
}

// ===== Permissions Handlers =====

// GetFilePermissions returns file permissions
func GetFilePermissions(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path parameter", nil))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	permissions, err := fileService.GetPermissions(ctx, path)
	if err != nil {
		logger.Error("Failed to get permissions", zap.String("path", path), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, permissions)
}

// ChangeFilePermissions changes file permissions
func ChangeFilePermissions(w http.ResponseWriter, r *http.Request) {
	var req files.PermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	if err := fileService.ChangePermissions(ctx, &req); err != nil {
		logger.Error("Failed to change permissions", zap.String("path", req.Path), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Permissions changed successfully",
	})
}

// GetDiskUsage returns disk usage information
func GetDiskUsage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path parameter", nil))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	usage, err := fileService.GetDiskUsage(ctx, path)
	if err != nil {
		logger.Error("Failed to get disk usage", zap.String("path", path), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, usage)
}

// ===== Archive Handlers =====

// CreateArchive creates a compressed archive
func CreateArchive(w http.ResponseWriter, r *http.Request) {
	var req files.ArchiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	if err := fileService.CreateArchive(ctx, &req); err != nil {
		logger.Error("Failed to create archive", zap.Strings("paths", req.Paths), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Archive created successfully",
	})
}

// ExtractArchive extracts a compressed archive
func ExtractArchive(w http.ResponseWriter, r *http.Request) {
	var req files.ExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	ctx, err := getSecurityContext(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	if err := fileService.ExtractArchive(ctx, &req); err != nil {
		logger.Error("Failed to extract archive", zap.String("archivePath", req.ArchivePath), zap.Error(err))
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Archive extracted successfully",
	})
}
