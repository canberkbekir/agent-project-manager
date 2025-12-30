package state

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// WorkflowRepository defines database operations for Workflows
type WorkflowRepository interface {
	CreateWorkflow(workflow *Workflow) error
	GetWorkflow(name string) (*Workflow, error)
	ListWorkflows() ([]*Workflow, error)
	UpdateWorkflow(workflow *Workflow) error
	DeleteWorkflow(name string) error
}

// CreateWorkflow creates a new workflow
func (r *postgresRepository) CreateWorkflow(workflow *Workflow) error {
	now := time.Now()
	workflow.CreatedAt = now
	workflow.UpdatedAt = now

	schemaJSON, _ := json.Marshal(workflow.Schema)

	query := `INSERT INTO workflows (name, description, schema, version, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, workflow.Name, workflow.Description, string(schemaJSON),
		workflow.Version, workflow.CreatedAt, workflow.UpdatedAt)
	return err
}

// GetWorkflow retrieves a workflow by name
func (r *postgresRepository) GetWorkflow(name string) (*Workflow, error) {
	workflow := &Workflow{}
	var schemaJSON string

	query := `SELECT name, description, schema, version, created_at, updated_at
	          FROM workflows WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(
		&workflow.Name, &workflow.Description, &schemaJSON,
		&workflow.Version, &workflow.CreatedAt, &workflow.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("workflow not found: %s", name)
		}
		return nil, err
	}

	json.Unmarshal([]byte(schemaJSON), &workflow.Schema)
	return workflow, nil
}

// ListWorkflows lists all workflows
func (r *postgresRepository) ListWorkflows() ([]*Workflow, error) {
	rows, err := r.db.Query(`SELECT name, description, schema, version, created_at, updated_at FROM workflows ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workflows := []*Workflow{}
	for rows.Next() {
		workflow := &Workflow{}
		var schemaJSON string

		err := rows.Scan(&workflow.Name, &workflow.Description, &schemaJSON,
			&workflow.Version, &workflow.CreatedAt, &workflow.UpdatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(schemaJSON), &workflow.Schema)
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// UpdateWorkflow updates an existing workflow
func (r *postgresRepository) UpdateWorkflow(workflow *Workflow) error {
	workflow.UpdatedAt = time.Now()
	schemaJSON, _ := json.Marshal(workflow.Schema)

	query := `UPDATE workflows SET description = $1, schema = $2, version = $3, updated_at = $4 WHERE name = $5`
	_, err := r.db.Exec(query, workflow.Description, string(schemaJSON), workflow.Version, workflow.UpdatedAt, workflow.Name)
	return err
}

// DeleteWorkflow deletes a workflow by name
func (r *postgresRepository) DeleteWorkflow(name string) error {
	_, err := r.db.Exec("DELETE FROM workflows WHERE name = $1", name)
	return err
}
