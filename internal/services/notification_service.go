package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"project-studio/internal/auth"
	"project-studio/internal/db"
	"project-studio/internal/models"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type NotificationService struct {
	ctx context.Context
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *NotificationService) GetMyNotifications(limit int, offset int) models.NotificationListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.NotificationListResponse{Success: false, Message: "Unauthenticated"}
	}

	// Default pagination
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT id, user_id, title, message, is_read, type, related_id, created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.DB.Query(query, userID, limit, offset)
	if err != nil {
		return models.NotificationListResponse{Success: false, Message: "Gagal mengambil notifikasi."}
	}
	defer rows.Close()

	var notifs []models.Notification
	layout := "2006-01-02 15:04:05"

	for rows.Next() {
		var n models.Notification
		var createdAtStr string
		err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.IsRead, &n.Type, &n.RelatedID, &createdAtStr)
		if err != nil {
			return models.NotificationListResponse{Success: false, Message: "Gagal mem-parsing notifikasi."}
		}
		if pt, err := time.Parse(layout, createdAtStr); err == nil {
			n.CreatedAt = pt
		}
		notifs = append(notifs, n)
	}

	// Count unread
	var unreadCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0", userID).Scan(&unreadCount)

	return models.NotificationListResponse{
		Success:       true,
		Message:       "Notifikasi berhasil dimuat.",
		Notifications: notifs,
		UnreadCount:   unreadCount,
	}
}

func (s *NotificationService) MarkAsRead(id string) models.NotificationResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.NotificationResponse{Success: false, Message: "Unauthenticated"}
	}

	_, err = db.DB.Exec("UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return models.NotificationResponse{Success: false, Message: "Gagal menandai notifikasi."}
	}

	return models.NotificationResponse{
		Success: true,
		Message: "Notifikasi ditandai sudah dibaca.",
	}
}

func (s *NotificationService) MarkAllAsRead() models.NotificationResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.NotificationResponse{Success: false, Message: "Unauthenticated"}
	}

	_, err = db.DB.Exec("UPDATE notifications SET is_read = 1 WHERE user_id = ?", userID)
	if err != nil {
		return models.NotificationResponse{Success: false, Message: "Gagal menandai semua notifikasi."}
	}

	// Trigger global unread update event
	runtime.EventsEmit(s.ctx, "rt:notification:count:"+userID, 0)

	return models.NotificationResponse{
		Success: true,
		Message: "Semua notifikasi ditandai sudah dibaca.",
	}
}

func (s *NotificationService) ClearAll() models.NotificationResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.NotificationResponse{Success: false, Message: "Unauthenticated"}
	}

	_, err = db.DB.Exec("DELETE FROM notifications WHERE user_id = ?", userID)
	if err != nil {
		return models.NotificationResponse{Success: false, Message: "Gagal membersihkan notifikasi."}
	}

	// Trigger global unread update event
	runtime.EventsEmit(s.ctx, "rt:notification:count:"+userID, 0)

	return models.NotificationResponse{
		Success: true,
		Message: "Semua notifikasi dibersihkan.",
	}
}

// CreateNotification is an internal helper function called by other services
func (s *NotificationService) CreateNotification(userID, title, message, notifType, relatedID string) error {
	notif := models.Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		Message:   message,
		IsRead:    false,
		Type:      notifType,
		RelatedID: relatedID,
	}

	query := `
		INSERT INTO notifications (id, user_id, title, message, is_read, type, related_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.DB.Exec(query, notif.ID, notif.UserID, notif.Title, notif.Message, notif.IsRead, notif.Type, notif.RelatedID)
	if err != nil {
		return err
	}

	notif.CreatedAt = time.Now()

	// 1. Emit Wails event for real-time in-app UI update
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "rt:notification:new:"+userID, notif)
	}

	// 2. Mock Email Sending (SMTP simulation)
	s.mockSendEmail(userID, title, message)

	return nil
}

func (s *NotificationService) mockSendEmail(userID, title, message string) {
	// Simulate fetching user email
	var email string
	err := db.DB.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&email)
	if err == nil && email != "" {
		log.Printf("\n[MOCK EMAIL SMTP] Sending to: %s\nSubject: %s\nBody: %s\n", email, title, message)
		fmt.Printf("[MOCK EMAIL SMTP] Email sent to %s successfully.\n", email)
	}
}
