// Revision: 2025-11-16 | Author: Claude | Version: 1.0.0
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"api-audit/internal/backend"
	"api-audit/internal/client"
	"api-audit/internal/discovery"
	"api-audit/internal/report"
	"api-audit/internal/testing"
)

const Version = "1.0.0"

func main() {
	// Parse flags
	var (
		url              = flag.String("url", "http://localhost:8080", "Base URL of NAS API")
		token            = flag.String("token", "", "JWT auth token")
		outputDir        = flag.String("output", "./audit_output", "Output directory for reports")
		format           = flag.String("format", "both", "Output format: json, md, both")
		endpointsOnly    = flag.Bool("endpoints-only", false, "Only test API endpoints")
		backendOnly      = flag.Bool("backend-only", false, "Only test backend functions")
		metricsOnly      = flag.Bool("metrics-only", false, "Only test Prometheus metrics")
		forceDestructive = flag.Bool("force-destructive", false, "Include destructive tests")
		category         = flag.String("category", "all", "Category to test: storage, sharing, network, users, all")
		timeout          = flag.Duration("timeout", 10*time.Second, "Request timeout")
		verbose          = flag.Bool("verbose", false, "Verbose output")
		help             = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	if *url == "" {
		fmt.Println("Error: --url is required")
		os.Exit(1)
	}

	fmt.Printf("StumpfWorks NAS API Audit Tool v%s\n", Version)
	fmt.Printf("Target: %s\n", *url)
	fmt.Println(strings.Repeat("=", 60))

	startTime := time.Now()

	// Create HTTP client
	httpClient := client.NewClient(*url, *token, *timeout)

	// Initialize report
	auditReport := &report.AuditReport{
		AuditInfo: report.AuditInfo{
			ToolVersion: Version,
			Timestamp:   startTime,
			TargetURL:   *url,
			NASVersion:  "1.1.1",
		},
		Endpoints:     []report.EndpointResult{},
		IssuesFound:   []report.Issue{},
	}

	// Test endpoints
	if !*backendOnly && !*metricsOnly {
		fmt.Println("\n📡 Testing API Endpoints...")
		endpoints := discovery.KnownEndpoints()
		filteredEndpoints := discovery.FilterEndpoints(endpoints, *forceDestructive, *category)

		tester := testing.NewEndpointTester(httpClient, *verbose)
		results := tester.TestAllEndpoints(filteredEndpoints)

		auditReport.Endpoints = results

		// Calculate summary
		auditReport.Summary.TotalEndpoints = len(endpoints)
		auditReport.Summary.TestedEndpoints = len(filteredEndpoints)

		for _, result := range results {
			switch result.TestResults.Status {
			case "PASS":
				auditReport.Summary.Passed++
			case "FAIL":
				auditReport.Summary.Failed++
			case "SKIP":
				auditReport.Summary.Skipped++
			}
		}

		if auditReport.Summary.TestedEndpoints > 0 {
			auditReport.Summary.CoveragePercent = float64(auditReport.Summary.Passed) / float64(auditReport.Summary.TestedEndpoints) * 100
		}

		fmt.Printf("✅ Passed: %d | ❌ Failed: %d | ⚠️ Skipped: %d\n",
			auditReport.Summary.Passed, auditReport.Summary.Failed, auditReport.Summary.Skipped)
	}

	// Test backend functions
	if !*endpointsOnly && !*metricsOnly {
		fmt.Println("\n🔧 Testing Backend Functions...")
		systemTester := backend.NewSystemTester(httpClient)
		auditReport.BackendFunctions = systemTester.TestBackendFunctions()
		fmt.Println("✅ Backend function tests complete")
	}

	// Test Prometheus metrics
	if !*endpointsOnly && !*backendOnly {
		fmt.Println("\n📊 Testing Prometheus Metrics...")
		systemTester := backend.NewSystemTester(httpClient)
		auditReport.PrometheusMetrics = systemTester.TestPrometheusMetrics()

		fmt.Printf("✅ Found: %d | ❌ Missing: %d\n",
			len(auditReport.PrometheusMetrics.MetricsFound),
			len(auditReport.PrometheusMetrics.MetricsMissing))
	}

	// Calculate performance metrics
	calculatePerformanceMetrics(auditReport)

	// Find issues
	findIssues(auditReport)

	// Set duration
	auditReport.AuditInfo.DurationSeconds = time.Since(startTime).Seconds()

	// Generate reports
	fmt.Printf("\n📝 Generating Reports in %s...\n", *outputDir)
	gen := report.NewGenerator(auditReport)

	if *format == "json" || *format == "both" {
		if err := gen.SaveJSON(*outputDir); err != nil {
			fmt.Printf("Error saving JSON: %v\n", err)
		} else {
			fmt.Println("✅ audit_report.json")
		}
	}

	if *format == "md" || *format == "both" {
		if err := gen.SaveMarkdown(*outputDir); err != nil {
			fmt.Printf("Error saving Markdown: %v\n", err)
		} else {
			fmt.Println("✅ AUDIT_API_REPORT.md")
		}
	}

	// Summary
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("✅ Audit complete in %.1fs\n", auditReport.AuditInfo.DurationSeconds)
	fmt.Printf("📊 Coverage: %.1f%%\n", auditReport.Summary.CoveragePercent)
	fmt.Printf("⚠️ Issues found: %d\n", len(auditReport.IssuesFound))

	// Exit code
	if auditReport.Summary.Failed > 0 {
		os.Exit(1)
	}
}

func calculatePerformanceMetrics(r *report.AuditReport) {
	if len(r.Endpoints) == 0 {
		return
	}

	var total int64
	var fastest, slowest *report.EndpointTiming

	for _, ep := range r.Endpoints {
		if ep.TestResults.Status != "PASS" {
			continue
		}

		total += ep.TestResults.ResponseTimeMs

		timing := &report.EndpointTiming{
			Path:   ep.Path,
			TimeMs: ep.TestResults.ResponseTimeMs,
		}

		if fastest == nil || timing.TimeMs < fastest.TimeMs {
			fastest = timing
		}
		if slowest == nil || timing.TimeMs > slowest.TimeMs {
			slowest = timing
		}
	}

	passedCount := int64(r.Summary.Passed)
	if passedCount > 0 {
		r.Performance.AverageResponseTimeMs = total / passedCount
	}

	r.Performance.FastestEndpoint = fastest
	r.Performance.SlowestEndpoint = slowest
}

func findIssues(r *report.AuditReport) {
	// Find 404 endpoints
	for _, ep := range r.Endpoints {
		if ep.TestResults.ResponseCode == 404 {
			r.IssuesFound = append(r.IssuesFound, report.Issue{
				Severity:       "HIGH",
				Type:           "MISSING_ENDPOINT",
				Description:    fmt.Sprintf("%s %s not implemented", ep.Method, ep.Path),
				Recommendation: "Implement missing endpoint",
			})
		}
	}

	// Find missing Prometheus metrics
	if len(r.PrometheusMetrics.MetricsMissing) > 0 {
		r.IssuesFound = append(r.IssuesFound, report.Issue{
			Severity:       "MEDIUM",
			Type:           "MISSING_METRIC",
			Description:    fmt.Sprintf("%d Prometheus metrics missing", len(r.PrometheusMetrics.MetricsMissing)),
			Recommendation: "Add missing NAS-specific metrics to /metrics endpoint",
		})
	}

	// Find slow endpoints (>1000ms)
	if r.Performance.SlowestEndpoint != nil && r.Performance.SlowestEndpoint.TimeMs > 1000 {
		r.IssuesFound = append(r.IssuesFound, report.Issue{
			Severity:       "LOW",
			Type:           "SLOW_ENDPOINT",
			Description:    fmt.Sprintf("%s responds in %dms (>1000ms)", r.Performance.SlowestEndpoint.Path, r.Performance.SlowestEndpoint.TimeMs),
			Recommendation: "Optimize endpoint performance",
		})
	}
}

func printHelp() {
	fmt.Println("StumpfWorks NAS API Audit Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  audit-tool [flags]")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  # Full audit")
	fmt.Println("  ./audit-tool --url http://localhost:8080 --token \"jwt_token\"")
	fmt.Println("\n  # Only test endpoints")
	fmt.Println("  ./audit-tool --url http://localhost:8080 --token \"jwt_token\" --endpoints-only")
	fmt.Println("\n  # Test specific category")
	fmt.Println("  ./audit-tool --url http://localhost:8080 --token \"jwt_token\" --category storage")
}

// Simple strings.Repeat implementation
var strings = struct {
	Repeat func(s string, count int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}
