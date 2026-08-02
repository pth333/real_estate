package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/global"
)

type AIRepository struct {
	httpClient *http.Client
}

type IAIRepository interface {
	GenerateContent(ctx context.Context, prompt string) (dto.AIContentResponse, error)
}

func NewAIRepository() IAIRepository {
	return &AIRepository{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *AIRepository) GenerateContent(ctx context.Context, prompt string) (dto.AIContentResponse, error) {
	if global.Config.AI.APIKey == "" {
		return dto.AIContentResponse{}, fmt.Errorf("AI API key chưa được cấu hình (RE_AI_API_KEY)")
	}

	body := map[string]any{
		"model": global.Config.AI.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Bạn là trợ lý viết tin đăng bất động sản tiếng Việt.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return dto.AIContentResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return dto.AIContentResponse{}, fmt.Errorf("tạo request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+global.Config.AI.APIKey)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return dto.AIContentResponse{}, fmt.Errorf("gọi API AI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return dto.AIContentResponse{}, fmt.Errorf("API AI trả lỗi %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response OpenAI: choices[0].message.content
	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return dto.AIContentResponse{}, fmt.Errorf("parse response AI: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return dto.AIContentResponse{}, fmt.Errorf("API AI trả về không có nội dung")
	}

	// Prompt yêu cầu AI trả JSON { title, description } — parse nếu có
	return parseAIContent(chatResp.Choices[0].Message.Content)
}

// parseAIContent tách nội dung AI trả về thành title + description riêng biệt.
// Xử lý các trường hợp AI trả về:
//   1. JSON thuần:  {"title": "...", "description": "..."}
//   2. JSON bọc trong markdown code block: ```json ... ```
//   3. Văn bản thường (không JSON): title = dòng đầu, description = phần còn lại
func parseAIContent(content string) (dto.AIContentResponse, error) {
	trimmed := strings.TrimSpace(content)

	// Bỏ khung markdown code block nếu có (```json ... ```)
	trimmed = stripCodeBlock(trimmed)

	var result dto.AIContentResponse
	if err := json.Unmarshal([]byte(trimmed), &result); err == nil {
		// Trường hợp 1 & 2: có JSON hợp lệ → dùng luôn
		return result, nil
	}

	// Trường hợp 3: không phải JSON → tách theo dòng
	return fallbackSplitContent(trimmed), nil
}

// stripCodeBlock bỏ khung ```json ... ``` bao quanh JSON nếu AI trả dạng markdown
func stripCodeBlock(content string) string {
	lines := strings.Split(content, "\n")
	// Bỏ dòng đầu nếu là ``` hoặc ```json
	if len(lines) > 0 {
		first := strings.TrimSpace(lines[0])
		if strings.HasPrefix(first, "```") {
			lines = lines[1:]
		}
	}
	// Bỏ dòng cuối nếu là ```
	if len(lines) > 0 {
		last := strings.TrimSpace(lines[len(lines)-1])
		if strings.HasPrefix(last, "```") {
			lines = lines[:len(lines)-1]
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// fallbackSplitContent tách văn bản thường: dòng đầu = title, còn lại = description
func fallbackSplitContent(content string) dto.AIContentResponse {
	lines := strings.Split(content, "\n")
	var meaningful []string
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			meaningful = append(meaningful, trimmed)
		}
	}

	if len(meaningful) == 0 {
		return dto.AIContentResponse{}
	}

	// Dòng đầu tiên là tiêu đề, các dòng còn lại nối thành mô tả
	title := meaningful[0]
	description := strings.Join(meaningful[1:], "\n")
	return dto.AIContentResponse{
		Title:       title,
		Description: description,
	}
}
