package api

import (
	"github.com/google/uuid"
)

// generateID generates a UUID string for use in IDs
func generateID() string {
	return uuid.New().String()
}

