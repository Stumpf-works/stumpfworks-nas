# StumpfWorks NAS Apps Repository Structure

Dieses Dokument beschreibt die Struktur für das separate Plugin-Repository `stumpfworks-nas-apps`.

## 📦 Repository: https://github.com/Stumpf-works/stumpfworks-nas-apps

### Ordnerstruktur

```
stumpfworks-nas-apps/
├── README.md                    # Übersicht über alle Plugins
├── CONTRIBUTING.md              # Anleitung für Plugin-Entwickler
├── .github/
│   └── workflows/
│       ├── validate-plugins.yml # CI für Plugin-Validierung
│       └── update-registry.yml  # Auto-Update registry.json
│
├── registry.json                # Haupt-Registry-Datei (auto-generated)
│
├── templates/                   # Plugin-Templates
│   ├── basic-plugin/
│   ├── docker-plugin/
│   └── full-stack-plugin/
│
└── plugins/                     # Alle Plugins
    ├── asterisk-voip/
    │   ├── plugin.json
    │   ├── README.md
    │   ├── CHANGELOG.md
    │   ├── releases/
    │   │   ├── v1.0.0.tar.gz
    │   │   ├── v1.1.0.tar.gz
    │   │   └── latest.tar.gz -> v1.1.0.tar.gz
    │   └── ... (plugin files)
    │
    ├── minio-storage/
    ├── plex-media-server/
    ├── nextcloud/
    └── ... (weitere Plugins)
```

---

## 📄 registry.json Schema

Die `registry.json` wird automatisch aus allen `plugin.json` Dateien generiert.

```json
{
  "version": "1.0.0",
  "updated": "2024-12-01T10:30:00Z",
  "repository": "https://github.com/Stumpf-works/stumpfworks-nas-apps",
  "plugins": [
    {
      "id": "com.stumpfworks.asterisk-voip",
      "name": "Asterisk VoIP PBX",
      "version": "1.0.0",
      "author": "StumpfWorks Team",
      "description": "Complete VoIP telephone system",
      "icon": "📞",
      "category": "communication",
      "repository_url": "https://github.com/Stumpf-works/stumpfworks-nas-apps/tree/main/plugins/asterisk-voip",
      "download_url": "https://github.com/Stumpf-works/stumpfworks-nas-apps/releases/download/asterisk-voip-v1.0.0/asterisk-voip-v1.0.0.tar.gz",
      "homepage": "https://docs.stumpfworks.com/plugins/asterisk-voip",
      "min_nas_version": "0.1.0",
      "require_docker": true,
      "required_ports": [5060, 5061, 8088],
      "screenshots": [
        "https://raw.githubusercontent.com/Stumpf-works/stumpfworks-nas-apps/main/plugins/asterisk-voip/screenshots/dashboard.png"
      ],
      "tags": ["voip", "pbx", "telephony", "sip", "asterisk"]
    }
  ]
}
```

---

## 🔄 Workflow

### Plugin hinzufügen

1. **Ordner erstellen**
   ```bash
   mkdir -p plugins/my-plugin
   cd plugins/my-plugin
   ```

2. **plugin.json erstellen**
   ```json
   {
     "id": "com.company.my-plugin",
     "name": "My Plugin",
     "version": "1.0.0",
     ...
   }
   ```

3. **Plugin-Code entwickeln**

4. **Release erstellen**
   ```bash
   # Automatisch via GitHub Actions bei Tag-Push
   git tag my-plugin-v1.0.0
   git push origin my-plugin-v1.0.0
   ```

5. **registry.json wird automatisch aktualisiert**

---

## 🤖 GitHub Actions

### validate-plugins.yml

Validiert alle Plugins bei jedem Push:
- plugin.json Syntax
- Required fields
- Version format
- Docker Compose Syntax (falls vorhanden)

### update-registry.yml

Aktualisiert registry.json automatisch:
- Bei Tag-Push (Release)
- Einmal täglich (Scheduled)
- Manuell via Workflow Dispatch

---

## 🎯 Kategorien

Plugins werden in folgende Kategorien eingeteilt:

- **storage** - Storage & Backup (MinIO, Syncthing, etc.)
- **media** - Media Server (Plex, Jellyfin, etc.)
- **communication** - Communication (Asterisk, Matrix, etc.)
- **development** - Developer Tools (Gitea, Jenkins, etc.)
- **monitoring** - Monitoring (Prometheus, Grafana, etc.)
- **networking** - Network Tools (Pi-hole, VPN, etc.)
- **productivity** - Productivity (Nextcloud, Bitwarden, etc.)
- **security** - Security Tools
- **utilities** - Utilities & Tools

---

## 📋 Plugin-Anforderungen

Jedes Plugin muss enthalten:

- ✅ `plugin.json` (Manifest)
- ✅ `README.md` (Dokumentation)
- ✅ `CHANGELOG.md` (Versionshistorie)
- ✅ `LICENSE` (Lizenz-Datei)
- ✅ Mindestens ein Release als `.tar.gz`

Optional aber empfohlen:
- 📸 Screenshots in `screenshots/`
- 📝 Detaillierte Docs in `docs/`
- ✅ Tests
- 🐳 `docker-compose.yml` für Container-Plugins

---

## 🔐 Sicherheit

Alle Plugins werden überprüft:
- Static Code Analysis
- Malware Scan
- Dependency Check
- Docker Image Scan (falls Docker genutzt wird)

---

## 📊 Statistiken

Das Registry-System trackt:
- Download-Zahlen
- Ratings (zukünftig)
- Kompatibilität mit NAS-Versionen
- Update-Frequenz

---

## 📞 Support

- Issues: https://github.com/Stumpf-works/stumpfworks-nas-apps/issues
- Discussions: https://github.com/Stumpf-works/stumpfworks-nas-apps/discussions
- Wiki: https://github.com/Stumpf-works/stumpfworks-nas-apps/wiki
