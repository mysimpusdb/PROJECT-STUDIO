package services

import (
	"context"
	"os"
	"strings"
	"time"

	"project-studio/internal/auth"
	"project-studio/internal/db"
	"project-studio/internal/models"
	"project-studio/internal/repository"

	"github.com/google/uuid"
)

type ProjectService struct {
	ctx             context.Context
	projectRepo       *repository.ProjectRepository
	projectMemberRepo *repository.ProjectMemberRepository
	activitySvc       *ActivityService
}

func NewProjectService(activitySvc *ActivityService) *ProjectService {
	return &ProjectService{
		projectRepo:       repository.NewProjectRepository(db.DB),
		projectMemberRepo: repository.NewProjectMemberRepository(db.DB),
		activitySvc:       activitySvc,
	}
}

func (s *ProjectService) Startup(ctx context.Context) {
	s.ctx = ctx
	s.projectRepo = repository.NewProjectRepository(db.DB)
	s.projectMemberRepo = repository.NewProjectMemberRepository(db.DB)
}

// CreateProject makes a new project in the currently active workspace
func (s *ProjectService) CreateProject(req models.CreateProjectRequest) models.ProjectResponse {
	// 1. Ensure user is logged in
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	// 2. Fetch the default workspace for this user
	// For now, we will query the first workspace the user belongs to
	var workspaceID string
	err = db.DB.QueryRow(`
		SELECT workspace_id FROM workspace_members WHERE user_id = ? LIMIT 1
	`, userID).Scan(&workspaceID)
	
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Workspace tidak ditemukan."}
	}

	// 3. Validate input and create the project object
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return models.ProjectResponse{Success: false, Message: "Nama proyek wajib diisi."}
	}
	if req.TotalBudget < 0 {
		return models.ProjectResponse{Success: false, Message: "Total budget tidak boleh negatif."}
	}
	startDate, _, err := parseDateInput(req.StartDate)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Tanggal mulai: " + err.Error()}
	}
	endDate, _, err := parseDateInput(req.EndDate)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Tanggal selesai: " + err.Error()}
	}
	if err := checkDateOrder(startDate, endDate); err != nil {
		return models.ProjectResponse{Success: false, Message: "Tanggal tidak valid: " + err.Error() + "."}
	}

	project := models.Project{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Name:        name,
		Description: req.Description,
		Color:       req.Color,
		Status:      "active",
		StartDate:   startDate,
		EndDate:     endDate,
		TotalBudget: req.TotalBudget,
		Currency:    req.Currency,
		CreatedBy:   userID,
	}

	// 4. Save to DB
	err = s.projectRepo.Create(&project)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Gagal membuat proyek: " + err.Error()}
	}

	// 5. Add current user as project member (owner)
	_, err = db.DB.Exec(`
		INSERT INTO project_members (project_id, user_id, role)
		VALUES (?, ?, 'owner')
	`, project.ID, userID)
	if err != nil {
		// Log the error but don't fail the whole request
	}

	// 6. Inject template tasks if requested
	if req.Template != "" && req.Template != "blank" {
		s.injectTemplateTasks(project.ID, req.Template, project.StartDate)
	}

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "created", "project", project.ID, "Proyek dibuat: "+project.Name)
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Proyek berhasil dibuat.",
		Project: &project,
	}
}

// GetProjects returns all projects for the current user's workspace
func (s *ProjectService) GetProjects() models.ProjectListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectListResponse{Success: false, Message: "Unauthenticated"}
	}

	var workspaceID, workspaceRole string
	err = db.DB.QueryRow(`
		SELECT workspace_id, role FROM workspace_members WHERE user_id = ? LIMIT 1
	`, userID).Scan(&workspaceID, &workspaceRole)

	if err != nil {
		return models.ProjectListResponse{Success: false, Message: "Workspace tidak ditemukan."}
	}

	// Workspace admins see every project; other users only see the projects
	// they created or were added to.
	var projects []models.Project
	if workspaceRole == "admin" {
		projects, err = s.projectRepo.GetByWorkspaceID(workspaceID)
	} else {
		projects, err = s.projectRepo.GetAccessibleByWorkspaceID(workspaceID, userID)
	}
	if err != nil {
		return models.ProjectListResponse{Success: false, Message: "Gagal memuat proyek."}
	}

	return models.ProjectListResponse{
		Success:  true,
		Message:  "Berhasil memuat proyek.",
		Projects: projects,
	}
}

// GetProjectByID retrieves a single project
func (s *ProjectService) GetProjectByID(id string) models.ProjectResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, id, permRead); err != nil {
		return models.ProjectResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	project, err := s.projectRepo.GetByID(id)
	if err != nil || project == nil {
		return models.ProjectResponse{Success: false, Message: "Proyek tidak ditemukan."}
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Berhasil memuat proyek.",
		Project: project,
	}
}

// UpdateProject modifies an existing project
func (s *ProjectService) UpdateProject(req models.UpdateProjectRequest) models.ProjectResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, req.ID, permManage); err != nil {
		return models.ProjectResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	project, err := s.projectRepo.GetByID(req.ID)
	if err != nil || project == nil {
		return models.ProjectResponse{Success: false, Message: "Proyek tidak ditemukan."}
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return models.ProjectResponse{Success: false, Message: "Nama proyek wajib diisi."}
	}
	project.Name = name
	project.Description = req.Description
	project.Color = req.Color

	if strings.TrimSpace(req.Status) != "" {
		status, ok := normalizeProjectStatus(req.Status)
		if !ok {
			return models.ProjectResponse{Success: false, Message: "Status proyek tidak valid."}
		}
		project.Status = status
	}
	if req.TotalBudget < 0 {
		return models.ProjectResponse{Success: false, Message: "Total budget tidak boleh negatif."}
	}
	project.TotalBudget = req.TotalBudget
	if strings.TrimSpace(req.Currency) != "" {
		project.Currency = strings.TrimSpace(req.Currency)
	}

	if start, set, err := parseDateInput(req.StartDate); err != nil {
		return models.ProjectResponse{Success: false, Message: "Tanggal mulai: " + err.Error()}
	} else if set {
		project.StartDate = start
	}
	if end, set, err := parseDateInput(req.EndDate); err != nil {
		return models.ProjectResponse{Success: false, Message: "Tanggal selesai: " + err.Error()}
	} else if set {
		project.EndDate = end
	}
	if err := checkDateOrder(project.StartDate, project.EndDate); err != nil {
		return models.ProjectResponse{Success: false, Message: "Tanggal tidak valid: " + err.Error() + "."}
	}

	err = s.projectRepo.Update(project)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Gagal memperbarui proyek."}
	}

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "updated", "project", project.ID, "Proyek diperbarui: "+project.Name)
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Proyek berhasil diperbarui.",
		Project: project,
	}
}

// ArchiveProject sets a project status to 'archived'
func (s *ProjectService) ArchiveProject(id string) models.ProjectResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, id, permManage); err != nil {
		return models.ProjectResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	project, err := s.projectRepo.GetByID(id)
	if err != nil || project == nil {
		return models.ProjectResponse{Success: false, Message: "Proyek tidak ditemukan."}
	}

	project.Status = "archived"
	err = s.projectRepo.Update(project)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Gagal mengarsipkan proyek."}
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Proyek berhasil diarsipkan.",
		Project: project,
	}
}

// DeleteProject permanently removes a project
func (s *ProjectService) DeleteProject(id string) models.ProjectResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, id, permManage); err != nil {
		return models.ProjectResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	files, err := s.projectRepo.Delete(id)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Gagal menghapus proyek."}
	}
	for _, f := range files {
		if isInsideAttachmentsDir(f) {
			_ = os.Remove(f)
		}
	}

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "deleted", "project", id, "Proyek dihapus")
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Proyek berhasil dihapus.",
	}
}

// DuplicateProject creates a copy of an existing project and its tasks
func (s *ProjectService) DuplicateProject(id string) models.ProjectResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, id, permRead); err != nil {
		return models.ProjectResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	// 1. Get original project
	originalProject, err := s.projectRepo.GetByID(id)
	if err != nil || originalProject == nil {
		return models.ProjectResponse{Success: false, Message: "Proyek tidak ditemukan."}
	}

	// 2. Create new project
	newProject := *originalProject
	newProject.ID = uuid.New().String()
	newProject.Name = originalProject.Name + " (Copy)"
	newProject.Status = "active"
	newProject.IsPublic = false
	newProject.ShareToken = nil
	newProject.CreatedBy = userID

	err = s.projectRepo.Create(&newProject)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Gagal menduplikasi proyek: " + err.Error()}
	}

	// 3. Add current user as project member (owner)
	_, err = db.DB.Exec(`
		INSERT INTO project_members (project_id, user_id, role)
		VALUES (?, ?, 'owner')
	`, newProject.ID, userID)
	if err != nil {
		// Log the error but don't fail the whole request
	}

	// 4. Duplicate tasks in two passes so parent IDs and dependency IDs can
	// always be remapped regardless of row order. We sort parents before
	// children so foreign key constraints on parent_id are satisfied.
	taskRepo := repository.NewTaskRepository(db.DB)
	originalTasks, err := taskRepo.GetByProjectID(originalProject.ID)
	if err == nil && len(originalTasks) > 0 {
		idMap := make(map[string]string, len(originalTasks))
		for _, oldTask := range originalTasks {
			idMap[oldTask.ID] = uuid.New().String()
		}

		inserted := make(map[string]bool, len(originalTasks))
		remaining := originalTasks
		for len(remaining) > 0 {
			var deferred []models.Task
			progressed := false
			for _, oldTask := range remaining {
				if oldTask.ParentID != nil && *oldTask.ParentID != "" {
					if _, existsInProject := idMap[*oldTask.ParentID]; existsInProject && !inserted[*oldTask.ParentID] {
						deferred = append(deferred, oldTask)
						continue
					}
				}

				newTask := oldTask
				newTask.ID = idMap[oldTask.ID]
				newTask.ProjectID = newProject.ID
				newTask.ParentID = nil
				if oldTask.ParentID != nil && *oldTask.ParentID != "" {
					if newParentID, exists := idMap[*oldTask.ParentID]; exists {
						newTask.ParentID = &newParentID
					}
				}
				newTask.Dependencies = remapIDList(oldTask.Dependencies, idMap)

				_ = taskRepo.Create(&newTask)
				inserted[oldTask.ID] = true
				progressed = true
			}
			if !progressed {
				// Break out if the source tasks have a parent cycle.
				for _, oldTask := range deferred {
					newTask := oldTask
					newTask.ID = idMap[oldTask.ID]
					newTask.ProjectID = newProject.ID
					newTask.ParentID = nil
					newTask.Dependencies = remapIDList(oldTask.Dependencies, idMap)
					_ = taskRepo.Create(&newTask)
				}
				break
			}
			remaining = deferred
		}
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Proyek berhasil diduplikasi.",
		Project: &newProject,
	}
}

func remapIDList(csv string, idMap map[string]string) string {
	if strings.TrimSpace(csv) == "" {
		return ""
	}
	var out []string
	for _, part := range strings.Split(csv, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if mapped, ok := idMap[part]; ok {
			out = append(out, mapped)
		}
	}
	return strings.Join(out, ", ")
}

// GetProjectMembers retrieves all members of a project
func (s *ProjectService) GetProjectMembers(projectID string) models.ProjectMemberListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectMemberListResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, projectID, permRead); err != nil {
		return models.ProjectMemberListResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	members, err := s.projectMemberRepo.GetByProjectID(projectID)
	if err != nil {
		return models.ProjectMemberListResponse{Success: false, Message: "Gagal memuat anggota proyek."}
	}

	return models.ProjectMemberListResponse{
		Success: true,
		Message: "Berhasil memuat anggota.",
		Members: members,
	}
}

// AddProjectMember adds a user to a project
func (s *ProjectService) AddProjectMember(req models.ProjectMemberRequest) models.ProjectMemberResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectMemberResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, req.ProjectID, permManage); err != nil {
		return models.ProjectMemberResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	// Validate role
	if req.Role != roleManager && req.Role != roleMember {
		return models.ProjectMemberResponse{Success: false, Message: "Peran tidak valid."}
	}

	// The invited user must belong to the project's workspace.
	if ok, err := isWorkspaceMemberOfProject(req.UserID, req.ProjectID); err != nil || !ok {
		return models.ProjectMemberResponse{Success: false, Message: "Pengguna bukan anggota workspace proyek ini."}
	}

	err = s.projectMemberRepo.AddMember(req.ProjectID, req.UserID, req.Role)
	if err != nil {
		return models.ProjectMemberResponse{Success: false, Message: "Gagal menambahkan anggota. Mungkin sudah menjadi anggota."}
	}

	// Try to return the newly added member details
	var member models.ProjectMember
	members, _ := s.projectMemberRepo.GetByProjectID(req.ProjectID)
	for _, m := range members {
		if m.UserID == req.UserID {
			member = m
			break
		}
	}

	return models.ProjectMemberResponse{
		Success: true,
		Message: "Anggota berhasil ditambahkan.",
		Member:  &member,
	}
}

// UpdateProjectMemberRole updates the role of a project member
func (s *ProjectService) UpdateProjectMemberRole(req models.ProjectMemberRequest) models.ProjectMemberResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectMemberResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, req.ProjectID, permManage); err != nil {
		return models.ProjectMemberResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	// Validate role
	if req.Role != roleManager && req.Role != roleMember {
		return models.ProjectMemberResponse{Success: false, Message: "Peran tidak valid."}
	}

	// The project owner's role cannot be changed.
	if s.isProjectOwner(req.ProjectID, req.UserID) {
		return models.ProjectMemberResponse{Success: false, Message: "Peran pemilik proyek tidak dapat diubah."}
	}

	err = s.projectMemberRepo.UpdateRole(req.ProjectID, req.UserID, req.Role)
	if err != nil {
		return models.ProjectMemberResponse{Success: false, Message: "Gagal memperbarui peran anggota."}
	}

	return models.ProjectMemberResponse{
		Success: true,
		Message: "Peran anggota berhasil diperbarui.",
	}
}

// RemoveProjectMember removes a user from a project
func (s *ProjectService) RemoveProjectMember(projectID string, memberUserID string) models.ProjectMemberResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectMemberResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, projectID, permManage); err != nil {
		return models.ProjectMemberResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	// The project owner cannot be removed.
	if s.isProjectOwner(projectID, memberUserID) {
		return models.ProjectMemberResponse{Success: false, Message: "Pemilik proyek tidak dapat dihapus."}
	}

	err = s.projectMemberRepo.RemoveMember(projectID, memberUserID)
	if err != nil {
		return models.ProjectMemberResponse{Success: false, Message: "Gagal menghapus anggota."}
	}

	return models.ProjectMemberResponse{
		Success: true,
		Message: "Anggota berhasil dihapus.",
	}
}

// isProjectOwner reports whether userID holds the 'owner' role in the project.
func (s *ProjectService) isProjectOwner(projectID, userID string) bool {
	var role string
	err := db.DB.QueryRow(
		`SELECT role FROM project_members WHERE project_id = ? AND user_id = ?`,
		projectID, userID,
	).Scan(&role)
	return err == nil && role == roleOwner
}

// GetWorkspaceUsers gets all users in the current user's workspace
func (s *ProjectService) GetWorkspaceUsers() models.UserListResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.UserListResponse{Success: false, Message: "Unauthenticated"}
	}

	var workspaceID string
	err = db.DB.QueryRow(`
		SELECT workspace_id FROM workspace_members WHERE user_id = ? LIMIT 1
	`, userID).Scan(&workspaceID)
	
	if err != nil {
		return models.UserListResponse{Success: false, Message: "Workspace tidak ditemukan."}
	}

	users, err := s.projectMemberRepo.GetWorkspaceUsers(workspaceID)
	if err != nil {
		return models.UserListResponse{Success: false, Message: "Gagal memuat pengguna workspace."}
	}

	return models.UserListResponse{
		Success: true,
		Message: "Berhasil memuat pengguna.",
		Users:   users,
	}
}

// UpdateMemberCapacity updates the daily capacity hours for a project member
func (s *ProjectService) UpdateMemberCapacity(req models.UpdateCapacityRequest) models.ProjectResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, req.ProjectID, permManage); err != nil {
		return models.ProjectResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}
	if req.DailyCapacityHours < 0 || req.DailyCapacityHours > 24 {
		return models.ProjectResponse{Success: false, Message: "Kapasitas harian harus antara 0 dan 24 jam."}
	}

	query := "UPDATE project_members SET daily_capacity_hours = ? WHERE project_id = ? AND user_id = ?"
	_, err = db.DB.Exec(query, req.DailyCapacityHours, req.ProjectID, req.UserID)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Gagal memperbarui kapasitas harian."}
	}

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "updated", "project_member", req.UserID, "Memperbarui kapasitas anggota")
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Kapasitas harian berhasil diperbarui.",
	}
}

// injectTemplateTasks creates a set of preset tasks based on the template
func (s *ProjectService) injectTemplateTasks(projectID string, template string, projectStartDate *time.Time) {
	taskRepo := repository.NewTaskRepository(db.DB)
	
	baseDate := time.Now()
	if projectStartDate != nil {
		baseDate = *projectStartDate
	}

	type templateTask struct {
		id          string
		name        string
		duration    int
		daysOffset  int
		deps        string
		isMilestone bool
	}

	var tTasks []templateTask

	switch template {
	case "software-development":
		tTasks = []templateTask{
			{"t1", "Requirement Analysis", 3, 0, "", false},
			{"t2", "System Design", 4, 3, "t1", false},
			{"t3", "Phase 1 Sign-off", 1, 7, "t2", true},
			{"t4", "Frontend Development", 7, 7, "t2", false},
			{"t5", "Backend Development", 7, 7, "t2", false},
			{"t6", "Integration & Testing", 5, 14, "t4, t5", false},
			{"t7", "Deployment", 2, 19, "t6", true},
		}
	case "marketing-campaign":
		tTasks = []templateTask{
			{"t1", "Market Research", 4, 0, "", false},
			{"t2", "Campaign Strategy", 3, 4, "t1", false},
			{"t3", "Content Creation", 5, 7, "t2", false},
			{"t4", "Ad Setup & Targeting", 3, 12, "t3", false},
			{"t5", "Campaign Launch", 1, 15, "t4", true},
		}
	case "event-planning":
		tTasks = []templateTask{
			{"t1", "Venue Selection", 5, 0, "", false},
			{"t2", "Vendor Booking", 4, 5, "t1", false},
			{"t3", "Invitations Sent", 1, 9, "t2", true},
			{"t4", "Catering Finalization", 3, 9, "t2", false},
			{"t5", "Event Day Execution", 1, 12, "t4", true},
		}
	case "construction-infrastructure":
		tTasks = []templateTask{
			{"t1", "Site Preparation", 5, 0, "", false},
			{"t2", "Foundation Works", 10, 5, "t1", false},
			{"t3", "Foundation Complete", 1, 15, "t2", true},
			{"t4", "Framing & Structure", 14, 15, "t2", false},
			{"t5", "Roofing & Siding", 7, 29, "t4", false},
			{"t6", "Interior Finishing", 10, 36, "t5", false},
			{"t7", "Final Inspection", 2, 46, "t6", true},
		}
	case "general-project", "general":
		tTasks = []templateTask{
			{"t1", "Project Kickoff & Planning", 2, 0, "", false},
			{"t2", "Execution Phase 1", 5, 2, "t1", false},
			{"t3", "Mid-Project Review", 1, 7, "t2", true},
			{"t4", "Execution Phase 2", 5, 8, "t3", false},
			{"t5", "Final Review & Handover", 2, 13, "t4", true},
		}
	default:
		return // Unknown template, do nothing
	}

	// Map generic template IDs to real UUIDs for dependency tracking
	idMap := make(map[string]string, len(tTasks))
	for _, tt := range tTasks {
		idMap[tt.id] = uuid.New().String()
	}

	for i, tt := range tTasks {
		start := baseDate.AddDate(0, 0, tt.daysOffset)
		end := start.AddDate(0, 0, tt.duration)
		deps := remapIDList(tt.deps, idMap)

		hours := tt.duration * 8
		if tt.isMilestone {
			hours = 0
		}

		task := models.Task{
			ID:             idMap[tt.id],
			ProjectID:      projectID,
			Name:           tt.name,
			Status:         "todo",
			Priority:       "medium",
			Duration:       tt.duration,
			EstimatedHours: hours,
			StartDate:      &start,
			EndDate:        &end,
			SortOrder:      i,
			Dependencies:   deps,
			IsMilestone:    tt.isMilestone,
		}
		_ = taskRepo.Create(&task)
	}
}

// GenerateShareLink creates a public token for sharing
func (s *ProjectService) GenerateShareLink(id string) models.ProjectResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, id, permManage); err != nil {
		return models.ProjectResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	project, err := s.projectRepo.GetByID(id)
	if err != nil || project == nil {
		return models.ProjectResponse{Success: false, Message: "Proyek tidak ditemukan."}
	}

	token := uuid.New().String()
	project.ShareToken = &token
	project.IsPublic = true

	err = s.projectRepo.Update(project)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Gagal membuat link tautan."}
	}

	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "updated", "project", project.ID, "Membuat tautan publik untuk proyek: "+project.Name)
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Tautan publik berhasil dibuat.",
		Project: project,
	}
}

// ToggleProjectPrivacy changes whether a project is accessible via public link
func (s *ProjectService) ToggleProjectPrivacy(id string, isPublic bool) models.ProjectResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.ProjectResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, id, permManage); err != nil {
		return models.ProjectResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	project, err := s.projectRepo.GetByID(id)
	if err != nil || project == nil {
		return models.ProjectResponse{Success: false, Message: "Proyek tidak ditemukan."}
	}

	project.IsPublic = isPublic
	err = s.projectRepo.Update(project)
	if err != nil {
		return models.ProjectResponse{Success: false, Message: "Gagal mengubah privasi proyek."}
	}

	statusStr := "Pribadi"
	if isPublic {
		statusStr = "Publik"
	}
	if s.activitySvc != nil {
		s.activitySvc.RecordActivity(userID, "updated", "project", project.ID, "Mengubah status privasi proyek menjadi: "+statusStr)
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Privasi proyek berhasil diubah.",
		Project: project,
	}
}

// GetPublicProject retrieves a project using only the share token, bypassing auth
func (s *ProjectService) GetPublicProject(token string) models.ProjectResponse {
	if token == "" {
		return models.ProjectResponse{Success: false, Message: "Token tautan tidak valid."}
	}

	project, err := s.projectRepo.GetByShareToken(token)
	if err != nil || project == nil {
		return models.ProjectResponse{Success: false, Message: "Proyek tidak ditemukan atau tautan sudah tidak aktif."}
	}

	return models.ProjectResponse{
		Success: true,
		Message: "Berhasil memuat proyek publik.",
		Project: project,
	}
}
