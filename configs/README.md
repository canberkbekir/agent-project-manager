# Configuration Guide

This application follows the [12-Factor App](https://12factor.net/config) principles for configuration management.

## Configuration Hierarchy

1. **Environment Variables** (highest priority) - Used in production and Docker
2. **Config File** (defaults) - Used for development and local testing

## Environment Variables

All configuration can be set via environment variables. The naming convention follows the pattern: `SECTION_SUBSECTION_KEY`.

### API Configuration

```bash
API_ADDR=:3333                    # API server address
API_BASE_URL=http://localhost:3333  # Base URL (alternative to API_ADDR)
```

### Database Configuration

```bash
# Primary option
STATE_CONNECTION_STRING=postgres://user:password@host:5432/database?sslmode=disable

# Alternative (common convention)
DATABASE_URL=postgres://user:password@host:5432/database?sslmode=disable

# Legacy (backward compatibility)
STATE_CONNECTIONSTRING=postgres://user:password@host:5432/database?sslmode=disable
```

### Queue Configuration

```bash
QUEUE_WORKERS=2  # Number of worker threads
```

### Artifacts Configuration

```bash
ARTIFACTS_WORK_DIR=/app/data/workdir  # Working directory for artifacts
```

### Logger Configuration

```bash
LOGGER_LEVEL=info          # debug, info, warn, error, fatal
LOGGER_FORMAT=json         # text, json
LOGGER_OUTPUT=stdout       # stdout, stderr, or file path
LOGGER_REPORT_CALLER=false # true/false
```

### LLM Configuration

```bash
LLM_PROVIDER=ollama                    # ollama or openai
LLM_OPENAI_API_KEY=your-api-key       # OpenAI API key
LLM_OPENAI_MODEL=gpt-4o-mini          # OpenAI model
LLM_OLLAMA_BASE_URL=http://localhost:11434  # Ollama base URL
LLM_OLLAMA_MODEL=qwen2.5-coder:7b     # Ollama model
```

### Observability - Tracing

```bash
OBS_TRACING_ENABLED=true              # true/false
OBS_TRACING_ENDPOINT=localhost:4318    # OTLP endpoint
```

### Observability - Metrics

```bash
OBS_METRICS_ENABLED=true                    # true/false
OBS_METRICS_ENDPOINT=localhost:4318          # OTLP endpoint
OBS_METRICS_PROMETHEUS_ENABLED=true          # Enable Prometheus endpoint
OBS_METRICS_PROMETHEUS_PATH=/metrics        # Prometheus metrics path
```

### Auth Configuration

```bash
AUTH_TOKEN=your-auth-token  # Authentication token
```

## Config File

Config files are optional and used primarily for development. They provide default values that can be overridden by environment variables.

### Default Config File

- **Path**: `configs/config.yaml`
- **Usage**: Local development
- **Override**: Set `CONFIG_PATH` environment variable to use a different file

### Docker Config File

- **Path**: `configs/config.docker.yaml`
- **Usage**: Example Docker configuration (not required when using environment variables)
- **Note**: In Docker Compose, environment variables are preferred over config files

## Docker Usage

### Using Environment Variables (Recommended)

```yaml
services:
  api:
    environment:
      - STATE_CONNECTION_STRING=postgres://user:password@database:5432/db?sslmode=disable
      - API_ADDR=:3333
      - LOGGER_LEVEL=info
```

### Using Config File (Optional)

```yaml
services:
  api:
    environment:
      - CONFIG_PATH=/app/configs/config.docker.yaml
    volumes:
      - ./configs:/app/configs:ro
```

## Best Practices

1. **Production**: Use environment variables only (no config files)
2. **Development**: Use config files for convenience
3. **Docker**: Prefer environment variables in docker-compose.yml
4. **Secrets**: Never commit secrets to config files; use environment variables or secrets management
5. **Validation**: All required config values are validated on startup

## Required Configuration

The following must be set (via environment variable or config file):

- `API_ADDR` or `API_BASE_URL` - API server address
- `STATE_CONNECTION_STRING` or `DATABASE_URL` - Database connection string

## Examples

### Local Development

```bash
# Use default config file
./agentd

# Or specify custom config
CONFIG_PATH=configs/config.local.yaml ./agentd
```

### Docker Compose

```yaml
environment:
  - STATE_CONNECTION_STRING=postgres://user:password@database:5432/db?sslmode=disable
  - API_ADDR=:3333
```

### Production (Kubernetes)

```yaml
env:
  - name: STATE_CONNECTION_STRING
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: connection-string
  - name: API_ADDR
    value: ":3333"
```

## Troubleshooting

### Config file not found

If you see "config file does not exist", it means:
- The config file is optional - you can use environment variables instead
- Or set `CONFIG_PATH` to point to an existing file

### Environment variable not working

- Check the variable name (case-sensitive, UPPER_SNAKE_CASE)
- Ensure no typos in the variable name
- Environment variables override config file values

### Database connection issues

- Verify `STATE_CONNECTION_STRING` or `DATABASE_URL` is set correctly
- Check that the database hostname is correct (use service name in Docker)
- Ensure database is accessible and credentials are correct

