package article

import (
	"fmt"
	"os"
	"path/filepath"
	"robotv1/config"
)

// FetchLocalArticle 获取本地文章
func FetchLocalArticle(date string) (string, error) {
	filePath := filepath.Join(config.AppConfig.LocalArticlePath, date+".md")

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("今日文章不存在: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取本地文章失败: %w", err)
	}

	return string(content), nil
}
