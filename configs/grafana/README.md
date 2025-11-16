# StumpfWorks NAS Grafana Dashboards

This directory contains pre-configured Grafana dashboard templates for monitoring your StumpfWorks NAS system.

## Available Dashboards

### 1. System Overview Dashboard
**File:** `system-overview-dashboard.json`

Provides a comprehensive view of system resources:
- CPU Usage (%)
- Memory Usage (%)
- Memory Breakdown (gauge)
- Disk I/O (read/write rates)
- Network Traffic (RX/TX)
- System Uptime
- Load Average (1m, 5m, 15m)

**Refresh Rate:** 30 seconds

### 2. ZFS Pools Dashboard
**File:** `zfs-pools-dashboard.json`

Monitors ZFS storage pools in detail:
- Pool Health Status (ONLINE/DEGRADED/OFFLINE)
- Pool Capacity (Total, Used, Free)
- Pool Usage Percentage (gauge)
- Fragmentation levels
- Deduplication ratios

**Refresh Rate:** 1 minute

### 3. Services & Shares Dashboard
**File:** `services-dashboard.json`

Tracks service status and share connections:
- Service Status (running/stopped)
- Samba Connections
- NFS Connections
- Total Active Connections
- Service Uptime
- Connection Rate (5m average)

**Refresh Rate:** 30 seconds

### 4. Disk Health (SMART) Dashboard
**File:** `disk-health-dashboard.json`

Monitors disk health via SMART metrics:
- Disk Health Status (healthy/unhealthy)
- Disk Temperatures (time series)
- Current Temperature by Disk (gauges)
- Disk Health Summary (table)
- Unhealthy Disk Alerts

**Refresh Rate:** 5 minutes

## Installation

### Import via Grafana UI

1. Log in to your Grafana instance
2. Navigate to **Dashboards** → **Import**
3. Click **Upload JSON file**
4. Select one of the dashboard JSON files
5. Choose your Prometheus data source
6. Click **Import**

### Import via Grafana API

```bash
# Set your Grafana credentials
GRAFANA_URL="http://localhost:3000"
GRAFANA_API_KEY="your-api-key-here"

# Import all dashboards
for dashboard in *.json; do
  curl -X POST \
    -H "Authorization: Bearer $GRAFANA_API_KEY" \
    -H "Content-Type: application/json" \
    -d @"$dashboard" \
    "$GRAFANA_URL/api/dashboards/db"
done
```

### Import via Provisioning

1. Copy dashboard files to Grafana provisioning directory:
   ```bash
   sudo cp *.json /etc/grafana/provisioning/dashboards/
   ```

2. Create a provisioning configuration file `/etc/grafana/provisioning/dashboards/stumpfworks.yaml`:
   ```yaml
   apiVersion: 1

   providers:
     - name: 'StumpfWorks NAS'
       orgId: 1
       folder: 'StumpfWorks NAS'
       type: file
       disableDeletion: false
       updateIntervalSeconds: 10
       allowUiUpdates: true
       options:
         path: /etc/grafana/provisioning/dashboards
         foldersFromFilesStructure: true
   ```

3. Restart Grafana:
   ```bash
   sudo systemctl restart grafana-server
   ```

## Prerequisites

### Prometheus Configuration

Ensure your Prometheus instance is configured to scrape metrics from StumpfWorks NAS:

```yaml
scrape_configs:
  - job_name: 'stumpfworks-nas'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

### Required Metrics

The dashboards expect the following metrics to be exported from `/metrics` endpoint:

**System Metrics:**
- `stumpfworks_cpu_usage_percent`
- `stumpfworks_memory_total_bytes`
- `stumpfworks_memory_used_bytes`
- `stumpfworks_disk_read_bytes`
- `stumpfworks_disk_write_bytes`
- `stumpfworks_network_rx_bytes`
- `stumpfworks_network_tx_bytes`
- `stumpfworks_uptime_seconds`
- `stumpfworks_load_average_1m`
- `stumpfworks_load_average_5m`
- `stumpfworks_load_average_15m`

**ZFS Metrics:**
- `stumpfworks_zfs_pool_health` (labels: pool, health)
- `stumpfworks_zfs_pool_total_bytes` (label: pool)
- `stumpfworks_zfs_pool_used_bytes` (label: pool)
- `stumpfworks_zfs_pool_free_bytes` (label: pool)
- `stumpfworks_zfs_pool_fragmentation_percent` (label: pool)
- `stumpfworks_zfs_pool_dedup_ratio` (label: pool)

**Service Metrics:**
- `stumpfworks_service_status` (label: service)
- `stumpfworks_share_samba_connections`
- `stumpfworks_share_nfs_connections`

**SMART Metrics:**
- `stumpfworks_disk_smart_healthy` (label: device)
- `stumpfworks_disk_smart_temperature` (label: device)

## Customization

All dashboards can be customized after import:
- Adjust refresh rates
- Modify thresholds and alerts
- Add custom panels
- Change color schemes
- Configure notification channels

## Alerting

To set up alerts:

1. Configure a notification channel in Grafana (**Alerting** → **Notification channels**)
2. Edit dashboard panels to add alert rules
3. Set thresholds (e.g., CPU > 90%, Disk health = 0)
4. Assign notification channels

Example alert conditions:
- CPU usage > 90% for 5 minutes
- Memory usage > 95% for 5 minutes
- ZFS pool health ≠ ONLINE
- Disk temperature > 55°C
- Disk health = unhealthy

## Support

For issues or feature requests, please visit:
- GitHub: https://github.com/Stumpf-works/stumpfworks-nas
- Documentation: https://github.com/Stumpf-works/stumpfworks-nas/wiki

## License

Copyright © 2025 Stumpf.Works. All rights reserved.
