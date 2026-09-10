package ai

import "gorm.io/gorm"

// ChatRepository - interface untuk operasi DB chat
type ChatRepository interface {
	CreateSession(session *ChatSession) error
	GetSessionByID(id string) (*ChatSession, error)
	GetSessionsByUserID(userID string) ([]ChatSession, error)
	DeleteSession(id string) error

	SaveMessage(message *ChatMessage) error
	GetMessagesBySessionID(sessionID string, limit int) ([]ChatMessage, error)
}

type chatRepository struct {
	db *gorm.DB
}

// NewChatRepository creates a new chat repository
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateSession(session *ChatSession) error {
	return r.db.Create(session).Error
}

func (r *chatRepository) GetSessionByID(id string) (*ChatSession, error) {
	var session ChatSession
	err := r.db.Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *chatRepository) GetSessionsByUserID(userID string) ([]ChatSession, error) {
	var sessions []ChatSession
	err := r.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&sessions).Error
	return sessions, err
}

func (r *chatRepository) DeleteSession(id string) error {
	return r.db.Where("id = ?", id).Delete(&ChatSession{}).Error
}

func (r *chatRepository) SaveMessage(message *ChatMessage) error {
	return r.db.Create(message).Error
}

func (r *chatRepository) GetMessagesBySessionID(sessionID string, limit int) ([]ChatMessage, error) {
	var messages []ChatMessage
	// Ambil pesan terbaru sejumlah limit, order descending untuk limit, tapi return order ascending
	err := r.db.Where("session_id = ?", sessionID).Order("created_at desc").Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	
	// Reverse supaya urutan jadi dari lama ke baru
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
