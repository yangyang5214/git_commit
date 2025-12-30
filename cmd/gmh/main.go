package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yangyang5214/git_commit/internal/config"
	"github.com/yangyang5214/git_commit/internal/git"
	"github.com/yangyang5214/git_commit/internal/provider"
	"github.com/yangyang5214/git_commit/internal/ui"
)

func main() {
	// 1. 加载配置
	appConf, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "配置加载错误: %v\n", err)
		os.Exit(1)
	}

	pConfig, err := appConf.GetCurrentProviderConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法确定使用的 Provider: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("使用 Provider: %s (Model: %s)\n", pConfig.Name, pConfig.Model)

	// 2. 初始化 Provider
	// 目前主要是 OpenAI 兼容的，未来可以根据 pConfig.Name 做 switch
	aiProvider := provider.NewOpenAICompatibleProvider(*pConfig)

	// 3. 主循环
	for {
		// 获取 Diff
		diff, err := git.GetDiff()
		if err != nil {
			fmt.Println(err)
			if strings.Contains(err.Error(), "暂存区为空") {
				os.Exit(0)
			}
			os.Exit(1)
		}

		// 生成 Message
		fmt.Println("正在生成提交信息...")
		message, err := aiProvider.GenerateCommitMessage(diff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
			os.Exit(1)
		}

		// 交互
		action := ui.InteractiveLoop(message)
		if action != ui.ActionRegenerate {
			break
		}
	}
}
