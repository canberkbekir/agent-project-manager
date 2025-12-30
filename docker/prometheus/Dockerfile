# Prometheus Dockerfile
FROM prom/prometheus:latest

# Set working directory
WORKDIR /etc/prometheus

# Copy Prometheus configuration
# Build context is project root, so use relative path
COPY docker/prometheus/prometheus.yml /etc/prometheus/prometheus.yml

# Create directory for data
# Prometheus image already handles permissions correctly
RUN mkdir -p /prometheus

# Expose Prometheus port
EXPOSE 9090

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9090/-/healthy || exit 1

# Use default prometheus entrypoint
# The prometheus image will use the config file at /etc/prometheus/prometheus.yml

