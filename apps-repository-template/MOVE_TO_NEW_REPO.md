# 🚀 Schnellanleitung: Apps Repository erstellen

So bewegst du dieses Template ins neue GitHub Repository.

---

## ⚡ Schnellste Methode (mit Script)

```bash
cd /home/user/stumpfworks-nas/apps-repository-template
./setup-apps-repo.sh
```

Das Script:
1. Fragt dich nach dem Zielordner
2. Kopiert alle Dateien
3. Initialisiert Git
4. Erstellt Initial Commit
5. Zeigt dir die nächsten Schritte

---

## 📋 Manuelle Methode (Schritt für Schritt)

### Schritt 1: GitHub Repository erstellen

1. Gehe zu: https://github.com/Stumpf-works
2. Klicke **New Repository**
3. Name: `stumpfworks-nas-apps`
4. Description: `Official plugin repository for StumpfWorks NAS`
5. **Public** ✓
6. Klicke **Create repository**

### Schritt 2: Template kopieren

```bash
# Kopiere Template in neuen Ordner
cp -r /home/user/stumpfworks-nas/apps-repository-template /home/user/stumpfworks-nas-apps

# Wechsle ins neue Verzeichnis
cd /home/user/stumpfworks-nas-apps

# Lösche diese Anleitung (wird nicht gebraucht)
rm MOVE_TO_NEW_REPO.md
rm README_TEMPLATE.md
rm setup-apps-repo.sh
```

### Schritt 3: Git initialisieren

```bash
# Init
git init

# Alle Dateien hinzufügen
git add .

# Initial Commit
git commit -m "Initial commit: StumpfWorks NAS Apps repository"
```

### Schritt 4: Zu GitHub pushen

```bash
# Remote hinzufügen
git remote add origin https://github.com/Stumpf-works/stumpfworks-nas-apps.git

# Branch umbenennen (falls nötig)
git branch -M main

# Pushen
git push -u origin main
```

### Schritt 5: Verifizieren

Öffne: https://github.com/Stumpf-works/stumpfworks-nas-apps

Du solltest sehen:
- ✅ README.md
- ✅ CONTRIBUTING.md
- ✅ Alle Prompt-Dateien
- ✅ .github/workflows/
- ✅ scripts/
- ✅ templates/

---

## 🔌 Asterisk Plugin hinzufügen (Optional)

```bash
cd /home/user/stumpfworks-nas-apps

# Plugins-Ordner erstellen
mkdir -p plugins

# Asterisk Plugin kopieren
cp -r /home/user/stumpfworks-nas/plugins/asterisk-voip plugins/

# Release erstellen
mkdir -p releases
cd plugins/asterisk-voip
tar czf ../../releases/asterisk-voip-v1.0.0-beta.tar.gz \
  --exclude=".git" \
  --exclude="node_modules" \
  --exclude="*.log" \
  .
cd ../..

# Committen
git add plugins/ releases/
git commit -m "Add: Asterisk VoIP Plugin v1.0.0-beta"
git push

# Tag erstellen
git tag -a asterisk-voip-v1.0.0-beta -m "Release v1.0.0-beta"
git push --tags
```

### GitHub Release erstellen

1. Gehe zu **Releases** auf GitHub
2. **Draft new release**
3. Choose tag: `asterisk-voip-v1.0.0-beta`
4. Upload: `releases/asterisk-voip-v1.0.0-beta.tar.gz`
5. **Publish release**

---

## 📊 Registry generieren

```bash
cd /home/user/stumpfworks-nas-apps

# Registry generieren
python3 scripts/generate-registry.py

# Committen
git add registry.json
git commit -m "chore: generate initial registry.json"
git push
```

---

## ✅ Testen vom StumpfWorks NAS

```bash
# Registry syncen
curl -X POST http://localhost:8080/api/v1/store/sync \
  -H "Authorization: Bearer $TOKEN"

# Plugins auflisten
curl http://localhost:8080/api/v1/store/plugins

# Asterisk Plugin installieren
curl -X POST http://localhost:8080/api/v1/store/plugins/com.stumpfworks.asterisk-voip/install \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🎉 Fertig!

Das neue Repository ist jetzt live und funktionsfähig!

### Was jetzt funktioniert:

- ✅ StumpfWorks NAS kann Plugins aus dem Registry abrufen
- ✅ One-Click Installation über UI/API
- ✅ Automatische Updates
- ✅ GitHub Actions validieren neue Plugins
- ✅ Registry wird automatisch aktualisiert

### Nächste Schritte:

1. Weitere Plugins hinzufügen
2. Community einladen beizutragen
3. Plugin Store UI entwickeln
4. Dokumentation erweitern

---

**Fragen?** Siehe START_HERE.md im Template!
