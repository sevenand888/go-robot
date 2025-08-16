// config/config.go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WebhookURL       string `yaml:"webhook_url"`
	ArticleSource    string `yaml:"article_source"`
	LocalArticlePath string `yaml:"local_article_path"`

	// DeepSeek配置
	DeepSeek struct {
		APIKey      string  `yaml:"api_key"`
		Model       string  `yaml:"model"`
		APIBase     string  `yaml:"api_base"`
		Temperature float64 `yaml:"temperature"`
		MaxTokens   int     `yaml:"max_tokens"`
	} `yaml:"deepseek"`
}

var AppConfig Config

func LoadConfig() error {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}
