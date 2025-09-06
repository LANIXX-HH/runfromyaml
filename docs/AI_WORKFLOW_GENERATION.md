# AI-Basierte Workflow-Generierung

## Problem

Die ursprüngliche MCP-Implementierung verwendete nur einfaches Pattern-Matching für die Workflow-Generierung. Dies führte dazu, dass:

1. **Der MCP-Agent (Amazon Q) die meiste Arbeit machen musste** - Der Server lieferte nur primitive Vorlagen
2. **Keine echte AI-Integration** - Nur statische String-Suche nach Schlüsselwörtern
3. **Begrenzte Flexibilität** - Workflows waren vorhersagbar und nicht an spezifische Anforderungen angepasst

## Lösung: Echte AI-Integration

### Neue Architektur

```
User Prompt → MCP Server → OpenAI API → Intelligente YAML-Generierung → Validierung → Ausführung
```

### Kernkomponenten

#### 1. AIWorkflowGenerator (`pkg/mcp/ai_workflow_generator.go`)

```go
type AIWorkflowGenerator struct {
    openaiClient *openai.Client
    enabled      bool
}
```

**Funktionen:**
- `GenerateWorkflowFromDescription()` - Hauptfunktion für AI-basierte Generierung
- `createWorkflowPrompt()` - Erstellt detaillierte Prompts für OpenAI
- `parseAIResponse()` - Parst und validiert AI-Antworten
- `validateAndEnhanceWorkflow()` - Stellt sicher, dass generierte Workflows gültig sind
- `generateFallbackWorkflow()` - Fallback auf Pattern-Matching bei AI-Fehlern

#### 2. Erweiterte MCP-Server Integration

Der MCP-Server wurde erweitert um:
- AI-Workflow-Generator-Instanz
- Automatische Konfiguration basierend auf verfügbaren AI-Parametern
- Intelligente Fallback-Mechanismen

### Verwendung

#### 1. MCP-Server mit AI starten

```bash
# Mit OpenAI API Key
runfromyaml --mcp --ai-key "sk-your-api-key" --ai-model "gpt-4" --debug

# Ohne AI (Fallback auf Pattern-Matching)
runfromyaml --mcp
```

#### 2. AI-Prompt-Engineering

Der AI-Generator verwendet detaillierte Prompts:

```
Generate a complete runfromyaml workflow YAML based on this description: "Setup a Node.js web application with Docker"

REQUIREMENTS:
1. Generate valid YAML that follows the runfromyaml schema
2. Include appropriate logging configuration
3. Add environment variables if needed
4. Create detailed command blocks with proper types
5. Use descriptive names and descriptions
6. Include error handling where appropriate

AVAILABLE BLOCK TYPES:
- shell: Execute shell commands
- exec: Execute system commands directly  
- docker: Run Docker containers
- docker-compose: Docker Compose operations
- ssh: Remote SSH commands
- conf: Create configuration files
```

#### 3. Intelligente Workflow-Generierung

**Beispiel-Input:** "Setup a PostgreSQL database with Docker Compose and create initial tables"

**AI-Generated Output:**
```yaml
logging:
  - level: info
  - output: stdout

env:
  - key: POSTGRES_DB
    value: myapp_db
  - key: POSTGRES_USER
    value: postgres
  - key: POSTGRES_PASSWORD
    value: secure_password

cmd:
  - type: conf
    name: docker-compose-config
    desc: "Create Docker Compose configuration for PostgreSQL"
    confdest: ./docker-compose.yml
    confperm: 0644
    confdata: |
      version: '3.8'
      services:
        postgres:
          image: postgres:15
          environment:
            POSTGRES_DB: ${POSTGRES_DB}
            POSTGRES_USER: ${POSTGRES_USER}
            POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
          ports:
            - "5432:5432"
          volumes:
            - postgres_data:/var/lib/postgresql/data
            - ./init.sql:/docker-entrypoint-initdb.d/init.sql
      volumes:
        postgres_data:

  - type: conf
    name: init-sql
    desc: "Create initial database schema"
    confdest: ./init.sql
    confperm: 0644
    confdata: |
      CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );
      
      CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        title VARCHAR(200) NOT NULL,
        content TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );

  - type: docker-compose
    name: start-database
    desc: "Start PostgreSQL database with Docker Compose"
    expandenv: true
    dcoptions:
      - -f
      - docker-compose.yml
    command: up
    cmdoptions:
      - -d
    service: ""
    values: []

  - type: shell
    name: verify-database
    desc: "Verify database connection and tables"
    expandenv: true
    values:
      - sleep 10
      - docker exec $(docker-compose ps -q postgres) psql -U $POSTGRES_USER -d $POSTGRES_DB -c "\dt"
      - echo "✅ Database setup complete!"
```

### Vorteile der neuen Implementierung

#### 1. **Echte AI-Intelligenz**
- Kontextbewusste Workflow-Generierung
- Anpassung an spezifische Anforderungen
- Realistische Konfigurationen und Befehle

#### 2. **Robuste Fallback-Mechanismen**
- Automatischer Fallback auf Pattern-Matching bei AI-Fehlern
- Validierung und Verbesserung von AI-generierten Workflows
- Fehlerbehandlung auf mehreren Ebenen

#### 3. **Flexibilität**
- Unterstützung verschiedener OpenAI-Modelle
- Konfigurierbare AI-Parameter
- Funktioniert mit und ohne AI-Integration

#### 4. **Produktionsreife**
- Umfassende Validierung generierter Workflows
- Sicherheitsüberprüfungen
- Detaillierte Fehlerberichterstattung

### Konfiguration

#### Environment Variables
```bash
export OPENAI_API_KEY="sk-your-api-key"
export OPENAI_MODEL="gpt-4"
```

#### YAML-Konfiguration
```yaml
options:
  - key: "ai-key"
    value: "sk-your-api-key"
  - key: "ai-model"
    value: "gpt-4"
  - key: "mcp"
    value: true
```

#### CLI-Parameter
```bash
runfromyaml --mcp --ai-key "sk-key" --ai-model "gpt-4" --debug
```

### Debugging und Monitoring

#### Debug-Modus aktivieren
```bash
runfromyaml --mcp --debug --ai-key "sk-key"
```

#### Log-Output
```
🚀 Starting MCP server 'runfromyaml-workflow-server' v1.0.0
🤖 AI Workflow Generator initialized with model: gpt-4
📡 Will use stdio transport (default for MCP)
🔧 AI-powered workflow generation enabled
```

### Fehlerbehandlung

#### AI-Service nicht verfügbar
```
⚠️  AI service unavailable, falling back to pattern matching
✅ Workflow generated using fallback method
```

#### Ungültiges AI-Response
```
❌ AI generated invalid YAML, using fallback
🔧 Workflow enhanced with validation fixes
```

#### API-Limits erreicht
```
⚠️  OpenAI API limit reached, using cached patterns
✅ Workflow generated successfully
```

### Best Practices

#### 1. **Prompt-Optimierung**
- Verwenden Sie spezifische, detaillierte Beschreibungen
- Geben Sie Kontext und Anforderungen an
- Erwähnen Sie gewünschte Technologien explizit

#### 2. **AI-Model-Auswahl**
- `gpt-4`: Beste Qualität, höhere Kosten
- `gpt-3.5-turbo`: Gute Balance zwischen Qualität und Kosten
- `gpt-4-turbo`: Optimiert für längere Workflows

#### 3. **Fallback-Strategien**
- Immer Fallback-Mechanismen aktiviert lassen
- Pattern-Matching als Backup konfigurieren
- Validierung nie überspringen

### Beispiele

#### Einfacher Web-Server
```
Prompt: "Create a simple Express.js web server with Docker"
```

#### Komplexe Microservice-Architektur
```
Prompt: "Setup a microservice architecture with Node.js API, PostgreSQL database, Redis cache, and Nginx reverse proxy using Docker Compose"
```

#### CI/CD Pipeline
```
Prompt: "Create a CI/CD pipeline that builds a React app, runs tests, builds Docker image, and deploys to staging"
```

### Troubleshooting

#### Problem: AI generiert ungültiges YAML
**Lösung:** Fallback-Mechanismus wird automatisch aktiviert

#### Problem: OpenAI API-Fehler
**Lösung:** Überprüfen Sie API-Key und Kontingent

#### Problem: Workflow funktioniert nicht wie erwartet
**Lösung:** Verwenden Sie `--debug` für detaillierte Logs

### Zukunftserweiterungen

1. **Multi-AI-Provider-Support** (Claude, Gemini)
2. **Workflow-Templates-Learning** 
3. **Kontext-basierte Verbesserungen**
4. **Automatische Optimierung basierend auf Feedback**
