package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"project-studio/internal/db"
)

// Authorization model
//
// A user's effective role in a project is resolved as follows:
//  1. their row in project_members (owner, manager, member, viewer);
//  2. otherwise "owner" if they created the project;
//  3. otherwise "manager" if they are an admin of the project's workspace;
//  4. otherwise no access.
//
// Roles map to three permission levels:
//   - permRead:   see the project, its tasks, comments, logs and reports
//   - permWrite:  create/edit/delete tasks, comments, attachments and time logs
//   - permManage: edit/archive/delete the project, manage members and sharing

type permission int

const (
	permRead permission = iota
	permWrite
	permManage
)

const (
	roleOwner   = "owner"
	roleManager = "manager"
	roleMember  = "member"
	roleViewer  = "viewer"
)

var (
	errNotFound  = errors.New("not found")
	errForbidden = errors.New("forbidden")
)

const forbiddenMessage = "Anda tidak memiliki izin untuk melakukan tindakan ini."

func roleAllows(role string, perm permission) bool {
	switch role {
	case roleOwner, roleManager:
		return true
	case roleMember:
		return perm <= permWrite
	case roleViewer:
		return perm == permRead
	}
	return false
}

// projectRole resolves the effective role of userID in projectID.
// It returns errNotFound if the project does not exist and errForbidden if the
// user has no access to it.
func projectRole(userID, projectID string) (string, error) {
	var role string
	err := db.DB.QueryRow(
		`SELECT role FROM project_members WHERE project_id = ? AND user_id = ?`,
		projectID, userID,
	).Scan(&role)
	if err == nil {
		return role, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	var createdBy, workspaceID string
	err = db.DB.QueryRow(
		`SELECT created_by, workspace_id FROM projects WHERE id = ?`, projectID,
	).Scan(&createdBy, &workspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errNotFound
	}
	if err != nil {
		return "", err
	}
	if createdBy == userID {
		return roleOwner, nil
	}

	var workspaceRole string
	err = db.DB.QueryRow(
		`SELECT role FROM workspace_members WHERE workspace_id = ? AND user_id = ?`,
		workspaceID, userID,
	).Scan(&workspaceRole)
	if err == nil && workspaceRole == "admin" {
		return roleManager, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	return "", errForbidden
}

// authorizeProject checks that userID holds at least perm on projectID and
// returns the effective role.
func authorizeProject(userID, projectID string, perm permission) (string, error) {
	role, err := projectRole(userID, projectID)
	if err != nil {
		return "", err
	}
	if !roleAllows(role, perm) {
		return "", errForbidden
	}
	return role, nil
}

func projectIDQuery(query, id string) (string, error) {
	var projectID string
	err := db.DB.QueryRow(query, id).Scan(&projectID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errNotFound
	}
	return projectID, err
}

// authorizeTask checks access to the project that owns taskID.
func authorizeTask(userID, taskID string, perm permission) (projectID string, err error) {
	projectID, err = projectIDQuery(`SELECT project_id FROM tasks WHERE id = ?`, taskID)
	if err != nil {
		return "", err
	}
	_, err = authorizeProject(userID, projectID, perm)
	return projectID, err
}

// authorizeComment checks access to the project that owns the comment.
func authorizeComment(userID, commentID string, perm permission) (projectID string, err error) {
	projectID, err = projectIDQuery(`
		SELECT t.project_id FROM task_comments c
		JOIN tasks t ON t.id = c.task_id
		WHERE c.id = ?`, commentID)
	if err != nil {
		return "", err
	}
	_, err = authorizeProject(userID, projectID, perm)
	return projectID, err
}

// authorizeAttachment checks access to the project that owns the attachment
// and returns the effective role so callers can apply uploader-only rules.
func authorizeAttachment(userID, attachmentID string, perm permission) (role string, err error) {
	projectID, err := projectIDQuery(`
		SELECT t.project_id FROM task_attachments a
		JOIN tasks t ON t.id = a.task_id
		WHERE a.id = ?`, attachmentID)
	if err != nil {
		return "", err
	}
	return authorizeProject(userID, projectID, perm)
}

// authorizeTimeLog checks access to the project that owns the time log.
func authorizeTimeLog(userID, logID string, perm permission) (projectID string, err error) {
	projectID, err = projectIDQuery(`
		SELECT t.project_id FROM time_logs l
		JOIN tasks t ON t.id = l.task_id
		WHERE l.id = ?`, logID)
	if err != nil {
		return "", err
	}
	_, err = authorizeProject(userID, projectID, perm)
	return projectID, err
}

// isWorkspaceMemberOfProject reports whether userID belongs to the workspace
// that owns projectID.
func isWorkspaceMemberOfProject(userID, projectID string) (bool, error) {
	var n int
	err := db.DB.QueryRow(`
		SELECT COUNT(*) FROM projects p
		JOIN workspace_members wm ON wm.workspace_id = p.workspace_id
		WHERE p.id = ? AND wm.user_id = ?`, projectID, userID).Scan(&n)
	return n > 0, err
}

// accessMessage converts an authorization error to a user-facing message.
func accessMessage(err error, notFoundMessage string) string {
	switch {
	case errors.Is(err, errNotFound):
		return notFoundMessage
	case errors.Is(err, errForbidden):
		return forbiddenMessage
	default:
		return "Terjadi kesalahan sistem."
	}
}

// ---- Input validation helpers ----

const dateLayout = "2006-01-02"

// normalizeStatus returns the canonical task status. Legacy underscore
// variants ("in_progress", "on_hold") are mapped to the hyphenated form.
func normalizeStatus(s string) (string, bool) {
	s = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(s)), "_", "-")
	switch s {
	case "todo", "in-progress", "on-hold", "done", "cancelled":
		return s, true
	}
	return "", false
}

// normalizePriority returns the canonical task priority. Legacy values
// ("normal", "critical") are mapped to "medium" and "urgent".
func normalizePriority(p string) (string, bool) {
	p = strings.ToLower(strings.TrimSpace(p))
	switch p {
	case "low", "medium", "high", "urgent":
		return p, true
	case "normal":
		return "medium", true
	case "critical":
		return "urgent", true
	}
	return "", false
}

func normalizeProjectStatus(s string) (string, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "active", "archived", "completed":
		return s, true
	}
	return "", false
}

// parseDateInput parses an optional YYYY-MM-DD string.
// set is false when the pointer is nil (keep the current value); a non-nil
// empty string means "clear" and returns set=true with a nil time.
func parseDateInput(v *string) (t *time.Time, set bool, err error) {
	if v == nil {
		return nil, false, nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil, true, nil
	}
	if len(s) > 10 {
		s = s[:10]
	}
	parsed, err := time.Parse(dateLayout, s)
	if err != nil {
		return nil, true, fmt.Errorf("format tanggal tidak valid (gunakan YYYY-MM-DD): %q", *v)
	}
	return &parsed, true, nil
}

// checkDateOrder returns an error if both dates are present and end < start.
func checkDateOrder(start, end *time.Time) error {
	if start != nil && end != nil && end.Before(*start) {
		return errors.New("tanggal selesai tidak boleh lebih awal dari tanggal mulai")
	}
	return nil
}

// spreadsheetSafe neutralises cells that spreadsheet applications would
// interpret as formulas (CSV/Excel formula injection).
func spreadsheetSafe(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + s
	}
	return s
}
