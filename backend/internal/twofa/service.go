package twofa

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	// BackupCodeLength is the length of backup codes
	BackupCodeLength = 8
	// BackupCodeCount is the number of backup codes to generate
	BackupCodeCount = 10
	// MaxFailedAttempts is the maximum number of failed 2FA attempts before lockout
	MaxFailedAttempts = 5
	// AttemptWindow is the time window to check for failed attempts
	AttemptWindow = 15 * time.Minute
)

// Service manages two-factor authentication
type Service struct {
	db *gorm.DB
	mu sync.RWMutex
}

var (
	globalService *Service
	once          sync.Once
)

// SetupRequest contains information to set up 2FA for a user
type SetupRequest struct {
	UserID uint
	Issuer string
}

// SetupResponse contains the 2FA setup information
type SetupResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qrCodeUrl"`
	BackupCodes []string `json:"backupCodes"`
}

// VerifyRequest contains information to verify a 2FA code
type VerifyRequest struct {
	UserID uint
	Code   string
	IsBackupCode bool
}

// Initialize initializes the 2FA service
func Initialize() (*Service, error) {
	var initErr error
	once.Do(func() {
		db := database.GetDB()
		if db == nil {
			initErr = fmt.Errorf("database not initialized")
			return
		}

		globalService = &Service{
			db: db,
		}

		logger.Info("Two-Factor Authentication service initialized")
	})

	return globalService, initErr
}

// GetService returns the global 2FA service
func GetService() *Service {
	if globalService == nil {
		globalService, _ = Initialize()
	}
	return globalService
}

// IsEnabled checks if 2FA is enabled for a user
func (s *Service) IsEnabled(ctx context.Context, userID uint) (bool, error) {
	var twoFA models.TwoFactorAuth
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&twoFA).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return twoFA.Enabled, nil
}

// SetupTwoFactor initiates 2FA setup for a user
func (s *Service) SetupTwoFactor(ctx context.Context, req SetupRequest) (*SetupResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get user
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, req.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      req.Issuer,
		AccountName: user.Email,
		SecretSize:  32,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	secret := key.Secret()

	// Generate backup codes
	backupCodes, hashedCodes, err := s.generateBackupCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	// Store 2FA configuration (not enabled yet)
	twoFA := models.TwoFactorAuth{
		UserID:  req.UserID,
		Secret:  secret,
		Enabled: false, // Will be enabled after verification
	}

	// Delete existing 2FA config if present
	s.db.WithContext(ctx).Where("user_id = ?", req.UserID).Delete(&models.TwoFactorAuth{})
	s.db.WithContext(ctx).Where("user_id = ?", req.UserID).Delete(&models.TwoFactorBackupCode{})

	if err := s.db.WithContext(ctx).Create(&twoFA).Error; err != nil {
		return nil, fmt.Errorf("failed to store 2FA config: %w", err)
	}

	// Store backup codes
	for _, hashedCode := range hashedCodes {
		backupCode := models.TwoFactorBackupCode{
			UserID: req.UserID,
			Code:   hashedCode,
			Used:   false,
		}
		if err := s.db.WithContext(ctx).Create(&backupCode).Error; err != nil {
			logger.Error("Failed to store backup code", zap.Error(err))
		}
	}

	logger.Info("2FA setup initiated", zap.Uint("userId", req.UserID))

	return &SetupResponse{
		Secret:      secret,
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

// EnableTwoFactor enables 2FA for a user after verifying the initial code
func (s *Service) EnableTwoFactor(ctx context.Context, userID uint, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get 2FA config
	var twoFA models.TwoFactorAuth
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&twoFA).Error; err != nil {
		return fmt.Errorf("2FA not set up: %w", err)
	}

	// Verify the code
	valid := totp.Validate(code, twoFA.Secret)
	if !valid {
		return fmt.Errorf("invalid verification code")
	}

	// Enable 2FA
	twoFA.Enabled = true
	if err := s.db.WithContext(ctx).Save(&twoFA).Error; err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	logger.Info("2FA enabled", zap.Uint("userId", userID))
	return nil
}

// DisableTwoFactor disables 2FA for a user
func (s *Service) DisableTwoFactor(ctx context.Context, userID uint, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get 2FA config
	var twoFA models.TwoFactorAuth
	if err := s.db.WithContext(ctx).Where("user_id = ? AND enabled = ?", userID, true).First(&twoFA).Error; err != nil {
		return fmt.Errorf("2FA not enabled: %w", err)
	}

	// Verify the code
	if !s.verifyCode(twoFA.Secret, code) {
		return fmt.Errorf("invalid verification code")
	}

	// Delete 2FA config and backup codes
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.TwoFactorAuth{}).Error; err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.TwoFactorBackupCode{}).Error; err != nil {
		logger.Error("Failed to delete backup codes", zap.Error(err))
	}

	logger.Info("2FA disabled", zap.Uint("userId", userID))
	return nil
}

// VerifyCode verifies a TOTP code or backup code
func (s *Service) VerifyCode(ctx context.Context, req VerifyRequest) (bool, error) {
	// Check rate limiting
	if err := s.checkRateLimit(ctx, req.UserID); err != nil {
		return false, err
	}

	// Get 2FA config
	var twoFA models.TwoFactorAuth
	if err := s.db.WithContext(ctx).Where("user_id = ? AND enabled = ?", req.UserID, true).First(&twoFA).Error; err != nil {
		return false, fmt.Errorf("2FA not enabled: %w", err)
	}

	var valid bool

	if req.IsBackupCode {
		// Verify backup code
		valid = s.verifyBackupCode(ctx, req.UserID, req.Code)
	} else {
		// Verify TOTP code
		valid = s.verifyCode(twoFA.Secret, req.Code)
	}

	// Record attempt
	attempt := models.TwoFactorAttempt{
		UserID:      req.UserID,
		IPAddress:   "", // Will be set by handler
		Success:     valid,
		AttemptedAt: time.Now(),
	}
	s.db.WithContext(ctx).Create(&attempt)

	if !valid {
		logger.Warn("Failed 2FA attempt", zap.Uint("userId", req.UserID))
	}

	return valid, nil
}

// verifyCode verifies a TOTP code
func (s *Service) verifyCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// verifyBackupCode verifies and consumes a backup code
func (s *Service) verifyBackupCode(ctx context.Context, userID uint, code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get all unused backup codes for user
	var backupCodes []models.TwoFactorBackupCode
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND used = ?", userID, false).
		Find(&backupCodes).Error; err != nil {
		logger.Error("Failed to get backup codes", zap.Error(err))
		return false
	}

	// Check each backup code
	for _, backupCode := range backupCodes {
		err := bcrypt.CompareHashAndPassword([]byte(backupCode.Code), []byte(code))
		if err == nil {
			// Mark code as used
			now := time.Now()
			backupCode.Used = true
			backupCode.UsedAt = &now
			if err := s.db.WithContext(ctx).Save(&backupCode).Error; err != nil {
				logger.Error("Failed to mark backup code as used", zap.Error(err))
			}
			logger.Info("Backup code used", zap.Uint("userId", userID))
			return true
		}
	}

	return false
}

// checkRateLimit checks if user has exceeded failed attempt limit
func (s *Service) checkRateLimit(ctx context.Context, userID uint) error {
	cutoff := time.Now().Add(-AttemptWindow)

	var failedCount int64
	if err := s.db.WithContext(ctx).
		Model(&models.TwoFactorAttempt{}).
		Where("user_id = ? AND success = ? AND attempted_at > ?", userID, false, cutoff).
		Count(&failedCount).Error; err != nil {
		return err
	}

	if failedCount >= MaxFailedAttempts {
		return fmt.Errorf("too many failed attempts, please try again later")
	}

	return nil
}

// generateBackupCodes generates a set of backup codes
func (s *Service) generateBackupCodes() ([]string, []string, error) {
	codes := make([]string, BackupCodeCount)
	hashedCodes := make([]string, BackupCodeCount)

	for i := 0; i < BackupCodeCount; i++ {
		code, err := generateRandomCode(BackupCodeLength)
		if err != nil {
			return nil, nil, err
		}

		// Hash the code for storage
		hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, err
		}

		codes[i] = code
		hashedCodes[i] = string(hashedCode)
	}

	return codes, hashedCodes, nil
}

// generateRandomCode generates a random alphanumeric code
func generateRandomCode(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}

	return string(result), nil
}

// GenerateSecret generates a new TOTP secret
func (s *Service) GenerateSecret() (string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), nil
}

// GenerateQRCode generates a QR code URL for TOTP setup
func (s *Service) GenerateQRCode(secret, issuer, accountName string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, accountName, secret, issuer)
}

// GetBackupCodes retrieves remaining backup codes for a user (count only, not actual codes)
func (s *Service) GetBackupCodes(ctx context.Context, userID uint) (int, error) {
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&models.TwoFactorBackupCode{}).
		Where("user_id = ? AND used = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// RegenerateBackupCodes generates new backup codes for a user
func (s *Service) RegenerateBackupCodes(ctx context.Context, userID uint, verificationCode string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify 2FA is enabled
	var twoFA models.TwoFactorAuth
	if err := s.db.WithContext(ctx).Where("user_id = ? AND enabled = ?", userID, true).First(&twoFA).Error; err != nil {
		return nil, fmt.Errorf("2FA not enabled: %w", err)
	}

	// Verify the code
	if !s.verifyCode(twoFA.Secret, verificationCode) {
		return nil, fmt.Errorf("invalid verification code")
	}

	// Delete old backup codes
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.TwoFactorBackupCode{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete old backup codes: %w", err)
	}

	// Generate new backup codes
	backupCodes, hashedCodes, err := s.generateBackupCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	// Store new backup codes
	for _, hashedCode := range hashedCodes {
		backupCode := models.TwoFactorBackupCode{
			UserID: userID,
			Code:   hashedCode,
			Used:   false,
		}
		if err := s.db.WithContext(ctx).Create(&backupCode).Error; err != nil {
			logger.Error("Failed to store backup code", zap.Error(err))
		}
	}

	logger.Info("Backup codes regenerated", zap.Uint("userId", userID))
	return backupCodes, nil
}

// CleanupOldAttempts removes old 2FA attempt records
func (s *Service) CleanupOldAttempts(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)

	result := s.db.WithContext(ctx).
		Where("attempted_at < ?", cutoff).
		Delete(&models.TwoFactorAttempt{})

	if result.Error != nil {
		return result.Error
	}

	logger.Info("Cleaned up old 2FA attempts", zap.Int64("deleted", result.RowsAffected))
	return nil
}

// FormatBackupCode formats a backup code for display (e.g., XXXX-XXXX)
func FormatBackupCode(code string) string {
	if len(code) != BackupCodeLength {
		return code
	}
	return fmt.Sprintf("%s-%s", code[:4], code[4:])
}

// UnformatBackupCode removes formatting from a backup code
func UnformatBackupCode(code string) string {
	return strings.ReplaceAll(code, "-", "")
}
