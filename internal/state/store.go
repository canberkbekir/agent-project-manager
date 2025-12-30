package state

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

// Repository is the interface for database operations
// It embeds all model-specific repositories
type Repository interface {
	JobRepository
	RunRepository
	WorkflowRepository
	StepRepository
	EventRepository
	ArtifactRepository
	AgentRepository
	QueueRepository
	
	// Migration
	Migrate(migrationsPath string) error

	// GetDB returns the underlying database connection
	GetDB() *sql.DB

	// Close closes the database connection
	Close() error
}

// Store is an alias for Repository for backward compatibility
type Store = Repository

// QueueStats represents queue statistics
type QueueStats struct {
	Pending int
	Leased  int
	Done    int
	Dead    int
	Total   int
}

// postgresRepository implements Repository using PostgreSQL
type postgresRepository struct {
	db *sql.DB
}

// Compile-time interface implementation checks
var (
	_ Repository         = (*postgresRepository)(nil)
	_ JobRepository      = (*postgresRepository)(nil)
	_ RunRepository      = (*postgresRepository)(nil)
	_ WorkflowRepository = (*postgresRepository)(nil)
	_ StepRepository     = (*postgresRepository)(nil)
	_ EventRepository    = (*postgresRepository)(nil)
	_ ArtifactRepository = (*postgresRepository)(nil)
	_ AgentRepository    = (*postgresRepository)(nil)
	_ QueueRepository    = (*postgresRepository)(nil)
)

// NewRepository creates a new PostgreSQL repository
func NewRepository(connectionString string) (Repository, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	repo := &postgresRepository{db: db}
	return repo, nil
}

// NewStore is an alias for NewRepository for backward compatibility
func NewStore(connectionString string) (Store, error) {
	return NewRepository(connectionString)
}

// Migrate runs database migrations
func (r *postgresRepository) Migrate(migrationsPath string) error {
	// Read migration file
	migrationSQL, err := os.ReadFile(filepath.Join(migrationsPath, "0001_init.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := r.db.Exec(string(migrationSQL)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

// GetDB returns the underlying database connection
func (r *postgresRepository) GetDB() *sql.DB {
	return r.db
}

// Close closes the database connection
func (r *postgresRepository) Close() error {
	return r.db.Close()
}
