package main

import (
	"fmt"
	"os"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/sysutil"
)

func main() {
	fmt.Println("StumpfWorks NAS - System Health Check")
	fmt.Println("======================================")
	fmt.Println()

	// Perform health check
	report := sysutil.PerformSystemHealthCheck()

	// Print report
	report.PrintReport()

	// Save JSON report
	jsonReport, err := report.ToJSON()
	if err != nil {
		fmt.Printf("Error generating JSON: %v\n", err)
		os.Exit(1)
	}

	filename := "health-check.json"
	if err := os.WriteFile(filename, []byte(jsonReport), 0644); err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("JSON report saved to: %s\n", filename)

	// Exit with appropriate code
	if report.OverallStatus == "unhealthy" {
		os.Exit(1)
	}
}
