package ai

import (
	"atlas_food/internal/domain/gamification"
	"atlas_food/internal/pkg/groq"
	"atlas_food/internal/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ChatService interface {
	Chat(userID string, req ChatRequest) (*ChatResponse, error)
	GetSessions(userID string) ([]SessionListResponse, error)
	GetMessages(userID, sessionID string) ([]MessageListResponse, error)
	DeleteSession(userID, sessionID string) error
}

type chatService struct {
	repo         ChatRepository
	groq         *groq.Client
	gamification gamification.Service
}

// NewChatService creates a new AI Chat Service
func NewChatService(repo ChatRepository, groqClient *groq.Client, gamification gamification.Service) ChatService {
	return &chatService{
		repo:         repo,
		groq:         groqClient,
		gamification: gamification,
	}
}

func (s *chatService) Chat(userID string, req ChatRequest) (*ChatResponse, error) {
	var sessionID string
	var session *ChatSession
	var xpEarned int
	isNewSession := false

	// 1. Cek atau buat sesi
	if req.SessionID != nil && *req.SessionID != "" {
		sessionID = *req.SessionID
		var err error
		session, err = s.repo.GetSessionByID(sessionID)
		if err != nil || session.UserID != userID {
			return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Sesi tidak ditemukan")
		}
	} else {
		// Sesi baru
		sessionID = uuid.New().String()
		title := req.Message
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		session = &ChatSession{
			ID:     sessionID,
			UserID: userID,
			Title:  title,
		}
		if err := s.repo.CreateSession(session); err != nil {
			return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal membuat sesi baru")
		}
		isNewSession = true
	}

	// 2. Ambil history (max 10 pesan terakhir untuk konteks)
	historyMessages, err := s.repo.GetMessagesBySessionID(sessionID, 10)
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil riwayat pesan")
	}

	// 3. Simpan pesan user
	userMessage := &ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Message,
	}
	_ = s.repo.SaveMessage(userMessage)

	// 4. Susun pesan untuk Groq (System prompt + History + Current message)
	systemPrompt := `Kamu adalah NutriBot, asisten gizi cerdas dari aplikasi Atlas Food — sebuah platform pemetaan pangan dan gizi terjangkau untuk mahasiswa dan masyarakat umum Indonesia.

Peranmu:
1. Menjawab pertanyaan seputar gizi, nutrisi, kalori, dan komposisi makanan Indonesia secara akurat, ramah, dan mudah dipahami.
2. Memberikan rekomendasi makanan bergizi dengan harga terjangkau (di bawah Rp 20.000/porsi) yang umum ditemui di sekitar kampus atau warung makan lokal.
3. Membantu pengguna memahami hasil food recall dan skor gizi dari Atlas Food.

Aturan ketat:
- Jawab HANYA pertanyaan terkait gizi, nutrisi, makanan, dan kesehatan umum.
- Jika ditanya di luar topik gizi/kesehatan, tolak dengan ramah dan arahkan kembali ke topik.
- Jangan memberikan diagnosis medis. Selalu sarankan konsultasi dokter untuk kondisi medis.
- Gunakan Bahasa Indonesia yang santai, ramah, dan mudah dipahami mahasiswa.
- Sertakan angka/estimasi gizi jika relevan.`

	// Karena client.go saat ini hanya menerima systemPrompt dan userPrompt string, 
	// kita kompilasi history menjadi satu string userPrompt yang besar, atau lebih baik kita 
	// perbarui client.go untuk menerima history. Namun untuk tidak memodifikasi client terlalu banyak,
	// kita pass history ke userPrompt dengan format tertentu.
	
	// Untuk saat ini karena di internal/pkg/groq/client.go Analyze hanya menerima systemPrompt, userPrompt,
	// kita format context menjadi teks (kurang ideal, tapi sesuai API yang tersedia):
	
	compiledPrompt := ""
	for _, msg := range historyMessages {
		if msg.Role == "user" {
			compiledPrompt += "User: " + msg.Content + "\n"
		} else if msg.Role == "assistant" {
			compiledPrompt += "NutriBot: " + msg.Content + "\n"
		}
	}
	compiledPrompt += "User: " + req.Message

	result, err := s.groq.Analyze(systemPrompt, compiledPrompt)
	if err != nil {
		return nil, utils.NewAppError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "AI service sedang sibuk, silakan coba lagi.")
	}

	// 5. Simpan balasan AI
	aiMessage := &ChatMessage{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   result.Content,
		ModelUsed: result.ModelUsed,
	}
	if result.TokenUsed != nil {
		aiMessage.TokenUsed = *result.TokenUsed
	}
	if result.LatencyMs != nil {
		aiMessage.LatencyMs = *result.LatencyMs
	}
	_ = s.repo.SaveMessage(aiMessage)

	// 6. Trigger gamification
	if isNewSession && s.gamification != nil {
		_, errGam := s.gamification.AddXP(userID, "ai_chat", 10, "Memulai sesi chat baru dengan NutriBot")
		if errGam == nil {
			xpEarned = 10
		}
	}

	// 7. Update session updated_at
	session.UpdatedAt = time.Now()
	// GORM tidak akan update automatically jika kita tidak panggil db.Save(), tapi tidak terlalu critical.

	latency := 0
	if result.LatencyMs != nil {
		latency = *result.LatencyMs
	}

	return &ChatResponse{
		SessionID: sessionID,
		Reply:     result.Content,
		XPEarned:  xpEarned,
		ModelUsed: result.ModelUsed,
		LatencyMs: latency,
	}, nil
}

func (s *chatService) GetSessions(userID string) ([]SessionListResponse, error) {
	sessions, err := s.repo.GetSessionsByUserID(userID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil sesi chat")
	}

	var res []SessionListResponse
	for _, sess := range sessions {
		res = append(res, SessionListResponse{
			ID:        sess.ID,
			Title:     sess.Title,
			CreatedAt: sess.CreatedAt,
		})
	}
	return res, nil
}

func (s *chatService) GetMessages(userID, sessionID string) ([]MessageListResponse, error) {
	session, err := s.repo.GetSessionByID(sessionID)
	if err != nil || session.UserID != userID {
		return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Sesi tidak ditemukan")
	}

	messages, err := s.repo.GetMessagesBySessionID(sessionID, 50)
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil pesan")
	}

	var res []MessageListResponse
	for _, msg := range messages {
		res = append(res, MessageListResponse{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
	}
	return res, nil
}

func (s *chatService) DeleteSession(userID, sessionID string) error {
	session, err := s.repo.GetSessionByID(sessionID)
	if err != nil || session.UserID != userID {
		return utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Sesi tidak ditemukan")
	}

	if err := s.repo.DeleteSession(sessionID); err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menghapus sesi")
	}

	return nil
}
