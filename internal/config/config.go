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
	SqlitePath string `yaml:"sqlitePath"`
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
	Endpoint string `yaml:"endpoint"` // OTLP endpoint (e.g., "localhost:4318"), empty or "none" to disable
	Enabled  bool   `yaml:"enabled"`  // Enable/disable metrics (default: false)
}

// Load loads from the fixed default path.
func Load() (Config, error) {
	return loadFrom(DefaultPath)
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
	return c, nil
}

func (c Config) Validate() error {
	if c.API.Addr == "" && c.API.BaseURL == "" {
		return errors.New("api.addr or api.baseURL must be set")
	}
	return nil
}
