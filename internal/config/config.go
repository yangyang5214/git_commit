package config

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yangyang5214/git_commit/internal/provider"
)

//go:embed default_config.ini
var demoConfigContent string

type AppConfig struct {
	CurrentProvider string
	Providers       map[string]*provider.ProviderConfig
}

func LoadConfig() (*AppConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("无法获取用户主目录: %v", err)
	}

	configPath := filepath.Join(homeDir, ".git_commit")
	appConfig := &AppConfig{
		Providers: make(map[string]*provider.ProviderConfig),
	}

	// 初始化默认 provider 容器
	appConfig.Providers["default"] = &provider.ProviderConfig{Name: "default"}

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return createDemoConfig(configPath)
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.Contains(line, "=") {
			// 兼容旧格式：单行 key
			if appConfig.Providers["default"].APIKey == "" {
				appConfig.Providers["default"].APIKey = line
			}
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		fullKey := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		if fullKey == "current_provider" || fullKey == "provider" {
			appConfig.CurrentProvider = value
			continue
		}

		// 解析 provider.key
		var providerName, key string
		if strings.Contains(fullKey, ".") {
			keyParts := strings.SplitN(fullKey, ".", 2)
			providerName = keyParts[0]
			key = keyParts[1]
		} else {
			providerName = "default"
			key = fullKey
		}

		// 确保 ProviderConfig 存在
		if _, ok := appConfig.Providers[providerName]; !ok {
			appConfig.Providers[providerName] = &provider.ProviderConfig{Name: providerName}
		}
		pConfig := appConfig.Providers[providerName]

		switch key {
		case "api_key":
			pConfig.APIKey = value
		case "base_url":
			pConfig.BaseURL = value
		case "model":
			pConfig.Model = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 环境变量覆盖 (仅针对当前/默认 provider)
	// 如果用户设置了 OPENAI_API_BASE，我们不知道他想覆盖哪个 provider，
	// 假设是覆盖 default 或者 current。
	// 这里为了简单，如果设置了 env，覆盖 default。
	if envBaseURL := os.Getenv("OPENAI_API_BASE"); envBaseURL != "" {
		if appConfig.Providers["default"] != nil {
			appConfig.Providers["default"].BaseURL = envBaseURL
		}
	}

	return appConfig, nil
}

func createDemoConfig(path string) (*AppConfig, error) {
	if err := os.WriteFile(path, []byte(demoConfigContent), 0644); err != nil {
		return nil, fmt.Errorf("创建示例配置文件失败: %v", err)
	}
	return nil, fmt.Errorf("配置文件不存在，已自动创建示例文件于 %s，请填写配置后重试", path)
}

// GetCurrentProviderConfig 获取当前激活的 Provider 配置
func (c *AppConfig) GetCurrentProviderConfig() (*provider.ProviderConfig, error) {
	// 1. 如果指定了 CurrentProvider
	if c.CurrentProvider != "" {
		if p, ok := c.Providers[c.CurrentProvider]; ok {
			// 检查必填项
			if p.APIKey == "" {
				return nil, fmt.Errorf("Provider '%s' 缺少 api_key", c.CurrentProvider)
			}
			return p, nil
		}
		return nil, fmt.Errorf("未找到 Provider '%s' 的配置", c.CurrentProvider)
	}

	// 2. 如果没指定，尝试使用 default
	if p, ok := c.Providers["default"]; ok && p.APIKey != "" {
		// 填充默认值
		if p.BaseURL == "" {
			p.BaseURL = "https://api.openai.com/v1/chat/completions"
		}
		if p.Model == "" {
			p.Model = "gpt-3.5-turbo"
		}
		return p, nil
	}

	// 3. 都没有，找第一个配置了 API Key 的 Provider
	for _, p := range c.Providers {
		if p.APIKey != "" {
			return p, nil
		}
	}

	return nil, fmt.Errorf("未找到有效的 Provider 配置，请检查配置文件")
}
