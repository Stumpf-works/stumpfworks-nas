// Revision: 2025-12-01 | Author: Claude | Version: 1.0.0
package network

import (
	"fmt"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateBridgePersistent creates a bridge and saves it to the database
// This ensures the bridge will be restored after reboot
func CreateBridgePersistent(name string, ports []string, ipAddress string, gateway string, description string) error {
	// Step 1: Check if bridge already exists in database
	var existing models.NetworkBridge
	if err := database.DB.Where("name = ?", name).First(&existing).Error; err == nil {
		return fmt.Errorf("bridge %s already exists in database", name)
	}

	// Step 2: Create the bridge in the system
	if err := CreateBridge(name, ports); err != nil {
		return fmt.Errorf("failed to create bridge in system: %w", err)
	}

	// Step 3: Apply IP configuration if specified
	if ipAddress != "" {
		// Parse CIDR to get address and netmask
		parts := strings.Split(ipAddress, "/")
		if len(parts) == 2 {
			addr := parts[0]
			netmask := parts[1]
			if err := ConfigureStaticIP(name, addr, convertCIDRToNetmask(netmask), gateway); err != nil {
				logger.Warn("Failed to apply IP configuration to bridge", zap.Error(err), zap.String("bridge", name))
			}
		}
	}

	// Step 4: Save to database for persistence
	bridge := models.NetworkBridge{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Ports:       strings.Join(ports, ","),
		IPAddress:   ipAddress,
		Gateway:     gateway,
		Autostart:   true,
		Status:      "active",
	}

	if err := database.DB.Create(&bridge).Error; err != nil {
		// Rollback: Delete the bridge from system
		DeleteBridge(name)
		return fmt.Errorf("failed to save bridge to database: %w", err)
	}

	logger.Info("Bridge created and saved to database",
		zap.String("name", name),
		zap.String("id", bridge.ID),
		zap.Strings("ports", ports))

	return nil
}

// DeleteBridgePersistent deletes a bridge from both system and database
func DeleteBridgePersistent(name string) error {
	// Step 1: Delete from system
	if err := DeleteBridge(name); err != nil {
		logger.Warn("Failed to delete bridge from system (may not exist)", zap.Error(err), zap.String("bridge", name))
	}

	// Step 2: Delete from database
	result := database.DB.Where("name = ?", name).Delete(&models.NetworkBridge{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete bridge from database: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("bridge %s not found in database", name)
	}

	logger.Info("Bridge deleted from system and database", zap.String("name", name))
	return nil
}

// RestoreAllBridges restores all bridges from database on system startup
// This should be called during system initialization
func RestoreAllBridges() error {
	var bridges []models.NetworkBridge
	if err := database.DB.Where("autostart = ?", true).Find(&bridges).Error; err != nil {
		return fmt.Errorf("failed to load bridges from database: %w", err)
	}

	logger.Info("Restoring network bridges from database", zap.Int("count", len(bridges)))

	for _, bridge := range bridges {
		logger.Info("Restoring bridge", zap.String("name", bridge.Name), zap.String("id", bridge.ID))

		// Parse ports
		var ports []string
		if bridge.Ports != "" {
			ports = strings.Split(bridge.Ports, ",")
		}

		// Create the bridge
		if err := CreateBridge(bridge.Name, ports); err != nil {
			logger.Error("Failed to restore bridge", zap.Error(err), zap.String("name", bridge.Name))
			// Update status in database
			database.DB.Model(&bridge).Updates(map[string]interface{}{
				"status":     "error",
				"last_error": err.Error(),
			})
			continue
		}

		// Apply IP configuration if specified
		if bridge.IPAddress != "" {
			parts := strings.Split(bridge.IPAddress, "/")
			if len(parts) == 2 {
				addr := parts[0]
				netmask := parts[1]
				if err := ConfigureStaticIP(bridge.Name, addr, convertCIDRToNetmask(netmask), bridge.Gateway); err != nil {
					logger.Warn("Failed to apply IP configuration to restored bridge",
						zap.Error(err),
						zap.String("bridge", bridge.Name))
				}
			}
		}

		// Update status in database
		database.DB.Model(&bridge).Updates(map[string]interface{}{
			"status":     "active",
			"last_error": "",
		})

		logger.Info("Bridge restored successfully", zap.String("name", bridge.Name))
	}

	return nil
}

// GetPersistedBridges returns all bridges from the database
func GetPersistedBridges() ([]models.NetworkBridge, error) {
	var bridges []models.NetworkBridge
	if err := database.DB.Find(&bridges).Error; err != nil {
		return nil, fmt.Errorf("failed to load bridges from database: %w", err)
	}
	return bridges, nil
}

// convertCIDRToNetmask converts CIDR notation to netmask (e.g., "24" -> "255.255.255.0")
func convertCIDRToNetmask(cidr string) string {
	// Common CIDR to netmask conversions
	cidrMap := map[string]string{
		"8":  "255.0.0.0",
		"16": "255.255.0.0",
		"24": "255.255.255.0",
		"25": "255.255.255.128",
		"26": "255.255.255.192",
		"27": "255.255.255.224",
		"28": "255.255.255.240",
		"29": "255.255.255.248",
		"30": "255.255.255.252",
		"32": "255.255.255.255",
	}

	if netmask, ok := cidrMap[cidr]; ok {
		return netmask
	}

	// Default to /24 if unknown
	return "255.255.255.0"
}
