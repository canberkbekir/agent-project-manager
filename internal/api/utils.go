package api

import (
	"agent-project-manager/internal/state"
)

// generateID generates a UUID string for use in IDs
func generateID() string {
	return state.NewUUID()
}

