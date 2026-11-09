package services

import (
	"context"
	"time"

	"project-studio/internal/auth"
	"project-studio/internal/db"
	"project-studio/internal/models"
)

type WorkloadService struct {
	ctx context.Context
}

func NewWorkloadService() *WorkloadService {
	return &WorkloadService{}
}

func (s *WorkloadService) Startup(ctx context.Context) {
	s.ctx = ctx
}

// GetProjectWorkload calculates the allocated hours per day for members of a project
// between startDate and endDate (format YYYY-MM-DD).
func (s *WorkloadService) GetProjectWorkload(projectID, startDateStr, endDateStr string) models.WorkloadReportResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.WorkloadReportResponse{Success: false, Message: "Unauthenticated"}
	}

	if _, err := authorizeProject(userID, projectID, permRead); err != nil {
		return models.WorkloadReportResponse{Success: false, Message: accessMessage(err, "Proyek tidak ditemukan.")}
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return models.WorkloadReportResponse{Success: false, Message: "Format startDate tidak valid."}
	}
	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return models.WorkloadReportResponse{Success: false, Message: "Format endDate tidak valid."}
	}

	if endDate.Before(startDate) {
		return models.WorkloadReportResponse{Success: false, Message: "endDate harus setelah startDate."}
	}

	// 1. Ambil member dan kapasitas mereka
	memberQuery := `SELECT user_id, u.full_name, daily_capacity_hours 
					FROM project_members pm
					JOIN users u ON pm.user_id = u.id
					WHERE project_id = ?`
	rows, err := db.DB.Query(memberQuery, projectID)
	if err != nil {
		return models.WorkloadReportResponse{Success: false, Message: "Gagal mengambil data anggota proyek."}
	}
	
	type memberInfo struct {
		name     string
		capacity float64
	}
	members := make(map[string]memberInfo)
	for rows.Next() {
		var uid, name string
		var cap float64
		rows.Scan(&uid, &name, &cap)
		members[uid] = memberInfo{name: name, capacity: cap}
	}
	rows.Close()

	// 2. Ambil semua tasks yang beririsan dengan rentang waktu, dan hitung alokasi harian
	taskQuery := `
		SELECT assignee_id, start_date, duration, estimated_hours
		FROM tasks 
		WHERE project_id = ? 
		  AND assignee_id IS NOT NULL 
		  AND status != 'done' 
		  AND status != 'cancelled'
		  AND start_date IS NOT NULL
	`
	rows, err = db.DB.Query(taskQuery, projectID)
	if err != nil {
		return models.WorkloadReportResponse{Success: false, Message: "Gagal mengambil data tasks."}
	}
	defer rows.Close()

	// Map: UserID -> DateString -> AllocatedHours
	workloadMap := make(map[string]map[string]float64)

	for rows.Next() {
		var assigneeID string
		var tStartStr string
		var duration, estimatedHours int

		err := rows.Scan(&assigneeID, &tStartStr, &duration, &estimatedHours)
		if err != nil {
			continue
		}

		if _, exists := members[assigneeID]; !exists {
			continue // assignee is not part of project members anymore, ignore
		}

		if len(tStartStr) < 10 {
			continue
		}
		tStart, err := time.Parse(layout, tStartStr[:10])
		if err != nil {
			continue
		}

		if duration <= 0 {
			duration = 1 // Default at least 1 day
		}
		if estimatedHours <= 0 {
			// Default to 8h per day when the task has no explicit estimate.
			estimatedHours = duration * 8
		}

		hoursPerDay := float64(estimatedHours) / float64(duration)

		// Sebarkan hoursPerDay dari tStart sampai tStart + duration - 1 hari
		for i := 0; i < duration; i++ {
			currentDay := tStart.AddDate(0, 0, i)
			
			// Hanya pedulikan yang ada di dalam rentang query
			if currentDay.Before(startDate) || currentDay.After(endDate) {
				continue
			}

			dayStr := currentDay.Format(layout)
			if workloadMap[assigneeID] == nil {
				workloadMap[assigneeID] = make(map[string]float64)
			}
			workloadMap[assigneeID][dayStr] += hoursPerDay
		}
	}

	// 3. Susun respons (Flat array untuk UI)
	var result []models.MemberWorkload
	for uid, info := range members {
		// Populate every day in the range even if 0
		daysRange := int(endDate.Sub(startDate).Hours()/24) + 1
		for i := 0; i < daysRange; i++ {
			currentDay := startDate.AddDate(0, 0, i)
			dayStr := currentDay.Format(layout)
			
			allocated := 0.0
			if userMap, ok := workloadMap[uid]; ok {
				allocated = userMap[dayStr]
			}

			result = append(result, models.MemberWorkload{
				UserID:             uid,
				UserName:           info.name,
				Date:               dayStr,
				AllocatedHours:     allocated,
				DailyCapacityHours: info.capacity,
			})
		}
	}

	return models.WorkloadReportResponse{
		Success:   true,
		Message:   "Workload report generated",
		Workloads: result,
	}
}
