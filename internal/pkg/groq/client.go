package groq

import (
	"atlas_food/internal/pkg/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client - wrapper Groq chat completions
type Client struct {
	apiKey     string
	model      string
	baseURL    string
	timeout    time.Duration
	maxTokens  int
	httpClient *http.Client
}

// Result - hasil mentah dari Groq
type Result struct {
	Content     string          `json:"content"`
	RawResponse json.RawMessage `json:"raw_response"`
	ModelUsed   string          `json:"model_used"`
	TokenUsed   *int            `json:"token_used,omitempty"`
	LatencyMs   *int            `json:"latency_ms,omitempty"`
}

// NewClient - buat client Groq
func NewClient(apiKey, model, baseURL string, timeoutSecs, maxTokens int) *Client {
	if model == "" {
		model = "llama3-8b-8192"
	}
	if baseURL == "" {
		baseURL = "https://api.groq.com/openai/v1"
	}
	if timeoutSecs <= 0 {
		timeoutSecs = 15
	}
	if maxTokens <= 0 {
		maxTokens = 512
	}
	return &Client{
		apiKey:     apiKey,
		model:      model,
		baseURL:    strings.TrimRight(baseURL, "/"),
		timeout:    time.Duration(timeoutSecs) * time.Second,
		maxTokens:  maxTokens,
		httpClient: &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
	}
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model          string            `json:"model"`
	Messages       []groqMessage     `json:"messages"`
	Temperature    float64           `json:"temperature"`
	MaxTokens      int               `json:"max_tokens"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Analyze - panggil Groq dengan prompt siap pakai
func (c *Client) Analyze(systemPrompt, userPrompt string) (*Result, error) {
	if c.apiKey == "" {
		return nil, utils.NewAppError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "GROQ_API_KEY belum dikonfigurasi")
	}

	systemPrompt = strings.TrimSpace(systemPrompt)
	userPrompt = strings.TrimSpace(userPrompt)

	fullSystemPrompt := `You are a nutrition analyst for meal recall data. Return ONLY valid JSON matching this schema:
{
  "overall_status": "good|less|excess",
  "overall_message": "string",
  "nutritional_analysis": [{"label":"Calories|Protein|Balance","status":"low|partial|good|high","description":"string"}],
  "ai_recommendation": "string",
  "recommended_foods": ["string"],
  "health_insight": {"title":"string","description":"string"},
  "suggested_activities": ["string"]
}
Do not add markdown fences. Keep recommendations practical and concise.`

	reqBody := groqRequest{
		Model: c.model,
		Messages: []groqMessage{
			{Role: "system", Content: fullSystemPrompt + "\n\n" + systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature:    0.2,
		MaxTokens:      c.maxTokens,
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	started := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, utils.NewAppError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", fmt.Sprintf("Groq error: %s", strings.TrimSpace(string(responseBytes))))
	}

	var parsed groqResponse
	if err := json.Unmarshal(responseBytes, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("groq response is empty")
	}

	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	content = trimCodeFence(content)

	latency := int(time.Since(started).Milliseconds())
	tokens := parsed.Usage.TotalTokens
	return &Result{
		Content:     content,
		RawResponse: responseBytes,
		ModelUsed:   c.model,
		TokenUsed:   &tokens,
		LatencyMs:   &latency,
	}, nil
}

func trimCodeFence(content string) string {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content)
}
