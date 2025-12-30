package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"agent-project-manager/internal/repository"
)

// handleListAgents handles GET /agents
// @Summary      List agents
// @Description  Get a list of all registered agents/workers
// @Tags         agents
// @Accept       json
// @Produce      json
// @Success      200  {object}  AgentListResponse
// @Router       /agents [get]
func handleListAgents(repo repository.IAgentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	// Get agents from repository
	stateAgents, err := repo.ListAgents()
	if err != nil {
		http.Error(w, "Failed to list agents: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert state models to API models
	agents := make([]Agent, len(stateAgents))
	for i, sa := range stateAgents {
		agents[i] = Agent{
			ID:        sa.ID,
			Name:      sa.Name,
			Status:    sa.Status,
			Metadata:  map[string]interface{}(sa.Metadata),
			LastSeen:  sa.LastSeen,
			CreatedAt: sa.CreatedAt,
		}
	}

	response := AgentListResponse{
		Agents: agents,
	}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// handleGetAgent handles GET /agents/{agentId}
// @Summary      Get agent details
// @Description  Get detailed information about a specific agent
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        agentId  path      string  true  "Agent ID"
// @Success      200      {object}  Agent
// @Failure      404      {string}  string  "Agent not found"
// @Router       /agents/{agentId} [get]
func handleGetAgent(repo repository.IAgentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agentID := chi.URLParam(r, "agentId")

		// Get agent from repository
		sa, err := repo.GetAgent(agentID)
		if err != nil {
			http.Error(w, "Agent not found", http.StatusNotFound)
			return
		}

		// Convert state model to API model
		agent := Agent{
			ID:        sa.ID,
			Name:      sa.Name,
			Status:    sa.Status,
			Metadata:  map[string]interface{}(sa.Metadata),
			LastSeen:  sa.LastSeen,
			CreatedAt: sa.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(agent)
	}
}

// handleGetAgentStatus handles GET /agents/{agentId}/status
// @Summary      Get agent status
// @Description  Get status information for a specific agent
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        agentId  path      string  true  "Agent ID"
// @Success      200      {object}  AgentStatus
// @Failure      404      {string}  string  "Agent not found"
// @Router       /agents/{agentId}/status [get]
func handleGetAgentStatus(repo repository.IAgentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agentID := chi.URLParam(r, "agentId")

		// TODO: Implement agent status retrieval
		_ = agentID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// handleDrainAgent handles POST /agents/{agentId}/drain
// @Summary      Drain agent
// @Description  Stop an agent from accepting new work (drain mode)
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        agentId  path      string  true  "Agent ID"
// @Success      200      {string}  string  "OK"
// @Failure      404      {string}  string  "Agent not found"
// @Router       /agents/{agentId}/drain [post]
func handleDrainAgent(repo repository.IAgentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agentID := chi.URLParam(r, "agentId")

		// TODO: Implement agent drain logic (stop taking new work)
		_ = agentID

		w.WriteHeader(http.StatusOK)
	}
}

// handleResumeAgent handles POST /agents/{agentId}/resume
// @Summary      Resume agent
// @Description  Resume an agent to accept new work
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        agentId  path      string  true  "Agent ID"
// @Success      200      {string}  string  "OK"
// @Failure      404      {string}  string  "Agent not found"
// @Router       /agents/{agentId}/resume [post]
func handleResumeAgent(repo repository.IAgentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agentID := chi.URLParam(r, "agentId")

		// TODO: Implement agent resume logic
		_ = agentID

		w.WriteHeader(http.StatusOK)
	}
}

