package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"agent-project-manager/internal/state"
)

// IRunRepository defines database operations for Runs
type IRunRepository interface {
	CreateRun(run *state.Run) error
	GetRun(id string) (*state.Run, error)
	ListRuns(jobID string, limit int, cursor string) ([]*state.Run, string, error)
	UpdateRun(run *state.Run) error
	DeleteRun(id string) error
}

// RunRepository implements IRunRepository
type RunRepository struct {
	db *sql.DB
}

// CreateRun creates a new run
func (r *RunRepository) CreateRun(run *state.Run) error {
	if run.ID == "" {
		run.ID = state.NewUUID()
	}
	now := time.Now()
	run.CreatedAt = now
	run.UpdatedAt = now

	paramsJSON, _ := json.Marshal(run.Params)

	query := `INSERT INTO runs (id, job_id, status, params, created_at, updated_at, started_at, completed_at, error)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query, run.ID, run.JobID, run.Status, string(paramsJSON),
		run.CreatedAt, run.UpdatedAt, run.StartedAt, run.CompletedAt, run.Error)
	return err
}

// GetRun retrieves a run by ID
func (r *RunRepository) GetRun(id string) (*state.Run, error) {
	run := &state.Run{}
	var paramsJSON string
	var startedAt, completedAt sql.NullTime

	query := `SELECT id, job_id, status, params, created_at, updated_at, started_at, completed_at, error
	          FROM runs WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&run.ID, &run.JobID, &run.Status, &paramsJSON,
		&run.CreatedAt, &run.UpdatedAt, &startedAt, &completedAt, &run.Error)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("run not found: %s", id)
		}
		return nil, err
	}

	json.Unmarshal([]byte(paramsJSON), &run.Params)
	if startedAt.Valid {
		run.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		run.CompletedAt = &completedAt.Time
	}

	return run, nil
}

// ListRuns lists runs with pagination and optional job filtering
func (r *RunRepository) ListRuns(jobID string, limit int, cursor string) ([]*state.Run, string, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT id, job_id, status, params, created_at, updated_at, started_at, completed_at, error
	          FROM runs WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if jobID != "" {
		query += fmt.Sprintf(" AND job_id = $%d", argPos)
		args = append(args, jobID)
		argPos++
	}
	if cursor != "" {
		query += fmt.Sprintf(" AND id > $%d", argPos)
		args = append(args, cursor)
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argPos)
	args = append(args, limit+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	runs := []*state.Run{}
	for rows.Next() {
		run := &state.Run{}
		var paramsJSON string
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(
			&run.ID, &run.JobID, &run.Status, &paramsJSON,
			&run.CreatedAt, &run.UpdatedAt, &startedAt, &completedAt, &run.Error)
		if err != nil {
			return nil, "", err
		}

		json.Unmarshal([]byte(paramsJSON), &run.Params)
		if startedAt.Valid {
			run.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			run.CompletedAt = &completedAt.Time
		}

		runs = append(runs, run)
	}

	nextCursor := ""
	if len(runs) > limit {
		nextCursor = runs[limit].ID
		runs = runs[:limit]
	}

	return runs, nextCursor, nil
}

// UpdateRun updates an existing run
func (r *RunRepository) UpdateRun(run *state.Run) error {
	run.UpdatedAt = time.Now()

	paramsJSON, _ := json.Marshal(run.Params)

	query := `UPDATE runs SET job_id = $1, status = $2, params = $3, updated_at = $4, 
	          started_at = $5, completed_at = $6, error = $7 WHERE id = $8`
	_, err := r.db.Exec(query, run.JobID, run.Status, string(paramsJSON),
		run.UpdatedAt, run.StartedAt, run.CompletedAt, run.Error, run.ID)
	return err
}

// DeleteRun deletes a run by ID
func (r *RunRepository) DeleteRun(id string) error {
	_, err := r.db.Exec("DELETE FROM runs WHERE id = $1", id)
	return err
}

// NewRunRepository creates a new RunRepository
func NewRunRepository(db *sql.DB) IRunRepository {
	return &RunRepository{db: db}
}
