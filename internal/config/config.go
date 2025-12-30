package config

import (
	"errors"
	"os"

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
	OpenAI   OpenAIConfig `yaml:"openai"`
	Ollama   OllamaConfig `yaml:"ollama"`
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

// Load loads configuration from file.
// If CONFIG_PATH environment variable is set, it uses that path.
// Otherwise, it falls back to the default path (configs/config.yaml).
func Load() (Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = DefaultPath
	}
	return loadFrom(configPath)
}

func loadFrom(path string) (Config, error) {
	var c Config
	b, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return c, err
	}
	
	// Override with environment variables if set
	if envAddr := os.Getenv("API_ADDR"); envAddr != "" {
		c.API.Addr = envAddr
	}
	if envConnStr := os.Getenv("STATE_CONNECTIONSTRING"); envConnStr != "" {
		c.State.ConnectionString = envConnStr
	}
	if envLogLevel := os.Getenv("LOGGER_LEVEL"); envLogLevel != "" {
		c.Logger.Level = envLogLevel
	}
	if envLogFormat := os.Getenv("LOGGER_FORMAT"); envLogFormat != "" {
		c.Logger.Format = envLogFormat
	}
	if envTracingEnabled := os.Getenv("OBS_TRACING_ENABLED"); envTracingEnabled != "" {
		c.Obs.Tracing.Enabled = envTracingEnabled == "true"
	}
	if envTracingEndpoint := os.Getenv("OBS_TRACING_ENDPOINT"); envTracingEndpoint != "" {
		c.Obs.Tracing.Endpoint = envTracingEndpoint
	}
	if envMetricsEnabled := os.Getenv("OBS_METRICS_ENABLED"); envMetricsEnabled != "" {
		c.Obs.Metrics.Enabled = envMetricsEnabled == "true"
	}
	if envMetricsEndpoint := os.Getenv("OBS_METRICS_ENDPOINT"); envMetricsEndpoint != "" {
		c.Obs.Metrics.Endpoint = envMetricsEndpoint
	}
	if envPrometheusEnabled := os.Getenv("OBS_METRICS_PROMETHEUSENABLED"); envPrometheusEnabled != "" {
		c.Obs.Metrics.PrometheusEnabled = envPrometheusEnabled == "true"
	}
	if envPrometheusPath := os.Getenv("OBS_METRICS_PROMETHEUSPATH"); envPrometheusPath != "" {
		c.Obs.Metrics.PrometheusPath = envPrometheusPath
	}
	
	return c, nil
}

func (c Config) Validate() error {
	if c.API.Addr == "" && c.API.BaseURL == "" {
		return errors.New("api.addr or api.baseURL must be set")
	}
	return nil
}
