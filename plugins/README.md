# StumpfWorks NAS - Plugins

Dieser Ordner enthält offizielle und Community-Plugins für StumpfWorks NAS.

## 📦 Verfügbare Plugins

### Offiziell
- **[asterisk-voip](./asterisk-voip/)** - Vollständige VoIP-Telefonanlage mit Asterisk PBX

## 🎯 Was sind Plugins?

Plugins erweitern StumpfWorks NAS um zusätzliche Funktionalität ohne den Core zu verändern. Sie laufen als separate Prozesse und kommunizieren über die StumpfWorks API.

## 🏗️ Plugin-Architektur

### Plugin-Typen

1. **Service Plugins** (z.B. asterisk-voip)
   - Stellen zusätzliche Dienste bereit
   - Laufen als dauerhafte Hintergrundprozesse
   - Haben eigene APIs und UIs

2. **Integration Plugins** (z.B. Synology Migration)
   - Integrieren externe Systeme
   - Meist event-basiert oder on-demand

3. **Utility Plugins** (z.B. Backup-Tools)
   - Bieten zusätzliche Werkzeuge
   - Können periodisch oder manuell laufen

### Plugin-Struktur

Jedes Plugin hat folgende Struktur:

```
plugin-name/
├── plugin.json              # Plugin-Manifest (erforderlich)
├── README.md                # Dokumentation
├── IMPLEMENTATION_PLAN.md   # Implementierungsplan (optional)
├── docker-compose.yml       # Docker Services (optional)
├── backend/                 # Backend-Code
│   ├── main.go             # Entry Point
│   ├── go.mod              # Go Dependencies
│   └── ...
├── frontend/               # Frontend-Code (optional)
│   └── ...
└── config/                 # Konfigurationsdateien
    └── ...
```

## 📋 Plugin-Manifest (plugin.json)

Jedes Plugin benötigt eine `plugin.json` Datei:

```json
{
  "id": "com.company.plugin-name",
  "name": "Plugin Display Name",
  "version": "1.0.0",
  "author": "Author Name",
  "description": "Plugin description",
  "icon": "🔌",
  "entryPoint": "executable-name",
  "requires": {
    "docker": true,
    "ports": [8080, 9000],
    "storage": "1GB",
    "minNasVersion": "0.1.0"
  },
  "config": {
    "key": "default-value"
  }
}
```

### Manifest-Felder

- **id**: Eindeutige Plugin-ID (reverse domain notation)
- **name**: Anzeigename für UI
- **version**: Semantic Versioning (x.y.z)
- **author**: Plugin-Entwickler
- **description**: Kurzbeschreibung
- **icon**: Emoji oder Icon-Pfad
- **entryPoint**: Ausführbare Datei (relativ zum Plugin-Ordner)
- **requires**: Systemanforderungen
  - `docker`: Benötigt Docker
  - `ports`: Benötigte Netzwerk-Ports
  - `storage`: Minimaler Speicherbedarf
  - `minNasVersion`: Minimale StumpfWorks NAS Version

## 🔌 Plugin API

### Umgebungsvariablen

Plugins erhalten folgende Umgebungsvariablen:

```bash
PLUGIN_ID=com.company.plugin-name
PLUGIN_DIR=/var/lib/stumpfworks/plugins/plugin-name
NAS_API_URL=http://localhost:8080/api/v1
NAS_API_TOKEN=<auth-token>
```

### StumpfWorks API Zugriff

Plugins können die StumpfWorks REST API nutzen:

```go
// Beispiel: User abrufen
resp, err := http.Get(os.Getenv("NAS_API_URL") + "/users")
```

Verfügbare APIs:
- `/api/v1/users` - Benutzerverwaltung
- `/api/v1/storage` - Storage-Operationen
- `/api/v1/docker` - Docker-Management
- `/api/v1/syslib` - System Library (ZFS, Samba, etc.)

## 🚀 Plugin-Entwicklung

Siehe [DEVELOPMENT.md](./DEVELOPMENT.md) für eine detaillierte Anleitung.

### Quick Start

1. **Plugin-Ordner erstellen**
   ```bash
   mkdir -p plugins/my-plugin
   cd plugins/my-plugin
   ```

2. **plugin.json erstellen**
   ```bash
   cat > plugin.json <<EOF
   {
     "id": "com.mycompany.my-plugin",
     "name": "My Plugin",
     "version": "1.0.0",
     "entryPoint": "my-plugin"
   }
   EOF
   ```

3. **Backend entwickeln**
   ```bash
   mkdir backend && cd backend
   go mod init my-plugin
   # ... entwickeln
   go build -o ../my-plugin
   ```

4. **Plugin installieren**
   ```bash
   cp -r . /var/lib/stumpfworks/plugins/my-plugin/
   # Plugin über API oder UI aktivieren
   ```

## 📚 Weitere Ressourcen

- [Plugin Development Guide](./DEVELOPMENT.md)
- [Asterisk VoIP Plugin PoC](./asterisk-voip/)
- [StumpfWorks API Dokumentation](../docs/API.md)
- [Plugin SDK Reference](../docs/PLUGIN_SDK.md)

## 🤝 Beitragen

Möchtest du ein Plugin beitragen?

1. Fork das Repository
2. Erstelle dein Plugin in `plugins/your-plugin/`
3. Dokumentiere es ausführlich
4. Erstelle einen Pull Request

## 📄 Lizenz

Plugins können eigene Lizenzen haben. Bitte beachte die jeweilige LICENSE-Datei im Plugin-Ordner.
