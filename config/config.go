package config

import (
	"os"
	"strconv"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config - Application configuration
type Config struct {
	Log    string `koanf:"log"`
	Debug  bool   `koanf:"debug"`
	Gemini struct {
		APIKey         string `koanf:"api_key"`
		ModelName      string `koanf:"model_name"`
		MaxTokens      int    `koanf:"max_tokens"`
		QueryTemplate  string `koanf:"query_template"`
		ThinkingLevel  string `koanf:"thinking_level"`
		ThinkingBudget *int   `koanf:"thinking_budget"`
	} `koanf:"gemini"`
	HTTP struct {
		Port             int      `koanf:"port"`
		EndpointPath     string   `koanf:"endpoint_path"`
		AuthToken        string   `koanf:"auth_token"`
		AllowedOrigins   []string `koanf:"allowed_origins"`
		HeartbeatSeconds int      `koanf:"heartbeat_seconds"`
	} `koanf:"http"`
}

func defaultValues() map[string]any {
	return map[string]any{
		"log":                    "",
		"debug":                  false,
		"gemini.model_name":      "gemini-3.1-pro-preview",
		"gemini.max_tokens":      5000,
		"gemini.thinking_level":  "",
		"http.port":              8080,
		"http.endpoint_path":     "/mcp",
		"http.heartbeat_seconds": 30,
	}
}

func envOverrides() map[string]any {
	m := map[string]any{}
	if v := os.Getenv("LOG_PATH"); v != "" {
		m["log"] = v
	}
	if v := os.Getenv("DEBUG"); v != "" {
		m["debug"] = v == "true" || v == "1"
	}
	if v := os.Getenv("GEMINI_API_KEY"); v != "" {
		m["gemini.api_key"] = v
	}
	if v := os.Getenv("GEMINI_MODEL_NAME"); v != "" {
		m["gemini.model_name"] = v
	}
	if v := os.Getenv("GEMINI_MAX_TOKENS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			m["gemini.max_tokens"] = n
		}
	}
	if v := os.Getenv("GEMINI_QUERY_TEMPLATE"); v != "" {
		m["gemini.query_template"] = v
	}
	if v := os.Getenv("GEMINI_THINKING_LEVEL"); v != "" {
		m["gemini.thinking_level"] = v
	}
	return m
}

// LoadConfig - Load configuration file
func LoadConfig(path string) (*Config, error) {
	k := koanf.New(".")

	k.Load(confmap.Provider(defaultValues(), "."), nil)

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, err
	}

	if overrides := envOverrides(); len(overrides) > 0 {
		k.Load(confmap.Provider(overrides, "."), nil)
	}

	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	// Handle ThinkingBudget separately since koanf doesn't support *int natively
	if v := os.Getenv("GEMINI_THINKING_BUDGET"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Gemini.ThinkingBudget = &n
		}
	} else if k.Exists("gemini.thinking_budget") {
		n := k.Int("gemini.thinking_budget")
		cfg.Gemini.ThinkingBudget = &n
	}

	return cfg, nil
}
