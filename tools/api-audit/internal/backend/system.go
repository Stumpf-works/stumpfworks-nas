// Revision: 2025-11-16 | Author: Claude | Version: 1.0.0
package backend

import (
	"encoding/json"

	"api-audit/internal/client"
	"api-audit/internal/report"
)

// SystemTester tests backend system functions
type SystemTester struct {
	client *client.Client
}

// NewSystemTester creates a new system tester
func NewSystemTester(c *client.Client) *SystemTester {
	return &SystemTester{client: c}
}

// TestBackendFunctions tests all backend functions
func (t *SystemTester) TestBackendFunctions() report.BackendFunctions {
	return report.BackendFunctions{
		SystemLibrary: report.SystemLibrary{
			Storage: t.testStorage(),
			Sharing: t.testSharing(),
			Network: t.testNetwork(),
			Users:   t.testUsers(),
		},
	}
}

func (t *SystemTester) testStorage() report.Storage {
	storage := report.Storage{}

	// Test ZFS
	resp, err := t.client.Get("/api/v1/syslib/zfs/pools")
	if err == nil && resp.StatusCode == 200 {
		storage.ZFSAvailable = true
		var data struct {
			Data []interface{} `json:"data"`
		}
		if json.Unmarshal(resp.Body, &data) == nil {
			storage.PoolsCount = len(data.Data)
		}
	}

	// Test SMART
	resp, err = t.client.Get("/api/v1/syslib/smart/disks")
	if err == nil && resp.StatusCode == 200 {
		storage.SmartAvailable = true
		var data struct {
			Data []interface{} `json:"data"`
		}
		if json.Unmarshal(resp.Body, &data) == nil {
			storage.DisksCount = len(data.Data)
		}
	}

	return storage
}

func (t *SystemTester) testSharing() report.Sharing {
	sharing := report.Sharing{}

	// Test Samba status
	resp, err := t.client.Get("/api/v1/syslib/samba/status")
	if err == nil && resp.StatusCode == 200 {
		sharing.SambaRunning = true
		var data struct {
			Data struct {
				Version string `json:"version"`
			} `json:"data"`
		}
		if json.Unmarshal(resp.Body, &data) == nil {
			sharing.SambaVersion = data.Data.Version
		}
	}

	// Test Samba shares
	resp, err = t.client.Get("/api/v1/syslib/samba/shares")
	if err == nil && resp.StatusCode == 200 {
		var data struct {
			Data []interface{} `json:"data"`
		}
		if json.Unmarshal(resp.Body, &data) == nil {
			sharing.SharesCount = len(data.Data)
		}
	}

	return sharing
}

func (t *SystemTester) testNetwork() report.Network {
	network := report.Network{}

	// Test network interfaces
	resp, err := t.client.Get("/api/v1/syslib/network/interfaces")
	if err == nil && resp.StatusCode == 200 {
		var data struct {
			Data []interface{} `json:"data"`
		}
		if json.Unmarshal(resp.Body, &data) == nil {
			network.InterfacesCount = len(data.Data)
		}
	}

	// Test firewall
	resp, err = t.client.Get("/api/v1/syslib/network/firewall")
	if err == nil && resp.StatusCode == 200 {
		var data struct {
			Data struct {
				Active bool `json:"active"`
			} `json:"data"`
		}
		if json.Unmarshal(resp.Body, &data) == nil {
			network.FirewallActive = data.Data.Active
		}
	}

	return network
}

func (t *SystemTester) testUsers() report.Users {
	users := report.Users{}

	// Test users list
	resp, err := t.client.Get("/api/v1/users")
	if err == nil && resp.StatusCode == 200 {
		var data struct {
			Data []interface{} `json:"data"`
		}
		if json.Unmarshal(resp.Body, &data) == nil {
			users.LocalUsersCount = len(data.Data)
		}
	}

	// Test groups list
	resp, err = t.client.Get("/api/v1/groups")
	if err == nil && resp.StatusCode == 200 {
		var data struct {
			Data []interface{} `json:"data"`
		}
		if json.Unmarshal(resp.Body, &data) == nil {
			users.GroupsCount = len(data.Data)
		}
	}

	return users
}

// TestPrometheusMetrics tests Prometheus metrics endpoint
func (t *SystemTester) TestPrometheusMetrics() report.PrometheusMetrics {
	metrics := report.PrometheusMetrics{
		MetricsFound:   []string{},
		MetricsMissing: []string{},
	}

	resp, err := t.client.Get("/metrics")
	if err != nil || resp.StatusCode != 200 {
		metrics.EndpointReachable = false
		return metrics
	}

	metrics.EndpointReachable = true

	// Expected metrics
	expected := []string{
		"stumpfworks_cpu_usage_percent",
		"stumpfworks_memory_total_bytes",
		"stumpfworks_memory_used_bytes",
		"stumpfworks_zfs_pool_health",
		"stumpfworks_disk_smart_healthy",
		"stumpfworks_service_status",
		"stumpfworks_share_samba_connections",
		"stumpfworks_share_nfs_connections",
	}

	body := string(resp.Body)
	for _, metric := range expected {
		found := false
		for i := 0; i <= len(body)-len(metric); i++ {
			if body[i:i+len(metric)] == metric {
				found = true
				break
			}
		}

		if found {
			metrics.MetricsFound = append(metrics.MetricsFound, metric)
		} else {
			metrics.MetricsMissing = append(metrics.MetricsMissing, metric)
		}
	}

	return metrics
}
