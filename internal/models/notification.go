package models

import "time"

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	Type      string    `json:"type"`       // 'task_assigned', 'mention', 'comment', etc.
	RelatedID string    `json:"related_id"` // ID of the task/project
	CreatedAt time.Time `json:"created_at"`
}

type NotificationListResponse struct {
	Success       bool           `json:"success"`
	Message       string         `json:"message"`
	Notifications []Notification `json:"notifications,omitempty"`
	UnreadCount   int            `json:"unread_count"`
}

type NotificationResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	Notification *Notification `json:"notification,omitempty"`
}

type NotificationActionRequest struct {
	ID string `json:"id"`
}
