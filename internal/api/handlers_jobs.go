package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"agent-project-manager/internal/repository"
	"agent-project-manager/internal/state"
)

// handleCreateJob handles POST /jobs
// @Summary      Create a new job
// @Description  Submit a new job to be processed
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        job  body      CreateJobRequest  true  "Job creation request"
// @Success      201  {object}  CreateJobResponse
// @Failure      400  {string}  string  "Invalid request body"
// @Router       /jobs [post]
func handleCreateJob(repo repository.IJobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	
	// Convert API model to state model
	job := &state.Job{
		Workflow: req.Workflow,
		Status:   string(JobStatusQueued),
		Input:    state.JSONMap(req.Input),
		Meta:     state.JSONMap(req.Meta),
	}

	if err := repo.CreateJob(job); err != nil {
		http.Error(w, "Failed to create job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateJobResponse{
		ID: job.ID,
	}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// handleListJobs handles GET /jobs
// @Summary      List jobs
// @Description  Get a paginated list of jobs with optional filtering
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        limit     query     int     false  "Maximum number of jobs to return"
// @Param        cursor    query     string  false  "Cursor for pagination"
// @Param        status    query     string  false  "Filter by status (queued|running|succeeded|failed)"
// @Param        workflow  query     string  false  "Filter by workflow name"
// @Success      200       {object}  JobListResponse
// @Router       /jobs [get]
func handleListJobs(repo repository.IJobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		cursor := r.URL.Query().Get("cursor")
		status := r.URL.Query().Get("status")
		workflow := r.URL.Query().Get("workflow")

		limit := 50 // default
		if limitStr != "" {
			if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		// Get jobs from repository
		stateJobs, nextCursor, err := repo.ListJobs(limit, cursor, status, workflow)
		if err != nil {
			http.Error(w, "Failed to list jobs: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert state models to API models
		jobs := make([]Job, len(stateJobs))
		for i, sj := range stateJobs {
			status, _ := JobStatusFromString(sj.Status)
			jobs[i] = Job{
				ID:          sj.ID,
				Workflow:    sj.Workflow,
				Status:      status,
				Input:       map[string]interface{}(sj.Input),
				Meta:        map[string]interface{}(sj.Meta),
				CreatedAt:   sj.CreatedAt,
				UpdatedAt:   sj.UpdatedAt,
				StartedAt:   sj.StartedAt,
				CompletedAt: sj.CompletedAt,
				Error:       sj.Error,
			}
		}

		response := JobListResponse{
			Jobs:    jobs,
			Cursor:  nextCursor,
			HasMore: nextCursor != "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// handleGetJob handles GET /jobs/{jobId}
// @Summary      Get job details
// @Description  Get detailed information about a specific job
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobId   path      string  true  "Job ID"
// @Success      200     {object}  Job
// @Failure      404     {string}  string  "Job not found"
// @Router       /jobs/{jobId} [get]
func handleGetJob(repo repository.IJobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")

		// Get job from repository
		sj, err := repo.GetJob(jobID)
		if err != nil {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}

		// Convert state model to API model
		status, _ := JobStatusFromString(sj.Status)
		job := Job{
			ID:          sj.ID,
			Workflow:    sj.Workflow,
			Status:      status,
			Input:       map[string]interface{}(sj.Input),
			Meta:        map[string]interface{}(sj.Meta),
			CreatedAt:   sj.CreatedAt,
			UpdatedAt:   sj.UpdatedAt,
			StartedAt:   sj.StartedAt,
			CompletedAt: sj.CompletedAt,
			Error:       sj.Error,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(job)
	}
}

// handleDeleteJob handles DELETE /jobs/{jobId}
// @Summary      Cancel a job
// @Description  Cancel a running or queued job
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobId   path      string  true  "Job ID"
// @Success      202     {string}  string  "Accepted"
// @Failure      404     {string}  string  "Job not found"
// @Router       /jobs/{jobId} [delete]
func handleDeleteJob(repo repository.IJobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")

	// Delete job from repository
	if err := repo.DeleteJob(jobID); err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

		w.WriteHeader(http.StatusAccepted)
	}
}

// handleRetryJob handles POST /jobs/{jobId}/retry
// @Summary      Retry a job
// @Description  Retry a failed or cancelled job
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobId   path      string  true  "Job ID"
// @Success      200     {object}  CreateJobResponse
// @Failure      404     {string}  string  "Job not found"
// @Router       /jobs/{jobId}/retry [post]
func handleRetryJob(repo repository.IJobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")

		// TODO: Implement job retry logic
		_ = jobID
		response := CreateJobResponse{
			ID: generateID(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// handleJobEvents handles GET /jobs/{jobId}/events
// @Summary      Stream job events
// @Description  Stream job status updates via Server-Sent Events (SSE)
// @Tags         jobs
// @Accept       json
// @Produce      text/event-stream
// @Param        jobId   path      string  true  "Job ID"
// @Success      200     {string}  string  "Event stream"
// @Router       /jobs/{jobId}/events [get]
func handleJobEvents(repo repository.IJobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")

		// TODO: Implement SSE (Server-Sent Events) streaming for status updates
		_ = jobID

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
	}
}

// handleJobLogs handles GET /jobs/{jobId}/logs
// @Summary      Get job logs
// @Description  Get aggregated logs for a job
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobId   path      string  true   "Job ID"
// @Param        tail    query     int     false  "Number of log lines to return"
// @Param        since   query     string  false  "RFC3339 timestamp to filter logs from"
// @Success      200     {object}  JobLogsResponse
// @Failure      404     {string}  string  "Job not found"
// @Router       /jobs/{jobId}/logs [get]
func handleJobLogs(repo repository.IJobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")
		// TODO: Parse query parameters: tail, since
		// tail := r.URL.Query().Get("tail")
		// since := r.URL.Query().Get("since")

		// TODO: Implement aggregated logs retrieval
		_ = jobID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// handleJobSteps handles GET /jobs/{jobId}/steps
// @Summary      List job steps
// @Description  Get a list of workflow steps for a job
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobId   path      string  true  "Job ID"
// @Success      200     {array}   JobStep
// @Failure      404     {string}  string  "Job not found"
// @Router       /jobs/{jobId}/steps [get]
func handleJobSteps(repo repository.IStepRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")

		// Get steps from repository
		stateSteps, err := repo.ListSteps(jobID)
		if err != nil {
			http.Error(w, "Failed to list steps: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert state models to API models
		steps := make([]JobStep, len(stateSteps))
		for i, ss := range stateSteps {
			status, _ := StepStatusFromString(ss.Status)
			steps[i] = JobStep{
				ID:          ss.ID,
				JobID:       ss.JobID,
				Name:        ss.Name,
				Status:      status,
				Input:       map[string]interface{}(ss.Input),
				Output:      map[string]interface{}(ss.Output),
				CreatedAt:   ss.CreatedAt,
				UpdatedAt:   ss.UpdatedAt,
				StartedAt:   ss.StartedAt,
				CompletedAt: ss.CompletedAt,
				Error:       ss.Error,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(steps)
	}
}

// handleGetStep handles GET /jobs/{jobId}/steps/{stepId}
// @Summary      Get step details
// @Description  Get detailed information about a workflow step
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobId   path      string  true  "Job ID"
// @Param        stepId  path      string  true  "Step ID"
// @Success      200     {object}  JobStep
// @Failure      404     {string}  string  "Step not found"
// @Router       /jobs/{jobId}/steps/{stepId} [get]
func handleGetStep(repo repository.IStepRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stepID := chi.URLParam(r, "stepId")

		// Get step from repository
		ss, err := repo.GetStep(stepID)
		if err != nil {
			http.Error(w, "Step not found", http.StatusNotFound)
			return
		}

		// Convert state model to API model
		status, _ := StepStatusFromString(ss.Status)
		step := JobStep{
			ID:          ss.ID,
			JobID:       ss.JobID,
			Name:        ss.Name,
			Status:      status,
			Input:       map[string]interface{}(ss.Input),
			Output:      map[string]interface{}(ss.Output),
			CreatedAt:   ss.CreatedAt,
			UpdatedAt:   ss.UpdatedAt,
			StartedAt:   ss.StartedAt,
			CompletedAt: ss.CompletedAt,
			Error:       ss.Error,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(step)
	}
}

// handleStepLogs handles GET /jobs/{jobId}/steps/{stepId}/logs
// @Summary      Get step logs
// @Description  Get logs for a specific workflow step
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobId   path      string  true  "Job ID"
// @Param        stepId  path      string  true  "Step ID"
// @Success      200     {object}  JobLogsResponse
// @Failure      404     {string}  string  "Step not found"
// @Router       /jobs/{jobId}/steps/{stepId}/logs [get]
func handleStepLogs(repo repository.IStepRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")
		stepID := chi.URLParam(r, "stepId")

		// TODO: Implement step logs retrieval
		_, _ = jobID, stepID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// handleJobResult handles GET /jobs/{jobId}/result
// @Summary      Get job result
// @Description  Get the latest result summary for a job
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobId   path      string  true  "Job ID"
// @Success      200     {object}  JobResult
// @Failure      404     {string}  string  "Job not found"
// @Router       /jobs/{jobId}/result [get]
func handleJobResult(repo repository.IJobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")

		// TODO: Implement result summary retrieval
		_ = jobID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

