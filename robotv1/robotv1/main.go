package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"robotv1/config"
	"robotv1/pkg/article"
	"robotv1/pkg/utils"
	"robotv1/pkg/wechat"
	"time"
)

// 全局日志记录器
var logger *log.Logger

func main() {
	// 初始化日志系统
	if err := initLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer closeLogger()

	logger.Println("======= 机器人服务启动 =======")

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		logger.Fatalf("加载配置失败: %v", err)
	}
	logger.Printf("当前配置: 文章源=%s, 本地文章路径=%s",
		config.AppConfig.ArticleSource, config.AppConfig.LocalArticlePath)

	// 确保文章目录存在
	if config.AppConfig.ArticleSource == "local" {
		if err := os.MkdirAll(config.AppConfig.LocalArticlePath, 0755); err != nil {
			logger.Printf("创建文章目录失败: %v", err)
		}
	}

	var content string
	var err error
	today := utils.TodayDateString()

	// 优先尝试获取本地文章
	logger.Println("尝试获取本地文章...")
	content, err = article.FetchLocalArticle(today)

	// 如果本地文章不存在或为空，尝试生成新文章
	if err != nil || content == "" {
		if os.IsNotExist(err) || content == "" {
			if content == "" {
				logger.Printf("本地文章为空: %s.md", today)
			} else {
				logger.Printf("本地文章不存在: %s.md", today)
			}

			// 检查是否配置了 DeepSeek API
			if config.AppConfig.DeepSeek.APIKey != "" {
				logger.Println("尝试使用 DeepSeek 生成文章...")
				content, err = article.GenerateArticleWithDeepSeek(logger, today)
				if err != nil {
					logger.Fatalf("生成文章失败: %v", err)
				}

				// 再次检查内容是否为空
				if content == "" {
					logger.Fatalf("生成的文章内容为空")
				}
			} else {
				logger.Fatalf("无法获取文章: 本地文章不存在或为空且未配置 DeepSeek API")
			}
		} else {
			logger.Fatalf("获取本地文章失败: %v", err)
		}
	} else {
		logger.Println("成功从本地获取文章")
	}

	logger.Printf("成功获取文章, 内容长度: %d 字符", len(content))

	// 发送到企业微信
	// 在main函数中替换发送消息的调用
	logger.Println("开始发送消息到企业微信...")
	if err := wechat.SendCardMessage(logger, content, today); err != nil {
		logger.Fatalf("消息发送失败: %v", err)
	}

	logger.Println("推送成功！")
	logger.Println("======= 机器人服务完成 =======")

}

// 启动文件服务器
func startFileServer() {
	// 确保文章目录存在
	articleDir := config.AppConfig.LocalArticlePath
	if articleDir == "" {
		articleDir = "articles"
	}
	if err := os.MkdirAll(articleDir, 0755); err != nil {
		logger.Printf("创建文章目录失败: %v", err)
	}

	// 设置文件服务器
	http.Handle("/articles/", http.StripPrefix("/articles/", http.FileServer(http.Dir(articleDir))))

	port := "8080"
	logger.Printf("启动文件服务器，端口: %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Printf("文件服务器启动失败: %v", err)
	}
}

// 初始化日志系统
func initLogger() error {
	// 确保日志目录存在
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 创建日志文件 (按日期命名)
	logPath := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %w", err)
	}

	// 创建多输出日志器 (同时输出到文件和控制台)
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(multiWriter, "ROBOT: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

// 关闭日志文件
func closeLogger() {
	if file, ok := logger.Writer().(*os.File); ok {
		file.Close()
	}
}
