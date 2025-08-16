package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"robotv1/config"
)

// CreateTodayArticle 创建今日文章（如果需要）
func CreateTodayArticle(logger *log.Logger) error {
	date := TodayDateString()
	filePath := filepath.Join(config.AppConfig.LocalArticlePath, date+".md")

	// 如果文件已存在则跳过
	if _, err := os.Stat(filePath); err == nil {
		if logger != nil {
			logger.Printf("今日文章已存在: %s", filePath)
		}
		return nil
	}

	if logger != nil {
		logger.Printf("创建今日文章: %s", filePath)
	}

	// 从模板复制
	templatePath := filepath.Join(config.AppConfig.LocalArticlePath, "template.md")

	// 检查模板文件是否存在
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		if logger != nil {
			logger.Printf("模板文件不存在: %s", templatePath)
		}
		// 创建空文件
		if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
			if logger != nil {
				logger.Printf("创建空文章失败: %v", err)
			}
			return err
		}
		return nil
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		if logger != nil {
			logger.Printf("读取模板失败: %v", err)
		}
		return fmt.Errorf("读取模板失败: %w", err)
	}

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		if logger != nil {
			logger.Printf("写入文章失败: %v", err)
		}
		return err
	}

	if logger != nil {
		logger.Printf("今日文章创建成功: %s", filePath)
	}
	return nil
}
