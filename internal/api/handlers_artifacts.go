package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"agent-project-manager/internal/repository"
)

// handleListArtifacts handles GET /artifacts
// @Summary      List artifacts
// @Description  Get a paginated list of artifacts with optional filtering
// @Tags         artifacts
// @Accept       json
// @Produce      json
// @Param        jobId   query     string  false  "Filter by job ID"
// @Param        runId   query     string  false  "Filter by run ID"
// @Param        limit   query     int     false  "Maximum number of artifacts to return"
// @Param        cursor  query     string  false  "Cursor for pagination"
// @Param        type    query     string  false  "Filter by artifact type (pdf|pptx|zip|log)"
// @Success      200     {object}  ArtifactListResponse
// @Router       /artifacts [get]
func handleListArtifacts(repo repository.IArtifactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		jobID := r.URL.Query().Get("jobId")
		runID := r.URL.Query().Get("runId")
		limitStr := r.URL.Query().Get("limit")
		cursor := r.URL.Query().Get("cursor")

		limit := 50 // default
		if limitStr != "" {
			if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		// Get artifacts from repository
		stateArtifacts, nextCursor, err := repo.ListArtifacts(jobID, runID, limit, cursor)
		if err != nil {
			http.Error(w, "Failed to list artifacts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert state models to API models
		artifacts := make([]Artifact, len(stateArtifacts))
		for i, sa := range stateArtifacts {
			artifactType, _ := ArtifactTypeFromString(sa.Type)
			artifacts[i] = Artifact{
				ID:        sa.ID,
				JobID:     sa.JobID,
				RunID:     sa.RunID,
				Type:      artifactType,
				Name:      sa.Name,
				Size:      sa.Size,
				Path:      sa.Path,
				CreatedAt: sa.CreatedAt,
			}
		}

		response := ArtifactListResponse{
			Artifacts: artifacts,
			Cursor:    nextCursor,
			HasMore:   nextCursor != "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// handleGetArtifact handles GET /artifacts/{artifactId}
// @Summary      Get artifact metadata
// @Description  Get metadata for a specific artifact
// @Tags         artifacts
// @Accept       json
// @Produce      json
// @Param        artifactId  path      string  true  "Artifact ID"
// @Success      200         {object}  Artifact
// @Failure      404         {string}  string  "Artifact not found"
// @Router       /artifacts/{artifactId} [get]
func handleGetArtifact(repo repository.IArtifactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artifactId := chi.URLParam(r, "artifactId")

		// Get artifact from repository
		sa, err := repo.GetArtifact(artifactId)
		if err != nil {
			http.Error(w, "Artifact not found", http.StatusNotFound)
			return
		}

		// Convert state model to API model
		artifactType, _ := ArtifactTypeFromString(sa.Type)
		artifact := Artifact{
			ID:        sa.ID,
			JobID:     sa.JobID,
			RunID:     sa.RunID,
			Type:      artifactType,
			Name:      sa.Name,
			Size:      sa.Size,
			Path:      sa.Path,
			CreatedAt: sa.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(artifact)
	}
}

// handleDownloadArtifact handles GET /artifacts/{artifactId}/download
// @Summary      Download artifact
// @Description  Download the artifact file
// @Tags         artifacts
// @Accept       json
// @Produce      application/octet-stream
// @Param        artifactId  path      string  true  "Artifact ID"
// @Success      200         {file}    file    "Artifact file"
// @Failure      404         {string}  string  "Artifact not found"
// @Router       /artifacts/{artifactId}/download [get]
func handleDownloadArtifact(repo repository.IArtifactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artifactId := chi.URLParam(r, "artifactId")

		// Get artifact from repository
		sa, err := repo.GetArtifact(artifactId)
		if err != nil {
			http.Error(w, "Artifact not found", http.StatusNotFound)
			return
		}

		// TODO: Implement artifact file streaming from sa.Path
		_ = sa

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
	}
}

// handleDeleteArtifact handles DELETE /artifacts/{artifactId}
// @Summary      Delete artifact
// @Description  Delete an artifact
// @Tags         artifacts
// @Accept       json
// @Produce      json
// @Param        artifactId  path      string  true  "Artifact ID"
// @Success      204         {string}  string  "No Content"
// @Failure      404         {string}  string  "Artifact not found"
// @Router       /artifacts/{artifactId} [delete]
func handleDeleteArtifact(repo repository.IArtifactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artifactId := chi.URLParam(r, "artifactId")

		// Delete artifact from repository
		if err := repo.DeleteArtifact(artifactId); err != nil {
			http.Error(w, "Artifact not found", http.StatusNotFound)
			return
		}

		// TODO: Implement artifact file cleanup logic

		w.WriteHeader(http.StatusNoContent)
	}
}

