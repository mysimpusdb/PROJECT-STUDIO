package services

import (
	"testing"

	"project-studio/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskService_PartialUpdateAndWorkload(t *testing.T) {
	userID := setupTestEnv(t)

	actSvc := NewActivityService()
	projSvc := NewProjectService(actSvc)
	taskSvc := NewTaskService(NewNotificationService(), actSvc)
	workloadSvc := NewWorkloadService()

	pResp := projSvc.CreateProject(models.CreateProjectRequest{
		Name:     "Workload Project",
		Template: "blank",
	})
	require.True(t, pResp.Success)

	startDate := "2026-11-10"
	endDate := "2026-11-12"
	cResp := taskSvc.CreateTask(models.CreateTaskRequest{
		ProjectID:      pResp.Project.ID,
		Name:           "Backend API",
		Status:         "todo",
		Priority:       "high",
		StartDate:      &startDate,
		EndDate:        &endDate,
		Duration:       2,
		EstimatedHours: 12,
		Labels:         `["api","go"]`,
		Checklists:     `[{"text":"unit test","done":false}]`,
		EstimatedCost:  500000,
	})
	require.True(t, cResp.Success, cResp.Message)
	require.NotNil(t, cResp.Task)
	assert.Equal(t, 12, cResp.Task.EstimatedHours)

	// Partial update (e.g. Kanban status drag) must not wipe labels, checklists, or estimated_hours.
	progress := 50
	uResp := taskSvc.UpdateTask(models.UpdateTaskRequest{
		ID:         cResp.Task.ID,
		Name:       cResp.Task.Name,
		Status:     "in-progress",
		Priority:   cResp.Task.Priority,
		Progress:   &progress,
		Duration:   cResp.Task.Duration,
		AssigneeID: &userID,
	})
	require.True(t, uResp.Success, uResp.Message)
	assert.Equal(t, "in-progress", uResp.Task.Status)
	assert.Equal(t, 12, uResp.Task.EstimatedHours)
	assert.Equal(t, `["api","go"]`, uResp.Task.Labels)
	assert.Equal(t, `[{"text":"unit test","done":false}]`, uResp.Task.Checklists)
	assert.Equal(t, 500000.0, uResp.Task.EstimatedCost)

	// Workload calculation should allocate 12h / 2d = 6h per day.
	wResp := workloadSvc.GetProjectWorkload(pResp.Project.ID, "2026-11-10", "2026-11-11")
	require.True(t, wResp.Success, wResp.Message)
	require.Len(t, wResp.Workloads, 2)
	assert.Equal(t, 6.0, wResp.Workloads[0].AllocatedHours)
	assert.Equal(t, 6.0, wResp.Workloads[1].AllocatedHours)
}
