// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package network

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"strings"
)

// FirewallManager manages firewall rules (iptables/nftables)
type FirewallManager struct {
	shell   executor.ShellExecutor
	enabled bool
	backend string // "iptables" or "nftables"
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	Chain       string `json:"chain"`       // INPUT, OUTPUT, FORWARD
	Protocol    string `json:"protocol"`    // tcp, udp, icmp, all
	Source      string `json:"source"`      // source IP/CIDR
	Destination string `json:"destination"` // destination IP/CIDR
	SourcePort  string `json:"source_port"`
	DestPort    string `json:"dest_port"`
	Action      string `json:"action"`      // ACCEPT, DROP, REJECT
	Comment     string `json:"comment"`
}

// NewFirewallManager creates a new firewall manager
func NewFirewallManager(shell executor.ShellExecutor) (*FirewallManager, error) {
	fm := &FirewallManager{
		shell: shell,
	}

	// Determine which backend is available
	if shell.CommandExists("nft") {
		fm.backend = "nftables"
		fm.enabled = true
	} else if shell.CommandExists("iptables") {
		fm.backend = "iptables"
		fm.enabled = true
	} else {
		return nil, fmt.Errorf("neither iptables nor nftables found")
	}

	return fm, nil
}

// IsEnabled returns whether firewall management is available
func (f *FirewallManager) IsEnabled() bool {
	return f.enabled
}

// GetBackend returns the firewall backend being used
func (f *FirewallManager) GetBackend() string {
	return f.backend
}

// ListRules lists all firewall rules
func (f *FirewallManager) ListRules(chain string) ([]FirewallRule, error) {
	if !f.enabled {
		return nil, fmt.Errorf("firewall not available")
	}

	if f.backend == "iptables" {
		return f.listIPTablesRules(chain)
	}

	return f.listNFTablesRules(chain)
}

// listIPTablesRules lists iptables rules
func (f *FirewallManager) listIPTablesRules(chain string) ([]FirewallRule, error) {
	args := []string{"-L"}
	if chain != "" {
		args = append(args, chain)
	}
	args = append(args, "-n", "-v", "--line-numbers")

	result, err := f.shell.Execute("iptables", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list iptables rules: %w", err)
	}

	var rules []FirewallRule
	lines := strings.Split(result.Stdout, "\n")
	currentChain := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Detect chain
		if strings.HasPrefix(line, "Chain") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				currentChain = fields[1]
			}
			continue
		}

		// Skip header lines
		if strings.HasPrefix(line, "num") || strings.HasPrefix(line, "pkts") || line == "" {
			continue
		}

		// Parse rule line
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		rule := FirewallRule{
			Chain:    currentChain,
			Action:   fields[3],
			Protocol: fields[4],
			Source:   fields[8],
		}

		if len(fields) > 9 {
			rule.Destination = fields[9]
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// listNFTablesRules lists nftables rules
func (f *FirewallManager) listNFTablesRules(chain string) ([]FirewallRule, error) {
	result, err := f.shell.Execute("nft", "list", "ruleset")
	if err != nil {
		return nil, fmt.Errorf("failed to list nftables rules: %w", err)
	}

	// Simplified parsing - nftables output is more complex
	var rules []FirewallRule
	lines := strings.Split(result.Stdout, "\n")
	currentChain := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "chain") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				currentChain = fields[1]
			}
			continue
		}

		if strings.Contains(line, "accept") || strings.Contains(line, "drop") {
			rule := FirewallRule{
				Chain: currentChain,
			}

			if strings.Contains(line, "accept") {
				rule.Action = "ACCEPT"
			} else if strings.Contains(line, "drop") {
				rule.Action = "DROP"
			}

			rules = append(rules, rule)
		}
	}

	return rules, nil
}

// AddRule adds a firewall rule
func (f *FirewallManager) AddRule(rule FirewallRule) error {
	if !f.enabled {
		return fmt.Errorf("firewall not available")
	}

	if f.backend == "iptables" {
		return f.addIPTablesRule(rule)
	}

	return f.addNFTablesRule(rule)
}

// addIPTablesRule adds an iptables rule
func (f *FirewallManager) addIPTablesRule(rule FirewallRule) error {
	args := []string{"-A", rule.Chain}

	if rule.Protocol != "" && rule.Protocol != "all" {
		args = append(args, "-p", rule.Protocol)
	}

	if rule.Source != "" {
		args = append(args, "-s", rule.Source)
	}

	if rule.Destination != "" {
		args = append(args, "-d", rule.Destination)
	}

	if rule.SourcePort != "" {
		args = append(args, "--sport", rule.SourcePort)
	}

	if rule.DestPort != "" {
		args = append(args, "--dport", rule.DestPort)
	}

	if rule.Comment != "" {
		args = append(args, "-m", "comment", "--comment", rule.Comment)
	}

	args = append(args, "-j", rule.Action)

	_, err := f.shell.Execute("iptables", args...)
	if err != nil {
		return fmt.Errorf("failed to add iptables rule: %w", err)
	}

	return nil
}

// addNFTablesRule adds an nftables rule
func (f *FirewallManager) addNFTablesRule(rule FirewallRule) error {
	// Build nftables rule
	ruleStr := ""

	if rule.Protocol != "" && rule.Protocol != "all" {
		ruleStr += fmt.Sprintf("ip protocol %s ", rule.Protocol)
	}

	if rule.Source != "" {
		ruleStr += fmt.Sprintf("ip saddr %s ", rule.Source)
	}

	if rule.Destination != "" {
		ruleStr += fmt.Sprintf("ip daddr %s ", rule.Destination)
	}

	if rule.DestPort != "" {
		ruleStr += fmt.Sprintf("%s dport %s ", rule.Protocol, rule.DestPort)
	}

	ruleStr += strings.ToLower(rule.Action)

	_, err := f.shell.Execute("nft", "add", "rule", "inet", "filter", rule.Chain, ruleStr)
	if err != nil {
		return fmt.Errorf("failed to add nftables rule: %w", err)
	}

	return nil
}

// DeleteRule deletes a firewall rule
func (f *FirewallManager) DeleteRule(chain string, ruleNum int) error {
	if !f.enabled {
		return fmt.Errorf("firewall not available")
	}

	if f.backend == "iptables" {
		_, err := f.shell.Execute("iptables", "-D", chain, fmt.Sprintf("%d", ruleNum))
		if err != nil {
			return fmt.Errorf("failed to delete iptables rule: %w", err)
		}
		return nil
	}

	// nftables delete is more complex - would need handle
	return fmt.Errorf("nftables delete not implemented")
}

// FlushChain flushes all rules in a chain
func (f *FirewallManager) FlushChain(chain string) error {
	if !f.enabled {
		return fmt.Errorf("firewall not available")
	}

	if f.backend == "iptables" {
		_, err := f.shell.Execute("iptables", "-F", chain)
		if err != nil {
			return fmt.Errorf("failed to flush chain: %w", err)
		}
		return nil
	}

	_, err := f.shell.Execute("nft", "flush", "chain", "inet", "filter", chain)
	if err != nil {
		return fmt.Errorf("failed to flush chain: %w", err)
	}

	return nil
}

// SetDefaultPolicy sets the default policy for a chain
func (f *FirewallManager) SetDefaultPolicy(chain string, policy string) error {
	if !f.enabled {
		return fmt.Errorf("firewall not available")
	}

	if policy != "ACCEPT" && policy != "DROP" {
		return fmt.Errorf("invalid policy: %s (must be ACCEPT or DROP)", policy)
	}

	if f.backend == "iptables" {
		_, err := f.shell.Execute("iptables", "-P", chain, policy)
		if err != nil {
			return fmt.Errorf("failed to set policy: %w", err)
		}
		return nil
	}

	_, err := f.shell.Execute("nft", "add", "chain", "inet", "filter", chain,
		fmt.Sprintf("{ policy %s ; }", strings.ToLower(policy)))
	if err != nil {
		return fmt.Errorf("failed to set policy: %w", err)
	}

	return nil
}

// AllowPort creates a rule to allow traffic on a specific port
func (f *FirewallManager) AllowPort(port int, protocol string, comment string) error {
	rule := FirewallRule{
		Chain:    "INPUT",
		Protocol: protocol,
		DestPort: fmt.Sprintf("%d", port),
		Action:   "ACCEPT",
		Comment:  comment,
	}

	return f.AddRule(rule)
}

// BlockIP blocks all traffic from an IP address
func (f *FirewallManager) BlockIP(ip string, comment string) error {
	rule := FirewallRule{
		Chain:   "INPUT",
		Source:  ip,
		Action:  "DROP",
		Comment: comment,
	}

	return f.AddRule(rule)
}

// SaveRules saves current rules to file (persistence)
func (f *FirewallManager) SaveRules() error {
	if !f.enabled {
		return fmt.Errorf("firewall not available")
	}

	if f.backend == "iptables" {
		if f.shell.CommandExists("iptables-save") {
			result, err := f.shell.Execute("iptables-save")
			if err != nil {
				return fmt.Errorf("failed to save iptables rules: %w", err)
			}

			// Write to file
			_, err = f.shell.Execute("bash", "-c",
				fmt.Sprintf("echo '%s' > /etc/iptables/rules.v4", result.Stdout))
			if err != nil {
				return fmt.Errorf("failed to write rules file: %w", err)
			}
		}
		return nil
	}

	// nftables
	if f.shell.CommandExists("nft") {
		result, err := f.shell.Execute("nft", "list", "ruleset")
		if err != nil {
			return fmt.Errorf("failed to save nftables rules: %w", err)
		}

		_, err = f.shell.Execute("bash", "-c",
			fmt.Sprintf("echo '%s' > /etc/nftables.conf", result.Stdout))
		if err != nil {
			return fmt.Errorf("failed to write rules file: %w", err)
		}
	}

	return nil
}

// RestoreRules restores rules from file
func (f *FirewallManager) RestoreRules() error {
	if !f.enabled {
		return fmt.Errorf("firewall not available")
	}

	if f.backend == "iptables" {
		if f.shell.CommandExists("iptables-restore") {
			_, err := f.shell.Execute("iptables-restore", "/etc/iptables/rules.v4")
			if err != nil {
				return fmt.Errorf("failed to restore iptables rules: %w", err)
			}
		}
		return nil
	}

	_, err := f.shell.Execute("nft", "-f", "/etc/nftables.conf")
	if err != nil {
		return fmt.Errorf("failed to restore nftables rules: %w", err)
	}

	return nil
}
