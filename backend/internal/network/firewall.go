package network

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// FirewallStatus represents the firewall status
type FirewallStatus struct {
	Enabled      bool           `json:"enabled"`
	DefaultIncoming string      `json:"defaultIncoming"`
	DefaultOutgoing string      `json:"defaultOutgoing"`
	DefaultRouted   string      `json:"defaultRouted"`
	Rules        []FirewallRule `json:"rules"`
}

// GetFirewallStatus returns the current firewall status
func GetFirewallStatus() (*FirewallStatus, error) {
	// Check if ufw is installed
	if _, err := exec.LookPath("ufw"); err != nil {
		return nil, fmt.Errorf("ufw is not installed")
	}

	// Get status
	cmd := exec.Command("ufw", "status", "verbose")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get firewall status: %s", string(output))
	}

	status := &FirewallStatus{
		Rules: []FirewallRule{},
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "Status:") {
			status.Enabled = strings.Contains(line, "active")
		} else if strings.HasPrefix(line, "Default:") {
			// Parse default policies: "Default: deny (incoming), allow (outgoing), disabled (routed)"
			parts := strings.Split(line, ",")
			for _, part := range parts {
				if strings.Contains(part, "incoming") {
					if strings.Contains(part, "allow") {
						status.DefaultIncoming = "allow"
					} else {
						status.DefaultIncoming = "deny"
					}
				} else if strings.Contains(part, "outgoing") {
					if strings.Contains(part, "allow") {
						status.DefaultOutgoing = "allow"
					} else {
						status.DefaultOutgoing = "deny"
					}
				} else if strings.Contains(part, "routed") {
					if strings.Contains(part, "allow") {
						status.DefaultRouted = "allow"
					} else {
						status.DefaultRouted = "deny"
					}
				}
			}
		}
	}

	// Get numbered rules
	cmd = exec.Command("ufw", "status", "numbered")
	output, err = cmd.CombinedOutput()
	if err == nil {
		status.Rules = parseFirewallRules(string(output))
	}

	return status, nil
}

// parseFirewallRules parses ufw status numbered output
func parseFirewallRules(output string) []FirewallRule {
	var rules []FirewallRule
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip header lines
		if !strings.HasPrefix(line, "[") {
			continue
		}

		// Parse line: "[ 1] 22/tcp                     ALLOW IN    Anywhere"
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		// Extract rule number
		numStr := strings.Trim(parts[0], "[]")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}

		rule := FirewallRule{
			Number: num,
		}

		// Parse port/protocol
		portProto := parts[1]
		if strings.Contains(portProto, "/") {
			pp := strings.Split(portProto, "/")
			rule.Port = pp[0]
			rule.Protocol = pp[1]
		} else {
			rule.Port = portProto
		}

		// Parse action (ALLOW/DENY/REJECT)
		rule.Action = strings.ToLower(parts[2])

		// Parse direction and source/dest
		if len(parts) >= 5 {
			direction := parts[3]
			if direction == "IN" {
				rule.To = "Anywhere"
				if len(parts) >= 5 {
					rule.From = strings.Join(parts[4:], " ")
				}
			} else if direction == "OUT" {
				rule.From = "Anywhere"
				if len(parts) >= 5 {
					rule.To = strings.Join(parts[4:], " ")
				}
			}
		}

		rules = append(rules, rule)
	}

	return rules
}

// EnableFirewall enables the firewall
func EnableFirewall() error {
	cmd := exec.Command("ufw", "--force", "enable")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable firewall: %s", string(output))
	}
	return nil
}

// DisableFirewall disables the firewall
func DisableFirewall() error {
	cmd := exec.Command("ufw", "disable")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to disable firewall: %s", string(output))
	}
	return nil
}

// AddFirewallRule adds a new firewall rule
func AddFirewallRule(action, port, protocol, from, to string) error {
	args := []string{}

	// Build command
	if action == "allow" {
		args = append(args, "allow")
	} else if action == "deny" {
		args = append(args, "deny")
	} else if action == "reject" {
		args = append(args, "reject")
	} else {
		return fmt.Errorf("invalid action: %s", action)
	}

	// Add from/to
	if from != "" && from != "any" {
		args = append(args, "from", from)
	}
	if to != "" && to != "any" {
		args = append(args, "to", to)
	}

	// Add port and protocol
	if port != "" {
		args = append(args, "port", port)
	}
	if protocol != "" {
		args = append(args, "proto", protocol)
	}

	cmd := exec.Command("ufw", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add rule: %s", string(output))
	}

	return nil
}

// DeleteFirewallRule deletes a firewall rule by number
func DeleteFirewallRule(ruleNumber int) error {
	cmd := exec.Command("ufw", "--force", "delete", strconv.Itoa(ruleNumber))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete rule: %s", string(output))
	}
	return nil
}

// SetDefaultPolicy sets the default policy for incoming/outgoing/routed traffic
func SetDefaultPolicy(direction, policy string) error {
	if direction != "incoming" && direction != "outgoing" && direction != "routed" {
		return fmt.Errorf("invalid direction: %s", direction)
	}
	if policy != "allow" && policy != "deny" && policy != "reject" {
		return fmt.Errorf("invalid policy: %s", policy)
	}

	cmd := exec.Command("ufw", "default", policy, direction)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set default policy: %s", string(output))
	}

	return nil
}

// ResetFirewall resets all firewall rules
func ResetFirewall() error {
	cmd := exec.Command("ufw", "--force", "reset")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reset firewall: %s", string(output))
	}
	return nil
}

// AllowService allows a predefined service
func AllowService(service string) error {
	cmd := exec.Command("ufw", "allow", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to allow service: %s", string(output))
	}
	return nil
}

// DenyService denies a predefined service
func DenyService(service string) error {
	cmd := exec.Command("ufw", "deny", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to deny service: %s", string(output))
	}
	return nil
}
