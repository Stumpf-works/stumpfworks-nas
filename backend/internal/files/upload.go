package files

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

const (
	// MaxUploadSize is the maximum file size for uploads (10GB default)
	MaxUploadSize = 10 * 1024 * 1024 * 1024
	// ChunkSize is the default chunk size for chunked uploads (10MB)
	ChunkSize = 10 * 1024 * 1024
	// UploadSessionTimeout is how long to keep upload sessions alive
	UploadSessionTimeout = 24 * time.Hour
	// MinimumFreeSpace is the minimum free space required (1GB buffer)
	MinimumFreeSpace = 1 * 1024 * 1024 * 1024
)

// UploadManager manages file uploads
type UploadManager struct {
	sessions map[string]*UploadSession
	mu       sync.RWMutex
	tempDir  string
}

// NewUploadManager creates a new upload manager
func NewUploadManager(tempDir string) *UploadManager {
	return &UploadManager{
		sessions: make(map[string]*UploadSession),
		tempDir:  tempDir,
	}
}

// StartUploadSession starts a new chunked upload session
func (um *UploadManager) StartUploadSession(fileName string, totalSize int64) (*UploadSession, error) {
	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, errors.InternalServerError("Failed to generate session ID", err)
	}

	// Calculate number of chunks
	numChunks := int(totalSize / ChunkSize)
	if totalSize%ChunkSize != 0 {
		numChunks++
	}

	session := &UploadSession{
		ID:           sessionID,
		FileName:     fileName,
		TotalSize:    totalSize,
		UploadedSize: 0,
		ChunkSize:    ChunkSize,
		Chunks:       make([]bool, numChunks),
		StartTime:    time.Now(),
		LastUpdate:   time.Now(),
	}

	um.mu.Lock()
	um.sessions[sessionID] = session
	um.mu.Unlock()

	logger.Info("Upload session started", zap.String("sessionID", sessionID), zap.String("fileName", fileName))
	return session, nil
}

// UploadChunk uploads a chunk of a file
func (um *UploadManager) UploadChunk(sessionID string, chunkIndex int, reader io.Reader) error {
	um.mu.Lock()
	session, exists := um.sessions[sessionID]
	um.mu.Unlock()

	if !exists {
		return errors.NotFound("Upload session not found", nil)
	}

	// Validate chunk index
	if chunkIndex < 0 || chunkIndex >= len(session.Chunks) {
		return errors.BadRequest("Invalid chunk index", nil)
	}

	// Check if chunk already uploaded
	if session.Chunks[chunkIndex] {
		return errors.Conflict("Chunk already uploaded", nil)
	}

	// Create temp directory for this session
	sessionDir := filepath.Join(um.tempDir, sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return errors.InternalServerError("Failed to create temp directory", err)
	}

	// Write chunk to temp file
	chunkPath := filepath.Join(sessionDir, fmt.Sprintf("chunk_%d", chunkIndex))
	chunkFile, err := os.Create(chunkPath)
	if err != nil {
		return errors.InternalServerError("Failed to create chunk file", err)
	}
	defer chunkFile.Close()

	written, err := io.Copy(chunkFile, reader)
	if err != nil {
		return errors.InternalServerError("Failed to write chunk", err)
	}

	// Update session
	um.mu.Lock()
	session.Chunks[chunkIndex] = true
	session.UploadedSize += written
	session.LastUpdate = time.Now()
	um.mu.Unlock()

	logger.Debug("Chunk uploaded", zap.String("sessionID", sessionID), zap.Int("chunkIndex", chunkIndex))
	return nil
}

// FinalizeUpload combines all chunks into the final file
func (um *UploadManager) FinalizeUpload(sessionID, destinationPath string) error {
	um.mu.Lock()
	session, exists := um.sessions[sessionID]
	um.mu.Unlock()

	if !exists {
		return errors.NotFound("Upload session not found", nil)
	}

	// Check if all chunks are uploaded
	for i, uploaded := range session.Chunks {
		if !uploaded {
			return errors.BadRequest(fmt.Sprintf("Missing chunk: %d", i), nil)
		}
	}

	sessionDir := filepath.Join(um.tempDir, sessionID)

	// Create final file
	finalFile, err := os.Create(destinationPath)
	if err != nil {
		return errors.InternalServerError("Failed to create final file", err)
	}
	defer finalFile.Close()

	// Combine all chunks
	for i := 0; i < len(session.Chunks); i++ {
		chunkPath := filepath.Join(sessionDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return errors.InternalServerError(fmt.Sprintf("Failed to open chunk %d", i), err)
		}

		if _, err := io.Copy(finalFile, chunkFile); err != nil {
			chunkFile.Close()
			return errors.InternalServerError(fmt.Sprintf("Failed to copy chunk %d", i), err)
		}
		chunkFile.Close()
	}

	// Cleanup temp files
	os.RemoveAll(sessionDir)

	// Remove session
	um.mu.Lock()
	delete(um.sessions, sessionID)
	um.mu.Unlock()

	logger.Info("Upload finalized", zap.String("sessionID", sessionID), zap.String("destination", destinationPath))
	return nil
}

// CancelUpload cancels an upload session
func (um *UploadManager) CancelUpload(sessionID string) error {
	um.mu.Lock()
	_, exists := um.sessions[sessionID]
	if exists {
		delete(um.sessions, sessionID)
	}
	um.mu.Unlock()

	if !exists {
		return errors.NotFound("Upload session not found", nil)
	}

	// Cleanup temp files
	sessionDir := filepath.Join(um.tempDir, sessionID)
	os.RemoveAll(sessionDir)

	logger.Info("Upload cancelled", zap.String("sessionID", sessionID))
	return nil
}

// GetUploadSession returns information about an upload session
func (um *UploadManager) GetUploadSession(sessionID string) (*UploadSession, error) {
	um.mu.RLock()
	session, exists := um.sessions[sessionID]
	um.mu.RUnlock()

	if !exists {
		return nil, errors.NotFound("Upload session not found", nil)
	}

	return session, nil
}

// CleanupExpiredSessions removes expired upload sessions
func (um *UploadManager) CleanupExpiredSessions() {
	um.mu.Lock()
	defer um.mu.Unlock()

	now := time.Now()
	for sessionID, session := range um.sessions {
		if now.Sub(session.LastUpdate) > UploadSessionTimeout {
			// Cleanup temp files
			sessionDir := filepath.Join(um.tempDir, sessionID)
			os.RemoveAll(sessionDir)

			delete(um.sessions, sessionID)
			logger.Info("Expired upload session cleaned up", zap.String("sessionID", sessionID))
		}
	}
}

// UploadSingleFile handles a simple single-file upload
func (s *Service) UploadSingleFile(ctx *SecurityContext, destinationDir string, file multipart.File, header *multipart.FileHeader) error {
	// Validate filename
	if err := ValidateFileName(header.Filename); err != nil {
		return err
	}

	// Validate and sanitize destination
	cleanDest, err := s.validator.ValidateAndSanitize(destinationDir)
	if err != nil {
		return err
	}

	// Check write permissions
	if err := s.permissions.CanWrite(ctx, cleanDest); err != nil {
		return err
	}

	// Check file size limit
	if header.Size > MaxUploadSize {
		return errors.BadRequest(fmt.Sprintf("File too large: max size is %d bytes", MaxUploadSize), nil)
	}

	// Build destination path
	destPath := filepath.Join(cleanDest, header.Filename)

	// Check if file already exists
	if _, err := os.Stat(destPath); err == nil {
		return errors.Conflict("File already exists", nil)
	}

	// Check disk space before uploading (file size + 1GB buffer)
	if err := CheckDiskSpace(cleanDest, header.Size+MinimumFreeSpace); err != nil {
		return err
	}

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return errors.InternalServerError("Failed to create file", err)
	}
	defer destFile.Close()

	// Copy file data
	written, err := io.Copy(destFile, file)
	if err != nil {
		os.Remove(destPath) // Cleanup on error
		return errors.InternalServerError("Failed to write file", err)
	}

	logger.Info("File uploaded", zap.String("path", destPath), zap.Int64("size", written), zap.String("user", ctx.User.Username))
	return nil
}

// Helper: generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
