package services

import (
	"testing"

	"project-studio/internal/auth"
	"project-studio/internal/db"
	"project-studio/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func setupTestEnv(t *testing.T) string {
	t.Helper()
	keyring.MockInit()
	require.NoError(t, db.InitDB(":memory:"))
	require.NoError(t, db.SeedDefaultAdmin())

	var userID string
	require.NoError(t, db.DB.QueryRow(`SELECT id FROM users LIMIT 1`).Scan(&userID))
	require.NoError(t, auth.SaveToken(userID, "test-token"))

	t.Cleanup(func() {
		_ = auth.DeleteToken()
		db.Close()
	})
	return userID
}

func TestProjectService_CreateDuplicateDelete(t *testing.T) {
	setupTestEnv(t)

	activitySvc := NewActivityService()
	projectSvc := NewProjectService(activitySvc)
	taskSvc := NewTaskService(NewNotificationService(), activitySvc)

	startDate := "2026-11-01"
	resp := projectSvc.CreateProject(models.CreateProjectRequest{
		Name:        "Alpha Project",
		Description: "Test with software template",
		Color:       "#6366F1",
		Template:    "software-development",
		StartDate:   &startDate,
		TotalBudget: 15000000,
		Currency:    "IDR",
	})
	require.True(t, resp.Success, resp.Message)
	require.NotNil(t, resp.Project)

	tasksResp := taskSvc.GetTasksByProject(resp.Project.ID)
	require.True(t, tasksResp.Success)
	assert.Len(t, tasksResp.Tasks, 7)

	// Duplicate project and ensure task dependencies are remapped to new task IDs.
	dupResp := projectSvc.DuplicateProject(resp.Project.ID)
	require.True(t, dupResp.Success, dupResp.Message)
	dupTasks := taskSvc.GetTasksByProject(dupResp.Project.ID)
	require.Len(t, dupTasks.Tasks, 7)

	origIDs := make(map[string]bool)
	for _, tk := range tasksResp.Tasks {
		origIDs[tk.ID] = true
	}
	for _, tk := range dupTasks.Tasks {
		assert.False(t, origIDs[tk.ID], "duplicated task kept original ID")
		for _, dep := range origIDs {
			_ = dep
		}
		if tk.Dependencies != "" {
			for origID := range origIDs {
				assert.NotContains(t, tk.Dependencies, origID, "dependency still points to original project task")
			}
		}
	}

	// Delete original project and verify cascade cleanup.
	delResp := projectSvc.DeleteProject(resp.Project.ID)
	require.True(t, delResp.Success)
	var remaining int
	require.NoError(t, db.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE project_id = ?`, resp.Project.ID).Scan(&remaining))
	assert.Equal(t, 0, remaining)
}
