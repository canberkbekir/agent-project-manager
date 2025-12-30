package state

import (
	"database/sql"
	"fmt"
	"time"
)

// ArtifactRepository defines database operations for Artifacts
type ArtifactRepository interface {
	CreateArtifact(artifact *Artifact) error
	GetArtifact(id string) (*Artifact, error)
	ListArtifacts(jobID string, runID string, limit int, cursor string) ([]*Artifact, string, error)
	DeleteArtifact(id string) error
}

// CreateArtifact creates a new artifact
func (r *postgresRepository) CreateArtifact(artifact *Artifact) error {
	if artifact.ID == "" {
		artifact.ID = NewUUID()
	}
	artifact.CreatedAt = time.Now()

	query := `INSERT INTO artifacts (id, job_id, run_id, type, name, size, path, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, artifact.ID, artifact.JobID, artifact.RunID,
		artifact.Type, artifact.Name, artifact.Size, artifact.Path, artifact.CreatedAt)
	return err
}

// GetArtifact retrieves an artifact by ID
func (r *postgresRepository) GetArtifact(id string) (*Artifact, error) {
	artifact := &Artifact{}

	query := `SELECT id, job_id, run_id, type, name, size, path, created_at FROM artifacts WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&artifact.ID, &artifact.JobID, &artifact.RunID, &artifact.Type,
		&artifact.Name, &artifact.Size, &artifact.Path, &artifact.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artifact not found: %s", id)
		}
		return nil, err
	}

	return artifact, nil
}

// ListArtifacts lists artifacts with pagination and optional filtering
func (r *postgresRepository) ListArtifacts(jobID string, runID string, limit int, cursor string) ([]*Artifact, string, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT id, job_id, run_id, type, name, size, path, created_at FROM artifacts WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if jobID != "" {
		query += fmt.Sprintf(" AND job_id = $%d", argPos)
		args = append(args, jobID)
		argPos++
	}
	if runID != "" {
		query += fmt.Sprintf(" AND run_id = $%d", argPos)
		args = append(args, runID)
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

	artifacts := []*Artifact{}
	for rows.Next() {
		artifact := &Artifact{}

		err := rows.Scan(&artifact.ID, &artifact.JobID, &artifact.RunID, &artifact.Type,
			&artifact.Name, &artifact.Size, &artifact.Path, &artifact.CreatedAt)
		if err != nil {
			return nil, "", err
		}

		artifacts = append(artifacts, artifact)
	}

	nextCursor := ""
	if len(artifacts) > limit {
		nextCursor = artifacts[limit].ID
		artifacts = artifacts[:limit]
	}

	return artifacts, nextCursor, nil
}

// DeleteArtifact deletes an artifact by ID
func (r *postgresRepository) DeleteArtifact(id string) error {
	_, err := r.db.Exec("DELETE FROM artifacts WHERE id = $1", id)
	return err
}
