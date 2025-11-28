# VSCode Claude Code Implementation Prompt

**Ziel:** Implementiere alle dokumentierten Features aus den Audit-Dokumenten schrittweise

---

## 📋 Verfügbare Dokumentation

Die folgenden Dokumente im Repository enthalten alle technischen Spezifikationen:

1. **FRONTEND_INTEGRATION_AUDIT.md** - Frontend-Lücken Audit (AD DC, Network Manager, Dock/UX)
2. **PHASE_6_ENTERPRISE_AUDIT.md** - Enterprise Features (ACLs, Quotas, HA) - Vollständige Tech-Specs
3. **TODO_PHASE_6.md** - Sprint-Checklisten für Phase 6 Implementation
4. **CLUSTER_INTEGRATION.md** - Cluster Features (Optional, GlusterFS, HAProxy, Swarm)
5. **CLUSTER_ACTIVATION_WIZARD.md** - Easy Cluster Setup Flow

---

## 🎯 Implementation-Reihenfolge

### **PRIORITÄT 1: SOFORT (Diese Woche)** 🔴

Bitte implementiere **Sprint 0** aus `TODO_PHASE_6.md`:

**System Manager App Refactor** (2-3 Tage)

**Aufgabe:** Konsolidiere Tasks, Backups und Monitoring in eine App

**Dateien zu erstellen:**
- `/frontend/src/apps/SystemManager/SystemManager.tsx`
- `/frontend/src/apps/SystemManager/tabs/MonitoringDashboard.tsx`
- `/frontend/src/apps/SystemManager/tabs/ScheduledTasks.tsx`
- `/frontend/src/apps/SystemManager/tabs/BackupManager.tsx`

**Dateien zu ändern:**
- `/frontend/src/apps/index.tsx` - Remove `tasks`, `backups`; Add `system`

**Features:**
1. **Monitoring Dashboard Tab:**
   - Health Score Widget (groß und prominent)
   - System Metrics Cards (CPU, RAM, Disk, Network)
   - Live Charts mit Auto-Refresh (5s interval)
     - CPU Usage Chart (react-chartjs-2 oder recharts)
     - Memory Usage Chart
     - Disk I/O Chart
     - Network Traffic Chart
   - Active Alerts Section
   - Link zu `/metrics` (Prometheus)

2. **Tasks Tab:**
   - Import bestehende Tasks Component
   - Wrapper für Integration

3. **Backups Tab:**
   - Import bestehende BackupManager Component
   - Wrapper für Integration

**Icon-Änderungen:**
- Backups: 💾 → ⏱️ (wenn noch separat verwendet)
- System Manager: ⚙️ oder 🖥️

**Resultat:**
- Dock reduziert von 15 auf 13 Apps
- Monitoring Dashboard mit Live-Charts
- Bessere Benutzerfreundlichkeit

---

### **PRIORITÄT 2: KRITISCH (Woche 1-2)** 🔴

Bitte implementiere **Sprint 1** aus `TODO_PHASE_6.md`:

**Filesystem ACLs** (4-5 Tage)

**Vollständige Specs in:** `PHASE_6_ENTERPRISE_AUDIT.md` Abschnitt 1

**Backend-Dateien zu erstellen:**
- `/backend/internal/system/filesystem/acl.go`
- `/backend/internal/api/handlers/filesystem_acl.go`

**Backend-Dateien zu ändern:**
- `/backend/internal/api/router.go` - Add ACL routes
- `/backend/internal/dependencies/checker.go` - Add `acl` package check

**Frontend-Dateien zu erstellen:**
- `/frontend/src/api/filesystem-acl.ts`
- `/frontend/src/apps/FileManager/components/ACLDialog.tsx`

**Frontend-Dateien zu ändern:**
- `/frontend/src/apps/FileManager/FileManager.tsx` - Add context menu "Manage ACLs..."

**Implementation-Details:**
- Siehe `PHASE_6_ENTERPRISE_AUDIT.md` Zeilen 93-400 für vollständige Code-Beispiele
- 6 API Endpoints (GET, POST, DELETE für ACLs)
- `getfacl` und `setfacl` Linux-Commands nutzen
- UI mit ACL Entry List + Add/Remove

---

### **PRIORITÄT 3: HOCH (Woche 2-3)** 🟡

Bitte implementiere **Sprint 2** aus `TODO_PHASE_6.md`:

**Disk Quotas** (5-6 Tage)

**Vollständige Specs in:** `PHASE_6_ENTERPRISE_AUDIT.md` Abschnitt 2

**Backend-Dateien zu erstellen:**
- `/backend/internal/system/filesystem/quota.go`
- `/backend/internal/api/handlers/quota.go`

**Backend-Dateien zu ändern:**
- `/backend/internal/system/storage/zfs.go` - Add Quota methods
- `/backend/internal/api/router.go` - Add Quota routes
- `/backend/internal/dependencies/checker.go` - Add `quota` package

**Frontend-Dateien zu erstellen:**
- `/frontend/src/api/quota.ts`
- `/frontend/src/apps/QuotaManager/QuotaManager.tsx`

**Frontend-Dateien zu ändern:**
- `/frontend/src/apps/index.tsx` - Add QuotaManager app
- `/frontend/src/apps/UserManager/UserManager.tsx` - Add "Set Disk Quota..." context menu

**Implementation-Details:**
- Siehe `PHASE_6_ENTERPRISE_AUDIT.md` Zeilen 402-700
- 7 API Endpoints
- User/Group Quotas für ext4 + ZFS
- `quota`, `setquota`, `repquota` Commands
- UI mit Progress Bars für Usage

---

### **PRIORITÄT 4: MITTEL (Woche 3-6)** 🟡

Bitte implementiere **Sprint 3** aus `TODO_PHASE_6.md`:

**Active Directory Domain Controller UI** (4-5 Tage)

**Vollständige Specs in:** `FRONTEND_INTEGRATION_AUDIT.md` Abschnitt 1

**Backend:** Bereits vollständig implementiert in `/backend/internal/api/handlers/ad_dc.go` (1.443 Zeilen, 40+ Endpoints)

**Frontend-Dateien zu erstellen:**
- `/frontend/src/api/addc.ts` - API Client mit 40+ Methoden
- `/frontend/src/apps/ADDomainController/ADDomainController.tsx`
- `/frontend/src/apps/ADDomainController/components/DomainStatus.tsx`
- `/frontend/src/apps/ADDomainController/components/UserManagement.tsx`
- `/frontend/src/apps/ADDomainController/components/GroupManagement.tsx`
- `/frontend/src/apps/ADDomainController/components/ComputerManagement.tsx`
- `/frontend/src/apps/ADDomainController/components/OUManagement.tsx`
- `/frontend/src/apps/ADDomainController/components/GPOManagement.tsx`
- `/frontend/src/apps/ADDomainController/components/DNSManagement.tsx`
- `/frontend/src/apps/ADDomainController/components/FSMOManagement.tsx`

**Frontend-Dateien zu ändern:**
- `/frontend/src/apps/index.tsx` - Add AD DC app

**Tabs in AD DC App:**
1. Domain Status & Provisioning
2. User Management (Create, Delete, Enable, Disable, Password, Expiry)
3. Group Management (Create, Delete, Members)
4. Computer Management
5. Organizational Units
6. Group Policy Objects
7. DNS Zone Management
8. FSMO Roles

**Implementation-Details:**
- Siehe `FRONTEND_INTEGRATION_AUDIT.md` Zeilen 1-210
- Router bereits registriert, nur Frontend fehlt
- Icon: 🏢 oder 🔐

---

### **PRIORITÄT 5: OPTIONAL (Woche 4-10)** 🟢

Bitte implementiere **Sprint 3-5** aus `TODO_PHASE_6.md`:

**High Availability Features** (20 Tage)

**Vollständige Specs in:** `PHASE_6_ENTERPRISE_AUDIT.md` Abschnitt 3

**Komponenten:**
1. **DRBD** (7 Tage) - Block Replication
2. **Pacemaker/Corosync** (9 Tage) - Cluster Management
3. **Keepalived** (4 Tage) - Virtual IP

**Backend-Dateien zu erstellen:**
- `/backend/internal/system/ha/drbd.go`
- `/backend/internal/system/ha/cluster.go`
- `/backend/internal/system/ha/keepalived.go`
- `/backend/internal/api/handlers/ha_drbd.go`
- `/backend/internal/api/handlers/ha_cluster.go`
- `/backend/internal/api/handlers/ha_vip.go`

**Frontend-Dateien zu erstellen:**
- `/frontend/src/api/ha.ts`
- `/frontend/src/apps/HighAvailability/HighAvailability.tsx`
- `/frontend/src/apps/HighAvailability/tabs/DRBDPanel.tsx`
- `/frontend/src/apps/HighAvailability/tabs/ClusterStatus.tsx`
- `/frontend/src/apps/HighAvailability/tabs/VIPManager.tsx`

**Implementation-Details:**
- Siehe `PHASE_6_ENTERPRISE_AUDIT.md` Zeilen 702-1200
- Benötigt 2-Node Setup für Testing
- Dependencies: drbd-utils, pacemaker, corosync, keepalived

---

### **PRIORITÄT 6: ZUKUNFT (Optional)** 🔵

**Cluster Integration** (v1.3.0, 15 Wochen)

**Vollständige Specs in:**
- `CLUSTER_INTEGRATION.md` - Tech-Details
- `CLUSTER_ACTIVATION_WIZARD.md` - UX Flow

**Features:**
1. Cluster Activation Wizard (One-Click Setup)
2. GlusterFS Distributed Storage
3. HAProxy Load Balancing
4. Docker Swarm Orchestration
5. etcd Distributed Configuration
6. Cluster Manager UI

**Hinweis:** Erst NACH Phase 6 (HA) implementieren!

---

## 🔧 Implementation-Workflow

Für jede Aufgabe folge diesem Workflow:

### Schritt 1: Lese die Dokumentation
```
Lese die relevante Section aus:
- PHASE_6_ENTERPRISE_AUDIT.md (für Tech-Specs)
- TODO_PHASE_6.md (für Checkliste)
- FRONTEND_INTEGRATION_AUDIT.md (für Frontend-Gaps)
```

### Schritt 2: Backend Implementation
```
1. Erstelle Manager-Klasse (z.B. ACLManager, QuotaManager)
2. Implementiere alle Methoden mit Shell-Commands
3. Erstelle API Handler
4. Registriere Routes in router.go
5. Update dependencies/checker.go
6. Teste mit curl oder Postman
```

### Schritt 3: Frontend Implementation
```
1. Erstelle API Client (z.B. filesystem-acl.ts)
2. Erstelle UI Component (z.B. ACLDialog.tsx)
3. Integriere in bestehende App
4. Teste im Browser
```

### Schritt 4: Testing
```
1. Funktionalität testen
2. Error Handling testen
3. UI/UX prüfen
4. Edge Cases testen
```

### Schritt 5: Commit
```
git add .
git commit -m "feat: Implement [Feature Name]

- Add backend: [Manager + Handler]
- Add frontend: [API Client + Component]
- Add routes and dependencies
- Tests: [What was tested]

Refs: [Document Name] Section [X]"
```

---

## 📝 Wichtige Hinweise

### Dependencies
Alle benötigten Debian-Packages sind dokumentiert in:
- `PHASE_6_ENTERPRISE_AUDIT.md` - Dependencies für ACLs, Quotas, HA
- `CLUSTER_INTEGRATION.md` - Dependencies für Cluster

**Installation:**
```bash
# ACLs
apt-get install acl

# Quotas
apt-get install quota

# HA
apt-get install drbd-utils pacemaker corosync keepalived

# Cluster (später)
apt-get install glusterfs-server haproxy etcd
```

### Code-Beispiele
Vollständige Code-Beispiele findest du in:
- `PHASE_6_ENTERPRISE_AUDIT.md` - Go + TypeScript Code
- `CLUSTER_INTEGRATION.md` - Go + TypeScript Code

**Beispiel-Pfade:**
- ACLs: Zeilen 93-400
- Quotas: Zeilen 402-700
- HA DRBD: Zeilen 702-900
- HA Cluster: Zeilen 902-1100

### Testing Checklisten
Detaillierte Testing-Schritte in:
- `PHASE_6_ENTERPRISE_AUDIT.md` Zeilen 1300-1500
- `TODO_PHASE_6.md` - Checkboxen für jeden Sprint

---

## 🎯 Schnellstart für diese Sitzung

**Empfohlener Fokus:**

```
Option A: System Manager App (SOFORT)
  → Reduziert Dock-Clutter
  → Live Monitoring Dashboard
  → Bessere UX
  → 2-3 Tage

Option B: Filesystem ACLs (KRITISCH)
  → Enterprise Feature
  → Vollständig dokumentiert
  → 4-5 Tage

Option C: AD Domain Controller UI (HOCH)
  → Backend existiert bereits
  → Nur Frontend fehlt
  → 4-5 Tage
```

**Meine Empfehlung:**

1. **HEUTE:** Start mit System Manager App (Sprint 0)
2. **DIESE WOCHE:** Finish System Manager + Start ACLs
3. **NÄCHSTE WOCHE:** Quotas + AD DC UI

---

## 📚 Dokumentations-Referenz

| Feature | Primary Doc | Backup Doc | Code Examples |
|---------|-------------|------------|---------------|
| **System Manager** | TODO_PHASE_6.md Sprint 0 | PHASE_6_ENTERPRISE_AUDIT.md | Zeilen 1-50 |
| **ACLs** | PHASE_6_ENTERPRISE_AUDIT.md § 1 | TODO_PHASE_6.md Sprint 1 | Zeilen 93-400 |
| **Quotas** | PHASE_6_ENTERPRISE_AUDIT.md § 2 | TODO_PHASE_6.md Sprint 2 | Zeilen 402-700 |
| **AD DC UI** | FRONTEND_INTEGRATION_AUDIT.md § 1 | - | Zeilen 1-210 |
| **HA** | PHASE_6_ENTERPRISE_AUDIT.md § 3 | TODO_PHASE_6.md Sprint 3-5 | Zeilen 702-1200 |
| **Cluster** | CLUSTER_INTEGRATION.md | CLUSTER_ACTIVATION_WIZARD.md | Alle |

---

## ✅ Success Criteria

Nach Completion dieser Implementation:

- ✅ Dock reduziert auf 13 Apps (von 15)
- ✅ Live Monitoring Dashboard mit Charts
- ✅ Filesystem ACLs voll funktionsfähig
- ✅ User/Group Disk Quotas implementiert
- ✅ AD Domain Controller vollständig verwaltbar
- ✅ Optional: 2-Node HA mit DRBD + Pacemaker
- ✅ Optional: 3+ Node Cluster mit GlusterFS

**README.md Update:**
- Phase 6: 10% → 100% ✅
- v1.2.0 Release Ready

---

## 🚀 Los geht's!

**Dein erster Command:**

```
Bitte implementiere Sprint 0 aus TODO_PHASE_6.md:
System Manager App Refactor

Erstelle:
1. /frontend/src/apps/SystemManager/SystemManager.tsx
2. /frontend/src/apps/SystemManager/tabs/MonitoringDashboard.tsx
3. /frontend/src/apps/SystemManager/tabs/ScheduledTasks.tsx
4. /frontend/src/apps/SystemManager/tabs/BackupManager.tsx

Ändere:
- /frontend/src/apps/index.tsx (Remove tasks/backups, Add system)

Features:
- Monitoring Dashboard mit Live-Charts (CPU, RAM, Disk, Network)
- Health Score prominent anzeigen
- Auto-Refresh alle 5 Sekunden
- Tasks & Backups Tabs als Wrapper

Nutze react-chartjs-2 für Charts.
Nutze bestehende monitoringApi für Daten.

Starte jetzt!
```

---

**Viel Erfolg! 🚀**
