# 🚀 Stumpf.Works NAS - TODO & ROADMAP

**Last Updated:** 2025-11-13
**Current Version:** v0.4.0
**Status:** 99% Feature Complete - Production-Ready für Single-Host

---

## 🎯 PROJECT GOAL

**Become a TRUE Unraid/TrueNAS alternative:**
- ✅ Easy to use Web UI
- ✅ SMB/NFS Share Management
- ✅ Docker Container Management
- ✅ User & Permission Management
- ✅ Monitoring & Alerting
- ⚠️ Advanced Features (ACLs, Quotas, HA)
- ❌ Multi-Site Replication
- ❌ ZFS Advanced Features

**Target Users:**
- Home Lab Enthusiasts ✅
- Small Business (< 50 Users) ✅
- Dev/Test Environments ✅
- Enterprise (1000+ Users) ❌ (Not Yet)

---

## 🔴 CRITICAL BUGS (FIX IMMEDIATELY!)

### 1. Upload Permission Security Hole 🔒
**Priority:** P0 - CRITICAL
**File:** `backend/internal/api/handlers/files.go:242`
**Issue:**
```go
FinalizeUpload() → ⚠️ TODO: Implement permission check
UploadChunk() → ❌ NO Permission Check!
```
User kann Chunks hochladen ohne dass Destination-Path Permission geprüft wird!

**Fix:**
```go
func FinalizeUpload(w http.ResponseWriter, r *http.Request) {
    // 1. Get SecurityContext
    ctx, err := getSecurityContext(r)
    if err != nil {
        utils.RespondError(w, err)
        return
    }

    // 2. Check destination path BEFORE moving file
    destDir := filepath.Dir(req.DestinationPath)
    if err := fileService.permissions.CanWrite(ctx, destDir); err != nil {
        uploadManager.CancelUpload(req.SessionID) // Cleanup!
        utils.RespondError(w, err)
        return
    }

    // 3. NOW finalize
    if err := uploadManager.FinalizeUpload(req.SessionID, req.DestinationPath); err != nil {
        ...
    }
}
```

**Estimated Time:** 30 minutes
**Testing:** Try to upload to restricted share, should fail

---

### 2. Samba ValidUsers Keine Validierung
**Priority:** P1 - HIGH
**File:** `backend/internal/storage/shares.go`
**Issue:** User können non-existent usernames in ValidUsers eintragen

**Fix:**
```go
func CreateShare(req *CreateShareRequest) (*Share, error) {
    // Validate ValidUsers exist
    for _, username := range req.ValidUsers {
        if _, err := users.GetUserByUsername(username); err != nil {
            return nil, fmt.Errorf("user '%s' does not exist", username)
        }
    }

    // ... rest of creation
}
```

**Estimated Time:** 1 hour
**Testing:** Try to create share with invalid username

---

### 3. User/Group Name Lookup Missing
**Priority:** P1 - HIGH
**File:** `backend/internal/files/permissions.go:103,159,169`
**Issue:** Multiple TODOs for UID/GID → Username/Groupname lookup

**Fix:**
```go
func getUsername(uid int) string {
    u, err := user.LookupId(strconv.Itoa(uid))
    if err != nil {
        return fmt.Sprintf("uid:%d", uid)
    }
    return u.Username
}

func getGroupname(gid int) string {
    g, err := user.LookupGroupId(strconv.Itoa(gid))
    if err != nil {
        return fmt.Sprintf("gid:%d", gid)
    }
    return g.Name
}
```

**Estimated Time:** 1 hour
**Testing:** Check FileInfo owner/group displays correctly

---

## 🟡 HIGH PRIORITY (Diese/Nächste Woche)

### 4. Group Management System
**Status:** ❌ Not Implemented
**Impact:** Users können nicht in Groups organisiert werden

**Implementation:**
1. **Database Model:** `models/group.go`
   ```go
   type Group struct {
       ID          uint
       Name        string `gorm:"unique"`
       Description string
       Members     []User `gorm:"many2many:user_groups;"`
   }
   ```

2. **API Endpoints:**
   ```
   GET    /api/v1/groups
   POST   /api/v1/groups
   PUT    /api/v1/groups/{id}
   DELETE /api/v1/groups/{id}
   POST   /api/v1/groups/{id}/members
   DELETE /api/v1/groups/{id}/members/{userId}
   ```

3. **Share Integration:**
   - Add `AllowedGroups []string` to Share Model
   - Update `GetAllowedPathsForUser()` to check groups
   - Update Samba config to write `@groupname`

4. **Frontend UI:**
   - New App: GroupManager
   - User-Picker Component
   - Group Assignment in ShareManager

**Estimated Time:** 2-3 days
**Benefits:**
- Easier permission management
- Scales better with many users
- Standard feature in all NAS systems

---

### 5. Permission Audit Logging
**Status:** ⚠️ Partial (nur Basic Ops geloggt)
**Impact:** Compliance & Security

**Implementation:**
1. Extend AuditLog with Permission-specific fields:
   ```go
   type PermissionChange struct {
       FileOp      string  // chmod, chown
       OldMode     string  // "755"
       NewMode     string  // "644"
       OldOwner    string  // "user1"
       NewOwner    string  // "user2"
   }
   ```

2. Log on every permission change:
   - `ChangePermissions()` → Audit
   - `CreateDirectory()` with custom perms → Audit
   - Share `ValidUsers` change → Audit

3. Frontend: Permission Change History View

**Estimated Time:** 1 day
**Benefits:** Security, Compliance, Debugging

---

### 6. Warning für Unsichere Permissions
**Status:** ❌ Not Implemented
**Impact:** Security Risk (777, etc.)

**Implementation:**
```go
func validatePermissionSecurity(mode os.FileMode) []string {
    warnings := []string{}

    if mode & 0o002 != 0 {
        warnings = append(warnings, "World-writable: Anyone can modify this file!")
    }
    if mode & 0o007 == 0o007 {
        warnings = append(warnings, "Full permissions for others (rwx)!")
    }

    return warnings
}
```

Frontend: Show warnings in PermissionsModal
Backend: Return warnings in API response

**Estimated Time:** 2 hours

---

### 7. Permission Templates
**Status:** ❌ Not Implemented
**Impact:** Usability

**Implementation:**
1. Presets in Frontend (Already exists partially)
2. Backend: Store Custom Templates
   ```go
   type PermissionTemplate struct {
       Name        string  // "Public Files"
       Description string
       Mode        string  // "755"
       Recursive   bool
   }
   ```

3. UI: "Apply Template" Button in FileManager

**Estimated Time:** 1 day

---

## 🟢 MEDIUM PRIORITY (Nächste 2 Wochen)

### 8. User Quota Management
**Status:** ❌ Not Implemented
**Why Important:** Disk voll → System crash

**Implementation:**
```go
type UserQuota struct {
    UserID      uint
    Limit       int64   // Bytes
    Used        int64   // Current usage
    Warning     int64   // Warning threshold (80%)
}
```

- Track per-User disk usage
- Enforce on File Operations
- Alert when >80%
- Admin UI for Quota Management

**Estimated Time:** 3-4 days

---

### 9. ACL Support (Access Control Lists)
**Status:** ❌ Not Implemented
**Why Important:** Fine-grained permissions

**Scope:**
- Linux Extended Attributes (xattr)
- NFSv4 ACLs
- NOT Windows ACLs (zu komplex)

**Implementation:**
```go
type ACL struct {
    Type        string  // "user" | "group"
    Name        string  // username/groupname
    Permissions string  // "rwx"
}

func SetACL(path string, acls []ACL) error
func GetACL(path string) ([]ACL, error)
```

**Estimated Time:** 1 week (complex!)

---

### 10. Share Performance Optimizations
**Status:** ⚠️ Works but not optimal

**Issues:**
- Share list loads ALL shares on every request
- Permission checks iterate all shares
- No caching

**Optimizations:**
```go
// Cache share permissions
type ShareCache struct {
    userShares map[uint][]string  // userID → allowed paths
    mu         sync.RWMutex
    ttl        time.Duration
}

func (c *ShareCache) Get(userID uint) ([]string, bool)
func (c *ShareCache) Invalidate(userID uint)
```

**Estimated Time:** 2 days

---

## 🔵 LOW PRIORITY (Nice-to-Have)

### 11. Dark/Light Theme Toggle
**Status:** ⚠️ Partial (Theme store exists, but no toggle)

**Implementation:**
- Add Theme Toggle in TopBar
- Persist preference in LocalStorage
- Already works via store

**Estimated Time:** 30 minutes

---

### 12. Multi-Language Support (i18n)
**Status:** ❌ Not Implemented
**Why:** International users

**Scope:**
- react-i18next integration
- German translation
- English (default)

**Estimated Time:** 3-4 days

---

### 13. Mobile-Responsive Design
**Status:** ⚠️ Partially Responsive

**Issues:**
- Dashboard works
- FileManager table not optimal on mobile
- Modals too wide

**Estimated Time:** 2-3 days

---

### 14. WebSocket für Live-Updates
**Status:** ⚠️ Backend vorhanden, nicht überall genutzt

**Scope:**
- File-Upload Progress
- Live Metrics (statt polling)
- Container Status Changes
- Notifications

**Estimated Time:** 2 days

---

### 15. Plugin System Completion
**Status:** ⚠️ Loading works, Execute unclear

**Investigation:** Is plugin execution implemented?
**If not:** Implement safe sandboxed execution

**Estimated Time:** Unknown (needs investigation)

---

## 🌟 FUTURE FEATURES (3+ Months)

### 16. Monitoring Dashboard mit Charts
**Status:** 🔄 In Progress
**Backend:** ✅ Done (Metrics Collection)
**Frontend:** ❌ Charts missing

**Implementation:**
- Install recharts: `npm install recharts`
- Create Chart Components:
  - CPU Line Chart
  - Memory Line Chart
  - Network Area Chart
  - Health Score Gauge

**Estimated Time:** 1-2 days

---

### 17. Snapshot Management (Btrfs/ZFS)
**Status:** ⚠️ Backend checks exist, no management

**Scope:**
- List snapshots
- Create snapshot
- Restore from snapshot
- Schedule automatic snapshots

**Requirements:** Btrfs or ZFS filesystem
**Estimated Time:** 1 week

---

### 18. Replication & Sync
**Status:** ❌ Not Implemented

**Scope:**
- Rsync-based sync to remote NAS
- Schedule sync jobs
- Conflict resolution
- Bandwidth limiting

**Estimated Time:** 2 weeks

---

### 19. S3-Compatible Object Storage
**Status:** ❌ Not Implemented

**Scope:**
- MinIO integration
- S3 API compatibility
- Bucket management UI

**Estimated Time:** 2 weeks

---

### 20. VM Management (KVM/QEMU)
**Status:** ❌ Not Implemented

**Scope:**
- Like Proxmox/Unraid VMs
- VM creation/start/stop
- VNC console
- Resource allocation

**Estimated Time:** 3-4 weeks (major feature)

---

## 📦 TESTING & QA

### Required Tests:
- [ ] Unit Tests (Backend) - Currently ~60%
- [ ] Integration Tests (APIs)
- [ ] E2E Tests (Frontend)
- [ ] Security Audit (Penetration Testing)
- [ ] Performance Testing (Load/Stress)
- [ ] Multi-Browser Testing (Chrome, Firefox, Safari, Edge)
- [ ] Windows SMB Client Testing
- [ ] macOS SMB Client Testing
- [ ] Linux NFS Client Testing

**Estimated Time:** 2 weeks dedicated QA

---

## 📚 DOCUMENTATION

### Missing Docs:
- [ ] API Documentation (Swagger/OpenAPI)
- [ ] User Manual
- [ ] Admin Guide
- [ ] Developer Guide
- [ ] Architecture Diagrams
- [ ] Deployment Guide (Docker/Bare Metal)
- [ ] Backup/Recovery Guide
- [ ] Troubleshooting Guide

**Estimated Time:** 1 week technical writing

---

## 🔧 DEVOPS & INFRASTRUCTURE

### CI/CD Pipeline:
- [ ] GitHub Actions workflow
- [ ] Automated Testing
- [ ] Build & Release automation
- [ ] Docker Image builds
- [ ] Changelog generation

### Deployment:
- [ ] Docker Compose setup
- [ ] Kubernetes manifests
- [ ] Ansible playbooks
- [ ] Systemd service files

**Estimated Time:** 3-4 days

---

## 📊 METRICS & ROADMAP

### Completion Status:
```
✅ Core Features:        159/161 = 99%
⚠️ Bug Fixes Needed:    3 Critical
🔄 In Progress:         1 Feature
❌ Missing Features:     20+ Identified
📚 Documentation:        40%
🧪 Test Coverage:        60%

OVERALL: 85% Production-Ready
```

### Timeline:
```
Week 1-2:   🔴 Critical Bugs + 🟡 High Priority (Items 1-7)
Week 3-4:   🟢 Medium Priority (Items 8-10)
Month 2:    🔵 Low Priority (Items 11-15)
Month 3+:   🌟 Future Features (Items 16-20)
```

---

## 🎉 SUCCESS CRITERIA

**Version 1.0 (Production Release) Requires:**
- ✅ All Critical Bugs fixed
- ✅ Groups & Quota implemented
- ✅ ACL Support (basic)
- ✅ Monitoring Charts
- ✅ >80% Test Coverage
- ✅ Complete Documentation
- ✅ Docker Deployment
- ✅ Security Audit passed

**Estimated Time to v1.0:** 2-3 Months

---

## 🤝 CONTRIBUTING

Want to help? Pick a task from the list!
See `CONTRIBUTING.md` for guidelines.

**Easy First Issues:**
- Dark Theme Toggle (#11)
- Warning für 777 Permissions (#6)
- Permission Templates (#7)

**Medium Difficulty:**
- Group Management (#4)
- Monitoring Charts (#16)
- User Quotas (#8)

**Hard (Expert):**
- ACL Support (#9)
- Replication (#18)
- VM Management (#20)

---

**Last Updated:** 2025-11-13
**Maintainer:** Stumpf.Works Team
**Status:** Active Development 🚀
