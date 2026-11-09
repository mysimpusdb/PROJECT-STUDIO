package models

import "time"

type Project struct {
	ID          string     `json:"id"`
	WorkspaceID string     `json:"workspace_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Color       string     `json:"color"`
	Status      string     `json:"status"` // e.g., 'active', 'archived', 'completed'
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	IsPublic    bool       `json:"is_public"`
	ShareToken  *string    `json:"share_token,omitempty"`
	TotalBudget float64    `json:"total_budget"`
	Currency    string     `json:"currency"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateProjectRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	Template    string  `json:"template,omitempty"`
	StartDate   *string `json:"start_date,omitempty"` // format: YYYY-MM-DD
	EndDate     *string `json:"end_date,omitempty"`   // format: YYYY-MM-DD
	TotalBudget float64 `json:"total_budget,omitempty"`
	Currency    string  `json:"currency,omitempty"`
}

type UpdateProjectRequest struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	Status      string  `json:"status"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
	TotalBudget float64 `json:"total_budget,omitempty"`
	Currency    string  `json:"currency,omitempty"`
}

type ProjectResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Project *Project `json:"project,omitempty"`
}

type ProjectListResponse struct {
	Success  bool       `json:"success"`
	Message  string     `json:"message"`
	Projects []Project  `json:"projects,omitempty"`
}

type ProjectMember struct {
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	DailyCapacityHours float64 `json:"daily_capacity_hours"`
	
	// Included for frontend convenience
	UserEmail    string  `json:"user_email,omitempty"`
	UserFullName string  `json:"user_full_name,omitempty"`
	UserAvatar   *string `json:"user_avatar,omitempty"`
}

type ProjectMemberRequest struct {
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"` // 'manager' or 'member'
}

type ProjectMemberListResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Members []ProjectMember `json:"members,omitempty"`
}

type ProjectMemberResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Member  *ProjectMember `json:"member,omitempty"`
}

type UpcomingDeadlineItem struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	ProjectName string     `json:"project_name"`
	Name        string     `json:"name"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
}

type DashboardSummaryResponse struct {
	Success        bool                   `json:"success"`
	Message        string                 `json:"message"`
	TotalProjects  int                    `json:"total_projects"`
	TotalTasks     int                    `json:"total_tasks"`
	OverdueTasks   int                    `json:"overdue_tasks"`
	CompletedToday int                    `json:"completed_today"`
	UpcomingTasks  []UpcomingDeadlineItem `json:"upcoming_tasks"`
}

type SearchItem struct {
	Type      string `json:"type"` // "project" or "task"
	ID        string `json:"id"`
	ProjectID string `json:"project_id,omitempty"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle,omitempty"`
	Status    string `json:"status,omitempty"`
	Color     string `json:"color,omitempty"`
}

type GlobalSearchResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Items   []SearchItem `json:"items"`
}
