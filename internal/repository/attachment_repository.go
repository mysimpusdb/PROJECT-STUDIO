package repository

import (
	"database/sql"
	"time"

	"project-studio/internal/models"
)

type AttachmentRepository struct {
	db *sql.DB
}

func NewAttachmentRepository(db *sql.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(attachment *models.TaskAttachment) error {
	query := `
		INSERT INTO task_attachments (id, task_id, user_id, file_name, file_path, file_size, uploaded_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, attachment.ID, attachment.TaskID, attachment.UserID, attachment.FileName, attachment.FilePath, attachment.FileSize, attachment.UploadedAt)
	return err
}

func (r *AttachmentRepository) GetByTaskID(taskID string) ([]models.TaskAttachment, error) {
	query := `
		SELECT a.id, a.task_id, a.user_id, u.full_name, a.file_name, a.file_path, a.file_size, a.uploaded_at
		FROM task_attachments a
		JOIN users u ON a.user_id = u.id
		WHERE a.task_id = ?
		ORDER BY a.uploaded_at DESC
	`
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.TaskAttachment
	layout := "2006-01-02 15:04:05"

	for rows.Next() {
		var a models.TaskAttachment
		var uploadedAtStr string
		err := rows.Scan(&a.ID, &a.TaskID, &a.UserID, &a.UserName, &a.FileName, &a.FilePath, &a.FileSize, &uploadedAtStr)
		if err != nil {
			return nil, err
		}
		if pt, err := time.Parse(layout, uploadedAtStr); err == nil {
			a.UploadedAt = pt
		} else if pt, err := time.Parse(time.RFC3339, uploadedAtStr); err == nil {
			a.UploadedAt = pt
		}
		attachments = append(attachments, a)
	}

	return attachments, nil
}

func (r *AttachmentRepository) GetByID(id string) (*models.TaskAttachment, error) {
	query := `
		SELECT id, task_id, user_id, file_name, file_path, file_size, uploaded_at
		FROM task_attachments
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var a models.TaskAttachment
	var uploadedAtStr string
	err := row.Scan(&a.ID, &a.TaskID, &a.UserID, &a.FileName, &a.FilePath, &a.FileSize, &uploadedAtStr)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02 15:04:05"
	if pt, err := time.Parse(layout, uploadedAtStr); err == nil {
		a.UploadedAt = pt
	} else if pt, err := time.Parse(time.RFC3339, uploadedAtStr); err == nil {
		a.UploadedAt = pt
	}

	return &a, nil
}

func (r *AttachmentRepository) Delete(id string) error {
	query := "DELETE FROM task_attachments WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
