# Grafana Dockerfile
FROM grafana/grafana:latest

# Set environment variables
ENV GF_SECURITY_ADMIN_USER=admin
ENV GF_SECURITY_ADMIN_PASSWORD=admin
ENV GF_INSTALL_PLUGINS=""

# Copy custom configuration if needed
# COPY grafana.ini /etc/grafana/grafana.ini

# Copy provisioning files for datasources and dashboards
COPY provisioning/ /etc/grafana/provisioning/

# Expose Grafana port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

# Use default grafana entrypoint

