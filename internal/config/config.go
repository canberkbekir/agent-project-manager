package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultPath = "configs/config.yaml"

type Config struct {
	API       APIConfig       `yaml:"api"`
	State     StateConfig     `yaml:"state"`
	Queue     QueueConfig     `yaml:"queue"`
	Artifacts ArtifactsConfig `yaml:"artifacts"`
	LLM       LLMConfig       `yaml:"llm"`
	Auth      AuthConfig      `yaml:"auth"`
	Logger    LoggerConfig    `yaml:"logger"`
	Obs       ObsConfig       `yaml:"obs"`
}

type APIConfig struct {
	Addr    string `yaml:"addr"`
	BaseURL string `yaml:"baseURL"`
}

type AuthConfig struct {
	Token string `yaml:"token"`
}

type StateConfig struct {
	ConnectionString string `yaml:"connectionString"`
}

type QueueConfig struct {
	Workers int `yaml:"workers"`
}

type ArtifactsConfig struct {
	WorkDir string `yaml:"workDir"`
}

type LLMConfig struct {
	Provider string       `yaml:"provider"`
	OpenAI   OpenAIConfig  `yaml:"openai"`
	Ollama   OllamaConfig  `yaml:"ollama"`
}

type OpenAIConfig struct {
	APIKey string `yaml:"apiKey"`
	Model  string `yaml:"model"`
}

type OllamaConfig struct {
	BaseURL string `yaml:"baseURL"`
	Model   string `yaml:"model"`
}

type LoggerConfig struct {
	Level        string `yaml:"level"`        // debug, info, warn, error, fatal (default: info)
	Format       string `yaml:"format"`       // text, json (default: text)
	Output       string `yaml:"output"`       // stdout, stderr, or file path (default: stdout)
	ReportCaller bool   `yaml:"reportCaller"` // include caller information (default: false)
}

type ObsConfig struct {
	Tracing TracingConfig `yaml:"tracing"`
	Metrics MetricsConfig `yaml:"metrics"`
}

type TracingConfig struct {
	Endpoint string `yaml:"endpoint"` // OTLP endpoint (e.g., "localhost:4318"), empty or "none" to disable
	Enabled  bool   `yaml:"enabled"`  // Enable/disable tracing (default: false)
}

type MetricsConfig struct {
	Endpoint          string `yaml:"endpoint"`          // OTLP endpoint (e.g., "localhost:4318"), empty or "none" to disable
	PrometheusEnabled bool   `yaml:"prometheusEnabled"` // Enable Prometheus metrics endpoint for Grafana (default: false)
	PrometheusPath    string `yaml:"prometheusPath"`    // HTTP path for Prometheus metrics endpoint (default: "/metrics")
	Enabled           bool   `yaml:"enabled"`          // Enable/disable metrics (default: false)
}

// Load loads configuration following 12-Factor App principles:
// 1. Environment variables take precedence over config file
// 2. Config file is optional (used for defaults in development)
// 3. If CONFIG_PATH is set, use that file; otherwise use default
func Load() (Config, error) {
	var c Config

	// Try to load from config file (optional)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = DefaultPath
	}

	// Load config file if it exists (not an error if it doesn't)
	if _, err := os.Stat(configPath); err == nil {
		if err := loadFromFile(configPath, &c); err != nil {
			return c, fmt.Errorf("failed to load config file %s: %w", configPath, err)
		}
	}

	// Override with environment variables (12-Factor App: config in environment)
	applyEnvironmentOverrides(&c)

	return c, nil
}

// loadFromFile loads configuration from a YAML file
func loadFromFile(path string, c *Config) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	if err := yaml.Unmarshal(b, c); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	return nil
}

// applyEnvironmentOverrides applies environment variables to config
// Environment variables follow the pattern: SECTION_SUBSECTION_KEY
// Examples:
//   - API_ADDR -> api.addr
//   - STATE_CONNECTION_STRING -> state.connectionString
//   - LOGGER_LEVEL -> logger.level
//   - OBS_TRACING_ENABLED -> obs.tracing.enabled
func applyEnvironmentOverrides(c *Config) {
	// API
	if v := os.Getenv("API_ADDR"); v != "" {
		c.API.Addr = v
	}
	if v := os.Getenv("API_BASE_URL"); v != "" {
		c.API.BaseURL = v
	}

	// State (Database)
	if v := os.Getenv("STATE_CONNECTION_STRING"); v != "" {
		c.State.ConnectionString = v
	}
	// Alternative naming for backward compatibility
	if v := os.Getenv("STATE_CONNECTIONSTRING"); v != "" {
		c.State.ConnectionString = v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		c.State.ConnectionString = v
	}

	// Queue
	if v := os.Getenv("QUEUE_WORKERS"); v != "" {
		if workers, err := strconv.Atoi(v); err == nil {
			c.Queue.Workers = workers
		}
	}

	// Artifacts
	if v := os.Getenv("ARTIFACTS_WORK_DIR"); v != "" {
		c.Artifacts.WorkDir = v
	}

	// LLM
	if v := os.Getenv("LLM_PROVIDER"); v != "" {
		c.LLM.Provider = v
	}
	if v := os.Getenv("LLM_OPENAI_API_KEY"); v != "" {
		c.LLM.OpenAI.APIKey = v
	}
	if v := os.Getenv("LLM_OPENAI_MODEL"); v != "" {
		c.LLM.OpenAI.Model = v
	}
	if v := os.Getenv("LLM_OLLAMA_BASE_URL"); v != "" {
		c.LLM.Ollama.BaseURL = v
	}
	if v := os.Getenv("LLM_OLLAMA_MODEL"); v != "" {
		c.LLM.Ollama.Model = v
	}

	// Auth
	if v := os.Getenv("AUTH_TOKEN"); v != "" {
		c.Auth.Token = v
	}

	// Logger
	if v := os.Getenv("LOGGER_LEVEL"); v != "" {
		c.Logger.Level = v
	}
	if v := os.Getenv("LOGGER_FORMAT"); v != "" {
		c.Logger.Format = v
	}
	if v := os.Getenv("LOGGER_OUTPUT"); v != "" {
		c.Logger.Output = v
	}
	if v := os.Getenv("LOGGER_REPORT_CALLER"); v != "" {
		c.Logger.ReportCaller = strings.ToLower(v) == "true"
	}

	// Observability - Tracing
	if v := os.Getenv("OBS_TRACING_ENABLED"); v != "" {
		c.Obs.Tracing.Enabled = strings.ToLower(v) == "true"
	}
	if v := os.Getenv("OBS_TRACING_ENDPOINT"); v != "" {
		c.Obs.Tracing.Endpoint = v
	}

	// Observability - Metrics
	if v := os.Getenv("OBS_METRICS_ENABLED"); v != "" {
		c.Obs.Metrics.Enabled = strings.ToLower(v) == "true"
	}
	if v := os.Getenv("OBS_METRICS_ENDPOINT"); v != "" {
		c.Obs.Metrics.Endpoint = v
	}
	if v := os.Getenv("OBS_METRICS_PROMETHEUS_ENABLED"); v != "" {
		c.Obs.Metrics.PrometheusEnabled = strings.ToLower(v) == "true"
	}
	// Alternative naming for backward compatibility
	if v := os.Getenv("OBS_METRICS_PROMETHEUSENABLED"); v != "" {
		c.Obs.Metrics.PrometheusEnabled = strings.ToLower(v) == "true"
	}
	if v := os.Getenv("OBS_METRICS_PROMETHEUS_PATH"); v != "" {
		c.Obs.Metrics.PrometheusPath = v
	}
	// Alternative naming for backward compatibility
	if v := os.Getenv("OBS_METRICS_PROMETHEUSPATH"); v != "" {
		c.Obs.Metrics.PrometheusPath = v
	}
}

func (c Config) Validate() error {
	if c.API.Addr == "" && c.API.BaseURL == "" {
		return errors.New("api.addr or api.baseURL must be set (or API_ADDR/API_BASE_URL environment variable)")
	}
	if c.State.ConnectionString == "" {
		return errors.New("state.connectionString must be set (or STATE_CONNECTION_STRING/DATABASE_URL environment variable)")
	}
	return nil
}
