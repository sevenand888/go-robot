package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"robotv1/config"
	"strings"
	"time"
)

// 企业微信卡片消息结构
type CardMessage struct {
	MsgType string `json:"msgtype"`
	News    News   `json:"news"`
}

type News struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PicURL      string `json:"picurl"`
}

// SendCardMessage 发送企业微信卡片消息（三个链接）
func SendCardMessage(logger *log.Logger, content string, articleDate string) error {
	// 从内容中提取标题
	title := extractTitle(content)
	if title == "" {
		title = "AI每日文章推送 - " + articleDate
	}

	// 创建文章摘要
	summary := createSummary(content)

	// 创建三个文章项
	articles := []Article{
		{
			Title:       title,
			Description: summary,
			#这里是指定存放文章连接的服务器网页地址，端口可以自己指定或者不加端口都是可以的
			URL:         "http://服务器ip地址:8080/a.php", // 固定链接
			PicURL:      "https://img.icons8.com/color/96/000000/artificial-intelligence.png",
		},
		{
			Title:       "📖 下载全文",
			Description: "点击查看完整文章内容",
			URL:         fmt.Sprintf("http://服务器ip地址:8080/articles/%s.md", articleDate),
		},
		{
			Title:       "💬 反馈意见",
			Description: "点击提供反馈建议",
			URL:         "https://我的另一个服务器页面，使用wp的表单服务专门负责收集信息",
		},
	}

	// 创建消息体
	message := CardMessage{
		MsgType: "news",
		News: News{
			Articles: articles,
		},
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		if logger != nil {
			logger.Printf("JSON序列化失败: %v", err)
		}
		return fmt.Errorf("JSON序列化失败: %w", err)
	}

	if logger != nil {
		logger.Printf("发送企业微信卡片消息, 标题: %s", title)
		logger.Printf("消息包含 %d 个链接", len(articles))
	}

	// 发送POST请求
	resp, err := http.Post(
		config.AppConfig.WebhookURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		if logger != nil {
			logger.Printf("请求失败: %v", err)
		}
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if logger != nil {
			logger.Printf("企业微信返回错误状态: %d", resp.StatusCode)
			logger.Printf("错误响应: %s", string(body))
		}
		return fmt.Errorf("企业微信返回错误状态: %d", resp.StatusCode)
	}

	if logger != nil {
		logger.Println("消息发送成功")
	}
	return nil
}

// 从内容中提取标题
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}

	// 备选创意标题
	creativeTitles := []string{
		"🤖 AI提示词设计秘籍：让机器懂你心",
		"🚀 提示词工程：与AI对话的艺术",
		"💡 这样设计提示词，效果翻倍！",
		"🎯 精准提示：解锁AI潜能的金钥匙",
		"🧠 大脑与机器的完美对话指南",
	}

	// 使用随机种子确保标题多样性
	rand.Seed(time.Now().UnixNano())
	return creativeTitles[rand.Intn(len(creativeTitles))]
}

// 创建摘要
func createSummary(content string) string {
	// 取前三行非空行作为摘要
	lines := strings.Split(content, "\n")
	var summaryLines []string
	count := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			summaryLines = append(summaryLines, trimmed)
			count++
			if count >= 3 {
				break
			}
		}
	}

	return strings.Join(summaryLines, "\n")
}
