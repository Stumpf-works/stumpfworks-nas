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
- ✅ Installiert Samba
- ✅ Erstellt die Basis-Konfiguration
- ✅ Richtet das `/etc/samba/shares.d/` Verzeichnis ein
- ✅ Startet den Samba-Dienst
- ✅ Konfiguriert automatisches Include für dynamische Shares

## Manuelle Installation

Falls du das manuell machen möchtest:

### 1. Samba installieren

```bash
sudo apt-get update
sudo apt-get install -y samba samba-common-bin
```

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

### 4. Samba neu starten

```bash
sudo systemctl enable smbd nmbd
sudo systemctl restart smbd nmbd
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

## Windows verbinden

### Methode 1: Netzwerklaufwerk hinzufügen

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
