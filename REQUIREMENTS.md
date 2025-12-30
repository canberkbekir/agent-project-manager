# System Requirements Checklist

Based on the architecture diagram, here's what you need to implement to complete the AI-powered code generation and review system:

## ✅ Already Implemented

- **API Gateway (HTTP)**: HTTP API server with routing, middleware, and Swagger documentation
- **CLI Client (`agentctl`)**: Basic CLI client structure exists
- **Configuration System**: YAML-based config with support for LLM providers
- **Observability**: Logging, metrics (Prometheus), and tracing (OpenTelemetry) infrastructure
- **Authentication**: Basic auth endpoints structure

## ❌ Missing Core Components

### 1. **Orchestrator (Workflow Engine)**
   - [ ] State machine for workflow execution
   - [ ] Workflow definition parser/validator
   - [ ] Step execution and coordination
   - [ ] Retry logic and error handling
   - [ ] Workflow state persistence
   - [ ] Integration with job queue and agents

### 2. **Job Queue (In-Process)**
   - [ ] In-memory job queue implementation
   - [ ] Worker pool for concurrent job processing
   - [ ] Job prioritization and scheduling
   - [ ] Queue persistence (optional, for recovery)
   - [ ] Job status tracking

### 3. **Shared State (SQLite/Files)**
   - [ ] SQLite database schema and migrations
   - [ ] Database models for:
     - Jobs
     - Runs
     - Workflows
     - Agents
     - Artifacts
     - Steps/Events
   - [ ] Database store implementation (CRUD operations)
   - [ ] Connection pooling and transaction management

### 4. **Artifact Store (Repo/Workdir)**
   - [ ] Workspace directory management
   - [ ] Artifact storage and retrieval
   - [ ] File system operations for generated code
   - [ ] Artifact versioning/metadata
   - [ ] Cleanup and garbage collection

### 5. **LLM Gateway (Provider Adapter)**
   - [ ] LLM client interface/abstraction
   - [ ] OpenAI adapter implementation
   - [ ] Ollama adapter implementation
   - [ ] LM Studio adapter (optional)
   - [ ] Provider switching logic
   - [ ] Request/response handling
   - [ ] Error handling and retries
   - [ ] Rate limiting (if needed)

### 6. **AI Agents**

   #### **Architect Agent**
   - [ ] Agent implementation
   - [ ] Integration with LLM gateway
   - [ ] Tool integration for:
     - Specification checks
     - Template management
   - [ ] Architecture decision logic

   #### **Code Gen Agent**
   - [ ] Agent implementation
   - [ ] Integration with LLM gateway
   - [ ] Tool integration for:
     - Project scaffolding
     - Build processes
     - Test execution
   - [ ] Code generation logic

   #### **Code Review Agent**
   - [ ] Agent implementation
   - [ ] Integration with LLM gateway
   - [ ] Tool integration for:
     - Linters (golangci-lint, etc.)
     - Static analysis tools
   - [ ] Review generation logic

### 7. **Tools Integration**

   #### **Git Tools**
   - [ ] Git repository operations
   - [ ] Clone/checkout functionality
   - [ ] Commit/push operations
   - [ ] Branch management

   #### **Build Tools**
   - [ ] Go build execution
   - [ ] Multi-language build support (if needed)
   - [ ] Build artifact collection

   #### **Test Tools**
   - [ ] Test execution (go test, etc.)
   - [ ] Test result parsing
   - [ ] Coverage reporting

   #### **Linting/Static Analysis**
   - [ ] Linter execution (golangci-lint, etc.)
   - [ ] Result parsing and formatting
   - [ ] Integration with review agent

   #### **Spec/Template Tools**
   - [ ] Template engine
   - [ ] Specification validation
   - [ ] Template management

### 8. **Client Interfaces**

   #### **CLI Client (`agentctl`)**
   - [ ] Complete command implementations:
     - Job submission
     - Job status checking
     - Artifact download
     - Workflow management
     - Agent management
   - [ ] Configuration management
   - [ ] Output formatting (JSON, table, etc.)

   #### **Web Interface (`agentweb`)**
   - [ ] Web UI implementation
   - [ ] Dashboard for job monitoring
   - [ ] Workflow visualization
   - [ ] Artifact browser
   - [ ] Real-time updates (WebSocket or polling)

   #### **Discord Bot (`agentbot-discord`)**
   - [ ] Discord bot integration
   - [ ] Command handlers
   - [ ] Job submission via Discord
   - [ ] Status notifications
   - [ ] Interactive job management

   #### **Git Hooks (`agenthook`)**
   - [ ] Pre-commit hook implementation
   - [ ] Pre-push hook implementation
   - [ ] Integration with code review agent
   - [ ] Hook configuration management

## 🔧 Infrastructure & Setup

### Hardware/Environment
- [ ] Raspberry Pi 4GB setup
- [ ] Network configuration
- [ ] Storage setup for artifacts and state

### External Services
- [ ] Local LLM setup (Ollama/LM Studio on your PC)
  - [ ] Ollama installation and configuration
  - [ ] Model downloads
  - [ ] Network connectivity from Pi to PC
- [ ] External LLM provider setup (optional)
  - [ ] OpenAI API key configuration
  - [ ] Other provider credentials

### Development Tools
- [ ] Go toolchain (already have)
- [ ] Database migration tool (goose or similar)
- [ ] Testing framework setup
- [ ] CI/CD setup (optional)

## 📋 Implementation Priority

### 🎯 **START HERE: Phase 1 - Core Foundation**

**Start with #1 - Shared State (SQLite)** - This is your foundation. Everything else depends on it.

#### **1. Shared State (SQLite)** ⭐ START HERE
   - **Why first?** Every component needs persistence (jobs, runs, workflows, artifacts, agents)
   - **Dependencies:** None (just SQLite driver)
   - **What to build:**
     - Database schema (migrations)
     - Models (Job, Run, Workflow, Artifact, Agent, Step, Event)
     - Store interface and implementation (CRUD operations)
   - **Can test immediately:** Yes - write unit tests for store operations

#### **2. Artifact Store**
   - **Why second?** Independent, needed early, relatively simple
   - **Dependencies:** None (just file system)
   - **What to build:**
     - Workspace directory management
     - Artifact storage/retrieval
     - File operations
   - **Can test immediately:** Yes - test file operations

#### **3. LLM Gateway**
   - **Why third?** Independent, needed by agents (which come later)
   - **Dependencies:** None (just HTTP clients)
   - **What to build:**
     - Client interface/abstraction
     - OpenAI adapter
     - Ollama adapter
     - Provider switching
   - **Can test immediately:** Yes - mock HTTP responses, test with real providers

#### **4. Tools (Basic)**
   - **Why fourth?** Needed by agents, but can build incrementally
   - **Dependencies:** None (just external command execution)
   - **What to build first:**
     - Git tools (clone, checkout, commit)
     - Exec wrapper for running commands
   - **Can test immediately:** Yes - test with real git repos

#### **5. Job Queue**
   - **Why fifth?** Needed by orchestrator
   - **Dependencies:** Shared State (to persist queue items)
   - **What to build:**
     - In-memory queue
     - Worker pool
     - Job status tracking
   - **Can test immediately:** Yes - test queue operations

#### **6. Orchestrator**
   - **Why sixth?** Coordinates everything, needs queue and state
   - **Dependencies:** Job Queue, Shared State, Artifact Store
   - **What to build:**
     - Workflow engine
     - State machine
     - Step execution
     - Retry logic
   - **Can test immediately:** Yes - test with simple workflows

#### **7. Agents**
   - **Why last?** Need all above components
   - **Dependencies:** LLM Gateway, Tools, Orchestrator, Shared State
   - **What to build:**
     - Architect Agent
     - Code Gen Agent
     - Code Review Agent
   - **Can test immediately:** Yes - test with mock LLM responses

### Phase 2: Client Interfaces
8. Complete CLI client (needs working API)
9. Web interface (needs working API)
10. Discord bot (needs working API)
11. Git hooks (needs tools and review agent)

### Phase 3: Polish & Optimization
12. Error handling improvements
13. Performance optimization
14. Documentation
15. Monitoring and alerting

## 🚀 Quick Start Guide

**Step 1:** Implement `internal/state/` package
   - Create database models
   - Write migrations
   - Implement store interface
   - Test with simple CRUD operations

**Step 2:** Wire state into API handlers
   - Update `handleCreateJob` to use state store
   - Update `handleListJobs` to query from state
   - Test via API endpoints

**Step 3:** Continue with Artifact Store, then LLM Gateway, etc.

## 📝 Notes

- Many directories exist but are empty (agents/, tools/, orchestrator/, queue/, state/, artifact/, llm/)
- API handlers exist but many are stubs (TODO comments)
- Discord bot, web interface, and git hooks are placeholder implementations
- The system architecture is well-defined, but most core logic needs implementation

