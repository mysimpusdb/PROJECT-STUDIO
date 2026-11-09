package services

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"time"

	"project-studio/internal/auth"
	"project-studio/internal/db"
	"project-studio/internal/models"
	"project-studio/internal/repository"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
)

// maxAttachmentSize is the largest file accepted as a task attachment.
const maxAttachmentSize = 100 << 20 // 100 MB

type TaskService struct {
	ctx          context.Context
	taskRepo     *repository.TaskRepository
	attRepo      *repository.AttachmentRepository
	notifService *NotificationService
	activitySvc  *ActivityService
}

func NewTaskService(notifService *NotificationService, activitySvc *ActivityService) *TaskService {
	return &TaskService{
		taskRepo:     repository.NewTaskRepository(db.DB),
		attRepo:      repository.NewAttachmentRepository(db.DB),
		notifService: notifService,
		activitySvc:  activitySvc,
	}
}

func (s *TaskService) Startup(ctx context.Context) {
	s.ctx = ctx
	s.taskRepo = repository.NewTaskRepository(db.DB)
	s.attRepo = repository.NewAttachmentRepository(db.DB)
}

// emit sends a realtime event to the frontend when a Wails context is available.
func (s *TaskService) emit(event string, data ...interface{}) {
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, event, data...)
	}
}

// validateAssignee ensures the assignee belongs to the project's workspace.
func validateAssignee(assigneeID, projectID string) error {
	ok, err := isWorkspaceMemberOfProject(assigneeID, projectID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("penerima tugas bukan anggota workspace proyek ini")
	}
	return nil
}

func (s *TaskService) CreateTask(req models.CreateTaskRequest) models.TaskResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TaskResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, req.ProjectID, permWrite); err != nil {
		return models.TaskResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return models.TaskResponse{Success: false, Message: "Nama tugas wajib diisi."}
	}

	status := "todo"
	if strings.TrimSpace(req.Status) != "" {
		var ok bool
		if status, ok = normalizeStatus(req.Status); !ok {
			return models.TaskResponse{Success: false, Message: "Status tugas tidak valid."}
		}
	}
	priority := "medium"
	if strings.TrimSpace(req.Priority) != "" {
		var ok bool
		if priority, ok = normalizePriority(req.Priority); !ok {
			return models.TaskResponse{Success: false, Message: "Prioritas tugas tidak valid."}
		}
	}
	if req.EstimatedHours < 0 || req.EstimatedCost < 0 {
		return models.TaskResponse{Success: false, Message: "Estimasi jam dan biaya tidak boleh negatif."}
	}

	startDate, _, err := parseDateInput(req.StartDate)
	if err != nil {
		return models.TaskResponse{Success: false, Message: "Tanggal mulai: " + err.Error()}
	}
	endDate, _, err := parseDateInput(req.EndDate)
	if err != nil {
		return models.TaskResponse{Success: false, Message: "Tanggal selesai: " + err.Error()}
	}
	if err := checkDateOrder(startDate, endDate); err != nil {
		return models.TaskResponse{Success: false, Message: "Tanggal tidak valid: " + err.Error() + "."}
	}

	// A parent task must live in the same project.
	if req.ParentID != nil && *req.ParentID != "" {
		parentProject, err := projectIDQuery(`SELECT project_id FROM tasks WHERE id = ?`, *req.ParentID)
		if err != nil || parentProject != req.ProjectID {
			return models.TaskResponse{Success: false, Message: "Tugas induk tidak ditemukan di proyek ini."}
		}
	}

	task := models.Task{
		ID:             uuid.New().String(),
		ProjectID:      req.ProjectID,
		ParentID:       req.ParentID,
		Name:           name,
		Description:    req.Description,
		Status:         status,
		Priority:       priority,
		Duration:       req.Duration,
		EstimatedHours: req.EstimatedHours,
		EstimatedCost:  req.EstimatedCost,
		Progress:       0,
		Labels:         req.Labels,
		Checklists:     req.Checklists,
		StartDate:      startDate,
		EndDate:        endDate,
	}
	if task.Duration <= 0 {
		task.Duration = 1
	}

	taskRepo := s.taskRepo
	if err := taskRepo.Create(&task); err != nil {
		return models.TaskResponse{Success: false, Message: "Gagal membuat tugas: " + err.Error()}
	}

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "created", "task", task.ID, "Dibuat tugas baru: "+task.Name)
	}

	return models.TaskResponse{
		Success: true,
		Message: "Tugas berhasil dibuat.",
		Task:    &task,
	}
}

func (s *TaskService) GetTasksByProject(projectID string) models.TaskListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TaskListResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, projectID, permRead); err != nil {
		return models.TaskListResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	taskRepo := s.taskRepo
	tasks, err := taskRepo.GetByProjectID(projectID)
	if err != nil {
		return models.TaskListResponse{Success: false, Message: "Gagal memuat tugas."}
	}

	return models.TaskListResponse{
		Success: true,
		Message: "Berhasil memuat tugas.",
		Tasks:   tasks,
	}
}

// UpdateTask applies a partial update: fields that are omitted keep the value
// currently stored for the task.
func (s *TaskService) UpdateTask(req models.UpdateTaskRequest) models.TaskResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TaskResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTask(userID, req.ID, permWrite); err != nil {
		return models.TaskResponse{Success: false, Message: accessMessage(err, "Tugas tidak ditemukan.")}
	}

	taskRepo := s.taskRepo
	existing, err := taskRepo.GetByID(req.ID)
	if err != nil {
		return models.TaskResponse{Success: false, Message: "Gagal memuat tugas."}
	}
	if existing == nil {
		return models.TaskResponse{Success: false, Message: "Tugas tidak ditemukan."}
	}

	task := *existing

	if name := strings.TrimSpace(req.Name); name != "" {
		task.Name = name
	}
	if strings.TrimSpace(req.Status) != "" {
		status, ok := normalizeStatus(req.Status)
		if !ok {
			return models.TaskResponse{Success: false, Message: "Status tugas tidak valid."}
		}
		task.Status = status
	}
	if strings.TrimSpace(req.Priority) != "" {
		priority, ok := normalizePriority(req.Priority)
		if !ok {
			return models.TaskResponse{Success: false, Message: "Prioritas tugas tidak valid."}
		}
		task.Priority = priority
	}
	if req.Duration > 0 {
		task.Duration = req.Duration
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Progress != nil {
		if *req.Progress < 0 || *req.Progress > 100 {
			return models.TaskResponse{Success: false, Message: "Progress harus antara 0 dan 100."}
		}
		task.Progress = *req.Progress
	}
	if req.SortOrder != nil {
		task.SortOrder = *req.SortOrder
	}
	if req.EstimatedHours != nil {
		if *req.EstimatedHours < 0 {
			return models.TaskResponse{Success: false, Message: "Estimasi jam tidak boleh negatif."}
		}
		task.EstimatedHours = *req.EstimatedHours
	}
	if req.Dependencies != nil {
		task.Dependencies = *req.Dependencies
	}
	if req.IsMilestone != nil {
		task.IsMilestone = *req.IsMilestone
	}
	if req.Labels != nil {
		task.Labels = *req.Labels
	}
	if req.Checklists != nil {
		task.Checklists = *req.Checklists
	}
	if req.EstimatedCost != nil {
		if *req.EstimatedCost < 0 {
			return models.TaskResponse{Success: false, Message: "Estimasi biaya tidak boleh negatif."}
		}
		task.EstimatedCost = *req.EstimatedCost
	}
	if req.ActualCost != nil {
		if *req.ActualCost < 0 {
			return models.TaskResponse{Success: false, Message: "Biaya aktual tidak boleh negatif."}
		}
		task.ActualCost = *req.ActualCost
	}

	if start, set, err := parseDateInput(req.StartDate); err != nil {
		return models.TaskResponse{Success: false, Message: "Tanggal mulai: " + err.Error()}
	} else if set {
		task.StartDate = start
	}
	if end, set, err := parseDateInput(req.EndDate); err != nil {
		return models.TaskResponse{Success: false, Message: "Tanggal selesai: " + err.Error()}
	} else if set {
		task.EndDate = end
	}
	if err := checkDateOrder(task.StartDate, task.EndDate); err != nil {
		return models.TaskResponse{Success: false, Message: "Tanggal tidak valid: " + err.Error() + "."}
	}

	// Assignee: nil keeps, "" unassigns, otherwise must belong to the workspace.
	assigneeChanged := false
	if req.AssigneeID != nil {
		newAssignee := strings.TrimSpace(*req.AssigneeID)
		if newAssignee == "" {
			assigneeChanged = task.AssigneeID != nil
			task.AssigneeID = nil
		} else {
			if err := validateAssignee(newAssignee, task.ProjectID); err != nil {
				return models.TaskResponse{Success: false, Message: "Gagal menugaskan: " + err.Error() + "."}
			}
			assigneeChanged = task.AssigneeID == nil || *task.AssigneeID != newAssignee
			task.AssigneeID = &newAssignee
		}
	}

	if err := taskRepo.Update(&task); err != nil {
		return models.TaskResponse{Success: false, Message: "Gagal memperbarui tugas."}
	}

	if assigneeChanged && task.AssigneeID != nil && *task.AssigneeID != userID && s.notifService != nil {
		s.notifService.CreateNotification(
			*task.AssigneeID,
			"Tugas Baru Untuk Anda",
			"Tugas '"+task.Name+"' ditugaskan kepada Anda.",
			"task_assigned",
			task.ID,
		)
	}

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "updated", "task", task.ID, "Memperbarui tugas: "+task.Name)
	}

	s.emit("rt:task:"+task.ID, "update", task)

	return models.TaskResponse{
		Success: true,
		Message: "Tugas berhasil diperbarui.",
		Task:    &task,
	}
}

func (s *TaskService) DeleteTask(id string) models.TaskResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TaskResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTask(userID, id, permWrite); err != nil {
		return models.TaskResponse{Success: false, Message: accessMessage(err, "Tugas tidak ditemukan.")}
	}

	taskRepo := s.taskRepo
	files, err := taskRepo.Delete(id)
	if err != nil {
		return models.TaskResponse{Success: false, Message: "Gagal menghapus tugas."}
	}
	for _, f := range files {
		if isInsideAttachmentsDir(f) {
			_ = os.Remove(f)
		}
	}

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "deleted", "task", id, "Tugas dihapus")
	}

	return models.TaskResponse{
		Success: true,
		Message: "Tugas berhasil dihapus.",
	}
}

func (s *TaskService) GetTaskComments(taskID string) models.CommentListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.CommentListResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTask(userID, taskID, permRead); err != nil {
		return models.CommentListResponse{Success: false, Message: accessMessage(err, "Tugas tidak ditemukan.")}
	}

	query := `
		SELECT c.id, c.task_id, c.user_id, u.full_name, c.content, IFNULL(c.reactions, ''), c.created_at
		FROM task_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.task_id = ?
		ORDER BY c.created_at ASC
	`
	rows, err := db.DB.Query(query, taskID)
	if err != nil {
		return models.CommentListResponse{Success: false, Message: "Gagal memuat komentar."}
	}
	defer rows.Close()

	var comments []models.TaskComment
	layout := "2006-01-02 15:04:05"

	for rows.Next() {
		var c models.TaskComment
		var createdAtStr string
		err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.UserName, &c.Content, &c.Reactions, &createdAtStr)
		if err != nil {
			return models.CommentListResponse{Success: false, Message: "Gagal mem-parsing komentar."}
		}
		if pt, err := time.Parse(layout, createdAtStr); err == nil {
			c.CreatedAt = pt
		}
		comments = append(comments, c)
	}

	return models.CommentListResponse{
		Success:  true,
		Message:  "Komentar dimuat",
		Comments: comments,
	}
}

func (s *TaskService) AddComment(req models.CreateCommentRequest) models.CommentResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.CommentResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTask(userID, req.TaskID, permWrite); err != nil {
		return models.CommentResponse{Success: false, Message: accessMessage(err, "Tugas tidak ditemukan.")}
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return models.CommentResponse{Success: false, Message: "Komentar tidak boleh kosong."}
	}

	comment := models.TaskComment{
		ID:      uuid.New().String(),
		TaskID:  req.TaskID,
		UserID:  userID,
		Content: content,
	}

	query := `
		INSERT INTO task_comments (id, task_id, user_id, content)
		VALUES (?, ?, ?, ?)
	`
	_, err = db.DB.Exec(query, comment.ID, comment.TaskID, comment.UserID, comment.Content)
	if err != nil {
		return models.CommentResponse{Success: false, Message: "Gagal menambahkan komentar."}
	}

	// Fetch username for response
	var userName string
	db.DB.QueryRow("SELECT full_name FROM users WHERE id = ?", userID).Scan(&userName)
	comment.UserName = userName
	comment.CreatedAt = time.Now()

	s.emit("rt:comment:"+req.TaskID, "add", comment)

	// Send notification to task assignee (if it's not the comment author)
	var assigneeID *string
	var taskName string
	err = db.DB.QueryRow("SELECT assignee_id, name FROM tasks WHERE id = ?", req.TaskID).Scan(&assigneeID, &taskName)
	if err == nil && assigneeID != nil && *assigneeID != "" && *assigneeID != userID && s.notifService != nil {
		s.notifService.CreateNotification(
			*assigneeID,
			"Komentar Baru",
			userName+" mengomentari tugas: "+taskName,
			"comment",
			req.TaskID,
		)
	}

	return models.CommentResponse{
		Success: true,
		Message: "Komentar ditambahkan",
		Comment: &comment,
	}
}

func (s *TaskService) UpdateComment(req models.UpdateCommentRequest) models.CommentResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.CommentResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeComment(userID, req.ID, permWrite); err != nil {
		return models.CommentResponse{Success: false, Message: accessMessage(err, "Komentar tidak ditemukan.")}
	}

	var taskID, ownerID string
	err = db.DB.QueryRow("SELECT task_id, user_id FROM task_comments WHERE id = ?", req.ID).Scan(&taskID, &ownerID)
	if err != nil {
		return models.CommentResponse{Success: false, Message: "Komentar tidak ditemukan."}
	}

	if ownerID != userID {
		return models.CommentResponse{Success: false, Message: "Anda tidak dapat mengubah komentar ini."}
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return models.CommentResponse{Success: false, Message: "Komentar tidak boleh kosong."}
	}
	req.Content = content

	_, err = db.DB.Exec("UPDATE task_comments SET content = ? WHERE id = ?", content, req.ID)
	if err != nil {
		return models.CommentResponse{Success: false, Message: "Gagal memperbarui komentar."}
	}

	s.emit("rt:comment:"+taskID, "update", req)

	return models.CommentResponse{
		Success: true,
		Message: "Komentar berhasil diperbarui",
	}
}

func (s *TaskService) DeleteComment(id string) models.TaskResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TaskResponse{Success: false, Message: "Unauthenticated"}
	}

	role, err := func() (string, error) {
		projectID, err := projectIDQuery(`
			SELECT t.project_id FROM task_comments c
			JOIN tasks t ON t.id = c.task_id WHERE c.id = ?`, id)
		if err != nil {
			return "", err
		}
		return authorizeProject(userID, projectID, permWrite)
	}()
	if err != nil {
		return models.TaskResponse{Success: false, Message: accessMessage(err, "Komentar tidak ditemukan.")}
	}

	var taskID, ownerID string
	err = db.DB.QueryRow("SELECT task_id, user_id FROM task_comments WHERE id = ?", id).Scan(&taskID, &ownerID)
	if err != nil {
		return models.TaskResponse{Success: false, Message: "Komentar tidak ditemukan."}
	}

	// The author or a project owner/manager may delete a comment.
	if ownerID != userID && role != roleOwner && role != roleManager {
		return models.TaskResponse{Success: false, Message: "Anda tidak memiliki izin menghapus komentar ini."}
	}

	_, err = db.DB.Exec("DELETE FROM task_comments WHERE id = ?", id)
	if err != nil {
		return models.TaskResponse{Success: false, Message: "Gagal menghapus komentar."}
	}

	s.emit("rt:comment:"+taskID, "delete", id)

	return models.TaskResponse{
		Success: true,
		Message: "Komentar dihapus",
	}
}

// attachmentsDir is where uploaded attachments are stored.
func attachmentsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".project-studio", "attachments"), nil
}

// isInsideAttachmentsDir guards file operations against paths that were
// tampered with in the database and point outside the attachments folder.
func isInsideAttachmentsDir(path string) bool {
	dir, err := attachmentsDir()
	if err != nil || path == "" {
		return false
	}
	rel, err := filepath.Rel(dir, filepath.Clean(path))
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func (s *TaskService) SelectAndUploadAttachment(taskID string) models.AttachmentResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.AttachmentResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTask(userID, taskID, permWrite); err != nil {
		return models.AttachmentResponse{Success: false, Message: accessMessage(err, "Tugas tidak ditemukan.")}
	}

	if s.ctx == nil {
		return models.AttachmentResponse{Success: false, Message: "Aplikasi belum siap."}
	}

	selectedFilePath, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Pilih File Lampiran",
	})
	if err != nil || selectedFilePath == "" {
		return models.AttachmentResponse{Success: false, Message: "Tidak ada file yang dipilih."}
	}

	fileInfo, err := os.Stat(selectedFilePath)
	if err != nil || fileInfo.IsDir() {
		return models.AttachmentResponse{Success: false, Message: "Gagal membaca file."}
	}
	if fileInfo.Size() > maxAttachmentSize {
		return models.AttachmentResponse{Success: false, Message: "Ukuran file melebihi batas 100 MB."}
	}

	destDir, err := attachmentsDir()
	if err != nil {
		return models.AttachmentResponse{Success: false, Message: "Gagal menentukan folder lampiran."}
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return models.AttachmentResponse{Success: false, Message: "Gagal membuat folder lampiran."}
	}

	fileName := fileInfo.Name()
	fileID := uuid.New().String()
	destPath := filepath.Join(destDir, fileID+"_"+fileName)

	sourceFile, err := os.Open(selectedFilePath)
	if err != nil {
		return models.AttachmentResponse{Success: false, Message: "Gagal membuka file sumber."}
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return models.AttachmentResponse{Success: false, Message: "Gagal membuat file tujuan."}
	}

	_, copyErr := io.Copy(destFile, sourceFile)
	closeErr := destFile.Close()
	if copyErr != nil || closeErr != nil {
		os.Remove(destPath)
		return models.AttachmentResponse{Success: false, Message: "Gagal menyalin file."}
	}

	attachment := models.TaskAttachment{
		ID:         fileID,
		TaskID:     taskID,
		UserID:     userID,
		FileName:   fileName,
		FilePath:   destPath,
		FileSize:   fileInfo.Size(),
		UploadedAt: time.Now(),
	}

	attRepo := s.attRepo
	if err := attRepo.Create(&attachment); err != nil {
		os.Remove(destPath) // Rollback if DB insert fails
		return models.AttachmentResponse{Success: false, Message: "Gagal menyimpan metadata ke database."}
	}

	var userName string
	db.DB.QueryRow("SELECT full_name FROM users WHERE id = ?", userID).Scan(&userName)
	attachment.UserName = userName

	return models.AttachmentResponse{
		Success:    true,
		Message:    "File berhasil diunggah",
		Attachment: &attachment,
	}
}

func (s *TaskService) GetTaskAttachments(taskID string) models.AttachmentListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.AttachmentListResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeTask(userID, taskID, permRead); err != nil {
		return models.AttachmentListResponse{Success: false, Message: accessMessage(err, "Tugas tidak ditemukan.")}
	}

	attRepo := s.attRepo
	attachments, err := attRepo.GetByTaskID(taskID)
	if err != nil {
		return models.AttachmentListResponse{Success: false, Message: "Gagal memuat lampiran."}
	}

	return models.AttachmentListResponse{
		Success:     true,
		Message:     "Lampiran dimuat",
		Attachments: attachments,
	}
}

func (s *TaskService) DeleteAttachment(id string) models.AttachmentResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.AttachmentResponse{Success: false, Message: "Unauthenticated"}
	}

	role, err := authorizeAttachment(userID, id, permWrite)
	if err != nil {
		return models.AttachmentResponse{Success: false, Message: accessMessage(err, "Lampiran tidak ditemukan.")}
	}

	attRepo := s.attRepo
	att, err := attRepo.GetByID(id)
	if err != nil || att == nil {
		return models.AttachmentResponse{Success: false, Message: "Lampiran tidak ditemukan."}
	}

	// The uploader or a project owner/manager may delete an attachment.
	if att.UserID != userID && role != roleOwner && role != roleManager {
		return models.AttachmentResponse{Success: false, Message: forbiddenMessage}
	}

	if err := attRepo.Delete(id); err != nil {
		return models.AttachmentResponse{Success: false, Message: "Gagal menghapus dari database."}
	}

	if isInsideAttachmentsDir(att.FilePath) {
		os.Remove(att.FilePath)
	}

	return models.AttachmentResponse{
		Success: true,
		Message: "Lampiran berhasil dihapus.",
	}
}

// openWithSystem opens a file with the operating system's default application.
func openWithSystem(path string) error {
	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() { _ = cmd.Wait() }()
	return nil
}

func (s *TaskService) OpenAttachment(id string) models.AttachmentResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.AttachmentResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeAttachment(userID, id, permRead); err != nil {
		return models.AttachmentResponse{Success: false, Message: accessMessage(err, "Lampiran tidak ditemukan.")}
	}

	attRepo := s.attRepo
	att, err := attRepo.GetByID(id)
	if err != nil || att == nil {
		return models.AttachmentResponse{Success: false, Message: "Lampiran tidak ditemukan."}
	}

	if !isInsideAttachmentsDir(att.FilePath) {
		return models.AttachmentResponse{Success: false, Message: "Lokasi file lampiran tidak valid."}
	}
	if _, err := os.Stat(att.FilePath); err != nil {
		return models.AttachmentResponse{Success: false, Message: "File lampiran tidak ditemukan di disk."}
	}

	if err := openWithSystem(att.FilePath); err != nil {
		return models.AttachmentResponse{Success: false, Message: "Gagal membuka file: " + err.Error()}
	}

	return models.AttachmentResponse{
		Success: true,
		Message: "File dibuka.",
	}
}

// ExportProjectTasksToExcel exports the project's tasks to an Excel file and prompts the user to save it.
func (s *TaskService) ExportProjectTasksToExcel(projectID string) models.TaskResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.TaskResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, projectID, permRead); err != nil {
		return models.TaskResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	taskRepo := s.taskRepo
	tasks, err := taskRepo.GetByProjectID(projectID)
	if err != nil {
		return models.TaskResponse{Success: false, Message: "Gagal memuat tugas."}
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Tasks"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"ID", "Nama Tugas", "Deskripsi", "Status", "Prioritas", "Mulai", "Selesai", "Progress (%)", "Estimated Cost", "Actual Cost", "Milestone"}
	for i, header := range headers {
		f.SetCellValue(sheetName, getCell(i+1, 1), header)
	}

	for i, t := range tasks {
		row := i + 2
		startDate := ""
		if t.StartDate != nil {
			startDate = t.StartDate.Format("2006-01-02")
		}
		endDate := ""
		if t.EndDate != nil {
			endDate = t.EndDate.Format("2006-01-02")
		}

		isMilestone := "Tidak"
		if t.IsMilestone {
			isMilestone = "Ya"
		}

		f.SetCellValue(sheetName, getCell(1, row), t.ID)
		f.SetCellValue(sheetName, getCell(2, row), spreadsheetSafe(t.Name))
		f.SetCellValue(sheetName, getCell(3, row), spreadsheetSafe(t.Description))
		f.SetCellValue(sheetName, getCell(4, row), t.Status)
		f.SetCellValue(sheetName, getCell(5, row), t.Priority)
		f.SetCellValue(sheetName, getCell(6, row), startDate)
		f.SetCellValue(sheetName, getCell(7, row), endDate)
		f.SetCellValue(sheetName, getCell(8, row), t.Progress)
		f.SetCellValue(sheetName, getCell(9, row), t.EstimatedCost)
		f.SetCellValue(sheetName, getCell(10, row), t.ActualCost)
		f.SetCellValue(sheetName, getCell(11, row), isMilestone)
	}

	if s.ctx == nil {
		return models.TaskResponse{Success: false, Message: "Aplikasi belum siap."}
	}
	savePath, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		Title:           "Simpan Export Excel",
		DefaultFilename: "project_tasks.xlsx",
		Filters: []runtime.FileFilter{
			{DisplayName: "Excel Files (*.xlsx)", Pattern: "*.xlsx"},
		},
	})

	if err != nil || savePath == "" {
		return models.TaskResponse{Success: false, Message: "Export dibatalkan."}
	}

	if err := f.SaveAs(savePath); err != nil {
		return models.TaskResponse{Success: false, Message: "Gagal menyimpan file Excel: " + err.Error()}
	}

	return models.TaskResponse{
		Success: true,
		Message: "File Excel berhasil disimpan di " + savePath,
	}
}

func getCell(col, row int) string {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	return cell
}
