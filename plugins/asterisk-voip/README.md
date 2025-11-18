# Asterisk VoIP Plugin für StumpfWorks NAS 📞

Vollständige VoIP-Telefonanlage (PBX) basierend auf Asterisk - Baue deine eigene professionelle Telefonanlage!

![Status](https://img.shields.io/badge/status-beta-yellow)
![Version](https://img.shields.io/badge/version-1.0.0--beta-blue)
![License](https://img.shields.io/badge/license-MIT-green)

---

## 🎯 Was ist das?

Das Asterisk VoIP Plugin verwandelt dein StumpfWorks NAS in eine **vollwertige Telefonanlage** mit professionellen Features:

- 📱 **SIP Extensions** - Interne Telefonnummern für deine Mitarbeiter
- 🌐 **SIP Trunks** - Anbindung an VoIP-Provider (sipgate, easybell, etc.)
- 📧 **Voicemail** - Anrufbeantworter mit Email-Benachrichtigung
- 🎙️ **Call Recording** - Gesprächsaufzeichnung
- 🤖 **IVR** - Interaktive Sprachmenüs ("Drücken Sie 1 für...")
- 👥 **Conference Rooms** - Telefonkonferenzen
- 📊 **Call Analytics** - Statistiken und Auswertungen
- 🌍 **WebRTC** - Telefonieren direkt im Browser

---

## 🚀 Quick Start

### Voraussetzungen

- StumpfWorks NAS v0.1.0+
- Docker & Docker Compose
- Mindestens 2GB RAM
- Offene Ports: 5060 (SIP), 10000-10099 (RTP)

### Installation

#### Über StumpfWorks UI (empfohlen)

1. Öffne StumpfWorks NAS UI
2. Navigiere zu **Plugins**
3. Suche nach "Asterisk VoIP"
4. Klicke auf **Installieren**
5. Warte auf Installation (ca. 2-3 Minuten)
6. Plugin aktivieren

#### Manuell (für Entwickler)

```bash
# 1. Plugin-Ordner kopieren
sudo cp -r plugins/asterisk-voip /var/lib/stumpfworks/plugins/

# 2. In Plugin-Ordner wechseln
cd /var/lib/stumpfworks/plugins/asterisk-voip

# 3. AMI Secret ändern (wichtig für Produktion!)
nano docker-compose.yml
# Ändere AMI_SECRET auf einen sicheren Wert

# 4. Docker Compose starten
docker-compose up -d

# 5. Logs prüfen
docker-compose logs -f
```

### Erste Schritte

1. **Dashboard öffnen**
   - In StumpfWorks UI auf das Telefon-Icon im Dock klicken

2. **Erste Extension erstellen**
   ```
   - Extension: 1000
   - Name: Administrator
   - Secret: (wird generiert)
   - Context: internal
   ```

3. **SIP-Client konfigurieren**
   - **Server**: IP deines NAS
   - **Port**: 5060
   - **Username**: 1000
   - **Passwort**: Aus Extension-Details kopieren
   - **Transport**: UDP

4. **Ersten Anruf tätigen**
   - Echo-Test: Wähle `*43`
   - Voicemail: Wähle `*97`

---

## 📖 Features im Detail

### 1. SIP Extensions

Interne Durchwahlen für deine Mitarbeiter/Abteilungen.

**Features:**
- Auto-generate sichere Passwörter
- CallerID Management
- Codec-Auswahl (ulaw, alaw, opus, gsm)
- Voicemail-Integration
- Status-Anzeige (online/offline/beschäftigt)

**Beispiel-Konfiguration:**
```json
{
  "id": "1000",
  "name": "Max Mustermann",
  "secret": "kH9mP2xL4qN7",
  "context": "internal",
  "caller_id": "Max Mustermann <1000>",
  "mailbox": "1000@default"
}
```

### 2. SIP Trunks

Verbindung zu externen VoIP-Providern.

**Unterstützte Provider:**
- sipgate
- easybell
- 1&1
- Telekom
- ... und alle anderen SIP-Provider

**Beispiel-Konfiguration (sipgate):**
```json
{
  "id": "sipgate",
  "name": "sipgate Trunk",
  "type": "friend",
  "host": "sipgate.de",
  "username": "1234567e0",
  "secret": "dein-passwort",
  "from_domain": "sipgate.de",
  "context": "from-trunk"
}
```

### 3. Dialplan

Anruf-Routing konfigurieren.

**Beispiel: Interne Anrufe**
```
Extension 1000 ruft 1001 an
→ 1001 klingelt 30 Sekunden
→ Falls keine Antwort: Voicemail
```

**Beispiel: Externe Anrufe**
```
Extension 1000 wählt 0123456789
→ Ausgehend über sipgate-Trunk
→ Rufnummer wird gewählt
```

**Beispiel: Eingehende Anrufe**
```
Anruf von extern auf +49123456789
→ IVR: "Willkommen bei Firma XYZ"
→ "Für Vertrieb drücken Sie 1"
→ "Für Support drücken Sie 2"
```

### 4. Voicemail

Anrufbeantworter mit Email-Benachrichtigung.

**Features:**
- Individuelle Begrüßungstexte
- Email mit WAV-Anhang
- Web-Interface zum Abhören
- MWI (Message Waiting Indicator)

**Verwendung:**
- Voicemail abhören: `*97`
- Eigene Voicemail: `*97` → PIN eingeben
- Fremde Voicemail: `*97` → `1001` → PIN

### 5. Call Recording

Gesprächsaufzeichnung für Qualitätssicherung.

**Modi:**
- **On-Demand**: Während Gespräch `*1` drücken
- **Automatisch**: Alle Gespräche aufzeichnen
- **Selektiv**: Nur bestimmte Extensions/Trunks

**Speicherort:** `/var/spool/asterisk/monitor`

**Format:** WAV, MP3 (komprimiert)

### 6. IVR (Sprachmenüs)

Interaktive Sprachmenüs erstellen.

**Beispiel:**
```
"Willkommen bei Musterfirma"
  1 → Vertrieb (Extension 1001)
  2 → Support (Extension 1002)
  3 → Buchhaltung (Extension 1003)
  0 → Zentrale (Extension 1000)
  * → Untermenü
```

**Visual IVR Builder** (geplant in Phase 4):
- Drag-and-drop Editor
- Flow-Chart-Darstellung
- Vorschau und Test

### 7. Conference Rooms

Telefonkonferenzen mit mehreren Teilnehmern.

**Features:**
- PIN-geschützt
- Moderator-Controls
- Aufzeichnung
- Teilnehmer-Liste
- Mute/Unmute

**Verwendung:**
```
Raum 8000 anrufen
→ PIN eingeben (falls konfiguriert)
→ Konferenz beitreten
```

### 8. Call Queues

Warteschlangen für Call-Center.

**Features:**
- Ring-Strategien (Ring All, Round Robin, etc.)
- Wartemusik
- Position in Warteschlange ansagen
- Agent-Login/Logout
- Statistiken

### 9. WebRTC

Im Browser telefonieren - ohne zusätzliche Software!

**Features:**
- Click-to-Call
- Softphone im Browser
- Video-Calls (optional)
- Screen Sharing (optional)

**Voraussetzungen:**
- HTTPS (Let's Encrypt empfohlen)
- STUN/TURN Server (für NAT-Traversal)

---

## 🔧 Konfiguration

### AMI (Asterisk Manager Interface)

Standardmäßig konfiguriert in `asterisk-config/manager.conf`:

```ini
[admin]
secret = stumpfworks2024  # ⚠️ IN PRODUKTION ÄNDERN!
permit = 127.0.0.1/255.255.255.0
permit = 192.168.0.0/255.255.0.0
```

**Wichtig:** Ändere das AMI Secret in Produktion!

```bash
# In docker-compose.yml:
environment:
  - AMI_SECRET=dein-sicheres-passwort-hier
```

### SIP Ports

Standard-Ports:
- **5060/UDP** - SIP
- **5060/TCP** - SIP (optional)
- **5061/TCP** - SIP TLS (empfohlen für Produktion)
- **10000-10099/UDP** - RTP (Audio/Video)

**Firewall öffnen:**
```bash
sudo ufw allow 5060/udp
sudo ufw allow 5060/tcp
sudo ufw allow 5061/tcp
sudo ufw allow 10000:10099/udp
```

**Router Port-Forwarding:**
```
External 5060 → NAS IP:5060 (UDP/TCP)
External 10000-10099 → NAS IP:10000-10099 (UDP)
```

### NAT Configuration

Wenn dein NAS hinter einem Router ist:

In `asterisk-config/sip.conf`:
```ini
externip = DEINE_ÖFFENTLICHE_IP
localnet = 192.168.1.0/255.255.255.0
nat = force_rport,comedia
```

**Oder automatisch via STUN:**
```ini
stunaddr = stun.l.google.com:19302
```

---

## 📊 API Dokumentation

Das Plugin bietet eine REST API auf Port **8090**.

### Basis-URL

```
http://nas-ip:8090/api/v1
```

### Authentifizierung

```bash
# Via StumpfWorks Token
curl -H "Authorization: Bearer $NAS_API_TOKEN" \
  http://localhost:8090/api/v1/extensions
```

### Wichtige Endpoints

#### Extensions

```bash
# Liste aller Extensions
GET /api/v1/extensions

# Extension erstellen
POST /api/v1/extensions
{
  "id": "1001",
  "name": "John Doe",
  "context": "internal"
}

# Extension löschen
DELETE /api/v1/extensions/1001
```

#### Aktive Anrufe

```bash
# Liste aktiver Anrufe
GET /api/v1/calls

# Anruf aufhängen
POST /api/v1/calls/{channel}/hangup

# Anruf starten
POST /api/v1/calls/originate
{
  "channel": "SIP/1000",
  "extension": "1001",
  "context": "internal"
}
```

#### Recordings

```bash
# Alle Aufzeichnungen
GET /api/v1/recordings

# Aufzeichnung abspielen
GET /api/v1/recordings/{id}

# Aufzeichnung löschen
DELETE /api/v1/recordings/{id}
```

**Vollständige API Docs:** [docs/API.md](./docs/API.md)

---

## 🔒 Sicherheit

### Wichtige Sicherheitsmaßnahmen

1. **AMI Secret ändern**
   ```bash
   # In docker-compose.yml
   AMI_SECRET=sehr-langes-zufälliges-passwort
   ```

2. **Starke SIP-Passwörter**
   ```
   Mindestens 12 Zeichen, zufällig generiert
   ```

3. **Firewall konfigurieren**
   ```bash
   # Nur benötigte Ports öffnen
   sudo ufw allow from 192.168.1.0/24 to any port 5060
   ```

4. **Fail2Ban aktivieren** (zukünftig)
   ```
   Automatisches Blockieren nach fehlgeschlagenen Login-Versuchen
   ```

5. **TLS/SRTP verwenden** (Produktion)
   ```ini
   # sip.conf
   [general]
   tlsenable=yes
   tlsbindaddr=0.0.0.0:5061
   ```

6. **VPN empfohlen**
   ```
   Greife nur über VPN auf dein NAS zu
   ```

---

## 🐛 Troubleshooting

### Container startet nicht

```bash
# Logs prüfen
docker-compose logs asterisk
docker-compose logs asterisk-manager

# Container neu starten
docker-compose restart
```

### AMI Connection Failed

```bash
# Prüfe AMI-Port
telnet localhost 5038

# Prüfe Credentials
cat docker-compose.yml | grep AMI_

# Manager Config prüfen
cat asterisk-config/manager.conf
```

### SIP-Client kann sich nicht registrieren

```bash
# Asterisk CLI öffnen
docker exec -it stumpfworks-asterisk asterisk -rvvv

# SIP Peers anzeigen
sip show peers

# SIP Debug aktivieren
sip set debug on

# Registrierungsversuch wiederholen und Logs beobachten
```

### Kein Audio bei Anrufen

**Häufigste Ursachen:**
1. **RTP-Ports nicht offen** (10000-10099/UDP)
2. **NAT-Konfiguration fehlt**
3. **Codec-Mismatch**

**Lösung:**
```bash
# RTP Ports prüfen
sudo netstat -tulpn | grep 100

# NAT in sip.conf setzen
nat=force_rport,comedia
externip=DEINE_IP

# Asterisk neu laden
docker exec stumpfworks-asterisk asterisk -rx "sip reload"
```

### WebRTC funktioniert nicht

**Voraussetzungen:**
- HTTPS (WSS benötigt TLS)
- STUN/TURN konfiguriert
- Browser unterstützt WebRTC

---

## 📈 Performance

### Systemanforderungen

| Szenario | CPU | RAM | Storage |
|----------|-----|-----|---------|
| Kleine Firma (5-10 Extensions) | 1 Core | 2 GB | 10 GB |
| Mittlere Firma (20-50 Extensions) | 2 Cores | 4 GB | 50 GB |
| Große Firma (100+ Extensions) | 4+ Cores | 8+ GB | 200+ GB |

### Gleichzeitige Anrufe

**Faustregel:** 1 Core = ~50 gleichzeitige Anrufe (ohne Transcoding)

**Mit Transcoding:** 1 Core = ~10 gleichzeitige Anrufe

### Storage

**Voicemail:** ~1 MB pro Minute
**Recordings:** ~1 MB pro Minute (WAV), ~100 KB (MP3)

**Beispiel:**
- 100 Anrufe/Tag à 5 Minuten
- Mit Recording
- = 500 MB/Tag = 15 GB/Monat

---

## 🗺️ Roadmap

### Phase 1: PoC ✅ (Aktuell)
- [x] Asterisk Container
- [x] AMI Integration
- [x] REST API Grundgerüst
- [x] Docker Compose Setup

### Phase 2: Core Features (Q2 2024)
- [ ] Extensions Management (UI)
- [ ] Trunks Management (UI)
- [ ] Dialplan Management
- [ ] Call History/CDR
- [ ] Voicemail UI

### Phase 3: Advanced Features (Q3 2024)
- [ ] Call Recording UI
- [ ] IVR Builder
- [ ] Conference Rooms
- [ ] Call Queues
- [ ] WebRTC Integration

### Phase 4: UI/UX (Q4 2024)
- [ ] Dashboard
- [ ] Visual Dialplan Builder
- [ ] Real-time Call Monitor
- [ ] Analytics/Reports
- [ ] Mobile App (optional)

### Phase 5: Production (Q1 2025)
- [ ] Security Hardening
- [ ] Fail2Ban Integration
- [ ] Monitoring/Alerting
- [ ] Backup/Restore
- [ ] Multi-Tenant Support

**Vollständiger Plan:** [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md)

---

## 🤝 Support & Community

### Hilfe erhalten

- **GitHub Issues**: https://github.com/stumpf-works/stumpfworks-nas/issues
- **Discussions**: https://github.com/stumpf-works/stumpfworks-nas/discussions
- **Discord**: https://discord.gg/stumpfworks (coming soon)

### Beitragen

Du möchtest mithelfen? Super! 🎉

1. Fork das Repository
2. Erstelle einen Feature-Branch
3. Implementiere deine Änderungen
4. Erstelle einen Pull Request

**Siehe:** [../../CONTRIBUTING.md](../../CONTRIBUTING.md)

---

## 📚 Dokumentation

- **[IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md)** - Detaillierter Implementierungsplan
- **[docs/API.md](./docs/API.md)** - API Referenz (coming soon)
- **[docs/QUICK_START.md](./docs/QUICK_START.md)** - Quick Start Guide (coming soon)
- **[docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)** - Architektur-Details (coming soon)

---

## 🙏 Credits

Dieses Plugin basiert auf:

- **Asterisk** - https://www.asterisk.org/
- **Docker Image** - https://github.com/andrius/asterisk
- **SIP.js** - https://sipjs.com/ (für WebRTC)
- **StumpfWorks NAS** - https://github.com/stumpf-works/stumpfworks-nas

---

## 📄 Lizenz

MIT License - Siehe [LICENSE](../../LICENSE)

---

## ⚠️ Haftungsausschluss

Dieses Plugin befindet sich in **Beta**. Verwende es nicht in produktionskritischen Umgebungen ohne ausreichende Tests!

**Hinweis:** Die Nutzung von VoIP-Diensten unterliegt gesetzlichen Regelungen. Informiere dich über die Gesetze in deinem Land, insbesondere bezüglich:
- Gesprächsaufzeichnung (Einwilligung erforderlich!)
- Notrufnummern (110, 112 müssen funktionieren!)
- Datenschutz (DSGVO)

---

**Viel Spaß mit deiner eigenen Telefonanlage! 📞🚀**
