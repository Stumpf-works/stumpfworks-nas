// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/metrics"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// MetricsHandler handles metrics-related HTTP requests
type MetricsHandler struct {
	service *metrics.Service
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		service: metrics.GetService(),
	}
}

// GetMetricsHistory returns historical metrics
func (h *MetricsHandler) GetMetricsHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	limitStr := r.URL.Query().Get("limit")

	// Default to last 24 hours
	end := time.Now()
	start := end.Add(-24 * time.Hour)

	if startStr != "" {
		if ts, err := time.Parse(time.RFC3339, startStr); err == nil {
			start = ts
		}
	}

	if endStr != "" {
		if ts, err := time.Parse(time.RFC3339, endStr); err == nil {
			end = ts
		}
	}

	limit := 1000 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get metrics
	metricsData, err := h.service.GetMetrics(ctx, start, end, limit)
	if err != nil {
		logger.Error("Failed to get metrics", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve metrics", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"metrics": metricsData,
		"start":   start,
		"end":     end,
		"count":   len(metricsData),
	})
}

// GetLatestMetric returns the most recent metric
func (h *MetricsHandler) GetLatestMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metric, err := h.service.GetLatestMetric(ctx)
	if err != nil {
		logger.Error("Failed to get latest metric", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve latest metric", err))
		return
	}

	utils.RespondSuccess(w, metric)
}

// GetHealthScores returns historical health scores
func (h *MetricsHandler) GetHealthScores(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	limitStr := r.URL.Query().Get("limit")

	// Default to last 7 days
	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)

	if startStr != "" {
		if ts, err := time.Parse(time.RFC3339, startStr); err == nil {
			start = ts
		}
	}

	if endStr != "" {
		if ts, err := time.Parse(time.RFC3339, endStr); err == nil {
			end = ts
		}
	}

	limit := 1000 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get health scores
	scores, err := h.service.GetHealthScores(ctx, start, end, limit)
	if err != nil {
		logger.Error("Failed to get health scores", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve health scores", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"scores": scores,
		"start":  start,
		"end":    end,
		"count":  len(scores),
	})
}

// GetLatestHealthScore returns the most recent health score
func (h *MetricsHandler) GetLatestHealthScore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	score, err := h.service.GetLatestHealthScore(ctx)
	if err != nil {
		logger.Error("Failed to get latest health score", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve latest health score", err))
		return
	}

	utils.RespondSuccess(w, score)
}

// GetTrends returns trend analysis for key metrics
func (h *MetricsHandler) GetTrends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse duration parameter (default to 1 hour)
	durationStr := r.URL.Query().Get("duration")
	duration := 1 * time.Hour

	if durationStr != "" {
		if d, err := time.ParseDuration(durationStr); err == nil {
			duration = d
		}
	}

	// Get trends
	trends, err := h.service.GetTrends(ctx, duration)
	if err != nil {
		logger.Error("Failed to get trends", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve trends", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"trends":   trends,
		"duration": duration.String(),
	})
}

// PrometheusMetricsHandler handles GET /metrics for Prometheus scraping
// This endpoint exposes system metrics in Prometheus text format
func PrometheusMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Get system library instance
	lib := system.Get()
	if lib == nil || lib.Metrics == nil {
		logger.Warn("System metrics collector not available")
		http.Error(w, "Metrics collector not initialized", http.StatusServiceUnavailable)
		return
	}

	// Get current metrics
	current := lib.Metrics.GetCurrent()
	if current == nil {
		logger.Warn("No system metrics available")
		http.Error(w, "No metrics available", http.StatusServiceUnavailable)
		return
	}

	// Convert to Prometheus format
	prometheusOutput := current.ToPrometheusFormat()

	// Set content type for Prometheus
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Write metrics
	_, err := w.Write([]byte(prometheusOutput))
	if err != nil {
		logger.Error("Failed to write Prometheus metrics response", zap.Error(err))
	}

	logger.Debug("Prometheus metrics exported",
		zap.String("remote_addr", r.RemoteAddr),
		zap.Int("response_size", len(prometheusOutput)))
}

// SystemHealthHandler handles GET /api/v1/system/health
// Returns overall system health status from centralized system library
func SystemHealthHandler(w http.ResponseWriter, r *http.Request) {
	lib := system.Get()
	if lib == nil {
		utils.RespondError(w, errors.InternalServerError("System library not initialized", nil))
		return
	}

	health, err := lib.HealthCheck()
	if err != nil {
		logger.Error("Failed to get system health", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get system health", err))
		return
	}

	utils.RespondSuccess(w, health)
}
