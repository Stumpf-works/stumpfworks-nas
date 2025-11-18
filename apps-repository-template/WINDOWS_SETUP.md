# 🪟 Windows Setup Guide

Anleitung zum Erstellen des Apps-Repository von einem Windows PC aus.

---

## 📋 Voraussetzungen

### Git installieren

1. Download: https://git-scm.com/download/win
2. Installieren mit Standardeinstellungen
3. Nach Installation: Terminal neu starten

**Prüfen**:
```cmd
git --version
```

### Python installieren (optional, aber empfohlen)

1. Download: https://www.python.org/downloads/
2. **Wichtig**: Haken bei "Add Python to PATH" setzen!
3. Installieren

**Prüfen**:
```cmd
python --version
```

---

## 🚀 Setup-Methoden

### Methode 1: Batch-Script (CMD) ⚡

**Am einfachsten!**

1. **Template vom Server holen**:
   ```cmd
   # Falls du das Repo geklont hast:
   cd C:\Path\To\stumpfworks-nas\apps-repository-template

   # Oder per SSH von Server holen:
   scp -r user@nas-ip:/home/user/stumpfworks-nas/apps-repository-template C:\Temp\
   cd C:\Temp\apps-repository-template
   ```

2. **Script ausführen**:
   ```cmd
   setup-windows.bat
   ```

3. **Folge den Anweisungen** im Script!

---

### Methode 2: PowerShell-Script 🔷

**Modernere Alternative**

1. **PowerShell als Administrator öffnen**

2. **Execution Policy setzen** (einmalig):
   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

3. **Template holen** (siehe oben)

4. **Script ausführen**:
   ```powershell
   cd C:\Path\To\apps-repository-template
   .\setup-windows.ps1
   ```

---

### Methode 3: Manuell (Git Bash)

**Für volle Kontrolle**

1. **Git Bash öffnen** (Rechtsklick im Ordner → "Git Bash Here")

2. **Template kopieren**:
   ```bash
   cp -r apps-repository-template ../stumpfworks-nas-apps
   cd ../stumpfworks-nas-apps
   ```

3. **Aufräumen**:
   ```bash
   rm README_TEMPLATE.md setup-*.bat setup-*.sh setup-*.ps1 MOVE_TO_NEW_REPO.md
   ```

4. **Git initialisieren**:
   ```bash
   git init
   git add .
   git commit -m "Initial commit: StumpfWorks NAS Apps repository"
   ```

---

## 📤 Zu GitHub pushen

### Schritt 1: GitHub Repository erstellen

1. Browser öffnen: https://github.com/Stumpf-works
2. **New Repository** klicken
3. Einstellungen:
   - **Name**: `stumpfworks-nas-apps`
   - **Description**: `Official plugin repository for StumpfWorks NAS`
   - **Public** ✓
   - **Keine README, .gitignore, License hinzufügen** (haben wir schon!)
4. **Create repository** klicken

### Schritt 2: Repository verknüpfen und pushen

**In CMD/PowerShell/Git Bash**:

```bash
# Ins Repo-Verzeichnis wechseln
cd C:\Path\To\stumpfworks-nas-apps

# Remote hinzufügen
git remote add origin https://github.com/Stumpf-works/stumpfworks-nas-apps.git

# Branch umbenennen
git branch -M main

# Pushen
git push -u origin main
```

**Bei Passwort-Abfrage**: Verwende Personal Access Token (nicht Passwort!)
- Token erstellen: https://github.com/settings/tokens
- Scopes: `repo` (Full control)
- Token kopieren und als Passwort eingeben

---

## 🔌 Asterisk Plugin hinzufügen

### Vom NAS-Server holen

**Via SSH**:
```cmd
# Plugin vom Server kopieren
scp -r user@nas-ip:/home/user/stumpfworks-nas/plugins/asterisk-voip C:\Path\To\stumpfworks-nas-apps\plugins\

# Oder wenn Repo lokal liegt:
xcopy /E /I C:\Path\To\stumpfworks-nas\plugins\asterisk-voip C:\Path\To\stumpfworks-nas-apps\plugins\asterisk-voip
```

### Release erstellen

**In PowerShell** (einfacher):
```powershell
cd C:\Path\To\stumpfworks-nas-apps

# Releases-Ordner erstellen
New-Item -ItemType Directory -Path releases -Force

# Release erstellen (braucht tar - in Git Bash nutzen!)
# Wechsel zu Git Bash:
```

**In Git Bash**:
```bash
cd /c/Path/To/stumpfworks-nas-apps
mkdir -p releases

# Release erstellen
tar czf releases/asterisk-voip-v1.0.0-beta.tar.gz \
  -C plugins/asterisk-voip \
  --exclude=".git" \
  --exclude="node_modules" \
  --exclude="*.log" \
  .
```

### Committen und Pushen

```bash
git add plugins/ releases/
git commit -m "Add: Asterisk VoIP Plugin v1.0.0-beta"
git push

# Tag erstellen
git tag -a asterisk-voip-v1.0.0-beta -m "Release v1.0.0-beta"
git push --tags
```

### GitHub Release erstellen

1. Gehe zu: https://github.com/Stumpf-works/stumpfworks-nas-apps/releases
2. **Draft a new release**
3. **Choose a tag**: `asterisk-voip-v1.0.0-beta`
4. **Release title**: `Asterisk VoIP Plugin v1.0.0-beta`
5. **Describe this release**: Kurze Beschreibung
6. **Attach files**: `releases/asterisk-voip-v1.0.0-beta.tar.gz` hochladen
7. **Publish release**

---

## 📊 Registry generieren

```cmd
cd C:\Path\To\stumpfworks-nas-apps

# Registry generieren
python scripts\generate-registry.py

# Committen
git add registry.json
git commit -m "chore: generate initial registry.json"
git push
```

---

## 🧪 Testen

### Vom Windows PC (wenn NAS erreichbar):

```cmd
# Registry syncen
curl -X POST http://nas-ip:8080/api/v1/store/sync ^
  -H "Authorization: Bearer YOUR_TOKEN"

# Plugins auflisten
curl http://nas-ip:8080/api/v1/store/plugins

# Plugin installieren
curl -X POST http://nas-ip:8080/api/v1/store/plugins/com.stumpfworks.asterisk-voip/install ^
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Oder im Browser**: http://nas-ip:8080 → Plugin Store

---

## 🔧 Troubleshooting

### Git nicht gefunden

**Lösung**:
1. Git installieren: https://git-scm.com/download/win
2. Terminal neu starten
3. `git --version` prüfen

### Python nicht gefunden

**Lösung**:
1. Python installieren: https://www.python.org/downloads/
2. **Wichtig**: "Add Python to PATH" aktivieren!
3. Terminal neu starten
4. `python --version` prüfen

### PowerShell Script blockiert

**Fehlermeldung**: "execution of scripts is disabled"

**Lösung**:
```powershell
# Als Administrator:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Tar-Befehl nicht gefunden

**Lösung**: Nutze Git Bash statt CMD/PowerShell für `tar`

Oder installiere 7-Zip und nutze:
```cmd
7z a -tgzip releases\asterisk-voip-v1.0.0-beta.tar.gz plugins\asterisk-voip\*
```

### Git fragt nach Passwort

**Lösung**: Personal Access Token nutzen
1. https://github.com/settings/tokens
2. **Generate new token** (classic)
3. Scope: `repo`
4. Token kopieren
5. Als Passwort eingeben beim `git push`

**Speichern für später**:
```bash
git config --global credential.helper wincred
```

---

## 📚 Wichtige Dateien

Nach Setup solltest du haben:

```
stumpfworks-nas-apps/
├── START_HERE.md              ← Start hier!
├── README.md                  ← GitHub-Übersicht
├── CONTRIBUTING.md            ← Für Entwickler
├── HOW_TO_USE_PROMPTS.md      ← Claude Code nutzen
├── QUICK_START_PROMPT.txt     ← Session starten
├── SESSION_PROMPT.md          ← Workflows
├── CLAUDE_CODE_MASTER_PROMPT.md ← NAS Architektur
├── registry.json              ← Plugin-Index
├── .github/workflows/         ← CI/CD
├── scripts/                   ← Python-Scripts
├── templates/                 ← Plugin-Templates
└── plugins/                   ← Deine Plugins
```

---

## ✅ Checkliste

- [ ] Git installiert
- [ ] Python installiert
- [ ] Template heruntergeladen
- [ ] Setup-Script ausgeführt
- [ ] GitHub Repository erstellt
- [ ] Initial Commit gepusht
- [ ] Asterisk Plugin hinzugefügt (optional)
- [ ] Release erstellt (optional)
- [ ] Registry generiert
- [ ] Von NAS getestet

---

## 🎓 Nächste Schritte

1. **Plugins entwickeln**
   - Siehe `CONTRIBUTING.md`
   - Nutze Templates in `templates/`

2. **Claude Code nutzen**
   - Lies `HOW_TO_USE_PROMPTS.md`
   - Starte mit `QUICK_START_PROMPT.txt`

3. **Community einladen**
   - Share GitHub Repository
   - Dokumentation erweitern
   - Neue Plugins hinzufügen

---

## 💬 Support

- **Setup-Probleme**: Siehe Troubleshooting oben
- **Plugin-Entwicklung**: `CONTRIBUTING.md`
- **Claude Code**: `HOW_TO_USE_PROMPTS.md`
- **GitHub Issues**: Nach Repo-Erstellung

---

**Viel Erfolg! 🚀**

*Built with ❤️ for Windows users*
