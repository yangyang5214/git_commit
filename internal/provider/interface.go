package provider

// ProviderConfig 定义了 Provider 所需的通用配置
type ProviderConfig struct {
	Name     string
	APIKey   string
	BaseURL  string
	Model    string
	Language string // zh, en
}

// Provider 定义了 AI 服务的通用接口
type Provider interface {
	// GenerateCommitMessage 根据 git diff 生成提交信息
	GenerateCommitMessage(diff string) (string, error)
}
