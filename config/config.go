package config

import (
	"github.com/jinzhu/configor"
)

// Config - Application configuration
type Config struct {
	Log    string `yaml:"log" default:"" env:"LOG_PATH"`
	Debug  bool   `yaml:"debug" default:"false" env:"DEBUG"`
	Gemini struct {
		APIKey        string `yaml:"api_key" env:"GEMINI_API_KEY" required:"true"`
		ModelName     string `yaml:"model_name" default:"gemini-2.5-flash" env:"GEMINI_MODEL_NAME"`
		MaxTokens     int    `yaml:"max_tokens" default:"5000" env:"GEMINI_MAX_TOKENS"`
		QueryTemplate string `yaml:"query_template" env:"GEMINI_QUERY_TEMPLATE"`
	} `yaml:"gemini"`
}

// LoadConfig - Load configuration file
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := configor.New(&configor.Config{
		Debug:      false,
		Verbose:    false,
		Silent:     true,
		AutoReload: false,
	}).Load(cfg, path)
	return cfg, err
}
