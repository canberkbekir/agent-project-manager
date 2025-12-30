# Docker Setup Guide

This directory contains Docker containers for the agent-project-manager project.

## Services

1. **API (agentd)**: Main HTTP API service
2. **Database (PostgreSQL)**: Database service
3. **Grafana**: Metrics visualization
4. **Prometheus**: Metrics collection and storage

## Quick Start

### Start all services

```bash
docker-compose up -d
```

### Stop services

```bash
docker-compose down
```

### Stop services and remove volumes

```bash
docker-compose down -v
```

## Service Access Information

- **API**: http://localhost:3333
  - Swagger UI: http://localhost:3333/swagger/
  - Health Check: http://localhost:3333/v1/healthz
  - Metrics: http://localhost:3333/metrics

- **PostgreSQL**: localhost:5432
  - Database: `agent_project_manager`
  - User: `user`
  - Password: `password`

- **Grafana**: http://localhost:3000
  - Username: `admin`
  - Password: `admin`

- **Prometheus**: http://localhost:9090

## Configuration

### API Service

The API service uses the `configs/config.yaml` file. This file is mounted as a volume with Docker Compose.

To update the database connection string, edit the `configs/config.yaml` file:

```yaml
state:
  connectionString: "postgres://user:password@database:5432/agent_project_manager?sslmode=disable"
```

**Note**: Since the service name is `database` in Docker Compose, use `database` instead of `localhost` in the connection string.

### Database

Default settings for PostgreSQL:
- Database: `agent_project_manager`
- User: `user`
- Password: `password`

To change these settings, update the environment variables in the `docker-compose.yml` file.

### Grafana

Grafana automatically configures Prometheus as a data source. You will be prompted to change the password on first login.

## Viewing Logs

### All services logs

```bash
docker-compose logs -f
```

### Specific service logs

```bash
docker-compose logs -f api
docker-compose logs -f database
docker-compose logs -f grafana
docker-compose logs -f prometheus
```

## Development

### Rebuild only the API service

```bash
docker-compose build api
docker-compose up -d api
```

### Rebuild all services

```bash
docker-compose build --no-cache
docker-compose up -d
```

## Data Persistence

Docker Compose creates the following volumes:
- `postgres_data`: PostgreSQL database data
- `api_data`: Working directory and artifacts for the API service
- `prometheus_data`: Prometheus metrics data
- `grafana_data`: Grafana configuration and dashboards

Data persists as long as these volumes are not deleted.

## Troubleshooting

### Services not starting

1. Make sure ports are available:
   ```bash
   lsof -i :3333  # API
   lsof -i :5432  # PostgreSQL
   lsof -i :3000  # Grafana
   lsof -i :9090  # Prometheus
   ```

2. Check logs:
   ```bash
   docker-compose logs
   ```

### Database connection error

For the API service to connect to the database:
1. Make sure the database service is healthy: `docker-compose ps`
2. Ensure the connection string is correct (host: `database`, port: `5432`)
3. Wait for the database service to fully start (until health check passes)

### Metrics not showing in Grafana

1. Make sure Prometheus is collecting metrics from the API service:
   - In Prometheus UI: http://localhost:9090/targets
   - Check that the `agentd-api` target is in `UP` state

2. Make sure the Prometheus data source is correctly configured in Grafana:
   - Configuration > Data Sources > Prometheus
   - URL should be: `http://prometheus:9090`

## Production Usage

Before using in production:

1. **Security**: Change all default passwords
2. **SSL/TLS**: Add a reverse proxy (nginx/traefik) for HTTPS
3. **Backup**: Create a regular backup strategy for the database
4. **Monitoring**: Add additional monitoring and alerting tools
5. **Resource Limits**: Set resource limits in docker-compose.yml
6. **Secrets Management**: Use Docker secrets or external secrets management for sensitive information
