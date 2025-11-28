# Phase 6: Enterprise Features - Audit & Implementation Plan

**Datum:** 2025-11-27
**Version:** 1.2.0 (geplant)
**Status:** ⏳ 10% (nur partielle iSCSI ACLs vorhanden)
**Branch:** `claude/audit-frontend-integration-017ykh14kLotgaQ7URfFanDX`

---

## Executive Summary

Phase 6 bringt **Enterprise-grade Features** in Stumpf.Works NAS:
- **ACLs (Access Control Lists)**: Granulare Dateisystem-Berechtigungen
- **Disk Quotas**: Benutzer- und Gruppen-basierte Speicherlimits
- **High Availability**: Failover, DRBD, Cluster-Support

**Aktueller Status:**
- ✅ iSCSI ACLs (CHAP Authentication) - IMPLEMENTIERT
- ✅ ZFS Quota Felder - DEFINIERT, aber nicht verwaltet
- ❌ Filesystem ACLs (setfacl/getfacl) - KOMPLETT FEHLEND
- ❌ User/Group Disk Quotas - KOMPLETT FEHLEND
- ❌ High Availability - KOMPLETT FEHLEND

---

## 1. Access Control Lists (ACLs)

### 1.1 Was sind ACLs?

ACLs erweitern das traditionelle Unix-Berechtigungsmodell (rwx für owner/group/others) um:
- **Mehrere Benutzer** mit unterschiedlichen Rechten auf derselben Datei
- **Mehrere Gruppen** mit unterschiedlichen Rechten
- **Default ACLs** für neu erstellte Dateien in einem Verzeichnis
- **Vererbung** von Berechtigungen

**Use Case:**
```bash
# Standard: Nur 1 Owner, 1 Group
-rw-r--r-- alice:developers file.txt

# Mit ACL: Mehrere Benutzer/Gruppen
-rw-r--r--+ alice:developers file.txt
  user:bob:rw-      # Bob kann lesen/schreiben
  user:charlie:r--  # Charlie nur lesen
  group:sales:r--   # Sales-Gruppe nur lesen
```

### 1.2 Was bereits existiert: iSCSI ACLs ✅

**Backend:** `/backend/internal/system/sharing/iscsi.go`

**Implementierte Funktionen:**
```go
// iSCSI Target ACLs (Initiator-basierte Zugriffskontrolle)
func (i *ISCSIManager) AddACL(targetIQN, initiatorIQN string) error
func (i *ISCSIManager) RemoveACL(targetIQN, initiatorIQN string) error
func (i *ISCSIManager) SetCHAPAuth(targetIQN, initiatorIQN, username, password string) error
func (i *ISCSIManager) DisableCHAPAuth(targetIQN, initiatorIQN string) error
```

**Frontend:** Bereits in `StorageManager` → iSCSI Tab integriert ✅

**Limitation:** Dies sind **iSCSI-spezifische ACLs**, NICHT Filesystem ACLs!

### 1.3 Was fehlt: Filesystem ACLs ❌

#### Backend-Implementation erforderlich

**Neue Datei:** `/backend/internal/system/filesystem/acl.go`

```go
package filesystem

import "os/exec"

type ACLManager struct {
    shell *executor.ShellExecutor
}

// ACL Entry
type ACLEntry struct {
    Type       string   `json:"type"`        // user, group, mask, other
    Name       string   `json:"name"`        // username or groupname
    Permissions string   `json:"permissions"` // rwx, r-x, etc.
}

// ACL Operations
func (a *ACLManager) GetACL(path string) ([]ACLEntry, error)
func (a *ACLManager) SetACL(path string, entries []ACLEntry) error
func (a *ACLManager) RemoveACL(path string, entryType, name string) error
func (a *ACLManager) SetDefaultACL(dirPath string, entries []ACLEntry) error
func (a *ACLManager) RemoveAllACLs(path string) error

// Bulk Operations
func (a *ACLManager) CopyACL(sourcePath, destPath string) error
func (a *ACLManager) ApplyRecursive(dirPath string, entries []ACLEntry) error
```

**Implementation Details:**

```go
// GetACL - Liest ACLs einer Datei/Verzeichnis
func (a *ACLManager) GetACL(path string) ([]ACLEntry, error) {
    output, err := a.shell.Execute("getfacl", "--omit-header", "--numeric", path)
    if err != nil {
        return nil, fmt.Errorf("failed to get ACL: %w", err)
    }

    // Parse output:
    // user::rw-
    // user:1001:rw-
    // group::r--
    // group:1005:r-x
    // mask::rwx
    // other::---

    var entries []ACLEntry
    for _, line := range strings.Split(output, "\n") {
        // Parse line und erstelle ACLEntry
    }

    return entries, nil
}

// SetACL - Setzt ACLs auf Datei/Verzeichnis
func (a *ACLManager) SetACL(path string, entries []ACLEntry) error {
    // Build setfacl commands
    var args []string
    for _, entry := range entries {
        // user:alice:rwx
        aclStr := fmt.Sprintf("%s:%s:%s", entry.Type, entry.Name, entry.Permissions)
        args = append(args, "-m", aclStr)
    }
    args = append(args, path)

    _, err := a.shell.Execute("setfacl", args...)
    return err
}

// SetDefaultACL - Setzt Default-ACLs für neue Dateien in Verzeichnis
func (a *ACLManager) SetDefaultACL(dirPath string, entries []ACLEntry) error {
    var args []string
    for _, entry := range entries {
        // default:user:alice:rwx
        aclStr := fmt.Sprintf("default:%s:%s:%s", entry.Type, entry.Name, entry.Permissions)
        args = append(args, "-m", aclStr)
    }
    args = append(args, dirPath)

    _, err := a.shell.Execute("setfacl", args...)
    return err
}
```

#### API Handler erforderlich

**Neue Datei:** `/backend/internal/api/handlers/filesystem_acl.go`

```go
package handlers

type FilesystemACLHandler struct {
    aclManager *filesystem.ACLManager
}

// Endpoints:
// GET    /api/v1/filesystem/acl?path=/path/to/file
// POST   /api/v1/filesystem/acl
// DELETE /api/v1/filesystem/acl
// POST   /api/v1/filesystem/acl/default
// POST   /api/v1/filesystem/acl/recursive
// DELETE /api/v1/filesystem/acl/all
```

**Request/Response Beispiele:**

```json
// GET /api/v1/filesystem/acl?path=/data/shared
{
  "success": true,
  "data": {
    "path": "/data/shared",
    "entries": [
      {"type": "user", "name": "alice", "permissions": "rwx"},
      {"type": "user", "name": "bob", "permissions": "r-x"},
      {"type": "group", "name": "developers", "permissions": "rw-"},
      {"type": "mask", "name": "", "permissions": "rwx"},
      {"type": "other", "name": "", "permissions": "---"}
    ]
  }
}

// POST /api/v1/filesystem/acl
{
  "path": "/data/shared",
  "entries": [
    {"type": "user", "name": "charlie", "permissions": "r--"}
  ]
}

// DELETE /api/v1/filesystem/acl
{
  "path": "/data/shared",
  "type": "user",
  "name": "bob"
}
```

#### Frontend API Client erforderlich

**Neue Datei:** `/frontend/src/api/filesystem-acl.ts`

```typescript
export interface ACLEntry {
  type: 'user' | 'group' | 'mask' | 'other';
  name: string;
  permissions: string; // rwx format
}

export interface ACLInfo {
  path: string;
  entries: ACLEntry[];
}

export const filesystemACLApi = {
  async getACL(path: string): Promise<ApiResponse<ACLInfo>>,
  async setACL(path: string, entries: ACLEntry[]): Promise<ApiResponse<any>>,
  async removeACL(path: string, type: string, name: string): Promise<ApiResponse<any>>,
  async setDefaultACL(dirPath: string, entries: ACLEntry[]): Promise<ApiResponse<any>>,
  async removeAllACLs(path: string): Promise<ApiResponse<any>>,
  async applyRecursive(dirPath: string, entries: ACLEntry[]): Promise<ApiResponse<any>>,
};
```

#### Frontend UI erforderlich

**Integration in FileManager:**

`/frontend/src/apps/FileManager/components/ACLDialog.tsx`

```tsx
interface ACLDialogProps {
  filePath: string;
  onClose: () => void;
}

export function ACLDialog({ filePath, onClose }: ACLDialogProps) {
  const [entries, setEntries] = useState<ACLEntry[]>([]);

  // Load ACLs
  useEffect(() => {
    filesystemACLApi.getACL(filePath).then(resp => {
      if (resp.success) setEntries(resp.data.entries);
    });
  }, [filePath]);

  return (
    <Dialog>
      <h2>ACL Permissions: {filePath}</h2>

      {/* ACL Entry List */}
      <ACLEntryList entries={entries} onChange={setEntries} />

      {/* Add New Entry */}
      <AddACLEntry onAdd={(entry) => setEntries([...entries, entry])} />

      {/* Default ACL Toggle (nur für Verzeichnisse) */}
      <DefaultACLSection dirPath={filePath} />

      {/* Actions */}
      <Button onClick={handleSave}>Apply ACLs</Button>
      <Button onClick={handleRemoveAll}>Remove All ACLs</Button>
    </Dialog>
  );
}
```

**UI Design:**

```
┌─────────────────────────────────────────────────┐
│  ACL Permissions: /data/shared/project          │
├─────────────────────────────────────────────────┤
│                                                 │
│  Current ACL Entries:                           │
│  ┌───────────────────────────────────────────┐  │
│  │ Type    Name           Permissions        │  │
│  ├───────────────────────────────────────────┤  │
│  │ User    alice          rwx  [Edit] [Del]  │  │
│  │ User    bob            r-x  [Edit] [Del]  │  │
│  │ Group   developers     rw-  [Edit] [Del]  │  │
│  │ Mask                   rwx                │  │
│  │ Other                  ---                │  │
│  └───────────────────────────────────────────┘  │
│                                                 │
│  Add New Entry:                                 │
│  [User ▼] [Select User...▼] [rwx ▼] [Add]      │
│                                                 │
│  ☐ Apply recursively to all files/folders      │
│  ☐ Set as default ACL for new files             │
│                                                 │
│  [Apply ACLs]  [Remove All ACLs]  [Cancel]      │
└─────────────────────────────────────────────────┘
```

**FileManager Context Menu:**

Rechtsklick auf Datei/Ordner → **"Manage ACLs..."** → Öffnet ACLDialog

#### Dependencies Check

**Debian Packages:**
```bash
apt-get install acl
```

**Dependency Checker Update:**

```go
// backend/internal/dependencies/checker.go
{
    Name:    "acl",
    Command: "getfacl",
    Package: "acl",
    Purpose: "POSIX ACL support for granular file permissions",
},
```

---

## 2. Disk Quotas

### 2.1 Was sind Disk Quotas?

Quotas limitieren den Speicherplatz, den Benutzer/Gruppen verwenden können:
- **User Quotas**: Pro Benutzer (z.B. alice: max 50GB)
- **Group Quotas**: Pro Gruppe (z.B. developers: max 500GB)
- **Soft Limit**: Warnung bei Überschreitung (Grace Period)
- **Hard Limit**: Keine weiteren Schreibvorgänge möglich

**Use Case:**
```
User alice:
  Soft Limit: 45 GB (Warnung)
  Hard Limit: 50 GB (Blockierung)
  Current Usage: 48 GB ⚠️ (über Soft Limit, Grace Period: 7 Tage)
```

### 2.2 Was bereits existiert: ZFS Quotas (unverwaltet) ⚠️

**Backend:** `/backend/internal/system/storage/zfs.go:46`

```go
type ZFSDataset struct {
    // ...
    Quota       uint64 `json:"quota"`       // DEFINIERT
    Reservation uint64 `json:"reservation"` // DEFINIERT
}
```

**Problem:** Diese Felder werden NUR gelesen, aber **NIE gesetzt**!

**Fehlende Funktionen:**
```go
func (z *ZFSManager) SetDatasetQuota(dataset string, quota uint64) error  // FEHLT
func (z *ZFSManager) SetUserQuota(dataset, user string, quota uint64) error  // FEHLT
func (z *ZFSManager) SetGroupQuota(dataset, group string, quota uint64) error  // FEHLT
func (z *ZFSManager) GetUserQuotaUsage(dataset, user string) (uint64, error)  // FEHLT
```

### 2.3 Was fehlt: Filesystem Quotas ❌

#### Backend-Implementation erforderlich

**Neue Datei:** `/backend/internal/system/filesystem/quota.go`

```go
package filesystem

type QuotaManager struct {
    shell *executor.ShellExecutor
}

// Quota Types
type QuotaType string
const (
    UserQuota  QuotaType = "user"
    GroupQuota QuotaType = "group"
)

// Quota Info
type QuotaInfo struct {
    Type        QuotaType `json:"type"`
    Name        string    `json:"name"`        // username or groupname
    Filesystem  string    `json:"filesystem"`  // /dev/sda1 or zfs pool
    BlocksUsed  uint64    `json:"blocks_used"`
    BlocksSoft  uint64    `json:"blocks_soft"`
    BlocksHard  uint64    `json:"blocks_hard"`
    InodesUsed  uint64    `json:"inodes_used"`
    InodesSoft  uint64    `json:"inodes_soft"`
    InodesHard  uint64    `json:"inodes_hard"`
    GracePeriod string    `json:"grace_period"` // "7days" oder "" wenn nicht überschritten
}

// Quota Operations
func (q *QuotaManager) EnableQuotas(filesystem string, quotaType QuotaType) error
func (q *QuotaManager) DisableQuotas(filesystem string, quotaType QuotaType) error
func (q *QuotaManager) GetQuota(filesystem string, quotaType QuotaType, name string) (*QuotaInfo, error)
func (q *QuotaManager) SetQuota(filesystem string, quotaType QuotaType, name string, softLimit, hardLimit uint64) error
func (q *QuotaManager) RemoveQuota(filesystem string, quotaType QuotaType, name string) error
func (q *QuotaManager) ListQuotas(filesystem string, quotaType QuotaType) ([]QuotaInfo, error)
func (q *QuotaManager) GetQuotaReport(filesystem string) (map[string][]QuotaInfo, error)
```

**Implementation Details:**

```go
// EnableQuotas - Aktiviert Quotas auf Filesystem
func (q *QuotaManager) EnableQuotas(filesystem string, quotaType QuotaType) error {
    // 1. Update /etc/fstab mit usrquota/grpquota Option
    // 2. Remount filesystem
    // 3. quotacheck -cug /mount/point
    // 4. quotaon /mount/point

    var quotaOpt string
    switch quotaType {
    case UserQuota:
        quotaOpt = "usrquota"
    case GroupQuota:
        quotaOpt = "grpquota"
    }

    // Update fstab
    // ... (fstab manipulation)

    // Remount
    _, err := q.shell.Execute("mount", "-o", "remount", filesystem)
    if err != nil {
        return fmt.Errorf("failed to remount: %w", err)
    }

    // Initialize quota files
    _, err = q.shell.Execute("quotacheck", "-cug", filesystem)
    if err != nil {
        return fmt.Errorf("quotacheck failed: %w", err)
    }

    // Enable quotas
    _, err = q.shell.Execute("quotaon", filesystem)
    return err
}

// SetQuota - Setzt Quota für User/Group
func (q *QuotaManager) SetQuota(filesystem string, quotaType QuotaType, name string, softLimit, hardLimit uint64) error {
    // setquota -u alice 45G 50G 0 0 /data
    //           ^  ^     ^    ^   ^  ^  ^
    //           |  |     |    |   |  |  filesystem
    //           |  |     |    |   |  inode hard
    //           |  |     |    |   inode soft
    //           |  |     |    block hard (50GB)
    //           |  |     block soft (45GB)
    //           |  username
    //           user quota type

    typeFlag := "-u"
    if quotaType == GroupQuota {
        typeFlag = "-g"
    }

    // Convert bytes to blocks (1 block = 1KB)
    softBlocks := softLimit / 1024
    hardBlocks := hardLimit / 1024

    _, err := q.shell.Execute("setquota", typeFlag, name,
        fmt.Sprintf("%d", softBlocks),
        fmt.Sprintf("%d", hardBlocks),
        "0", "0", // inode limits (0 = unlimited)
        filesystem,
    )

    return err
}

// GetQuota - Liest Quota-Info für User/Group
func (q *QuotaManager) GetQuota(filesystem string, quotaType QuotaType, name string) (*QuotaInfo, error) {
    typeFlag := "-u"
    if quotaType == GroupQuota {
        typeFlag = "-g"
    }

    // quota -u alice
    output, err := q.shell.Execute("quota", typeFlag, name)
    if err != nil {
        return nil, fmt.Errorf("failed to get quota: %w", err)
    }

    // Parse output:
    // Filesystem  blocks   quota   limit   grace   files   quota   limit   grace
    // /dev/sda1   48234567 46137344 51200000 6days  12456   0       0

    info := &QuotaInfo{
        Type:       quotaType,
        Name:       name,
        Filesystem: filesystem,
    }

    // Parse output und fülle info struct

    return info, nil
}

// ListQuotas - Listet alle Quotas auf Filesystem
func (q *QuotaManager) ListQuotas(filesystem string, quotaType QuotaType) ([]QuotaInfo, error) {
    // repquota -u /data
    typeFlag := "-u"
    if quotaType == GroupQuota {
        typeFlag = "-g"
    }

    output, err := q.shell.Execute("repquota", typeFlag, filesystem)
    if err != nil {
        return nil, err
    }

    // Parse repquota output
    var quotas []QuotaInfo
    // ... parse lines

    return quotas, nil
}
```

#### ZFS-spezifische Quotas erweitern

**Update:** `/backend/internal/system/storage/zfs.go`

```go
// SetDatasetQuota - Setzt Quota auf ZFS Dataset
func (z *ZFSManager) SetDatasetQuota(dataset string, quota uint64) error {
    quotaStr := fmt.Sprintf("%d", quota)
    if quota == 0 {
        quotaStr = "none" // Unlimitiert
    }

    _, err := z.shell.Execute("zfs", "set", "quota="+quotaStr, dataset)
    return err
}

// SetUserQuota - Setzt User-Quota auf ZFS Dataset
func (z *ZFSManager) SetUserQuota(dataset, user string, quota uint64) error {
    quotaStr := fmt.Sprintf("%d", quota)
    if quota == 0 {
        quotaStr = "none"
    }

    property := fmt.Sprintf("userquota@%s=%s", user, quotaStr)
    _, err := z.shell.Execute("zfs", "set", property, dataset)
    return err
}

// SetGroupQuota - Setzt Group-Quota auf ZFS Dataset
func (z *ZFSManager) SetGroupQuota(dataset, group string, quota uint64) error {
    quotaStr := fmt.Sprintf("%d", quota)
    if quota == 0 {
        quotaStr = "none"
    }

    property := fmt.Sprintf("groupquota@%s=%s", group, quotaStr)
    _, err := z.shell.Execute("zfs", "set", property, dataset)
    return err
}

// GetUserQuotaUsage - Liest User-Quota-Verbrauch
func (z *ZFSManager) GetUserQuotaUsage(dataset, user string) (used, quota uint64, err error) {
    // zfs userspace -H -p dataset
    output, err := z.shell.Execute("zfs", "userspace", "-H", "-p", dataset)
    if err != nil {
        return 0, 0, err
    }

    // Parse output für user
    // Columns: type name used quota
    // user  alice 48234567 51200000

    for _, line := range strings.Split(output, "\n") {
        fields := strings.Fields(line)
        if len(fields) >= 4 && fields[0] == "POSIX User" && fields[1] == user {
            used, _ = strconv.ParseUint(fields[2], 10, 64)
            quota, _ = strconv.ParseUint(fields[3], 10, 64)
            return used, quota, nil
        }
    }

    return 0, 0, fmt.Errorf("user not found")
}

// GetGroupQuotaUsage - Liest Group-Quota-Verbrauch
func (z *ZFSManager) GetGroupQuotaUsage(dataset, group string) (used, quota uint64, err error) {
    // zfs groupspace -H -p dataset
    // ... ähnlich wie GetUserQuotaUsage
}
```

#### API Handler erforderlich

**Neue Datei:** `/backend/internal/api/handlers/quota.go`

```go
// Endpoints:
// GET    /api/v1/quotas?filesystem=/dev/sda1&type=user
// GET    /api/v1/quotas/:filesystem/:type/:name
// POST   /api/v1/quotas/enable
// POST   /api/v1/quotas/disable
// POST   /api/v1/quotas
// DELETE /api/v1/quotas/:filesystem/:type/:name
// GET    /api/v1/quotas/report?filesystem=/dev/sda1
```

#### Frontend API Client erforderlich

**Neue Datei:** `/frontend/src/api/quota.ts`

```typescript
export type QuotaType = 'user' | 'group';

export interface QuotaInfo {
  type: QuotaType;
  name: string;
  filesystem: string;
  blocksUsed: number;
  blocksSoft: number;
  blocksHard: number;
  gracePeriod?: string;
}

export const quotaApi = {
  async enableQuotas(filesystem: string, type: QuotaType): Promise<ApiResponse<any>>,
  async disableQuotas(filesystem: string, type: QuotaType): Promise<ApiResponse<any>>,
  async getQuota(filesystem: string, type: QuotaType, name: string): Promise<ApiResponse<QuotaInfo>>,
  async setQuota(filesystem: string, type: QuotaType, name: string, softLimit: number, hardLimit: number): Promise<ApiResponse<any>>,
  async removeQuota(filesystem: string, type: QuotaType, name: string): Promise<ApiResponse<any>>,
  async listQuotas(filesystem: string, type: QuotaType): Promise<ApiResponse<QuotaInfo[]>>,
  async getQuotaReport(filesystem: string): Promise<ApiResponse<{ users: QuotaInfo[], groups: QuotaInfo[] }>>,
};
```

#### Frontend UI erforderlich

**Neue App:** `/frontend/src/apps/QuotaManager/QuotaManager.tsx`

```tsx
export function QuotaManager() {
  const [filesystem, setFilesystem] = useState('/dev/sda1');
  const [quotaType, setQuotaType] = useState<QuotaType>('user');
  const [quotas, setQuotas] = useState<QuotaInfo[]>([]);

  return (
    <div>
      <h1>Disk Quota Management</h1>

      {/* Filesystem Selector */}
      <FilesystemSelect value={filesystem} onChange={setFilesystem} />

      {/* Type Tabs */}
      <Tabs>
        <Tab label="User Quotas" onClick={() => setQuotaType('user')} />
        <Tab label="Group Quotas" onClick={() => setQuotaType('group')} />
      </Tabs>

      {/* Quota List */}
      <QuotaTable quotas={quotas} onEdit={handleEdit} onDelete={handleDelete} />

      {/* Add Quota */}
      <Button onClick={() => setShowAddDialog(true)}>Add Quota</Button>

      {/* Enable/Disable */}
      <QuotaStatusToggle filesystem={filesystem} type={quotaType} />
    </div>
  );
}
```

**UI Design:**

```
┌─────────────────────────────────────────────────────────────────┐
│  Disk Quota Management                                          │
├─────────────────────────────────────────────────────────────────┤
│  Filesystem: [/dev/sda1 ▼]   Status: ✅ Enabled [Disable]       │
├─────────────────────────────────────────────────────────────────┤
│  [User Quotas] [Group Quotas]                                   │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ User       Used      Soft Limit  Hard Limit  Status    │    │
│  ├─────────────────────────────────────────────────────────┤    │
│  │ alice     48.2 GB    45 GB       50 GB      ⚠️ Warning  │    │
│  │           ████████▓░ 96%                     (6d grace) │    │
│  │                                              [Edit][Del]│    │
│  ├─────────────────────────────────────────────────────────┤    │
│  │ bob       12.5 GB    20 GB       25 GB      ✅ OK       │    │
│  │           ██████░░░░ 63%                     [Edit][Del]│    │
│  ├─────────────────────────────────────────────────────────┤    │
│  │ charlie   23.8 GB    20 GB       25 GB      ❌ Over    │    │
│  │           ██████████ 100%                    [Edit][Del]│    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  [+ Add User Quota]  [Import from CSV]  [Export Report]         │
└─────────────────────────────────────────────────────────────────┘
```

**Integration in UserManager:**

Rechtsklick auf User → **"Set Disk Quota..."** → Öffnet Quota-Dialog

#### Dependencies Check

**Debian Packages:**
```bash
apt-get install quota
```

**Dependency Checker Update:**

```go
{
    Name:    "quota",
    Command: "quota",
    Package: "quota",
    Purpose: "Disk quota management for user/group limits",
},
```

---

## 3. High Availability (HA)

### 3.1 Was ist High Availability?

HA sorgt für **kontinuierliche Verfügbarkeit** durch:
- **Failover**: Automatischer Wechsel zu Backup-Server bei Ausfall
- **Replication**: Daten-Synchronisation zwischen Servern (DRBD, rsync)
- **Cluster Management**: Mehrere Nodes arbeiten zusammen (Pacemaker, Corosync)
- **Virtual IP (VIP)**: Floating IP die automatisch zum aktiven Node wandert
- **Split-Brain Prevention**: Verhindert, dass beide Nodes gleichzeitig aktiv sind

**Use Case:**
```
Primary Node (Active)    Secondary Node (Standby)
┌──────────────┐         ┌──────────────┐
│ VIP: 10.0.0.10│◄─────►│              │
│ DRBD: Primary │ Sync   │ DRBD: Sec.   │
│ Services: ON  │        │ Services: OFF│
└──────────────┘         └──────────────┘

Bei Ausfall von Primary:
┌──────────────┐         ┌──────────────┐
│ DOWN ❌      │         │ VIP: 10.0.0.10│
│              │         │ DRBD: Primary │
│              │         │ Services: ON  │
└──────────────┘         └──────────────┘
```

### 3.2 Was fehlt: HA komplett ❌

**Status:** NICHTS implementiert (0%)

#### Komponenten erforderlich:

##### 1. DRBD (Distributed Replicated Block Device)

**Was ist DRBD?**
- Echtzeit-Block-Replikation zwischen zwei Servern
- Ähnlich wie RAID-1, aber über Netzwerk
- Kernel-Modul für transparente Synchronisation

**Backend:** `/backend/internal/system/ha/drbd.go`

```go
package ha

type DRBDManager struct {
    shell *executor.ShellExecutor
}

type DRBDResource struct {
    Name          string   `json:"name"`
    Device        string   `json:"device"`        // /dev/drbd0
    Disk          string   `json:"disk"`          // /dev/sda1
    Role          string   `json:"role"`          // Primary/Secondary
    ConnectionState string `json:"connection_state"` // Connected/Disconnected
    DiskState     string   `json:"disk_state"`    // UpToDate/Outdated
    PeerRole      string   `json:"peer_role"`
    PeerDiskState string   `json:"peer_disk_state"`
    SyncProgress  float64  `json:"sync_progress"` // 0-100%
}

// DRBD Operations
func (d *DRBDManager) CreateResource(name, localDisk, remoteDisk, remoteIP string, port int) error
func (d *DRBDManager) DeleteResource(name string) error
func (d *DRBDManager) GetResourceStatus(name string) (*DRBDResource, error)
func (d *DRBDManager) ListResources() ([]DRBDResource, error)
func (d *DRBDManager) PromoteToPrimary(name string) error
func (d *DRBDManager) DemoteToSecondary(name string) error
func (d *DRBDManager) ForcePrimary(name string) error // Bei Split-Brain
func (d *DRBDManager) Disconnect(name string) error
func (d *DRBDManager) Connect(name string) error
func (d *DRBDManager) StartSync(name string) error
```

##### 2. Pacemaker & Corosync (Cluster Management)

**Was ist Pacemaker/Corosync?**
- **Corosync**: Cluster-Kommunikation und Membership
- **Pacemaker**: Resource Manager (startet/stoppt Services)
- **Heartbeat**: Überwacht Node-Verfügbarkeit

**Backend:** `/backend/internal/system/ha/cluster.go`

```go
package ha

type ClusterManager struct {
    shell *executor.ShellExecutor
}

type ClusterNode struct {
    Name        string `json:"name"`
    IP          string `json:"ip"`
    Status      string `json:"status"`      // online/offline/standby
    Resources   []ClusterResource `json:"resources"`
}

type ClusterResource struct {
    ID          string `json:"id"`
    Type        string `json:"type"`        // IP, filesystem, service
    Node        string `json:"node"`        // Welcher Node führt Resource aus
    Status      string `json:"status"`      // started/stopped/failed
    Target      string `json:"target"`      // Target Node
}

// Cluster Operations
func (c *ClusterManager) InitializeCluster(clusterName string, nodes []string) error
func (c *ClusterManager) AddNode(nodeName, nodeIP string) error
func (c *ClusterManager) RemoveNode(nodeName string) error
func (c *ClusterManager) GetClusterStatus() ([]ClusterNode, error)
func (c *ClusterManager) StandbyNode(nodeName string) error
func (c *ClusterManager) OnlineNode(nodeName string) error

// Resource Management
func (c *ClusterManager) AddVirtualIP(ip, netmask string) error
func (c *ClusterManager) AddFilesystem(device, mountpoint, fstype string) error
func (c *ClusterManager) AddService(serviceName string) error
func (c *ClusterManager) MoveResource(resourceID, targetNode string) error
func (c *ClusterManager) StopResource(resourceID string) error
func (c *ClusterManager) StartResource(resourceID string) error
```

##### 3. Keepalived (Virtual IP Management)

**Was ist Keepalived?**
- VRRP (Virtual Router Redundancy Protocol) Implementation
- Verwaltet Floating IPs zwischen Nodes
- Leichtgewichtige Alternative zu Pacemaker

**Backend:** `/backend/internal/system/ha/keepalived.go`

```go
type KeepaliveManager struct {
    shell *executor.ShellExecutor
}

type VirtualIP struct {
    IP          string `json:"ip"`
    Interface   string `json:"interface"`
    State       string `json:"state"`       // MASTER/BACKUP
    Priority    int    `json:"priority"`    // 1-255
    VirtualRouter int  `json:"virtual_router"` // VRID
}

func (k *KeepaliveManager) ConfigureVIP(vip, iface string, priority int) error
func (k *KeepaliveManager) RemoveVIP(vip string) error
func (k *KeepaliveManager) GetVIPStatus(vip string) (*VirtualIP, error)
func (k *KeepaliveManager) ListVIPs() ([]VirtualIP, error)
func (k *KeepaliveManager) SetPriority(vip string, priority int) error
```

#### API Handler erforderlich

**Neue Datei:** `/backend/internal/api/handlers/ha.go`

```go
// DRBD Endpoints:
// GET    /api/v1/ha/drbd/resources
// GET    /api/v1/ha/drbd/resources/:name
// POST   /api/v1/ha/drbd/resources
// DELETE /api/v1/ha/drbd/resources/:name
// POST   /api/v1/ha/drbd/resources/:name/promote
// POST   /api/v1/ha/drbd/resources/:name/demote
// POST   /api/v1/ha/drbd/resources/:name/force-primary

// Cluster Endpoints:
// GET    /api/v1/ha/cluster/status
// POST   /api/v1/ha/cluster/init
// POST   /api/v1/ha/cluster/nodes
// DELETE /api/v1/ha/cluster/nodes/:name
// POST   /api/v1/ha/cluster/resources/vip
// POST   /api/v1/ha/cluster/resources/filesystem
// POST   /api/v1/ha/cluster/resources/service
// POST   /api/v1/ha/cluster/resources/:id/move

// Keepalived Endpoints:
// GET    /api/v1/ha/vips
// POST   /api/v1/ha/vips
// DELETE /api/v1/ha/vips/:ip
// PUT    /api/v1/ha/vips/:ip/priority
```

#### Frontend API Client erforderlich

**Neue Datei:** `/frontend/src/api/ha.ts`

```typescript
export interface DRBDResource {
  name: string;
  device: string;
  role: 'Primary' | 'Secondary';
  connectionState: string;
  syncProgress: number;
}

export interface ClusterNode {
  name: string;
  ip: string;
  status: 'online' | 'offline' | 'standby';
  resources: ClusterResource[];
}

export const haApi = {
  // DRBD
  async listDRBDResources(): Promise<ApiResponse<DRBDResource[]>>,
  async createDRBDResource(data: CreateDRBDRequest): Promise<ApiResponse<any>>,
  async promoteDRBD(name: string): Promise<ApiResponse<any>>,
  async demoteDRBD(name: string): Promise<ApiResponse<any>>,

  // Cluster
  async getClusterStatus(): Promise<ApiResponse<ClusterNode[]>>,
  async addNode(name: string, ip: string): Promise<ApiResponse<any>>,
  async addVirtualIP(ip: string, netmask: string): Promise<ApiResponse<any>>,
  async moveResource(resourceId: string, targetNode: string): Promise<ApiResponse<any>>,

  // VIP
  async listVIPs(): Promise<ApiResponse<VirtualIP[]>>,
  async addVIP(ip: string, iface: string, priority: number): Promise<ApiResponse<any>>,
};
```

#### Frontend UI erforderlich

**Neue App:** `/frontend/src/apps/HighAvailability/HighAvailability.tsx`

```tsx
export function HighAvailability() {
  const [activeTab, setActiveTab] = useState<'cluster' | 'drbd' | 'vip'>('cluster');

  return (
    <div>
      <h1>High Availability Management</h1>

      <Tabs>
        <Tab label="Cluster Status" active={activeTab === 'cluster'} />
        <Tab label="DRBD Replication" active={activeTab === 'drbd'} />
        <Tab label="Virtual IPs" active={activeTab === 'vip'} />
      </Tabs>

      {activeTab === 'cluster' && <ClusterStatus />}
      {activeTab === 'drbd' && <DRBDManager />}
      {activeTab === 'vip' && <VIPManager />}
    </div>
  );
}
```

**UI Design (Cluster Status):**

```
┌─────────────────────────────────────────────────────────┐
│  High Availability - Cluster Status                     │
├─────────────────────────────────────────────────────────┤
│  [Cluster Status] [DRBD Replication] [Virtual IPs]      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Cluster: production-nas                                │
│  Status: ✅ Healthy (2/2 nodes online)                  │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ Node            IP            Status  Resources │   │
│  ├─────────────────────────────────────────────────┤   │
│  │ ✅ nas-primary   10.0.0.11     Online  3 active  │   │
│  │    • VIP: 10.0.0.10           (MASTER)          │   │
│  │    • /data (DRBD Primary)                        │   │
│  │    • stumpfworks-nas service                     │   │
│  │                                      [Standby]    │   │
│  ├─────────────────────────────────────────────────┤   │
│  │ ⏸️ nas-backup    10.0.0.12    Standby 0 active  │   │
│  │    • VIP: (passive)           (BACKUP)          │   │
│  │    • /data (DRBD Secondary)   ████████ 100%     │   │
│  │    • stumpfworks-nas stopped                     │   │
│  │                                      [Online]     │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  [+ Add Node]  [Trigger Failover]  [Remove Node]        │
└─────────────────────────────────────────────────────────┘
```

#### Dependencies Check

**Debian Packages:**
```bash
apt-get install drbd-utils pacemaker corosync keepalived
```

**Dependency Checker Update:**

```go
{
    Name:    "drbd",
    Command: "drbdadm",
    Package: "drbd-utils",
    Purpose: "DRBD block device replication for HA",
},
{
    Name:    "pacemaker",
    Command: "crm",
    Package: "pacemaker",
    Purpose: "Cluster resource management",
},
{
    Name:    "corosync",
    Command: "corosync-cfgtool",
    Package: "corosync",
    Purpose: "Cluster communication layer",
},
{
    Name:    "keepalived",
    Command: "keepalived",
    Package: "keepalived",
    Purpose: "Virtual IP management (VRRP)",
},
```

---

## 4. System Management App - Konsolidierung

### Konzept

**Problem:** Aktuell haben wir separate Apps für:
- 📅 Scheduled Tasks (Cron Jobs)
- ⏱️ Backups (Backup Jobs & Snapshots)
- 📊 Monitoring (nur Config, keine Metriken-Visualisierung)

**Lösung:** Eine **System Management App** die alle drei vereint.

### Neue App-Struktur

**Datei:** `/frontend/src/apps/SystemManager/SystemManager.tsx`

```tsx
type SystemTab = 'tasks' | 'backups' | 'monitoring';

export function SystemManager() {
  const [activeTab, setActiveTab] = useState<SystemTab>('monitoring');

  return (
    <div className="system-manager">
      <Header title="System Management" />

      <Tabs>
        <Tab
          id="monitoring"
          icon="📊"
          label="Monitoring"
          active={activeTab === 'monitoring'}
          onClick={() => setActiveTab('monitoring')}
        />
        <Tab
          id="tasks"
          icon="📅"
          label="Scheduled Tasks"
          active={activeTab === 'tasks'}
          onClick={() => setActiveTab('tasks')}
        />
        <Tab
          id="backups"
          icon="⏱️"
          label="Backups"
          active={activeTab === 'backups'}
          onClick={() => setActiveTab('backups')}
        />
      </Tabs>

      <TabContent>
        {activeTab === 'monitoring' && <MonitoringDashboard />}
        {activeTab === 'tasks' && <ScheduledTasks />}
        {activeTab === 'backups' && <BackupManager />}
      </TabContent>
    </div>
  );
}
```

### Tab-Komponenten

#### 1. Monitoring Dashboard

**Komponente:** `/frontend/src/apps/SystemManager/tabs/MonitoringDashboard.tsx`

**Features:**
- Real-time Metriken (CPU, RAM, Disk, Network)
- Health Score prominent anzeigen
- Charts mit Auto-Refresh (react-chartjs-2 oder recharts)
- Alert-Integration (zeige aktive Alerts)
- Prometheus Metrics Link

**UI:**
```
┌─────────────────────────────────────────────────┐
│  System Health Score: 89% ✅                    │
├─────────────────────────────────────────────────┤
│  ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │ CPU 45%  │ │ RAM 72%  │ │ Disk 68% │       │
│  │ ████▓░░  │ │ ███████░ │ │ ██████▓░ │       │
│  └──────────┘ └──────────┘ └──────────┘       │
│                                                 │
│  CPU Usage (Last 24h)                          │
│  [Line Chart]                                   │
│                                                 │
│  Memory Usage (Last 24h)                       │
│  [Area Chart]                                   │
│                                                 │
│  Active Alerts (2)                             │
│  ⚠️ Disk /data at 95% capacity                 │
│  ⚠️ High CPU on container "database"           │
└─────────────────────────────────────────────────┘
```

#### 2. Scheduled Tasks

**Komponente:** `/frontend/src/apps/SystemManager/tabs/ScheduledTasks.tsx`

**Features:**
- Existierende Tasks-Komponente einbinden
- Cron Job Management
- Task History & Logs

#### 3. Backups

**Komponente:** `/frontend/src/apps/SystemManager/tabs/BackupManager.tsx`

**Features:**
- Existierende BackupManager-Komponente einbinden
- Backup Jobs, Snapshots, Recovery

### App Registration

**Update:** `/frontend/src/apps/index.tsx`

```typescript
// ENTFERNEN (aus Dock):
// - tasks (id: 'tasks')
// - backups (id: 'backups')
// - (monitoring existiert noch nicht als App)

// NEU HINZUFÜGEN:
{
  id: 'system',
  name: 'System',
  icon: '⚙️',  // oder '🖥️'
  component: SystemManager,
  defaultSize: { width: 1400, height: 900 },
  minSize: { width: 1000, height: 700 },
}

// Settings kann bleiben (für NAS-Konfiguration)
// System ist für Monitoring/Tasks/Backups
```

### Vorteile der Konsolidierung

✅ **Weniger Dock-Icons** - Von 15 auf 13 Apps (3 Apps → 1 App)
✅ **Logische Gruppierung** - Alles was "System-Management" ist an einem Ort
✅ **Bessere UX** - Zusammenhängende Features sind zusammen
✅ **Konsistentes Design** - Einheitliches Tab-Layout

---

## 5. Zusammenfassung & Priorisierung

### Feature-Übersicht

| Feature | Backend | Frontend | Priorität | Aufwand |
|---------|---------|----------|-----------|---------|
| **ACLs** | ❌ filesystem/acl.go | ❌ ACLDialog.tsx | HOCH | 3-4 Tage |
| **Quotas** | ❌ filesystem/quota.go<br>⚠️ ZFS erweitern | ❌ QuotaManager.tsx | HOCH | 4-5 Tage |
| **HA - DRBD** | ❌ ha/drbd.go | ❌ DRBDPanel.tsx | MITTEL | 5-7 Tage |
| **HA - Cluster** | ❌ ha/cluster.go | ❌ ClusterStatus.tsx | MITTEL | 5-7 Tage |
| **HA - VIP** | ❌ ha/keepalived.go | ❌ VIPManager.tsx | MITTEL | 2-3 Tage |
| **System App** | ✅ Existiert | ⚠️ Refactor | HOCH | 2 Tage |

### Implementierungs-Reihenfolge

#### Sprint 1 (Woche 1-2): ACLs & Quotas
**Ziel:** Grundlegende Enterprise-Features

1. **Tag 1-2: Filesystem ACLs Backend**
   - `backend/internal/system/filesystem/acl.go`
   - `backend/internal/api/handlers/filesystem_acl.go`
   - API Routes registrieren
   - Testing

2. **Tag 3-4: Filesystem ACLs Frontend**
   - `frontend/src/api/filesystem-acl.ts`
   - `frontend/src/apps/FileManager/components/ACLDialog.tsx`
   - FileManager Integration (Rechtsklick → "Manage ACLs")
   - Testing

3. **Tag 5-6: Filesystem Quotas Backend**
   - `backend/internal/system/filesystem/quota.go`
   - `backend/internal/api/handlers/quota.go`
   - API Routes registrieren
   - Testing

4. **Tag 7-8: ZFS Quotas erweitern**
   - ZFS SetDatasetQuota, SetUserQuota, SetGroupQuota
   - ZFS GetUserQuotaUsage, GetGroupQuotaUsage
   - API Handler erweitern
   - Testing

5. **Tag 9-10: Quotas Frontend**
   - `frontend/src/api/quota.ts`
   - `frontend/src/apps/QuotaManager/QuotaManager.tsx`
   - UserManager Integration (Rechtsklick → "Set Quota")
   - App registrieren
   - Testing

**Deliverables:**
- ✅ Filesystem ACLs vollständig
- ✅ Filesystem Quotas vollständig
- ✅ ZFS Quotas vollständig
- ✅ Frontend-UIs für beide

#### Sprint 2 (Woche 3): System Management App Refactor
**Ziel:** UX-Verbesserung, Dock aufräumen

1. **Tag 1-2: System Manager App erstellen**
   - `frontend/src/apps/SystemManager/SystemManager.tsx`
   - Tabs: Monitoring, Tasks, Backups
   - Monitoring Dashboard mit Charts implementieren
   - Tasks & Backups Komponenten einbinden

2. **Tag 3: App Registration & Testing**
   - Alte Apps (Tasks, Backups) aus Dock entfernen
   - System Manager registrieren
   - Comprehensive Testing
   - Icon-Anpassungen (Backups: 💾 → ⏱️)

**Deliverables:**
- ✅ System Manager App mit 3 Tabs
- ✅ Dock von 15 auf 13 Apps reduziert
- ✅ Monitoring Dashboard mit Live-Charts

#### Sprint 3 (Woche 4-6): High Availability - DRBD
**Ziel:** Block Replication

1. **Tag 1-3: DRBD Backend**
   - `backend/internal/system/ha/drbd.go`
   - CreateResource, Promote, Demote, GetStatus
   - Testing mit 2-Node-Setup

2. **Tag 4-5: DRBD API Handler**
   - `backend/internal/api/handlers/ha_drbd.go`
   - API Routes
   - Testing

3. **Tag 6-7: DRBD Frontend**
   - `frontend/src/api/ha.ts` (DRBD Teil)
   - `frontend/src/apps/HighAvailability/tabs/DRBDPanel.tsx`
   - Testing

**Deliverables:**
- ✅ DRBD Resource Management
- ✅ Promote/Demote Funktionalität
- ✅ Sync Status Monitoring

#### Sprint 4 (Woche 7-9): High Availability - Cluster
**Ziel:** Cluster Management & Failover

1. **Tag 1-4: Cluster Backend**
   - `backend/internal/system/ha/cluster.go`
   - Pacemaker/Corosync Integration
   - InitializeCluster, AddNode, GetStatus
   - Resource Management (VIP, Filesystem, Service)
   - Testing mit 2-Node-Cluster

2. **Tag 5-6: Cluster API Handler**
   - `backend/internal/api/handlers/ha_cluster.go`
   - API Routes
   - Testing

3. **Tag 7-9: Cluster Frontend**
   - `frontend/src/api/ha.ts` (Cluster Teil)
   - `frontend/src/apps/HighAvailability/tabs/ClusterStatus.tsx`
   - Node-Management UI
   - Resource-Management UI
   - Failover UI
   - Testing

**Deliverables:**
- ✅ Cluster Node Management
- ✅ Resource Management (VIP, FS, Services)
- ✅ Failover Funktionalität
- ✅ Cluster Status Dashboard

#### Sprint 5 (Woche 10): High Availability - VIP
**Ziel:** Virtual IP Management

1. **Tag 1-2: Keepalived Backend**
   - `backend/internal/system/ha/keepalived.go`
   - ConfigureVIP, GetStatus, SetPriority
   - Testing

2. **Tag 3: Keepalived API Handler**
   - `backend/internal/api/handlers/ha_vip.go`
   - API Routes
   - Testing

3. **Tag 4-5: VIP Frontend**
   - `frontend/src/apps/HighAvailability/tabs/VIPManager.tsx`
   - VIP List, Add, Remove, Priority
   - Testing

**Deliverables:**
- ✅ VIP Management komplett
- ✅ VRRP Status Monitoring

#### Sprint 6 (Woche 11): Integration & Testing
**Ziel:** End-to-End Testing, Dokumentation

1. **Tag 1-2: Integration Testing**
   - Alle Features zusammen testen
   - ACLs + Quotas
   - HA Failover Szenarien
   - Performance Testing

2. **Tag 3-4: Dokumentation**
   - API Dokumentation updaten
   - User Guide für ACLs/Quotas/HA
   - Deployment Guide für HA-Setup
   - README updaten (Phase 6: 100%)

3. **Tag 5: Finalisierung**
   - Bug Fixes
   - UI Polish
   - Release vorbereiten (v1.2.0)

**Deliverables:**
- ✅ Comprehensive Testing
- ✅ Complete Documentation
- ✅ v1.2.0 Release Ready

---

## 6. Dependency Matrix

### Packages erforderlich

```bash
# ACLs
apt-get install acl

# Quotas
apt-get install quota

# High Availability
apt-get install drbd-utils pacemaker corosync keepalived

# Optional: Testing Tools
apt-get install drbdmanage-doc pacemaker-doc corosync-doc
```

### Kernel Modules

```bash
# DRBD
modprobe drbd

# Quota
# (Meist bereits im Kernel)
```

### Service Dependencies

```bash
# Enable services
systemctl enable drbd
systemctl enable pacemaker
systemctl enable corosync
systemctl enable keepalived

# Start services
systemctl start corosync
systemctl start pacemaker
systemctl start keepalived
```

---

## 7. Testing Checkliste

### ACLs Testing
- [ ] Setze User ACL auf Datei
- [ ] Setze Group ACL auf Datei
- [ ] Setze Default ACL auf Verzeichnis
- [ ] Verifiziere Vererbung (neue Dateien haben Default ACL)
- [ ] Entferne einzelnen ACL Entry
- [ ] Entferne alle ACLs
- [ ] Teste Recursive Apply
- [ ] Teste ACL mit SMB Share
- [ ] Teste ACL mit NFS Export

### Quotas Testing
- [ ] Enable User Quotas auf ext4 Filesystem
- [ ] Enable Group Quotas auf ext4 Filesystem
- [ ] Setze User Quota (Soft + Hard Limit)
- [ ] Überschreite Soft Limit → Warnung
- [ ] Überschreite Hard Limit → Blockierung
- [ ] Grace Period testen
- [ ] ZFS User Quota setzen
- [ ] ZFS Group Quota setzen
- [ ] ZFS Dataset Quota setzen
- [ ] Quota Report generieren

### HA Testing (DRBD)
- [ ] Create DRBD Resource
- [ ] Initial Sync durchführen
- [ ] Promote zu Primary
- [ ] Demote zu Secondary
- [ ] Schreibe Daten auf Primary → Sync zu Secondary
- [ ] Simuliere Disconnect → Reconnect
- [ ] Split-Brain Szenario → Force Primary
- [ ] Delete DRBD Resource

### HA Testing (Cluster)
- [ ] Initialize 2-Node Cluster
- [ ] Add Virtual IP Resource
- [ ] Add Filesystem Resource
- [ ] Add Service Resource
- [ ] Trigger Manual Failover
- [ ] Simuliere Node-Ausfall → Automatic Failover
- [ ] VIP wandert zu Backup Node
- [ ] Services starten auf Backup Node
- [ ] Primary Node Recovery → Failback
- [ ] Remove Node from Cluster

### HA Testing (VIP)
- [ ] Configure VIP mit Keepalived
- [ ] Verifiziere MASTER/BACKUP State
- [ ] Ping VIP (sollte funktionieren)
- [ ] Stoppe Keepalived auf MASTER → VIP zu BACKUP
- [ ] Starte Keepalived wieder → Failback
- [ ] Priority ändern → VIP wechselt Node

---

## 8. Gesamtaufwand & Timeline

### Aufwandschätzung

| Feature | Backend | Frontend | Testing | Gesamt |
|---------|---------|----------|---------|--------|
| **ACLs** | 2 Tage | 2 Tage | 0.5 Tage | **4.5 Tage** |
| **Quotas** | 3 Tage | 2 Tage | 0.5 Tage | **5.5 Tage** |
| **System App** | - | 2 Tage | 0.5 Tage | **2.5 Tage** |
| **HA - DRBD** | 3 Tage | 2 Tage | 2 Tage | **7 Tage** |
| **HA - Cluster** | 4 Tage | 3 Tage | 2 Tage | **9 Tage** |
| **HA - VIP** | 2 Tage | 1 Tag | 1 Tag | **4 Tage** |
| **Integration & Docs** | - | - | 5 Tage | **5 Tage** |
| **GESAMT** | **14 Tage** | **12 Tage** | **11.5 Tage** | **37.5 Tage** |

### Timeline

**Phase 6 komplett:** ~**8-10 Wochen** (2 Monate)

- **Woche 1-2:** ACLs & Quotas (10 Tage)
- **Woche 3:** System Manager Refactor (2.5 Tage) + Buffer
- **Woche 4-6:** HA - DRBD (7 Tage) + Testing
- **Woche 7-9:** HA - Cluster (9 Tage) + Testing
- **Woche 10:** HA - VIP (4 Tage) + Buffer
- **Woche 11:** Integration & Docs (5 Tage)

**Realistisch mit Puffer:** 10-12 Wochen

---

## 9. Release Plan

### v1.2.0 - Enterprise Features

**Target:** Q2 2025

**Milestones:**

#### v1.2.0-alpha (Woche 2)
- ✅ Filesystem ACLs
- ✅ Filesystem Quotas
- ⏳ ZFS Quotas (in Progress)

#### v1.2.0-beta (Woche 5)
- ✅ ACLs & Quotas komplett
- ✅ System Manager App
- ✅ DRBD Backend
- ⏳ DRBD Frontend (in Progress)

#### v1.2.0-rc1 (Woche 9)
- ✅ DRBD komplett
- ✅ Cluster Backend
- ✅ Cluster Frontend
- ⏳ VIP (in Progress)

#### v1.2.0-rc2 (Woche 10)
- ✅ Alle Features komplett
- ⏳ Testing & Bug Fixes

#### v1.2.0 Final (Woche 11)
- ✅ Production Ready
- ✅ Dokumentation komplett
- ✅ Phase 6: 100% ✅

---

## 10. Risiken & Herausforderungen

### Technische Risiken

#### ACLs
- ⚠️ **Filesystem-Support**: Nicht alle Filesystems unterstützen ACLs (ext4 ✅, xfs ✅, zfs ❌ native ACLs)
- **Mitigation**: ZFS hat eigenes ACL-System (`zfs allow`), dokumentieren

#### Quotas
- ⚠️ **Kernel-Config**: Quotas müssen im Kernel aktiviert sein
- ⚠️ **Filesystem-Remount**: Erfordert Remount mit usrquota/grpquota Option
- **Mitigation**: Dependency Checker erweitern, Auto-Config in postinst

#### High Availability
- 🔴 **Komplexität**: HA ist sehr komplex und fehleranfällig
- 🔴 **Split-Brain**: Katastrophales Szenario wenn beide Nodes gleichzeitig Primary werden
- 🔴 **Network**: Erfordert dediziertes Heartbeat-Netzwerk (separates Interface)
- 🔴 **Testing**: Benötigt mindestens 2 physische/virtuelle Maschinen
- **Mitigation**:
  - Comprehensive Documentation
  - Split-Brain Detection & Prevention (STONITH)
  - Wizard für HA-Setup
  - Extensive Testing mit VMs

### Deployment-Risiken

- ⚠️ **Breaking Changes**: Neue Kernel-Module (drbd) könnten Neustart erfordern
- ⚠️ **Downtime**: HA-Setup erfordert temporären Service-Ausfall
- **Mitigation**:
  - Migrations-Guide
  - Dry-Run Mode für HA-Setup
  - Rollback-Strategie

---

## 11. Dokumentation erforderlich

### User Guides

1. **ACL Management Guide**
   - Was sind ACLs?
   - Use Cases
   - Best Practices
   - Troubleshooting

2. **Quota Management Guide**
   - Was sind Quotas?
   - User vs Group Quotas
   - Grace Periods
   - ZFS vs ext4 Quotas
   - Troubleshooting

3. **High Availability Guide**
   - HA Konzepte
   - 2-Node vs 3-Node Cluster
   - DRBD Setup
   - Pacemaker/Corosync Setup
   - Failover Testing
   - Split-Brain Recovery
   - Production Deployment

### Admin Guides

1. **ACL Administration**
   - CLI Commands (getfacl, setfacl)
   - Batch Operations
   - ACL Migration
   - Permission Debugging

2. **Quota Administration**
   - CLI Commands (quota, setquota, repquota)
   - Enable/Disable Quotas
   - Quota Enforcement
   - Performance Impact

3. **HA Administration**
   - DRBD Administration (drbdadm)
   - Cluster Administration (crm, pcs)
   - Node Maintenance
   - Emergency Procedures
   - Monitoring & Alerting

### API Documentation

Update `/docs/API.md`:
- ACL Endpoints
- Quota Endpoints
- HA Endpoints (DRBD, Cluster, VIP)

---

## 12. Nächste Schritte

### Sofort (Diese Woche):
1. ✅ System Manager App Refactor (2-3 Tage)
   - Monitoring Dashboard mit Charts
   - Tasks & Backups einbinden
   - Dock aufräumen

### Sprint 1 (Woche 1-2):
2. ⏳ Filesystem ACLs implementieren (4 Tage Backend + Frontend)
3. ⏳ Filesystem Quotas implementieren (5 Tage Backend + Frontend)

### Sprint 2 (Woche 3-6):
4. ⏳ DRBD implementieren (7 Tage)
5. ⏳ Cluster Management implementieren (9 Tage)

### Sprint 3 (Woche 7-10):
6. ⏳ VIP Management implementieren (4 Tage)
7. ⏳ Integration Testing & Docs (5 Tage)

### Release:
8. 🎯 v1.2.0 Release (Woche 11)

---

**Erstellt von:** Claude
**Branch:** `claude/audit-frontend-integration-017ykh14kLotgaQ7URfFanDX`
**Repository:** Stumpf-works/stumpfworks-nas
