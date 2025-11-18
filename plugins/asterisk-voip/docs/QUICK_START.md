# Quick Start Guide - Asterisk VoIP Plugin

Schnellanleitung, um in 15 Minuten deine erste Telefonanlage zum Laufen zu bringen!

## 🎯 Ziel

Am Ende dieses Guides kannst du:
- ✅ Zwischen zwei Telefonen intern telefonieren
- ✅ Voicemail nutzen
- ✅ Einen Echo-Test durchführen

---

## 📋 Voraussetzungen

- StumpfWorks NAS installiert und läuft
- Docker ist verfügbar
- Mindestens 2GB freier RAM
- 2 SIP-fähige Geräte (z.B. Softphones auf deinem Handy)

**Empfohlene SIP-Clients:**
- **Android**: Linphone, Zoiper
- **iOS**: Linphone, Zoiper
- **Desktop**: Linphone, Jitsi, MicroSIP (Windows)

---

## 🚀 Schritt-für-Schritt

### Schritt 1: Plugin starten (2 Minuten)

```bash
# In StumpfWorks NAS Verzeichnis wechseln
cd /home/user/stumpfworks-nas/plugins/asterisk-voip

# Docker Container starten
docker-compose up -d

# Warten bis Asterisk gestartet ist (ca. 30 Sekunden)
docker-compose logs -f asterisk

# Warte auf: "Asterisk Ready"
# Dann STRG+C drücken
```

### Schritt 2: Extensions erstellen (3 Minuten)

**Extension 1000 (bereits vorhanden):**
```bash
curl -X GET http://localhost:8090/api/v1/extensions
```

**Extension 1001 erstellen:**
```bash
curl -X POST http://localhost:8090/api/v1/extensions \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1001",
    "name": "Test User",
    "secret": "test1001",
    "context": "internal",
    "caller_id": "Test User <1001>",
    "mailbox": "1001@default"
  }'
```

**Oder manuell in Asterisk-Config:**

1. Öffne `asterisk-config/sip.conf`
2. Füge hinzu:
```ini
[1001]
type = friend
secret = test1001
host = dynamic
context = internal
qualify = yes
dtmfmode = rfc2833
canreinvite = no
mailbox = 1001@default
```

3. Config neu laden:
```bash
docker exec stumpfworks-asterisk asterisk -rx "sip reload"
```

### Schritt 3: SIP-Clients konfigurieren (5 Minuten)

**Gerät 1 (Extension 1000):**
1. SIP-App öffnen (z.B. Linphone)
2. **Neues Konto hinzufügen**:
   - **Benutzername**: `1000`
   - **Passwort**: `CHANGE_ME` (siehe sip.conf)
   - **Domain/Server**: `[IP deines NAS]`
   - **Port**: `5060`
   - **Transport**: `UDP`
3. **Speichern** und warten auf "Registriert"

**Gerät 2 (Extension 1001):**
1. Gleiche Schritte wie oben
2. Aber mit:
   - **Benutzername**: `1001`
   - **Passwort**: `test1001`

### Schritt 4: Ersten Anruf testen (2 Minuten)

**Von Extension 1000 zu 1001:**
1. In SIP-App auf Gerät 1
2. Nummer wählen: `1001`
3. **Anrufen** drücken
4. Auf Gerät 2 sollte es klingeln!
5. **Annehmen** und sprechen

**Gratulation! 🎉 Dein erster VoIP-Anruf!**

### Schritt 5: Echo-Test (1 Minute)

**Test ob Audio funktioniert:**
1. Von beliebiger Extension
2. Nummer wählen: `*43`
3. Du solltest eine Ansage hören
4. Dann wird alles was du sagst zurückgespielt (Echo)
5. **Aufhängen** mit der roten Taste

### Schritt 6: Voicemail testen (2 Minuten)

**Voicemail-Box konfigurieren:**

1. Öffne `asterisk-config/voicemail.conf`
2. Extension 1001 sollte bereits existieren:
```ini
[default]
1001 => 1234,Test User,test@example.com
```

3. Falls nicht, hinzufügen und Asterisk neu laden:
```bash
docker exec stumpfworks-asterisk asterisk -rx "module reload app_voicemail.so"
```

**Voicemail hinterlassen:**
1. Von Extension 1000 anrufen: `1001`
2. **NICHT** annehmen auf 1001
3. Nach 30 Sekunden sollte Voicemail aktivieren
4. Nachricht sprechen
5. Auflegen

**Voicemail abhören:**
1. Von Extension 1001
2. Wähle `*97`
3. Mailbox eingeben: `1001`
4. PIN eingeben: `1234`
5. Nachricht anhören

---

## ✅ Checkliste

- [ ] Docker Container laufen
- [ ] Extension 1000 registriert
- [ ] Extension 1001 registriert
- [ ] Anruf von 1000 → 1001 funktioniert
- [ ] Echo-Test funktioniert (*43)
- [ ] Voicemail funktioniert (*97)

---

## 🔍 Troubleshooting

### Problem: SIP-Client registriert sich nicht

**Lösung 1: Firewall prüfen**
```bash
sudo ufw allow 5060/udp
sudo ufw allow 5060/tcp
```

**Lösung 2: Asterisk Logs prüfen**
```bash
docker exec -it stumpfworks-asterisk asterisk -rvvv
# Dann in Asterisk CLI:
sip set debug on
# Registrierungsversuch wiederholen und Logs beobachten
```

**Lösung 3: SIP-Config prüfen**
```bash
docker exec stumpfworks-asterisk asterisk -rx "sip show peers"
# 1000 und 1001 sollten auftauchen
```

### Problem: Kein Audio bei Anrufen

**Lösung: RTP-Ports öffnen**
```bash
sudo ufw allow 10000:10099/udp
```

**In sip.conf NAT-Settings prüfen:**
```ini
[general]
nat=force_rport,comedia
externip=DEINE_ÖFFENTLICHE_IP  # oder weglassen wenn nur lokal
```

### Problem: Echo-Test funktioniert nicht

**Asterisk CLI öffnen und prüfen:**
```bash
docker exec -it stumpfworks-asterisk asterisk -rvvv
# Dann Echo-Test anrufen und Logs beobachten
```

**Extensions.conf prüfen:**
```bash
cat asterisk-config/extensions.conf | grep -A5 "exten => \*43"
```

---

## 📚 Nächste Schritte

Jetzt wo die Basics funktionieren:

1. **Externe Telefonie einrichten**
   - VoIP-Provider-Account besorgen (z.B. sipgate)
   - Trunk konfigurieren
   - Externe Anrufe tätigen/empfangen

2. **IVR (Sprachmenü) erstellen**
   - Siehe [IMPLEMENTATION_PLAN.md](../IMPLEMENTATION_PLAN.md#32-ivr)

3. **Anrufaufzeichnung aktivieren**
   - Siehe [README.md](../README.md#5-call-recording)

4. **WebRTC Softphone nutzen**
   - Direkt im Browser telefonieren
   - Siehe Phase 3.5 im Implementation Plan

---

## 🎓 Lernressourcen

**Asterisk Basics:**
- [Asterisk Dokumentation](https://www.asterisk.org/documentation/)
- [Asterisk: The Definitive Guide (kostenlos)](https://www.asteriskdocs.org/)

**SIP Protocol:**
- [SIP für Einsteiger](https://www.voip-info.org/sip/)

**Dialplan:**
- [Dialplan Basics](https://www.voip-info.org/asterisk-dialplan/)

---

## 💬 Hilfe benötigt?

- **GitHub Issues**: https://github.com/stumpf-works/stumpfworks-nas/issues
- **Discussions**: https://github.com/stumpf-works/stumpfworks-nas/discussions

---

**Viel Erfolg! 📞**
