package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UUID is a type alias for uuid.UUID to use throughout the API
type UUID = uuid.UUID

// Enum types

// JobStatus represents the status of a job
type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusSucceeded JobStatus = "succeeded"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// String returns the string representation of JobStatus
func (s JobStatus) String() string {
	return string(s)
}

// IsValid checks if the JobStatus value is valid
func (s JobStatus) IsValid() bool {
	switch s {
	case JobStatusQueued, JobStatusRunning, JobStatusSucceeded, JobStatusFailed, JobStatusCancelled:
		return true
	default:
		return false
	}
}

// JobStatusFromString parses a string into a JobStatus
func JobStatusFromString(s string) (JobStatus, bool) {
	status := JobStatus(s)
	return status, status.IsValid()
}

// AllJobStatuses returns all valid JobStatus values
func AllJobStatuses() []JobStatus {
	return []JobStatus{
		JobStatusQueued,
		JobStatusRunning,
		JobStatusSucceeded,
		JobStatusFailed,
		JobStatusCancelled,
	}
}

// UnmarshalJSON implements json.Unmarshaler for JobStatus
func (s *JobStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	status, ok := JobStatusFromString(str)
	if !ok {
		return fmt.Errorf("invalid JobStatus: %s", str)
	}
	*s = status
	return nil
}

// RunStatus represents the status of a run
type RunStatus string

const (
	RunStatusPending   RunStatus = "pending"
	RunStatusRunning   RunStatus = "running"
	RunStatusSucceeded RunStatus = "succeeded"
	RunStatusFailed    RunStatus = "failed"
	RunStatusCancelled RunStatus = "cancelled"
)

// String returns the string representation of RunStatus
func (s RunStatus) String() string {
	return string(s)
}

// IsValid checks if the RunStatus value is valid
func (s RunStatus) IsValid() bool {
	switch s {
	case RunStatusPending, RunStatusRunning, RunStatusSucceeded, RunStatusFailed, RunStatusCancelled:
		return true
	default:
		return false
	}
}

// RunStatusFromString parses a string into a RunStatus
func RunStatusFromString(s string) (RunStatus, bool) {
	status := RunStatus(s)
	return status, status.IsValid()
}

// AllRunStatuses returns all valid RunStatus values
func AllRunStatuses() []RunStatus {
	return []RunStatus{
		RunStatusPending,
		RunStatusRunning,
		RunStatusSucceeded,
		RunStatusFailed,
		RunStatusCancelled,
	}
}

// UnmarshalJSON implements json.Unmarshaler for RunStatus
func (s *RunStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	status, ok := RunStatusFromString(str)
	if !ok {
		return fmt.Errorf("invalid RunStatus: %s", str)
	}
	*s = status
	return nil
}

// StepStatus represents the status of a workflow step
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusSucceeded StepStatus = "succeeded"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped"
)

// String returns the string representation of StepStatus
func (s StepStatus) String() string {
	return string(s)
}

// IsValid checks if the StepStatus value is valid
func (s StepStatus) IsValid() bool {
	switch s {
	case StepStatusPending, StepStatusRunning, StepStatusSucceeded, StepStatusFailed, StepStatusSkipped:
		return true
	default:
		return false
	}
}

// StepStatusFromString parses a string into a StepStatus
func StepStatusFromString(s string) (StepStatus, bool) {
	status := StepStatus(s)
	return status, status.IsValid()
}

// AllStepStatuses returns all valid StepStatus values
func AllStepStatuses() []StepStatus {
	return []StepStatus{
		StepStatusPending,
		StepStatusRunning,
		StepStatusSucceeded,
		StepStatusFailed,
		StepStatusSkipped,
	}
}

// UnmarshalJSON implements json.Unmarshaler for StepStatus
func (s *StepStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	status, ok := StepStatusFromString(str)
	if !ok {
		return fmt.Errorf("invalid StepStatus: %s", str)
	}
	*s = status
	return nil
}

// QueueState represents the state of a queue item
type QueueState string

const (
	QueueStatePending QueueState = "pending"
	QueueStateLeased  QueueState = "leased"
	QueueStateDone    QueueState = "done"
	QueueStateDead    QueueState = "dead"
)

// String returns the string representation of QueueState
func (s QueueState) String() string {
	return string(s)
}

// IsValid checks if the QueueState value is valid
func (s QueueState) IsValid() bool {
	switch s {
	case QueueStatePending, QueueStateLeased, QueueStateDone, QueueStateDead:
		return true
	default:
		return false
	}
}

// QueueStateFromString parses a string into a QueueState
func QueueStateFromString(s string) (QueueState, bool) {
	state := QueueState(s)
	return state, state.IsValid()
}

// AllQueueStates returns all valid QueueState values
func AllQueueStates() []QueueState {
	return []QueueState{
		QueueStatePending,
		QueueStateLeased,
		QueueStateDone,
		QueueStateDead,
	}
}

// UnmarshalJSON implements json.Unmarshaler for QueueState
func (s *QueueState) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	state, ok := QueueStateFromString(str)
	if !ok {
		return fmt.Errorf("invalid QueueState: %s", str)
	}
	*s = state
	return nil
}

// ArtifactType represents the type of an artifact
type ArtifactType string

const (
	ArtifactTypePDF  ArtifactType = "pdf"
	ArtifactTypePPTX ArtifactType = "pptx"
	ArtifactTypeZIP  ArtifactType = "zip"
	ArtifactTypeLog  ArtifactType = "log"
)

// String returns the string representation of ArtifactType
func (t ArtifactType) String() string {
	return string(t)
}

// IsValid checks if the ArtifactType value is valid
func (t ArtifactType) IsValid() bool {
	switch t {
	case ArtifactTypePDF, ArtifactTypePPTX, ArtifactTypeZIP, ArtifactTypeLog:
		return true
	default:
		return false
	}
}

// ArtifactTypeFromString parses a string into an ArtifactType
func ArtifactTypeFromString(s string) (ArtifactType, bool) {
	artifactType := ArtifactType(s)
	return artifactType, artifactType.IsValid()
}

// AllArtifactTypes returns all valid ArtifactType values
func AllArtifactTypes() []ArtifactType {
	return []ArtifactType{
		ArtifactTypePDF,
		ArtifactTypePPTX,
		ArtifactTypeZIP,
		ArtifactTypeLog,
	}
}

// UnmarshalJSON implements json.Unmarshaler for ArtifactType
func (t *ArtifactType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	artifactType, ok := ArtifactTypeFromString(str)
	if !ok {
		return fmt.Errorf("invalid ArtifactType: %s", str)
	}
	*t = artifactType
	return nil
}

// System Models

// VersionResponse represents the version information response
type VersionResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

// Auth Models

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

// RefreshRequest represents a refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// Job Models

// CreateJobRequest represents a job creation request
type CreateJobRequest struct {
	Workflow string                 `json:"workflow"`
	Input    map[string]interface{} `json:"input"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

// CreateJobResponse represents a job creation response
type CreateJobResponse struct {
	ID string `json:"id"`
}

// Job represents a job entity
type Job struct {
	ID        string                 `json:"id"`
	Workflow  string                 `json:"workflow"`
	Status    JobStatus              `json:"status"`
	Input     map[string]interface{} `json:"input"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	StartedAt *time.Time             `json:"startedAt,omitempty"`
	CompletedAt *time.Time           `json:"completedAt,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// JobListResponse represents a paginated list of jobs
type JobListResponse struct {
	Jobs  []Job  `json:"jobs"`
	Cursor string `json:"cursor,omitempty"`
	HasMore bool  `json:"hasMore"`
}

// JobStep represents a workflow step
type JobStep struct {
	ID        string                 `json:"id"`
	JobID     string                 `json:"jobId"`
	Name      string                 `json:"name"`
	Status    StepStatus             `json:"status"`
	Input     map[string]interface{} `json:"input,omitempty"`
	Output    map[string]interface{} `json:"output,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	StartedAt *time.Time             `json:"startedAt,omitempty"`
	CompletedAt *time.Time           `json:"completedAt,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// JobResult represents the result summary of a job
type JobResult struct {
	JobID     string                 `json:"jobId"`
	Status    JobStatus              `json:"status"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Artifacts []string               `json:"artifacts,omitempty"` // artifact IDs
	CompletedAt *time.Time           `json:"completedAt,omitempty"`
}

// JobLogsResponse represents job logs
type JobLogsResponse struct {
	JobID string   `json:"jobId"`
	Logs  []LogEntry `json:"logs"`
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	StepID    string    `json:"stepId,omitempty"`
}

// Run Models

// CreateRunRequest represents a run creation request
type CreateRunRequest struct {
	JobID  string                 `json:"jobId"`
	Params map[string]interface{} `json:"params"`
}

// CreateRunResponse represents a run creation response
type CreateRunResponse struct {
	ID string `json:"id"`
}

// Run represents a run entity
type Run struct {
	ID        string                 `json:"id"`
	JobID     string                 `json:"jobId"`
	Status    RunStatus              `json:"status"`
	Params    map[string]interface{} `json:"params"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	StartedAt *time.Time             `json:"startedAt,omitempty"`
	CompletedAt *time.Time           `json:"completedAt,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// RunListResponse represents a paginated list of runs
type RunListResponse struct {
	Runs    []Run  `json:"runs"`
	Cursor  string `json:"cursor,omitempty"`
	HasMore bool   `json:"hasMore"`
}

// Workflow Models

// Workflow represents a workflow definition
type Workflow struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
	Version     string                 `json:"version,omitempty"`
}

// WorkflowListResponse represents a list of workflows
type WorkflowListResponse struct {
	Workflows []Workflow `json:"workflows"`
}

// ValidateWorkflowRequest represents a workflow validation request
type ValidateWorkflowRequest struct {
	Workflow string                 `json:"workflow"`
	Input    map[string]interface{} `json:"input"`
}

// ValidateWorkflowResponse represents workflow validation response
type ValidateWorkflowResponse struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

// Artifact Models

// Artifact represents an artifact entity
type Artifact struct {
	ID        string      `json:"id"`
	JobID     string      `json:"jobId,omitempty"`
	RunID     string      `json:"runId,omitempty"`
	Type      ArtifactType `json:"type"`
	Name      string      `json:"name"`
	Size      int64       `json:"size"`
	Path      string      `json:"path,omitempty"`
	CreatedAt time.Time   `json:"createdAt"`
}

// ArtifactListResponse represents a paginated list of artifacts
type ArtifactListResponse struct {
	Artifacts []Artifact `json:"artifacts"`
	Cursor    string     `json:"cursor,omitempty"`
	HasMore   bool       `json:"hasMore"`
}

// Agent Models

// Agent represents an agent/worker entity
type Agent struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Status    string                 `json:"status"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	LastSeen  time.Time              `json:"lastSeen"`
	CreatedAt time.Time              `json:"createdAt"`
}

// AgentListResponse represents a list of agents
type AgentListResponse struct {
	Agents []Agent `json:"agents"`
}

// AgentStatus represents agent status details
type AgentStatus struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	IsDraining    bool      `json:"isDraining"`
	ActiveJobs    int       `json:"activeJobs"`
	CompletedJobs int64     `json:"completedJobs"`
	LastSeen      time.Time `json:"lastSeen"`
}

// Queue Models

// QueueStats represents queue statistics
type QueueStats struct {
	Pending int `json:"pending"`
	Leased  int `json:"leased"`
	Done    int `json:"done"`
	Dead    int `json:"dead"`
	Total   int `json:"total"`
}

// QueueItem represents a queue item
type QueueItem struct {
	ID        string                 `json:"id"`
	JobID     string                 `json:"jobId"`
	State     QueueState             `json:"state"`
	Data      map[string]interface{} `json:"data,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	LeasedAt  *time.Time             `json:"leasedAt,omitempty"`
	CompletedAt *time.Time           `json:"completedAt,omitempty"`
}

// QueueItemListResponse represents a paginated list of queue items
type QueueItemListResponse struct {
	Items   []QueueItem `json:"items"`
	Cursor  string      `json:"cursor,omitempty"`
	HasMore bool        `json:"hasMore"`
}

// RequeueRequest represents a requeue request
type RequeueRequest struct {
	JobID  string `json:"jobId"`
	Reason string `json:"reason"`
}

// Helper functions

// NewUUID generates a new UUID string
func NewUUID() string {
	return uuid.New().String()
}

// ParseUUID parses a UUID string
func ParseUUID(s string) (UUID, error) {
	return uuid.Parse(s)
}

