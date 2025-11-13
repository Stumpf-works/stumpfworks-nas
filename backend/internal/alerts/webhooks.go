package alerts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// DiscordWebhook represents a Discord webhook message
type DiscordWebhook struct {
	Username  string         `json:"username,omitempty"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Content   string         `json:"content,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

// DiscordEmbed represents a Discord embed
type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Color       int                 `json:"color,omitempty"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"`
}

// DiscordEmbedField represents a Discord embed field
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// SlackWebhook represents a Slack webhook message
type SlackWebhook struct {
	Username    string             `json:"username,omitempty"`
	IconURL     string             `json:"icon_url,omitempty"`
	Text        string             `json:"text,omitempty"`
	Attachments []SlackAttachment  `json:"attachments,omitempty"`
}

// SlackAttachment represents a Slack attachment
type SlackAttachment struct {
	Color      string        `json:"color,omitempty"`
	Title      string        `json:"title,omitempty"`
	Text       string        `json:"text,omitempty"`
	Fields     []SlackField  `json:"fields,omitempty"`
	Footer     string        `json:"footer,omitempty"`
	Timestamp  int64         `json:"ts,omitempty"`
}

// SlackField represents a Slack field
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// sendWebhook sends a webhook notification
func (s *Service) sendWebhook(ctx context.Context, config *models.AlertConfig, subject, body, alertType string) error {
	if !config.WebhookEnabled || config.WebhookURL == "" {
		return fmt.Errorf("webhook not configured")
	}

	var payload []byte
	var err error

	switch config.WebhookType {
	case models.WebhookTypeDiscord:
		payload, err = s.buildDiscordWebhook(config, subject, body, alertType)
	case models.WebhookTypeSlack:
		payload, err = s.buildSlackWebhook(config, subject, body, alertType)
	case models.WebhookTypeCustom:
		payload, err = s.buildCustomWebhook(config, subject, body, alertType)
	default:
		return fmt.Errorf("unsupported webhook type: %s", config.WebhookType)
	}

	if err != nil {
		return fmt.Errorf("failed to build webhook payload: %w", err)
	}

	// Send HTTP POST request
	req, err := http.NewRequestWithContext(ctx, "POST", config.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	// Log the webhook
	alertLog := &models.AlertLog{
		AlertType: alertType,
		Channel:   models.AlertChannelWebhook,
		Subject:   subject,
		Body:      body,
		Recipient: config.WebhookURL,
		Status:    "sent",
	}

	s.db.WithContext(ctx).Create(alertLog)

	logger.Info("Webhook sent successfully",
		zap.String("type", alertType),
		zap.String("webhook_type", config.WebhookType))

	return nil
}

// buildDiscordWebhook builds a Discord webhook payload
func (s *Service) buildDiscordWebhook(config *models.AlertConfig, subject, body, alertType string) ([]byte, error) {
	// Determine color based on alert type
	color := 0x5865F2 // Discord blue
	switch alertType {
	case models.AlertTypeFailedLogin:
		color = 0xF0B232 // Warning orange
	case models.AlertTypeIPBlock:
		color = 0xED4245 // Error red
	case models.AlertTypeCriticalEvent:
		color = 0xED4245 // Error red
	}

	username := config.WebhookUsername
	if username == "" {
		username = "Stumpf.Works NAS"
	}

	webhook := DiscordWebhook{
		Username:  username,
		AvatarURL: config.WebhookAvatarURL,
		Embeds: []DiscordEmbed{
			{
				Title:       subject,
				Description: body,
				Color:       color,
				Timestamp:   time.Now().Format(time.RFC3339),
			},
		},
	}

	return json.Marshal(webhook)
}

// buildSlackWebhook builds a Slack webhook payload
func (s *Service) buildSlackWebhook(config *models.AlertConfig, subject, body, alertType string) ([]byte, error) {
	// Determine color based on alert type
	color := "good"
	switch alertType {
	case models.AlertTypeFailedLogin:
		color = "warning"
	case models.AlertTypeIPBlock:
		color = "danger"
	case models.AlertTypeCriticalEvent:
		color = "danger"
	}

	username := config.WebhookUsername
	if username == "" {
		username = "Stumpf.Works NAS"
	}

	webhook := SlackWebhook{
		Username: username,
		IconURL:  config.WebhookAvatarURL,
		Attachments: []SlackAttachment{
			{
				Color:     color,
				Title:     subject,
				Text:      body,
				Footer:    "Stumpf.Works NAS Alert System",
				Timestamp: time.Now().Unix(),
			},
		},
	}

	return json.Marshal(webhook)
}

// buildCustomWebhook builds a generic webhook payload
func (s *Service) buildCustomWebhook(config *models.AlertConfig, subject, body, alertType string) ([]byte, error) {
	payload := map[string]interface{}{
		"alert_type": alertType,
		"subject":    subject,
		"body":       body,
		"timestamp":  time.Now().Unix(),
		"source":     "stumpfworks-nas",
	}

	if config.WebhookUsername != "" {
		payload["username"] = config.WebhookUsername
	}

	return json.Marshal(payload)
}

// TestWebhook sends a test webhook
func (s *Service) TestWebhook(ctx context.Context, config *models.AlertConfig) error {
	if config.WebhookURL == "" {
		return fmt.Errorf("webhook URL is required")
	}

	subject := "Stumpf.Works NAS - Test Webhook"
	body := fmt.Sprintf("This is a test webhook from your Stumpf.Works NAS system.\n\nTime: %s\nWebhook Type: %s",
		time.Now().Format("2006-01-02 15:04:05"),
		config.WebhookType)

	return s.sendWebhook(ctx, config, subject, body, models.AlertTypeSystemError)
}
