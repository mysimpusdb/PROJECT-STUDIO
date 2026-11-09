package models

import "time"

type TimeLog struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"task_id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name,omitempty"` // for UI
	LoggedDate  string    `json:"logged_date"`         // YYYY-MM-DD
	Hours       float64   `json:"hours"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateTimeLogRequest struct {
	TaskID      string  `json:"task_id"`
	LoggedDate  string  `json:"logged_date"`
	Hours       float64 `json:"hours"`
	Description string  `json:"description,omitempty"`
}

type UpdateTimeLogRequest struct {
	ID          string  `json:"id"`
	LoggedDate  string  `json:"logged_date"`
	Hours       float64 `json:"hours"`
	Description string  `json:"description,omitempty"`
}

type TimeLogResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	TimeLog *TimeLog `json:"time_log,omitempty"`
}

type TimeLogListResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	TimeLogs []TimeLog `json:"time_logs,omitempty"`
}
