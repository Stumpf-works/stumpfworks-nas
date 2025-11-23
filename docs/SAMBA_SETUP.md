# Samba Setup für Stumpf.Works NAS

Diese Anleitung erklärt, wie du Samba für Windows-Netzlaufwerke einrichtest.

## Warum Samba?

Samba ist ein Open-Source-Implementation des SMB/CIFS-Protokolls, das von Windows für Netzwerkfreigaben verwendet wird. Damit können Windows-PCs auf deine NAS-Shares zugreifen.

## Automatische Installation

Der einfachste Weg ist unser Setup-Script:

```bash
cd /home/user/stumpfworks-nas
sudo ./scripts/setup-samba.sh
```

Das Script:
- ✅ Installiert Samba und erforderliche Tools
- ✅ Installiert WSDD für Windows 10/11 Network Discovery
- ✅ Installiert Avahi für macOS/Linux Discovery
- ✅ Erstellt die Basis-Konfiguration mit Windows Network Discovery
- ✅ Richtet das `/etc/samba/shares.d/` Verzeichnis ein
- ✅ Startet alle erforderlichen Dienste (smbd, nmbd, wsdd, avahi)
- ✅ Konfiguriert automatisches Include für dynamische Shares

Nach der Installation sollte dein NAS automatisch sichtbar sein in:
- **Windows Explorer** → Netzwerk (als "STUMPFWORKS-NAS")
- **macOS Finder** → Netzwerk
- **Linux** Dateimanager (smb://)

## Manuelle Installation

Falls du das manuell machen möchtest:

### 1. Samba und Network Discovery Tools installieren

```bash
sudo apt-get update
sudo apt-get install -y samba samba-common-bin wsdd avahi-daemon
```

**Was wird installiert:**
- `samba` - SMB/CIFS Server
- `samba-common-bin` - Samba Utilities
- `wsdd` - Web Service Discovery für Windows 10/11
- `avahi-daemon` - Bonjour/Zeroconf für macOS/Linux

### 2. Verzeichnis für Shares erstellen

```bash
sudo mkdir -p /etc/samba/shares.d
sudo chmod 755 /etc/samba/shares.d
```

### 3. smb.conf konfigurieren

Bearbeite `/etc/samba/smb.conf` und füge in der `[global]` Sektion hinzu:

```ini
include = /etc/samba/shares.d/*.conf
```

### 4. Alle Services starten

```bash
# Samba Services
sudo systemctl enable smbd nmbd
sudo systemctl restart smbd nmbd

# Windows 10/11 Network Discovery
sudo systemctl enable wsdd
sudo systemctl start wsdd

# Avahi für macOS/Linux Discovery
sudo systemctl enable avahi-daemon
sudo systemctl start avahi-daemon
```

## Benutzer hinzufügen

Wenn du Benutzer über die Web-Oberfläche erstellst, werden diese automatisch als Samba-Benutzer angelegt. Falls du manuell einen hinzufügen möchtest:

```bash
# Linux-Benutzer erstellen (falls noch nicht vorhanden)
sudo useradd -M -s /bin/false username

# Samba-Passwort setzen
sudo smbpasswd -a username
```

## Shares erstellen

Shares werden über die Stumpf.Works Web-Oberfläche erstellt:

1. Gehe zu **Storage** → **Shares**
2. Klicke auf **New Share**
3. Wähle **Type: SMB**
4. Konfiguriere:
   - **Name**: Der Name, der in Windows angezeigt wird (z.B. `stumpfs`)
   - **Path**: Der lokale Pfad (z.B. `/mnt/data/stumpfs`)
   - **Valid Users**: Komma-getrennte Liste von Benutzernamen
   - **Guest OK**: Erlaubt Zugriff ohne Authentifizierung
   - **Read Only**: Share ist nur lesbar

Das Backend erstellt automatisch eine Konfigurationsdatei in `/etc/samba/shares.d/<name>.conf` und lädt Samba neu.

## Windows Network Discovery (Automatische Sichtbarkeit)

Dein NAS ist nach der Installation automatisch im Windows Netzwerk sichtbar! 🎉

### So findest du dein NAS in Windows:

**Windows 10/11:**
1. Öffne den **Windows Explorer** (Win+E)
2. Klicke links auf **Netzwerk**
3. Dein NAS sollte als **STUMPFWORKS-NAS** erscheinen
4. Doppelklick darauf, um verfügbare Shares zu sehen

**Warum funktioniert das?**

Das Setup-Script konfiguriert mehrere Discovery-Mechanismen:
- **WSDD** - Web Service Discovery für Windows 10/11
- **NetBIOS/WINS** - Für ältere Windows-Versionen (7/8)
- **Avahi** - Für macOS und Linux

**Wichtige Samba-Einstellungen für Network Discovery:**
```ini
local master = yes       # NAS wird Master Browser
preferred master = yes   # Bevorzugt als Master
os level = 65           # Hohe Priorität im Netzwerk
wins support = yes      # WINS-Server aktiviert
netbios name = STUMPFWORKS-NAS  # Name im Netzwerk
```

### Troubleshooting: NAS nicht sichtbar?

**1. Prüfe ob WSDD läuft (wichtig für Windows 10/11):**
```bash
sudo systemctl status wsdd
```

Falls nicht gestartet:
```bash
sudo systemctl enable wsdd
sudo systemctl start wsdd
```

**2. Prüfe NetBIOS (für ältere Windows-Versionen):**
```bash
sudo systemctl status nmbd
```

**3. Windows Network Discovery aktivieren:**
- Öffne **Einstellungen** → **Netzwerk & Internet**
- Klicke auf **Erweiterte Netzwerkeinstellungen**
- Aktiviere **Netzwerkerkennung** und **Datei- und Druckerfreigabe**

**4. Firewall-Ports prüfen:**
```bash
# Samba-Ports
sudo ufw allow 139/tcp   # NetBIOS Session
sudo ufw allow 445/tcp   # SMB
sudo ufw allow 137/udp   # NetBIOS Name Service
sudo ufw allow 138/udp   # NetBIOS Datagram

# WSDD-Port
sudo ufw allow 3702/udp  # WS-Discovery
```

**5. Manuell per IP verbinden:**

Falls die automatische Erkennung nicht funktioniert, kannst du dich direkt per IP verbinden:
```
\\192.168.178.42
```
(Ersetze mit deiner NAS-IP)

## Windows verbinden

### Methode 1: Automatisch (empfohlen)

1. Öffne **Windows Explorer** (Win+E)
2. Klicke links auf **Netzwerk**
3. Doppelklick auf **STUMPFWORKS-NAS**
4. Doppelklick auf einen Share
5. Gib bei Bedarf Benutzername und Passwort ein

### Methode 2: Netzwerklaufwerk hinzufügen

1. Öffne **Dieser PC** / **This PC**
2. Rechtsklick → **Netzlaufwerk verbinden**
3. Gib die Adresse ein:
   ```
   \\192.168.178.42\stumpfs
   ```
   (Ersetze die IP mit deiner NAS-IP und `stumpfs` mit dem Share-Namen)
4. Aktiviere **"Mit anderen Anmeldeinformationen verbinden"** falls nötig
5. Gib Benutzername und Passwort ein

### Methode 2: Explorer-Adressleiste

1. Öffne den Windows Explorer
2. Gib in die Adressleiste ein:
   ```
   \\192.168.178.42\stumpfs
   ```
3. Drücke Enter
4. Gib bei Aufforderung Benutzername und Passwort ein

## Problembehebung

### Share wird nicht angezeigt

```bash
# Prüfe ob Samba läuft
sudo systemctl status smbd

# Teste die Konfiguration
sudo testparm

# Liste verfügbare Shares
smbclient -L localhost -U%

# Prüfe die Share-Konfiguration
ls -la /etc/samba/shares.d/
cat /etc/samba/shares.d/<share-name>.conf
```

### Authentifizierung schlägt fehl

```bash
# Prüfe ob der Samba-Benutzer existiert
sudo pdbedit -L

# Setze das Passwort neu
sudo smbpasswd -a username

# Prüfe die ValidUsers in der Share-Config
cat /etc/samba/shares.d/<share-name>.conf
```

### Berechtigungsprobleme

```bash
# Prüfe Dateisystem-Berechtigungen
ls -ld /pfad/zum/share
ls -la /pfad/zum/share

# Setze korrekte Berechtigungen
sudo chown -R username:username /pfad/zum/share
sudo chmod -R 755 /pfad/zum/share
```

### Firewall

Stelle sicher, dass die Samba-Ports offen sind:

```bash
# Für UFW (Ubuntu Firewall)
sudo ufw allow Samba

# Oder manuell
sudo ufw allow 139/tcp
sudo ufw allow 445/tcp
sudo ufw allow 137/udp
sudo ufw allow 138/udp
```

### Logs prüfen

```bash
# Samba Hauptlog
sudo tail -f /var/log/samba/log.smbd

# Spezifisches Client-Log
sudo tail -f /var/log/samba/log.<client-ip>

# System-Log
sudo journalctl -u smbd -f
```

## Erweiterte Konfiguration

### Performance-Tuning

Die Standard-Konfiguration enthält bereits Performance-Optimierungen:

```ini
socket options = TCP_NODELAY IPTOS_LOWDELAY SO_RCVBUF=131072 SO_SNDBUF=131072
read raw = yes
write raw = yes
max xmit = 65535
use sendfile = yes
aio read size = 16384
aio write size = 16384
```

### Sicherheit

Für erhöhte Sicherheit kannst du in `/etc/samba/smb.conf` hinzufügen:

```ini
[global]
   # Nur sichere SMB-Versionen erlauben
   server min protocol = SMB2
   client min protocol = SMB2

   # IP-Zugriffsbeschränkung
   hosts allow = 192.168.178.0/24
   hosts deny = ALL

   # Audit-Logging
   vfs objects = full_audit
   full_audit:prefix = %u|%I|%m|%S
   full_audit:success = mkdir rename unlink rmdir pwrite
   full_audit:failure = all
```

## Nützliche Befehle

```bash
# Konfiguration testen
sudo testparm

# Konfiguration neu laden (ohne Neustart)
sudo systemctl reload smbd

# Samba komplett neu starten
sudo systemctl restart smbd nmbd

# Alle Samba-Benutzer anzeigen
sudo pdbedit -L -v

# Share-Verbindungen anzeigen
sudo smbstatus

# Samba-Version prüfen
smbd --version
```

## Automatische Integration

Das Stumpf.Works NAS Backend integriert automatisch mit Samba:

- ✅ **User-Sync**: Neue Benutzer werden automatisch als Samba-User angelegt
- ✅ **Share-Config**: Shares werden automatisch in `/etc/samba/shares.d/` konfiguriert
- ✅ **Auto-Reload**: Samba wird automatisch neu geladen nach Änderungen
- ✅ **Fehlerbehandlung**: Wenn Samba nicht installiert ist, werden Shares nur in der DB gespeichert

## Weitere Ressourcen

- [Samba Official Documentation](https://www.samba.org/samba/docs/)
- [Samba Wiki](https://wiki.samba.org/)
- [Ubuntu Samba Guide](https://ubuntu.com/server/docs/samba-file-server)
