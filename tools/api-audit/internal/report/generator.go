// Revision: 2025-11-16 | Author: Claude | Version: 1.0.0
package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generator generates audit reports
type Generator struct {
	report *AuditReport
}

// NewGenerator creates a new report generator
func NewGenerator(report *AuditReport) *Generator {
	return &Generator{report: report}
}

// SaveJSON saves the report as JSON
func (g *Generator) SaveJSON(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	data, err := json.MarshalIndent(g.report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	filename := filepath.Join(outputDir, "audit_report.json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("write JSON file: %w", err)
	}

	return nil
}

// SaveMarkdown saves the report as Markdown
func (g *Generator) SaveMarkdown(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	md := g.generateMarkdown()

	filename := filepath.Join(outputDir, "AUDIT_API_REPORT.md")
	if err := os.WriteFile(filename, []byte(md), 0644); err != nil {
		return fmt.Errorf("write markdown file: %w", err)
	}

	return nil
}

func (g *Generator) generateMarkdown() string {
	var md strings.Builder

	// Header
	md.WriteString("# StumpfWorks NAS - API Audit Report\n\n")

	// Executive Summary
	md.WriteString("## Executive Summary\n\n")
	md.WriteString(fmt.Sprintf("- **Datum:** %s\n", g.report.AuditInfo.Timestamp.Format("2006-01-02 15:04:05")))
	md.WriteString(fmt.Sprintf("- **NAS Version:** %s\n", g.report.AuditInfo.NASVersion))
	md.WriteString(fmt.Sprintf("- **Tool Version:** %s\n", g.report.AuditInfo.ToolVersion))
	md.WriteString(fmt.Sprintf("- **Dauer:** %.1fs\n", g.report.AuditInfo.DurationSeconds))
	md.WriteString(fmt.Sprintf("- **Getestete Endpoints:** %d\n", g.report.Summary.TestedEndpoints))
	md.WriteString(fmt.Sprintf("- **Erfolgsrate:** %.1f%%\n", g.report.Summary.CoveragePercent))
	md.WriteString(fmt.Sprintf("- **Durchschnittliche Response Time:** %dms\n\n", g.report.Performance.AverageResponseTimeMs))

	// Endpoint Overview
	md.WriteString("## Endpoint Übersicht\n\n")

	// Passed tests
	md.WriteString(fmt.Sprintf("### ✅ Funktionierende Endpoints (%d)\n\n", g.report.Summary.Passed))
	md.WriteString("| Method | Path | Response Time | Status |\n")
	md.WriteString("|--------|------|---------------|--------|\n")
	for _, ep := range g.report.Endpoints {
		if ep.TestResults.Status == "PASS" {
			md.WriteString(fmt.Sprintf("| %s | %s | %dms | %d |\n",
				ep.Method, ep.Path, ep.TestResults.ResponseTimeMs, ep.TestResults.ResponseCode))
		}
	}
	md.WriteString("\n")

	// Failed tests
	if g.report.Summary.Failed > 0 {
		md.WriteString(fmt.Sprintf("### ❌ Fehlgeschlagene Tests (%d)\n\n", g.report.Summary.Failed))
		md.WriteString("| Method | Path | Issue | Details |\n")
		md.WriteString("|--------|------|-------|----------|\n")
		for _, ep := range g.report.Endpoints {
			if ep.TestResults.Status == "FAIL" {
				errors := strings.Join(ep.TestResults.Errors, ", ")
				md.WriteString(fmt.Sprintf("| %s | %s | %s | %d |\n",
					ep.Method, ep.Path, errors, ep.TestResults.ResponseCode))
			}
		}
		md.WriteString("\n")
	}

	// Skipped tests
	if g.report.Summary.Skipped > 0 {
		md.WriteString(fmt.Sprintf("### ⚠️ Übersprungene Tests (%d)\n\n", g.report.Summary.Skipped))
		for _, ep := range g.report.Endpoints {
			if ep.TestResults.Status == "SKIP" {
				md.WriteString(fmt.Sprintf("- %s %s (destructive, needs --force flag)\n", ep.Method, ep.Path))
			}
		}
		md.WriteString("\n")
	}

	// Backend Functions
	md.WriteString("## Backend Functions\n\n")

	md.WriteString("### Storage\n\n")
	md.WriteString(fmt.Sprintf("- ✅ ZFS verfügbar: %v\n", g.report.BackendFunctions.SystemLibrary.Storage.ZFSAvailable))
	md.WriteString(fmt.Sprintf("- ✅ Pools gefunden: %d\n", g.report.BackendFunctions.SystemLibrary.Storage.PoolsCount))
	md.WriteString(fmt.Sprintf("- ✅ SMART verfügbar: %v\n", g.report.BackendFunctions.SystemLibrary.Storage.SmartAvailable))
	md.WriteString(fmt.Sprintf("- ✅ Disks erkannt: %d\n\n", g.report.BackendFunctions.SystemLibrary.Storage.DisksCount))

	md.WriteString("### Sharing\n\n")
	md.WriteString(fmt.Sprintf("- ✅ Samba läuft: %v", g.report.BackendFunctions.SystemLibrary.Sharing.SambaRunning))
	if g.report.BackendFunctions.SystemLibrary.Sharing.SambaVersion != "" {
		md.WriteString(fmt.Sprintf(" (v%s)", g.report.BackendFunctions.SystemLibrary.Sharing.SambaVersion))
	}
	md.WriteString("\n")
	md.WriteString(fmt.Sprintf("- ✅ Aktive Shares: %d\n", g.report.BackendFunctions.SystemLibrary.Sharing.SharesCount))
	md.WriteString(fmt.Sprintf("- ✅ Verbindungen: %d\n\n", g.report.BackendFunctions.SystemLibrary.Sharing.ConnectionsCount))

	md.WriteString("### Network\n\n")
	md.WriteString(fmt.Sprintf("- ✅ Interfaces: %d\n", g.report.BackendFunctions.SystemLibrary.Network.InterfacesCount))
	md.WriteString(fmt.Sprintf("- ✅ Firewall aktiv: %v\n\n", g.report.BackendFunctions.SystemLibrary.Network.FirewallActive))

	// Prometheus Metrics
	md.WriteString("## Prometheus Metrics\n\n")
	md.WriteString(fmt.Sprintf("**Endpoint erreichbar:** %v\n\n", g.report.PrometheusMetrics.EndpointReachable))

	if len(g.report.PrometheusMetrics.MetricsFound) > 0 {
		md.WriteString(fmt.Sprintf("**Vorhanden (%d):**\n\n", len(g.report.PrometheusMetrics.MetricsFound)))
		for _, metric := range g.report.PrometheusMetrics.MetricsFound {
			md.WriteString(fmt.Sprintf("- ✅ %s\n", metric))
		}
		md.WriteString("\n")
	}

	if len(g.report.PrometheusMetrics.MetricsMissing) > 0 {
		md.WriteString(fmt.Sprintf("**Fehlend (%d):**\n\n", len(g.report.PrometheusMetrics.MetricsMissing)))
		for _, metric := range g.report.PrometheusMetrics.MetricsMissing {
			md.WriteString(fmt.Sprintf("- ❌ %s\n", metric))
		}
		md.WriteString("\n")
	}

	// Performance
	md.WriteString("## Performance\n\n")
	if g.report.Performance.FastestEndpoint != nil {
		md.WriteString(fmt.Sprintf("- **Schnellster:** %s (%dms)\n",
			g.report.Performance.FastestEndpoint.Path, g.report.Performance.FastestEndpoint.TimeMs))
	}
	if g.report.Performance.SlowestEndpoint != nil {
		md.WriteString(fmt.Sprintf("- **Langsamster:** %s (%dms)\n",
			g.report.Performance.SlowestEndpoint.Path, g.report.Performance.SlowestEndpoint.TimeMs))
	}
	md.WriteString(fmt.Sprintf("- **Durchschnitt:** %dms\n\n", g.report.Performance.AverageResponseTimeMs))

	// Issues
	if len(g.report.IssuesFound) > 0 {
		md.WriteString("## Kritische Issues\n\n")
		for _, issue := range g.report.IssuesFound {
			emoji := "🟡"
			if issue.Severity == "HIGH" {
				emoji = "🔴"
			} else if issue.Severity == "LOW" {
				emoji = "🟢"
			}

			md.WriteString(fmt.Sprintf("### %s %s: %s\n\n", emoji, issue.Severity, issue.Type))
			md.WriteString(fmt.Sprintf("**Problem:** %s\n\n", issue.Description))
			md.WriteString(fmt.Sprintf("**Empfehlung:** %s\n\n", issue.Recommendation))
		}
	}

	// Footer
	md.WriteString("---\n\n")
	md.WriteString(fmt.Sprintf("*Generiert von StumpfWorks API Audit Tool v%s*\n", g.report.AuditInfo.ToolVersion))

	return md.String()
}
