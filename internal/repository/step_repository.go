package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"agent-project-manager/internal/state"
)

// IStepRepository defines database operations for Steps
type IStepRepository interface {
	CreateStep(step *state.Step) error
	GetStep(id string) (*state.Step, error)
	ListSteps(jobID string) ([]*state.Step, error)
	UpdateStep(step *state.Step) error
	DeleteStep(id string) error
}

// StepRepository implements IStepRepository
type StepRepository struct {
	db *sql.DB
}

// CreateStep creates a new step
func (r *StepRepository) CreateStep(step *state.Step) error {
	if step.ID == "" {
		step.ID = state.NewUUID()
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
func (r *StepRepository) GetStep(id string) (*state.Step, error) {
	step := &state.Step{}
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

// ListSteps lists steps for a job
func (r *StepRepository) ListSteps(jobID string) ([]*state.Step, error) {
	query := `SELECT id, job_id, name, status, input, output, created_at, updated_at, started_at, completed_at, error
	          FROM steps WHERE job_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := []*state.Step{}
	for rows.Next() {
		step := &state.Step{}
		var inputJSON, outputJSON string
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(
			&step.ID, &step.JobID, &step.Name, &step.Status, &inputJSON, &outputJSON,
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
func (r *StepRepository) UpdateStep(step *state.Step) error {
	step.UpdatedAt = time.Now()

	inputJSON, _ := json.Marshal(step.Input)
	outputJSON, _ := json.Marshal(step.Output)

	query := `UPDATE steps SET name = $1, status = $2, input = $3, output = $4, updated_at = $5, 
	          started_at = $6, completed_at = $7, error = $8 WHERE id = $9`
	_, err := r.db.Exec(query, step.Name, step.Status, string(inputJSON), string(outputJSON),
		step.UpdatedAt, step.StartedAt, step.CompletedAt, step.Error, step.ID)
	return err
}

// DeleteStep deletes a step by ID
func (r *StepRepository) DeleteStep(id string) error {
	_, err := r.db.Exec("DELETE FROM steps WHERE id = $1", id)
	return err
}

// NewStepRepository creates a new StepRepository
func NewStepRepository(db *sql.DB) IStepRepository {
	return &StepRepository{db: db}
}
