package collab

// Service handles collaboration business logic
type Service struct {
	hub *Hub
}

// NewService creates a new collaboration service
func NewService(hub *Hub) *Service {
	return &Service{
		hub: hub,
	}
}

// GetHub returns the hub instance
func (s *Service) GetHub() *Hub {
	return s.hub
}
