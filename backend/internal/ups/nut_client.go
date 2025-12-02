// Revision: 2025-12-02 | Author: Claude | Version: 1.2.0
package ups

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
)

// QueryUPS queries the UPS via NUT protocol and returns current status
func (s *Service) QueryUPS(config *models.UPSConfig) (*UPSStatus, error) {
	// Connect to NUT server
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", config.UPSHost, config.UPSPort),
		5*time.Second)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to NUT server: %w", err)
	}
	defer conn.Close()

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	// Authenticate if credentials provided
	if config.UPSUsername != "" {
		if err := s.nutAuth(conn, config.UPSUsername, config.UPSPassword); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Query UPS variables
	vars, err := s.nutListVars(conn, config.UPSName)
	if err != nil {
		return nil, fmt.Errorf("failed to query UPS variables: %w", err)
	}

	// Parse status
	status := &UPSStatus{
		LastUpdate: time.Now(),
	}

	// Parse standard UPS variables
	if val, ok := vars["ups.status"]; ok {
		status.Status = val
		// OL = Online, OB = On Battery, LB = Low Battery
		status.Online = strings.Contains(val, "OL")
	}

	if val, ok := vars["battery.charge"]; ok {
		if charge, err := strconv.Atoi(val); err == nil {
			status.BatteryCharge = charge
		}
	}

	if val, ok := vars["battery.runtime"]; ok {
		if runtime, err := strconv.Atoi(val); err == nil {
			status.Runtime = runtime
		}
	}

	if val, ok := vars["ups.load"]; ok {
		if load, err := strconv.Atoi(val); err == nil {
			status.LoadPercent = load
		}
	}

	if val, ok := vars["input.voltage"]; ok {
		if voltage, err := strconv.ParseFloat(val, 64); err == nil {
			status.InputVoltage = voltage
		}
	}

	if val, ok := vars["output.voltage"]; ok {
		if voltage, err := strconv.ParseFloat(val, 64); err == nil {
			status.OutputVoltage = voltage
		}
	}

	if val, ok := vars["ups.temperature"]; ok {
		if temp, err := strconv.ParseFloat(val, 64); err == nil {
			status.Temperature = temp
		}
	}

	if val, ok := vars["ups.model"]; ok {
		status.Model = val
	}

	if val, ok := vars["ups.mfr"]; ok {
		status.Manufacturer = val
	}

	if val, ok := vars["ups.serial"]; ok {
		status.Serial = val
	}

	return status, nil
}

// nutAuth authenticates with the NUT server
func (s *Service) nutAuth(conn net.Conn, username, password string) error {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send USERNAME command
	if _, err := fmt.Fprintf(writer, "USERNAME %s\n", username); err != nil {
		return err
	}
	writer.Flush()

	// Read response
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasPrefix(response, "OK") {
		return fmt.Errorf("username not accepted: %s", response)
	}

	// Send PASSWORD command
	if _, err := fmt.Fprintf(writer, "PASSWORD %s\n", password); err != nil {
		return err
	}
	writer.Flush()

	// Read response
	response, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasPrefix(response, "OK") {
		return fmt.Errorf("password not accepted: %s", response)
	}

	return nil
}

// nutListVars lists all variables for a UPS
func (s *Service) nutListVars(conn net.Conn, upsName string) (map[string]string, error) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send LIST VAR command
	if _, err := fmt.Fprintf(writer, "LIST VAR %s\n", upsName); err != nil {
		return nil, err
	}
	writer.Flush()

	vars := make(map[string]string)

	// Read responses
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		// End of list
		if line == "BEGIN LIST VAR "+upsName || line == "" {
			continue
		}
		if strings.HasPrefix(line, "END LIST VAR") {
			break
		}

		// Parse variable line: VAR ups.name "value"
		if strings.HasPrefix(line, "VAR "+upsName) {
			parts := strings.SplitN(line, " ", 3)
			if len(parts) >= 3 {
				varName := parts[1]
				varValue := strings.Trim(parts[2], "\"")
				// Remove UPS name prefix from variable name
				varName = strings.TrimPrefix(varName, upsName+".")
				vars[varName] = varValue
			}
		}
	}

	return vars, nil
}

// nutGetVar gets a single variable from the UPS
func (s *Service) nutGetVar(conn net.Conn, upsName, varName string) (string, error) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send GET VAR command
	if _, err := fmt.Fprintf(writer, "GET VAR %s %s\n", upsName, varName); err != nil {
		return "", err
	}
	writer.Flush()

	// Read response: VAR ups.name variable.name "value"
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "VAR") {
		parts := strings.SplitN(response, " ", 3)
		if len(parts) >= 3 {
			value := strings.Trim(parts[2], "\"")
			return value, nil
		}
	}

	if strings.HasPrefix(response, "ERR") {
		return "", fmt.Errorf("NUT error: %s", response)
	}

	return "", fmt.Errorf("unexpected response: %s", response)
}
