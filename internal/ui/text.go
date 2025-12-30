package ui

import "strings"

func GetText(language, key string) string {
	isEn := strings.HasPrefix(strings.ToLower(language), "en")
	
	if isEn {
		switch key {
		case "options":
			return "Options: [y]Submit / [e]Edit / [r]Regenerate / [q]Cancel: "
		case "commit_success":
			return "Commit successful!"
		case "commit_fail":
			return "Commit failed: %v\n"
		case "enter_new_message":
			return "Enter new Commit Message: "
		case "canceled":
			return "Operation canceled."
		case "invalid_option":
			return "Invalid option, please try again."
		case "suggested_message":
			return "Suggested Commit Message:"
		case "using_provider":
			return "Using Provider: %s (Model: %s)\n"
		case "generating_message":
			return "Generating commit message..."
		case "config_load_error":
			return "Configuration load error: %v\n"
		case "provider_determine_error":
			return "Unable to determine provider: %v\n"
		case "generation_failed":
			return "Generation failed: %v\n"
		case "staging_area_empty":
			return "Staging area is empty."
		case "not_git_repo":
			return "Not a git repository."
		}
	} else {
		// Default Chinese
		switch key {
		case "options":
			return "选项: [y]提交 / [e]编辑 / [r]重新生成 / [q]取消: "
		case "commit_success":
			return "提交成功！"
		case "commit_fail":
			return "提交失败: %v\n"
		case "enter_new_message":
			return "请输入新的 Commit Message: "
		case "canceled":
			return "操作已取消。"
		case "invalid_option":
			return "无效的选项，请重试。"
		case "suggested_message":
			return "建议的 Commit Message:"
		case "using_provider":
			return "使用 Provider: %s (Model: %s)\n"
		case "generating_message":
			return "正在生成提交信息..."
		case "config_load_error":
			return "配置加载错误: %v\n"
		case "provider_determine_error":
			return "无法确定使用的 Provider: %v\n"
		case "generation_failed":
			return "生成失败: %v\n"
		case "staging_area_empty":
			return "暂存区为空。"
		case "not_git_repo":
			return "不是 git 仓库。"
		}
	}
	return key
}

