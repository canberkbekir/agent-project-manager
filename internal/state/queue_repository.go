package state

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// QueueRepository defines database operations for Queue
type QueueRepository interface {
	CreateQueueItem(item *QueueItem) error
	GetQueueItem(id string) (*QueueItem, error)
	ListQueueItems(state string, limit int, cursor string) ([]*QueueItem, string, error)
	UpdateQueueItem(item *QueueItem) error
	DeleteQueueItem(id string) error
	GetQueueStats() (*QueueStats, error)
}

// CreateQueueItem creates a new queue item
func (r *postgresRepository) CreateQueueItem(item *QueueItem) error {
	if item.ID == "" {
		item.ID = NewUUID()
	}
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	dataJSON, _ := json.Marshal(item.Data)

	query := `INSERT INTO queue_items (id, job_id, state, data, created_at, updated_at, leased_at, completed_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, item.ID, item.JobID, item.State, string(dataJSON),
		item.CreatedAt, item.UpdatedAt, item.LeasedAt, item.CompletedAt)
	return err
}

// GetQueueItem retrieves a queue item by ID
func (r *postgresRepository) GetQueueItem(id string) (*QueueItem, error) {
	item := &QueueItem{}
	var dataJSON string
	var leasedAt, completedAt sql.NullTime

	query := `SELECT id, job_id, state, data, created_at, updated_at, leased_at, completed_at
	          FROM queue_items WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.JobID, &item.State, &dataJSON,
		&item.CreatedAt, &item.UpdatedAt, &leasedAt, &completedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("queue item not found: %s", id)
		}
		return nil, err
	}

	json.Unmarshal([]byte(dataJSON), &item.Data)
	if leasedAt.Valid {
		item.LeasedAt = &leasedAt.Time
	}
	if completedAt.Valid {
		item.CompletedAt = &completedAt.Time
	}

	return item, nil
}

// ListQueueItems lists queue items with pagination and optional state filtering
func (r *postgresRepository) ListQueueItems(state string, limit int, cursor string) ([]*QueueItem, string, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT id, job_id, state, data, created_at, updated_at, leased_at, completed_at
	          FROM queue_items WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if state != "" {
		query += fmt.Sprintf(" AND state = $%d", argPos)
		args = append(args, state)
		argPos++
	}
	if cursor != "" {
		query += fmt.Sprintf(" AND id > $%d", argPos)
		args = append(args, cursor)
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY created_at ASC LIMIT $%d", argPos)
	args = append(args, limit+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	items := []*QueueItem{}
	for rows.Next() {
		item := &QueueItem{}
		var dataJSON string
		var leasedAt, completedAt sql.NullTime

		err := rows.Scan(&item.ID, &item.JobID, &item.State, &dataJSON,
			&item.CreatedAt, &item.UpdatedAt, &leasedAt, &completedAt)
		if err != nil {
			return nil, "", err
		}

		json.Unmarshal([]byte(dataJSON), &item.Data)
		if leasedAt.Valid {
			item.LeasedAt = &leasedAt.Time
		}
		if completedAt.Valid {
			item.CompletedAt = &completedAt.Time
		}

		items = append(items, item)
	}

	nextCursor := ""
	if len(items) > limit {
		nextCursor = items[limit].ID
		items = items[:limit]
	}

	return items, nextCursor, nil
}

// UpdateQueueItem updates an existing queue item
func (r *postgresRepository) UpdateQueueItem(item *QueueItem) error {
	item.UpdatedAt = time.Now()
	dataJSON, _ := json.Marshal(item.Data)

	query := `UPDATE queue_items SET job_id = $1, state = $2, data = $3, updated_at = $4, 
	          leased_at = $5, completed_at = $6 WHERE id = $7`
	_, err := r.db.Exec(query, item.JobID, item.State, string(dataJSON),
		item.UpdatedAt, item.LeasedAt, item.CompletedAt, item.ID)
	return err
}

// DeleteQueueItem deletes a queue item by ID
func (r *postgresRepository) DeleteQueueItem(id string) error {
	_, err := r.db.Exec("DELETE FROM queue_items WHERE id = $1", id)
	return err
}

// GetQueueStats retrieves queue statistics
func (r *postgresRepository) GetQueueStats() (*QueueStats, error) {
	stats := &QueueStats{}

	query := `SELECT 
		SUM(CASE WHEN state = 'pending' THEN 1 ELSE 0 END) as pending,
		SUM(CASE WHEN state = 'leased' THEN 1 ELSE 0 END) as leased,
		SUM(CASE WHEN state = 'done' THEN 1 ELSE 0 END) as done,
		SUM(CASE WHEN state = 'dead' THEN 1 ELSE 0 END) as dead,
		COUNT(*) as total
		FROM queue_items`
	
	err := r.db.QueryRow(query).Scan(&stats.Pending, &stats.Leased, &stats.Done, &stats.Dead, &stats.Total)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
