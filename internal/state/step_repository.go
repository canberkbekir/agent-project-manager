package state

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// StepRepository defines database operations for Steps
type StepRepository interface {
	CreateStep(step *Step) error
	GetStep(id string) (*Step, error)
	ListSteps(jobID string) ([]*Step, error)
	UpdateStep(step *Step) error
	DeleteStep(id string) error
}

// CreateStep creates a new step
func (r *postgresRepository) CreateStep(step *Step) error {
	if step.ID == "" {
		step.ID = NewUUID()
	}
	now := time.Now()
	step.CreatedAt = now
	step.UpdatedAt = now

	inputJSON, _ := json.Marshal(step.Input)
	outputJSON, _ := json.Marshal(step.Output)

	query := `INSERT INTO steps (id, job_id, name, status, input, output, created_at, updated_at, started_at, completed_at, error)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(query, step.ID, step.JobID, step.Name, step.Status,
		string(inputJSON), string(outputJSON), step.CreatedAt, step.UpdatedAt,
		step.StartedAt, step.CompletedAt, step.Error)
	return err
}

// GetStep retrieves a step by ID
func (r *postgresRepository) GetStep(id string) (*Step, error) {
	step := &Step{}
	var inputJSON, outputJSON string
	var startedAt, completedAt sql.NullTime

	query := `SELECT id, job_id, name, status, input, output, created_at, updated_at, started_at, completed_at, error
	          FROM steps WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&step.ID, &step.JobID, &step.Name, &step.Status, &inputJSON, &outputJSON,
		&step.CreatedAt, &step.UpdatedAt, &startedAt, &completedAt, &step.Error)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("step not found: %s", id)
		}
		return nil, err
	}

	json.Unmarshal([]byte(inputJSON), &step.Input)
	json.Unmarshal([]byte(outputJSON), &step.Output)
	if startedAt.Valid {
		step.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		step.CompletedAt = &completedAt.Time
	}

	return step, nil
}

// ListSteps lists all steps for a job
func (r *postgresRepository) ListSteps(jobID string) ([]*Step, error) {
	rows, err := r.db.Query(`SELECT id, job_id, name, status, input, output, created_at, updated_at, started_at, completed_at, error
	                          FROM steps WHERE job_id = $1 ORDER BY created_at`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := []*Step{}
	for rows.Next() {
		step := &Step{}
		var inputJSON, outputJSON string
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(&step.ID, &step.JobID, &step.Name, &step.Status, &inputJSON, &outputJSON,
			&step.CreatedAt, &step.UpdatedAt, &startedAt, &completedAt, &step.Error)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(inputJSON), &step.Input)
		json.Unmarshal([]byte(outputJSON), &step.Output)
		if startedAt.Valid {
			step.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			step.CompletedAt = &completedAt.Time
		}

		steps = append(steps, step)
	}

	return steps, nil
}

// UpdateStep updates an existing step
func (r *postgresRepository) UpdateStep(step *Step) error {
	step.UpdatedAt = time.Now()
	inputJSON, _ := json.Marshal(step.Input)
	outputJSON, _ := json.Marshal(step.Output)

	query := `UPDATE steps SET job_id = $1, name = $2, status = $3, input = $4, output = $5, updated_at = $6, 
	          started_at = $7, completed_at = $8, error = $9 WHERE id = $10`
	_, err := r.db.Exec(query, step.JobID, step.Name, step.Status,
		string(inputJSON), string(outputJSON), step.UpdatedAt,
		step.StartedAt, step.CompletedAt, step.Error, step.ID)
	return err
}

// DeleteStep deletes a step by ID
func (r *postgresRepository) DeleteStep(id string) error {
	_, err := r.db.Exec("DELETE FROM steps WHERE id = $1", id)
	return err
}
