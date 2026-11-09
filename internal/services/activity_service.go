package services

import (
	"context"
	"time"

	"project-studio/internal/auth"
	"project-studio/internal/db"
	"project-studio/internal/models"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ActivityService struct {
	ctx context.Context
}

func NewActivityService() *ActivityService {
	return &ActivityService{}
}

func (s *ActivityService) Startup(ctx context.Context) {
	s.ctx = ctx
}

// RecordActivity is called internally by other services
func (s *ActivityService) RecordActivity(userID, action, entityType, entityID, details string) error {
	log := models.ActivityLog{
		ID:         uuid.New().String(),
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
	}

	query := `
		INSERT INTO activity_logs (id, user_id, action, entity_type, entity_id, details, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.DB.Exec(query, log.ID, log.UserID, log.Action, log.EntityType, log.EntityID, log.Details, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	// Emit event for real-time dashboard updates if needed
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "rt:activity:new", log)
	}

	return nil
}

func (s *ActivityService) GetProjectActivities(projectID string, limit int, offset int) models.ActivityListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ActivityListResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, projectID, permRead); err != nil {
		return models.ActivityListResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	if limit <= 0 {
		limit = 50
	}

	// This is a basic implementation. We'll need a way to filter activities by project.
	// Since activity_logs doesn't have a project_id column, we can do a hacky check or we'd need to alter the schema.
	// We'll alter the query to look at tasks that belong to the project.
	query := `
		SELECT a.id, a.user_id, u.full_name, a.action, a.entity_type, IFNULL(a.entity_id, ''), IFNULL(a.details, ''), a.created_at
		FROM activity_logs a
		JOIN users u ON a.user_id = u.id
		WHERE (a.entity_type = 'project' AND a.entity_id = ?)
		   OR (a.entity_type IN ('task', 'comment') AND a.entity_id IN (SELECT id FROM tasks WHERE project_id = ?))
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.DB.Query(query, projectID, projectID, limit, offset)
	if err != nil {
		return models.ActivityListResponse{Success: false, Message: "Gagal mengambil log aktivitas."}
	}
	defer rows.Close()

	var activities []models.ActivityLog
	layout := "2006-01-02 15:04:05"

	for rows.Next() {
		var a models.ActivityLog
		var createdAtStr string
		err := rows.Scan(&a.ID, &a.UserID, &a.UserName, &a.Action, &a.EntityType, &a.EntityID, &a.Details, &createdAtStr)
		if err != nil {
			return models.ActivityListResponse{Success: false, Message: "Gagal mem-parsing log aktivitas."}
		}
		if pt, err := time.Parse(layout, createdAtStr); err == nil {
			a.CreatedAt = pt
		}
		activities = append(activities, a)
	}

	return models.ActivityListResponse{
		Success:    true,
		Message:    "Log aktivitas berhasil dimuat.",
		Activities: activities,
	}
}

func (s *ActivityService) GetGlobalActivities(limit int, offset int) models.ActivityListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ActivityListResponse{Success: false, Message: "Unauthenticated"}
	}

	if limit <= 0 {
		limit = 50
	}

	// Fetch activities for projects the user is part of
	query := `
		SELECT a.id, a.user_id, u.full_name, a.action, a.entity_type, IFNULL(a.entity_id, ''), IFNULL(a.details, ''), a.created_at
		FROM activity_logs a
		JOIN users u ON a.user_id = u.id
		WHERE (a.entity_type = 'project' AND a.entity_id IN (SELECT project_id FROM project_members WHERE user_id = ?))
		   OR (a.entity_type IN ('task', 'comment') AND a.entity_id IN (SELECT id FROM tasks WHERE project_id IN (SELECT project_id FROM project_members WHERE user_id = ?)))
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.DB.Query(query, userID, userID, limit, offset)
	if err != nil {
		return models.ActivityListResponse{Success: false, Message: "Gagal mengambil log aktivitas."}
	}
	defer rows.Close()

	var activities []models.ActivityLog
	layout := "2006-01-02 15:04:05"

	for rows.Next() {
		var a models.ActivityLog
		var createdAtStr string
		err := rows.Scan(&a.ID, &a.UserID, &a.UserName, &a.Action, &a.EntityType, &a.EntityID, &a.Details, &createdAtStr)
		if err != nil {
			continue // Skip problematic rows instead of failing everything
		}
		if pt, err := time.Parse(layout, createdAtStr); err == nil {
			a.CreatedAt = pt
		}
		activities = append(activities, a)
	}

	return models.ActivityListResponse{
		Success:    true,
		Message:    "Log aktivitas berhasil dimuat.",
		Activities: activities,
	}
}
