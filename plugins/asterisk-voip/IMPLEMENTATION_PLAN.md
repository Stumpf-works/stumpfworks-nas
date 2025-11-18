# Asterisk VoIP Plugin - Implementation Plan

Detaillierter Plan zur vollständigen Implementierung des Asterisk VoIP Plugins für StumpfWorks NAS.

## 📋 Inhaltsverzeichnis

- [Übersicht](#übersicht)
- [Phase 1: Proof of Concept](#phase-1-proof-of-concept-poc-✅-erledigt)
- [Phase 2: Core Features](#phase-2-core-features)
- [Phase 3: Advanced Features](#phase-3-advanced-features)
- [Phase 4: UI/UX](#phase-4-uiux)
- [Phase 5: Production Readiness](#phase-5-production-readiness)
- [Technische Details](#technische-details)
- [Timeline](#timeline)

---

## 🎯 Übersicht

### Projektziele

1. **Vollständige VoIP-Telefonanlage** in StumpfWorks NAS integrieren
2. **Benutzerfreundliche UI** für nicht-technische Benutzer
3. **Enterprise-Features** (IVR, Call Recording, Conferencing)
4. **WebRTC-Integration** für Browser-basierte Telefonie
5. **Multi-Tenant-Fähigkeit** für mehrere Organisationen

### Technologie-Stack

**Backend:**
- Asterisk 20 (PBX Core)
- Go 1.21+ (Management Backend)
- AMI (Asterisk Manager Interface)
- ARI (Asterisk REST Interface) für erweiterte Features
- SQLite (Konfigurationsdatenbank)

**Frontend:**
- React 18+ (UI Framework)
- TypeScript (Typsicherheit)
- TailwindCSS (Styling)
- WebRTC (Browser-Telefonie)
- WebSocket (Echtzeit-Updates)

**Infrastructure:**
- Docker & Docker Compose
- STUN/TURN Server (NAT-Traversal)
- Let's Encrypt (TLS/SRTP)

---

## Phase 1: Proof of Concept (PoC) ✅ ERLEDIGT

### Status: **Abgeschlossen**

### Deliverables

- ✅ Plugin-Manifest (`plugin.json`)
- ✅ Docker Compose Setup
- ✅ Asterisk Basis-Konfiguration
  - ✅ `sip.conf` - SIP Configuration
  - ✅ `extensions.conf` - Dialplan
  - ✅ `voicemail.conf` - Voicemail
  - ✅ `manager.conf` - AMI
  - ✅ `http.conf` - WebRTC/WebSocket
  - ✅ `rtp.conf` - Media Transport
- ✅ Go Backend
  - ✅ AMI Client Implementation
  - ✅ REST API Server
  - ✅ Basic API Endpoints
  - ✅ Dockerfile
- ✅ Dokumentation
  - ✅ README
  - ✅ Implementierungsplan

### Was funktioniert im PoC

- Asterisk Container startet und läuft
- AMI-Verbindung wird hergestellt
- REST API ist erreichbar
- Basis-Endpoints für Extensions, Trunks, Calls verfügbar

### Nächste Schritte

➡️ **Phase 2**: Core Features vollständig implementieren

---

## Phase 2: Core Features

### Dauer: 3-4 Wochen

### 2.1 Configuration Management (Woche 1)

#### Ziele
- Asterisk-Konfigurationsdateien programmatisch verwalten
- Änderungen persistent speichern
- Automatisches Reload nach Änderungen

#### Tasks

**Config Parser/Writer**
```go
// backend/config/parser.go
type SIPConfigParser struct {
    FilePath string
}

func (p *SIPConfigParser) Parse() (*SIPConfig, error)
func (p *SIPConfigParser) Write(config *SIPConfig) error
func (p *SIPConfigParser) Reload() error
```

**Implementierung:**
1. INI-Parser für Asterisk-Configs (gopkg.in/ini.v1)
2. CRUD-Operationen für Peers/Trunks/Extensions
3. Validierung vor dem Schreiben
4. Atomic File Updates (write temp, then rename)
5. Automatisches Backup vor Änderungen

**API Endpoints:**
- `POST /api/v1/config/validate` - Konfiguration validieren
- `POST /api/v1/config/backup` - Backup erstellen
- `POST /api/v1/config/restore` - Backup wiederherstellen

#### Deliverables
- [ ] INI Config Parser/Writer
- [ ] SIP Config Management
- [ ] Extensions Config Management
- [ ] Voicemail Config Management
- [ ] Automatic Reload on Changes
- [ ] Config Backup/Restore

---

### 2.2 Extensions Management (Woche 2)

#### Ziele
- Vollständiges CRUD für SIP Extensions
- Echtzeit-Status (online/offline/busy)
- Integration mit StumpfWorks User Management

#### Tasks

**Extension CRUD:**
```go
type ExtensionManager struct {
    configPath string
    amiClient  *ami.Client
}

func (m *ExtensionManager) Create(ext *Extension) error
func (m *ExtensionManager) Update(id string, ext *Extension) error
func (m *ExtensionManager) Delete(id string) error
func (m *ExtensionManager) Get(id string) (*Extension, error)
func (m *ExtensionManager) List() ([]*Extension, error)
func (m *ExtensionManager) GetStatus(id string) (*ExtensionStatus, error)
```

**Features:**
- Auto-generate strong SIP passwords
- CallerID Management
- Mailbox Assignment
- Context Assignment
- Codec Selection per Extension
- NAT Settings

**AMI Integration:**
- `SIPpeers` - Liste aller Extensions
- `SIPshowpeer` - Details zu Extension
- `SIPnotify` - Notifications senden

**StumpfWorks Integration:**
- Extension-User-Mapping in SQLite
- Sync Extensions mit StumpfWorks Users
- Permission-Based Access

#### Deliverables
- [ ] Extension CRUD Implementation
- [ ] Real-time Status via AMI
- [ ] StumpfWorks User Integration
- [ ] Password Generation
- [ ] Bulk Import/Export (CSV)

---

### 2.3 Trunk Management (Woche 2)

#### Ziele
- SIP Trunks zu VoIP-Providern
- Trunk-Status-Monitoring
- Failover zwischen Trunks

#### Tasks

**Trunk Types:**
1. **SIP Registration** (bei Provider registrieren)
2. **SIP Peer** (Provider connectet zu uns)
3. **SIP Friend** (bidirektional)

**Features:**
- Multiple Trunks pro Provider
- Prefix-based Routing
- Least-Cost Routing (LCR)
- Failover/Load Balancing

**Monitoring:**
- Registration Status
- Latency/Jitter
- Call Quality Metrics
- Usage Statistics

#### Deliverables
- [ ] Trunk CRUD Implementation
- [ ] Registration Monitoring
- [ ] Outbound Route Management
- [ ] Inbound Route Management
- [ ] Trunk Failover Logic

---

### 2.4 Dialplan Management (Woche 3)

#### Ziele
- Dialplan visuell konfigurieren
- Extension Routing
- Outbound Routing

#### Tasks

**Dialplan Builder:**
```json
{
  "context": "internal",
  "rules": [
    {
      "pattern": "_1XXX",
      "actions": [
        {"type": "Dial", "target": "SIP/${EXTEN}", "timeout": 30},
        {"type": "VoiceMail", "box": "${EXTEN}@default"}
      ]
    }
  ]
}
```

**Features:**
- Pattern Matching (_XXXX, _0., etc.)
- Time-based Routing
- Geographic Routing
- Custom Variables
- Goto/GotoIf Logic

**Visual Editor (Phase 4):**
- Drag-and-drop Dialplan Builder
- Flow-Chart-Darstellung

#### Deliverables
- [ ] JSON-based Dialplan Definition
- [ ] Dialplan Generator (JSON → extensions.conf)
- [ ] Dialplan Validator
- [ ] Time-based Routing
- [ ] Emergency Number Handling

---

### 2.5 Call Management (Woche 3-4)

#### Ziele
- Active Calls anzeigen
- Calls steuern (Hangup, Transfer, etc.)
- Call History/CDR

#### Tasks

**Active Calls:**
- Real-time Call List via AMI Events
- Call Details (Caller, Callee, Duration)
- Call State (Ringing, Answered, etc.)

**Call Control:**
- Hangup
- Transfer (Blind/Attended)
- Park
- Listen/Whisper/Barge (Supervision)

**Call Detail Records (CDR):**
```sql
CREATE TABLE cdr (
    id INTEGER PRIMARY KEY,
    call_id TEXT,
    caller_id TEXT,
    callee_id TEXT,
    start_time INTEGER,
    answer_time INTEGER,
    end_time INTEGER,
    duration INTEGER,
    billsec INTEGER,
    disposition TEXT, -- ANSWERED, NO ANSWER, BUSY, FAILED
    recording_path TEXT
);
```

#### Deliverables
- [ ] Active Calls List (WebSocket)
- [ ] Call Control Functions
- [ ] CDR Database Schema
- [ ] CDR Import from Asterisk
- [ ] Call History API
- [ ] Call Statistics

---

### 2.6 Voicemail (Woche 4)

#### Ziele
- Voicemail Boxes verwalten
- Messages abrufen
- Email Notifications

#### Tasks

**Voicemail Management:**
- Create/Update/Delete Mailboxes
- List Messages per Mailbox
- Play/Download Messages
- Delete Messages
- Mark as Read/Unread

**Email Integration:**
- SMTP Configuration
- Attach WAV file to email
- Transcription (future)

**Features:**
- Greeting Messages (Unavailable/Busy)
- PIN Protection
- MWI (Message Waiting Indicator)

#### Deliverables
- [ ] Voicemail Box CRUD
- [ ] Message Management API
- [ ] Audio File Streaming
- [ ] Email Notifications
- [ ] MWI via SIP NOTIFY

---

## Phase 3: Advanced Features

### Dauer: 4-6 Wochen

### 3.1 Call Recording (Woche 5)

#### Features
- On-demand Recording (via DTMF *1)
- Automatic Recording per Extension/Trunk
- Recording Storage Management
- Playback/Download via API
- Retention Policies

#### Implementation
```go
type RecordingManager struct {
    storagePath string
}

func (m *RecordingManager) List(filters RecordingFilters) ([]*Recording, error)
func (m *RecordingManager) Get(id string) (*Recording, error)
func (m *RecordingManager) Stream(id string) (io.ReadCloser, error)
func (m *RecordingManager) Delete(id string) error
func (m *RecordingManager) ApplyRetentionPolicy() error
```

#### Deliverables
- [ ] Recording Management API
- [ ] File Storage (with compression)
- [ ] Streaming Endpoint
- [ ] Retention Policy Engine
- [ ] Disk Usage Monitoring

---

### 3.2 IVR (Interactive Voice Response) (Woche 6)

#### Features
- Multi-level Menus
- Text-to-Speech (TTS)
- Audio File Upload
- Time-based IVR
- Language Selection

#### IVR Builder
```json
{
  "id": "main-ivr",
  "name": "Main Menu",
  "greeting": "welcome.wav",
  "timeout": 5,
  "options": [
    {"digit": "1", "action": "dial", "target": "SIP/sales"},
    {"digit": "2", "action": "dial", "target": "SIP/support"},
    {"digit": "0", "action": "dial", "target": "SIP/operator"},
    {"digit": "*", "action": "ivr", "target": "sub-menu"}
  ],
  "invalid": "invalid.wav",
  "timeout_action": "voicemail"
}
```

#### Deliverables
- [ ] IVR Definition Schema
- [ ] IVR to Dialplan Generator
- [ ] Audio File Management
- [ ] TTS Integration (Google/AWS/Azure)
- [ ] Visual IVR Builder (Frontend)

---

### 3.3 Conference Rooms (Woche 7)

#### Features
- Conference Room Creation
- PIN Protection
- Moderator Controls
- Recording
- Max Participants Limit

#### Implementation
```asterisk
[conference-8000]
exten => 8000,1,ConfBridge(8000,default_bridge,default_user)
```

**ConfBridge Configuration:**
- Admin/User Roles
- Mute/Unmute
- Kick Participants
- Announce Join/Leave

#### Deliverables
- [ ] Conference Room Management
- [ ] Real-time Participant List
- [ ] Moderator Controls API
- [ ] Conference Recording
- [ ] WebRTC Conference Support

---

### 3.4 Call Queue (Woche 8)

#### Features
- Queue Management
- Agent Login/Logout
- Queue Statistics
- Call Distribution Strategies
- Overflow Handling

#### Queue Strategies
- Ring All
- Round Robin
- Least Recent
- Fewest Calls
- Random

#### Deliverables
- [ ] Queue CRUD
- [ ] Agent Management
- [ ] Real-time Queue Status
- [ ] Queue Statistics/Reports
- [ ] Overflow/Timeout Handling

---

### 3.5 WebRTC Integration (Woche 9-10)

#### Ziele
- Browser-basierte Softphone
- Click-to-Call
- Video Calls (optional)

#### Architecture
```
Browser (WebRTC) ←→ WebSocket ←→ Asterisk (chan_pjsip + res_http_websocket)
```

#### Implementation

**Asterisk pjsip.conf:**
```ini
[webrtc-transport]
type=transport
protocol=wss
bind=0.0.0.0:8089

[webrtc-template](!)
type=endpoint
webrtc=yes
context=internal
dtls_auto_generate_cert=yes
```

**Frontend (React):**
```typescript
import { SIPSession } from 'sip.js';

class WebRTCPhone {
  connect(server: string, user: string, password: string)
  call(extension: string)
  answer()
  hangup()
  mute()
  hold()
}
```

#### Deliverables
- [ ] PJSIP WebRTC Configuration
- [ ] WebSocket Endpoint
- [ ] TLS/SRTP Setup
- [ ] React WebRTC Component
- [ ] SIP.js Integration
- [ ] Video Call Support (optional)

---

## Phase 4: UI/UX

### Dauer: 4-5 Wochen

### 4.1 Dashboard (Woche 11)

#### Features
- System Status
- Active Calls Counter
- Extension Status Grid
- Recent Calls
- Call Statistics Charts

#### Components
```tsx
<Dashboard>
  <StatusCard title="Active Calls" value={5} />
  <ExtensionGrid extensions={extensions} />
  <CallChart data={callStats} />
  <RecentCalls calls={recentCalls} />
</Dashboard>
```

#### Deliverables
- [ ] Dashboard Layout
- [ ] Real-time Updates (WebSocket)
- [ ] Charts (Recharts/Chart.js)
- [ ] Responsive Design

---

### 4.2 Extensions UI (Woche 12)

#### Features
- Extension List (Table)
- Create/Edit Modal
- Status Indicators
- Quick Actions (Call, SMS)
- Bulk Operations

#### Deliverables
- [ ] Extension Table Component
- [ ] Extension Form
- [ ] Status Badges
- [ ] Search/Filter
- [ ] CSV Import/Export UI

---

### 4.3 Call Management UI (Woche 13)

#### Features
- Active Calls List
- Call History
- Call Details Modal
- Call Control Buttons

#### Deliverables
- [ ] Active Calls Component
- [ ] Call History Table
- [ ] Call Player (for recordings)
- [ ] Call Control UI

---

### 4.4 Settings UI (Woche 14)

#### Features
- General Settings
- SIP Settings
- Codec Configuration
- Email Settings
- Backup/Restore

#### Deliverables
- [ ] Settings Tabs
- [ ] Form Validation
- [ ] Save/Reset Buttons
- [ ] Backup/Restore UI

---

### 4.5 Visual Dialplan/IVR Builder (Woche 15)

#### Features
- Drag-and-drop Flow Builder
- Node Types (Dial, Voicemail, IVR, etc.)
- Connections/Routing
- Preview/Test

#### Libraries
- React Flow / ReactFlow
- DnD Kit

#### Deliverables
- [ ] Flow Builder Component
- [ ] Node Palette
- [ ] Connection Logic
- [ ] Export to Dialplan

---

## Phase 5: Production Readiness

### Dauer: 3-4 Wochen

### 5.1 Security (Woche 16)

#### Tasks
- [ ] TLS/SRTP für SIP
- [ ] WebSocket WSS
- [ ] Fail2Ban Integration
- [ ] SIP Authentication
- [ ] Rate Limiting
- [ ] IP Whitelisting/Blacklisting
- [ ] Audit Logging

---

### 5.2 Monitoring & Logging (Woche 17)

#### Tasks
- [ ] Asterisk Logs Collection
- [ ] AMI Event Logging
- [ ] Metrics (Prometheus)
- [ ] Alerting (Email/Slack/SMS)
- [ ] Health Checks
- [ ] Performance Monitoring

---

### 5.3 Backup & Restore (Woche 18)

#### Tasks
- [ ] Configuration Backup
- [ ] Database Backup
- [ ] Voicemail Backup
- [ ] Recording Backup
- [ ] Automated Backups (Cron)
- [ ] Restore Wizard

---

### 5.4 Documentation (Woche 19)

#### Tasks
- [ ] User Manual
- [ ] Administrator Guide
- [ ] API Documentation (OpenAPI/Swagger)
- [ ] Troubleshooting Guide
- [ ] Video Tutorials
- [ ] FAQ

---

### 5.5 Testing (Woche 19)

#### Tasks
- [ ] Unit Tests (Backend)
- [ ] Integration Tests
- [ ] E2E Tests (Frontend)
- [ ] Load Testing (SIPp)
- [ ] Security Audit
- [ ] Penetration Testing

---

## 📊 Technische Details

### Database Schema

**SQLite Database:** `/data/asterisk-manager.db`

```sql
-- Extensions
CREATE TABLE extensions (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    secret TEXT NOT NULL,
    context TEXT DEFAULT 'internal',
    caller_id TEXT,
    mailbox TEXT,
    user_id TEXT, -- StumpfWorks User ID
    created_at INTEGER,
    updated_at INTEGER
);

-- Trunks
CREATE TABLE trunks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- peer, friend, user
    host TEXT NOT NULL,
    port INTEGER DEFAULT 5060,
    username TEXT,
    secret TEXT,
    context TEXT DEFAULT 'from-trunk',
    from_domain TEXT,
    created_at INTEGER,
    updated_at INTEGER
);

-- Call Detail Records
CREATE TABLE cdr (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    call_id TEXT UNIQUE,
    caller_id TEXT,
    caller_name TEXT,
    callee_id TEXT,
    trunk_id TEXT,
    start_time INTEGER,
    answer_time INTEGER,
    end_time INTEGER,
    duration INTEGER,
    billsec INTEGER,
    disposition TEXT,
    recording_path TEXT
);

-- IVR Definitions
CREATE TABLE ivr_menus (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    greeting_file TEXT,
    timeout INTEGER DEFAULT 5,
    invalid_file TEXT,
    definition JSON, -- Full IVR tree
    created_at INTEGER,
    updated_at INTEGER
);

-- Voicemail Boxes
CREATE TABLE voicemail_boxes (
    id TEXT PRIMARY KEY,
    context TEXT DEFAULT 'default',
    pin TEXT NOT NULL,
    name TEXT,
    email TEXT,
    extension_id TEXT,
    created_at INTEGER,
    updated_at INTEGER
);

-- Conference Rooms
CREATE TABLE conference_rooms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    pin TEXT,
    max_participants INTEGER,
    record BOOLEAN DEFAULT FALSE,
    created_at INTEGER,
    updated_at INTEGER
);

-- Call Queues
CREATE TABLE call_queues (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    strategy TEXT DEFAULT 'ringall',
    timeout INTEGER DEFAULT 300,
    max_callers INTEGER DEFAULT 0,
    created_at INTEGER,
    updated_at INTEGER
);

-- Queue Members
CREATE TABLE queue_members (
    queue_id TEXT,
    extension_id TEXT,
    penalty INTEGER DEFAULT 0,
    paused BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (queue_id, extension_id)
);
```

---

### API Endpoints (Complete)

#### Status & Info
- `GET /health` - Health check
- `GET /version` - Asterisk version
- `GET /api/v1/status` - System status
- `GET /api/v1/ami/status` - AMI connection status

#### Extensions
- `GET /api/v1/extensions` - List extensions
- `POST /api/v1/extensions` - Create extension
- `GET /api/v1/extensions/:id` - Get extension
- `PUT /api/v1/extensions/:id` - Update extension
- `DELETE /api/v1/extensions/:id` - Delete extension
- `GET /api/v1/extensions/:id/status` - Real-time status

#### Trunks
- `GET /api/v1/trunks` - List trunks
- `POST /api/v1/trunks` - Create trunk
- `GET /api/v1/trunks/:id` - Get trunk
- `PUT /api/v1/trunks/:id` - Update trunk
- `DELETE /api/v1/trunks/:id` - Delete trunk
- `GET /api/v1/trunks/:id/status` - Registration status

#### Calls
- `GET /api/v1/calls` - List active calls
- `POST /api/v1/calls/originate` - Originate call
- `POST /api/v1/calls/:id/hangup` - Hangup call
- `POST /api/v1/calls/:id/transfer` - Transfer call
- `POST /api/v1/calls/:id/park` - Park call
- `GET /api/v1/cdr` - Call history (CDR)

#### Voicemail
- `GET /api/v1/voicemail/boxes` - List boxes
- `POST /api/v1/voicemail/boxes` - Create box
- `GET /api/v1/voicemail/boxes/:id/messages` - List messages
- `GET /api/v1/voicemail/messages/:id` - Get message audio
- `DELETE /api/v1/voicemail/messages/:id` - Delete message

#### Recordings
- `GET /api/v1/recordings` - List recordings
- `GET /api/v1/recordings/:id` - Stream recording
- `DELETE /api/v1/recordings/:id` - Delete recording

#### IVR
- `GET /api/v1/ivr` - List IVR menus
- `POST /api/v1/ivr` - Create IVR
- `GET /api/v1/ivr/:id` - Get IVR
- `PUT /api/v1/ivr/:id` - Update IVR
- `DELETE /api/v1/ivr/:id` - Delete IVR

#### Conference
- `GET /api/v1/conference/rooms` - List rooms
- `POST /api/v1/conference/rooms` - Create room
- `GET /api/v1/conference/rooms/:id/participants` - List participants
- `POST /api/v1/conference/rooms/:id/kick/:participant` - Kick participant

#### Queues
- `GET /api/v1/queues` - List queues
- `POST /api/v1/queues` - Create queue
- `GET /api/v1/queues/:id/members` - List members
- `POST /api/v1/queues/:id/members` - Add member
- `DELETE /api/v1/queues/:id/members/:extension` - Remove member

#### Configuration
- `GET /api/v1/config` - Get configuration
- `PUT /api/v1/config` - Update configuration
- `POST /api/v1/config/reload` - Reload Asterisk config
- `POST /api/v1/config/backup` - Create backup
- `POST /api/v1/config/restore` - Restore backup

---

## 📅 Timeline

### Gesamt: ~19 Wochen (ca. 4-5 Monate)

| Phase | Wochen | Beschreibung |
|-------|--------|--------------|
| Phase 1: PoC | ✅ Fertig | Grundgerüst, Docker, AMI Client |
| Phase 2: Core | 4 Wochen | Extensions, Trunks, Dialplan, CDR |
| Phase 3: Advanced | 6 Wochen | Recording, IVR, Conference, WebRTC |
| Phase 4: UI/UX | 5 Wochen | Dashboard, Forms, Visual Builder |
| Phase 5: Production | 4 Wochen | Security, Monitoring, Testing |

### Meilensteine

1. **M1 - PoC** ✅ (Woche 0)
   - Asterisk läuft
   - AMI-Verbindung
   - Basis-API

2. **M2 - Alpha** (Woche 4)
   - Extensions voll funktional
   - Trunks voll funktional
   - Calls anzeigen/steuern
   - Voicemail Basic

3. **M3 - Beta** (Woche 10)
   - Recording
   - IVR
   - Conference
   - WebRTC
   - Queues

4. **M4 - RC** (Woche 15)
   - Komplette UI
   - Visual Builder
   - Alle Features implementiert

5. **M5 - Release** (Woche 19)
   - Production-ready
   - Dokumentation
   - Tests
   - Security Audit

---

## 🚀 Quick Start (für Entwickler)

### PoC lokal testen

```bash
# 1. Repository klonen
cd /path/to/stumpfworks-nas

# 2. Plugin-Ordner wechseln
cd plugins/asterisk-voip

# 3. Go Dependencies laden
cd backend
go mod download

# 4. Docker Compose starten
cd ..
docker-compose up -d

# 5. Logs prüfen
docker-compose logs -f

# 6. API testen
curl http://localhost:8090/health
curl http://localhost:8090/api/v1/status

# 7. AMI testen (direkt)
telnet localhost 5038
# Login mit admin / stumpfworks2024
```

---

## 📚 Weitere Dokumentation

- [README.md](./README.md) - Übersicht und Installation
- [docs/API.md](./docs/API.md) - API Dokumentation
- [docs/QUICK_START.md](./docs/QUICK_START.md) - Quick Start Guide
- [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) - Architektur-Details

---

## 🤝 Beitragen

Du möchtest mithelfen? Super!

1. Fork das Repository
2. Erstelle einen Feature-Branch (`git checkout -b feature/amazing`)
3. Commit deine Änderungen
4. Push zum Branch
5. Erstelle einen Pull Request

---

## 📄 Lizenz

Siehe [LICENSE](../../LICENSE)

---

## ❓ Fragen?

- GitHub Issues: https://github.com/stumpf-works/stumpfworks-nas/issues
- Discussions: https://github.com/stumpf-works/stumpfworks-nas/discussions
