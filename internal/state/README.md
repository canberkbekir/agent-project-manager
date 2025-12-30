# State Package

This package provides the SQLite-based persistence layer for the agent-project-manager system.

## Components

### Models (`models.go`)
Database models for all entities:
- `Job` - Workflow execution jobs
- `Run` - Individual run instances
- `Workflow` - Workflow definitions
- `Step` - Workflow step execution
- `Event` - Event logging
- `Artifact` - Generated artifacts
- `Agent` - Agent/worker instances
- `QueueItem` - Queue items

### Store (`store.go`)
The `Store` interface and SQLite implementation providing:
- Full CRUD operations for all entities
- Pagination support (cursor-based)
- JSON field handling for flexible data storage
- Connection pooling and transaction management

### Migrations (`migrations/0001_init.sql`)
Initial database schema with:
- All tables with proper foreign keys
- Indexes for performance
- JSON storage for flexible fields

## Usage

```go
import "agent-project-manager/internal/state"

// Create a new store
store, err := state.NewStore("data/state.db")
if err != nil {
    log.Fatal(err)
}
defer store.Close()

// Run migrations
if err := state.RunMigrations(store, "migrations"); err != nil {
    log.Fatal(err)
}

// Create a job
job := &state.Job{
    ID:       state.NewUUID(),
    Workflow: "test-workflow",
    Status:   "queued",
    Input:    state.JSONMap{"key": "value"},
    Meta:     state.JSONMap{},
}
err = store.CreateJob(job)

// Get a job
job, err := store.GetJob("job-id")

// List jobs with pagination
jobs, cursor, err := store.ListJobs(50, "", "queued", "")
```

## Database Schema

The schema includes:
- **jobs** - Main job table
- **runs** - Run instances (linked to jobs)
- **workflows** - Workflow definitions
- **steps** - Workflow steps (linked to jobs)
- **events** - Event log (linked to jobs/steps)
- **artifacts** - Artifacts (linked to jobs/runs)
- **agents** - Agent registry
- **queue_items** - Queue management

All tables use proper foreign keys and indexes for performance.

