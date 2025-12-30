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
func InteractiveLoop(initialMessage string, language string) Action {
	message := initialMessage
	scanner := bufio.NewScanner(os.Stdin)

	// 简单的多语言支持
	getText := func(key string) string {
		return GetText(language, key)
	}

	for {
		printMessage(message, getText("suggested_message"))

		fmt.Print(getText("options"))
		if !scanner.Scan() {
			return ActionQuit
		}
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))

		switch choice {
		case "y":
			if err := git.Commit(message); err != nil {
				fmt.Printf(getText("commit_fail"), err)
				// 失败后继续循环，允许用户重试或编辑
			} else {
				fmt.Println(getText("commit_success"))
				return ActionCommit
			}

		case "e":
			fmt.Print(getText("enter_new_message"))
			if scanner.Scan() {
				newMessage := strings.TrimSpace(scanner.Text())
				if newMessage != "" {
					message = newMessage
				}
			}

		case "r":
			return ActionRegenerate

		case "q":
			fmt.Println(getText("canceled"))
			return ActionQuit

		default:
			fmt.Println(getText("invalid_option"))
		}
	}
}

func printMessage(message, title string) {
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println(title)
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println(message)
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println(strings.Repeat("=", 40) + "\n")
}
