package files

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// CreateArchive creates a compressed archive from specified files
func (s *Service) CreateArchive(ctx *SecurityContext, req *ArchiveRequest) error {
	// Validate paths
	sourcePaths, err := s.validator.ValidatePaths(req.Paths)
	if err != nil {
		return err
	}

	// Check read permissions for all sources
	for _, path := range sourcePaths {
		if err := s.permissions.CanAccess(ctx, path); err != nil {
			return err
		}
	}

	// Validate output path
	outputPath, err := s.validator.ValidateAndSanitize(req.OutputPath)
	if err != nil {
		return err
	}

	// Check write permissions for output
	outputDir := filepath.Dir(outputPath)
	if err := s.permissions.CanWrite(ctx, outputDir); err != nil {
		return err
	}

	// Create archive based on format
	switch req.Format {
	case "zip":
		err = s.createZipArchive(sourcePaths, outputPath)
	case "tar":
		err = s.createTarArchive(sourcePaths, outputPath, false)
	case "tar.gz", "tgz":
		err = s.createTarArchive(sourcePaths, outputPath, true)
	default:
		return errors.BadRequest("Unsupported archive format (supported: zip, tar, tar.gz)", nil)
	}

	if err != nil {
		return err
	}

	logger.Info("Archive created", zap.String("output", outputPath), zap.String("format", req.Format), zap.String("user", ctx.User.Username))
	return nil
}

// ExtractArchive extracts a compressed archive
func (s *Service) ExtractArchive(ctx *SecurityContext, req *ExtractRequest) error {
	// Validate archive path
	archivePath, err := s.validator.ValidateAndSanitize(req.ArchivePath)
	if err != nil {
		return err
	}

	// Check read permissions
	if err := s.permissions.CanAccess(ctx, archivePath); err != nil {
		return err
	}

	// Validate destination path
	destPath, err := s.validator.ValidateAndSanitize(req.Destination)
	if err != nil {
		return err
	}

	// Check write permissions
	if err := s.permissions.CanWrite(ctx, destPath); err != nil {
		return err
	}

	// Detect archive format from extension
	ext := strings.ToLower(filepath.Ext(archivePath))

	switch ext {
	case ".zip":
		err = s.extractZipArchive(archivePath, destPath)
	case ".tar":
		err = s.extractTarArchive(archivePath, destPath, false)
	case ".gz", ".tgz":
		err = s.extractTarArchive(archivePath, destPath, true)
	default:
		return errors.BadRequest("Unsupported archive format (supported: .zip, .tar, .tar.gz, .tgz)", nil)
	}

	if err != nil {
		return err
	}

	logger.Info("Archive extracted", zap.String("archive", archivePath), zap.String("destination", destPath), zap.String("user", ctx.User.Username))
	return nil
}

// Helper: createZipArchive creates a ZIP archive
func (s *Service) createZipArchive(sourcePaths []string, outputPath string) error {
	// Create output file
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return errors.InternalServerError("Failed to create archive file", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add each source to the archive
	for _, sourcePath := range sourcePaths {
		if err := s.addToZip(zipWriter, sourcePath, ""); err != nil {
			os.Remove(outputPath) // Cleanup on error
			return err
		}
	}

	return nil
}

// Helper: addToZip adds a file or directory to a ZIP archive
func (s *Service) addToZip(zipWriter *zip.Writer, sourcePath, baseInZip string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return errors.InternalServerError("Failed to stat file", err)
	}

	if info.IsDir() {
		// Add directory recursively
		entries, err := os.ReadDir(sourcePath)
		if err != nil {
			return errors.InternalServerError("Failed to read directory", err)
		}

		for _, entry := range entries {
			entryPath := filepath.Join(sourcePath, entry.Name())
			entryBase := filepath.Join(baseInZip, filepath.Base(sourcePath), entry.Name())
			if err := s.addToZip(zipWriter, entryPath, entryBase); err != nil {
				return err
			}
		}
	} else {
		// Add file
		file, err := os.Open(sourcePath)
		if err != nil {
			return errors.InternalServerError("Failed to open file", err)
		}
		defer file.Close()

		// Create entry in ZIP
		zipPath := filepath.Join(baseInZip, filepath.Base(sourcePath))
		if baseInZip == "" {
			zipPath = filepath.Base(sourcePath)
		}

		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return errors.InternalServerError("Failed to create ZIP entry", err)
		}

		if _, err := io.Copy(writer, file); err != nil {
			return errors.InternalServerError("Failed to write to ZIP", err)
		}
	}

	return nil
}

// Helper: extractZipArchive extracts a ZIP archive
func (s *Service) extractZipArchive(archivePath, destPath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return errors.InternalServerError("Failed to open ZIP archive", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Validate path to prevent zip slip
		filePath := filepath.Join(destPath, file.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return errors.BadRequest("Archive contains invalid path (zip slip detected)", nil)
		}

		if file.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return errors.InternalServerError("Failed to create directory", err)
			}
		} else {
			// Create file
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return errors.InternalServerError("Failed to create parent directory", err)
			}

			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return errors.InternalServerError("Failed to create file", err)
			}

			rc, err := file.Open()
			if err != nil {
				outFile.Close()
				return errors.InternalServerError("Failed to open ZIP entry", err)
			}

			_, err = io.Copy(outFile, rc)
			rc.Close()
			outFile.Close()

			if err != nil {
				return errors.InternalServerError("Failed to extract file", err)
			}
		}
	}

	return nil
}

// Helper: createTarArchive creates a TAR archive (optionally gzipped)
func (s *Service) createTarArchive(sourcePaths []string, outputPath string, gzipped bool) error {
	// Create output file
	tarFile, err := os.Create(outputPath)
	if err != nil {
		return errors.InternalServerError("Failed to create archive file", err)
	}
	defer tarFile.Close()

	var tarWriter *tar.Writer

	if gzipped {
		gzipWriter := gzip.NewWriter(tarFile)
		defer gzipWriter.Close()
		tarWriter = tar.NewWriter(gzipWriter)
	} else {
		tarWriter = tar.NewWriter(tarFile)
	}
	defer tarWriter.Close()

	// Add each source to the archive
	for _, sourcePath := range sourcePaths {
		if err := s.addToTar(tarWriter, sourcePath, ""); err != nil {
			os.Remove(outputPath) // Cleanup on error
			return err
		}
	}

	return nil
}

// Helper: addToTar adds a file or directory to a TAR archive
func (s *Service) addToTar(tarWriter *tar.Writer, sourcePath, baseInTar string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return errors.InternalServerError("Failed to stat file", err)
	}

	// Create tar header
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return errors.InternalServerError("Failed to create tar header", err)
	}

	// Set name in archive
	if baseInTar != "" {
		header.Name = filepath.Join(baseInTar, filepath.Base(sourcePath))
	} else {
		header.Name = filepath.Base(sourcePath)
	}

	// Write header
	if err := tarWriter.WriteHeader(header); err != nil {
		return errors.InternalServerError("Failed to write tar header", err)
	}

	if info.IsDir() {
		// Add directory recursively
		entries, err := os.ReadDir(sourcePath)
		if err != nil {
			return errors.InternalServerError("Failed to read directory", err)
		}

		for _, entry := range entries {
			entryPath := filepath.Join(sourcePath, entry.Name())
			entryBase := header.Name
			if err := s.addToTar(tarWriter, entryPath, entryBase); err != nil {
				return err
			}
		}
	} else {
		// Add file content
		file, err := os.Open(sourcePath)
		if err != nil {
			return errors.InternalServerError("Failed to open file", err)
		}
		defer file.Close()

		if _, err := io.Copy(tarWriter, file); err != nil {
			return errors.InternalServerError("Failed to write to tar", err)
		}
	}

	return nil
}

// Helper: extractTarArchive extracts a TAR archive (optionally gzipped)
func (s *Service) extractTarArchive(archivePath, destPath string, gzipped bool) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return errors.InternalServerError("Failed to open archive", err)
	}
	defer file.Close()

	var tarReader *tar.Reader

	if gzipped {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return errors.InternalServerError("Failed to create gzip reader", err)
		}
		defer gzipReader.Close()
		tarReader = tar.NewReader(gzipReader)
	} else {
		tarReader = tar.NewReader(file)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.InternalServerError("Failed to read tar header", err)
		}

		// Validate path to prevent tar slip
		targetPath := filepath.Join(destPath, header.Name)
		if !strings.HasPrefix(targetPath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return errors.BadRequest("Archive contains invalid path (tar slip detected)", nil)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return errors.InternalServerError("Failed to create directory", err)
			}
		case tar.TypeReg:
			// Create file
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return errors.InternalServerError("Failed to create parent directory", err)
			}

			outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return errors.InternalServerError("Failed to create file", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return errors.InternalServerError("Failed to extract file", err)
			}
			outFile.Close()
		default:
			logger.Warn("Unsupported tar entry type", zap.String("name", header.Name), zap.Uint8("type", header.Typeflag))
		}
	}

	return nil
}
