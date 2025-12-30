# AI Git Commit (gmh)

[中文](./README.md) | [English](./README_EN.md)

A command-line tool that uses AI (DeepSeek, OpenAI, Gemini, etc.) to automatically generate Git Commit Messages, named `gmh`.

## Features

- **Multi-model Support**: Perfect support for DeepSeek V3, OpenAI GPT-4, Google Gemini, and more.
- **Smart Detection**: Automatically reads changes from the staging area (`git diff --cached`).
- **Interactive**: Allows users to preview, edit, or regenerate the message before committing.
- **Flexible Configuration**: Supports switching between multiple provider configurations.

## Installation

### Method 1: Use go install (Recommended)

If you have a Go environment installed (1.25+), you can install directly from GitHub:

```bash
go install github.com/yangyang5214/git_commit/cmd/gmh@latest
```

After installation, simply type `gmh` in your terminal to use it.

### Method 2: Local Installation

If you have downloaded the source code:

```bash
# Run in the project root directory
go install ./cmd/gmh
```

### Method 3: Manual Compilation

```bash
# 1. Compile
go build -o gmh ./cmd/gmh

# 2. Move to system path (Optional)
mv gmh /usr/local/bin/
```

## Configuration Guide

Configuration file path: `~/.git_commit`

Supports configuring multiple AI providers and quickly switching via `current_provider`.

### Example Configuration (`~/.git_commit`)

```properties
# ========================
# Global Settings
# ========================
# Specify the current provider (corresponds to the prefix in the configuration blocks below)
current_provider=deepseek

# ========================
# DeepSeek (Recommended)
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
# Note: Gemini requires an OpenAI-compatible Endpoint or use an adapter
gemini.api_key=AIzaSy...
gemini.base_url=https://generativelanguage.googleapis.com/v1beta/openai/
gemini.model=gemini-1.5-flash

# ========================
# Ollama (Local Model)
# ========================
ollama.api_key=ollama # Can be anything
ollama.base_url=http://localhost:11434/v1/chat/completions
ollama.model=llama3
```


## Usage

1.  **Stage Code**:
    ```bash
    git add .
    ```

2.  **Run Tool**:
    ```bash
    gmh
    ```

3.  **Interactive Operations**:
    - `y`: Confirm commit.
    - `e`: Edit the generated message.
    - `r`: Not satisfied? Let AI regenerate.
    - `q`: Quit.

