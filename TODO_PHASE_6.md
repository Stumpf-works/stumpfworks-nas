# Phase 6: Enterprise Features - TODO Liste

**Version:** 1.2.0 (geplant)
**Status:** ⏳ 10% → Target: 100%
**Timeline:** 10-12 Wochen

---

## 🎯 Überblick

Phase 6 fügt Enterprise-Grade Features hinzu:
- ✅ **System Manager App** (Konsolidierung von Tasks/Backup/Monitoring)
- ⏳ **ACLs** (Access Control Lists)
- ⏳ **Disk Quotas** (User/Group Limits)
- ⏳ **High Availability** (Failover, DRBD, Cluster)

---

## Sprint 0: System Manager App Refactor (SOFORT)

**Ziel:** Dock aufräumen, Monitoring Dashboard verbessern
**Dauer:** 2-3 Tage
**Priorität:** HOCH 🔴

### Backend (keine Änderungen erforderlich)
- [ ] ✅ Monitoring API existiert bereits
- [ ] ✅ Tasks API existiert bereits
- [ ] ✅ Backup API existiert bereits

### Frontend

#### Neue App erstellen
- [ ] **Datei:** `/frontend/src/apps/SystemManager/SystemManager.tsx`
  - [ ] Tab-basierte UI (Monitoring, Tasks, Backups)
  - [ ] Header mit "System Management" Titel
  - [ ] Tab-Switching State Management

#### Monitoring Dashboard Tab
- [ ] **Datei:** `/frontend/src/apps/SystemManager/tabs/MonitoringDashboard.tsx`
  - [ ] Health Score Widget (groß und prominent)
  - [ ] System Metrics Cards (CPU, RAM, Disk, Network)
  - [ ] Live Charts (React-Chartjs-2 oder Recharts)
    - [ ] CPU Usage Chart (Line Chart, Last 24h)
    - [ ] Memory Usage Chart (Area Chart, Last 24h)
    - [ ] Disk I/O Chart (Bar Chart)
    - [ ] Network Traffic Chart (Line Chart)
  - [ ] Auto-Refresh (5s Interval)
  - [ ] Active Alerts Section
  - [ ] Link zu Prometheus Metrics (`/metrics`)

#### Tasks Tab
- [ ] **Datei:** `/frontend/src/apps/SystemManager/tabs/ScheduledTasks.tsx`
  - [ ] Import bestehende Tasks Component (`/apps/Tasks/Tasks.tsx`)
  - [ ] Minimal Wrapper für Integration

#### Backups Tab
- [ ] **Datei:** `/frontend/src/apps/SystemManager/tabs/BackupManager.tsx`
  - [ ] Import bestehende BackupManager Component (`/apps/BackupManager/BackupManager.tsx`)
  - [ ] Minimal Wrapper für Integration

#### App Registration
- [ ] **Datei:** `/frontend/src/apps/index.tsx`
  - [ ] Remove: `tasks` (id: 'tasks')
  - [ ] Remove: `backups` (id: 'backups')
  - [ ] Add: `system` (id: 'system', icon: '⚙️' oder '🖥️')
  - [ ] Change: `backups` icon von 💾 zu ⏱️ (falls noch verwendet)

### Testing
- [ ] Alle 3 Tabs öffnen und funktionsfähig testen
- [ ] Charts refreshen korrekt
- [ ] Tasks CRUD Operations
- [ ] Backup CRUD Operations
- [ ] Monitoring Alerts anzeigen

**Deliverable:** ✅ System Manager App mit Live-Monitoring-Dashboard

---

## Sprint 1: ACLs (Woche 1-2)

**Ziel:** Granulare Dateisystem-Berechtigungen
**Dauer:** 4-5 Tage
**Priorität:** HOCH 🔴

### Backend

#### Filesystem ACL Manager
- [ ] **Datei:** `/backend/internal/system/filesystem/acl.go`
  - [ ] Package erstellen: `package filesystem`
  - [ ] Struct: `ACLManager` mit ShellExecutor
  - [ ] Struct: `ACLEntry` (Type, Name, Permissions)
  - [ ] Method: `GetACL(path string) ([]ACLEntry, error)`
    - [ ] `getfacl` Command ausführen
    - [ ] Output parsen (user::rw-, user:alice:rwx, etc.)
    - [ ] ACLEntry Array zurückgeben
  - [ ] Method: `SetACL(path string, entries []ACLEntry) error`
    - [ ] `setfacl -m` Command bauen
    - [ ] Für jeden Entry: `-m type:name:permissions`
  - [ ] Method: `RemoveACL(path string, entryType, name string) error`
    - [ ] `setfacl -x type:name`
  - [ ] Method: `SetDefaultACL(dirPath string, entries []ACLEntry) error`
    - [ ] `setfacl -m default:type:name:permissions`
  - [ ] Method: `RemoveAllACLs(path string) error`
    - [ ] `setfacl -b`
  - [ ] Method: `ApplyRecursive(dirPath string, entries []ACLEntry) error`
    - [ ] `setfacl -R -m`

#### API Handler
- [ ] **Datei:** `/backend/internal/api/handlers/filesystem_acl.go`
  - [ ] Struct: `FilesystemACLHandler` mit ACLManager
  - [ ] Handler: `GetACL(w, r)` → GET /api/v1/filesystem/acl?path=...
  - [ ] Handler: `SetACL(w, r)` → POST /api/v1/filesystem/acl
  - [ ] Handler: `RemoveACL(w, r)` → DELETE /api/v1/filesystem/acl
  - [ ] Handler: `SetDefaultACL(w, r)` → POST /api/v1/filesystem/acl/default
  - [ ] Handler: `ApplyRecursive(w, r)` → POST /api/v1/filesystem/acl/recursive
  - [ ] Handler: `RemoveAllACLs(w, r)` → DELETE /api/v1/filesystem/acl/all

#### Router Registration
- [ ] **Datei:** `/backend/internal/api/router.go`
  - [ ] Route Group: `/api/v1/filesystem/acl`
  - [ ] GET `/api/v1/filesystem/acl`
  - [ ] POST `/api/v1/filesystem/acl`
  - [ ] DELETE `/api/v1/filesystem/acl`
  - [ ] POST `/api/v1/filesystem/acl/default`
  - [ ] POST `/api/v1/filesystem/acl/recursive`
  - [ ] DELETE `/api/v1/filesystem/acl/all`

#### Dependencies
- [ ] **Datei:** `/backend/internal/dependencies/checker.go`
  - [ ] Add: `acl` Package (`getfacl` command)

### Frontend

#### API Client
- [ ] **Datei:** `/frontend/src/api/filesystem-acl.ts`
  - [ ] Interface: `ACLEntry` (type, name, permissions)
  - [ ] Interface: `ACLInfo` (path, entries)
  - [ ] `filesystemACLApi.getACL(path)`
  - [ ] `filesystemACLApi.setACL(path, entries)`
  - [ ] `filesystemACLApi.removeACL(path, type, name)`
  - [ ] `filesystemACLApi.setDefaultACL(dirPath, entries)`
  - [ ] `filesystemACLApi.removeAllACLs(path)`
  - [ ] `filesystemACLApi.applyRecursive(dirPath, entries)`

#### UI Component
- [ ] **Datei:** `/frontend/src/apps/FileManager/components/ACLDialog.tsx`
  - [ ] Props: `filePath`, `onClose`
  - [ ] State: `entries[]`, `loading`, `error`
  - [ ] useEffect: Load ACLs via API
  - [ ] Component: `ACLEntryList` (Tabelle mit Entries)
  - [ ] Component: `AddACLEntry` (Formular: Type Dropdown, Name Input, Permissions)
  - [ ] Component: `DefaultACLSection` (nur für Directories)
  - [ ] Button: "Apply ACLs" → `setACL()`
  - [ ] Button: "Remove All ACLs" → `removeAllACLs()`
  - [ ] Checkbox: "Apply Recursively"
  - [ ] Permissions Editor (rwx Checkboxen oder Dropdown)

#### FileManager Integration
- [ ] **Datei:** `/frontend/src/apps/FileManager/FileManager.tsx`
  - [ ] Rechtsklick Context Menu erweitern
  - [ ] Menu Item: "Manage ACLs..." → `setShowACLDialog(true)`
  - [ ] ACLDialog einbinden `{showACLDialog && <ACLDialog ... />}`

### Testing
- [ ] Setze User ACL auf Datei
- [ ] Setze Group ACL auf Datei
- [ ] Setze Default ACL auf Verzeichnis
- [ ] Verifiziere neue Datei hat Default ACL
- [ ] Entferne einzelnen Entry
- [ ] Entferne alle ACLs
- [ ] Recursive Apply testen
- [ ] ACL mit SMB Share testen
- [ ] Error Handling (ungültige Permissions, nicht-existierende User)

**Deliverable:** ✅ Vollständige ACL-Verwaltung in FileManager

---

## Sprint 2: Disk Quotas (Woche 2-3)

**Ziel:** User/Group Speicherlimits
**Dauer:** 5-6 Tage
**Priorität:** HOCH 🔴

### Backend

#### Filesystem Quota Manager
- [ ] **Datei:** `/backend/internal/system/filesystem/quota.go`
  - [ ] Package erstellen: `package filesystem`
  - [ ] Type: `QuotaType` (UserQuota, GroupQuota)
  - [ ] Struct: `QuotaManager` mit ShellExecutor
  - [ ] Struct: `QuotaInfo` (Type, Name, Filesystem, BlocksUsed, BlocksSoft, BlocksHard, GracePeriod)
  - [ ] Method: `EnableQuotas(filesystem, quotaType) error`
    - [ ] Update /etc/fstab mit usrquota/grpquota
    - [ ] `mount -o remount`
    - [ ] `quotacheck -cug`
    - [ ] `quotaon`
  - [ ] Method: `DisableQuotas(filesystem, quotaType) error`
    - [ ] `quotaoff`
    - [ ] Update /etc/fstab (remove quota option)
  - [ ] Method: `GetQuota(filesystem, quotaType, name) (*QuotaInfo, error)`
    - [ ] `quota -u user` oder `quota -g group`
    - [ ] Parse Output
  - [ ] Method: `SetQuota(filesystem, quotaType, name, softLimit, hardLimit) error`
    - [ ] `setquota -u user soft hard 0 0 filesystem`
    - [ ] Convert bytes to blocks (1 block = 1KB)
  - [ ] Method: `RemoveQuota(filesystem, quotaType, name) error`
    - [ ] `setquota -u user 0 0 0 0 filesystem`
  - [ ] Method: `ListQuotas(filesystem, quotaType) ([]QuotaInfo, error)`
    - [ ] `repquota -u filesystem` oder `repquota -g`
    - [ ] Parse Output
  - [ ] Method: `GetQuotaReport(filesystem) (map[string][]QuotaInfo, error)`
    - [ ] Beide User und Group Quotas
    - [ ] Return map["user"] = [...], map["group"] = [...]

#### ZFS Quota erweitern
- [ ] **Datei:** `/backend/internal/system/storage/zfs.go`
  - [ ] Method: `SetDatasetQuota(dataset string, quota uint64) error`
    - [ ] `zfs set quota=XG dataset` (oder "none" für unlimitiert)
  - [ ] Method: `SetUserQuota(dataset, user string, quota uint64) error`
    - [ ] `zfs set userquota@user=XG dataset`
  - [ ] Method: `SetGroupQuota(dataset, group string, quota uint64) error`
    - [ ] `zfs set groupquota@group=XG dataset`
  - [ ] Method: `GetUserQuotaUsage(dataset, user string) (used, quota uint64, error)`
    - [ ] `zfs userspace -H -p dataset`
    - [ ] Parse für user
  - [ ] Method: `GetGroupQuotaUsage(dataset, group string) (used, quota uint64, error)`
    - [ ] `zfs groupspace -H -p dataset`

#### API Handler
- [ ] **Datei:** `/backend/internal/api/handlers/quota.go`
  - [ ] Struct: `QuotaHandler` mit QuotaManager + ZFSManager
  - [ ] Handler: `ListQuotas(w, r)` → GET /api/v1/quotas?filesystem=...&type=user
  - [ ] Handler: `GetQuota(w, r)` → GET /api/v1/quotas/:filesystem/:type/:name
  - [ ] Handler: `EnableQuotas(w, r)` → POST /api/v1/quotas/enable
  - [ ] Handler: `DisableQuotas(w, r)` → POST /api/v1/quotas/disable
  - [ ] Handler: `SetQuota(w, r)` → POST /api/v1/quotas
  - [ ] Handler: `RemoveQuota(w, r)` → DELETE /api/v1/quotas/:filesystem/:type/:name
  - [ ] Handler: `GetQuotaReport(w, r)` → GET /api/v1/quotas/report?filesystem=...

#### Router Registration
- [ ] **Datei:** `/backend/internal/api/router.go`
  - [ ] Route Group: `/api/v1/quotas`
  - [ ] GET `/api/v1/quotas`
  - [ ] GET `/api/v1/quotas/:filesystem/:type/:name`
  - [ ] POST `/api/v1/quotas/enable`
  - [ ] POST `/api/v1/quotas/disable`
  - [ ] POST `/api/v1/quotas`
  - [ ] DELETE `/api/v1/quotas/:filesystem/:type/:name`
  - [ ] GET `/api/v1/quotas/report`

#### Dependencies
- [ ] **Datei:** `/backend/internal/dependencies/checker.go`
  - [ ] Add: `quota` Package (`quota`, `setquota`, `repquota` commands)

### Frontend

#### API Client
- [ ] **Datei:** `/frontend/src/api/quota.ts`
  - [ ] Type: `QuotaType = 'user' | 'group'`
  - [ ] Interface: `QuotaInfo`
  - [ ] `quotaApi.enableQuotas(filesystem, type)`
  - [ ] `quotaApi.disableQuotas(filesystem, type)`
  - [ ] `quotaApi.getQuota(filesystem, type, name)`
  - [ ] `quotaApi.setQuota(filesystem, type, name, softLimit, hardLimit)`
  - [ ] `quotaApi.removeQuota(filesystem, type, name)`
  - [ ] `quotaApi.listQuotas(filesystem, type)`
  - [ ] `quotaApi.getQuotaReport(filesystem)`

#### UI Component - QuotaManager App
- [ ] **Datei:** `/frontend/src/apps/QuotaManager/QuotaManager.tsx`
  - [ ] State: `filesystem`, `quotaType`, `quotas[]`, `quotaEnabled`
  - [ ] Component: `FilesystemSelect` (Dropdown mit Filesystems)
  - [ ] Component: `QuotaTypeTabs` (User Quotas / Group Quotas)
  - [ ] Component: `QuotaTable`
    - [ ] Columns: Name, Used, Soft Limit, Hard Limit, Status, Grace Period, Actions
    - [ ] Progress Bar für Usage %
    - [ ] Status Icons (✅ OK, ⚠️ Warning, ❌ Over Limit)
    - [ ] Edit Button → Edit Dialog
    - [ ] Delete Button → Remove Quota
  - [ ] Component: `AddQuotaDialog`
    - [ ] Select User/Group
    - [ ] Input Soft Limit (GB)
    - [ ] Input Hard Limit (GB)
    - [ ] Button: "Add Quota"
  - [ ] Component: `QuotaStatusToggle`
    - [ ] Zeige: "Quotas Enabled ✅ [Disable]" oder "Quotas Disabled ❌ [Enable]"
  - [ ] Button: "Export Report" → CSV Download

#### App Registration
- [ ] **Datei:** `/frontend/src/apps/index.tsx`
  - [ ] Add: `quota-manager` (id: 'quota-manager', icon: '📊' oder '💿', name: 'Quotas')

#### UserManager Integration
- [ ] **Datei:** `/frontend/src/apps/UserManager/UserManager.tsx`
  - [ ] Rechtsklick auf User
  - [ ] Menu Item: "Set Disk Quota..." → Öffnet Quota Dialog für diesen User

### Testing
- [ ] Enable User Quotas auf ext4
- [ ] Enable Group Quotas auf ext4
- [ ] Set User Quota (alice: 45GB soft, 50GB hard)
- [ ] Überschreite Soft Limit → Warnung anzeigen
- [ ] Grace Period testen (7 Tage)
- [ ] Überschreite Hard Limit → Schreiben blockiert
- [ ] ZFS Dataset Quota setzen
- [ ] ZFS User Quota setzen
- [ ] ZFS Group Quota setzen
- [ ] Quota Report exportieren (CSV)
- [ ] Error Handling (Quota disabled, ungültige Limits)

**Deliverable:** ✅ Vollständige Quota-Verwaltung für ext4 und ZFS

---

## Sprint 3: High Availability - DRBD (Woche 4-6)

**Ziel:** Block-Level Replication
**Dauer:** 7 Tage
**Priorität:** MITTEL 🟡

### Backend

#### DRBD Manager
- [ ] **Datei:** `/backend/internal/system/ha/drbd.go`
  - [ ] Package erstellen: `package ha`
  - [ ] Struct: `DRBDManager` mit ShellExecutor
  - [ ] Struct: `DRBDResource` (Name, Device, Role, ConnectionState, DiskState, SyncProgress)
  - [ ] Method: `CreateResource(name, localDisk, remoteDisk, remoteIP, port) error`
    - [ ] Create /etc/drbd.d/resourcename.res
    - [ ] `drbdadm create-md resourcename`
    - [ ] `drbdadm up resourcename`
  - [ ] Method: `DeleteResource(name) error`
    - [ ] `drbdadm down resourcename`
    - [ ] Remove config file
  - [ ] Method: `GetResourceStatus(name) (*DRBDResource, error)`
    - [ ] `drbdadm status resourcename`
    - [ ] Parse Output
  - [ ] Method: `ListResources() ([]DRBDResource, error)`
    - [ ] `drbdadm status all`
  - [ ] Method: `PromoteToPrimary(name) error`
    - [ ] `drbdadm primary resourcename`
  - [ ] Method: `DemoteToSecondary(name) error`
    - [ ] `drbdadm secondary resourcename`
  - [ ] Method: `ForcePrimary(name) error`
    - [ ] `drbdadm primary --force resourcename` (Split-Brain)
  - [ ] Method: `Disconnect(name) error`
    - [ ] `drbdadm disconnect resourcename`
  - [ ] Method: `Connect(name) error`
    - [ ] `drbdadm connect resourcename`
  - [ ] Method: `StartSync(name) error`
    - [ ] `drbdadm invalidate resourcename` (force resync)

#### API Handler
- [ ] **Datei:** `/backend/internal/api/handlers/ha_drbd.go`
  - [ ] Struct: `HADRBDHandler` mit DRBDManager
  - [ ] Handler: `ListResources(w, r)` → GET /api/v1/ha/drbd/resources
  - [ ] Handler: `GetResource(w, r)` → GET /api/v1/ha/drbd/resources/:name
  - [ ] Handler: `CreateResource(w, r)` → POST /api/v1/ha/drbd/resources
  - [ ] Handler: `DeleteResource(w, r)` → DELETE /api/v1/ha/drbd/resources/:name
  - [ ] Handler: `Promote(w, r)` → POST /api/v1/ha/drbd/resources/:name/promote
  - [ ] Handler: `Demote(w, r)` → POST /api/v1/ha/drbd/resources/:name/demote
  - [ ] Handler: `ForcePrimary(w, r)` → POST /api/v1/ha/drbd/resources/:name/force-primary
  - [ ] Handler: `Disconnect(w, r)` → POST /api/v1/ha/drbd/resources/:name/disconnect
  - [ ] Handler: `Connect(w, r)` → POST /api/v1/ha/drbd/resources/:name/connect

#### Router Registration
- [ ] **Datei:** `/backend/internal/api/router.go`
  - [ ] Route Group: `/api/v1/ha/drbd`
  - [ ] Alle DRBD Routes registrieren

#### Dependencies
- [ ] **Datei:** `/backend/internal/dependencies/checker.go`
  - [ ] Add: `drbd-utils` Package (`drbdadm` command)

### Frontend

#### API Client (HA)
- [ ] **Datei:** `/frontend/src/api/ha.ts`
  - [ ] Interface: `DRBDResource`
  - [ ] `haApi.listDRBDResources()`
  - [ ] `haApi.getDRBDResource(name)`
  - [ ] `haApi.createDRBDResource(data)`
  - [ ] `haApi.deleteDRBDResource(name)`
  - [ ] `haApi.promoteDRBD(name)`
  - [ ] `haApi.demoteDRBD(name)`
  - [ ] `haApi.forcePrimaryDRBD(name)`
  - [ ] `haApi.disconnectDRBD(name)`
  - [ ] `haApi.connectDRBD(name)`

#### UI Component - HA App
- [ ] **Datei:** `/frontend/src/apps/HighAvailability/HighAvailability.tsx`
  - [ ] Tab-based UI: Cluster, DRBD, VIP
  - [ ] State: `activeTab`

- [ ] **Datei:** `/frontend/src/apps/HighAvailability/tabs/DRBDPanel.tsx`
  - [ ] Component: `DRBDResourceList`
    - [ ] Columns: Name, Device, Role, Connection State, Disk State, Sync Progress
    - [ ] Progress Bar für Sync%
    - [ ] Actions: Promote, Demote, Disconnect, Connect, Delete
  - [ ] Component: `CreateDRBDDialog`
    - [ ] Input: Resource Name
    - [ ] Input: Local Disk (/dev/sda1)
    - [ ] Input: Remote IP
    - [ ] Input: Remote Disk (/dev/sda1)
    - [ ] Input: Port (default: 7789)
    - [ ] Button: "Create Resource"
  - [ ] Button: "Force Primary" (mit Warning)

#### App Registration
- [ ] **Datei:** `/frontend/src/apps/index.tsx`
  - [ ] Add: `high-availability` (id: 'ha', icon: '🔄', name: 'High Availability')

### Testing (2-Node Setup erforderlich!)
- [ ] Create DRBD Resource
- [ ] Initial Sync durchführen (100%)
- [ ] Promote Node A zu Primary
- [ ] Schreibe Daten auf Node A
- [ ] Verifiziere Sync zu Node B
- [ ] Demote Node A zu Secondary
- [ ] Promote Node B zu Primary
- [ ] Simuliere Disconnect → Reconnect
- [ ] Split-Brain Szenario → Force Primary
- [ ] Delete DRBD Resource

**Deliverable:** ✅ DRBD Block Replication vollständig funktional

---

## Sprint 4: High Availability - Cluster (Woche 7-9)

**Ziel:** Pacemaker/Corosync Cluster Management
**Dauer:** 9 Tage
**Priorität:** MITTEL 🟡

### Backend

#### Cluster Manager
- [ ] **Datei:** `/backend/internal/system/ha/cluster.go`
  - [ ] Struct: `ClusterManager` mit ShellExecutor
  - [ ] Struct: `ClusterNode` (Name, IP, Status, Resources)
  - [ ] Struct: `ClusterResource` (ID, Type, Node, Status)
  - [ ] Method: `InitializeCluster(clusterName, nodes) error`
    - [ ] `pcs cluster setup --name clusterName node1 node2`
    - [ ] `pcs cluster start --all`
    - [ ] `pcs cluster enable --all`
  - [ ] Method: `AddNode(nodeName, nodeIP) error`
    - [ ] `pcs cluster node add nodeName`
  - [ ] Method: `RemoveNode(nodeName) error`
    - [ ] `pcs cluster node remove nodeName`
  - [ ] Method: `GetClusterStatus() ([]ClusterNode, error)`
    - [ ] `pcs status xml` (oder `crm_mon -X`)
    - [ ] Parse XML Output
  - [ ] Method: `StandbyNode(nodeName) error`
    - [ ] `pcs node standby nodeName`
  - [ ] Method: `OnlineNode(nodeName) error`
    - [ ] `pcs node unstandby nodeName`
  - [ ] Method: `AddVirtualIP(ip, netmask) error`
    - [ ] `pcs resource create vip IPaddr2 ip=X.X.X.X cidr_netmask=24`
  - [ ] Method: `AddFilesystem(device, mountpoint, fstype) error`
    - [ ] `pcs resource create fs Filesystem device=/dev/drbd0 directory=/data fstype=ext4`
  - [ ] Method: `AddService(serviceName) error`
    - [ ] `pcs resource create service systemd:stumpfworks-nas`
  - [ ] Method: `MoveResource(resourceID, targetNode) error`
    - [ ] `pcs resource move resourceID targetNode`
  - [ ] Method: `StopResource(resourceID) error`
    - [ ] `pcs resource disable resourceID`
  - [ ] Method: `StartResource(resourceID) error`
    - [ ] `pcs resource enable resourceID`

#### API Handler
- [ ] **Datei:** `/backend/internal/api/handlers/ha_cluster.go`
  - [ ] Struct: `HAClusterHandler` mit ClusterManager
  - [ ] Handler: `GetClusterStatus(w, r)` → GET /api/v1/ha/cluster/status
  - [ ] Handler: `InitializeCluster(w, r)` → POST /api/v1/ha/cluster/init
  - [ ] Handler: `AddNode(w, r)` → POST /api/v1/ha/cluster/nodes
  - [ ] Handler: `RemoveNode(w, r)` → DELETE /api/v1/ha/cluster/nodes/:name
  - [ ] Handler: `StandbyNode(w, r)` → POST /api/v1/ha/cluster/nodes/:name/standby
  - [ ] Handler: `OnlineNode(w, r)` → POST /api/v1/ha/cluster/nodes/:name/online
  - [ ] Handler: `AddVirtualIP(w, r)` → POST /api/v1/ha/cluster/resources/vip
  - [ ] Handler: `AddFilesystem(w, r)` → POST /api/v1/ha/cluster/resources/filesystem
  - [ ] Handler: `AddService(w, r)` → POST /api/v1/ha/cluster/resources/service
  - [ ] Handler: `MoveResource(w, r)` → POST /api/v1/ha/cluster/resources/:id/move
  - [ ] Handler: `StopResource(w, r)` → POST /api/v1/ha/cluster/resources/:id/stop
  - [ ] Handler: `StartResource(w, r)` → POST /api/v1/ha/cluster/resources/:id/start

#### Router Registration
- [ ] **Datei:** `/backend/internal/api/router.go`
  - [ ] Route Group: `/api/v1/ha/cluster`
  - [ ] Alle Cluster Routes registrieren

#### Dependencies
- [ ] **Datei:** `/backend/internal/dependencies/checker.go`
  - [ ] Add: `pacemaker` Package (`pcs` oder `crm` command)
  - [ ] Add: `corosync` Package (`corosync-cfgtool` command)

### Frontend

#### API Client Update
- [ ] **Datei:** `/frontend/src/api/ha.ts`
  - [ ] Interface: `ClusterNode`
  - [ ] Interface: `ClusterResource`
  - [ ] `haApi.getClusterStatus()`
  - [ ] `haApi.initializeCluster(clusterName, nodes)`
  - [ ] `haApi.addNode(name, ip)`
  - [ ] `haApi.removeNode(name)`
  - [ ] `haApi.standbyNode(name)`
  - [ ] `haApi.onlineNode(name)`
  - [ ] `haApi.addVirtualIP(ip, netmask)`
  - [ ] `haApi.addFilesystem(device, mountpoint, fstype)`
  - [ ] `haApi.addService(serviceName)`
  - [ ] `haApi.moveResource(resourceId, targetNode)`
  - [ ] `haApi.stopResource(resourceId)`
  - [ ] `haApi.startResource(resourceId)`

#### UI Component
- [ ] **Datei:** `/frontend/src/apps/HighAvailability/tabs/ClusterStatus.tsx`
  - [ ] Component: `ClusterOverview`
    - [ ] Cluster Name
    - [ ] Status: Healthy/Degraded/Failed
    - [ ] Node Count (X/Y online)
  - [ ] Component: `NodeList`
    - [ ] Columns: Node Name, IP, Status, Resources (count), Actions
    - [ ] Status Icons: ✅ Online, ⏸️ Standby, ❌ Offline
    - [ ] Actions: Standby, Online, Remove
  - [ ] Component: `ResourceList`
    - [ ] Columns: Resource ID, Type, Current Node, Status, Actions
    - [ ] Type Icons: 🌐 VIP, 💾 Filesystem, ⚙️ Service
    - [ ] Actions: Move, Stop, Start, Delete
  - [ ] Component: `InitClusterDialog`
    - [ ] Input: Cluster Name
    - [ ] Multi-Input: Nodes (Name + IP)
    - [ ] Button: "Initialize Cluster"
  - [ ] Component: `AddNodeDialog`
    - [ ] Input: Node Name
    - [ ] Input: Node IP
    - [ ] Button: "Add Node"
  - [ ] Component: `AddResourceDialog`
    - [ ] Tabs: VIP, Filesystem, Service
    - [ ] VIP: IP + Netmask
    - [ ] Filesystem: Device + Mountpoint + FSType
    - [ ] Service: Service Name
    - [ ] Button: "Add Resource"
  - [ ] Button: "Trigger Failover" → Move all resources zu anderem Node

### Testing (2-Node Cluster erforderlich!)
- [ ] Initialize 2-Node Cluster
- [ ] Add Virtual IP (10.0.0.10)
- [ ] Add Filesystem (/dev/drbd0 → /data)
- [ ] Add Service (stumpfworks-nas)
- [ ] Verify alle Resources auf Node A
- [ ] Trigger Manual Failover → Node B
- [ ] Verify VIP gewandert
- [ ] Verify Filesystem umount/mount
- [ ] Verify Service gestoppt/gestartet
- [ ] Simulate Node A Crash → Automatic Failover
- [ ] Node A Recovery → Failback
- [ ] Remove Node
- [ ] Error Handling (Split-Brain, Quorum Loss)

**Deliverable:** ✅ Pacemaker/Corosync Cluster vollständig funktional

---

## Sprint 5: High Availability - VIP (Woche 10)

**Ziel:** Keepalived Virtual IP Management
**Dauer:** 4 Tage
**Priorität:** MITTEL 🟡

### Backend

#### Keepalived Manager
- [ ] **Datei:** `/backend/internal/system/ha/keepalived.go`
  - [ ] Struct: `KeepaliveManager` mit ShellExecutor
  - [ ] Struct: `VirtualIP` (IP, Interface, State, Priority, VirtualRouter)
  - [ ] Method: `ConfigureVIP(vip, iface, priority) error`
    - [ ] Create /etc/keepalived/keepalived.conf
    - [ ] vrrp_instance definition
    - [ ] `systemctl restart keepalived`
  - [ ] Method: `RemoveVIP(vip) error`
    - [ ] Remove vrrp_instance from config
    - [ ] `systemctl restart keepalived`
  - [ ] Method: `GetVIPStatus(vip) (*VirtualIP, error)`
    - [ ] Parse /var/run/keepalived.data (oder `ip addr show`)
    - [ ] Determine State (MASTER/BACKUP)
  - [ ] Method: `ListVIPs() ([]VirtualIP, error)`
    - [ ] Parse keepalived.conf
    - [ ] Get Status für jede VIP
  - [ ] Method: `SetPriority(vip, priority) error`
    - [ ] Update config
    - [ ] Restart keepalived

#### API Handler
- [ ] **Datei:** `/backend/internal/api/handlers/ha_vip.go`
  - [ ] Struct: `HAVIPHandler` mit KeepaliveManager
  - [ ] Handler: `ListVIPs(w, r)` → GET /api/v1/ha/vips
  - [ ] Handler: `GetVIP(w, r)` → GET /api/v1/ha/vips/:ip
  - [ ] Handler: `AddVIP(w, r)` → POST /api/v1/ha/vips
  - [ ] Handler: `RemoveVIP(w, r)` → DELETE /api/v1/ha/vips/:ip
  - [ ] Handler: `SetPriority(w, r)` → PUT /api/v1/ha/vips/:ip/priority

#### Router Registration
- [ ] **Datei:** `/backend/internal/api/router.go`
  - [ ] Route Group: `/api/v1/ha/vips`
  - [ ] Alle VIP Routes registrieren

#### Dependencies
- [ ] **Datei:** `/backend/internal/dependencies/checker.go`
  - [ ] Add: `keepalived` Package (`keepalived` command)

### Frontend

#### API Client Update
- [ ] **Datei:** `/frontend/src/api/ha.ts`
  - [ ] Interface: `VirtualIP`
  - [ ] `haApi.listVIPs()`
  - [ ] `haApi.getVIP(ip)`
  - [ ] `haApi.addVIP(ip, iface, priority)`
  - [ ] `haApi.removeVIP(ip)`
  - [ ] `haApi.setVIPPriority(ip, priority)`

#### UI Component
- [ ] **Datei:** `/frontend/src/apps/HighAvailability/tabs/VIPManager.tsx`
  - [ ] Component: `VIPList`
    - [ ] Columns: IP, Interface, State (MASTER/BACKUP), Priority, Actions
    - [ ] State Icons: 👑 MASTER, 🔄 BACKUP
    - [ ] Actions: Edit Priority, Remove
  - [ ] Component: `AddVIPDialog`
    - [ ] Input: Virtual IP (e.g., 10.0.0.10)
    - [ ] Select: Interface (eth0, eth1, etc.)
    - [ ] Input: Priority (1-255, default: 100)
    - [ ] Info: Higher priority = MASTER
    - [ ] Button: "Add VIP"
  - [ ] Component: `EditPriorityDialog`
    - [ ] Input: New Priority
    - [ ] Warning: "Changing priority may trigger failover"
    - [ ] Button: "Update Priority"

### Testing (2-Node Setup erforderlich!)
- [ ] Configure VIP 10.0.0.10 on Node A (Priority 100)
- [ ] Configure VIP 10.0.0.10 on Node B (Priority 50)
- [ ] Verify Node A is MASTER
- [ ] Ping 10.0.0.10 (sollte funktionieren)
- [ ] Stop keepalived on Node A → VIP zu Node B
- [ ] Verify Node B is now MASTER
- [ ] Start keepalived on Node A → Failback
- [ ] Change Priority (Node B = 150) → VIP zu Node B
- [ ] Error Handling (Invalid IP, Interface doesn't exist)

**Deliverable:** ✅ Keepalived VIP Management vollständig funktional

---

## Sprint 6: Integration & Dokumentation (Woche 11)

**Ziel:** Testing, Docs, Release vorbereiten
**Dauer:** 5 Tage
**Priorität:** HOCH 🔴

### Integration Testing
- [ ] **Tag 1-2: End-to-End Testing**
  - [ ] ACLs + Quotas zusammen testen
    - [ ] User mit Quota + ACL auf /data
    - [ ] Quota enforcement + ACL permissions
  - [ ] HA Failover Szenarien
    - [ ] DRBD Sync + Cluster Failover + VIP Migration
    - [ ] Simulate Node Crash → Full Recovery
  - [ ] Performance Testing
    - [ ] DRBD Sync Performance (MB/s)
    - [ ] Quota Overhead Measurement
    - [ ] ACL Performance Impact
  - [ ] Load Testing
    - [ ] Multiple concurrent Quota checks
    - [ ] Multiple concurrent ACL operations

### Dokumentation
- [ ] **Tag 3-4: Documentation**
  - [ ] User Guides
    - [ ] `/docs/ACL_GUIDE.md`
      - [ ] Was sind ACLs?
      - [ ] Use Cases & Examples
      - [ ] Best Practices
      - [ ] Troubleshooting
    - [ ] `/docs/QUOTA_GUIDE.md`
      - [ ] Was sind Quotas?
      - [ ] User vs Group Quotas
      - [ ] ZFS vs ext4 Quotas
      - [ ] Grace Periods
      - [ ] Troubleshooting
    - [ ] `/docs/HA_GUIDE.md`
      - [ ] HA Konzepte (Failover, Replication, VIP)
      - [ ] DRBD Setup Tutorial
      - [ ] Cluster Setup Tutorial
      - [ ] VIP Setup Tutorial
      - [ ] Failover Testing
      - [ ] Split-Brain Recovery
      - [ ] Production Deployment Checklist
  - [ ] Admin Guides
    - [ ] `/docs/ACL_ADMIN.md`
      - [ ] CLI Commands (getfacl, setfacl)
      - [ ] Batch Operations
      - [ ] ACL Migration Scripts
    - [ ] `/docs/QUOTA_ADMIN.md`
      - [ ] CLI Commands (quota, setquota, repquota)
      - [ ] Enable/Disable Quotas
      - [ ] Quota Reports
    - [ ] `/docs/HA_ADMIN.md`
      - [ ] DRBD Administration (drbdadm)
      - [ ] Cluster Administration (pcs, crm)
      - [ ] Node Maintenance Procedures
      - [ ] Emergency Procedures
      - [ ] Monitoring & Alerting Setup
  - [ ] API Documentation
    - [ ] Update `/docs/API.md`
      - [ ] ACL Endpoints (6 endpoints)
      - [ ] Quota Endpoints (7 endpoints)
      - [ ] HA DRBD Endpoints (9 endpoints)
      - [ ] HA Cluster Endpoints (13 endpoints)
      - [ ] HA VIP Endpoints (5 endpoints)
  - [ ] README Updates
    - [ ] Update `/README.md`
      - [ ] Phase 6: 100% ✅
      - [ ] Enterprise Features Section
      - [ ] HA Support Badge

### Finalisierung
- [ ] **Tag 5: Release Preparation**
  - [ ] Bug Fixes
    - [ ] Review all TODOs from testing
    - [ ] Fix critical bugs
    - [ ] Fix UI/UX issues
  - [ ] UI Polish
    - [ ] Consistent styling
    - [ ] Loading states
    - [ ] Error messages
    - [ ] Tooltips & Help text
  - [ ] CHANGELOG Update
    - [ ] Add v1.2.0 entry
    - [ ] List all new features
    - [ ] Breaking changes (if any)
    - [ ] Migration guide
  - [ ] Version Bump
    - [ ] Update version to 1.2.0
    - [ ] Update dependency versions
  - [ ] Build & Test
    - [ ] Build Debian package
    - [ ] Test installation on fresh system
    - [ ] Test upgrade from v1.1.0

**Deliverable:** ✅ v1.2.0 Production Ready

---

## Zusammenfassung

### Gesamte Phase 6 Timeline

| Sprint | Dauer | Features | Status |
|--------|-------|----------|--------|
| **Sprint 0** | 2-3 Tage | System Manager App | ⏳ SOFORT |
| **Sprint 1** | 4-5 Tage | ACLs | ⏳ Woche 1-2 |
| **Sprint 2** | 5-6 Tage | Disk Quotas | ⏳ Woche 2-3 |
| **Sprint 3** | 7 Tage | HA - DRBD | ⏳ Woche 4-6 |
| **Sprint 4** | 9 Tage | HA - Cluster | ⏳ Woche 7-9 |
| **Sprint 5** | 4 Tage | HA - VIP | ⏳ Woche 10 |
| **Sprint 6** | 5 Tage | Integration & Docs | ⏳ Woche 11 |

**Gesamt:** ~37.5 Tage = **10-12 Wochen**

### Deliverables v1.2.0

✅ **System Manager App** - Monitoring Dashboard mit Live-Charts
✅ **Filesystem ACLs** - Granulare Berechtigungen (setfacl/getfacl)
✅ **Disk Quotas** - User/Group Limits (ext4 + ZFS)
✅ **DRBD** - Block-Level Replication
✅ **Pacemaker/Corosync** - Cluster Management
✅ **Keepalived** - Virtual IP Management
✅ **Comprehensive Documentation** - User & Admin Guides

### Success Metrics

- [ ] README.md: Phase 6: 100% ✅
- [ ] All Enterprise Features funktional
- [ ] 2-Node HA Setup erfolgreich getestet
- [ ] Complete User & Admin Documentation
- [ ] v1.2.0 Release auf APT Repository deployed

---

**Next Action:** 🚀 Start Sprint 0 (System Manager App) - SOFORT!
