package timemachine

import (
	"fmt"
	"os"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

const (
	avahiServiceFile = "/etc/avahi/services/timemachine.service"
)

// enableAvahi enables Avahi service discovery for Time Machine
func (m *Manager) enableAvahi() error {
	logger.Info("Enabling Avahi service discovery for Time Machine")

	config, err := m.GetConfig()
	if err != nil {
		return err
	}

	// Create Avahi service file
	serviceXML := fmt.Sprintf(`<?xml version="1.0" standalone='no'?>
<!DOCTYPE service-group SYSTEM "avahi-service.dtd">
<service-group>
  <name replace-wildcards="yes">%%h - Time Machine</name>
  <service>
    <type>_smb._tcp</type>
    <port>445</port>
  </service>
  <service>
    <type>_device-info._tcp</type>
    <port>0</port>
    <txt-record>model=TimeCapsule8,119</txt-record>
  </service>
  <service>
    <type>_adisk._tcp</type>
    <port>9</port>
    <txt-record>dk0=adVN=%s,adVF=0x82</txt-record>
    <txt-record>sys=waMA=0,adVF=0x100</txt-record>
  </service>
</service-group>
`, config.ShareName)

	// Write service file
	if err := os.WriteFile(avahiServiceFile, []byte(serviceXML), 0644); err != nil {
		return fmt.Errorf("failed to write Avahi service file: %w", err)
	}

	// Restart Avahi daemon
	if _, err := m.shell.Execute("systemctl", "restart", "avahi-daemon"); err != nil {
		logger.Warn("Failed to restart Avahi daemon", zap.Error(err))
		// Non-fatal error
	}

	logger.Info("Avahi service discovery enabled")
	return nil
}

// disableAvahi disables Avahi service discovery for Time Machine
func (m *Manager) disableAvahi() error {
	logger.Info("Disabling Avahi service discovery for Time Machine")

	// Remove service file
	if err := os.Remove(avahiServiceFile); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove Avahi service file: %w", err)
		}
	}

	// Restart Avahi daemon
	if _, err := m.shell.Execute("systemctl", "restart", "avahi-daemon"); err != nil {
		logger.Warn("Failed to restart Avahi daemon", zap.Error(err))
	}

	logger.Info("Avahi service discovery disabled")
	return nil
}
