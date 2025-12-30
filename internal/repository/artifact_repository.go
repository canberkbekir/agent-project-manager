package repository

import (
	"database/sql"
	"fmt"
	"time"

	"agent-project-manager/internal/state"
)

// IArtifactRepository defines database operations for Artifacts
type IArtifactRepository interface {
	CreateArtifact(artifact *state.Artifact) error
	GetArtifact(id string) (*state.Artifact, error)
	ListArtifacts(jobID string, runID string, limit int, cursor string) ([]*state.Artifact, string, error)
	DeleteArtifact(id string) error
}

// ArtifactRepository implements IArtifactRepository
type ArtifactRepository struct {
	db *sql.DB
}

// CreateArtifact creates a new artifact
func (r *ArtifactRepository) CreateArtifact(artifact *state.Artifact) error {
	if artifact.ID == "" {
		artifact.ID = state.NewUUID()
	}
	artifact.CreatedAt = time.Now()

	query := `INSERT INTO artifacts (id, job_id, run_id, type, name, size, path, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, artifact.ID, artifact.JobID, artifact.RunID,
		artifact.Type, artifact.Name, artifact.Size, artifact.Path, artifact.CreatedAt)
	return err
}

// GetArtifact retrieves an artifact by ID
func (r *ArtifactRepository) GetArtifact(id string) (*state.Artifact, error) {
	artifact := &state.Artifact{}

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
func (r *ArtifactRepository) ListArtifacts(jobID string, runID string, limit int, cursor string) ([]*state.Artifact, string, error) {
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

	artifacts := []*state.Artifact{}
	for rows.Next() {
		artifact := &state.Artifact{}

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
func (r *ArtifactRepository) DeleteArtifact(id string) error {
	_, err := r.db.Exec("DELETE FROM artifacts WHERE id = $1", id)
	return err
}

// NewArtifactRepository creates a new ArtifactRepository
func NewArtifactRepository(db *sql.DB) IArtifactRepository {
	return &ArtifactRepository{db: db}
}
