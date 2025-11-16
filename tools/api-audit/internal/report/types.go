// Revision: 2025-11-16 | Author: Claude | Version: 1.0.0
package report

import "time"

// AuditReport represents the complete audit report
type AuditReport struct {
	AuditInfo         AuditInfo                `json:"audit_info"`
	Summary           Summary                  `json:"summary"`
	Endpoints         []EndpointResult         `json:"endpoints"`
	BackendFunctions  BackendFunctions         `json:"backend_functions"`
	PrometheusMetrics PrometheusMetrics        `json:"prometheus_metrics"`
	SecurityTests     SecurityTests            `json:"security_tests"`
	Performance       Performance              `json:"performance"`
	IssuesFound       []Issue                  `json:"issues_found"`
}

// AuditInfo contains metadata about the audit run
type AuditInfo struct {
	ToolVersion     string    `json:"tool_version"`
	Timestamp       time.Time `json:"timestamp"`
	DurationSeconds float64   `json:"duration_seconds"`
	TargetURL       string    `json:"target_url"`
	NASVersion      string    `json:"nas_version"`
}

// Summary contains high-level statistics
type Summary struct {
	TotalEndpoints   int     `json:"total_endpoints"`
	TestedEndpoints  int     `json:"tested_endpoints"`
	Passed           int     `json:"passed"`
	Failed           int     `json:"failed"`
	Skipped          int     `json:"skipped"`
	CoveragePercent  float64 `json:"coverage_percent"`
}

// EndpointResult represents test results for a single endpoint
type EndpointResult struct {
	Method       string      `json:"method"`
	Path         string      `json:"path"`
	Description  string      `json:"description"`
	AuthRequired bool        `json:"auth_required"`
	TestResults  TestResults `json:"test_results"`
}

// TestResults contains detailed test results
type TestResults struct {
	Status            string                 `json:"status"` // PASS, FAIL, SKIP
	ResponseCode      int                    `json:"response_code,omitempty"`
	ResponseTimeMs    int64                  `json:"response_time_ms,omitempty"`
	ResponseBody      interface{}            `json:"response_body,omitempty"`
	Headers           map[string]string      `json:"headers,omitempty"`
	Errors            []string               `json:"errors,omitempty"`
	ValidRequestTest  *RequestTest           `json:"valid_request_test,omitempty"`
	InvalidRequestTest *RequestTest          `json:"invalid_request_test,omitempty"`
}

// RequestTest represents a single request test
type RequestTest struct {
	RequestBody    interface{}       `json:"request_body,omitempty"`
	ResponseCode   int               `json:"response_code"`
	ResponseTimeMs int64             `json:"response_time_ms"`
	ErrorMessage   string            `json:"error_message,omitempty"`
}

// BackendFunctions contains results from backend function tests
type BackendFunctions struct {
	SystemLibrary SystemLibrary `json:"system_library"`
}

// SystemLibrary contains system library test results
type SystemLibrary struct {
	Storage Storage `json:"storage"`
	Sharing Sharing `json:"sharing"`
	Network Network `json:"network"`
	Users   Users   `json:"users"`
}

// Storage contains storage subsystem results
type Storage struct {
	ZFSAvailable  bool `json:"zfs_available"`
	PoolsCount    int  `json:"pools_count"`
	SmartAvailable bool `json:"smart_available"`
	DisksCount    int  `json:"disks_count"`
}

// Sharing contains sharing subsystem results
type Sharing struct {
	SambaRunning      bool   `json:"samba_running"`
	SambaVersion      string `json:"samba_version,omitempty"`
	SharesCount       int    `json:"shares_count"`
	ConnectionsCount  int    `json:"connections_count"`
}

// Network contains network subsystem results
type Network struct {
	InterfacesCount int  `json:"interfaces_count"`
	FirewallActive  bool `json:"firewall_active"`
}

// Users contains user subsystem results
type Users struct {
	LocalUsersCount int `json:"local_users_count"`
	GroupsCount     int `json:"groups_count"`
}

// PrometheusMetrics contains Prometheus metrics check results
type PrometheusMetrics struct {
	EndpointReachable bool     `json:"endpoint_reachable"`
	MetricsFound      []string `json:"metrics_found"`
	MetricsMissing    []string `json:"metrics_missing"`
}

// SecurityTests contains security test results
type SecurityTests struct {
	AuthWithoutToken       string `json:"auth_without_token"`
	AuthWithInvalidToken   string `json:"auth_with_invalid_token"`
	SQLInjectionAttempt    string `json:"sql_injection_attempt"`
	XSSAttempt             string `json:"xss_attempt"`
}

// Performance contains performance metrics
type Performance struct {
	AverageResponseTimeMs int64            `json:"average_response_time_ms"`
	SlowestEndpoint       *EndpointTiming  `json:"slowest_endpoint,omitempty"`
	FastestEndpoint       *EndpointTiming  `json:"fastest_endpoint,omitempty"`
}

// EndpointTiming represents timing for a specific endpoint
type EndpointTiming struct {
	Path   string `json:"path"`
	TimeMs int64  `json:"time_ms"`
}

// Issue represents a problem found during audit
type Issue struct {
	Severity       string `json:"severity"` // HIGH, MEDIUM, LOW
	Type           string `json:"type"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
}
