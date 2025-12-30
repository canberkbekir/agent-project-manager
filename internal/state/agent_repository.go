package state

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// AgentRepository defines database operations for Agents
type AgentRepository interface {
	CreateAgent(agent *Agent) error
	GetAgent(id string) (*Agent, error)
	ListAgents() ([]*Agent, error)
	UpdateAgent(agent *Agent) error
	DeleteAgent(id string) error
}

// CreateAgent creates a new agent
func (r *postgresRepository) CreateAgent(agent *Agent) error {
	if agent.ID == "" {
		agent.ID = NewUUID()
	}
	now := time.Now()
	agent.CreatedAt = now
	agent.UpdatedAt = now
	agent.LastSeen = now

	metadataJSON, _ := json.Marshal(agent.Metadata)

	query := `INSERT INTO agents (id, name, status, metadata, last_seen, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, agent.ID, agent.Name, agent.Status, string(metadataJSON),
		agent.LastSeen, agent.CreatedAt, agent.UpdatedAt)
	return err
}

// GetAgent retrieves an agent by ID
func (r *postgresRepository) GetAgent(id string) (*Agent, error) {
	agent := &Agent{}
	var metadataJSON string

	query := `SELECT id, name, status, metadata, last_seen, created_at, updated_at FROM agents WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&agent.ID, &agent.Name, &agent.Status, &metadataJSON,
		&agent.LastSeen, &agent.CreatedAt, &agent.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agent not found: %s", id)
		}
		return nil, err
	}

	json.Unmarshal([]byte(metadataJSON), &agent.Metadata)
	return agent, nil
}

// ListAgents lists all agents
func (r *postgresRepository) ListAgents() ([]*Agent, error) {
	rows, err := r.db.Query(`SELECT id, name, status, metadata, last_seen, created_at, updated_at FROM agents ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := []*Agent{}
	for rows.Next() {
		agent := &Agent{}
		var metadataJSON string

		err := rows.Scan(&agent.ID, &agent.Name, &agent.Status, &metadataJSON,
			&agent.LastSeen, &agent.CreatedAt, &agent.UpdatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(metadataJSON), &agent.Metadata)
		agents = append(agents, agent)
	}

	return agents, nil
}

// UpdateAgent updates an existing agent
func (r *postgresRepository) UpdateAgent(agent *Agent) error {
	agent.UpdatedAt = time.Now()
	agent.LastSeen = time.Now()
	metadataJSON, _ := json.Marshal(agent.Metadata)

	query := `UPDATE agents SET name = $1, status = $2, metadata = $3, last_seen = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.Exec(query, agent.Name, agent.Status, string(metadataJSON),
		agent.LastSeen, agent.UpdatedAt, agent.ID)
	return err
}

// DeleteAgent deletes an agent by ID
func (r *postgresRepository) DeleteAgent(id string) error {
	_, err := r.db.Exec("DELETE FROM agents WHERE id = $1", id)
	return err
}
