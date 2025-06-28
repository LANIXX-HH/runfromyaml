# SSH Key Issue Fix - macOS Extended Attributes Problem

## Problem

Nach dem Aktivieren von SSH auf macOS erscheint folgende Warnung:
```
Warning: Identity file  /Users/anatoli.lichii/.ssh/id_rsa-localhost not accessible: No such file or directory.
```

Obwohl die Datei existiert und die richtigen Berechtigungen hat.

## Root Cause

Das Problem liegt an **Extended Attributes** (erweiterte Attribute) die macOS automatisch zu SSH-Schlüsseln hinzufügt:

```bash
$ ls -la ~/.ssh/id_rsa-localhost
-rw-------@ 1 anatoli.lichii  staff  3422 Jun 28 18:52 /Users/anatoli.lichii/.ssh/id_rsa-localhost
#         ^ Das @ zeigt erweiterte Attribute an

$ xattr -l ~/.ssh/id_rsa-localhost
com.apple.provenance: \u0001\u0002
```

Das `com.apple.provenance` Attribut blockiert SSH-Zugriff auf den Schlüssel.

## Lösungen

### Lösung 1: Extended Attributes entfernen (Empfohlen)

```bash
# Entferne alle erweiterten Attribute
xattr -c ~/.ssh/id_rsa-localhost

# Oder spezifisch das provenance Attribut
xattr -d com.apple.provenance ~/.ssh/id_rsa-localhost
```

### Lösung 2: SSH-Schlüssel neu erstellen

```bash
# Alten Schlüssel entfernen
rm -f ~/.ssh/id_rsa-localhost*

# Neuen Schlüssel erstellen
ssh-keygen -t rsa -b 4096 -N '' -f ~/.ssh/id_rsa-localhost

# Extended Attributes sofort entfernen
xattr -c ~/.ssh/id_rsa-localhost

# Öffentlichen Schlüssel zu authorized_keys hinzufügen
cat ~/.ssh/id_rsa-localhost.pub >> ~/.ssh/authorized_keys
```

### Lösung 3: SSH Remote Login aktivieren

**Via Terminal (benötigt Full Disk Access):**
```bash
sudo systemsetup -setremotelogin on
```

**Via System Preferences (Einfacher):**
1. System Preferences → Sharing
2. Remote Login ✅ aktivieren
3. Benutzer auswählen die SSH-Zugriff haben sollen

## Verification

Nach dem Fix sollte SSH funktionieren:

```bash
# Test SSH-Verbindung
ssh -o ConnectTimeout=5 -i ~/.ssh/id_rsa-localhost localhost echo "SSH works!"

# Test mit runfromyaml
./runfromyaml --file commands.yaml
```

## runfromyaml commands.yaml Update

Um das Problem zu vermeiden, kann man die commands.yaml anpassen:

```yaml
cmd:
  - type: shell
    expandenv: true
    desc: "erstelle SSH-Schlüssel und entferne Extended Attributes"
    name: "create-ssh-key-clean"
    values:
      - ls $HOME/.ssh/id_rsa-localhost || ssh-keygen -t rsa -b 4096 -N '' -f $HOME/.ssh/id_rsa-localhost
      - xattr -c $HOME/.ssh/id_rsa-localhost 2>/dev/null || true

  - type: shell
    expandenv: true
    desc: "füge public key zu authorized_keys hinzu"
    name: "setup-authorized-keys"
    values:
      - grep -f $HOME/.ssh/id_rsa-localhost.pub $HOME/.ssh/authorized_keys || cat $HOME/.ssh/id_rsa-localhost.pub >> $HOME/.ssh/authorized_keys

  - type: "ssh"
    expandenv: true
    name: "ssh-run"
    desc: "run command via ssh connection"
    user: $USER
    host: localhost
    port: 22
    options:
      - -i $HOME/.ssh/id_rsa-localhost
      - -o StrictHostKeyChecking=no
    values:
      - uname -a
      - pwd
```

## Warum passiert das?

macOS fügt automatisch das `com.apple.provenance` Attribut zu Dateien hinzu, die von bestimmten Programmen erstellt werden. SSH ist sehr sicherheitsbewusst und verweigert den Zugriff auf Schlüssel mit unbekannten Extended Attributes.

## Prevention

Um das Problem in Zukunft zu vermeiden:

1. **Immer Extended Attributes nach SSH-Schlüssel-Erstellung entfernen**
2. **SSH-Schlüssel außerhalb von automatisierten Tools erstellen**
3. **Regelmäßig SSH-Konfiguration testen**

## Status

- ✅ **SSH expandenv Fix funktioniert**: Umgebungsvariablen werden korrekt expandiert
- ⚠️ **SSH-Verbindung Problem**: Extended Attributes blockieren SSH-Zugriff
- 🔧 **Lösung verfügbar**: Extended Attributes entfernen löst das Problem
