package state

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Database models - these represent the database schema

// Job represents a job in the database
type Job struct {
	ID          string    `db:"id"`
	Workflow    string    `db:"workflow"`
	Status      string    `db:"status"`
	Input       JSONMap   `db:"input"`
	Meta        JSONMap   `db:"meta"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	StartedAt   *time.Time `db:"started_at"`
	CompletedAt *time.Time `db:"completed_at"`
	Error       string    `db:"error"`
}

// Run represents a run in the database
type Run struct {
	ID          string    `db:"id"`
	JobID       string    `db:"job_id"`
	Status      string    `db:"status"`
	Params      JSONMap   `db:"params"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	StartedAt   *time.Time `db:"started_at"`
	CompletedAt *time.Time `db:"completed_at"`
	Error       string    `db:"error"`
}

// Workflow represents a workflow definition in the database
type Workflow struct {
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Schema      JSONMap   `db:"schema"`
	Version     string    `db:"version"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Step represents a workflow step in the database
type Step struct {
	ID          string    `db:"id"`
	JobID       string    `db:"job_id"`
	Name        string    `db:"name"`
	Status      string    `db:"status"`
	Input       JSONMap   `db:"input"`
	Output      JSONMap   `db:"output"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	StartedAt   *time.Time `db:"started_at"`
	CompletedAt *time.Time `db:"completed_at"`
	Error       string    `db:"error"`
}

// Event represents an event in the database
type Event struct {
	ID        string    `db:"id"`
	JobID     string    `db:"job_id"`
	StepID    string    `db:"step_id"`
	Type      string    `db:"type"`
	Message   string    `db:"message"`
	Data      JSONMap   `db:"data"`
	CreatedAt time.Time `db:"created_at"`
}

// Artifact represents an artifact in the database
type Artifact struct {
	ID        string    `db:"id"`
	JobID     string    `db:"job_id"`
	RunID     string    `db:"run_id"`
	Type      string    `db:"type"`
	Name      string    `db:"name"`
	Size      int64     `db:"size"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
}

// Agent represents an agent in the database
type Agent struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Status    string    `db:"status"`
	Metadata  JSONMap   `db:"metadata"`
	LastSeen  time.Time `db:"last_seen"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// QueueItem represents a queue item in the database
type QueueItem struct {
	ID          string    `db:"id"`
	JobID       string    `db:"job_id"`
	State       string    `db:"state"`
	Data        JSONMap   `db:"data"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	LeasedAt    *time.Time `db:"leased_at"`
	CompletedAt *time.Time `db:"completed_at"`
}

// JSONMap is a type alias for map[string]interface{} that implements
// sql/driver.Valuer and sql.Scanner for JSON storage in SQLite
type JSONMap map[string]interface{}

// Value implements driver.Valuer
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

// NewUUID generates a new UUID string
func NewUUID() string {
	return uuid.New().String()
}

