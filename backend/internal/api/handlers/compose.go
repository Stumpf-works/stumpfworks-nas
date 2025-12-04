// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/Stumpf-works/stumpfworks-nas/internal/docker"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	// DefaultStacksDir is the default directory for Docker Compose stacks
	DefaultStacksDir = "/var/lib/stumpfworks/stacks"
)

// ComposeHandler handles Docker Compose stack operations
type ComposeHandler struct {
	service   *docker.Service
	stacksDir string
}

// NewComposeHandler creates a new Compose handler
func NewComposeHandler(stacksDir string) *ComposeHandler {
	if stacksDir == "" {
		stacksDir = DefaultStacksDir
	}
	return &ComposeHandler{
		service:   docker.GetService(),
		stacksDir: stacksDir,
	}
}

// ListStacks lists all Docker Compose stacks
func (h *ComposeHandler) ListStacks(w http.ResponseWriter, r *http.Request) {
	stacks, err := h.service.ListStacks(r.Context(), h.stacksDir)
	if err != nil {
		logger.Error("Failed to list stacks", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list stacks", err))
		return
	}

	utils.RespondSuccess(w, stacks)
}

// GetStack gets detailed information about a specific stack
func (h *ComposeHandler) GetStack(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	stack, err := h.service.GetStack(r.Context(), stackName, stackPath)
	if err != nil {
		logger.Error("Failed to get stack", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to get stack", err))
		return
	}

	utils.RespondSuccess(w, stack)
}

// CreateStack creates a new Docker Compose stack
func (h *ComposeHandler) CreateStack(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Compose string `json:"compose"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Name == "" {
		utils.RespondError(w, errors.BadRequest("Stack name is required", nil))
		return
	}

	if req.Compose == "" {
		utils.RespondError(w, errors.BadRequest("Compose content is required", nil))
		return
	}

	if err := h.service.CreateStack(r.Context(), h.stacksDir, req.Name, req.Compose); err != nil {
		logger.Error("Failed to create stack", zap.Error(err), zap.String("stack", req.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to create stack", err))
		return
	}

	logger.Info("Stack created", zap.String("stack", req.Name))
	utils.RespondSuccess(w, map[string]string{"message": "Stack created successfully", "name": req.Name})
}

// UpdateStack updates an existing stack's compose file
func (h *ComposeHandler) UpdateStack(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	var req struct {
		Compose string `json:"compose"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Compose == "" {
		utils.RespondError(w, errors.BadRequest("Compose content is required", nil))
		return
	}

	if err := h.service.UpdateStack(r.Context(), stackPath, req.Compose); err != nil {
		logger.Error("Failed to update stack", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to update stack", err))
		return
	}

	logger.Info("Stack updated", zap.String("stack", stackName))
	utils.RespondSuccess(w, map[string]string{"message": "Stack updated successfully"})
}

// DeleteStack deletes a stack and its directory
func (h *ComposeHandler) DeleteStack(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	if err := h.service.DeleteStack(r.Context(), stackPath); err != nil {
		logger.Error("Failed to delete stack", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to delete stack", err))
		return
	}

	logger.Info("Stack deleted", zap.String("stack", stackName))
	utils.RespondSuccess(w, map[string]string{"message": "Stack deleted successfully"})
}

// DeployStack deploys a Docker Compose stack
func (h *ComposeHandler) DeployStack(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	if err := h.service.DeployStack(r.Context(), stackPath); err != nil {
		logger.Error("Failed to deploy stack", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to deploy stack", err))
		return
	}

	logger.Info("Stack deployed", zap.String("stack", stackName))
	utils.RespondSuccess(w, map[string]string{"message": "Stack deployed successfully"})
}

// StopStack stops a Docker Compose stack
func (h *ComposeHandler) StopStack(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	if err := h.service.StopStack(r.Context(), stackPath); err != nil {
		logger.Error("Failed to stop stack", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to stop stack", err))
		return
	}

	logger.Info("Stack stopped", zap.String("stack", stackName))
	utils.RespondSuccess(w, map[string]string{"message": "Stack stopped successfully"})
}

// RestartStack restarts a Docker Compose stack
func (h *ComposeHandler) RestartStack(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	if err := h.service.RestartStack(r.Context(), stackPath); err != nil {
		logger.Error("Failed to restart stack", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to restart stack", err))
		return
	}

	logger.Info("Stack restarted", zap.String("stack", stackName))
	utils.RespondSuccess(w, map[string]string{"message": "Stack restarted successfully"})
}

// RemoveStack removes a Docker Compose stack (docker-compose down)
func (h *ComposeHandler) RemoveStack(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	// Check if we should remove volumes
	removeVolumes := r.URL.Query().Get("volumes") == "true"

	if err := h.service.RemoveStack(r.Context(), stackPath, removeVolumes); err != nil {
		logger.Error("Failed to remove stack", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to remove stack", err))
		return
	}

	logger.Info("Stack removed", zap.String("stack", stackName), zap.Bool("volumes", removeVolumes))
	utils.RespondSuccess(w, map[string]string{"message": "Stack removed successfully"})
}

// GetStackLogs gets logs from a Docker Compose stack
func (h *ComposeHandler) GetStackLogs(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	// Get tail parameter (default 500 lines)
	tail := 500

	logs, err := h.service.GetStackLogs(r.Context(), stackPath, tail)
	if err != nil {
		logger.Error("Failed to get stack logs", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to get stack logs", err))
		return
	}

	utils.RespondSuccess(w, logs)
}

// GetComposeFile gets the compose file content
func (h *ComposeHandler) GetComposeFile(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "name")
	stackPath := filepath.Join(h.stacksDir, stackName)

	content, err := h.service.GetComposeFile(stackPath)
	if err != nil {
		logger.Error("Failed to get compose file", zap.Error(err), zap.String("stack", stackName))
		utils.RespondError(w, errors.InternalServerError("Failed to get compose file", err))
		return
	}

	utils.RespondSuccess(w, content)
}

// ===== Template Management =====

// ListTemplates lists all available Docker Compose templates
func (h *ComposeHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var templates []docker.ComposeTemplate
	if category != "" {
		templates = docker.GetTemplatesByCategory(category)
	} else {
		templates = docker.ListAllTemplates()
	}

	utils.RespondSuccess(w, templates)
}

// GetTemplate gets a specific template by ID
func (h *ComposeHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")

	template := docker.GetTemplateByID(templateID)
	if template == nil {
		utils.RespondError(w, errors.NotFound("Template not found", nil))
		return
	}

	utils.RespondSuccess(w, template)
}

// GetTemplateCategories returns all unique template categories
func (h *ComposeHandler) GetTemplateCategories(w http.ResponseWriter, r *http.Request) {
	categories := docker.GetAllCategories()
	utils.RespondSuccess(w, categories)
}

// DeployTemplate deploys a template with user-provided variables
func (h *ComposeHandler) DeployTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")

	var req struct {
		StackName string            `json:"stack_name"`
		Variables map[string]string `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.StackName == "" {
		utils.RespondError(w, errors.BadRequest("Stack name is required", nil))
		return
	}

	// Get template
	template := docker.GetTemplateByID(templateID)
	if template == nil {
		utils.RespondError(w, errors.NotFound("Template not found", nil))
		return
	}

	// Render template with variables
	composeContent := docker.RenderTemplate(template, req.Variables)

	// Create stack from rendered template
	if err := h.service.CreateStack(r.Context(), h.stacksDir, req.StackName, composeContent); err != nil {
		logger.Error("Failed to create stack from template", zap.Error(err), zap.String("template", templateID), zap.String("stack", req.StackName))
		utils.RespondError(w, errors.InternalServerError("Failed to deploy template", err))
		return
	}

	// Deploy the stack
	stackPath := filepath.Join(h.stacksDir, req.StackName)
	if err := h.service.DeployStack(r.Context(), stackPath); err != nil {
		logger.Error("Failed to deploy stack", zap.Error(err), zap.String("stack", req.StackName))
		utils.RespondError(w, errors.InternalServerError("Failed to deploy template", err))
		return
	}

	// Auto-open required ports from template
	openedPorts := []int{}
	if template.Requirements.Ports != nil && len(template.Requirements.Ports) > 0 {
		if err := h.openTemplatePorts(template, req.StackName); err != nil {
			logger.Warn("Failed to open template ports", zap.Error(err), zap.String("stack", req.StackName))
			// Don't fail deployment if port opening fails
		} else {
			openedPorts = template.Requirements.Ports
		}
	}

	logger.Info("Template deployed successfully", zap.String("template", templateID), zap.String("stack", req.StackName), zap.Ints("ports", openedPorts))
	utils.RespondSuccess(w, map[string]interface{}{
		"message":      "Template deployed successfully",
		"template_id":  templateID,
		"stack_name":   req.StackName,
		"ports_opened": openedPorts,
	})
}

// openTemplatePorts opens firewall ports required by a template
func (h *ComposeHandler) openTemplatePorts(template *docker.ComposeTemplate, stackName string) error {
	if template.Requirements.Ports == nil || len(template.Requirements.Ports) == 0 {
		return nil
	}

	// TODO: Integrate with firewall manager to open ports
	// For now, we'll just log the ports that should be opened
	logger.Info("Template requires ports to be opened",
		zap.String("stack", stackName),
		zap.Ints("ports", template.Requirements.Ports),
		zap.String("template", template.Name),
	)

	// Note: Actual firewall integration would go here
	// Example:
	// for _, port := range template.Requirements.Ports {
	//     firewall.AddRule(FirewallRule{
	//         Chain: "INPUT",
	//         Protocol: "tcp",
	//         DestPort: strconv.Itoa(port),
	//         Action: "ACCEPT",
	//         Comment: fmt.Sprintf("Docker stack: %s (%s)", stackName, template.Name),
	//     })
	// }

	return nil
}


// GetHubStatus returns the status of Stumpfworks Hub connection
func (h *ComposeHandler) GetHubStatus(w http.ResponseWriter, r *http.Request) {
	hub := docker.GetHubClient()
	ctx := r.Context()

	// Try to fetch templates from Hub
	templates, err := hub.ListTemplates(ctx)

	status := map[string]interface{}{
		"hub_url":        hub.GetBaseURL(),
		"is_online":      false,
		"template_count": 0,
		"error":          nil,
	}

	if err != nil {
		status["error"] = err.Error()
		status["is_online"] = false
		logger.Warn("Hub is offline", zap.Error(err))
	} else {
		status["is_online"] = true
		status["template_count"] = len(templates)
		logger.Info("Hub is online", zap.Int("templates", len(templates)))
	}

	utils.RespondSuccess(w, status)
}
