package article

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"robotv1/config"
	"strings"
	"time"
)

// 定义模型返回的JSON结构
type DeepSeekResponse struct {
	Content string `json:"content"`
}

// GenerateArticleWithDeepSeek 使用DeepSeek生成文章
func GenerateArticleWithDeepSeek(logger *log.Logger, date string) (string, error) {
	cfg := config.AppConfig.DeepSeek

	// 确保API基础地址正确
	apiBase := cfg.APIBase
	if apiBase == "" {
		apiBase = "https://api.deepseek.com/v1" // 修正URL格式
	}

	if logger != nil {
		logger.Printf("🚀 开始生成文章: 模型=%s, 温度=%.1f", cfg.Model, cfg.Temperature)
	}

	// 关键修改：提示词中的换行符必须用\\n转义
	// 修改提示词增加多样性
	requestBody := map[string]interface{}{
		"model": cfg.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": `现在是2025年，你是一位精通爆款文章写作的AI专家，特别擅长撰写吸引眼球的AI技术类文章。\\n请严格按JSON格式输出，包含一个键为"content"的值，值为文章内容（Markdown格式）。\\n\\n文章风格要求：\\n1. 使用震惊体标题\\n2. 每小节必须使用不同的emoji符号开头\\n3. 段落之间必须用空行分隔（两个换行符）\\n4. 包含具体案例对比（优化前VS优化后），使用清晰对比格式\\n5. 使用活泼语气和夸张表达（如"效果炸裂"、"老板加薪"等）\\n6. 结尾用鼓励性语言+emoji\\n\\n严格遵循以下结构但使用不同的章节标题：\\n# [标题emoji+震惊体标题]\\n\\n## [新emoji] 引言\\n[网络流行语开头]\\n\\n## [新emoji] 核心原则\\n[分条列出4项原则，使用不同描述]\\n\\n## [新emoji] 实用技巧\\n### [新emoji] 技巧1\\n- 描述\\n- 示例\\n\\n## [新emoji] 真实场景\\n### 场景1：[新名称]\\n**优化前**：...\\n**优化后**：...\\n\\n## [新emoji] 高级策略\\n[分条列出技巧]\\n\\n## [新emoji] 未来展望\\n[预测性语言]\\n\\n## [新emoji] 结语\\n[鼓励语+emoji]`,
			},
			{
				"role":    "user",
				"content": `请生成全新《AI提示词工程》文章，要求：\\n1. 标题：使用[新比喻]+emoji（禁用"惊爆/震惊"）\\n2. 结构：\\n   ## 引言（用当周流行语）\\n   ## [新emoji] 核心原则（重组顺序）\\n   ## [新emoji] 实用技巧（选3种新组合）\\n   ## [新emoji] 真实场景（2个新领域）\\n   ## [新emoji] 高级策略\\n   ## [新emoji] 未来展望\\n   ## [新emoji] 结语（新鼓励语）\\n3. 内容：\\n   - 原则：从{清晰,目标,开放,迭代}选4项但用不同表述\\n   - 技巧：从{任务导向,角色扮演,分步指导}选3种但创新描述\\n   - 场景：从{互联网,游戏,金融,电商,医疗}选2领域\\n   - 趋势：从{多模态,个性化,自动化,伦理安全}选3项\\n4. 格式：\\n   - 每章节用不同emoji\\n   - 案例用**优化前/后**格式\\n   - 禁用数字列表\\n\\n必须：\\n- 使用最新网络流行语\\n- 案例参考2025年行业报告\\n- 每200字≥3个不同emoji`,
			},
		},
		"temperature": 0.9, // 提高温度增加创造性
		"max_tokens":  cfg.MaxTokens + 500,
		"response_format": map[string]string{
			"type": "json_object",
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		if logger != nil {
			logger.Printf("❌ JSON序列化失败: %v", err)
		}
		return "", fmt.Errorf("JSON序列化失败: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		apiBase+"/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		if logger != nil {
			logger.Printf("❌ 创建API请求失败: %v", err)
		}
		return "", fmt.Errorf("创建API请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	if logger != nil {
		logger.Printf("📤 发送请求到DeepSeek API，主题: 《AI提示词工程》")
	}

	startTime := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		if logger != nil {
			logger.Printf("❌ API请求失败: %v", err)
			logger.Printf("⏱️ 耗时: %v", duration)
		}
		return "", fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if logger != nil {
			logger.Printf("❌ API错误状态: %d", resp.StatusCode)
			logger.Printf("⏱️ 耗时: %v", duration)
			body, _ := io.ReadAll(resp.Body)
			logger.Printf("📋 错误响应: %s", string(body))
		}
		return "", fmt.Errorf("API错误状态: %d", resp.StatusCode)
	} else if logger != nil {
		logger.Printf("✅ 请求成功! 状态码: %d", resp.StatusCode)
		logger.Printf("⏱️ 耗时: %v", duration)
	}

	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if logger != nil {
			logger.Printf("❌ 读取响应失败: %v", err)
		}
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if logger != nil {
		logger.Printf("🔄 解析API响应")
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		if logger != nil {
			logger.Printf("❌ 解析响应失败: %v", err)
			logger.Printf("📋 原始响应: %s", string(body))
		}
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(apiResponse.Choices) == 0 {
		errMsg := "❌ 生成失败"
		if apiResponse.Error.Message != "" {
			errMsg += ": " + apiResponse.Error.Message
		}
		if logger != nil {
			logger.Println(errMsg)
		}
		return "", fmt.Errorf(errMsg)
	}

	var contentResp DeepSeekResponse
	contentStr := apiResponse.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(contentStr), &contentResp); err != nil {
		if logger != nil {
			logger.Printf("❌ 解析内容失败: %v", err)
			logger.Printf("📋 原始内容: %s", contentStr)
		}
		return "", fmt.Errorf("解析内容失败: %w", err)
	}

	content := contentResp.Content
	if content != "" {
		filePath := filepath.Join(config.AppConfig.LocalArticlePath, date+".md")
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			if logger != nil {
				logger.Printf("❌ 保存失败: %v", err)
			}
		} else if logger != nil {
			logger.Printf("💾 文章保存成功: %s", filePath)
			logger.Printf("📝 标题: %s", extractTitle(content))
		}
	} else if logger != nil {
		logger.Println("⚠️ 内容为空")
	}

	return content, nil
}

// 提取标题（保持不变）
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return "无标题"
}
