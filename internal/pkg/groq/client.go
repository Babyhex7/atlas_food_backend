package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ErrorKind - jenis kegagalan pemanggilan Groq. Pemanggil memetakannya ke
// status HTTP dan pesan untuk user; paket ini sengaja tidak tahu soal HTTP API
// milik aplikasi supaya bisa dipakai domain lain.
type ErrorKind string

const (
	// ErrNotConfigured - GROQ_API_KEY kosong
	ErrNotConfigured ErrorKind = "not_configured"
	// ErrRateLimited - Groq membalas 429 dan percobaan ulang sudah habis
	ErrRateLimited ErrorKind = "rate_limited"
	// ErrTimeout - batas waktu terlampaui
	ErrTimeout ErrorKind = "timeout"
	// ErrUnavailable - gangguan jaringan atau Groq membalas 5xx
	ErrUnavailable ErrorKind = "unavailable"
	// ErrRejected - Groq menolak permintaan (4xx selain 429): kunci salah, model tidak ada, dsb.
	ErrRejected ErrorKind = "rejected"
	// ErrEmptyResponse - balasan tidak memuat isi yang bisa dipakai
	ErrEmptyResponse ErrorKind = "empty_response"
)

// Error - kegagalan pemanggilan Groq beserta jenisnya.
// Detail mentah dari Groq hanya dicatat di log server, tidak pernah dibawa di sini,
// supaya tidak bocor ke response API.
type Error struct {
	Kind       ErrorKind
	StatusCode int
}

// Error - implementasi interface error
func (e *Error) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("groq: %s (status %d)", e.Kind, e.StatusCode)
	}
	return fmt.Sprintf("groq: %s", e.Kind)
}

// KindOf - ambil jenis kegagalan dari sebuah error; ErrUnavailable untuk error yang tidak dikenal
func KindOf(err error) ErrorKind {
	var gErr *Error
	if errors.As(err, &gErr) {
		return gErr.Kind
	}
	return ErrUnavailable
}

// maxAttempts - jumlah percobaan total untuk kegagalan yang layak diulang (429, 5xx, jaringan)
const maxAttempts = 3

// Client - wrapper Groq chat completions (mode JSON)
type Client struct {
	apiKey     string
	model      string
	baseURL    string
	timeout    time.Duration
	maxTokens  int
	httpClient *http.Client
}

// Result - hasil satu pemanggilan chat completion
type Result struct {
	Content          string
	ModelUsed        string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	LatencyMs        int
}

// NewClient - buat client Groq
func NewClient(apiKey, model, baseURL string, timeoutSecs, maxTokens int) *Client {
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}
	if baseURL == "" {
		baseURL = "https://api.groq.com/openai/v1"
	}
	if timeoutSecs <= 0 {
		timeoutSecs = 45
	}
	if maxTokens <= 0 {
		maxTokens = 2048
	}
	return &Client{
		apiKey:    strings.TrimSpace(apiKey),
		model:     model,
		baseURL:   strings.TrimRight(baseURL, "/"),
		timeout:   time.Duration(timeoutSecs) * time.Second,
		maxTokens: maxTokens,
		// Batas waktu dipegang context per pemanggilan (mencakup seluruh percobaan ulang),
		// bukan http.Client.Timeout yang berlaku per request.
		httpClient: &http.Client{},
	}
}

// Configured - true kalau API key tersedia
func (c *Client) Configured() bool {
	return c != nil && c.apiKey != ""
}

// Model - nama model yang dipakai client ini
func (c *Client) Model() string {
	return c.model
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model          string            `json:"model"`
	Messages       []chatMessage     `json:"messages"`
	Temperature    float64           `json:"temperature"`
	MaxTokens      int               `json:"max_tokens"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

type chatResponse struct {
	Model   string `json:"model"`
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

// CompleteJSON - kirim satu percakapan system+user dan minta jawaban berupa objek JSON.
// Kegagalan sementara (429, 5xx, jaringan) diulang dengan jeda bertambah selama
// batas waktu client belum habis. Error yang dikembalikan selalu bertipe *Error.
func (c *Client) CompleteJSON(ctx context.Context, systemPrompt, userPrompt string) (*Result, error) {
	if !c.Configured() {
		return nil, &Error{Kind: ErrNotConfigured}
	}

	body, err := json.Marshal(chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: strings.TrimSpace(systemPrompt)},
			{Role: "user", Content: strings.TrimSpace(userPrompt)},
		},
		Temperature:    0.2,
		MaxTokens:      c.maxTokens,
		ResponseFormat: map[string]string{"type": "json_object"},
	})
	if err != nil {
		return nil, &Error{Kind: ErrRejected}
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	started := time.Now()
	var lastErr *Error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, callErr := c.call(ctx, body)
		if callErr == nil {
			result.LatencyMs = int(time.Since(started).Milliseconds())
			return result, nil
		}
		lastErr = callErr
		if !retryable(callErr.Kind) || attempt == maxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return nil, &Error{Kind: ErrTimeout}
		case <-time.After(time.Duration(attempt) * 700 * time.Millisecond):
		}
	}
	return nil, lastErr
}

// retryable - true untuk kegagalan yang mungkin hilang bila dicoba lagi
func retryable(kind ErrorKind) bool {
	return kind == ErrRateLimited || kind == ErrUnavailable
}

// call - satu request HTTP ke Groq
func (c *Client) call(ctx context.Context, body []byte) (*Result, *Error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, &Error{Kind: ErrRejected}
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, &Error{Kind: ErrTimeout}
		}
		log.Printf("[groq] request gagal: %v", err)
		return nil, &Error{Kind: ErrUnavailable}
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		if ctx.Err() != nil {
			return nil, &Error{Kind: ErrTimeout}
		}
		return nil, &Error{Kind: ErrUnavailable}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Isi balasan Groq bisa memuat detail akun/model — cukup di log server.
		log.Printf("[groq] status %d: %s", resp.StatusCode, truncate(string(responseBytes), 500))
		switch {
		case resp.StatusCode == http.StatusTooManyRequests:
			return nil, &Error{Kind: ErrRateLimited, StatusCode: resp.StatusCode}
		case resp.StatusCode >= 500:
			return nil, &Error{Kind: ErrUnavailable, StatusCode: resp.StatusCode}
		default:
			return nil, &Error{Kind: ErrRejected, StatusCode: resp.StatusCode}
		}
	}

	var parsed chatResponse
	if err := json.Unmarshal(responseBytes, &parsed); err != nil || len(parsed.Choices) == 0 {
		return nil, &Error{Kind: ErrEmptyResponse}
	}

	content := trimCodeFence(parsed.Choices[0].Message.Content)
	if content == "" {
		return nil, &Error{Kind: ErrEmptyResponse}
	}

	model := parsed.Model
	if model == "" {
		model = c.model
	}
	return &Result{
		Content:          content,
		ModelUsed:        model,
		PromptTokens:     parsed.Usage.PromptTokens,
		CompletionTokens: parsed.Usage.CompletionTokens,
		TotalTokens:      parsed.Usage.TotalTokens,
	}, nil
}

// trimCodeFence - buang pembungkus markdown ```json ... ``` dari jawaban LLM supaya bisa di-parse sebagai JSON
func trimCodeFence(content string) string {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content)
}

// truncate - potong string panjang untuk keperluan log
func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
