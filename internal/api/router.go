package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	// EnableTracing controls whether OpenTelemetry tracing middleware is applied
	EnableTracing bool
)

// Router returns a new HTTP router with all routes configured.
func Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	
	// OpenTelemetry tracing middleware (if enabled)
	if EnableTracing {
		r.Use(func(next http.Handler) http.Handler {
			return TracingMiddleware(next)
		})
	}
	
	r.Use(RequestLogger()) // Custom logger using application logger
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler())

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// System endpoints
		r.Get("/healthz", handleHealthz)
		r.Get("/readyz", handleReadyz)
		r.Get("/version", handleVersion)

		// Auth endpoints (optional)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handleLogin)
			r.Post("/refresh", handleRefresh)
			r.Post("/logout", handleLogout)
		})

		// Jobs endpoints
		r.Route("/jobs", func(r chi.Router) {
			r.Post("/", handleCreateJob)
			r.Get("/", handleListJobs)
			r.Get("/{jobId}", handleGetJob)
			r.Delete("/{jobId}", handleDeleteJob)
			r.Post("/{jobId}/retry", handleRetryJob)
			r.Get("/{jobId}/events", handleJobEvents)
			r.Get("/{jobId}/logs", handleJobLogs)
			r.Get("/{jobId}/result", handleJobResult)

			// Job steps
			r.Get("/{jobId}/steps", handleJobSteps)
			r.Get("/{jobId}/steps/{stepId}", handleGetStep)
			r.Get("/{jobId}/steps/{stepId}/logs", handleStepLogs)
		})

		// Runs endpoints
		r.Route("/runs", func(r chi.Router) {
			r.Post("/", handleCreateRun)
			r.Get("/", handleListRuns)
			r.Get("/{runId}", handleGetRun)
			r.Delete("/{runId}", handleDeleteRun)
		})

		// Workflows endpoints
		r.Route("/workflows", func(r chi.Router) {
			r.Get("/", handleListWorkflows)
			r.Get("/{name}", handleGetWorkflow)
			r.Post("/validate", handleValidateWorkflow)
		})

		// Artifacts endpoints
		r.Route("/artifacts", func(r chi.Router) {
			r.Get("/", handleListArtifacts)
			r.Get("/{artifactId}", handleGetArtifact)
			r.Get("/{artifactId}/download", handleDownloadArtifact)
			r.Delete("/{artifactId}", handleDeleteArtifact)
		})

		// Agents endpoints
		r.Route("/agents", func(r chi.Router) {
			r.Get("/", handleListAgents)
			r.Get("/{agentId}", handleGetAgent)
			r.Get("/{agentId}/status", handleGetAgentStatus)
			r.Post("/{agentId}/drain", handleDrainAgent)
			r.Post("/{agentId}/resume", handleResumeAgent)
		})

		// Queue endpoints
		r.Route("/queue", func(r chi.Router) {
			r.Get("/", handleGetQueue)
			r.Get("/items", handleListQueueItems)
			r.Post("/requeue", handleRequeue)
		})
	})

	return r
}
