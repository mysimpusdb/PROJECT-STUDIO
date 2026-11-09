package models

import "time"

type Task struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	ParentID    *string    `json:"parent_id,omitempty"` // For WBS subtasks
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`   // 'todo', 'in-progress', 'done'
	Priority    string     `json:"priority"` // 'low', 'medium', 'high', 'urgent'
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Duration    int        `json:"duration"`       // Number of days
	EstimatedHours int     `json:"estimated_hours"` // Number of hours expected
	Progress    int        `json:"progress"`       // 0-100
	AssigneeID  *string    `json:"assignee_id,omitempty"`
	SortOrder   int        `json:"sort_order"` // For drag & drop reordering
	Dependencies string    `json:"dependencies"` // Comma-separated Task IDs
	IsMilestone  bool      `json:"is_milestone"`
	Labels       string    `json:"labels,omitempty"`
	Checklists   string    `json:"checklists,omitempty"`
	EstimatedCost float64  `json:"estimated_cost"`
	ActualCost    float64  `json:"actual_cost"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	ProjectID   string  `json:"project_id"`
	ParentID    *string `json:"parent_id,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
	Duration    int     `json:"duration"`
	EstimatedHours int  `json:"estimated_hours"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	Labels      string  `json:"labels,omitempty"`
	Checklists  string  `json:"checklists,omitempty"`
	EstimatedCost float64 `json:"estimated_cost,omitempty"`
}

// UpdateTaskRequest is a partial update: any field that is omitted (nil or the
// empty string for Name/Status/Priority, 0 for Duration) keeps its stored value.
type UpdateTaskRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Duration int    `json:"duration"`
	// Dates: nil = keep existing value, "" = clear, "YYYY-MM-DD" = set.
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	// AssigneeID: nil = keep existing assignee, "" = unassign, otherwise the user ID.
	AssigneeID *string `json:"assignee_id,omitempty"`
	// nil = keep the value currently stored in DB.
	Description    *string  `json:"description,omitempty"`
	Progress       *int     `json:"progress,omitempty"`
	SortOrder      *int     `json:"sort_order,omitempty"`
	EstimatedHours *int     `json:"estimated_hours,omitempty"`
	Dependencies   *string  `json:"dependencies,omitempty"`
	IsMilestone    *bool    `json:"is_milestone,omitempty"`
	Labels         *string  `json:"labels,omitempty"`
	Checklists     *string  `json:"checklists,omitempty"`
	EstimatedCost  *float64 `json:"estimated_cost,omitempty"`
	ActualCost     *float64 `json:"actual_cost,omitempty"`
}

type TaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Task    *Task  `json:"task,omitempty"`
}

type TaskListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Tasks   []Task `json:"tasks,omitempty"`
}

type TaskComment struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"` // For display
	Content   string    `json:"content"`
	Reactions string    `json:"reactions,omitempty"` // JSON string
	CreatedAt time.Time `json:"created_at"`
}

type UpdateCommentRequest struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type CreateCommentRequest struct {
	TaskID  string `json:"task_id"`
	Content string `json:"content"`
}

type CommentResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Comment *TaskComment `json:"comment,omitempty"`
}

type CommentListResponse struct {
	Success  bool          `json:"success"`
	Message  string        `json:"message"`
	Comments []TaskComment `json:"comments,omitempty"`
}
