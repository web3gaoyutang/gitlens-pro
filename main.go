package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	insertCode = `e={user:{id:"88888888-8888-8888-8888-888888888888",name:"Neo",email:"x@x.com",status:"activated",createdDate:"2000-01-01T00:00:00.000Z"},licenses:{paidLicenses:{},effectiveLicenses:{"gitlens-pro":{organizationId:"Linux",latestStatus:"active",latestStartDate:"2024-01-01",latestEndDate:"2999-01-01",reactivationCount:99,nextOptInDate:"2999-01-01"}}},nextOptInDate:"2999-01-01"};`
)

var defaultExtDir = ".vscode"

func getExtensionsDir() string {
	// 优先使用命令行参数
	if len(os.Args) > 2 && os.Args[1] == "--ext-dir" {
		defaultExtDir = os.Args[2]
	}

	// 其次使用环境变量
	if envDir := os.Getenv("VSCODE_EXTENSIONS_DIR"); envDir != "" {
		return envDir
	}

	// 最后使用默认值
	home, err := os.UserHomeDir()

	fmt.Println(home, defaultExtDir)
	if err != nil {
		return defaultExtDir
	}
	return filepath.Join(home, defaultExtDir)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "restore" {
		if err := Restore(); err != nil {
			fmt.Printf("恢复失败: %v\n", err)
			return
		}
		fmt.Println("恢复完成! 请重启 VS Code。")
		return
	}
	// 获取扩展目录
	extensionsDir := filepath.Join(getExtensionsDir(), "extensions")

	fmt.Println(extensionsDir)

	// 获取最新版本的 GitLens 目录
	extensionPath, err := getLatestGitLensPath(extensionsDir)
	if err != nil {
		fmt.Printf("查找 GitLens 扩展失败: %v\n", err)
		return
	}

	// 需要修改的文件列表
	filesToModify := []string{
		filepath.Join(extensionPath, "dist", "gitlens.js"),
		filepath.Join(extensionPath, "dist", "browser", "gitlens.js"),
	}

	// 处理每个文件
	for _, file := range filesToModify {
		if err := processFile(file); err != nil {
			fmt.Printf("处理文件 %s 失败: %v\n", file, err)
			continue
		}
	}

	fmt.Println("\n激活完成! 请重启 VS Code 以使更改生效。")
}

func getLatestGitLensPath(extensionsDir string) (string, error) {
	entries, err := os.ReadDir(extensionsDir)
	if err != nil {
		return "", fmt.Errorf("读取扩展目录失败: %v", err)
	}

	var gitLensDirs []string
	pattern := regexp.MustCompile(`^eamodio\.gitlens-\d+\.\d+\.\d+$`)

	for _, entry := range entries {
		if entry.IsDir() && pattern.MatchString(entry.Name()) {
			gitLensDirs = append(gitLensDirs, entry.Name())
		}
	}

	if len(gitLensDirs) == 0 {
		return "", fmt.Errorf("未找到 GitLens 扩展")
	}

	// 按版本号排序
	sort.Sort(sort.Reverse(sort.StringSlice(gitLensDirs)))
	return filepath.Join(extensionsDir, gitLensDirs[0]), nil
}

// 添加版本检测函数
func isVersion15(dirName string) bool {
	pattern := regexp.MustCompile(`^eamodio\.gitlens-15\.\d+\.\d+$`)
	return pattern.MatchString(dirName)
}

// 修改 processFile 函数
func processFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 创建备份
	backupPath := filePath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		if err := os.WriteFile(backupPath, content, 0644); err != nil {
			return fmt.Errorf("创建备份失败: %v", err)
		}
		fmt.Printf("已创建备份: %s\n", backupPath)
	}

	// 获取目录名以检查版本
	dirName := filepath.Base(filepath.Dir(filepath.Dir(filePath)))
	if isVersion15(dirName) {
		return processVersion15File(filePath, content)
	}
	return processVersion16File(filePath, content)
}

// 处理 15.x 版本的文件
func processVersion15File(filePath string, content []byte) error {
	contentStr := string(content)

	// 替换模式
	replacements := map[string]string{
		"qn.CommunityWithAccount": "qn.Enterprise",
		"qn.Community":            "qn.Enterprise",
		"qn.Pro":                  "qn.Enterprise",
	}

	modified := false
	for old, new := range replacements {
		if strings.Contains(contentStr, old) {
			contentStr = strings.ReplaceAll(contentStr, old, new)
			modified = true
			fmt.Printf("替换 %s 为 %s\n", old, new)
		}
	}

	if !modified {
		return fmt.Errorf("未找到需要替换的内容")
	}

	// 写入修改后的内容
	if err := os.WriteFile(filePath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("成功修改文件: %s\n", filePath)
	return nil
}

// 处理 16.x 版本的文件
func processVersion16File(filePath string, content []byte) error {
	// 查找匹配模式
	pattern := regexp.MustCompile(`let ([a-zA-Z])={id:e\.user\.id,name:`)
	matches := pattern.FindStringSubmatch(string(content))

	if len(matches) < 2 {
		return fmt.Errorf("未找到匹配模式")
	}

	matchedLetter := matches[1]
	exactMatch := fmt.Sprintf("let %s={id:e.user.id,name:", matchedLetter)

	// 注入激活代码
	newContent := strings.Replace(string(content), exactMatch, insertCode+exactMatch, 1)

	// 写入修改后的内容
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("成功修改文件: %s\n", filePath)
	fmt.Printf("匹配的变量: %s\n", matchedLetter)

	return nil
}
