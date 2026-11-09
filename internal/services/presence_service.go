package services

import (
	"context"
	"sync"
	"time"

	"project-studio/internal/auth"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PresenceService struct {
	ctx           context.Context
	activeUsers   map[string]time.Time
	mu            sync.RWMutex
	checkInterval time.Duration
	timeout       time.Duration
}

func NewPresenceService() *PresenceService {
	return &PresenceService{
		activeUsers:   make(map[string]time.Time),
		checkInterval: 10 * time.Second,
		timeout:       30 * time.Second,
	}
}

func (s *PresenceService) Startup(ctx context.Context) {
	s.ctx = ctx
	go s.monitorPresence()
}

// Heartbeat is called by the frontend periodically (e.g. every 15s)
func (s *PresenceService) Heartbeat() {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return
	}

	s.mu.Lock()
	isNew := false
	if _, exists := s.activeUsers[userID]; !exists {
		isNew = true
	}
	s.activeUsers[userID] = time.Now()
	s.mu.Unlock()

	// If user just came online, emit event
	if isNew && s.ctx != nil {
		runtime.EventsEmit(s.ctx, "rt:presence:update", map[string]interface{}{
			"user_id": userID,
			"status":  "online",
		})
	}
}

// GetOnlineUsers returns a list of currently online user IDs
func (s *PresenceService) GetOnlineUsers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []string
	now := time.Now()
	for id, lastSeen := range s.activeUsers {
		if now.Sub(lastSeen) <= s.timeout {
			users = append(users, id)
		}
	}
	return users
}

func (s *PresenceService) monitorPresence() {
	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now()
			for id, lastSeen := range s.activeUsers {
				if now.Sub(lastSeen) > s.timeout {
					// User went offline
					delete(s.activeUsers, id)
					if s.ctx != nil {
						runtime.EventsEmit(s.ctx, "rt:presence:update", map[string]interface{}{
							"user_id": id,
							"status":  "offline",
						})
					}
				}
			}
			s.mu.Unlock()
		case <-s.ctx.Done():
			return
		}
	}
}
