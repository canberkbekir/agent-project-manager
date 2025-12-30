# Grafana Setup Guide

This guide will help you set up Grafana to visualize metrics from your agent-project-manager application.

## Prerequisites

- Your application is running with `prometheusEnabled: true` in `configs/config.yaml`
- The metrics endpoint is available at `http://localhost:3333/metrics`

## Option 1: Using Docker (Recommended - Easiest)

### Step 1: Run Grafana with Docker

```bash
docker run -d \
  --name=grafana \
  -p 3000:3000 \
  -e "GF_SECURITY_ADMIN_PASSWORD=admin" \
  grafana/grafana:latest
```

### Step 2: Access Grafana

1. Open your browser and go to `http://localhost:3000`
2. Login with:
   - Username: `admin`
   - Password: `admin` (you'll be prompted to change it)

### Step 3: Add Prometheus Data Source

1. Click on the **Configuration** (gear icon) in the left sidebar
2. Select **Data Sources**
3. Click **Add data source**
4. Select **Prometheus**
5. Configure:
   - **URL**: `http://host.docker.internal:3333` (if using Docker) or `http://localhost:3333`
   - Click **Save & Test**

**Note**: If Grafana is running in Docker and your app is on the host, use `http://host.docker.internal:3333`. If both are on the same machine outside Docker, use `http://localhost:3333`.

## Option 2: Direct Installation (Windows)

### Step 1: Download Grafana

1. Go to https://grafana.com/grafana/download?platform=windows
2. Download the Windows installer
3. Run the installer and follow the setup wizard

### Step 2: Start Grafana

Grafana will start as a Windows service automatically, or you can start it manually from the Start menu.

### Step 3: Access Grafana

1. Open your browser and go to `http://localhost:3000`
2. Login with:
   - Username: `admin`
   - Password: `admin` (you'll be prompted to change it)

### Step 4: Add Prometheus Data Source

1. Click on the **Configuration** (gear icon) in the left sidebar
2. Select **Data Sources**
3. Click **Add data source**
4. Select **Prometheus**
5. Configure:
   - **URL**: `http://localhost:3333`
   - Click **Save & Test**

## Option 3: Using Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana

volumes:
  grafana-storage:
```

Then run:
```bash
docker-compose up -d
```

## Verify Metrics Endpoint

Before configuring Grafana, verify your metrics endpoint is working:

```bash
# Using curl (if available)
curl http://localhost:3333/metrics

# Or open in browser
# http://localhost:3333/metrics
```

You should see Prometheus-formatted metrics output.

## Create Your First Dashboard

1. In Grafana, click the **+** icon in the left sidebar
2. Select **Create Dashboard**
3. Click **Add visualization**
4. Select your Prometheus data source
5. In the query editor, try some example queries:

   - **All metrics**: `{__name__=~".+"}`
   - **Counter metrics**: `{__name__=~".*_total"}`
   - **Gauge metrics**: `{__name__=~".*_gauge"}`

6. Click **Apply** to see your visualization

## Example Queries

Since your metrics use the `agentd` namespace, queries will look like:

- `agentd_*` - All metrics from your application
- `rate(agentd_*_total[5m])` - Rate of counter metrics over 5 minutes
- `agentd_*_gauge` - All gauge metrics

## Troubleshooting

### Grafana can't connect to Prometheus endpoint

1. **Check if your app is running**: Make sure `agentd` is running and listening on port 3333
2. **Check the URL**: 
   - If Grafana is in Docker: use `http://host.docker.internal:3333`
   - If both are on host: use `http://localhost:3333`
3. **Check firewall**: Ensure port 3333 is accessible
4. **Test the endpoint**: Open `http://localhost:3333/metrics` in your browser

### No metrics showing up

1. **Verify metrics are enabled**: Check `configs/config.yaml` has `prometheusEnabled: true`
2. **Check application logs**: Look for "Prometheus metrics endpoint enabled" message
3. **Verify endpoint**: Visit `http://localhost:3333/metrics` directly
4. **Wait a bit**: Some metrics only appear after activity

## Next Steps

- Create custom dashboards for your specific metrics
- Set up alerts based on metric thresholds
- Explore Grafana's visualization options (graphs, tables, stat panels, etc.)

