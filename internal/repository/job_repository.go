package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"agent-project-manager/internal/state"
)

// IJobRepository defines database operations for Jobs
type IJobRepository interface {
	CreateJob(job *state.Job) error
	GetJob(id string) (*state.Job, error)
	ListJobs(limit int, cursor string, status string, workflow string) ([]*state.Job, string, error)
	UpdateJob(job *state.Job) error
	DeleteJob(id string) error
}

// JobRepository implements IJobRepository
type JobRepository struct {
	db *sql.DB
}

// CreateJob creates a new job in the database
func (r *JobRepository) CreateJob(job *state.Job) error {
	if job.ID == "" {
		job.ID = state.NewUUID()
	}
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now

	inputJSON, _ := json.Marshal(job.Input)
	metaJSON, _ := json.Marshal(job.Meta)

	query := `INSERT INTO jobs (id, workflow, status, input, meta, created_at, updated_at, started_at, completed_at, error)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(query, job.ID, job.Workflow, job.Status, string(inputJSON), string(metaJSON),
		job.CreatedAt, job.UpdatedAt, job.StartedAt, job.CompletedAt, job.Error)
	return err
}

// GetJob retrieves a job by ID from the database
func (r *JobRepository) GetJob(id string) (*state.Job, error) {
	job := &state.Job{}
	var inputJSON, metaJSON string
	var startedAt, completedAt sql.NullTime

	query := `SELECT id, workflow, status, input, meta, created_at, updated_at, started_at, completed_at, error
	          FROM jobs WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&job.ID, &job.Workflow, &job.Status, &inputJSON, &metaJSON,
		&job.CreatedAt, &job.UpdatedAt, &startedAt, &completedAt, &job.Error)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job not found: %s", id)
		}
		return nil, err
	}

	json.Unmarshal([]byte(inputJSON), &job.Input)
	json.Unmarshal([]byte(metaJSON), &job.Meta)
	if startedAt.Valid {
		job.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}

	return job, nil
}

// ListJobs lists jobs from the database with pagination and filtering
func (r *JobRepository) ListJobs(limit int, cursor string, status string, workflow string) ([]*state.Job, string, error) {
	query := `SELECT id, workflow, status, input, meta, created_at, updated_at, started_at, completed_at, error
	          FROM jobs WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, status)
		argPos++
	}
	if workflow != "" {
		query += fmt.Sprintf(" AND workflow = $%d", argPos)
		args = append(args, workflow)
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

	jobs := []*state.Job{}
	for rows.Next() {
		job := &state.Job{}
		var inputJSON, metaJSON string
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(&job.ID, &job.Workflow, &job.Status, &inputJSON, &metaJSON,
			&job.CreatedAt, &job.UpdatedAt, &startedAt, &completedAt, &job.Error)
		if err != nil {
			return nil, "", err
		}

		json.Unmarshal([]byte(inputJSON), &job.Input)
		json.Unmarshal([]byte(metaJSON), &job.Meta)
		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}

		jobs = append(jobs, job)
	}

	nextCursor := ""
	if len(jobs) > limit {
		nextCursor = jobs[limit].ID
		jobs = jobs[:limit]
	}

	return jobs, nextCursor, nil
}

// UpdateJob updates an existing job
func (r *JobRepository) UpdateJob(job *state.Job) error {
	job.UpdatedAt = time.Now()

	inputJSON, _ := json.Marshal(job.Input)
	metaJSON, _ := json.Marshal(job.Meta)

	query := `UPDATE jobs SET workflow = $1, status = $2, input = $3, meta = $4, updated_at = $5, 
	          started_at = $6, completed_at = $7, error = $8 WHERE id = $9`
	_, err := r.db.Exec(query, job.Workflow, job.Status, string(inputJSON), string(metaJSON),
		job.UpdatedAt, job.StartedAt, job.CompletedAt, job.Error, job.ID)
	return err
}

// DeleteJob deletes a job by ID
func (r *JobRepository) DeleteJob(id string) error {
	_, err := r.db.Exec("DELETE FROM jobs WHERE id = $1", id)
	return err
}

// NewJobRepository creates a new JobRepository
func NewJobRepository(db *sql.DB) IJobRepository {
	return &JobRepository{db: db}
}
