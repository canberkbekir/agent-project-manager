package state

import (
	"encoding/json"
	"fmt"
	"time"
)

// EventRepository defines database operations for Events
type EventRepository interface {
	CreateEvent(event *Event) error
	ListEvents(jobID string, stepID string, limit int) ([]*Event, error)
	DeleteEvents(jobID string) error
}

// CreateEvent creates a new event
func (r *postgresRepository) CreateEvent(event *Event) error {
	if event.ID == "" {
		event.ID = NewUUID()
	}
	event.CreatedAt = time.Now()

	dataJSON, _ := json.Marshal(event.Data)

	query := `INSERT INTO events (id, job_id, step_id, type, message, data, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, event.ID, event.JobID, event.StepID, event.Type,
		event.Message, string(dataJSON), event.CreatedAt)
	return err
}

// ListEvents lists events with optional filtering
func (r *postgresRepository) ListEvents(jobID string, stepID string, limit int) ([]*Event, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `SELECT id, job_id, step_id, type, message, data, created_at FROM events WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if jobID != "" {
		query += fmt.Sprintf(" AND job_id = $%d", argPos)
		args = append(args, jobID)
		argPos++
	}
	if stepID != "" {
		query += fmt.Sprintf(" AND step_id = $%d", argPos)
		args = append(args, stepID)
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argPos)
	args = append(args, limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}
	for rows.Next() {
		event := &Event{}
		var dataJSON string

		err := rows.Scan(&event.ID, &event.JobID, &event.StepID, &event.Type,
			&event.Message, &dataJSON, &event.CreatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(dataJSON), &event.Data)
		events = append(events, event)
	}

	return events, nil
}

// DeleteEvents deletes all events for a job
func (r *postgresRepository) DeleteEvents(jobID string) error {
	_, err := r.db.Exec("DELETE FROM events WHERE job_id = $1", jobID)
	return err
}
