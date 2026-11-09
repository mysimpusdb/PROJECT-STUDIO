package models

import "time"

type ActivityLog struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name,omitempty"` // For frontend
	Action     string    `json:"action"`              // 'created', 'updated', 'deleted'
	EntityType string    `json:"entity_type"`         // 'task', 'project', 'comment'
	EntityID   string    `json:"entity_id,omitempty"`
	Details    string    `json:"details,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ActivityListResponse struct {
	Success    bool          `json:"success"`
	Message    string        `json:"message"`
	Activities []ActivityLog `json:"activities,omitempty"`
}
