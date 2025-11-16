package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// PluginLogger is an io.Writer that logs plugin output
type PluginLogger struct {
	pluginID string
	stream   string // "stdout" or "stderr"
}

// NewPluginLogger creates a new plugin logger
func NewPluginLogger(pluginID, stream string) *PluginLogger {
	return &PluginLogger{
		pluginID: pluginID,
		stream:   stream,
	}
}

// Write implements io.Writer interface
func (p *PluginLogger) Write(data []byte) (n int, err error) {
	// Convert bytes to string and remove trailing newline
	message := strings.TrimRight(string(data), "\n")

	// Log with appropriate level based on stream
	if p.stream == "stderr" {
		Error(fmt.Sprintf("[Plugin:%s] %s", p.pluginID, message))
	} else {
		Info(fmt.Sprintf("[Plugin:%s] %s", p.pluginID, message),
			zap.String("pluginID", p.pluginID))
	}

	return len(data), nil
}
