# Stumpf.Works NAS - Frontend Integration Audit

**Datum:** 2025-11-26
**Projekt:** Stumpf.Works NAS
**Branch:** `claude/audit-frontend-integration-017ykh14kLotgaQ7URfFanDX`

---

## Executive Summary

Dieses Audit identifiziert Backend-Funktionalitäten, die keine oder unvollständige Frontend-Integration haben, sowie kritische UX/UI-Verbesserungen für eine produktionsreife Anwendung.

**Kritische Findings:**
- ❌ **Active Directory Domain Controller**: 40+ Backend-Endpoints ohne jegliche UI (KRITISCH)
- ⚠️ **Network Manager**: Teilweise defekt/unvollständig (vom User bestätigt)
- ⚠️ **Dock & App Gallery**: Keine Personalisierung, überladen, kein Launchpad-Feature
- ℹ️ **Monitoring Dashboard**: Nur Konfiguration, keine Metriken-Visualisierung

---

## 1. KRITISCH: Active Directory Domain Controller

### Status: KOMPLETT FEHLEND ❌

**Backend-Implementation:**
- **Handler:** `/backend/internal/api/handlers/ad_dc.go` (1.443 Zeilen)
- **Endpoints:** 40+ vollständig implementierte API-Endpunkte
- **Router:** Registriert in `/backend/internal/api/router.go` (Zeilen 460-540)

### Fehlende Endpoints (alle unter `/api/v1/ad-dc/*`):

#### Domain Controller Management (8 Endpoints)
- `GET /status` - GetDCStatus
- `GET /config` - GetDCConfig
- `PUT /config` - UpdateDCConfig
- `POST /provision` - ProvisionDomain
- `POST /demote` - DemoteDomain
- `GET /domain-level` - GetDomainLevel
- `POST /raise-domain-level` - RaiseDomainLevel
- `POST /restart` - RestartService

#### User Management (7 Endpoints)
- `GET /users` - ListUsers
- `POST /users` - CreateUser
- `DELETE /users/{username}` - DeleteUser
- `POST /users/{username}/enable` - EnableUser
- `POST /users/{username}/disable` - DisableUser
- `POST /users/{username}/password` - SetUserPassword
- `POST /users/{username}/expiry` - SetUserExpiry

#### Group Management (6 Endpoints)
- `GET /groups` - ListGroups
- `POST /groups` - CreateGroup
- `DELETE /groups/{groupname}` - DeleteGroup
- `GET /groups/{groupname}/members` - ListGroupMembers
- `POST /groups/{groupname}/members` - AddGroupMember
- `DELETE /groups/{groupname}/members/{username}` - RemoveGroupMember

#### Computer Management (3 Endpoints)
- `GET /computers` - ListComputers
- `POST /computers` - CreateComputer
- `DELETE /computers/{computername}` - DeleteComputer

#### Organizational Unit Management (3 Endpoints)
- `GET /ous` - ListOUs
- `POST /ous` - CreateOU
- `DELETE /ous` - DeleteOU

#### Group Policy Objects (5 Endpoints)
- `GET /gpos` - ListGPOs
- `POST /gpos` - CreateGPO
- `DELETE /gpos/{gponame}` - DeleteGPO
- `POST /gpos/{gponame}/link` - LinkGPO
- `POST /gpos/{gponame}/unlink` - UnlinkGPO

#### DNS Management (6 Endpoints)
- `GET /dns/zones` - ListDNSZones
- `POST /dns/zones` - CreateDNSZone
- `DELETE /dns/zones/{zone}` - DeleteDNSZone
- `GET /dns/zones/{zone}/records` - ListDNSRecords
- `POST /dns/zones/{zone}/records` - AddDNSRecord
- `DELETE /dns/zones/{zone}/records/{recordid}` - DeleteDNSRecord

#### FSMO Roles Management (3 Endpoints)
- `GET /fsmo` - ShowFSMORoles
- `POST /fsmo/transfer` - TransferFSMORoles
- `POST /fsmo/seize` - SeizeFSMORoles

#### Utility Functions (3 Endpoints)
- `POST /test` - TestConfiguration
- `GET /dbcheck` - ShowDBCheck
- `POST /backup` - BackupOnline

### Was fehlt im Frontend:

#### 1. API Client fehlt komplett
**Datei:** `/frontend/src/api/addc.ts` (NICHT VORHANDEN)

#### 2. Keine UI-Komponente
**Verzeichnis:** `/frontend/src/apps/ADDomainController/` (NICHT VORHANDEN)

#### 3. Nicht in Apps registriert
**Datei:** `/frontend/src/apps/index.tsx` - Keine AD DC App registriert

### Business Impact: KRITISCH

Unternehmen, die Stumpf.Works als Active Directory Domain Controller nutzen möchten, können:
- ❌ Keine Domäne provisionieren
- ❌ Keine Benutzer/Gruppen verwalten
- ❌ Keine Gruppenrichtlinien erstellen
- ❌ Keine DNS-Zonen verwalten
- ❌ Keine Computer-Konten verwalten
- ❌ Keine FSMO-Rollen transferieren

**Alle Operationen erfordern CLI/SSH-Zugriff** - inakzeptabel für eine Enterprise-NAS-Lösung.

### Empfohlene Lösung:

#### Phase 1: API Client (1 Tag)
```typescript
// frontend/src/api/addc.ts
export const addcApi = {
  // Domain Controller
  async getDCStatus(): Promise<ApiResponse<DCStatus>>
  async getDCConfig(): Promise<ApiResponse<DCConfig>>
  async provisionDomain(data: ProvisionRequest): Promise<ApiResponse<any>>
  // ... alle 40+ Methoden
}
```

#### Phase 2: UI-Komponente (2-3 Tage)
```
frontend/src/apps/ADDomainController/
├── ADDomainController.tsx (Haupt-App mit Tabs)
├── components/
│   ├── DomainStatus.tsx (Status & Provisionierung)
│   ├── UserManagement.tsx (AD-Benutzer verwalten)
│   ├── GroupManagement.tsx (AD-Gruppen verwalten)
│   ├── ComputerManagement.tsx (Computer-Konten)
│   ├── OUManagement.tsx (Organizational Units)
│   ├── GPOManagement.tsx (Gruppenrichtlinien)
│   ├── DNSManagement.tsx (DNS-Zonen & Einträge)
│   └── FSMOManagement.tsx (FSMO-Rollen)
└── index.tsx
```

#### Phase 3: Integration & Testing (1 Tag)
- App in `/frontend/src/apps/index.tsx` registrieren
- Icon: 🏢 oder 🔐
- Umfassende Tests aller Endpoints

**Gesamtaufwand:** 4-5 Tage
**Priorität:** KRITISCH

---

## 2. HOCH: Network Manager - Teilweise defekt

### Status: TEILWEISE DEFEKT ⚠️

**User-Feedback:** "Besonders der Reiter Netzwerk ist noch nicht voll funktionsfähig"

**Backend:** `/backend/internal/api/handlers/network.go` (347 Zeilen, 20+ Endpoints)
**Frontend API:** `/frontend/src/api/network.ts` (160 Zeilen) ✅
**Frontend App:** `/frontend/src/apps/NetworkManager/` (6 Komponenten, 2.442 Zeilen Code)

### Identifizierte Probleme:

#### 2.1 Interface-Konfiguration (InterfaceManager.tsx)

**Problem:** DHCP/Static IP Konfiguration möglicherweise defekt

**Backend-Endpoint:** `POST /api/v1/network/interfaces/{name}/configure`
```go
// Erwartet:
{
  "mode": "static" | "dhcp",
  "address": "192.168.1.100",  // für static
  "netmask": "255.255.255.0",  // für static
  "gateway": "192.168.1.1"     // optional
}
```

**Frontend-Code:** `/frontend/src/apps/NetworkManager/components/InterfaceManager.tsx:67-96`

**Mögliche Issues:**
- ✅ API-Call sieht korrekt aus
- ⚠️ Fehlerbehandlung könnte verbessert werden
- ⚠️ Keine Validierung der IP-Adressen (Regex)
- ⚠️ Keine Bestätigungsdialog bei kritischen Änderungen
- ⚠️ Netzwerk-Neustart nach Konfiguration fehlt möglicherweise

**Testing erforderlich:**
- Test DHCP → Static IP Wechsel
- Test Static IP → DHCP Wechsel
- Test ungültige IP-Adressen
- Test Gateway-Konfiguration
- Test Interface Up/Down nach Konfigurationsänderung

#### 2.2 Network Bonding & VLAN (NetworkConfig.tsx - "Advanced" Tab)

**Status:** Implementiert, aber ungetestet

**Backend-Endpoints:**
- `POST /api/v1/syslib/network/bond` - CreateBondInterface
- `POST /api/v1/syslib/network/vlan` - CreateVLANInterface

**Frontend:** `/frontend/src/apps/NetworkManager/components/NetworkConfig.tsx` (549 Zeilen)

**Implementierte Features:**
- ✅ Bond-Creation mit 7 Modi (balance-rr, active-backup, balance-xor, broadcast, 802.3ad, balance-tlb, balance-alb)
- ✅ VLAN-Creation mit Parent-Interface-Auswahl
- ✅ Interface-Auswahl für Bonding (min. 2 Interfaces)

**Mögliche Issues:**
- ⚠️ Keine Bond-Deletion implementiert
- ⚠️ Keine VLAN-Deletion implementiert
- ⚠️ Keine Bond-/VLAN-Bearbeitung
- ⚠️ Keine Validierung ob Interfaces bereits in Bond verwendet werden

**Backend prüfen:**
```bash
# Sind Delete-Endpoints implementiert?
grep -n "DeleteBond\|DeleteVLAN" backend/internal/api/handlers/syslib.go
```

#### 2.3 Firewall Manager (FirewallManager.tsx)

**Status:** Implementiert, könnte Verbesserungen brauchen

**Implementiert:**
- ✅ Firewall Enable/Disable
- ✅ Add/Delete Rules
- ✅ Set Default Policy
- ✅ Reset Firewall

**Mögliche Verbesserungen:**
- Regel-Priorisierung
- Regel-Bearbeitung (derzeit nur Add/Delete)
- Regel-Import/Export
- Template für häufige Regeln

#### 2.4 DNS & Routes (DNSSettings.tsx)

**Status:** Funktional ✅

**Implementiert:**
- ✅ DNS Nameserver Management
- ✅ Search Domains Management
- ✅ Routing Table anzeigen (read-only)

**Fehlende Features:**
- ❌ Statische Routen hinzufügen/löschen (Backend hat nur GET /routes)

**Backend erweitern:**
```go
// Benötigt:
POST   /api/v1/network/routes      // AddRoute
DELETE /api/v1/network/routes/{id} // DeleteRoute
```

#### 2.5 Diagnostics (DiagnosticsTool.tsx)

**Status:** Funktional ✅

**Implementiert:**
- ✅ Ping
- ✅ Traceroute
- ✅ Netstat
- ✅ Wake-on-LAN

#### 2.6 Bandwidth Monitor (BandwidthMonitor.tsx)

**Status:** Funktional ✅ (nutzt GetInterfaceStats)

### Empfohlene Maßnahmen für Network Manager:

#### Sofort (1-2 Tage):
1. ✅ InterfaceManager testen - DHCP/Static IP Wechsel
2. ✅ IP-Adress-Validierung hinzufügen
3. ✅ Bestätigungsdialoge für kritische Änderungen
4. ✅ NetworkConfig testen - Bond/VLAN Creation

#### Kurzfristig (3-5 Tage):
5. ⚠️ Bond/VLAN Deletion implementieren (Backend + Frontend)
6. ⚠️ Bond/VLAN Editing implementieren
7. ⚠️ Statische Routen Management (Backend + Frontend)
8. ⚠️ Erweiterte Firewall-Regel-Bearbeitung

**Priorität:** HOCH

---

## 3. HOCH: Dock & App Gallery UX-Verbesserungen

### Status: FUNKTIONAL, ABER UX-PROBLEME ⚠️

**Aktuelle Implementation:** `/frontend/src/layout/Dock.tsx` (80 Zeilen)
**Registrierte Apps:** `/frontend/src/apps/index.tsx` (15 Apps)

### Identifizierte Probleme:

#### 3.1 Dock ist überladen

**Aktuell:** ALLE 15 Apps werden im Dock angezeigt
```typescript
registeredApps.map((app) => <DockIcon ... />)
```

**Apps im Dock (aktuell):**
1. 📊 Dashboard
2. 💾 Storage
3. 📁 Files
4. 👥 Users
5. 🔒 Audit Logs
6. 🛡️ Security
7. 🔔 Alerts
8. 📅 Scheduled Tasks
9. 🌐 Network
10. 🐳 Docker
11. 🔌 Plugins
12. 🛒 App Store
13. 💻 Terminal
14. 💾 Backups (DUPLIZIERTES ICON wie Storage!)
15. ⚙️ Settings

**Problem:**
- ❌ Zu viele Icons → unübersichtlich
- ❌ Keine Personalisierung möglich
- ❌ Duplizierte Icons (Storage und Backups beide 💾)
- ❌ Keine logische Gruppierung

#### 3.2 Keine Rechtsklick-Funktionalität

**Aktuell:** Dock-Icons haben nur onClick

**Fehlende Features:**
- ❌ Rechtsklick-Menü
- ❌ "Aus Dock entfernen"
- ❌ "Optionen" / "Einstellungen"
- ❌ "Im Dock behalten" Toggle

#### 3.3 Keine App Gallery / Launchpad

**Fehlend:** Kein macOS-ähnliches Launchpad

**Was fehlt:**
- ❌ Vollbild-App-Übersicht (wie macOS Launchpad)
- ❌ Drag & Drop Apps ins Dock
- ❌ App-Kategorien (System, Management, Tools, etc.)
- ❌ App-Suche

### Empfohlene Lösung:

#### Phase 1: Dock-Personalisierung (1-2 Tage)

**1.1 User Dock Preferences Store**
```typescript
// frontend/src/store/dockStore.ts
interface DockStore {
  dockApps: string[];  // Array von App-IDs
  addToDock: (appId: string) => void;
  removeFromDock: (appId: string) => void;
  reorderDock: (from: number, to: number) => void;
}

// Default Dock Apps (z.B. nur 7-8 wichtigste)
const defaultDockApps = [
  'dashboard',
  'files',
  'storage',
  'network',
  'docker',
  'terminal',
  'settings'
];
```

**1.2 Rechtsklick-Menü implementieren**
```typescript
// In DockIcon.tsx
<ContextMenu
  onRemoveFromDock={() => dockStore.removeFromDock(app.id)}
  onOptions={() => openAppSettings(app.id)}
/>
```

**1.3 Dock Icons differenzieren**
```typescript
// Backups Icon ändern von 💾 zu ⏱️ oder 💿
{
  id: 'backups',
  name: 'Backups',
  icon: '⏱️',  // GEÄNDERT: war 💾 (Konflikt mit Storage)
  ...
}
```

#### Phase 2: App Gallery / Launchpad (2-3 Tage)

**2.1 App Gallery Komponente**
```typescript
// frontend/src/components/AppGallery.tsx
interface AppGalleryProps {
  isOpen: boolean;
  onClose: () => void;
}

// Features:
// - Vollbild-Overlay mit Grid-Layout
// - Kategorien: System, Management, Tools, Development
// - Suchfunktion
// - Drag & Drop zu Dock
// - Click to Launch
```

**2.2 Kategorisierung**
```typescript
// frontend/src/apps/index.tsx
export const appCategories = {
  system: ['dashboard', 'settings', 'terminal'],
  management: ['users', 'network', 'storage'],
  security: ['security', 'audit-logs', 'alerts'],
  tools: ['files', 'backups', 'tasks'],
  development: ['docker', 'plugins', 'app-store']
};
```

**2.3 Launchpad-Trigger**
```typescript
// Launcher Icon im Dock (linke Seite)
<DockIcon
  app={{
    id: 'launchpad',
    name: 'App Gallery',
    icon: '⊞',  // oder '🚀'
    ...
  }}
  onClick={() => setShowAppGallery(true)}
/>
```

**2.4 App Gallery Layout**
```
┌─────────────────────────────────────────┐
│         🚀 App Gallery                  │
├─────────────────────────────────────────┤
│  [Suche...]                             │
├─────────────────────────────────────────┤
│                                         │
│  System                                 │
│  [📊] [⚙️] [💻]                         │
│                                         │
│  Management                             │
│  [👥] [🌐] [💾] [💿]                    │
│                                         │
│  Security                               │
│  [🛡️] [🔒] [🔔]                        │
│                                         │
│  Tools                                  │
│  [📁] [⏱️] [📅] [🐳] [🔌] [🛒]        │
│                                         │
└─────────────────────────────────────────┘
```

#### Phase 3: Erweiterte Features (1-2 Tage)

**3.1 Dock-Ordner (Stacks)**
```typescript
// Gruppierung ähnlicher Apps
{
  id: 'system-stack',
  name: 'System',
  icon: '📦',
  type: 'stack',
  apps: ['audit-logs', 'security', 'alerts', 'tasks']
}
```

**3.2 Persistierung**
```typescript
// Dock-Preferences im Backend speichern
POST /api/v1/users/me/preferences
{
  "dock_apps": ["dashboard", "files", ...],
  "dock_order": [0, 1, 2, ...]
}
```

**3.3 Drag & Drop Reordering**
```typescript
// react-beautiful-dnd oder dnd-kit
<Droppable droppableId="dock">
  {dockApps.map((app, index) => (
    <Draggable key={app.id} draggableId={app.id} index={index}>
      <DockIcon ... />
    </Draggable>
  ))}
</Droppable>
```

### Logische Gruppierungsvorschläge:

#### Standard-Dock (8 Apps):
1. 📊 Dashboard
2. 📁 Files
3. 💾 Storage
4. 🌐 Network
5. 🐳 Docker
6. 💻 Terminal
7. ⚙️ Settings
8. 🚀 App Gallery (Launchpad)

#### Erweitert (Power User):
- Zusätzlich: Users, Security, Backups, Tasks

**Priorität:** HOCH
**Gesamtaufwand:** 5-7 Tage

---

## 4. MITTEL: Monitoring Dashboard fehlt

### Status: NUR KONFIGURATION, KEINE VISUALISIERUNG ⚠️

**Backend:** `/backend/internal/api/handlers/monitoring.go`
**Frontend API:** `/frontend/src/api/monitoring.ts` ✅
**Frontend UI:** `/frontend/src/apps/Settings/sections/MonitoringSection.tsx` (nur Config)

### Was existiert:

#### Backend:
- ✅ `GET /api/v1/monitoring/config` - GetMonitoringConfig
- ✅ `PUT /api/v1/monitoring/config` - UpdateMonitoringConfig
- ✅ `GET /metrics` - Prometheus Metrics Endpoint

#### Frontend:
- ✅ Monitoring-Konfiguration in Settings
- ✅ Enable/Disable Monitoring
- ❌ **KEINE Metriken-Visualisierung**

### Was fehlt:

#### 1. Metriken-Dashboard App fehlt komplett
```
frontend/src/apps/Monitoring/
├── MonitoringDashboard.tsx
├── components/
│   ├── MetricsOverview.tsx
│   ├── CPUMetrics.tsx
│   ├── MemoryMetrics.tsx
│   ├── DiskMetrics.tsx
│   ├── NetworkMetrics.tsx
│   └── HealthScore.tsx
└── index.tsx
```

#### 2. Real-time Metrics API fehlt
**Backend erweitern:**
```go
GET /api/v1/monitoring/metrics/current  // Aktuelle Metriken als JSON
GET /api/v1/monitoring/metrics/history  // Historische Daten
GET /api/v1/monitoring/health           // Health Score
```

#### 3. Grafana/Prometheus Integration

**Optionen:**
- **Option A:** Eingebettete Grafana-Dashboards (iframe)
- **Option B:** Eigene Chart-Library (recharts, visx, chart.js)
- **Option C:** Hybrid: Eigene Übersicht + Grafana für Details

### Empfohlene Lösung:

#### Phase 1: Metriken-API (1 Tag)
```go
// backend/internal/api/handlers/monitoring.go
func GetCurrentMetrics(w http.ResponseWriter, r *http.Request)
func GetMetricsHistory(w http.ResponseWriter, r *http.Request)
func GetHealthScore(w http.ResponseWriter, r *http.Request)
```

#### Phase 2: Frontend Dashboard (2-3 Tage)
- Real-time Metriken mit Auto-Refresh
- Charts für CPU, RAM, Disk, Network
- Health Score prominent anzeigen
- Alert-Integration (Warnung bei kritischen Werten)

**Priorität:** MITTEL
**Gesamtaufwand:** 3-4 Tage

---

## 5. Weitere Findings

### 5.1 Icon-Duplikate

**Problem:** Storage (💾) und Backups (💾) haben dasselbe Icon

**Lösung:**
```typescript
// Backups Icon ändern
{
  id: 'backups',
  icon: '⏱️',  // oder 💿, 📼, ⏮️
}
```

### 5.2 Terminal - Simulation Mode

**Status:** Terminal.tsx ist im "Simulation Mode"

**Datei:** `/frontend/src/apps/Terminal/Terminal.tsx`

**Problem:**
- ⚠️ WebSocket-Verbindung könnte defekt sein
- ⚠️ Nur Simulation ohne echte Backend-Verbindung?

**Testing erforderlich:**
- Echte Terminal-Session testen
- WebSocket-Verbindung prüfen
- TTY-Allokation prüfen

### 5.3 Samba Share - Doppelte Implementation?

**Potentielles Problem:**
- `handlers/syslib.go` hat Samba-Operationen
- `handlers/storage.go` könnte auch Samba haben

**Prüfen:**
```bash
grep -n "Samba\|smb" backend/internal/api/handlers/storage.go
```

---

## Zusammenfassung & Priorisierung

### TIER 1 - KRITISCH (Sofort angehen)

| Feature | Aufwand | Impact | Dateien |
|---------|---------|--------|---------|
| **AD Domain Controller UI** | 4-5 Tage | KRITISCH | `frontend/src/api/addc.ts`, `frontend/src/apps/ADDomainController/` |
| **Network Manager Fixes** | 2-3 Tage | HOCH | `frontend/src/apps/NetworkManager/components/InterfaceManager.tsx` |
| **Dock Personalisierung** | 2-3 Tage | HOCH | `frontend/src/layout/Dock.tsx`, `frontend/src/store/dockStore.ts` |

**Gesamt TIER 1:** 8-11 Tage

### TIER 2 - HOCH (Nächste Iteration)

| Feature | Aufwand | Impact |
|---------|---------|--------|
| **App Gallery / Launchpad** | 2-3 Tage | HOCH |
| **Monitoring Dashboard** | 3-4 Tage | MITTEL |
| **Bond/VLAN Deletion** | 1-2 Tage | MITTEL |
| **Static Routes Management** | 1-2 Tage | MITTEL |

**Gesamt TIER 2:** 7-11 Tage

### TIER 3 - MITTEL (Backlog)

- Terminal WebSocket Testing & Fix (1-2 Tage)
- Erweiterte Firewall-Regel-Bearbeitung (1-2 Tage)
- Samba Duplication Audit (1 Tag)
- Icon-Duplikate beheben (0.5 Tage)
- Dock Stacks/Ordner (2-3 Tage)

---

## Detaillierte Datei-Referenzen

### Backend Handler
```
/home/user/stumpfworks-nas/backend/internal/api/handlers/
├── ad_dc.go (1.443 Zeilen) ⚠️ KEINE UI
├── network.go (347 Zeilen) ⚠️ TEILWEISE DEFEKT
├── syslib.go (Bonding/VLAN implementiert)
├── monitoring.go (nur Config, keine Metriken-API)
└── router.go (Route-Definitionen)
```

### Frontend API Clients
```
/home/user/stumpfworks-nas/frontend/src/api/
├── addc.ts ❌ FEHLT KOMPLETT
├── network.ts ✅ (160 Zeilen)
├── syslib.ts ✅ (276 Zeilen)
└── monitoring.ts ✅ (nur Config)
```

### Frontend Apps
```
/home/user/stumpfworks-nas/frontend/src/apps/
├── ADDomainController/ ❌ FEHLT KOMPLETT
├── NetworkManager/ ⚠️ (2.442 Zeilen, teilweise defekt)
│   ├── NetworkManager.tsx
│   ├── components/
│   │   ├── InterfaceManager.tsx (379 Zeilen) ⚠️
│   │   ├── NetworkConfig.tsx (549 Zeilen) ⚠️
│   │   ├── DNSSettings.tsx (328 Zeilen) ✅
│   │   ├── FirewallManager.tsx (449 Zeilen) ✅
│   │   ├── DiagnosticsTool.tsx (324 Zeilen) ✅
│   │   └── BandwidthMonitor.tsx (413 Zeilen) ✅
├── Monitoring/ ❌ FEHLT (nur Settings-Section)
└── index.tsx (15 Apps registriert)
```

### Layout & UX
```
/home/user/stumpfworks-nas/frontend/src/layout/
├── Dock.tsx ⚠️ (80 Zeilen - keine Personalisierung)
└── AppGallery.tsx ❌ FEHLT KOMPLETT
```

---

## Empfohlene Reihenfolge

### Sprint 1 (Woche 1): Kritische Basics
1. Tag 1-2: Network Manager Fixes + Testing
2. Tag 3-4: Dock Personalisierung (Rechtsklick, Remove)
3. Tag 5: Icon-Duplikate, Testing, Bug Fixes

### Sprint 2 (Woche 2): AD DC Foundation
1. Tag 1: AD DC API Client (`addc.ts`)
2. Tag 2-3: AD DC UI - Domain Management & User Management
3. Tag 4-5: AD DC UI - Groups, Computers, OUs

### Sprint 3 (Woche 3): AD DC Complete + UX
1. Tag 1-2: AD DC UI - GPO, DNS, FSMO
2. Tag 3-4: App Gallery / Launchpad
3. Tag 5: Integration Testing, Polish

### Sprint 4 (Woche 4): Monitoring & Network Advanced
1. Tag 1-2: Monitoring Dashboard
2. Tag 3-4: Bond/VLAN Deletion, Static Routes
3. Tag 5: Release Testing

**Gesamt:** ~4 Wochen für vollständige Integration aller kritischen Features

---

## Testing Checkliste

### AD Domain Controller
- [ ] Domain provisionieren
- [ ] Benutzer erstellen/löschen/deaktivieren
- [ ] Gruppen erstellen und Mitglieder verwalten
- [ ] Computer-Konten verwalten
- [ ] OUs erstellen und löschen
- [ ] GPOs erstellen, linken, unlinken
- [ ] DNS-Zonen und Records verwalten
- [ ] FSMO-Rollen anzeigen und transferieren

### Network Manager
- [ ] Interface DHCP → Static IP
- [ ] Interface Static IP → DHCP
- [ ] Interface Up/Down
- [ ] DNS-Server ändern
- [ ] Firewall-Regeln hinzufügen/löschen
- [ ] Bond mit 2 Interfaces erstellen (alle 7 Modi)
- [ ] VLAN erstellen
- [ ] Ping/Traceroute/Netstat
- [ ] Wake-on-LAN

### Dock & App Gallery
- [ ] Rechtsklick → "Aus Dock entfernen"
- [ ] App Gallery öffnen
- [ ] App per Drag & Drop ins Dock
- [ ] Dock-Reihenfolge ändern
- [ ] Dock-Preferences persistieren
- [ ] Nach Reload: Dock-Einstellungen bleiben

### Monitoring
- [ ] Monitoring aktivieren/deaktivieren
- [ ] Live-Metriken anzeigen
- [ ] Charts für CPU/RAM/Disk/Network
- [ ] Health Score berechnen
- [ ] Prometheus /metrics Endpoint

---

## Schlussfolgerung

**Hauptprobleme:**
1. ❌ **Active Directory DC**: Komplette UI fehlt (40+ Endpoints ungenutzt)
2. ⚠️ **Network Manager**: Teilweise defekt, unvollständig
3. ⚠️ **Dock/UX**: Keine Personalisierung, überladen
4. ℹ️ **Monitoring**: Nur Config, keine Visualisierung

**Gesamtaufwand für vollständige Integration:**
- TIER 1 (Kritisch): 8-11 Tage
- TIER 2 (Hoch): 7-11 Tage
- **GESAMT: ~3-4 Wochen**

**Nächste Schritte:**
1. Network Manager Fixes testen und deployen (2-3 Tage)
2. Dock Personalisierung implementieren (2-3 Tage)
3. AD DC vollständige UI entwickeln (4-5 Tage)
4. App Gallery / Launchpad (2-3 Tage)
5. Monitoring Dashboard (3-4 Tage)

---

**Erstellt von:** Claude
**Branch:** `claude/audit-frontend-integration-017ykh14kLotgaQ7URfFanDX`
**Repository:** Stumpf-works/stumpfworks-nas
