package main

import (
	"errors"
	"fmt"
	"os"

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
		fmt.Fprintf(os.Stderr, ui.GetText(appConf.Language, "provider_determine_error"), err)
		os.Exit(1)
	}

	fmt.Printf(ui.GetText(appConf.Language, "using_provider"), pConfig.Name, pConfig.Model)

	// 2. 初始化 Provider
	// 目前主要是 OpenAI 兼容的，未来可以根据 pConfig.Name 做 switch
	aiProvider := provider.NewOpenAICompatibleProvider(*pConfig)

	// 3. 主循环
	for {
		// 获取 Diff
		diff, err := git.GetDiff()
		if err != nil {
			if errors.Is(err, git.ErrStagingEmpty) {
				fmt.Println(ui.GetText(appConf.Language, "staging_area_empty"))
				os.Exit(0)
			}
			if errors.Is(err, git.ErrNotGitRepo) {
				fmt.Println(ui.GetText(appConf.Language, "not_git_repo"))
				os.Exit(1)
			}
			fmt.Println(err)
			os.Exit(1)
		}

		// 生成 Message
		fmt.Println(ui.GetText(appConf.Language, "generating_message"))
		message, err := aiProvider.GenerateCommitMessage(diff)
		if err != nil {
			fmt.Fprintf(os.Stderr, ui.GetText(appConf.Language, "generation_failed"), err)
			os.Exit(1)
		}

		// 交互
		action := ui.InteractiveLoop(message, appConf.Language)
		if action != ui.ActionRegenerate {
			break
		}
	}
}
