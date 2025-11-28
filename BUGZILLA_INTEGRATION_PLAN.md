# Bugzilla Error Reporting Integration - Implementierungsplan

**Projekt:** StumpfWorks NAS
**Feature:** Optional aktivierbares Fehlerberichts-System
**Datum:** 2025-01-27
**Status:** Planung (keine Implementierung)

---

## 📋 Executive Summary

Dieses Dokument beschreibt die Architektur und den Implementierungsplan für ein **opt-in basiertes Fehlerberichts-System** mit Bugzilla-Integration für StumpfWorks NAS.

**Kernprinzipien:**
- ✅ Opt-In erforderlich (Standard: **OFF**)
- ✅ Nur technische Fehler (Panics, API-Fehler, System-Fehler)
- ✅ Keine Nutzerdaten oder Inhalte
- ✅ Vollständig deaktivierbar
- ✅ Transparenz durch Open Source
- ✅ Minimalprinzip bei Datenerfassung

---

## 🎯 1. ANALYSE DER FEHLERQUELLEN

### 1.1 Aktuelles Error-Handling

**Stärken:**
- ✅ Konsistentes Error-Handling-Pattern in allen API-Handlers
- ✅ Strukturiertes Logging mit Uber Zap
- ✅ Panic Recovery via Chi Middleware
- ✅ Custom `AppError`-Typ mit HTTP-Status-Mapping
- ✅ Existierende Alert-/Webhook-Infrastruktur

**Schwächen:**
- ❌ Keine zentrale Fehlererfassung oder -aggregation
- ❌ Keine Error-Tracking über Requests hinweg
- ❌ Keine automatische Erkennung von Error-Spikes
- ❌ Duplikate in Error-Packages (`internal/errors` vs `pkg/errors`)

### 1.2 Identifizierte Fehlerquellen

#### **A. API-Layer (35+ Handler-Dateien)**
**Ort:** `backend/internal/api/handlers/*.go`

**Fehlertypen:**
- Request-Validierung (400 Bad Request)
- Authentifizierung/Autorisierung (401/403)
- Business-Logic-Fehler (409 Conflict)
- Service-Unavailability (503)
- Interne Server-Fehler (500)

**Muster:**
```go
if err != nil {
    logger.Error("Operation failed", zap.Error(err))
    utils.RespondError(w, errors.InternalServerError("...", err))
    return
}
```

**Integration-Point:** `backend/pkg/utils/response.go:RespondError()`

---

#### **B. System-Operationen**
**Ort:** `backend/internal/system/*.go`

**Kritische Operationen:**
- **Disk-Operationen:** Formatierung, Partitionierung, Mount/Unmount
- **RAID-Management:** mdadm-Operationen, RAID-Recovery
- **ZFS/Btrfs-Operations:** Pool-Erstellung, Snapshots
- **Network-Config:** Interface-Konfiguration, Routing, Bonding/VLANs
- **Service-Management:** Samba, NFS, iSCSI Restarts

**Fehlertypen:**
- Command execution failures (via `ShellExecutor`)
- Permission errors (`ErrNotRoot`)
- Hardware-Fehler (Disk failures, SMART errors)
- Konfigurationsfehler

**Integration-Point:** `backend/internal/system/shell.go:RunCommand()`

---

#### **C. Externe Service-Integrationen**

##### **Docker-Service**
**Ort:** `backend/internal/docker/docker.go`

**Fehlerquellen:**
- Docker daemon nicht verfügbar
- Container-Operationen (start, stop, remove)
- Image-Pull-Fehler
- Volume/Network-Operationen

##### **Active Directory**
**Ort:** `backend/internal/ad/*.go`

**Fehlerquellen:**
- LDAP-Verbindungsfehler
- Samba-Tool-Execution-Failures
- Domain-Join/Leave-Fehler
- DNS/Kerberos-Probleme

##### **Backup-Service**
**Ort:** `backend/internal/backup/backup.go`

**Fehlerquellen:**
- Snapshot-Erstellung fehlgeschlagen
- Restore-Operationen fehlgeschlagen
- Storage-Full-Fehler

---

#### **D. Database-Operationen**
**Ort:** `backend/internal/database/db.go`

**Fehlertypen:**
- Verbindungsfehler (fatal → Server stoppt)
- Query-Fehler (GORM errors)
- Record not found (wird oft ignoriert)
- Transaction-Fehler

**Besonderheit:** Database-Init-Fehler sind **fatal** und stoppen den Server.

---

#### **E. Panics**
**Ort:** Überall (gefangen von Chi's `middleware.Recoverer`)

**Quellen:**
- Nil-Pointer-Dereferences
- Index-Out-of-Bounds
- Type-Assertions
- Explizite `panic()` Aufrufe

**Aktuelles Verhalten:**
- Chi's Recoverer fängt Panic ab
- Loggt Stack-Trace
- Gibt 500 Internal Server Error zurück
- Server läuft weiter

**Problem:** Panics werden **nicht persistent gespeichert** oder gemeldet.

---

### 1.3 Zentrale vs. Verteilte Fehlerbehandlung

**Aktuell: Hybrid-Ansatz**

**Zentralisierte Komponenten:**
- ✅ `backend/pkg/errors/errors.go` - Error-Typen
- ✅ `backend/pkg/utils/response.go` - HTTP-Response-Handling
- ✅ `backend/pkg/logger/logger.go` - Logging
- ✅ `backend/internal/api/router.go` - Middleware (Recoverer)

**Verteilte Komponenten:**
- ❌ Error-Handling in jedem Handler individuell
- ❌ Service-Layer gibt Errors nach oben weiter
- ❌ Keine zentrale Error-Aggregation

**Empfehlung:** Zentralisieren durch **Error-Reporter-Service**

---

## 🔄 2. ABLAUF-/FUNKTIONSPLAN

### 2.1 Error-Detection-Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    ERROR DETECTION                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │  Error Source   │
                    └─────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
  ┌──────────┐        ┌──────────┐         ┌──────────┐
  │ API      │        │ System   │         │ Panic    │
  │ Error    │        │ Error    │         │ Recovery │
  └──────────┘        └──────────┘         └──────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Error Collector │
                    │  (New Service)  │
                    └─────────────────┘
```

### 2.2 Detaillierter Ablauf

#### **Schritt 1: Fehler Erkennen**

**API-Fehler:**
```
HTTP Handler → Error → RespondError() → Error Reporter Hook
```

**System-Fehler:**
```
System Operation → Command Fails → Error Return → Error Reporter Hook
```

**Panic:**
```
Any Code → panic() → Custom Recovery Middleware → Error Reporter Hook
```

---

#### **Schritt 2: Opt-In Prüfen**

```
┌──────────────────────────────────────────────────────────┐
│  Error Reporter Hook                                     │
└──────────────────────────────────────────────────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │ Check Opt-In     │
              │ (Config Service) │
              └──────────────────┘
                        │
        ┌───────────────┴───────────────┐
        │                               │
        ▼                               ▼
  ┌──────────┐                  ┌──────────┐
  │ Enabled  │                  │ Disabled │
  └──────────┘                  └──────────┘
        │                               │
        ▼                               ▼
  Continue                       Log Only
  Reporting                      (No Report)
```

**Config-Check:**
```go
config := errorReporter.GetConfig()
if !config.Enabled {
    return // No reporting
}
```

---

#### **Schritt 3: Fehlerdaten Aufbereiten**

**Gesammelte Daten (MINIMAL):**

```go
type ErrorReport struct {
    // Technical Metadata
    Timestamp   time.Time `json:"timestamp"`
    AppVersion  string    `json:"app_version"`
    GoVersion   string    `json:"go_version"`
    OS          string    `json:"os"`
    Arch        string    `json:"arch"`

    // Error Details
    ErrorType   string    `json:"error_type"`    // "api_error", "system_error", "panic"
    ErrorCode   int       `json:"error_code"`    // HTTP status or custom code
    ErrorMsg    string    `json:"error_message"` // Error message
    Stacktrace  string    `json:"stacktrace,omitempty"`

    // Context (KEINE USER-DATEN!)
    Component   string    `json:"component"`     // "api", "docker", "system"
    Operation   string    `json:"operation"`     // "list_containers", "format_disk"

    // Fingerprint für Deduplication
    Fingerprint string    `json:"fingerprint"`
}
```

**Explizit NICHT enthalten:**
- ❌ User-IDs, Usernames
- ❌ IP-Adressen
- ❌ Dateinamen, Pfade (außer interne Code-Pfade)
- ❌ Request-Bodies
- ❌ Konfigurationswerte
- ❌ Secrets, API-Keys

**Fingerprint-Generierung (für Deduplication):**
```
SHA256(ErrorType + ErrorCode + Component + StacktraceFirstLine)
```

---

#### **Schritt 4: Bugzilla-Client Aufrufen**

```
┌────────────────────────────────────────────────────────────┐
│  Prepared Error Report                                     │
└────────────────────────────────────────────────────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │ Bugzilla Client  │
              │ (REST API)       │
              └──────────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │ Check Cache      │
              │ (Deduplication)  │
              └──────────────────┘
                        │
        ┌───────────────┴───────────────┐
        │                               │
        ▼                               ▼
  ┌──────────┐                  ┌──────────┐
  │ New      │                  │ Duplicate│
  │ Error    │                  │ (24h)    │
  └──────────┘                  └──────────┘
        │                               │
        ▼                               ▼
  Create Bug                      Skip
  in Bugzilla                     (Already Reported)
```

**Bugzilla REST API Call:**
```http
POST https://bugzilla.stumpfworks.de/rest/bug
Content-Type: application/json
Authorization: Bearer <API_KEY>

{
  "product": "StumpfWorks NAS",
  "component": "Backend",
  "summary": "[Auto] API Error 500: Failed to list containers",
  "version": "1.2.1",
  "severity": "major",
  "description": "... formatted error report ..."
}
```

---

#### **Schritt 5: Doppelte Fehler Vermeiden**

**Deduplication-Strategie:**

1. **In-Memory-Cache** (24 Stunden):
   ```go
   type ErrorCache struct {
       mu          sync.RWMutex
       reported    map[string]time.Time // fingerprint -> timestamp
       maxAge      time.Duration        // 24h
   }
   ```

2. **Fingerprint-Check:**
   ```go
   if cache.WasReportedRecently(fingerprint) {
       logger.Debug("Error already reported", zap.String("fingerprint", fingerprint))
       return nil
   }
   ```

3. **Rate-Limiting:**
   - Max 10 Fehler pro Minute
   - Max 100 Fehler pro Stunde
   - Bei Überschreitung: Batch-Reporting

---

#### **Schritt 6: Offline-Fallback**

**Problem:** Bugzilla nicht erreichbar (Netzwerk, Server down)

**Lösung: Persistent Queue**

```
┌────────────────────────────────────────────────────────────┐
│  Bugzilla API Call Failed                                  │
└────────────────────────────────────────────────────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │ Persist to Disk  │
              │ (SQLite Queue)   │
              └──────────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │ Background Worker│
              │ (Retry Loop)     │
              └──────────────────┘
                        │
                        ▼
        ┌───────────────┴───────────────┐
        │                               │
        ▼                               ▼
  ┌──────────┐                  ┌──────────┐
  │ Success  │                  │ Max      │
  │          │                  │ Retries  │
  └──────────┘                  └──────────┘
        │                               │
        ▼                               ▼
  Delete from Queue              Discard (Log Error)
```

**Retry-Strategie:**
- Retry 1: nach 1 Minute
- Retry 2: nach 5 Minuten
- Retry 3: nach 15 Minuten
- Max 3 Retries, dann verwerfen

**Queue-Limits:**
- Max 1000 Einträge
- Max 7 Tage Aufbewahrung

---

#### **Schritt 7: Logging und Version-Bezug**

**Jeder Fehler wird geloggt (unabhängig von Opt-In):**

```go
logger.Error("Error occurred",
    zap.String("error_type", errorType),
    zap.Int("error_code", errorCode),
    zap.String("fingerprint", fingerprint),
    zap.Bool("reported_to_bugzilla", reportedSuccess),
)
```

**Version-Tracking:**
- Fehler werden mit `app_version` getaggt
- Bugzilla-Feld: `Version` = Git-Tag (z.B. "1.2.1-101-g6f02c45")
- Ermöglicht: "Fehler in Version X behoben?" Tracking

---

## 🏗️ 3. ÄNDERUNGS- & IMPLEMENTIERUNGSPLAN

### 3.1 Neue Komponenten

#### **A. Error Reporter Service**

**Datei:** `backend/internal/errorreporter/service.go` (NEU)

**Verantwortlichkeiten:**
- Error-Report-Aggregation
- Fingerprint-Generierung
- Deduplication-Cache
- Bugzilla-Client-Verwaltung
- Offline-Queue-Management

**Public API:**
```go
type Service struct {
    config        *Config
    bugzillaClient *BugzillaClient
    cache         *ErrorCache
    queue         *PersistentQueue
}

func Initialize(cfg *Config) (*Service, error)
func (s *Service) ReportError(err error, ctx ErrorContext) error
func (s *Service) ReportPanic(panicValue interface{}, stack []byte) error
func (s *Service) GetConfig() *Config
func (s *Service) UpdateConfig(cfg *Config) error
func (s *Service) GetStats() Stats
```

---

#### **B. Bugzilla Client**

**Datei:** `backend/internal/errorreporter/bugzilla.go` (NEU)

**Verantwortlichkeiten:**
- REST-API-Kommunikation
- Authentication (API-Key)
- Bug-Erstellung
- Connection-Pooling

**API:**
```go
type BugzillaClient struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

func NewBugzillaClient(baseURL, apiKey string) *BugzillaClient
func (c *BugzillaClient) CreateBug(report *ErrorReport) (bugID int, err error)
func (c *BugzillaClient) TestConnection() error
```

**Bugzilla REST API Endpoints:**
```
POST   /rest/bug              - Create bug
GET    /rest/bug/{id}         - Get bug details
GET    /rest/version          - Test connection
```

---

#### **C. Error Configuration**

**Datei:** `backend/internal/errorreporter/config.go` (NEU)

**Config-Struktur:**
```go
type Config struct {
    Enabled         bool   `json:"enabled" yaml:"enabled"`                   // Default: false
    BugzillaURL     string `json:"bugzilla_url" yaml:"bugzilla_url"`        // https://bugzilla.stumpfworks.de
    BugzillaAPIKey  string `json:"bugzilla_api_key" yaml:"bugzilla_api_key"` // Encrypted
    Product         string `json:"product" yaml:"product"`                  // "StumpfWorks NAS"
    Component       string `json:"component" yaml:"component"`              // "Backend"

    // Rate Limiting
    MaxReportsPerMinute int `json:"max_reports_per_minute" yaml:"max_reports_per_minute"` // 10
    MaxReportsPerHour   int `json:"max_reports_per_hour" yaml:"max_reports_per_hour"`     // 100

    // Deduplication
    DeduplicationWindow time.Duration `json:"deduplication_window" yaml:"deduplication_window"` // 24h

    // Queue
    MaxQueueSize    int           `json:"max_queue_size" yaml:"max_queue_size"`       // 1000
    QueueRetention  time.Duration `json:"queue_retention" yaml:"queue_retention"`     // 7 days
}
```

**Storage:** In `backend/internal/database/models/config.go` als `ErrorReporterConfig`

---

#### **D. Persistent Queue**

**Datei:** `backend/internal/errorreporter/queue.go` (NEU)

**Implementierung:** SQLite-basierte Queue (getrennt von Haupt-DB)

**Schema:**
```sql
CREATE TABLE error_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fingerprint TEXT NOT NULL,
    report_data TEXT NOT NULL,  -- JSON
    created_at TIMESTAMP NOT NULL,
    retry_count INTEGER DEFAULT 0,
    last_retry TIMESTAMP,
    status TEXT DEFAULT 'pending'  -- pending, retrying, failed
);

CREATE INDEX idx_status ON error_queue(status);
CREATE INDEX idx_created_at ON error_queue(created_at);
```

**API:**
```go
type PersistentQueue struct {
    db *sql.DB
}

func NewPersistentQueue(dbPath string) (*PersistentQueue, error)
func (q *PersistentQueue) Enqueue(report *ErrorReport) error
func (q *PersistentQueue) Dequeue() (*ErrorReport, error)
func (q *PersistentQueue) MarkAsProcessed(id int) error
func (q *PersistentQueue) GetPendingCount() int
func (q *PersistentQueue) Cleanup(olderThan time.Time) error
```

---

#### **E. Global Error Hook**

**Datei:** `backend/pkg/utils/response.go` (MODIFIKATION)

**Änderung:**
```go
// VORHER:
func RespondError(w http.ResponseWriter, err error) {
    // ... existing error response logic ...
    logger.Error("Request failed", zap.Error(err))
}

// NACHHER:
func RespondError(w http.ResponseWriter, err error) {
    // ... existing error response logic ...
    logger.Error("Request failed", zap.Error(err))

    // NEW: Report to error tracking service
    if appErr, ok := err.(*errors.AppError); ok && appErr.Code >= 500 {
        if reporter := errorreporter.GetService(); reporter != nil {
            ctx := ErrorContext{
                Component: "api",
                Operation: extractOperationFromRequest(r),
            }
            go reporter.ReportError(appErr, ctx)  // Non-blocking
        }
    }
}
```

---

#### **F. Custom Panic Recovery Middleware**

**Datei:** `backend/internal/api/middleware/recovery.go` (NEU)

**Implementierung:**
```go
func PanicRecovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if panicValue := recover(); panicValue != nil {
                stack := debug.Stack()

                // Log panic
                logger.Error("Panic recovered",
                    zap.Any("panic", panicValue),
                    zap.ByteString("stack", stack),
                )

                // Report to error tracking (if enabled)
                if reporter := errorreporter.GetService(); reporter != nil {
                    go reporter.ReportPanic(panicValue, stack)
                }

                // Return 500 response
                http.Error(w, "Internal Server Error", 500)
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

**Integration:** `backend/internal/api/router.go`:
```go
r.Use(middleware.PanicRecovery)  // NEW - Reports panics
r.Use(middleware.Recoverer)      // EXISTING - Fallback
```

---

### 3.2 Anpassungen an bestehenden Dateien

#### **A. Router Middleware-Stack**

**Datei:** `backend/internal/api/router.go`

**Änderung:**
```go
// VORHER:
r.Use(middleware.Recoverer)

// NACHHER:
r.Use(middleware.PanicRecovery)  // NEW - Custom recovery mit Error Reporting
r.Use(middleware.Recoverer)      // Fallback
```

---

#### **B. Database Models**

**Datei:** `backend/internal/database/models/config.go`

**Neue Tabelle:**
```go
type ErrorReporterConfig struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Enabled   bool      `gorm:"not null;default:false" json:"enabled"`
    BugzillaURL string  `gorm:"type:varchar(255)" json:"bugzilla_url"`
    BugzillaAPIKey string `gorm:"type:text" json:"bugzilla_api_key"` // Encrypted
    Product   string    `gorm:"type:varchar(100)" json:"product"`
    Component string    `gorm:"type:varchar(100)" json:"component"`
    MaxReportsPerMinute int `gorm:"default:10" json:"max_reports_per_minute"`
    MaxReportsPerHour   int `gorm:"default:100" json:"max_reports_per_hour"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (ErrorReporterConfig) TableName() string {
    return "error_reporter_config"
}
```

**Migration:** In `backend/internal/database/db.go:initializeDatabase()`:
```go
&models.ErrorReporterConfig{},
```

---

#### **C. Main Server Initialization**

**Datei:** `backend/cmd/stumpfworks-server/main.go`

**Änderung:**
```go
// Nach Database-Init und vor Router-Init

// Initialize Error Reporter (optional service, non-fatal)
errorReporterService, err := errorreporter.Initialize()
if err != nil {
    logger.Warn("Failed to initialize error reporter", zap.Error(err))
} else {
    logger.Info("Error reporter initialized")
}
```

---

#### **D. System Shell Execution**

**Datei:** `backend/internal/system/shell.go`

**Änderung (optional):**
```go
func (e *RealShellExecutor) RunCommand(name string, args ...string) (string, error) {
    output, err := cmd.CombinedOutput()
    if err != nil {
        cmdErr := fmt.Errorf("%s failed: %s: %w", name, string(output), err)

        // Report critical system errors
        if reporter := errorreporter.GetService(); reporter != nil {
            if isCriticalCommand(name) {
                ctx := ErrorContext{
                    Component: "system",
                    Operation: name,
                }
                go reporter.ReportError(cmdErr, ctx)
            }
        }

        return "", cmdErr
    }
    return string(output), nil
}
```

**Kritische Commands:** `mkfs`, `mdadm`, `zpool`, `mount`

---

### 3.3 Frontend-Erweiterungen (React)

#### **A. Settings-Page Erweiterung**

**Datei:** `frontend/src/apps/Settings/Settings.tsx`

**Neuer Tab:** "Error Reporting"

**UI-Komponenten:**

```tsx
interface ErrorReportingSettings {
  enabled: boolean;
  bugzillaUrl: string;
  product: string;
  component: string;
}

function ErrorReportingTab() {
  const [config, setConfig] = useState<ErrorReportingSettings>(...);
  const [testStatus, setTestStatus] = useState<'idle' | 'testing' | 'success' | 'error'>('idle');

  return (
    <div className="space-y-6">
      {/* Privacy Notice */}
      <Card className="border-blue-200 bg-blue-50">
        <div className="p-4">
          <h3 className="font-semibold mb-2">🔒 Datenschutz-Hinweis</h3>
          <ul className="text-sm space-y-1">
            <li>✅ Nur technische Fehler werden gemeldet</li>
            <li>✅ Keine Benutzerdaten oder Inhalte</li>
            <li>✅ Jederzeit deaktivierbar</li>
            <li>✅ Keine Hintergrund-Übertragung</li>
          </ul>
        </div>
      </Card>

      {/* Enable/Disable Toggle */}
      <Card>
        <div className="p-4">
          <label className="flex items-center gap-3">
            <input
              type="checkbox"
              checked={config.enabled}
              onChange={(e) => setConfig({...config, enabled: e.target.checked})}
            />
            <div>
              <div className="font-medium">Automatische Fehlerberichte aktivieren</div>
              <div className="text-sm text-gray-600">
                Hilft uns, Fehler schneller zu beheben
              </div>
            </div>
          </label>
        </div>
      </Card>

      {/* Configuration (nur wenn enabled) */}
      {config.enabled && (
        <Card>
          <div className="p-4 space-y-4">
            <h3 className="font-semibold">Bugzilla-Konfiguration</h3>

            <div>
              <label className="block text-sm font-medium mb-1">Bugzilla URL</label>
              <input
                type="url"
                value={config.bugzillaUrl}
                onChange={(e) => setConfig({...config, bugzillaUrl: e.target.value})}
                className="w-full px-3 py-2 border rounded-lg"
                placeholder="https://bugzilla.stumpfworks.de"
              />
            </div>

            <button
              onClick={handleTestConnection}
              disabled={testStatus === 'testing'}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              {testStatus === 'testing' ? 'Teste...' : 'Verbindung testen'}
            </button>
          </div>
        </Card>
      )}

      {/* Statistics */}
      <Card>
        <div className="p-4">
          <h3 className="font-semibold mb-3">Statistiken</h3>
          <div className="grid grid-cols-3 gap-4">
            <div>
              <div className="text-2xl font-bold">{stats.reportedToday}</div>
              <div className="text-sm text-gray-600">Heute gemeldet</div>
            </div>
            <div>
              <div className="text-2xl font-bold">{stats.queuedReports}</div>
              <div className="text-sm text-gray-600">In Warteschlange</div>
            </div>
            <div>
              <div className="text-2xl font-bold">{stats.totalReported}</div>
              <div className="text-sm text-gray-600">Gesamt gemeldet</div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}
```

---

#### **B. API-Client**

**Datei:** `frontend/src/api/errorreporting.ts` (NEU)

```typescript
export interface ErrorReportingConfig {
  enabled: boolean;
  bugzilla_url: string;
  product: string;
  component: string;
  max_reports_per_minute: number;
  max_reports_per_hour: number;
}

export interface ErrorReportingStats {
  reported_today: number;
  queued_reports: number;
  total_reported: number;
  last_report_at?: string;
}

export const errorReportingApi = {
  getConfig: async () => {
    const response = await client.get<ApiResponse<ErrorReportingConfig>>('/error-reporting/config');
    return response.data;
  },

  updateConfig: async (config: ErrorReportingConfig) => {
    const response = await client.put<ApiResponse<{ message: string }>>('/error-reporting/config', config);
    return response.data;
  },

  testConnection: async () => {
    const response = await client.post<ApiResponse<{ message: string }>>('/error-reporting/test');
    return response.data;
  },

  getStats: async () => {
    const response = await client.get<ApiResponse<ErrorReportingStats>>('/error-reporting/stats');
    return response.data;
  },
};
```

---

#### **C. Backend API-Handler**

**Datei:** `backend/internal/api/handlers/errorreporting.go` (NEU)

**Endpoints:**
```go
type ErrorReportingHandler struct {
    service *errorreporter.Service
}

// GET /api/v1/error-reporting/config
func (h *ErrorReportingHandler) GetConfig(w http.ResponseWriter, r *http.Request)

// PUT /api/v1/error-reporting/config
func (h *ErrorReportingHandler) UpdateConfig(w http.ResponseWriter, r *http.Request)

// POST /api/v1/error-reporting/test
func (h *ErrorReportingHandler) TestConnection(w http.ResponseWriter, r *http.Request)

// GET /api/v1/error-reporting/stats
func (h *ErrorReportingHandler) GetStats(w http.ResponseWriter, r *http.Request)
```

**Router Registration:** `backend/internal/api/router.go`:
```go
// Error Reporting
errorReportingHandler := handlers.NewErrorReportingHandler(errorReporterService)
r.Route("/error-reporting", func(r chi.Router) {
    r.Get("/config", errorReportingHandler.GetConfig)
    r.Put("/config", errorReportingHandler.UpdateConfig)
    r.Post("/test", errorReportingHandler.TestConnection)
    r.Get("/stats", errorReportingHandler.GetStats)
})
```

---

### 3.4 Bugzilla Installation

**Server:** `46.4.25.15` (apt.stumpfworks.de)

**Installation:**

```bash
# 1. Install dependencies
apt-get update
apt-get install -y apache2 mysql-server libapache2-mod-perl2 \
    libdbd-mysql-perl libtemplate-perl libemail-sender-perl \
    libdatetime-perl libcgi-pm-perl libmath-random-isaac-perl \
    libhtml-scrubber-perl libjson-xs-perl liblist-moreutils-perl \
    libwww-perl liburi-perl libxml-twig-perl

# 2. Download Bugzilla
cd /var/www
wget https://ftp.mozilla.org/pub/mozilla.org/webtools/bugzilla-5.0.6.tar.gz
tar -xzf bugzilla-5.0.6.tar.gz
mv bugzilla-5.0.6 bugzilla
cd bugzilla

# 3. Run checksetup
./checksetup.pl --check-modules
./install-module.sh --all

# 4. Configure database
mysql -u root -p
CREATE DATABASE bugs;
CREATE USER 'bugs'@'localhost' IDENTIFIED BY '<secure-password>';
GRANT ALL PRIVILEGES ON bugs.* TO 'bugs'@'localhost';
FLUSH PRIVILEGES;

# 5. Configure Bugzilla
cp localconfig.sample localconfig
vim localconfig
# Set: $db_name = 'bugs';
#      $db_user = 'bugs';
#      $db_pass = '<secure-password>';

./checksetup.pl
# Follow prompts for admin user

# 6. Configure Apache
cat > /etc/apache2/sites-available/bugzilla.conf <<'EOF'
<VirtualHost *:80>
    ServerName bugzilla.stumpfworks.de
    DocumentRoot /var/www/bugzilla

    <Directory /var/www/bugzilla>
        AddHandler cgi-script .cgi
        Options +ExecCGI
        DirectoryIndex index.cgi index.html
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>
EOF

a2ensite bugzilla
a2enmod cgi headers expires rewrite
systemctl restart apache2

# 7. Setup SSL (Let's Encrypt)
apt-get install -y certbot python3-certbot-apache
certbot --apache -d bugzilla.stumpfworks.de
```

**Bugzilla-Konfiguration:**

1. **Product erstellen:** "StumpfWorks NAS"
2. **Components:** Backend, Frontend, System
3. **API-Key generieren:** User Preferences → API Keys
4. **REST-API aktivieren:** editparams.cgi → REST API

---

## 🎯 4. BEST-PRACTICE-EMPFEHLUNGEN

### 4.1 Datenschutz & Minimalprinzip

#### **✅ DO:**

1. **Nur technische Daten sammeln:**
   - Error-Typ, Message, Stacktrace
   - App-Version, OS, Architektur
   - Component, Operation

2. **Anonymisierung:**
   - Keine IP-Adressen
   - Keine User-IDs
   - Keine Dateipfade (außer Code-Pfade)

3. **Transparenz:**
   - Open-Source-Implementierung
   - Dokumentation was gesendet wird
   - Opt-In mit klarer Erklärung

4. **User-Kontrolle:**
   - Standard: OFF
   - Jederzeit deaktivierbar
   - Sichtbare Statistiken

#### **❌ DON'T:**

1. Keine automatische Aktivierung
2. Keine Sammlung von Request-Bodies
3. Keine Speicherung von Credentials
4. Keine Tracking-IDs über Requests hinweg
5. Keine Hintergrund-Übertragung bei deaktiviertem Feature

---

### 4.2 Testbarkeit

#### **Mock-Interfaces:**

```go
// Bugzilla Client Interface
type BugzillaClientInterface interface {
    CreateBug(report *ErrorReport) (int, error)
    TestConnection() error
}

// Mock Implementation
type MockBugzillaClient struct {
    CreateBugFunc      func(*ErrorReport) (int, error)
    TestConnectionFunc func() error
}

func (m *MockBugzillaClient) CreateBug(report *ErrorReport) (int, error) {
    if m.CreateBugFunc != nil {
        return m.CreateBugFunc(report)
    }
    return 123, nil // Default mock
}
```

#### **Test-Szenarien:**

```go
func TestErrorReporter_ReportError_OptInDisabled(t *testing.T) {
    // Arrange
    config := &Config{Enabled: false}
    reporter := NewService(config)

    // Act
    err := reporter.ReportError(testError, testContext)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 0, mockClient.CreateBugCallCount)
}

func TestErrorReporter_Deduplication(t *testing.T) {
    // Test dass gleicher Fehler nicht 2x reported wird
}

func TestBugzillaClient_RetryOnFailure(t *testing.T) {
    // Test Offline-Queue und Retry-Logik
}
```

---

### 4.3 Skalierbarkeit

#### **Performance-Considerations:**

1. **Non-Blocking Reporting:**
   ```go
   go reporter.ReportError(err, ctx)  // Goroutine
   ```

2. **Rate-Limiting:**
   - Sliding-Window-Algorithmus
   - Verhindert Bugzilla-Overload

3. **Batch-Reporting (optional):**
   - Bei Error-Spikes: Batch von 10 Fehlern zusammenfassen
   - "10 identical errors occurred"

4. **Connection-Pooling:**
   ```go
   httpClient := &http.Client{
       Timeout: 10 * time.Second,
       Transport: &http.Transport{
           MaxIdleConns:        10,
           IdleConnTimeout:     30 * time.Second,
           MaxIdleConnsPerHost: 5,
       },
   }
   ```

#### **Resource-Limits:**

- Queue-Größe: Max 1000 Einträge
- Memory-Cache: Max 10.000 Fingerprints (ca. 1MB)
- Disk-Queue: Max 100MB

---

### 4.4 Logging-Konzept

#### **Logging-Levels:**

```go
// Normaler Error (immer loggen)
logger.Error("API request failed", zap.Error(err))

// Error Reporting aktiviert
logger.Info("Error reported to Bugzilla",
    zap.String("fingerprint", fp),
    zap.Int("bug_id", bugID),
)

// Error Reporting fehlgeschlagen
logger.Warn("Failed to report error to Bugzilla",
    zap.Error(err),
    zap.String("fingerprint", fp),
)

// Deduplication
logger.Debug("Error already reported, skipping",
    zap.String("fingerprint", fp),
)
```

#### **Structured Logging:**

```go
logger.Error("System command failed",
    zap.String("command", "mkfs.ext4"),
    zap.String("component", "system"),
    zap.String("operation", "format_disk"),
    zap.Error(err),
    zap.Bool("reported_to_bugzilla", true),
)
```

---

### 4.5 Versionierung

#### **Fehler pro Build/Version:**

- Jeder Error-Report enthält Git-Tag/Commit
- Bugzilla-Feld: `Version` = "1.2.1-101-g6f02c45"
- Ermöglicht: "In Version X behoben" Tracking

#### **Changelog-Integration:**

```markdown
## [1.2.2] - 2025-01-30

### Fixed
- Fixed panic in Docker container listing (Bug #123)
- Fixed disk formatting failure on ext4 (Bug #124)
```

#### **Automated Bug-Closing:**

- Git-Commit-Message: `Fixes #123`
- CI/CD-Integration könnte automatisch Bugzilla-Bug schließen

---

### 4.6 Clean Architecture

#### **Dependency-Injection:**

```go
type Service struct {
    config         *Config
    bugzillaClient BugzillaClientInterface  // Interface, nicht Struct
    cache          CacheInterface
    queue          QueueInterface
    logger         *zap.Logger
}

func NewService(
    config *Config,
    client BugzillaClientInterface,
    cache CacheInterface,
    queue QueueInterface,
    logger *zap.Logger,
) *Service {
    return &Service{...}
}
```

#### **Separation of Concerns:**

```
errorreporter/
├── service.go          # Service-Logik
├── bugzilla.go         # Bugzilla-Client
├── cache.go            # Deduplication-Cache
├── queue.go            # Persistent Queue
├── config.go           # Configuration
├── fingerprint.go      # Fingerprint-Generierung
├── models.go           # Data-Modelle
└── interfaces.go       # Interfaces für Testing
```

#### **No Global State (außer Singleton-Service):**

```go
// GOOD
reporter := errorreporter.GetService()
if reporter != nil {
    reporter.ReportError(...)
}

// BAD
errorreporter.ReportError(...)  // Direkter globaler Aufruf
```

---

## 📊 5. ARCHITEKTUR-/SEQUENZDIAGRAMM

### 5.1 Komponenten-Architektur

```
┌────────────────────────────────────────────────────────────────────┐
│                      STUMPFWORKS NAS                               │
└────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                          API LAYER                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐             │
│  │   Docker     │  │   Storage    │  │    Users     │  ...        │
│  │   Handler    │  │   Handler    │  │   Handler    │             │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘             │
│         │                  │                  │                     │
│         └──────────────────┼──────────────────┘                     │
│                            │                                        │
│                            ▼                                        │
│                ┌──────────────────────┐                            │
│                │  Response Utilities  │                            │
│                │  (RespondError)      │◄───── Global Error Hook    │
│                └──────────┬───────────┘                            │
└───────────────────────────┼────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    ERROR REPORTER SERVICE                           │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Service Coordinator                                         │  │
│  │  • Opt-In Check                                              │  │
│  │  • Fingerprint Generation                                    │  │
│  │  • Rate Limiting                                             │  │
│  └──────────────────────────────────────────────────────────────┘  │
│           │                    │                    │               │
│           ▼                    ▼                    ▼               │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐         │
│  │ Deduplication│    │   Bugzilla   │    │  Persistent  │         │
│  │    Cache     │    │    Client    │    │    Queue     │         │
│  │ (In-Memory)  │    │  (REST API)  │    │  (SQLite)    │         │
│  └──────────────┘    └──────┬───────┘    └──────────────┘         │
└─────────────────────────────┼──────────────────────────────────────┘
                              │
                              ▼ HTTPS
                   ┌────────────────────┐
                   │     Bugzilla       │
                   │  bugzilla.         │
                   │  stumpfworks.de    │
                   └────────────────────┘
```

---

### 5.2 Fehler-Reporting-Sequenz (API-Error)

```
┌────────┐     ┌────────┐     ┌──────────┐     ┌─────────┐     ┌──────────┐     ┌──────────┐
│ Client │     │  API   │     │ Response │     │  Error  │     │ Bugzilla │     │ Bugzilla │
│        │     │Handler │     │  Utils   │     │Reporter │     │  Client  │     │  Server  │
└───┬────┘     └───┬────┘     └────┬─────┘     └────┬────┘     └────┬─────┘     └────┬─────┘
    │              │               │                │               │                │
    │  HTTP        │               │                │               │                │
    │  Request     │               │                │               │                │
    ├─────────────►│               │                │               │                │
    │              │               │                │               │                │
    │              │  Execute      │                │               │                │
    │              │  Operation    │                │               │                │
    │              │  (Error!)     │                │               │                │
    │              │               │                │               │                │
    │              │  RespondError │                │               │                │
    │              ├──────────────►│                │               │                │
    │              │               │                │               │                │
    │              │               │  Check Opt-In  │               │                │
    │              │               ├───────────────►│               │                │
    │              │               │                │               │                │
    │              │               │  (Enabled)     │               │                │
    │              │               │◄───────────────┤               │                │
    │              │               │                │               │                │
    │              │               │  Generate      │               │                │
    │              │               │  Fingerprint   │               │                │
    │              │               ├───────────────►│               │                │
    │              │               │                │               │                │
    │              │               │  Check Cache   │               │                │
    │              │               ├───────────────►│               │                │
    │              │               │  (Not Seen)    │               │                │
    │              │               │◄───────────────┤               │                │
    │              │               │                │               │                │
    │              │               │  Prepare       │               │                │
    │              │               │  Report        │               │                │
    │              │               ├───────────────►│               │                │
    │              │               │                │               │                │
    │              │               │  CreateBug     │               │                │
    │              │               │                ├──────────────►│                │
    │              │               │                │               │                │
    │              │               │                │  POST /rest/  │                │
    │              │               │                │  bug          │                │
    │              │               │                │               ├───────────────►│
    │              │               │                │               │                │
    │              │               │                │               │  Bug Created   │
    │              │               │                │               │  (ID: 123)     │
    │              │               │                │               │◄───────────────┤
    │              │               │                │               │                │
    │              │               │                │  Bug ID: 123  │                │
    │              │               │                │◄──────────────┤                │
    │              │               │                │               │                │
    │              │               │  Log Success   │               │                │
    │              │               │◄───────────────┤               │                │
    │              │               │                │               │                │
    │              │  HTTP 500     │                │               │                │
    │◄─────────────┤               │                │               │                │
    │              │               │                │               │                │
```

---

### 5.3 Panic-Recovery-Sequenz

```
┌────────┐     ┌────────┐     ┌──────────┐     ┌─────────┐     ┌──────────┐
│ Client │     │  HTTP  │     │  Panic   │     │  Error  │     │ Bugzilla │
│        │     │Handler │     │ Recovery │     │Reporter │     │  Client  │
└───┬────┘     └───┬────┘     └────┬─────┘     └────┬────┘     └────┬─────┘
    │              │               │                │               │
    │  HTTP        │               │                │               │
    │  Request     │               │                │               │
    ├─────────────►│               │                │               │
    │              │               │                │               │
    │              │  Execute      │                │               │
    │              │  Handler      │                │               │
    │              │               │                │               │
    │              │  PANIC!       │                │               │
    │              ├───────────────┼───────────────X│               │
    │              │               │                │               │
    │              │               │  recover()     │               │
    │              │               │◄───────────────┤               │
    │              │               │                │               │
    │              │               │  Capture       │               │
    │              │               │  Stack Trace   │               │
    │              │               │                │               │
    │              │               │  Log Panic     │               │
    │              │               │                │               │
    │              │               │  ReportPanic() │               │
    │              │               ├───────────────►│               │
    │              │               │                │               │
    │              │               │                │  CreateBug    │
    │              │               │                ├──────────────►│
    │              │               │                │               │
    │              │               │  HTTP 500      │               │
    │              │◄──────────────┤                │               │
    │              │               │                │               │
    │  HTTP 500    │               │                │               │
    │◄─────────────┤               │                │               │
    │              │               │                │               │
```

---

### 5.4 Offline-Queue-Sequenz

```
┌─────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Error  │     │ Bugzilla │     │ Persistent│     │Background│
│Reporter │     │  Client  │     │   Queue   │     │  Worker  │
└────┬────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │               │                │                │
     │  CreateBug    │                │                │
     ├──────────────►│                │                │
     │               │                │                │
     │               │  HTTPS Request │                │
     │               ├───────────────X│  (Network      │
     │               │                │   Error)       │
     │               │                │                │
     │  Error        │                │                │
     │◄──────────────┤                │                │
     │               │                │                │
     │  Enqueue      │                │                │
     ├───────────────────────────────►│                │
     │               │                │                │
     │               │                │  INSERT INTO   │
     │               │                │  error_queue   │
     │               │                │                │
     │  Success      │                │                │
     │◄───────────────────────────────┤                │
     │               │                │                │
     │               │                │  (Every 5min)  │
     │               │                │                │
     │               │                │  Dequeue       │
     │               │                │◄───────────────┤
     │               │                │                │
     │               │                │  Report Data   │
     │               │                ├───────────────►│
     │               │                │                │
     │               │                │  CreateBug     │
     │               │◄───────────────────────────────┤│
     │               │                │                │
     │               │  (Success)     │                │
     │               ├───────────────────────────────►││
     │               │                │                │
     │               │                │  MarkProcessed │
     │               │                │◄───────────────┤
     │               │                │                │
     │               │                │  DELETE FROM   │
     │               │                │  error_queue   │
     │               │                │                │
```

---

## 📝 6. IMPLEMENTIERUNGS-CHECKLISTE

### Phase 1: Basis-Infrastruktur (Tag 1)

- [ ] Bugzilla auf Server installieren
- [ ] DNS-Konfiguration (bugzilla.stumpfworks.de)
- [ ] SSL-Zertifikat (Let's Encrypt)
- [ ] Bugzilla konfigurieren (Product, Components, API-Key)

### Phase 2: Backend-Core (Tag 2-3)

- [ ] Neue Package-Struktur: `backend/internal/errorreporter/`
- [ ] `service.go` - Core-Service-Logik
- [ ] `bugzilla.go` - Bugzilla-REST-Client
- [ ] `cache.go` - Deduplication-Cache
- [ ] `queue.go` - Persistent Queue (SQLite)
- [ ] `config.go` - Configuration-Management
- [ ] `fingerprint.go` - Fingerprint-Generierung
- [ ] `models.go` - Data-Models
- [ ] Database-Migration für ErrorReporterConfig
- [ ] Unit-Tests für alle Komponenten

### Phase 3: Integration (Tag 4)

- [ ] Custom Panic Recovery Middleware
- [ ] RespondError()-Hook für Error-Reporting
- [ ] System-Shell-Execution-Hook (optional)
- [ ] Router-Integration (Middleware, API-Endpoints)
- [ ] Main-Server-Initialization

### Phase 4: Frontend (Tag 5)

- [ ] Settings-Tab "Error Reporting"
- [ ] API-Client (`frontend/src/api/errorreporting.ts`)
- [ ] UI-Komponenten (Toggle, Config, Stats)
- [ ] Privacy-Notice-Card

### Phase 5: Testing & Dokumentation (Tag 6)

- [ ] Integration-Tests (End-to-End)
- [ ] Mock-Tests (Bugzilla-Client)
- [ ] Load-Tests (Rate-Limiting)
- [ ] User-Dokumentation
- [ ] Admin-Dokumentation (Bugzilla-Setup)

### Phase 6: Deployment (Tag 7)

- [ ] Build & Package
- [ ] Deployment auf Test-System
- [ ] QA-Testing
- [ ] Production-Deployment
- [ ] Monitoring-Setup

---

## ⚠️ WICHTIGE HINWEISE

### ❌ Was NICHT implementiert werden soll

1. **Automatische Aktivierung** - Standard muss OFF sein
2. **Tracking-IDs** - Keine User-Verfolgung über Requests
3. **Hintergrund-Übertragung** - Nur bei aktiven Fehlern
4. **User-Daten-Sammlung** - Keine IPs, Usernames, Inhalte
5. **Closed-Source-Komponenten** - Alles Open Source

### ✅ Must-Have-Features

1. **Opt-In-Requirement** - Explizite User-Zustimmung
2. **Deaktivierbarkeit** - Jederzeit ausschaltbar
3. **Transparenz** - Klare Dokumentation was gesendet wird
4. **Minimalprinzip** - Nur absolut notwendige Daten
5. **Offline-Fähigkeit** - Queue bei Netzwerkausfall

### 🔒 Datenschutz-Compliance

- **DSGVO-konform:** Opt-In, Minimalprinzip, Transparenz
- **Keine personenbezogenen Daten:** Nur technische Fehler
- **User-Kontrolle:** Vollständig deaktivierbar
- **Open-Source:** Code ist einsehbar

---

## 📈 ERFOLGSKRITERIEN

### Funktionale Kriterien

- ✅ Panics werden automatisch gemeldet (bei Opt-In)
- ✅ 500er-Errors werden automatisch gemeldet (bei Opt-In)
- ✅ System-Fehler werden gemeldet (bei Opt-In)
- ✅ Duplikate werden erkannt und nicht doppelt gemeldet
- ✅ Offline-Queue funktioniert bei Netzwerkausfall
- ✅ Frontend zeigt Opt-In/Opt-Out-Einstellung
- ✅ Statistiken sind sichtbar

### Non-Funktionale Kriterien

- ✅ Keine Performance-Einbußen (Non-blocking)
- ✅ Keine False-Positives (nur echte Fehler)
- ✅ Keine Datenschutz-Verletzungen
- ✅ 100% Test-Coverage für Core-Logik
- ✅ Dokumentation vollständig

---

## 🎯 ZEITSCHÄTZUNG

**Gesamt: 6-7 Arbeitstage**

| Phase | Dauer | Beschreibung |
|-------|-------|--------------|
| Phase 1 | 0.5 Tage | Bugzilla-Installation |
| Phase 2 | 2 Tage | Backend-Core-Komponenten |
| Phase 3 | 1 Tag | Integration in bestehendes System |
| Phase 4 | 1 Tag | Frontend-Entwicklung |
| Phase 5 | 1 Tag | Testing & Dokumentation |
| Phase 6 | 0.5 Tage | Deployment |

**Puffer:** +1 Tag für unvorhergesehene Probleme

---

## 📚 REFERENZEN

### Technische Dokumentation

- **Bugzilla REST API:** https://bugzilla.readthedocs.io/en/latest/api/
- **Go Error Wrapping:** https://go.dev/blog/go1.13-errors
- **Uber Zap Logger:** https://pkg.go.dev/go.uber.org/zap
- **Chi Router Middleware:** https://github.com/go-chi/chi

### Best Practices

- **OWASP Top 10:** https://owasp.org/www-project-top-ten/
- **DSGVO-Guidance:** https://gdpr.eu/
- **Error Tracking Best Practices:** https://sentry.io/best-practices/

---

**Ende des Implementierungsplans**

*Erstellt: 2025-01-27*
*Autor: Claude (Anthropic)*
*Version: 1.0*
