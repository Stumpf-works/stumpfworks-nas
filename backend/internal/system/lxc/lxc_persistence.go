// Revision: 2025-12-01 | Author: Claude | Version: 1.0.0
package lxc

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateContainerPersistent creates an LXC container and saves it to the database
func (lm *LXCManager) CreateContainerPersistent(req ContainerCreateRequest) error {
	// Step 1: Check if container already exists in database
	var existing models.LXCContainer
	if err := database.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		return fmt.Errorf("container %s already exists in database", req.Name)
	}

	// Step 2: Create the container using existing LXC logic
	if err := lm.CreateContainer(req); err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Step 3: Save to database for persistence
	container := models.LXCContainer{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Template:     req.Template,
		Release:      req.Release,
		Architecture: req.Architecture,
		MemoryLimit:  req.MemoryLimit,
		CPULimit:     req.CPULimit,
		Autostart:    req.Autostart,
		NetworkMode:  req.NetworkMode,
		Bridge:       req.Bridge,
		Status:       "stopped", // Just created, not started yet
	}

	if err := database.DB.Create(&container).Error; err != nil {
		// Rollback: Delete the container from system
		lm.DeleteContainer(req.Name)
		return fmt.Errorf("failed to save container to database: %w", err)
	}

	logger.Info("LXC container created and saved to database",
		zap.String("name", req.Name),
		zap.String("id", container.ID),
		zap.String("template", req.Template))

	return nil
}

// DeleteContainerPersistent deletes an LXC container from both system and database
func (lm *LXCManager) DeleteContainerPersistent(name string) error {
	// Step 1: Delete from system
	if err := lm.DeleteContainer(name); err != nil {
		logger.Warn("Failed to delete container from system (may not exist)", zap.Error(err), zap.String("container", name))
	}

	// Step 2: Delete from database
	result := database.DB.Where("name = ?", name).Delete(&models.LXCContainer{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete container from database: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("container %s not found in database", name)
	}

	logger.Info("LXC container deleted from system and database", zap.String("name", name))
	return nil
}

// UpdateContainerStatus updates the container status in the database
func (lm *LXCManager) UpdateContainerStatus(name string, status string, ipv4 string, ipv6 string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if ipv4 != "" {
		updates["ipv4"] = ipv4
	}
	if ipv6 != "" {
		updates["ipv6"] = ipv6
	}

	result := database.DB.Model(&models.LXCContainer{}).Where("name = ?", name).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update container status: %w", result.Error)
	}

	return nil
}

// RestoreAllContainers restores all autostart containers from database on system startup
func (lm *LXCManager) RestoreAllContainers() error {
	if !lm.enabled {
		logger.Info("LXC not enabled, skipping container restoration")
		return nil
	}

	var containers []models.LXCContainer
	if err := database.DB.Where("autostart = ?", true).Find(&containers).Error; err != nil {
		return fmt.Errorf("failed to load containers from database: %w", err)
	}

	logger.Info("Restoring LXC containers from database", zap.Int("count", len(containers)))

	for _, container := range containers {
		logger.Info("Restoring container", zap.String("name", container.Name), zap.String("id", container.ID))

		// Check if container exists in system
		systemContainer, err := lm.GetContainer(container.Name)
		if err != nil {
			// Container doesn't exist in system, need to recreate it
			logger.Warn("Container not found in system, cannot auto-recreate",
				zap.String("name", container.Name),
				zap.Error(err))

			// Update status in database
			database.DB.Model(&container).Updates(map[string]interface{}{
				"status":     "error",
				"last_error": "Container not found in system. Please recreate manually.",
			})
			continue
		}

		// Container exists, start it if autostart is enabled
		if container.Autostart {
			if systemContainer.State != "RUNNING" {
				logger.Info("Starting autostart container", zap.String("name", container.Name))
				if err := lm.StartContainer(container.Name); err != nil {
					logger.Error("Failed to start container", zap.Error(err), zap.String("name", container.Name))
					database.DB.Model(&container).Updates(map[string]interface{}{
						"status":     "error",
						"last_error": err.Error(),
					})
					continue
				}
			}

			// Update status in database
			lm.UpdateContainerStatus(container.Name, "running", systemContainer.IPv4, systemContainer.IPv6)
			logger.Info("Container started successfully", zap.String("name", container.Name))
		}
	}

	return nil
}

// GetPersistedContainers returns all containers from the database
func (lm *LXCManager) GetPersistedContainers() ([]models.LXCContainer, error) {
	var containers []models.LXCContainer
	if err := database.DB.Find(&containers).Error; err != nil {
		return nil, fmt.Errorf("failed to load containers from database: %w", err)
	}
	return containers, nil
}
