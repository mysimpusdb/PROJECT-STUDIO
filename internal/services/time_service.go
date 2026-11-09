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

type TimeTrackingService struct {
	ctx         context.Context
	activitySvc *ActivityService
}

func NewTimeTrackingService(activitySvc *ActivityService) *TimeTrackingService {
	return &TimeTrackingService{
		activitySvc: activitySvc,
	}
}

func (s *TimeTrackingService) Startup(ctx context.Context) {
	s.ctx = ctx
}

// emit sends a realtime event to the frontend when a Wails context is available.
func (s *TimeTrackingService) emit(event string, data ...interface{}) {
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, event, data...)
	}
}

// validateTimeLogInput returns a user-facing error message, or "" when valid.
func validateTimeLogInput(loggedDate string, hours float64) string {
	if hours <= 0 || hours > 24 {
		return "Jumlah jam harus lebih dari 0 dan maksimal 24."
	}
	if _, err := time.Parse(dateLayout, loggedDate); err != nil {
		return "Format tanggal tidak valid (gunakan YYYY-MM-DD)."
	}
	return ""
}

// LogTime adds a new time log entry
func (s *TimeTrackingService) LogTime(req models.CreateTimeLogRequest) models.TimeLogResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TimeLogResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTask(userID, req.TaskID, permWrite); err != nil {
		return models.TimeLogResponse{Success: false, Message: accessMessage(err, "Tugas tidak ditemukan.")}
	}
	if msg := validateTimeLogInput(req.LoggedDate, req.Hours); msg != "" {
		return models.TimeLogResponse{Success: false, Message: msg}
	}

	logEntry := models.TimeLog{
		ID:          uuid.New().String(),
		TaskID:      req.TaskID,
		UserID:      userID,
		LoggedDate:  req.LoggedDate,
		Hours:       req.Hours,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	query := `
		INSERT INTO time_logs (id, task_id, user_id, logged_date, hours, description, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.DB.Exec(query, logEntry.ID, logEntry.TaskID, logEntry.UserID, logEntry.LoggedDate, logEntry.Hours, logEntry.Description, logEntry.CreatedAt)
	if err != nil {
		return models.TimeLogResponse{Success: false, Message: "Gagal menyimpan time log."}
	}

	// Fetch username for response
	var userName string
	db.DB.QueryRow("SELECT full_name FROM users WHERE id = ?", userID).Scan(&userName)
	logEntry.UserName = userName

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "created", "task", req.TaskID, "Mencatat waktu kerja")
	}

	s.emit("rt:timelog:"+req.TaskID, "add", logEntry)

	return models.TimeLogResponse{
		Success: true,
		Message: "Waktu kerja berhasil dicatat.",
		TimeLog: &logEntry,
	}
}

// GetTaskTimeLogs fetches time logs for a specific task
func (s *TimeTrackingService) GetTaskTimeLogs(taskID string) models.TimeLogListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TimeLogListResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTask(userID, taskID, permRead); err != nil {
		return models.TimeLogListResponse{Success: false, Message: accessMessage(err, "Tugas tidak ditemukan.")}
	}

	query := `
		SELECT t.id, t.task_id, t.user_id, u.full_name, t.logged_date, t.hours, t.description, t.created_at
		FROM time_logs t
		JOIN users u ON t.user_id = u.id
		WHERE t.task_id = ?
		ORDER BY t.created_at DESC
	`
	rows, err := db.DB.Query(query, taskID)
	if err != nil {
		return models.TimeLogListResponse{Success: false, Message: "Gagal mengambil data time log."}
	}
	defer rows.Close()

	var logs []models.TimeLog
	layout := "2006-01-02 15:04:05"

	for rows.Next() {
		var l models.TimeLog
		var createdAtStr string
		err := rows.Scan(&l.ID, &l.TaskID, &l.UserID, &l.UserName, &l.LoggedDate, &l.Hours, &l.Description, &createdAtStr)
		if err != nil {
			return models.TimeLogListResponse{Success: false, Message: "Gagal mem-parsing time log."}
		}
		if pt, err := time.Parse(layout, createdAtStr); err == nil {
			l.CreatedAt = pt
		}
		logs = append(logs, l)
	}

	return models.TimeLogListResponse{
		Success:  true,
		Message:  "Time log berhasil dimuat.",
		TimeLogs: logs,
	}
}

// UpdateTimeLog updates an existing time log
func (s *TimeTrackingService) UpdateTimeLog(req models.UpdateTimeLogRequest) models.TimeLogResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TimeLogResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTimeLog(userID, req.ID, permWrite); err != nil {
		return models.TimeLogResponse{Success: false, Message: accessMessage(err, "Rekaman tidak ditemukan.")}
	}
	if msg := validateTimeLogInput(req.LoggedDate, req.Hours); msg != "" {
		return models.TimeLogResponse{Success: false, Message: msg}
	}

	// Verify ownership
	var ownerID, taskID string
	err = db.DB.QueryRow("SELECT user_id, task_id FROM time_logs WHERE id = ?", req.ID).Scan(&ownerID, &taskID)
	if err != nil || ownerID != userID {
		return models.TimeLogResponse{Success: false, Message: "Tidak memiliki akses untuk mengubah rekaman ini."}
	}

	query := `UPDATE time_logs SET logged_date = ?, hours = ?, description = ? WHERE id = ?`
	_, err = db.DB.Exec(query, req.LoggedDate, req.Hours, req.Description, req.ID)
	if err != nil {
		return models.TimeLogResponse{Success: false, Message: "Gagal memperbarui time log."}
	}

	s.emit("rt:timelog:"+taskID, "update", req)

	return models.TimeLogResponse{
		Success: true,
		Message: "Time log diperbarui.",
	}
}

// DeleteTimeLog removes a time log
func (s *TimeTrackingService) DeleteTimeLog(id string) models.TimeLogResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TimeLogResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTimeLog(userID, id, permWrite); err != nil {
		return models.TimeLogResponse{Success: false, Message: accessMessage(err, "Rekaman tidak ditemukan.")}
	}

	var ownerID, taskID string
	err = db.DB.QueryRow("SELECT user_id, task_id FROM time_logs WHERE id = ?", id).Scan(&ownerID, &taskID)
	if err != nil || ownerID != userID {
		return models.TimeLogResponse{Success: false, Message: "Tidak memiliki akses untuk menghapus rekaman ini."}
	}

	_, err = db.DB.Exec("DELETE FROM time_logs WHERE id = ?", id)
	if err != nil {
		return models.TimeLogResponse{Success: false, Message: "Gagal menghapus time log."}
	}

	s.emit("rt:timelog:"+taskID, "delete", id)

	return models.TimeLogResponse{
		Success: true,
		Message: "Time log dihapus.",
	}
}
