// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package sysutil

import (
	"net"
	"strings"
)

// ValidateIP checks if a string is a valid IP address (IPv4 or IPv6)
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ValidateIPv4 checks if a string is a valid IPv4 address
func ValidateIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.To4() != nil
}

// ValidateIPv6 checks if a string is a valid IPv6 address
func ValidateIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.To4() == nil
}

// ValidateCIDR checks if a string is a valid CIDR notation
func ValidateCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

// IsPrivateIP checks if an IP address is in a private range
func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check private ranges
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",    // Loopback
		"169.254.0.0/16", // Link-local
		"fc00::/7",       // IPv6 Unique Local
		"fe80::/10",      // IPv6 Link-local
		"::1/128",        // IPv6 Loopback
	}

	for _, cidr := range privateRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(parsedIP) {
			return true
		}
	}

	return false
}

// IsLoopbackIP checks if an IP address is a loopback address
func IsLoopbackIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.IsLoopback()
}

// NormalizeIP normalizes an IP address string
// Returns the canonical form of the IP address
func NormalizeIP(ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ""
	}
	return parsedIP.String()
}

// ParseHostPort splits a host:port string
// Handles IPv6 addresses correctly (e.g., [::1]:8080)
func ParseHostPort(hostport string) (host string, port string, err error) {
	host, port, err = net.SplitHostPort(hostport)
	if err != nil {
		return "", "", err
	}
	return host, port, nil
}

// IsValidPort checks if a port number is valid (1-65535)
func IsValidPort(port int) bool {
	return port > 0 && port <= 65535
}

// IsValidHostname checks if a string is a valid hostname
func IsValidHostname(hostname string) bool {
	// Basic hostname validation
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}

	// Check if it's an IP address
	if ValidateIP(hostname) {
		return true
	}

	// Check hostname format
	for _, label := range strings.Split(hostname, ".") {
		if len(label) == 0 || len(label) > 63 {
			return false
		}

		// Check for valid characters (alphanumeric and hyphen)
		for i, c := range label {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
				(c >= '0' && c <= '9') || (c == '-' && i > 0 && i < len(label)-1)) {
				return false
			}
		}
	}

	return true
}
