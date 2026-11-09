package models

import "time"

type TaskAttachment struct {
	ID         string    `json:"id"`
	TaskID     string    `json:"task_id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name,omitempty"` // For UI
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	FileSize   int64     `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type AttachmentResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Attachment *TaskAttachment `json:"attachment,omitempty"`
}

type AttachmentListResponse struct {
	Success     bool             `json:"success"`
	Message     string           `json:"message"`
	Attachments []TaskAttachment `json:"attachments,omitempty"`
}
