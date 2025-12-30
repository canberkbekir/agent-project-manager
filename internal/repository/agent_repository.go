package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"agent-project-manager/internal/state"
)

// IAgentRepository defines database operations for Agents
type IAgentRepository interface {
	CreateAgent(agent *state.Agent) error
	GetAgent(id string) (*state.Agent, error)
	ListAgents() ([]*state.Agent, error)
	UpdateAgent(agent *state.Agent) error
	DeleteAgent(id string) error
}

// AgentRepository implements IAgentRepository
type AgentRepository struct {
	db *sql.DB
}

// CreateAgent creates a new agent
func (r *AgentRepository) CreateAgent(agent *state.Agent) error {
	if agent.ID == "" {
		agent.ID = state.NewUUID()
	}
	agent.CreatedAt = time.Now()

	metadataJSON, _ := json.Marshal(agent.Metadata)

	query := `INSERT INTO agents (id, name, status, metadata, last_seen, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, agent.ID, agent.Name, agent.Status, string(metadataJSON),
		agent.LastSeen, agent.CreatedAt)
	return err
}

// GetAgent retrieves an agent by ID
func (r *AgentRepository) GetAgent(id string) (*state.Agent, error) {
	agent := &state.Agent{}
	var metadataJSON string
	var lastSeen sql.NullTime

	query := `SELECT id, name, status, metadata, last_seen, created_at FROM agents WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&agent.ID, &agent.Name, &agent.Status, &metadataJSON, &lastSeen, &agent.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agent not found: %s", id)
		}
		return nil, err
	}

	json.Unmarshal([]byte(metadataJSON), &agent.Metadata)
	if lastSeen.Valid {
		agent.LastSeen = lastSeen.Time
	}

	return agent, nil
}

// ListAgents lists all agents
func (r *AgentRepository) ListAgents() ([]*state.Agent, error) {
	query := `SELECT id, name, status, metadata, last_seen, created_at FROM agents ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := []*state.Agent{}
	for rows.Next() {
		agent := &state.Agent{}
		var metadataJSON string
		var lastSeen sql.NullTime

		err := rows.Scan(
			&agent.ID, &agent.Name, &agent.Status, &metadataJSON, &lastSeen, &agent.CreatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(metadataJSON), &agent.Metadata)
		if lastSeen.Valid {
			agent.LastSeen = lastSeen.Time
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

// UpdateAgent updates an existing agent
func (r *AgentRepository) UpdateAgent(agent *state.Agent) error {
	metadataJSON, _ := json.Marshal(agent.Metadata)

	query := `UPDATE agents SET name = $1, status = $2, metadata = $3, last_seen = $4 WHERE id = $5`
	_, err := r.db.Exec(query, agent.Name, agent.Status, string(metadataJSON),
		agent.LastSeen, agent.ID)
	return err
}

// DeleteAgent deletes an agent by ID
func (r *AgentRepository) DeleteAgent(id string) error {
	_, err := r.db.Exec("DELETE FROM agents WHERE id = $1", id)
	return err
}

// NewAgentRepository creates a new AgentRepository
func NewAgentRepository(db *sql.DB) IAgentRepository {
	return &AgentRepository{db: db}
}
