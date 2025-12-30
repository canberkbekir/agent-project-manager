# Build Instructions

## PostgreSQL Setup

The application now uses PostgreSQL as the database backend. Make sure you have PostgreSQL installed and running.

### PostgreSQL Installation

#### Windows:
1. Download PostgreSQL from: https://www.postgresql.org/download/windows/
2. Install PostgreSQL (default port: 5432)
3. Create a database:
   ```sql
   CREATE DATABASE agent_project_manager;
   ```

#### Linux/macOS:
```bash
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib

# macOS (Homebrew)
brew install postgresql
brew services start postgresql

# Create database
createdb agent_project_manager
```

### Configuration

Update `configs/config.yaml` with your PostgreSQL connection string:

```yaml
state:
  connectionString: "postgres://user:password@localhost:5432/agent_project_manager?sslmode=disable"
```

Connection string format:
- `postgres://username:password@host:port/database?sslmode=disable`
- For local development, `sslmode=disable` is acceptable
- For production, use `sslmode=require` or `sslmode=verify-full`

### Build Commands

#### Windows PowerShell:
```powershell
# Build
go build -o agentd.exe ./cmd/agentd

# Or use make (if you have make installed)
make build-agentd
```

#### Windows CMD:
```cmd
go build -o agentd.exe ./cmd/agentd
```

#### Linux/macOS:
```bash
go build -o agentd ./cmd/agentd
# Or
make build-agentd
```

### Running the Application

1. Ensure PostgreSQL is running
2. Update `configs/config.yaml` with correct connection string
3. Run the application:
   ```bash
   ./agentd --config configs/config.yaml
   ```

The application will automatically run database migrations on startup.

### Troubleshooting

#### Error: "connection refused"
- Ensure PostgreSQL is running
- Check connection string in config.yaml
- Verify database exists

#### Error: "authentication failed"
- Check username and password in connection string
- Verify PostgreSQL user has access to the database

#### Error: "database does not exist"
- Create the database: `CREATE DATABASE agent_project_manager;`
