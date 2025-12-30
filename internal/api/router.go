package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	
	"agent-project-manager/internal/obs"
	"agent-project-manager/internal/repository"
	"agent-project-manager/internal/state"
)

var (
	// EnableTracing controls whether OpenTelemetry tracing middleware is applied
	EnableTracing bool
	// PrometheusMetricsPath is the HTTP path for Prometheus metrics endpoint (empty to disable)
	PrometheusMetricsPath string
)

// Router returns a new HTTP router with all routes configured.
// It accepts a repository for database operations.
func Router(repo state.Repository) http.Handler {
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

	// Prometheus metrics endpoint (for Grafana)
	if PrometheusMetricsPath != "" {
		r.Get(PrometheusMetricsPath, obs.GetPrometheusHandler().ServeHTTP)
	}

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

		// Create repositories from database connection
		db := repo.GetDB()
		jobRepo := repository.NewJobRepository(db)
		runRepo := repository.NewRunRepository(db)
		agentRepo := repository.NewAgentRepository(db)
		stepRepo := repository.NewStepRepository(db)
		workflowRepo := repository.NewWorkflowRepository(db)
		artifactRepo := repository.NewArtifactRepository(db)
		queueRepo := repository.NewQueueRepository(db)

		// Jobs endpoints
		r.Route("/jobs", func(r chi.Router) {
			r.Post("/", handleCreateJob(jobRepo))
			r.Get("/", handleListJobs(jobRepo))
			r.Get("/{jobId}", handleGetJob(jobRepo))
			r.Delete("/{jobId}", handleDeleteJob(jobRepo))
			r.Post("/{jobId}/retry", handleRetryJob(jobRepo))
			r.Get("/{jobId}/events", handleJobEvents(jobRepo))
			r.Get("/{jobId}/logs", handleJobLogs(jobRepo))
			r.Get("/{jobId}/result", handleJobResult(jobRepo))

			// Job steps
			r.Get("/{jobId}/steps", handleJobSteps(stepRepo))
			r.Get("/{jobId}/steps/{stepId}", handleGetStep(stepRepo))
			r.Get("/{jobId}/steps/{stepId}/logs", handleStepLogs(stepRepo))
		})

		// Runs endpoints
		r.Route("/runs", func(r chi.Router) {
			r.Post("/", handleCreateRun(runRepo))
			r.Get("/", handleListRuns(runRepo))
			r.Get("/{runId}", handleGetRun(runRepo))
			r.Delete("/{runId}", handleDeleteRun(runRepo))
		})

		// Workflows endpoints
		r.Route("/workflows", func(r chi.Router) {
			r.Get("/", handleListWorkflows(workflowRepo))
			r.Get("/{name}", handleGetWorkflow(workflowRepo))
			r.Post("/validate", handleValidateWorkflow(workflowRepo))
		})

		// Artifacts endpoints
		r.Route("/artifacts", func(r chi.Router) {
			r.Get("/", handleListArtifacts(artifactRepo))
			r.Get("/{artifactId}", handleGetArtifact(artifactRepo))
			r.Get("/{artifactId}/download", handleDownloadArtifact(artifactRepo))
			r.Delete("/{artifactId}", handleDeleteArtifact(artifactRepo))
		})

		// Agents endpoints
		r.Route("/agents", func(r chi.Router) {
			r.Get("/", handleListAgents(agentRepo))
			r.Get("/{agentId}", handleGetAgent(agentRepo))
			r.Get("/{agentId}/status", handleGetAgentStatus(agentRepo))
			r.Post("/{agentId}/drain", handleDrainAgent(agentRepo))
			r.Post("/{agentId}/resume", handleResumeAgent(agentRepo))
		})

		// Queue endpoints
		r.Route("/queue", func(r chi.Router) {
			r.Get("/", handleGetQueue(queueRepo))
			r.Get("/items", handleListQueueItems(queueRepo))
			r.Post("/requeue", handleRequeue(queueRepo))
		})
	})

	return r
}
