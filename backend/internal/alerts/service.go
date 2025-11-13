package alerts

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service handles alerting functionality
type Service struct {
	db              *gorm.DB
	mu              sync.RWMutex
	lastAlertTimes  map[string]time.Time // Rate limiting by alert type
}

var (
	globalService *Service
	once          sync.Once
)

// Initialize initializes the alert service
func Initialize() (*Service, error) {
	var initErr error
	once.Do(func() {
		db := database.GetDB()
		if db == nil {
			initErr = fmt.Errorf("database not initialized")
			return
		}

		globalService = &Service{
			db:             db,
			lastAlertTimes: make(map[string]time.Time),
		}

		logger.Info("Alert service initialized")
	})

	return globalService, initErr
}

// GetService returns the global alert service
func GetService() *Service {
	if globalService == nil {
		globalService, _ = Initialize()
	}
	return globalService
}

// GetConfig retrieves the alert configuration
func (s *Service) GetConfig(ctx context.Context) (*models.AlertConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var config models.AlertConfig
	result := s.db.WithContext(ctx).First(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return default config
			return &models.AlertConfig{
				Enabled:              false,
				SMTPPort:             587,
				SMTPUseTLS:           true,
				OnFailedLogin:        true,
				OnIPBlock:            true,
				OnCriticalEvent:      true,
				FailedLoginThreshold: 3,
				RateLimitMinutes:     15,
			}, nil
		}
		return nil, result.Error
	}

	return &config, nil
}

// UpdateConfig updates the alert configuration
func (s *Service) UpdateConfig(ctx context.Context, config *models.AlertConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var existingConfig models.AlertConfig
	result := s.db.WithContext(ctx).First(&existingConfig)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new config
		return s.db.WithContext(ctx).Create(config).Error
	}

	// Update existing config
	config.ID = existingConfig.ID
	config.CreatedAt = existingConfig.CreatedAt
	return s.db.WithContext(ctx).Save(config).Error
}

// TestEmail sends a test email
func (s *Service) TestEmail(ctx context.Context, config *models.AlertConfig) error {
	if config.AlertRecipient == "" {
		return fmt.Errorf("alert recipient email is required")
	}

	subject := "Stumpf.Works NAS - Test Alert"
	body := fmt.Sprintf(`
<html>
<body>
<h2>Test Alert</h2>
<p>This is a test alert from your Stumpf.Works NAS system.</p>
<p>If you received this email, your alert configuration is working correctly.</p>
<p><strong>Time:</strong> %s</p>
</body>
</html>
`, time.Now().Format("2006-01-02 15:04:05"))

	return s.sendEmail(ctx, config, subject, body, models.AlertTypeSystemError)
}

// SendFailedLoginAlert sends an alert for failed login attempts
func (s *Service) SendFailedLoginAlert(ctx context.Context, username, ipAddress string, attemptCount int) error {
	config, err := s.GetConfig(ctx)
	if err != nil || !config.Enabled || !config.OnFailedLogin {
		return nil // Silently skip if not enabled
	}

	// Check threshold
	if attemptCount < config.FailedLoginThreshold {
		return nil
	}

	// Check rate limiting
	if !s.shouldSendAlert(models.AlertTypeFailedLogin, config.RateLimitMinutes) {
		logger.Debug("Skipping alert due to rate limiting",
			zap.String("type", models.AlertTypeFailedLogin))
		return nil
	}

	subject := fmt.Sprintf("⚠️ Failed Login Alert - %d Attempts Detected", attemptCount)
	body := fmt.Sprintf(`
<html>
<body>
<h2>Failed Login Alert</h2>
<p><strong>Multiple failed login attempts have been detected on your system.</strong></p>
<ul>
<li><strong>Username:</strong> %s</li>
<li><strong>IP Address:</strong> %s</li>
<li><strong>Attempt Count:</strong> %d</li>
<li><strong>Time:</strong> %s</li>
</ul>
<p>If this was not you, please review your security settings immediately.</p>
</body>
</html>
`, username, ipAddress, attemptCount, time.Now().Format("2006-01-02 15:04:05"))

	return s.sendEmail(ctx, config, subject, body, models.AlertTypeFailedLogin)
}

// SendIPBlockAlert sends an alert when an IP is blocked
func (s *Service) SendIPBlockAlert(ctx context.Context, ipAddress string, reason string, attempts int) error {
	config, err := s.GetConfig(ctx)
	if err != nil || !config.Enabled || !config.OnIPBlock {
		return nil
	}

	// Check rate limiting
	if !s.shouldSendAlert(models.AlertTypeIPBlock, config.RateLimitMinutes) {
		logger.Debug("Skipping alert due to rate limiting",
			zap.String("type", models.AlertTypeIPBlock))
		return nil
	}

	subject := fmt.Sprintf("🛡️ IP Blocked - Security Alert")
	body := fmt.Sprintf(`
<html>
<body>
<h2>IP Block Alert</h2>
<p><strong>An IP address has been automatically blocked due to suspicious activity.</strong></p>
<ul>
<li><strong>IP Address:</strong> %s</li>
<li><strong>Reason:</strong> %s</li>
<li><strong>Failed Attempts:</strong> %d</li>
<li><strong>Time:</strong> %s</li>
</ul>
<p>The IP address will remain blocked for 15 minutes. You can manually unblock it from the Security dashboard.</p>
</body>
</html>
`, ipAddress, reason, attempts, time.Now().Format("2006-01-02 15:04:05"))

	return s.sendEmail(ctx, config, subject, body, models.AlertTypeIPBlock)
}

// SendCriticalEventAlert sends an alert for critical security events
func (s *Service) SendCriticalEventAlert(ctx context.Context, action, username, ipAddress, message string) error {
	config, err := s.GetConfig(ctx)
	if err != nil || !config.Enabled || !config.OnCriticalEvent {
		return nil
	}

	// Check rate limiting
	if !s.shouldSendAlert(models.AlertTypeCriticalEvent, config.RateLimitMinutes) {
		logger.Debug("Skipping alert due to rate limiting",
			zap.String("type", models.AlertTypeCriticalEvent))
		return nil
	}

	subject := fmt.Sprintf("🚨 Critical Security Event - %s", action)
	body := fmt.Sprintf(`
<html>
<body>
<h2>Critical Security Event</h2>
<p><strong>A critical security event has been detected on your system.</strong></p>
<ul>
<li><strong>Action:</strong> %s</li>
<li><strong>User:</strong> %s</li>
<li><strong>IP Address:</strong> %s</li>
<li><strong>Message:</strong> %s</li>
<li><strong>Time:</strong> %s</li>
</ul>
<p>Please review the audit logs for more details.</p>
</body>
</html>
`, action, username, ipAddress, message, time.Now().Format("2006-01-02 15:04:05"))

	return s.sendEmail(ctx, config, subject, body, models.AlertTypeCriticalEvent)
}

// shouldSendAlert checks if an alert should be sent based on rate limiting
func (s *Service) shouldSendAlert(alertType string, rateLimitMinutes int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	lastTime, exists := s.lastAlertTimes[alertType]
	if !exists {
		s.lastAlertTimes[alertType] = time.Now()
		return true
	}

	if time.Since(lastTime) > time.Duration(rateLimitMinutes)*time.Minute {
		s.lastAlertTimes[alertType] = time.Now()
		return true
	}

	return false
}

// sendEmail sends an email alert
func (s *Service) sendEmail(ctx context.Context, config *models.AlertConfig, subject, body, alertType string) error {
	// Validate config
	if config.SMTPHost == "" || config.AlertRecipient == "" {
		return fmt.Errorf("SMTP host and recipient are required")
	}

	from := config.SMTPFromEmail
	if from == "" {
		from = config.SMTPUsername
	}

	fromName := config.SMTPFromName
	if fromName == "" {
		fromName = "Stumpf.Works NAS"
	}

	// Prepare email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
	headers["To"] = config.AlertRecipient
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Send email
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	addr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)

	var err error
	if config.SMTPUseTLS {
		err = s.sendEmailTLS(addr, auth, from, []string{config.AlertRecipient}, []byte(message))
	} else {
		err = smtp.SendMail(addr, auth, from, []string{config.AlertRecipient}, []byte(message))
	}

	// Log the alert
	alertLog := &models.AlertLog{
		AlertType: alertType,
		Subject:   subject,
		Body:      body,
		Recipient: config.AlertRecipient,
		Status:    "sent",
	}

	if err != nil {
		alertLog.Status = "failed"
		alertLog.Error = err.Error()
		logger.Error("Failed to send alert email",
			zap.Error(err),
			zap.String("type", alertType),
			zap.String("recipient", config.AlertRecipient))
	} else {
		logger.Info("Alert email sent",
			zap.String("type", alertType),
			zap.String("recipient", config.AlertRecipient))
	}

	s.db.WithContext(ctx).Create(alertLog)

	return err
}

// sendEmailTLS sends email with TLS
func (s *Service) sendEmailTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName:         addr[:len(addr)-4], // Remove :port
		InsecureSkipVerify: false,
	}

	// Connect to SMTP server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, tlsConfig.ServerName)
	if err != nil {
		return err
	}
	defer client.Close()

	// Authenticate
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	// Set sender and recipients
	if err = client.Mail(from); err != nil {
		return err
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return err
		}
	}

	// Send message
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

// GetAlertLogs retrieves recent alert logs
func (s *Service) GetAlertLogs(ctx context.Context, limit int) ([]models.AlertLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var logs []models.AlertLog
	result := s.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs)

	return logs, result.Error
}
