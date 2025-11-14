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

**Methode:**
1. Alle Backend-Routen auflisten (aus `internal/api/router.go`)
2. Für jeden Endpunkt prüfen:
   - Gibt es einen API-Call im Frontend? (`frontend/src/api/*.ts`)
   - Wird er irgendwo verwendet? (Komponenten checken)
   - Funktioniert er? (404 Fehler? Type-Fehler?)
3. Liste erstellen mit:
   - ✅ Implementiert und funktioniert
   - ⚠️ Implementiert aber mit Fehlern
   - ❌ Nicht implementiert

**Wichtige Bereiche:**
- `/api/v1/auth/*` - Authentication
- `/api/v1/users/*` - User Management
- `/api/v1/groups/*` - User Groups ✅ (neu implementiert)
- `/api/v1/storage/*` - Storage, Shares, Disks, Volumes
- `/api/v1/docker/*` - Docker Management
- `/api/v1/system/*` - System Metrics, Updates, Health
- `/api/v1/audit/*` - Audit Logs ❌ (404)
- `/api/v1/security/*` - Security Settings ❌ (404)
- `/api/v1/scheduler/*` oder `/api/v1/tasks/*` - Scheduled Tasks ⚠️
- `/api/v1/alerts/*` - Alerting System
- `/api/v1/backup/*` - Backup Management
- `/api/v1/plugins/*` - Plugin System
- `/api/v1/ad/*` - Active Directory Integration

### 3️⃣ FEHLENDE FEATURES IDENTIFIZIEREN

Prüfe welche großen Features noch fehlen oder unvollständig sind:
- Backup-Funktionalität vollständig?
- Plugin-System vollständig?
- Alerting vollständig?
- Docker-Features vollständig?
- AD-Integration vollständig?

### 4️⃣ TODO-LISTE AUFRÄUMEN

Prüfe das Projekt auf offene TODOs:
```bash
grep -r "TODO\|FIXME\|XXX\|HACK" backend/ frontend/ --include="*.go" --include="*.ts" --include="*.tsx"
```

Erstelle eine priorisierte Liste:
- **P0 - Blocker für 1.0:** Muss behoben werden
- **P1 - Wichtig:** Sollte behoben werden
- **P2 - Nice-to-have:** Kann warten

### 5️⃣ GITHUB ACTIONS RELEASE-SYSTEM VORBEREITEN

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

### 6️⃣ RELEASE 1.0 CHECKLISTE

**VOR dem Release:**
- [ ] Alle kritischen Bugs behoben (Schritt 1️⃣)
- [ ] Backend ↔ Frontend Mapping vollständig (Schritt 2️⃣)
- [ ] Wichtige TODOs erledigt (Schritt 4️⃣)
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
3. Danach: **Vollständigkeits-Check** durchführen
4. TODO-Liste erstellen und priorisieren
5. **GitHub Actions Release-System vorbereiten:**
   - Workflows prüfen (`.github/workflows/release.yml` und `ci.yml`)
   - Testweise einen Build laufen lassen (lokal)
   - Version auf 1.0.0 bumpen
6. **Release 1.0 erstellen:**
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
