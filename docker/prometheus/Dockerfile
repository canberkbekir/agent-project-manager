# Prometheus Dockerfile
FROM prom/prometheus:latest

# Set working directory
WORKDIR /etc/prometheus

# Copy Prometheus configuration
COPY prometheus.yml /etc/prometheus/prometheus.yml

# Create directory for data
RUN mkdir -p /prometheus && \
    chown -R nobody:nobody /etc/prometheus /prometheus

# Expose Prometheus port
EXPOSE 9090

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9090/-/healthy || exit 1

# Use default prometheus entrypoint
# The prometheus image will use the config file at /etc/prometheus/prometheus.yml

