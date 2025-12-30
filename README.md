# AI Git Commit (gmh)

一个使用 AI (DeepSeek, OpenAI, Gemini 等) 自动生成 Git Commit Message 的命令行工具，命令名为 `gmh`。

## 特性

- **多模型支持**: 完美支持 DeepSeek V3, OpenAI GPT-4, Google Gemini 等多种模型。
- **智能检测**: 自动读取暂存区 (`git diff --cached`) 的变更。
- **交互式**: 允许用户在提交前预览、编辑或重新生成 Message。
- **配置灵活**: 支持多 Provider 配置切换。

## 安装

### 方式 1: 使用 go install (推荐)

如果你已安装 Go 环境 (1.25+)，可以直接从 GitHub 安装：

```bash
go install github.com/yangyang5214/git_commit/cmd/gmh@latest
```

安装完成后，直接在终端输入 `gmh` 即可使用。

### 方式 2: 本地安装

如果你下载了源码：

```bash
# 在项目根目录下运行
go install ./cmd/gmh
```

### 方式 3: 手动编译

```bash
# 1. 编译
go build -o gmh ./cmd/gmh

# 2. 移动到系统路径 (可选)
mv gmh /usr/local/bin/
```

## 配置指南

配置文件路径: `~/.git_commit`

支持配置多个 AI 提供商，并通过 `current_provider` 快速切换。

### 示例配置 (`~/.git_commit`)

```properties
# ========================
# 全局设置
# ========================
# 指定当前使用的提供商 (对应下面配置块的前缀)
current_provider=deepseek

# ========================
# DeepSeek (推荐)
# ========================
deepseek.api_key=sk-xxxxxxxxxxxxxxxxxxxxxxxx
deepseek.base_url=https://api.deepseek.com
deepseek.model=deepseek-chat

# ========================
# OpenAI
# ========================
openai.api_key=sk-yyyyyyyyyyyyyyyyyyyyyyyy
openai.base_url=https://api.openai.com/v1/chat/completions
openai.model=gpt-4

# ========================
# Google Gemini
# ========================
# 注意: Gemini 需要使用兼容 OpenAI 的 Endpoint 或者通过适配器
gemini.api_key=AIzaSy...
gemini.base_url=https://generativelanguage.googleapis.com/v1beta/openai/
gemini.model=gemini-1.5-flash

# ========================
# Ollama (本地模型)
# ========================
ollama.api_key=ollama # 随便填
ollama.base_url=http://localhost:11434/v1/chat/completions
ollama.model=llama3
```

### 环境变量覆盖

你也可以临时通过环境变量覆盖当前使用的 Base URL (通常用于代理)：

```bash
export OPENAI_API_BASE="https://proxy.example.com/v1/chat/completions"
```

## 使用方法

1.  **Stage 代码**:
    ```bash
    git add .
    ```

2.  **运行工具**:
    ```bash
    gmh
    ```

3.  **交互操作**:
    - `y`: 确认提交。
    - `e`: 编辑生成的消息。
    - `r`: 不满意？让 AI 重新生成。
    - `q`: 退出。
