package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"agent-project-manager/internal/repository"
	"agent-project-manager/internal/state"
)

// handleCreateRun handles POST /runs
// @Summary      Create a new run
// @Description  Start a new run for a job
// @Tags         runs
// @Accept       json
// @Produce      json
// @Param        run   body      CreateRunRequest  true  "Run creation request"
// @Success      201   {object}  CreateRunResponse
// @Failure      400   {string}  string  "Invalid request"
// @Router       /runs [post]
func handleCreateRun(repo repository.IRunRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Convert API model to state model
		run := &state.Run{
			JobID:  req.JobID,
			Status: string(RunStatusPending),
			Params: state.JSONMap(req.Params),
		}

		if err := repo.CreateRun(run); err != nil {
			http.Error(w, "Failed to create run: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := CreateRunResponse{
			ID: run.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// handleListRuns handles GET /runs
// @Summary      List runs
// @Description  Get a paginated list of runs with optional filtering
// @Tags         runs
// @Accept       json
// @Produce      json
// @Param        limit   query     int     false  "Maximum number of runs to return"
// @Param        cursor  query     string  false  "Cursor for pagination"
// @Param        status  query     string  false  "Filter by status"
// @Param        jobId   query     string  false  "Filter by job ID"
// @Success      200     {object}  RunListResponse
// @Router       /runs [get]
func handleListRuns(repo repository.IRunRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	cursor := r.URL.Query().Get("cursor")
	jobID := r.URL.Query().Get("jobId")

	limit := 50 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	// Get runs from repository
	stateRuns, nextCursor, err := repo.ListRuns(jobID, limit, cursor)
	if err != nil {
		http.Error(w, "Failed to list runs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert state models to API models
	runs := make([]Run, len(stateRuns))
	for i, sr := range stateRuns {
		status, _ := RunStatusFromString(sr.Status)
		runs[i] = Run{
			ID:          sr.ID,
			JobID:       sr.JobID,
			Status:      status,
			Params:      map[string]interface{}(sr.Params),
			CreatedAt:   sr.CreatedAt,
			UpdatedAt:   sr.UpdatedAt,
			StartedAt:   sr.StartedAt,
			CompletedAt: sr.CompletedAt,
			Error:       sr.Error,
		}
	}

	response := RunListResponse{
		Runs:    runs,
		Cursor:  nextCursor,
		HasMore: nextCursor != "",
	}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// handleGetRun handles GET /runs/{runId}
// @Summary      Get run details
// @Description  Get detailed information about a specific run
// @Tags         runs
// @Accept       json
// @Produce      json
// @Param        runId   path      string  true  "Run ID"
// @Success      200     {object}  Run
// @Failure      404     {string}  string  "Run not found"
// @Router       /runs/{runId} [get]
func handleGetRun(repo repository.IRunRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runID := chi.URLParam(r, "runId")

		// Get run from repository
		sr, err := repo.GetRun(runID)
		if err != nil {
			http.Error(w, "Run not found", http.StatusNotFound)
			return
		}

		// Convert state model to API model
		status, _ := RunStatusFromString(sr.Status)
		run := Run{
			ID:          sr.ID,
			JobID:       sr.JobID,
			Status:      status,
			Params:      map[string]interface{}(sr.Params),
			CreatedAt:   sr.CreatedAt,
			UpdatedAt:   sr.UpdatedAt,
			StartedAt:   sr.StartedAt,
			CompletedAt: sr.CompletedAt,
			Error:       sr.Error,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(run)
	}
}

// handleDeleteRun handles DELETE /runs/{runId}
// @Summary      Cancel a run
// @Description  Cancel a running or pending run
// @Tags         runs
// @Accept       json
// @Produce      json
// @Param        runId   path      string  true  "Run ID"
// @Success      202     {string}  string  "Accepted"
// @Failure      404     {string}  string  "Run not found"
// @Router       /runs/{runId} [delete]
func handleDeleteRun(repo repository.IRunRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runID := chi.URLParam(r, "runId")

		// Delete run from repository
		if err := repo.DeleteRun(runID); err != nil {
			http.Error(w, "Run not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

