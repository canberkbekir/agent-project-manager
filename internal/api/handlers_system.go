package api

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
)

var (
	// These can be set at build time via ldflags
	appVersion = "0.0.0"
	appCommit  = "unknown"
)

// handleHealthz returns a liveness health check response
// @Summary      Health check (liveness)
// @Description  Returns 200 OK if the service is alive
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "OK"
// @Router       /healthz [get]
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// handleReadyz returns a readiness health check response
// @Summary      Readiness check
// @Description  Returns 200 OK if the service is ready to accept requests
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "OK"
// @Router       /readyz [get]
func handleReadyz(w http.ResponseWriter, r *http.Request) {
	// TODO: Add readiness checks (e.g., database connectivity, queue status)
	w.WriteHeader(http.StatusOK)
}

// handleVersion returns version information
// @Summary      Get version information
// @Description  Returns the service name, version, and commit information
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  VersionResponse
// @Router       /version [get]
func handleVersion(w http.ResponseWriter, r *http.Request) {
	commit := appCommit
	if commit == "unknown" {
		if info, ok := debug.ReadBuildInfo(); ok {
			commit = extractCommitFromBuildInfo(info)
		}
	}

	response := VersionResponse{
		Name:    "agentd",
		Version: appVersion,
		Commit:  commit,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// extractCommitFromBuildInfo attempts to extract commit from build info
func extractCommitFromBuildInfo(info *debug.BuildInfo) string {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			if len(setting.Value) > 7 {
				return setting.Value[:7]
			}
			return setting.Value
		}
	}
	return "unknown"
}

