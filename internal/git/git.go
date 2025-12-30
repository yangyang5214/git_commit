package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	ErrNotGitRepo   = errors.New("not a git repository")
	ErrStagingEmpty = errors.New("staging area is empty")
)

// GetDiff 获取暂存区的 git diff
func GetDiff() (string, error) {
	// 检查是否在 git 仓库中
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return "", ErrNotGitRepo
	}

	// 获取暂存区 diff
	cmd = exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %v", err)
	}

	diff := strings.TrimSpace(string(output))
	if diff == "" {
		return "", ErrStagingEmpty
	}

	return diff, nil
}

// Commit 提交代码
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	// 将 stdout/stderr 连接到当前进程，以便用户看到 git 的输出
	// 或者我们可以捕获输出并返回
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %v\nOutput: %s", err, string(output))
	}
	fmt.Print(string(output))
	return nil
}

