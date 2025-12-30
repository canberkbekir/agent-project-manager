# PostgreSQL Database Dockerfile
FROM postgres:16-alpine

# Set environment variables
ENV POSTGRES_DB=agent_project_manager
ENV POSTGRES_USER=user
ENV POSTGRES_PASSWORD=password
ENV PGDATA=/var/lib/postgresql/data/pgdata

# Create directory for custom initialization scripts
RUN mkdir -p /docker-entrypoint-initdb.d

# Copy initialization scripts if needed
# COPY init.sql /docker-entrypoint-initdb.d/

# Expose PostgreSQL port
EXPOSE 5432

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 \
    CMD pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} || exit 1

# Use default postgres entrypoint
# The postgres image will automatically run any scripts in /docker-entrypoint-initdb.d

