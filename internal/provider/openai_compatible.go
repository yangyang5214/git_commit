package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAICompatibleProvider 实现了 OpenAI API 格式的 Provider
// 适用于 OpenAI, DeepSeek, Moonshot, 以及其他兼容 OpenAI 接口的模型
type OpenAICompatibleProvider struct {
	config ProviderConfig
}

func NewOpenAICompatibleProvider(config ProviderConfig) *OpenAICompatibleProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1/chat/completions"
	}
	// 确保 BaseURL 指向 chat/completions (如果用户只给了 host)
	// 这里做一个简单的处理，如果用户给了 https://api.deepseek.com 而没给路径，
	// 不同的厂商路径可能不同，所以最好还是要在配置里写全，或者这里不做过顶推断，
	// 依然假设 BaseURL 是完整的 endpoint 或者基础 URL。
	// 为了兼容性，我们建议配置完整的 Endpoint URL，或者我们在代码里拼接。
	// 通常 OpenAI SDK 是 base_url + /chat/completions。
	// 但为了灵活性，我们让 BaseURL 就是完整的 URL (Endpoint)。
	// 如果用户配置的是 "https://api.openai.com/v1"，我们需要拼接。
	if !strings.HasSuffix(config.BaseURL, "/chat/completions") && !strings.Contains(config.BaseURL, "googleapis") {
		// googleapis (Gemini) 的路径比较特殊，不做处理
		// 其他情况尝试智能拼接，但最稳妥是让 BaseURL 等于 Endpoint
		config.BaseURL = strings.TrimRight(config.BaseURL, "/") + "/chat/completions"
	}

	return &OpenAICompatibleProvider{
		config: config,
	}
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	Stream      bool            `json:"stream"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string      `json:"message"`
		Code    interface{} `json:"code"`
	} `json:"error,omitempty"`
}

func (p *OpenAICompatibleProvider) GenerateCommitMessage(diff string) (string, error) {
	// 截断 diff 防止 token 溢出，根据模型不同限制可能不同，这里取一个保守值
	if len(diff) > 12000 {
		diff = diff[:12000] + "\n...(truncated)"
	}

	var systemPrompt string
	if strings.HasPrefix(p.config.Language, "en") {
		systemPrompt = "You are a professional developer tool. Based on the following git diff content, " +
			"generate a concise and standardized git commit message. " +
			"Format requirements: The first line should be a short summary (within 50 characters). " +
			"If necessary, leave a blank line followed by a detailed description. " +
			"The response should contain ONLY the commit message itself, without any explanation or other text. " +
			"Please answer in English."
	} else {
		systemPrompt = "你是一个专业的开发者工具。请根据以下的 git diff 内容，" +
			"生成一个简洁、规范的 git commit message。" +
			"格式要求：第一行是简短的摘要（50字符以内），" +
			"如果需要，空一行后可以跟详细的描述。" +
			"回答中只包含 commit message 本身，不要包含任何解释或其他文字。" +
			"请使用中文回答。"
	}

	reqBody := openAIRequest{
		Model: p.config.Model,
		Messages: []openAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: "Diff content:\n" + diff},
		},
		Temperature: 0.7,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON 编码失败: %v", err)
	}

	req, err := http.NewRequest("POST", p.config.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API 请求失败 (状态码 %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(openAIResp.Choices) == 0 {
		if openAIResp.Error.Message != "" {
			return "", fmt.Errorf("API 返回错误: %s", openAIResp.Error.Message)
		}
		return "", fmt.Errorf("API 返回结果为空")
	}

	content := strings.TrimSpace(openAIResp.Choices[0].Message.Content)

	// 清理 Markdown 代码块
	content = strings.TrimPrefix(content, "```git commit")
	content = strings.TrimPrefix(content, "```commit")
	content = strings.TrimPrefix(content, "```text")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")

	return strings.TrimSpace(content), nil
}
