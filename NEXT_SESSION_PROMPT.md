# Stumpf.Works NAS - Session Briefing für Release 1.0

## Projekt-Status
Wir arbeiten am **Stumpf.Works NAS** Projekt - ein NAS-System mit Web-UI im macOS-Stil.

**Aktueller Branch:** `claude/fix-errors-015RUcvbj7wCZgJaP2KPj3Me`
**Aktuelle Version:** 0.3.0 (bereit für Release)
**Ziel:** Release 1.0 vorbereiten

## Letzte Session - Was wurde erreicht:

✅ **User Groups Feature** komplett implementiert (Backend + Frontend)
✅ **GitHub Actions** für automatische Releases und CI eingerichtet
✅ **TypeScript Build-Fehler** alle behoben
✅ **Update-Checker 404-Fehler** gefixt
✅ **Samba-Konfiguration** verbessert (direkt in smb.conf schreiben)

## AUFGABEN FÜR DIESE SESSION

### 1️⃣ KRITISCHE BUGS BEHEBEN

#### A) Audit Log - 404 Fehler
- **Problem:** Request failed with status code 404
- **Location:** Frontend Audit Log Bereich
- **TODO:**
  - Backend-Endpunkt `/api/v1/audit` prüfen ob vorhanden
  - Falls nicht: Endpunkt implementieren
  - Frontend API-Call prüfen und fixen

#### B) Security - 404 Fehler
- **Problem:** Request failed with status code 404
- **Location:** Frontend Security Bereich
- **TODO:**
  - Backend-Endpunkt prüfen (vermutlich `/api/v1/security` oder `/api/v1/system/security`)
  - Falls nicht: Endpunkt implementieren
  - Frontend API-Call prüfen und fixen

#### C) Scheduled Tasks - Fehler
- **Problem:** Nicht näher spezifiziert, aber funktioniert nicht
- **Location:** Scheduled Tasks Bereich
- **TODO:**
  - Backend-Endpunkt `/api/v1/tasks` oder `/api/v1/scheduler` prüfen
  - Frontend prüfen auf Fehler in Console
  - Verbindung Backend ↔ Frontend fixen

#### D) Docker Networks - UI Crash
- **Problem:** Wenn man auf "Networks" klickt wird das gesamte Web-Interface schwarz
- **Location:** Docker Manager → Networks Tab
- **TODO:**
  - JavaScript-Fehler in Browser Console identifizieren
  - Komponente `DockerManager` oder `NetworkManager` prüfen
  - Vermutlich: undefined/null Fehler oder fehlender Error Handler
  - Fix implementieren

#### E) Dashboard Advanced Monitoring - Zeigt nur Metrics
- **Problem:** Im Advanced Monitoring Bereich werden nur Metrics angezeigt
- **Location:** Dashboard → Advanced Monitoring
- **TODO:**
  - Prüfen was SOLLTE angezeigt werden (Logs? Alerts? Weitere Daten?)
  - Backend-Endpunkte für erweiterte Monitoring-Daten prüfen
  - Frontend erweitern um zusätzliche Daten anzuzeigen

#### F) Disks - Partitionen nicht sichtbar
- **Problem:** Partitionen werden nicht angezeigt (außer bei System-Disk)
- **Location:** Storage Manager → Disks
- **TODO:**
  - Backend prüfen ob Partitionen korrekt gelesen werden (lsblk, fdisk, etc.)
  - API-Response prüfen ob `partitions[]` Array gefüllt ist
  - Frontend `DiskManager.tsx` prüfen ob Partitionen gerendert werden
  - Fix: Backend oder Frontend oder beides

### 2️⃣ VOLLSTÄNDIGKEITS-CHECK: Backend ↔ Frontend

**Ziel:** Sicherstellen dass ALLE Backend-Endpunkte im Frontend integriert sind

#### Phase 1: Backend-Routen vollständig auflisten

```bash
# 1. ALLE Routen aus router.go extrahieren
cd backend
grep -rn "r\.Get\|r\.Post\|r\.Put\|r\.Delete\|r\.Patch" internal/api/router.go > /tmp/backend_routes.txt

# 2. Auch Handler-Dateien prüfen
find internal/api/handlers -name "*.go" -exec grep -l "func.*http\.ResponseWriter" {} \;

# 3. Ausgabe: Komplette Liste aller Backend-Endpunkte
# Format: Method + Path + Handler
```

#### Phase 2: Frontend API-Clients auflisten

```bash
# 1. Alle API-Client-Dateien finden
cd frontend/src
find api -name "*.ts" | sort

# 2. Für jede Datei: Welche Endpunkte werden aufgerufen?
grep -rn "client\.\(get\|post\|put\|delete\|patch\)" api/*.ts

# 3. Welche Komponenten nutzen welche APIs?
grep -rn "Api\.\|import.*from.*@/api" apps/ components/
```

#### Phase 3: Systematischer Abgleich

Erstelle eine **vollständige Mapping-Tabelle**:

| Backend Endpunkt | Method | Handler | Frontend API | Komponente | Status |
|------------------|--------|---------|--------------|------------|--------|
| `/api/v1/auth/login` | POST | LoginHandler | `authApi.login()` | Login.tsx | ✅ OK |
| `/api/v1/audit/logs` | GET | GetAuditLogs | - | - | ❌ FEHLT |
| ... | ... | ... | ... | ... | ... |

**Status-Codes:**
- ✅ **OK** - Backend implementiert, Frontend integriert, funktioniert
- ⚠️ **FEHLER** - Implementiert aber 404/500 Fehler oder Type-Fehler
- ❌ **FEHLT** - Backend existiert, aber Frontend nutzt es nicht
- 🔨 **TODO** - Backend fehlt, muss implementiert werden

**Wichtige Bereiche die zu prüfen sind:**

**Authentication & Users:**
- `/api/v1/auth/*` - Login, Logout, Refresh, Register, 2FA
- `/api/v1/users/*` - List, Get, Create, Update, Delete
- `/api/v1/groups/*` - User Groups (neu implementiert) ✅

**Storage:**
- `/api/v1/storage/disks/*` - List, Get, Format, SMART
- `/api/v1/storage/volumes/*` - List, Create, Delete, Expand
- `/api/v1/storage/shares/*` - List, Get, Create, Update, Delete, Enable/Disable
- `/api/v1/storage/stats` - Storage Statistics
- `/api/v1/storage/health` - Health Status

**Docker:**
- `/api/v1/docker/containers/*` - List, Start, Stop, Logs, Stats
- `/api/v1/docker/images/*` - List, Pull, Remove
- `/api/v1/docker/networks/*` - List, Create, Remove ⚠️ (UI Crash!)
- `/api/v1/docker/volumes/*` - List, Create, Remove
- `/api/v1/docker/compose/*` oder `/stacks/*` - Stack Management

**System:**
- `/api/v1/system/info` - System Information
- `/api/v1/system/metrics` - CPU, RAM, Network
- `/api/v1/system/health` - Health Check
- `/api/v1/system/check-updates` - Update Checker ✅ (gefixt)
- `/api/v1/system/logs` - System Logs

**Audit & Security:**
- `/api/v1/audit/*` - Audit Logs ❌ (404 Fehler)
- `/api/v1/security/*` - Security Settings ❌ (404 Fehler)

**Scheduled Tasks:**
- `/api/v1/scheduler/*` oder `/api/v1/tasks/*` - ⚠️ (Fehler)

**Alerts & Notifications:**
- `/api/v1/alerts/*` - Alerting System
- `/api/v1/notifications/*` - Notifications

**Backup:**
- `/api/v1/backup/*` - Backup Jobs, Restore

**Plugins:**
- `/api/v1/plugins/*` - Plugin Management

**Active Directory:**
- `/api/v1/ad/*` - AD Integration, User Sync

#### Phase 4: Für fehlende/fehlerhafte Endpunkte

**Für jeden ❌ FEHLT oder ⚠️ FEHLER Eintrag:**

1. **Entscheiden:**
   - Ist dieses Feature wichtig für 1.0? (P0/P1/P2)
   - Sollte es implementiert werden oder kann es warten?

2. **Wenn wichtig:**
   - Backend implementieren (falls fehlt)
   - Frontend API-Client erstellen
   - Komponente bauen die es nutzt
   - Testen

3. **Dokumentieren:**
   - Was fehlt noch bis 1.0
   - Was kann auf 1.1 oder später verschoben werden

### 3️⃣ MARKDOWN-DATEIEN UND DOKUMENTATION PRÜFEN

**Ziel:** ALLE Markdown-Dateien im Projekt auf TODOs, Feature-Requests und versprochene Funktionalität durchsuchen.

#### Schritt 1: Alle Markdown-Dateien finden

```bash
# Finde alle .md Dateien im Projekt
find . -name "*.md" -type f | grep -v node_modules | sort

# Sollte mindestens zeigen:
# - README.md
# - .github/RELEASE.md
# - NEXT_SESSION_PROMPT.md (diese Datei)
# - Weitere Dokumentation (falls vorhanden)
```

#### Schritt 2: Nach TODOs und Features suchen

```bash
# Suche in allen Markdown-Dateien nach:
# - TODO / FIXME - Offene Aufgaben
# - Features - Geplante Features
# - Roadmap - Zukünftige Entwicklung
# - Implementation - Zu implementierende Dinge
grep -rn "TODO\|FIXME\|Feature\|TODO:\|TASK\|IMPLEMENT\|Roadmap" \
  *.md .github/*.md docs/*.md --color=always 2>/dev/null
```

#### Schritt 3: Systematisch jede .md Datei durchgehen

**Für jede Markdown-Datei:**

1. **README.md**
   - Prüfe "Features" Sektion: Sind alle aufgelisteten Features implementiert?
   - Prüfe "Roadmap" Sektion: Was ist für 1.0 geplant?
   - Prüfe API-Dokumentation: Stimmt sie mit tatsächlichen Endpunkten überein?
   - Suche nach Begriffen wie "coming soon", "planned", "future"

2. **.github/RELEASE.md**
   - Prüfe Release-Checklist: Alle Punkte erfüllt?
   - Prüfe Pre-Release Anforderungen
   - Prüfe Manual Release Schritte (als Fallback)

3. **NEXT_SESSION_PROMPT.md** (diese Datei)
   - Alle Bugs (A-F) behoben?
   - Backend ↔ Frontend Mapping vollständig?
   - Alle genannten Bereiche geprüft?

4. **Weitere Dokumentation**
   - docs/API.md (falls vorhanden) - API vollständig dokumentiert?
   - docs/DEPLOYMENT.md (falls vorhanden) - Deployment-Schritte getestet?
   - CHANGELOG.md (falls vorhanden) - Aktuell?

#### Schritt 4: Feature-Vergleich mit vorherigen Sessions

**Prüfe welche Features in früheren Sessions besprochen/versprochen wurden:**

**Aus Session-History bekannte Features:**
- ✅ User Groups (implementiert)
- ✅ GitHub Actions (implementiert)
- ✅ Samba Shares (implementiert)
- ✅ Update Checker (gefixt)
- ❓ Audit Logging (404 Fehler - zu prüfen!)
- ❓ Security Settings (404 Fehler - zu prüfen!)
- ❓ Scheduled Tasks (Fehler - zu prüfen!)
- ❓ Docker Networks (Crash - zu fixen!)
- ❓ Advanced Monitoring (unvollständig - zu erweitern!)
- ❓ Disk Partitions (nicht sichtbar - zu fixen!)

**Prüfe diese weiteren Bereiche systematisch:**
```bash
# Im Code suchen nach implementierten aber ungenutzten Features
grep -r "func.*Handler\|type.*Service" backend/internal/*/

# Im Frontend suchen nach unvollständigen Komponenten
find frontend/src/apps -name "*.tsx" -exec grep -l "TODO\|FIXME\|WIP" {} \;
```

#### Schritt 5: Mapping erstellen - Was steht in Docs vs. was ist implementiert

Erstelle Tabelle:

| Feature/Bereich | In Docs erwähnt? | Backend implementiert? | Frontend implementiert? | Status | Action |
|-----------------|------------------|------------------------|-------------------------|--------|--------|
| User Management | ✓ README.md | ✅ Ja | ✅ Ja | OK | - |
| User Groups | ✓ README.md | ✅ Ja | ✅ Ja | OK | - |
| Audit Logs | ? | ? | ❌ 404 | FEHLT | Implementieren |
| Security Settings | ? | ? | ❌ 404 | FEHLT | Implementieren |
| ... | ... | ... | ... | ... | ... |

**Status-Codes:**
- ✅ **OK** - Dokumentiert, implementiert, funktioniert
- ⚠️ **UNVOLLSTÄNDIG** - Teil-implementiert oder hat Bugs
- ❌ **FEHLT** - Dokumentiert aber nicht implementiert
- 📝 **UNDOKUMENTIERT** - Implementiert aber nicht dokumentiert
- 🚫 **DEPRECATED** - Veraltet, sollte entfernt werden

#### Schritt 6: Priorisierung für Release 1.0

**Für jedes identifizierte fehlende/kaputte Feature:**

**Frage 1: Ist es kritisch für 1.0?**
- **JA** → P0 - Muss behoben werden
- **TEILWEISE** → P1 - Sollte behoben werden
- **NEIN** → P2 - Kann auf 1.1+ verschoben werden

**Frage 2: Wie viel Aufwand?**
- **Low** - < 1 Stunde
- **Medium** - 1-4 Stunden
- **High** - > 4 Stunden

**Entscheidungsmatrix:**
| Priorität | Aufwand Low | Aufwand Medium | Aufwand High |
|-----------|-------------|----------------|--------------|
| P0 | ✅ Sofort machen | ✅ Machen | ✅ Machen (aber simplify) |
| P1 | ✅ Machen | ⚠️ Wenn Zeit | ⏭️ Verschieben |
| P2 | ⚠️ Wenn Zeit | ⏭️ Verschieben | ⏭️ Verschieben |

### 4️⃣ FEHLENDE FEATURES IDENTIFIZIEREN

Prüfe welche großen Features noch fehlen oder unvollständig sind:
- Backup-Funktionalität vollständig?
- Plugin-System vollständig?
- Alerting vollständig?
- Docker-Features vollständig?
- AD-Integration vollständig?

### 5️⃣ TODO-LISTE AUFRÄUMEN

Prüfe das Projekt auf offene TODOs:
```bash
grep -r "TODO\|FIXME\|XXX\|HACK" backend/ frontend/ --include="*.go" --include="*.ts" --include="*.tsx"
```

Erstelle eine priorisierte Liste:
- **P0 - Blocker für 1.0:** Muss behoben werden
- **P1 - Wichtig:** Sollte behoben werden
- **P2 - Nice-to-have:** Kann warten

### 6️⃣ GITHUB ACTIONS RELEASE-SYSTEM VORBEREITEN

**WICHTIG:** Das Projekt hat automatische Release-Erstellung via GitHub Actions!

#### Wie das Release-System funktioniert:

1. **GitHub Actions Workflows** (bereits implementiert):
   - `.github/workflows/release.yml` - Erstellt Releases automatisch
   - `.github/workflows/ci.yml` - Testet Code bei jedem Push

2. **Release-Workflow macht automatisch:**
   - Baut Backend für Linux AMD64 und ARM64
   - Baut Frontend und packt als Tarball
   - Generiert Changelog aus Git-Commits
   - Erstellt GitHub Release mit allen Binaries
   - Generiert SHA256-Checksums

3. **Trigger:** Release wird erstellt wenn ein Tag gepusht wird
   - Tag-Format: `v1.0.0`, `v1.0.1`, etc.
   - Pre-Releases: `v1.0.0-beta.1`, `v1.0.0-rc.1`

#### Schritte für Release 1.0:

**VORBEREITUNG:**
1. Alle Bugs aus Schritt 1️⃣ behoben
2. Alle Tests laufen durch
3. Backend und Frontend bauen ohne Fehler

**RELEASE ERSTELLEN:**

```bash
# 1. Version in Code aktualisieren
# Dateien:
# - backend/cmd/stumpfworks-server/main.go → AppVersion = "1.0.0"
# - backend/internal/updates/update_service.go → CurrentVersion = "v1.0.0"
# - frontend/package.json → "version": "1.0.0"

# 2. Änderungen committen
git add backend/cmd/stumpfworks-server/main.go \
        backend/internal/updates/update_service.go \
        frontend/package.json
git commit -m "chore: bump version to 1.0.0"
git push origin <current-branch>

# 3. WICHTIG: Branch muss in main gemergt werden!
# Das Git-System erlaubt nur Tag-Push auf main Branch
# Optionen:
#   a) Pull Request erstellen (empfohlen)
#   b) Oder direkt mergen wenn du Zugriff hast

# 4. Nach Merge in main: Tag erstellen und pushen
git checkout main
git pull origin main
git tag v1.0.0 -m "Release v1.0.0 - Production Ready"
git push origin v1.0.0

# 5. GitHub Actions startet automatisch!
# Prüfe: https://github.com/Stumpf-works/stumpfworks-nas/actions
```

#### Was passiert nach Tag-Push:

1. **GitHub Actions Workflow startet** (~5-10 Minuten)
2. **Baut alle Binaries:**
   - `stumpfworks-nas-linux-amd64`
   - `stumpfworks-nas-linux-arm64`
   - `stumpfworks-nas-frontend.tar.gz`
3. **Erstellt Release auf GitHub** mit:
   - Automatischer Changelog
   - Download-Links für alle Binaries
   - SHA256 Checksums
   - Installation-Anleitung

4. **Update-Checker findet Release** automatisch!
   - Dashboard zeigt "Update available" an
   - Keine 404-Fehler mehr

#### Troubleshooting Release-Erstellung:

**Problem: "403 Forbidden" beim Tag-Push**
- **Grund:** Tag kann nur auf main Branch gepusht werden
- **Lösung:** Branch erst in main mergen, dann von main aus Tag pushen

**Problem: "Workflow not found"**
- **Grund:** `.github/workflows/release.yml` fehlt
- **Lösung:** Workflows sind bereits committed, sicherstellen dass sie auf main sind

**Problem: "Build failed"**
- **Grund:** TypeScript oder Go Build-Fehler
- **Lösung:** Lokal testen mit `cd frontend && npm run build` und `cd backend && go build`

**Problem: "No releases found" im Dashboard**
- **Grund:** Release noch nicht erstellt oder Workflow läuft noch
- **Lösung:** Warten bis Workflow fertig ist (5-10 Min), dann Dashboard refreshen

#### Validierung nach Release:

```bash
# 1. Release auf GitHub prüfen
# https://github.com/Stumpf-works/stumpfworks-nas/releases

# 2. Download-Links testen
curl -L https://github.com/Stumpf-works/stumpfworks-nas/releases/download/v1.0.0/stumpfworks-nas-linux-amd64

# 3. Update-Checker testen
# Im Dashboard → System → Check for Updates
# Sollte jetzt "v1.0.0" finden statt 404-Fehler
```

### 7️⃣ RELEASE 1.0 CHECKLISTE

**VOR dem Release:**
- [ ] Alle kritischen Bugs behoben (Schritt 1️⃣)
- [ ] Backend ↔ Frontend Mapping vollständig (Schritt 2️⃣)
- [ ] Alle Markdown-Dateien geprüft und Docs aktuell (Schritt 3️⃣)
- [ ] Fehlende Features identifiziert und priorisiert (Schritt 4️⃣)
- [ ] Wichtige TODOs erledigt (Schritt 5️⃣)
- [ ] Backend Build erfolgreich: `cd backend && go build ./cmd/stumpfworks-server`
- [ ] Frontend Build erfolgreich: `cd frontend && npm run build`
- [ ] Version auf 1.0.0 gebumpt in allen 3 Dateien
- [ ] Changelog/Release Notes vorbereitet

**WÄHREND des Releases:**
- [ ] Branch in main gemergt (oder PR erstellt)
- [ ] Tag v1.0.0 erstellt und gepusht
- [ ] GitHub Actions Workflow läuft erfolgreich
- [ ] Release auf GitHub sichtbar

**NACH dem Release:**
- [ ] Binaries herunterladbar
- [ ] Update-Checker findet v1.0.0
- [ ] Dokumentation aktualisiert
- [ ] User informieren 🎉

## TECHNISCHE DETAILS

### Repository
- **Path:** `/home/user/stumpfworks-nas`
- **Backend:** Go 1.23, `/backend/`
- **Frontend:** React + TypeScript + Vite, `/frontend/`
- **Database:** SQLite mit GORM

### Backend Struktur
```
backend/
├── cmd/stumpfworks-server/main.go  # Entry point
├── internal/
│   ├── api/                        # HTTP API
│   │   ├── router.go              # ALLE Routen hier!
│   │   └── handlers/              # Handler für Endpunkte
│   ├── storage/                    # Storage Management
│   ├── docker/                     # Docker Integration
│   ├── users/                      # User Management
│   ├── usergroups/                # User Groups (neu)
│   ├── audit/                      # Audit Logging
│   ├── scheduler/                  # Task Scheduling
│   └── ...
```

### Frontend Struktur
```
frontend/src/
├── api/                           # API Clients
│   ├── client.ts                  # Axios base client
│   ├── auth.ts
│   ├── users.ts
│   ├── groups.ts                  # User Groups (neu)
│   ├── storage.ts
│   ├── docker.ts
│   └── ...
├── apps/                          # UI Komponenten
│   ├── Dashboard/
│   ├── StorageManager/
│   ├── DockerManager/
│   ├── UserManager/               # Mit Groups Tab
│   └── ...
```

### Bekannte Systeme
- **Samba:** Shares werden direkt in `/etc/samba/smb.conf` geschrieben
- **User Groups:** Synchronisiert mit Unix-Gruppen via groupadd/usermod
- **GitHub Actions:** Automatische Releases bei Tag-Push

## ERWARTETE AUSGABE

Am Ende dieser Session sollten wir haben:
1. ✅ Alle 404-Fehler behoben (Audit, Security)
2. ✅ Docker Networks ohne Crash
3. ✅ Disks zeigen alle Partitionen
4. ✅ Dashboard Advanced Monitoring vollständig
5. ✅ Scheduled Tasks funktionieren
6. ✅ Vollständige Übersicht: Backend ↔ Frontend Mapping
7. ✅ Priorisierte TODO-Liste für 1.0
8. ✅ GitHub Actions Release-System getestet und funktionsfähig
9. ✅ System production-ready für Release 1.0
10. ✅ Release 1.0 erstellt und auf GitHub veröffentlicht

## WORKFLOW

**Schritt-für-Schritt:**
1. Beginne mit **kritischen Bugs** (A-F oben)
2. Für jeden Bug:
   - Identifiziere root cause (Backend oder Frontend?)
   - Implementiere Fix
   - Teste
   - Committe mit sinnvoller Message
3. Danach: **Backend ↔ Frontend Vollständigkeits-Check** durchführen
4. **Markdown-Dateien und Dokumentation prüfen**
5. **Fehlende Features identifizieren und priorisieren**
6. **TODO-Liste erstellen und priorisieren**
7. **GitHub Actions Release-System vorbereiten:**
   - Workflows prüfen (`.github/workflows/release.yml` und `ci.yml`)
   - Testweise einen Build laufen lassen (lokal)
   - Version auf 1.0.0 bumpen
8. **Release 1.0 erstellen:**
   - Branch in main mergen (oder PR erstellen)
   - Tag v1.0.0 erstellen und pushen
   - GitHub Actions Workflow überwachen
   - Release auf GitHub verifizieren
   - Update-Checker im Dashboard testen

## DEBUGGING TIPPS

### Frontend Fehler finden:
```bash
cd frontend
npm run build  # Zeigt TypeScript-Fehler
# Browser Console checken für Runtime-Fehler
```

### Backend Endpunkte auflisten:
```bash
cd backend
grep -r "r.Get\|r.Post\|r.Put\|r.Delete\|r.Route" internal/api/router.go
```

### 404 Fehler debuggen:
1. Backend: Prüfe ob Route in `router.go` registriert ist
2. Frontend: Prüfe API-Call (richtiger Path? richtige Methode?)
3. Network Tab in Browser: Exakte Request URL checken

### GitHub Actions Workflows testen:

**Lokal Build testen (simuliert was GitHub Actions macht):**
```bash
# Backend für Linux AMD64
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags "-s -w" \
  -o /tmp/stumpfworks-nas-linux-amd64 \
  ./cmd/stumpfworks-server

# Backend für Linux ARM64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
  -ldflags "-s -w" \
  -o /tmp/stumpfworks-nas-linux-arm64 \
  ./cmd/stumpfworks-server

# Frontend
cd frontend
npm ci
npm run build
tar -czf /tmp/stumpfworks-nas-frontend.tar.gz -C dist .

# Wenn alles erfolgreich: GitHub Actions wird auch funktionieren!
```

**Workflow-Status prüfen:**
```bash
# Nach Tag-Push, prüfe GitHub Actions Status:
# https://github.com/Stumpf-works/stumpfworks-nas/actions

# Oder via gh CLI (falls verfügbar):
gh run list --workflow=release.yml
gh run view <run-id> --log
```

**Häufige Workflow-Fehler:**
- **npm ci fails:** `package-lock.json` ist out-of-date → `npm install` lokal laufen lassen
- **go build fails:** Dependencies fehlen → `go mod tidy`
- **Permission denied:** GitHub Settings → Actions → Workflow permissions → "Read and write"
- **Tag already exists:** Alten Tag löschen mit `git tag -d v1.0.0 && git push origin :refs/tags/v1.0.0`

## WICHTIGE COMMITS AUS LETZTER SESSION

- `92391f9` - Version bump auf 0.3.0
- `8100726` - TypeScript Build-Fehler behoben
- `fb1db96` - GitHub Actions Workflows
- `fcbe5ed` - User Groups Backend
- `a9ea028` - User Groups Frontend

---

**START COMMAND:**
"Ich möchte das Stumpf.Works NAS Projekt für Release 1.0 vorbereiten. Bitte arbeite die kritischen Bugs ab und erstelle dann eine vollständige Übersicht aller Backend-Endpunkte und deren Frontend-Integration."

**GOAL:** Ein production-ready System ohne kritische Bugs, bereit für Release 1.0!
