package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetDiff 获取暂存区的 git diff
func GetDiff() (string, error) {
	// 检查是否在 git 仓库中
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("当前目录不是 git 仓库")
	}

	// 获取暂存区 diff
	cmd = exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取 git diff 失败: %v", err)
	}

	diff := strings.TrimSpace(string(output))
	if diff == "" {
		return "", fmt.Errorf("暂存区为空 (没有 staged 的文件)。请先 git add 文件")
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
		return fmt.Errorf("git commit 失败: %v\nOutput: %s", err, string(output))
	}
	fmt.Print(string(output))
	return nil
}

