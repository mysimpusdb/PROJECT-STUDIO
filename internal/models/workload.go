package models

type MemberWorkload struct {
	UserID             string  `json:"user_id"`
	UserName           string  `json:"user_name"`
	Date               string  `json:"date"` // YYYY-MM-DD
	AllocatedHours     float64 `json:"allocated_hours"`
	DailyCapacityHours float64 `json:"daily_capacity_hours"`
}

type WorkloadReportResponse struct {
	Success   bool             `json:"success"`
	Message   string           `json:"message"`
	Workloads []MemberWorkload `json:"workloads,omitempty"`
}

type UpdateCapacityRequest struct {
	ProjectID          string  `json:"project_id"`
	UserID             string  `json:"user_id"`
	DailyCapacityHours float64 `json:"daily_capacity_hours"`
}
