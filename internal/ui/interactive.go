package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yangyang5214/git_commit/internal/git"
)

type Action string

const (
	ActionCommit     Action = "COMMIT"
	ActionEdit       Action = "EDIT"
	ActionRegenerate Action = "REGENERATE"
	ActionQuit       Action = "QUIT"
)

// InteractiveLoop 进入交互式循环
func InteractiveLoop(initialMessage string) Action {
	message := initialMessage
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMessage(message)

		fmt.Print("选项: [y]提交 / [e]编辑 / [r]重新生成 / [q]取消: ")
		if !scanner.Scan() {
			return ActionQuit
		}
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))

		switch choice {
		case "y":
			if err := git.Commit(message); err != nil {
				fmt.Printf("提交失败: %v\n", err)
				// 失败后继续循环，允许用户重试或编辑
			} else {
				fmt.Println("提交成功！")
				return ActionCommit
			}

		case "e":
			fmt.Print("请输入新的 Commit Message: ")
			if scanner.Scan() {
				newMessage := strings.TrimSpace(scanner.Text())
				if newMessage != "" {
					message = newMessage
				}
			}

		case "r":
			return ActionRegenerate

		case "q":
			fmt.Println("操作已取消。")
			return ActionQuit

		default:
			fmt.Println("无效的选项，请重试。")
		}
	}
}

func printMessage(message string) {
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("建议的 Commit Message:")
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println(message)
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println(strings.Repeat("=", 40) + "\n")
}

