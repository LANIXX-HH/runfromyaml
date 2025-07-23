# Amazon Q Integration für runfromyaml

Diese Dokumentation beschreibt, wie der runfromyaml MCP Server mit Amazon Q integriert werden kann, um Workflows durch natürliche Sprache zu generieren und zu verwalten.

## Übersicht

Der runfromyaml MCP Server wurde speziell für die Integration mit Amazon Q optimiert, um folgende Funktionen bereitzustellen:

- Generierung von Workflows aus natürlichsprachlichen Beschreibungen
- Validierung von Workflow-Strukturen
- Erklärung von Workflows vor der Ausführung
- Template-basierte Workflow-Erstellung
- Ausführung von Workflows (mit Sicherheitsbestätigungen)

## Installation und Konfiguration

### 1. Installation von runfromyaml

Stelle sicher, dass runfromyaml installiert ist:

```bash
# Installation über Go
go install github.com/lanixx-hh/runfromyaml@latest

# Oder über direkten Download
curl -L https://github.com/LANIXX-HH/runfromyaml/releases/latest/download/runfromyaml-$(uname -s)-$(uname -m) -o /usr/local/bin/runfromyaml
chmod +x /usr/local/bin/runfromyaml
```

### 2. Amazon Q Konfiguration

Füge den folgenden Eintrag zur Datei `~/.aws/amazonq/mcp.json` hinzu:

```json
{
  "mcpServers": {
    "runfromyaml-workflow-server": {
      "command": "runfromyaml",
      "args": ["--mcp", "--no-file"],
      "env": {
        "DEBUG": "false"
      },
      "disabled": false,
      "autoApprove": ["generate_workflow", "validate_workflow", "explain_workflow"]
    }
  }
}
```

Falls die Datei bereits andere MCP-Server enthält, füge nur den `"runfromyaml-workflow-server"` Eintrag zum `mcpServers` Objekt hinzu.

### 3. Konfigurationserklärung

- **Server-Name**: `runfromyaml-workflow-server` - Der Name des MCP-Servers
- **Befehl**: `runfromyaml` - Der Befehl zum Starten des Tools
- **Argumente**: 
  - `--mcp`: Startet den MCP-Server-Modus
  - `--no-file`: Deaktiviert die Suche nach einer Standard-YAML-Datei
- **Umgebungsvariablen**:
  - `DEBUG`: "false" - Deaktiviert den Debug-Modus für eine sauberere Ausgabe
- **Status**: `disabled: false` - Der Server ist aktiviert
- **AutoApprove**: Sichere Tools, die ohne Bestätigung ausgeführt werden können:
  - `generate_workflow`: Generiert Workflows ohne Ausführung
  - `validate_workflow`: Validiert Workflow-Syntax
  - `explain_workflow`: Erklärt, was ein Workflow tun würde

## Verfügbare Tools

Der runfromyaml MCP Server stellt folgende Tools für Amazon Q bereit:

### 1. generate_workflow

Generiert Workflow-YAML aus natürlichsprachlicher Beschreibung ohne Ausführung.

**Beispiel-Anfrage:**
"Erstelle einen Workflow für eine Docker-basierte Web-Anwendung mit PostgreSQL"

### 2. validate_workflow

Validiert die Struktur und Syntax eines Workflow-YAML.

**Beispiel-Anfrage:**
"Validiere diesen Workflow: [YAML einfügen]"

### 3. explain_workflow

Erklärt, was ein Workflow tun wird, ohne ihn auszuführen.

**Beispiel-Anfrage:**
"Erkläre, was dieser Workflow tut: [YAML einfügen]"

### 4. workflow_from_template

Generiert einen Workflow aus einer vordefinierten Vorlage.

**Beispiel-Anfrage:**
"Erstelle einen Web-App-Workflow mit Port 3000"

### 5. generate_and_execute_workflow (erfordert Bestätigung)

Generiert einen Workflow aus natürlichsprachlicher Beschreibung und führt ihn aus.

**Beispiel-Anfrage:**
"Erstelle und führe einen Workflow aus, der eine Docker-Umgebung einrichtet"

### 6. execute_existing_workflow (erfordert Bestätigung)

Führt einen bestehenden YAML-Workflow aus.

**Beispiel-Anfrage:**
"Führe diesen Workflow aus: [YAML einfügen]"

## Sicherheitshinweise

- Die Tools `generate_and_execute_workflow` und `execute_existing_workflow` sind nicht für automatische Genehmigung konfiguriert, da sie tatsächlich Befehle ausführen. Amazon Q wird vor deren Ausführung um Bestätigung bitten.
- Überprüfe generierte Workflows immer, bevor du sie ausführst, besonders wenn sie Systemänderungen vornehmen.
- Verwende `explain_workflow`, um zu verstehen, was ein Workflow tun wird, bevor du ihn ausführst.

## Beispiel-Workflows

### Docker-basierte Web-Anwendung

```yaml
logging:
  - level: info
  - output: stdout
env:
  - key: APP_PORT
    value: "3000"
cmd:
  - type: docker-compose
    name: docker-compose-setup
    desc: Docker Compose setup for web application
    expandenv: true
    dcoptions:
      - -f
      - docker-compose.yml
    command: up
    cmdoptions:
      - -d
    service: ""
    values: []
  - type: conf
    name: web-config
    desc: Web application configuration
    confdest: ./app.conf
    confperm: 0644
    confdata: |
      # Web Application Configuration
      port=3000
      host=0.0.0.0
```

### Datenbank-Setup

```yaml
logging:
  - level: info
  - output: stdout
env:
  - key: DB_HOST
    value: localhost
  - key: DB_PORT
    value: "5432"
cmd:
  - type: docker
    name: postgres-setup
    desc: Setup PostgreSQL database
    expandenv: true
    command: run
    container: postgres:13
    values:
      - echo 'PostgreSQL container started'
      - psql -U postgres -c 'CREATE DATABASE myapp;'
```

## Fehlerbehebung

### Häufige Probleme

1. **Server wird nicht erkannt**
   - Überprüfe, ob runfromyaml im PATH verfügbar ist
   - Stelle sicher, dass die mcp.json korrekt formatiert ist
   - Starte Amazon Q neu

2. **Berechtigungsprobleme**
   - Stelle sicher, dass runfromyaml ausführbare Berechtigungen hat
   - Überprüfe, ob der Benutzer die notwendigen Berechtigungen hat

3. **Workflow-Ausführungsfehler**
   - Aktiviere den Debug-Modus mit `"DEBUG": "true"` in der Konfiguration
   - Überprüfe die generierten Workflows auf Syntax- oder Logikfehler

## Weitere Ressourcen

- **Workflow-Templates**: Der Server bietet vordefinierte Templates für gängige Szenarien
- **Beispiel-Workflows**: Verfügbar über die Ressource `workflow://examples`
- **Best Practices**: Verfügbar über die Ressource `workflow://best-practices`
- **Schema-Validierung**: Verfügbar über die Ressource `workflow://schema`
