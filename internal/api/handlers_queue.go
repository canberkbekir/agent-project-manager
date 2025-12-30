package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"agent-project-manager/internal/repository"
)

// handleGetQueue handles GET /queue
// @Summary      Get queue statistics
// @Description  Get queue statistics and metrics
// @Tags         queue
// @Accept       json
// @Produce      json
// @Success      200  {object}  QueueStats
// @Router       /queue [get]
func handleGetQueue(repo repository.IQueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get queue stats from repository
		stats, err := repo.GetQueueStats()
		if err != nil {
			http.Error(w, "Failed to get queue stats: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := QueueStats{
			Pending: stats.Pending,
			Leased:  stats.Leased,
			Done:    stats.Done,
			Dead:    stats.Dead,
			Total:   stats.Total,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// handleListQueueItems handles GET /queue/items
// @Summary      List queue items
// @Description  Get a paginated list of queue items with optional filtering
// @Tags         queue
// @Accept       json
// @Produce      json
// @Param        limit   query     int     false  "Maximum number of items to return"
// @Param        cursor  query     string  false  "Cursor for pagination"
// @Param        state   query     string  false  "Filter by state (pending|leased|done|dead)"
// @Success      200     {object}  QueueItemListResponse
// @Router       /queue/items [get]
func handleListQueueItems(repo repository.IQueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		cursor := r.URL.Query().Get("cursor")
		state := r.URL.Query().Get("state")

		limit := 50 // default
		if limitStr != "" {
			if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		// Get queue items from repository
		stateItems, nextCursor, err := repo.ListQueueItems(state, limit, cursor)
		if err != nil {
			http.Error(w, "Failed to list queue items: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert state models to API models
		items := make([]QueueItem, len(stateItems))
		for i, si := range stateItems {
			queueState, _ := QueueStateFromString(si.State)
			items[i] = QueueItem{
				ID:          si.ID,
				JobID:       si.JobID,
				State:       queueState,
				Data:        map[string]interface{}(si.Data),
				CreatedAt:   si.CreatedAt,
				UpdatedAt:   si.UpdatedAt,
				LeasedAt:    si.LeasedAt,
				CompletedAt: si.CompletedAt,
			}
		}

		response := QueueItemListResponse{
			Items:   items,
			Cursor:  nextCursor,
			HasMore: nextCursor != "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// handleRequeue handles POST /queue/requeue
// @Summary      Requeue a job
// @Description  Requeue a job in the queue
// @Tags         queue
// @Accept       json
// @Produce      json
// @Param        requeue  body      RequeueRequest  true  "Requeue request"
// @Success      200      {string}  string  "OK"
// @Failure      400      {string}  string  "Invalid request"
// @Router       /queue/requeue [post]
func handleRequeue(repo repository.IQueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequeueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get queue item
		item, err := repo.GetQueueItem(req.JobID)
		if err != nil {
			http.Error(w, "Queue item not found", http.StatusNotFound)
			return
		}

		// Update state to pending to requeue
		item.State = "pending"
		item.LeasedAt = nil
		item.CompletedAt = nil

		if err := repo.UpdateQueueItem(item); err != nil {
			http.Error(w, "Failed to requeue: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

